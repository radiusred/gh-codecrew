package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	codecrew "github.com/radiusred/gh-codecrew"
	"github.com/radiusred/gh-codecrew/internal/config"
	"testing/fstest"
)

var fakeContracts = fstest.MapFS{
	config.RolesDir + "/implementer.md": {Data: []byte("# Role: implementer\n")},
	config.RolesDir + "/qa.md":          {Data: []byte("# Role: qa\n")},
}

// rolesPath names a file in the scaffolded contracts directory the way the
// scaffold's own keys do — host separators, under .codecrew/roles/.
func rolesPath(name string) string {
	return filepath.Join(filepath.FromSlash(config.RolesDir), name)
}

func TestScaffoldHub(t *testing.T) {
	dir := t.TempDir()
	written, skipped, err := scaffold(dir, "self", fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	if len(skipped) != 0 {
		t.Errorf("fresh dir skipped %v", skipped)
	}
	for _, want := range []string{filepath.FromSlash(config.Pointer), filepath.FromSlash(config.AgentsFile), "ROADMAP.md", "AGENTS.md", "CLAUDE.md", rolesPath("qa.md"), rolesPath("qa" + localSuffix), rolesPath("implementer" + localSuffix)} {
		if !slices.Contains(written, want) {
			t.Errorf("missing %s from written %v", want, written)
		}
		if _, err := os.Stat(filepath.Join(dir, want)); err != nil {
			t.Errorf("%s not on disk: %v", want, err)
		}
	}
	cfg, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(config.Pointer)))
	if err != nil {
		t.Fatal(err)
	}
	for _, role := range []string{"implementer", "reviewer", "qa", "doc-synthesizer", "coordinator"} {
		if !strings.Contains(string(cfg), role) {
			t.Errorf("hub config missing role %q", role)
		}
	}
}

// A spoke gets the pointer and the entry point and nothing else: an agent is
// dispatched into a spoke exactly as into a hub, so it needs the same file to
// land on, while the roadmap, the contracts and their extensions stay in the
// hub that owns them (M13-R3).
func TestScaffoldSpoke(t *testing.T) {
	dir := t.TempDir()
	written, _, err := scaffold(dir, "org/hub", fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{filepath.FromSlash(config.AgentsFile), filepath.FromSlash(config.Pointer), "AGENTS.md", "CLAUDE.md"}
	slices.Sort(want)
	if !slices.Equal(written, want) {
		t.Fatalf("spoke mode wrote %v, want %v", written, want)
	}
	cfg, _ := os.ReadFile(filepath.Join(dir, filepath.FromSlash(config.Pointer)))
	if !strings.Contains(string(cfg), "hub: org/hub") {
		t.Errorf("pointer = %q", cfg)
	}
}

// repoRoot finds the repository from anywhere inside it and reports
// nothing outside one — a real repository, not a stray .git entry.
func TestRepoRoot(t *testing.T) {
	if root := repoRoot(t.TempDir()); root != "" {
		t.Errorf("bare temp dir read as a repo: %q", root)
	}
	dir := gitRepo(t)
	if root := repoRoot(dir); !sameDir(root, dir) {
		t.Errorf("repoRoot = %q, want %q", root, dir)
	}
}

func TestScaffoldIdempotent(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := scaffold(dir, "self", fakeContracts); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(dir, "ROADMAP.md")
	if err := os.WriteFile(marker, []byte("edited"), 0o644); err != nil {
		t.Fatal(err)
	}
	written, skipped, err := scaffold(dir, "self", fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 0 {
		t.Errorf("re-run wrote %v", written)
	}
	if len(skipped) == 0 {
		t.Error("re-run reported nothing skipped")
	}
	if data, _ := os.ReadFile(marker); string(data) != "edited" {
		t.Error("re-run clobbered an existing file")
	}
}

// TestScaffoldedAgentsCarriesDispatchAuthorization: harness guardrails
// defer to AGENTS.md, so the scaffold must state CodeCrew's dispatch
// expectation explicitly — conditionally, so a platform-dispatched role
// agent reads a prohibition, not a licence (finding 9 on #73).
func TestScaffoldedAgentsCarriesDispatchAuthorization(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := scaffold(dir, "self", fakeContracts); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(config.AgentsFile)))
	if err != nil {
		t.Fatal(err)
	}
	flat := strings.Join(strings.Fields(string(data)), " ")
	for _, want := range []string{
		"Dispatch authorization",
		"operator's primary session",
		"authorized and expected",
		"never dispatches another role",
		"Reviews are model reviews",
	} {
		if !strings.Contains(flat, want) {
			t.Errorf("scaffolded %s missing %q", config.AgentsFile, want)
		}
	}
}

