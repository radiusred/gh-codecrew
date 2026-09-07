package tracker

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/gh"
)

// scriptGH stands one scripted gh behind gh.Command: every call prints out
// on stdout and exits with code, and every call's arguments are appended to
// a file the test reads back. Like fakeGH it is this test binary re-entered
// — no script on disk, no PATH shim — but scripted rather than routed, so
// each test says what gh answered and then asserts what it was asked.
func scriptGH(t *testing.T, out string, code int) (args func() []string) {
	t.Helper()
	log := filepath.Join(t.TempDir(), "calls")
	orig := gh.Command
	gh.Command = func(_ string, a ...string) *exec.Cmd {
		cmd := exec.Command(os.Args[0], append([]string{"-test.run=^TestHelperVenueGH$", "--"}, a...)...)
		cmd.Env = append(os.Environ(),
			"VENUE_GH=1",
			"VENUE_GH_OUT="+out,
			"VENUE_GH_EXIT="+strconv.Itoa(code),
			"VENUE_GH_LOG="+log,
		)
		return cmd
	}
	t.Cleanup(func() { gh.Command = orig })
	return func() []string {
		data, err := os.ReadFile(log)
		if err != nil {
			return nil
		}
		return strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	}
}

// TestHelperVenueGH is the scripted gh's body; inert unless scriptGH
// re-enters the binary with VENUE_GH set.
func TestHelperVenueGH(t *testing.T) {
	if os.Getenv("VENUE_GH") != "1" {
		return
	}
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	if log := os.Getenv("VENUE_GH_LOG"); log != "" {
		f, err := os.OpenFile(log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		if err == nil {
			f.WriteString(strings.Join(args, " ") + "\n")
			f.Close()
		}
	}
	os.Stdout.WriteString(os.Getenv("VENUE_GH_OUT"))
	if os.Getenv("VENUE_GH_EXIT") != "0" {
		os.Exit(1)
	}
	os.Exit(0)
}

// The venue's URL grammar. It is knowledge, not I/O — there is nothing to
// script and nothing to fake — so it is package functions beside the
// interface, and this is where it is pinned.
func TestRecordAPIPath(t *testing.T) {
	for _, tc := range []struct {
		url, path string
		ok        bool
	}{
		{"https://github.com/o/r/issues/68", "repos/o/r/issues/68", true},
		{"https://github.com/o/r/pull/100", "repos/o/r/issues/100", true},
		{"https://github.com/o/r/issues/68#issuecomment-99", "repos/o/r/issues/68", true},
		{"https://github.com/o/r/commit/abc123", "repos/o/r/commits/abc123", true},
		{"https://github.com/o/r/blob/main/docs/x.md", "repos/o/r/contents/docs/x.md?ref=main", true},
		{"https://github.com/apps/some-app", "", false}, // App pages: plain HTTP
		{"https://github.com/o/r", "", false},           // bare repo page
		{"https://docs.github.com/en/apps", "", false},  // not github.com
		{"https://example.com/o/r/issues/1", "", false},
	} {
		path, ok := RecordAPIPath(tc.url)
		if path != tc.path || ok != tc.ok {
			t.Errorf("RecordAPIPath(%q) = %q,%v want %q,%v", tc.url, path, ok, tc.path, tc.ok)
		}
	}
}

func TestIsRecordLink(t *testing.T) {
	for url, want := range map[string]bool{
		"https://github.com/o/r/issues/1":      true,
		"https://github.com/apps/some-app":     true,
		"http://github.com/o/r":                true,
		"https://docs.github.com/en/apps":      false,
		"https://github.community/x":           false,
		"https://example.com/github.com/o/r":   false,
		"https://codecrew.works/docs/protocol": false,
	} {
		if got := IsRecordLink(url); got != want {
			t.Errorf("IsRecordLink(%q) = %v, want %v", url, got, want)
		}
	}
}

func TestAppAndRepoURLs(t *testing.T) {
	for _, tc := range []struct{ got, want string }{
		{RepoURL("radiusred/gh-codecrew"), "https://github.com/radiusred/gh-codecrew"},
		{AppInstallURL("radiusred-cody"), "https://github.com/apps/radiusred-cody/installations/new"},
		{AppManifestURL("radiusred", "Organization"), "https://github.com/organizations/radiusred/settings/apps/new"},
		{AppManifestURL("davison", "User"), "https://github.com/settings/apps/new"},
		{AppSettingsURL("radiusred", "Organization", "radiusred-reviewy"), "https://github.com/organizations/radiusred/settings/apps/radiusred-reviewy"},
		{AppSettingsURL("davison", "User", "davison-reviewy"), "https://github.com/settings/apps/davison-reviewy"},
	} {
		if tc.got != tc.want {
			t.Errorf("URL = %q, want %q", tc.got, tc.want)
		}
	}
}

// RepoInDir asks in a named directory, not the process's own — init runs
// it against the tree it is scaffolding — and answers with one error for
// every incomplete answer, so the caller has one condition to handle
// rather than four.
func TestRepoInDirAnswersOrErrors(t *testing.T) {
	calls := scriptGH(t, `{"nameWithOwner":"o/r","defaultBranchRef":{"name":"main"}}`, 0)
	repo, branch, err := GitHub{}.RepoInDir(t.TempDir())
	if err != nil || repo != "o/r" || branch != "main" {
		t.Fatalf("RepoInDir = %q,%q,%v", repo, branch, err)
	}
	if got := calls(); len(got) != 1 || got[0] != "repo view --json nameWithOwner,defaultBranchRef" {
		t.Errorf("asked gh %v", got)
	}

	for _, tc := range []struct {
		name, out string
		code      int
	}{
		{"the client failed", "", 1},
		{"the output does not parse", "not json", 0},
		{"no repository in the answer", `{"defaultBranchRef":{"name":"main"}}`, 0},
		{"no default branch in the answer", `{"nameWithOwner":"o/r"}`, 0},
	} {
		scriptGH(t, tc.out, tc.code)
		if repo, branch, err := (GitHub{}).RepoInDir(t.TempDir()); err == nil {
			t.Errorf("%s: RepoInDir = %q,%q,nil — an incomplete answer must be an error", tc.name, repo, branch)
		}
	}
}

func TestBranchRuleTypes(t *testing.T) {
	calls := scriptGH(t, `[{"type":"deletion"},{"type":"pull_request"}]`, 0)
	types, err := GitHub{}.BranchRuleTypes("o/r", "main")
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 2 || types[0] != "deletion" || types[1] != "pull_request" {
		t.Errorf("types = %v", types)
	}
	if got := calls(); len(got) != 1 || !strings.Contains(got[0], "repos/o/r/rules/branches/main") {
		t.Errorf("asked gh %v", got)
	}
}

func TestTeamMembersAndAccountType(t *testing.T) {
	calls := scriptGH(t, `[{"login":"alice"},{"login":"bob"}]`, 0)
	logins, err := GitHub{}.TeamMembers("myorg", "review-crew")
	if err != nil || len(logins) != 2 || logins[0] != "alice" || logins[1] != "bob" {
		t.Fatalf("TeamMembers = %v, %v", logins, err)
	}
	if got := calls(); len(got) != 1 || !strings.Contains(got[0], "--paginate /orgs/myorg/teams/review-crew/members") {
		t.Errorf("asked gh %v — the listing must be paginated", got)
	}

	// The login is escaped: an App's account login carries a `[bot]`
	// suffix, and the brackets are not path characters.
	calls = scriptGH(t, `{"type":"Bot"}`, 0)
	kind, err := GitHub{}.AccountType("radiusred-cody[bot]")
	if err != nil || kind != "Bot" {
		t.Fatalf("AccountType = %q, %v", kind, err)
	}
	if got := calls(); len(got) != 1 || !strings.Contains(got[0], "users/radiusred-cody%5Bbot%5D") {
		t.Errorf("asked gh %v — the login must be path-escaped", got)
	}
}

// APIReachable answers whether the path resolved and nothing else: the body
// is never read, so a caller cannot come to depend on it.
func TestAPIReachable(t *testing.T) {
	calls := scriptGH(t, `{"number":68}`, 0)
	if err := (GitHub{}).APIReachable("repos/o/r/issues/68"); err != nil {
		t.Fatal(err)
	}
	if got := calls(); len(got) != 1 || got[0] != "api repos/o/r/issues/68" {
		t.Errorf("asked gh %v", got)
	}
	scriptGH(t, "", 1)
	if err := (GitHub{}).APIReachable("repos/o/r/issues/404"); err == nil {
		t.Error("a dead path resolved")
	}
}

func TestAppManifestConversion(t *testing.T) {
	calls := scriptGH(t, `{"id":42,"slug":"myorg-coder","client_id":"Iv1.x","pem":"KEY"}`, 0)
	creds, err := GitHub{}.AppManifestConversion("code-123")
	if err != nil {
		t.Fatal(err)
	}
	if creds.ID != 42 || creds.Slug != "myorg-coder" || creds.PEM != "KEY" {
		t.Errorf("creds = %+v", creds)
	}
	if got := calls(); len(got) != 1 || got[0] != "api --method POST app-manifests/code-123/conversions" {
		t.Errorf("asked gh %v", got)
	}
}

// AppRequest is the one transport that is not gh — an App JWT is neither a
// user nor an installation credential. It signs the call, reads the body,
// and names no condition of its own: the status comes back for the caller
// to classify, a 401 included.
func TestAppRequestSignsAndClassifiesNothing(t *testing.T) {
	type seen struct {
		method, path, auth, accept, contentType, body string
	}
	var got seen
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = seen{r.Method, r.URL.Path, r.Header.Get("Authorization"), r.Header.Get("Accept"), r.Header.Get("Content-Type"), string(body)}
		if r.URL.Path == "/unauthorized" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message":"A JSON web token could not be decoded"}`))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"token":"ghs_abc"}`))
	}))
	defer srv.Close()
	prev := AppAPIBase
	AppAPIBase = srv.URL
	defer func() { AppAPIBase = prev }()

	status, data, err := GitHub{}.AppRequest(srv.Client(), "the.jwt", "PATCH", "/app/hook/config", map[string]string{"secret": "s"})
	if err != nil {
		t.Fatal(err)
	}
	if status != http.StatusCreated || string(data) != `{"token":"ghs_abc"}` {
		t.Errorf("status/body = %d %q", status, data)
	}
	want := seen{"PATCH", "/app/hook/config", "Bearer the.jwt", "application/vnd.github+json", "application/json", `{"secret":"s"}`}
	if got != want {
		t.Errorf("request = %+v, want %+v", got, want)
	}

	// No body: no Content-Type, and nothing sent.
	if _, _, err := (GitHub{}).AppRequest(srv.Client(), "the.jwt", "GET", "/app/installations?per_page=100", nil); err != nil {
		t.Fatal(err)
	}
	if got.contentType != "" || got.body != "" {
		t.Errorf("a bodyless call sent %q / %q", got.contentType, got.body)
	}

	// A 401 is a status, not a refusal: the mint and the webhook verbs each
	// name their own condition from it, and this method names none.
	status, data, err = GitHub{}.AppRequest(srv.Client(), "bad", "GET", "/unauthorized", nil)
	if err != nil {
		t.Fatalf("AppRequest classified a 401: %v", err)
	}
	if status != http.StatusUnauthorized || !strings.Contains(string(data), "could not be decoded") {
		t.Errorf("401 = %d %q", status, data)
	}
}

