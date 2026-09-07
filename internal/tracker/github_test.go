package tracker

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/gh"
)

// fakeGH stands a gh behind gh.Command for one test: `pr view` answers
// with a minimal open PR, `pr checks` fails with checksErr on stderr the
// way gh reports a GraphQL refusal, and the task-branch ref listing answers
// the way GitHub does. The fake is this test binary re-entered at
// TestHelperGH — no script, no PATH shim.
func fakeGH(t *testing.T, checksErr string, env ...string) {
	t.Helper()
	orig := gh.Command
	gh.Command = func(_ string, args ...string) *exec.Cmd {
		cmd := exec.Command(os.Args[0], append([]string{"-test.run=^TestHelperGH$", "--"}, args...)...)
		cmd.Env = append(os.Environ(), append([]string{"GH_HELPER=1", "GH_HELPER_CHECKS_STDERR=" + checksErr}, env...)...)
		return cmd
	}
	t.Cleanup(func() { gh.Command = orig })
}

// TestHelperGH is the fake gh's body; it is inert unless fakeGH re-enters
// the binary with GH_HELPER set.
func TestHelperGH(t *testing.T) {
	if os.Getenv("GH_HELPER") != "1" {
		return
	}
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	if len(args) >= 2 && args[0] == "pr" {
		switch args[1] {
		case "view":
			fmt.Print(`{"state":"OPEN","reviewDecision":"","headRefName":"task/1-x","headRefOid":"abc","isCrossRepository":false,"mergedAt":"","author":{"login":"app/myorg-coder"},"reviews":[]}`)
			os.Exit(0)
		case "checks":
			fmt.Fprint(os.Stderr, os.Getenv("GH_HELPER_CHECKS_STDERR"))
			os.Exit(1)
		}
	}
	// The ref listing, answered as GitHub answers it: each node's name has
	// the queried prefix removed. GH_HELPER_HAS_NEXT makes the page a
	// partial one.
	if len(args) >= 2 && args[0] == "api" && args[1] == "graphql" && slices.Contains(args, "prefix=refs/heads/task/") {
		hasNext := "false"
		if os.Getenv("GH_HELPER_HAS_NEXT") == "1" {
			hasNext = "true"
		}
		fmt.Printf(`{"data":{"repository":{"refs":{"nodes":[{"name":"8-cycle-1-record"},{"name":"14-cycle-2-record"},{"name":""}],"pageInfo":{"hasNextPage":%s}}}}}`, hasNext)
		os.Exit(0)
	}
	// The open-PR-by-head listing: the fake echoes the head filter it was
	// given, so the test can assert the `<owner>:<ref>` grammar reached gh.
	// The path is matched wherever it sits in the argument list, because
	// the call carries --paginate ahead of it.
	if len(args) >= 2 && args[0] == "api" && slices.ContainsFunc(args, func(a string) bool { return strings.HasSuffix(a, "/pulls") }) {
		if slices.Contains(args, "head=radiusred:task/8-cycle-1-record") {
			fmt.Print(`[{"number":41},{"number":42}]`)
			os.Exit(0)
		}
		fmt.Print(`[]`)
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, "fake gh: unexpected call %v", args)
	os.Exit(2)
}

