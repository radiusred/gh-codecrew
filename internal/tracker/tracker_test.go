package tracker

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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

func TestMilestoneNumberHolder(t *testing.T) {
	self := IssueRef{Repo: "o/hub", Number: 7}
	milestones := []TitledIssue{{Ref: IssueRef{Repo: "o/hub", Number: 1}, Title: "M1: One"}}
	recent := []TitledIssue{
		{Ref: IssueRef{Repo: "o/hub", Number: 9}, Title: "M3 rules the roost"}, // no colon: not a milestone title
		{Ref: self, Title: "M3: Ours"},
		{Ref: IssueRef{Repo: "o/hub", Number: 6}, Title: "M3: Theirs"},
	}
	if got := MilestoneNumberHolder(3, self, milestones, recent); got == nil || got.Ref.Number != 6 {
		t.Errorf("holder of M3 = %v, want o/hub#6 (self and the colonless title skipped)", got)
	}
	if got := MilestoneNumberHolder(2, self, milestones, recent); got != nil {
		t.Errorf("holder of M2 = %v, want none", got)
	}
	if got := Titles(milestones, recent); len(got) != 4 || got[0] != "M1: One" || got[3] != "M3: Theirs" {
		t.Errorf("Titles = %v", got)
	}
	if got := NextMilestoneNumber(Titles(milestones, recent)); got != 4 {
		t.Errorf("next over both listings = %d, want 4", got)
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

// A comment typed in GitHub's web UI arrives with CRLF line endings; the
// paragraph split must still find a record after other text.
func TestExtractRecordsCRLF(t *testing.T) {
	records := ExtractRecords(IssueRef{"o/r", 1}, []Comment{{Body: "**Review findings addressed:** done.\r\n\r\n**Deviation:** skipped W.\r\n**Why:** unnecessary.\r\n"}})
	if len(records) != 1 || records[0].Kind != "Deviation" || !strings.Contains(records[0].Body, "**Why:** unnecessary.") {
		t.Errorf("CRLF comment: got %+v", records)
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

// The protocol-2.0 gate grammar (M13-R6): **Gate raised:** is read per
// paragraph like every other label, only **Gate resolved:** resolves, and
// it resolves the gates open before it — never the ones after.
// The tightened rules must not reclassify the hub's own history. The M5
// corpus carries the project's one hand-raised gate (#68) and its
// **Gate resolved:** answer: it was resolved before protocol 2.0 and stays
// resolved after it.
func TestUnresolvedGatesM5Corpus(t *testing.T) {
	data, err := os.ReadFile("testdata/m5-corpus.json")
	if err != nil {
		t.Fatal(err)
	}
	var corpus struct {
		Sources []struct {
			Ref      string    `json:"ref"`
			Comments []Comment `json:"comments"`
		} `json:"sources"`
	}
	if err := json.Unmarshal(data, &corpus); err != nil {
		t.Fatal(err)
	}
	raised := 0
	for _, src := range corpus.Sources {
		for _, c := range src.Comments {
			for _, para := range paragraphs(c.Body) {
				if gateRaisedLabel.MatchString(para) {
					raised++
				}
			}
		}
		if got := UnresolvedGates(src.Comments); len(got) != 0 {
			t.Errorf("%s: %d gate(s) newly unresolved", src.Ref, len(got))
		}
	}
	if raised != 1 {
		t.Errorf("the corpus should carry the one recorded gate, found %d", raised)
	}
}

func TestUnresolvedGates(t *testing.T) {
	gate := Comment{Author: "cody", Body: "**Gate raised:** rename or split?"}
	gateLate := Comment{Author: "cody", Body: "I have pushed the rebase and the tests are green.\n\n**Gate raised:** rename or split?"}
	gateTwice := Comment{Author: "cody", Body: "**Gate raised:** rename or split?\n\n**Gate raised:** and which name?"}
	resolved := Comment{Author: "human", Body: "**Gate resolved:** rename"}
	resolvedLate := Comment{Author: "human", Body: "Read both branches.\n\n**Gate resolved:** rename"}
	decision := Comment{Author: "human", Body: "**Decision:** rename"}
	deviation := Comment{Author: "cody", Body: "**Deviation:** skipped W.\n**Why:** unnecessary."}
	quoted := Comment{Author: "human", Body: "The convention is a paragraph opening **Gate raised:** like this one mid-line."}
	chat := Comment{Author: "human", Body: "thinking about it"}
	cases := []struct {
		name     string
		comments []Comment
		want     int
	}{
		{"no gates", []Comment{chat, decision}, 0},
		{"gate then resolution", []Comment{gate, resolved}, 0},
		{"a bare Decision no longer resolves", []Comment{gate, decision}, 1},
		{"a Deviation never resolved", []Comment{gate, deviation}, 1},
		{"gate with only chat after", []Comment{gate, chat}, 1},
		{"resolution before gate does not count", []Comment{resolved, gate}, 1},
		{"second gate after resolution is unresolved", []Comment{gate, resolved, gate}, 1},
		{"one resolution covers every gate open before it", []Comment{gate, gate, resolved}, 0},
		{"a gate raised in a later paragraph counts", []Comment{gateLate}, 1},
		{"a gate raised in a later paragraph is resolvable", []Comment{gateLate, resolved}, 0},
		{"a resolution in a later paragraph resolves", []Comment{gate, resolvedLate}, 0},
		{"two gates in one comment are two gates", []Comment{gateTwice}, 2},
		{"two gates in one comment, one resolution", []Comment{gateTwice, resolved}, 0},
		{"a label mid-line is not a gate", []Comment{quoted}, 0},
		{"a resolution then a later gate leaves the later one", []Comment{gate, resolved, gateLate}, 1},
	}
	for _, c := range cases {
		if got := len(UnresolvedGates(c.comments)); got != c.want {
			t.Errorf("%s: got %d unresolved, want %d", c.name, got, c.want)
		}
	}
	// The comment carrying the gate is what task finish links to.
	if got := UnresolvedGates([]Comment{gateLate}); len(got) != 1 || got[0].Author != "cody" {
		t.Errorf("the raising comment is returned: %+v", got)
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

// radiusred/numberguess#3, as the orchestrator run wrote it (#119 finding
// 28, #144): six bold IDs under "## Goal", the scaffold's placeholder left
// under "## Requirements". The parser is section-scoped on purpose (the
// placeholder guard above), so this body yields nothing — and a close over
// it must refuse rather than verify zero requirements.
func TestRequirementIDsIgnoresIDsOutsideTheSection(t *testing.T) {
	body, err := os.ReadFile("testdata/numberguess-m1-body.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "**M1-R6**") {
		t.Fatal("fixture lost its bold IDs")
	}
	if ids := RequirementIDs(string(body)); len(ids) != 0 {
		t.Errorf("IDs outside the Requirements section counted: %v", ids)
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

// Verdict supersession, protocol 2.0 (M13-R6): per comment, not per
// match. Within a comment the first verdict for an ID counts, so a QA
// comment that restates an earlier verdict below its own no longer
// supersedes itself; across comments the latest comment carrying a verdict
// for the ID wins, which is what milestone close tallies.
func TestParseVerdictsPerComment(t *testing.T) {
	cases := []struct {
		name     string
		comments []Comment
		want     []Verdict
	}{
		{
			"the first match for an ID in a comment counts",
			[]Comment{{Author: "testy", Body: "**M2-R1 — not satisfied.** the probe fails\n\nMy earlier **M2-R1 — satisfied.** was wrong."}},
			[]Verdict{{ID: "M2-R1", State: "not satisfied", Author: "testy"}},
		},
		{
			"a later comment supersedes",
			[]Comment{
				{Author: "testy", Body: "**M2-R1 — not satisfied.** the probe fails"},
				{Author: "testy", Body: "**M2-R1 — satisfied.** fixed on the follow-up PR"},
			},
			[]Verdict{
				{ID: "M2-R1", State: "not satisfied", Author: "testy"},
				{ID: "M2-R1", State: "satisfied", Author: "testy"},
			},
		},
		{
			"a verdict quoted in a fenced block is content",
			[]Comment{{Author: "testy", Body: "The form is:\n\n```markdown\n**M2-R1 — satisfied.** <evidence>\n```\n\nMine follows.\n\n**M2-R1 — not satisfied.** the probe fails"}},
			[]Verdict{{ID: "M2-R1", State: "not satisfied", Author: "testy"}},
		},
		{
			"a verdict quoted in a code span is content",
			[]Comment{{Author: "testy", Body: "I withdraw `**M2-R1 — satisfied.**` and record:\n\n**M2-R1 — untestable.** no environment"}},
			[]Verdict{{ID: "M2-R1", State: "untestable", Author: "testy"}},
		},
		{
			"a fence quoting one ID does not hide another",
			[]Comment{{Author: "testy", Body: "```\n**M2-R1 — satisfied.**\n```\n\n**M2-R2 — satisfied.** it holds"}},
			[]Verdict{{ID: "M2-R2", State: "satisfied", Author: "testy"}},
		},
		{
			// testy's probe on #260, verbatim: the quoted line is sample
			// output, and only the verdict in prose below it counts.
			"a verdict quoted in an indented block is content",
			[]Comment{{Author: "testy", Body: "Quoted example:\n\n    **M2-R1 — satisfied.** this is sample output, not a verdict\n\nActual prose.\n\n**M2-R1 — not satisfied.** the probe fails"}},
			[]Verdict{{ID: "M2-R1", State: "not satisfied", Author: "testy"}},
		},
		{
			// testy's re-verdict probe on #254, verbatim (#288): an
			// indented block needs no blank line after a heading.
			"a verdict quoted in an indented block after a heading is content",
			[]Comment{{Author: "testy", Body: "### Example\n    **M2-R1 — satisfied.** sample output, not a verdict\n\n**M2-R1 — not satisfied.** actual verdict"}},
			[]Verdict{{ID: "M2-R1", State: "not satisfied", Author: "testy"}},
		},
		{
			"a verdict quoted in an indented block after a thematic break is content",
			[]Comment{{Author: "testy", Body: "---\n    **M2-R1 — satisfied.** sample output\n\n**M2-R1 — untestable.** no environment"}},
			[]Verdict{{ID: "M2-R1", State: "untestable", Author: "testy"}},
		},
		{
			"a verdict on an indented paragraph continuation is a verdict",
			[]Comment{{Author: "testy", Body: "I reran the suite and\n    **M2-R1 — satisfied.** is the record"}},
			[]Verdict{{ID: "M2-R1", State: "satisfied", Author: "testy"}},
		},
		{
			"a verdict on a list item's indented continuation is a verdict",
			[]Comment{{Author: "testy", Body: "- reran on merged main:\n    **M2-R2 — satisfied.** green, and the probe held"}},
			[]Verdict{{ID: "M2-R2", State: "satisfied", Author: "testy"}},
		},
		{
			"different IDs in one comment all count",
			[]Comment{{Author: "testy", Body: "- **M2-R1 — satisfied.** a\n- **M2-R2 — not satisfied.** b\n"}},
			[]Verdict{
				{ID: "M2-R1", State: "satisfied", Author: "testy"},
				{ID: "M2-R2", State: "not satisfied", Author: "testy"},
			},
		},
	}
	for _, c := range cases {
		got := ParseVerdicts(c.comments)
		if len(got) != len(c.want) {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
			continue
		}
		for i := range c.want {
			if got[i] != c.want[i] {
				t.Errorf("%s: verdict[%d] = %v, want %v", c.name, i, got[i], c.want[i])
			}
		}
	}
}

// The tally milestone close performs over that stream: last write per ID
// wins, which under the per-comment rule means the latest comment's first
// match for the ID.
func TestVerdictTallyIsPerComment(t *testing.T) {
	comments := []Comment{
		{Author: "testy", Body: "- **M2-R1 — satisfied.** a\n- **M2-R2 — satisfied.** b\n"},
		{Author: "testy", Body: "**M2-R2 — not satisfied.** the hairbrush case\n\nEarlier I wrote **M2-R2 — satisfied.**, which this supersedes."},
	}
	latest := map[string]string{}
	for _, v := range ParseVerdicts(comments) {
		latest[v.ID] = v.State
	}
	if latest["M2-R1"] != "satisfied" || latest["M2-R2"] != "not satisfied" {
		t.Errorf("tally = %v, want M2-R1 satisfied and M2-R2 not satisfied", latest)
	}
}

// SPEC §4: a requirement ID is M<milestone>-R<k>. IDs outside the
// Requirements section are not requirements and are not checked.
func TestMismatchedRequirementIDs(t *testing.T) {
	cases := []struct {
		name string
		body string
		n    int
		want []string
	}{
		{"a matching set", "## Requirements\n- **M13-R1** — a\n- **M13-R2** — b\n", 13, nil},
		{"a foreign ID", "## Requirements\n- **M13-R1** — a\n- **M12-R3** — b\n", 13, []string{"M12-R3"}},
		{"every ID foreign", "## Requirements\n- **M1-R1** — a\n- **M2-R1** — b\n", 13, []string{"M1-R1", "M2-R1"}},
		{"a prefix is not a match", "## Requirements\n- **M1-R1** — a\n", 13, []string{"M1-R1"}},
		{"M1 does not swallow M13", "## Requirements\n- **M13-R1** — a\n", 1, []string{"M13-R1"}},
		{"IDs outside the section are not requirements", "## Requirements\n- **M13-R1** — a\n\n## Gates\n- **M9-R9** — a gate\n", 13, nil},
		{"no section", "## Goal\nnothing here", 13, nil},
	}
	for _, c := range cases {
		got := MismatchedRequirementIDs(c.body, c.n)
		if fmt.Sprint(got) != fmt.Sprint(c.want) {
			t.Errorf("%s: MismatchedRequirementIDs = %v, want %v", c.name, got, c.want)
		}
	}
}

// One code rule for the whole binary: the citation walk and the verdict
// scan read a comment through StripCode, so what is content for one is
// content for the other (M13-R6). All three Markdown code forms are
// stripped — spans, fences and indented blocks (#285) — and an indented
// block opens wherever it would not interrupt a paragraph, which is
// CommonMark's only restriction on the form (#288).
func TestStripCode(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{"a span is blanked", "keep `drop` keep", "keep   keep"},
		{"a fence is blanked", "a\n```\ndrop\n```\nb\n", "a\nb\n"},
		{"a tilde fence is blanked", "a\n~~~\ndrop\n~~~\nb\n", "a\nb\n"},
		{"an unclosed fence runs to the end", "a\n```\ndrop\ndrop\n", "a\n"},
		{"an unclosed backtick run is literal", "a ` b\n", "a ` b\n"},
		// The finding's probe (#260, comment 5560080971): four spaces after
		// a blank line is a code block, so what it holds is content.
		{"an indented block is blanked", "Quoted example:\n\n    **M13-R1 — satisfied.** sample output\n\nActual prose.\n", "Quoted example:\n\nActual prose.\n"},
		{"a tab-indented block is blanked", "Quoted example:\n\n\t**M13-R1 — satisfied.** sample output\n\nActual prose.\n", "Quoted example:\n\nActual prose.\n"},
		{"a blank line inside an indented block belongs to it", "a\n\n    one\n\n    two\n\nb\n", "a\n\nb\n"},
		{"an indented block opens at the start of the text", "    drop\nb\n", "b\n"},
		{"an indented block ends at the first line under four columns", "a\n\n    drop\n  b\n", "a\n\n  b\n"},
		{"an indented paragraph continuation is not a block", "The report says\n    **M13-R1 — satisfied.** and means it\n", "The report says\n    **M13-R1 — satisfied.** and means it\n"},
		{"a list item's indented continuation is not a block", "- the probe ran and\n    **M13-R1 — satisfied.** is the record\n", "- the probe ran and\n    **M13-R1 — satisfied.** is the record\n"},
		{"a fence indented four columns is code either way", "a\n\n    ```\n    drop\n    ```\nb\n", "a\n\nb\n"},
		// The openers CommonMark allows with no blank line after them: a
		// heading, a thematic break, a fence's close (#288). The heading
		// shape is testy's probe on #285, comment 5560244463.
		{"an indented block opens after an ATX heading", "### Example\n    **M13-R6 — satisfied.** sample output\n\nprose\n", "### Example\nprose\n"},
		{"an indented block opens after a closed ATX heading", "# Title #\n    drop\nprose\n", "# Title #\nprose\n"},
		{"an indented block opens after every heading level", "###### h6\n    drop\nprose\n", "###### h6\nprose\n"},
		{"seven hashes are not a heading", "####### not a heading\n    kept\n", "####### not a heading\n    kept\n"},
		{"a hash run with no space is not a heading", "#tag\n    kept\n", "#tag\n    kept\n"},
		{"an indented block opens after a thematic break", "---\n    drop\nprose\n", "---\nprose\n"},
		{"a thematic break of asterisks or underscores opens one too", "* * *\n    drop\n___\n    drop\nprose\n", "* * *\n___\nprose\n"},
		{"a bullet is not a thematic break", "- item\n    kept\n", "- item\n    kept\n"},
		{"two dashes are not a thematic break", "--\n    kept\n", "--\n    kept\n"},
		{"mixed break characters are not a thematic break", "*-*\n    kept\n", "*-*\n    kept\n"},
		// The "nothing else on the line" clause: without it, a heading-like
		// aside would open a block and swallow the line under it — an
		// over-strip, which costs a verdict rather than a strip (#289).
		{"three dashes with trailing text are not a thematic break", "--- x\n    kept\n", "--- x\n    kept\n"},
		{"a break's characters around other text are not a thematic break", "*** NOTE ***\n    kept\n", "*** NOTE ***\n    kept\n"},
		{"an indented block opens after a fence closes", "```\ndrop\n```\n    drop\nprose\n", "prose\n"},
		{"a heading indented four columns is code, not an opener", "a\n\n    ### in code\n    drop\nprose\n", "a\n\nprose\n"},
		{"a heading lazily continuing a paragraph is not an opener", "prose\n    ### not a heading\n    kept\n", "prose\n    ### not a heading\n    kept\n"},
		{"a prose line still prevents opening", "The report says\n    **M13-R6 — satisfied.** and means it\n", "The report says\n    **M13-R6 — satisfied.** and means it\n"},
	}
	for _, c := range cases {
		if got := StripCode(c.in); got != c.want {
			t.Errorf("%s: StripCode(%q) = %q, want %q", c.name, c.in, got, c.want)
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

// task start posts a **Started by** comment on every start (App
// identities are not assignable, so the comment is the record, not the
// assignment); StartedBy reads it back, latest comment first, so task
// finish can hold the seat that started a task to finishing it (#165).
func TestStartedBy(t *testing.T) {
	started := func(login string) Comment { return Comment{Author: login, Body: StartRecord(login)} }
	cases := []struct {
		name string
		cs   []Comment
		want string
	}{
		{"nothing recorded", nil, ""},
		{"comment", []Comment{{Body: "plan…"}, started("radiusred-wordy[bot]")}, "radiusred-wordy[bot]"},
		{"latest comment wins", []Comment{started("radiusred-cody[bot]"), started("radiusred-wordy[bot]")}, "radiusred-wordy[bot]"},
		{"prose mentioning the phrase is not a record", []Comment{{Author: "x", Body: "The **Started by** @x comment is how we know."}}, ""},
		{"a record with trailing prose is not a record", []Comment{started("radiusred-cody[bot]"), {Author: "radiusred-checky[bot]", Body: "**Started by** @radiusred-checky[bot]. This is an example."}}, "radiusred-cody[bot]"},
		{"a record naming someone its author is not", []Comment{started("radiusred-cody[bot]"), {Author: "radiusred-checky[bot]", Body: StartRecord("radiusred-cody[bot]")}}, "radiusred-cody[bot]"},
		{"author matches with the suffix ignored", []Comment{{Author: "radiusred-cody", Body: StartRecord("radiusred-cody[bot]")}}, "radiusred-cody[bot]"},
		{"human restart: latest record wins", []Comment{started("alice"), started("bob")}, "bob"},
	}
	for _, c := range cases {
		if got := StartedBy(c.cs); got != c.want {
			t.Errorf("%s: StartedBy = %q, want %q", c.name, got, c.want)
		}
	}
}

// The 1.0 shim: an assigned task with no start record had the first
// assignee for an owner, "for tasks started before the record existed".
// Deleted in protocol 2.0 (M13-R7) — assignment is not a start, an
// assignee never chose to hold the seat, and the record task start posts
// is the only thing that says a task was started.
func TestStartedByIgnoresAssignees(t *testing.T) {
	if got := StartedBy(nil); got != "" {
		t.Errorf("an unstarted task has owner %q, want none", got)
	}
	// A task whose only record is prose is unstarted too: nothing about
	// who is assigned to it can make it started.
	if got := StartedBy([]Comment{{Author: "alice", Body: "Assigned to me, picking it up tomorrow."}}); got != "" {
		t.Errorf("prose gave the task owner %q, want none", got)
	}
}

func TestSameLogin(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want bool
	}{
		{"radiusred-cody[bot]", "radiusred-cody[bot]", true},
		{"radiusred-cody[bot]", "radiusred-cody", true},
		{"@radiusred-cody[bot]", "Radiusred-Cody", true},
		{"davison", "Davison", true},
		{"radiusred-cody[bot]", "radiusred-wordy[bot]", false},
		{"davison", "radiusred-cody[bot]", false},
	} {
		if got := SameLogin(c.a, c.b); got != c.want {
			t.Errorf("SameLogin(%q, %q) = %v", c.a, c.b, got)
		}
	}
}

// The two shapes #198 recorded, in the order the field met them: the
// rollup itself needs checks: read; once granted, the workflow run under
// each suite needs actions: read. The same message on any other path, and
// every other failure, is not this — it must surface raw.
func TestMissingChecksPermission(t *testing.T) {
	cases := map[string]string{
		"gh pr: GraphQL: Resource not accessible by integration (node.statusCheckRollup.nodes.0.commit.statusCheckRollup)":                                                                                                                                           "checks: read",
		"gh pr: GraphQL: Resource not accessible by integration (node.statusCheckRollup.nodes.0.commit.statusCheckRollup.contexts.nodes.0.checkSuite.workflowRun)":                                                                                                   "actions: read",
		"gh pr: GraphQL: Resource not accessible by integration (repository.collaborators)":                                                                                                                                                                          "",
		"gh pr: GraphQL: Resource not accessible by integration (node.statusCheckRollup.nodes.0.commit.statusCheckRollup.contexts.nodes.0.checkSuite.workflowRun), Resource not accessible by integration (node.statusCheckRollup.nodes.0.commit.statusCheckRollup)": "actions: read",
		"gh pr: GraphQL: Something went wrong while executing your query (node.statusCheckRollup)":                                                                                                                                                                   "",
		"gh pr: HTTP 502": "",
		"gh pr: no checks reported on the 'probe' branch": "",
	}
	for msg, want := range cases {
		if got := MissingChecksPermission(fmt.Errorf("%s", msg)); got != want {
			t.Errorf("%q: got %q, want %q", msg, got, want)
		}
	}
	if got := MissingChecksPermission(nil); got != "" {
		t.Errorf("nil: %q", got)
	}
}

// The ## Adopts section is read one list line at a time, and only the ref
// at the head of a line counts: the prose after it is the capture's title,
// so a `#42` inside a title is not an adoption, and neither is a ref
// written anywhere else in the body. Bare numbers resolve against the
// task's own repo, full refs against themselves, and a repeat collapses.
func TestAdoptedRefs(t *testing.T) {
	body := `## Goal
Delivers #900, which is not an adoption.

## Requirements
M14-R1

## Adopts
- #193 — task new/close: carry an adopted backlog issue
* o/other#7 — the tidy verb, fixing #42 on the way
- #193 — the same capture again
- not a ref at all
- #8

## Plan
Adopts #999 in prose, which is not the section.
`
	want := []IssueRef{
		{Repo: "o/hub", Number: 193},
		{Repo: "o/other", Number: 7},
		{Repo: "o/hub", Number: 8},
	}
	if got := AdoptedRefs(body, "o/hub"); !reflect.DeepEqual(got, want) {
		t.Errorf("AdoptedRefs = %v, want %v", got, want)
	}
	if got := AdoptedRefs("## Goal\nNo section here.\n", "o/hub"); got != nil {
		t.Errorf("a body with no ## Adopts section adopted %v", got)
	}
}

// AdoptsBlock writes what AdoptedRefs reads, and writes nothing at all for
// a task that adopts nothing — the body task new has always written.
func TestAdoptsBlockRoundTrips(t *testing.T) {
	if got := AdoptsBlock("o/hub", nil); got != "" {
		t.Errorf("no adoptions rendered %q", got)
	}
	adopted := []Adoption{
		{Ref: IssueRef{Repo: "o/hub", Number: 193}, Title: "carry an adopted backlog issue"},
		{Ref: IssueRef{Repo: "o/other", Number: 7}, Title: ""},
	}
	block := AdoptsBlock("o/hub", adopted)
	if want := "\n## Adopts\n- #193 — carry an adopted backlog issue\n- o/other#7\n"; block != want {
		t.Fatalf("AdoptsBlock = %q, want %q", block, want)
	}
	want := []IssueRef{adopted[0].Ref, adopted[1].Ref}
	if got := AdoptedRefs(block+"\n## Plan\n", "o/hub"); !reflect.DeepEqual(got, want) {
		t.Errorf("the block reads back as %v, want %v", got, want)
	}
}
