package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

func TestBranchAction(t *testing.T) {
	cases := []struct {
		merged, open bool
		ahead        int
		del          bool
	}{
		{true, false, 5, true},   // merged: content is on main, whatever the rebase did
		{true, true, 0, true},    // merged wins over a stale "open" (cannot happen live; the rule is total)
		{false, true, 0, false},  // open PR: never
		{false, false, 0, true},  // no PR, empty: the release-task case
		{false, false, 3, false}, // unmerged work: keep
	}
	for _, c := range cases {
		if del, _ := branchAction(c.merged, c.open, c.ahead); del != c.del {
			t.Errorf("branchAction(%v, %v, %d) = %v, want %v", c.merged, c.open, c.ahead, del, c.del)
		}
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
	titles  map[int]string
	deleted []string
	failDel string
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
func (f *fakeTracker) BranchAhead(_, b string) (int, error) {
	n, ok := f.ahead[b]
	if !ok {
		return 0, errors.New("branch not found")
	}
	return n, nil
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
		// The relation survives only for open tasks: #1's merged branch is
		// known solely from its PR's head, #2's only from the convention.
		linked: map[int][]string{4: {"task/4-open"}, 6: {"task/6-gone"}},
		titles: map[int]string{2: "Release v0.3.0 of the gh extension", 3: "Unmerged work"},
		prs:    map[int][]int{1: {11}, 4: {14}, 5: {15}, 7: {17}},
		info: map[int]tracker.PR{
			11: {HeadRef: "task/1-merged", Merged: true},
			14: {HeadRef: "task/4-open", Open: true},
			15: {HeadRef: "task/5-closed-pr"},
			17: {HeadRef: "task/7-forbidden", Merged: true},
		},
		ahead: map[string]int{
			"task/1-merged": 4, "task/2-release-v0-3-0-of-the-gh-extension": 0,
			"task/3-unmerged-work": 2, "task/4-open": 1, "task/5-closed-pr": 1, "task/7-forbidden": 3,
		},
		failDel: "task/7-forbidden",
	}
	m := &tracker.Milestone{}
	for n := 1; n <= 7; n++ {
		m.Tasks = append(m.Tasks, tracker.IssueRef{Repo: "o/r", Number: n})
	}
	var out bytes.Buffer
	deleted, err := sweepBranches(&out, ft, m)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(deleted, ",") != "task/1-merged,task/2-release-v0-3-0-of-the-gh-extension" {
		t.Errorf("deleted = %v", deleted)
	}
	for _, want := range []string{
		"task/1-merged: deleted (PR merged)", // 4 ahead — ancestry never mattered
		"task/2-release-v0-3-0-of-the-gh-extension: deleted (no PR, nothing beyond the default branch)",
		"task/3-unmerged-work: kept (2 commit(s) not on the default branch, no merged PR)",
		"task/4-open: kept (open PR)",
		"task/5-closed-pr: kept (1 commit(s)",
		"task/7-forbidden: kept (PR merged; delete failed",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output missing %q:\n%s", want, out.String())
		}
	}
	if strings.Contains(out.String(), "task/6-gone") {
		t.Errorf("a branch that no longer exists must be skipped silently:\n%s", out.String())
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
	deleteHead(&out, ft, tracker.PR{Repo: "o/r"}) // no head known: silent
	if out.String() != "" {
		t.Errorf("empty head produced output: %q", out.String())
	}
}