// PRInfo carries the permission an installation token lacks as
// PR.ChecksUnreadable instead of the raw GraphQL error — the two verbatim
// shapes #198 recorded — and returns every other failure raw, the same
// message on an unrelated path included. The dry-run test injects the
// field through a fake tracker; this is the wiring that sets it.
func TestPRInfoMapsMissingChecksPermission(t *testing.T) {
	cases := []struct {
		stderr, want string
	}{
		{"GraphQL: Resource not accessible by integration (node.statusCheckRollup.nodes.0.commit.statusCheckRollup)", "checks: read"},
		{"GraphQL: Resource not accessible by integration (node.statusCheckRollup.nodes.0.commit.statusCheckRollup.contexts.nodes.0.checkSuite.workflowRun)", "actions: read"},
	}
	for _, tc := range cases {
		fakeGH(t, tc.stderr)
		pr, err := GitHub{}.PRInfo("o/r", 9)
		if err != nil {
			t.Fatalf("%q: PRInfo returned the raw error: %v", tc.want, err)
		}
		if pr.ChecksUnreadable != tc.want {
			t.Errorf("ChecksUnreadable = %q, want %q", pr.ChecksUnreadable, tc.want)
		}
		if pr.NoChecks || pr.ChecksOK || pr.ChecksPending {
			t.Errorf("%q: check state must stay unknown, got %+v", tc.want, pr)
		}
		if pr.Number != 9 || pr.Author != "myorg-coder" || !pr.Open {
			t.Errorf("%q: the PR's own state must survive the mapping: %+v", tc.want, pr)
		}
	}

	raw := "GraphQL: Resource not accessible by integration (repository.collaborators)"
	fakeGH(t, raw)
	pr, err := GitHub{}.PRInfo("o/r", 9)
	if err == nil || !strings.Contains(err.Error(), raw) {
		t.Fatalf("an unrelated refusal must surface raw, got err %v pr %+v", err, pr)
	}
	if pr.ChecksUnreadable != "" {
		t.Errorf("unrelated refusal classified as %q", pr.ChecksUnreadable)
	}

	fakeGH(t, "no checks reported on the 'task/1-x' branch")
	pr, err = GitHub{}.PRInfo("o/r", 9)
	if err != nil || !pr.NoChecks || pr.ChecksUnreadable != "" {
		t.Errorf("the checkless shape still maps to NoChecks: err %v pr %+v", err, pr)
	}
}

// The stale sweep's listing is filtered at the server by ref prefix, and
// GitHub returns each ref's name with that prefix removed — so TaskBranches
// puts it back, and what the sweep reads is a branch name it can delete
// rather than a bare slug. The shape was read off radiusred/numberguess,
// the repo #167 was captured from. A nameless node is not a branch.
func TestTaskBranchesRestoresThePrefix(t *testing.T) {
	fakeGH(t, "")
	got, truncated, err := GitHub{}.TaskBranches("radiusred/numberguess")
	if err != nil {
		t.Fatalf("TaskBranches: %v", err)
	}
	want := []string{"task/8-cycle-1-record", "task/14-cycle-2-record"}
	if !slices.Equal(got, want) {
		t.Errorf("TaskBranches = %v, want %v", got, want)
	}
	if truncated {
		t.Error("a listing GitHub says has no next page must not read as truncated")
	}
	// A repo with more task branches than the page holds says so, and still
	// hands back the page it got.
	fakeGH(t, "", "GH_HELPER_HAS_NEXT=1")
	got, truncated, err = GitHub{}.TaskBranches("radiusred/numberguess")
	if err != nil || !truncated || !slices.Equal(got, want) {
		t.Errorf("truncated listing: %v %v %v", got, truncated, err)
	}
	_, _, err = GitHub{}.TaskBranches("numberguess")
	if err == nil {
		t.Error("a repo ref with no owner must refuse before the API call")
	}
}

// OpenPRsForBranch is the relation ClosingPRs cannot see — a PR open on a
// branch that references no issue — and the sweep will not delete a branch
// under one. The filter's grammar is `<owner>:<ref>`, and a branch name's
// slashes go through gh's own query building untouched.
func TestOpenPRsForBranch(t *testing.T) {
	fakeGH(t, "")
	got, err := GitHub{}.OpenPRsForBranch("radiusred/numberguess", "task/8-cycle-1-record")
	if err != nil || !slices.Equal(got, []int{41, 42}) {
		t.Errorf("OpenPRsForBranch = %v, %v", got, err)
	}
	got, err = GitHub{}.OpenPRsForBranch("radiusred/numberguess", "task/14-cycle-2-record")
	if err != nil || len(got) != 0 {
		t.Errorf("a branch with no open PR = %v, %v", got, err)
	}
	_, err = GitHub{}.OpenPRsForBranch("numberguess", "task/8-x")
	if err == nil {
		t.Error("a repo ref with no owner must refuse before the API call")
	}
}

