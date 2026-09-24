package cli

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/tracker"
	"github.com/radiusred/gh-codecrew/internal/tracker/faketracker"
)

// strikeWorld is one open milestone, M8 on o/r#5, with one task, o/r#6,
// and the comments each issue carries. qa routes to an App, so a verdict
// and a strike have different authors; coordinator is the operator.
type strikeWorld struct {
	body     string
	comments map[int][]tracker.Comment
	viewer   string
	posted   []string
	closed   []string
	hasDoc   bool
}

var (
	strikeMilestone = tracker.IssueRef{Repo: "o/r", Number: 5}
	strikeTask      = tracker.IssueRef{Repo: "o/r", Number: 6}
)

func commentOn(n, id int) string {
	return fmt.Sprintf("https://github.com/o/r/issues/%d#issuecomment-%d", n, id)
}

func newStrikeWorld() *strikeWorld {
	return &strikeWorld{
		body: "## Goal\nx\n\n## Requirements\n- **M8-R1** — one\n- **M8-R2** — two\n",
		comments: map[int][]tracker.Comment{
			5: {{Author: "davison", URL: commentOn(5, 50), Body: "**Decision:** strike M8-R2: the store needs a domain we do not own.\n\n**Rejected:** a false satisfied."}},
			6: {
				{Author: "davison", URL: commentOn(6, 60), Body: "**Gate resolved:** M8-R2 is out of scope for this milestone."},
				{Author: "davison", URL: commentOn(6, 61), Body: "**Decision:** use the other schema."},
				{Author: "davison", URL: commentOn(6, 62), Body: "**Deviation:** M8-R2 slipped.\n\n**Why:** reasons."},
				{Author: "davison", URL: commentOn(6, 63), Body: "Just a remark about M8-R2."},
			},
		},
		viewer: "davison",
		hasDoc: true,
	}
}

func (s *strikeWorld) ctx(t *testing.T) *ctx {
	t.Helper()
	v := &faketracker.Venue{
		OpenMilestonesFn: func(string) ([]tracker.Milestone, error) {
			return []tracker.Milestone{{Ref: strikeMilestone, Title: "M8: Eight", Tasks: []tracker.IssueRef{strikeTask}}}, nil
		},
		IssueBodyFn:   func(tracker.IssueRef) (string, error) { return s.body, nil },
		IssueLabelsFn: func(tracker.IssueRef) ([]string, error) { return []string{tracker.LabelMilestone}, nil },
		TaskFn:        func(tracker.IssueRef) (tracker.Task, error) { return tracker.Task{Closed: true}, nil },
		CommentsFn: func(ref tracker.IssueRef) ([]tracker.Comment, error) {
			return s.comments[ref.Number], nil
		},
		CommentFn: func(ref tracker.IssueRef, body string) error {
			s.posted = append(s.posted, fmt.Sprintf("%s: %s", ref, body))
			return nil
		},
		CloseIssueFn: func(_ tracker.IssueRef, c string) error {
			s.closed = append(s.closed, c)
			return nil
		},
		ViewerFn:          func() (string, error) { return s.viewer, nil },
		HasMilestoneDocFn: func(string, int) (bool, error) { return s.hasDoc, nil },
		RepoInfoFn: func(string) (tracker.RepoInfo, error) {
			return tracker.RepoInfo{DefaultBranch: "main", DeleteBranchOnMerge: true}, nil
		},
	}
	cfg := &config.Config{Codecrew: "2.0", Hub: "self", Dir: t.TempDir(), Roles: map[string]config.Role{
		"implementer":     {Identity: config.Identity{Kind: config.KindApp, Value: "radiusred-cody"}},
		"reviewer":        {},
		"qa":              {Identity: config.Identity{Kind: config.KindApp, Value: "radiusred-testy"}},
		"doc-synthesizer": {},
		"coordinator":     {},
	}}
	return &ctx{cfg: cfg, roles: cfg, current: "o/r", hub: "o/r", t: v}
}

// strike runs the verb's plan and, when it has something to post, the post.
func (s *strikeWorld) strike(t *testing.T, id, link string, reinstate bool) (string, error) {
	t.Helper()
	c := s.ctx(t)
	p, run, err := planStrike(c, 8, id, link, reinstate)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if p.refusal != nil {
		return "", p.refusal
	}
	for _, n := range p.notes {
		fmt.Fprintln(&out, n)
	}
	if run != nil {
		if err := run(&out); err != nil {
			t.Fatal(err)
		}
	}
	return out.String(), nil
}

func strikeCode(err error) string {
	var r refusal
	if errors.As(err, &r) {
		return r.Code
	}
	return ""
}

