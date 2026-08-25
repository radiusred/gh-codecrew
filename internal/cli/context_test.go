package cli

import (
	"fmt"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/config"
)

func teamCtx(t *testing.T) *ctx {
	t.Helper()
	cfg, err := config.Parse([]byte(`
codecrew: "0.1"
hub: self
roles:
  implementer: { identity: myorg-coder }
  reviewer: { identity: myorg/review-crew }
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

func TestTeamIdentityParsing(t *testing.T) {
	for _, tc := range []struct {
		in      string
		org, tm string
		ok      bool
	}{
		{"myorg/review-crew", "myorg", "review-crew", true},
		{"alice", "", "", false},
		{"myorg-coder", "", "", false},
		{"/team", "", "", false},
		{"org/", "", "", false},
		{"", "", "", false},
	} {
		org, tm, ok := config.TeamIdentity(tc.in)
		if org != tc.org || tm != tc.tm || ok != tc.ok {
			t.Errorf("TeamIdentity(%q) = %q,%q,%v", tc.in, org, tm, ok)
		}
	}
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
