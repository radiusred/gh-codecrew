// Package faketracker is the one fake venue: a tracker.Tracker that records
// every call made through it and answers from whatever the test scripted.
//
// It exists because the venue is one interface (M16-R1, #194), so a fake of
// it is worth writing once. It is an ordinary package rather than a
// _test.go file so that internal/cli can import it; nothing outside a test
// has any reason to.
//
// A test sets only the fields for the methods it exercises. Every other
// method still answers — the zero value of its results, and no error — and
// every call, scripted or not, is recorded in order with its arguments, so
// a test can assert what was asked as well as what came back. The venue's
// pure knowledge (tracker.RecordAPIPath and the rest of the URL grammar,
// tracker.Unreachable) is deliberately not here: it is not I/O, there is
// nothing to script, and a fake that could lie about it would let a test
// pass against a grammar the venue does not have.
package faketracker

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// Call is one call made through the fake: the method's name and its
// arguments, in order.
type Call struct {
	Method string
	Args   []any
}

// String renders a call the way a failure message wants to read it —
// `Task(o/r#7)`.
func (c Call) String() string {
	parts := make([]string, len(c.Args))
	for i, a := range c.Args {
		parts[i] = fmt.Sprint(a)
	}
	return c.Method + "(" + strings.Join(parts, ", ") + ")"
}

// Venue is the fake. Its zero value is usable: every method answers with
// zero values and records the call.
type Venue struct {
	mu    sync.Mutex
	calls []Call

	// One field per method: set, it answers; nil, the method returns zero
	// values. Named <Method>Fn so the pair is obvious at the call site.
	OpenMilestonesFn        func(string) ([]tracker.Milestone, error)
	MilestoneIssuesFn       func(string) ([]tracker.TitledIssue, error)
	RecentIssuesFn          func(string) ([]tracker.TitledIssue, error)
	TaskFn                  func(tracker.IssueRef) (tracker.Task, error)
	IssueBodyFn             func(tracker.IssueRef) (string, error)
	IssueLabelsFn           func(tracker.IssueRef) ([]string, error)
	CreateIssueFn           func(string, string, string, []string) (tracker.IssueRef, error)
	EditIssueFn             func(tracker.IssueRef, string, string) error
	AddSubIssueFn           func(tracker.IssueRef, tracker.IssueRef) error
	SubIssuesFn             func(tracker.IssueRef) ([]tracker.IssueRef, error)
	CommentFn               func(tracker.IssueRef, string) error
	AddLabelFn              func(tracker.IssueRef, string) error
	LabelsFn                func(string) ([]tracker.Label, error)
	CreateLabelFn           func(string, tracker.Label) error
	UpdateLabelFn           func(string, tracker.Label) error
	AssignFn                func(tracker.IssueRef, string) error
	ViewerFn                func() (string, error)
	DevelopBranchFn         func(tracker.IssueRef, string) error
	ClosingPRsFn            func(tracker.IssueRef, bool) ([]int, error)
	PRInfoFn                func(string, int) (tracker.PR, error)
	ClosingReferencesFn     func(string, int) ([]tracker.TitledIssue, error)
	MergePRFn               func(string, int) error
	MergePRBypassFn         func(string, int) error
	MergeCommitFn           func(string, int) (string, error)
	CloseIssueFn            func(tracker.IssueRef, string) error
	CommentsFn              func(tracker.IssueRef) ([]tracker.Comment, error)
	HasMilestoneDocFn       func(string, int) (bool, error)
	FileContentFn           func(string, string) ([]byte, error)
	LinkedBranchesFn        func(tracker.IssueRef) ([]string, error)
	TaskBranchesFn          func(string) ([]string, bool, error)
	OpenPRsForBranchFn      func(string, string) ([]int, error)
	BranchAheadFn           func(string, string) (int, string, error)
	DeleteBranchFn          func(string, string) error
	RepoInfoFn              func(string) (tracker.RepoInfo, error)
	ClientVersionFn         func() (string, error)
	CurrentRepoFn           func() (string, error)
	RepoInDirFn             func(string) (string, string, error)
	BranchRuleTypesFn       func(string, string) ([]string, error)
	TeamMembersFn           func(string, string) ([]string, error)
	AccountTypeFn           func(string) (string, error)
	APIReachableFn          func(string) error
	AppManifestConversionFn func(string) (*tracker.AppCredentials, error)
	AppRequestFn            func(*http.Client, string, string, string, any) (int, []byte, error)
}

