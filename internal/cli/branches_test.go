package cli

import (
	"bytes"
	"errors"
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
	linked  map[int][]string
	prs     map[int][]int
	info    map[int]tracker.PR
	ahead   map[string]int
	tips    map[string]string
	titles  map[int]string
	deleted []string
	failDel string
	repoErr error
}

func (f *fakeTracker) RepoInfo(string) (tracker.RepoInfo, error) {
	return tracker.RepoInfo{DefaultBranch: "main"}, f.repoErr
}
func (f *fakeTracker) Task(ref tracker.IssueRef) (tracker.Task, error) {
	return tracker.Task{Title: f.titles[ref.Number]}, nil
}
func (f *fakeTracker) LinkedBranches(ref tracker.IssueRef) ([]string, error) {
	return f.linked[ref.Number], nil
}
func (f *fakeTracker) ClosingPRs(ref tracker.IssueRef, _ bool) ([]int, error) {
	return f.prs[ref.Number], nil
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
		prs:    map[int][]int{1: {11}, 4: {14, 15}, 5: {16}, 7: {17}, 8: {18}, 11: {19}},
		info: map[int]tracker.PR{
			11: {HeadRef: "task/1-merged", Merged: true, HeadSHA: "m1"},
			14: {HeadRef: "task/4-open", Merged: true, HeadSHA: "m4"}, // an earlier merged PR…
			15: {HeadRef: "task/4-open", Open: true},                  // …and a live one on the same head
			16: {HeadRef: "task/5-closed-pr"},
			17: {HeadRef: "task/7-forbidden", Merged: true, HeadSHA: "m7"},
			18: {HeadRef: "task/8-moved", Merged: true, HeadSHA: "m8"},
			19: {HeadRef: "task/11-fork", Merged: true, CrossRepo: true},
		},
		ahead: map[string]int{
			"task/1-merged": 4, "task/2-release-v0-3-0-of-the-gh-extension": 0, "task/3-unmerged-work": 2,
			"task/4-open": 1, "task/5-closed-pr": 0, "task/7-forbidden": 3, "task/8-moved": 9,
			"task/9-linked-only": 0, "main": 0, "task/10-main": 0, "task/11-fork": 2,
		},
		tips:    map[string]string{"task/1-merged": "m1", "task/7-forbidden": "m7", "task/8-moved": "later"},
		failDel: "task/7-forbidden",
	}
	m := &tracker.Milestone{}
	for n := 1; n <= 11; n++ {
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
