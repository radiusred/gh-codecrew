package tracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/gh"
)

// This file is the rest of the venue seam: the reads and writes that are
// not about issues and pull requests at all — the client's own version,
// the repository gh resolves for a directory, an account's type, a team's
// members, the App-manifest exchange and the App-JWT transport — plus the
// venue's URL grammar, which is knowledge rather than I/O and so lives in
// package functions beside the interface, as NoChecksReported and
// MissingChecksPermission do. They are here because internal/cli reached
// past the interface for every one of them, and internal/gh now has
// exactly one importer: this package (M16-R1, #194).
//
// GitHub is the only implementation, by design. A second venue is backlog
// until a community asks for one (#194), so nothing here — no config key,
// no flag, no constructor — names another.

// AppAPIBase is the REST base for the App-JWT calls that cannot go through
// gh: gh authenticates as a user or an installation, and an App signing
// its own JWT is neither. A variable so tests point it at a local server.
var AppAPIBase = "https://api.github.com"

// AppCredentials is what the one-hour App-manifest conversion returns.
// client_secret and webhook_secret are printed once by the caller and
// never written to disk.
type AppCredentials struct {
	ID            int64  `json:"id"`
	Slug          string `json:"slug"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	WebhookSecret string `json:"webhook_secret"`
	PEM           string `json:"pem"`
	HTMLURL       string `json:"html_url"`
}

// Unreachable reports whether err is the venue never having been reached —
// offline, DNS-less, or unauthenticated — as opposed to it answering with
// a refusal of its own. A classifier, like NoChecksReported: the caller
// turns the answer into its own condition.
func Unreachable(err error) bool { return gh.Unreachable(err) }

// CompareVersions orders two of the venue client's x.y.z versions
// numerically: -1, 0 or 1.
func CompareVersions(a, b string) int { return gh.CompareVersions(a, b) }

// ClientVersion is the installed gh's release, which the verbs hold a
// floor on.
func (GitHub) ClientVersion() (string, error) { return gh.Version() }

// CurrentRepo returns the owner/repo of the repository in the working
// directory, per the venue client's own resolution of the origin remote.
func (GitHub) CurrentRepo() (string, error) { return gh.CurrentRepo() }

// RepoInDir is CurrentRepo for a named directory, with the default branch
// alongside: the one read that cannot use the process's own working
// directory, because init runs it against the tree it is scaffolding. Any
// answer that is not a complete one — the client failed, the output did
// not parse, either field is empty — is an error, so the caller has one
// condition to handle rather than four.
func (GitHub) RepoInDir(dir string) (repo, defaultBranch string, err error) {
	var view struct {
		NameWithOwner    string `json:"nameWithOwner"`
		DefaultBranchRef struct {
			Name string `json:"name"`
		} `json:"defaultBranchRef"`
	}
	cmd := gh.Command("gh", "repo", "view", "--json", "nameWithOwner,defaultBranchRef")
	cmd.Dir = dir
	data, err := cmd.Output()
	if err != nil {
		return "", "", err
	}
	if err := json.Unmarshal(data, &view); err != nil {
		return "", "", fmt.Errorf("gh repo view: unexpected output: %w", err)
	}
	if view.NameWithOwner == "" || view.DefaultBranchRef.Name == "" {
		return "", "", fmt.Errorf("gh repo view: no repository and default branch in %s", dir)
	}
	return view.NameWithOwner, view.DefaultBranchRef.Name, nil
}

// BranchRuleTypes lists the types of the rules that apply to a branch —
// "pull_request" among them when the branch may only be written through a
// pull request.
func (GitHub) BranchRuleTypes(repo, branch string) ([]string, error) {
	var rules []struct {
		Type string `json:"type"`
	}
	if err := gh.JSON(&rules, "api", fmt.Sprintf("repos/%s/rules/branches/%s", repo, branch)); err != nil {
		return nil, err
	}
	types := make([]string, 0, len(rules))
	for _, r := range rules {
		types = append(types, r.Type)
	}
	return types, nil
}

// TeamMembers returns a team's member logins, child-team members included,
// per the API contract.
func (GitHub) TeamMembers(org, team string) ([]string, error) {
	var members []struct {
		Login string `json:"login"`
	}
	if err := gh.JSON(&members, "api", "--paginate", fmt.Sprintf("/orgs/%s/teams/%s/members", org, team)); err != nil {
		return nil, err
	}
	logins := make([]string, 0, len(members))
	for _, m := range members {
		logins = append(logins, m.Login)
	}
	return logins, nil
}

// AccountType reports what the venue holds at users/<login> — "User" for a
// human, "Bot" for an App's account, "Organization". A login that is
// nobody is the venue's own not-found error, which the caller classifies:
// absent and unreadable are different conditions.
func (GitHub) AccountType(login string) (string, error) {
	var acct struct {
		Type string `json:"type"`
	}
	if err := gh.JSON(&acct, "api", "users/"+url.PathEscape(login)); err != nil {
		return "", err
	}
	return acct.Type, nil
}

// APIReachable resolves one API path with the caller's own credentials —
// how a record link is checked without an anonymous page fetch, private
// repositories included. The body is not read: whether the path answers is
// the whole question.
func (GitHub) APIReachable(path string) error {
	_, err := gh.Run("api", path)
	return err
}

// AppManifestConversion exchanges the temporary code the venue redirects
// back with for the created App's credentials. Unauthenticated by design
// on GitHub's side; routed through gh per the founding decision.
func (GitHub) AppManifestConversion(code string) (*AppCredentials, error) {
	var creds AppCredentials
	if err := gh.JSON(&creds, "api", "--method", "POST", "app-manifests/"+code+"/conversions"); err != nil {
		return nil, err
	}
	return &creds, nil
}

// AppRequest performs one REST call signed as the App itself. It is the
// only transport in the codebase that is not gh: an App JWT is neither a
// user nor an installation credential, and gh has no way to carry one. The
// status and the body come back unread — the mint and the webhook verbs
// each name their own conditions from them, and this method names none.
func (GitHub) AppRequest(client *http.Client, jwt, method, path string, body any) (status int, data []byte, err error) {
	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		payload = bytes.NewReader(encoded)
	}
	req, _ := http.NewRequest(method, AppAPIBase+path, payload)
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
	data, _ = io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	return resp.StatusCode, data, nil
}

// RepoURL is a repository's page — the venue's own address for an
// owner/repo.
func RepoURL(repo string) string { return "https://github.com/" + repo }

// AppInstallURL is where an App is installed on an account.
func AppInstallURL(slug string) string {
	return "https://github.com/apps/" + slug + "/installations/new"
}

// AppManifestURL is the page an App-creation manifest form posts to: the
// App-creation endpoint of the owning account, which differs for
// organizations and personal accounts.
func AppManifestURL(owner, ownerType string) string {
	if ownerType == "Organization" {
		return "https://github.com/organizations/" + owner + "/settings/apps/new"
	}
	return "https://github.com/settings/apps/new"
}

// AppSettingsURL is an App's settings page — where the steps the API
// cannot take (the avatar, the event subscriptions) are taken by hand.
func AppSettingsURL(owner, ownerType, slug string) string {
	if ownerType == "Organization" {
		return "https://github.com/organizations/" + owner + "/settings/apps/" + slug
	}
	return "https://github.com/settings/apps/" + slug
}

// IsRecordLink says whether a citation points at the venue — the record
// itself: issues, pull requests, comments, commits, blobs, runs, App
// pages. Those are what QA tests against, so one that does not resolve is
// a refusal; any other host is external content.
func IsRecordLink(link string) bool {
	return strings.HasPrefix(link, "https://github.com/") || strings.HasPrefix(link, "http://github.com/")
}

// RecordAPIPath maps a record link — issue, pull request, comment anchor,
// commit, blob — to the API path that answers for it, so reachability is
// checked with the caller's auth instead of an anonymous page fetch. A
// venue knows its own URL grammar; links this one cannot map are checked
// as plain HTTP by the caller.
func RecordAPIPath(link string) (string, bool) {
	rest, ok := strings.CutPrefix(link, "https://github.com/")
	if !ok {
		return "", false
	}
	rest, _, _ = strings.Cut(rest, "#") // comment anchors resolve via the issue
	parts := strings.Split(rest, "/")
	if len(parts) < 4 {
		return "", false
	}
	owner, repo, kind := parts[0], parts[1], parts[2]
	tail := parts[3:]
	switch kind {
	case "issues", "pull":
		return fmt.Sprintf("repos/%s/%s/issues/%s", owner, repo, tail[0]), true
	case "commit":
		return fmt.Sprintf("repos/%s/%s/commits/%s", owner, repo, tail[0]), true
	case "blob", "tree":
		if len(tail) < 2 {
			return "", false
		}
		ref, path := tail[0], strings.Join(tail[1:], "/")
		return fmt.Sprintf("repos/%s/%s/contents/%s?ref=%s", owner, repo, path, ref), true
	}
	return "", false
}
