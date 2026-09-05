package cli

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

func TestExtractURLs(t *testing.T) {
	text := "See [the gate](https://github.com/radiusred/gh-codecrew/issues/68#issuecomment-1) and " +
		"https://example.com/page. Also <https://example.com/page> again, **https://docs.github.com/x** bold."
	got := extractURLs(text)
	want := []string{
		"https://github.com/radiusred/gh-codecrew/issues/68#issuecomment-1",
		"https://example.com/page",
		"https://docs.github.com/x",
	}
	if len(got) != len(want) {
		t.Fatalf("extractURLs = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("url[%d] = %q, want %q (dedup and trim must hold)", i, got[i], want[i])
		}
	}
}

func TestGithubAPIPath(t *testing.T) {
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
		path, ok := githubAPIPath(tc.url)
		if path != tc.path || ok != tc.ok {
			t.Errorf("githubAPIPath(%q) = %q,%v want %q,%v", tc.url, path, ok, tc.path, tc.ok)
		}
	}
}

func TestCheckURLRouting(t *testing.T) {
	origAPI, origHTTP := checkAPI, checkHTTP
	t.Cleanup(func() { checkAPI, checkHTTP = origAPI, origHTTP })
	var apiPaths, httpURLs []string
	checkAPI = func(path string) error { apiPaths = append(apiPaths, path); return nil }
	checkHTTP = func(url string) error { httpURLs = append(httpURLs, url); return fmt.Errorf("HTTP 404") }

	if err := checkURL("https://github.com/o/r/issues/5"); err != nil {
		t.Errorf("API-mapped link errored: %v", err)
	}
	if err := checkURL("https://example.com/gone"); err == nil {
		t.Error("dead HTTP link resolved")
	}
	if len(apiPaths) != 1 || apiPaths[0] != "repos/o/r/issues/5" {
		t.Errorf("API routing = %v", apiPaths)
	}
	if len(httpURLs) != 1 || httpURLs[0] != "https://example.com/gone" {
		t.Errorf("HTTP routing = %v", httpURLs)
	}
}

// TestCheckHTTPAgainstRealServer exercises the real HTTP checker — the
// plan promised an httptest round and the first commit shipped only
// stubs (checky's PR #102 finding, in the spirit of the very obligation
// this task adds).
func TestCheckHTTPAgainstRealServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			w.WriteHeader(http.StatusOK)
		case "/redirect":
			http.Redirect(w, r, "/ok", http.StatusFound)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	if err := checkHTTP(srv.URL + "/ok"); err != nil {
		t.Errorf("200 reported unreachable: %v", err)
	}
	if err := checkHTTP(srv.URL + "/redirect"); err != nil {
		t.Errorf("redirect-to-200 reported unreachable: %v", err)
	}
	if err := checkHTTP(srv.URL + "/gone"); err == nil {
		t.Error("404 reported reachable")
	} else if !strings.Contains(err.Error(), "404") {
		t.Errorf("404 classification lost: %v", err)
	}
	srv.Close()
	if err := checkHTTP(srv.URL + "/ok"); err == nil {
		t.Error("dead server reported reachable")
	}
}

