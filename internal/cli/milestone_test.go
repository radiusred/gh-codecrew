package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// The regression guard renders the real template: if the Requirements
// placeholder ever matches the requirement parser again (finding 12 on
// #73 — a bold example became a phantom requirement at a live milestone
// close), this fails no matter how the template is reworded.
func TestMilestoneTemplatePlaceholderYieldsNoRequirements(t *testing.T) {
	body, err := milestoneBody("A goal.", 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	if ids := tracker.RequirementIDs(body); len(ids) != 0 {
		t.Errorf("fresh milestone-new body yields requirements: %v", ids)
	}
}

// The orchestrator run closed numberguess M1 over a "not satisfied"
// verdict because its requirements sat under Goal and the gate iterated
// zero IDs (#119 finding 28, #144). Zero requirements now refuses.
func TestCloseRefusesWhenRequirementsSectionIsEmpty(t *testing.T) {
	body, err := os.ReadFile("../tracker/testdata/numberguess-m1-body.md")
	if err != nil {
		t.Fatal(err)
	}
	ref := tracker.IssueRef{Repo: "radiusred/numberguess", Number: 3}
	err = requireRequirements(ref, tracker.RequirementIDs(string(body)))
	if err == nil {
		t.Fatal("expected refusal, got nil")
	}
	if !strings.Contains(err.Error(), "refused[NO_REQUIREMENTS]") || !strings.Contains(err.Error(), "radiusred/numberguess#3") {
		t.Errorf("unexpected refusal: %v", err)
	}
	if err := requireRequirements(ref, []string{"M1-R1"}); err != nil {
		t.Errorf("populated section must pass, got %v", err)
	}
}

func TestRequirementsNote(t *testing.T) {
	if n := requirementsNote(nil); !strings.Contains(n, "NO_REQUIREMENTS") || !strings.Contains(n, "## Requirements") {
		t.Errorf("empty note must name the section and the refusal: %q", n)
	}
	if n := requirementsNote([]string{"M1-R1", "M1-R2"}); n != "requirements counted: M1-R1, M1-R2 (2)" {
		t.Errorf("counted note = %q", n)
	}
	// The scaffolded body always starts empty, so `milestone new` always notes it.
	body, _ := milestoneBody("g", 1, nil)
	if ids := tracker.RequirementIDs(body); len(ids) != 0 {
		t.Errorf("template must yield no IDs, got %v", ids)
	}
}

// The close's DOC_MISSING detail once said "dispatch the doc-synthesizer,
// merge its PR, rerun" — and an orchestrator did exactly that, planning a
// by-hand merge with an identity that could not. The detail names the task
// path (#119 finding 27).
func TestDocMissingNamesTheTaskPath(t *testing.T) {
	err := docMissing(6, "radiusred/gh-codecrew")
	if err == nil {
		t.Fatal("expected refusal, got nil")
	}
	msg := err.Error()
	for _, want := range []string{"refused[DOC_MISSING]", "docs/milestones/6-*.md", "radiusred/gh-codecrew", "task start", "Closes #", "task finish"} {
		if !strings.Contains(msg, want) {
			t.Errorf("DOC_MISSING detail missing %q: %s", want, msg)
		}
	}
	if strings.Contains(msg, "merge its PR") {
		t.Errorf("DOC_MISSING detail still sends the coordinator to merge by hand: %s", msg)
	}
}

// Both milestones the orchestrator opened carried their requirement IDs
// under Goal, because --goal was the only text input the verb offered
// (#147; #119 findings 19a, 28, 32). --requirement writes them where the
// parser reads them, numbered by the CLI.
func TestMilestoneBodyWritesRequirementsWhereTheParserReads(t *testing.T) {
	body, err := milestoneBody("Play from the terminal.", 2, []string{
		"npm start launches an interactive game",
		"invalid input is rejected with a hint",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := tracker.RequirementIDs(body), []string{"M2-R1", "M2-R2"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("RequirementIDs = %v, want %v\n%s", got, want, body)
	}
	if !strings.Contains(body, "- **M2-R2** — invalid input is rejected with a hint") {
		t.Errorf("requirement line not rendered:\n%s", body)
	}
	if strings.Contains(body, "this placeholder") {
		t.Errorf("placeholder survives alongside real requirements:\n%s", body)
	}
	if !strings.Contains(body, "## Goal\nPlay from the terminal.\n\n## Requirements\n") || !strings.Contains(body, "\n\n## Gates\n") {
		t.Errorf("section layout changed:\n%s", body)
	}
}

func TestMilestoneBodyRefusesTextThatCarriesAnID(t *testing.T) {
	for _, r := range []string{"M2-R1 — already numbered", "**M1-R3** — bold and numbered", "  M9-R1: spaced"} {
		if _, err := milestoneBody("g", 2, []string{"fine", r}); err == nil {
			t.Errorf("%q accepted; the CLI numbers requirements", r)
		} else if !strings.Contains(err.Error(), "carries an ID") {
			t.Errorf("%q: unexpected error %v", r, err)
		}
	}
	if _, err := milestoneBody("g", 2, []string{"fine", "  "}); err == nil {
		t.Error("empty requirement accepted")
	}
}

// numberguess #11 was titled "M2 — …" while the CLI counted 3 and created
// "M3: M2 — …", closed as a duplicate. A prefix that disagrees is refused;
// one that agrees is stripped, never doubled.
func TestMilestoneTitleAgreesWithTheNumber(t *testing.T) {
	cases := []struct {
		title string
		n     int
		want  string
		err   bool
	}{
		{"Interactive CLI", 3, "M3: Interactive CLI", false},
		{"M3: Interactive CLI", 3, "M3: Interactive CLI", false},
		{"M3 — Interactive CLI", 3, "M3: Interactive CLI", false},
		{"M3 - Interactive CLI", 3, "M3: Interactive CLI", false},
		{"M2 — Interactive CLI", 3, "", true},
		{"M2: Interactive CLI", 3, "", true},
		{"M3:", 3, "", true},
		{"M3 rules the roost", 3, "M3: M3 rules the roost", false}, // no separator: not a prefix
	}
	for _, c := range cases {
		got, err := milestoneTitle(c.title, c.n)
		if c.err {
			if err == nil {
				t.Errorf("%q n=%d: expected an error, got %q", c.title, c.n, got)
			} else if !strings.Contains(err.Error(), fmt.Sprintf("next milestone number is %d", c.n)) && !strings.Contains(err.Error(), "only a milestone number") {
				t.Errorf("%q: unexpected error %v", c.title, err)
			}
			continue
		}
		if err != nil || got != c.want {
			t.Errorf("%q n=%d = %q, %v; want %q", c.title, c.n, got, err, c.want)
		}
	}
}

// newFake is the slice of the tracker milestone new reads, plus a record
// of every issue it creates and edits. The embedded nil Tracker makes any
// other call a panic: the verb is an API call and nothing else.
//
// Both listings are eventually consistent in the field (#195), so each is
// a sequence of answers: call k gets the k-th, and the last one repeats.
// A fake whose listings never learn of the issue it created is exactly the
// lag that produced three M2s.
type newFake struct {
	tracker.Tracker
	milestones [][]tracker.TitledIssue // successive MilestoneIssues answers
	recent     [][]tracker.TitledIssue // successive RecentIssues answers
	created    []string
	edits      []string // "<title>\n<body>" per EditIssue
	editErr    error
	mCalls     int
	rCalls     int
}

func nth(seq [][]tracker.TitledIssue, k int) []tracker.TitledIssue {
	if len(seq) == 0 {
		return nil
	}
	if k >= len(seq) {
		k = len(seq) - 1
	}
	return seq[k]
}

func (f *newFake) MilestoneIssues(string) ([]tracker.TitledIssue, error) {
	f.mCalls++
	return nth(f.milestones, f.mCalls-1), nil
}
func (f *newFake) RecentIssues(string) ([]tracker.TitledIssue, error) {
	f.rCalls++
	return nth(f.recent, f.rCalls-1), nil
}
func (f *newFake) CreateIssue(repo, title, _ string, _ []string) (tracker.IssueRef, error) {
	f.created = append(f.created, title)
	return tracker.IssueRef{Repo: repo, Number: 207}, nil
}
func (f *newFake) EditIssue(_ tracker.IssueRef, title, body string) error {
	if f.editErr != nil {
		return f.editErr
	}
	f.edits = append(f.edits, title+"\n"+body)
	return nil
}

// issue is a listing entry in the hub.
func issue(n int, title string) tracker.TitledIssue {
	return tracker.TitledIssue{Ref: tracker.IssueRef{Repo: "o/hub", Number: n}, Title: title}
}

// listing is a listing that never changes: the same answer to every call.
func listing(issues ...tracker.TitledIssue) [][]tracker.TitledIssue {
	return [][]tracker.TitledIssue{issues}
}

// milestoneNewCtx is a hub whose pointer directory holds a ROADMAP.md, so
// a verb that still wrote the file would be caught.
func milestoneNewCtx(t *testing.T, f *newFake) (*ctx, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "ROADMAP.md")
	if err := os.WriteFile(path, []byte(roadmapScaffold), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{Codecrew: "1.0", Hub: "self", Dir: dir}
	return &ctx{cfg: cfg, roles: cfg, current: "o/hub", hub: "o/hub", t: f}, path
}

// The row has no PR to ride in when a milestone's tasks all live in spokes
// (#197, hit three times in the field), so milestone new no longer edits
// ROADMAP.md or tells any seat to commit a row: the doc-synthesizer adds it
// as Done in the record PR (M10-R1).
func TestMilestoneNewLeavesTheRoadmapAlone(t *testing.T) {
	for _, dryRun := range []bool{false, true} {
		f := &newFake{milestones: listing(issue(1, "M1: One"), issue(2, "M2: Two"))}
		c, path := milestoneNewCtx(t, f)
		var out bytes.Buffer
		if err := runMilestoneNew(c, &out, "Bookkeeping", "Close the gaps.", []string{"the row belongs to the record PR"}, dryRun); err != nil {
			t.Fatalf("dryRun=%v: %v", dryRun, err)
		}
		got := out.String()
		if strings.Contains(strings.ToLower(got), "roadmap") {
			t.Errorf("dryRun=%v: milestone new still talks about the roadmap:\n%s", dryRun, got)
		}
		if !strings.Contains(got, "M3: Bookkeeping") || !strings.Contains(got, "requirements counted: M3-R1 (1)") {
			t.Errorf("dryRun=%v: number, title or requirement IDs missing:\n%s", dryRun, got)
		}
		after, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(after) != roadmapScaffold {
			t.Errorf("dryRun=%v: ROADMAP.md was edited:\n%s", dryRun, after)
		}
		if dryRun {
			if len(f.created) != 0 || !strings.Contains(got, "dry run: nothing written") {
				t.Errorf("dry run created %v:\n%s", f.created, got)
			}
			continue
		}
		if fmt.Sprint(f.created) != "[M3: Bookkeeping]" || !strings.Contains(got, "created milestone M3: Bookkeeping (o/hub#207)") {
			t.Errorf("live run created %v:\n%s", f.created, got)
		}
	}
}

// The floor is the max over both listings: the unfiltered recent listing
// knows the milestone created seconds ago that the label-filtered one has
// not indexed yet, so the number is right before anything is created and
// the dry run says the same number the live run would use (M10-R2, #195).
func TestMilestoneNewFloorReadsBothListings(t *testing.T) {
	for _, dryRun := range []bool{true, false} {
		f := &newFake{
			milestones: listing(issue(1, "M1: One"), issue(2, "M2: Two")),
			recent:     listing(issue(9, "a plain issue"), issue(3, "M3: Created seconds ago"), issue(2, "M2: Two")),
		}
		c, _ := milestoneNewCtx(t, f)
		var out bytes.Buffer
		if err := runMilestoneNew(c, &out, "Bookkeeping", "g", []string{"r"}, dryRun); err != nil {
			t.Fatalf("dryRun=%v: %v", dryRun, err)
		}
		got := out.String()
		if !strings.Contains(got, "M4: Bookkeeping") || !strings.Contains(got, "requirements counted: M4-R1 (1)") || strings.Contains(got, "M3: Bookkeeping") {
			t.Errorf("dryRun=%v: the recent listing's M3 did not raise the floor:\n%s", dryRun, got)
		}
		if dryRun {
			if f.mCalls != 1 || f.rCalls != 1 || len(f.created) != 0 {
				t.Errorf("dry run read %d+%d listings and created %v", f.mCalls, f.rCalls, f.created)
			}
			continue
		}
		if len(f.edits) != 0 || !strings.Contains(got, "number check: M4 is carried by o/hub#207 alone") {
			t.Errorf("a clean creation edited %v or skipped the check:\n%s", f.edits, got)
		}
		if f.mCalls != 2 || f.rCalls != 2 {
			t.Errorf("live run read the listings %d+%d times, want 2+2 (floor, then the post-create check)", f.mCalls, f.rCalls)
		}
	}
}

// Both listings lag at derivation time, so the verb creates a second M3 —
// the #195 failure — and the post-create check catches it: the new issue
// is renumbered, title and every requirement ID, and the output says so.
// The renumbered title survives a caller who wrote the (then correct)
// prefix themselves.
func TestMilestoneNewRepairsATakenNumber(t *testing.T) {
	for _, title := range []string{"Bookkeeping", "M3: Bookkeeping", "M3 — Bookkeeping"} {
		stale := listing(issue(1, "M1: One"), issue(2, "M2: Two"))
		f := &newFake{
			milestones: stale,
			recent: [][]tracker.TitledIssue{
				nil, // the floor: nothing indexed yet
				{issue(3, "M3: Rival"), issue(207, "M3: Bookkeeping")}, // the check: the rival and ourselves
			},
		}
		c, _ := milestoneNewCtx(t, f)
		var out bytes.Buffer
		if err := runMilestoneNew(c, &out, title, "g", []string{"first", "second"}, false); err != nil {
			t.Fatalf("%q: %v", title, err)
		}
		got := out.String()
		if fmt.Sprint(f.created) != "[M3: Bookkeeping]" {
			t.Errorf("%q: created %v", title, f.created)
		}
		if len(f.edits) != 1 || !strings.HasPrefix(f.edits[0], "M4: Bookkeeping\n") {
			t.Fatalf("%q: edits %q", title, f.edits)
		}
		body := f.edits[0][strings.Index(f.edits[0], "\n")+1:]
		if ids := tracker.RequirementIDs(body); fmt.Sprint(ids) != "[M4-R1 M4-R2]" {
			t.Errorf("%q: the repaired body's IDs are %v:\n%s", title, ids, body)
		}
		if strings.Contains(body, "M3-R") {
			t.Errorf("%q: the repaired body still carries M3 IDs:\n%s", title, body)
		}
		for _, line := range []string{
			"created milestone M3: Bookkeeping (o/hub#207)",
			"renumbered: M3 was already o/hub#3 (M3: Rival), so o/hub#207 is now M4: Bookkeeping",
			"number check: M4 is carried by o/hub#207 alone",
			"requirements counted: M4-R1, M4-R2 (2)",
		} {
			if !strings.Contains(got, line) {
				t.Errorf("%q: output missing %q:\n%s", title, line, got)
			}
		}
		if strings.Contains(got, "requirements counted: M3-R1") {
			t.Errorf("%q: the IDs reported are the ones the repair replaced:\n%s", title, got)
		}
	}
}

// A collision the verb cannot repair is refused with a code that names both
// issues and the hand fix — never left as two milestones sharing a prefix.
func TestMilestoneNewRefusesWhenTheRepairFails(t *testing.T) {
	f := &newFake{
		milestones: listing(issue(1, "M1: One"), issue(2, "M2: Two")),
		recent:     [][]tracker.TitledIssue{nil, {issue(3, "M3: Rival")}},
		editErr:    errors.New("HTTP 403"),
	}
	c, _ := milestoneNewCtx(t, f)
	var out bytes.Buffer
	err := runMilestoneNew(c, &out, "Bookkeeping", "g", []string{"r"}, false)
	var r refusal
	if !errors.As(err, &r) || r.Code != "MILESTONE_NUMBER_TAKEN" {
		t.Fatalf("err = %v, want refused[MILESTONE_NUMBER_TAKEN]", err)
	}
	for _, s := range []string{"o/hub#3", "M3: Rival", "o/hub#207", "renumbering to M4 failed", "HTTP 403", "gh issue edit --title", "M3-R<k>"} {
		if !strings.Contains(r.Detail, s) {
			t.Errorf("detail missing %q: %s", s, r.Detail)
		}
	}
	if strings.Contains(out.String(), "requirements counted") {
		t.Errorf("a refused creation still reported IDs:\n%s", out.String())
	}
}

// The repair is bounded: a listing that keeps producing a holder for every
// number the verb moves to is refused after milestoneNumberRepairs rounds,
// not chased forever, and every renumbering it did made it to GitHub.
func TestMilestoneNewBoundsTheRepair(t *testing.T) {
	var recent [][]tracker.TitledIssue
	recent = append(recent, nil) // the floor
	var rivals []tracker.TitledIssue
	for n := 3; n <= 3+milestoneNumberRepairs; n++ {
		rivals = append(rivals, issue(n, fmt.Sprintf("M%d: Rival", n)))
		recent = append(recent, append([]tracker.TitledIssue(nil), rivals...))
	}
	f := &newFake{milestones: listing(issue(2, "M2: Two")), recent: recent}
	c, _ := milestoneNewCtx(t, f)
	var out bytes.Buffer
	err := runMilestoneNew(c, &out, "Bookkeeping", "g", nil, false)
	var r refusal
	if !errors.As(err, &r) || r.Code != "MILESTONE_NUMBER_TAKEN" || !strings.Contains(r.Detail, fmt.Sprintf("still taken after %d renumberings", milestoneNumberRepairs)) {
		t.Fatalf("err = %v", err)
	}
	if len(f.edits) != milestoneNumberRepairs {
		t.Errorf("%d edits, want %d:\n%s", len(f.edits), milestoneNumberRepairs, out.String())
	}
	if want := fmt.Sprintf("M%d: Bookkeeping\n", 3+milestoneNumberRepairs); !strings.HasPrefix(f.edits[len(f.edits)-1], want) {
		t.Errorf("last edit %q, want prefix %q", f.edits[len(f.edits)-1], want)
	}
}

// The top-level help is a surface the runMilestoneNew test never reaches:
// it kept advertising "the number and row it would get" after the verb
// stopped printing a row (checky on PR #216). The dry-run line under
// `milestone new` must say what the verb prints and never name a row.
func TestUsageDescribesMilestoneNewDryRunWithoutARow(t *testing.T) {
	var block []string
	in := false
	for _, line := range strings.Split(usage, "\n") {
		switch {
		case strings.HasPrefix(line, "  milestone new "):
			in = true
		case in && !strings.HasPrefix(line, "           "):
			in = false
		}
		if in {
			block = append(block, line)
		}
	}
	if len(block) == 0 {
		t.Fatal("usage has no `milestone new` block")
	}
	text := strings.Join(block, "\n")
	var dry string
	for _, line := range block {
		if strings.Contains(line, "--dry-run") {
			dry = line
		}
	}
	if dry == "" {
		t.Fatalf("no --dry-run line under milestone new:\n%s", text)
	}
	for _, want := range []string{"number", "title", "requirement IDs", "create nothing"} {
		if !strings.Contains(dry, want) {
			t.Errorf("--dry-run line lacks %q: %s", want, dry)
		}
	}
	if lower := strings.ToLower(text); strings.Contains(lower, "row") || strings.Contains(lower, "roadmap") {
		t.Errorf("milestone new help still speaks of a row:\n%s", text)
	}
}
