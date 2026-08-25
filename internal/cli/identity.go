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
	"regexp"
	"strings"
	"time"

	"github.com/radiusred/gh-codecrew/internal/config"
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
func buildManifest(role, name, homepage, redirectURL string, withWebhook bool, webhookURL string, withApproval bool) (map[string]any, error) {
	perms, ok := rolePermissions[role]
	if !ok {
		return nil, fmt.Errorf("unknown role %q — one of implementer, reviewer, qa, doc-synthesizer", role)
	}
	if withApproval {
		// The one permission that makes an App's approvals count toward
		// required-review rules (#73, superseding Decision). Reviewer only:
		// approvals gate merges nowhere else, so anywhere else the flag
		// would be privilege without meaning. Never the default — privilege
		// is not acquired by default (the #59 rule, inverted).
		if role != "reviewer" {
			return nil, fmt.Errorf("--with-approval-permission applies only to the reviewer role — approvals gate merges nowhere else")
		}
		granted := map[string]string{}
		for k, v := range perms {
			granted[k] = v
		}
		granted["contents"] = "write"
		perms = granted
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

// appSettingsURL is the created App's settings page — where the manual
// display steps live (the manifest has no avatar field and no API uploads
// one, so the crew logo is added by hand under Display information).
func appSettingsURL(owner, ownerType, slug string) string {
	if ownerType == "Organization" {
		return "https://github.com/organizations/" + owner + "/settings/apps/" + slug
	}
	return "https://github.com/settings/apps/" + slug
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

// stubPath is the credential stub beside the key: the non-secret half of
// the credential set (App ID, client ID), persisted at minting so
// codecrew-token can mint installation tokens with the App's own JWT and
// no user-credential API lookup (finding 11 on #73).
func stubPath(configDir, slug string) string {
	return filepath.Join(configDir, "codecrew", slug+".json")
}

// writeCredentials stores what the conversion returned and must persist:
// the private key (0600) and the credential stub. Secrets other than the
// key are never written.
func writeCredentials(configDir string, creds *appCreds, now time.Time) (key, stub string, err error) {
	key = pemPath(configDir, creds.Slug, now)
	if err := os.MkdirAll(filepath.Dir(key), 0o755); err != nil {
		return "", "", err
	}
	if err := os.WriteFile(key, []byte(creds.PEM), 0o600); err != nil {
		return "", "", err
	}
	stub = stubPath(configDir, creds.Slug)
	data, err := json.MarshalIndent(map[string]any{
		"slug": creds.Slug, "app_id": creds.ID, "client_id": creds.ClientID,
	}, "", "  ")
	if err != nil {
		return "", "", err
	}
	if err := os.WriteFile(stub, append(data, '\n'), 0o644); err != nil {
		return "", "", err
	}
	return key, stub, nil
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

// routeRole rewrites role's identity to slug in the .codecrew.yml at path,
// by line surgery so the file's comments and layout survive. Both table
// shapes are handled: the scaffold's inline `role: { identity: ~ }` and a
// nested `identity:` line under the role key (whose trailing comment, if
// any, is dropped — it described the old routing). The result must
// re-parse with the role routed to slug, or nothing is written.
func routeRole(path, role, slug string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	inRoles, done := false, false
	roleIndent := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		indent := len(line) - len(strings.TrimLeft(line, " "))
		switch {
		case !inRoles:
			if trimmed == "roles:" {
				inRoles = true
			}
		case roleIndent == -1:
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				continue
			}
			if indent == 0 {
				return fmt.Errorf("role %q not found in the routing table", role)
			}
			if strings.HasPrefix(trimmed, role+":") {
				rest := strings.TrimSpace(strings.TrimPrefix(trimmed, role+":"))
				if strings.HasPrefix(rest, "{") {
					// Inline shape: rewrite the identity value in the braces.
					re := regexp.MustCompile(`(identity:\s*)[^,}]*`)
					if !re.MatchString(rest) {
						return fmt.Errorf("no identity key on the %s role's line", role)
					}
					lines[i] = line[:indent] + role + ": " + re.ReplaceAllString(rest, "${1}"+slug)
					done = true
				} else {
					roleIndent = indent // nested shape: identity on a deeper line
				}
			}
		default: // scanning inside the role's nested block
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				continue
			}
			if indent <= roleIndent {
				return fmt.Errorf("no identity key under the %s role", role)
			}
			if strings.HasPrefix(trimmed, "identity:") {
				lines[i] = line[:indent] + "identity: " + slug
				done = true
			}
		}
		if done {
			break
		}
	}
	if !done {
		return fmt.Errorf("role %q not found in the routing table", role)
	}
	out := strings.Join(lines, "\n")
	cfg, err := config.Parse([]byte(out))
	if err != nil {
		return fmt.Errorf("routing edit produced unparseable YAML: %w", err)
	}
	if cfg.Roles[role].Identity != slug {
		return fmt.Errorf("routing edit did not take")
	}
	return os.WriteFile(path, []byte(out), 0o644)
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
	noRoute := fs.Bool("no-route", false, "print the routing step instead of writing it into the hub's .codecrew.yml")
	withWebhook := fs.Bool("with-webhook", false, "subscribe the App to protocol-traffic events (platform users)")
	webhookURL := fs.String("webhook-url", "", "receiver for --with-webhook deliveries")
	withApproval := fs.Bool("with-approval-permission", false, "reviewer only: grant contents: write so the App's approvals satisfy required-review rules")
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
	if _, err := buildManifest(role, *name, "", "", *withWebhook, *webhookURL, *withApproval); err != nil {
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

	manifest, err := buildManifest(role, *name, "https://github.com/"+c.hub, local+"/callback", *withWebhook, *webhookURL, *withApproval)
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
	keyPath, stub, err := writeCredentials(configDir, creds, time.Now())
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "\ncreated %s (App ID %d, client ID %s)\n", creds.Slug, creds.ID, creds.ClientID)
	if *withApproval {
		fmt.Fprintln(w, "contents: write granted — this App's approvals satisfy required-review rules;")
		fmt.Fprintln(w, "its contract forbids editing code: the write grant exists only to make its judgment count")
	}
	fmt.Fprintf(w, "private key: %s\n", keyPath)
	fmt.Fprintf(w, "credential stub: %s (lets codecrew-token mint without any account lookup)\n", stub)
	if creds.WebhookSecret != "" {
		fmt.Fprintf(w, "webhook secret (shown once, not stored): %s\n", creds.WebhookSecret)
	}
	fmt.Fprintf(w, "\nnext:\n")
	fmt.Fprintf(w, "  1. install it: https://github.com/apps/%s/installations/new\n", creds.Slug)
	fmt.Fprintln(w, "     (installations are per-account — repeat for any other account it must reach)")
	routed := false
	if !*noRoute && c.current == c.hub {
		routed = routeRole(filepath.Join(c.cfg.Dir, ".codecrew.yml"), role, creds.Slug) == nil
	}
	if routed {
		fmt.Fprintf(w, "  2. routed %s → %s in .codecrew.yml — commit it via your next PR (--no-route to skip)\n", role, creds.Slug)
	} else {
		fmt.Fprintf(w, "  2. route the role in the hub's .codecrew.yml: roles.%s.identity: %s\n", role, creds.Slug)
	}
	fmt.Fprintf(w, "  3. optional: give it the crew logo under Display information: %s\n", appSettingsURL(*owner, acct.Type, creds.Slug))
	return nil
}
