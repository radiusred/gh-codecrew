// Package tracker defines the backend interface, shaped by the workflow
// verbs rather than by any tracker's feature set (SPEC.md §10), and the pure
// protocol logic: task-ref parsing and state inference.
package tracker

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// IssueRef identifies an issue by repo and number.
type IssueRef struct {
	Repo   string // owner/repo
	Number int
}

func (r IssueRef) String() string {
	return fmt.Sprintf("%s#%d", r.Repo, r.Number)
}

// Milestone is a cc:milestone tracking issue in the hub.
type Milestone struct {
	Ref   IssueRef
	Title string
	Tasks []IssueRef
}

// TitledIssue is an issue's ref and title — all that milestone numbering
// reads from a listing.
type TitledIssue struct {
	Ref   IssueRef
	Title string
}

// Task is a cc:task issue in a spoke.
type Task struct {
	Ref          IssueRef
	Title        string
	Closed       bool
	Assignees    []string
	Labels       []string
	OpenLinkedPR bool
}

// State is an inferred task lifecycle state (SPEC.md §4).
type State string

const (
	Ready      State = "ready"
	InProgress State = "in progress"
	Gated      State = "gated"
	InReview   State = "in review"
	Done       State = "done"
)

// LabelNeedsDecision marks a raised human gate.
const LabelNeedsDecision = "cc:needs-decision"

// LabelMilestone marks a milestone tracking issue in the hub.
const LabelMilestone = "cc:milestone"

// Comment is one issue or PR comment.
type Comment struct {
	Author string
	Body   string
	URL    string
}

// PR is the review-surface state task finish gates on.
type PR struct {
	// ReviewDecision is GitHub's own verdict on required reviews: empty
	// when no rule applies, "APPROVED" when satisfied, "REVIEW_REQUIRED"
	// when the rule is not met by counted approvals — which come only from
	// write-access principals; a read-only App's review does not count
	// (the superseding Decision on #73).
	ReviewDecision string
	Repo           string
	Number         int
	Author         string
	HeadRef        string // the PR's head branch name
	HeadSHA        string // the head commit as GitHub last saw it — frozen at merge
	CrossRepo      bool   // head lives in another repo (a fork)
	Open           bool
	Merged         bool
	NoChecks       bool // zero CI checks reported — the deterministic gate cannot be satisfied by absence
	// ChecksUnreadable is the permission the installation token lacks to
	// read the checks at all ("checks: read" or "actions: read"), set in
	// place of the raw GraphQL error so the gate can name it (#198).
	ChecksUnreadable string
	ChecksPending    bool
	ChecksOK         bool
	ApprovedBy       []string
}

// NoChecksReported classifies gh's "no checks reported" failure mode:
// `gh pr checks` on a checkless PR prints nothing parseable and exits
// nonzero, so it surfaces as an error rather than an empty list.
func NoChecksReported(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no checks reported")
}

// MissingChecksPermission classifies the failure an installation token
// meets on a private repo when its App lacks a permission `gh pr checks`
// reads with: GitHub answers "Resource not accessible by integration" and
// names the GraphQL path it refused. Under statusCheckRollup that is
// checks: read — or actions: read when the path reaches
// checkSuite.workflowRun, the field gh reads for the workflow column
// (#198). It returns the missing permission, or "" for every other error:
// the same message on an unrelated path is not this failure.
func MissingChecksPermission(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if !strings.Contains(msg, "Resource not accessible by integration") {
		return ""
	}
	i := strings.Index(msg, "statusCheckRollup")
	if i < 0 {
		return ""
	}
	if strings.Contains(msg[i:], "checkSuite.workflowRun") {
		return "actions: read"
	}
	return "checks: read"
}

