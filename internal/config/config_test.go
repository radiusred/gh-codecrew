package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Load walks upward for .codecrew/config.yml and reports Dir as the
// directory that *contains* .codecrew/ — the repo root — so every caller's
// cfg.Dir keeps meaning "root", whatever subdirectory the verb ran in.
func TestLoadWalksUpwardToTheRoot(t *testing.T) {
	root := t.TempDir()
	writePointer(t, root, "codecrew: \"2.0\"\nhub: self\n")
	deep := filepath.Join(root, "internal", "cli")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(deep)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Hub != "self" || cfg.Codecrew != "2.0" {
		t.Errorf("parsed %+v", cfg)
	}
	// t.TempDir can sit under a symlinked /tmp; compare resolved paths.
	wantDir, _ := filepath.EvalSymlinks(root)
	gotDir, _ := filepath.EvalSymlinks(cfg.Dir)
	if gotDir != wantDir {
		t.Errorf("Dir = %q, want the directory holding .codecrew/ (%q)", gotDir, wantDir)
	}
}

// The pointer wins wherever it sits above the caller. A 2.0 repo with an
// ordinary nested roles/ — a playbook's, a project's own — loads from
// inside that directory: the layout move exists to end that collision, and
// a walk that judged the 1.x layout level by level would have kept it
// (checky's finding on PR #277).
func TestLoadIgnoresANestedRolesUnderA20Pointer(t *testing.T) {
	root := t.TempDir()
	writePointer(t, root, "codecrew: \"2.0\"\nhub: self\n")
	nested := filepath.Join(root, "playbook", "roles")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "qa.md"), []byte("# Role: qa\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, from := range []string{filepath.Join(root, "playbook"), nested} {
		cfg, err := Load(from)
		if err != nil {
			t.Fatalf("Load(%s) = %v, want the 2.0 pointer above it", from, err)
		}
		wantDir, _ := filepath.EvalSymlinks(root)
		gotDir, _ := filepath.EvalSymlinks(cfg.Dir)
		if gotDir != wantDir {
			t.Errorf("Load(%s): Dir = %q, want %q", from, gotDir, wantDir)
		}
	}
}

// With no pointer anywhere, the 1.x layout is looked for at one level: the
// repo root when the caller is in a repository, the starting directory
// otherwise. So a 1.x repo refuses from any subdirectory, and a nested
// roles/ with nothing at the root is somebody else's directory.
func TestLoadLooksForTheLegacyLayoutAtTheRepoRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	deep := filepath.Join(root, "internal", "cli")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}
	// Nothing 1.x at the root yet: a nested roles/ is not the layout.
	nested := filepath.Join(deep, "roles")
	os.MkdirAll(nested, 0o755)
	os.WriteFile(filepath.Join(nested, "qa.md"), []byte("# Role: qa\n"), 0o644)
	_, err := Load(deep)
	var legacy *LegacyLayoutError
	if errors.As(err, &legacy) {
		t.Fatalf("a nested roles/ read as the repo's layout: %v", err)
	}
	if err == nil || !strings.Contains(err.Error(), Pointer) {
		t.Fatalf("err = %v, want the not-a-CodeCrew-repo error", err)
	}
	// Now the repo really is 1.x: refused from the subdirectory, named at
	// the root.
	os.WriteFile(filepath.Join(root, ".codecrew.yml"), []byte("codecrew: \"1.0\"\nhub: self\n"), 0o644)
	if _, err := Load(deep); !errors.As(err, &legacy) {
		t.Fatalf("Load(%s) = %v, want a *LegacyLayoutError", deep, err)
	}
	wantDir, _ := filepath.EvalSymlinks(root)
	gotDir, _ := filepath.EvalSymlinks(legacy.Dir)
	if gotDir != wantDir {
		t.Errorf("refusal names %q, want the repo root %q", gotDir, wantDir)
	}
}