// recordGH stands a gh behind gh.Command that records the arguments and
// answers with stdout, so a call's shape can be asserted without a
// GitHub. The command is `true`, which succeeds and prints nothing, with
// the answer piped in through the helper binary only where one is needed.
func recordGH(t *testing.T, stdout string) *[][]string {
	t.Helper()
	orig := gh.Command
	var calls [][]string
	gh.Command = func(_ string, args ...string) *exec.Cmd {
		calls = append(calls, args)
		return exec.Command("printf", "%s", stdout)
	}
	t.Cleanup(func() { gh.Command = orig })
	return &calls
}

// The two label calls are pass-throughs, so what is asserted is their
// shape: the listing is paginated, because a repository with more than a
// page of labels must not read as missing the protocol's, and the colour
// goes to the API as six hex digits with no leading "#", which the API
// rejects.
func TestLabelCalls(t *testing.T) {
	calls := recordGH(t, `[{"name":"bug","color":"d73a4a","description":"Something is broken"},{"name":"cc:task","color":"ededed","description":""}]`)
	got, err := GitHub{}.Labels("o/r")
	if err != nil {
		t.Fatal(err)
	}
	want := []Label{{"bug", "d73a4a", "Something is broken"}, {"cc:task", "ededed", ""}}
	if !slices.Equal(got, want) {
		t.Errorf("Labels = %+v, want %+v", got, want)
	}
	if len(*calls) != 1 {
		t.Fatalf("calls = %v", *calls)
	}
	line := strings.Join((*calls)[0], " ")
	for _, want := range []string{"--paginate", "repos/o/r/labels?per_page=100"} {
		if !strings.Contains(line, want) {
			t.Errorf("the listing call %q is missing %q", line, want)
		}
	}

	calls = recordGH(t, "")
	l, _ := ProtocolLabel(LabelNeedsDecision)
	if err := (GitHub{}).CreateLabel("o/r", l); err != nil {
		t.Fatal(err)
	}
	line = strings.Join((*calls)[0], " ")
	for _, want := range []string{"-X POST", "repos/o/r/labels", "name=" + l.Name, "color=" + l.Color, "description=" + l.Description} {
		if !strings.Contains(line, want) {
			t.Errorf("the create call %q is missing %q", line, want)
		}
	}
	if strings.Contains(line, "color=#") {
		t.Errorf("the colour reached the API with a leading #: %q", line)
	}

	// A restyle is the shape the call site builds: the name as the
	// repository spells it, carrying the protocol's colour and
	// description. It goes in the path segment and no new_name is sent,
	// so the label keeps its spelling — restyling is not renaming.
	calls = recordGH(t, "")
	repoSpelling := Label{Name: "CC:Needs-Decision", Color: l.Color, Description: l.Description}
	if err := (GitHub{}).UpdateLabel("o/r", repoSpelling); err != nil {
		t.Fatal(err)
	}
	line = strings.Join((*calls)[0], " ")
	for _, want := range []string{"-X PATCH", "repos/o/r/labels/CC:Needs-Decision", "color=" + l.Color, "description=" + l.Description} {
		if !strings.Contains(line, want) {
			t.Errorf("the restyle call %q is missing %q", line, want)
		}
	}
	if strings.Contains(line, "new_name") {
		t.Errorf("the restyle renamed the label: %q", line)
	}
}

// listingPage is the hundred GitHub's own first page holds, and what the
// paging fake truncates an unpaginated read to.
const listingPage = 100

