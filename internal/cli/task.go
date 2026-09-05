package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

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
	n, ok := tracker.MilestoneNumber("M" + strings.TrimPrefix(*milestoneArg, "M") + ":")
	if !ok {
		return fmt.Errorf("bad milestone number %q", *milestoneArg)
	}
	c, err := load()
	if err != nil {
		return err
	}
	target := *repo
	if target == "" {
		target = c.current
	}
	return runTaskNew(c, w, n, target, *title, *goal, *requirements)
}

// runTaskNew creates the task issue in target and attaches it to milestone
// M<n> as a sub-issue.
func runTaskNew(c *ctx, w io.Writer, n int, target, title, goal, requirements string) error {
	milestone, err := resolveMilestone(c, w, n)
	if err != nil {
		return err
	}
	body := fmt.Sprintf(taskTemplate, goal, requirements, tracker.PlanPlaceholder)
	ref, err := c.t.CreateIssue(target, title, body, []string{"cc:task"})
	if err != nil {
		return err
	}
	if err := c.t.AddSubIssue(milestone.Ref, ref); err != nil {
		return err
	}
	fmt.Fprintf(w, "created task %s as a sub-issue of %s\n", ref, milestone.Ref)
	return nil
}

// milestoneLookupAttempts bounds how many times task new reads the listings
// for a milestone that is not there yet, and milestoneLookupWait is the
// pause between reads. Three refusals in a row hit a milestone created
// seconds earlier before the fourth call found it (#234); a few seconds is
// the window the listing has been seen to lag by (#195).
const (
	milestoneLookupAttempts = 3
	milestoneLookupWait     = 2 * time.Second
)

// sleep is a func var so tests record the waits instead of taking them.
var sleep = time.Sleep

// resolveMilestone finds the open milestone M<n> in the hub. The open
// milestone listing is label-filtered and eventually consistent: it can lag
// an issue created seconds ago (#195, #234), so a miss is not yet NOT_FOUND.
// The unfiltered newest issues are read for an `M<n>:` title — accepted
// when the issue itself carries cc:milestone and is open — and, failing
// that, the reads repeat after a short wait, bounded. A milestone found by
// either route is reported on w, so the record shows the listing lagged.
func resolveMilestone(c *ctx, w io.Writer, n int) (*tracker.Milestone, error) {
	for attempt := 1; ; attempt++ {
		milestones, err := c.t.OpenMilestones(c.hub)
		if err != nil {
			return nil, err
		}
		for i := range milestones {
			if got, ok := tracker.MilestoneNumber(milestones[i].Title); ok && got == n {
				if attempt > 1 {
					fmt.Fprintf(w, "milestone M%d (%s) appeared in the listing on read %d — the label-filtered listing lags a milestone created seconds ago\n", n, milestones[i].Ref, attempt)
				}
				return &milestones[i], nil
			}
		}
		m, err := recentMilestone(c, n)
		if err != nil {
			return nil, err
		}
		if m != nil {
			fmt.Fprintf(w, "milestone M%d (%s) is not in the open-milestone listing yet, found among the hub's newest issues — the label-filtered listing lags a milestone created seconds ago\n", n, m.Ref)
			return m, nil
		}
		if attempt == milestoneLookupAttempts {
			return nil, refuse("NOT_FOUND", "no open milestone M%d in %s after %d reads of the listing %v apart — the listing can lag a milestone created seconds ago, so retry in a moment if milestone new just created it; gh codecrew status lists the open milestones", n, c.hub, attempt, milestoneLookupWait)
		}
		sleep(milestoneLookupWait)
	}
}

// recentMilestone looks for M<n> among the hub's newest issues regardless
// of label — the second source milestone new reads for the same reason
// (#209). A title match is confirmed on the issue itself: it must carry
// cc:milestone and be open, since the listing is unfiltered and shows
// closed issues too. Nil when there is no such milestone.
func recentMilestone(c *ctx, n int) (*tracker.Milestone, error) {
	recent, err := c.t.RecentIssues(c.hub)
	if err != nil {
		return nil, err
	}
	for _, is := range recent {
		got, ok := tracker.MilestoneNumber(is.Title)
		if !ok || got != n {
			continue
		}
		t, err := c.t.Task(is.Ref)
		if err != nil {
			return nil, err
		}
		if t.Closed || !tracker.ContainsLabel(t.Labels, tracker.LabelMilestone) {
			continue
		}
		return &tracker.Milestone{Ref: is.Ref, Title: is.Title}, nil
	}
	return nil, nil
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
	return raiseGate(w, c, ref, *question)
}

