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
	// No open milestone replaces the board, not the report: the two checks
	// below are local and have nothing to do with milestone state, and the
	// quiet period between milestones is exactly when an operator
	// reconciles a .codecrew/roles/ fork against a new release (#253).
	if len(milestones) == 0 {
		fmt.Fprintf(w, "no open milestones in %s\n", c.hub)
	} else if err := milestoneBoard(w, c, milestones); err != nil {
		return err
	}

	// The repo's own branch hygiene setting: advisory, like routing — the
	// verbs clean up regardless (task finish deletes the merged head,
	// milestone close sweeps), so an unreadable setting is skipped.
	if info, err := c.t.RepoInfo(c.current); err == nil && !info.DeleteBranchOnMerge {
		fmt.Fprintf(w, "note: %s does not delete branches on merge (GitHub setting) — task finish and milestone close clean up task branches; enable it for other PRs\n", c.current)
	}

	// Contract drift: purely local — the embedded contracts ride the
	// binary, so status can say when a hub's .codecrew/roles/ fork has
	// diverged from the installed release without touching the network.
	if drifted, err := contractDrift(c.cfg.Dir, codecrew.Roles); err == nil && len(drifted) > 0 {
		fmt.Fprintln(w)
		for _, role := range drifted {
			fmt.Fprintf(w, "contract drift: %s differs from the embedded %s contract — gh codecrew roles diff %s\n", contractPath(role), version, role)
		}
	}

	return nil
}

// milestoneBoard prints the board itself — every open milestone with its
// tasks, then the gates raised across them. It is the part of status that
// a hub between milestones has nothing to say for.
func milestoneBoard(w io.Writer, c *ctx, milestones []tracker.Milestone) error {
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
			// The Requirements section, read the way milestone close
			// reads it: empty, or declaring an ID that is not this
			// milestone's (M13-R6). status reports both as a line and
			// carries on — it is the board, not a gate.
			num, numbered := tracker.MilestoneNumber(m.Title)
			if ids := tracker.RequirementIDs(body); len(ids) == 0 {
				fmt.Fprintf(w, "  %s\n", requirementsNote(ids))
			} else if numbered {
				if bad := tracker.MismatchedRequirementIDs(body, num); len(bad) > 0 {
					fmt.Fprintf(w, "  %s\n", requirementIDMismatchNote(num, bad))
				}
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
				if login := taskHolder(c, ref, task); login != "" {
					who = " @" + login
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

	return nil
}

// taskHolder names the seat holding a task that is in flight: the login
// from its latest `**Started by**` record, which is the only thing that
// says a task was started and the same login task finish holds to
// (NOT_OWNER). Reading the record rather than the assignee list is what
// lets an App-run task name its holder at all — GitHub does not accept an
// App as an assignee, so those tasks have no assignee to show (#287) — and
// after a handover it names the seat that took the task over rather than
// whoever was assigned first.
//
// The first assignee remains the fallback, for display only: a task that
// nothing records a start on is in progress precisely because it carries an
// assignee (SPEC §4's lifecycle table), so a state with no name beside it
// would contradict the signal that produced it. It is not an ownership
// signal — task finish still refuses a task with no start record.
//
// The comments read is one per task and only in the two states that have a
// holder; a read that fails is not worth failing the board for, so it falls
// through to the assignee.
func taskHolder(c *ctx, ref tracker.IssueRef, task tracker.Task) string {
	if comments, err := c.t.Comments(ref); err == nil {
		if login := tracker.StartedBy(comments); login != "" {
			return login
		}
	}
	if len(task.Assignees) > 0 {
		return task.Assignees[0]
	}
	return ""
}
