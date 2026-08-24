package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestBuildManifestPermissionsPerRole(t *testing.T) {
	for role, want := range rolePermissions {
		m, err := buildManifest(role, "myorg-crew", "https://github.com/o/r", "http://127.0.0.1:1/callback", false, "")
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
	if _, err := buildManifest("navigator", "x", "u", "r", false, ""); err == nil {
		t.Error("unknown role accepted")
	}
	if _, err := buildManifest("qa", "", "u", "r", false, ""); err == nil {
		t.Error("empty name accepted")
	}
	if _, err := buildManifest("qa", "qa", "u", "r", false, ""); err == nil {
		t.Error("role-named App accepted — identities outlive role reassignments")
	}
	if _, err := buildManifest("qa", "myorg-testy", "u", "r", true, ""); err == nil {
		t.Error("--with-webhook without --webhook-url accepted")
	}
}

func TestBuildManifestWebhookDefaultsOff(t *testing.T) {
	m, _ := buildManifest("reviewer", "myorg-reviewy", "u", "r", false, "")
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
	m, err := buildManifest("reviewer", "myorg-reviewy", "u", "r", true, "https://platform.example/hook")
	if err != nil {
		t.Fatal(err)
	}
	hook := m["hook_attributes"].(map[string]any)
	if hook["active"] != true || hook["url"] != "https://platform.example/hook" {
		t.Errorf("hook_attributes = %v", hook)
	}
	events := m["default_events"].([]string)
	if len(events) == 0 {
		t.Fatal("no events with --with-webhook")
	}
	for _, want := range []string{"issues", "pull_request_review"} {
		found := false
		for _, e := range events {
			if e == want {
				found = true
			}
		}
		if !found {
			t.Errorf("protocol-traffic event %q missing from %v", want, events)
		}
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