// The 1.x layout is refused, never read: with no pointer anywhere above,
// a root .codecrew.yml or a root roles/ holding one of the five contracts
// is the repo's layout (M13-R1).
func TestLoadRefusesTheLegacyLayout(t *testing.T) {
	for _, c := range []struct {
		name  string
		write func(dir string)
		found string
	}{
		{"the 1.x pointer", func(dir string) {
			os.WriteFile(filepath.Join(dir, ".codecrew.yml"), []byte("codecrew: \"1.0\"\nhub: self\n"), 0o644)
		}, ".codecrew.yml"},
		{"a 1.x roles/ with no pointer", func(dir string) {
			os.MkdirAll(filepath.Join(dir, "roles"), 0o755)
			os.WriteFile(filepath.Join(dir, "roles", "qa.md"), []byte("# Role: qa\n"), 0o644)
		}, "roles/qa.md"},
	} {
		dir := t.TempDir()
		c.write(dir)
		_, err := Load(dir)
		var legacy *LegacyLayoutError
		if !errors.As(err, &legacy) {
			t.Fatalf("%s: err = %v, want a *LegacyLayoutError", c.name, err)
		}
		if len(legacy.Found) != 1 || legacy.Found[0] != c.found {
			t.Errorf("%s: found %v, want [%s]", c.name, legacy.Found, c.found)
		}
		if !strings.Contains(legacy.Error(), c.found) {
			t.Errorf("%s: message %q does not name the file", c.name, legacy.Error())
		}
	}
	// A 2.0 pointer wins at its own level even with 1.x leftovers beside it:
	// the walk refuses only where there is nothing current to read.
	dir := t.TempDir()
	writePointer(t, dir, "codecrew: \"2.0\"\nhub: self\n")
	os.WriteFile(filepath.Join(dir, ".codecrew.yml"), []byte("codecrew: \"1.0\"\nhub: self\n"), 0o644)
	if _, err := Load(dir); err != nil {
		t.Errorf("a 2.0 pointer beside a 1.x leftover: %v", err)
	}
}

// Neither layout: the plain not-a-CodeCrew-repo error, naming the file the
// caller is missing. A roles/ that holds none of the contracts — Ansible's,
// say — is not the 1.x layout.
func TestLoadWithoutEitherLayout(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "roles", "webserver"), 0o755)
	os.WriteFile(filepath.Join(dir, "roles", "webserver", "main.yml"), []byte("- name: x\n"), 0o644)
	_, err := Load(dir)
	var legacy *LegacyLayoutError
	if errors.As(err, &legacy) {
		t.Fatalf("an unrelated roles/ tree read as the 1.x layout: %v", err)
	}
	if err == nil || !strings.Contains(err.Error(), Pointer) {
		t.Errorf("err = %v, want it to name %s", err, Pointer)
	}
}

