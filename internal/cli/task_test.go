package cli

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// mergeGate implements the #73 Decisions: GitHub counts approvals only
// from write-access principals — a read-only App's review and an operator
// confirmation do not count — so a REVIEW_REQUIRED decision refuses
// unless the operator asked for the recorded bypass.
func TestMergeGate(t *testing.T) {
	for _, tc := range []struct {
		decision string
		bypass   bool
		admin    bool
		refused  string
	}{
		{"", false, false, ""},                                  // no rule
		{"APPROVED", false, false, ""},                          // counted approval satisfies GitHub
		{"REVIEW_REQUIRED", false, false, "REVIEW_NOT_COUNTED"}, // the R1 configuration, no bypass asked
		{"REVIEW_REQUIRED", true, true, ""},                     // recorded bypass path
		{"CHANGES_REQUESTED", false, false, ""},                 // no approving review exists here in the common case, so the approval gate refuses earlier; with a separate non-author approval, GitHub itself refuses the merge
	} {
		admin, err := mergeGate(tc.decision, tc.bypass)
		if tc.refused == "" {
			if err != nil {
				t.Errorf("mergeGate(%q,%v) unexpected error %v", tc.decision, tc.bypass, err)
			}
			if admin != tc.admin {
				t.Errorf("mergeGate(%q,%v) admin = %v, want %v", tc.decision, tc.bypass, admin, tc.admin)
			}
			continue
		}
		var r refusal
		if !errors.As(err, &r) || r.Code != tc.refused {
			t.Errorf("mergeGate(%q,%v) = %v, want refused[%s]", tc.decision, tc.bypass, err, tc.refused)
		}
		if !strings.Contains(err.Error(), "--bypass") {
			t.Errorf("refusal does not name the supported path: %v", err)
		}
	}
}

func TestHolderReviewed(t *testing.T) {
	holds := func(login string) bool { return login == "radiusred-reviewy" || login == "radiusred-reviewy[bot]" }
	if holderReviewed([]string{"davison", "someone"}, holds) {
		t.Error("non-holder approvals satisfied the holder gate")
	}
	if !holderReviewed([]string{"davison", "radiusred-reviewy"}, holds) {
		t.Error("holder approval not recognized among others")
	}
	if holderReviewed(nil, holds) {
		t.Error("no approvals satisfied the holder gate")
	}
}

// The ownership gate holds a task to the seat that started it — the same
// login, or the same routed seat (a team-held role is any member). An
// owner who has left resolves to no role and no longer matches; and since
// 2.0 (M13-R7) no owner recorded matches nobody, because an unstarted
// task is one to run task start on, not one anyone may finish
// (#165, operator's questions on #175).
func TestSameSeat(t *testing.T) {
	roleFor := func(login string) string {
		switch strings.ToLower(strings.TrimSuffix(login, "[bot]")) {
		case "radiusred-cody":
			return "implementer"
		case "radiusred-wordy":
			return "doc-synthesizer"
		case "alice", "bob": // members of the team the reviewer seat routes to
			return "reviewer"
		}
		return ""
	}
	for _, c := range []struct {
		owner, viewer string
		want          bool
	}{
		{"", "anyone", false},                                // no start record: the gate refuses, it does not wave through
		{"radiusred-cody[bot]", "radiusred-cody[bot]", true}, // same login
		{"radiusred-cody[bot]", "Radiusred-Cody", true},      // suffix and case
		{"alice", "bob", true},                               // same team-held seat
		{"bob", "alice", true},
		{"alice", "radiusred-cody[bot]", false}, // different seats
		{"radiusred-wordy[bot]", "radiusred-cody[bot]", false},
		{"carol", "alice", false},                 // starter no longer holds any seat
		{"davison", "radiusred-cody[bot]", false}, // operator-held task, a crew finisher
		{"radiusred-cody[bot]", "davison", false}, // the operator is not exempt
		{"davison", "davison", true},
	} {
		if got := sameSeat(c.owner, c.viewer, roleFor); got != c.want {
			t.Errorf("sameSeat(%q, %q) = %v, want %v", c.owner, c.viewer, got, c.want)
		}
	}
}

