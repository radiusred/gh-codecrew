package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// finishFake is the slice of the tracker task finish reads, plus a record
// of every write — a dry run must leave that record empty.
type finishFake struct {
	tracker.Tracker
	task     tracker.Task
	comments []tracker.Comment
	viewer   string
	prs      []int
	pr       tracker.PR
	writes   []string
}

func (f *finishFake) Task(tracker.IssueRef) (tracker.Task, error)          { return f.task, nil }
func (f *finishFake) Comments(tracker.IssueRef) ([]tracker.Comment, error) { return f.comments, nil }
func (f *finishFake) Viewer() (string, error)                              { return f.viewer, nil }
func (f *finishFake) ClosingPRs(tracker.IssueRef, bool) ([]int, error)     { return f.prs, nil }
func (f *finishFake) PRInfo(string, int) (tracker.PR, error)               { return f.pr, nil }
func (f *finishFake) Comment(_ tracker.IssueRef, body string) error {
	f.writes = append(f.writes, "comment: "+firstLine(body))
	return nil
}
func (f *finishFake) MergePR(_ string, n int) error { f.writes = append(f.writes, "merge"); return nil }
func (f *finishFake) MergePRBypass(_ string, n int) error {
	f.writes = append(f.writes, "merge-bypass")
	return nil
}
func (f *finishFake) DeleteBranch(_, b string) error {
	f.writes = append(f.writes, "delete "+b)
	return nil
}
func (f *finishFake) RepoInfo(string) (tracker.RepoInfo, error) {
	return tracker.RepoInfo{DefaultBranch: "main"}, nil
}

func finishCtx(f *finishFake, roles map[string]config.Role) *ctx {
	cfg := &config.Config{Codecrew: "1.0", Hub: "self", Roles: roles}
	return &ctx{cfg: cfg, roles: cfg, current: "o/r", hub: "o/r", t: f}
}

var crewRoles = map[string]config.Role{
	"implementer":     {Identity: "myorg-coder"},
	"reviewer":        {Identity: "myorg-reviewy"},
	"qa":              {},
	"doc-synthesizer": {},
}

func cleanFinish() *finishFake {
	return &finishFake{
		task:     tracker.Task{Ref: tracker.IssueRef{Repo: "o/r", Number: 7}, Labels: []string{"cc:task"}},
		comments: []tracker.Comment{{Author: "myorg-coder[bot]", Body: tracker.StartRecord("myorg-coder[bot]")}},
		viewer:   "myorg-coder[bot]",
		prs:      []int{9},
		pr: tracker.PR{Repo: "o/r", Number: 9, Author: "myorg-coder[bot]", HeadRef: "task/7-x", HeadSHA: "abc",
			Open: true, ChecksOK: true, ApprovedBy: []string{"myorg-reviewy[bot]"}, ReviewDecision: "APPROVED"},
	}
}

func gateLine(p *plan, name string) string {
	var buf bytes.Buffer
	p.print(&buf)
	for _, l := range strings.Split(buf.String(), "\n") {
		if strings.HasPrefix(l, "gate "+name+":") {
			return l
		}
	}
	return ""
}

// A clean pass: every gate ok, the actions named, nothing written until
// run — and the dry run writes nothing at all.
func TestPlanFinishCleanPassAndDryRunWritesNothing(t *testing.T) {
	f := cleanFinish()
	c := finishCtx(f, crewRoles)
	p, run, err := planFinish(c, f.task.Ref, false, false)
	if err != nil || p.refusal != nil {
		t.Fatalf("plan: err %v refusal %v", err, p.refusal)
	}
	for _, g := range []string{"task open", "no gate raised", "gates resolved", "owner", "closing PR", "CI checks", "review", "GitHub's required review"} {
		if l := gateLine(p, g); !strings.HasSuffix(l, ": ok") {
			t.Errorf("%s: %q", g, l)
		}
	}
	for _, g := range []string{"operator confirmation", "bypass actor"} {
		if l := gateLine(p, g); !strings.HasSuffix(l, ": not applicable") {
			t.Errorf("%s: %q", g, l)
		}
	}
	var buf bytes.Buffer
	p.print(&buf)
	for _, want := range []string{"would merge PR #9 (rebase); o/r#7 closes", "would delete head task/7-x", "dry run: nothing written"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("report lacks %q:\n%s", want, buf.String())
		}
	}
	if len(f.writes) != 0 {
		t.Errorf("planning wrote: %v", f.writes)
	}
	if err := run(&buf); err != nil {
		t.Fatal(err)
	}
	if strings.Join(f.writes, ",") != "merge,delete task/7-x" {
		t.Errorf("run wrote %v", f.writes)
	}
}

