package cli

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// defaultWebhookEvents is what --with-webhook subscribes an App to: the
// two transitions a platform routes to seats (a PR opened or updated → the
// reviewer; a review posted → the implementer). The 1.0 set also carried
// issues, issue_comment and check_suite, which on a platform are wakes —
// and woke the coordinator for nothing (#119; #164 findings 46, 53).
// `issues` is the opt-in for a receiver that wants gates (cc:needs-decision).
var defaultWebhookEvents = []string{"pull_request", "pull_request_review"}

// webhookEventPermission maps each event an App may subscribe to onto the
// permission GitHub requires for it — a manifest that subscribes to an
// event its permissions cannot receive is refused by GitHub after the
// browser flow, which is too late to be told.
var webhookEventPermission = map[string]string{
	"pull_request":                "pull_requests",
	"pull_request_review":         "pull_requests",
	"pull_request_review_comment": "pull_requests",
	"issues":                      "issues",
	"issue_comment":               "issues",
	"check_suite":                 "checks",
	"check_run":                   "checks",
	"push":                        "contents",
}

// webhookEvents validates a subscription list against a role's permission
// set; an empty list is the default pair.
func webhookEvents(role string, requested []string) ([]string, error) {
	if len(requested) == 0 {
		requested = defaultWebhookEvents
	}
	perms := rolePermissions[role]
	var out []string
	seen := map[string]bool{}
	for _, e := range requested {
		e = strings.TrimSpace(e)
		if e == "" || seen[e] {
			continue
		}
		need, known := webhookEventPermission[e]
		if !known {
			keys := make([]string, 0, len(webhookEventPermission))
			for k := range webhookEventPermission {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			return nil, fmt.Errorf("unknown webhook event %q — one of %s", e, strings.Join(keys, ", "))
		}
		if _, ok := perms[need]; !ok {
			return nil, fmt.Errorf("event %q needs the %s permission, which the %s role's set does not carry (docs/identities.md)", e, need, role)
		}
		seen[e] = true
		out = append(out, e)
	}
	return out, nil
}

// hookConfig is GitHub's App webhook configuration (GET/PATCH
// /app/hook/config, under the App's JWT). The secret comes back masked.
type hookConfig struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	Secret      string `json:"secret"`
	InsecureSSL string `json:"insecure_ssl"`
}

// noWebhook is the refusal for an App minted without a webhook: GitHub
// keeps no hook configuration for it (GET and PATCH /app/hook/config both
// answer 404 — verified live on this org's reviewer App) and offers no
// endpoint to create one, so activation is a settings-page act; the verb
// takes over from there.
func noWebhook(settings string) error {
	return refuse("NO_WEBHOOK", "this App was minted without a webhook and GitHub's API cannot create one — activate it under Webhook on the App's settings page (%s), give it the receiver's URL there, then set the secret with identity webhook --secret and tick the events under Subscribe to events", settings)
}

func appRequest(client *http.Client, jwt, method, path string, body any) (int, []byte, error) {
	var payload io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		payload = bytes.NewReader(data)
	}
	req, _ := http.NewRequest(method, githubAPI+path, payload)
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Accept", "application/vnd.github+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode == http.StatusUnauthorized {
		return resp.StatusCode, data, refuse("BAD_CREDENTIALS", "GitHub rejected the App JWT (401): the private key and the App id do not belong to the same App, or the key was revoked — check the id against gh api /apps/<slug> --jq .id; retrying will not help")
	}
	return resp.StatusCode, data, nil
}

// getHookConfig reads the App's webhook configuration.
func getHookConfig(client *http.Client, jwt, settings string) (hookConfig, error) {
	status, data, err := appRequest(client, jwt, "GET", "/app/hook/config", nil)
	if err != nil {
		return hookConfig{}, err
	}
	if status == http.StatusNotFound {
		return hookConfig{}, noWebhook(settings)
	}
	if status != http.StatusOK {
		return hookConfig{}, fmt.Errorf("GET /app/hook/config: HTTP %d: %s", status, strings.TrimSpace(string(data)))
	}
	var cfg hookConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return hookConfig{}, err
	}
	return cfg, nil
}

// patchHookConfig updates the fields given (url, secret) and returns the
// configuration GitHub now holds.
func patchHookConfig(client *http.Client, jwt, settings string, patch map[string]string) (hookConfig, error) {
	status, data, err := appRequest(client, jwt, "PATCH", "/app/hook/config", patch)
	if err != nil {
		return hookConfig{}, err
	}
	if status == http.StatusNotFound {
		return hookConfig{}, noWebhook(settings)
	}
	if status != http.StatusOK {
		return hookConfig{}, fmt.Errorf("PATCH /app/hook/config: HTTP %d: %s", status, strings.TrimSpace(string(data)))
	}
	var cfg hookConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return hookConfig{}, err
	}
	return cfg, nil
}