// gateFake is the slice of the tracker checkpoint touches: it reads the
// target's labels through IssueLabels — the endpoint that serves PRs — and
// records the comment and label it writes. Task is deliberately not
// defined: a checkpoint that reached for the GraphQL issue query would
// panic on the embedded nil interface, which is how the pull-request case
// below proves the target stays reachable.
type gateFake struct {
	tracker.Tracker
	labels   []string
	comments []string
	added    []string
}

func (f *gateFake) IssueLabels(tracker.IssueRef) ([]string, error) { return f.labels, nil }
func (f *gateFake) Comment(_ tracker.IssueRef, body string) error {
	f.comments = append(f.comments, body)
	return nil
}
func (f *gateFake) AddLabel(_ tracker.IssueRef, label string) error {
	f.added = append(f.added, label)
	return nil
}

// checkpoint accepts a milestone ref — a question about a requirement has
// no task to carry it — and says what holds there: status lists the gate
// and milestone close refuses on it (#200, #219). A task keeps the task
// finish wording, and so does a pull request — the scaffold PR carries the
// pre-milestone gate (.codecrew/roles/coordinator.md) and must stay a valid target
// (checky's finding on PR #218).
func TestRaiseGateWordingByTarget(t *testing.T) {
	for _, tc := range []struct {
		name    string
		labels  []string
		comment string
		receipt string
		absent  string
	}{
		{"task", []string{"cc:task"}, "`task finish` refuses while the label is present", "gate raised on o/r#6 — blocked until a human removes cc:needs-decision\n", "(milestone issue)"},
		{"pull request", nil, "`task finish` refuses while the label is present", "gate raised on o/r#6 — blocked until a human removes cc:needs-decision\n", "(milestone issue)"},
		{"milestone", []string{tracker.LabelMilestone}, "`status` lists this gate beside the tasks' gates and `milestone close` refuses while the label is present", "gate raised on o/r#6 (milestone issue) — status lists it and milestone close refuses until a human removes cc:needs-decision\n", "`task finish` refuses"},
	} {
		f := &gateFake{labels: tc.labels}
		cfg := &config.Config{Codecrew: "1.0", Hub: "self"}
		c := &ctx{cfg: cfg, roles: cfg, current: "o/r", hub: "o/r", t: f}
		var out bytes.Buffer
		if err := raiseGate(&out, c, tracker.IssueRef{Repo: "o/r", Number: 6}, "which way?"); err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if out.String() != tc.receipt {
			t.Errorf("%s: receipt = %q, want %q", tc.name, out.String(), tc.receipt)
		}
		if len(f.comments) != 1 || !strings.HasPrefix(f.comments[0], "**Gate raised:** which way?\n\n") {
			t.Fatalf("%s: comments = %q", tc.name, f.comments)
		}
		if !strings.Contains(f.comments[0], tc.comment) || strings.Contains(f.comments[0], tc.absent) {
			t.Errorf("%s: comment must say %q and not %q:\n%s", tc.name, tc.comment, tc.absent, f.comments[0])
		}
		if len(f.added) != 1 || f.added[0] != tracker.LabelNeedsDecision {
			t.Errorf("%s: labels added = %q", tc.name, f.added)
		}
	}
}

// taskNewFake answers OpenMilestones and RecentIssues from successive
// listings, the last one repeating — a listing that catches up between
// reads, or never does.
type taskNewFake struct {
	tracker.Tracker
	open    [][]tracker.Milestone   // successive OpenMilestones answers
	recent  [][]tracker.TitledIssue // successive RecentIssues answers
	issues  map[int]tracker.Task    // what Task answers for a hub issue number
	oCalls  int
	rCalls  int
	created []string
	body    string   // the body of the issue CreateIssue was last given
	linked  []string // "<parent> <- <child>" per AddSubIssue
	posted  []string // "<ref>: <body>" per Comment
}

