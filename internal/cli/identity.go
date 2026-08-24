package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/radiusred/gh-codecrew/internal/gh"
)

// rolePermissions is the minimal per-role permission set, mirroring the
// table in docs/identities.md. Metadata: read is implicit in the manifest
// flow. Workflows is deliberately absent — an implementer that will touch
// .github/workflows/ files gains it in the App's settings afterwards.
var rolePermissions = map[string]map[string]string{
	"implementer":     {"contents": "write", "issues": "write", "pull_requests": "write", "checks": "read"},
	"reviewer":        {"contents": "read", "issues": "write", "pull_requests": "write", "checks": "read"},
	"qa":              {"contents": "read", "issues": "write", "pull_requests": "write", "checks": "read"},
	"doc-synthesizer": {"contents": "write", "issues": "write", "pull_requests": "write", "checks": "read"},
}

// webhookEvents is the protocol-traffic event set --with-webhook subscribes
// an identity to: everything the CodeCrew verbs write that an orchestrator
// watching the seam (SPEC §9) would want delivered rather than polled.
var webhookEvents = []string{
	"issues", "issue_comment", "pull_request", "pull_request_review", "check_suite",
}

// buildManifest assembles the GitHub App manifest for a role. The App is
// named for a crew member, never the role (identities outlive role
// reassignments), so a name equal to the role is refused. A webhook-less
// manifest carries no hook_attributes at all — GitHub requires
// hook_attributes.url whenever the object is present, regardless of
// active (the #73 live-fire finding) — which registers the App with no
// webhook: a crew App acts, it never listens. withWebhook adds the object,
// active, with the protocol-traffic events delivered to webhookURL (the
// platform's receiver).
func buildManifest(role, name, homepage, redirectURL string, withWebhook bool, webhookURL string) (map[string]any, error) {
	perms, ok := rolePermissions[role]
	if !ok {
		return nil, fmt.Errorf("unknown role %q — one of implementer, reviewer, qa, doc-synthesizer", role)
	}
	if name == "" {
		return nil, fmt.Errorf("--name is required: Apps are named for crew members (myorg-coder), not roles")
	}
	if name == role {
		return nil, fmt.Errorf("name the App for a crew member, not the role %q — identities outlive role reassignments", role)
	}
	m := map[string]any{
		"name":                name,
		"url":                 homepage,
		"redirect_url":        redirectURL,
		"public":              false,
		"default_permissions": perms,
	}
	if withWebhook {
		if webhookURL == "" {
			return nil, fmt.Errorf("--with-webhook requires --webhook-url (the receiver events are delivered to)")
		}
		m["hook_attributes"] = map[string]any{"active": true, "url": webhookURL}
		m["default_events"] = webhookEvents
	}
	return m, nil
}

// manifestTarget is the page the manifest form posts to: the App-creation
// endpoint of the owning account, which differs for orgs and personal
// accounts.
func manifestTarget(owner, ownerType string) string {
	if ownerType == "Organization" {
		return "https://github.com/organizations/" + owner + "/settings/apps/new"
	}
	return "https://github.com/settings/apps/new"
}