// appSubscriptions reads the App's slug, owner and subscribed events — the
// events an existing App has are readable here and changeable only on its
// settings page (there is no endpoint for them; verified live).
func appSubscriptions(client *http.Client, jwt string) (slug, ownerLogin, ownerType string, events []string, err error) {
	status, data, err := appRequest(client, jwt, "GET", "/app", nil)
	if err != nil {
		return "", "", "", nil, err
	}
	if status != http.StatusOK {
		return "", "", "", nil, fmt.Errorf("GET /app: HTTP %d: %s", status, strings.TrimSpace(string(data)))
	}
	var app struct {
		Slug   string   `json:"slug"`
		Events []string `json:"events"`
		Owner  struct {
			Login string `json:"login"`
			Type  string `json:"type"`
		} `json:"owner"`
	}
	if err := json.Unmarshal(data, &app); err != nil {
		return "", "", "", nil, err
	}
	return app.Slug, app.Owner.Login, app.Owner.Type, app.Events, nil
}

// newWebhookSecret is 32 random bytes as hex — what --rotate-secret sets.
func newWebhookSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// identityWebhook is `gh codecrew identity webhook <slug> [--show] [--url U]
// [--secret S | --rotate-secret]`: the App's hook, read and set under its
// own key — the URL and secret any time; the events only where GitHub
// lets them be set, which is the settings page, so the verb says where.
func identityWebhook(w io.Writer, args []string) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	return runIdentityWebhook(w, os.Getenv, configDir, http.DefaultClient, time.Now(), args)
}

func runIdentityWebhook(w io.Writer, getenv func(string) string, configDir string, client *http.Client, now time.Time, args []string) error {
	fs := flag.NewFlagSet("identity webhook", flag.ContinueOnError)
	fs.SetOutput(w)
	slug, args := splitLeadingRef(args)
	show := fs.Bool("show", false, "print the hook's URL, content type, whether a secret is set, and the App's subscribed events")
	url := fs.String("url", "", "set the receiver URL")
	secret := fs.String("secret", "", "set the hook secret to the receiver's (never stored; never printed)")
	rotate := fs.Bool("rotate-secret", false, "generate a new secret, set it, print it once")
	events := fs.String("events", "", "not settable after creation — the verb prints where they change")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if slug == "" && fs.NArg() == 1 {
		slug = fs.Arg(0)
	}
	if (slug != "" && fs.NArg() > 0) || fs.NArg() > 1 {
		return fmt.Errorf("usage: gh codecrew identity webhook <slug> [--show] [--url U] [--secret S | --rotate-secret]")
	}
	if *secret != "" && *rotate {
		return fmt.Errorf("identity webhook: --secret and --rotate-secret are exclusive")
	}
	cred, err := resolveCredential(getenv, configDir, slug)
	if err != nil {
		return err
	}
	jwt, err := signAppJWT(cred.Key, cred.Issuer, now)
	if err != nil {
		return refuse("BAD_CREDENTIALS", "%s: %v", cred.Source, err)
	}
	appSlug, ownerLogin, ownerType, subscribed, err := appSubscriptions(client, jwt)
	if err != nil {
		return err
	}
	settings := appSettingsURL(ownerLogin, ownerType, appSlug)
	if *events != "" {
		return fmt.Errorf("an existing App's event subscriptions cannot be set through the API (GitHub offers no endpoint; they are set by the manifest at creation, or by hand) — tick them under Subscribe to events at %s/permissions; subscribed now: %s", settings, strings.Join(subscribed, ", "))
	}
	patch := map[string]string{}
	if *url != "" {
		patch["url"] = *url
	}
	var minted string
	switch {
	case *secret != "":
		patch["secret"] = *secret
	case *rotate:
		if minted, err = newWebhookSecret(); err != nil {
			return err
		}
		patch["secret"] = minted
	}
	var cfg hookConfig
	if len(patch) > 0 {
		if cfg, err = patchHookConfig(client, jwt, settings, patch); err != nil {
			return err
		}
		if patch["url"] != "" {
			fmt.Fprintf(w, "hook URL set: %s\n", cfg.URL)
		}
		if *secret != "" {
			fmt.Fprintln(w, "hook secret set to the one you supplied (not stored, not printed)")
		}
		if minted != "" {
			fmt.Fprintf(w, "hook secret rotated (shown once, not stored): %s\n", minted)
		}
	}
	if *show || len(patch) == 0 {
		if len(patch) == 0 {
			if cfg, err = getHookConfig(client, jwt, settings); err != nil {
				return err
			}
		}
		fmt.Fprintf(w, "%s webhook\n", appSlug)
		fmt.Fprintf(w, "  url:          %s\n", orNone(cfg.URL))
		fmt.Fprintf(w, "  content type: %s\n", orNone(cfg.ContentType))
		if cfg.Secret != "" {
			fmt.Fprintln(w, "  secret:       set (masked by GitHub; rotate with --rotate-secret)")
		} else {
			fmt.Fprintln(w, "  secret:       none — set the receiver's with --secret")
		}
		fmt.Fprintf(w, "  events:       %s\n", orNone(strings.Join(subscribed, ", ")))
		fmt.Fprintf(w, "  events change only on the App's settings page: %s/permissions\n", settings)
	}
	return nil
}

func orNone(s string) string {
	if s == "" {
		return "none"
	}
	return s
}
