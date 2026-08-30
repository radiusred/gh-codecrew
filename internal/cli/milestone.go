package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// Tasks have no body section: they are attached as GitHub sub-issues, which
// the milestone issue tracks natively.
const milestoneTemplate = `## Goal
%s

## Requirements
%s

## Gates
_What "done" means beyond CI: e2e suites, manual UAT, sign-offs._
`

// requirementsPlaceholder is what the Requirements section holds when
// `milestone new` is given no --requirement; it yields no IDs by design.
func requirementsPlaceholder(n int) string {
	return fmt.Sprintf(`_One line per requirement, its ID in bold: M%d-R1, M%d-R2, … — only
bolded IDs in this section count as requirements, so this placeholder
does not, and neither do IDs written under Goal or Gates._`, n, n)
}

// multiFlag collects a repeatable string flag in the order given.
type multiFlag []string

func (m *multiFlag) String() string     { return strings.Join(*m, "; ") }
func (m *multiFlag) Set(v string) error { *m = append(*m, v); return nil }

// titlePrefix matches a title that carries its own milestone number.
var titlePrefix = regexp.MustCompile(`^M(\d+)\s*(?::|—|–|-)\s*`)

// milestoneTitle writes "M<n>: <title>". The CLI derives n, so a title that
// names a different number is refused rather than doubled: the orchestrator
// run titled a milestone "M2 — …" while the CLI counted 3 and created
// "M3: M2 — …", closed as a duplicate (#147). A prefix that agrees is
// stripped.
func milestoneTitle(title string, n int) (string, error) {
	title = strings.TrimSpace(title)
	m := titlePrefix.FindStringSubmatch(title)
	if m == nil {
		return fmt.Sprintf("M%d: %s", n, title), nil
	}
	k, _ := strconv.Atoi(m[1])
	if k != n {
		return "", fmt.Errorf("milestone new: the title says M%d but the next milestone number is %d — drop the prefix; the CLI writes \"M%d: <title>\"", k, n, n)
	}
	rest := strings.TrimSpace(title[len(m[0]):])
	if rest == "" {
		return "", fmt.Errorf("milestone new: the title is only a milestone number")
	}
	return fmt.Sprintf("M%d: %s", n, rest), nil
}

// requirementPrefix matches requirement text that already carries an ID.
var requirementPrefix = regexp.MustCompile(`^\**M\d+-R\d+\**`)

// milestoneBody renders the tracking issue. Each --requirement becomes one
// bold-ID line under ## Requirements, numbered M<n>-R<i> in the order given
// — the section the parser reads, so a coordinator never has to write IDs
// under Goal (#147; #119 findings 19a and 32). Text that brings its own ID
// is refused: two numbering schemes in one body is how #144 happened.
func milestoneBody(goal string, n int, reqs []string) (string, error) {
	if len(reqs) == 0 {
		return fmt.Sprintf(milestoneTemplate, goal, requirementsPlaceholder(n)), nil
	}
	lines := make([]string, 0, len(reqs))
	for i, r := range reqs {
		r = strings.TrimSpace(r)
		if r == "" {
			return "", fmt.Errorf("milestone new: --requirement %d is empty", i+1)
		}
		if requirementPrefix.MatchString(r) {
			return "", fmt.Errorf("milestone new: --requirement %q carries an ID; the CLI numbers requirements M%d-R1, M%d-R2, … in the order given — pass the text only", r, n, n)
		}
		lines = append(lines, fmt.Sprintf("- **M%d-R%d** — %s", n, i+1, r))
	}
	return fmt.Sprintf(milestoneTemplate, goal, strings.Join(lines, "\n")), nil
}