// A refused gate: its code, the gates after it not reached, and no actions.
func TestPlanFinishRefusalStopsInOrder(t *testing.T) {
	f := cleanFinish()
	f.viewer = "myorg-reviewy[bot]" // not the seat that started it
	p, run, err := planFinish(finishCtx(f, crewRoles), f.task.Ref, false, false)
	if err != nil || run != nil {
		t.Fatalf("err %v, run %v", err, run != nil)
	}
	var r refusal
	if !errors.As(p.refusal, &r) || r.Code != "NOT_OWNER" {
		t.Fatalf("refusal = %v", p.refusal)
	}
	if l := gateLine(p, "owner"); !strings.Contains(l, "refused[NOT_OWNER]") {
		t.Errorf("owner line: %q", l)
	}
	for _, g := range []string{"closing PR", "CI checks", "review", "GitHub's required review", "operator confirmation", "bypass actor"} {
		if l := gateLine(p, g); !strings.HasSuffix(l, ": not reached") {
			t.Errorf("%s: %q", g, l)
		}
	}
	var buf bytes.Buffer
	p.print(&buf)
	if strings.Contains(buf.String(), "would ") || !strings.Contains(buf.String(), "stops at the first refusal") {
		t.Errorf("a refused plan must list no actions:\n%s", buf.String())
	}
	if len(f.writes) != 0 {
		t.Errorf("planning wrote: %v", f.writes)
	}
}

// --bypass previewed by an operator on another seat's task: the owner
// bypass comment is an action, the merge path follows GitHub's decision.
func TestPlanFinishPreviewsBypassAndConfirmation(t *testing.T) {
	f := cleanFinish()
	f.viewer = "davison"
	f.pr.ReviewDecision = "REVIEW_REQUIRED"
	p, _, err := planFinish(finishCtx(f, crewRoles), f.task.Ref, false, true)
	if err != nil || p.refusal != nil {
		t.Fatalf("err %v refusal %v", err, p.refusal)
	}
	var buf bytes.Buffer
	p.print(&buf)
	for _, want := range []string{"would comment on PR #9: **Owner bypass:** o/r#7 was started by @myorg-coder[bot]; finished by @davison", "would comment on PR #9: **Merge bypass:**", "would merge PR #9 via the ruleset's administrator bypass"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("lacks %q:\n%s", want, buf.String())
		}
	}
	// The same preview as a crew identity refuses before any write.
	f.viewer = "myorg-reviewy[bot]"
	p, _, _ = planFinish(finishCtx(f, crewRoles), f.task.Ref, false, true)
	var r refusal
	if !errors.As(p.refusal, &r) || r.Code != "CREW_BYPASS" {
		t.Errorf("crew bypass preview = %v", p.refusal)
	}
	// Solo tier: the operator's own task, no approval, --operator-confirm.
	f = cleanFinish()
	f.viewer, f.pr.Author, f.pr.ApprovedBy, f.pr.ReviewDecision = "davison", "davison", nil, ""
	f.comments = []tracker.Comment{{Author: "davison", Body: tracker.StartRecord("davison")}}
	p, _, _ = planFinish(finishCtx(f, nil), f.task.Ref, true, false)
	if p.refusal != nil {
		t.Fatalf("solo confirm: %v", p.refusal)
	}
	buf.Reset()
	p.print(&buf)
	if !strings.Contains(buf.String(), "gate operator confirmation: ok") || !strings.Contains(buf.String(), "would comment on PR #9: **Operator confirmation:** reviewed and accepted by @davison as both author and operator") {
		t.Errorf("solo preview:\n%s", buf.String())
	}
	if len(f.writes) != 0 {
		t.Errorf("planning wrote: %v", f.writes)
	}
}