// raiseGate posts the **Gate raised:** comment and applies the label. The
// ref may be a task, a milestone issue or a pull request: a question about
// a requirement has no task to carry it, so it is raised on the milestone
// issue, where nothing mechanical blocks on the label — status lists the
// gate instead (#200); before the first milestone the gate is recorded on
// the scaffold PR (roles/coordinator.md). The labels come from the REST
// issues endpoint, which serves PRs — Task's GraphQL issue query does not
// (checky's finding on PR #218). The comment and the receipt say which of
// the two wordings holds; a PR gets the task's.
func raiseGate(w io.Writer, c *ctx, ref tracker.IssueRef, question string) error {
	labels, err := c.t.IssueLabels(ref)
	if err != nil {
		return err
	}
	onMilestone := tracker.ContainsLabel(labels, tracker.LabelMilestone)
	msg := "**Gate raised:** " + question + "\n\n" +
		"Resolve by replying with a comment starting `**Gate resolved:**` — the decision, " +
		"and the trade-off if one was weighed — then removing the `cc:needs-decision` label. " +
		"The resolution is gathered into the milestone record; "
	if onMilestone {
		msg += "a question about a requirement has no task to carry it, so it is raised here, " +
			"and `status` lists this gate beside the tasks' gates while the label is present."
	} else {
		msg += "`task finish` refuses while the label is present or the gate lacks a resolution comment."
	}
	if err := c.t.Comment(ref, msg); err != nil {
		return err
	}
	if err := c.t.AddLabel(ref, tracker.LabelNeedsDecision); err != nil {
		return err
	}
	if onMilestone {
		fmt.Fprintf(w, "gate raised on %s (milestone issue) — status lists it beside the tasks' gates until a human removes %s\n", ref, tracker.LabelNeedsDecision)
	} else {
		fmt.Fprintf(w, "gate raised on %s — blocked until a human removes %s\n", ref, tracker.LabelNeedsDecision)
	}
	return nil
}

// finishGates, in the order task finish checks them — the dry run prints
// them in this order, refused, ok, not reached or not applicable.
var finishGates = []string{"task open", "no gate raised", "gates resolved", "owner", "closing PR", "CI checks", "review", "GitHub's required review", "operator confirmation", "bypass actor"}

func taskFinish(w io.Writer, args []string) error {
	refArg, args := splitLeadingRef(args)
	fs := flag.NewFlagSet("task finish", flag.ContinueOnError)
	operatorConfirm := fs.Bool("operator-confirm", false,
		"solo tier: record explicit operator confirmation in place of a non-doer approval")
	bypass := fs.Bool("bypass", false,
		"merge with the ruleset's administrator bypass when GitHub does not count the recorded approval (recorded on the PR)")
	dryRun := fs.Bool("dry-run", false, "print every gate and what the verb would do; write nothing; exit with the first refusal's code")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if refArg == "" && fs.NArg() == 1 {
		refArg = fs.Arg(0)
	}
	if refArg == "" {
		return fmt.Errorf("usage: gh codecrew task finish <ref> [--operator-confirm] [--bypass] [--dry-run]")
	}
	c, err := load()
	if err != nil {
		return err
	}
	ref, err := tracker.ParseRef(refArg, c.current)
	if err != nil {
		return err
	}
	p, run, err := planFinish(c, ref, *operatorConfirm, *bypass)
	if err != nil {
		return err
	}
	if *dryRun {
		p.print(w)
		return p.refusal
	}
	if p.refusal != nil {
		return p.refusal
	}
	return run(w)
}

