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
	task      tracker.Task
	body      string               // the task issue body — its ## Adopts section is what finish reads
	issues    map[int]tracker.Task // adopted captures, by number; anything else answers f.task
	mergeSHA  string
	comments  []tracker.Comment
	viewer    string
	prs       []int
	pr        tracker.PR
	closes    []tracker.TitledIssue // what GitHub says the PR would close
	closesErr error                 // when set, what the closing-references read fails with
	writes    []string
	closed    []string // the closing comment of every CloseIssue, in order
}

func (f *finishFake) Task(ref tracker.IssueRef) (tracker.Task, error) {
	if issue, ok := f.issues[ref.Number]; ok {
		return issue, nil
	}
	return f.task, nil
}
func (f *finishFake) IssueBody(tracker.IssueRef) (string, error) { return f.body, nil }
func (f *finishFake) MergeCommit(string, int) (string, error)    { return f.mergeSHA, nil }
func (f *finishFake) CloseIssue(ref tracker.IssueRef, c string) error {
	f.writes = append(f.writes, "close "+ref.String())
	f.closed = append(f.closed, c)
	return nil
}
func (f *finishFake) Comments(tracker.IssueRef) ([]tracker.Comment, error) { return f.comments, nil }
func (f *finishFake) Viewer() (string, error)                              { return f.viewer, nil }
func (f *finishFake) ClosingPRs(tracker.IssueRef, bool) ([]int, error)     { return f.prs, nil }
func (f *finishFake) PRInfo(string, int) (tracker.PR, error)               { return f.pr, nil }
func (f *finishFake) ClosingReferences(string, int) ([]tracker.TitledIssue, error) {
	if f.closesErr != nil {
		return nil, f.closesErr
	}
	return f.closes, nil
}
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
	"implementer":     {Identity: config.Identity{Kind: config.KindApp, Value: "myorg-coder"}},
	"reviewer":        {Identity: config.Identity{Kind: config.KindApp, Value: "myorg-reviewy"}},
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
		// What a well-formed PR body leaves GitHub holding: the task, alone.
		closes: []tracker.TitledIssue{{Ref: tracker.IssueRef{Repo: "o/r", Number: 7}, Title: "Seven"}},
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

