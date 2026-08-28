package cli

import (
	"fmt"
	"os"
	"strings"
	"testing"

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
