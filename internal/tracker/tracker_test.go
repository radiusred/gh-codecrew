package tracker

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestParseRef(t *testing.T) {
	cases := []struct {
		in   string
		want IssueRef
		ok   bool
	}{
		{"12", IssueRef{"o/def", 12}, true},
		{"#12", IssueRef{"o/def", 12}, true},
		{"radiusred/codecrew#7", IssueRef{"radiusred/codecrew", 7}, true},
		{"nonsense", IssueRef{}, false},
		{"owner/repo", IssueRef{}, false},
	}
	for _, c := range cases {
		got, err := ParseRef(c.in, "o/def")
		if c.ok && (err != nil || got != c.want) {
			t.Errorf("ParseRef(%q) = %v, %v; want %v", c.in, got, err, c.want)
		}
		if !c.ok && err == nil {
			t.Errorf("ParseRef(%q) should fail", c.in)
		}
	}
}

func TestNextMilestoneNumber(t *testing.T) {
	titles := []string{"M1: First", "M3: Skipped ahead", "not a milestone", "M2: Second"}
	if got := NextMilestoneNumber(titles); got != 4 {
		t.Errorf("NextMilestoneNumber = %d, want 4", got)
	}
	if got := NextMilestoneNumber(nil); got != 1 {
		t.Errorf("NextMilestoneNumber(empty) = %d, want 1", got)
	}
}

func TestPlanPresent(t *testing.T) {
	planless := "## Goal\nX\n\n## Plan\n" + PlanPlaceholder + "\n\n## Ask-the-human points\nNone."
	if PlanPresent(planless) {
		t.Error("placeholder plan should not count as present")
	}
	if PlanPresent("## Goal\nX\n\n## Ask-the-human points\nNone.") {
		t.Error("missing Plan section should not count as present")
	}
	planned := "## Goal\nX\n\n## Plan\n- change the thing\n\n## Ask-the-human points\nNone."
	if !PlanPresent(planned) {
		t.Error("real plan should count as present")
	}
}

func TestExtractRecords(t *testing.T) {
	src := IssueRef{"o/r", 4}
	comments := []Comment{
		{Author: "cody", Body: "**Decision:** use X\n**Trade-off:** Y\n**Rejected:** Z"},
		{Author: "human", Body: "just a chat comment mentioning **Deviation:** midway"},
		{Author: "cody", Body: "  **Deviation:** skipped W\n**Why:** unnecessary"},
		{Author: "human", Body: "**Gate resolved:** rename the repo\n**Trade-off:** redirects vs a second repo"},
	}
	records := ExtractRecords(src, comments)
	if len(records) != 3 {
		t.Fatalf("got %d records, want 3", len(records))
	}
	if records[0].Kind != "Decision" || records[1].Kind != "Deviation" {
		t.Errorf("kinds = %s, %s", records[0].Kind, records[1].Kind)
	}
	if records[2].Kind != "Decision" {
		t.Errorf("gate resolution kind = %s, want Decision", records[2].Kind)
	}
	if records[0].Source != "o/r#4" {
		t.Errorf("source = %q", records[0].Source)
	}
	if records[0].Label != "**Decision:**" || records[0].Body != "**Decision:** use X\n**Trade-off:** Y\n**Rejected:** Z" {
		t.Errorf("one-record comment altered: label %q body %q", records[0].Label, records[0].Body)
	}
}

