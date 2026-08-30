package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

const taskTemplate = `## Goal
%s

## Requirements
%s

## Plan
%s

## Ask-the-human points
_None identified yet — the implementer adds any before starting._
`

func taskNew(w io.Writer, args []string) error {
	fs := flag.NewFlagSet("task new", flag.ContinueOnError)
	milestoneArg := fs.String("milestone", "", "milestone number to link into (required)")
	title := fs.String("title", "", "task title (required)")
	goal := fs.String("goal", "_To be written._", "what this task delivers")
	requirements := fs.String("requirements", "None directly.", "comma-separated requirement IDs")
	repo := fs.String("repo", "", "spoke repo for the task (default: current repo)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *title == "" || *milestoneArg == "" {
		return fmt.Errorf("task new: --title and --milestone are required")
	}
	c, err := load()
	if err != nil {
		return err
	}
	target := *repo
	if target == "" {
		target = c.current
	}
	n, ok := tracker.MilestoneNumber("M" + strings.TrimPrefix(*milestoneArg, "M") + ":")
	if !ok {
		return fmt.Errorf("bad milestone number %q", *milestoneArg)
	}
	milestones, err := c.t.OpenMilestones(c.hub)
	if err != nil {
		return err
	}
	var milestone *tracker.Milestone
	for i := range milestones {
		if got, ok := tracker.MilestoneNumber(milestones[i].Title); ok && got == n {
			milestone = &milestones[i]
			break
		}
	}
	if milestone == nil {
		return refuse("NOT_FOUND", "no open milestone M%d in %s", n, c.hub)
	}

	body := fmt.Sprintf(taskTemplate, *goal, *requirements, tracker.PlanPlaceholder)
	ref, err := c.t.CreateIssue(target, *title, body, []string{"cc:task"})
	if err != nil {
		return err
	}
	if err := c.t.AddSubIssue(milestone.Ref, ref); err != nil {
		return err
	}
	fmt.Fprintf(w, "created task %s as a sub-issue of %s\n", ref, milestone.Ref)
	return nil
}

func taskStart(w io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: gh codecrew task start <ref>")
	}
	c, err := load()
	if err != nil {
		return err
	}
	ref, err := tracker.ParseRef(args[0], c.current)
	if err != nil {
		return err
	}
	task, err := c.t.Task(ref)
	if err != nil {
		return err
	}
	if task.Closed {
		return refuse("CLOSED", "%s is already closed", ref)
	}
	if !tracker.HasLabel(task, "cc:task") {
		return refuse("NOT_A_TASK", "%s is not labeled cc:task", ref)
	}
	body, err := c.t.IssueBody(ref)
	if err != nil {
		return err
	}
	if !tracker.PlanPresent(body) {
		return refuse("NO_PLAN", "%s has no plan — write the Plan section before starting (SPEC §4)", ref)
	}

	viewer, err := c.t.Viewer()
	if err != nil {
		return err
	}
	if err := c.t.Assign(ref, viewer); err != nil {
		// Bot identities are not assignable; the assignment is a courtesy
		// for humans, the record below is the fact.
		fmt.Fprintf(w, "note: could not assign @%s (%v)\n", viewer, err)
	}
	// Every start posts the record, so the latest one is the owner task
	// finish holds to (NOT_OWNER) — across a restart or a handover.
	if err := c.t.Comment(ref, tracker.StartRecord(viewer)); err != nil {
		return err
	}

	if role := c.roleFor(viewer); role == "qa" || role == "reviewer" {
		fmt.Fprintf(w, "role %s does not commit (roles/%s.md); no linked branch created\n", role, role)
	} else {
		branch := fmt.Sprintf("task/%d-%s", ref.Number, slug(task.Title))
		if err := c.t.DevelopBranch(ref, branch); err != nil {
			fmt.Fprintf(w, "note: could not create linked branch (%v); create %q manually\n", err, branch)
		} else {
			fmt.Fprintf(w, "linked branch %s created\n", branch)
			fmt.Fprintf(w, "locally: git fetch && git switch %s\n", branch)
		}
	}
	fmt.Fprintf(w, "started %s as @%s\n", ref, viewer)
	return nil
}

