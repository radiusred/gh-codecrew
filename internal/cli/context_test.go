package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/config"
)

func teamCtx(t *testing.T) *ctx {
	t.Helper()
	cfg, err := config.Parse([]byte(`
codecrew: "0.1"
hub: self
roles:
  implementer: { identity: app:myorg-coder }
  reviewer: { identity: "team:myorg/review-crew" }
  qa: { identity: ~ }
`))
	if err != nil {
		t.Fatal(err)
	}
	return &ctx{cfg: cfg, roles: cfg}
}

func stubTeams(t *testing.T, members map[string]bool) *int {
	t.Helper()
	calls := 0
	orig := teamMembers
	teamMembers = func(org, team string) (map[string]bool, error) {
		calls++
		if org != "myorg" || team != "review-crew" {
			return nil, fmt.Errorf("unexpected team %s/%s", org, team)
		}
		return members, nil
	}
	t.Cleanup(func() { teamMembers = orig })
	return &calls
}

func TestTeamHeldRole(t *testing.T) {
	c := teamCtx(t)
	calls := stubTeams(t, map[string]bool{"alice": true, "bob": true})

	// Any member holds the role (#44); non-members and bots do not.
	if !c.holdsRole("alice", "reviewer") || !c.holdsRole("bob", "reviewer") {
		t.Error("team members do not hold the team-held role")
	}
	if c.holdsRole("mallory", "reviewer") {
		t.Error("non-member holds the team-held role")
	}
	if c.holdsRole("alice[bot]", "reviewer") {
		t.Error("bot login matched team membership")
	}
	// Membership makes a login crew: roleFor resolves through the team.
	if got := c.roleFor("alice"); got != "reviewer" {
		t.Errorf("roleFor(alice) = %q, want reviewer", got)
	}
	if got := c.roleFor("mallory"); got != "" {
		t.Errorf("roleFor(mallory) = %q, want none", got)
	}
	// Pure matching untouched for App identities.
	if !c.holdsRole("myorg-coder[bot]", "implementer") {
		t.Error("App identity matching broke")
	}
	// A team member must NOT hold unrouted roles: the operator fallback
	// is crew-free through the team-aware lens (checky's PR #101 finding —
	// a reviewer-team member's verdict must not count as qa).
	if c.holdsRole("alice", "qa") {
		t.Error("team member holds the unrouted qa role")
	}
	if !c.holdsRole("mallory", "qa") {
		t.Error("uncrewed human lost the unrouted-role fallback")
	}
	// Memoized: many checks, one fetch.
	if *calls != 1 {
		t.Errorf("team fetched %d times, want 1 (memoized per run)", *calls)
	}
}

func TestUnreadableTeamFailsClosed(t *testing.T) {
	c := teamCtx(t)
	orig := teamMembers
	teamMembers = func(org, team string) (map[string]bool, error) { return nil, fmt.Errorf("403") }
	t.Cleanup(func() { teamMembers = orig })
	if c.holdsRole("alice", "reviewer") {
		t.Error("unreadable team granted the role — must fail closed")
	}
}

// The pointer's protocol major gates every verb that loads it (SPEC §5):
// another major refuses PROTOCOL_MISMATCH; 0.1 and a missing field proceed.
func TestLoadConfigChecksProtocol(t *testing.T) {
	for _, c := range []struct {
		yml     string
		refused bool
	}{
		{"codecrew: \"1.0\"\nhub: self\n", false},
		{"codecrew: \"0.1\"\nhub: self\n", false},
		{"hub: self\n", false},
		{"codecrew: \"2.0\"\nhub: self\n", true},
	} {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, ".codecrew.yml"), []byte(c.yml), 0o644)
		var notes bytes.Buffer
		_, err := loadConfig(dir, &notes)
		wantNote := !strings.Contains(c.yml, `"1.0"`) && !c.refused
		if gotNote := strings.HasPrefix(notes.String(), "note:"); gotNote != wantNote {
			t.Errorf("%q: note = %q, want note %v", c.yml, notes.String(), wantNote)
		}
		var r refusal
		got := errors.As(err, &r) && r.Code == "PROTOCOL_MISMATCH"
		if got != c.refused || (err != nil && !c.refused) {
			t.Errorf("%q: err = %v, refused = %v, want refused %v", c.yml, err, got, c.refused)
		}
	}
}

// A 1.0 routing table read by a 2.0 binary refuses at load, naming the row
// and the four forms — the same split as PROTOCOL_MISMATCH: config detects
// it, the CLI gives it its code (M13-R4).
func TestLoadConfigRefusesUntypedIdentity(t *testing.T) {
	for _, c := range []struct {
		yml     string
		refused bool
	}{
		{"hub: self\nroles:\n  reviewer: { identity: myorg-reviewy }\n", true},
		{"hub: self\nroles:\n  reviewer: { identity: myorg/review-crew }\n", true},
		{"hub: self\nroles:\n  reviewer: { identity: app:myorg-reviewy }\n", false},
		{"hub: self\nroles:\n  reviewer: { identity: user:alice }\n", false},
		{"hub: self\nroles:\n  reviewer: { identity: \"team:myorg/review-crew\" }\n", false},
		{"hub: self\nroles:\n  reviewer: { identity: ~ }\n", false},
		{"hub: self\n", false}, // a pointer-only spoke declares no table
	} {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, ".codecrew.yml"), []byte(c.yml), 0o644); err != nil {
			t.Fatal(err)
		}
		var notes bytes.Buffer
		_, err := loadConfig(dir, &notes)
		var r refusal
		got := errors.As(err, &r) && r.Code == "IDENTITY_UNTYPED"
		if got != c.refused {
			t.Errorf("%q: err = %v, want IDENTITY_UNTYPED %v", c.yml, err, c.refused)
		}
		if c.refused && !strings.Contains(r.Detail, "roles.reviewer.identity") {
			t.Errorf("%q: detail does not name the row: %s", c.yml, r.Detail)
		}
	}
}