// Tracker is the backend interface, shaped by the workflow verbs. GitHub is
// the only implementation; the seam exists so a future backend stays possible.
type Tracker interface {
	// OpenMilestones returns the open cc:milestone issues in the hub repo.
	OpenMilestones(hub string) ([]Milestone, error)
	// MilestoneIssues returns every cc:milestone issue in the hub, open or
	// closed, with its title — the listing milestone-number derivation and
	// the post-create collision check read. The listing is label-filtered
	// and eventually consistent: it can lag an issue created seconds ago
	// (#195).
	MilestoneIssues(hub string) ([]TitledIssue, error)
	// RecentIssues returns the newest issues in repo, newest first, with no
	// label filter and pull requests excluded — the first page only. It
	// catches a milestone the label-filtered listing has not indexed yet.
	RecentIssues(repo string) ([]TitledIssue, error)
	// Task fetches one task issue.
	Task(ref IssueRef) (Task, error)
	// IssueBody fetches an issue's body text.
	IssueBody(ref IssueRef) (string, error)
	// IssueLabels fetches the labels on an issue or a pull request — the
	// REST issues endpoint serves both, where Task's GraphQL issue query
	// answers NOT_FOUND for a PR; a gate may be recorded on the scaffold
	// PR (.codecrew/roles/coordinator.md), so checkpoint reads labels this way.
	IssueLabels(ref IssueRef) ([]string, error)
	// CreateIssue opens an issue and returns its ref.
	CreateIssue(repo, title, body string, labels []string) (IssueRef, error)
	// EditIssue replaces an issue's title and body.
	EditIssue(ref IssueRef, title, body string) error
	// AddSubIssue attaches child to parent as a GitHub sub-issue; the parent
	// tracks progress natively, so nothing is hand-maintained.
	AddSubIssue(parent, child IssueRef) error
	// SubIssues lists the refs attached to parent as sub-issues.
	SubIssues(parent IssueRef) ([]IssueRef, error)
	// Comment posts an issue (or PR) comment.
	Comment(ref IssueRef, body string) error
	// AddLabel applies a label.
	AddLabel(ref IssueRef, label string) error
	// Assign assigns a login to an issue.
	Assign(ref IssueRef, login string) error
	// Viewer returns the login the current credentials act as.
	Viewer() (string, error)
	// DevelopBranch creates a branch linked to the issue.
	DevelopBranch(ref IssueRef, name string) error
	// ClosingPRs returns numbers of PRs that will close (or closed) the
	// issue, in the issue's own repo.
	ClosingPRs(ref IssueRef, includeClosed bool) ([]int, error)
	// PRInfo fetches the gate-relevant state of one PR.
	PRInfo(repo string, number int) (PR, error)
	// MergePR rebase-merges a PR.
	MergePR(repo string, number int) error
	// MergePRBypass rebase-merges with the ruleset's administrator bypass
	// (task finish --bypass; fails with GitHub's own error when the caller
	// is not a bypass actor).
	MergePRBypass(repo string, number int) error
	// MergeCommit returns the commit a merged PR left on the base branch,
	// or "" when it has none yet. A rebase merge rewrites the head, so the
	// head SHA is not that commit — task finish names this one when it
	// closes the captures the task adopted.
	MergeCommit(repo string, number int) (string, error)
	// CloseIssue closes an issue with a closing comment.
	CloseIssue(ref IssueRef, comment string) error
	// Comments lists issue (or PR) comments.
	Comments(ref IssueRef) ([]Comment, error)
	// HasMilestoneDoc reports whether docs/milestones/<n>-*.md exists on the
	// default branch of repo.
	HasMilestoneDoc(repo string, n int) (bool, error)
	// FileContent fetches a file from the default branch of repo.
	FileContent(repo, path string) ([]byte, error)
	// LinkedBranches lists the branch names linked to an issue — the
	// relation task start creates through gh issue develop.
	LinkedBranches(ref IssueRef) ([]string, error)
	// TaskBranches lists repo's `task/<n>-<slug>` branch names, filtered at
	// the server by ref prefix: one listing per repo, and no branch outside
	// the protocol's own naming can enter a sweep's candidate set.
	TaskBranches(repo string) ([]string, error)
	// BranchAhead reports how many commits branch carries beyond repo's
	// default branch and the branch's current tip; an error when the branch
	// does not exist.
	BranchAhead(repo, branch string) (ahead int, sha string, err error)
	// DeleteBranch deletes a branch ref.
	DeleteBranch(repo, branch string) error
	// RepoInfo fetches the repo settings the verbs consult.
	RepoInfo(repo string) (RepoInfo, error)
}

