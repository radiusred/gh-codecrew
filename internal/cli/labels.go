package cli

import (
	"fmt"
	"io"

	"github.com/radiusred/gh-codecrew/internal/gh"
	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// ensureLabels creates the labels in want that repo does not already
// define, and reports each creation. It answers the condition that made
// the first `cc:needs-decision` in a repository an untested path (#267):
// nothing created the protocol's labels, so `gh issue create --label` and
// `POST /issues/N/labels` created them implicitly, with whatever colour
// GitHub generated and no description at all.
//
// Two rules govern it, and both are about not owning what the project
// owns:
//
//   - A label that already exists is left exactly as it is — its colour
//     and its description too. A project may have restyled it
//     deliberately, and nothing here can tell that from a GitHub default
//     (SPEC §4).
//   - Nothing it does is a failure. Labels are a GitHub write standing
//     beside verbs that are not, so a repository that is offline, not on
//     GitHub yet, or held by a token without `issues: write` gets a
//     `note:` line and the verb carries on. It returns nothing, so no
//     caller can treat it as fatal.
func ensureLabels(w io.Writer, t tracker.Tracker, repo string, want []tracker.Label) {
	existing, err := t.Labels(repo)
	if err != nil {
		fmt.Fprintf(w, "note: could not read %s's labels (%v) — the protocol's labels are created on first use instead, with GitHub's own colour\n", repo, err)
		return
	}
	for _, l := range want {
		if tracker.ContainsLabel(existing, l.Name) {
			continue
		}
		if err := t.CreateLabel(repo, l); err != nil {
			fmt.Fprintf(w, "note: could not create the %s label in %s (%v) — it is created on first use instead, with GitHub's own colour\n", l.Name, repo, err)
			continue
		}
		fmt.Fprintf(w, "created label %s (#%s) — %s\n", l.Name, l.Color, l.Description)
	}
}

// labelTarget resolves what ensureLabels writes through and into: the
// tracker, and the owner/repo `gh` says the working directory belongs to.
// A func var for the reason defaultRequiresPR is one — the verbs that call
// it are tested without a GitHub behind them.
var labelTarget = func() (tracker.Tracker, string, error) {
	repo, err := gh.CurrentRepo()
	return tracker.GitHub{}, repo, err
}

// initLabels is init's whole label step, and it is deliberately the last
// thing the verb does that touches GitHub: the scaffold is files and a
// local commit, the labels are a remote write, and no failure of the
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
	t, repo, err := labelTarget()
	if err != nil {
		fmt.Fprintf(w, "note: could not ask GitHub which repository this is (%v) — the %s labels are created on first use instead, with GitHub's own colour\n", err, labelPrefix)
		return
	}
	ensureLabels(w, t, repo, tracker.ProtocolLabels)
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
