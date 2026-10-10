package cli

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// bootstrapBranch is where init commits when the default branch requires
// pull requests: that PR is the one merge the operator does by hand, the
// pre-milestone gate (#164 finding 52), and delete-on-merge cleans it.
const bootstrapBranch = "codecrew-bootstrap"

const scaffoldSubject = "chore: scaffold codecrew"

// bodyLineMax is the longest line the commit lint accepts in a body; every
// body a verb writes is held to it, so the commit it makes never fails
// the check the project's PRs run.
const bodyLineMax = 100

// bodyFilesMax is how many files a housekeeping body lists one per line
// before it summarises them as a count: past it the list stops being
// something a reviewer reads, and the commit's own file list says it.
const bodyFilesMax = 20

// housekeepingTrailer closes every body a verb writes on a commit of its
// own: the line a reviewer looks for to know the light path was meant.
const housekeepingTrailer = "Housekeeping (SPEC §4); no task."

// housekeepingBody is the commit body SPEC §4 asks of a housekeeping
// change, in the one shape init, migrate and roles sync all write (#389):
// the tool that defined the target — the invocation as the verb read it,
// the release and the protocol the binary implements — then what it
// wrote, one line each (a count past bodyFilesMax), then the trailer.
// done is the participle that opens it ("Scaffolded", "Migrated",
// "Written"); invocation is what follows "gh codecrew". The tool line
// wraps when it would pass bodyLineMax (wrapLine), so no owner/repo GitHub
// accepts — up to 39 and 100 characters — makes a line the lint refuses.
func housekeepingBody(done, invocation string, files []string) string {
	tool := fmt.Sprintf("%s by gh codecrew %s (gh-codecrew %s, protocol %s):", done, invocation, version, protocolVersion)
	lines := wrapLine(tool, bodyLineMax)
	if len(files) > bodyFilesMax {
		lines = append(lines, fmt.Sprintf("%d files; git show --stat lists them", len(files)))
	} else {
		lines = append(lines, files...)
	}
	lines = append(lines, housekeepingTrailer)
	return strings.Join(lines, "\n")
}

// wrapLine breaks s into lines of at most max characters, at spaces. A
// word that cannot fit on a line of its own — a maximum-length owner/repo
// is 140 characters — is split rather than left to break the bound: just
// after a slash where one falls within the room, so owner/ ends one line
// and the repository name opens the next, and only a word with no usable
// slash is cut at the limit. A split continues on the next line with no
// space, so the word reads back whole by joining the two lines; the hub's
// own name is in .codecrew/config.yml either way.
func wrapLine(s string, max int) []string {
	var lines []string
	var line []rune
	for _, field := range strings.Fields(s) {
		word := []rune(field)
		switch {
		case len(line) == 0 && len(word) <= max:
			line = word
			continue
		case len(line) > 0 && len(line)+1+len(word) <= max:
			line = append(append(line, ' '), word...)
			continue
		case len(word) <= max:
			lines = append(lines, string(line))
			line = word
			continue
		}
		// Too long for any line: put as much as fits after the current
		// line's words, cut after a slash if one fits, then carry on.
		if len(line) > 0 {
			if cut := slashCut(word, max-len(line)-1); cut > 0 {
				lines = append(lines, string(line)+" "+string(word[:cut]))
				word = word[cut:]
			} else {
				lines = append(lines, string(line))
			}
			line = nil
		}
		for len(word) > max {
			cut := slashCut(word, max)
			if cut == 0 {
				cut = max
			}
			lines = append(lines, string(word[:cut]))
			word = word[cut:]
		}
		line = word
	}
	if len(line) > 0 {
		lines = append(lines, string(line))
	}
	return lines
}

// slashCut is where to split word so its first part, ending in a slash, is
// at most room characters: just after the last slash that fits, or 0 when
// none does (a slash at the very end leaves nothing to carry, so it is no
// cut).
func slashCut(word []rune, room int) int {
	for i := min(room, len(word)-1) - 1; i >= 0; i-- {
		if word[i] == '/' {
			return i + 1
		}
	}
	return 0
}