func milestoneNew(w io.Writer, args []string) error {
	fs := flag.NewFlagSet("milestone new", flag.ContinueOnError)
	title := fs.String("title", "", "milestone title (required)")
	goal := fs.String("goal", "_To be written._", "one-paragraph goal")
	var reqs multiFlag
	fs.Var(&reqs, "requirement", "a requirement's text; repeatable, numbered M<n>-R1, R2, … in order")
	dryRun := fs.Bool("dry-run", false, "print the number, title, requirement IDs and ROADMAP row the milestone would get; create nothing")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *title == "" {
		return fmt.Errorf("milestone new: --title is required")
	}
	c, err := load()
	if err != nil {
		return err
	}
	titles, err := c.t.AllMilestoneTitles(c.hub)
	if err != nil {
		return err
	}
	n := tracker.NextMilestoneNumber(titles)
	fullTitle, err := milestoneTitle(*title, n)
	if err != nil {
		return err
	}
	body, err := milestoneBody(*goal, n, reqs)
	if err != nil {
		return err
	}
	if *dryRun {
		// The number is the point: it is assigned after the requirements
		// are written, and a closed duplicate still counts toward it — so
		// prose can be written knowing it instead of guessing (M7-R5).
		fmt.Fprintf(w, "would create milestone M%d in %s: %s\n", n, c.hub, fullTitle)
		fmt.Fprintln(w, requirementsNote(tracker.RequirementIDs(body)))
		fmt.Fprintf(w, "would add to ROADMAP.md: | M%d | %s | #<issue> | Open |\n", n, strings.TrimPrefix(fullTitle, fmt.Sprintf("M%d: ", n)))
		fmt.Fprintln(w, "dry run: nothing written")
		return nil
	}
	ref, err := c.t.CreateIssue(c.hub, fullTitle, body, []string{"cc:milestone"})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "created milestone %s (%s)\n", fullTitle, ref)
	fmt.Fprintln(w, requirementsNote(tracker.RequirementIDs(body)))

	row := fmt.Sprintf("| M%d | %s | [#%d](https://github.com/%s/issues/%d) | Open |",
		n, strings.TrimPrefix(fullTitle, fmt.Sprintf("M%d: ", n)), ref.Number, ref.Repo, ref.Number)
	if c.current == c.hub {
		if appendRoadmapRow(filepath.Join(c.cfg.Dir, "ROADMAP.md"), row) == nil {
			fmt.Fprintln(w, "ROADMAP.md updated locally — it rides in this milestone's first PR (the implementer's)")
			return nil
		}
	}
	fmt.Fprintf(w, "add to the hub's ROADMAP.md:\n%s\n", row)
	return nil
}

// appendRoadmapRow adds the row after the roadmap table's last row.
func appendRoadmapRow(path, row string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	last := -1
	for i, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "|") {
			last = i
		}
	}
	if last < 0 {
		return fmt.Errorf("no table in %s", path)
	}
	lines = append(lines[:last+1], append([]string{row}, lines[last+1:]...)...)
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644)
}

// closeGates, in the order milestone close checks them.
var closeGates = []string{"milestone open", "tasks closed", "requirements declared", "QA verdicts", "milestone document"}

func milestoneClose(w io.Writer, args []string) error {
	fs := flag.NewFlagSet("milestone close", flag.ContinueOnError)
	dryRun := fs.Bool("dry-run", false, "print every gate, the branches the sweep would delete or keep, and the closing comment; write nothing; exit with the first refusal's code")
	nArg, args := splitLeadingRef(args)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if nArg == "" && fs.NArg() == 1 {
		nArg = fs.Arg(0)
	}
	if nArg == "" || fs.NArg() > 1 {
		return fmt.Errorf("usage: gh codecrew milestone close <milestone number> [--dry-run]")
	}
	n, err := strconv.Atoi(strings.TrimPrefix(nArg, "M"))
	if err != nil {
		return fmt.Errorf("bad milestone number %q", nArg)
	}
	c, err := load()
	if err != nil {
		return err
	}
	p, run, err := planClose(c, n, *dryRun, w)
	if err != nil {
		return err
	}
	if *dryRun {
		p.print(w)
		return p.refusal
	}
	if p.refusal != nil {
		return p.refusal
	}
	return run(w)
}

