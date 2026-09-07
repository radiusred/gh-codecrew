package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// statusFake is the slice of the tracker status reads; issues holds the
// milestone issue and the tasks alike, keyed by number.
type statusFake struct {
	tracker.Tracker
	milestones   []tracker.Milestone
	issues       map[int]tracker.Task
	comments     map[int][]tracker.Comment // what Comments answers per task
	commentsErr  error                     // when set, what every Comments call fails with
	read         []int                     // the task numbers Comments was called for
	body         string                    // the milestone body; a well-formed one by default
	keepBranches bool                      // the repo does NOT delete branches on merge

	// The branch side, read only by the stale-branch report: the repo's
	// task/<n>-… listing and the lookups staleBranchAction makes off it.
	branches  []string
	truncated bool
	branchErr error
	taskErr   map[int]error
	prs       map[int][]int
	prInfo    map[int]tracker.PR
	ahead     map[string]int
	tips      map[string]string
	openPRs   map[string][]int
}

func (f *statusFake) OpenMilestones(string) ([]tracker.Milestone, error) { return f.milestones, nil }
func (f *statusFake) Task(ref tracker.IssueRef) (tracker.Task, error) {
	if err := f.taskErr[ref.Number]; err != nil {
		return tracker.Task{}, err
	}
	return f.issues[ref.Number], nil
}
func (f *statusFake) TaskBranches(string) ([]string, bool, error) {
	return f.branches, f.truncated, f.branchErr
}
func (f *statusFake) ClosingPRs(ref tracker.IssueRef, _ bool) ([]int, error) {
	return f.prs[ref.Number], nil
}
func (f *statusFake) PRInfo(_ string, n int) (tracker.PR, error) { return f.prInfo[n], nil }
func (f *statusFake) OpenPRsForBranch(_, branch string) ([]int, error) {
	return f.openPRs[branch], nil
}
func (f *statusFake) BranchAhead(_, b string) (int, string, error) {
	n, ok := f.ahead[b]
	if !ok {
		return 0, "", errors.New("branch not found")
	}
	return n, f.tips[b], nil
}
func (f *statusFake) IssueBody(tracker.IssueRef) (string, error) {
	if f.body != "" {
		return f.body, nil
	}
	return "## Requirements\n- **M2-R1** — a thing\n", nil
}
func (f *statusFake) Comments(ref tracker.IssueRef) ([]tracker.Comment, error) {
	f.read = append(f.read, ref.Number)
	if f.commentsErr != nil {
		return nil, f.commentsErr
	}
	return f.comments[ref.Number], nil
}
func (f *statusFake) RepoInfo(string) (tracker.RepoInfo, error) {
	return tracker.RepoInfo{DefaultBranch: "main", DeleteBranchOnMerge: !f.keepBranches}, nil
}

func statusCtx(t *testing.T, f tracker.Tracker) *ctx {
	t.Helper()
	cfg := &config.Config{Codecrew: "1.0", Hub: "self", Dir: t.TempDir()}
	return &ctx{cfg: cfg, roles: cfg, current: "o/r", hub: "o/r", t: f}
}

// A gate raised on the milestone issue itself — a requirement-level
// question no task carries — shows on the board beside the task gates,
// marked so the reader can tell the two apart (#200).
func TestStatusListsMilestoneGate(t *testing.T) {
	f := &statusFake{
		milestones: []tracker.Milestone{{
			Ref: tracker.IssueRef{Repo: "o/r", Number: 5}, Title: "M2: Two",
			Tasks: []tracker.IssueRef{{Repo: "o/r", Number: 6}, {Repo: "o/r", Number: 7}},
		}},
		issues: map[int]tracker.Task{
			5: {Title: "M2: Two", Labels: []string{tracker.LabelMilestone, tracker.LabelNeedsDecision}},
			6: {Title: "Six", Labels: []string{"cc:task", tracker.LabelNeedsDecision}},
			7: {Title: "Seven", Labels: []string{"cc:task"}},
		},
	}
	var out bytes.Buffer
	if err := statusReport(&out, statusCtx(t, f)); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{
		"M2: Two (o/r#5) — gate raised on the milestone issue\n",
		"gates raised:\n  o/r#5 — M2: Two (milestone)\n  o/r#6 — Six\n",
		"[gated      ] o/r#6",
		"[ready      ] o/r#7",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("status output lacks %q:\n%s", want, got)
		}
	}
	if strings.Count(got, "(milestone)") != 1 {
		t.Errorf("only the milestone's gate is marked (milestone):\n%s", got)
	}
}

