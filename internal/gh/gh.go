// Package gh is a thin exec wrapper over the GitHub CLI. Wrapping gh (rather
// than speaking REST directly) is a founding decision: auth, base URLs, and
// enterprise quirks come for free (docs/founding-decisions.md).
package gh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Command builds the gh process every call runs through. It is a variable
// so tests can stand a fake gh behind the adapters that parse its output
// and errors (the pattern cli's ghVersion set): production never
// reassigns it.
var Command = exec.Command

// Run executes gh with args and returns stdout.
func Run(args ...string) ([]byte, error) {
	cmd := Command("gh", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := bytes.TrimSpace(stderr.Bytes())
		if len(msg) == 0 {
			return nil, fmt.Errorf("gh %s: %w", args[0], err)
		}
		return nil, fmt.Errorf("gh %s: %s", args[0], msg)
	}
	return stdout.Bytes(), nil
}

// JSON executes gh and unmarshals its stdout into v.
func JSON(v any, args ...string) error {
	out, err := Run(args...)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(out, v); err != nil {
		return fmt.Errorf("gh %s: unexpected output: %w", args[0], err)
	}
	return nil
}

// JSONLoose is JSON for gh commands that exit nonzero while still printing a
// valid result (e.g. `gh pr checks` with failing checks). The exit error is
// ignored when stdout unmarshals cleanly.
func JSONLoose(v any, args ...string) error {
	cmd := Command("gh", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()
	if jsonErr := json.Unmarshal(stdout.Bytes(), v); jsonErr == nil {
		return nil
	}
	if runErr != nil {
		return fmt.Errorf("gh %s: %s", args[0], bytes.TrimSpace(stderr.Bytes()))
	}
	return fmt.Errorf("gh %s: unexpected output", args[0])
}

// CurrentRepo returns the owner/repo of the repository in the working
// directory, per gh's own resolution of the origin remote.
func CurrentRepo() (string, error) {
	var repo struct {
		NameWithOwner string `json:"nameWithOwner"`
	}
	if err := JSON(&repo, "repo", "view", "--json", "nameWithOwner"); err != nil {
		return "", err
	}
	return repo.NameWithOwner, nil
}

// unreachableMarkers are what gh prints when the API was never answered —
// the transport failed, or the caller carries no credentials GitHub will
// look at. Each was read off the installed gh (2.x) rather than guessed:
//
//	dial tcp 127.0.0.1:443: connect: connection refused   (no route, refused, timed out)
//	error connecting to nonexistent.invalid               (the host does not resolve)
//	check your internet connection or https://githubstatus.com
//	To get started with GitHub CLI, please run:  gh auth login
//	Alternatively, populate the GH_TOKEN environment variable …
//	gh: Bad credentials (HTTP 401)
//
// The list is deliberately narrow. An HTTP 403 or 404 means GitHub
// answered — the caller lacks access, or the path is absent — and those
// are the caller's own conditions to name, not this one.
var unreachableMarkers = []string{
	"dial tcp",
	"error connecting to",
	"check your internet connection",
	"gh auth login",
	"GH_TOKEN environment variable",
	"Bad credentials",
	"(HTTP 401)",
}

// Unreachable reports whether err is gh failing to reach GitHub at all —
// offline, DNS-less, or unauthenticated — as opposed to GitHub answering
// with a refusal of its own. Callers that fetch across repos use it to
// name the condition instead of reporting a missing file (SPEC §6).
func Unreachable(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, m := range unreachableMarkers {
		if strings.Contains(msg, m) {
			return true
		}
	}
	return false
}
