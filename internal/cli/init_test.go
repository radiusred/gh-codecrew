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
	for _, want := range []string{".codecrew.yml", "ROADMAP.md", "AGENTS.md", filepath.Join("roles", "qa.md")} {
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
	for _, role := range []string{"implementer", "reviewer", "qa", "doc-synthesizer"} {
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

func TestInGitRepo(t *testing.T) {
	dir := t.TempDir()
	if inGitRepo(dir) {
		t.Error("bare temp dir should not read as a git repo")
	}
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if !inGitRepo(dir) {
		t.Error(".git directory should read as a git repo")
	}
	worktree := t.TempDir()
	if err := os.WriteFile(filepath.Join(worktree, ".git"), []byte("gitdir: elsewhere"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !inGitRepo(worktree) {
		t.Error(".git file (worktree) should read as a git repo")
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
