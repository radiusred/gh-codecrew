package cli

import (
	"fmt"
	"io"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// labelStep is what a verb does to a label that is not already the
// protocol's: define it, or restyle the one that is there.
type labelStep struct {
	label   tracker.Label // the protocol's default — what the label becomes
	was     tracker.Label // the label as the repository has it, for a restyle
	restyle bool
}

// planLabels compares what the repository defines with what the protocol
// wants and returns the steps, in the order of want. A missing label is a
// creation. A label that exists is a restyle only when restyle is asked
// for and it is not already wearing the protocol's colour and description;
// otherwise it is left out of the plan entirely, which is how `init`'s
// never-restyle rule is expressed (SPEC §4).
func planLabels(existing, want []tracker.Label, restyle bool) []labelStep {
	var steps []labelStep
	for _, l := range want {
		got, ok := tracker.FindLabel(existing, l.Name)
		switch {
		case !ok:
			steps = append(steps, labelStep{label: l})
		case restyle && !got.Styled(l):
			steps = append(steps, labelStep{label: l, was: got, restyle: true})
		}
	}
	return steps
}

// applyLabels is the label step of every verb that has one. It answers the
// condition that made the first `cc:needs-decision` in a repository an
// untested path (#267): nothing created the protocol's labels, so
// `gh issue create --label` and `POST /issues/N/labels` created them
// implicitly, with whatever colour GitHub generated and no description.
//
// Three rules govern it:
//
//   - Restyling is one verb's business only. `migrate` is the one-shot
//     move to the 2.0 layout and its goal is a repository indistinguishable
//     from a fresh 2.0 `init`, so it sets the three to the protocol's
//     colour and description whatever they were wearing. `init` and
//     `checkpoint` create what is missing and touch nothing else: they run
//     on repositories that may have restyled a label deliberately, and
//     nothing there can tell that from a GitHub default (the operator's
//     Decision on #283).
//   - A dry run reads and prints and writes nothing.
//   - Nothing it does is a failure. Labels are a GitHub write standing
//     beside verbs whose real work is local, so a repository that is
//     offline, not on GitHub yet, or held by a token without
//     `issues: write` gets a `note:` line and the verb carries on. It
//     returns nothing, so no caller can treat it as fatal.
//
// It returns what it planned and whether it got to read the listing at
// all — never an error, so no caller can turn a GitHub refusal into a
// failure. The pair is enough for a caller to say "nothing to do", which
// only the rerun path needs (migrateLabelsRerun): planned is what the
// comparison asked for, not what succeeded, so a refused write leaves the
// note standing rather than reading as an empty plan.
func applyLabels(w io.Writer, t tracker.Tracker, repo string, want []tracker.Label, restyle, dryRun bool) (planned int, read bool) {
	existing, err := t.Labels(repo)
	if err != nil {
		fmt.Fprintf(w, "note: could not read %s's labels (%v) — the protocol's labels are created on first use instead, with GitHub's own colour\n", repo, err)
		return 0, false
	}
	steps := planLabels(existing, want, restyle)
	for _, s := range steps {
		verb := "created"
		if dryRun {
			verb = "would create"
		}
		if s.restyle {
			verb = "restyled"
			if dryRun {
				verb = "would restyle"
			}
		}
		if !dryRun {
			// A restyle addresses the label by the name the repository
			// spells it with — a recased `CC:Task` is the same label to
			// GitHub, and restyling is not renaming — carrying the
			// protocol's colour and description. A creation uses the
			// protocol's own spelling, there being nothing else.
			write, l := t.CreateLabel, s.label
			if s.restyle {
				write, l = t.UpdateLabel, tracker.Label{Name: s.was.Name, Color: s.label.Color, Description: s.label.Description}
			}
			if err := write(repo, l); err != nil {
				what := "create"
				if s.restyle {
					what = "restyle"
				}
				fmt.Fprintf(w, "note: could not %s the %s label in %s (%v)\n", what, s.label.Name, repo, err)
				continue
			}
		}
		if s.restyle {
			fmt.Fprintf(w, "%s label %s (#%s -> #%s) — %s\n", verb, s.label.Name, s.was.Color, s.label.Color, s.label.Description)
			continue
		}
		fmt.Fprintf(w, "%s label %s (#%s) — %s\n", verb, s.label.Name, s.label.Color, s.label.Description)
	}
	return len(steps), true
}

// ensureLabels creates the labels in want that repo does not already
// define, and leaves every existing one exactly as it is — colour and
// description included. It is what `init` and `checkpoint` do.
func ensureLabels(w io.Writer, t tracker.Tracker, repo string, want []tracker.Label) {
	applyLabels(w, t, repo, want, false, false)
}

// labelTarget resolves what the label step writes through and into: the
// venue, and the owner/repo it says the working directory belongs to.
// A func var for the reason defaultRequiresPR is one — the verbs that call
// it are tested without a GitHub behind them.
var labelTarget = func() (tracker.Tracker, string, error) {
	v := tracker.GitHub{}
	repo, err := v.CurrentRepo()
	return v, repo, err
}

// initLabels is init's whole label step, and it is deliberately the only
// thing the verb does to GitHub after the commit: the scaffold is files
// and a local commit, the labels are a remote write, and no failure of the
// second may reach the first. So it runs after the commit, reports every
// failure as a note, and skips outright in the two states init is
// explicitly built to work in — a directory that is not a git repository,
// and one whose repository GitHub does not know (no origin, no network, no
// credentials). Both are ordinary at scaffold time: `init` prints the
// `git init` and `gh repo create` lines for the first, and a rerun once
// the repository exists creates the labels then.
func initLabels(w io.Writer, dir string) {
	if repoRoot(dir) == "" {
		fmt.Fprintf(w, "note: the %s labels are created once this is a repository on GitHub — rerun gh codecrew init there\n", labelPrefix)
		return
	}
	withLabelTarget(w, func(t tracker.Tracker, repo string) {
		ensureLabels(w, t, repo, tracker.ProtocolLabels)
	})
}

// migrateLabels is migrate's label step: create the missing and restyle
// the rest, because the migration's whole promise is a repository that
// looks like a fresh 2.0 scaffold, and a 1.x repo's labels were all
// created implicitly with a random colour and no description. Like every
// other part of the step it is a note when GitHub says no.
//
// It runs on every path where the move reached disk — committed, detached,
// or the commit refused — so "the files moved" and "the labels were done"
// are never two different answers (the Decision on #283, from checky's
// finding 2 on PR #291).
func migrateLabels(w io.Writer, dryRun bool) {
	withLabelTarget(w, func(t tracker.Tracker, repo string) {
		applyLabels(w, t, repo, tracker.ProtocolLabels, true, dryRun)
	})
}

// migrateLabelsRerun is the same step on a repository already on the 2.0
// layout, and it exists because a rerun is the documented recovery from a
// label step that could not reach GitHub. Without it `migrate` returned at
// the idempotence check and that repository stayed on 2.0 wearing GitHub's
// grey with nothing in the tool that would ever fix it — the exact state
// M14-R5 exists to end (checky's finding 1 on PR #291). It says so when
// there is nothing to do, because on this path silence would be
// indistinguishable from the step not having run.
func migrateLabelsRerun(w io.Writer, dryRun bool) {
	withLabelTarget(w, func(t tracker.Tracker, repo string) {
		if planned, read := applyLabels(w, t, repo, tracker.ProtocolLabels, true, dryRun); read && planned == 0 {
			fmt.Fprintln(w, "labels already at the protocol defaults")
		}
	})
}

// withLabelTarget resolves the target and runs do, or prints the one note
// that covers every way gh can fail to name the repository — no origin, no
// network, no credentials.
func withLabelTarget(w io.Writer, do func(tracker.Tracker, string)) {
	t, repo, err := labelTarget()
	if err != nil {
		fmt.Fprintf(w, "note: could not ask GitHub which repository this is (%v) — the %s labels are created on first use instead, with GitHub's own colour\n", err, labelPrefix)
		return
	}
	do(t, repo)
}

// labelPrefix names the label family in the notes, so a message that has
// no one label to name still says which labels it means.
const labelPrefix = "cc:"

// needsDecisionLabel is the one-element want list checkpoint passes: the
// gate label alone, since raising a gate is no reason to define the other
// two in a repository that may never see them.
func needsDecisionLabel() []tracker.Label {
	l, _ := tracker.ProtocolLabel(tracker.LabelNeedsDecision)
	return []tracker.Label{l}
}
