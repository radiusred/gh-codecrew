package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// TestMain pins the clone lookup to "no clone" for the whole package: the
// local cleanup would otherwise read whatever checkout the test binary
// happens to be running in, and the finish tests are about the gates. A
// test that wants a clone stands one up with stubClone.
func TestMain(m *testing.M) {
	cloneRoot = func() string { return "" }
	os.Exit(m.Run())
}

func stubClone(t *testing.T, dir string) {
	t.Helper()
	prev := cloneRoot
	cloneRoot = func() string { return dir }
	t.Cleanup(func() { cloneRoot = prev })
}

// repoFake answers the one tracker question the cleanup asks.
type repoFake struct {
	tracker.Tracker
	def string
}

func (r repoFake) RepoInfo(string) (tracker.RepoInfo, error) {
	return tracker.RepoInfo{DefaultBranch: r.def}, nil
}

func commitFile(t *testing.T, dir, name, content, msg string) string {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", name}, {"commit", "-q", "-m", msg}} {
		if _, err := git(dir, args...); err != nil {
			t.Fatal(err)
		}
	}
	sha, err := git(dir, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	return sha
}

func mustGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := git(dir, args...)
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return out
}

// taskClone builds the shape task finish meets: a remote holding main and
// a task branch, and a clone of it with both branches local. It returns the
// clone, the remote and the task branch's tip — the commit GitHub would
// report as the PR's head.
func taskClone(t *testing.T, branch string) (dir, remote, tip string) {
	t.Helper()
	remote = gitRepo(t)
	commitFile(t, remote, "a.txt", "one\n", "feat: one")
	mustGit(t, remote, "switch", "-q", "-c", branch)
	tip = commitFile(t, remote, "b.txt", "two\n", "feat: two")
	mustGit(t, remote, "switch", "-q", "main")

	dir = filepath.Join(t.TempDir(), "clone")
	mustGit(t, t.TempDir(), "clone", "-q", remote, dir)
	for _, args := range [][]string{{"config", "user.email", "t@example.com"}, {"config", "user.name", "t"}, {"config", "commit.gpgsign", "false"}, {"switch", "-q", "-c", branch, "origin/" + branch}} {
		mustGit(t, dir, args...)
	}
	return dir, remote, tip
}

// merged is the state task finish leaves on the remote: main carrying a
// commit that is not the branch's (a rebase-merge rewrites them) and the
// head branch gone.
func merged(t *testing.T, remote, branch string) string {
	t.Helper()
	sha := commitFile(t, remote, "b.txt", "two\n", "feat: two (rebased)")
	mustGit(t, remote, "branch", "-q", "-D", branch)
	return sha
}

func mergedPR(headRef, headSHA string) tracker.PR {
	return tracker.PR{Repo: "o/r", Number: 9, HeadRef: headRef, HeadSHA: headSHA, Merged: true}
}

func cleanupOutput(t *testing.T, pr tracker.PR) string {
	t.Helper()
	var buf bytes.Buffer
	planClone(repoFake{def: "main"}, pr, "o/r", true).run(&buf)
	return buf.String()
}

func branchExists(t *testing.T, dir, branch string) bool {
	t.Helper()
	_, err := git(dir, "rev-parse", "--verify", "--quiet", "refs/heads/"+branch)
	return err == nil
}

func wantLines(t *testing.T, out string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(out, w) {
			t.Errorf("output lacks %q:\n%s", w, out)
		}
	}
}

// Standing on the task branch when it merges: off it, main caught up with
// the merge, the branch gone.
func TestCloneCleanupOnTaskBranch(t *testing.T) {
	const branch = "task/9-x"
	dir, remote, tip := taskClone(t, branch)
	stubClone(t, dir)
	head := merged(t, remote, branch)

	out := cleanupOutput(t, mergedPR(branch, tip))
	wantLines(t, out,
		"local: switched to main",
		"local: fast-forwarded main to "+shortSHA(head),
		"local: deleted branch "+branch+" (at the merged head "+shortSHA(tip)+")")
	if got := mustGit(t, dir, "symbolic-ref", "--short", "HEAD"); got != "main" {
		t.Errorf("HEAD is %q", got)
	}
	if got := mustGit(t, dir, "rev-parse", "main"); got != head {
		t.Errorf("main is at %s, want the merge %s", got, head)
	}
	if branchExists(t, dir, branch) {
		t.Error("the task branch survived")
	}
	if _, err := git(dir, "rev-parse", "--verify", "--quiet", "refs/remotes/origin/"+branch); err == nil {
		t.Error("the remote-tracking ref survived the prune")
	}
}

