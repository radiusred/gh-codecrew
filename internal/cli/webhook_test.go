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
		if events == nil {
			events = []string{}
		}
		switch r.Method + " " + r.URL.Path {
		case "GET /app":
			fmt.Fprintf(w, `{"slug":"myorg-reviewy","events":%s,"owner":{"login":"myorg","type":"Organization"}}`, mustJSON(events))
		case "GET /app/hook/config", "PATCH /app/hook/config":
			if cfg == nil { // an App minted without a webhook: GitHub's real answer
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"message":"Not Found","documentation_url":"https://docs.github.com/rest/apps/webhooks","status":"404"}`)
				return
			}
			if r.Method == "PATCH" {
				var p map[string]string
				json.NewDecoder(r.Body).Decode(&p)
				patches = append(patches, p)
				if v, ok := p["url"]; ok {
					cfg.URL = v
				}
				if v, ok := p["secret"]; ok {
					cfg.Secret = v
				}
			}
			json.NewEncoder(w).Encode(masked())
		case "never":
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

// An App minted without a webhook has no hook configuration: GitHub
// answers 404 to GET and PATCH alike (verified live on this org's reviewer
// App) and has no endpoint to create one — so every path refuses
// NO_WEBHOOK naming the settings page, and writes nothing.
func TestIdentityWebhookNoHookAndBadKey(t *testing.T) {
	srv, patches := fakeAppAPI(t, nil, nil)
	for _, args := range [][]string{{"--show"}, {}, {"--url", "https://r.example/fire"}, {"--secret", "x"}, {"--rotate-secret"}} {
		_, err := runWebhook(t, srv, args...)
		var r refusal
		if !errors.As(err, &r) || r.Code != "NO_WEBHOOK" || !strings.Contains(r.Detail, "/settings/apps/myorg-reviewy") {
			t.Errorf("%v: err = %v, want NO_WEBHOOK naming the settings page", args, err)
		}
	}
	if len(*patches) != 0 {
		t.Errorf("a hookless App was written to: %v", *patches)
	}
	var err error
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

// identity new after the browser flow, with --webhook-secret: the
// receiver's secret is set under the key just stored and never printed;
// when the set fails, GitHub's generated secret is shown once with the
// recovery command. The conversion itself is the browser flow and is not
// exercised here — only what follows it.
func TestStoreAndReportWebhookSecret(t *testing.T) {
	prev := setWebhookSecret
	defer func() { setWebhookSecret = prev }()
	creds := &appCreds{ID: 77, Slug: "myorg-reviewy", ClientID: "Iv1.x", PEM: string(pkcs1PEM(testKey)), WebhookSecret: "gh-generated"}
	var got []string
	setWebhookSecret = func(c *appCreds, s string) error {
		got = append(got, fmt.Sprintf("%d:%s", c.ID, s))
		return nil
	}
	var out bytes.Buffer
	if _, stub, err := storeAndReport(&out, creds, t.TempDir(), "receiver-secret", false, time.Unix(0, 0)); err != nil || !strings.HasSuffix(stub, "myorg-reviewy.json") {
		t.Fatalf("err %v stub %q", err, stub)
	}
	if strings.Join(got, ",") != "77:receiver-secret" {
		t.Errorf("set calls = %v", got)
	}
	if !strings.Contains(out.String(), "webhook secret set to the one you supplied (not stored, not printed)") || strings.Contains(out.String(), "receiver-secret") || strings.Contains(out.String(), "gh-generated") {
		t.Errorf("success output:\n%s", out.String())
	}
	// The set fails: the generated secret is handed over once, with the way back.
	setWebhookSecret = func(*appCreds, string) error { return errors.New("HTTP 404") }
	out.Reset()
	if _, _, err := storeAndReport(&out, creds, t.TempDir(), "receiver-secret", false, time.Unix(0, 0)); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"setting yours failed (HTTP 404)", "gh codecrew identity webhook myorg-reviewy --secret <yours>", "GitHub generated (shown once, not stored): gh-generated"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("failure output lacks %q:\n%s", want, out.String())
		}
	}
	if strings.Contains(out.String(), "receiver-secret") {
		t.Error("the supplied secret was printed")
	}
	// No --webhook-secret: GitHub's is shown once with the alternative; nothing is set.
	got = nil
	setWebhookSecret = func(c *appCreds, s string) error { got = append(got, s); return nil }
	out.Reset()
	storeAndReport(&out, creds, t.TempDir(), "", true, time.Unix(0, 0))
	if len(got) != 0 || !strings.Contains(out.String(), "webhook secret (shown once, not stored): gh-generated — or set the receiver's") || !strings.Contains(out.String(), "contents: write granted") {
		t.Errorf("default output:\n%s (set calls %v)", out.String(), got)
	}
	// An App minted without a webhook has no secret to report.
	creds.WebhookSecret = ""
	out.Reset()
	storeAndReport(&out, creds, t.TempDir(), "", false, time.Unix(0, 0))
	if strings.Contains(out.String(), "webhook secret") {
		t.Errorf("hookless App reported a secret:\n%s", out.String())
	}
}