// The sweep plan names every branch and its fate and deletes nothing;
// executing it deletes exactly the planned ones (the live close's path).
func TestPlanSweepDeletesNothing(t *testing.T) {
	f := &fakeTracker{
		linked: map[int][]string{1: {"task/1-a"}},
		prs:    map[int][]int{1: {10}},
		info:   map[int]tracker.PR{10: {Number: 10, Merged: true, HeadRef: "task/1-a", HeadSHA: "s1"}},
		ahead:  map[string]int{"task/1-a": 0, "task/2-b": 3},
		tips:   map[string]string{"task/1-a": "s1", "task/2-b": "s2"},
		titles: map[int]string{1: "a", 2: "b"},
	}
	m := &tracker.Milestone{Tasks: []tracker.IssueRef{{Repo: "o/r", Number: 1}, {Repo: "o/r", Number: 2}}}
	items := planSweep(f, m)
	if len(f.deleted) != 0 {
		t.Fatalf("plan deleted %v", f.deleted)
	}
	got := map[string]bool{}
	for _, it := range items {
		got[it.Name] = it.Delete
	}
	if !got["task/1-a"] || got["task/2-b"] {
		t.Errorf("plan = %+v", items)
	}
	var buf bytes.Buffer
	if deleted := executeSweep(&buf, f, items); strings.Join(deleted, ",") != "task/1-a" || strings.Join(f.deleted, ",") != "task/1-a" {
		t.Errorf("execute deleted %v / %v", deleted, f.deleted)
	}
}

func TestPlanPrintStates(t *testing.T) {
	p := &plan{}
	p.gate("a", nil)
	p.na("b")
	p.gate("c", refuse("X", "why"))
	p.stop([]string{"a", "b", "c", "d", "e"})
	var buf bytes.Buffer
	p.print(&buf)
	want := "gate a: ok\ngate b: not applicable\ngate c: refused[X]: why\ngate d: not reached\ngate e: not reached\ndry run: nothing written — the live verb stops at the first refusal above\n"
	if buf.String() != want {
		t.Errorf("print =\n%q\nwant\n%q", buf.String(), want)
	}
	var r refusal
	if !errors.As(p.refusal, &r) || r.Code != "X" {
		t.Errorf("refusal = %v", p.refusal)
	}
}

// closeFake is the slice of the tracker milestone close reads.
type closeFake struct {
	tracker.Tracker
	milestone tracker.Milestone
	labels    []string // on the milestone issue itself
	tasks     map[int]tracker.Task
	body      string
	comments  map[int][]tracker.Comment
	hasDoc    bool
	writes    []string
}

func (f *closeFake) OpenMilestones(string) ([]tracker.Milestone, error) {
	return []tracker.Milestone{f.milestone}, nil
}
func (f *closeFake) Task(ref tracker.IssueRef) (tracker.Task, error) { return f.tasks[ref.Number], nil }
func (f *closeFake) IssueLabels(ref tracker.IssueRef) ([]string, error) {
	if ref != f.milestone.Ref {
		return nil, errors.New("labels read off the milestone issue")
	}
	return f.labels, nil
}
func (f *closeFake) IssueBody(tracker.IssueRef) (string, error) { return f.body, nil }
func (f *closeFake) Comments(ref tracker.IssueRef) ([]tracker.Comment, error) {
	return f.comments[ref.Number], nil
}
func (f *closeFake) ClosingPRs(tracker.IssueRef, bool) ([]int, error) { return nil, nil }
func (f *closeFake) HasMilestoneDoc(string, int) (bool, error)        { return f.hasDoc, nil }
func (f *closeFake) RepoInfo(string) (tracker.RepoInfo, error) {
	return tracker.RepoInfo{DefaultBranch: "main"}, nil
}
func (f *closeFake) LinkedBranches(tracker.IssueRef) ([]string, error) { return nil, nil }
func (f *closeFake) BranchAhead(_, b string) (int, string, error) {
	return 0, "", errors.New("no such branch")
}
func (f *closeFake) CloseIssue(_ tracker.IssueRef, c string) error {
	f.writes = append(f.writes, "close: "+c)
	return nil
}