// The verb posts exactly the line R4 names, on the milestone issue, and
// never edits the body; a Decision on the milestone issue and a Gate
// resolved on a task both count.
func TestMilestoneStrikePostsTheLine(t *testing.T) {
	s := newStrikeWorld()
	out, err := s.strike(t, "M8-R2", commentOn(5, 50), false)
	if err != nil {
		t.Fatal(err)
	}
	want := "o/r#5: **M8-R2 — struck.** " + commentOn(5, 50)
	if len(s.posted) != 1 || s.posted[0] != want {
		t.Fatalf("posted %v, want [%s]", s.posted, want)
	}
	if !strings.Contains(out, "posted on o/r#5 (M8: Eight): **M8-R2 — struck.**") {
		t.Errorf("receipt:\n%s", out)
	}
	s2 := newStrikeWorld()
	if _, err := s2.strike(t, "M8-R2", commentOn(6, 60), false); err != nil || len(s2.posted) != 1 {
		t.Errorf("a **Gate resolved:** on a task must count: err %v posted %v", err, s2.posted)
	}
}

// Every way the record can fail refuses before anything is posted.
func TestMilestoneStrikeRefusals(t *testing.T) {
	for _, tc := range []struct {
		name, id, link, code, detail string
	}{
		{"undeclared", "M8-R3", commentOn(5, 50), "REQUIREMENT_UNDECLARED", "M8-R3 is not declared"},
		{"another milestone's ID", "M7-R2", commentOn(5, 50), "REQUIREMENT_UNDECLARED", "M7-R2"},
		{"not a comment URL", "M8-R2", "https://github.com/o/r/issues/5", "DECISION_UNRECORDED", "not an issue-comment URL"},
		{"a foreign issue", "M8-R2", commentOn(9, 90), "DECISION_UNRECORDED", "o/r#9 is neither the milestone issue nor one of its tasks"},
		{"no such comment", "M8-R2", commentOn(6, 99), "DECISION_UNRECORDED", "no such comment on o/r#6"},
		{"a Decision that does not name the ID", "M8-R2", commentOn(6, 61), "DECISION_UNRECORDED", "does not name M8-R2"},
		{"a Deviation", "M8-R2", commentOn(6, 62), "DECISION_UNRECORDED", "carries no **Decision:** or **Gate resolved:** record"},
		{"no record at all", "M8-R2", commentOn(6, 63), "DECISION_UNRECORDED", "carries no **Decision:**"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := newStrikeWorld()
			_, err := s.strike(t, tc.id, tc.link, false)
			if strikeCode(err) != tc.code || !strings.Contains(err.Error(), tc.detail) {
				t.Errorf("err = %v, want %s containing %q", err, tc.code, tc.detail)
			}
			if len(s.posted) != 0 {
				t.Errorf("a refused strike posted %v", s.posted)
			}
		})
	}
	s := newStrikeWorld()
	c := s.ctx(t)
	p, _, err := planStrike(c, 9, "M9-R1", commentOn(5, 50), false)
	if err != nil || strikeCode(p.refusal) != "NOT_FOUND" {
		t.Errorf("no open M9: err %v refusal %v", err, p.refusal)
	}
}

// --dry-run prints every gate and the line, writes nothing, and a refused
// dry run marks the rest not reached.
func TestMilestoneStrikeDryRun(t *testing.T) {
	s := newStrikeWorld()
	c := s.ctx(t)
	p, _, err := planStrike(c, 8, "M8-R2", commentOn(5, 50), false)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	p.print(&out)
	want := "gate milestone open: ok\ngate requirement declared: ok\ngate decision recorded: ok\nwould post on o/r#5 (M8: Eight): **M8-R2 — struck.** " + commentOn(5, 50) + "\ndry run: nothing written\n"
	if out.String() != want {
		t.Errorf("dry run =\n%s\nwant\n%s", out.String(), want)
	}
	p, _, _ = planStrike(c, 8, "M8-R3", commentOn(5, 50), false)
	out.Reset()
	p.print(&out)
	if !strings.Contains(out.String(), "gate requirement declared: refused[REQUIREMENT_UNDECLARED]") || !strings.Contains(out.String(), "gate decision recorded: not reached") {
		t.Errorf("refused dry run:\n%s", out.String())
	}
	if len(s.posted) != 0 {
		t.Errorf("a dry run posted %v", s.posted)
	}
}