// Both scaffolds carry the protocol version this binary implements — a
// hub pointer its own binary would refuse must be impossible to scaffold.
func TestScaffoldsCarryProtocolVersion(t *testing.T) {
	for _, hub := range []string{"self", "o/hub"} {
		dir := t.TempDir()
		if _, _, err := scaffold(dir, hub, fakeContracts); err != nil {
			t.Fatal(err)
		}
		data, _ := os.ReadFile(filepath.Join(dir, filepath.FromSlash(config.Pointer)))
		if !strings.HasPrefix(string(data), "codecrew: \""+protocolVersion+"\"") {
			t.Errorf("hub=%s: pointer starts %q, want codecrew: %q", hub, strings.SplitN(string(data), "\n", 2)[0], protocolVersion)
		}
	}
}

// The scaffolded routing table teaches the typed identity grammar (SPEC
// §5): a hub told to write values its own binary refuses with
// IDENTITY_UNTYPED must be impossible to scaffold, and the guidance must
// not drift from the grammar again (checky's finding on PR #276).
func TestScaffoldedRoutingTeachesTypedIdentities(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := scaffold(dir, "self", fakeContracts); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(config.Pointer)))
	if err != nil {
		t.Fatal(err)
	}
	pointer := string(data)
	for _, want := range []string{"app:<slug>", "user:<login>", "team:<org>/<slug>", "IDENTITY_UNTYPED"} {
		if !strings.Contains(pointer, want) {
			t.Errorf("scaffolded routing guidance does not name %q:\n%s", want, pointer)
		}
	}
	// The 1.0 grammar it replaced, in any of its wordings.
	for _, gone := range []string{"org/team-slug", "a GitHub App slug or a username"} {
		if strings.Contains(pointer, gone) {
			t.Errorf("scaffolded routing guidance still teaches the 1.0 grammar (%q)", gone)
		}
	}
	// And what it scaffolds parses, with every seat the operator's.
	cfg, err := config.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	for role, r := range cfg.Roles {
		if !r.Identity.Operator() {
			t.Errorf("scaffolded %s is %v, want ~", role, r.Identity)
		}
	}
	if len(cfg.Roles) != 5 {
		t.Errorf("scaffolded %d roles, want 5", len(cfg.Roles))
	}
}

// Nothing init writes may point at a file the adopter's repo does not
// have: the protocol and its docs live upstream (Codex pre-launch scan,
// #131).
func TestScaffoldReferencesResolveUpstream(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := scaffold(dir, "self", codecrew.Roles); err != nil {
		t.Fatal(err)
	}
	var files []string
	if err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, p)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		data, _ := os.ReadFile(f)
		for _, bad := range []string{"the hub's docs/", "hub's `SPEC.md`", "hub's SPEC.md", "`scripts/codecrew-token"} {
			if strings.Contains(string(data), bad) {
				t.Errorf("%s contains hub-relative reference %q", strings.TrimPrefix(f, dir), bad)
			}
		}
	}
}

func TestHelpIsNotAnError(t *testing.T) {
	for _, args := range [][]string{
		{"task", "finish", "--help"}, {"task", "new", "--help"}, {"init", "--help"},
		{"roles", "show", "x", "--help"}, {"milestone", "close", "--help"}, {"checkpoint", "-h"},
		{"task", "new", "--milestone", "999", "--title", "--help"}, // help wins before the verb runs
	} {
		if err := Run(args); err != nil {
			t.Errorf("%v: %v", args, err)
		}
	}
	// A genuine bad argument is still a failure — help must not swallow it.
	for _, args := range [][]string{{"milestone", "close", "notanumber"}, {"bogusverb"}, {"roles", "show"}} {
		if err := Run(args); err == nil {
			t.Errorf("%v: accepted", args)
		}
	}
}

