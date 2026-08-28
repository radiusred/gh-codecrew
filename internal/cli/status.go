package cli

import (
	"fmt"
	"io"
	"os"

	codecrew "github.com/radiusred/gh-codecrew"

	"github.com/radiusred/gh-codecrew/internal/gh"
	"github.com/radiusred/gh-codecrew/internal/tracker"
)

func status(w io.Writer) error {
	cfg, err := loadConfig(".", os.Stderr)
	if err != nil {
		return err
	}
	current, err := gh.CurrentRepo()
	if err != nil {
		return err
	}
	hub := cfg.HubRepo(current)

	var t tracker.Tracker = tracker.GitHub{}
	milestones, err := t.OpenMilestones(hub)
	if err != nil {
		return err
	}
	if len(milestones) == 0 {
		fmt.Fprintf(w, "no open milestones in %s\n", hub)
		return nil
	}

	var gated []tracker.Task
	for _, m := range milestones {
		fmt.Fprintf(w, "%s (%s)\n", m.Title, m.Ref)
		if body, err := t.IssueBody(m.Ref); err == nil {
			if ids := tracker.RequirementIDs(body); len(ids) == 0 {
				fmt.Fprintf(w, "  %s\n", requirementsNote(ids))
			}
		}
		if len(m.Tasks) == 0 {
			fmt.Fprintln(w, "  no tasks yet")
		}
		for _, ref := range m.Tasks {
			task, err := t.Task(ref)
			if err != nil {
				return err
			}
			state := tracker.InferState(task)
			who := ""
			if state == tracker.InProgress || state == tracker.InReview {
				if len(task.Assignees) > 0 {
					who = " @" + task.Assignees[0]
				}
			}
			fmt.Fprintf(w, "  [%-11s] %-28s %s%s\n", state, ref, task.Title, who)
			if state == tracker.Gated {
				gated = append(gated, task)
			}
		}
		fmt.Fprintln(w)
	}

	if len(gated) == 0 {
		fmt.Fprintln(w, "gates raised: none")
	} else {
		fmt.Fprintln(w, "gates raised:")
		for _, task := range gated {
			fmt.Fprintf(w, "  %s — %s\n", task.Ref, task.Title)
		}
	}

	// The repo's own branch hygiene setting: advisory, like routing — the
	// verbs clean up regardless (task finish deletes the merged head,
	// milestone close sweeps), so an unreadable setting is skipped.
	if info, err := t.RepoInfo(current); err == nil && !info.DeleteBranchOnMerge {
		fmt.Fprintf(w, "note: %s does not delete branches on merge (GitHub setting) — task finish and milestone close clean up task branches; enable it for other PRs\n", current)
	}

	// Contract drift: purely local — the embedded contracts ride the
	// binary, so status can say when a hub's roles/ fork has diverged
	// from the installed release without touching the network.
	if drifted, err := contractDrift(cfg.Dir, codecrew.Roles); err == nil && len(drifted) > 0 {
		fmt.Fprintln(w)
		for _, role := range drifted {
			fmt.Fprintf(w, "contract drift: roles/%s.md differs from the embedded %s contract — codecrew roles diff %s\n", role, version, role)
		}
	}

	return nil
}