// The classifiers the CLI reads through the venue rather than through
// internal/gh, which now has exactly one importer: this package.
func TestVenueClassifiers(t *testing.T) {
	if !Unreachable(errTest("gh api: dial tcp: lookup api.github.com: no such host")) {
		t.Error("an offline error was not classified as unreachable")
	}
	if Unreachable(errTest("gh api: Not Found (HTTP 404)")) {
		t.Error("a 404 was classified as unreachable — GitHub answered it")
	}
	if CompareVersions("2.46.0", "2.50.0") >= 0 || CompareVersions("2.50.0", "2.50.0") != 0 {
		t.Error("CompareVersions does not order the gh floor")
	}
}

type errTest string

func (e errTest) Error() string { return string(e) }

// The App-credential shape is the venue's, and the mint reads GitHub's
// field names off it.
func TestAppCredentialsFieldNames(t *testing.T) {
	var creds AppCredentials
	if err := json.Unmarshal([]byte(`{"id":7,"slug":"s","client_id":"c","client_secret":"cs","webhook_secret":"ws","pem":"p","html_url":"h"}`), &creds); err != nil {
		t.Fatal(err)
	}
	want := AppCredentials{ID: 7, Slug: "s", ClientID: "c", ClientSecret: "cs", WebhookSecret: "ws", PEM: "p", HTMLURL: "h"}
	if creds != want {
		t.Errorf("creds = %+v, want %+v", creds, want)
	}
}