// planFinish evaluates task finish's gates in order and decides its
// actions without writing anything; run performs them. A non-gate error
// (the API) is returned as such — it is not a refusal.
func planFinish(c *ctx, ref tracker.IssueRef, operatorConfirm, bypass bool) (*plan, func(io.Writer) error, error) {
	p := &plan{}
	crew := func(login string) bool { return strings.HasSuffix(login, "[bot]") || c.roleFor(login) != "" }
	task, err := c.t.Task(ref)
	if err != nil {
		return nil, nil, err
	}
	var e error
	if task.Closed {
		e = refuse("CLOSED", "%s is already closed", ref)
	}
	if !p.gate("task open", e) {
		return p.stop(finishGates), nil, nil
	}
	e = nil
	if tracker.HasLabel(task, tracker.LabelNeedsDecision) {
		e = refuse("GATED", "%s has %s raised — a human must resolve it and remove the label", ref, tracker.LabelNeedsDecision)
	}
	if !p.gate("no gate raised", e) {
		return p.stop(finishGates), nil, nil
	}
	comments, err := c.t.Comments(ref)
	if err != nil {
		return nil, nil, err
	}
	e = nil
	if unresolved := tracker.UnresolvedGates(comments); len(unresolved) > 0 {
		e = refuse("GATE_UNRECORDED", "%d gate(s) on %s lack a resolution record — reply with **Gate resolved:** stating the decision (SPEC §8): %s",
			len(unresolved), ref, unresolved[0].URL)
	}
	if !p.gate("gates resolved", e) {
		return p.stop(finishGates), nil, nil
	}
	// The seat that started a task finishes it (#165: the coordinator's
	// table once sent "approved → the implementer" regardless of whose
	// PR it was, and the implementer merged the doc-synthesizer's
	// document). The operator's own auth is not exempt — finishing a
	// routed seat's task is the same misattribution; --bypass is the
	// recorded escape hatch, and it stays an operator's act.
	viewer, err := c.t.Viewer()
	if err != nil {
		return nil, nil, err
	}
	var ownerBypass string
	owner := tracker.StartedBy(task, comments)
	e = nil
	if !sameSeat(owner, viewer, c.roleFor) {
		switch {
		case !bypass:
			e = refuse("NOT_OWNER", "%s was started by @%s, not @%s — the seat that started a task finishes it: dispatch it, hand it over with task start, or an operator overrides on the record with --bypass (SPEC §8)", ref, owner, viewer)
		case crew(viewer):
			e = refuse("CREW_BYPASS", "--bypass requires a human operator; @%s is a crew identity", viewer)
		default:
			ownerBypass = owner
		}
	}
	if !p.gate("owner", e) {
		return p.stop(finishGates), nil, nil
	}

	prs, err := c.t.ClosingPRs(ref, false)
	if err != nil {
		return nil, nil, err
	}
	e = nil
	if len(prs) == 0 {
		e = refuse("NO_PR", "no open PR closes %s", ref)
	}
	if !p.gate("closing PR", e) {
		return p.stop(finishGates), nil, nil
	}
	pr, err := c.t.PRInfo(ref.Repo, prs[0])
	if err != nil {
		return nil, nil, err
	}
	e = nil
	switch {
	case pr.ChecksUnreadable != "":
		// The token could not read the checks at all — a private repo
		// refuses the rollup to an App without the permission, where a
		// public one reads it freely (#198). The App is the viewer.
		e = refuse("NO_CHECKS_PERMISSION", "PR #%d's checks are unreadable by %s: the installation token lacks `%s`, which a private repo requires to read the status check rollup (a public repo reads it without) — add the permission on the App's settings page (Permissions & events → Repository permissions), then accept the change on the installation (Installed GitHub Apps → Configure); GitHub exposes neither through the API (docs/identities.md)", pr.Number, seatName(viewer), pr.ChecksUnreadable)
	case pr.NoChecks:
		e = refuse("NO_CHECKS", "PR #%d has no CI checks — the deterministic gate cannot be satisfied by absence (SPEC §8); add a workflow that runs on pull_request, let it report, then re-run", pr.Number)
	case pr.ChecksPending:
		e = refuse("CHECKS_PENDING", "PR #%d checks still running", pr.Number)
	case !pr.ChecksOK:
		e = refuse("CHECKS_FAILING", "PR #%d has failing checks", pr.Number)
	}
	if !p.gate("CI checks", e) {
		return p.stop(finishGates), nil, nil
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
	e = nil
	if reviewerHolder, err := holder(c.rolesConfig().Roles, "reviewer"); err == nil && reviewerHolder != "~" {
		if !holderReviewed(pr.ApprovedBy, func(login string) bool { return c.holdsRole(login, "reviewer") }) {
			e = refuse("NO_HOLDER_REVIEW", "PR #%d has no approving review from the reviewer role's holder (%s) — the role defines whose review counts; dispatch the reviewer (docs/identities.md)", pr.Number, reviewerHolder)
		}
	} else if !approved && !operatorConfirm {
		e = refuse("NO_NONDOER_APPROVAL", "PR #%d has no approving review from a non-author (solo tier: rerun with --operator-confirm)", pr.Number)
	}
	if !p.gate("review", e) {
		return p.stop(finishGates), nil, nil
	}
	// Decide the merge path before writing anything: a REVIEW_NOT_COUNTED
	// refusal must land before the operator-confirm comment, or a rerun
	// with --bypass would record the confirmation twice (testy finding on
	// PR #89).
	admin, e := mergeGate(pr.ReviewDecision, bypass)
	if !p.gate("GitHub's required review", e) {
		return p.stop(finishGates), nil, nil
	}

	prRef := tracker.IssueRef{Repo: pr.Repo, Number: pr.Number}
	var posts []string
	if ownerBypass != "" {
		// Recorded before any merge path: the override is the fact, whatever
		// GitHub then does with the merge.
		posts = append(posts, fmt.Sprintf("**Owner bypass:** %s was started by @%s; finished by @%s with --bypass (SPEC §8 — the seat that started a task finishes it; this is the operator's recorded override).", ref, ownerBypass, viewer))
	}
	if !approved {
		// Crew identities can never waive review; a human operator can —
		// including on their own PR, where no distinct principal exists
		// (pure solo tier, SPEC §5) — and the record says so.
		e = nil
		if crew(viewer) {
			e = refuse("SELF_CONFIRM", "--operator-confirm requires a human operator; @%s is a crew identity", viewer)
		}
		if !p.gate("operator confirmation", e) {
			return p.stop(finishGates), nil, nil
		}
		msg := fmt.Sprintf("**Operator confirmation:** reviewed and accepted by @%s in place of a formal approval (solo tier, SPEC §6).", viewer)
		if viewer == pr.Author {
			msg = fmt.Sprintf("**Operator confirmation:** reviewed and accepted by @%s as both author and operator (pure solo tier, SPEC §6) — no independent principal exists in this project.", viewer)
		}
		posts = append(posts, msg)
	} else {
		p.na("operator confirmation")
	}
	if admin {
		// A bypass is an operator act: crew identities never hold it.
		e = nil
		if crew(viewer) {
			e = refuse("CREW_BYPASS", "--bypass requires a human operator; @%s is a crew identity", viewer)
		}
		if !p.gate("bypass actor", e) {
			return p.stop(finishGates), nil, nil
		}
		// Worded as the act, not the outcome: GitHub enforces bypass
		// eligibility after this comment posts, and a false "merged"
		// record must not survive a refused bypass (testy finding on
		// PR #89).
		posts = append(posts, fmt.Sprintf("**Merge bypass:** the protocol's review gate is satisfied, but GitHub's "+
			"required-review rule does not count the recorded approvals — approvals count only from "+
			"principals with write access (superseding Decision, radiusred/gh-codecrew#73). Merging with "+
			"the ruleset's administrator bypass as @%s; GitHub enforces eligibility — if this PR remains "+
			"unmerged, the bypass was refused.", viewer))
	} else {
		p.na("bypass actor")
	}
	for _, m := range posts {
		p.would("comment on PR #%d: %s", pr.Number, firstLine(m))
	}
	if admin {
		p.would("merge PR #%d via the ruleset's administrator bypass", pr.Number)
	} else {
		p.would("merge PR #%d (rebase); %s closes via its closing keyword", pr.Number, ref)
	}
	if pr.HeadRef != "" {
		p.would("delete head %s", pr.HeadRef)
	}
	run := func(w io.Writer) error {
		for _, m := range posts {
			if err := c.t.Comment(prRef, m); err != nil {
				return err
			}
		}
		if admin {
			if err := c.t.MergePRBypass(pr.Repo, pr.Number); err != nil {
				return err
			}
			fmt.Fprintf(w, "merged PR #%d via administrator bypass (recorded); %s closes via its closing keyword\n", pr.Number, ref)
		} else {
			if err := c.t.MergePR(pr.Repo, pr.Number); err != nil {
				return err
			}
			fmt.Fprintf(w, "merged PR #%d; %s closes via its closing keyword\n", pr.Number, ref)
		}
		deleteHead(w, c.t, pr)
		return nil
	}
	return p, run, nil
}

// seatName names the identity a refusal is about: the App behind a
// `[bot]` login by its slug, anyone else by their login.
func seatName(viewer string) string {
	if slug, ok := strings.CutSuffix(viewer, "[bot]"); ok {
		return "App " + slug
	}
	return "@" + viewer
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