// planClose evaluates milestone close's gates in order, plans the sweep and
// the closing comment, and writes nothing; run performs them. The
// doc-synthesizer's raw material is printed to w on DOC_MISSING by the live
// verb only — the dry run names the gate.
func planClose(c *ctx, n int, dryRun bool, w io.Writer) (*plan, func(io.Writer) error, error) {
	p := &plan{}
	milestones, err := c.t.OpenMilestones(c.hub)
	if err != nil {
		return nil, nil, err
	}
	var milestone *tracker.Milestone
	for i := range milestones {
		if got, ok := tracker.MilestoneNumber(milestones[i].Title); ok && got == n {
			milestone = &milestones[i]
			break
		}
	}
	var e error
	if milestone == nil {
		e = refuse("NOT_FOUND", "no open milestone M%d in %s", n, c.hub)
	}
	if !p.gate("milestone open", e) {
		return p.stop(closeGates), nil, nil
	}

	// Gate 1: every task closed.
	var open []string
	for _, ref := range milestone.Tasks {
		task, err := c.t.Task(ref)
		if err != nil {
			return nil, nil, err
		}
		if !task.Closed {
			open = append(open, fmt.Sprintf("%s (%s)", ref, tracker.InferState(task)))
		}
	}
	e = nil
	if len(open) > 0 {
		e = refuse("OPEN_TASKS", "tasks not closed: %s", strings.Join(open, ", "))
	}
	if !p.gate("tasks closed", e) {
		return p.stop(closeGates), nil, nil
	}

	// Gate 2: the milestone declares requirements where the parser reads
	// them. A close that would verify zero requirements is not a close —
	// numberguess M1 closed over an unsatisfied verdict this way (#144).
	body, err := c.t.IssueBody(milestone.Ref)
	if err != nil {
		return nil, nil, err
	}
	ids := tracker.RequirementIDs(body)
	if !p.gate("requirements declared", requireRequirements(milestone.Ref, ids)) {
		return p.stop(closeGates), nil, nil
	}

	// Gate 3: every requirement carries a satisfied QA verdict (a later
	// verdict supersedes an earlier one; only the qa role's holder counts —
	// its routed identity, or the human operator when the role is unrouted).
	comments, err := c.t.Comments(milestone.Ref)
	if err != nil {
		return nil, nil, err
	}
	latest := map[string]string{}
	for _, v := range tracker.ParseVerdicts(comments) {
		if c.holdsRole(v.Author, "qa") {
			latest[v.ID] = v.State
		}
	}
	var missing, unsatisfied []string
	for _, id := range ids {
		state, ok := latest[id]
		switch {
		case !ok:
			missing = append(missing, id)
		case state != "satisfied":
			unsatisfied = append(unsatisfied, fmt.Sprintf("%s (%s)", id, state))
		}
	}
	e = nil
	switch {
	case len(missing) > 0:
		e = refuse("VERDICT_MISSING", "no QA verdict on %s for: %s — dispatch QA (roles/qa.md)", milestone.Ref, strings.Join(missing, ", "))
	case len(unsatisfied) > 0:
		e = refuse("VERDICT_UNSATISFIED", "latest QA verdict not satisfied for: %s — remedy and re-dispatch QA", strings.Join(unsatisfied, ", "))
	}
	if !p.gate("QA verdicts", e) {
		return p.stop(closeGates), nil, nil
	}
	if len(latest) > 0 && c.rolesConfig().Roles["qa"].Identity == "" {
		// Said here, before the document gate, as it always was — a live
		// close refused DOC_MISSING still says it (checky's finding on
		// PR #179); the dry run shows it under the gate.
		note := "note: qa is unrouted — verdicts counted from the human operator holding the role; declare role routing in the hub's .codecrew.yml at onboarding (SPEC §5)"
		if dryRun {
			p.note(note)
		} else {
			fmt.Fprintln(w, note)
		}
	}

	// Gather Decision/Deviation raw material for the doc-synthesizer.
	records, summaries, err := gatherRecords(c, milestone)
	if err != nil {
		return nil, nil, err
	}

	// Gate 4: the synthesized milestone document is merged.
	hasDoc, err := c.t.HasMilestoneDoc(c.hub, n)
	if err != nil {
		return nil, nil, err
	}
	e = nil
	if !hasDoc {
		if !dryRun {
			fmt.Fprintf(w, "raw material for docs/milestones/%d-<slug>.md (%d records):\n\n", n, len(records))
			writeRecords(w, records)
			fmt.Fprintf(w, "task PRs (open and merged — their descriptions are the task summaries, read directly, never gathered; SPEC §4):\n")
			for _, line := range summaries {
				fmt.Fprintf(w, "  %s\n", line)
			}
			fmt.Fprintln(w)
		}
		e = docMissing(n, c.hub)
	}
	if !p.gate("milestone document", e) {
		return p.stop(closeGates), nil, nil
	}

	// Every gate has passed: sweep the tasks' branches — merged or empty
	// ones go, anything else is reported — so the successful close, and its
	// closing comment, carry what was removed (M6-R8). The sweep runs
	// before the close on purpose: it only ever deletes a merged or empty
	// branch, so a close that fails after it loses nothing (#133).
	items := planSweep(c.t, milestone)
	var wouldDelete []string
	for _, it := range items {
		switch {
		case it.Note != "":
			p.would("skip: %s", strings.TrimPrefix(it.Note, "note: "))
		case it.Delete:
			wouldDelete = append(wouldDelete, it.Name)
			p.would("delete branch %s (%s)", it.Name, it.Reason)
		default:
			p.would("keep branch %s (%s)", it.Name, it.Reason)
		}
	}
	closing := func(swept []string) string {
		comment := fmt.Sprintf("Closed by `gh codecrew milestone close %d`: all %d tasks done, milestone document merged.", n, len(milestone.Tasks))
		if len(swept) > 0 {
			comment += fmt.Sprintf(" Swept %d task branch(es): %s.", len(swept), strings.Join(swept, ", "))
		}
		return comment
	}
	p.would("close %s (%s) with: %s", milestone.Title, milestone.Ref, closing(wouldDelete))
	run := func(w io.Writer) error {
		swept := executeSweep(w, c.t, items)
		if err := c.t.CloseIssue(milestone.Ref, closing(swept)); err != nil {
			return err
		}
		fmt.Fprintf(w, "closed %s (%s)\n", milestone.Title, milestone.Ref)
		return nil
	}
	return p, run, nil
}

