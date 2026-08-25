package cli

import (
	"errors"
	"strings"
	"testing"
)

// mergeGate implements the R1 Decision (radiusred/gh-codecrew#73): GitHub
// never counts App reviews or operator confirmations toward a
// required-review rule, so a REVIEW_REQUIRED decision refuses unless the
// operator asked for the recorded bypass.
func TestMergeGate(t *testing.T) {
	for _, tc := range []struct {
		decision string
		bypass   bool
		admin    bool
		refused  string
	}{
		{"", false, false, ""},                                  // no rule
		{"APPROVED", false, false, ""},                          // counted approval satisfies GitHub
		{"REVIEW_REQUIRED", false, false, "REVIEW_NOT_COUNTED"}, // the R1 configuration, no bypass asked
		{"REVIEW_REQUIRED", true, true, ""},                     // recorded bypass path
		{"CHANGES_REQUESTED", false, false, ""},                 // upstream review gates catch this first
	} {
		admin, err := mergeGate(tc.decision, tc.bypass)
		if tc.refused == "" {
			if err != nil {
				t.Errorf("mergeGate(%q,%v) unexpected error %v", tc.decision, tc.bypass, err)
			}
			if admin != tc.admin {
				t.Errorf("mergeGate(%q,%v) admin = %v, want %v", tc.decision, tc.bypass, admin, tc.admin)
			}
			continue
		}
		var r refusal
		if !errors.As(err, &r) || r.Code != tc.refused {
			t.Errorf("mergeGate(%q,%v) = %v, want refused[%s]", tc.decision, tc.bypass, err, tc.refused)
		}
		if !strings.Contains(err.Error(), "--bypass") {
			t.Errorf("refusal does not name the supported path: %v", err)
		}
	}
}
