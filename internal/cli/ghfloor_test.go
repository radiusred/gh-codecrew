package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The crew's container carried Debian's gh 2.46: task finish died inside
// `gh pr checks --json` and milestone close skipped its sweep for the same
// reason, and the agent found the floor by the failure (#119 findings 21,
// 30; #149). The floor is checked once, up front, with a code.
func TestCheckGHRefusesBelowTheFloor(t *testing.T) {
	defer func(orig func() (string, error)) { ghVersion = orig }(ghVersion)

	ghVersion = func() (string, error) { return "2.46.0", nil }
	var notes strings.Builder
	err := checkGH(&notes)
	if err == nil {
		t.Fatal("gh 2.46.0 passed the floor")
	}
	for _, want := range []string{"refused[GH_TOO_OLD]", "2.46.0", ghFloor, "pr checks --json"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("refusal missing %q: %v", want, err)
		}
	}

	for _, v := range []string{ghFloor, "2.98.0", "3.0.0"} {
		ghVersion = func() (string, error) { return v, nil }
		if err := checkGH(&notes); err != nil {
			t.Errorf("gh %s refused: %v", v, err)
		}
	}
	if notes.Len() != 0 {
		t.Errorf("parseable versions produced notes: %q", notes.String())
	}
}

// An unrecognised banner is reported, not refused: a build the parser has
// not met must not lock the operator out of every verb.
func TestCheckGHNotesAnUnparseableVersion(t *testing.T) {
	defer func(orig func() (string, error)) { ghVersion = orig }(ghVersion)
	ghVersion = func() (string, error) { return "", errors.New("gh --version: unrecognised output \"gh (homebrew)\"") }
	var notes strings.Builder
	if err := checkGH(&notes); err != nil {
		t.Fatalf("unparseable banner refused: %v", err)
	}
	if !strings.Contains(notes.String(), "note:") || !strings.Contains(notes.String(), ghFloor) {
		t.Errorf("note missing or does not name the floor: %q", notes.String())
	}
}

// Every verb that reads the pointer meets the check before its first gh
// call — status and roles included, which read the pointer without
// building a ctx (the reviewer of #153 found them bypassing it).
func TestEveryPointerReadingVerbMeetsTheFloor(t *testing.T) {
	defer func(orig func() (string, error)) { ghVersion = orig }(ghVersion)
	ghVersion = func() (string, error) { return "2.46.0", nil }

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".codecrew.yml"), []byte("codecrew: \"1.0\"\nhub: self\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	wd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	var out strings.Builder
	verbs := map[string]func() error{
		"status":     func() error { return status(&out) },
		"roles show": func() error { return rolesCmd(&out, []string{"show", "qa"}) },
		"load()":     func() error { _, err := load(); return err },
	}
	for name, run := range verbs {
		err := run()
		if err == nil || !strings.Contains(err.Error(), "refused[GH_TOO_OLD]") {
			t.Errorf("%s under gh 2.46.0: got %v, want refused[GH_TOO_OLD]", name, err)
		}
	}
}
