package cli

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	codecrew "github.com/radiusred/gh-codecrew"
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
	diskRead := func(path string) ([]byte, error) { return os.ReadFile(filepath.Join(dir, path)) }
	if err := rolesShow(&buf, "qa", true, fakeContracts, diskRead, ""); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "# Role: qa\n" {
		t.Errorf("show --latest = %q", buf.String())
	}
	// Without --latest, the hub's own (stamped) file is what a session loads.
	buf.Reset()
	if err := rolesShow(&buf, "qa", false, fakeContracts, diskRead, ""); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(buf.String(), stampPrefix) || !strings.Contains(buf.String(), "# Role: qa\n") {
		t.Errorf("show (composed, no extension) = %q", buf.String())
	}
	if err := rolesShow(&buf, "navigator", false, fakeContracts, diskRead, ""); err == nil {
		t.Error("show of a role with no hub contract accepted")
	}
	if err := rolesDiff(&buf, dir, fakeContracts, "navigator"); err == nil {
		t.Error("unknown role accepted")
	}
}

func TestComposeContractOrder(t *testing.T) {
	base := "# Role: qa\nBody.\n"
	if got := composeContract(base, nil); got != base {
		t.Errorf("no extensions altered the contract: %q", got)
	}
	got := composeContract(base, []localPart{
		{Source: "roles/qa.local.md (hub)", Body: "hub line"},
		{Source: "roles/qa.local.md (spoke)", Body: "   \n"}, // blank: skipped
		{Source: "roles/qa.local.md (spoke)", Body: "spoke line\n"},
	})
	want := "# Role: qa\nBody.\n\n<!-- extension: roles/qa.local.md (hub) -->\n\nhub line\n\n<!-- extension: roles/qa.local.md (spoke) -->\n\nspoke line\n"
	if got != want {
		t.Errorf("composeContract =\n%q\nwant\n%q", got, want)
	}
}

func TestLocalExtensionIsNotDrift(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := scaffold(dir, "self", fakeContracts); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(dir, "roles", "qa"+localSuffix), []byte("- House style.\n"), 0o644)
	drifted, err := contractDrift(dir, fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	if len(drifted) != 0 {
		t.Errorf("a local extension reported as drift: %v — extensions have no embedded counterpart and must stay invisible to the drift check", drifted)
	}
}

func TestComposedContractLayersHubThenSpoke(t *testing.T) {
	hub, spoke := t.TempDir(), t.TempDir()
	os.MkdirAll(filepath.Join(hub, "roles"), 0o755)
	os.MkdirAll(filepath.Join(spoke, "roles"), 0o755)
	os.WriteFile(filepath.Join(hub, "roles", "qa.md"), []byte("# Role: qa\n"), 0o644)
	os.WriteFile(filepath.Join(hub, "roles", "qa"+localSuffix), []byte("hub voice\n"), 0o644)
	os.WriteFile(filepath.Join(spoke, "roles", "qa"+localSuffix), []byte("spoke voice\n"), 0o644)
	// From a spoke the hub is read through a fetcher (the tracker in
	// production); here it reads the hub directory.
	hubRead := func(path string) ([]byte, error) { return os.ReadFile(filepath.Join(hub, path)) }

	got, err := composedContract("qa", hubRead, spoke)
	if err != nil {
		t.Fatal(err)
	}
	h, s := strings.Index(got, "hub voice"), strings.Index(got, "spoke voice")
	if !strings.HasPrefix(got, "# Role: qa\n") || h < 0 || s < 0 || h > s {
		t.Errorf("layering wrong (contract, hub, spoke):\n%s", got)
	}
	for _, marker := range []string{"(hub)", "(spoke)"} {
		if !strings.Contains(got, "<!-- extension: roles/qa.local.md "+marker+" -->") {
			t.Errorf("missing %s marker:\n%s", marker, got)
		}
	}
	// In the hub itself there is no spoke layer.
	got, err = composedContract("qa", hubRead, "")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "spoke voice") || !strings.Contains(got, "hub voice") {
		t.Errorf("hub-mode composition wrong:\n%s", got)
	}
	// A role with no contract in the hub is an error, extension or not.
	if _, err := composedContract("navigator", hubRead, spoke); err == nil {
		t.Error("missing contract accepted")
	}
}

