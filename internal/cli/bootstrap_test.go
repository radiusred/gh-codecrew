package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func gitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@example.com"}, {"config", "user.name", "t"}, {"config", "commit.gpgsign", "false"}} {
		if _, err := git(dir, args...); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func stubProtection(t *testing.T, required, known bool) {
	t.Helper()
	prev := defaultRequiresPR
	defaultRequiresPR = func(string) (bool, bool, string) { return required, known, "main" }
	t.Cleanup(func() { defaultRequiresPR = prev })
}

func committedFiles(t *testing.T, dir, rev string) string {
	t.Helper()
	out, err := git(dir, "show", "--name-only", "--format=%s", rev)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// A fresh repository: the scaffold is the root commit, holding exactly
// the files init wrote, on the branch the operator was on.
func TestCommitScaffoldRootCommit(t *testing.T) {
	stubProtection(t, false, true)
	dir := gitRepo(t)
	written, _, err := scaffold(dir, "self", fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	commitScaffold(&out, dir, written)
	if !strings.Contains(out.String(), "committed ") || !strings.Contains(out.String(), " on main: \"chore: scaffold codecrew\"") {
		t.Fatalf("output:\n%s", out.String())
	}
	show := committedFiles(t, dir, "HEAD")
	if !strings.HasPrefix(show, "chore: scaffold codecrew\n") {
		t.Errorf("subject: %q", show)
	}
	for _, f := range written {
		if !strings.Contains(show, f) {
			t.Errorf("%s not in the commit:\n%s", f, show)
		}
	}
	if n, _ := git(dir, "rev-list", "--count", "HEAD"); n != "1" {
		t.Errorf("%s commits, want the root commit only", n)
	}
	if st, _ := git(dir, "status", "--porcelain"); st != "" {
		t.Errorf("tree not clean after the commit: %q", st)
	}
	// A rerun writes nothing and commits nothing.
	out.Reset()
	written, _, _ = scaffold(dir, "self", fakeContracts)
	commitScaffold(&out, dir, written)
	if out.Len() != 0 {
		t.Errorf("rerun output: %q", out.String())
	}
	if n, _ := git(dir, "rev-list", "--count", "HEAD"); n != "1" {
		t.Errorf("rerun added a commit")
	}
}

// The operator's own work is untouched: a modified tracked file stays
// modified and unstaged, a staged unrelated file stays staged and out of
// the scaffold commit — no stash, no -A.
func TestCommitScaffoldLeavesTheOperatorsWorkAlone(t *testing.T) {
	stubProtection(t, false, true)
	dir := gitRepo(t)
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("one\n"), 0o644)
	git(dir, "add", "README.md")
	git(dir, "commit", "-q", "-m", "initial")
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("two\n"), 0o644) // modified, unstaged
	os.WriteFile(filepath.Join(dir, "staged.txt"), []byte("s\n"), 0o644)
	git(dir, "add", "staged.txt") // staged, unrelated
	written, _, _ := scaffold(dir, "self", fakeContracts)
	var out bytes.Buffer
	commitScaffold(&out, dir, written)
	show := committedFiles(t, dir, "HEAD")
	if strings.Contains(show, "README.md") || strings.Contains(show, "staged.txt") {
		t.Errorf("the operator's files rode along:\n%s", show)
	}
	unstaged, _ := git(dir, "diff", "--name-only")
	staged, _ := git(dir, "diff", "--cached", "--name-only")
	if unstaged != "README.md" || staged != "staged.txt" {
		t.Errorf("after: unstaged %q, staged %q — want README.md modified-unstaged and staged.txt still staged", unstaged, staged)
	}
	if n, _ := git(dir, "rev-list", "--count", "HEAD"); n != "2" {
		t.Errorf("%s commits, want 2", n)
	}
}