// A live close that stops at DOC_MISSING still says qa is unrouted, before
// the raw material, as it did before the plan/execute split (checky's
// finding on PR #179); the dry run shows the note under its gate.
func TestPlanCloseUnroutedNoteSurvivesDocMissing(t *testing.T) {
	f := &closeFake{
		milestone: tracker.Milestone{Ref: tracker.IssueRef{Repo: "o/r", Number: 5}, Title: "M2: Two", Tasks: []tracker.IssueRef{{Repo: "o/r", Number: 6}}},
		tasks:     map[int]tracker.Task{6: {Closed: true}},
		body:      "## Goal\nx\n\n## Requirements\n- **M2-R1** — a thing\n",
		comments:  map[int][]tracker.Comment{5: {{Author: "davison", Body: "**M2-R1 — satisfied.** ran it"}}},
	}
	cfg := &config.Config{Codecrew: "1.0", Hub: "self", Roles: map[string]config.Role{"implementer": {}, "reviewer": {}, "qa": {}, "doc-synthesizer": {}}}
	c := &ctx{cfg: cfg, roles: cfg, current: "o/r", hub: "o/r", t: f}
	var live bytes.Buffer
	p, run, err := planClose(c, 2, false, &live)
	if err != nil || run != nil {
		t.Fatalf("err %v run %v", err, run != nil)
	}
	var r refusal
	if !errors.As(p.refusal, &r) || r.Code != "DOC_MISSING" {
		t.Fatalf("refusal = %v", p.refusal)
	}
	out := live.String()
	noteAt, rawAt := strings.Index(out, "note: qa is unrouted"), strings.Index(out, "raw material for docs/milestones/2-")
	if noteAt < 0 || rawAt < 0 || noteAt > rawAt {
		t.Errorf("live output must say the note before the raw material:\n%s", out)
	}
	var dry bytes.Buffer
	p, _, _ = planClose(c, 2, true, &dry)
	if dry.Len() != 0 {
		t.Errorf("dry run printed live output: %q", dry.String())
	}
	p.print(&dry)
	if !strings.Contains(dry.String(), "gate QA verdicts: ok\n  note: qa is unrouted") || !strings.Contains(dry.String(), "gate milestone document: refused[DOC_MISSING]") {
		t.Errorf("dry report:\n%s", dry.String())
	}
	if len(f.writes) != 0 {
		t.Errorf("planning wrote: %v", f.writes)
	}
	// With the document merged, the plan names the close and run performs it.
	f.hasDoc = true
	p, run, err = planClose(c, 2, true, &dry)
	if err != nil || p.refusal != nil {
		t.Fatalf("clean: err %v refusal %v", err, p.refusal)
	}
	if err := run(&live); err != nil || strings.Join(f.writes, ";") != "close: Closed by `gh codecrew milestone close 2`: all 1 tasks done, milestone document merged." {
		t.Errorf("run: err %v writes %v", err, f.writes)
	}
}