// The shipped embed must hold contracts only: an extension inside the
// binary would be scaffolded into every new project and counted as drift
// (checky's finding on PR #123). Bound to the real codecrew.Roles, not a
// fixture — a fixture cannot contain what the glob would have caught.
func TestEmbeddedRolesHoldNoExtensions(t *testing.T) {
	entries, err := fs.ReadDir(codecrew.Roles, "roles")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("no embedded contracts")
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), localSuffix) {
			t.Errorf("embedded %s: extensions must not ship in the binary", e.Name())
		}
	}
	dir := t.TempDir()
	if _, _, err := scaffold(dir, "self", codecrew.Roles); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "roles", "doc-synthesizer"+localSuffix)); err == nil {
		t.Error("init scaffolded this project's editorial voice into a new project")
	}
}

// Even if an extension found its way into an embed, scaffold and drift
// ignore it by name.
func TestExtensionInEmbedIsSkipped(t *testing.T) {
	polluted := fstest.MapFS{
		"roles/qa.md":       {Data: []byte("# Role: qa\n")},
		"roles/qa.local.md": {Data: []byte("house style\n")},
	}
	dir := t.TempDir()
	written, _, err := scaffold(dir, "self", polluted)
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range written {
		if strings.HasSuffix(w, localSuffix) {
			t.Errorf("scaffold wrote %s", w)
		}
	}
	os.WriteFile(filepath.Join(dir, "roles", "qa"+localSuffix), []byte("customised\n"), 0o644)
	drifted, err := contractDrift(dir, polluted)
	if err != nil {
		t.Fatal(err)
	}
	if len(drifted) != 0 {
		t.Errorf("drift = %v, want none", drifted)
	}
}

// A fetch failure is not an absent file: the session must not run on a
// partial contract because the hub was unreachable.
func TestComposedContractSurfacesReadFailures(t *testing.T) {
	base := "# Role: qa\n"
	flaky := func(fail string) func(string) ([]byte, error) {
		return func(path string) ([]byte, error) {
			switch {
			case path == fail:
				return nil, errors.New("gh api: HTTP 500")
			case strings.HasSuffix(path, ".local.md"):
				return nil, fs.ErrNotExist
			default:
				return []byte(base), nil
			}
		}
	}
	if _, err := composedContract("qa", flaky("roles/qa.local.md"), ""); err == nil || !strings.Contains(err.Error(), "HTTP 500") {
		t.Errorf("extension fetch failure swallowed: %v", err)
	}
	if _, err := composedContract("qa", flaky("roles/qa.md"), ""); err == nil || !strings.Contains(err.Error(), "reading roles/qa.md") {
		t.Errorf("contract fetch failure misreported: %v", err)
	}
	if got, err := composedContract("qa", flaky("none"), ""); err != nil || got != base {
		t.Errorf("absent extension should be skipped: %q, %v", got, err)
	}
}

// The five contracts ship in the binary: the coordinator's is composed,
// diffed and scaffolded like the crew's, with no special-casing.
func TestCoordinatorContractIsEmbedded(t *testing.T) {
	data, err := fs.ReadFile(codecrew.Roles, "roles/coordinator.md")
	if err != nil {
		t.Fatalf("coordinator contract not embedded: %v", err)
	}
	if !strings.HasPrefix(string(data), "# Role: coordinator\n") {
		t.Errorf("embedded coordinator contract starts %q", string(data[:40]))
	}
	for _, must := range []string{"roles.coordinator.identity", "never contents: write", "task finish", "checkpoint", "--requirement"} {
		if !strings.Contains(string(data), must) {
			t.Errorf("coordinator contract lacks %q", must)
		}
	}
	dir := t.TempDir()
	written, _, err := scaffold(dir, "self", codecrew.Roles)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(written, "roles/coordinator.md") {
		t.Errorf("init did not scaffold the coordinator contract: %v", written)
	}
	var buf bytes.Buffer
	if err := rolesDiff(&buf, dir, codecrew.Roles, "coordinator"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "matches the embedded") {
		t.Errorf("fresh coordinator contract reported as drift: %q", buf.String())
	}
}
