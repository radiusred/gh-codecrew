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

// defaultRequiresPR asks GitHub whether the working repository's default
// branch has a pull_request rule. known is false when the question could
// not be asked — no remote, no gh auth, offline — and the caller assumes
// the branch: a commit stranded on a protected default branch is the
// worse outcome (#172). A repository with no remote is not protected.
var defaultRequiresPR = func(dir string) (required, known bool) {
	if _, err := git(dir, "remote", "get-url", "origin"); err != nil {
		return false, true // nothing on GitHub yet: the root commit belongs here
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
		return false, false
	}
	var rules []struct {
		Type string `json:"type"`
	}
	if err := gh.JSON(&rules, "api", fmt.Sprintf("repos/%s/rules/branches/%s", repo.NameWithOwner, repo.DefaultBranchRef.Name)); err != nil {
		return false, false
	}
	for _, r := range rules {
		if r.Type == "pull_request" {
			return true, true
		}
	}
	return false, true
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
	required, known := defaultRequiresPR(dir)
	branch, _ := git(dir, "symbolic-ref", "--short", "HEAD")
	onBootstrap := false
	if required || !known {
		if _, err := git(dir, "rev-parse", "--verify", "--quiet", "HEAD"); err == nil {
			if _, err := git(dir, "switch", "-c", bootstrapBranch); err != nil {
				fmt.Fprintf(w, "note: could not create %s (%v); the scaffold is left uncommitted — commit it on a branch and open a PR\n", bootstrapBranch, err)
				return
			}
			onBootstrap = true
			branch = bootstrapBranch
		}
	}
	if _, err := git(dir, append([]string{"add", "--"}, written...)...); err != nil {
		fmt.Fprintf(w, "note: could not stage the scaffold (%v) — commit it by hand: git add %s && git commit -m %q\n", err, strings.Join(written, " "), scaffoldSubject)
		return
	}
	if _, err := git(dir, append([]string{"commit", "--only", "--quiet", "-m", scaffoldSubject, "--"}, written...)...); err != nil {
		fmt.Fprintf(w, "note: could not commit the scaffold (%v) — commit it by hand: git commit -m %q -- %s\n", err, scaffoldSubject, strings.Join(written, " "))
		return
	}
	sha, _ := git(dir, "rev-parse", "--short", "HEAD")
	if branch == "" {
		branch = "HEAD"
	}
	fmt.Fprintf(w, "committed %s on %s: %q — the scaffold only; your other changes are as they were\n", sha, branch, scaffoldSubject)
	switch {
	case onBootstrap && !known:
		fmt.Fprintf(w, "note: could not ask GitHub whether the default branch requires pull requests, so the commit is on %s — push it and open a PR (or merge it locally if the branch is unprotected)\n", bootstrapBranch)
	case onBootstrap:
		fmt.Fprintf(w, "the default branch requires pull requests: push %s and open the scaffold PR — the one merge you do yourself, since no task exists yet for task finish; delete-on-merge cleans the branch\n", bootstrapBranch)
	default:
		fmt.Fprintln(w, "push when ready; init never pushes")
	}
}
