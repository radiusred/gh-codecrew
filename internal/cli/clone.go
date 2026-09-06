package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// cloneRoot is the repository the verb is running in, "" outside one. A
// func var so a test can stand a temp clone under the verb rather than
// chdir the whole test binary into one (the ghVersion pattern).
var cloneRoot = func() string { return repoRoot(".") }

// cloneCleanup is task finish's local half: what the clone the verb runs in
// needs once the merge has landed on GitHub and the remote head is gone.
// The protocol knows the branch it just merged, so the operator should not
// have to sweep it by hand (#192) — but only in a clone of the task's own
// repo, and never at the cost of a commit that is not on the merged head.
type cloneCleanup struct {
	dir      string // the clone's top level
	def      string // the default branch's name
	branch   string // the task branch, as the merged PR's head
	head     string // the branch checked out when the verb ran ("" if detached)
	switchTo bool   // HEAD is on the task branch and must leave it first
	ff       ffState
	ffTo     string // what the default branch moves to, for the report
	del      bool
	reason   string // why the branch goes, or why it stays
	fetchErr error
}

// ffState is what the local default branch can do with the fetched one.
type ffState int

const (
	ffNone     ffState = iota // already there, or no such branch to move
	ffAhead                   // fast-forwardable: the local branch is an ancestor
	ffDiverged                // carries commits the remote does not: left alone
)

// planClone decides the whole of the local cleanup from the clone's own
// state and writes nothing but the fetch — the one preparation the decision
// needs, and the reason the dry run passes fetch false: it reads the clone
// as it is, and before the merge there is nothing on the remote to fetch
// anyway. nil means there is nothing to do and nothing to print, which is
// the answer for every clone that is not this task's: no git repository, a
// different repo (a hub's clone met by a spoke's task), a fork's head, or
// no local branch of that name. Cheapest questions first, so "somewhere
// else" costs one git call and no API call.
func planClone(t tracker.Tracker, pr tracker.PR, current string, fetch bool) *cloneCleanup {
	if pr.HeadRef == "" || pr.CrossRepo || current != pr.Repo {
		return nil
	}
	dir := cloneRoot()
	if dir == "" {
		return nil
	}
	tip, err := git(dir, "rev-parse", "--verify", "--quiet", "refs/heads/"+pr.HeadRef)
	if err != nil {
		return nil
	}
	cc := &cloneCleanup{dir: dir, branch: pr.HeadRef}
	cc.def = defaultBranch(t, dir, pr.Repo)
	if cc.def == "" || cc.def == pr.HeadRef {
		return nil
	}
	if fetch {
		// Last of the conditions, so a clone with nothing to clean is left
		// as untouched as one in another repo. --prune retires the
		// remote-tracking ref the head deletion just orphaned, which is
		// the other half of the mess #192 reported.
		_, cc.fetchErr = git(dir, "fetch", "--prune", "origin")
	}
	cc.head, _ = git(dir, "symbolic-ref", "--short", "--quiet", "HEAD")
	cc.switchTo = cc.head == pr.HeadRef

	remote, rerr := git(dir, "rev-parse", "--verify", "--quiet", "refs/remotes/origin/"+cc.def)
	local, lerr := git(dir, "rev-parse", "--verify", "--quiet", "refs/heads/"+cc.def)
	switch {
	case rerr != nil || lerr != nil:
		// No remote-tracking ref, or no local default branch at all: there
		// is nothing to fast-forward, and a switch creates the branch at
		// origin's tip by itself.
		cc.ff = ffNone
	case local == remote && fetch:
		cc.ff = ffNone
	case !isAncestor(dir, local, remote):
		cc.ff = ffDiverged
	case local == remote:
		// The dry run has not fetched, so the merge is not in origin's tip
		// yet; the step it names is the one the live verb will take.
		cc.ff, cc.ffTo = ffAhead, "the merge"
	default:
		cc.ff, cc.ffTo = ffAhead, shortSHA(remote)
	}
	cc.del, cc.reason = localBranchAction(dir, tip, pr.HeadSHA, remote, cc.def)
	return cc
}

// defaultBranch is the branch the cleanup fast-forwards: GitHub's answer,
// or the clone's own origin/HEAD when the API cannot be asked. Empty when
// neither knows, which stops the cleanup rather than guessing at "main".
func defaultBranch(t tracker.Tracker, dir, repo string) string {
	if info, err := t.RepoInfo(repo); err == nil && info.DefaultBranch != "" {
		return info.DefaultBranch
	}
	if head, err := git(dir, "rev-parse", "--abbrev-ref", "origin/HEAD"); err == nil {
		return strings.TrimPrefix(head, "origin/")
	}
	return ""
}

