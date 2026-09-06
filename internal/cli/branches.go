package cli

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// branchAction decides whether a task branch may be deleted. Two grounds
// only: the branch's PR merged and the branch still sits at the commit
// that merged (rebase-merge rewrites commits, so ancestry of the tip is
// never the test — but a tip that moved after the merge is new work), or
// no PR is open and the branch carries nothing beyond the default branch.
// Everything else is kept, with the reason.
func branchAction(pr tracker.PR, hasPR bool, ahead int, tip string) (del bool, reason string) {
	switch {
	case hasPR && pr.Open:
		return false, "open PR"
	case hasPR && pr.Merged && pr.HeadSHA != "" && tip == pr.HeadSHA:
		return true, "PR merged"
	case hasPR && pr.Merged:
		return false, "commits pushed after the PR merged"
	case ahead == 0 && hasPR:
		return true, "PR closed unmerged, nothing beyond the default branch"
	case ahead == 0:
		return true, "no PR, nothing beyond the default branch"
	default:
		return false, fmt.Sprintf("%d commit(s) not on the default branch, no merged PR", ahead)
	}
}

// deleteHead removes a merged PR's head branch — the counterpart of task
// start creating it. The merge has already happened and been reported, so
// a failure here is a note, never an error. A fork's head is not ours to
// delete.
func deleteHead(w io.Writer, t tracker.Tracker, pr tracker.PR) {
	if pr.HeadRef == "" || pr.CrossRepo {
		return
	}
	if err := t.DeleteBranch(pr.Repo, pr.HeadRef); err != nil {
		fmt.Fprintf(w, "note: could not delete branch %s (%v); delete it by hand\n", pr.HeadRef, err)
		return
	}
	fmt.Fprintf(w, "deleted branch %s\n", pr.HeadRef)
}

// prByHead indexes a task's PRs by head branch. When several PRs share a
// head, the one that forbids deletion wins: open over merged over closed —
// never the last one the API happened to list.
func prByHead(prs []tracker.PR) map[string]tracker.PR {
	rank := func(pr tracker.PR) int {
		switch {
		case pr.Open:
			return 2
		case pr.Merged:
			return 1
		}
		return 0
	}
	byHead := map[string]tracker.PR{}
	for _, pr := range prs {
		if pr.CrossRepo || pr.HeadRef == "" {
			continue // a fork's branch is not in this repo
		}
		if cur, ok := byHead[pr.HeadRef]; !ok || rank(pr) > rank(cur) {
			byHead[pr.HeadRef] = pr
		}
	}
	return byHead
}

// sweepBranches deletes the task branches a milestone leaves behind that
// branchAction allows and reports every other one that still exists, so
// the close output records what was removed and what was left and why.
// Candidates are the head refs of the task's PRs (same-repo only), the
// conventional name task start used (task/<n>-<slug>, for tasks that never
// had a PR), and GitHub's linked-branch relation — which is dropped once a
// PR attaches, so it only ever contributes for PR-less tasks. The default
// branch is never a candidate. A candidate that no longer exists is
// skipped silently. Failures are reported and never abort: the caller has
// passed every gate by the time it sweeps.
func sweepBranches(w io.Writer, t tracker.Tracker, m *tracker.Milestone) (deleted []string) {
	return executeSweep(w, t, planSweep(t, m))
}

// sweepItem is one branch the close would consider — delete it, or keep
// it and say why — or a note about a lookup it had to skip, in the order
// the sweep meets them. Stale marks the second pass's items: a branch no
// milestone of this close's own is responsible for.
type sweepItem struct {
	Repo, Name string
	Delete     bool
	Reason     string
	Note       string
	Stale      bool
}

// planSweep decides, without touching anything, what the sweep would do
// to every branch a milestone's tasks left. The dry run prints the plan;
// the live close executes it.
func planSweep(t tracker.Tracker, m *tracker.Milestone) (items []sweepItem) {
	defaults := map[string]string{}
	for _, task := range m.Tasks {
		if _, ok := defaults[task.Repo]; !ok {
			info, err := t.RepoInfo(task.Repo)
			if err != nil {
				items = append(items, sweepItem{Note: fmt.Sprintf("note: branch sweep skipped for %s (%v)", task.Repo, err)})
				defaults[task.Repo] = ""
				continue
			}
			defaults[task.Repo] = info.DefaultBranch
		}
		if defaults[task.Repo] == "" {
			continue
		}
		nums, err := t.ClosingPRs(task, true)
		if err != nil {
			items = append(items, sweepItem{Note: fmt.Sprintf("note: branch sweep skipped for %s (%v)", task, err)})
			continue
		}
		var prs []tracker.PR
		for _, num := range nums {
			pr, err := t.PRInfo(task.Repo, num)
			if err != nil {
				items = append(items, sweepItem{Note: fmt.Sprintf("note: branch sweep skipped for %s (%v)", task, err)})
				prs = nil
				break
			}
			prs = append(prs, pr)
		}
		if prs == nil && len(nums) > 0 {
			continue
		}
		byHead := prByHead(prs)
		var names []string
		seen := map[string]bool{}
		add := func(n string) {
			if n != "" && n != defaults[task.Repo] && !seen[n] {
				seen[n] = true
				names = append(names, n)
			}
		}
		for n := range byHead {
			add(n)
		}
		if info, err := t.Task(task); err == nil {
			add(fmt.Sprintf("task/%d-%s", task.Number, slug(info.Title)))
		}
		if linked, err := t.LinkedBranches(task); err == nil {
			for _, n := range linked {
				add(n)
			}
		}
		for _, name := range names {
			ahead, tip, err := t.BranchAhead(task.Repo, name)
			if err != nil {
				continue // gone already, or never made
			}
			pr, hasPR := byHead[name]
			del, reason := branchAction(pr, hasPR, ahead, tip)
			items = append(items, sweepItem{Repo: task.Repo, Name: name, Delete: del, Reason: reason})
		}
	}
	return items
}

