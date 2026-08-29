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
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/gh"
)

// githubAPI is the REST base the mint talks to — github.com only in 1.x
// (SPEC §10); a variable so tests point it at a local server.
var githubAPI = "https://api.github.com"

// appCredential is what a mint needs: the App's private key and the id it
// signs as (the numeric App ID or the client ID — GitHub accepts either
// as the JWT issuer), plus where they came from for the stderr receipt
// and an installation hint that is never trusted blind (#119 finding 35).
type appCredential struct {
	Issuer string
	Key    []byte
	Slug   string // known only on the file path
	Source string
	Hint   string
}

// resolveCredential walks SPEC §5's order and stops at the first hit: the
// environment under the names platforms bind (`GITHUB_APP_ID` or
// `GITHUB_CLIENT_ID`; `GITHUB_PRIVATE_KEY` or `GITHUB_PEM`, as PEM text or
// a path), then the `~/.config/codecrew/` stub and key for the slug. There
// is no third path — the operator's own `gh` auth is not a mint, and a
// verb that silently handed back another principal's token would be the
// misattribution the identity tiers exist to prevent.
func resolveCredential(getenv func(string) string, configDir, slug string) (*appCredential, error) {
	issuer, issuerVar := firstEnv(getenv, "GITHUB_APP_ID", "GITHUB_CLIENT_ID")
	keyVal, keyVar := firstEnv(getenv, "GITHUB_PRIVATE_KEY", "GITHUB_PEM")
	switch {
	case issuer != "" && keyVal != "":
		key, err := keyMaterial(keyVal)
		if err != nil {
			return nil, refuse("NO_CREDENTIALS", "%s is set but is neither PEM text nor a readable key file: %v", keyVar, err)
		}
		return &appCredential{Issuer: issuer, Key: key, Source: issuerVar + " + " + keyVar, Hint: getenv("GITHUB_INSTALLATION_ID")}, nil
	case issuer != "":
		return nil, refuse("NO_CREDENTIALS", "%s is set but no private key is bound under GITHUB_PRIVATE_KEY or GITHUB_PEM (PEM text or a file path)", issuerVar)
	case keyVal != "":
		return nil, refuse("NO_CREDENTIALS", "%s is set but no App id is bound under GITHUB_APP_ID or GITHUB_CLIENT_ID", keyVar)
	}
	if slug == "" {
		return nil, refuse("NO_CREDENTIALS", "nothing to mint with: no GITHUB_APP_ID/GITHUB_CLIENT_ID + GITHUB_PRIVATE_KEY/GITHUB_PEM in the environment, and no <slug> given for the ~/.config/codecrew/ key and stub")
	}
	dir := filepath.Join(configDir, "codecrew")
	keys, _ := filepath.Glob(filepath.Join(dir, slug+".*.private-key.pem"))
	if len(keys) == 0 {
		return nil, refuse("NO_CREDENTIALS", "no private key for %q: expected %s (identity new writes it; the slug is the full App name, myorg-coder not coder)", slug, filepath.Join(dir, slug+".<date>.private-key.pem"))
	}
	sort.Strings(keys) // the date in the name sorts; the newest key is last
	key, err := os.ReadFile(keys[len(keys)-1])
	if err != nil {
		return nil, refuse("NO_CREDENTIALS", "reading %s: %v", keys[len(keys)-1], err)
	}
	stub := stubPath(configDir, slug)
	issuer, err = stubIssuer(stub)
	if err != nil {
		return nil, refuse("NO_CREDENTIALS", "key found for %q but %v — write the stub by hand: {\"slug\":%q,\"app_id\":<numeric App ID>} at %s (the App ID is on the App's settings page, or gh api /apps/%s --jq .id)", slug, err, slug, stub, slug)
	}
	return &appCredential{Issuer: issuer, Key: key, Slug: slug, Source: filepath.Base(keys[len(keys)-1]) + " + " + filepath.Base(stub), Hint: getenv("GITHUB_INSTALLATION_ID")}, nil
}

func firstEnv(getenv func(string) string, names ...string) (value, name string) {
	for _, n := range names {
		if v := strings.TrimSpace(getenv(n)); v != "" {
			return v, n
		}
	}
	return "", ""
}

