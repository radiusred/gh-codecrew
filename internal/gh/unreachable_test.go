package gh

import (
	"errors"
	"fmt"
	"testing"
)

// The strings below were read off the installed gh, not invented: each was
// produced by running `gh api repos/.../contents/...` under the condition
// named. Unreachable must separate "GitHub never answered" from "GitHub
// answered and said no" — the second is the caller's own condition to
// name, and calling it offline would send an operator to check a network
// that is fine (M13-R5).
func TestUnreachable(t *testing.T) {
	for _, c := range []struct {
		name string
		msg  string
		want bool
	}{
		{"a refused connection", `gh api: Get "https://127.0.0.1/api/v3/repos/acme/hub/contents/x": dial tcp 127.0.0.1:443: connect: connection refused`, true},
		{"a host that does not resolve", "gh api: error connecting to nonexistent.invalid\ncheck your internet connection or https://githubstatus.com", true},
		{"no credentials at all", "gh api: To get started with GitHub CLI, please run:  gh auth login\nAlternatively, populate the GH_TOKEN environment variable with a GitHub API authentication token.", true},
		{"credentials GitHub rejects", "gh api: gh: Bad credentials (HTTP 401)", true},

		{"a missing file", "gh api: gh: Not Found (HTTP 404)", false},
		{"no access to the repo", "gh api: gh: Resource not accessible by integration (HTTP 403)", false},
		{"a rate limit", "gh api: gh: API rate limit exceeded (HTTP 403)", false},
		{"a body that does not parse", "gh api: unexpected output: invalid character 'x'", false},
	} {
		if got := Unreachable(errors.New(c.msg)); got != c.want {
			t.Errorf("%s: Unreachable = %v, want %v", c.name, got, c.want)
		}
	}
	if Unreachable(nil) {
		t.Error("Unreachable(nil) = true")
	}
	// It reads through a wrap, because the adapters wrap before the
	// classification happens.
	wrapped := fmt.Errorf("reading the hub's pointer: %w", errors.New("gh api: dial tcp: lookup api.github.com: no such host"))
	if !Unreachable(wrapped) {
		t.Error("a wrapped transport failure did not classify")
	}
}
