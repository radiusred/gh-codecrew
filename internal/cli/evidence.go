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

// urlPattern matches https?:// links in issue bodies and comments: a URL
// ends at the first character that cannot be in one — whitespace, quotes,
// backticks, brackets, and anything outside ASCII (a prose ellipsis glued
// to a link once made milestone evidence refuse a real citation, #138).
// Parentheses are matched so a Wikipedia-style path survives; extractURLs
// trims the markdown and punctuation that ride along.
var urlPattern = regexp.MustCompile(`https?://[A-Za-z0-9\-._~:/?#@!$&+,;=%()]+`)

// extractURLs pulls the deduplicated, ordered citations from a piece of
// record text. A citation is a URL in prose or in a Markdown link outside
// code; a URL inside an inline code span or a fenced code block is content
// — a probe target meant to be unreachable, a command transcript, an error
// string — and is not a citation (#222: two NXDOMAIN-by-design hostnames
// in a survey comment refused a complete record). Trailing punctuation
// that prose attaches (., ,, ;, :) is trimmed, and a closing parenthesis
// is kept only when the URL opened it — a markdown link's `)` and a
// sentence's `)` are not part of the address.
func extractURLs(text string) []string {
	seen := map[string]bool{}
	var urls []string
	for _, u := range urlPattern.FindAllString(stripCode(text), -1) {
		u = trimURL(u)
		if u != "" && !seen[u] {
			seen[u] = true
			urls = append(urls, u)
		}
	}
	return urls
}

// stripCode blanks Markdown code out of record text: fenced blocks (a line
// opening with three or more backticks or tildes, closed by a fence of the
// same character at least as long) and inline spans (a backtick run closed
// by a run of the same length, as CommonMark reads them). An unclosed fence
// runs to the end of the text; an unclosed backtick run is literal text.
// Replaced with a space so words on either side do not fuse.
func stripCode(text string) string {
	var out strings.Builder
	var fence string // the opening fence of the block being skipped
	for _, line := range strings.SplitAfter(text, "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		if fence != "" {
			if strings.HasPrefix(trimmed, fence) && strings.Trim(trimmed, fence[:1]+" \t\r\n") == "" {
				fence = ""
			}
			continue
		}
		if f := fenceOpener(trimmed); f != "" {
			fence = f
			continue
		}
		out.WriteString(stripSpans(line))
	}
	return out.String()
}

// fenceOpener returns the fence run that opens a code block on this line
// — three or more backticks or tildes at its start, an info string such as
// a language tag allowed after them — or "" if none. A fence may open on
// the same line as a list marker (`- ```sh`): the item's content starts
// after the marker, so the marker is skipped first.
func fenceOpener(line string) string {
	line = strings.TrimLeft(skipListMarker(line), " \t")
	if line == "" || (line[0] != '`' && line[0] != '~') {
		return ""
	}
	n := 0
	for n < len(line) && line[n] == line[0] {
		n++
	}
	if n < 3 {
		return ""
	}
	return line[:n]
}

// skipListMarker drops a leading bullet (-, *, +) or ordered marker (1.,
// 1)) and the space after it; a line with no marker is returned as is.
func skipListMarker(line string) string {
	rest := line
	if len(rest) > 0 && strings.ContainsRune("-*+", rune(rest[0])) {
		rest = rest[1:]
	} else {
		n := 0
		for n < len(rest) && rest[n] >= '0' && rest[n] <= '9' {
			n++
		}
		if n == 0 || n > 9 || n >= len(rest) || (rest[n] != '.' && rest[n] != ')') {
			return line
		}
		rest = rest[n+1:]
	}
	if rest == "" || (rest[0] != ' ' && rest[0] != '\t') {
		return line
	}
	return rest
}

// stripSpans replaces every closed inline code span on one line with a
// space; a backtick run with no matching closer stays as it is.
func stripSpans(line string) string {
	var out strings.Builder
	for i := 0; i < len(line); {
		if line[i] != '`' {
			out.WriteByte(line[i])
			i++
			continue
		}
		n := 0
		for i+n < len(line) && line[i+n] == '`' {
			n++
		}
		run := line[i : i+n]
		end := -1
		for j := i + n; j < len(line); {
			k := strings.Index(line[j:], run)
			if k < 0 {
				break
			}
			j += k
			m := j
			for m < len(line) && line[m] == '`' {
				m++
			}
			if m-j == n {
				end = j
				break
			}
			j = m
		}
		if end < 0 {
			out.WriteString(run)
			i += n
			continue
		}
		out.WriteByte(' ')
		i = end + n
	}
	return out.String()
}

