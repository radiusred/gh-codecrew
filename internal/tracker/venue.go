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
func (g GitHub) BranchRuleTypes(repo, branch string) ([]string, error) {
	rules, err := g.BranchRules(repo, branch)
	if err != nil {
		return nil, err
	}
	types := make([]string, 0, len(rules))
	for _, r := range rules {
		types = append(types, r.Type)
	}
	return types, nil
}

// BranchRules reads the rules in force on a branch — every active ruleset
// that targets it, the organization's included — through the endpoint a
// read-only token may call. A plan without rulesets answers it with an
// error of its own (404, or 403 on a private repository), which the caller
// treats as no policy to report on.
func (GitHub) BranchRules(repo, branch string) ([]BranchRule, error) {
	var raw []struct {
		Type       string `json:"type"`
		Parameters struct {
			RequiredStatusChecks []struct {
				Context string `json:"context"`
			} `json:"required_status_checks"`
		} `json:"parameters"`
	}
	if err := gh.JSON(&raw, "api", "--paginate", fmt.Sprintf("repos/%s/rules/branches/%s", repo, branch)); err != nil {
		return nil, err
	}
	rules := make([]BranchRule, 0, len(raw))
	for _, r := range raw {
		rule := BranchRule{Type: r.Type}
		for _, c := range r.Parameters.RequiredStatusChecks {
			rule.Checks = append(rule.Checks, c.Context)
		}
		rules = append(rules, rule)
	}
	return rules, nil
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

// recordShape is one row of the venue's URL grammar: a github.com page the
// record cites, and the REST path that answers for it under the caller's
// own auth. Path is the URL's path after the host, by segment: a literal
// matches itself, {x} one segment, and a final {x*} the one or more that
// remain. Anchor, when set, is the fragment the row requires — a literal
// prefix and one {x} capture — and a row without one matches whatever
// fragment the link carries, or none. API is the REST path with the
// captures put back.
type recordShape struct {
	Path, Anchor, API string
}

// recordShapes is every github.com shape the fleet cites, mapped onto the
// API (#388, M21-R3): first match wins, so an anchored row comes before
// the plain row for the same page. A private repository answers 404 to an
// anonymous fetch of any of its pages, so a shape left off this table is
// one the verbs cannot verify there; adding the row is the fix.
var recordShapes = []recordShape{
	{"{o}/{r}", "", "repos/{o}/{r}"},
	{"{o}/{r}/issues/{n}", "issuecomment-{id}", "repos/{o}/{r}/issues/comments/{id}"},
	{"{o}/{r}/issues/{n}", "", "repos/{o}/{r}/issues/{n}"},
	{"{o}/{r}/pull/{n}", "issuecomment-{id}", "repos/{o}/{r}/issues/comments/{id}"},
	{"{o}/{r}/pull/{n}", "pullrequestreview-{id}", "repos/{o}/{r}/pulls/{n}/reviews/{id}"},
	{"{o}/{r}/pull/{n}", "discussion_r{id}", "repos/{o}/{r}/pulls/comments/{id}"},
	{"{o}/{r}/pull/{n}/{tab}", "discussion_r{id}", "repos/{o}/{r}/pulls/comments/{id}"},
	{"{o}/{r}/pull/{n}/{tab}", "r{id}", "repos/{o}/{r}/pulls/comments/{id}"},
	{"{o}/{r}/pull/{n}/commits/{sha}", "", "repos/{o}/{r}/commits/{sha}"},
	{"{o}/{r}/pull/{n}", "", "repos/{o}/{r}/issues/{n}"},
	{"{o}/{r}/pull/{n}/{tab*}", "", "repos/{o}/{r}/issues/{n}"},
	{"{o}/{r}/commit/{sha}", "", "repos/{o}/{r}/commits/{sha}"},
	{"{o}/{r}/commits/{ref}", "", "repos/{o}/{r}/commits/{ref}"},
	{"{o}/{r}/tree/{ref}", "", "repos/{o}/{r}/commits/{ref}"},
	{"{o}/{r}/tree/{ref}/{path*}", "", "repos/{o}/{r}/contents/{path}?ref={ref}"},
	{"{o}/{r}/blob/{ref}/{path*}", "", "repos/{o}/{r}/contents/{path}?ref={ref}"},
	{"{o}/{r}/compare/{basehead*}", "", "repos/{o}/{r}/compare/{basehead}"},
	{"{o}/{r}/releases", "", "repos/{o}/{r}/releases"},
	{"{o}/{r}/releases/latest", "", "repos/{o}/{r}/releases/latest"},
	{"{o}/{r}/releases/tag/{tag}", "", "repos/{o}/{r}/releases/tags/{tag}"},
	{"{o}/{r}/actions/runs/{run}", "", "repos/{o}/{r}/actions/runs/{run}"},
	{"{o}/{r}/actions/runs/{run}/job/{job}", "", "repos/{o}/{r}/actions/jobs/{job}"},
	{"{o}/{r}/actions/runs/{run}/jobs/{job}", "", "repos/{o}/{r}/actions/jobs/{job}"},
	{"{o}/{r}/actions/runs/{run}/attempts/{attempt}", "", "repos/{o}/{r}/actions/runs/{run}/attempts/{attempt}"},
	{"{o}/{r}/actions/workflows/{file}", "", "repos/{o}/{r}/actions/workflows/{file}"},
	{"{o}/{r}/milestone/{n}", "", "repos/{o}/{r}/milestones/{n}"},
}

// numericCaptures are the captures that are only ever a number on the
// venue: issues/new is a form, not issue "new".
var numericCaptures = map[string]bool{"n": true, "id": true, "run": true, "job": true, "attempt": true}

// reservedOwners are github.com's own top-level pages: a link under one of
// them is an App, an organization's settings, a user list — never a
// repository, so nothing under it maps to repos/.
var reservedOwners = map[string]bool{
	"about": true, "apps": true, "codespaces": true, "collections": true,
	"contact": true, "customer-stories": true, "dashboard": true,
	"enterprise": true, "enterprises": true, "events": true, "explore": true,
	"features": true, "issues": true, "login": true, "logout": true,
	"marketplace": true, "new": true, "notifications": true, "orgs": true,
	"organizations": true, "pricing": true, "pulls": true, "search": true,
	"security": true, "sessions": true, "settings": true, "site": true,
	"sponsors": true, "stars": true, "team": true, "topics": true,
	"trending": true, "users": true, "watching": true,
}

// splitRecordLink takes a github.com link apart: the path segments, a
// repository's `.git` suffix and any trailing slash dropped, and the
// fragment. A query (`?pr=12`, `?tab=readme`) is the page's view, not
// what it is, and is dropped. ok is false for any other host.
func splitRecordLink(link string) (segs []string, fragment string, ok bool) {
	rest, ok := strings.CutPrefix(link, "https://github.com/")
	if !ok {
		if rest, ok = strings.CutPrefix(link, "http://github.com/"); !ok {
			return nil, "", false
		}
	}
	rest, fragment, _ = strings.Cut(rest, "#")
	rest, _, _ = strings.Cut(rest, "?")
	rest = strings.TrimRight(rest, "/")
	if rest == "" {
		return nil, fragment, true
	}
	segs = strings.Split(rest, "/")
	if len(segs) >= 2 {
		segs[1] = strings.TrimSuffix(segs[1], ".git")
	}
	return segs, fragment, true
}

// match fits a link's segments and fragment to one row, returning the
// captures.
func (s recordShape) match(segs []string, fragment string) (map[string]string, bool) {
	caps := map[string]string{}
	pat := strings.Split(s.Path, "/")
	for i, p := range pat {
		if name, ok := strings.CutPrefix(p, "{"); ok {
			name = strings.TrimSuffix(name, "}")
			if rest, ok := strings.CutSuffix(name, "*"); ok && i == len(pat)-1 {
				if len(segs) <= i {
					return nil, false
				}
				caps[rest] = strings.Join(segs[i:], "/")
				return caps, s.capturesValid(caps) && s.matchAnchor(fragment, caps)
			}
			if i >= len(segs) || segs[i] == "" {
				return nil, false
			}
			caps[name] = segs[i]
			continue
		}
		if i >= len(segs) || segs[i] != p {
			return nil, false
		}
	}
	if len(segs) != len(pat) {
		return nil, false
	}
	return caps, s.capturesValid(caps) && s.matchAnchor(fragment, caps)
}

func (s recordShape) matchAnchor(fragment string, caps map[string]string) bool {
	if s.Anchor == "" {
		return true
	}
	prefix, name, _ := strings.Cut(s.Anchor, "{")
	name = strings.TrimSuffix(name, "}")
	id, ok := strings.CutPrefix(fragment, prefix)
	if !ok || !isDigits(id) {
		return false
	}
	caps[name] = id
	return true
}

func (s recordShape) capturesValid(caps map[string]string) bool {
	for name, v := range caps {
		if numericCaptures[name] && !isDigits(v) {
			return false
		}
	}
	if bh, ok := caps["basehead"]; ok && !strings.Contains(bh, "..") {
		return false
	}
	return true
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// RecordAPIPath maps a record link to the API path that answers for it —
// recordShapes, first match — so reachability is checked with the
// caller's auth instead of an anonymous page fetch. A venue knows its own
// URL grammar; a link no row maps is the caller's to classify: RecordRepo
// says which repository it lies under, if any.
func RecordAPIPath(link string) (string, bool) {
	_, path, ok := resolveShape(link)
	return path, ok
}

// resolveShape is RecordAPIPath with the row that matched, so the test
// can hold every row of the table to a link.
func resolveShape(link string) (row int, path string, ok bool) {
	segs, fragment, ok := splitRecordLink(link)
	if !ok || len(segs) < 2 || reservedOwners[segs[0]] {
		return -1, "", false
	}
	for i, s := range recordShapes {
		caps, ok := s.match(segs, fragment)
		if !ok {
			continue
		}
		path := s.API
		for name, v := range caps {
			path = strings.ReplaceAll(path, "{"+name+"}", v)
		}
		return i, path, true
	}
	return -1, "", false
}

// RecordRepo names the repository a github.com link lies under, and the
// link's shape there — the first path segment after the repository, the
// part of the address that decides what kind of page it is. ok is false
// for another host and for github.com's own pages (an App, a user, an
// organization's listing), which lie under no repository.
func RecordRepo(link string) (repo, shape string, ok bool) {
	segs, _, ok := splitRecordLink(link)
	if !ok || len(segs) < 2 || reservedOwners[segs[0]] {
		return "", "", false
	}
	if len(segs) > 2 {
		shape = segs[2]
	}
	return segs[0] + "/" + segs[1], shape, true
}
