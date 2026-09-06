package cli

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

func TestBranchAction(t *testing.T) {
	merged := tracker.PR{Merged: true, HeadSHA: "abc"}
	open := tracker.PR{Open: true}
	closed := tracker.PR{}
	cases := []struct {
		name  string
		pr    tracker.PR
		hasPR bool
		ahead int
		tip   string
		del   bool
	}{
		{"merged, tip untouched", merged, true, 5, "abc", true}, // ancestry never mattered
		{"merged, commits after", merged, true, 7, "def", false},
		{"merged, no SHA known", tracker.PR{Merged: true}, true, 5, "", false}, // empty == empty must not read as untouched
		{"open PR", open, true, 0, "x", false},
		{"closed unmerged, empty", closed, true, 0, "x", true},
		{"closed unmerged, work", closed, true, 1, "x", false},
		{"no PR, empty", tracker.PR{}, false, 0, "x", true},
		{"no PR, work", tracker.PR{}, false, 3, "x", false},
	}
	for _, c := range cases {
		if del, _ := branchAction(c.pr, c.hasPR, c.ahead, c.tip); del != c.del {
			t.Errorf("%s: delete = %v, want %v", c.name, del, c.del)
		}
	}
}

func TestPRByHeadOpenWins(t *testing.T) {
	open := tracker.PR{HeadRef: "task/1-x", Open: true}
	merged := tracker.PR{HeadRef: "task/1-x", Merged: true, HeadSHA: "abc"}
	for _, order := range [][]tracker.PR{{open, merged}, {merged, open}} {
		if got := prByHead(order)["task/1-x"]; !got.Open {
			t.Errorf("order %v: open PR must win over merged", order)
		}
	}
	fork := tracker.PR{HeadRef: "task/2-y", Merged: true, CrossRepo: true}
	if _, ok := prByHead([]tracker.PR{fork})["task/2-y"]; ok {
		t.Error("a fork's head must not become a candidate in the base repo")
	}
}

// fakeTracker embeds the interface so only the methods a test exercises
// need defining; any other call panics, which is the point.
type fakeTracker struct {
	tracker.Tracker
	linked    map[int][]string
	prs       map[int][]int
	info      map[int]tracker.PR
	ahead     map[string]int
	tips      map[string]string
	titles    map[int]string
	closed    map[int]bool
	branches  map[string][]string
	calls     []string // every issue and PR lookup, in order, so a bound can be asserted
	deleted   []string
	failDel   string
	repoErr   error
	branchErr error
}

func (f *fakeTracker) RepoInfo(string) (tracker.RepoInfo, error) {
	return tracker.RepoInfo{DefaultBranch: "main"}, f.repoErr
}
func (f *fakeTracker) Task(ref tracker.IssueRef) (tracker.Task, error) {
	f.calls = append(f.calls, fmt.Sprintf("task %d", ref.Number))
	return tracker.Task{Title: f.titles[ref.Number], Closed: f.closed[ref.Number]}, nil
}
func (f *fakeTracker) LinkedBranches(ref tracker.IssueRef) ([]string, error) {
	return f.linked[ref.Number], nil
}
func (f *fakeTracker) ClosingPRs(ref tracker.IssueRef, _ bool) ([]int, error) {
	f.calls = append(f.calls, fmt.Sprintf("prs %d", ref.Number))
	return f.prs[ref.Number], nil
}
func (f *fakeTracker) TaskBranches(repo string) ([]string, error) {
	return f.branches[repo], f.branchErr
}
func (f *fakeTracker) PRInfo(_ string, n int) (tracker.PR, error) { return f.info[n], nil }
func (f *fakeTracker) BranchAhead(_, b string) (int, string, error) {
	n, ok := f.ahead[b]
	if !ok {
		return 0, "", errors.New("branch not found")
	}
	return n, f.tips[b], nil
}
func (f *fakeTracker) DeleteBranch(_, b string) error {
	if b == f.failDel {
		return errors.New("403")
	}
	f.deleted = append(f.deleted, b)
	return nil
}