func writePointer(t *testing.T, dir, yml string) {
	t.Helper()
	path := filepath.Join(dir, filepath.FromSlash(Pointer))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(yml), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestParse(t *testing.T) {
	cfg, err := Parse([]byte(`
codecrew: "0.1"
hub: self
roles:
  implementer:
    harness: claude-code
    model: claude-fable-5
    identity: app:radiusred-cody
  reviewer:
    harness: codex
    model: gpt-5.6-sol
    identity: ~
`))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Hub != "self" {
		t.Errorf("Hub = %q, want self", cfg.Hub)
	}
	if got := (cfg.Roles["implementer"].Identity); got.Kind != KindApp || got.Value != "radiusred-cody" {
		t.Errorf("implementer identity = %v", got)
	}
	if got := cfg.Roles["reviewer"].Identity; !got.Operator() {
		t.Errorf("nil identity should parse as the operator, got %v", got)
	}
	if got := cfg.Roles["reviewer"].Model; got != "gpt-5.6-sol" {
		t.Errorf("model on a codex row should load, got %q", got)
	}
}

func TestParseMissingHub(t *testing.T) {
	if _, err := Parse([]byte(`codecrew: "0.1"`)); err == nil {
		t.Error("expected error for missing hub")
	}
}

func TestRoleFor(t *testing.T) {
	cfg := &Config{Roles: map[string]Role{
		"implementer": {Identity: ParseIdentity("app:radiusred-cody")},
		"qa":          {Identity: ParseIdentity("app:radiusred-testy")},
		"reviewer":    {}, // identity: ~ — human operator during bootstrap
	}}
	cases := []struct {
		login, want string
	}{
		{"radiusred-cody[bot]", "implementer"},
		{"radiusred-testy[bot]", "qa"},
		{"radiusred-testy", "qa"},
		{"davison", ""},
		{"", ""}, // empty login must not match the reviewer's empty identity
	}
	for _, c := range cases {
		if got := cfg.RoleFor(c.login); got != c.want {
			t.Errorf("RoleFor(%q) = %q, want %q", c.login, got, c.want)
		}
	}
	if got := (&Config{}).RoleFor("radiusred-cody[bot]"); got != "" {
		t.Errorf("no roles configured should resolve to \"\", got %q", got)
	}
}

func TestHoldsRole(t *testing.T) {
	routed := &Config{Roles: map[string]Role{
		"implementer": {Identity: ParseIdentity("app:radiusred-cody")},
		"qa":          {Identity: ParseIdentity("app:radiusred-testy")},
	}}
	human := &Config{Roles: map[string]Role{
		"implementer": {Identity: ParseIdentity("app:radiusred-cody")},
		"qa":          {Identity: ParseIdentity("user:alice")},
	}}
	team := &Config{Roles: map[string]Role{
		"implementer": {Identity: ParseIdentity("app:radiusred-cody")},
		"qa":          {Identity: ParseIdentity("team:myorg/qa-crew")},
	}}
	unrouted := &Config{Roles: map[string]Role{
		"implementer": {Identity: ParseIdentity("app:radiusred-cody")},
		"qa":          {}, // identity: ~ — held by the human operator
	}}
	empty := &Config{}
	cases := []struct {
		name  string
		cfg   *Config
		login string
		want  bool
	}{
		{"routed: its bot holds", routed, "radiusred-testy[bot]", true},
		{"routed: bare slug holds", routed, "radiusred-testy", true},
		{"routed: a human does not", routed, "davison", false},
		{"routed: another crew bot does not", routed, "radiusred-cody[bot]", false},
		{"routed to a human: that human holds", human, "alice", true},
		{"routed to a human: another human does not", human, "bob", false},
		{"routed to a human: a bot of the same name does not", human, "alice[bot]", false},
		{"routed to a team: membership is not this table's question", team, "alice", false},
		{"unrouted: the operator holds", unrouted, "davison", true},
		{"unrouted: a crew bot does not", unrouted, "radiusred-cody[bot]", false},
		{"unrouted: an unrelated bot does not", unrouted, "somebot[bot]", false},
		{"unrouted: an identity routed elsewhere does not", unrouted, "radiusred-cody", false},
		{"no roles at all: the operator holds", empty, "davison", true},
		{"empty login never holds", unrouted, "", false},
	}
	for _, c := range cases {
		if got := c.cfg.HoldsRole(c.login, "qa"); got != c.want {
			t.Errorf("%s: HoldsRole(%q) = %v, want %v", c.name, c.login, got, c.want)
		}
	}
}

func TestHubRepo(t *testing.T) {
	self := &Config{Hub: "self"}
	if got := self.HubRepo("radiusred/spoke"); got != "radiusred/spoke" {
		t.Errorf("self hub = %q, want current repo", got)
	}
	remote := &Config{Hub: "radiusred/hub"}
	if got := remote.HubRepo("radiusred/spoke"); got != "radiusred/hub" {
		t.Errorf("named hub = %q, want radiusred/hub", got)
	}
}

func TestCompatible(t *testing.T) {
	cases := []struct {
		pointer  string
		wantNote bool
		wantErr  bool
	}{
		{"1.0", false, false},
		{"1.4", false, false}, // same major, later minor
		{"0.1", true, false},  // the pre-1.0 form of 1.0
		{"", true, false},     // missing: assumed, noted
		{"2.0", false, true},  // another major
		{"0.2", false, true},  // not the frozen form
	}
	for _, c := range cases {
		note, err := Compatible(c.pointer, "1.0")
		if (err != nil) != c.wantErr || (note != "") != c.wantNote {
			t.Errorf("Compatible(%q, 1.0) = note %q, err %v", c.pointer, note, err)
		}
	}
}

// The 2.0 grammar: every form the routing table's identity value may take,
// and everything else refused rather than guessed at (SPEC §5, M13-R4).
func TestParseIdentity(t *testing.T) {
	cases := []struct {
		in    string
		kind  IdentityKind
		value string
		str   string
	}{
		{"~", KindOperator, "", "~"},
		{"", KindOperator, "", "~"},
		{"  ", KindOperator, "", "~"},
		{"app:myorg-coder", KindApp, "myorg-coder", "app:myorg-coder"},
		{"user:alice", KindUser, "alice", "user:alice"},
		{"team:myorg/review-crew", KindTeam, "myorg/review-crew", "team:myorg/review-crew"},
		{" app:myorg-coder ", KindApp, "myorg-coder", "app:myorg-coder"},
		// Untyped: the 1.0 forms, and malformed typed ones.
		{"myorg-coder", KindUntyped, "myorg-coder", "myorg-coder"},
		{"myorg/review-crew", KindUntyped, "myorg/review-crew", "myorg/review-crew"},
		{"app:", KindUntyped, "app:", "app:"},
		{"user:", KindUntyped, "user:", "user:"},
		{"app:myorg/coder", KindUntyped, "app:myorg/coder", "app:myorg/coder"},
		{"team:review-crew", KindUntyped, "team:review-crew", "team:review-crew"},
		{"team:myorg/", KindUntyped, "team:myorg/", "team:myorg/"},
		{"team:/review-crew", KindUntyped, "team:/review-crew", "team:/review-crew"},
		{"team:a/b/c", KindUntyped, "team:a/b/c", "team:a/b/c"},
		{"bot:myorg-coder", KindUntyped, "bot:myorg-coder", "bot:myorg-coder"},
	}
	for _, c := range cases {
		got := ParseIdentity(c.in)
		if got.Kind != c.kind || got.Value != c.value || got.String() != c.str {
			t.Errorf("ParseIdentity(%q) = %+v (%q), want kind %v value %q string %q",
				c.in, got, got.String(), c.kind, c.value, c.str)
		}
	}
}

// Login is the handle a review request can name: user: and team: yield one
// (a team in the org/slug form `gh pr create --reviewer` wants), an App and
// the operator yield nothing — the implementer contract's whole branch.
func TestIdentityLogin(t *testing.T) {
	cases := map[string]string{
		"~":                      "",
		"app:myorg-reviewy":      "",
		"user:alice":             "alice",
		"team:myorg/review-crew": "myorg/review-crew",
	}
	for in, want := range cases {
		if got := ParseIdentity(in).Login(); got != want {
			t.Errorf("ParseIdentity(%q).Login() = %q, want %q", in, got, want)
		}
	}
}

func TestIdentityTeam(t *testing.T) {
	org, team, ok := ParseIdentity("team:myorg/review-crew").Team()
	if org != "myorg" || team != "review-crew" || !ok {
		t.Errorf("Team() = %q, %q, %v", org, team, ok)
	}
	for _, other := range []string{"~", "app:myorg-coder", "user:alice"} {
		if _, _, ok := ParseIdentity(other).Team(); ok {
			t.Errorf("%q reported as a team", other)
		}
	}
}

// A 1.0 table read by a 2.0 binary: the untyped row is named, and the
// refusal is deterministic — role names in order, so the same row is
// always the one reported.
func TestParseRefusesUntypedIdentity(t *testing.T) {
	_, err := Parse([]byte(`
hub: self
roles:
  reviewer: { identity: myorg-reviewy }
  implementer: { identity: myorg-coder }
  qa: { identity: ~ }
`))
	var untyped *UntypedIdentityError
	if !errors.As(err, &untyped) {
		t.Fatalf("err = %v, want *UntypedIdentityError", err)
	}
	if untyped.Role != "implementer" || untyped.Value != "myorg-coder" {
		t.Errorf("named row = %q / %q, want implementer / myorg-coder", untyped.Role, untyped.Value)
	}
	for _, want := range []string{"roles.implementer.identity", "app:<slug>", "user:<login>", "team:<org>/<slug>", "`~`"} {
		if !strings.Contains(untyped.Error(), want) {
			t.Errorf("detail does not name %q: %s", want, untyped.Error())
		}
	}
}

// A fully typed table parses, and a routing row that names no identity at
// all is the operator — not an untyped value.
func TestParseTypedTable(t *testing.T) {
	cfg, err := Parse([]byte(`
codecrew: "2.0"
hub: self
roles:
  implementer: { harness: claude-code, identity: app:myorg-coder }
  reviewer: { identity: "team:myorg/review-crew" }
  qa: { identity: user:alice }
  doc-synthesizer: { harness: claude-code }
  coordinator: { identity: ~ }
`))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"implementer":     "app:myorg-coder",
		"reviewer":        "team:myorg/review-crew",
		"qa":              "user:alice",
		"doc-synthesizer": "~",
		"coordinator":     "~",
	}
	for role, w := range want {
		if got := cfg.Roles[role].Identity.String(); got != w {
			t.Errorf("%s identity = %q, want %q", role, got, w)
		}
	}
}

