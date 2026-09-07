package tracker

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// The fixtures below are written with LF and read three ways: as they
// stand, and with every line ending turned into CRLF — what GitHub's web
// editor saves when a human edits an issue, a comment or a PR body — each
// of those through the GitHub reader that fetches the text, and then the
// CRLF one again straight into the scanner with no reader in the path.
//
// The three readings pin the two layers separately, which is the whole
// reason there are two (NormalizeLineEndings):
//
//   - through the reader — the path a body actually takes into the
//     package. Delete the normalisation in github.go and these fail.
//   - straight into the scanner — the seam Tracker is: a string reaching
//     an exported scanner from another backend, a fake tracker or a
//     caller's own hand. Delete the normalisation at the scanner entries
//     in tracker.go and these fail, github.go untouched.
//
// Measured, one layer removed at a time:
//
//   - github.go's two calls removed, tracker.go untouched:
//     TestReadersNormaliseAtTheBoundary fails — the package is handed a
//     body carrying CR, and everything downstream of the readers, the
//     citation walk in internal/cli included, sees two shapes of body.
//   - tracker.go's seven calls removed, github.go untouched: the
//     AdoptedRefs, ExtractRecords and UnresolvedGates rows fail on their
//     direct reading.
//
// Three of the eight rows, because only three of the scans can be
// defeated by a CR at all. The `## Adopts` heading is matched by a
// `(?m)…$` line and Go's `$` matches only before `\n`, so a browser-edited
// task body yielded no adoptions and `task finish` would have shut none of
// the captures the task carried (#296); the paragraph split collapses such
// a comment into one paragraph, so a Decision swallows the gate raised
// after it. The other five scans — the two `## `-section cuts, the verdict
// line and the start record — read a CRLF body correctly already, by not
// being anchored to a line's end. Their rows are guards rather than pins:
// they assert the property the second layer exists to make unconditional,
// and they are what fails if one of those scans grows a line-end anchor.

func toCRLF(s string) string { return strings.ReplaceAll(s, "\n", "\r\n") }

const taskBodyFixture = `## Goal
Read a web-editor body as the tracker reads an LF one.

## Adopts
- #296 — the line-anchored regexps miss CRLF bodies
- o/other#7 — a capture in another repo

## Plan
Normalise where a body enters the package, and at every scanner's entry.

## Ask-the-human points
_None._
`

const plannedNothingFixture = `## Goal
A task nobody has planned yet.

## Plan
` + PlanPlaceholder + `

## Ask-the-human points
_None._
`

const milestoneBodyFixture = `## Goal
The milestone.

## Requirements
- **M15-R4** — normalise CRLF at the tracker boundary
- **M15-R5** — the stale-branch sweep
- **M1-R1** — a stray ID from another milestone

## Gates
_None._
`

// The comment fixtures, in the order they were posted: a start record, a
// gate raised in a comment's second paragraph, its resolution, a verdict
// beside one quoted in a code span, and a Decision with a **Trade-off:**
// continuation followed by a gate that nothing answers.
var commentFixtures = []struct{ author, body string }{
	{"radiusred-cody", "**Started by** @radiusred-cody."},
	{"darrendavison", "Round-up of the review.\n\nWhy the second layer?"},
	{"darrendavison", "**Gate raised:** does the normalisation belong at the boundary alone?"},
	{"darrendavison", "**Gate resolved (operator):** at the boundary and at each scanner's entry."},
	{"radiusred-testy", "**M15-R4 — satisfied.** The table test drives every scanner twice.\n\nA verdict quoted in code — `**M15-R5 — satisfied.**` — is content, not a verdict."},
	{"radiusred-cody", "**Decision:** normalise at the boundary and at each exported scanner's entry.\n**Trade-off:** two passes over a body instead of one.\n\n**Gate raised:** does SPEC §4 need a sentence about line endings?"},
}

// rawComments builds the fixtures into Comment values directly, with
// ending applied to each body — a scanner reached without a reader.
func rawComments(ending func(string) string) []Comment {
	out := make([]Comment, len(commentFixtures))
	for i, c := range commentFixtures {
		out[i] = Comment{
			Author: c.author,
			Body:   ending(c.body),
			URL:    "https://github.com/o/r/issues/1#issuecomment-" + string(rune('a'+i)),
		}
	}
	return out
}

// issueBodyVia serves body as the issue JSON and reads it back through
// GitHub.IssueBody.
func issueBodyVia(t *testing.T, body string) string {
	t.Helper()
	payload, err := json.Marshal(struct {
		Body string `json:"body"`
	}{body})
	if err != nil {
		t.Fatal(err)
	}
	fakeGH(t, "", "GH_HELPER_ISSUE_JSON="+string(payload))
	got, err := GitHub{}.IssueBody(IssueRef{Repo: "o/r", Number: 1})
	if err != nil {
		t.Fatalf("IssueBody: %v", err)
	}
	return got
}

// commentsVia serves the fixtures as the comment listing, with lineEnding
// applied to each body, and reads them back through GitHub.Comments.
func commentsVia(t *testing.T, ending func(string) string) []Comment {
	t.Helper()
	type raw struct {
		Body string `json:"body"`
		URL  string `json:"html_url"`
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	}
	items := make([]raw, len(commentFixtures))
	for i, c := range commentFixtures {
		items[i].Body = ending(c.body)
		items[i].URL = "https://github.com/o/r/issues/1#issuecomment-" + string(rune('a'+i))
		items[i].User.Login = c.author
	}
	payload, err := json.Marshal(items)
	if err != nil {
		t.Fatal(err)
	}
	fakeGH(t, "", "GH_HELPER_COMMENTS_JSON="+string(payload))
	got, err := GitHub{}.Comments(IssueRef{Repo: "o/r", Number: 1})
	if err != nil {
		t.Fatalf("Comments: %v", err)
	}
	return got
}