func TestSweepBranches(t *testing.T) {
	ft := &fakeTracker{
		// The relation only ever contributes for PR-less tasks: #1's merged
		// branch is known solely from its PR's head, #2's only from the
		// convention, #9's only from the relation.
		linked: map[int][]string{6: {"task/6-gone"}, 9: {"task/9-linked-only"}},
		titles: map[int]string{2: "Release v0.3.0 of the gh extension", 3: "Unmerged work", 10: "Main"},
		prs:    map[int][]int{1: {11}, 4: {14, 15}, 5: {16}, 7: {17}, 8: {18}, 11: {19}, 12: {20}},
		info: map[int]tracker.PR{
			11: {HeadRef: "task/1-merged", Merged: true, HeadSHA: "m1"},
			14: {HeadRef: "task/4-open", Merged: true, HeadSHA: "m4"}, // an earlier merged PR…
			15: {HeadRef: "task/4-open", Open: true},                  // …and a live one on the same head
			16: {HeadRef: "task/5-closed-pr"},
			17: {HeadRef: "task/7-forbidden", Merged: true, HeadSHA: "m7"},
			18: {HeadRef: "task/8-moved", Merged: true, HeadSHA: "m8"},
			19: {HeadRef: "task/11-fork", Merged: true, CrossRepo: true},
			20: {HeadRef: "main", Merged: true, HeadSHA: "mm"}, // a PR from the default branch itself
		},
		ahead: map[string]int{
			"task/1-merged": 4, "task/2-release-v0-3-0-of-the-gh-extension": 0, "task/3-unmerged-work": 2,
			"task/4-open": 1, "task/5-closed-pr": 0, "task/7-forbidden": 3, "task/8-moved": 9,
			"task/9-linked-only": 0, "main": 0, "task/10-main": 0, "task/11-fork": 2,
		},
		tips:    map[string]string{"task/1-merged": "m1", "task/7-forbidden": "m7", "task/8-moved": "later", "main": "mm"},
		failDel: "task/7-forbidden",
	}
	m := &tracker.Milestone{}
	for n := 1; n <= 12; n++ {
		m.Tasks = append(m.Tasks, tracker.IssueRef{Repo: "o/r", Number: n})
	}
	var out bytes.Buffer
	deleted := sweepBranches(&out, ft, m)
	want := "task/1-merged,task/2-release-v0-3-0-of-the-gh-extension,task/5-closed-pr,task/9-linked-only,task/10-main"
	if strings.Join(deleted, ",") != want {
		t.Errorf("deleted = %v\nwant %s\n%s", deleted, want, out.String())
	}
	for _, line := range []string{
		"task/1-merged: deleted (PR merged)",
		"task/2-release-v0-3-0-of-the-gh-extension: deleted (no PR, nothing beyond the default branch)",
		"task/3-unmerged-work: kept (2 commit(s) not on the default branch, no merged PR)",
		"task/4-open: kept (open PR)", // open wins over the earlier merged PR on the same head
		"task/5-closed-pr: deleted (PR closed unmerged, nothing beyond the default branch)",
		"task/7-forbidden: kept (PR merged; delete failed",
		"task/8-moved: kept (commits pushed after the PR merged)",
		"task/9-linked-only: deleted (no PR, nothing beyond the default branch)",
	} {
		if !strings.Contains(out.String(), line) {
			t.Errorf("output missing %q:\n%s", line, out.String())
		}
	}
	for _, never := range []string{"task/6-gone", "branch main:", "task/11-fork"} {
		if strings.Contains(out.String(), never) {
			t.Errorf("%q must not appear (gone / default branch / fork):\n%s", never, out.String())
		}
	}
	for _, d := range ft.deleted {
		if d == "main" || d == "task/11-fork" {
			t.Errorf("deleted %s", d)
		}
	}
}

func TestSweepBranchesFailuresAreNotes(t *testing.T) {
	ft := &fakeTracker{repoErr: errors.New("api down")}
	var out bytes.Buffer
	if got := sweepBranches(&out, ft, &tracker.Milestone{Tasks: []tracker.IssueRef{{Repo: "o/r", Number: 1}}}); len(got) != 0 || !strings.HasPrefix(out.String(), "note: branch sweep skipped") {
		t.Errorf("failure not a note: %v %q", got, out.String())
	}
}

func TestDeleteHeadIsANoteOnFailure(t *testing.T) {
	ft := &fakeTracker{failDel: "task/9-x"}
	var out bytes.Buffer
	deleteHead(&out, ft, tracker.PR{Repo: "o/r", HeadRef: "task/9-x"})
	if !strings.HasPrefix(out.String(), "note: could not delete branch task/9-x") {
		t.Errorf("failure not a note: %q", out.String())
	}
	out.Reset()
	deleteHead(&out, ft, tracker.PR{Repo: "o/r", HeadRef: "task/8-ok"})
	if out.String() != "deleted branch task/8-ok\n" || len(ft.deleted) != 1 {
		t.Errorf("success: %q, deleted %v", out.String(), ft.deleted)
	}
	out.Reset()
	deleteHead(&out, ft, tracker.PR{Repo: "o/r", HeadRef: "task/7-fork", CrossRepo: true}) // a fork's branch: silent
	deleteHead(&out, ft, tracker.PR{Repo: "o/r"})                                          // no head known: silent
	if out.String() != "" || len(ft.deleted) != 1 {
		t.Errorf("fork/empty head produced output or deletion: %q %v", out.String(), ft.deleted)
	}
}

func TestTaskNumber(t *testing.T) {
	cases := map[string]int{
		"task/273-milestone-close-sweeps-stale": 273,
		"task/8-cycle-1-milestone-record-m1-r6": 8,
		"task/9":                                9, // a title that slugged to nothing
		"main":                                  0,
		"codecrew-init":                         0,
		"task/":                                 0,
		"task/-x":                               0,
		"task/x-1":                              0,
		"task/0-zero":                           0,
		"task/007-padded":                       0, // task start never writes one
		"feature/task/5-not-ours":               0,
	}
	for branch, want := range cases {
		if got := taskNumber(branch); got != want {
			t.Errorf("taskNumber(%q) = %d, want %d", branch, got, want)
		}
	}
}