// record appends one call.
func (v *Venue) record(method string, args ...any) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.calls = append(v.calls, Call{Method: method, Args: args})
}

// Calls returns every call made so far, in order.
func (v *Venue) Calls() []Call {
	v.mu.Lock()
	defer v.mu.Unlock()
	return append([]Call(nil), v.calls...)
}

// CallsTo returns the calls made to one method, in order.
func (v *Venue) CallsTo(method string) []Call {
	var out []Call
	for _, c := range v.Calls() {
		if c.Method == method {
			out = append(out, c)
		}
	}
	return out
}

// Count is how many times a method was called.
func (v *Venue) Count(method string) int { return len(v.CallsTo(method)) }

// Reset forgets every recorded call and leaves the scripted answers in
// place.
func (v *Venue) Reset() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.calls = nil
}

// The compiler is the guard that the fake stays a whole venue as the
// interface grows.
var _ tracker.Tracker = (*Venue)(nil)

func (v *Venue) OpenMilestones(hub string) ([]tracker.Milestone, error) {
	v.record("OpenMilestones", hub)
	if v.OpenMilestonesFn != nil {
		return v.OpenMilestonesFn(hub)
	}
	return nil, nil
}

func (v *Venue) MilestoneIssues(hub string) ([]tracker.TitledIssue, error) {
	v.record("MilestoneIssues", hub)
	if v.MilestoneIssuesFn != nil {
		return v.MilestoneIssuesFn(hub)
	}
	return nil, nil
}

func (v *Venue) RecentIssues(repo string) ([]tracker.TitledIssue, error) {
	v.record("RecentIssues", repo)
	if v.RecentIssuesFn != nil {
		return v.RecentIssuesFn(repo)
	}
	return nil, nil
}

func (v *Venue) Task(ref tracker.IssueRef) (tracker.Task, error) {
	v.record("Task", ref)
	if v.TaskFn != nil {
		return v.TaskFn(ref)
	}
	return tracker.Task{}, nil
}

func (v *Venue) IssueBody(ref tracker.IssueRef) (string, error) {
	v.record("IssueBody", ref)
	if v.IssueBodyFn != nil {
		return v.IssueBodyFn(ref)
	}
	return "", nil
}

func (v *Venue) IssueLabels(ref tracker.IssueRef) ([]string, error) {
	v.record("IssueLabels", ref)
	if v.IssueLabelsFn != nil {
		return v.IssueLabelsFn(ref)
	}
	return nil, nil
}

func (v *Venue) CreateIssue(repo string, title string, body string, labels []string) (tracker.IssueRef, error) {
	v.record("CreateIssue", repo, title, body, labels)
	if v.CreateIssueFn != nil {
		return v.CreateIssueFn(repo, title, body, labels)
	}
	return tracker.IssueRef{}, nil
}

func (v *Venue) EditIssue(ref tracker.IssueRef, title string, body string) error {
	v.record("EditIssue", ref, title, body)
	if v.EditIssueFn != nil {
		return v.EditIssueFn(ref, title, body)
	}
	return nil
}

func (v *Venue) AddSubIssue(parent tracker.IssueRef, child tracker.IssueRef) error {
	v.record("AddSubIssue", parent, child)
	if v.AddSubIssueFn != nil {
		return v.AddSubIssueFn(parent, child)
	}
	return nil
}

func (v *Venue) SubIssues(parent tracker.IssueRef) ([]tracker.IssueRef, error) {
	v.record("SubIssues", parent)
	if v.SubIssuesFn != nil {
		return v.SubIssuesFn(parent)
	}
	return nil, nil
}