// An identity that is not a scalar at all is a parse failure, not a
// silently empty row.
func TestParseRefusesNonScalarIdentity(t *testing.T) {
	if _, err := Parse([]byte("hub: self\nroles:\n  qa: { identity: [a, b] }\n")); err == nil {
		t.Error("a sequence identity parsed")
	}
}

// The mismatch is not symmetric. A pointer ahead of the binary is the
// operator's extension being old; one behind it is the repo being old, and
// only the second names the migration. Neither ever tells anyone to edit
// the version field, which would make the file lie about the repo.
func TestCompatibleMismatchIsAsymmetric(t *testing.T) {
	_, err := Compatible("3.0", "2.0")
	if err == nil {
		t.Fatal("a newer pointer was accepted")
	}
	if !strings.Contains(err.Error(), "upgrade the extension") {
		t.Errorf("newer pointer: %v", err)
	}
	if strings.Contains(err.Error(), "migrate") {
		t.Errorf("newer pointer offered a migration: %v", err)
	}

	_, err = Compatible("1.0", "2.0")
	if err == nil {
		t.Fatal("an older pointer was accepted")
	}
	for _, want := range []string{"predates", "gh codecrew migrate"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("older pointer: %v, want it to name %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "upgrade the extension") {
		t.Errorf("older pointer told the operator to upgrade: %v", err)
	}

	// An unparseable major is not read as older: nothing suggests a migration.
	if _, err := Compatible("weird", "2.0"); err == nil || strings.Contains(err.Error(), "migrate") {
		t.Errorf("unparseable pointer: %v", err)
	}
}

// A spoke's pointer carries no routing table: the hub carries one table
// for the whole project, and a copy in a spoke either outranks it or goes
// stale (the Claude scan on #254, finding 8). Refused at parse, so no verb
// ever reads one (M13-R5).
func TestParseRefusesASpokeRoutingTable(t *testing.T) {
	for _, c := range []struct {
		name    string
		yml     string
		refused bool
	}{
		{"a spoke carrying a table", "hub: acme/hub\nroles:\n  reviewer: { identity: user:alice }\n", true},
		{"a spoke carrying only the pointer", "hub: acme/hub\n", false},
		{"a spoke with an empty table", "hub: acme/hub\nroles: {}\n", false},
		{"the hub's own table", "hub: self\nroles:\n  reviewer: { identity: user:alice }\n", false},
	} {
		_, err := Parse([]byte(c.yml))
		var spoke *SpokeRoutingError
		if got := errors.As(err, &spoke); got != c.refused {
			t.Errorf("%s: err = %v, want SpokeRoutingError %v", c.name, err, c.refused)
		}
		if !c.refused {
			continue
		}
		// The detail must name the hub that does carry the table, the row
		// found, and the file — an agent acts on it without reading code.
		for _, want := range []string{"acme/hub", "reviewer", Pointer} {
			if !strings.Contains(spoke.Error(), want) {
				t.Errorf("%s: detail %q does not name %q", c.name, spoke, want)
			}
		}
	}
	// Several rows are named in role order, so the refusal reads the same
	// on every run.
	_, err := Parse([]byte("hub: acme/hub\nroles:\n  reviewer: { identity: user:alice }\n  qa: { identity: ~ }\n"))
	var spoke *SpokeRoutingError
	if !errors.As(err, &spoke) {
		t.Fatalf("err = %v, want SpokeRoutingError", err)
	}
	if got := strings.Join(spoke.Roles, ","); got != "qa,reviewer" {
		t.Errorf("Roles = %q, want them sorted", got)
	}
}

// The hub's protocol major is checked from the spoke that reads its table,
// which is what makes the version check topology-wide (M13-R5). The two
// directions read as Compatible's do, and both sides are named.
func TestCompatibleHub(t *testing.T) {
	for _, c := range []struct {
		hubVersion string
		wantNote   bool
		wantErr    string
	}{
		{hubVersion: "2.0"},
		{hubVersion: "2.7"},                                   // same major, later minor
		{hubVersion: "", wantNote: true},                      // absent: assumed, noted
		{hubVersion: "1.0", wantErr: "gh codecrew migrate"},   // the hub is behind
		{hubVersion: "3.0", wantErr: "upgrade the extension"}, // the hub is ahead
		{hubVersion: "0.1", wantErr: "gh codecrew migrate"},   // two majors back
	} {
		note, err := CompatibleHub("acme/hub", c.hubVersion, "2.0")
		if (note != "") != c.wantNote {
			t.Errorf("%q: note = %q, want note %v", c.hubVersion, note, c.wantNote)
		}
		if (err != nil) != (c.wantErr != "") {
			t.Fatalf("%q: err = %v, want error %v", c.hubVersion, err, c.wantErr != "")
		}
		if err == nil {
			continue
		}
		// Both sides, so the reader knows which end to move.
		for _, want := range []string{"acme/hub", c.hubVersion, "2.0", c.wantErr} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("%q: %v does not name %q", c.hubVersion, err, want)
			}
		}
	}
}
