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
	ChecksPending  bool
	ChecksOK       bool
	ApprovedBy     []string
}

// NoChecksReported classifies gh's "no checks reported" failure mode:
// `gh pr checks` on a checkless PR prints nothing parseable and exits
// nonzero, so it surfaces as an error rather than an empty list.
func NoChecksReported(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no checks reported")
}

// Tracker is the backend interface, shaped by the workflow verbs. GitHub is
// the only implementation; the seam exists so a future backend stays possible.
type Tracker interface {
	// OpenMilestones returns the open cc:milestone issues in the hub repo.
	OpenMilestones(hub string) ([]Milestone, error)
	// AllMilestoneTitles returns titles of every cc:milestone issue, open or
	// closed, for milestone-number derivation.
	AllMilestoneTitles(hub string) ([]string, error)
	// Task fetches one task issue.
	Task(ref IssueRef) (Task, error)
	// IssueBody fetches an issue's body text.
	IssueBody(ref IssueRef) (string, error)
	// IssueLabels fetches the labels on an issue or a pull request — the
	// REST issues endpoint serves both, where Task's GraphQL issue query
	// answers NOT_FOUND for a PR; a gate may be recorded on the scaffold
	// PR (roles/coordinator.md), so checkpoint reads labels this way.
	IssueLabels(ref IssueRef) ([]string, error)
	// CreateIssue opens an issue and returns its ref.
	CreateIssue(repo, title, body string, labels []string) (IssueRef, error)
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

// PlanPlaceholder is the Plan section content task new writes; task start
// refuses while it is still in place.
const PlanPlaceholder = "_To be written by the implementer before the first commit._"

// PlanPresent reports whether the task body's Plan section has real content.
func PlanPresent(body string) bool {
	content := section(body, "## Plan")
	content = strings.ReplaceAll(content, PlanPlaceholder, "")
	return strings.TrimSpace(content) != ""
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
		body := strings.ReplaceAll(c.Body, "\r\n", "\n") // web-UI comments arrive CRLF
		for _, para := range paragraphBreak.Split(strings.TrimSpace(body), -1) {
			para = strings.TrimSpace(para)
			if para == "" {
				continue
			}
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

// Verdict is one QA requirement verdict found in a comment (roles/qa.md).
type Verdict struct {
	ID     string
	State  string // "satisfied", "not satisfied", or "untestable"
	Author string
}

// ParseVerdicts scans comments in order for verdict lines. Callers filter
// by author role and take the last entry per ID: a later verdict
// supersedes an earlier one.
func ParseVerdicts(comments []Comment) []Verdict {
	var verdicts []Verdict
	for _, c := range comments {
		for _, m := range verdictLine.FindAllStringSubmatch(c.Body, -1) {
			verdicts = append(verdicts, Verdict{
				ID:     m[1],
				State:  strings.ToLower(m[2]),
				Author: c.Author,
			})
		}
	}
	return verdicts
}

// UnresolvedGates returns the **Gate raised:** comments that have no later
// resolution record (a comment opening **Gate resolved:** or **Decision:**,
// qualified or bare). A single trailing
// resolution covers every gate raised before it — a human may answer several
// questions in one comment; the label removal remains the hard block.
func UnresolvedGates(comments []Comment) []Comment {
	var unresolved []Comment
	lastResolution := -1
	for i := len(comments) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(comments[i].Body)
		if m := recordLabel.FindStringSubmatch(trimmed); m != nil && m[1] != "Deviation" {
			lastResolution = i
			break
		}
	}
	for i, c := range comments {
		if strings.HasPrefix(strings.TrimSpace(c.Body), "**Gate raised:**") && i > lastResolution {
			unresolved = append(unresolved, c)
		}
	}
	return unresolved
}

// StartRecord is the exact comment task start posts for every start; a
// record is only what the named login posted itself.
func StartRecord(login string) string { return "**Started by** @" + login + "." }

// StartedBy returns the login that started a task, from the record task
// start leaves: the latest `**Started by** @<login>.` comment — accepted
// only when its body is exactly that record and its author is the login
// it names (a comment that merely begins with the phrase, or names
// someone else, is prose, not a record: checky's finding on PR #176).
// Every start posts one, so the latest is the current owner across
// restarts and handovers. The first assignee is the legacy fallback for
// a task started before the record existed. Empty when nothing records a
// start, so the ownership gate has nothing to compare.
func StartedBy(t Task, comments []Comment) string {
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
	if len(t.Assignees) > 0 {
		return t.Assignees[0]
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