// A URL ends at the first character that cannot be in one — the prose
// ellipsis that cost milestone evidence a real citation (#138), quotes, a
// backtick, brackets — and a closing parenthesis belongs to it only when
// it opened one; trailing sentence punctuation is trimmed.
func TestExtractURLsStopsWhereAURLCannotContinue(t *testing.T) {
	for _, tc := range []struct {
		text string
		want []string
	}{
		{"upstream URLs (`…/blob/main/…`) resolve from https://github.com/radiusred/gh-codecrew/blob/main/… as", []string{"https://github.com/radiusred/gh-codecrew/blob/main/"}},
		{"see https://github.com/o/r/issues/1…", []string{"https://github.com/o/r/issues/1"}},
		{"in `code` then https://github.com/o/r/pull/2 and \"https://example.com/a\" or 'https://example.com/b'", []string{"https://github.com/o/r/pull/2", "https://example.com/a", "https://example.com/b"}},
		{"an unclosed backtick still ends the URL https://github.com/o/r/pull/2`x", []string{"https://github.com/o/r/pull/2"}},
		{"(see https://github.com/o/r/issues/3).", []string{"https://github.com/o/r/issues/3"}},
		{"[link](https://github.com/o/r/issues/4), then [wiki](https://en.wikipedia.org/wiki/Foo_(bar)).", []string{"https://github.com/o/r/issues/4", "https://en.wikipedia.org/wiki/Foo_(bar)"}},
		{"ends the sentence https://example.com/x: next https://example.com/y; and https://example.com/z,", []string{"https://example.com/x", "https://example.com/y", "https://example.com/z"}},
		{"anchor https://github.com/o/r/issues/5#issuecomment-99 — dash", []string{"https://github.com/o/r/issues/5#issuecomment-99"}},
		{"query https://example.com/p?a=1&b=2%20c stays whole", []string{"https://example.com/p?a=1&b=2%20c"}},
		{"nothing here", nil},
	} {
		got := extractURLs(tc.text)
		if len(got) != len(tc.want) {
			t.Errorf("%q: got %v, want %v", tc.text, got, tc.want)
			continue
		}
		for i := range tc.want {
			if got[i] != tc.want[i] {
				t.Errorf("%q: url[%d] = %q, want %q", tc.text, i, got[i], tc.want[i])
			}
		}
	}
}

// A URL inside code is content, not a citation (#222): probe targets that
// are NXDOMAIN by design, command transcripts and error strings live in
// inline spans and fenced blocks. A Markdown link outside code, and a bare
// URL in prose, remain citations. Spans follow CommonMark: a run of n
// backticks closes only on a run of exactly n; an unclosed run is literal.
func TestExtractURLsSkipsCode(t *testing.T) {
	for _, tc := range []struct {
		name, text string
		want       []string
	}{
		{"inline span", "probed `curl -sI https://hooks.example.test/` — NXDOMAIN, as expected; see [the survey](https://github.com/o/r/issues/64#issuecomment-1)", []string{"https://github.com/o/r/issues/64#issuecomment-1"}},
		{"double-backtick span holding a backtick", "the error was ``fetch `https://zoo.example.test/` failed`` and nothing else", nil},
		{"fenced block", "before https://github.com/o/r/pull/9\n```\n$ curl https://hooks.example.test/\ncurl: (6) Could not resolve host\n```\nafter https://example.com/after", []string{"https://github.com/o/r/pull/9", "https://example.com/after"}},
		{"tilde fence with info string", "~~~sh\ncurl https://zoo.example.test/\n~~~\nsee https://github.com/o/r/issues/1", []string{"https://github.com/o/r/issues/1"}},
		{"indented fence", "  ```\n  https://example.com/in-fence\n  ```\nhttps://example.com/out", []string{"https://example.com/out"}},
		{"longer closer closes a shorter fence", "````\nhttps://example.com/a\n`````\nhttps://example.com/b", []string{"https://example.com/b"}},
		{"shorter closer does not close", "````\nhttps://example.com/a\n```\nhttps://example.com/b", nil},
		{"unclosed fence runs to the end", "```\nhttps://example.com/a\nhttps://example.com/b", nil},
		{"unclosed span is literal", "a lone ` before https://example.com/a and after", []string{"https://example.com/a"}},
		{"span across words does not fuse them", "`x`https://example.com/a`y`", []string{"https://example.com/a"}},
		{"markdown link outside code", "[the run](https://github.com/o/r/actions/runs/7) after `https://example.com/not-cited`", []string{"https://github.com/o/r/actions/runs/7"}},
	} {
		got := extractURLs(tc.text)
		if len(got) != len(tc.want) {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
			continue
		}
		for i := range tc.want {
			if got[i] != tc.want[i] {
				t.Errorf("%s: url[%d] = %q, want %q", tc.name, i, got[i], tc.want[i])
			}
		}
	}
}