// A task nobody started. Until 2.0 StartedBy fell back to the first
// assignee and, failing that, an empty owner waved the gate through — so
// an assigned-but-never-started task could be finished by anyone. The
// shim is gone (M13-R7): the gate refuses NOT_OWNER, and its detail names
// task start rather than an owner that does not exist. --bypass still
// overrides it, and its record says there was no start rather than naming
// a phantom starter.
func TestPlanFinishRefusesATaskNobodyStarted(t *testing.T) {
	f := cleanFinish()
	f.task.Assignees = []string{"davison"} // an assignee is not a start record
	f.comments = []tracker.Comment{{Author: "davison", Body: "Plan looks right, picking this up."}}
	p, run, err := planFinish(finishCtx(f, crewRoles), f.task.Ref, false, false)
	if err != nil || run != nil {
		t.Fatalf("err %v, run %v", err, run != nil)
	}
	var r refusal
	if !errors.As(p.refusal, &r) || r.Code != "NOT_OWNER" {
		t.Fatalf("refusal = %v", p.refusal)
	}
	for _, want := range []string{"nothing records a start on o/r#7", "gh codecrew task start 7"} {
		if !strings.Contains(r.Detail, want) {
			t.Errorf("detail %q lacks %q", r.Detail, want)
		}
	}
	if strings.Contains(r.Detail, "@davison") {
		t.Errorf("the refusal named the assignee as an owner: %q", r.Detail)
	}

	// The operator's override: recorded, and worded for a task with no
	// start record.
	f.viewer = "davison"
	p, _, err = planFinish(finishCtx(f, crewRoles), f.task.Ref, false, true)
	if err != nil || p.refusal != nil {
		t.Fatalf("bypass: err %v refusal %v", err, p.refusal)
	}
	var buf bytes.Buffer
	p.print(&buf)
	if want := "would comment on PR #9: **Owner bypass:** nothing records a start on o/r#7; finished by @davison"; !strings.Contains(buf.String(), want) {
		t.Errorf("lacks %q:\n%s", want, buf.String())
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
	branches  map[string][]string // the stale sweep's per-repo listing
	openPRs   map[string][]int    // open PRs by head branch, whatever they close
	ahead     map[string]int
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
func (f *closeFake) TaskBranches(repo string) ([]string, bool, error) {
	return f.branches[repo], false, nil
}
func (f *closeFake) OpenPRsForBranch(_, branch string) ([]int, error) { return f.openPRs[branch], nil }
func (f *closeFake) BranchAhead(_, b string) (int, string, error) {
	n, ok := f.ahead[b]
	if !ok {
		return 0, "", errors.New("no such branch")
	}
	return n, "", nil
}
func (f *closeFake) DeleteBranch(_, b string) error {
	f.writes = append(f.writes, "delete: "+b)
	return nil
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

// A `user:`-typed holder is a human who happens to hold a seat: the
// operator's acts stay open to them. A 1.0 table could not say which kind
// of principal a routed login was, so every one of them was refused as a
// crew identity — the bug M13-R4 exists to fix (#254, Claude scan
// finding 4).
func TestHumanSeatHolderIsNotCrew(t *testing.T) {
	humanReviewer := map[string]config.Role{
		"implementer": {Identity: config.ParseIdentity("app:myorg-coder")},
		"reviewer":    {Identity: config.ParseIdentity("user:alice")},
		"qa":          {Identity: config.ParseIdentity("team:myorg/qa-crew")},
	}
	// --bypass: alice holds the reviewer seat and is still an operator.
	f := cleanFinish()
	f.viewer = "alice"
	f.pr.ReviewDecision = "REVIEW_REQUIRED"
	f.pr.ApprovedBy = []string{"alice"}
	p, _, err := planFinish(finishCtx(f, humanReviewer), f.task.Ref, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if p.refusal != nil {
		t.Errorf("a user:-typed holder was refused --bypass: %v", p.refusal)
	}
	// The App-typed holder is still crew, on the same table.
	f = cleanFinish()
	f.viewer = "myorg-coder[bot]"
	f.pr.ReviewDecision = "REVIEW_REQUIRED"
	f.pr.ApprovedBy = []string{"alice"}
	p, _, _ = planFinish(finishCtx(f, humanReviewer), f.task.Ref, false, true)
	var r refusal
	if !errors.As(p.refusal, &r) || r.Code != "CREW_BYPASS" {
		t.Errorf("app:-typed holder's bypass = %v, want CREW_BYPASS", p.refusal)
	}
	// --operator-confirm: the same split on the confirmation gate, with
	// the reviewer seat unrouted so the solo path is the one under test.
	soloReview := map[string]config.Role{
		"implementer": {Identity: config.ParseIdentity("app:myorg-coder")},
		"reviewer":    {},
		"qa":          {Identity: config.ParseIdentity("user:alice")},
	}
	f = cleanFinish()
	f.viewer, f.pr.Author, f.pr.ApprovedBy, f.pr.ReviewDecision = "alice", "alice", nil, ""
	f.comments = []tracker.Comment{{Author: "alice", Body: tracker.StartRecord("alice")}}
	p, _, _ = planFinish(finishCtx(f, soloReview), f.task.Ref, true, false)
	if p.refusal != nil {
		t.Errorf("a user:-typed holder was refused --operator-confirm: %v", p.refusal)
	}
	f.viewer, f.pr.Author = "myorg-coder[bot]", "myorg-coder[bot]"
	f.comments = []tracker.Comment{{Author: "myorg-coder[bot]", Body: tracker.StartRecord("myorg-coder[bot]")}}
	p, _, _ = planFinish(finishCtx(f, soloReview), f.task.Ref, true, false)
	if !errors.As(p.refusal, &r) || r.Code != "SELF_CONFIRM" {
		t.Errorf("app:-typed holder's confirmation = %v, want SELF_CONFIRM", p.refusal)
	}
}

// The captures a task adopted are closed by the merge that delivers them,
// so no PR body has to remember a Closes line (#193). The dry run lists
// each one beside the branch cleanup and writes nothing; the live run
// closes the open capture with a comment naming the task, the PR and the
// commit the merge left, and reports the already-closed one as a note
// rather than failing after a merge that has happened.
func TestPlanFinishClosesTheAdoptedCaptures(t *testing.T) {
	f := cleanFinish()
	f.body = "## Goal\ng\n\n## Adopts\n- #193 — the capture this task adopts\n- #194 — one closed since\n\n## Plan\np\n"
	f.issues = map[int]tracker.Task{194: {Closed: true}}
	f.mergeSHA = "d1ff00d"
	c := finishCtx(f, crewRoles)

	p, run, err := planFinish(c, f.task.Ref, false, false)
	if err != nil || p.refusal != nil {
		t.Fatalf("plan: err %v refusal %v", err, p.refusal)
	}
	var buf bytes.Buffer
	p.print(&buf)
	for _, want := range []string{"would close o/r#193 (adopted)", "would skip o/r#194: already closed (adopted)", "would delete head task/7-x"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("report lacks %q:\n%s", want, buf.String())
		}
	}
	if len(f.writes) != 0 {
		t.Fatalf("planning wrote: %v", f.writes)
	}

	buf.Reset()
	if err := run(&buf); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(f.writes, ","); got != "merge,close o/r#193,delete task/7-x" {
		t.Errorf("run wrote %v", f.writes)
	}
	if len(f.closed) != 1 || f.closed[0] != tracker.AdoptionClose(f.task.Ref, tracker.IssueRef{Repo: "o/r", Number: 9}, "d1ff00d") {
		t.Fatalf("closing comment %q", f.closed)
	}
	for _, want := range []string{"o/r#7", "o/r#9", "d1ff00d"} {
		if !strings.Contains(f.closed[0], want) {
			t.Errorf("closing comment does not name %s: %s", want, f.closed[0])
		}
	}
	for _, want := range []string{"closed o/r#193 (adopted by o/r#7)\n", "note: o/r#194 was already closed (adopted by o/r#7)\n"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, buf.String())
		}
	}
}

// A task that adopts nothing reads its body, finds no section, and does
// nothing further: no close, and no merge-commit read either.
func TestPlanFinishWithoutAdoptionsClosesNothing(t *testing.T) {
	f := cleanFinish()
	f.body = "## Goal\ng\n\n## Plan\np\n"
	p, run, err := planFinish(finishCtx(f, crewRoles), f.task.Ref, false, false)
	if err != nil || p.refusal != nil {
		t.Fatalf("plan: err %v refusal %v", err, p.refusal)
	}
	var buf bytes.Buffer
	if err := run(&buf); err != nil {
		t.Fatal(err)
	}
	if strings.Join(f.writes, ",") != "merge,delete task/7-x" {
		t.Errorf("run wrote %v", f.writes)
	}
	if strings.Contains(buf.String(), "adopted") {
		t.Errorf("output mentions adoption:\n%s", buf.String())
	}
}

// GitHub reports no merge commit for a PR it has only just merged, and
// returns it as an empty string rather than an error. The closing comments
// go out without a SHA, and the output says so — the same note a failed
// read gets, rather than a comment quietly missing half its back-reference.
func TestCloseAdoptedNotesAMergeCommitItCannotName(t *testing.T) {
	f := cleanFinish()
	f.body = "## Adopts\n- #193 — the capture\n\n## Plan\np\n"
	f.mergeSHA = ""
	_, run, err := planFinish(finishCtx(f, crewRoles), f.task.Ref, false, false)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := run(&buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "note: GitHub reports no merge commit for o/r#9 yet") {
		t.Errorf("output lacks the note:\n%s", buf.String())
	}
	if len(f.closed) != 1 || strings.Contains(f.closed[0], "merged as") {
		t.Errorf("closing comment %q", f.closed)
	}
}

// A close sweeps the task branches earlier closes left behind as well as
// its own (#167, M14-R4): a closed task's empty branch goes and is named in
// the closing comment under its own sentence, a closed task's unmerged work
// and an open task's branch are named and left, and --dry-run lists all of
// it beside the milestone's own while writing nothing.
func TestPlanCloseSweepsStaleBranchesFromEarlierCloses(t *testing.T) {
	f := &closeFake{
		milestone: tracker.Milestone{Ref: tracker.IssueRef{Repo: "o/r", Number: 5}, Title: "M2: Two", Tasks: []tracker.IssueRef{{Repo: "o/r", Number: 6}}},
		tasks: map[int]tracker.Task{
			6: {Closed: true},  // this milestone's own task
			7: {Closed: true},  // an earlier milestone's, its branch empty
			8: {Closed: false}, // still open
			9: {Closed: true},  // closed, but the branch carries work
		},
		body:     "## Goal\nx\n\n## Requirements\n- **M2-R1** — a thing\n",
		comments: map[int][]tracker.Comment{5: {{Author: "davison", Body: "**M2-R1 — satisfied.** ran it"}}},
		hasDoc:   true,
		branches: map[string][]string{"o/r": {"task/7-empty", "task/8-open", "task/9-work"}},
		ahead:    map[string]int{"task/7-empty": 0, "task/8-open": 0, "task/9-work": 2},
	}
	cfg := &config.Config{Codecrew: "1.0", Hub: "self", Roles: map[string]config.Role{"implementer": {}, "reviewer": {}, "qa": {}, "doc-synthesizer": {}}}
	c := &ctx{cfg: cfg, roles: cfg, current: "o/r", hub: "o/r", t: f}

	var dry bytes.Buffer
	p, _, err := planClose(c, 2, true, &dry)
	if err != nil || p.refusal != nil {
		t.Fatalf("dry run: err %v refusal %v", err, p.refusal)
	}
	p.print(&dry)
	for _, want := range []string{
		"would delete stale branch task/7-empty (no PR, nothing beyond the default branch)",
		"would keep stale branch task/8-open (o/r#8 is open)",
		"would keep stale branch task/9-work (2 commit(s) not on the default branch, no merged PR)",
		"Swept from earlier closes: task/7-empty.",
	} {
		if !strings.Contains(dry.String(), want) {
			t.Errorf("dry report must contain %q:\n%s", want, dry.String())
		}
	}
	if len(f.writes) != 0 {
		t.Errorf("the dry run wrote: %v", f.writes)
	}

	var live bytes.Buffer
	p, run, err := planClose(c, 2, false, &live)
	if err != nil || p.refusal != nil {
		t.Fatalf("live: err %v refusal %v", err, p.refusal)
	}
	if err := run(&live); err != nil {
		t.Fatalf("run: %v", err)
	}
	for _, want := range []string{
		"swept from earlier closes:\n",
		"stale branch task/7-empty: deleted (no PR, nothing beyond the default branch)",
		"stale branch task/8-open: kept (o/r#8 is open)",
	} {
		if !strings.Contains(live.String(), want) {
			t.Errorf("live output must contain %q:\n%s", want, live.String())
		}
	}
	want := "delete: task/7-empty;close: Closed by `gh codecrew milestone close 2`: all 1 tasks done, milestone document merged. Swept from earlier closes: task/7-empty."
	if got := strings.Join(f.writes, ";"); got != want {
		t.Errorf("writes =\n%s\nwant\n%s", got, want)
	}
}

// GitHub parses the PR body itself, as prose, so a closing keyword near an
// example ref hands the merge an issue nobody meant to close — PR #294
// shipped two of them (#303, M15-R5). task finish says so before it merges,
// and the dry run says the same thing in the same place. It is a note and
// never a gate: the body was parsed long before the verb ran.
func TestFinishNotesClosingReferencesBeyondTheTask(t *testing.T) {
	f := cleanFinish()
	f.closes = []tracker.TitledIssue{
		{Ref: tracker.IssueRef{Repo: "o/r", Number: 7}, Title: "Seven"},
		{Ref: tracker.IssueRef{Repo: "o/r", Number: 42}, Title: "An example ref in the prose"},
	}
	p, run, err := planFinish(finishCtx(f, crewRoles), f.task.Ref, false, false)
	if err != nil || p.refusal != nil {
		t.Fatalf("an unintended reference is a note, not a refusal: err %v refusal %v", err, p.refusal)
	}
	want := "note: this PR would also close o/r#42 (An example ref in the prose) — not the task\n"
	var dry bytes.Buffer
	p.print(&dry)
	if !strings.Contains(dry.String(), want) {
		t.Errorf("the dry run's listing lacks %q:\n%s", want, dry.String())
	}
	if strings.Contains(dry.String(), "o/r#7 (Seven)") {
		t.Errorf("the task's own reference is what the PR is for and must not be noted:\n%s", dry.String())
	}
	if len(f.writes) != 0 {
		t.Errorf("planning wrote: %v", f.writes)
	}
	var live bytes.Buffer
	if err := run(&live); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(live.String(), want) {
		t.Errorf("the live run lacks %q:\n%s", want, live.String())
	}
	// Before the merge, so the operator reads it while it still matters.
	if i, j := strings.Index(live.String(), want), strings.Index(live.String(), "merged PR #9"); i < 0 || j < 0 || i > j {
		t.Errorf("the note must print before the merge line:\n%s", live.String())
	}
}

// A PR that closes only its task says nothing at all — in either mode.
func TestFinishSaysNothingWhenThePRClosesOnlyTheTask(t *testing.T) {
	f := cleanFinish()
	p, run, err := planFinish(finishCtx(f, crewRoles), f.task.Ref, false, false)
	if err != nil || p.refusal != nil {
		t.Fatalf("plan: err %v refusal %v", err, p.refusal)
	}
	var dry bytes.Buffer
	p.print(&dry)
	var live bytes.Buffer
	if err := run(&live); err != nil {
		t.Fatal(err)
	}
	for name, out := range map[string]string{"dry run": dry.String(), "live run": live.String()} {
		if strings.Contains(out, "would also close") {
			t.Errorf("%s: nothing beyond the task is closed, so nothing is said:\n%s", name, out)
		}
	}
}

// The read behind the note is advisory, and so is its failure: every gate
// has passed by the time it runs, and a GraphQL hiccup on an informational
// line must not abort a merge nothing else objects to. It says what it
// could not read and names the command that answers it by hand (checky,
// PR #317).
func TestFinishNotesAnUnreadableClosingReference(t *testing.T) {
	f := cleanFinish()
	f.closesErr = errors.New("502")
	p, run, err := planFinish(finishCtx(f, crewRoles), f.task.Ref, false, false)
	if err != nil || p.refusal != nil || run == nil {
		t.Fatalf("an advisory read must not abort the finish: err %v refusal %v", err, p.refusal)
	}
	want := "note: could not read what else PR #9 would close (502) — check with gh pr view 9 --json closingIssuesReferences\n"
	var dry bytes.Buffer
	p.print(&dry)
	if !strings.Contains(dry.String(), want) {
		t.Errorf("the dry run's listing lacks %q:\n%s", want, dry.String())
	}
	var live bytes.Buffer
	if err := run(&live); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(live.String(), want) {
		t.Errorf("the live run lacks %q:\n%s", want, live.String())
	}
	if !strings.Contains(live.String(), "merged PR #9") {
		t.Errorf("the merge must still happen:\n%s", live.String())
	}
}