// RepoInfo is the slice of repository settings the verbs read.
type RepoInfo struct {
	DefaultBranch       string
	DeleteBranchOnMerge bool
}

// InferState derives a task's lifecycle state from tracker signals, most
// terminal first: Done > Gated > In review > In progress > Ready.
func InferState(t Task) State {
	switch {
	case t.Closed:
		return Done
	case hasLabel(t, LabelNeedsDecision):
		return Gated
	case t.OpenLinkedPR:
		return InReview
	case len(t.Assignees) > 0:
		return InProgress
	default:
		return Ready
	}
}

func hasLabel(t Task, name string) bool { return ContainsLabel(t.Labels, name) }

// HasLabel reports whether the task carries the label.
func HasLabel(t Task, name string) bool { return hasLabel(t, name) }

// ContainsLabel reports whether name is among labels, case-insensitively,
// as GitHub compares label names.
func ContainsLabel(labels []string, name string) bool {
	for _, l := range labels {
		if strings.EqualFold(l, name) {
			return true
		}
	}
	return false
}

var refPattern = regexp.MustCompile(`^(?:([\w.-]+/[\w.-]+))?#?(\d+)$`)

// ParseRef parses "12", "#12", or "owner/repo#12"; bare and short forms
// resolve against defaultRepo.
func ParseRef(s, defaultRepo string) (IssueRef, error) {
	m := refPattern.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return IssueRef{}, fmt.Errorf("bad issue ref %q (want N, #N, or owner/repo#N)", s)
	}
	repo := m[1]
	if repo == "" {
		repo = defaultRepo
	}
	n, _ := strconv.Atoi(m[2])
	return IssueRef{Repo: repo, Number: n}, nil
}

var milestoneTitle = regexp.MustCompile(`^M(\d+)\s*:`)

// MilestoneNumber extracts n from a milestone title of the form "M<n>: ...".
func MilestoneNumber(title string) (int, bool) {
	m := milestoneTitle.FindStringSubmatch(strings.TrimSpace(title))
	if m == nil {
		return 0, false
	}
	n, _ := strconv.Atoi(m[1])
	return n, true
}

// NextMilestoneNumber derives the next milestone number from existing
// milestone issue titles.
func NextMilestoneNumber(titles []string) int {
	max := 0
	for _, t := range titles {
		if n, ok := MilestoneNumber(t); ok && n > max {
			max = n
		}
	}
	return max + 1
}

// MilestoneNumberHolder returns the first issue in the listings, other than
// self, whose title carries "M<n>:" — the collision `milestone new` checks
// for after creating — or nil when the number is self's alone.
func MilestoneNumberHolder(n int, self IssueRef, listings ...[]TitledIssue) *TitledIssue {
	for _, issues := range listings {
		for i := range issues {
			if issues[i].Ref == self {
				continue
			}
			if got, ok := MilestoneNumber(issues[i].Title); ok && got == n {
				return &issues[i]
			}
		}
	}
	return nil
}

// Titles flattens listings into the titles NextMilestoneNumber reads.
func Titles(listings ...[]TitledIssue) []string {
	var titles []string
	for _, issues := range listings {
		for _, is := range issues {
			titles = append(titles, is.Title)
		}
	}
	return titles
}

// PlanPlaceholder is the Plan section content task new writes; task start
// refuses while it is still in place.
const PlanPlaceholder = "_To be written by the implementer before the first commit._"