// The M5 close missed three of #73's four Decisions because their labels
// carried a parenthetical qualifier, and PR #95's Deviation because it was
// the second paragraph of a review round-up (#113).
func TestExtractRecordsVariants(t *testing.T) {
	src := IssueRef{"o/r", 4}
	comments := []Comment{
		{Body: "**Decision (superseding the R1 counting claim; proven live):** write-access Apps count.\n**Trade-off:** privilege for a gated merge."},
		{Body: "**Gate resolved (operator, 2026-08-25):** guided creation."},
		{Body: "**Review findings addressed:** (1) fixed.\n\n**Deviation:** plan step 2 promised tests; none possible.\n**Why:** no fake yet.\n\nOn the tally: noted."},
		{Body: "**Decision:** first.\n\n**Finding 3:** unrelated.\n\n**Deviation:** second, in the same comment.\n\n**Why:** because."},
		{Body: "Prose paragraph.\n\nMore prose that mentions **Decision:** mid-line only."},
	}
	records := ExtractRecords(src, comments)
	want := []struct{ kind, label string }{
		{"Decision", "**Decision (superseding the R1 counting claim; proven live):**"},
		{"Decision", "**Gate resolved (operator, 2026-08-25):**"},
		{"Deviation", "**Deviation:**"},
		{"Decision", "**Decision:**"},
		{"Deviation", "**Deviation:**"},
	}
	if len(records) != len(want) {
		t.Fatalf("got %d records, want %d: %+v", len(records), len(want), records)
	}
	for i, w := range want {
		if records[i].Kind != w.kind || records[i].Label != w.label {
			t.Errorf("record %d = %s %q, want %s %q", i, records[i].Kind, records[i].Label, w.kind, w.label)
		}
	}
	// The mid-comment Deviation keeps its Why and the trailing prose, not the round-up before it.
	if b := records[2].Body; !strings.HasPrefix(b, "**Deviation:**") || !strings.Contains(b, "**Why:** no fake yet.") || !strings.Contains(b, "On the tally") || strings.Contains(b, "Review findings") {
		t.Errorf("mid-comment record body wrong: %q", b)
	}
	// A different bold label ends the record; the next record label starts another.
	if b := records[3].Body; strings.Contains(b, "Finding 3") || strings.Contains(b, "second") {
		t.Errorf("record not closed by the other label: %q", b)
	}
	if b := records[4].Body; !strings.Contains(b, "**Why:** because.") {
		t.Errorf("continuation not attached: %q", b)
	}
}

func TestUnresolvedGatesAcceptsQualifiedResolution(t *testing.T) {
	comments := []Comment{
		{Body: "**Gate raised:** which path?"},
		{Body: "**Gate resolved (operator, 2026-08-25):** the manifest flow."},
	}
	if got := UnresolvedGates(comments); len(got) != 0 {
		t.Errorf("qualified resolution not recognised: %d unresolved", len(got))
	}
	comments = append(comments[:1], Comment{Body: "**Deviation (small):** not a resolution."})
	if got := UnresolvedGates(comments); len(got) != 1 {
		t.Errorf("a Deviation must not resolve a gate: %d unresolved", len(got))
	}
}

// The real M5 comment corpus (every sub-issue of #67 and its closing PRs,
// fetched 2026-08-26): the gather must yield all seven records, including
// the four the M5 close missed.
func TestExtractRecordsM5Corpus(t *testing.T) {
	data, err := os.ReadFile("testdata/m5-corpus.json")
	if err != nil {
		t.Fatal(err)
	}
	var corpus struct {
		Sources []struct {
			Ref        string    `json:"ref"`
			Comments   []Comment `json:"comments"`
			ClosingPRs []struct {
				Ref      string    `json:"ref"`
				Comments []Comment `json:"comments"`
			} `json:"closing_prs"`
		} `json:"sources"`
	}
	if err := json.Unmarshal(data, &corpus); err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, s := range corpus.Sources {
		ref, _ := ParseRef(s.Ref, "radiusred/gh-codecrew")
		for _, r := range ExtractRecords(ref, s.Comments) {
			got = append(got, r.Source+" "+r.Kind+" "+r.Label)
		}
		for _, pr := range s.ClosingPRs {
			prRef, _ := ParseRef(pr.Ref, "radiusred/gh-codecrew")
			for _, r := range ExtractRecords(prRef, pr.Comments) {
				got = append(got, r.Source+" "+r.Kind+" "+r.Label)
			}
		}
	}
	want := []string{
		"radiusred/gh-codecrew#68 Decision **Gate resolved:**",
		"radiusred/gh-codecrew#73 Decision **Decision:**",
		"radiusred/gh-codecrew#73 Decision **Decision (M5-R1, the required-review question — answered empirically on the numberguess M3 proof run):**",
		"radiusred/gh-codecrew#73 Decision **Decision (superseding the R1 Decision's counting claim; operator hypothesis 2026-08-25, proven live):**",
		"radiusred/gh-codecrew#73 Decision **Decision (operator direction, 2026-08-25):**",
		"radiusred/gh-codecrew#76 Deviation **Deviation:**",
		"radiusred/gh-codecrew#95 Deviation **Deviation:**",
	}
	if len(got) != len(want) {
		t.Fatalf("gathered %d records from the M5 corpus, want %d:\n%s", len(got), len(want), strings.Join(got, "\n"))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("record %d:\n got %s\nwant %s", i, got[i], want[i])
		}
	}
}

