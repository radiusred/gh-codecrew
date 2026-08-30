package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/radiusred/gh-codecrew/internal/config"
)

func TestBuildManifestPermissionsPerRole(t *testing.T) {
	for role, want := range rolePermissions {
		m, err := buildManifest(role, "myorg-crew", "https://github.com/o/r", "http://127.0.0.1:1/callback", false, "", nil, false)
		if err != nil {
			t.Fatalf("%s: %v", role, err)
		}
		got := m["default_permissions"].(map[string]string)
		for perm, level := range want {
			if got[perm] != level {
				t.Errorf("%s: permission %s = %q, want %q", role, perm, got[perm], level)
			}
		}
		if len(got) != len(want) {
			t.Errorf("%s: %d permissions, want %d (minimal set only)", role, len(got), len(want))
		}
	}
}

func TestBuildManifestRefusals(t *testing.T) {
	if _, err := buildManifest("navigator", "x", "u", "r", false, "", nil, false); err == nil {
		t.Error("unknown role accepted")
	}
	if _, err := buildManifest("qa", "", "u", "r", false, "", nil, false); err == nil {
		t.Error("empty name accepted")
	}
	if _, err := buildManifest("qa", "qa", "u", "r", false, "", nil, false); err == nil {
		t.Error("role-named App accepted — identities outlive role reassignments")
	}
	if _, err := buildManifest("qa", "myorg-testy", "u", "r", true, "", nil, false); err == nil {
		t.Error("--with-webhook without --webhook-url accepted")
	}
}

func TestBuildManifestWebhookDefaultsOff(t *testing.T) {
	m, _ := buildManifest("reviewer", "myorg-reviewy", "u", "r", false, "", nil, false)
	// GitHub requires hook_attributes.url whenever the object is present,
	// regardless of active (the #73 finding) — a webhook-less manifest
	// must omit the object entirely.
	if _, ok := m["hook_attributes"]; ok {
		t.Error("hook_attributes present without --with-webhook — GitHub refuses it urlless")
	}
	if _, ok := m["default_events"]; ok {
		t.Error("events subscribed without --with-webhook")
	}
}

func TestBuildManifestWithWebhook(t *testing.T) {
	m, err := buildManifest("reviewer", "myorg-reviewy", "u", "r", true, "https://platform.example/hook", nil, false)
	if err != nil {
		t.Fatal(err)
	}
	hook := m["hook_attributes"].(map[string]any)
	if hook["active"] != true || hook["url"] != "https://platform.example/hook" {
		t.Errorf("hook_attributes = %v", hook)
	}
	// The default is the two transitions a platform routes to seats; the
	// 1.0 set's issues/issue_comment/check_suite were wakes for nothing on
	// a platform (#164 findings 46, 53) and are opt-in now.
	if events := m["default_events"].([]string); strings.Join(events, ",") != "pull_request,pull_request_review" {
		t.Errorf("default events = %v", events)
	}
	m, err = buildManifest("reviewer", "myorg-reviewy", "u", "r", true, "https://platform.example/hook", []string{"pull_request", "issues", " pull_request "}, false)
	if err != nil {
		t.Fatal(err)
	}
	if events := m["default_events"].([]string); strings.Join(events, ",") != "pull_request,issues" {
		t.Errorf("--events = %v (deduplicated, trimmed, in order)", events)
	}
	// An event GitHub does not have, or one the role's permissions cannot
	// receive (the coordinator has no checks), is refused before the
	// browser flow; push needs only contents: read, which every set has.
	if _, err := buildManifest("reviewer", "myorg-reviewy", "u", "r", true, "https://platform.example/hook", []string{"release"}, false); err == nil {
		t.Error("unknown event accepted")
	}
	if _, err := buildManifest("coordinator", "myorg-loopy", "u", "r", true, "https://platform.example/hook", []string{"check_suite"}, false); err == nil {
		t.Error("check_suite accepted for the coordinator, whose set has no checks")
	}
	if _, err := buildManifest("implementer", "myorg-coder", "u", "r", true, "https://platform.example/hook", []string{"push"}, false); err != nil {
		t.Errorf("push refused for the implementer, whose set carries contents: %v", err)
	}
}

func TestManifestTarget(t *testing.T) {
	if got := manifestTarget("radiusred", "Organization"); got != "https://github.com/organizations/radiusred/settings/apps/new" {
		t.Errorf("org target = %q", got)
	}
	if got := manifestTarget("davison", "User"); got != "https://github.com/settings/apps/new" {
		t.Errorf("personal target = %q", got)
	}
}

