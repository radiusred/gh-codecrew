package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// A struck requirement (SPEC §4, M18-R4): one withdrawn from the
// milestone's scope by a recorded Decision. The coordination layer strikes
// it — whoever holds coordination owns the scope change, human or agent —
// by posting `**M<n>-R<k> — struck.** <decision link>` on the milestone
// issue, never by editing the body. `milestone close` counts a verified
// strike as terminal; `--reinstate` posts the `reinstated` twin under the
// same check, after which QA verdicts count again.

// strikeGates, in the order milestone strike checks them.
var strikeGates = []string{"milestone open", "requirement declared", "decision recorded"}

// requirementIDArg is the whole ID the verb takes — what is typed is what
// is posted, so `R2` is not shorthand for anything.
var requirementIDArg = regexp.MustCompile(`^M(\d+)-R\d+$`)

func milestoneStrike(w io.Writer, args []string) error {
	fs := flag.NewFlagSet("milestone strike", flag.ContinueOnError)
	decision := fs.String("decision", "", "the issue-comment URL of the recorded Decision that strikes the requirement (required)")
	reinstate := fs.Bool("reinstate", false, "post the reinstated line instead: the requirement is back in scope and QA verdicts count again")
	dryRun := fs.Bool("dry-run", false, "print every gate and the line it would post; write nothing; exit with the first refusal's code")
	var pos []string
	for {
		var a string
		a, args = splitLeadingRef(args)
		if a == "" {
			break
		}
		pos = append(pos, a)
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	pos = append(pos, fs.Args()...)
	const usage = "usage: gh codecrew milestone strike <milestone number> <requirement ID> --decision <comment URL> [--reinstate] [--dry-run]"
	if len(pos) != 2 {
		return errors.New(usage)
	}
	n, err := strconv.Atoi(strings.TrimPrefix(pos[0], "M"))
	if err != nil {
		return fmt.Errorf("bad milestone number %q — %s", pos[0], usage)
	}
	id := pos[1]
	if !requirementIDArg.MatchString(id) {
		return fmt.Errorf("milestone strike: %q is not a requirement ID — write it in full, M%d-R<k>", id, n)
	}
	if *decision == "" {
		return fmt.Errorf("milestone strike: --decision is required — the URL of the comment carrying the Decision, on the milestone issue or one of its tasks")
	}
	c, err := load()
	if err != nil {
		return err
	}
	p, run, err := planStrike(c, n, id, *decision, *reinstate)
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
	for _, line := range p.notes {
		fmt.Fprintln(w, line)
	}
	if run == nil {
		return nil
	}
	return run(w)
}

// planStrike evaluates milestone strike's gates and writes nothing; run
// posts the line, and is nil when there is nothing to post — the ID is
// already in the state asked for, so a repeated wake is a no-op.
func planStrike(c *ctx, n int, id, decision string, reinstate bool) (*plan, func(io.Writer) error, error) {
	p := &plan{}
	milestone, err := findOpenMilestone(c, n)
	if err != nil {
		return nil, nil, err
	}
	var e error
	if milestone == nil {
		e = refuse("NOT_FOUND", "no open milestone M%d in %s", n, c.hub)
	}
	if !p.gate("milestone open", e) {
		return p.stop(strikeGates), nil, nil
	}

	body, err := c.t.IssueBody(milestone.Ref)
	if err != nil {
		return nil, nil, err
	}
	e = nil
	if !declared(tracker.RequirementIDs(body), id) || !requirementIDArg.MatchString(id) || idMilestone(id) != n {
		e = refuse("REQUIREMENT_UNDECLARED", "%s is not declared under ## Requirements of %s (M%d) — only a bold M%d-R<k> in that section is a requirement to strike", id, milestone.Ref, n, n)
	}
	if !p.gate("requirement declared", e) {
		return p.stop(strikeGates), nil, nil
	}

	e, err = verifyDecision(c, milestone, id, decision)
	if err != nil {
		return nil, nil, err
	}
	if !p.gate("decision recorded", e) {
		return p.stop(strikeGates), nil, nil
	}

	comments, err := c.t.Comments(milestone.Ref)
	if err != nil {
		return nil, nil, err
	}
	states, err := milestoneStrikes(c, milestone, comments)
	if err != nil {
		return nil, nil, err
	}
	s, has := states[id]
	struck := has && s.Word == tracker.Struck && s.Err == nil
	switch {
	case !reinstate && struck:
		p.remark("%s is already struck by %s (decision %s) — nothing to post", id, s.URL, s.Decision)
		return p, nil, nil
	case reinstate && !struck:
		p.remark("%s carries no strike to reinstate — nothing to post", id)
		return p, nil, nil
	}
	if viewer, err := c.t.Viewer(); err == nil && viewer != "" && !c.holdsRole(viewer, "coordinator") {
		p.remark("note: @%s does not hold the coordinator seat — milestone close counts struck and reinstated lines from the coordinator seat's holder only (SPEC §7); posting it anyway", viewer)
	}
	word := tracker.Struck
	if reinstate {
		word = tracker.Reinstated
	}
	line := tracker.StrikeLine(id, word, decision)
	p.would("post on %s (%s): %s", milestone.Ref, milestone.Title, line)
	run := func(w io.Writer) error {
		if err := c.t.Comment(milestone.Ref, line); err != nil {
			return err
		}
		fmt.Fprintf(w, "posted on %s (%s): %s\n", milestone.Ref, milestone.Title, line)
		return nil
	}
	return p, run, nil
}

// findOpenMilestone finds open milestone M<n> in the hub, or nil.
func findOpenMilestone(c *ctx, n int) (*tracker.Milestone, error) {
	milestones, err := c.t.OpenMilestones(c.hub)
	if err != nil {
		return nil, err
	}
	for i := range milestones {
		if got, ok := tracker.MilestoneNumber(milestones[i].Title); ok && got == n {
			return &milestones[i], nil
		}
	}
	return nil, nil
}

func declared(ids []string, id string) bool {
	for _, d := range ids {
		if d == id {
			return true
		}
	}
	return false
}

// idMilestone is the milestone number an ID carries, or -1.
func idMilestone(id string) int {
	m := requirementIDArg.FindStringSubmatch(id)
	if m == nil {
		return -1
	}
	k, _ := strconv.Atoi(m[1])
	return k
}

// verifyDecision is the one check behind every strike, the verb's and the
// close's alike: the link names a comment on the milestone issue or on one
// of its tasks, and that comment carries a Decision record — a
// **Decision:** or a **Gate resolved:**, read per paragraph as SPEC §4
// gathers them — that names the ID. The first result is the refusal
// (DECISION_UNRECORDED) when the record fails; the second is a read that
// failed, which is not the record's fault.
func verifyDecision(c *ctx, m *tracker.Milestone, id, link string) (error, error) {
	unrecorded := func(format string, args ...any) error {
		return refuse("DECISION_UNRECORDED", "%s for %s: %s — a requirement is struck or reinstated only on a recorded **Decision:** (or **Gate resolved:**) comment that names it, on %s or one of its tasks (SPEC §4)", linkOrNone(link), id, fmt.Sprintf(format, args...), m.Ref)
	}
	if link == "" {
		return unrecorded("the line carries no decision link"), nil
	}
	ref, commentID, ok := tracker.ParseCommentURL(link)
	if !ok {
		return unrecorded("not an issue-comment URL (%s)", tracker.CommentURLForm), nil
	}
	if !onMilestone(m, ref) {
		return unrecorded("%s is neither the milestone issue nor one of its tasks", ref), nil
	}
	comments, err := c.t.Comments(ref)
	if err != nil {
		return nil, err
	}
	for _, cm := range comments {
		if tracker.CommentID(cm.URL) != commentID {
			continue
		}
		decisions := 0
		for _, r := range tracker.ExtractRecords(ref, []tracker.Comment{cm}) {
			if r.Kind != "Decision" {
				continue
			}
			decisions++
			if tracker.NamesRequirement(r.Body, id) {
				return nil, nil
			}
		}
		if decisions == 0 {
			return unrecorded("the comment carries no **Decision:** or **Gate resolved:** record"), nil
		}
		return unrecorded("the comment's Decision does not name %s", id), nil
	}
	return unrecorded("no such comment on %s", ref), nil
}

func linkOrNone(link string) string {
	if link == "" {
		return "(no link)"
	}
	return link
}

// onMilestone reports whether ref is the milestone issue or one of its
// sub-issue tasks. Owner and repo names compare as GitHub compares them,
// without case.
func onMilestone(m *tracker.Milestone, ref tracker.IssueRef) bool {
	same := func(a, b tracker.IssueRef) bool {
		return a.Number == b.Number && strings.EqualFold(a.Repo, b.Repo)
	}
	if same(m.Ref, ref) {
		return true
	}
	for _, t := range m.Tasks {
		if same(t, ref) {
			return true
		}
	}
	return false
}

// strikeState is where a requirement's strike stands: the latest struck or
// reinstated line the coordinator seat's holder posted for it, and whether
// its decision link verifies (Err is DECISION_UNRECORDED when it does not).
type strikeState struct {
	tracker.Strike
	Err error
}

// milestoneStrikes reads the milestone issue's comments the way close
// counts them — only the coordinator seat's holder's lines, the latest
// comment per ID winning — and re-verifies each ID's latest line, so a
// line posted by hand is held to the verb's own check. close and status
// both read through it, so status never reports a strike close would
// refuse.
func milestoneStrikes(c *ctx, m *tracker.Milestone, comments []tracker.Comment) (map[string]strikeState, error) {
	latest := map[string]tracker.Strike{}
	var order []string
	for _, s := range tracker.ParseStrikes(comments) {
		if !c.holdsRole(s.Author, "coordinator") {
			continue
		}
		if _, ok := latest[s.ID]; !ok {
			order = append(order, s.ID)
		}
		latest[s.ID] = s
	}
	out := map[string]strikeState{}
	for _, id := range order {
		s := latest[id]
		refusal, err := verifyDecision(c, m, id, s.Decision)
		if err != nil {
			return nil, err
		}
		out[id] = strikeState{Strike: s, Err: refusal}
	}
	return out, nil
}

// struckIDs returns, in declaration order, the IDs whose latest line is a
// verified strike, and the details of those whose latest line does not
// verify.
func struckIDs(ids []string, states map[string]strikeState) (struck []string, unrecorded []string) {
	for _, id := range ids {
		s, ok := states[id]
		switch {
		case !ok:
		case s.Err != nil:
			var r refusal
			detail := s.Err.Error()
			if errors.As(s.Err, &r) {
				detail = r.Detail
			}
			unrecorded = append(unrecorded, fmt.Sprintf("%s %s line in %s: %s", id, s.Word, s.URL, detail))
		case s.Word == tracker.Struck:
			struck = append(struck, id)
		}
	}
	return struck, unrecorded
}
