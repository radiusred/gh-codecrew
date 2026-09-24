package cli

import (
	"fmt"
	"io"
	"io/fs"

	codecrew "github.com/radiusred/gh-codecrew"

	"github.com/radiusred/gh-codecrew/internal/config"
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
	// milestone close sweeps), so an unreadable setting is skipped. The
	// same read serves the stale-branch report below, which needs the
	// default branch's name.
	info, err := c.t.RepoInfo(c.current)
	if err == nil && !info.DeleteBranchOnMerge {
		fmt.Fprintf(w, "note: %s does not delete branches on merge (GitHub setting) — task finish and milestone close clean up task branches; enable it for other PRs\n", c.current)
	}
	if err == nil {
		staleBranches(w, c.t, c.current, info.DefaultBranch)
	} else {
		// The two are not symmetric on this error. The delete-on-merge
		// note's absence claims nothing — it prints only when the setting
		// is off. The report's absence claims the repo carries no stale
		// branch, and that claim must not be made by a read that failed.
		fmt.Fprintf(w, "note: stale task branches not listed for %s (%v)\n", c.current, err)
	}

	// Contract drift and absence: purely local — the embedded contracts
	// ride the binary, so status can say when a hub's .codecrew/roles/ has
	// diverged from the installed release, or lacks a contract the release
	// carries (#266), without touching the network. Only a hub holds
	// contracts; a spoke is pointer-only and has nothing to report.
	if c.cfg.Hub == "self" {
		contractReport(w, c.cfg.Dir, codecrew.Roles, contractHistory)
	}
	// The agents file is in hub and spoke alike, and so is its line (#372).
	agentsReport(w, c.cfg.Dir, agentsScaffold, contractHistory)

	return nil
}

// contractReport prints one line per hub contract that is not the
// embedded text, naming the verb that acts on it: roles sync for a missing
// contract and for an earlier release's text, which it can write; roles
// diff for a fork, which it never touches (SPEC §7). Nothing to report
// prints nothing, and a classification that fails is not worth failing
// the board for.
func contractReport(w io.Writer, dir string, contracts fs.FS, history []releasedContract) {
	statuses, err := classifyContracts(dir, contracts, history, nil)
	if err != nil {
		return
	}
	first := true
	for _, s := range statuses {
		if s.State == contractCurrent {
			continue
		}
		if first {
			fmt.Fprintln(w)
			first = false
		}
		p := contractPath(s.Role)
		switch s.State {
		case contractAbsent:
			fmt.Fprintf(w, "contract missing: %s — the embedded %s contract has no local copy; gh codecrew roles sync writes it\n", p, version)
		case contractRelease:
			fmt.Fprintf(w, "contract drift: %s is the %s text, behind the embedded %s contract — gh codecrew roles sync\n", p, s.Release, version)
		case contractForked:
			fmt.Fprintf(w, "contract drift: %s differs from the embedded %s contract and from every release's, a fork — gh codecrew roles diff %s\n", p, version, s.Role)
		}
	}
}

// agentsReport prints one line when .codecrew/AGENTS.md is not the
// scaffold this binary writes, naming the verb that acts on it: roles sync
// when it is absent or a release's scaffold, roles diff when it is the
// project's own — which roles sync never overwrites (#372).
func agentsReport(w io.Writer, dir, embedded string, history []releasedContract) {
	st, err := classifyAgents(dir, embedded, history)
	if err != nil || st.State == contractCurrent {
		return
	}
	fmt.Fprintln(w)
	p := config.AgentsFile
	switch st.State {
	case contractAbsent:
		fmt.Fprintf(w, "agents file missing: %s — gh codecrew roles sync writes the embedded %s scaffold\n", p, version)
	case contractRelease:
		fmt.Fprintf(w, "agents file drift: %s is the %s text, behind the embedded %s scaffold — gh codecrew roles sync\n", p, st.Release, version)
	case contractForked:
		fmt.Fprintf(w, "agents file drift: %s differs from the embedded %s scaffold and from every release's, the project's own — gh codecrew roles diff %s\n", p, version, p)
	}
}

// staleBranches reports the task branches this repo still carries whose
// task issue has closed — a sweep that never ran, or one that declined
// (#295). `milestone close` reaches these branches only when a milestone
// closes; between closes nothing looks, and this is the only place a solo
// operator with no open milestone would ever see one. The verdict and its
// reason come from staleBranchAction, the close's own second pass, so the
// report can never promise a deletion the close would decline.
//
// Every line is advisory, like the two notes it sits beside: a listing
// GitHub will not give up is one note: and the verb carries on. Nothing to
// report says nothing at all.
func staleBranches(w io.Writer, t tracker.Tracker, repo, defaultBranch string) {
	items, listed, truncated, err := planStaleReport(t, repo, defaultBranch)
	if err != nil {
		fmt.Fprintf(w, "note: stale task branches not listed for %s (%v)\n", repo, err)
		return
	}
	for _, it := range items {
		switch {
		case it.Note != "":
			// Worded by the sweep, printed unchanged: an unreadable task
			// issue is named the same way wherever it is met.
			fmt.Fprintln(w, it.Note)
		case it.TaskOpen:
			// Nothing has finished with this branch; it is not stale.
		case it.Delete:
			fmt.Fprintf(w, "stale branch: %s — %s is closed; %s — the next milestone close would delete it\n", it.Name, staleRef(it), it.Reason)
		default:
			fmt.Fprintf(w, "stale branch: %s — %s is closed; %s — kept by the next milestone close\n", it.Name, staleRef(it), it.Reason)
		}
	}
	if truncated {
		fmt.Fprintf(w, "note: %s carries more task branches than one listing holds; this report saw %d of them\n", repo, listed)
	}
}

// staleRef names the task issue a reported branch belongs to, read back
// out of the branch name by the one function the sweep reads it with.
func staleRef(it sweepItem) tracker.IssueRef {
	return tracker.IssueRef{Repo: it.Repo, Number: taskNumber(it.Name)}
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
			strikeLines(w, c, &m, body)
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

// strikeLines reports the milestone's strikes as milestone close counts
// them, through the same function (M18-R4): a verified strike as a line, an
// unverified one as the note naming the code close will refuse with, and a
// still-bold ID struck through in the body with no strike behind it as a
// note that the protocol does not read the body's strikethrough. status
// reports rather than gates, so a comments read that fails is a note too.
func strikeLines(w io.Writer, c *ctx, m *tracker.Milestone, body string) {
	ids := tracker.RequirementIDs(body)
	if len(ids) == 0 {
		return
	}
	comments, err := c.t.Comments(m.Ref)
	var states map[string]strikeState
	if err == nil {
		states, err = milestoneStrikes(c, m, comments)
	}
	if err != nil {
		fmt.Fprintf(w, "  note: strikes not read for %s (%v)\n", m.Ref, err)
		return
	}
	struck, unrecorded := struckIDs(ids, states)
	isStruck := map[string]bool{}
	for _, id := range struck {
		isStruck[id] = true
		fmt.Fprintf(w, "  struck: %s — decision %s\n", id, states[id].Decision)
	}
	for _, u := range unrecorded {
		fmt.Fprintf(w, "  note: a strike does not verify — milestone close refuses DECISION_UNRECORDED: %s\n", u)
	}
	for _, id := range tracker.StruckThroughIDs(body) {
		if !isStruck[id] {
			fmt.Fprintf(w, "  note: %s is struck through in the body, which the protocol does not read — it is still a requirement; strike it with gh codecrew milestone strike and a recorded Decision (SPEC §4)\n", id)
		}
	}
}