// A default branch that requires pull requests: the commit goes on
// codecrew-bootstrap, cut from HEAD, and the original branch is untouched;
// when the question cannot be asked, the branch is assumed and the note
// says so.
func TestCommitScaffoldProtectedDefaultBranch(t *testing.T) {
	for _, tc := range []struct {
		name            string
		required, known bool
		note            string
	}{
		{"protected", true, true, "requires pull requests: git push -u origin codecrew-bootstrap"},
		{"unknown", false, false, "could not ask GitHub"},
	} {
		stubProtection(t, tc.required, tc.known)
		dir := gitRepo(t)
		os.WriteFile(filepath.Join(dir, "README.md"), []byte("one\n"), 0o644)
		git(dir, "add", "README.md")
		git(dir, "commit", "-q", "-m", "initial")
		base, _ := git(dir, "rev-parse", "HEAD")
		written, _, _ := scaffold(dir, "self", fakeContracts)
		var out bytes.Buffer
		commitScaffold(&out, dir, written)
		if !strings.Contains(out.String(), " on codecrew-bootstrap: ") || !strings.Contains(out.String(), tc.note) {
			t.Errorf("%s: output:\n%s", tc.name, out.String())
		}
		if b, _ := git(dir, "symbolic-ref", "--short", "HEAD"); b != bootstrapBranch {
			t.Errorf("%s: on %s", tc.name, b)
		}
		if m, _ := git(dir, "rev-parse", "main"); m != base {
			t.Errorf("%s: main moved", tc.name)
		}
		if p, _ := git(dir, "rev-parse", "HEAD~1"); p != base {
			t.Errorf("%s: bootstrap not cut from HEAD", tc.name)
		}
	}
	// Protected and on a feature branch: the bootstrap branch is cut from
	// the default branch, so the scaffold PR carries the scaffold only.
	stubProtection(t, true, true)
	dir2 := gitRepo(t)
	os.WriteFile(filepath.Join(dir2, "README.md"), []byte("one\n"), 0o644)
	git(dir2, "add", "README.md")
	git(dir2, "commit", "-q", "-m", "initial")
	git(dir2, "switch", "-q", "-c", "feature")
	os.WriteFile(filepath.Join(dir2, "feature.txt"), []byte("f\n"), 0o644)
	git(dir2, "add", "feature.txt")
	git(dir2, "commit", "-q", "-m", "feature work")
	written2, _, _ := scaffold(dir2, "self", fakeContracts)
	var out2 bytes.Buffer
	commitScaffold(&out2, dir2, written2)
	if n, _ := git(dir2, "rev-list", "--count", "main..codecrew-bootstrap"); n != "1" {
		t.Errorf("feature branch: %s commits on the bootstrap branch beyond main, want the scaffold only", n)
	}
	if show := committedFiles(t, dir2, "codecrew-bootstrap"); strings.Contains(show, "feature.txt") {
		t.Errorf("feature work rode along:\n%s", show)
	}
	if !strings.Contains(out2.String(), "git push -u origin codecrew-bootstrap") {
		t.Errorf("push instruction: %q", out2.String())
	}
	// An unborn HEAD cannot branch: the root commit lands on the current branch even when protection is assumed.
	stubProtection(t, false, false)
	dir := gitRepo(t)
	written, _, _ := scaffold(dir, "self", fakeContracts)
	var out bytes.Buffer
	commitScaffold(&out, dir, written)
	if b, _ := git(dir, "symbolic-ref", "--short", "HEAD"); b != "main" || !strings.Contains(out.String(), " on main: ") {
		t.Errorf("unborn: on %s, output %q", b, out.String())
	}
}

// A commit that cannot be made is a note; the scaffold stays on disk.
func TestCommitScaffoldFailureIsANote(t *testing.T) {
	stubProtection(t, false, true)
	dir := gitRepo(t)
	git(dir, "config", "--unset", "user.email")
	git(dir, "config", "--unset", "user.name")
	t.Setenv("GIT_AUTHOR_NAME", "")
	t.Setenv("GIT_COMMITTER_NAME", "")
	t.Setenv("GIT_AUTHOR_EMAIL", "")
	t.Setenv("GIT_COMMITTER_EMAIL", "")
	t.Setenv("EMAIL", "")
	t.Setenv("HOME", t.TempDir()) // no global identity either
	written, _, _ := scaffold(dir, "self", fakeContracts)
	var out bytes.Buffer
	commitScaffold(&out, dir, written)
	if !strings.Contains(out.String(), "note: could not commit the scaffold") || !strings.Contains(out.String(), "commit it by hand") {
		t.Errorf("output:\n%s", out.String())
	}
	if _, err := os.Stat(filepath.Join(dir, ".codecrew.yml")); err != nil {
		t.Error("scaffold lost")
	}
	if _, err := git(dir, "rev-parse", "--verify", "HEAD"); err == nil {
		t.Error("a commit was made without an identity")
	}
}

// The real protection check treats a repository with no remote as
// unprotected and known: the root commit belongs on the current branch.
func TestDefaultRequiresPRNoRemote(t *testing.T) {
	dir := gitRepo(t)
	if required, known, def := defaultRequiresPR(dir); required || !known || def != "" {
		t.Errorf("no remote: required %v known %v default %q", required, known, def)
	}
}

// A detached HEAD would strand the scaffold on no branch: the files are
// written, the commit is left to the operator with the exact command.
func TestCommitScaffoldDetachedHead(t *testing.T) {
	stubProtection(t, false, true)
	dir := gitRepo(t)
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("one\n"), 0o644)
	git(dir, "add", "README.md")
	git(dir, "commit", "-q", "-m", "initial")
	sha, _ := git(dir, "rev-parse", "HEAD")
	git(dir, "switch", "-q", "--detach", sha)
	written, _, _ := scaffold(dir, "self", fakeContracts)
	var out bytes.Buffer
	commitScaffold(&out, dir, written)
	if !strings.Contains(out.String(), "HEAD is detached") || !strings.Contains(out.String(), "git commit -m \"chore: scaffold codecrew\" --") {
		t.Errorf("output:\n%s", out.String())
	}
	if h, _ := git(dir, "rev-parse", "HEAD"); h != sha {
		t.Error("a commit was made on a detached HEAD")
	}
	if _, err := os.Stat(filepath.Join(dir, ".codecrew.yml")); err != nil {
		t.Error("scaffold lost")
	}
}

// The pointer belongs at the repository root: init refuses a subdirectory
// before writing anything (checky's finding on PR #184 — inGitRepo saw no
// ./.git there and called it "not a repository").
func TestInitRefusesASubdirectory(t *testing.T) {
	dir := gitRepo(t)
	nested := filepath.Join(dir, "nested")
	os.MkdirAll(nested, 0o755)
	if root := repoRoot(nested); !sameDir(root, dir) {
		t.Fatalf("repoRoot(nested) = %q, want %q", root, dir)
	}
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir(nested)
	var out bytes.Buffer
	err := initCmd(&out, nil)
	if err == nil || !strings.Contains(err.Error(), "repository root") {
		t.Errorf("err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(nested, ".codecrew.yml")); err == nil {
		t.Error("init wrote into the subdirectory")
	}
}
