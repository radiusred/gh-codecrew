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
	// TaskOpen marks the one keep the second pass makes without judging
	// the branch at all: its task issue is still open, so nothing has
	// finished with it. The sweep reports it like any other keep; a caller
	// that only wants finished work — `status`'s report — tells it apart
	// here rather than by reading the reason prose.
	TaskOpen bool
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
			add(fmt.Sprintf("%s%d-%s", tracker.TaskBranchPrefix, task.Number, slug(info.Title)))
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

// taskNumber reads the task issue number out of a branch `task start` cut —
// `task/<n>-<slug>`. Anything else yields 0 and is no candidate for a
// sweep: another prefix, no number, or a number written with a leading zero
// that `task start` would never have produced.
func taskNumber(branch string) int {
	rest, ok := strings.CutPrefix(branch, tracker.TaskBranchPrefix)
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
		branches, truncated, err := t.TaskBranches(repo)
		if err != nil {
			items = append(items, sweepItem{Note: fmt.Sprintf("note: stale branch sweep skipped for %s (%v)", repo, err)})
			continue
		}
		if truncated {
			// A page is the whole listing, so a repo with more task
			// branches than one holds would otherwise be swept in silent
			// part. Successive closes make progress; the operator is told
			// there is more rather than left to infer it.
			items = append(items, sweepItem{Note: fmt.Sprintf("note: %s carries more task branches than one listing holds; this sweep saw %d of them and a later close reaches the rest", repo, len(branches))})
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
		return sweepItem{Repo: ref.Repo, Name: name, Reason: fmt.Sprintf("%s is open", ref), Stale: true, TaskOpen: true}, true
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
	if del && !hasPR {
		// branchAction's "no PR" arm is only as good as the PR set it was
		// handed, and ClosingPRs finds the PRs that close <n> — so a PR
		// opened on this branch that carries no `Closes #<n>` line is
		// invisible to it, and deleting the branch would close that PR.
		// Pass one's candidates come from the task's own PRs; pass two's
		// come from a repo-wide listing, so the exposure is wider than the
		// arm was written for (checky on PR #293). One lookup, and only for
		// a branch about to go on the strength of there being no PR at all.
		open, err := t.OpenPRsForBranch(ref.Repo, name)
		if err != nil {
			return skip(err)
		}
		if len(open) > 0 {
			del, reason = false, fmt.Sprintf("open PR #%d on this branch", open[0])
		}
	}
	return sweepItem{Repo: ref.Repo, Name: name, Delete: del, Reason: reason, Stale: true}, true
}

// planStaleReport judges every task branch one repo carries the way the
// close's second pass would, and writes nothing — the report `status`
// prints (#295). It is the same staleBranchAction, so a verdict here and a
// verdict at the next close can never disagree; only the wording around it
// differs, a report speaking of what a close would do rather than of what
// this run is about to.
//
// The cost is the reason it was ever an open question: one prefix-filtered
// listing, then one issue read per task branch — and the PR and comparison
// lookups only for the branches whose task has closed, which is what
// staleBranchAction's early return buys. listed is how many branches the
// listing held, and truncated carries GitHub's hasNextPage back to the
// caller unworded: the sweep's "a later close reaches the rest" is a
// promise a report cannot make.
func planStaleReport(t tracker.Tracker, repo, defaultBranch string) (items []sweepItem, listed int, truncated bool, err error) {
	branches, truncated, err := t.TaskBranches(repo)
	if err != nil {
		return nil, 0, false, err
	}
	for _, name := range branches {
		ref := tracker.IssueRef{Repo: repo, Number: taskNumber(name)}
		if ref.Number == 0 || name == defaultBranch {
			continue
		}
		if item, ok := staleBranchAction(t, ref, name); ok {
			items = append(items, item)
		}
	}
	return items, len(branches), truncated, nil
}