// Every scanner that reads an issue body: the CRLF reading must equal the
// LF one, and the LF one must be right.
func TestBodyScannersReadCRLFAsLF(t *testing.T) {
	cases := []struct {
		name string
		body string
		scan func(string) any
		want any
	}{
		{"AdoptedRefs", taskBodyFixture,
			func(b string) any { return AdoptedRefs(b, "o/r") },
			[]IssueRef{{Repo: "o/r", Number: 296}, {Repo: "o/other", Number: 7}}},
		{"PlanPresent", taskBodyFixture,
			func(b string) any { return PlanPresent(b) }, true},
		{"PlanPresent on the placeholder", plannedNothingFixture,
			func(b string) any { return PlanPresent(b) }, false},
		{"RequirementIDs", milestoneBodyFixture,
			func(b string) any { return RequirementIDs(b) },
			[]string{"M15-R4", "M15-R5", "M1-R1"}},
		{"MismatchedRequirementIDs", milestoneBodyFixture,
			func(b string) any { return MismatchedRequirementIDs(b, 15) },
			[]string{"M1-R1"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lf := tc.scan(issueBodyVia(t, tc.body))
			if !reflect.DeepEqual(lf, tc.want) {
				t.Fatalf("LF body: got %#v, want %#v", lf, tc.want)
			}
			crlf := tc.scan(issueBodyVia(t, toCRLF(tc.body)))
			if !reflect.DeepEqual(crlf, lf) {
				t.Errorf("CRLF body through the reader: got %#v, want the LF reading %#v", crlf, lf)
			}
			direct := tc.scan(toCRLF(tc.body))
			if !reflect.DeepEqual(direct, lf) {
				t.Errorf("CRLF body straight into the scanner: got %#v, want the LF reading %#v", direct, lf)
			}
		})
	}
}

// Every scanner that reads comments — which is where the paragraph split
// lives, so ExtractRecords and UnresolvedGates cover it. Each scanner is
// projected onto strings so a failure names what was read rather than a
// struct dump.
func TestCommentScannersReadCRLFAsLF(t *testing.T) {
	records := func(cs []Comment) any {
		var out []string
		for _, r := range ExtractRecords(IssueRef{Repo: "o/r", Number: 1}, cs) {
			out = append(out, r.Kind+"|"+r.Author+"|"+r.Body)
		}
		return out
	}
	// UnresolvedGates hands back the Comment values it was given rather
	// than normalised copies — scanning is its job, rewriting its input is
	// not — so the row projects what task finish reads off the result, the
	// author and the URL identifying which comment carries the open gate.
	// That the gate was found in that comment's *second* paragraph is the
	// assertion: it is the paragraph split being read right.
	gates := func(cs []Comment) any {
		var out []string
		for _, c := range UnresolvedGates(cs) {
			out = append(out, c.Author+"|"+c.URL)
		}
		return out
	}
	verdicts := func(cs []Comment) any {
		var out []string
		for _, v := range ParseVerdicts(cs) {
			out = append(out, v.ID+"|"+v.State+"|"+v.Author)
		}
		return out
	}
	cases := []struct {
		name string
		scan func([]Comment) any
		want any
	}{
		{"ExtractRecords", records, []string{
			"Decision|darrendavison|**Gate resolved (operator):** at the boundary and at each scanner's entry.",
			"Decision|radiusred-cody|**Decision:** normalise at the boundary and at each exported scanner's entry.\n**Trade-off:** two passes over a body instead of one.",
		}},
		{"UnresolvedGates", gates, []string{
			"radiusred-cody|https://github.com/o/r/issues/1#issuecomment-f",
		}},
		{"ParseVerdicts", verdicts, []string{"M15-R4|satisfied|radiusred-testy"}},
		{"StartedBy", func(cs []Comment) any { return StartedBy(cs) }, "radiusred-cody"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lf := tc.scan(commentsVia(t, func(s string) string { return s }))
			if !reflect.DeepEqual(lf, tc.want) {
				t.Fatalf("LF comments: got %#v, want %#v", lf, tc.want)
			}
			crlf := tc.scan(commentsVia(t, toCRLF))
			if !reflect.DeepEqual(crlf, lf) {
				t.Errorf("CRLF comments through the reader: got %#v, want the LF reading %#v", crlf, lf)
			}
			direct := tc.scan(rawComments(toCRLF))
			if !reflect.DeepEqual(direct, lf) {
				t.Errorf("CRLF comments straight into the scanner: got %#v, want the LF reading %#v", direct, lf)
			}
		})
	}
}

// The readers hand the package LF whatever GitHub stored, so everything
// downstream of them — the citation walk in internal/cli included — reads
// one shape of body.
func TestReadersNormaliseAtTheBoundary(t *testing.T) {
	if got := issueBodyVia(t, toCRLF(taskBodyFixture)); strings.Contains(got, "\r") {
		t.Errorf("IssueBody returned a body carrying CR: %q", got)
	}
	for _, c := range commentsVia(t, toCRLF) {
		if strings.Contains(c.Body, "\r") {
			t.Errorf("Comments returned a body carrying CR: %q", c.Body)
		}
	}
}