func (v *Venue) Comment(ref tracker.IssueRef, body string) error {
	v.record("Comment", ref, body)
	if v.CommentFn != nil {
		return v.CommentFn(ref, body)
	}
	return nil
}

func (v *Venue) AddLabel(ref tracker.IssueRef, label string) error {
	v.record("AddLabel", ref, label)
	if v.AddLabelFn != nil {
		return v.AddLabelFn(ref, label)
	}
	return nil
}

func (v *Venue) Labels(repo string) ([]tracker.Label, error) {
	v.record("Labels", repo)
	if v.LabelsFn != nil {
		return v.LabelsFn(repo)
	}
	return nil, nil
}

func (v *Venue) CreateLabel(repo string, label tracker.Label) error {
	v.record("CreateLabel", repo, label)
	if v.CreateLabelFn != nil {
		return v.CreateLabelFn(repo, label)
	}
	return nil
}

func (v *Venue) UpdateLabel(repo string, label tracker.Label) error {
	v.record("UpdateLabel", repo, label)
	if v.UpdateLabelFn != nil {
		return v.UpdateLabelFn(repo, label)
	}
	return nil
}

func (v *Venue) Assign(ref tracker.IssueRef, login string) error {
	v.record("Assign", ref, login)
	if v.AssignFn != nil {
		return v.AssignFn(ref, login)
	}
	return nil
}

func (v *Venue) Viewer() (string, error) {
	v.record("Viewer")
	if v.ViewerFn != nil {
		return v.ViewerFn()
	}
	return "", nil
}

func (v *Venue) DevelopBranch(ref tracker.IssueRef, name string) error {
	v.record("DevelopBranch", ref, name)
	if v.DevelopBranchFn != nil {
		return v.DevelopBranchFn(ref, name)
	}
	return nil
}

func (v *Venue) ClosingPRs(ref tracker.IssueRef, includeClosed bool) ([]int, error) {
	v.record("ClosingPRs", ref, includeClosed)
	if v.ClosingPRsFn != nil {
		return v.ClosingPRsFn(ref, includeClosed)
	}
	return nil, nil
}

func (v *Venue) PRInfo(repo string, number int) (tracker.PR, error) {
	v.record("PRInfo", repo, number)
	if v.PRInfoFn != nil {
		return v.PRInfoFn(repo, number)
	}
	return tracker.PR{}, nil
}

func (v *Venue) ClosingReferences(repo string, number int) ([]tracker.TitledIssue, error) {
	v.record("ClosingReferences", repo, number)
	if v.ClosingReferencesFn != nil {
		return v.ClosingReferencesFn(repo, number)
	}
	return nil, nil
}

func (v *Venue) MergePR(repo string, number int) error {
	v.record("MergePR", repo, number)
	if v.MergePRFn != nil {
		return v.MergePRFn(repo, number)
	}
	return nil
}

func (v *Venue) MergePRBypass(repo string, number int) error {
	v.record("MergePRBypass", repo, number)
	if v.MergePRBypassFn != nil {
		return v.MergePRBypassFn(repo, number)
	}
	return nil
}

func (v *Venue) MergeCommit(repo string, number int) (string, error) {
	v.record("MergeCommit", repo, number)
	if v.MergeCommitFn != nil {
		return v.MergeCommitFn(repo, number)
	}
	return "", nil
}

func (v *Venue) CloseIssue(ref tracker.IssueRef, comment string) error {
	v.record("CloseIssue", ref, comment)
	if v.CloseIssueFn != nil {
		return v.CloseIssueFn(ref, comment)
	}
	return nil
}

func (v *Venue) Comments(ref tracker.IssueRef) ([]tracker.Comment, error) {
	v.record("Comments", ref)
	if v.CommentsFn != nil {
		return v.CommentsFn(ref)
	}
	return nil, nil
}

func (v *Venue) HasMilestoneDoc(repo string, n int) (bool, error) {
	v.record("HasMilestoneDoc", repo, n)
	if v.HasMilestoneDocFn != nil {
		return v.HasMilestoneDocFn(repo, n)
	}
	return false, nil
}