// keyMaterial accepts the key as PEM text or as a path to a PEM file — a
// platform may bind either, and a path that does not exist is the error
// the caller reports, not a key.
func keyMaterial(v string) ([]byte, error) {
	if strings.Contains(v, "-----BEGIN") {
		return []byte(v), nil
	}
	data, err := os.ReadFile(v)
	if err != nil {
		return nil, err
	}
	if !bytes.Contains(data, []byte("-----BEGIN")) {
		return nil, fmt.Errorf("%s holds no PEM block", v)
	}
	return data, nil
}

// stubIssuer reads the id to sign as from the credential stub identity new
// writes beside the key: app_id first, client_id as the fallback.
func stubIssuer(path string) (string, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("no credential stub at %s", path)
	} else if err != nil {
		return "", err
	}
	var stub struct {
		AppID    json.Number `json:"app_id"`
		ClientID string      `json:"client_id"`
	}
	if err := json.Unmarshal(data, &stub); err != nil {
		return "", fmt.Errorf("stub %s is not JSON: %v", path, err)
	}
	if stub.AppID != "" {
		return stub.AppID.String(), nil
	}
	if stub.ClientID != "" {
		return stub.ClientID, nil
	}
	return "", fmt.Errorf("stub %s has neither app_id nor client_id", path)
}

// signAppJWT builds the RS256 App JWT GitHub exchanges for an installation
// token: iat a minute back for clock skew, exp nine minutes ahead (the
// ceiling is ten), iss the App ID or client ID. PKCS#1 and PKCS#8 keys both
// parse — GitHub hands out the former, platforms sometimes re-encode.
func signAppJWT(key []byte, issuer string, now time.Time) (string, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return "", fmt.Errorf("private key is not PEM")
	}
	var rsaKey *rsa.PrivateKey
	if k, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		rsaKey = k
	} else if k, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		var ok bool
		if rsaKey, ok = k.(*rsa.PrivateKey); !ok {
			return "", fmt.Errorf("private key is not RSA — GitHub App keys are RSA and the JWT is RS256")
		}
	} else {
		return "", fmt.Errorf("private key is neither PKCS#1 nor PKCS#8: %v", err)
	}
	enc := base64.RawURLEncoding
	header := enc.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	claims, _ := json.Marshal(map[string]any{"iat": now.Unix() - 60, "exp": now.Unix() + 540, "iss": issuer})
	signing := header + "." + enc.EncodeToString(claims)
	sum := sha256.Sum256([]byte(signing))
	sig, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, sum[:])
	if err != nil {
		return "", err
	}
	return signing + "." + enc.EncodeToString(sig), nil
}

// installation is the slice of GitHub's installation object the mint uses.
type installation struct {
	ID      int64  `json:"id"`
	AppSlug string `json:"app_slug"`
	Account struct {
		Login string `json:"login"`
	} `json:"account"`
}

// listInstallations asks the App itself where it is installed — the one
// lookup that needs no user credential and cannot be stale.
func listInstallations(client *http.Client, jwt string) ([]installation, error) {
	req, _ := http.NewRequest("GET", githubAPI+"/app/installations?per_page=100", nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusUnauthorized:
		return nil, refuse("BAD_CREDENTIALS", "GitHub rejected the App JWT (401): the private key and the App id do not belong to the same App, or the key was revoked — check the id against gh api /apps/<slug> --jq .id; retrying will not help")
	default:
		return nil, fmt.Errorf("GET /app/installations: HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var list []installation
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("GET /app/installations: %v", err)
	}
	return list, nil
}

// chooseInstallation picks the installation to mint against. A hint (the
// --installation flag, else GITHUB_INSTALLATION_ID) is honoured only when
// the App can see it — the platform-supplied id was stale once (#119
// finding 35) and the agent that ignored it was right; a stale hint is
// reported and discovery proceeds. With no usable hint: one installation
// is the answer, several narrow to the account that owns the working
// repo's hub, and anything still ambiguous refuses with the choices.
func chooseInstallation(list []installation, hint, owner string, notes io.Writer) (installation, error) {
	if len(list) == 0 {
		return installation{}, refuse("NO_INSTALLATION", "the App is installed on no account this key can see — install it on the account that owns the hub and its spokes (docs/identities.md step 4)")
	}
	if hint != "" {
		for _, in := range list {
			if fmt.Sprint(in.ID) == hint {
				return in, nil
			}
		}
		fmt.Fprintf(notes, "note: installation %s is not among the App's installations (stale hint) — discovering instead\n", hint)
	}
	if len(list) == 1 {
		return list[0], nil
	}
	if owner != "" {
		for _, in := range list {
			if strings.EqualFold(in.Account.Login, owner) {
				return in, nil
			}
		}
	}
	var choices []string
	for _, in := range list {
		choices = append(choices, fmt.Sprintf("%s (%d)", in.Account.Login, in.ID))
	}
	return installation{}, refuse("INSTALLATION_AMBIGUOUS", "the App is installed on %d accounts and none is the hub's owner — pass --installation <id> (or bind GITHUB_INSTALLATION_ID): %s", len(list), strings.Join(choices, ", "))
}

// mintInstallationToken exchanges the App JWT for an installation token —
// the one-hour credential every verb runs under.
func mintInstallationToken(client *http.Client, jwt string, id int64) (string, error) {
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/installations/%d/access_tokens", githubAPI, id), nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("POST /app/installations/%d/access_tokens: HTTP %d: %s", id, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var out struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(body, &out); err != nil || out.Token == "" {
		return "", fmt.Errorf("POST /app/installations/%d/access_tokens: no token in the response", id)
	}
	return out.Token, nil
}

