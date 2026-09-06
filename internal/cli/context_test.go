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
// another major refuses PROTOCOL_MISMATCH; a missing field proceeds. The
// two mismatch directions read differently — ahead of the binary, upgrade
// the extension; behind it, migrate — and neither invites a hand edit of
// the version field.
func TestLoadConfigChecksProtocol(t *testing.T) {
	for _, c := range []struct {
		yml      string
		refused  bool
		wantNote bool
		detail   string
	}{
		{yml: "codecrew: \"2.0\"\nhub: self\n"},
		{yml: "codecrew: \"2.4\"\nhub: self\n"},                                                 // same major, later minor
		{yml: "hub: self\n", wantNote: true},                                                    // missing: assumed, noted
		{yml: "codecrew: \"3.0\"\nhub: self\n", refused: true, detail: "upgrade the extension"}, // ahead of us
		{yml: "codecrew: \"1.0\"\nhub: self\n", refused: true, detail: "gh codecrew migrate"},   // behind us
		{yml: "codecrew: \"0.1\"\nhub: self\n", refused: true, detail: "predates the protocol"}, // the pre-1.0 form, two majors back
	} {
		dir := t.TempDir()
		writePointer(t, dir, c.yml)
		var notes bytes.Buffer
		_, err := loadConfig(dir, &notes)
		if gotNote := strings.HasPrefix(notes.String(), "note:"); gotNote != c.wantNote {
			t.Errorf("%q: note = %q, want note %v", c.yml, notes.String(), c.wantNote)
		}
		var r refusal
		got := errors.As(err, &r) && r.Code == "PROTOCOL_MISMATCH"
		if got != c.refused || (err != nil && !c.refused) {
			t.Errorf("%q: err = %v, refused = %v, want refused %v", c.yml, err, got, c.refused)
		}
		if c.refused {
			if !strings.Contains(r.Detail, c.detail) {
				t.Errorf("%q: detail = %q, want it to name %q", c.yml, r.Detail, c.detail)
			}
			if strings.Contains(r.Detail, "the pointer") || strings.Contains(r.Detail, "edit") {
				t.Errorf("%q: detail invites a hand edit of the version field: %q", c.yml, r.Detail)
			}
		}
	}
}

// writePointer writes a 2.0 pointer into dir, .codecrew/ and all.
func writePointer(t *testing.T, dir, yml string) {
	t.Helper()
	path := filepath.Join(dir, filepath.FromSlash(config.Pointer))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(yml), 0o644); err != nil {
		t.Fatal(err)
	}
}

// A repo still on the 1.x layout is refused, not read: the detail names
// what was found, the protocol this binary speaks, and the verb that moves
// the repo forward (M13-R1; migrate lands in #256).
func TestLoadConfigRefusesTheLegacyLayout(t *testing.T) {
	for _, c := range []struct {
		name  string
		write func(dir string)
		names string
	}{
		{"a root pointer", func(dir string) {
			os.WriteFile(filepath.Join(dir, ".codecrew.yml"), []byte("codecrew: \"1.0\"\nhub: self\n"), 0o644)
		}, ".codecrew.yml"},
		{"a root roles/ with a contract and no pointer", func(dir string) {
			os.MkdirAll(filepath.Join(dir, "roles"), 0o755)
			os.WriteFile(filepath.Join(dir, "roles", "qa.md"), []byte("# Role: qa\n"), 0o644)
		}, "roles/qa.md"},
	} {
		dir := t.TempDir()
		c.write(dir)
		_, err := loadConfig(dir, &bytes.Buffer{})
		var r refusal
		if !errors.As(err, &r) || r.Code != "LAYOUT_LEGACY" {
			t.Fatalf("%s: err = %v, want refused[LAYOUT_LEGACY]", c.name, err)
		}
		for _, want := range []string{c.names, "protocol " + protocolVersion, "gh codecrew migrate"} {
			if !strings.Contains(r.Detail, want) {
				t.Errorf("%s: detail %q does not name %q", c.name, r.Detail, want)
			}
		}
	}
	// A project's own roles/ holding none of the five contracts is not the
	// 1.x layout — it is somebody else's directory, and the walk passes over it.
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "roles", "webserver", "tasks"), 0o755)
	os.WriteFile(filepath.Join(dir, "roles", "webserver", "tasks", "main.yml"), []byte("- name: x\n"), 0o644)
	_, err := loadConfig(dir, &bytes.Buffer{})
	var r refusal
	if errors.As(err, &r) {
		t.Errorf("an unrelated roles/ tree refused %v", err)
	}
	if err == nil || !strings.Contains(err.Error(), config.Pointer) {
		t.Errorf("err = %v, want the not-a-CodeCrew-repo error naming %s", err, config.Pointer)
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
		writePointer(t, dir, c.yml)
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