// PlanPresent reports whether the task body's Plan section has real content.
func PlanPresent(body string) bool {
	content := section(body, "## Plan")
	content = strings.ReplaceAll(content, PlanPlaceholder, "")
	return strings.TrimSpace(content) != ""
}

// AdoptsHeading opens the task-body section listing the backlog issues a
// task adopts (SPEC §4). `task new --adopts` writes it; `task finish`
// reads it and closes each capture once the merge is done.
const AdoptsHeading = "## Adopts"

// Adoption is one backlog issue a task adopts: the capture's ref, and its
// title as it read at the moment of adoption.
type Adoption struct {
	Ref   IssueRef
	Title string
}

// AdoptsBlock renders the ## Adopts section for a task body in taskRepo —
// a leading blank line, the heading, one list line per capture — or the
// empty string when a task adopts nothing, so a body without adoptions is
// byte for byte the body task new has always written. A ref in the task's
// own repo is written short (`#193`), one elsewhere in full; the title
// follows it for the reader.
func AdoptsBlock(taskRepo string, adopted []Adoption) string {
	if len(adopted) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n" + AdoptsHeading + "\n")
	for _, a := range adopted {
		b.WriteString("- " + ShortRef(a.Ref, taskRepo))
		if t := strings.TrimSpace(a.Title); t != "" {
			b.WriteString(" — " + t)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// ShortRef writes a ref the way a body in repo reads it: `#N` in the same
// repo, `owner/repo#N` anywhere else.
func ShortRef(ref IssueRef, repo string) string {
	if ref.Repo == repo {
		return fmt.Sprintf("#%d", ref.Number)
	}
	return ref.String()
}

// adoptedLine matches one entry of the ## Adopts section: a list marker,
// then the ref. Only the ref at the head of the line is one — the prose
// after it is the capture's title, so a title carrying a `#42` of its own
// is not a second adoption, and nothing task finish closes was ever
// merely quoted.
var adoptedLine = regexp.MustCompile(`(?m)^[ \t]*[-*][ \t]+((?:[\w.-]+/[\w.-]+)?#\d+)\b`)

// AdoptedRefs returns the captures listed under a task body's ## Adopts
// section, in the order written, deduplicated; a bare `#N` resolves
// against defaultRepo, the task's own repo. No section, or none of its
// lines carrying a ref, is no adoptions.
//
// The body is read through StripCode first, as the verdict scan and the
// citation walk read a comment (M13-R6): a ref quoted in a fenced block,
// an indented block or an inline span is content, not an adoption. That
// rule matters more here than anywhere it already held — what this returns
// is what task finish closes after a merge, where nothing can refuse — so
// the section is found by heading rather than by substring too
// (adoptsSection).
func AdoptedRefs(body, defaultRepo string) []IssueRef {
	var refs []IssueRef
	seen := map[IssueRef]bool{}
	for _, m := range adoptedLine.FindAllStringSubmatch(adoptsSection(StripCode(body)), -1) {
		ref, err := ParseRef(m[1], defaultRepo)
		if err != nil || seen[ref] {
			continue
		}
		seen[ref] = true
		refs = append(refs, ref)
	}
	return refs
}

// adoptsHeadingLine matches the ## Adopts heading as a heading — at the
// head of a line, up to three columns of indentation as CommonMark allows,
// the text exact, trailing spaces ignored — and sectionEnd matches the
// heading of level 1 or 2 that ends it.
var (
	adoptsHeadingLine = regexp.MustCompile(`(?m)^ {0,3}` + AdoptsHeading + `[ \t]*$`)
	sectionEnd        = regexp.MustCompile(`(?m)^ {0,3}#{1,2} `)
)

// adoptsSection returns the text under the ## Adopts heading. section()
// finds its heading by substring, which is enough where a miss costs a
// presence check (PlanPresent) or an ID list (RequirementIDs); here it is
// not, because prose that merely quotes the heading would shadow the real
// section and decide what task finish closes — and a task about this
// feature quotes it (#270's own Goal does). So the heading is matched as a
// heading and the section runs to the next one of level 1 or 2.
func adoptsSection(body string) string {
	loc := adoptsHeadingLine.FindStringIndex(body)
	if loc == nil {
		return ""
	}
	rest := body[loc[1]:]
	if end := sectionEnd.FindStringIndex(rest); end != nil {
		rest = rest[:end[0]]
	}
	return rest
}

// AdoptionRecord is the comment task new posts on a capture, so the
// backlog issue itself says which task now carries it.
func AdoptionRecord(task IssueRef) string {
	return fmt.Sprintf("**Adopted by** %s — this issue is delivered as that task, and `gh codecrew task finish` closes it when the task's pull request merges (SPEC §4). The PR needs no `Closes` line for it.", task)
}

// AdoptionClose is the comment task finish closes a capture with: the task
// that adopted it, the pull request that delivered it, and the commit the
// merge left on the default branch when that could be read.
func AdoptionClose(task, pr IssueRef, sha string) string {
	msg := fmt.Sprintf("Closed by `gh codecrew task finish %d`: adopted by %s and delivered by %s", task.Number, task, pr)
	if sha != "" {
		msg += fmt.Sprintf(", merged as %s", sha)
	}
	return msg + "."
}

// section returns the text between the given heading and the next "## ".
func section(body, heading string) string {
	_, rest, found := strings.Cut(body, heading)
	if !found {
		return ""
	}
	if i := strings.Index(rest, "\n## "); i >= 0 {
		rest = rest[:i]
	}
	return rest
}

// Record is one Decision or Deviation captured in a comment.
type Record struct {
	Kind   string // "Decision" or "Deviation"
	Label  string // the label as written, e.g. "**Decision (superseding …):**"
	Source string // the issue/PR ref the comment was found on
	Author string
	Body   string
	URL    string
}

// recordLabel matches a record label at the start of a paragraph: the bare
// SPEC §4 form (**Decision:**) or a qualified one (**Decision (…):**) — the
// qualifier is kept verbatim in Record.Label and carries no semantics. A
// gate resolution (**Gate resolved:**, SPEC §8) is a decision made at a
// human gate, so it is captured as a Decision record.
var recordLabel = regexp.MustCompile(`^\*\*(Decision|Deviation|Gate resolved)(\s*\([^\n]*?\))?:\*\*`)

// gateRaisedLabel matches a raised gate at the start of a paragraph, and
// gateResolvedLabel its resolution — the same placement rule the record
// labels obey, and the same optional parenthetical qualifier (SPEC §4).
// **Gate raised:** is not a record: it states a question, and only its
// **Gate resolved:** answer is gathered as a Decision.
var (
	gateRaisedLabel   = regexp.MustCompile(`^\*\*Gate raised(\s*\([^\n]*?\))?:\*\*`)
	gateResolvedLabel = regexp.MustCompile(`^\*\*Gate resolved(\s*\([^\n]*?\))?:\*\*`)
)

// continuationLabel opens a paragraph that belongs to the record before it.
var continuationLabel = regexp.MustCompile(`^\*\*(Why|Trade-off|Rejected):\*\*`)

// otherLabel is any other bold label opening a paragraph — it ends the
// record before it without starting one.
var otherLabel = regexp.MustCompile(`^\*\*[^*\n]+:\*\*`)

var paragraphBreak = regexp.MustCompile(`\n[ \t]*\n`)

// ExtractRecords finds Decision/Deviation records in comments per the SPEC
// §4 convention, one record per labelled paragraph: a paragraph opening
// with a record label starts a record; the unlabelled and **Why:** /
// **Trade-off:** / **Rejected:** paragraphs after it belong to it until the
// next label. A comment that is one record is one record; a record written
// after other text in the same comment (a review round-up that ends with a
// Deviation) is still gathered. A label mentioned mid-line is not a record.
func ExtractRecords(source IssueRef, comments []Comment) []Record {
	var records []Record
	for _, c := range comments {
		var open *Record
		for _, para := range paragraphs(c.Body) {
			if m := recordLabel.FindStringSubmatch(para); m != nil {
				kind := m[1]
				if kind == "Gate resolved" {
					kind = "Decision"
				}
				records = append(records, Record{Kind: kind, Label: m[0], Source: source.String(), Author: c.Author, Body: para, URL: c.URL})
				open = &records[len(records)-1]
				continue
			}
			if open == nil {
				continue
			}
			if otherLabel.MatchString(para) && !continuationLabel.MatchString(para) {
				open = nil
				continue
			}
			open.Body += "\n\n" + para
		}
	}
	return records
}

// requirementID matches a bold requirement ID as written in milestone
// bodies (**M3-R1**); verdictLine matches the QA verdict convention with
// the state inside the bold (**M3-R1 — satisfied.**), so requirement
// definition lines never parse as verdicts.
var (
	requirementID = regexp.MustCompile(`\*\*(M\d+-R\d+)\*\*`)
	verdictLine   = regexp.MustCompile(`(?i)\*\*(M\d+-R\d+)\s*[—–-]+\s*(satisfied|not satisfied|untestable)\b`)
)

// RequirementIDs extracts the ordered, deduplicated requirement IDs from a
// milestone body's Requirements section.
func RequirementIDs(body string) []string {
	var ids []string
	seen := map[string]bool{}
	for _, m := range requirementID.FindAllStringSubmatch(section(body, "## Requirements"), -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			ids = append(ids, m[1])
		}
	}
	return ids
}

// MismatchedRequirementIDs returns the IDs under a milestone body's
// Requirements section whose milestone number is not n, in the order they
// appear — SPEC §4's grammar made enforceable: a requirement ID is
// M<milestone>-R<k>, and M13's requirements are M13-R1, M13-R2, … The
// callers refuse REQUIREMENT_ID_MISMATCH on a non-empty result (status
// prints it instead: it reports rather than gates). An ID whose number
// matches by a prefix only — M1-R1 under M13 — is a mismatch: the number
// is compared as a number, not as text.
func MismatchedRequirementIDs(body string, n int) []string {
	var bad []string
	for _, id := range RequirementIDs(body) {
		num, _, _ := strings.Cut(strings.TrimPrefix(id, "M"), "-")
		if got, err := strconv.Atoi(num); err != nil || got != n {
			bad = append(bad, id)
		}
	}
	return bad
}

// Verdict is one QA requirement verdict found in a comment
// (.codecrew/roles/qa.md).
type Verdict struct {
	ID     string
	State  string // "satisfied", "not satisfied", or "untestable"
	Author string
}

// ParseVerdicts scans comments in order for verdict lines, at most one per
// requirement ID per comment — the first match in the comment — with all
// three Markdown code forms stripped first (StripCode: spans, fenced blocks
// and blocks indented four columns anywhere they do not continue a
// paragraph), so a verdict quoted back inside code is content, not a
// verdict, exactly as a URL in code is not a citation.
//
// Supersession is therefore per comment: callers filter by author role and
// take the last entry per ID, which now means "the latest comment carrying
// a verdict for that ID wins, and within it the first match counts"
// (SPEC §6, M13-R6). Until protocol 2.0 the last match anywhere won, so a
// QA comment that restated an earlier verdict below its own superseded it.
func ParseVerdicts(comments []Comment) []Verdict {
	var verdicts []Verdict
	for _, c := range comments {
		seen := map[string]bool{}
		for _, m := range verdictLine.FindAllStringSubmatch(StripCode(c.Body), -1) {
			if seen[m[1]] {
				continue
			}
			seen[m[1]] = true
			verdicts = append(verdicts, Verdict{
				ID:     m[1],
				State:  strings.ToLower(m[2]),
				Author: c.Author,
			})
		}
	}
	return verdicts
}

// paragraphs splits a comment body the way ExtractRecords reads it: on a
// blank line, CRLF normalised first (web-UI comments arrive CRLF), each
// paragraph trimmed and the empty ones dropped.
func paragraphs(body string) []string {
	var out []string
	for _, para := range paragraphBreak.Split(strings.TrimSpace(strings.ReplaceAll(body, "\r\n", "\n")), -1) {
		if para = strings.TrimSpace(para); para != "" {
			out = append(out, para)
		}
	}
	return out
}

// UnresolvedGates returns, in the order they were raised, the comments
// carrying a **Gate raised:** paragraph that no later **Gate resolved:**
// answers — one entry per gate, so a comment raising two gates that are
// never answered appears twice.
//
// Two rules, both tightened in protocol 2.0 (M13-R6):
//
//   - Placement. A gate is recognised per paragraph, anywhere in a
//     comment, exactly as the record labels are (SPEC §4). Until 2.0 only
//     a comment whose body opened with the label counted, so a gate raised
//     as a comment's second paragraph did not exist to task finish.
//   - Resolution. Only a **Gate resolved:** paragraph resolves, and it
//     resolves every gate still open before it — a human may answer
//     several questions in one comment — never the gates raised after it.
//     A bare **Decision:** resolves nothing, which is SPEC §8's wording;
//     until 2.0 any Decision anywhere later cleared every gate at once.
//
// The cc:needs-decision label remains the hard block; this is the record
// behind it.
func UnresolvedGates(comments []Comment) []Comment {
	var open []Comment
	for _, c := range comments {
		for _, para := range paragraphs(c.Body) {
			switch {
			case gateRaisedLabel.MatchString(para):
				open = append(open, c)
			case gateResolvedLabel.MatchString(para):
				open = nil
			}
		}
	}
	return open
}

// StartRecord is the exact comment task start posts for every start; a
// record is only what the named login posted itself.
func StartRecord(login string) string { return "**Started by** @" + login + "." }

// StartedBy returns the login that started a task, from the one record
// task start leaves: the latest `**Started by** @<login>.` comment —
// accepted only when its body is exactly that record and its author is
// the login it names (a comment that merely begins with the phrase, or
// names someone else, is prose, not a record: checky's finding on PR
// #176). Every start posts one, so the latest is the current owner across
// restarts and handovers.
//
// Empty means nothing records a start, and under 2.0 that is a fact about
// the task rather than a gap to be filled: the first assignee was the 1.0
// fallback "for tasks started before the record existed", which gave every
// assigned-but-never-started task an implicit owner, and it is deleted
// with the rest of the shims (M13-R7). A task with no start record has no
// owner, and task finish refuses it — the fix is to run task start (SPEC
// §6, §8).
func StartedBy(comments []Comment) string {
	for i := len(comments) - 1; i >= 0; i-- {
		body := strings.TrimSpace(comments[i].Body)
		rest, ok := strings.CutPrefix(body, "**Started by** @")
		if !ok || !strings.HasSuffix(rest, ".") {
			continue
		}
		login := strings.TrimSuffix(rest, ".")
		if login == "" || strings.ContainsAny(login, " \n\t") || !SameLogin(comments[i].Author, login) {
			continue
		}
		return login
	}
	return ""
}

// SameLogin compares two GitHub logins the way the routing table does:
// the "[bot]" suffix an App token's viewer carries is not part of the
// identity, and logins are case-insensitive.
func SameLogin(a, b string) bool {
	norm := func(s string) string { return strings.ToLower(strings.TrimSuffix(strings.TrimPrefix(s, "@"), "[bot]")) }
	return norm(a) == norm(b)
}
