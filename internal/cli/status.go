package cli

import (
	"fmt"
	"io"

	codecrew "github.com/radiusred/gh-codecrew"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

func status(w io.Writer) error {
	c, err := load()
	if err != nil {
		return err
	}
	return statusReport(w, c)
}

// gate is one entry in status's "gates raised" list: a task carrying
// cc:needs-decision, or a milestone issue carrying it — a question about a
// requirement has no task to carry it, so it is raised on the milestone
// issue and must show on the board like any other gate (#200).
type gate struct {
	ref       tracker.IssueRef
	title     string
	milestone bool
}

func statusReport(w io.Writer, c *ctx) error {
	milestones, err := c.t.OpenMilestones(c.hub)
	if err != nil {
		return err
	}
	if len(milestones) == 0 {
		fmt.Fprintf(w, "no open milestones in %s\n", c.hub)
		return nil
	}

	var gated []gate
	for _, m := range milestones {
		// The milestone issue's own labels: task states never reflect a
		// gate raised there, so it is read directly (#200). Task is a plain
		// issue query, so it serves for the milestone issue too.
		issue, err := c.t.Task(m.Ref)
		if err != nil {
			return err
		}
		header := fmt.Sprintf("%s (%s)", m.Title, m.Ref)
		if tracker.HasLabel(issue, tracker.LabelNeedsDecision) {
			header += " — gate raised on the milestone issue"
			gated = append(gated, gate{ref: m.Ref, title: m.Title, milestone: true})
		}
		fmt.Fprintln(w, header)
		if body, err := c.t.IssueBody(m.Ref); err == nil {
			if ids := tracker.RequirementIDs(body); len(ids) == 0 {
				fmt.Fprintf(w, "  %s\n", requirementsNote(ids))
			}
		}
		if len(m.Tasks) == 0 {
			fmt.Fprintln(w, "  no tasks yet")
		}
		for _, ref := range m.Tasks {
			task, err := c.t.Task(ref)
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
				gated = append(gated, gate{ref: ref, title: task.Title})
			}
		}
		fmt.Fprintln(w)
	}

	if len(gated) == 0 {
		fmt.Fprintln(w, "gates raised: none")
	} else {
		fmt.Fprintln(w, "gates raised:")
		for _, g := range gated {
			mark := ""
			if g.milestone {
				mark = " (milestone)"
			}
			fmt.Fprintf(w, "  %s — %s%s\n", g.ref, g.title, mark)
		}
	}

	// The repo's own branch hygiene setting: advisory, like routing — the
	// verbs clean up regardless (task finish deletes the merged head,
	// milestone close sweeps), so an unreadable setting is skipped.
	if info, err := c.t.RepoInfo(c.current); err == nil && !info.DeleteBranchOnMerge {
		fmt.Fprintf(w, "note: %s does not delete branches on merge (GitHub setting) — task finish and milestone close clean up task branches; enable it for other PRs\n", c.current)
	}

	// Contract drift: purely local — the embedded contracts ride the
	// binary, so status can say when a hub's roles/ fork has diverged
	// from the installed release without touching the network.
	if drifted, err := contractDrift(c.cfg.Dir, codecrew.Roles); err == nil && len(drifted) > 0 {
		fmt.Fprintln(w)
		for _, role := range drifted {
			fmt.Fprintf(w, "contract drift: roles/%s.md differs from the embedded %s contract — gh codecrew roles diff %s\n", role, version, role)
		}
	}

	return nil
}