// Claude Code loads CLAUDE.md and never AGENTS.md, so the hub scaffold must
// write a CLAUDE.md whose first line imports the shared entry point — else
// "gh codecrew init && claude" starts blind (#141).
func TestScaffoldedClaudeImportsAgents(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := scaffold(dir, "self", fakeContracts); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	first, _, _ := strings.Cut(string(data), "\n")
	if first != "@AGENTS.md" {
		t.Errorf("CLAUDE.md first line = %q, want %q", first, "@AGENTS.md")
	}
	if strings.Contains(string(data), "`@AGENTS.md`") {
		t.Error("the import must not be wrapped in backticks — Claude Code treats that as literal text")
	}
}

// What stays hub-only: the roadmap and the contracts. A spoke that grew its
// own copy of either would be a second source of truth for what the hub owns.
func TestScaffoldSpokeWritesNoRoadmapOrContracts(t *testing.T) {
	dir := t.TempDir()
	written, _, err := scaffold(dir, "owner/hub", fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range written {
		if f == "ROADMAP.md" || strings.HasPrefix(f, filepath.FromSlash(config.RolesDir)) {
			t.Errorf("spoke scaffold wrote %s; the roadmap and the contracts live in the hub", f)
		}
	}
}

// The instructions live in .codecrew/AGENTS.md — CodeCrew's own directory,
// rewritable whole — and the root AGENTS.md is only the pointer at them, in
// both forms: the sentence a plain-markdown harness follows and the bare
// @-import Claude Code resolves (M13-R3).
func TestScaffoldWritesTheEntryPointAndARootPointer(t *testing.T) {
	for _, hub := range []string{"self", "org/hub"} {
		dir := t.TempDir()
		if _, _, err := scaffold(dir, hub, fakeContracts); err != nil {
			t.Fatal(err)
		}
		instructions, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(config.AgentsFile)))
		if err != nil {
			t.Fatalf("hub=%s: %v", hub, err)
		}
		if !strings.Contains(string(instructions), "gh codecrew roles show <role>") {
			t.Errorf("hub=%s: %s does not carry the instructions:\n%s", hub, config.AgentsFile, instructions)
		}
		root, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
		if err != nil {
			t.Fatalf("hub=%s: %v", hub, err)
		}
		// The pointer is short: a heading and the lines, nothing an adopter
		// has to re-merge when the instructions change.
		if len(strings.Split(strings.TrimSpace(string(root)), "\n")) > 6 {
			t.Errorf("hub=%s: root AGENTS.md is not a short pointer:\n%s", hub, root)
		}
		for _, want := range []string{config.AgentsFile, "@" + config.AgentsFile} {
			if !strings.Contains(string(root), want) {
				t.Errorf("hub=%s: root AGENTS.md does not carry %q:\n%s", hub, want, root)
			}
		}
		if strings.Contains(string(root), "`@"+config.AgentsFile+"`") {
			t.Errorf("hub=%s: the import is wrapped in backticks — Claude Code reads that as literal text", hub)
		}
		if !strings.Contains(string(root), entryPointLines) {
			t.Errorf("hub=%s: root AGENTS.md is not the block init prints:\n%s", hub, root)
		}
	}
}

