package cli

import (
	"fmt"
	"io"

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
	defaults := map[string]string{}
	for _, task := range m.Tasks {
		if _, ok := defaults[task.Repo]; !ok {
			info, err := t.RepoInfo(task.Repo)
			if err != nil {
				fmt.Fprintf(w, "note: branch sweep skipped for %s (%v)\n", task.Repo, err)
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
			fmt.Fprintf(w, "note: branch sweep skipped for %s (%v)\n", task, err)
			continue
		}
		var prs []tracker.PR
		for _, num := range nums {
			pr, err := t.PRInfo(task.Repo, num)
			if err != nil {
				fmt.Fprintf(w, "note: branch sweep skipped for %s (%v)\n", task, err)
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
			if !del {
				fmt.Fprintf(w, "branch %s: kept (%s)\n", name, reason)
				continue
			}
			if err := t.DeleteBranch(task.Repo, name); err != nil {
				fmt.Fprintf(w, "branch %s: kept (%s; delete failed: %v)\n", name, reason, err)
				continue
			}
			fmt.Fprintf(w, "branch %s: deleted (%s)\n", name, reason)
			deleted = append(deleted, name)
		}
	}
	return deleted
}
