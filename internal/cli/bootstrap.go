package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/gh"
)

// bootstrapBranch is where init commits when the default branch requires
// pull requests: that PR is the one merge the operator does by hand, the
// pre-milestone gate (#164 finding 52), and delete-on-merge cleans it.
const bootstrapBranch = "codecrew-bootstrap"

const scaffoldSubject = "chore: scaffold codecrew"

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
	var repo struct {
		NameWithOwner    string `json:"nameWithOwner"`
		DefaultBranchRef struct {
			Name string `json:"name"`
		} `json:"defaultBranchRef"`
	}
	cmd := exec.Command("gh", "repo", "view", "--json", "nameWithOwner,defaultBranchRef")
	cmd.Dir = dir
	data, err := cmd.Output()
	if err != nil || json.Unmarshal(data, &repo) != nil || repo.NameWithOwner == "" || repo.DefaultBranchRef.Name == "" {
		return false, false, ""
	}
	var rules []struct {
		Type string `json:"type"`
	}
	if err := gh.JSON(&rules, "api", fmt.Sprintf("repos/%s/rules/branches/%s", repo.NameWithOwner, repo.DefaultBranchRef.Name)); err != nil {
		return false, false, repo.DefaultBranchRef.Name
	}
	for _, r := range rules {
		if r.Type == "pull_request" {
			return true, true, repo.DefaultBranchRef.Name
		}
	}
	return false, true, repo.DefaultBranchRef.Name
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
// cannot be made is a note with the command to run by hand.
func commitScaffold(w io.Writer, dir string, written []string) {
	if len(written) == 0 {
		return
	}
	byHand := fmt.Sprintf("git add %s && git commit -m %q -- %s", strings.Join(written, " "), scaffoldSubject, strings.Join(written, " "))
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
	if _, err := git(dir, append([]string{"add", "--"}, written...)...); err != nil {
		fmt.Fprintf(w, "note: could not stage the scaffold (%v) — commit it by hand: %s\n", err, byHand)
		return
	}
	if _, err := git(dir, append([]string{"commit", "--only", "--quiet", "-m", scaffoldSubject, "--"}, written...)...); err != nil {
		fmt.Fprintf(w, "note: could not commit the scaffold (%v) — commit it by hand: %s\n", err, byHand)
		return
	}
	sha, _ := git(dir, "rev-parse", "--short", "HEAD")
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