// pagingGH stands a gh behind gh.Command that pages the way the real one
// does: with --paginate the whole listing comes back as one JSON array —
// gh joins a REST array endpoint's pages — and without it the read stops
// at the first hundred, exactly as GitHub's first page does. A reader that
// drops the flag therefore fails here rather than passing on a fixture
// that fits in one page (#264). items are the listing's JSON objects, in
// the order GitHub returns them; the recorded calls let a test assert the
// flag reached gh.
func pagingGH(t *testing.T, items []string) *[][]string {
	t.Helper()
	orig := gh.Command
	var calls [][]string
	gh.Command = func(_ string, args ...string) *exec.Cmd {
		calls = append(calls, args)
		page := items
		if !slices.Contains(args, "--paginate") && len(page) > listingPage {
			page = page[:listingPage]
		}
		return exec.Command("printf", "%s", "["+strings.Join(page, ",")+"]")
	}
	t.Cleanup(func() { gh.Command = orig })
	return &calls
}

// Comments walks the whole listing, in GitHub's oldest-first order. The
// capture is exact: latest-wins reads the newest verdict, so a milestone
// issue whose satisfied verdict landed past comment one hundred had
// `milestone close` refusing VERDICT_MISSING for a requirement QA had
// signed off (#264). The fixture is the 101-comment issue that capture
// asked for, carrying a superseded verdict on page one and the live one
// last.
func TestCommentsPaginatesInOrder(t *testing.T) {
	var raw []string
	for i := 1; i <= 101; i++ {
		body := fmt.Sprintf("comment %d", i)
		switch i {
		case 7:
			body = "**M15-R3 — not satisfied** the first pass"
		case 101:
			body = "**M15-R3 — satisfied** the re-run"
		}
		raw = append(raw, fmt.Sprintf(`{"body":%q,"html_url":"https://x/%d","user":{"login":"myorg-qa"}}`, body, i))
	}
	calls := pagingGH(t, raw)
	got, err := GitHub{}.Comments(IssueRef{Repo: "o/r", Number: 5})
	if err != nil {
		t.Fatalf("Comments: %v", err)
	}
	if len(got) != 101 {
		t.Fatalf("Comments returned %d of 101 comments — the listing stopped at a page boundary", len(got))
	}
	if got[0].Body != "comment 1" || got[len(got)-1].Body != "**M15-R3 — satisfied** the re-run" {
		t.Errorf("oldest-first order not preserved across pages: first %q, last %q", got[0].Body, got[len(got)-1].Body)
	}
	for i, c := range got {
		if c.Author != "myorg-qa" || c.URL != fmt.Sprintf("https://x/%d", i+1) {
			t.Fatalf("comment %d arrived out of order or half-mapped: %+v", i, c)
		}
	}
	// What the gate does with them: the newest verdict wins, which is only
	// true if the newest comment is in hand.
	latest := map[string]string{}
	for _, v := range ParseVerdicts(got) {
		latest[v.ID] = v.State
	}
	if latest["M15-R3"] != "satisfied" {
		t.Errorf("latest-wins verdict = %q, want satisfied — the newest verdict was dropped", latest["M15-R3"])
	}
	if len(*calls) != 1 {
		t.Fatalf("calls = %v", *calls)
	}
	line := strings.Join((*calls)[0], " ")
	for _, want := range []string{"--paginate", "repos/o/r/issues/5/comments?per_page=100"} {
		if !strings.Contains(line, want) {
			t.Errorf("the comments call %q is missing %q", line, want)
		}
	}
}

// SubIssues walks the whole listing too: a milestone past a page of tasks
// must not read as having only its first hundred, or every gate counting
// tasks counts the wrong set.
func TestSubIssuesPaginates(t *testing.T) {
	var raw []string
	for i := 1; i <= 150; i++ {
		repo := "https://api.github.com/repos/o/r"
		if i == 150 {
			repo = "https://api.github.com/repos/o/spoke"
		}
		raw = append(raw, fmt.Sprintf(`{"number":%d,"repository_url":%q}`, i, repo))
	}
	calls := pagingGH(t, raw)
	got, err := GitHub{}.SubIssues(IssueRef{Repo: "o/r", Number: 5})
	if err != nil {
		t.Fatalf("SubIssues: %v", err)
	}
	if len(got) != 150 {
		t.Fatalf("SubIssues returned %d of 150 sub-issues", len(got))
	}
	if got[0] != (IssueRef{Repo: "o/r", Number: 1}) {
		t.Errorf("first sub-issue = %+v", got[0])
	}
	// The last one is on the second page, and in a spoke: both survive.
	if got[149] != (IssueRef{Repo: "o/spoke", Number: 150}) {
		t.Errorf("last sub-issue = %+v, want o/spoke#150", got[149])
	}
	if line := strings.Join((*calls)[0], " "); !strings.Contains(line, "--paginate") {
		t.Errorf("the sub-issue call %q is not paginated", line)
	}
}