// The branch is held locally but not checked out: deleted just the same,
// and HEAD is not moved.
func TestCloneCleanupBranchNotCheckedOut(t *testing.T) {
	const branch = "task/9-x"
	dir, remote, tip := taskClone(t, branch)
	mustGit(t, dir, "switch", "-q", "main")
	stubClone(t, dir)
	head := merged(t, remote, branch)

	out := cleanupOutput(t, mergedPR(branch, tip))
	if strings.Contains(out, "switched to") {
		t.Errorf("switched while already on main:\n%s", out)
	}
	wantLines(t, out, "local: fast-forwarded main to "+shortSHA(head), "local: deleted branch "+branch)
	if got := mustGit(t, dir, "symbolic-ref", "--short", "HEAD"); got != "main" {
		t.Errorf("HEAD is %q", got)
	}
	if branchExists(t, dir, branch) {
		t.Error("the task branch survived")
	}
}

// Standing on a third branch: the default branch moves by its ref, since
// nothing has it checked out, and the third branch is left alone.
func TestCloneCleanupFromAnotherBranch(t *testing.T) {
	const branch = "task/9-x"
	dir, remote, tip := taskClone(t, branch)
	mustGit(t, dir, "switch", "-q", "-c", "scratch")
	stubClone(t, dir)
	head := merged(t, remote, branch)

	out := cleanupOutput(t, mergedPR(branch, tip))
	wantLines(t, out, "local: fast-forwarded main to "+shortSHA(head), "local: deleted branch "+branch)
	if got := mustGit(t, dir, "symbolic-ref", "--short", "HEAD"); got != "scratch" {
		t.Errorf("HEAD is %q", got)
	}
	if got := mustGit(t, dir, "rev-parse", "main"); got != head {
		t.Errorf("main is at %s, want the merge %s", got, head)
	}
}

// A commit the remote never saw is work, whatever the PR merged: the
// branch is named and kept.
func TestCloneCleanupKeepsUnpushedCommits(t *testing.T) {
	const branch = "task/9-x"
	dir, remote, tip := taskClone(t, branch)
	stubClone(t, dir)
	extra := commitFile(t, dir, "c.txt", "three\n", "feat: three")
	head := merged(t, remote, branch)

	out := cleanupOutput(t, mergedPR(branch, tip))
	wantLines(t, out,
		"local: switched to main",
		"local: fast-forwarded main to "+shortSHA(head),
		"local: kept branch "+branch+" (1 commit(s) beyond the merged head "+shortSHA(tip)+")")
	if !branchExists(t, dir, branch) {
		t.Fatal("the branch was deleted with work on it")
	}
	if got := mustGit(t, dir, "rev-parse", branch); got != extra {
		t.Errorf("the branch moved: %s", got)
	}
}

// A local default branch that is not an ancestor of the fetched one is not
// fast-forwardable: it is named and left as it is, and the merged task
// branch still goes.
func TestCloneCleanupDivergedDefaultBranch(t *testing.T) {
	const branch = "task/9-x"
	dir, remote, tip := taskClone(t, branch)
	mustGit(t, dir, "switch", "-q", "main")
	local := commitFile(t, dir, "local.txt", "mine\n", "chore: local")
	mustGit(t, dir, "switch", "-q", branch)
	stubClone(t, dir)
	merged(t, remote, branch)

	out := cleanupOutput(t, mergedPR(branch, tip))
	wantLines(t, out,
		"local: main has commits origin/main does not — left alone",
		"local: deleted branch "+branch)
	if strings.Contains(out, "fast-forwarded") {
		t.Errorf("fast-forwarded a diverged branch:\n%s", out)
	}
	if got := mustGit(t, dir, "rev-parse", "main"); got != local {
		t.Errorf("main moved to %s, want %s", got, local)
	}
}

// Nothing to do, nothing said: outside a repository, in a clone of another
// repo, when the branch is not held locally, and for a fork's head.
func TestCloneCleanupSilentElsewhere(t *testing.T) {
	const branch = "task/9-x"
	dir, _, tip := taskClone(t, branch)
	cases := map[string]struct {
		root string
		pr   tracker.PR
		cur  string
	}{
		"no repository":    {"", mergedPR(branch, tip), "o/r"},
		"another repo":     {dir, mergedPR(branch, tip), "o/other"},
		"no such branch":   {dir, mergedPR("task/9-gone", tip), "o/r"},
		"a fork's head":    {dir, tracker.PR{Repo: "o/r", HeadRef: branch, HeadSHA: tip, CrossRepo: true}, "o/r"},
		"no head at all":   {dir, tracker.PR{Repo: "o/r"}, "o/r"},
		"the head is main": {dir, mergedPR("main", tip), "o/r"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			stubClone(t, tc.root)
			cc := planClone(repoFake{def: "main"}, tc.pr, tc.cur, true)
			if cc != nil {
				t.Fatalf("planned %+v", cc)
			}
			var buf bytes.Buffer
			cc.run(&buf)
			p := &plan{}
			cc.would(p)
			if buf.Len() != 0 || len(p.actions) != 0 {
				t.Errorf("said %q / %v", buf.String(), p.actions)
			}
		})
	}
}

