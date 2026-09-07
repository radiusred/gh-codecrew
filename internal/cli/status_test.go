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
}

func (f *statusFake) OpenMilestones(string) ([]tracker.Milestone, error) { return f.milestones, nil }
func (f *statusFake) Task(ref tracker.IssueRef) (tracker.Task, error) {
	return f.issues[ref.Number], nil
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
