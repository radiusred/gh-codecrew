package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/tracker"
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

// hubFake stands in for the tracker on the one call resolveRoles makes:
// the hub's pointer. It records the paths asked for, so a test can prove a
// hub never fetches its own file.
type hubFake struct {
	tracker.Tracker
	data  []byte
	err   error
	paths []string
}

func (f *hubFake) FileContent(repo, path string) ([]byte, error) {
	f.paths = append(f.paths, repo+":"+path)
	return f.data, f.err
}

// spokeCtx is a spoke as load() builds one: a pointer-only config, the hub
// somewhere else, and a tracker that answers the pointer fetch.
func spokeCtx(t *testing.T, fake *hubFake) *ctx {
	t.Helper()
	cfg, err := config.Parse([]byte("codecrew: \"2.0\"\nhub: acme/hub\n"))
	if err != nil {
		t.Fatal(err)
	}
	return &ctx{cfg: cfg, current: "acme/spoke", hub: "acme/hub", t: fake}
}

// M13-R5: a spoke whose hub cannot be read refuses rather than falling
// back to its own empty table. The three failure shapes are told apart —
// unreadable, protocol skew, unreachable — and "the hub declares no table"
// is not a failure at all.
func TestResolveRolesFailsClosed(t *testing.T) {
	for _, c := range []struct {
		name    string
		data    string
		err     error
		code    string
		details []string
	}{
		{
			name:    "the hub's pointer 404s (an unmigrated hub, mid-window)",
			err:     errors.New("gh api: gh: Not Found (HTTP 404)"),
			code:    "HUB_UNREADABLE",
			details: []string{"acme/hub", config.Pointer, "404", "gh codecrew migrate"},
		},
		{
			name:    "the spoke cannot see the hub",
			err:     errors.New("gh api: gh: Resource not accessible by integration (HTTP 403)"),
			code:    "HUB_UNREADABLE",
			details: []string{"acme/hub", "403"},
		},
		{
			name:    "the hub's pointer does not parse",
			data:    "codecrew: \"2.0\"\nhub: self\nroles: [this, is, a, list]\n",
			code:    "HUB_UNREADABLE",
			details: []string{"acme/hub", "does not parse"},
		},
		{
			name:    "the hub's pointer has no hub: field",
			data:    "codecrew: \"2.0\"\n",
			code:    "HUB_UNREADABLE",
			details: []string{"acme/hub", "hub"},
		},
		{
			name:    "the hub's routing table is a 1.0 one",
			data:    "codecrew: \"2.0\"\nhub: self\nroles:\n  qa: { identity: acme-qa }\n",
			code:    "HUB_UNREADABLE",
			details: []string{"acme/hub", "does not parse"},
		},
		{
			name:    "the hub speaks an older protocol major",
			data:    "codecrew: \"1.0\"\nhub: self\n",
			code:    "PROTOCOL_MISMATCH",
			details: []string{"acme/hub", "1.0", protocolVersion, "gh codecrew migrate"},
		},
		{
			name:    "the hub speaks a newer protocol major",
			data:    "codecrew: \"3.0\"\nhub: self\n",
			code:    "PROTOCOL_MISMATCH",
			details: []string{"acme/hub", "3.0", protocolVersion, "upgrade the extension"},
		},
		{
			name:    "GitHub is not reachable",
			err:     errors.New(`gh api: Get "https://api.github.com/...": dial tcp: lookup api.github.com: no such host`),
			code:    "GH_UNREACHABLE",
			details: []string{"GitHub could not be reached", "gh auth status"},
		},
		{
			name:    "the caller holds no credentials",
			err:     errors.New("gh api: To get started with GitHub CLI, please run:  gh auth login"),
			code:    "GH_UNREACHABLE",
			details: []string{"GitHub could not be reached"},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			ctx := spokeCtx(t, &hubFake{data: []byte(c.data), err: c.err})
			err := ctx.resolveRoles(io.Discard)
			var r refusal
			if !errors.As(err, &r) || r.Code != c.code {
				t.Fatalf("err = %v, want refused[%s]", err, c.code)
			}
			for _, want := range c.details {
				if !strings.Contains(r.Detail, want) {
					t.Errorf("detail %q does not name %q", r.Detail, want)
				}
			}
			// Nothing was adopted: a refused resolution leaves no table
			// behind for a later call to read.
			if ctx.roles != nil {
				t.Error("a refused resolution left a routing table on the ctx")
			}
		})
	}
}

