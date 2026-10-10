package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/radiusred/gh-codecrew/internal/config"
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

// stubVersion stands in for the release a build stamps, so a test can see
// the body name the release rather than whatever the test binary carries.
func stubVersion(t *testing.T, v string) {
	t.Helper()
	prev := version
	version = v
	t.Cleanup(func() { version = prev })
}

// toolLine is the line a housekeeping body opens with, as the test expects
// to read it once any wrapping is undone.
func toolLine(done, invocation string) string {
	return fmt.Sprintf("%s by gh codecrew %s (gh-codecrew %s, protocol %s):", done, invocation, version, protocolVersion)
}

// headMessage is the whole message of HEAD: subject, blank line, body.
func headMessage(t *testing.T, dir string) string {
	t.Helper()
	out, err := git(dir, "log", "-1", "--format=%B")
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// assertHousekeepingMessage holds a commit message to the shape SPEC §4
// asks of a housekeeping commit and the lint accepts (#389): the subject,
// a blank line, the tool line — rejoined if it wrapped, with a space, or
// with none after a line wrapLine split at a slash — naming the verb,
// the release and the protocol, then exactly the files, then the trailer;
// and no line, subject included, past bodyLineMax.
func assertHousekeepingMessage(t *testing.T, msg, subject, tool string, files []string) {
	t.Helper()
	lines := strings.Split(strings.TrimRight(msg, "\n"), "\n")
	for _, l := range lines {
		if n := utf8.RuneCountInString(l); n > bodyLineMax {
			t.Errorf("a line of %d characters, past the lint's %d: %q", n, bodyLineMax, l)
		}
	}
	if len(lines) < 4 || lines[0] != subject || lines[1] != "" {
		t.Fatalf("want the subject %q, a blank line and a body:\n%s", subject, msg)
	}
	body := lines[2:]
	i := 0
	for i < len(body) && !strings.HasSuffix(body[i], "):") {
		i++
	}
	if i == len(body) {
		t.Fatalf("no tool line ending \"):\":\n%s", msg)
	}
	got := body[0]
	for _, l := range body[1 : i+1] {
		if strings.HasSuffix(got, "/") {
			got += l
		} else {
			got += " " + l
		}
	}
	if got != tool {
		t.Errorf("tool line = %q, want %q", got, tool)
	}
	want := append(append([]string{}, files...), housekeepingTrailer)
	if got := body[i+1:]; !slices.Equal(got, want) {
		t.Errorf("after the tool line:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// runPrinted runs a command a verb printed for the operator, as they would
// paste it, in dir — with an identity, since the operator has one.
func runPrinted(t *testing.T, dir, command string) {
	t.Helper()
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com",
		"GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=commit.gpgsign", "GIT_CONFIG_VALUE_0=false")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("the printed command failed: %v\n%s\n--- command:\n%s", err, out, command)
	}
}

// A fresh repository: the scaffold is the root commit, holding exactly
// the files init wrote, on the branch the operator was on — and its body
// names init, the release and the protocol, and every file (SPEC §4, #389).
func TestCommitScaffoldRootCommit(t *testing.T) {
	stubProtection(t, false, true)
	stubVersion(t, "v2.1.1")
	dir := gitRepo(t)
	written, _, err := scaffold(dir, "self", fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	commitScaffold(&out, dir, "self", written)
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
	assertHousekeepingMessage(t, headMessage(t, dir), scaffoldSubject,
		"Scaffolded by gh codecrew init (gh-codecrew v2.1.1, protocol "+protocolVersion+"):", written)
	if st, _ := git(dir, "status", "--porcelain"); st != "" {
		t.Errorf("tree not clean after the commit: %q", st)
	}
	// A rerun writes nothing and commits nothing.
	out.Reset()
	written, _, _ = scaffold(dir, "self", fakeContracts)
	commitScaffold(&out, dir, "self", written)
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
	commitScaffold(&out, dir, "self", written)
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
		commitScaffold(&out, dir, "self", written)
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
	commitScaffold(&out2, dir2, "self", written2)
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
	commitScaffold(&out, dir, "self", written)
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
	commitScaffold(&out, dir, "self", written)
	if !strings.Contains(out.String(), "note: could not commit the scaffold") || !strings.Contains(out.String(), "commit it by hand") {
		t.Errorf("output:\n%s", out.String())
	}
	if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(config.Pointer))); err != nil {
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
	commitScaffold(&out, dir, "self", written)
	if !strings.Contains(out.String(), "HEAD is detached") || !strings.Contains(out.String(), "git commit -m \"chore: scaffold codecrew\" -m 'Scaffolded by gh codecrew init (") {
		t.Errorf("output:\n%s", out.String())
	}
	if h, _ := git(dir, "rev-parse", "HEAD"); h != sha {
		t.Error("a commit was made on a detached HEAD")
	}
	if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(config.Pointer))); err != nil {
		t.Error("scaffold lost")
	}
	// The command it prints makes the commit init would have: run as
	// pasted on a branch, it lands the scaffold with the body (#389).
	_, command, ok := strings.Cut(strings.TrimRight(out.String(), "\n"), "switch to a branch and run: ")
	if !ok {
		t.Fatalf("no command to run:\n%s", out.String())
	}
	git(dir, "switch", "-q", "-c", "scaffold")
	runPrinted(t, dir, command)
	assertHousekeepingMessage(t, headMessage(t, dir), scaffoldSubject, toolLine("Scaffolded", "init"), written)
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
	if _, err := os.Stat(filepath.Join(nested, filepath.FromSlash(config.Pointer))); err == nil {
		t.Error("init wrote into the subdirectory")
	}
}

