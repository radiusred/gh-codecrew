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
	body := fmt.Sprintf(milestoneTemplate, "A goal.", 2, 2)
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
	if ids := tracker.RequirementIDs(fmt.Sprintf(milestoneTemplate, "g", 1, 1)); len(ids) != 0 {
		t.Errorf("template must yield no IDs, got %v", ids)
	}
}
