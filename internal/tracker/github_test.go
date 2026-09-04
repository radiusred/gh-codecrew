package tracker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/gh"
)

// fakeGH stands a gh behind gh.Command for one test: `pr view` answers
// with a minimal open PR, `pr checks` fails with checksErr on stderr the
// way gh reports a GraphQL refusal. The fake is this test binary
// re-entered at TestHelperGH — no script, no PATH shim.
func fakeGH(t *testing.T, checksErr string) {
	t.Helper()
	orig := gh.Command
	gh.Command = func(_ string, args ...string) *exec.Cmd {
		cmd := exec.Command(os.Args[0], append([]string{"-test.run=^TestHelperGH$", "--"}, args...)...)
		cmd.Env = append(os.Environ(), "GH_HELPER=1", "GH_HELPER_CHECKS_STDERR="+checksErr)
		return cmd
	}
	t.Cleanup(func() { gh.Command = orig })
}

// TestHelperGH is the fake gh's body; it is inert unless fakeGH re-enters
// the binary with GH_HELPER set.
func TestHelperGH(t *testing.T) {
	if os.Getenv("GH_HELPER") != "1" {
		return
	}
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	if len(args) >= 2 && args[0] == "pr" {
		switch args[1] {
		case "view":
			fmt.Print(`{"state":"OPEN","reviewDecision":"","headRefName":"task/1-x","headRefOid":"abc","isCrossRepository":false,"mergedAt":"","author":{"login":"app/myorg-coder"},"reviews":[]}`)
			os.Exit(0)
		case "checks":
			fmt.Fprint(os.Stderr, os.Getenv("GH_HELPER_CHECKS_STDERR"))
			os.Exit(1)
		}
	}
	fmt.Fprintf(os.Stderr, "fake gh: unexpected call %v", args)
	os.Exit(2)
}

// PRInfo carries the permission an installation token lacks as
// PR.ChecksUnreadable instead of the raw GraphQL error — the two verbatim
// shapes #198 recorded — and returns every other failure raw, the same
// message on an unrelated path included. The dry-run test injects the
// field through a fake tracker; this is the wiring that sets it.
func TestPRInfoMapsMissingChecksPermission(t *testing.T) {
	cases := []struct {
		stderr, want string
	}{
		{"GraphQL: Resource not accessible by integration (node.statusCheckRollup.nodes.0.commit.statusCheckRollup)", "checks: read"},
		{"GraphQL: Resource not accessible by integration (node.statusCheckRollup.nodes.0.commit.statusCheckRollup.contexts.nodes.0.checkSuite.workflowRun)", "actions: read"},
	}
	for _, tc := range cases {
		fakeGH(t, tc.stderr)
		pr, err := GitHub{}.PRInfo("o/r", 9)
		if err != nil {
			t.Fatalf("%q: PRInfo returned the raw error: %v", tc.want, err)
		}
		if pr.ChecksUnreadable != tc.want {
			t.Errorf("ChecksUnreadable = %q, want %q", pr.ChecksUnreadable, tc.want)
		}
		if pr.NoChecks || pr.ChecksOK || pr.ChecksPending {
			t.Errorf("%q: check state must stay unknown, got %+v", tc.want, pr)
		}
		if pr.Number != 9 || pr.Author != "myorg-coder" || !pr.Open {
			t.Errorf("%q: the PR's own state must survive the mapping: %+v", tc.want, pr)
		}
	}

	raw := "GraphQL: Resource not accessible by integration (repository.collaborators)"
	fakeGH(t, raw)
	pr, err := GitHub{}.PRInfo("o/r", 9)
	if err == nil || !strings.Contains(err.Error(), raw) {
		t.Fatalf("an unrelated refusal must surface raw, got err %v pr %+v", err, pr)
	}
	if pr.ChecksUnreadable != "" {
		t.Errorf("unrelated refusal classified as %q", pr.ChecksUnreadable)
	}

	fakeGH(t, "no checks reported on the 'task/1-x' branch")
	pr, err = GitHub{}.PRInfo("o/r", 9)
	if err != nil || !pr.NoChecks || pr.ChecksUnreadable != "" {
		t.Errorf("the checkless shape still maps to NoChecks: err %v pr %+v", err, pr)
	}
}
