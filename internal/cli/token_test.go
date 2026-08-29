package cli

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// testKey is one RSA key for the whole file — generation is the slow part.
var testKey, _ = rsa.GenerateKey(rand.Reader, 2048)

func pkcs1PEM(k *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)})
}

func pkcs8PEM(t *testing.T, k *rsa.PrivateKey) []byte {
	der, err := x509.MarshalPKCS8PrivateKey(k)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
}

func envOf(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

// verifyJWT checks an RS256 JWT against the test key's public half and
// returns its claims — the assertion that the mint signs what GitHub
// verifies, not merely something JWT-shaped.
func verifyJWT(t *testing.T, tok string) map[string]any {
	t.Helper()
	parts := strings.Split(tok, ".")
	if len(parts) != 3 {
		t.Fatalf("JWT has %d parts", len(parts))
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	if err := rsa.VerifyPKCS1v15(&testKey.PublicKey, crypto.SHA256, sum[:], sig); err != nil {
		t.Fatalf("signature does not verify: %v", err)
	}
	header, _ := base64.RawURLEncoding.DecodeString(parts[0])
	if string(header) != `{"alg":"RS256","typ":"JWT"}` {
		t.Errorf("header = %s", header)
	}
	claimsJSON, _ := base64.RawURLEncoding.DecodeString(parts[1])
	var claims map[string]any
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		t.Fatal(err)
	}
	return claims
}

func TestSignAppJWT(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	for name, key := range map[string][]byte{"pkcs1": pkcs1PEM(testKey), "pkcs8": pkcs8PEM(t, testKey)} {
		tok, err := signAppJWT(key, "4744165", now)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		claims := verifyJWT(t, tok)
		if claims["iss"] != "4744165" || claims["iat"] != float64(now.Unix()-60) || claims["exp"] != float64(now.Unix()+540) {
			t.Errorf("%s claims = %v", name, claims)
		}
	}
	if _, err := signAppJWT([]byte("not a key"), "1", now); err == nil {
		t.Error("non-PEM key signed")
	}
}

func TestResolveCredentialEnvFirst(t *testing.T) {
	dir := t.TempDir()
	keyFile := filepath.Join(dir, "key.pem")
	os.WriteFile(keyFile, pkcs1PEM(testKey), 0o600)
	cases := []struct {
		name   string
		env    map[string]string
		issuer string
		source string
	}{
		{"app id + PEM text", map[string]string{"GITHUB_APP_ID": "11", "GITHUB_PRIVATE_KEY": string(pkcs1PEM(testKey))}, "11", "GITHUB_APP_ID + GITHUB_PRIVATE_KEY"},
		{"client id + PEM path", map[string]string{"GITHUB_CLIENT_ID": "Iv1.abc", "GITHUB_PEM": keyFile}, "Iv1.abc", "GITHUB_CLIENT_ID + GITHUB_PEM"},
		{"app id wins over client id", map[string]string{"GITHUB_APP_ID": "11", "GITHUB_CLIENT_ID": "Iv1.abc", "GITHUB_PEM": keyFile}, "11", "GITHUB_APP_ID + GITHUB_PEM"},
	}
	for _, c := range cases {
		// A stub and key exist for the slug too: the environment must win.
		cred, err := resolveCredential(envOf(c.env), dir, "myorg-coder")
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if cred.Issuer != c.issuer || cred.Source != c.source || !bytes.Contains(cred.Key, []byte("-----BEGIN")) {
			t.Errorf("%s: issuer %q source %q", c.name, cred.Issuer, cred.Source)
		}
	}
	// Hint rides along from the environment.
	cred, _ := resolveCredential(envOf(map[string]string{"GITHUB_APP_ID": "11", "GITHUB_PEM": keyFile, "GITHUB_INSTALLATION_ID": "99"}), dir, "")
	if cred.Hint != "99" {
		t.Errorf("hint = %q", cred.Hint)
	}
}

func TestResolveCredentialFiles(t *testing.T) {
	dir := t.TempDir()
	cc := filepath.Join(dir, "codecrew")
	os.MkdirAll(cc, 0o755)
	os.WriteFile(filepath.Join(cc, "myorg-coder.2026-01-01.private-key.pem"), []byte("-----BEGIN OLD-----\n"), 0o600)
	os.WriteFile(filepath.Join(cc, "myorg-coder.2026-08-30.private-key.pem"), pkcs1PEM(testKey), 0o600)
	os.WriteFile(filepath.Join(cc, "myorg-coder.json"), []byte(`{"slug":"myorg-coder","app_id":3163997,"client_id":"Iv1.x"}`), 0o644)
	cred, err := resolveCredential(envOf(nil), dir, "myorg-coder")
	if err != nil {
		t.Fatal(err)
	}
	if cred.Issuer != "3163997" || cred.Slug != "myorg-coder" || !strings.Contains(cred.Source, "2026-08-30") {
		t.Errorf("issuer %q slug %q source %q — want the newest key and the numeric app_id", cred.Issuer, cred.Slug, cred.Source)
	}
	// client_id is the fallback issuer when the stub carries no app_id.
	os.WriteFile(filepath.Join(cc, "myorg-coder.json"), []byte(`{"slug":"myorg-coder","client_id":"Iv1.x"}`), 0o644)
	if cred, _ := resolveCredential(envOf(nil), dir, "myorg-coder"); cred.Issuer != "Iv1.x" {
		t.Errorf("client_id fallback: issuer = %q", cred.Issuer)
	}
}

func TestResolveCredentialRefusals(t *testing.T) {
	dir := t.TempDir()
	cc := filepath.Join(dir, "codecrew")
	os.MkdirAll(cc, 0o755)
	os.WriteFile(filepath.Join(cc, "stubless.2026-08-30.private-key.pem"), pkcs1PEM(testKey), 0o600)
	cases := []struct {
		name   string
		env    map[string]string
		slug   string
		detail string
	}{
		{"nothing at all", nil, "", "no <slug>"},
		{"id without key", map[string]string{"GITHUB_APP_ID": "1"}, "", "no private key is bound"},
		{"key without id", map[string]string{"GITHUB_PEM": "x"}, "", "no App id is bound"},
		{"key path missing", map[string]string{"GITHUB_APP_ID": "1", "GITHUB_PEM": filepath.Join(dir, "nope.pem")}, "", "neither PEM text nor a readable key file"},
		{"unknown slug", nil, "myorg-ghost", "no private key for"},
		{"key but no stub", nil, "stubless", "write the stub by hand"},
	}
	for _, c := range cases {
		_, err := resolveCredential(envOf(c.env), dir, c.slug)
		var r refusal
		if !errors.As(err, &r) || r.Code != "NO_CREDENTIALS" || !strings.Contains(r.Detail, c.detail) {
			t.Errorf("%s: err = %v, want NO_CREDENTIALS containing %q", c.name, err, c.detail)
		}
	}
}

// fakeGitHub serves the two endpoints the mint uses, verifying the App
// JWT on each call the way GitHub does.
func fakeGitHub(t *testing.T, list []installation, tokenFor map[int64]string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if auth == "" || auth == "bad" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"message":"A JSON web token could not be decoded"}`)
			return
		}
		claims := verifyJWT(t, auth)
		if claims["iss"] == "wrong" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch {
		case r.Method == "GET" && r.URL.Path == "/app/installations":
			json.NewEncoder(w).Encode(list)
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/access_tokens"):
			var id int64
			fmt.Sscanf(r.URL.Path, "/app/installations/%d/access_tokens", &id)
			tok, ok := tokenFor[id]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"token":%q,"expires_at":"2026-08-30T01:00:00Z"}`, tok)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func inst(id int64, login, slug string) installation {
	var in installation
	in.ID, in.AppSlug = id, slug
	in.Account.Login = login
	return in
}