func (f *taskNewFake) OpenMilestones(string) ([]tracker.Milestone, error) {
	f.oCalls++
	if len(f.open) == 0 {
		return nil, nil
	}
	k := f.oCalls - 1
	if k >= len(f.open) {
		k = len(f.open) - 1
	}
	return f.open[k], nil
}
func (f *taskNewFake) RecentIssues(string) ([]tracker.TitledIssue, error) {
	f.rCalls++
	return nth(f.recent, f.rCalls-1), nil
}
func (f *taskNewFake) Task(ref tracker.IssueRef) (tracker.Task, error) {
	t, ok := f.issues[ref.Number]
	if !ok {
		return tracker.Task{}, errors.New("issue not found")
	}
	return t, nil
}
func (f *taskNewFake) CreateIssue(repo, title, body string, _ []string) (tracker.IssueRef, error) {
	f.created = append(f.created, title)
	f.body = body
	return tracker.IssueRef{Repo: repo, Number: 21}, nil
}
func (f *taskNewFake) Comment(ref tracker.IssueRef, body string) error {
	f.posted = append(f.posted, ref.String()+": "+body)
	return nil
}
func (f *taskNewFake) AddSubIssue(parent, child tracker.IssueRef) error {
	f.linked = append(f.linked, parent.String()+" <- "+child.String())
	return nil
}

func openMilestone(n int, title string) tracker.Milestone {
	return tracker.Milestone{Ref: tracker.IssueRef{Repo: "o/hub", Number: n}, Title: title}
}

// recordSleeps swaps the sleeper for one that records the waits it was
// asked for, so the tests take none of them.
func recordSleeps(t *testing.T) *[]time.Duration {
	t.Helper()
	var waits []time.Duration
	prev := sleep
	sleep = func(d time.Duration) { waits = append(waits, d) }
	t.Cleanup(func() { sleep = prev })
	return &waits
}

func taskNewCtx(f *taskNewFake) *ctx {
	cfg := &config.Config{Codecrew: "1.0", Hub: "self"}
	return &ctx{cfg: cfg, roles: cfg, current: "o/hub", hub: "o/hub", t: f}
}

// The listing has the milestone: one read, no wait, no note.
func TestTaskNewFindsTheMilestoneFirstTime(t *testing.T) {
	waits := recordSleeps(t)
	f := &taskNewFake{open: [][]tracker.Milestone{{openMilestone(233, "M11: Housekeeping")}}}
	var out bytes.Buffer
	if err := runTaskNew(taskNewCtx(f), &out, 11, "o/spoke", "Cut the README", "g", "M11-R1", nil); err != nil {
		t.Fatal(err)
	}
	if f.oCalls != 1 || f.rCalls != 0 || len(*waits) != 0 {
		t.Errorf("a fresh listing was read %d times, recent %d, waits %v", f.oCalls, f.rCalls, *waits)
	}
	if got := out.String(); got != "created task o/spoke#21 as a sub-issue of o/hub#233\n" {
		t.Errorf("output:\n%s", got)
	}
	if len(f.linked) != 1 || f.linked[0] != "o/hub#233 <- o/spoke#21" {
		t.Errorf("linked %v", f.linked)
	}
}

// The #234 failure: milestone new created M11 seconds ago and neither
// listing has it on the first read. The verb waits and reads again instead
// of refusing, and says the listing lagged.
func TestTaskNewWaitsForAListingThatCatchesUp(t *testing.T) {
	waits := recordSleeps(t)
	f := &taskNewFake{
		open:   [][]tracker.Milestone{{openMilestone(200, "M10: Field fixes")}, {openMilestone(200, "M10: Field fixes")}, {openMilestone(200, "M10: Field fixes"), openMilestone(233, "M11: Housekeeping")}},
		recent: listing(issue(9, "a plain issue"), issue(200, "M10: Field fixes")),
	}
	var out bytes.Buffer
	if err := runTaskNew(taskNewCtx(f), &out, 11, "o/spoke", "Cut the README", "g", "M11-R1", nil); err != nil {
		t.Fatal(err)
	}
	if f.oCalls != 3 || f.rCalls != 2 {
		t.Errorf("read the open listing %d times and recent %d, want 3 and 2 (the third open read found it)", f.oCalls, f.rCalls)
	}
	if want := []time.Duration{milestoneLookupWait, milestoneLookupWait}; !reflect.DeepEqual(*waits, want) {
		t.Errorf("waits %v, want %v", *waits, want)
	}
	got := out.String()
	if !strings.Contains(got, "milestone M11 (o/hub#233) appeared in the listing on read 3") || !strings.Contains(got, "created task o/spoke#21 as a sub-issue of o/hub#233") {
		t.Errorf("output:\n%s", got)
	}
}

