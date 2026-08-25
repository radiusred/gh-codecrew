package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStampRoundTrip(t *testing.T) {
	stamped := contractStamp("implementer.md") + "# Role: implementer\nBody.\n"
	if got := stripStamp(stamped); got != "# Role: implementer\nBody.\n" {
		t.Errorf("stripStamp = %q", got)
	}
	// Unstamped content — the hub's own roles/ — passes through whole.
	if got := stripStamp("# Role: qa\n"); got != "# Role: qa\n" {
		t.Errorf("unstamped content altered: %q", got)
	}
}

func TestFreshScaffoldShowsNoDrift(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := scaffold(dir, "self", fakeContracts); err != nil {
		t.Fatal(err)
	}
	drifted, err := contractDrift(dir, fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	if len(drifted) != 0 {
		t.Errorf("fresh scaffold drifted: %v — the stamp must strip cleanly", drifted)
	}
	// The scaffolded file really is stamped.
	data, _ := os.ReadFile(filepath.Join(dir, "roles", "implementer.md"))
	if !strings.HasPrefix(string(data), stampPrefix) {
		t.Error("scaffolded contract missing the provenance stamp")
	}
}

func TestLocalEditIsDrift(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := scaffold(dir, "self", fakeContracts); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, "roles", "qa.md")
	data, _ := os.ReadFile(p)
	os.WriteFile(p, append(data, []byte("\n- Local convention: tests ride along.\n")...), 0o644)
	drifted, err := contractDrift(dir, fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	if len(drifted) != 1 || drifted[0] != "qa" {
		t.Errorf("drift = %v, want [qa]", drifted)
	}
}

func TestContractDriftSkipsSpokes(t *testing.T) {
	drifted, err := contractDrift(t.TempDir(), fakeContracts) // no roles/ dir at all
	if err != nil {
		t.Fatal(err)
	}
	if len(drifted) != 0 {
		t.Errorf("spoke (no roles/) reported drift: %v", drifted)
	}
}

func TestUnifiedDiff(t *testing.T) {
	got := unifiedDiff("a\nb\nc\n", "a\nx\nc\n")
	want := "  a\n- b\n+ x\n  c\n"
	if got != want {
		t.Errorf("unifiedDiff = %q, want %q", got, want)
	}
}

func TestRolesDiffAndShow(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := scaffold(dir, "self", fakeContracts); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := rolesDiff(&buf, dir, fakeContracts, "implementer"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "matches the embedded") {
		t.Errorf("clean diff output: %q", buf.String())
	}
	// Edit locally, expect the divergence and the reconciliation line.
	p := filepath.Join(dir, "roles", "implementer.md")
	data, _ := os.ReadFile(p)
	os.WriteFile(p, append(data, []byte("local line\n")...), 0o644)
	buf.Reset()
	if err := rolesDiff(&buf, dir, fakeContracts, "implementer"); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"- local line", "reconcile through a task and PR"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("diff output missing %q:\n%s", want, buf.String())
		}
	}
	buf.Reset()
	if err := rolesShow(&buf, fakeContracts, "qa", true); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "# Role: qa\n" {
		t.Errorf("show --latest = %q", buf.String())
	}
	if err := rolesShow(&buf, fakeContracts, "qa", false); err == nil {
		t.Error("show without --latest accepted")
	}
	if err := rolesDiff(&buf, dir, fakeContracts, "navigator"); err == nil {
		t.Error("unknown role accepted")
	}
}