func trimURL(u string) string {
	for {
		trimmed := strings.TrimRight(u, ".,;:")
		if strings.HasSuffix(trimmed, ")") && strings.Count(trimmed, ")") > strings.Count(trimmed, "(") {
			trimmed = strings.TrimSuffix(trimmed, ")")
		}
		if trimmed == u {
			return u
		}
		u = trimmed
	}
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

// isGitHubLink says whether a citation points at github.com — the record
// itself: issues, PRs, comments, commits, blobs, runs, App pages. Those
// are what QA tests against, so one that does not resolve refuses; any
// other host is external content whose death is reported as a warning
// for the qa seat to weigh (#222, the third shape).
func isGitHubLink(url string) bool {
	return strings.HasPrefix(url, "https://github.com/") || strings.HasPrefix(url, "http://github.com/")
}

// evidenceRecord is one issue's texts — body first, then comments — as
// the walk reads them.
type evidenceRecord struct {
	Ref   tracker.IssueRef
	Texts []string
}

// evidenceReport is what checkEvidence found: how many distinct citations
// the record makes, the github.com ones that do not resolve (a refusal),
// and the external ones that do not resolve (warnings).
type evidenceReport struct {
	Total       int
	Unreachable []string
	Warnings    []string
}

// checkEvidence resolves every citation across the records once, in the
// order first cited, and sorts the two failure lists for a stable report.
func checkEvidence(records []evidenceRecord) evidenceReport {
	var rep evidenceReport
	seen := map[string]bool{}
	for _, rec := range records {
		for _, text := range rec.Texts {
			for _, u := range extractURLs(text) {
				if seen[u] {
					continue
				}
				seen[u] = true
				rep.Total++
				err := checkURL(u)
				if err == nil {
					continue
				}
				line := fmt.Sprintf("%s (%v) — cited on %s", u, err, rec.Ref)
				if isGitHubLink(u) {
					rep.Unreachable = append(rep.Unreachable, line)
				} else {
					rep.Warnings = append(rep.Warnings, line)
				}
			}
		}
	}
	sort.Strings(rep.Unreachable)
	sort.Strings(rep.Warnings)
	return rep
}

// reportEvidence prints the walk's outcome: every warning, every
// unreachable github.com citation, then the refusal or the summary line.
func reportEvidence(w io.Writer, rep evidenceReport, issues int) error {
	for _, u := range rep.Warnings {
		fmt.Fprintf(w, "warning: external link does not resolve — %s\n", u)
	}
	for _, u := range rep.Unreachable {
		fmt.Fprintf(w, "unreachable: %s\n", u)
	}
	if len(rep.Unreachable) > 0 {
		return refuse("EVIDENCE_UNREACHABLE", "%d of %d cited links do not resolve — commit or repair the evidence before dispatching QA", len(rep.Unreachable), rep.Total)
	}
	if n := len(rep.Warnings); n > 0 {
		fmt.Fprintf(w, "%d of %d cited %s resolve across %d issues — every github.com citation does; the %d external %s above do not, and QA weighs them — evidence is reachable\n",
			rep.Total-n, rep.Total, plural("link", rep.Total), issues, n, plural("link", n))
		return nil
	}
	fmt.Fprintf(w, "all %d cited %s resolve across %d issues — evidence is reachable\n", rep.Total, plural("link", rep.Total), issues)
	return nil
}

func plural(word string, n int) string {
	if n == 1 {
		return word
	}
	return word + "s"
}

// milestoneEvidence walks a milestone's record — the tracking issue and
// every sub-issue, bodies and comments — and verifies every cited link
// resolves. The M4-R4 lesson made deterministic: a record that exists only
// in a working tree is not evidence, and QA must not be dispatched against
// citations that 404 (findings 5 and 7 on #73 repeated it twice). What
// counts as a citation, and which failures refuse, is checkEvidence's.
func milestoneEvidence(w io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: gh codecrew milestone evidence <n>")
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

	if body, err := c.t.IssueBody(milestone.Ref); err == nil {
		fmt.Fprintln(w, requirementsNote(tracker.RequirementIDs(body)))
	}

	refs := []tracker.IssueRef{milestone.Ref}
	subs, err := c.t.SubIssues(milestone.Ref)
	if err != nil {
		return err
	}
	refs = append(refs, subs...)

	var records []evidenceRecord
	for _, ref := range refs {
		rec := evidenceRecord{Ref: ref}
		if body, err := c.t.IssueBody(ref); err == nil {
			rec.Texts = append(rec.Texts, body)
		}
		if comments, err := c.t.Comments(ref); err == nil {
			for _, cm := range comments {
				rec.Texts = append(rec.Texts, cm.Body)
			}
		}
		records = append(records, rec)
	}
	return reportEvidence(w, checkEvidence(records), len(refs))
}
