package cli

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	codecrew "github.com/radiusred/gh-codecrew"
	"testing/fstest"
)

var fakeContracts = fstest.MapFS{
	"roles/implementer.md": {Data: []byte("# Role: implementer\n")},
	"roles/qa.md":          {Data: []byte("# Role: qa\n")},
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
	for _, want := range []string{".codecrew.yml", "ROADMAP.md", "AGENTS.md", "CLAUDE.md", filepath.Join("roles", "qa.md"), filepath.Join("roles", "qa"+localSuffix), filepath.Join("roles", "implementer"+localSuffix)} {
		if !slices.Contains(written, want) {
			t.Errorf("missing %s from written %v", want, written)
		}
		if _, err := os.Stat(filepath.Join(dir, want)); err != nil {
			t.Errorf("%s not on disk: %v", want, err)
		}
	}
	cfg, err := os.ReadFile(filepath.Join(dir, ".codecrew.yml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, role := range []string{"implementer", "reviewer", "qa", "doc-synthesizer", "coordinator"} {
		if !strings.Contains(string(cfg), role) {
			t.Errorf("hub config missing role %q", role)
		}
	}
}

func TestScaffoldSpoke(t *testing.T) {
	dir := t.TempDir()
	written, _, err := scaffold(dir, "org/hub", fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 1 || written[0] != ".codecrew.yml" {
		t.Fatalf("spoke mode wrote %v, want only the pointer", written)
	}
	cfg, _ := os.ReadFile(filepath.Join(dir, ".codecrew.yml"))
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
	data, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
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
			t.Errorf("scaffolded AGENTS.md missing %q", want)
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
		data, _ := os.ReadFile(filepath.Join(dir, ".codecrew.yml"))
		if !strings.HasPrefix(string(data), "codecrew: \""+protocolVersion+"\"") {
			t.Errorf("hub=%s: pointer starts %q, want codecrew: %q", hub, strings.SplitN(string(data), "\n", 2)[0], protocolVersion)
		}
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

func TestScaffoldSpokeWritesNoClaude(t *testing.T) {
	dir := t.TempDir()
	written, _, err := scaffold(dir, "owner/hub", fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range written {
		if f == "CLAUDE.md" || f == "AGENTS.md" {
			t.Errorf("spoke scaffold wrote %s; spokes get only the pointer", f)
		}
	}
}

func TestScaffoldKeepsExistingClaude(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "CLAUDE.md")
	if err := os.WriteFile(marker, []byte("mine\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, skipped, err := scaffold(dir, "self", fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(skipped, "CLAUDE.md") {
		t.Errorf("existing CLAUDE.md not reported as skipped: %v", skipped)
	}
	if data, _ := os.ReadFile(marker); string(data) != "mine\n" {
		t.Error("existing CLAUDE.md was overwritten")
	}
}

// init writes a blank extension beside every contract it scaffolds — the
// mechanism made visible at onboarding, holding only the comment that says
// what the file is for (M7-R4). Blank means comments-only: it composes to
// nothing, an existing extension is never touched, and a spoke gets none.
func TestScaffoldWritesBlankExtensions(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := scaffold(dir, "self", fakeContracts); err != nil {
		t.Fatal(err)
	}
	for _, role := range []string{"implementer", "qa"} {
		data, err := os.ReadFile(filepath.Join(dir, "roles", role+localSuffix))
		if err != nil {
			t.Fatalf("no blank extension for %s: %v", role, err)
		}
		s := string(data)
		for _, want := range []string{"roles/" + role + ".local.md", "roles/" + role + ".md", "gh codecrew roles show " + role, U + "/docs/extensions.md", U + "/SPEC.md"} {
			if !strings.Contains(s, want) {
				t.Errorf("%s extension lacks %q", role, want)
			}
		}
		if strings.TrimSpace(withoutHTMLComments(s)) != "" {
			t.Errorf("%s extension is not comments-only: %q", role, s)
		}
	}
	// Rerunning init keeps a written extension and reports it.
	p := filepath.Join(dir, "roles", "qa"+localSuffix)
	os.WriteFile(p, []byte("- House style.\n"), 0o644)
	written, skipped, err := scaffold(dir, "self", fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 0 || !slices.Contains(skipped, filepath.Join("roles", "qa"+localSuffix)) {
		t.Errorf("rerun wrote %v, skipped %v", written, skipped)
	}
	if data, _ := os.ReadFile(p); string(data) != "- House style.\n" {
		t.Error("rerun overwrote a written extension")
	}
	spoke := t.TempDir()
	written, _, _ = scaffold(spoke, "org/hub", fakeContracts)
	for _, w := range written {
		if strings.HasSuffix(w, localSuffix) {
			t.Errorf("spoke scaffold wrote an extension: %s", w)
		}
	}
}