// A repeated wake is a no-op: an ID already struck posts nothing and says
// so; reinstating posts the twin line, after which the ID can be struck
// again; reinstating what is not struck posts nothing.
func TestMilestoneStrikeIsIdempotentAndReinstates(t *testing.T) {
	s := newStrikeWorld()
	out, err := s.strike(t, "M8-R2", commentOn(5, 50), true)
	if err != nil || len(s.posted) != 0 || !strings.Contains(out, "carries no strike to reinstate") {
		t.Fatalf("reinstate unstruck: err %v posted %v out %q", err, s.posted, out)
	}
	s.comments[5] = append(s.comments[5], tracker.Comment{Author: "davison", URL: commentOn(5, 51), Body: "**M8-R2 — struck.** " + commentOn(5, 50)})
	out, err = s.strike(t, "M8-R2", commentOn(5, 50), false)
	if err != nil || len(s.posted) != 0 || !strings.Contains(out, "M8-R2 is already struck by "+commentOn(5, 51)) {
		t.Fatalf("re-strike: err %v posted %v out %q", err, s.posted, out)
	}
	if _, err = s.strike(t, "M8-R2", commentOn(6, 60), true); err != nil {
		t.Fatal(err)
	}
	if len(s.posted) != 1 || s.posted[0] != "o/r#5: **M8-R2 — reinstated.** "+commentOn(6, 60) {
		t.Errorf("reinstate posted %v", s.posted)
	}
}

// The verb checks the record, not the actor (settled on #348), but says
// when the line it posts is one close will not count.
func TestMilestoneStrikeNotesANonCoordinatorPost(t *testing.T) {
	s := newStrikeWorld()
	s.viewer = "radiusred-cody[bot]"
	out, err := s.strike(t, "M8-R2", commentOn(5, 50), false)
	if err != nil || len(s.posted) != 1 {
		t.Fatalf("err %v posted %v", err, s.posted)
	}
	if !strings.Contains(out, "note: @radiusred-cody[bot] does not hold the coordinator seat") {
		t.Errorf("no note:\n%s", out)
	}
}

// close runs milestone close's plan and returns the refusal, or closes.
func (s *strikeWorld) close(t *testing.T) (*plan, error) {
	t.Helper()
	c := s.ctx(t)
	var live bytes.Buffer
	p, run, err := planClose(c, 8, false, &live)
	if err != nil {
		t.Fatal(err)
	}
	if p.refusal != nil {
		return p, p.refusal
	}
	return p, run(&live)
}

func (s *strikeWorld) say(n int, author, body string) {
	s.comments[n] = append(s.comments[n], tracker.Comment{Author: author, URL: commentOn(n, 100+len(s.comments[n])), Body: body})
}

const testy = "radiusred-testy[bot]"

// A verified strike is terminal: it asks for no QA verdict, it outranks a
// not-satisfied verdict before it and a verdict after it, and the closing
// comment names it.
func TestCloseCountsAVerifiedStrikeAsTerminal(t *testing.T) {
	s := newStrikeWorld()
	s.say(5, testy, "**M8-R1 — satisfied.** ran it\n**M8-R2 — not satisfied.** broken")
	if _, err := s.close(t); strikeCode(err) != "VERDICT_UNSATISFIED" {
		t.Fatalf("before the strike: %v", err)
	}
	s.say(5, "davison", "**M8-R2 — struck.** "+commentOn(5, 50))
	if _, err := s.close(t); err != nil {
		t.Fatalf("struck after a not-satisfied verdict must close: %v", err)
	}
	if len(s.closed) != 1 || !strings.Contains(s.closed[0], "Struck by recorded decision: M8-R2.") {
		t.Errorf("closing comment: %v", s.closed)
	}
	s.say(5, testy, "**M8-R2 — not satisfied.** still broken")
	if _, err := s.close(t); err != nil {
		t.Errorf("a verdict after the strike does not undo it: %v", err)
	}
}

// Reinstated, the requirement is back in scope and wants its verdict.
func TestCloseAfterReinstateWantsAVerdict(t *testing.T) {
	s := newStrikeWorld()
	s.say(5, testy, "**M8-R1 — satisfied.** ran it")
	s.say(5, "davison", "**M8-R2 — struck.** "+commentOn(5, 50))
	s.say(6, "davison", "**Decision:** bring M8-R2 back; the domain arrived.")
	s.say(5, "davison", "**M8-R2 — reinstated.** "+commentOn(6, 100+len(s.comments[6])-1))
	_, err := s.close(t)
	if strikeCode(err) != "VERDICT_MISSING" || !strings.Contains(err.Error(), "M8-R2") {
		t.Errorf("reinstated without a verdict: %v", err)
	}
	s.say(5, testy, "**M8-R2 — satisfied.** ran it")
	if _, err := s.close(t); err != nil {
		t.Errorf("reinstated and satisfied: %v", err)
	}
}