// executeSweep performs a sweep plan, reporting each branch's fate and
// each skipped lookup where the sweep met it.
func executeSweep(w io.Writer, t tracker.Tracker, items []sweepItem) (deleted []string) {
	for _, it := range items {
		if it.Note != "" {
			fmt.Fprintln(w, it.Note)
			continue
		}
		kind := "branch"
		if it.Stale {
			kind = "stale branch"
		}
		if !it.Delete {
			fmt.Fprintf(w, "%s %s: kept (%s)\n", kind, it.Name, it.Reason)
			continue
		}
		if err := t.DeleteBranch(it.Repo, it.Name); err != nil {
			fmt.Fprintf(w, "%s %s: kept (%s; delete failed: %v)\n", kind, it.Name, it.Reason, err)
			continue
		}
		fmt.Fprintf(w, "%s %s: deleted (%s)\n", kind, it.Name, it.Reason)
		deleted = append(deleted, it.Name)
	}
	return deleted
}

// taskBranchPrefix opens every branch `task start` cuts. The stale sweep
// lists under it and reads the task number back out of it.
const taskBranchPrefix = "task/"

// taskNumber reads the task issue number out of a branch `task start` cut —
// `task/<n>-<slug>`. Anything else yields 0 and is no candidate for a
// sweep: another prefix, no number, or a number written with a leading zero
// that `task start` would never have produced.
func taskNumber(branch string) int {
	rest, ok := strings.CutPrefix(branch, taskBranchPrefix)
	if !ok {
		return 0
	}
	digits, _, _ := strings.Cut(rest, "-")
	if digits == "" || (len(digits) > 1 && digits[0] == '0') {
		return 0
	}
	n, err := strconv.Atoi(digits)
	if err != nil || n <= 0 {
		return 0
	}
	return n
}

// planStaleSweep decides, without touching anything, what the close would
// do to the task branches EARLIER closes left behind. `milestone close`
// sweeps only the closing milestone's own tasks, so a branch whose task
// shipped under a milestone that closed before the sweep worked — or whose
// delete failed once — is invisible to every later verb and accumulates
// (#167). The candidate set is the change and the delete conditions are
// not: every branch here meets the same branchAction the milestone's own
// meet, which is also the only test that catches a rebase-merged branch,
// whose commits are rewritten and so are never `ahead == 0`.
//
// The repos are the hub the milestone issue lives in and every repo its
// tasks name — one prefix-filtered listing each. Branches the milestone's
// own pass already considered, and its own tasks', are left to it, so a
// verdict is never restated or overridden.
func planStaleSweep(t tracker.Tracker, m *tracker.Milestone, own []sweepItem) (items []sweepItem) {
	seen := map[string]bool{}
	for _, it := range own {
		if it.Name != "" {
			seen[it.Repo+" "+it.Name] = true
		}
	}
	mine := map[tracker.IssueRef]bool{}
	repos := []string{m.Ref.Repo}
	for _, task := range m.Tasks {
		mine[task] = true
		if !slices.Contains(repos, task.Repo) {
			repos = append(repos, task.Repo)
		}
	}
	repos = slices.DeleteFunc(repos, func(r string) bool { return r == "" })
	for _, repo := range repos {
		info, err := t.RepoInfo(repo)
		if err != nil {
			items = append(items, sweepItem{Note: fmt.Sprintf("note: stale branch sweep skipped for %s (%v)", repo, err)})
			continue
		}
		branches, err := t.TaskBranches(repo)
		if err != nil {
			items = append(items, sweepItem{Note: fmt.Sprintf("note: stale branch sweep skipped for %s (%v)", repo, err)})
			continue
		}
		for _, name := range branches {
			ref := tracker.IssueRef{Repo: repo, Number: taskNumber(name)}
			if ref.Number == 0 || name == info.DefaultBranch || seen[repo+" "+name] || mine[ref] {
				continue
			}
			seen[repo+" "+name] = true
			if item, ok := staleBranchAction(t, ref, name); ok {
				items = append(items, item)
			}
		}
	}
	return items
}

// staleBranchAction judges one branch an earlier close left behind. The
// task issue's state is read first, so a branch whose task is still open
// costs one issue read and stops there; only a closed task's branch goes on
// to the PR and comparison lookups the delete conditions need. ok is false
// when the branch is not the sweep's to report at all — it went away
// between the listing and the check.
func staleBranchAction(t tracker.Tracker, ref tracker.IssueRef, name string) (sweepItem, bool) {
	skip := func(err error) (sweepItem, bool) {
		return sweepItem{Note: fmt.Sprintf("note: stale branch %s skipped (%v)", name, err)}, true
	}
	task, err := t.Task(ref)
	if err != nil {
		return skip(err)
	}
	if !task.Closed {
		return sweepItem{Repo: ref.Repo, Name: name, Reason: fmt.Sprintf("%s is open", ref), Stale: true}, true
	}
	nums, err := t.ClosingPRs(ref, true)
	if err != nil {
		return skip(err)
	}
	var prs []tracker.PR
	for _, num := range nums {
		pr, err := t.PRInfo(ref.Repo, num)
		if err != nil {
			return skip(err)
		}
		prs = append(prs, pr)
	}
	ahead, tip, err := t.BranchAhead(ref.Repo, name)
	if err != nil {
		return sweepItem{}, false
	}
	pr, hasPR := prByHead(prs)[name]
	del, reason := branchAction(pr, hasPR, ahead, tip)
	return sweepItem{Repo: ref.Repo, Name: name, Delete: del, Reason: reason, Stale: true}, true
}