func TestPemPathConvention(t *testing.T) {
	now := time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC)
	got := pemPath("/home/x/.config", "myorg-coder", now)
	want := "/home/x/.config/codecrew/myorg-coder.2026-08-24.private-key.pem"
	if got != want {
		t.Errorf("pemPath = %q, want %q", got, want)
	}
}

// TestServeFlow drives the loopback half end to end: the form page carries
// the manifest and target, and the callback hands the code back.
func TestServeFlow(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	manifest, _ := json.Marshal(map[string]any{"name": "myorg-testy"})
	target := "https://github.com/organizations/o/settings/apps/new"

	type result struct {
		code string
		err  error
	}
	done := make(chan result, 1)
	go func() {
		code, err := serveFlow(l, manifest, target, 5*time.Second)
		done <- result{code, err}
	}()

	base := fmt.Sprintf("http://%s", l.Addr())
	resp, err := http.Get(base + "/")
	if err != nil {
		t.Fatal(err)
	}
	page, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if !strings.Contains(string(page), target) {
		t.Errorf("form page missing target %q:\n%s", target, page)
	}
	if !strings.Contains(string(page), "myorg-testy") {
		t.Errorf("form page missing manifest content:\n%s", page)
	}

	if _, err := http.Get(base + "/callback?" + url.Values{"code": {"tempcode123"}}.Encode()); err != nil {
		t.Fatal(err)
	}
	r := <-done
	if r.err != nil {
		t.Fatal(r.err)
	}
	if r.code != "tempcode123" {
		t.Errorf("code = %q, want tempcode123", r.code)
	}
}

func TestServeFlowRejectsCodelessCallback(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go serveFlow(l, []byte("{}"), "t", 2*time.Second)
	resp, err := http.Get(fmt.Sprintf("http://%s/callback", l.Addr()))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("codeless callback status = %d, want 400", resp.StatusCode)
	}
}

const nestedYML = `codecrew: "0.1" # protocol version
hub: self

# Advisory role routing, read by whoever dispatches agents.
roles:
  implementer:
    harness: claude-code
    identity: myorg-coder
  reviewer:
    harness: codex
    identity: alice # the bootstrap human, by name
  qa:
    harness: codex
    identity: myorg-testy
`

const inlineYML = `codecrew: "0.1"
hub: self
roles:
  implementer: { identity: ~ }
  reviewer: { identity: ~ }
  qa: { identity: ~ }
  doc-synthesizer: { identity: ~ }
`

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), ".codecrew.yml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRouteRoleNestedShape(t *testing.T) {
	p := writeTemp(t, nestedYML)
	if err := routeRole(p, "reviewer", "myorg-reviewy"); err != nil {
		t.Fatal(err)
	}
	out, _ := os.ReadFile(p)
	cfg, err := config.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if got := cfg.Roles["reviewer"].Identity; got != "myorg-reviewy" {
		t.Errorf("reviewer identity = %q", got)
	}
	// The stale trailing comment described the old routing — gone.
	if strings.Contains(string(out), "bootstrap human") {
		t.Error("stale identity comment survived the rewrite")
	}
	// Everything else survives: siblings, harness lines, file comments.
	for _, keep := range []string{"# Advisory role routing", "harness: codex", "identity: myorg-coder", "identity: myorg-testy", `codecrew: "0.1" # protocol version`} {
		if !strings.Contains(string(out), keep) {
			t.Errorf("line lost in surgery: %q", keep)
		}
	}
}

func TestRouteRoleInlineShape(t *testing.T) {
	p := writeTemp(t, inlineYML)
	if err := routeRole(p, "reviewer", "myorg-reviewy"); err != nil {
		t.Fatal(err)
	}
	out, _ := os.ReadFile(p)
	cfg, err := config.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if got := cfg.Roles["reviewer"].Identity; got != "myorg-reviewy" {
		t.Errorf("reviewer identity = %q", got)
	}
	for _, other := range []string{"implementer", "qa", "doc-synthesizer"} {
		if cfg.Roles[other].Identity != "" {
			t.Errorf("%s identity changed to %q", other, cfg.Roles[other].Identity)
		}
	}
}

func TestRouteRoleErrors(t *testing.T) {
	unchanged := func(t *testing.T, p, want string) {
		t.Helper()
		out, _ := os.ReadFile(p)
		if string(out) != want {
			t.Error("file modified despite error")
		}
	}
	p := writeTemp(t, nestedYML)
	if err := routeRole(p, "doc-synthesizer", "myorg-wordy"); err == nil {
		t.Error("absent role routed")
	}
	unchanged(t, p, nestedYML)

	spoke := writeTemp(t, "codecrew: \"0.1\"\nhub: myorg/hub\n")
	if err := routeRole(spoke, "reviewer", "x"); err == nil {
		t.Error("routed into a pointer-only spoke config")
	}
}