// scaffoldBody is the body of init's commit: the invocation it ran as —
// --hub only in spoke mode — and the files it wrote.
func scaffoldBody(hub string, written []string) string {
	invocation := "init"
	if hub != "self" {
		invocation += " --hub " + hub
	}
	return housekeepingBody("Scaffolded", invocation, written)
}

func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(errb.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("git %s: %s", args[0], msg)
	}
	return strings.TrimSpace(out.String()), nil
}

// repoRoot is the repository's top level, or "" outside one. inGitRepo
// looked only for ./.git and so mistook a subdirectory for no repository
// (checky's finding on PR #184); the pointer belongs at the root, where
// config.Load's upward walk finds it first.
func repoRoot(dir string) string {
	root, err := git(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return ""
	}
	return root
}

// defaultRequiresPR asks GitHub whether the working repository's default
// branch has a pull_request rule, and which branch that is. known is
// false when the question could not be asked — no gh auth, offline — and
// the caller assumes the branch: a commit stranded on a protected default
// branch is the worse outcome (#172). A repository with no remote is not
// protected: nothing is on GitHub yet.
var defaultRequiresPR = func(dir string) (required, known bool, defaultBranch string) {
	if _, err := git(dir, "remote", "get-url", "origin"); err != nil {
		return false, true, ""
	}
	v := tracker.GitHub{}
	repo, defaultBranch, err := v.RepoInDir(dir)
	if err != nil {
		return false, false, ""
	}
	types, err := v.BranchRuleTypes(repo, defaultBranch)
	if err != nil {
		return false, false, defaultBranch
	}
	for _, t := range types {
		if t == "pull_request" {
			return true, true, defaultBranch
		}
	}
	return false, true, defaultBranch
}

// bootstrapBase is where codecrew-bootstrap is cut from: the default
// branch as the remote has it, else the local one, else HEAD — so the
// scaffold PR carries the scaffold and nothing a feature branch had
// (checky's finding on PR #184).
func bootstrapBase(dir, defaultBranch string) (ref string, fromDefault bool) {
	if defaultBranch != "" {
		for _, r := range []string{"refs/remotes/origin/" + defaultBranch, "refs/heads/" + defaultBranch} {
			if _, err := git(dir, "rev-parse", "--verify", "--quiet", r); err == nil {
				return r, true
			}
		}
	}
	return "HEAD", false
}

// commitScaffold commits exactly the files init wrote — a pathspec commit,
// so the operator's staged and unstaged work stays as it was (no stash,
// no -A) — on the current branch, or on codecrew-bootstrap when the
// default branch requires pull requests. It never pushes. A commit that
// cannot be made is a note with the command to run by hand. hub is init's
// --hub, which the commit body names (SPEC §4).
func commitScaffold(w io.Writer, dir, hub string, written []string) {
	if len(written) == 0 {
		return
	}
	body := scaffoldBody(hub, written)
	byHand := commitByHand(scaffoldSubject, body, written, written)
	branch, err := git(dir, "symbolic-ref", "--short", "HEAD")
	born := true
	if _, e := git(dir, "rev-parse", "--verify", "--quiet", "HEAD"); e != nil {
		born = false // an unborn branch: the scaffold will be the root commit
	}
	if err != nil && born {
		// Detached HEAD: a commit here would strand the scaffold on no
		// branch; the files are written, the commit is the operator's.
		fmt.Fprintf(w, "note: HEAD is detached — the scaffold is written but not committed; switch to a branch and run: %s\n", byHand)
		return
	}
	required, known, defaultBranch := defaultRequiresPR(dir)
	onBootstrap, fromDefault := false, false
	if (required || !known) && born {
		base, fd := bootstrapBase(dir, defaultBranch)
		if _, err := git(dir, "switch", "-c", bootstrapBranch, base); err != nil {
			fmt.Fprintf(w, "note: could not create %s from %s (%v) — the scaffold is written but not committed; on a branch cut from the default branch run: %s\n", bootstrapBranch, base, err, byHand)
			return
		}
		onBootstrap, fromDefault = true, fd
		branch = bootstrapBranch
	}
	sha, ok := pathspecCommit(w, dir, scaffoldSubject, body, "the scaffold", written, written)
	if !ok {
		return
	}
	fmt.Fprintf(w, "committed %s on %s: %q — the scaffold only; your other changes are as they were\n", sha, branch, scaffoldSubject)
	switch {
	case onBootstrap && !known:
		fmt.Fprintf(w, "note: could not ask GitHub whether the default branch requires pull requests, so the commit is on %s — git push -u origin %s, then open a PR (or merge it locally if the branch is unprotected)\n", bootstrapBranch, bootstrapBranch)
	case onBootstrap && !fromDefault:
		fmt.Fprintf(w, "the default branch requires pull requests: %s was cut from HEAD because the default branch was not found locally — check it carries only the scaffold, then git push -u origin %s and open the scaffold PR\n", bootstrapBranch, bootstrapBranch)
	case onBootstrap:
		fmt.Fprintf(w, "the default branch requires pull requests: git push -u origin %s and open the scaffold PR — the one merge you do yourself, since no task exists yet for task finish; delete-on-merge cleans the branch\n", bootstrapBranch)
	default:
		fmt.Fprintf(w, "git push -u origin %s when ready; init never pushes\n", branch)
	}
}