func (v *Venue) FileContent(repo string, path string) ([]byte, error) {
	v.record("FileContent", repo, path)
	if v.FileContentFn != nil {
		return v.FileContentFn(repo, path)
	}
	return nil, nil
}

func (v *Venue) LinkedBranches(ref tracker.IssueRef) ([]string, error) {
	v.record("LinkedBranches", ref)
	if v.LinkedBranchesFn != nil {
		return v.LinkedBranchesFn(ref)
	}
	return nil, nil
}

func (v *Venue) TaskBranches(repo string) ([]string, bool, error) {
	v.record("TaskBranches", repo)
	if v.TaskBranchesFn != nil {
		return v.TaskBranchesFn(repo)
	}
	return nil, false, nil
}

func (v *Venue) OpenPRsForBranch(repo string, branch string) ([]int, error) {
	v.record("OpenPRsForBranch", repo, branch)
	if v.OpenPRsForBranchFn != nil {
		return v.OpenPRsForBranchFn(repo, branch)
	}
	return nil, nil
}

func (v *Venue) BranchAhead(repo string, branch string) (int, string, error) {
	v.record("BranchAhead", repo, branch)
	if v.BranchAheadFn != nil {
		return v.BranchAheadFn(repo, branch)
	}
	return 0, "", nil
}

func (v *Venue) DeleteBranch(repo string, branch string) error {
	v.record("DeleteBranch", repo, branch)
	if v.DeleteBranchFn != nil {
		return v.DeleteBranchFn(repo, branch)
	}
	return nil
}

func (v *Venue) RepoInfo(repo string) (tracker.RepoInfo, error) {
	v.record("RepoInfo", repo)
	if v.RepoInfoFn != nil {
		return v.RepoInfoFn(repo)
	}
	return tracker.RepoInfo{}, nil
}

func (v *Venue) ClientVersion() (string, error) {
	v.record("ClientVersion")
	if v.ClientVersionFn != nil {
		return v.ClientVersionFn()
	}
	return "", nil
}

func (v *Venue) CurrentRepo() (string, error) {
	v.record("CurrentRepo")
	if v.CurrentRepoFn != nil {
		return v.CurrentRepoFn()
	}
	return "", nil
}

func (v *Venue) RepoInDir(dir string) (string, string, error) {
	v.record("RepoInDir", dir)
	if v.RepoInDirFn != nil {
		return v.RepoInDirFn(dir)
	}
	return "", "", nil
}

func (v *Venue) BranchRuleTypes(repo string, branch string) ([]string, error) {
	v.record("BranchRuleTypes", repo, branch)
	if v.BranchRuleTypesFn != nil {
		return v.BranchRuleTypesFn(repo, branch)
	}
	return nil, nil
}

func (v *Venue) TeamMembers(org string, team string) ([]string, error) {
	v.record("TeamMembers", org, team)
	if v.TeamMembersFn != nil {
		return v.TeamMembersFn(org, team)
	}
	return nil, nil
}

func (v *Venue) AccountType(login string) (string, error) {
	v.record("AccountType", login)
	if v.AccountTypeFn != nil {
		return v.AccountTypeFn(login)
	}
	return "", nil
}

func (v *Venue) APIReachable(path string) error {
	v.record("APIReachable", path)
	if v.APIReachableFn != nil {
		return v.APIReachableFn(path)
	}
	return nil
}

func (v *Venue) AppManifestConversion(code string) (*tracker.AppCredentials, error) {
	v.record("AppManifestConversion", code)
	if v.AppManifestConversionFn != nil {
		return v.AppManifestConversionFn(code)
	}
	return nil, nil
}

func (v *Venue) AppRequest(client *http.Client, jwt string, method string, path string, body any) (int, []byte, error) {
	v.record("AppRequest", client, jwt, method, path, body)
	if v.AppRequestFn != nil {
		return v.AppRequestFn(client, jwt, method, path, body)
	}
	return 0, nil, nil
}
