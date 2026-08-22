package cli

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
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