func TestUnresolvedGates(t *testing.T) {
	gate := Comment{Author: "cody", Body: "**Gate raised:** rename or split?"}
	resolved := Comment{Author: "human", Body: "**Gate resolved:** rename"}
	decision := Comment{Author: "human", Body: "**Decision:** rename"}
	chat := Comment{Author: "human", Body: "thinking about it"}
	cases := []struct {
		name     string
		comments []Comment
		want     int
	}{
		{"no gates", []Comment{chat, decision}, 0},
		{"gate then resolution", []Comment{gate, resolved}, 0},
		{"gate then decision counts", []Comment{gate, decision}, 0},
		{"gate with only chat after", []Comment{gate, chat}, 1},
		{"resolution before gate does not count", []Comment{resolved, gate}, 1},
		{"second gate after resolution is unresolved", []Comment{gate, resolved, gate}, 1},
		{"one trailing resolution covers earlier gates", []Comment{gate, gate, resolved}, 0},
	}
	for _, c := range cases {
		if got := len(UnresolvedGates(c.comments)); got != c.want {
			t.Errorf("%s: got %d unresolved, want %d", c.name, got, c.want)
		}
	}
}

func TestRequirementIDs(t *testing.T) {
	body := "## Goal\nX\n\n## Requirements\n" +
		"- **M3-R1** — gate resolutions are gathered\n" +
		"- **M3-R2** — task start is role-aware\n" +
		"- **M3-R1** — duplicated ID is deduped\n\n" +
		"## Gates\n- **M3-R9** — outside the Requirements section, ignored\n"
	got := RequirementIDs(body)
	want := []string{"M3-R1", "M3-R2"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("RequirementIDs = %v, want %v", got, want)
	}
	if ids := RequirementIDs("## Goal\nno requirements section"); len(ids) != 0 {
		t.Errorf("missing section should yield none, got %v", ids)
	}
}

func TestParseVerdicts(t *testing.T) {
	comments := []Comment{
		// M2's actual shape: bullet list, em-dash inside the bold, trailing period.
		{Author: "testy", Body: "## QA verdict\n\n- **M2-R1 — satisfied.** built and ran it\n- **M2-R2 — not satisfied.** stale claim\n"},
		{Author: "testy", Body: "**M2-R2 — satisfied.** This verdict supersedes my earlier one."},
		{Author: "cody", Body: "**M2-R4 — satisfied.** (not QA — caller filters by author)"},
		{Author: "human", Body: "- **M2-R3** — a requirement definition line must not parse as a verdict"},
	}
	got := ParseVerdicts(comments)
	want := []Verdict{
		{ID: "M2-R1", State: "satisfied", Author: "testy"},
		{ID: "M2-R2", State: "not satisfied", Author: "testy"},
		{ID: "M2-R2", State: "satisfied", Author: "testy"},
		{ID: "M2-R4", State: "satisfied", Author: "cody"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d verdicts %v, want %d", len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("verdict[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestNoChecksReported(t *testing.T) {
	if !NoChecksReported(fmt.Errorf("gh pr: no checks reported on the 'probe' branch")) {
		t.Error("gh's checkless failure should classify as no-checks")
	}
	if NoChecksReported(fmt.Errorf("gh pr: HTTP 502")) {
		t.Error("unrelated errors must not classify as no-checks")
	}
	if NoChecksReported(nil) {
		t.Error("nil error must not classify as no-checks")
	}
}

func TestInferState(t *testing.T) {
	cases := []struct {
		name string
		task Task
		want State
	}{
		{"open unassigned is ready", Task{}, Ready},
		{"assigned is in progress", Task{Assignees: []string{"cody"}}, InProgress},
		{"open PR is in review", Task{Assignees: []string{"cody"}, OpenLinkedPR: true}, InReview},
		{"needs-decision gates even in review", Task{Labels: []string{"cc:task", "cc:needs-decision"}, OpenLinkedPR: true}, Gated},
		{"closed is done regardless", Task{Closed: true, Labels: []string{"cc:needs-decision"}, OpenLinkedPR: true}, Done},
	}
	for _, c := range cases {
		if got := InferState(c.task); got != c.want {
			t.Errorf("%s: InferState = %q, want %q", c.name, got, c.want)
		}
	}
}