// The second pass reaches the branches earlier closes left behind (#167):
// a closed task's merged or empty branch goes, a closed task's unmerged
// work and an open task's branch are named and left, and the closing
// milestone's own branches belong to the first pass, not this one.
func TestPlanStaleSweep(t *testing.T) {
	ft := &fakeTracker{
		branches: map[string][]string{
			"o/hub": {
				"task/1-own",       // the closing milestone's own task: pass one's
				"task/20-merged",   // closed task, PR merged, tip untouched: swept
				"task/21-empty",    // closed task, no PR, nothing beyond main: swept
				"task/22-unmerged", // closed task, work not on main: kept and named
				"task/23-open",     // task still open: kept and named
				"task/24-gone",     // vanished between the listing and the check
				"task/0-nonsense",  // no number task start would have written
				"main",
			},
			"o/spoke": {"task/30-merged"},
		},
		closed: map[int]bool{20: true, 21: true, 22: true, 24: true, 30: true},
		prs:    map[int][]int{20: {200}, 22: {220}, 30: {300}},
		info: map[int]tracker.PR{
			200: {HeadRef: "task/20-merged", Merged: true, HeadSHA: "m20"},
			220: {HeadRef: "task/22-unmerged"},
			300: {HeadRef: "task/30-merged", Merged: true, HeadSHA: "m30"},
		},
		ahead: map[string]int{"task/20-merged": 3, "task/21-empty": 0, "task/22-unmerged": 2, "task/23-open": 0, "task/30-merged": 4},
		tips:  map[string]string{"task/20-merged": "m20", "task/30-merged": "m30"},
	}
	m := &tracker.Milestone{
		Ref:   tracker.IssueRef{Repo: "o/hub", Number: 9},
		Tasks: []tracker.IssueRef{{Repo: "o/hub", Number: 1}, {Repo: "o/spoke", Number: 2}},
	}
	own := []sweepItem{{Repo: "o/hub", Name: "task/1-own", Delete: true, Reason: "PR merged"}}
	items := planStaleSweep(ft, m, own)
	if len(ft.deleted) != 0 {
		t.Fatalf("planning deleted %v", ft.deleted)
	}
	for _, it := range items {
		if it.Note == "" && !it.Stale {
			t.Errorf("%s not marked stale: %+v", it.Name, it)
		}
	}
	var out bytes.Buffer
	deleted := executeSweep(&out, ft, items)
	if want := "task/20-merged,task/21-empty,task/30-merged"; strings.Join(deleted, ",") != want {
		t.Errorf("deleted = %v, want %s\n%s", deleted, want, out.String())
	}
	for _, line := range []string{
		"stale branch task/20-merged: deleted (PR merged)",
		"stale branch task/21-empty: deleted (no PR, nothing beyond the default branch)",
		"stale branch task/22-unmerged: kept (2 commit(s) not on the default branch, no merged PR)",
		"stale branch task/23-open: kept (o/hub#23 is open)",
		"stale branch task/30-merged: deleted (PR merged)",
	} {
		if !strings.Contains(out.String(), line) {
			t.Errorf("output missing %q:\n%s", line, out.String())
		}
	}
	// The closing milestone's own branches, the default branch, a name with
	// no task number and a branch that went away are none of this pass's.
	for _, never := range []string{"task/1-own", "main", "task/0-nonsense", "task/24-gone"} {
		if strings.Contains(out.String(), never) {
			t.Errorf("%q must not appear:\n%s", never, out.String())
		}
	}
	// Bounded: an open task's branch costs one issue read and stops there;
	// the milestone's own tasks are never re-read by this pass at all.
	for _, never := range []string{"prs 23", "task 1", "task 2"} {
		if slices.Contains(ft.calls, never) {
			t.Errorf("lookup %q is beyond the bound: %v", never, ft.calls)
		}
	}
	if n := slices.Index(ft.calls, "task 23"); n < 0 || slices.Index(ft.calls[n+1:], "task 23") >= 0 {
		t.Errorf("open task not read exactly once: %v", ft.calls)
	}
}

func TestPlanStaleSweepFailuresAreNotes(t *testing.T) {
	m := &tracker.Milestone{Ref: tracker.IssueRef{Repo: "o/hub", Number: 9}}
	for _, ft := range []*fakeTracker{
		{repoErr: errors.New("api down")},
		{branchErr: errors.New("api down")},
	} {
		items := planStaleSweep(ft, m, nil)
		var out bytes.Buffer
		if got := executeSweep(&out, ft, items); len(got) != 0 || !strings.HasPrefix(out.String(), "note: stale branch sweep skipped for o/hub") {
			t.Errorf("failure not a note: %v %q", got, out.String())
		}
	}
}