func checkpoint(w io.Writer, args []string) error {
	refArg, args := splitLeadingRef(args)
	fs := flag.NewFlagSet("checkpoint", flag.ContinueOnError)
	question := fs.String("question", "", "the judgment call a human must make (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if refArg == "" && fs.NArg() == 1 {
		refArg = fs.Arg(0)
	}
	if *question == "" || refArg == "" {
		return fmt.Errorf("usage: gh codecrew checkpoint <ref> --question \"...\"")
	}
	c, err := load()
	if err != nil {
		return err
	}
	ref, err := tracker.ParseRef(refArg, c.current)
	if err != nil {
		return err
	}
	msg := "**Gate raised:** " + *question + "\n\n" +
		"Resolve by replying with a comment starting `**Gate resolved:**` — the decision, " +
		"and the trade-off if one was weighed — then removing the `cc:needs-decision` label. " +
		"The resolution is gathered into the milestone record; `task finish` refuses while " +
		"the label is present or the gate lacks a resolution comment."
	if err := c.t.Comment(ref, msg); err != nil {
		return err
	}
	if err := c.t.AddLabel(ref, tracker.LabelNeedsDecision); err != nil {
		return err
	}
	fmt.Fprintf(w, "gate raised on %s — blocked until a human removes %s\n", ref, tracker.LabelNeedsDecision)
	return nil
}