// The unfiltered newest issues catch the milestone the label-filtered
// listing has not indexed yet, so no wait is needed — but only an open
// issue carrying cc:milestone counts: a closed one, or a plain issue whose
// title happens to start "M11:", is not the milestone.
func TestTaskNewFallsBackToTheRecentIssues(t *testing.T) {
	waits := recordSleeps(t)
	f := &taskNewFake{
		open:   [][]tracker.Milestone{{openMilestone(200, "M10: Field fixes")}},
		recent: listing(issue(235, "M11: a plain issue with the prefix"), issue(234, "M11: Housekeeping, closed by mistake"), issue(233, "M11: Housekeeping")),
		issues: map[int]tracker.Task{
			235: {Labels: []string{"cc:task"}},
			234: {Labels: []string{"cc:milestone"}, Closed: true},
			233: {Labels: []string{"cc:milestone"}},
		},
	}
	var out bytes.Buffer
	if err := runTaskNew(taskNewCtx(f), &out, 11, "o/spoke", "Cut the README", "g", "M11-R1", nil); err != nil {
		t.Fatal(err)
	}
	if f.oCalls != 1 || f.rCalls != 1 || len(*waits) != 0 {
		t.Errorf("read the open listing %d times and recent %d with waits %v; want one read each and no wait", f.oCalls, f.rCalls, *waits)
	}
	got := out.String()
	if !strings.Contains(got, "milestone M11 (o/hub#233) is not in the open-milestone listing yet, found among the hub's newest issues") || !strings.Contains(got, "created task o/spoke#21 as a sub-issue of o/hub#233") {
		t.Errorf("output:\n%s", got)
	}
}

// Never found: NOT_FOUND after the bounded reads, with exactly the waits
// between them, nothing created, and a detail that says the listing may lag.
func TestTaskNewRefusesNotFoundAfterBoundedRetries(t *testing.T) {
	waits := recordSleeps(t)
	f := &taskNewFake{
		open:   [][]tracker.Milestone{{openMilestone(200, "M10: Field fixes")}},
		recent: listing(issue(9, "a plain issue"), issue(200, "M10: Field fixes")),
	}
	var out bytes.Buffer
	err := runTaskNew(taskNewCtx(f), &out, 11, "o/spoke", "Cut the README", "g", "M11-R1", nil)
	var r refusal
	if !errors.As(err, &r) || r.Code != "NOT_FOUND" {
		t.Fatalf("err = %v, want refused[NOT_FOUND]", err)
	}
	if !strings.Contains(r.Detail, "no open milestone M11 in o/hub after 3 reads") || !strings.Contains(r.Detail, "lag a milestone created seconds ago") {
		t.Errorf("detail: %s", r.Detail)
	}
	if f.oCalls != milestoneLookupAttempts || f.rCalls != milestoneLookupAttempts {
		t.Errorf("read the open listing %d times and recent %d, want %d each", f.oCalls, f.rCalls, milestoneLookupAttempts)
	}
	if want := []time.Duration{milestoneLookupWait, milestoneLookupWait}; !reflect.DeepEqual(*waits, want) {
		t.Errorf("waits %v, want %v", *waits, want)
	}
	if len(f.created) != 0 || len(f.linked) != 0 || out.Len() != 0 {
		t.Errorf("a refusal created %v, linked %v, printed %q", f.created, f.linked, out.String())
	}
}

