package cli

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/radiusred/gh-codecrew/internal/gh"
	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// urlPattern matches https?:// links in issue bodies and comments;
// extractURLs trims the markdown and punctuation that ride along.
var urlPattern = regexp.MustCompile(`https?://[^\s<>"'` + "`" + `)\]]+`)

// extractURLs pulls the deduplicated, ordered links from a piece of record
// text. Trailing punctuation that prose attaches (., ,, ;, :) is trimmed.
func extractURLs(text string) []string {
	seen := map[string]bool{}
	var urls []string
	for _, u := range urlPattern.FindAllString(text, -1) {
		u = strings.TrimRight(u, ".,;:")
		if !seen[u] {
			seen[u] = true
			urls = append(urls, u)
		}
	}
	return urls
}

// githubAPIPath maps a github.com record link — issue, PR, comment anchor,
// commit, blob — to its API path, so reachability is checked with the
// caller's auth (private repos included) instead of an anonymous page
// fetch. Links it cannot map check as plain HTTP.
func githubAPIPath(url string) (string, bool) {
	rest, ok := strings.CutPrefix(url, "https://github.com/")
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

// checkAPI and checkHTTP are stubbable (the convertManifest pattern) so the
// verb's walk is tested without a network.
var checkAPI = func(path string) error {
	_, err := gh.Run("api", path)
	return err
}

var checkHTTP = func(url string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}

// checkURL resolves one link by the right mechanism.
func checkURL(url string) error {
	if path, ok := githubAPIPath(url); ok {
		return checkAPI(path)
	}
	return checkHTTP(url)
}

// milestoneEvidence walks a milestone's record — the tracking issue and
// every sub-issue, bodies and comments — and verifies every cited link
// resolves. The M4-R4 lesson made deterministic: a record that exists only
// in a working tree is not evidence, and QA must not be dispatched against
// citations that 404 (findings 5 and 7 on #73 repeated it twice).
func milestoneEvidence(w io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: codecrew milestone evidence <n>")
	}
	c, err := load()
	if err != nil {
		return err
	}
	milestones, err := c.t.OpenMilestones(c.hub)
	if err != nil {
		return err
	}
	var milestone *tracker.Milestone
	for i := range milestones {
		if got, ok := tracker.MilestoneNumber(milestones[i].Title); ok && fmt.Sprint(got) == args[0] {
			milestone = &milestones[i]
			break
		}
	}
	if milestone == nil {
		return refuse("NOT_FOUND", "no open milestone M%s in %s", args[0], c.hub)
	}

	refs := []tracker.IssueRef{milestone.Ref}
	subs, err := c.t.SubIssues(milestone.Ref)
	if err != nil {
		return err
	}
	refs = append(refs, subs...)

	seen := map[string]bool{}
	var unreachable []string
	total := 0
	for _, ref := range refs {
		var texts []string
		if body, err := c.t.IssueBody(ref); err == nil {
			texts = append(texts, body)
		}
		if comments, err := c.t.Comments(ref); err == nil {
			for _, cm := range comments {
				texts = append(texts, cm.Body)
			}
		}
		for _, text := range texts {
			for _, u := range extractURLs(text) {
				if seen[u] {
					continue
				}
				seen[u] = true
				total++
				if err := checkURL(u); err != nil {
					unreachable = append(unreachable, fmt.Sprintf("%s (%v) — cited on %s", u, err, ref))
				}
			}
		}
	}
	if len(unreachable) > 0 {
		sort.Strings(unreachable)
		for _, u := range unreachable {
			fmt.Fprintf(w, "unreachable: %s\n", u)
		}
		return refuse("EVIDENCE_UNREACHABLE", "%d of %d cited links do not resolve — commit or repair the evidence before dispatching QA", len(unreachable), total)
	}
	fmt.Fprintf(w, "all %d cited links resolve across %d issues — evidence is reachable\n", total, len(refs))
	return nil
}