func taskFinish(w io.Writer, args []string) error {
	refArg, args := splitLeadingRef(args)
	fs := flag.NewFlagSet("task finish", flag.ContinueOnError)
	operatorConfirm := fs.Bool("operator-confirm", false,
		"solo tier: record explicit operator confirmation in place of a non-doer approval")
	bypass := fs.Bool("bypass", false,
		"merge with the ruleset's administrator bypass when GitHub does not count the recorded approval (recorded on the PR)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if refArg == "" && fs.NArg() == 1 {
		refArg = fs.Arg(0)
	}
	if refArg == "" {
		return fmt.Errorf("usage: gh codecrew task finish <ref> [--operator-confirm] [--bypass]")
	}
	c, err := load()
	if err != nil {
		return err
	}
	ref, err := tracker.ParseRef(refArg, c.current)
	if err != nil {
		return err
	}
	task, err := c.t.Task(ref)
	if err != nil {
		return err
	}
	if task.Closed {
		return refuse("CLOSED", "%s is already closed", ref)
	}
	if tracker.HasLabel(task, tracker.LabelNeedsDecision) {
		return refuse("GATED", "%s has %s raised — a human must resolve it and remove the label", ref, tracker.LabelNeedsDecision)
	}
	comments, err := c.t.Comments(ref)
	if err != nil {
		return err
	}
	if unresolved := tracker.UnresolvedGates(comments); len(unresolved) > 0 {
		return refuse("GATE_UNRECORDED", "%d gate(s) on %s lack a resolution record — reply with **Gate resolved:** stating the decision (SPEC §8): %s",
			len(unresolved), ref, unresolved[0].URL)
	}
	// The seat that started a task finishes it (#165: the coordinator's
	// table once sent "approved → the implementer" regardless of whose
	// PR it was, and the implementer merged the doc-synthesizer's
	// document). The operator's own auth is not exempt — finishing a
	// routed seat's task is the same misattribution; --bypass is the
	// recorded escape hatch, and it stays an operator's act.
	viewer, err := c.t.Viewer()
	if err != nil {
		return err
	}
	var ownerBypass string
	owner := tracker.StartedBy(task, comments)
	if !sameSeat(owner, viewer, c.roleFor) {
		if !*bypass {
			return refuse("NOT_OWNER", "%s was started by @%s, not @%s — the seat that started a task finishes it: dispatch it, hand it over with task start, or an operator overrides on the record with --bypass (SPEC §8)", ref, owner, viewer)
		}
		if strings.HasSuffix(viewer, "[bot]") || c.roleFor(viewer) != "" {
			return refuse("CREW_BYPASS", "--bypass requires a human operator; @%s is a crew identity", viewer)
		}
		ownerBypass = owner
	}

	prs, err := c.t.ClosingPRs(ref, false)
	if err != nil {
		return err
	}
	if len(prs) == 0 {
		return refuse("NO_PR", "no open PR closes %s", ref)
	}
	pr, err := c.t.PRInfo(ref.Repo, prs[0])
	if err != nil {
		return err
	}
	if pr.NoChecks {
		return refuse("NO_CHECKS", "PR #%d has no CI checks — the deterministic gate cannot be satisfied by absence (SPEC §8); add a workflow that runs on pull_request, let it report, then re-run", pr.Number)
	}
	if pr.ChecksPending {
		return refuse("CHECKS_PENDING", "PR #%d checks still running", pr.Number)
	}
	if !pr.ChecksOK {
		return refuse("CHECKS_FAILING", "PR #%d has failing checks", pr.Number)
	}

	approved := false
	for _, login := range pr.ApprovedBy {
		if login != pr.Author {
			approved = true
		}
	}
	// The model-review mandate (#73 Decision): where a distinct principal
	// holds the reviewer seat, that holder's approving review is the gate —
	// other approvals coexist but do not satisfy it, and the solo
	// confirmation cannot stand in for a holder that exists.
	if reviewerHolder, err := holder(c.rolesConfig().Roles, "reviewer"); err == nil && reviewerHolder != "~" {
		if !holderReviewed(pr.ApprovedBy, func(login string) bool { return c.holdsRole(login, "reviewer") }) {
			return refuse("NO_HOLDER_REVIEW", "PR #%d has no approving review from the reviewer role's holder (%s) — the role defines whose review counts; dispatch the reviewer (docs/identities.md)", pr.Number, reviewerHolder)
		}
	} else if !approved && !*operatorConfirm {
		return refuse("NO_NONDOER_APPROVAL", "PR #%d has no approving review from a non-author (solo tier: rerun with --operator-confirm)", pr.Number)
	}
	// Decide the merge path before writing anything: a REVIEW_NOT_COUNTED
	// refusal must land before the operator-confirm comment, or a rerun
	// with --bypass would record the confirmation twice (testy finding on
	// PR #89).
	admin, err := mergeGate(pr.ReviewDecision, *bypass)
	if err != nil {
		return err
	}
	if ownerBypass != "" {
		// Recorded before any merge path: the override is the fact, whatever
		// GitHub then does with the merge.
		prRef := tracker.IssueRef{Repo: pr.Repo, Number: pr.Number}
		msg := fmt.Sprintf("**Owner bypass:** %s was started by @%s; finished by @%s with --bypass (SPEC §8 — the seat that started a task finishes it; this is the operator's recorded override).", ref, ownerBypass, viewer)
		if err := c.t.Comment(prRef, msg); err != nil {
			return err
		}
	}
	if !approved {
		// Crew identities can never waive review; a human operator can —
		// including on their own PR, where no distinct principal exists
		// (pure solo tier, SPEC §5) — and the record says so.
		if strings.HasSuffix(viewer, "[bot]") || c.roleFor(viewer) != "" {
			return refuse("SELF_CONFIRM", "--operator-confirm requires a human operator; @%s is a crew identity", viewer)
		}
		prRef := tracker.IssueRef{Repo: pr.Repo, Number: pr.Number}
		msg := fmt.Sprintf("**Operator confirmation:** reviewed and accepted by @%s in place of a formal approval (solo tier, SPEC §6).", viewer)
		if viewer == pr.Author {
			msg = fmt.Sprintf("**Operator confirmation:** reviewed and accepted by @%s as both author and operator (pure solo tier, SPEC §6) — no independent principal exists in this project.", viewer)
		}
		if err := c.t.Comment(prRef, msg); err != nil {
			return err
		}
	}

	if admin {
		// A bypass is an operator act: crew identities never hold it.
		if strings.HasSuffix(viewer, "[bot]") || c.roleFor(viewer) != "" {
			return refuse("CREW_BYPASS", "--bypass requires a human operator; @%s is a crew identity", viewer)
		}
		prRef := tracker.IssueRef{Repo: pr.Repo, Number: pr.Number}
		// Worded as the act, not the outcome: GitHub enforces bypass
		// eligibility after this comment posts, and a false "merged"
		// record must not survive a refused bypass (testy finding on
		// PR #89).
		msg := fmt.Sprintf("**Merge bypass:** the protocol's review gate is satisfied, but GitHub's "+
			"required-review rule does not count the recorded approvals — approvals count only from "+
			"principals with write access (superseding Decision, radiusred/gh-codecrew#73). Merging with "+
			"the ruleset's administrator bypass as @%s; GitHub enforces eligibility — if this PR remains "+
			"unmerged, the bypass was refused.", viewer)
		if err := c.t.Comment(prRef, msg); err != nil {
			return err
		}
		if err := c.t.MergePRBypass(pr.Repo, pr.Number); err != nil {
			return err
		}
		fmt.Fprintf(w, "merged PR #%d via administrator bypass (recorded); %s closes via its closing keyword\n", pr.Number, ref)
		deleteHead(w, c.t, pr)
		return nil
	}
	if err := c.t.MergePR(pr.Repo, pr.Number); err != nil {
		return err
	}
	fmt.Fprintf(w, "merged PR #%d; %s closes via its closing keyword\n", pr.Number, ref)
	deleteHead(w, c.t, pr)
	return nil
}

