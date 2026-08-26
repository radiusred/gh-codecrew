package cli

import (
	"fmt"
	"io"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// branchAction decides whether a task branch may be deleted. Only two
// facts justify deletion: the branch's PR merged (its content is on the
// default branch, whatever the rebase did to its commits), or no PR is
// open and the branch carries nothing beyond the default branch. Ancestry
// of the tip is never the test — rebase-merge rewrites the commits, so a
// fully merged branch's tip is usually not an ancestor of main.
func branchAction(merged, openPR bool, ahead int) (del bool, reason string) {
	switch {
	case merged:
		return true, "PR merged"
	case openPR:
		return false, "open PR"
	case ahead == 0:
		return true, "no PR, nothing beyond the default branch"
	default:
		return false, fmt.Sprintf("%d commit(s) not on the default branch, no merged PR", ahead)
	}
}

// deleteHead removes a merged PR's head branch — the counterpart of task
// start creating it. The merge has already happened and been reported, so
// a failure here is a note, never an error.
func deleteHead(w io.Writer, t tracker.Tracker, pr tracker.PR) {
	if pr.HeadRef == "" {
		return
	}
	if err := t.DeleteBranch(pr.Repo, pr.HeadRef); err != nil {
		fmt.Fprintf(w, "note: could not delete branch %s (%v); delete it by hand\n", pr.HeadRef, err)
		return
	}
	fmt.Fprintf(w, "deleted branch %s\n", pr.HeadRef)
}

// sweepBranches deletes the task branches a milestone leaves behind that
// branchAction allows and reports every other one that still exists, so
// the close output records what was removed and what was left and why.
// Candidates come from three places, because GitHub drops the linked-
// branch relation once an issue closes: the relation (open tasks), the
// head refs of the task's PRs, and the conventional name task start used
// (task/<n>-<slug>, for tasks that never had a PR). A candidate that no
// longer exists — deleted at task finish, or never created — is skipped
// silently. Returns the names deleted.
func sweepBranches(w io.Writer, t tracker.Tracker, m *tracker.Milestone) ([]string, error) {
	var deleted []string
	for _, task := range m.Tasks {
		linked, err := t.LinkedBranches(task)
		if err != nil {
			return deleted, err
		}
		prs, err := t.ClosingPRs(task, true)
		if err != nil {
			return deleted, err
		}
		byHead := map[string]tracker.PR{}
		var names []string
		seen := map[string]bool{}
		add := func(n string) {
			if n != "" && !seen[n] {
				seen[n] = true
				names = append(names, n)
			}
		}
		for _, n := range linked {
			add(n)
		}
		for _, num := range prs {
			pr, err := t.PRInfo(task.Repo, num)
			if err != nil {
				return deleted, err
			}
			byHead[pr.HeadRef] = pr
			add(pr.HeadRef)
		}
		if info, err := t.Task(task); err == nil {
			add(fmt.Sprintf("task/%d-%s", task.Number, slug(info.Title)))
		}
		for _, name := range names {
			ahead, err := t.BranchAhead(task.Repo, name)
			if err != nil {
				continue // gone already, or never made
			}
			pr, hasPR := byHead[name]
			del, reason := branchAction(hasPR && pr.Merged, hasPR && pr.Open, ahead)
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
	return deleted, nil
}