// commitByHand is the command that makes the same commit the verb would
// have — printed in every note that says the commit could not be made, so
// the operator is never left with written files and no way to land them,
// nor with a commit that lacks the body SPEC §4 asks for (#389). stage and
// paths are pathspecCommit's: what needs adding, and what the commit is
// limited to. They differ only for migrate, whose git mv has already
// staged each source's removal — and git add refuses a path that is gone.
func commitByHand(subject, body string, stage, paths []string) string {
	return fmt.Sprintf("git add %s && git commit %s -- %s", strings.Join(stage, " "), commitMessageArgs(subject, body), strings.Join(paths, " "))
}

// commitMessageArgs is the -m arguments of a printed git commit: the
// subject, then the body as a second -m — git makes each -m a paragraph,
// so the body's lines stay one paragraph under it. The body is
// single-quoted, newlines and all, which sh, bash, zsh and fish all take
// as one argument when the command is pasted.
func commitMessageArgs(subject, body string) string {
	args := fmt.Sprintf("-m %q", subject)
	if body != "" {
		args += " -m " + shellQuote(body)
	}
	return args
}

// shellQuote single-quotes s for a POSIX shell or fish, closing and
// reopening the quotes around any single quote it holds.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// pathspecCommit stages and commits exactly paths on the current branch —
// `git add --` then `git commit --only --`, never `-A` and never a stash —
// so every other staged and unstaged change is left as it was. It is the
// machinery the verbs that commit go through: init's scaffold, migrate's
// layout move and roles sync's contracts. body, when not empty, is the
// commit message's body — each of the three names the tool that defined
// the target there (SPEC §4, housekeeping; housekeepingBody). noun names what is being committed in the notes a
// failure prints; stage is what needs adding (a scaffold's whole pathspec,
// or only the destinations of a set of renames git already knows), and
// paths is the pathspec the commit is limited to.
func pathspecCommit(w io.Writer, dir, subject, body, noun string, stage, paths []string) (sha string, ok bool) {
	byHand := commitByHand(subject, body, stage, paths)
	if len(stage) > 0 {
		if _, err := git(dir, append([]string{"add", "--"}, stage...)...); err != nil {
			fmt.Fprintf(w, "note: could not stage %s (%v) — commit it by hand: %s\n", noun, err, byHand)
			return "", false
		}
	}
	args := []string{"commit", "--only", "--quiet", "-m", subject}
	if body != "" {
		args = append(args, "-m", body)
	}
	if _, err := git(dir, append(append(args, "--"), paths...)...); err != nil {
		fmt.Fprintf(w, "note: could not commit %s (%v) — commit it by hand: %s\n", noun, err, byHand)
		return "", false
	}
	sha, _ = git(dir, "rev-parse", "--short", "HEAD")
	return sha, true
}