// The dry run names the steps and touches nothing: no fetch, no switch, no
// deletion — every ref exactly where it was.
func TestCloneCleanupDryRunChangesNothing(t *testing.T) {
	const branch = "task/9-x"
	dir, remote, tip := taskClone(t, branch)
	stubClone(t, dir)
	merged(t, remote, branch)
	before := mustGit(t, dir, "for-each-ref", "--format=%(refname) %(objectname)")

	p := &plan{}
	planClone(repoFake{def: "main"}, mergedPR(branch, tip), "o/r", false).would(p)
	var buf bytes.Buffer
	p.print(&buf)
	wantLines(t, buf.String(),
		"would fetch origin in this clone ("+dir+")",
		"would switch this clone to main",
		"would fast-forward main in this clone to the merge",
		"would delete the local branch "+branch+" (at the merged head "+shortSHA(tip)+")")
	if after := mustGit(t, dir, "for-each-ref", "--format=%(refname) %(objectname)"); after != before {
		t.Errorf("the dry run moved refs:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if got := mustGit(t, dir, "symbolic-ref", "--short", "HEAD"); got != branch {
		t.Errorf("HEAD is %q", got)
	}
}

// task finish's own report: the local steps stand beside the remote ones
// in the dry run, and the live run performs them after the merge.
func TestFinishReportsLocalCleanup(t *testing.T) {
	const branch = "task/7-x"
	dir, remote, tip := taskClone(t, branch)
	stubClone(t, dir)

	f := cleanFinish()
	f.pr.HeadRef, f.pr.HeadSHA = branch, tip
	p, runFinish, err := planFinish(finishCtx(f, crewRoles), f.task.Ref, false, false)
	if err != nil || p.refusal != nil {
		t.Fatalf("plan: err %v refusal %v", err, p.refusal)
	}
	var buf bytes.Buffer
	p.print(&buf)
	wantLines(t, buf.String(), "would delete head "+branch, "would switch this clone to main", "would delete the local branch "+branch)

	head := merged(t, remote, branch)
	buf.Reset()
	if err := runFinish(&buf); err != nil {
		t.Fatal(err)
	}
	wantLines(t, buf.String(), "merged PR #9", "local: switched to main", "local: fast-forwarded main to "+shortSHA(head), "local: deleted branch "+branch)
	if branchExists(t, dir, branch) {
		t.Error("the task branch survived the finish")
	}
}

// A branch whose commits are on the fetched default branch anyway is the
// second ground for deleting: nothing is lost with it, whatever GitHub
// reported as the head.
func TestCloneCleanupDeletesBranchContainedInDefault(t *testing.T) {
	const branch = "task/9-x"
	dir, remote, _ := taskClone(t, branch)
	stubClone(t, dir)
	mustGit(t, remote, "switch", "-q", "main")
	mustGit(t, remote, "merge", "-q", "--no-ff", "-m", "merge: the branch", branch)
	mustGit(t, remote, "branch", "-q", "-D", branch)

	out := cleanupOutput(t, mergedPR(branch, "0123456789012345678901234567890123456789"))
	wantLines(t, out, "local: deleted branch "+branch+" (contained in origin/main)")
	if branchExists(t, dir, branch) {
		t.Error("the branch survived")
	}
}

// Neither ground holds when GitHub reported no head commit at all: the
// branch is kept rather than force-deleted on a guess.
func TestCloneCleanupKeepsBranchWithNoMergedHead(t *testing.T) {
	const branch = "task/9-x"
	dir, remote, _ := taskClone(t, branch)
	stubClone(t, dir)
	merged(t, remote, branch)

	out := cleanupOutput(t, mergedPR(branch, ""))
	wantLines(t, out, "local: kept branch "+branch+" (GitHub did not report the commit it merged)")
	if !branchExists(t, dir, branch) {
		t.Error("the branch was deleted on no ground at all")
	}
}
