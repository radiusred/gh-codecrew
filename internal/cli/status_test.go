package cli

import (
	"bytes"
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
	keepBranches bool // the repo does NOT delete branches on merge
}

func (f *statusFake) OpenMilestones(string) ([]tracker.Milestone, error) { return f.milestones, nil }
func (f *statusFake) Task(ref tracker.IssueRef) (tracker.Task, error) {
	return f.issues[ref.Number], nil
}
func (f *statusFake) IssueBody(tracker.IssueRef) (string, error) {
	return "## Requirements\n- **M2-R1** — a thing\n", nil
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
// check reads the hub's roles/ against the contracts embedded in the
// binary and never touches milestone state, and the quiet period is when a
// fork gets reconciled against a new release — so the no-open-milestones
// line replaces the board, not the two advisory checks under it (#253).
func TestStatusWithoutOpenMilestonesStillReportsDriftAndSetting(t *testing.T) {
	c := statusCtx(t, &statusFake{keepBranches: true})
	if err := os.MkdirAll(filepath.Join(c.cfg.Dir, "roles"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(c.cfg.Dir, "roles", "coordinator.md"), []byte("a fork of the contract\n"), 0o644); err != nil {
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
		"contract drift: roles/coordinator.md differs from the embedded",
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