// hubOwner is the account discovery prefers when the App is installed in
// several places: CODECREW_ORG when set, else the owner of the working
// repo's hub. Tolerant — the verb must mint from a bare shell with no
// repo at hand, so every failure here is just "no preference".
func hubOwner(getenv func(string) string) string {
	if o := strings.TrimSpace(getenv("CODECREW_ORG")); o != "" {
		return o
	}
	cfg, err := config.Load(".")
	if err != nil {
		return ""
	}
	repo := cfg.Hub
	if repo == "self" {
		if repo, err = gh.CurrentRepo(); err != nil {
			return ""
		}
	}
	owner, _, _ := strings.Cut(repo, "/")
	return owner
}

// identityToken is `gh codecrew identity token [<slug>] [--installation
// <id>]`: the token alone on stdout — `GH_TOKEN=$(gh codecrew identity
// token <slug>)` is the whole recipe — and one receipt line on stderr
// naming what was minted, which is what a dispatch's identity check reads
// (#139). It never writes gh's config: an installation token lives an hour,
// and a persisted one is the stale credential the next session inherits
// (#119 finding 10).
func identityToken(w io.Writer, args []string) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	return runIdentityToken(w, os.Stderr, os.Getenv, configDir, hubOwner(os.Getenv), http.DefaultClient, time.Now(), args)
}

func runIdentityToken(w, notes io.Writer, getenv func(string) string, configDir, owner string, client *http.Client, now time.Time, args []string) error {
	fs := flag.NewFlagSet("identity token", flag.ContinueOnError)
	fs.SetOutput(notes)
	slug, args := splitLeadingRef(args)
	hint := fs.String("installation", "", "installation id to prefer (a hint: used only when the App can see it)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if slug == "" && fs.NArg() == 1 {
		slug = fs.Arg(0)
	}
	if (slug != "" && fs.NArg() > 0) || fs.NArg() > 1 {
		return fmt.Errorf("usage: gh codecrew identity token [<slug>] [--installation <id>]")
	}
	cred, err := resolveCredential(getenv, configDir, slug)
	if err != nil {
		return err
	}
	if *hint != "" {
		cred.Hint = *hint
	}
	jwt, err := signAppJWT(cred.Key, cred.Issuer, now)
	if err != nil {
		return refuse("BAD_CREDENTIALS", "%s: %v", cred.Source, err)
	}
	list, err := listInstallations(client, jwt)
	if err != nil {
		return err
	}
	in, err := chooseInstallation(list, cred.Hint, owner, notes)
	if err != nil {
		return err
	}
	token, err := mintInstallationToken(client, jwt, in.ID)
	if err != nil {
		return err
	}
	name := in.AppSlug
	if name == "" {
		name = cred.Slug
	}
	fmt.Fprintf(notes, "minted for %s (App %s, from %s): installation %d on %s\n", name, cred.Issuer, cred.Source, in.ID, in.Account.Login)
	fmt.Fprintln(w, token)
	return nil
}