// The other half of the same requirement: a hub that reads fine governs
// the spoke's roles, and a hub that declares no table is legitimately `~`
// everywhere. "Read fine, no table" must never be confused with "could not
// read" — before M13-R5 they were the same value.
func TestResolveRolesAdoptsTheHubsTable(t *testing.T) {
	fake := &hubFake{data: []byte(`
codecrew: "2.0"
hub: self
roles:
  qa: { identity: app:acme-qa }
  reviewer: { identity: user:alice }
`)}
	c := spokeCtx(t, fake)
	var notes bytes.Buffer
	if err := c.resolveRoles(&notes); err != nil {
		t.Fatalf("resolveRoles: %v", err)
	}
	if got := fake.paths; len(got) != 1 || got[0] != "acme/hub:"+config.Pointer {
		t.Errorf("fetched %v, want the hub's pointer once", got)
	}
	if !c.holdsRole("acme-qa[bot]", "qa") || c.holdsRole("mallory", "qa") {
		t.Error("the hub's qa row does not govern the spoke")
	}
	id, err := holder(c.rolesConfig().Roles, "reviewer")
	if err != nil || id.String() != "user:alice" {
		t.Errorf("holder(reviewer) = %v, %v, want user:alice", id, err)
	}
	if notes.Len() != 0 {
		t.Errorf("unexpected note: %q", notes.String())
	}

	// A hub that declares no table: no refusal, and every seat is the
	// operator's, exactly as a hub's own empty table has always been.
	empty := spokeCtx(t, &hubFake{data: []byte("codecrew: \"2.0\"\nhub: self\n")})
	if err := empty.resolveRoles(io.Discard); err != nil {
		t.Fatalf("a hub declaring no table refused: %v", err)
	}
	id, err = holder(empty.rolesConfig().Roles, "reviewer")
	if err != nil || !id.Operator() {
		t.Errorf("holder(reviewer) = %v, %v, want ~", id, err)
	}
	if !empty.holdsRole("mallory", "qa") {
		t.Error("an unrouted seat is not the operator's")
	}

	// A hub pointer with no codecrew: field is assumed current, with the
	// note a local pointer gets.
	unversioned := spokeCtx(t, &hubFake{data: []byte("hub: self\n")})
	notes.Reset()
	if err := unversioned.resolveRoles(&notes); err != nil {
		t.Fatalf("an unversioned hub pointer refused: %v", err)
	}
	if !strings.HasPrefix(notes.String(), "note:") || !strings.Contains(notes.String(), "acme/hub") {
		t.Errorf("note = %q, want one naming the hub", notes.String())
	}
}

// A hub reads its own pointer off disk: the local file *is* the table, so
// `codecrew` in a hub needs no network to resolve a role, and a hub with
// no table stays `~` everywhere with nothing fetched.
func TestHubResolvesItsOwnTableWithoutTheNetwork(t *testing.T) {
	cfg, err := config.Parse([]byte("codecrew: \"2.0\"\nhub: self\nroles:\n  qa: { identity: app:acme-qa }\n"))
	if err != nil {
		t.Fatal(err)
	}
	fake := &hubFake{err: errors.New("gh api: dial tcp: lookup api.github.com: no such host")}
	c := &ctx{cfg: cfg, current: "acme/hub", hub: "acme/hub", t: fake}
	if err := c.resolveRoles(io.Discard); err != nil {
		t.Fatalf("a hub refused with the network down: %v", err)
	}
	if len(fake.paths) != 0 {
		t.Errorf("a hub fetched %v — its own pointer is already read", fake.paths)
	}
	if !c.holdsRole("acme-qa[bot]", "qa") {
		t.Error("the hub's own table did not govern")
	}
}

// The regression this requirement exists for. An unreadable hub used to
// leave the spoke's own empty table in place, and an empty table resolves
// every seat to `~` — which turns task finish's holder-review gate into
// "any non-author approved" and milestone close's verdict count into
// "anyone commented" (the Claude scan on #254, finding 2). The second half
// of this test is the old behaviour, spelled out: it passes, which is
// precisely why the first half must refuse before any gate sees a table.
func TestAnUnreadableHubNeverReachesTheGates(t *testing.T) {
	c := spokeCtx(t, &hubFake{err: errors.New("gh api: gh: Not Found (HTTP 404)")})
	var r refusal
	if err := c.resolveRoles(io.Discard); !errors.As(err, &r) || r.Code != "HUB_UNREADABLE" {
		t.Fatalf("err = %v, want refused[HUB_UNREADABLE]", err)
	}

	// What the gates would have been handed had the error been swallowed
	// and the local table kept — the pre-2.0 shape.
	failOpen := &ctx{cfg: c.cfg, roles: c.cfg}
	id, err := holder(failOpen.rolesConfig().Roles, "reviewer")
	if err != nil || !id.Operator() {
		t.Fatalf("holder(reviewer) on the empty table = %v, %v; the fail-open shape has changed and this test needs rewriting", id, err)
	}
	// `~` is the branch task finish takes to "any non-author approval"...
	if !failOpen.holdsRole("anyone-at-all", "reviewer") {
		t.Error("the empty table did not open the review gate; this test no longer describes the regression")
	}
	// ...and the branch milestone close counts a QA verdict from any
	// commenter on.
	if !failOpen.holdsRole("passer-by", "qa") {
		t.Error("the empty table did not open the verdict gate; this test no longer describes the regression")
	}
}

// A spoke pointer that carries a routing table is refused where every
// other pointer condition is — at load, by the code config raised and the
// CLI named (M13-R5).
func TestLoadConfigRefusesASpokeRoutingTable(t *testing.T) {
	for _, c := range []struct {
		yml     string
		refused bool
	}{
		{"codecrew: \"2.0\"\nhub: acme/hub\nroles:\n  reviewer: { identity: user:alice }\n", true},
		{"codecrew: \"2.0\"\nhub: acme/hub\n", false},
		{"codecrew: \"2.0\"\nhub: self\nroles:\n  reviewer: { identity: user:alice }\n", false},
	} {
		dir := t.TempDir()
		writePointer(t, dir, c.yml)
		_, err := loadConfig(dir, &bytes.Buffer{})
		var r refusal
		got := errors.As(err, &r) && r.Code == "SPOKE_ROUTING"
		if got != c.refused {
			t.Errorf("%q: err = %v, want SPOKE_ROUTING %v", c.yml, err, c.refused)
		}
		if c.refused && !strings.Contains(r.Detail, "acme/hub") {
			t.Errorf("%q: detail does not name the hub that carries the table: %s", c.yml, r.Detail)
		}
	}
}
