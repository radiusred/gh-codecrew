package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fakeAppAPI serves /app and /app/hook/config the way GitHub does for an
// App JWT: the secret comes back masked, PATCH merges the fields given.
func fakeAppAPI(t *testing.T, cfg *hookConfig, events []string) (*httptest.Server, *[]map[string]string) {
	t.Helper()
	var patches []map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		claims := verifyJWT(t, auth)
		if claims["iss"] == "wrong" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		masked := func() hookConfig {
			c := *cfg
			if c.Secret != "" {
				c.Secret = "********"
			}
			return c
		}
		switch r.Method + " " + r.URL.Path {
		case "GET /app":
			fmt.Fprintf(w, `{"slug":"myorg-reviewy","events":%s,"owner":{"login":"myorg","type":"Organization"}}`, mustJSON(events))
		case "GET /app/hook/config":
			json.NewEncoder(w).Encode(masked())
		case "PATCH /app/hook/config":
			var p map[string]string
			json.NewDecoder(r.Body).Decode(&p)
			patches = append(patches, p)
			if v, ok := p["url"]; ok {
				cfg.URL = v
			}
			if v, ok := p["secret"]; ok {
				cfg.Secret = v
			}
			json.NewEncoder(w).Encode(masked())
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, &patches
}

func mustJSON(v any) string { b, _ := json.Marshal(v); return string(b) }

func runWebhook(t *testing.T, srv *httptest.Server, args ...string) (string, error) {
	t.Helper()
	prev := githubAPI
	githubAPI = srv.URL
	defer func() { githubAPI = prev }()
	env := envOf(map[string]string{"GITHUB_APP_ID": "7", "GITHUB_PEM": string(pkcs1PEM(testKey))})
	var out bytes.Buffer
	err := runIdentityWebhook(&out, env, t.TempDir(), srv.Client(), time.Now(), args)
	return out.String(), err
}

func TestIdentityWebhookShowAndSet(t *testing.T) {
	cfg := &hookConfig{URL: "https://old.example/hook", ContentType: "json", Secret: "s3cr3t"}
	srv, patches := fakeAppAPI(t, cfg, []string{"pull_request", "pull_request_review"})
	// No flags: show. The secret is never printed, only its presence.
	out, err := runWebhook(t, srv)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"myorg-reviewy webhook", "url:          https://old.example/hook", "secret:       set (masked", "events:       pull_request, pull_request_review", "settings page: https://github.com/organizations/myorg/settings/apps/myorg-reviewy/permissions"} {
		if !strings.Contains(out, want) {
			t.Errorf("show lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "s3cr3t") || len(*patches) != 0 {
		t.Errorf("show leaked the secret or wrote: %q %v", out, *patches)
	}
	// --url and --secret patch exactly those fields; the secret is not echoed.
	out, err = runWebhook(t, srv, "--url", "https://new.example/fire", "--secret", "receiver-secret")
	if err != nil {
		t.Fatal(err)
	}
	if len(*patches) != 1 || (*patches)[0]["url"] != "https://new.example/fire" || (*patches)[0]["secret"] != "receiver-secret" {
		t.Errorf("patches = %v", *patches)
	}
	if !strings.Contains(out, "hook URL set: https://new.example/fire") || !strings.Contains(out, "secret set to the one you supplied") || strings.Contains(out, "receiver-secret") {
		t.Errorf("set output:\n%s", out)
	}
	// --rotate-secret mints 32 bytes of hex, sets it and prints it once.
	out, err = runWebhook(t, srv, "--rotate-secret")
	if err != nil {
		t.Fatal(err)
	}
	minted := (*patches)[1]["secret"]
	if len(minted) != 64 || !strings.Contains(out, "rotated (shown once, not stored): "+minted) {
		t.Errorf("rotate: secret %q, out %q", minted, out)
	}
	if cfg.Secret != minted {
		t.Error("rotate did not set the secret GitHub holds")
	}
	// The two are exclusive; --events is refused with the settings URL.
	if _, err := runWebhook(t, srv, "--secret", "x", "--rotate-secret"); err == nil {
		t.Error("--secret with --rotate-secret accepted")
	}
	_, err = runWebhook(t, srv, "--events", "issues")
	if err == nil || !strings.Contains(err.Error(), "no endpoint") || !strings.Contains(err.Error(), "/settings/apps/myorg-reviewy/permissions") {
		t.Errorf("--events: %v", err)
	}
}

func TestIdentityWebhookNoHookAndBadKey(t *testing.T) {
	cfg := &hookConfig{}
	srv, _ := fakeAppAPI(t, cfg, nil)
	out, err := runWebhook(t, srv, "--show")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "url:          none") || !strings.Contains(out, "secret:       none — set the receiver's with --secret") || !strings.Contains(out, "events:       none") {
		t.Errorf("empty hook:\n%s", out)
	}
	prev := githubAPI
	githubAPI = srv.URL
	defer func() { githubAPI = prev }()
	var buf bytes.Buffer
	err = runIdentityWebhook(&buf, envOf(map[string]string{"GITHUB_APP_ID": "wrong", "GITHUB_PEM": string(pkcs1PEM(testKey))}), t.TempDir(), srv.Client(), time.Now(), []string{"--show"})
	var r refusal
	if !errors.As(err, &r) || r.Code != "BAD_CREDENTIALS" {
		t.Errorf("bad key: %v", err)
	}
	// No credentials at all: the shared resolver's refusal.
	err = runIdentityWebhook(&buf, envOf(nil), t.TempDir(), srv.Client(), time.Now(), []string{"myorg-ghost"})
	if !errors.As(err, &r) || r.Code != "NO_CREDENTIALS" {
		t.Errorf("no creds: %v", err)
	}
}

func TestWebhookEventsTable(t *testing.T) {
	got, err := webhookEvents("coordinator", nil)
	if err != nil || strings.Join(got, ",") != "pull_request,pull_request_review" {
		t.Errorf("coordinator default = %v, %v", got, err)
	}
	if _, err := webhookEvents("coordinator", []string{"check_suite"}); err == nil {
		t.Error("check_suite accepted for a role without checks")
	}
	if got, _ := webhookEvents("coordinator", []string{"issues", "issue_comment"}); strings.Join(got, ",") != "issues,issue_comment" {
		t.Errorf("issues events = %v", got)
	}
}
