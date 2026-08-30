package cli

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
		{"in `https://github.com/o/r/pull/2` and \"https://example.com/a\" or 'https://example.com/b'", []string{"https://github.com/o/r/pull/2", "https://example.com/a", "https://example.com/b"}},
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