// --adopts is repeatable and comma-separated, resolves a bare number
// against the task's own repo and a full ref anywhere, collapses a
// duplicate, and leaves the capture's title in the body for the reader.
// The section it writes is the section task finish reads back.
func TestTaskNewAdoptsWritesTheSectionAndRecordsIt(t *testing.T) {
	f := &taskNewFake{
		open: [][]tracker.Milestone{{openMilestone(233, "M11: Housekeeping")}},
		issues: map[int]tracker.Task{
			193: {Title: "carry an adopted backlog issue"},
			194: {Title: "the tidy verb"},
			7:   {Title: "a capture in another repo"},
		},
	}
	var out bytes.Buffer
	if err := runTaskNew(taskNewCtx(f), &out, 11, "o/spoke", "Cut the README", "g", "M11-R1", []string{"193, #194", "o/other#7", "193"}); err != nil {
		t.Fatal(err)
	}
	want := "## Requirements\nM11-R1\n\n## Adopts\n- #193 — carry an adopted backlog issue\n- #194 — the tidy verb\n- o/other#7 — a capture in another repo\n\n## Plan\n"
	if !strings.Contains(f.body, want) {
		t.Errorf("body:\n%s\nwant the section:\n%s", f.body, want)
	}
	task := tracker.IssueRef{Repo: "o/spoke", Number: 21}
	refs := tracker.AdoptedRefs(f.body, "o/spoke")
	if want := []tracker.IssueRef{{Repo: "o/spoke", Number: 193}, {Repo: "o/spoke", Number: 194}, {Repo: "o/other", Number: 7}}; !reflect.DeepEqual(refs, want) {
		t.Errorf("the section reads back as %v, want %v", refs, want)
	}
	var wantPosts []string
	for _, ref := range refs {
		wantPosts = append(wantPosts, ref.String()+": "+tracker.AdoptionRecord(task))
	}
	if !reflect.DeepEqual(f.posted, wantPosts) {
		t.Errorf("comments posted:\n%v\nwant:\n%v", f.posted, wantPosts)
	}
	for _, line := range []string{"created task o/spoke#21 as a sub-issue of o/hub#233\n", "adopts o/spoke#193 — recorded on the capture\n", "adopts o/other#7 — recorded on the capture\n"} {
		if !strings.Contains(out.String(), line) {
			t.Errorf("output lacks %q:\n%s", line, out.String())
		}
	}
}

// A ref that is not an open issue refuses ADOPT_NOT_OPEN — unreadable and
// closed are the two ways, one code — and refuses before the task is
// created, so nothing is half-adopted.
func TestTaskNewRefusesAnAdoptionThatIsNotOpen(t *testing.T) {
	for _, tc := range []struct {
		name, ref, detail string
	}{
		{"unknown", "404", "could not be read"},
		{"closed", "193", "is closed"},
	} {
		f := &taskNewFake{
			open:   [][]tracker.Milestone{{openMilestone(233, "M11: Housekeeping")}},
			issues: map[int]tracker.Task{193: {Title: "already delivered", Closed: true}},
		}
		var out bytes.Buffer
		err := runTaskNew(taskNewCtx(f), &out, 11, "o/spoke", "Cut the README", "g", "M11-R1", []string{tc.ref})
		var r refusal
		if !errors.As(err, &r) || r.Code != "ADOPT_NOT_OPEN" {
			t.Fatalf("%s: err = %v, want refused[ADOPT_NOT_OPEN]", tc.name, err)
		}
		if !strings.Contains(r.Detail, "o/spoke#"+tc.ref) || !strings.Contains(r.Detail, tc.detail) {
			t.Errorf("%s: detail = %s", tc.name, r.Detail)
		}
		if len(f.created) != 0 || len(f.linked) != 0 || len(f.posted) != 0 || out.Len() != 0 {
			t.Errorf("%s: a refusal created %v, linked %v, posted %v, printed %q", tc.name, f.created, f.linked, f.posted, out.String())
		}
	}
}
