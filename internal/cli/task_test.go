package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// mergeGate implements the #73 Decisions: GitHub counts approvals only
// from write-access principals — a read-only App's review and an operator
// confirmation do not count — so a REVIEW_REQUIRED decision refuses
// unless the operator asked for the recorded bypass.
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
		{"CHANGES_REQUESTED", false, false, ""},                 // no approving review exists here in the common case, so the approval gate refuses earlier; with a separate non-author approval, GitHub itself refuses the merge
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

func TestHolderReviewed(t *testing.T) {
	holds := func(login string) bool { return login == "radiusred-reviewy" || login == "radiusred-reviewy[bot]" }
	if holderReviewed([]string{"davison", "someone"}, holds) {
		t.Error("non-holder approvals satisfied the holder gate")
	}
	if !holderReviewed([]string{"davison", "radiusred-reviewy"}, holds) {
		t.Error("holder approval not recognized among others")
	}
	if holderReviewed(nil, holds) {
		t.Error("no approvals satisfied the holder gate")
	}
}

// The ownership gate holds a task to the seat that started it — the same
// login, or the same routed seat (a team-held role is any member). An
// owner who has left resolves to no role and no longer matches; no owner
// recorded holds nobody (#165, operator's questions on #175).
func TestSameSeat(t *testing.T) {
	roleFor := func(login string) string {
		switch strings.ToLower(strings.TrimSuffix(login, "[bot]")) {
		case "radiusred-cody":
			return "implementer"
		case "radiusred-wordy":
			return "doc-synthesizer"
		case "alice", "bob": // members of the team the reviewer seat routes to
			return "reviewer"
		}
		return ""
	}
	for _, c := range []struct {
		owner, viewer string
		want          bool
	}{
		{"", "anyone", true}, // no start record
		{"radiusred-cody[bot]", "radiusred-cody[bot]", true}, // same login
		{"radiusred-cody[bot]", "Radiusred-Cody", true},      // suffix and case
		{"alice", "bob", true},                               // same team-held seat
		{"bob", "alice", true},
		{"alice", "radiusred-cody[bot]", false}, // different seats
		{"radiusred-wordy[bot]", "radiusred-cody[bot]", false},
		{"carol", "alice", false},                 // starter no longer holds any seat
		{"davison", "radiusred-cody[bot]", false}, // operator-held task, a crew finisher
		{"radiusred-cody[bot]", "davison", false}, // the operator is not exempt
		{"davison", "davison", true},
	} {
		if got := sameSeat(c.owner, c.viewer, roleFor); got != c.want {
			t.Errorf("sameSeat(%q, %q) = %v, want %v", c.owner, c.viewer, got, c.want)
		}
	}
}

// gateFake is the slice of the tracker checkpoint touches: it reads the
// target's labels and records the comment and label it writes.
type gateFake struct {
	tracker.Tracker
	labels   []string
	comments []string
	added    []string
}

func (f *gateFake) Task(ref tracker.IssueRef) (tracker.Task, error) {
	return tracker.Task{Ref: ref, Labels: f.labels}, nil
}
func (f *gateFake) Comment(_ tracker.IssueRef, body string) error {
	f.comments = append(f.comments, body)
	return nil
}
func (f *gateFake) AddLabel(_ tracker.IssueRef, label string) error {
	f.added = append(f.added, label)
	return nil
}

// checkpoint accepts a milestone ref — a question about a requirement has
// no task to carry it — and says so: nothing mechanical blocks on that
// label, status lists the gate instead. A task keeps the task finish
// wording (#200).
func TestRaiseGateWordingByTarget(t *testing.T) {
	for _, tc := range []struct {
		name, label, comment, receipt, absent string
	}{
		{"task", "cc:task", "`task finish` refuses while the label is present", "gate raised on o/r#6 — blocked until a human removes cc:needs-decision\n", "(milestone issue)"},
		{"milestone", tracker.LabelMilestone, "`status` lists this gate beside the tasks' gates", "gate raised on o/r#6 (milestone issue) — status lists it beside the tasks' gates until a human removes cc:needs-decision\n", "`task finish` refuses"},
	} {
		f := &gateFake{labels: []string{tc.label}}
		cfg := &config.Config{Codecrew: "1.0", Hub: "self"}
		c := &ctx{cfg: cfg, roles: cfg, current: "o/r", hub: "o/r", t: f}
		var out bytes.Buffer
		if err := raiseGate(&out, c, tracker.IssueRef{Repo: "o/r", Number: 6}, "which way?"); err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if out.String() != tc.receipt {
			t.Errorf("%s: receipt = %q, want %q", tc.name, out.String(), tc.receipt)
		}
		if len(f.comments) != 1 || !strings.HasPrefix(f.comments[0], "**Gate raised:** which way?\n\n") {
			t.Fatalf("%s: comments = %q", tc.name, f.comments)
		}
		if !strings.Contains(f.comments[0], tc.comment) || strings.Contains(f.comments[0], tc.absent) {
			t.Errorf("%s: comment must say %q and not %q:\n%s", tc.name, tc.comment, tc.absent, f.comments[0])
		}
		if len(f.added) != 1 || f.added[0] != tracker.LabelNeedsDecision {
			t.Errorf("%s: labels added = %q", tc.name, f.added)
		}
	}
}