func mint(t *testing.T, srv *httptest.Server, env map[string]string, configDir, owner string, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	prev := githubAPI
	githubAPI = srv.URL
	defer func() { githubAPI = prev }()
	var out, notes bytes.Buffer
	err = runIdentityToken(&out, &notes, envOf(env), configDir, owner, srv.Client(), time.Now(), args)
	return out.String(), notes.String(), err
}

func TestIdentityTokenMintsAndPrintsOnlyTheToken(t *testing.T) {
	srv := fakeGitHub(t, []installation{inst(156598567, "radiusred", "radiusred-cody")}, map[int64]string{156598567: "ghs_abc"})
	env := map[string]string{"GITHUB_APP_ID": "3163997", "GITHUB_PRIVATE_KEY": string(pkcs1PEM(testKey))}
	out, notes, err := mint(t, srv, env, t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	if out != "ghs_abc\n" {
		t.Errorf("stdout = %q — the token and nothing else", out)
	}
	for _, want := range []string{"minted for radiusred-cody", "App 3163997", "GITHUB_APP_ID + GITHUB_PRIVATE_KEY", "installation 156598567 on radiusred"} {
		if !strings.Contains(notes, want) {
			t.Errorf("stderr receipt lacks %q: %q", want, notes)
		}
	}
}

func TestIdentityTokenDiscovery(t *testing.T) {
	two := []installation{inst(1, "radiusred", "myorg-coder"), inst(2, "davison", "myorg-coder")}
	tokens := map[int64]string{1: "t1", 2: "t2"}
	env := func(extra map[string]string) map[string]string {
		m := map[string]string{"GITHUB_APP_ID": "7", "GITHUB_PEM": string(pkcs1PEM(testKey))}
		for k, v := range extra {
			m[k] = v
		}
		return m
	}
	// A hint the App can see is honoured.
	out, _, err := mint(t, fakeGitHub(t, two, tokens), env(map[string]string{"GITHUB_INSTALLATION_ID": "2"}), t.TempDir(), "radiusred")
	if err != nil || out != "t2\n" {
		t.Errorf("hint: out %q err %v", out, err)
	}
	// The flag outranks the environment's hint.
	out, _, err = mint(t, fakeGitHub(t, two, tokens), env(map[string]string{"GITHUB_INSTALLATION_ID": "2"}), t.TempDir(), "", "--installation", "1")
	if err != nil || out != "t1\n" {
		t.Errorf("flag: out %q err %v", out, err)
	}
	// A stale hint is reported and discovery proceeds (#119 finding 35).
	out, notes, err := mint(t, fakeGitHub(t, two, tokens), env(map[string]string{"GITHUB_INSTALLATION_ID": "157121314"}), t.TempDir(), "davison")
	if err != nil || out != "t2\n" || !strings.Contains(notes, "stale hint") {
		t.Errorf("stale hint: out %q notes %q err %v", out, notes, err)
	}
	// One installation needs no owner.
	out, _, err = mint(t, fakeGitHub(t, two[:1], tokens), env(nil), t.TempDir(), "")
	if err != nil || out != "t1\n" {
		t.Errorf("single: out %q err %v", out, err)
	}
	// Several narrow to the hub's owner, case-insensitively.
	out, _, err = mint(t, fakeGitHub(t, two, tokens), env(nil), t.TempDir(), "RadiusRed")
	if err != nil || out != "t1\n" {
		t.Errorf("owner: out %q err %v", out, err)
	}
	// Several, no owner: refuse with the choices, minting nothing.
	out, _, err = mint(t, fakeGitHub(t, two, tokens), env(nil), t.TempDir(), "")
	var r refusal
	if !errors.As(err, &r) || r.Code != "INSTALLATION_AMBIGUOUS" || !strings.Contains(r.Detail, "davison (2)") || out != "" {
		t.Errorf("ambiguous: out %q err %v", out, err)
	}
	// None: NO_INSTALLATION.
	_, _, err = mint(t, fakeGitHub(t, nil, tokens), env(nil), t.TempDir(), "radiusred")
	if !errors.As(err, &r) || r.Code != "NO_INSTALLATION" {
		t.Errorf("none: err %v", err)
	}
	// A key GitHub rejects: BAD_CREDENTIALS, no retry advice.
	_, _, err = mint(t, fakeGitHub(t, two, tokens), env(map[string]string{"GITHUB_APP_ID": "wrong"}), t.TempDir(), "radiusred")
	if !errors.As(err, &r) || r.Code != "BAD_CREDENTIALS" {
		t.Errorf("rejected JWT: err %v", err)
	}
	// Too many positional arguments is a usage error, not a mint.
	if _, _, err := mint(t, fakeGitHub(t, two, tokens), env(nil), t.TempDir(), "radiusred", "a", "b"); err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("two slugs: err %v", err)
	}
}

func TestIdentityTokenFromStub(t *testing.T) {
	dir := t.TempDir()
	cc := filepath.Join(dir, "codecrew")
	os.MkdirAll(cc, 0o755)
	os.WriteFile(filepath.Join(cc, "myorg-coder.2026-08-30.private-key.pem"), pkcs8PEM(t, testKey), 0o600)
	os.WriteFile(filepath.Join(cc, "myorg-coder.json"), []byte(`{"slug":"myorg-coder","app_id":7}`), 0o644)
	srv := fakeGitHub(t, []installation{inst(5, "myorg", "")}, map[int64]string{5: "ghs_stub"})
	out, notes, err := mint(t, srv, nil, dir, "", "myorg-coder")
	if err != nil || out != "ghs_stub\n" {
		t.Fatalf("out %q err %v", out, err)
	}
	// With no app_slug in the response the receipt falls back to the slug asked for.
	if !strings.Contains(notes, "minted for myorg-coder (App 7, from myorg-coder.2026-08-30.private-key.pem + myorg-coder.json)") {
		t.Errorf("receipt = %q", notes)
	}
}
