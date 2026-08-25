package cli

import (
	"fmt"
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