// The scaffold commits on whatever branch the operator is on, and the
// push instruction names it — nothing assumes main.
func TestCommitScaffoldOnABranchNotNamedMain(t *testing.T) {
	stubProtection(t, false, true)
	dir := t.TempDir()
	for _, args := range [][]string{{"init", "-q", "-b", "trunk"}, {"config", "user.email", "t@example.com"}, {"config", "user.name", "t"}, {"config", "commit.gpgsign", "false"}} {
		if _, err := git(dir, args...); err != nil {
			t.Fatal(err)
		}
	}
	written, _, _ := scaffold(dir, "self", fakeContracts)
	var out bytes.Buffer
	commitScaffold(&out, dir, "self", written)
	if b, _ := git(dir, "symbolic-ref", "--short", "HEAD"); b != "trunk" {
		t.Errorf("committed on %s", b)
	}
	if !strings.Contains(out.String(), " on trunk: ") || !strings.Contains(out.String(), "git push -u origin trunk when ready") {
		t.Errorf("output:\n%s", out.String())
	}
}

// A spoke's scaffold commit names the hub it was pointed at, as the
// operator typed it: the run a reviewer would repeat (#389).
func TestCommitScaffoldBodyNamesTheSpokesHub(t *testing.T) {
	stubProtection(t, false, true)
	dir := gitRepo(t)
	written, _, err := scaffold(dir, "myorg/platform", fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	commitScaffold(&out, dir, "myorg/platform", written)
	assertHousekeepingMessage(t, headMessage(t, dir), scaffoldSubject, toolLine("Scaffolded", "init --hub myorg/platform"), written)
}

// The body holds to the lint whatever it is given: a hub name long enough
// to push the tool line past 100 characters wraps it at a space, and a
// file list past bodyFilesMax becomes a count.
func TestHousekeepingBodyWrapsAndSummarises(t *testing.T) {
	stubVersion(t, "v2.1.1")
	hub := "a-rather-long-organisation/an-even-longer-repository-name"
	body := scaffoldBody(hub, []string{".codecrew/config.yml"})
	lines := strings.Split(body, "\n")
	if len(lines) != 4 {
		t.Fatalf("want the tool line wrapped onto two lines, the file and the trailer:\n%s", body)
	}
	assertHousekeepingMessage(t, scaffoldSubject+"\n\n"+body, scaffoldSubject, toolLine("Scaffolded", "init --hub "+hub), []string{".codecrew/config.yml"})

	var many []string
	for i := range bodyFilesMax + 1 {
		many = append(many, fmt.Sprintf("file%d.md", i))
	}
	assertHousekeepingMessage(t, "chore: x\n\n"+housekeepingBody("Written", "roles sync", many), "chore: x",
		toolLine("Written", "roles sync"), []string{fmt.Sprintf("%d files; git show --stat lists them", bodyFilesMax+1)})
	assertHousekeepingMessage(t, "chore: x\n\n"+housekeepingBody("Written", "roles sync", many[:bodyFilesMax]), "chore: x",
		toolLine("Written", "roles sync"), many[:bodyFilesMax])
}

// The printed command survives a paste: a body with a single quote in it
// is still one argument, and comes back byte for byte.
func TestCommitMessageArgsQuoteTheBody(t *testing.T) {
	dir := gitRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	body := "Written by gh codecrew roles sync (it's quoted):\nf.txt\n" + housekeepingTrailer
	runPrinted(t, dir, commitByHand("chore: quoted", body, []string{"f.txt"}, []string{"f.txt"}))
	if got := headMessage(t, dir); got != "chore: quoted\n\n"+body {
		t.Errorf("message = %q", got)
	}
}

// The longest owner/repo GitHub accepts — a 39-character owner and a
// 100-character repository, 140 in all — cannot sit on one line of the
// lint's 100, so the body splits it after the slash: owner/ ends the line
// init opens, the repository name fills the next, and nothing is lost
// (checky's finding on PR #404). Committed, the message holds to the bound.
func TestScaffoldBodyHoldsTheBoundForTheLongestHub(t *testing.T) {
	stubVersion(t, "v2.1.1")
	owner, repo := strings.Repeat("o", 39), strings.Repeat("r", 100)
	hub := owner + "/" + repo
	body := scaffoldBody(hub, []string{".codecrew/config.yml"})
	want := []string{
		"Scaffolded by gh codecrew init --hub " + owner + "/",
		repo,
		"(gh-codecrew v2.1.1, protocol " + protocolVersion + "):",
		".codecrew/config.yml",
		housekeepingTrailer,
	}
	if got := strings.Split(body, "\n"); !slices.Equal(got, want) {
		t.Errorf("body:\n%s\nwant:\n%s", body, strings.Join(want, "\n"))
	}

	stubProtection(t, false, true)
	dir := gitRepo(t)
	written, _, err := scaffold(dir, hub, fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	commitScaffold(&out, dir, hub, written)
	assertHousekeepingMessage(t, headMessage(t, dir), scaffoldSubject, toolLine("Scaffolded", "init --hub "+hub), written)
}

// The bound holds for any word: one with no slash to split after is cut at
// the limit, and the pieces read back as the word.
func TestWrapLineCutsAWordWithNoSlash(t *testing.T) {
	word := strings.Repeat("x", 250)
	lines := wrapLine("lead "+word+" tail", bodyLineMax)
	for _, l := range lines {
		if n := utf8.RuneCountInString(l); n > bodyLineMax {
			t.Errorf("a line of %d characters: %q", n, l)
		}
	}
	if got := strings.Join(lines, ""); got != "lead"+word+" tail" {
		t.Errorf("lines %q do not read back as the input", lines)
	}
}