func TestAppSettingsURL(t *testing.T) {
	if got := appSettingsURL("radiusred", "Organization", "radiusred-reviewy"); got != "https://github.com/organizations/radiusred/settings/apps/radiusred-reviewy" {
		t.Errorf("org settings URL = %q", got)
	}
	if got := appSettingsURL("davison", "User", "davison-reviewy"); got != "https://github.com/settings/apps/davison-reviewy" {
		t.Errorf("personal settings URL = %q", got)
	}
}

func TestWriteCredentials(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	creds := &appCreds{ID: 4704266, Slug: "davison-review-bot", ClientID: "Iv1.abc", PEM: "PEMDATA", ClientSecret: "never-on-disk", WebhookSecret: "nor-this"}
	key, stub, err := writeCredentials(dir, creds, now)
	if err != nil {
		t.Fatal(err)
	}
	if key != filepath.Join(dir, "codecrew", "davison-review-bot.2026-08-25.private-key.pem") {
		t.Errorf("key path = %q", key)
	}
	if stub != filepath.Join(dir, "codecrew", "davison-review-bot.json") {
		t.Errorf("stub path = %q", stub)
	}
	info, err := os.Stat(key)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("key mode = %v, want 0600", info.Mode().Perm())
	}
	data, _ := os.ReadFile(stub)
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got["app_id"] != float64(4704266) || got["slug"] != "davison-review-bot" || got["client_id"] != "Iv1.abc" {
		t.Errorf("stub content = %v", got)
	}
	for _, secret := range []string{"never-on-disk", "nor-this"} {
		if strings.Contains(string(data), secret) {
			t.Errorf("secret %q written to the stub", secret)
		}
	}
}

func TestBuildManifestWithApprovalPermission(t *testing.T) {
	m, err := buildManifest("reviewer", "myorg-reviewy", "u", "r", false, "", nil, true)
	if err != nil {
		t.Fatal(err)
	}
	got := m["default_permissions"].(map[string]string)
	if got["contents"] != "write" {
		t.Errorf("contents = %q, want write", got["contents"])
	}
	// Exactly one permission changes; the rest of the reviewer set stands.
	for perm, want := range rolePermissions["reviewer"] {
		if perm == "contents" {
			continue
		}
		if got[perm] != want {
			t.Errorf("permission %s = %q, want %q (only contents may change)", perm, got[perm], want)
		}
	}
	if len(got) != len(rolePermissions["reviewer"]) {
		t.Errorf("%d permissions, want %d — the flag adds no new scopes", len(got), len(rolePermissions["reviewer"]))
	}
	// The shared table itself must not be mutated by the grant.
	if rolePermissions["reviewer"]["contents"] != "read" {
		t.Error("rolePermissions table mutated — the grant must copy")
	}
	// Reviewer-only: approvals gate merges nowhere else.
	for _, role := range []string{"implementer", "qa", "doc-synthesizer", "coordinator"} {
		if _, err := buildManifest(role, "myorg-x", "u", "r", false, "", nil, true); err == nil {
			t.Errorf("--with-approval-permission accepted for %s", role)
		}
	}
	// Default path unchanged.
	m, _ = buildManifest("reviewer", "myorg-reviewy", "u", "r", false, "", nil, false)
	if m["default_permissions"].(map[string]string)["contents"] != "read" {
		t.Error("default reviewer manifest no longer read-only")
	}
}

// The coordinator seat mints the set #119 finding 16 specified: it opens
// issues, comments and labels and never pushes, reviews or merges — so
// contents and pull requests stay read, and no flag can widen them.
func TestCoordinatorManifestIsReadOnlyOnCode(t *testing.T) {
	m, err := buildManifest("coordinator", "myorg-loopy", "u", "r", false, "", nil, false)
	if err != nil {
		t.Fatal(err)
	}
	perms := m["default_permissions"].(map[string]string)
	want := map[string]string{"contents": "read", "issues": "write", "pull_requests": "read"}
	if len(perms) != len(want) {
		t.Errorf("coordinator permissions = %v, want exactly %v", perms, want)
	}
	for k, v := range want {
		if perms[k] != v {
			t.Errorf("coordinator %s = %q, want %q", k, perms[k], v)
		}
	}
	if _, ok := perms["checks"]; ok {
		t.Error("coordinator granted checks — it reads gate results through the verbs, not the API")
	}
	if _, err := buildManifest("coordinator", "coordinator", "u", "r", false, "", nil, false); err == nil {
		t.Error("App named after the role accepted")
	}
}