// appCreds is what the one-hour conversion exchange returns. client_secret
// and webhook_secret are printed once and never written to disk.
type appCreds struct {
	ID            int64  `json:"id"`
	Slug          string `json:"slug"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	WebhookSecret string `json:"webhook_secret"`
	PEM           string `json:"pem"`
	HTMLURL       string `json:"html_url"`
}

// convertManifest exchanges the temporary code GitHub redirects back with
// for the created App's credentials. Unauthenticated by design on GitHub's
// side; routed through gh per the founding decision. Stubbed in tests.
var convertManifest = func(code string) (*appCreds, error) {
	var creds appCreds
	if err := gh.JSON(&creds, "api", "--method", "POST", "app-manifests/"+code+"/conversions"); err != nil {
		return nil, err
	}
	return &creds, nil
}

// pemPath is the storage convention from docs/identities.md:
// ~/.config/codecrew/<slug>.<date>.private-key.pem.
func pemPath(configDir, slug string, now time.Time) string {
	return filepath.Join(configDir, "codecrew", fmt.Sprintf("%s.%s.private-key.pem", slug, now.Format("2006-01-02")))
}

// formPage serves the auto-submitting manifest form. The manifest flow is a
// form POST, not a link, so the "one-click URL" the verb prints is this
// loopback page; the button is the JavaScript-less fallback.
var formPage = template.Must(template.New("form").Parse(`<!DOCTYPE html>
<html><body onload="document.forms[0].submit()">
<form action="{{.Target}}" method="post">
<input type="hidden" name="manifest" value="{{.Manifest}}">
<noscript><button type="submit">Create the GitHub App</button></noscript>
</form>
</body></html>
`))

// serveFlow runs the loopback half of the manifest flow on l: serves the
// form at /, receives GitHub's redirect at /callback, and returns the
// temporary code. Blocks until the callback arrives or timeout passes.
func serveFlow(l net.Listener, manifestJSON []byte, target string, timeout time.Duration) (string, error) {
	codes := make(chan string, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		formPage.Execute(w, struct{ Target, Manifest string }{target, string(manifestJSON)})
	})
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "missing code", http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "App created — return to the terminal.")
		select {
		case codes <- code:
		default:
		}
	})
	srv := &http.Server{Handler: mux}
	defer srv.Close()
	go srv.Serve(l)
	select {
	case code := <-codes:
		return code, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for the browser flow (the manifest code is valid for one hour — rerun when ready)")
	}
}

// identityNew mints a role's App identity via the manifest flow: the path
// decided at the #68 gate. It builds the manifest, hands the operator a
// one-click loopback URL, exchanges the redirect code, stores the private
// key under the docs/identities.md convention, and prints what remains
// manual: installing the App (per-account — installations do not cross
// account boundaries) and routing the role.
func identityNew(w io.Writer, args []string) error {
	fs := flag.NewFlagSet("identity new", flag.ContinueOnError)
	role, args := splitLeadingRef(args)
	name := fs.String("name", "", "App name — a crew member (myorg-coder), not a role (required)")
	owner := fs.String("owner", "", "account to own the App (default: the hub's owner)")
	withWebhook := fs.Bool("with-webhook", false, "subscribe the App to protocol-traffic events (platform users)")
	webhookURL := fs.String("webhook-url", "", "receiver for --with-webhook deliveries")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if role == "" && fs.NArg() == 1 {
		role = fs.Arg(0)
	}
	if role == "" || fs.NArg() > 1 {
		return fmt.Errorf("usage: codecrew identity new <role> --name <app-name>")
	}
	// Validate the manifest inputs before any repo or API work, so bad
	// invocations refuse instantly; the real manifest (with the loopback
	// redirect) is rebuilt once the listener's address is known.
	if _, err := buildManifest(role, *name, "", "", *withWebhook, *webhookURL); err != nil {
		return err
	}

	c, err := load()
	if err != nil {
		return err
	}
	if *owner == "" {
		*owner, _, _ = strings.Cut(c.hub, "/")
	}
	var acct struct {
		Type string `json:"type"`
	}
	if err := gh.JSON(&acct, "api", "users/"+*owner); err != nil {
		return err
	}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	defer l.Close()
	local := fmt.Sprintf("http://%s", l.Addr())

	manifest, err := buildManifest(role, *name, "https://github.com/"+c.hub, local+"/callback", *withWebhook, *webhookURL)
	if err != nil {
		return err
	}
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "open %s — it submits the manifest for %q to %s and GitHub asks you to confirm\n", local, *name, *owner)
	fmt.Fprintln(w, "waiting for the browser flow…")
	code, err := serveFlow(l, manifestJSON, manifestTarget(*owner, acct.Type), time.Hour)
	if err != nil {
		return err
	}
	creds, err := convertManifest(code)
	if err != nil {
		return err
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	keyPath := pemPath(configDir, creds.Slug, time.Now())
	if err := os.MkdirAll(filepath.Dir(keyPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(keyPath, []byte(creds.PEM), 0o600); err != nil {
		return err
	}

	fmt.Fprintf(w, "\ncreated %s (App ID %d, client ID %s)\n", creds.Slug, creds.ID, creds.ClientID)
	fmt.Fprintf(w, "private key: %s\n", keyPath)
	if creds.WebhookSecret != "" {
		fmt.Fprintf(w, "webhook secret (shown once, not stored): %s\n", creds.WebhookSecret)
	}
	fmt.Fprintf(w, "\nnext:\n")
	fmt.Fprintf(w, "  1. install it: https://github.com/apps/%s/installations/new\n", creds.Slug)
	fmt.Fprintln(w, "     (installations are per-account — repeat for any other account it must reach)")
	fmt.Fprintf(w, "  2. route the role in the hub's .codecrew.yml: roles.%s.identity: %s\n", role, creds.Slug)
	return nil
}
