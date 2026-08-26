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

// sweepBranches deletes the linked branches of a milestone's tasks that
// branchAction allows and reports every other one, so the close output
// records what was removed and what was left and why. Returns the names
// deleted.
func sweepBranches(w io.Writer, t tracker.Tracker, m *tracker.Milestone) ([]string, error) {
	var deleted []string
	for _, task := range m.Tasks {
		branches, err := t.LinkedBranches(task)
		if err != nil {
			return deleted, err
		}
		if len(branches) == 0 {
			continue
		}
		// One lookup of the task's PRs serves every linked branch.
		prs, err := t.ClosingPRs(task, true)
		if err != nil {
			return deleted, err
		}
		byHead := map[string]tracker.PR{}
		for _, num := range prs {
			pr, err := t.PRInfo(task.Repo, num)
			if err != nil {
				return deleted, err
			}
			byHead[pr.HeadRef] = pr
		}
		for _, name := range branches {
			pr, hasPR := byHead[name]
			ahead := 0
			if !hasPR || (!pr.Merged && !pr.Open) {
				if ahead, err = t.BranchAhead(task.Repo, name); err != nil {
					fmt.Fprintf(w, "branch %s: kept (%v)\n", name, err)
					continue
				}
			}
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
