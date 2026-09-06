package cli

import (
	"bytes"
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// labelFake is the slice of the tracker the label step touches: what the
// repository already defines, and what it is asked to define.
type labelFake struct {
	tracker.Tracker
	defined  []tracker.Label
	readErr  error
	failOn   string // the one label name a write refuses
	created  []tracker.Label
	restyled []tracker.Label
	reads    int
}

func (f *labelFake) Labels(string) ([]tracker.Label, error) {
	f.reads++
	return f.defined, f.readErr
}

func (f *labelFake) CreateLabel(_ string, l tracker.Label) error {
	if l.Name == f.failOn {
		return errors.New("403")
	}
	f.created = append(f.created, l)
	return nil
}

func (f *labelFake) UpdateLabel(_ string, l tracker.Label) error {
	if l.Name == f.failOn {
		return errors.New("403")
	}
	f.restyled = append(f.restyled, l)
	return nil
}

// defaults is the protocol's own styling for the named labels — a
// repository already indistinguishable from a fresh scaffold.
func defaults(names ...string) []tracker.Label {
	var out []tracker.Label
	for _, n := range names {
		l, _ := tracker.ProtocolLabel(n)
		out = append(out, l)
	}
	return out
}

// implicit is a label as GitHub creates one that nobody defined: a
// generated colour and no description at all.
func implicit(name string) tracker.Label {
	return tracker.Label{Name: name, Color: "ededed"}
}

func names(labels []tracker.Label) []string {
	var out []string
	for _, l := range labels {
		out = append(out, l.Name)
	}
	return out
}

// The whole of ensureLabels: what it creates, what it leaves, and that
// nothing it meets is a failure.
func TestEnsureLabels(t *testing.T) {
	all := []string{tracker.LabelMilestone, tracker.LabelTask, tracker.LabelNeedsDecision}
	for _, tc := range []struct {
		name    string
		fake    *labelFake
		created []string
		lines   []string
		absent  []string
	}{
		{
			name:    "a repository with none of them",
			fake:    &labelFake{},
			created: all,
			lines: []string{
				"created label cc:milestone (#01d4ff) — CodeCrew: the milestone tracking issue in the hub (SPEC §4)\n",
				"created label cc:task (#92edff) — ",
				"created label cc:needs-decision (#f0aeff) — ",
			},
		},
		{
			// Only the missing one is created, and the two that exist are
			// left as they are — whatever colour the project gave them,
			// which this fake cannot even report, because nothing reads it.
			name:    "a repository that already has two",
			fake:    &labelFake{defined: []tracker.Label{implicit("bug"), implicit(tracker.LabelMilestone), implicit(tracker.LabelTask)}},
			created: []string{tracker.LabelNeedsDecision},
			lines:   []string{"created label cc:needs-decision (#f0aeff)"},
			absent:  []string{"created label cc:milestone", "created label cc:task"},
		},
		{
			// GitHub compares label names case-insensitively, and so does
			// the check for what already exists — and a label the project
			// restyled keeps its colour, which is init's whole rule.
			name:    "a repository that restyled and recased them",
			fake:    &labelFake{defined: []tracker.Label{implicit("CC:Milestone"), implicit("cc:TASK"), implicit("CC:NEEDS-DECISION")}},
			created: nil,
			absent:  []string{"created label", "restyled label"},
		},
		{
			name:    "the label listing cannot be read",
			fake:    &labelFake{readErr: errors.New("403")},
			created: nil,
			lines:   []string{"note: could not read o/r's labels (403) — the protocol's labels are created on first use instead"},
			absent:  []string{"created label"},
		},
		{
			// One refusal does not end the run: the labels after it are
			// still attempted, and only the one that failed is a note.
			name:    "one creation is refused",
			fake:    &labelFake{failOn: tracker.LabelTask},
			created: []string{tracker.LabelMilestone, tracker.LabelNeedsDecision},
			lines: []string{
				"note: could not create the cc:task label in o/r (403)",
				"created label cc:milestone",
				"created label cc:needs-decision",
			},
			absent: []string{"created label cc:task"},
		},
	} {
		var out bytes.Buffer
		ensureLabels(&out, tc.fake, "o/r", tracker.ProtocolLabels)
		if got := names(tc.fake.created); !slices.Equal(got, tc.created) {
			t.Errorf("%s: created %v, want %v", tc.name, got, tc.created)
		}
		// ensureLabels never restyles, whatever it finds.
		if len(tc.fake.restyled) != 0 {
			t.Errorf("%s: ensureLabels restyled %v", tc.name, names(tc.fake.restyled))
		}
		for _, want := range tc.lines {
			if !strings.Contains(out.String(), want) {
				t.Errorf("%s: output missing %q:\n%s", tc.name, want, out.String())
			}
		}
		for _, never := range tc.absent {
			if strings.Contains(out.String(), never) {
				t.Errorf("%s: output must not contain %q:\n%s", tc.name, never, out.String())
			}
		}
		// One listing read, however many labels are wanted.
		if tc.fake.reads != 1 {
			t.Errorf("%s: read the label listing %d times, want 1", tc.name, tc.fake.reads)
		}
	}
}

// stubLabelTarget stands in for gh's answer to "which repository is this,
// and what writes to it".
func stubLabelTarget(t *testing.T, f tracker.Tracker, repo string, err error) {
	t.Helper()
	prev := labelTarget
	labelTarget = func() (tracker.Tracker, string, error) { return f, repo, err }
	t.Cleanup(func() { labelTarget = prev })
}

// init writes the scaffold to disk and the labels to GitHub, and the two
// must not be able to hurt each other: the labels are created after the
// commit, and every way the GitHub side can fail is a note the scaffold
// survives.
func TestInitCreatesTheLabelsAndTheScaffoldSurvivesGitHub(t *testing.T) {
	for _, tc := range []struct {
		name      string
		fake      *labelFake
		targetErr error
		created   []string
		lines     []string
		absent    []string
	}{
		{
			name:    "a repository with none of them",
			fake:    &labelFake{},
			created: []string{tracker.LabelMilestone, tracker.LabelTask, tracker.LabelNeedsDecision},
			lines:   []string{"created label cc:milestone (#01d4ff)", "created label cc:task (#92edff)", "created label cc:needs-decision (#f0aeff)"},
		},
		{
			name:    "GitHub refuses the listing",
			fake:    &labelFake{readErr: errors.New("Resource not accessible by integration")},
			created: nil,
			lines:   []string{"note: could not read o/r's labels ("},
			absent:  []string{"created label"},
		},
		{
			name:      "gh cannot say which repository this is",
			fake:      &labelFake{},
			targetErr: errors.New("no git remotes found"),
			created:   nil,
			lines:     []string{"note: could not ask GitHub which repository this is (no git remotes found) — the cc: labels are created on first use instead"},
			absent:    []string{"created label"},
		},
	} {
		stubProtection(t, false, true)
		stubLabelTarget(t, tc.fake, "o/r", tc.targetErr)
		dir := gitRepo(t)
		wd, _ := os.Getwd()
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		var out bytes.Buffer
		err := initCmd(&out, nil)
		os.Chdir(wd)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if got := names(tc.fake.created); !slices.Equal(got, tc.created) {
			t.Errorf("%s: created %v, want %v", tc.name, got, tc.created)
		}
		for _, want := range tc.lines {
			if !strings.Contains(out.String(), want) {
				t.Errorf("%s: output missing %q:\n%s", tc.name, want, out.String())
			}
		}
		for _, never := range tc.absent {
			if strings.Contains(out.String(), never) {
				t.Errorf("%s: output must not contain %q:\n%s", tc.name, never, out.String())
			}
		}
		// Whatever GitHub said, the scaffold is on disk and committed.
		if !strings.Contains(out.String(), "committed ") {
			t.Errorf("%s: the scaffold was not committed:\n%s", tc.name, out.String())
		}
		if show := committedFiles(t, dir, "HEAD"); !strings.Contains(show, ".codecrew/config.yml") {
			t.Errorf("%s: the scaffold commit is missing the pointer:\n%s", tc.name, show)
		}
		// The labels come after the commit: nothing they do can reach it.
		if i, j := strings.Index(out.String(), "committed "), strings.Index(out.String(), tc.lines[0]); i > j {
			t.Errorf("%s: the label step ran before the commit:\n%s", tc.name, out.String())
		}
	}
}

// A spoke gets all three, cc:milestone included: R5 says so, and a spoke
// promoted to a hub later already carries it.
func TestInitCreatesTheLabelsInASpokeToo(t *testing.T) {
	stubProtection(t, false, true)
	f := &labelFake{}
	stubLabelTarget(t, f, "o/spoke", nil)
	dir := gitRepo(t)
	wd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	err := initCmd(&out, []string{"--hub", "o/hub"})
	os.Chdir(wd)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{tracker.LabelMilestone, tracker.LabelTask, tracker.LabelNeedsDecision}
	if got := names(f.created); !slices.Equal(got, want) {
		t.Errorf("spoke created %v, want %v", got, want)
	}
}

// Outside a git repository there is no repository on GitHub to hold a
// label, so the step is skipped with a note and nothing is asked of gh.
func TestInitSkipsTheLabelsOutsideAGitRepository(t *testing.T) {
	f := &labelFake{}
	stubLabelTarget(t, f, "o/r", nil)
	dir := t.TempDir()
	wd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	err := initCmd(&out, nil)
	os.Chdir(wd)
	if err != nil {
		t.Fatal(err)
	}
	if want := "note: the cc: labels are created once this is a repository on GitHub — rerun gh codecrew init there"; !strings.Contains(out.String(), want) {
		t.Errorf("output missing %q:\n%s", want, out.String())
	}
	if f.reads != 0 || len(f.created) != 0 {
		t.Errorf("gh was asked about a repository that does not exist: %d reads, created %v", f.reads, f.created)
	}
}

// migrate restyles where init leaves alone: its whole promise is a
// repository indistinguishable from a fresh 2.0 scaffold, and a 1.x
// repository's cc: labels were all created implicitly, with a generated
// colour and no description (the operator's Decision on #283). --dry-run
// lists exactly the same steps and writes nothing.
func TestApplyLabelsRestyles(t *testing.T) {
	for _, tc := range []struct {
		name     string
		defined  []tracker.Label
		created  []string
		restyled []string
		lines    []string
		absent   []string
	}{
		{
			name:    "a 1.x repository with none of them",
			created: []string{tracker.LabelMilestone, tracker.LabelTask, tracker.LabelNeedsDecision},
			lines:   []string{"created label cc:milestone (#01d4ff)"},
			absent:  []string{"restyle"},
		},
		{
			// The implicit creations: GitHub's grey and no description.
			name:     "a 1.x repository whose labels came from first use",
			defined:  []tracker.Label{implicit(tracker.LabelMilestone), implicit(tracker.LabelTask), implicit(tracker.LabelNeedsDecision)},
			restyled: []string{tracker.LabelMilestone, tracker.LabelTask, tracker.LabelNeedsDecision},
			lines: []string{
				"restyled label cc:milestone (#ededed -> #01d4ff) — CodeCrew: the milestone tracking issue in the hub (SPEC §4)",
				"restyled label cc:task (#ededed -> #92edff)",
				"restyled label cc:needs-decision (#ededed -> #f0aeff)",
			},
			absent: []string{"created label"},
		},
		{
			// Already the protocol's: nothing to do, so a rerun is silent.
			name:    "a repository already wearing the defaults",
			defined: defaults(tracker.LabelMilestone, tracker.LabelTask, tracker.LabelNeedsDecision),
			absent:  []string{"label cc:"},
		},
		{
			// The colour matches, case aside, and the description does
			// not — a label styled by hand before the defaults existed.
			// The restyle is addressed by the name the repository spells
			// the label with — restyling is not renaming — while the
			// receipt names the protocol's own spelling.
			name:     "the colour matches and the description does not",
			defined:  []tracker.Label{{Name: "CC:Task", Color: "92EDFF", Description: "tasks"}},
			created:  []string{tracker.LabelMilestone, tracker.LabelNeedsDecision},
			restyled: []string{"CC:Task"},
			lines:    []string{"restyled label cc:task (#92EDFF -> #92edff)"},
		},
		{
			name:     "one implicit label and the two that are missing",
			defined:  []tracker.Label{implicit(tracker.LabelTask)},
			created:  []string{tracker.LabelMilestone, tracker.LabelNeedsDecision},
			restyled: []string{tracker.LabelTask},
		},
	} {
		// The real run writes, and reports each step.
		f := &labelFake{defined: slices.Clone(tc.defined)}
		var out bytes.Buffer
		applyLabels(&out, f, "o/r", tracker.ProtocolLabels, true, false)
		if got := names(f.created); !slices.Equal(got, tc.created) {
			t.Errorf("%s: created %v, want %v", tc.name, got, tc.created)
		}
		if got := names(f.restyled); !slices.Equal(got, tc.restyled) {
			t.Errorf("%s: restyled %v, want %v", tc.name, got, tc.restyled)
		}
		for _, want := range tc.lines {
			if !strings.Contains(out.String(), want) {
				t.Errorf("%s: output missing %q:\n%s", tc.name, want, out.String())
			}
		}
		for _, never := range tc.absent {
			if strings.Contains(out.String(), never) {
				t.Errorf("%s: output must not contain %q:\n%s", tc.name, never, out.String())
			}
		}
		// A restyle carries the protocol's colour and description, and is
		// addressed by the name the repository spells the label with — the
		// one the listing returned, never the protocol's own spelling.
		for _, l := range f.restyled {
			want, ok := tracker.ProtocolLabel(l.Name)
			if !ok {
				t.Errorf("%s: restyled %q, which is not a protocol label", tc.name, l.Name)
				continue
			}
			if l.Color != want.Color || l.Description != want.Description {
				t.Errorf("%s: restyled to %+v, want the defaults %+v", tc.name, l, want)
			}
			if had, _ := tracker.FindLabel(tc.defined, l.Name); had.Name != l.Name {
				t.Errorf("%s: restyle addressed %q, want the repository's own %q", tc.name, l.Name, had.Name)
			}
		}

		// The dry run lists the same steps, in the same order, and writes
		// nothing.
		d := &labelFake{defined: slices.Clone(tc.defined)}
		var dry bytes.Buffer
		applyLabels(&dry, d, "o/r", tracker.ProtocolLabels, true, true)
		if len(d.created) != 0 || len(d.restyled) != 0 {
			t.Errorf("%s: the dry run wrote: created %v restyled %v", tc.name, names(d.created), names(d.restyled))
		}
		want := strings.NewReplacer("created ", "would create ", "restyled ", "would restyle ").Replace(out.String())
		if dry.String() != want {
			t.Errorf("%s: dry run said\n%s\nwant\n%s", tc.name, dry.String(), want)
		}
	}
}

// A restyle GitHub refuses is a note like any other, and the labels after
// it are still attempted.
func TestApplyLabelsRestyleFailureIsANote(t *testing.T) {
	f := &labelFake{
		defined: []tracker.Label{implicit(tracker.LabelMilestone), implicit(tracker.LabelTask)},
		failOn:  tracker.LabelMilestone,
	}
	var out bytes.Buffer
	applyLabels(&out, f, "o/r", tracker.ProtocolLabels, true, false)
	if want := "note: could not restyle the cc:milestone label in o/r (403)"; !strings.Contains(out.String(), want) {
		t.Errorf("output missing %q:\n%s", want, out.String())
	}
	if got := names(f.restyled); !slices.Equal(got, []string{tracker.LabelTask}) {
		t.Errorf("restyled %v, want the label after the failure", got)
	}
	if got := names(f.created); !slices.Equal(got, []string{tracker.LabelNeedsDecision}) {
		t.Errorf("created %v", got)
	}
}