// Without the label on the milestone issue the header is bare and the
// list is what the tasks contribute — none here.
func TestStatusMilestoneWithoutGate(t *testing.T) {
	f := &statusFake{
		milestones: []tracker.Milestone{{
			Ref: tracker.IssueRef{Repo: "o/r", Number: 5}, Title: "M2: Two",
			Tasks: []tracker.IssueRef{{Repo: "o/r", Number: 7}},
		}},
		issues: map[int]tracker.Task{
			5: {Title: "M2: Two", Labels: []string{tracker.LabelMilestone}},
			7: {Title: "Seven", Labels: []string{"cc:task"}},
		},
	}
	var out bytes.Buffer
	if err := statusReport(&out, statusCtx(t, f)); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	if !strings.Contains(got, "M2: Two (o/r#5)\n") || strings.Contains(got, "gate raised on the milestone issue") {
		t.Errorf("header must be bare without the label:\n%s", got)
	}
	if !strings.Contains(got, "gates raised: none\n") {
		t.Errorf("no gate anywhere must print none:\n%s", got)
	}
}

// Between milestones status still reports what it knows: the contract-drift
// check reads the hub's .codecrew/roles/ against the contracts embedded in
// the binary and never touches milestone state, and the quiet period is
// when a
// fork gets reconciled against a new release — so the no-open-milestones
// line replaces the board, not the two advisory checks under it (#253).
func TestStatusWithoutOpenMilestonesStillReportsDriftAndSetting(t *testing.T) {
	c := statusCtx(t, &statusFake{keepBranches: true})
	if err := os.MkdirAll(filepath.Join(c.cfg.Dir, filepath.FromSlash(config.RolesDir)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(c.cfg.Dir, rolesPath("coordinator.md")), []byte("a fork of the contract\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := statusReport(&out, c); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{
		"no open milestones in o/r\n",
		"note: o/r does not delete branches on merge",
		"contract drift: " + contractPath("coordinator") + " differs from the embedded",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("status output lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "gates raised") {
		t.Errorf("with no open milestone there is no gates section to print:\n%s", got)
	}
}

// The counterpart: an undrifted hub whose repo deletes branches prints the
// line and nothing else, so neither advisory check is free.
func TestStatusWithoutOpenMilestonesSaysNothingElseWhenClean(t *testing.T) {
	var out bytes.Buffer
	if err := statusReport(&out, statusCtx(t, &statusFake{})); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); got != "no open milestones in o/r\n" {
		t.Errorf("a clean hub between milestones prints one line, got:\n%s", got)
	}
}

// status reports the board; it does not gate it. A requirement ID that is
// not the milestone's own — the condition milestone close and milestone
// evidence refuse REQUIREMENT_ID_MISMATCH on — is printed as a line naming
// the code, and the rest of the report still runs (M13-R6).
func TestStatusReportsRequirementIDMismatch(t *testing.T) {
	f := &statusFake{
		milestones: []tracker.Milestone{{
			Ref: tracker.IssueRef{Repo: "o/r", Number: 5}, Title: "M2: Two",
			Tasks: []tracker.IssueRef{{Repo: "o/r", Number: 7}},
		}},
		issues: map[int]tracker.Task{
			5: {Title: "M2: Two", Labels: []string{tracker.LabelMilestone}},
			7: {Title: "Seven", Labels: []string{"cc:task"}},
		},
		body: "## Requirements\n- **M2-R1** — mine\n- **M9-R4** — another milestone's\n",
	}
	var out bytes.Buffer
	if err := statusReport(&out, statusCtx(t, f)); err != nil {
		t.Fatalf("status must report the mismatch, not die on it: %v", err)
	}
	got := out.String()
	for _, want := range []string{
		"REQUIREMENT_ID_MISMATCH",
		"M9-R4",
		"not M2's",
		"[ready      ] o/r#7",  // the report carries on past the note…
		"gates raised: none\n", // …to the end
	} {
		if !strings.Contains(got, want) {
			t.Errorf("status output lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "M2-R1") {
		t.Errorf("the milestone's own ID must not be named as a mismatch:\n%s", got)
	}
}

// A well-formed Requirements section prints no mismatch line at all.
func TestStatusSilentWhenRequirementIDsMatch(t *testing.T) {
	f := &statusFake{
		milestones: []tracker.Milestone{{
			Ref: tracker.IssueRef{Repo: "o/r", Number: 5}, Title: "M2: Two",
			Tasks: []tracker.IssueRef{{Repo: "o/r", Number: 7}},
		}},
		issues: map[int]tracker.Task{
			5: {Title: "M2: Two", Labels: []string{tracker.LabelMilestone}},
			7: {Title: "Seven", Labels: []string{"cc:task"}},
		},
	}
	var out bytes.Buffer
	if err := statusReport(&out, statusCtx(t, f)); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); strings.Contains(got, "REQUIREMENT_ID_MISMATCH") {
		t.Errorf("a milestone's own IDs must print no note:\n%s", got)
	}
}

// startedBy is the record task start posts, as the task's own seat.
func startedBy(login string) []tracker.Comment {
	return []tracker.Comment{{Author: login, Body: tracker.StartRecord(login)}}
}

// #287: an App-held seat has no assignee to show — GitHub does not accept
// an App as an issue assignee — so the board read the holder off an empty
// list and named nobody. It reads the **Started by** record instead, which
// is the only thing that says a task was started, and names the App.
func TestStatusNamesTheHolderFromTheStartRecord(t *testing.T) {
	f := &statusFake{
		milestones: []tracker.Milestone{{
			Ref: tracker.IssueRef{Repo: "o/r", Number: 5}, Title: "M2: Two",
			Tasks: []tracker.IssueRef{{Repo: "o/r", Number: 6}, {Repo: "o/r", Number: 7}},
		}},
		issues: map[int]tracker.Task{
			5: {Title: "M2: Two", Labels: []string{tracker.LabelMilestone}},
			// In review with no assignee at all: the App-run shape.
			6: {Title: "Six", Labels: []string{tracker.LabelTask}, OpenLinkedPR: true},
			// In progress, and the assignee is not who started it: a
			// handover, where the latest record is the current holder.
			7: {Title: "Seven", Labels: []string{tracker.LabelTask}, Assignees: []string{"davison"}},
		},
		comments: map[int][]tracker.Comment{
			6: startedBy("radiusred-cody[bot]"),
			7: startedBy("radiusred-checky[bot]"),
		},
	}
	var out bytes.Buffer
	if err := statusReport(&out, statusCtx(t, f)); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{
		"[in review  ] o/r#6", "Six @radiusred-cody[bot]\n",
		"[in progress] o/r#7", "Seven @radiusred-checky[bot]\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("status output lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "@davison") {
		t.Errorf("the assignee outranked the start record:\n%s", got)
	}
}

// The fallback is display-only and narrow: a task nothing records a start
// on is in progress precisely because it carries an assignee, so the state
// would otherwise sit there with no name to explain it. A task in neither
// state has no holder to name and costs no comments read.
func TestStatusFallsBackToTheAssigneeAndReadsOnlyTasksInFlight(t *testing.T) {
	f := &statusFake{
		milestones: []tracker.Milestone{{
			Ref: tracker.IssueRef{Repo: "o/r", Number: 5}, Title: "M2: Two",
			Tasks: []tracker.IssueRef{{Repo: "o/r", Number: 6}, {Repo: "o/r", Number: 7}},
		}},
		issues: map[int]tracker.Task{
			5: {Title: "M2: Two", Labels: []string{tracker.LabelMilestone}},
			6: {Title: "Six", Labels: []string{tracker.LabelTask}, Assignees: []string{"davison"}},
			7: {Title: "Seven", Labels: []string{tracker.LabelTask}},
		},
	}
	var out bytes.Buffer
	if err := statusReport(&out, statusCtx(t, f)); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	if !strings.Contains(got, "[in progress] o/r#6") || !strings.Contains(got, "Six @davison\n") {
		t.Errorf("no start record must fall back to the assignee:\n%s", got)
	}
	if !strings.Contains(got, "[ready      ] o/r#7") || !strings.Contains(got, "Seven\n") {
		t.Errorf("a ready task names no holder:\n%s", got)
	}
	if len(f.read) != 1 || f.read[0] != 6 {
		t.Errorf("comments were read for %v, want only the task in flight [6]", f.read)
	}
}

// The board is a report, not a gate: a comments read this seat cannot make
// falls through to the assignee and the rest of the board still prints.
func TestStatusSurvivesAnUnreadableCommentsList(t *testing.T) {
	f := &statusFake{
		milestones: []tracker.Milestone{{
			Ref: tracker.IssueRef{Repo: "o/r", Number: 5}, Title: "M2: Two",
			Tasks: []tracker.IssueRef{{Repo: "o/r", Number: 6}},
		}},
		issues: map[int]tracker.Task{
			5: {Title: "M2: Two", Labels: []string{tracker.LabelMilestone}},
			6: {Title: "Six", Labels: []string{tracker.LabelTask}, Assignees: []string{"davison"}},
		},
		commentsErr: errors.New("gh: Forbidden (HTTP 403)"),
	}
	var out bytes.Buffer
	if err := statusReport(&out, statusCtx(t, f)); err != nil {
		t.Fatalf("an unreadable comments list must not fail the board: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "Six @davison\n") || !strings.Contains(got, "gates raised: none\n") {
		t.Errorf("the board must carry on past the failed read:\n%s", got)
	}
}

// staleFake is a hub between milestones carrying four task branches: one
// whose task closed with a merged PR still at the tip, one whose task
// closed with unmerged work on it, one whose task is still open, and one
// whose task cannot be read at all.
func staleFake() *statusFake {
	return &statusFake{
		branches: []string{"task/7-merged", "task/8-unmerged", "task/9-open", "task/10-unreadable"},
		issues: map[int]tracker.Task{
			7: {Title: "Seven", Closed: true},
			8: {Title: "Eight", Closed: true},
			9: {Title: "Nine"},
		},
		taskErr: map[int]error{10: errors.New("404")},
		prs:     map[int][]int{7: {70}},
		prInfo:  map[int]tracker.PR{70: {HeadRef: "task/7-merged", Merged: true, HeadSHA: "abc"}},
		ahead:   map[string]int{"task/7-merged": 3, "task/8-unmerged": 2, "task/9-open": 1, "task/10-unreadable": 1},
		tips:    map[string]string{"task/7-merged": "abc"},
	}
}

// The report names a closed task's branch with the verdict and reason
// branchAction gives it — the same function milestone close's second sweep
// judges by, so the two can never disagree (#295, M15-R5) — says nothing
// at all about a branch whose task is still open, and repeats the sweep's
// own note: for one whose task it could not read.
func TestStatusReportsStaleBranchesWithTheSweepsVerdict(t *testing.T) {
	f := staleFake()
	var out bytes.Buffer
	if err := statusReport(&out, statusCtx(t, f)); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	// The reasons are not written out here: they are asked of the very
	// function the sweep asks, so a reworded verdict moves both at once.
	_, merged := branchAction(f.prInfo[70], true, 3, "abc")
	_, unmerged := branchAction(tracker.PR{}, false, 2, "")
	for _, want := range []string{
		"stale branch: task/7-merged — o/r#7 is closed; " + merged + " — the next milestone close would delete it\n",
		"stale branch: task/8-unmerged — o/r#8 is closed; " + unmerged + " — kept by the next milestone close\n",
		"note: stale branch task/10-unreadable skipped (404)\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("status output lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "task/9-open") {
		t.Errorf("a branch whose task is still open is not stale and must not be reported:\n%s", got)
	}
	if strings.Contains(got, "carries more task branches") {
		t.Errorf("a complete listing must not be called partial:\n%s", got)
	}
}

// The reason a report costs what it does: the verdict is only computed for
// a branch whose task has closed. An open task's branch costs the one issue
// read and stops there — no PR listing, no comparison.
func TestStatusStaleReportStopsAtAnOpenTask(t *testing.T) {
	f := staleFake()
	f.branches = []string{"task/9-open"}
	f.ahead = nil // a comparison would fail outright if one were made
	var out bytes.Buffer
	if err := statusReport(&out, statusCtx(t, f)); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); strings.Contains(got, "stale branch") {
		t.Errorf("nothing stale must print nothing:\n%s", got)
	}
}

// A listing GitHub could not give whole is reported as partial rather than
// left to read as the repo's complete state.
func TestStatusSaysWhenTheBranchListingWasPartial(t *testing.T) {
	f := staleFake()
	f.truncated = true
	var out bytes.Buffer
	if err := statusReport(&out, statusCtx(t, f)); err != nil {
		t.Fatal(err)
	}
	want := "note: o/r carries more task branches than one listing holds; this report saw 4 of them\n"
	if got := out.String(); !strings.Contains(got, want) {
		t.Errorf("status output lacks %q:\n%s", want, got)
	}
}

// Every check under the board is advisory: a listing that cannot be read is
// one note: and the rest of the report still prints.
func TestStatusStaleReportIsAdvisory(t *testing.T) {
	f := staleFake()
	f.branchErr = errors.New("403")
	f.keepBranches = true
	var out bytes.Buffer
	if err := statusReport(&out, statusCtx(t, f)); err != nil {
		t.Fatalf("an unreadable listing must not fail the verb: %v", err)
	}
	got := out.String()
	for _, want := range []string{
		"note: stale task branches not listed for o/r (403)\n",
		"note: o/r does not delete branches on merge",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("status output lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "stale branch:") {
		t.Errorf("no branch was read, so none may be reported:\n%s", got)
	}
}