// An existing root entry point is kept — and when it does not reach the
// instructions, init prints the exact lines to add to it, byte for byte the
// ones its own pointer carries, rather than reporting a skip that leaves the
// project with instructions nothing arrives at. A kept file that already
// reaches them asks for nothing: a rerun on what init wrote is idempotent in
// its output too, not only on disk (M13-R3, SPEC §6; checky on PR #278).
func TestInitPrintsTheLineToAddForAStrandedEntryPoint(t *testing.T) {
	for _, c := range []struct {
		name string
		// what is on disk before init runs
		existing map[string]string
		// the files init must name under "action needed"; none means the
		// block must not appear at all
		stranded []string
	}{
		{"a foreign AGENTS.md", map[string]string{"AGENTS.md": "# mine\n"}, []string{"AGENTS.md"}},
		{"a foreign CLAUDE.md", map[string]string{"CLAUDE.md": "# mine\n"}, []string{"CLAUDE.md"}},
		{"both foreign", map[string]string{"AGENTS.md": "# mine\n", "CLAUDE.md": "# mine\n"}, []string{"AGENTS.md", "CLAUDE.md"}},
		{"a fresh repo", nil, nil},
		{"a rerun on what init wrote", map[string]string{"AGENTS.md": agentsPointerScaffold, "CLAUDE.md": claudeScaffold}, nil},
		// CLAUDE.md's import is only as good as the file it lands on: init's
		// own CLAUDE.md beside somebody else's AGENTS.md reaches nothing.
		{"init's CLAUDE.md over a foreign AGENTS.md", map[string]string{"AGENTS.md": "# mine\n", "CLAUDE.md": claudeScaffold}, []string{"AGENTS.md", "CLAUDE.md"}},
	} {
		dir := t.TempDir()
		for f, content := range c.existing {
			if err := os.WriteFile(filepath.Join(dir, f), []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		wd, _ := os.Getwd()
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		var out bytes.Buffer
		err := initCmd(&out, nil)
		os.Chdir(wd)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		got := out.String()

		// Whatever was there is kept, byte for byte, and reported.
		for f, content := range c.existing {
			if !strings.Contains(got, "kept existing "+f) {
				t.Errorf("%s: init did not report keeping %s:\n%s", c.name, f, got)
			}
			if data, _ := os.ReadFile(filepath.Join(dir, f)); string(data) != content {
				t.Errorf("%s: init overwrote %s", c.name, f)
			}
		}

		if len(c.stranded) == 0 {
			if strings.Contains(got, "action needed") {
				t.Errorf("%s: init asked for an action it does not need:\n%s", c.name, got)
			}
			continue
		}
		for _, f := range c.stranded {
			if !strings.Contains(got, "Kept: ") || !strings.Contains(got, f) {
				t.Errorf("%s: the action-needed heading does not name %s:\n%s", c.name, f, got)
			}
		}
		// The lines it prints are the ones its own root pointer carries.
		if !strings.Contains(got, entryPointLines) {
			t.Errorf("%s: init printed no line to add:\n%s", c.name, got)
		}
		if !strings.Contains(agentsPointerScaffold, entryPointLines) {
			t.Error("the printed block is not what the root pointer scaffold contains")
		}
		if !strings.Contains(got, config.AgentsFile) {
			t.Errorf("%s: the report never names %s:\n%s", c.name, config.AgentsFile, got)
		}
	}
}

// init scaffolds rather than loading a pointer, so it never meets the
// pointer check — but a repo on the 1.x layout is still refused: a second
// layout written beside the first is the one outcome nobody can migrate
// from (M13-R1).
func TestInitRefusesTheLegacyLayout(t *testing.T) {
	for _, c := range []struct {
		name  string
		write func(dir string)
		found string
	}{
		{"the 1.x pointer", func(dir string) {
			os.WriteFile(filepath.Join(dir, ".codecrew.yml"), []byte("codecrew: \"1.0\"\nhub: self\n"), 0o644)
		}, ".codecrew.yml"},
		{"a 1.x roles/", func(dir string) {
			os.MkdirAll(filepath.Join(dir, "roles"), 0o755)
			os.WriteFile(filepath.Join(dir, "roles", "implementer.md"), []byte("# Role: implementer\n"), 0o644)
		}, "roles/implementer.md"},
	} {
		dir := t.TempDir()
		c.write(dir)
		wd, _ := os.Getwd()
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		var out bytes.Buffer
		err := initCmd(&out, nil)
		os.Chdir(wd)

		var r refusal
		if !errors.As(err, &r) || r.Code != "LAYOUT_LEGACY" {
			t.Fatalf("%s: err = %v, want refused[LAYOUT_LEGACY]", c.name, err)
		}
		for _, want := range []string{c.found, "gh codecrew migrate"} {
			if !strings.Contains(r.Detail, want) {
				t.Errorf("%s: detail %q does not name %q", c.name, r.Detail, want)
			}
		}
		if out.Len() != 0 {
			t.Errorf("%s: a refused init printed %q", c.name, out.String())
		}
		if _, err := os.Stat(filepath.Join(dir, ".codecrew")); err == nil {
			t.Errorf("%s: a refused init wrote a second layout", c.name)
		}
	}
}