// The milestone listing walks every page — the hundred-and-first milestone
// is not absent — and pull requests stay dropped wherever they fall. Its
// unpaginated sibling is RecentIssues, which wants the newest page and
// only that: the floor fallback of #195 would otherwise walk a repo's
// whole issue history for a number the first page already carries.
func TestMilestoneIssuesPaginateAndRecentIssuesDoesNot(t *testing.T) {
	var raw []string
	for i := 1; i <= 150; i++ {
		if i == 120 {
			raw = append(raw, fmt.Sprintf(`{"number":%d,"title":"a PR","pull_request":{}}`, i))
			continue
		}
		raw = append(raw, fmt.Sprintf(`{"number":%d,"title":"M%d"}`, i, i))
	}
	calls := pagingGH(t, raw)
	got, err := GitHub{}.MilestoneIssues("o/r")
	if err != nil {
		t.Fatalf("MilestoneIssues: %v", err)
	}
	if len(got) != 149 {
		t.Fatalf("MilestoneIssues returned %d of the 149 issues in a 150-row listing", len(got))
	}
	if got[len(got)-1].Ref.Number != 150 || got[len(got)-1].Title != "M150" {
		t.Errorf("last milestone = %+v, want #150", got[len(got)-1])
	}
	for _, g := range got {
		if g.Ref.Number == 120 {
			t.Error("a pull request on the second page reached the milestone listing")
		}
	}
	line := strings.Join((*calls)[0], " ")
	for _, want := range []string{"--paginate", "repos/o/r/issues?labels=cc:milestone&state=all&per_page=100"} {
		if !strings.Contains(line, want) {
			t.Errorf("the milestone listing call %q is missing %q", line, want)
		}
	}

	calls = pagingGH(t, raw)
	recent, err := GitHub{}.RecentIssues("o/r")
	if err != nil {
		t.Fatalf("RecentIssues: %v", err)
	}
	if len(recent) != listingPage {
		t.Errorf("RecentIssues read %d rows: it is the deliberate exception and reads the newest page only", len(recent))
	}
	if line := strings.Join((*calls)[0], " "); strings.Contains(line, "--paginate") {
		t.Errorf("RecentIssues paginated: %q", line)
	}
}

// The sweep keeps a branch that carries an open PR, so the PR listing it
// asks must be the whole one: a PR on page two is still a PR.
func TestOpenPRsForBranchPaginates(t *testing.T) {
	var raw []string
	for i := 1; i <= 150; i++ {
		raw = append(raw, fmt.Sprintf(`{"number":%d}`, i))
	}
	calls := pagingGH(t, raw)
	got, err := GitHub{}.OpenPRsForBranch("o/r", "task/8-x")
	if err != nil {
		t.Fatalf("OpenPRsForBranch: %v", err)
	}
	if len(got) != 150 || got[0] != 1 || got[149] != 150 {
		t.Fatalf("OpenPRsForBranch returned %d PRs (first %d, last %d)", len(got), got[0], got[len(got)-1])
	}
	line := strings.Join((*calls)[0], " ")
	for _, want := range []string{"--paginate", "head=o:task/8-x"} {
		if !strings.Contains(line, want) {
			t.Errorf("the pulls call %q is missing %q", line, want)
		}
	}
}