// localBranchAction decides whether the local task branch may go. Two
// grounds, and the reason is worded to read the same in the dry run's
// "would delete" and the live verb's "deleted": the branch sits at the
// commit GitHub merged, or it is contained in the fetched default branch
// anyway. Ancestry of the tip is never the test on its own — rebase-merge
// rewrites the commits, which is why `git branch -d` refuses a merged task
// branch and the deletion has to be `-D` (#192); pinning the tip to the
// merged head is what makes that force safe. Anything else is work the
// remote has not got, and it is named and kept.
func localBranchAction(dir, tip, merged, defaultSHA, def string) (del bool, reason string) {
	switch {
	case merged != "" && tip == merged:
		return true, "at the merged head " + shortSHA(merged)
	case defaultSHA != "" && isAncestor(dir, tip, defaultSHA):
		return true, "contained in origin/" + def
	case merged != "" && isAncestor(dir, merged, tip):
		return false, fmt.Sprintf("%d commit(s) beyond the merged head %s", commitsAhead(dir, merged, tip), shortSHA(merged))
	case merged == "":
		return false, "GitHub did not report the commit it merged"
	default:
		return false, fmt.Sprintf("its tip %s is not the merged head %s", shortSHA(tip), shortSHA(merged))
	}
}

// isAncestor reports whether commit a is reachable from b (true when equal).
func isAncestor(dir, a, b string) bool {
	_, err := git(dir, "merge-base", "--is-ancestor", a, b)
	return err == nil
}

// commitsAhead is how many commits b carries beyond a.
func commitsAhead(dir, a, b string) int {
	out, err := git(dir, "rev-list", "--count", a+".."+b)
	if err != nil {
		return 0
	}
	n := 0
	if _, err := fmt.Sscanf(out, "%d", &n); err != nil {
		return 0
	}
	return n
}

func shortSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
}

// would names the local steps in a task finish plan, in the order run
// takes them — the dry run's half of the one code path (#133).
func (cc *cloneCleanup) would(p *plan) {
	if cc == nil {
		return
	}
	p.would("fetch origin in this clone (%s)", cc.dir)
	if cc.switchTo {
		p.would("switch this clone to %s", cc.def)
	}
	switch cc.ff {
	case ffAhead:
		p.would("fast-forward %s in this clone to %s", cc.def, cc.ffTo)
	case ffDiverged:
		p.would("leave %s in this clone alone: it has commits origin/%s does not", cc.def, cc.def)
	}
	if cc.del {
		p.would("delete the local branch %s (%s)", cc.branch, cc.reason)
	} else {
		p.would("keep the local branch %s (%s)", cc.branch, cc.reason)
	}
}

// run performs the cleanup, reporting every step it takes. Nothing here is
// an error: the merge has happened and been reported, so a clone that will
// not cooperate is a note naming what the operator is left holding — the
// same standing deleteHead has for the remote branch.
func (cc *cloneCleanup) run(w io.Writer) {
	if cc == nil {
		return
	}
	if cc.fetchErr != nil {
		fmt.Fprintf(w, "local: could not fetch origin (%v) — reading %s as it stands\n", cc.fetchErr, cc.dir)
	} else {
		fmt.Fprintf(w, "local: fetched origin in %s\n", cc.dir)
	}
	if cc.switchTo {
		if _, err := git(cc.dir, "switch", cc.def); err != nil {
			fmt.Fprintf(w, "local: could not switch to %s (%v) — %s is left checked out\n", cc.def, err, cc.branch)
			return
		}
		fmt.Fprintf(w, "local: switched to %s\n", cc.def)
	}
	switch cc.ff {
	case ffAhead:
		var err error
		if cc.switchTo || cc.head == cc.def {
			_, err = git(cc.dir, "merge", "--ff-only", "origin/"+cc.def)
		} else {
			// Not checked out, so the ref moves on its own; branch -f
			// refuses a checked-out branch, which is the guard behind this.
			_, err = git(cc.dir, "branch", "-f", cc.def, "origin/"+cc.def)
		}
		if err != nil {
			fmt.Fprintf(w, "local: could not fast-forward %s (%v)\n", cc.def, err)
		} else {
			fmt.Fprintf(w, "local: fast-forwarded %s to %s\n", cc.def, cc.ffTo)
		}
	case ffDiverged:
		fmt.Fprintf(w, "local: %s has commits origin/%s does not — left alone\n", cc.def, cc.def)
	}
	if !cc.del {
		fmt.Fprintf(w, "local: kept branch %s (%s)\n", cc.branch, cc.reason)
		return
	}
	if _, err := git(cc.dir, "branch", "-D", cc.branch); err != nil {
		fmt.Fprintf(w, "local: could not delete branch %s (%v) — delete it by hand\n", cc.branch, err)
		return
	}
	fmt.Fprintf(w, "local: deleted branch %s (%s)\n", cc.branch, cc.reason)
}