// gatherRecords collects Decision/Deviation records from every task issue
// in the milestone and each task's PRs (open and closed), and lists each
// task's PRs so the doc-synthesizer has the summary pointers without
// another walk. PR bodies are never gathered (SPEC §4).
func gatherRecords(c *ctx, m *tracker.Milestone) ([]tracker.Record, []string, error) {
	var records []tracker.Record
	var summaries []string
	for _, ref := range m.Tasks {
		comments, err := c.t.Comments(ref)
		if err != nil {
			return nil, nil, err
		}
		records = append(records, tracker.ExtractRecords(ref, comments)...)
		prs, err := c.t.ClosingPRs(ref, true)
		if err != nil {
			return nil, nil, err
		}
		var prRefs []string
		for _, num := range prs {
			prRef := tracker.IssueRef{Repo: ref.Repo, Number: num}
			prComments, err := c.t.Comments(prRef)
			if err != nil {
				return nil, nil, err
			}
			records = append(records, tracker.ExtractRecords(prRef, prComments)...)
			prRefs = append(prRefs, prRef.String())
		}
		if len(prRefs) > 0 {
			summaries = append(summaries, fmt.Sprintf("%s: %s", ref, strings.Join(prRefs, ", ")))
		} else {
			summaries = append(summaries, fmt.Sprintf("%s: no PR", ref))
		}
	}
	return records, summaries, nil
}

func writeRecords(w io.Writer, records []tracker.Record) {
	for _, r := range records {
		label := r.Label
		if label == "" {
			label = r.Kind
		}
		fmt.Fprintf(w, "---\n%s on %s by @%s (%s)\n\n%s\n\n", label, r.Source, r.Author, r.URL, r.Body)
	}
}

// requireRequirements is milestone close's second gate: a body whose
// "## Requirements" section yields no bold IDs has nothing to verdict, so
// the close refuses instead of passing vacuously. IDs written anywhere
// else in the body are not requirements (the parser is section-scoped so
// the template's placeholder can never become a phantom requirement).
func requireRequirements(ref tracker.IssueRef, ids []string) error {
	if len(ids) > 0 {
		return nil
	}
	return refuse("NO_REQUIREMENTS", "the Requirements section of %s yields no bold IDs (bold IDs elsewhere in the body do not count) — put each requirement under ## Requirements as **M<n>-R<k>** and rerun", ref)
}

// requirementsNote is the line new, status and evidence print so the
// section's emptiness is visible long before close refuses over it.
func requirementsNote(ids []string) string {
	if len(ids) == 0 {
		return "note: the Requirements section yields no bold IDs — edit the issue and list each requirement under ## Requirements as **M<n>-R<k>** (IDs elsewhere in the body do not count); milestone close refuses NO_REQUIREMENTS otherwise"
	}
	return fmt.Sprintf("requirements counted: %s (%d)", strings.Join(ids, ", "), len(ids))
}

// docMissing is the close's last refusal. Its detail names the task path
// rather than "merge its PR": a document PR with no task behind it has no
// owner for its review loop and nothing that can merge it — the orchestrator
// run's coordinator, sent there by the old wording, planned to merge by hand
// with an App that could not (#119 finding 27).
func docMissing(n int, hub string) error {
	return refuse("DOC_MISSING", "docs/milestones/%d-*.md not on the default branch of %s — dispatch the doc-synthesizer as a task: it writes the plan, runs task start, opens the PR with Closes #<task>, and task finish merges it; then rerun", n, hub)
}