// A cc:needs-decision on the milestone issue itself — a requirement-level
// question raised there (#200) — once let the milestone close under it
// (#219). The close now refuses MILESTONE_GATED before it counts tasks, in
// the shape of task finish's "no gate raised" gate: the dry run prints the
// gate line and marks the rest not reached, nothing is written, and the
// same milestone without the label passes the gate (M12-R3).
func TestPlanCloseRefusesWhileMilestoneIssueIsGated(t *testing.T) {
	f := &closeFake{
		milestone: tracker.Milestone{Ref: tracker.IssueRef{Repo: "o/r", Number: 5}, Title: "M2: Two", Tasks: []tracker.IssueRef{{Repo: "o/r", Number: 6}}},
		labels:    []string{tracker.LabelMilestone, "CC:Needs-Decision"}, // GitHub compares label names case-insensitively
		tasks:     map[int]tracker.Task{6: {Closed: false}},              // an open task, so the ordering shows
		body:      "## Goal\nx\n\n## Requirements\n- **M2-R1** — a thing\n",
		comments:  map[int][]tracker.Comment{5: {{Author: "davison", Body: "**M2-R1 — satisfied.** ran it"}}},
		hasDoc:    true,
	}
	cfg := &config.Config{Codecrew: "1.0", Hub: "self", Roles: map[string]config.Role{"implementer": {}, "reviewer": {}, "qa": {}, "doc-synthesizer": {}}}
	c := &ctx{cfg: cfg, roles: cfg, current: "o/r", hub: "o/r", t: f}
	var live bytes.Buffer
	p, run, err := planClose(c, 2, false, &live)
	if err != nil || run != nil {
		t.Fatalf("err %v run %v", err, run != nil)
	}
	var r refusal
	if !errors.As(p.refusal, &r) || r.Code != "MILESTONE_GATED" {
		t.Fatalf("refusal = %v", p.refusal)
	}
	for _, want := range []string{"o/r#5", tracker.LabelNeedsDecision, "**Gate resolved:**"} {
		if !strings.Contains(p.refusal.Error(), want) {
			t.Errorf("detail must say %q: %v", want, p.refusal)
		}
	}
	if live.Len() != 0 {
		t.Errorf("a gated close printed: %q", live.String())
	}
	var dry bytes.Buffer
	p, _, _ = planClose(c, 2, true, &dry)
	p.print(&dry)
	for _, want := range []string{
		"gate milestone open: ok\ngate no gate raised: refused[MILESTONE_GATED]: o/r#5 has cc:needs-decision raised on the milestone issue",
		"\ngate tasks closed: not reached\n",
		"gate milestone document: not reached\ndry run: nothing written — the live verb stops at the first refusal above\n",
	} {
		if !strings.Contains(dry.String(), want) {
			t.Errorf("dry report must contain %q:\n%s", want, dry.String())
		}
	}
	if len(f.writes) != 0 {
		t.Errorf("planning wrote: %v", f.writes)
	}
	// Label gone: the gate passes and the close reaches the tasks, where
	// the open one is what refuses now.
	f.labels = []string{tracker.LabelMilestone}
	p, _, err = planClose(c, 2, true, &dry)
	if err != nil {
		t.Fatal(err)
	}
	if !errors.As(p.refusal, &r) || r.Code != "OPEN_TASKS" {
		t.Errorf("without the label the close must reach the tasks, got %v", p.refusal)
	}
	dry.Reset()
	p.print(&dry)
	if !strings.Contains(dry.String(), "gate no gate raised: ok\ngate tasks closed: refused[OPEN_TASKS]") {
		t.Errorf("dry report:\n%s", dry.String())
	}
	// And with the task closed too, the plan is clean and run closes.
	f.tasks[6] = tracker.Task{Closed: true}
	p, run, err = planClose(c, 2, false, &live)
	if err != nil || p.refusal != nil {
		t.Fatalf("clean: err %v refusal %v", err, p.refusal)
	}
	if err := run(&live); err != nil || len(f.writes) != 1 || !strings.HasPrefix(f.writes[0], "close: ") {
		t.Errorf("run: err %v writes %v", err, f.writes)
	}
}

// An App whose installation token cannot read the checks at all (a private
// repo, no checks: read or actions: read — #198): the CI checks gate
// refuses with a code naming the App and the permission, the gates after
// it are not reached, the dry run writes nothing — and the raw GraphQL
// error never reaches the operator.
func TestPlanFinishNoChecksPermissionNamesAppAndPermission(t *testing.T) {
	for _, perm := range []string{"checks: read", "actions: read"} {
		f := cleanFinish()
		f.pr.ChecksOK = false
		f.pr.ChecksUnreadable = perm
		p, run, err := planFinish(finishCtx(f, crewRoles), f.task.Ref, false, false)
		if err != nil || run != nil {
			t.Fatalf("%s: err %v, run %v", perm, err, run != nil)
		}
		var r refusal
		if !errors.As(p.refusal, &r) || r.Code != "NO_CHECKS_PERMISSION" {
			t.Fatalf("%s: refusal = %v", perm, p.refusal)
		}
		l := gateLine(p, "CI checks")
		for _, want := range []string{"refused[NO_CHECKS_PERMISSION]", "PR #9", "App myorg-coder", "`" + perm + "`", "private repo", "settings page", "installation", "docs/identities.md"} {
			if !strings.Contains(l, want) {
				t.Errorf("%s: CI checks line lacks %q: %q", perm, want, l)
			}
		}
		for _, g := range []string{"review", "GitHub's required review", "operator confirmation", "bypass actor"} {
			if gl := gateLine(p, g); !strings.HasSuffix(gl, ": not reached") {
				t.Errorf("%s: %s: %q", perm, g, gl)
			}
		}
		if len(f.writes) != 0 {
			t.Errorf("%s: planning wrote: %v", perm, f.writes)
		}
	}
	// A human's token meeting the same failure is named by login.
	if got := seatName("davison"); got != "@davison" {
		t.Errorf("seatName(human) = %q", got)
	}
}