func TestIsGitHubLink(t *testing.T) {
	for url, want := range map[string]bool{
		"https://github.com/o/r/issues/1":      true,
		"https://github.com/apps/some-app":     true,
		"http://github.com/o/r":                true,
		"https://docs.github.com/en/apps":      false,
		"https://github.community/x":           false,
		"https://example.com/github.com/o/r":   false,
		"https://codecrew.works/docs/protocol": false,
	} {
		if got := isGitHubLink(url); got != want {
			t.Errorf("isGitHubLink(%q) = %v, want %v", url, got, want)
		}
	}
}

// The walk over a record, with the network stubbed where it would reach
// GitHub and a real server standing in for the external host: URLs in code
// are never checked; a dead external link is a warning and the verb still
// passes; a dead github.com link is the refusal.
func TestCheckEvidenceClassifiesCitations(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/alive" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	origAPI := checkAPI
	t.Cleanup(func() { checkAPI = origAPI })
	var apiPaths []string
	checkAPI = func(path string) error {
		apiPaths = append(apiPaths, path)
		if strings.HasSuffix(path, "/issues/404") {
			return fmt.Errorf("HTTP 404")
		}
		return nil
	}

	ref := tracker.IssueRef{Repo: "o/r", Number: 5}
	hub := tracker.IssueRef{Repo: "o/r", Number: 1}
	codeOnly := []evidenceRecord{{Ref: hub, Texts: []string{
		"Baseline: `curl https://hooks.example.test/` and\n```\ncurl https://zoo.example.test/\n```\nare NXDOMAIN by design.",
	}}}
	rep := checkEvidence(codeOnly)
	if rep.Total != 0 || len(rep.Unreachable) != 0 || len(rep.Warnings) != 0 {
		t.Errorf("code-only record was scanned: %+v", rep)
	}

	external := []evidenceRecord{
		{Ref: hub, Texts: []string{"tracked in https://github.com/o/r/issues/7"}},
		{Ref: ref, Texts: []string{"see [the page](" + srv.URL + "/gone) and " + srv.URL + "/alive, plus `" + srv.URL + "/in-code`"}},
	}
	rep = checkEvidence(external)
	if rep.Total != 3 {
		t.Errorf("total = %d, want 3 (in-code URL must not count): %+v", rep.Total, rep)
	}
	if len(rep.Unreachable) != 0 {
		t.Errorf("external 404 refused: %v", rep.Unreachable)
	}
	if len(rep.Warnings) != 1 || !strings.HasPrefix(rep.Warnings[0], srv.URL+"/gone (HTTP 404) — cited on o/r#5") {
		t.Errorf("warnings = %v", rep.Warnings)
	}
	var out bytes.Buffer
	if err := reportEvidence(&out, rep, 2); err != nil {
		t.Errorf("a dead external link blocked the verb: %v", err)
	}
	if !strings.Contains(out.String(), "warning: external link does not resolve — "+srv.URL+"/gone") || !strings.Contains(out.String(), "evidence is reachable") {
		t.Errorf("report:\n%s", out.String())
	}

	github := append(external, evidenceRecord{Ref: ref, Texts: []string{"and [the gate](https://github.com/o/r/issues/404)"}})
	rep = checkEvidence(github)
	if len(rep.Unreachable) != 1 || !strings.HasPrefix(rep.Unreachable[0], "https://github.com/o/r/issues/404 (HTTP 404) — cited on o/r#5") {
		t.Errorf("unreachable = %v", rep.Unreachable)
	}
	out.Reset()
	err := reportEvidence(&out, rep, 2)
	var ref2 refusal
	if !errors.As(err, &ref2) || ref2.Code != "EVIDENCE_UNREACHABLE" {
		t.Fatalf("dead github.com link did not refuse: %v", err)
	}
	if !strings.Contains(ref2.Detail, "1 of 4 cited links") {
		t.Errorf("refusal counts the external warning as a failure or loses the total: %s", ref2.Detail)
	}
	if !strings.Contains(out.String(), "unreachable: https://github.com/o/r/issues/404") || !strings.Contains(out.String(), "warning: external link") {
		t.Errorf("report:\n%s", out.String())
	}
	for _, p := range apiPaths {
		if strings.Contains(p, "example.test") || strings.Contains(p, "in-code") {
			t.Errorf("a URL inside code reached the network: %s", p)
		}
	}
}