// mergeGate decides the merge path from GitHub's own review decision,
// after the protocol's review gate has passed. REVIEW_REQUIRED means the
// platform will refuse the normal merge — the protocol refuses first, with
// the supported paths, unless the operator asked for the recorded bypass.
func mergeGate(reviewDecision string, bypass bool) (admin bool, err error) {
	if reviewDecision != "REVIEW_REQUIRED" {
		return false, nil
	}
	if !bypass {
		return false, refuse("REVIEW_NOT_COUNTED",
			"GitHub's required-review rule is not satisfied by the recorded approvals — approvals count "+
				"only from principals with write access, so a read-only App's review and an operator "+
				"confirmation do not (superseding Decision, radiusred/gh-codecrew#73). Either a non-author "+
				"human approves on the reviewer's recommendation, grant the reviewer App write access so its "+
				"approvals count, or rerun with --bypass if the ruleset lists you as a bypass actor")
	}
	return true, nil
}

// holderReviewed reports whether any approving login holds the reviewer
// role per the routing table (holds wraps HoldsRole, which normalizes the
// [bot] suffix) — the mandate half of the model-review Decision on #73.
func holderReviewed(approvedBy []string, holds func(login string) bool) bool {
	for _, login := range approvedBy {
		if holds(login) {
			return true
		}
	}
	return false
}

// splitLeadingRef pulls a leading positional ref off the args so verbs
// accept both "<ref> --flag" and "--flag <ref>" orders (stdlib flag stops
// parsing at the first positional).
// sameSeat reports whether viewer may finish a task owner started: the
// same login, or the same routed seat — a role held by a GitHub team is
// held by any member, so a teammate finishing a teammate's task is the
// seat finishing its own. An owner who has since left the team resolves
// to no role and no longer matches: the task is handed over by running
// task start again (latest record wins) or finished by the operator with
// --bypass on the record. No owner recorded: nothing to hold anyone to.
func sameSeat(owner, viewer string, roleFor func(string) string) bool {
	if owner == "" || tracker.SameLogin(owner, viewer) {
		return true
	}
	role := roleFor(owner)
	return role != "" && role == roleFor(viewer)
}

func splitLeadingRef(args []string) (string, []string) {
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		return args[0], args[1:]
	}
	return "", args
}

func slug(title string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
		case b.Len() > 0 && !strings.HasSuffix(b.String(), "-"):
			b.WriteRune('-')
		}
	}
	s := strings.Trim(b.String(), "-")
	if len(s) > 40 {
		s = strings.Trim(s[:40], "-")
	}
	return s
}