// A line posted by hand is held to the verb's check: one whose link does
// not verify refuses DECISION_UNRECORDED ahead of the verdicts, and one
// that verifies counts.
func TestCloseReverifiesAHandPostedStrike(t *testing.T) {
	s := newStrikeWorld()
	s.say(5, testy, "**M8-R1 — satisfied.** ran it")
	s.say(5, "davison", "**M8-R2 — struck.** we decided, trust me")
	_, err := s.close(t)
	if strikeCode(err) != "DECISION_UNRECORDED" || !strings.Contains(err.Error(), "the line carries no decision link") {
		t.Fatalf("unlinked: %v", err)
	}
	s.say(5, "davison", "**M8-R2 — struck.** "+commentOn(6, 61))
	if _, err = s.close(t); strikeCode(err) != "DECISION_UNRECORDED" || !strings.Contains(err.Error(), "does not name M8-R2") {
		t.Fatalf("wrong decision: %v", err)
	}
	s.say(5, "davison", "Struck by hand, properly this time:\n\n**M8-R2 — struck.** [decision]("+commentOn(6, 60)+")")
	if _, err = s.close(t); err != nil {
		t.Errorf("a verified hand-posted strike must count: %v", err)
	}
}

// Only the coordinator seat's holder strikes: the same line from the qa
// App, or from the implementer, counts for nothing.
func TestCloseIgnoresAStrikeFromANonCoordinator(t *testing.T) {
	for _, author := range []string{testy, "radiusred-cody[bot]"} {
		s := newStrikeWorld()
		s.say(5, testy, "**M8-R1 — satisfied.** ran it")
		s.say(5, author, "**M8-R2 — struck.** "+commentOn(5, 50))
		if _, err := s.close(t); strikeCode(err) != "VERDICT_MISSING" {
			t.Errorf("a strike by %s counted: %v", author, err)
		}
	}
}

// A still-bold struck-through ID is still a requirement: with no strike
// line behind it, close wants its verdict (#348's test case).
func TestCloseStillCountsABodyStruckThroughID(t *testing.T) {
	s := newStrikeWorld()
	s.body = "## Requirements\n- **M8-R1** — one\n- ~~**M8-R2** — two~~\n"
	s.say(5, testy, "**M8-R1 — satisfied.** ran it")
	if _, err := s.close(t); strikeCode(err) != "VERDICT_MISSING" || !strings.Contains(err.Error(), "M8-R2") {
		t.Errorf("a body strikethrough must not strike: %v", err)
	}
}

// The raw material now carries the milestone issue's own records, where a
// strike decision lives (M18-R4).
func TestCloseGathersTheMilestoneIssuesRecords(t *testing.T) {
	s := newStrikeWorld()
	s.hasDoc = false
	s.say(5, testy, "**M8-R1 — satisfied.** ran it")
	s.say(5, "davison", "**M8-R2 — struck.** "+commentOn(5, 50))
	c := s.ctx(t)
	var live bytes.Buffer
	p, _, err := planClose(c, 8, false, &live)
	if err != nil || strikeCode(p.refusal) != "DOC_MISSING" {
		t.Fatalf("err %v refusal %v", err, p.refusal)
	}
	if !strings.Contains(live.String(), "**Decision:** on o/r#5 by @davison ("+commentOn(5, 50)+")") {
		t.Errorf("the milestone issue's Decision is not in the raw material:\n%s", live.String())
	}
	var dry bytes.Buffer
	p, _, _ = planClose(c, 8, true, &dry)
	p.print(&dry)
	if !strings.Contains(dry.String(), "gate QA verdicts: ok\n  struck: M8-R2 — decision "+commentOn(5, 50)) {
		t.Errorf("dry run does not show the strike:\n%s", dry.String())
	}
}

// status reports each strike as close counts it, notes one that does not
// verify, and notes a body strikethrough with no strike behind it.
func TestStatusReportsStrikes(t *testing.T) {
	s := newStrikeWorld()
	s.body = "## Requirements\n- **M8-R1** — one\n- **M8-R2** — two\n- ~~**M8-R3** — three~~\n- ~~**M8-R4**~~\n"
	s.say(5, "davison", "**M8-R2 — struck.** "+commentOn(5, 50))
	s.say(5, "davison", "**M8-R1 — struck.** no link")
	s.say(6, "davison", "**Decision:** M8-R4 is out.")
	s.say(5, "davison", "**M8-R4 — struck.** "+commentOn(6, 104))
	var out bytes.Buffer
	if err := statusReport(&out, s.ctx(t)); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{
		"  struck: M8-R2 — decision " + commentOn(5, 50) + "\n",
		"  struck: M8-R4 — decision " + commentOn(6, 104) + "\n",
		"  note: a strike does not verify — milestone close refuses DECISION_UNRECORDED: M8-R1 struck line in ",
		"  note: M8-R3 is struck through in the body, which the protocol does not read",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("status lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "note: M8-R4 is struck through") {
		t.Errorf("a struck-through ID that is struck needs no note:\n%s", got)
	}
}
