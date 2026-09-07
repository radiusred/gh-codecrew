package tracker

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/gh"
)

// GitHub implements Tracker over the gh CLI.
type GitHub struct{}

func (g GitHub) OpenMilestones(hub string) ([]Milestone, error) {
	milestones, err := listMilestones(hub, "open")
	if err != nil {
		return nil, err
	}
	for i := range milestones {
		tasks, err := g.SubIssues(milestones[i].Ref)
		if err != nil {
			return nil, err
		}
		milestones[i].Tasks = tasks
	}
	return milestones, nil
}

func listMilestones(hub, state string) ([]Milestone, error) {
	issues, err := listIssues(hub, fmt.Sprintf("labels=cc:milestone&state=%s", state))
	if err != nil {
		return nil, err
	}
	milestones := make([]Milestone, 0, len(issues))
	for _, is := range issues {
		milestones = append(milestones, Milestone{Ref: is.Ref, Title: is.Title})
	}
	return milestones, nil
}

// listIssues walks the whole of repo's issue listing under query, pull
// requests dropped: the REST listing returns both, marked by a
// pull_request object. A read that stopped at the first hundred would
// report the hundred-and-first milestone as absent (#264).
func listIssues(repo, query string) ([]TitledIssue, error) {
	return issueListing(repo, query, true)
}

// issueListing reads repo's issue listing under query. With paginate the
// read walks every page — gh joins a REST array endpoint's pages into one
// JSON array, so the whole listing unmarshals as it stands — and without
// it the read stops at the first page, which only RecentIssues wants.
func issueListing(repo, query string, paginate bool) ([]TitledIssue, error) {
	var items []struct {
		Number      int       `json:"number"`
		Title       string    `json:"title"`
		PullRequest *struct{} `json:"pull_request"`
	}
	args := []string{"api"}
	if paginate {
		args = append(args, "--paginate")
	}
	args = append(args, fmt.Sprintf("repos/%s/issues?%s&per_page=100", repo, query))
	if err := gh.JSON(&items, args...); err != nil {
		return nil, err
	}
	issues := make([]TitledIssue, 0, len(items))
	for _, it := range items {
		if it.PullRequest != nil {
			continue
		}
		issues = append(issues, TitledIssue{Ref: IssueRef{Repo: repo, Number: it.Number}, Title: it.Title})
	}
	return issues, nil
}

func (GitHub) AddSubIssue(parent, child IssueRef) error {
	var issue struct {
		ID int64 `json:"id"`
	}
	if err := gh.JSON(&issue, "api", fmt.Sprintf("repos/%s/issues/%d", child.Repo, child.Number)); err != nil {
		return err
	}
	_, err := gh.Run("api", "-X", "POST",
		fmt.Sprintf("repos/%s/issues/%d/sub_issues", parent.Repo, parent.Number),
		"-F", fmt.Sprintf("sub_issue_id=%d", issue.ID))
	return err
}

// SubIssues lists every sub-issue of parent, paginated: a milestone past
// a page of tasks must not read as having only the first hundred (#264).
func (GitHub) SubIssues(parent IssueRef) ([]IssueRef, error) {
	var subs []struct {
		Number        int    `json:"number"`
		RepositoryURL string `json:"repository_url"`
	}
	path := fmt.Sprintf("repos/%s/issues/%d/sub_issues?per_page=100", parent.Repo, parent.Number)
	if err := gh.JSON(&subs, "api", "--paginate", path); err != nil {
		return nil, err
	}
	refs := make([]IssueRef, 0, len(subs))
	for _, s := range subs {
		repo := parent.Repo
		if i := strings.Index(s.RepositoryURL, "/repos/"); i >= 0 {
			repo = s.RepositoryURL[i+len("/repos/"):]
		}
		refs = append(refs, IssueRef{Repo: repo, Number: s.Number})
	}
	return refs, nil
}

func (GitHub) MilestoneIssues(hub string) ([]TitledIssue, error) {
	return listIssues(hub, "labels=cc:milestone&state=all")
}

// RecentIssues reads the newest page of the unfiltered listing. The
// label-filtered listing lagged an issue created seconds earlier and three
// milestones came back as M2 (#195); the unfiltered listing is a second
// source for the floor, not a guarantee — the post-create check is that.
// It is deliberately unpaginated — one of the two reads here that stay on
// a single page, the other being TaskBranches: the newest page is the
// whole point, and walking a repo's entire issue history to raise a floor
// would cost a request per hundred issues for a number the first page
// already carries.
func (GitHub) RecentIssues(repo string) ([]TitledIssue, error) {
	return issueListing(repo, "state=all&sort=created&direction=desc", false)
}

// IssueBody returns the body with its line endings normalised: this is the
// boundary a body crosses into the package, and every scan below it reads
// line by line (NormalizeLineEndings).
func (GitHub) IssueBody(ref IssueRef) (string, error) {
	var issue struct {
		Body string `json:"body"`
	}
	err := gh.JSON(&issue, "api", fmt.Sprintf("repos/%s/issues/%d", ref.Repo, ref.Number))
	return NormalizeLineEndings(issue.Body), err
}

func (GitHub) IssueLabels(ref IssueRef) ([]string, error) {
	var issue struct {
		Labels []struct {
			Name string `json:"name"`
		} `json:"labels"`
	}
	if err := gh.JSON(&issue, "api", fmt.Sprintf("repos/%s/issues/%d", ref.Repo, ref.Number)); err != nil {
		return nil, err
	}
	labels := make([]string, 0, len(issue.Labels))
	for _, l := range issue.Labels {
		labels = append(labels, l.Name)
	}
	return labels, nil
}

func (GitHub) CreateIssue(repo, title, body string, labels []string) (IssueRef, error) {
	args := []string{"api", "-X", "POST", fmt.Sprintf("repos/%s/issues", repo),
		"-f", "title=" + title, "-f", "body=" + body}
	for _, l := range labels {
		args = append(args, "-f", "labels[]="+l)
	}
	var created struct {
		Number int `json:"number"`
	}
	if err := gh.JSON(&created, args...); err != nil {
		return IssueRef{}, err
	}
	return IssueRef{Repo: repo, Number: created.Number}, nil
}

func (GitHub) EditIssue(ref IssueRef, title, body string) error {
	_, err := gh.Run("api", "-X", "PATCH",
		fmt.Sprintf("repos/%s/issues/%d", ref.Repo, ref.Number),
		"-f", "title="+title, "-f", "body="+body)
	return err
}

func (GitHub) Comment(ref IssueRef, body string) error {
	_, err := gh.Run("api", "-X", "POST",
		fmt.Sprintf("repos/%s/issues/%d/comments", ref.Repo, ref.Number), "-f", "body="+body)
	return err
}

func (GitHub) AddLabel(ref IssueRef, label string) error {
	_, err := gh.Run("api", "-X", "POST",
		fmt.Sprintf("repos/%s/issues/%d/labels", ref.Repo, ref.Number), "-f", "labels[]="+label)
	return err
}

// Labels lists the labels repo defines, paginated: a repository with more
// than a page of labels must not read as missing the protocol's.
func (GitHub) Labels(repo string) ([]Label, error) {
	var items []struct {
		Name        string `json:"name"`
		Color       string `json:"color"`
		Description string `json:"description"`
	}
	if err := gh.JSON(&items, "api", "--paginate", fmt.Sprintf("repos/%s/labels?per_page=100", repo)); err != nil {
		return nil, err
	}
	labels := make([]Label, 0, len(items))
	for _, l := range items {
		labels = append(labels, Label{Name: l.Name, Color: l.Color, Description: l.Description})
	}
	return labels, nil
}

func (GitHub) CreateLabel(repo string, label Label) error {
	_, err := gh.Run("api", "-X", "POST", fmt.Sprintf("repos/%s/labels", repo),
		"-f", "name="+label.Name, "-f", "color="+label.Color, "-f", "description="+label.Description)
	return err
}

// UpdateLabel sends colour and description and no new_name: the label is
// addressed by the name the repository spells it with — which may differ
// in case, GitHub matching names case-insensitively — and restyling is not
// renaming.
func (GitHub) UpdateLabel(repo string, label Label) error {
	_, err := gh.Run("api", "-X", "PATCH",
		fmt.Sprintf("repos/%s/labels/%s", repo, url.PathEscape(label.Name)),
		"-f", "color="+label.Color, "-f", "description="+label.Description)
	return err
}

func (GitHub) Assign(ref IssueRef, login string) error {
	_, err := gh.Run("api", "-X", "POST",
		fmt.Sprintf("repos/%s/issues/%d/assignees", ref.Repo, ref.Number), "-f", "assignees[]="+login)
	return err
}

// Viewer resolves the current login. Installation tokens cannot call REST
// /user, so GraphQL viewer (which resolves to the bot user) is tried next.
func (GitHub) Viewer() (string, error) {
	var user struct {
		Login string `json:"login"`
	}
	if err := gh.JSON(&user, "api", "user"); err == nil {
		return user.Login, nil
	}
	var resp struct {
		Data struct {
			Viewer struct {
				Login string `json:"login"`
			} `json:"viewer"`
		} `json:"data"`
	}
	if err := gh.JSON(&resp, "api", "graphql", "-f", "query=query { viewer { login } }"); err != nil {
		return "", fmt.Errorf("cannot resolve current identity: %w", err)
	}
	return resp.Data.Viewer.Login, nil
}

func (GitHub) DevelopBranch(ref IssueRef, name string) error {
	_, err := gh.Run("issue", "develop", fmt.Sprint(ref.Number),
		"--repo", ref.Repo, "--name", name)
	return err
}

func (GitHub) ClosingPRs(ref IssueRef, includeClosed bool) ([]int, error) {
	owner, repo, ok := strings.Cut(ref.Repo, "/")
	if !ok {
		return nil, fmt.Errorf("bad repo ref %q", ref.Repo)
	}
	var resp struct {
		Data struct {
			Repository struct {
				Issue struct {
					Refs struct {
						Nodes []struct {
							Number int `json:"number"`
						} `json:"nodes"`
					} `json:"closedByPullRequestsReferences"`
				} `json:"issue"`
			} `json:"repository"`
		} `json:"data"`
	}
	query := fmt.Sprintf(`
query($owner: String!, $repo: String!, $num: Int!) {
  repository(owner: $owner, name: $repo) {
    issue(number: $num) {
      closedByPullRequestsReferences(first: 20, includeClosedPrs: %t) {
        nodes { number }
      }
    }
  }
}`, includeClosed)
	err := gh.JSON(&resp, "api", "graphql",
		"-f", "query="+query,
		"-f", "owner="+owner,
		"-f", "repo="+repo,
		"-F", fmt.Sprintf("num=%d", ref.Number))
	if err != nil {
		return nil, err
	}
	var numbers []int
	for _, n := range resp.Data.Repository.Issue.Refs.Nodes {
		numbers = append(numbers, n.Number)
	}
	return numbers, nil
}

// ClosingReferences asks GitHub what its own body parser made of the pull
// request: the issues the merge will close. One page is the whole answer —
// a PR's closing references are a handful by construction, and a body that
// somehow declared more than fifty has a larger problem than this note.
func (GitHub) ClosingReferences(repo string, number int) ([]TitledIssue, error) {
	owner, name, ok := strings.Cut(repo, "/")
	if !ok {
		return nil, fmt.Errorf("bad repo ref %q", repo)
	}
	var resp struct {
		Data struct {
			Repository struct {
				PullRequest struct {
					Refs struct {
						Nodes []struct {
							Number     int    `json:"number"`
							Title      string `json:"title"`
							Repository struct {
								NameWithOwner string `json:"nameWithOwner"`
							} `json:"repository"`
						} `json:"nodes"`
					} `json:"closingIssuesReferences"`
				} `json:"pullRequest"`
			} `json:"repository"`
		} `json:"data"`
	}
	query := `
query($owner: String!, $repo: String!, $num: Int!) {
  repository(owner: $owner, name: $repo) {
    pullRequest(number: $num) {
      closingIssuesReferences(first: 50) {
        nodes { number title repository { nameWithOwner } }
      }
    }
  }
}`
	if err := gh.JSON(&resp, "api", "graphql",
		"-f", "query="+query,
		"-f", "owner="+owner,
		"-f", "repo="+name,
		"-F", fmt.Sprintf("num=%d", number)); err != nil {
		return nil, err
	}
	var refs []TitledIssue
	for _, n := range resp.Data.Repository.PullRequest.Refs.Nodes {
		refs = append(refs, TitledIssue{
			Ref:   IssueRef{Repo: n.Repository.NameWithOwner, Number: n.Number},
			Title: n.Title,
		})
	}
	return refs, nil
}

func (GitHub) PRInfo(repo string, number int) (PR, error) {
	var view struct {
		State          string `json:"state"`
		ReviewDecision string `json:"reviewDecision"`
		HeadRefName    string `json:"headRefName"`
		HeadRefOid     string `json:"headRefOid"`
		IsCrossRepo    bool   `json:"isCrossRepository"`
		MergedAt       string `json:"mergedAt"`
		Author         struct {
			Login string `json:"login"`
		} `json:"author"`
		Reviews []struct {
			State  string `json:"state"`
			Author struct {
				Login string `json:"login"`
			} `json:"author"`
		} `json:"reviews"`
	}
	err := gh.JSON(&view, "pr", "view", fmt.Sprint(number), "--repo", repo,
		"--json", "state,author,reviews,reviewDecision,headRefName,headRefOid,isCrossRepository,mergedAt")
	if err != nil {
		return PR{}, err
	}
	pr := PR{
		Repo:      repo,
		Number:    number,
		Author:    strings.TrimPrefix(view.Author.Login, "app/"),
		HeadRef:   view.HeadRefName,
		HeadSHA:   view.HeadRefOid,
		CrossRepo: view.IsCrossRepo,
		Open:      view.State == "OPEN",
		Merged:    view.MergedAt != "",

		ReviewDecision: view.ReviewDecision,
	}
	latest := map[string]string{}
	for _, r := range view.Reviews {
		if r.State == "APPROVED" || r.State == "CHANGES_REQUESTED" || r.State == "DISMISSED" {
			latest[r.Author.Login] = r.State
		}
	}
	for login, state := range latest {
		if state == "APPROVED" {
			pr.ApprovedBy = append(pr.ApprovedBy, login)
		}
	}
	var checks []struct {
		Bucket string `json:"bucket"`
	}
	if err := gh.JSONLoose(&checks, "pr", "checks", fmt.Sprint(number), "--repo", repo, "--json", "bucket"); err != nil {
		if NoChecksReported(err) {
			pr.NoChecks = true
			return pr, nil
		}
		if perm := MissingChecksPermission(err); perm != "" {
			pr.ChecksUnreadable = perm
			return pr, nil
		}
		return PR{}, err
	}
	if len(checks) == 0 {
		pr.NoChecks = true
		return pr, nil
	}
	pr.ChecksOK = true
	for _, c := range checks {
		switch c.Bucket {
		case "pass", "skipping":
		case "pending":
			pr.ChecksPending = true
			pr.ChecksOK = false
		default:
			pr.ChecksOK = false
		}
	}
	return pr, nil
}

func (GitHub) MergePR(repo string, number int) error {
	_, err := gh.Run("pr", "merge", fmt.Sprint(number), "--repo", repo, "--rebase")
	return err
}

// MergePRBypass rebase-merges using the ruleset's administrator bypass.
// GitHub enforces eligibility: without a bypass actor covering the caller,
// this fails with the platform's own error, unmasked.
func (GitHub) MergePRBypass(repo string, number int) error {
	_, err := gh.Run("pr", "merge", fmt.Sprint(number), "--repo", repo, "--rebase", "--admin")
	return err
}

// MergeCommit reads the commit a merge left on the base branch. GitHub
// reports it for every merge method, the rebase included, where it is the
// last commit replayed onto the base — not the PR's head, which the rebase
// rewrote.
func (GitHub) MergeCommit(repo string, number int) (string, error) {
	var view struct {
		MergeCommit struct {
			OID string `json:"oid"`
		} `json:"mergeCommit"`
	}
	if err := gh.JSON(&view, "pr", "view", fmt.Sprint(number), "--repo", repo, "--json", "mergeCommit"); err != nil {
		return "", err
	}
	return view.MergeCommit.OID, nil
}

func (GitHub) CloseIssue(ref IssueRef, comment string) error {
	_, err := gh.Run("issue", "close", fmt.Sprint(ref.Number),
		"--repo", ref.Repo, "--comment", comment)
	return err
}

// Comments returns every comment on ref in GitHub's oldest-first order,
// paginated: latest-wins reads the newest verdict, and a milestone issue
// past a hundred comments carries it on a later page, where an unwalked
// listing left `milestone close` reporting VERDICT_MISSING for a
// requirement that was satisfied (#264). gh joins the pages in request
// order, so the concatenation is still oldest first.
func (GitHub) Comments(ref IssueRef) ([]Comment, error) {
	var raw []struct {
		Body string `json:"body"`
		URL  string `json:"html_url"`
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	}
	path := fmt.Sprintf("repos/%s/issues/%d/comments?per_page=100", ref.Repo, ref.Number)
	if err := gh.JSON(&raw, "api", "--paginate", path); err != nil {
		return nil, err
	}
	comments := make([]Comment, len(raw))
	for i, c := range raw {
		// The boundary, as in IssueBody: a comment body reaches the
		// record scans with LF line endings whatever GitHub stored.
		comments[i] = Comment{Author: c.User.Login, Body: NormalizeLineEndings(c.Body), URL: c.URL}
	}
	return comments, nil
}

func (GitHub) HasMilestoneDoc(repo string, n int) (bool, error) {
	var entries []struct {
		Name string `json:"name"`
	}
	err := gh.JSON(&entries, "api", fmt.Sprintf("repos/%s/contents/docs/milestones", repo))
	if err != nil {
		// A missing directory means no docs yet, not a failure.
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found") {
			return false, nil
		}
		return false, err
	}
	prefix := fmt.Sprintf("%d-", n)
	for _, e := range entries {
		if strings.HasPrefix(e.Name, prefix) && strings.HasSuffix(e.Name, ".md") {
			return true, nil
		}
	}
	return false, nil
}

func (GitHub) FileContent(repo, path string) ([]byte, error) {
	return gh.Run("api", fmt.Sprintf("repos/%s/contents/%s", repo, path),
		"-H", "Accept: application/vnd.github.raw+json")
}

const taskQuery = `
query($owner: String!, $repo: String!, $num: Int!) {
  repository(owner: $owner, name: $repo) {
    issue(number: $num) {
      title
      state
      assignees(first: 10) { nodes { login } }
      labels(first: 20) { nodes { name } }
      closedByPullRequestsReferences(first: 10, includeClosedPrs: false) {
        nodes { state }
      }
    }
  }
}`

func (GitHub) Task(ref IssueRef) (Task, error) {
	owner, repo, ok := strings.Cut(ref.Repo, "/")
	if !ok {
		return Task{}, fmt.Errorf("bad repo ref %q", ref.Repo)
	}
	var resp struct {
		Data struct {
			Repository struct {
				Issue *struct {
					Title     string `json:"title"`
					State     string `json:"state"`
					Assignees struct {
						Nodes []struct {
							Login string `json:"login"`
						} `json:"nodes"`
					} `json:"assignees"`
					Labels struct {
						Nodes []struct {
							Name string `json:"name"`
						} `json:"nodes"`
					} `json:"labels"`
					ClosedByPullRequestsReferences struct {
						Nodes []struct {
							State string `json:"state"`
						} `json:"nodes"`
					} `json:"closedByPullRequestsReferences"`
				} `json:"issue"`
			} `json:"repository"`
		} `json:"data"`
	}
	err := gh.JSON(&resp, "api", "graphql",
		"-f", "query="+taskQuery,
		"-f", "owner="+owner,
		"-f", "repo="+repo,
		"-F", fmt.Sprintf("num=%d", ref.Number))
	if err != nil {
		return Task{}, err
	}
	issue := resp.Data.Repository.Issue
	if issue == nil {
		return Task{}, fmt.Errorf("%s: issue not found", ref)
	}
	t := Task{
		Ref:    ref,
		Title:  issue.Title,
		Closed: issue.State == "CLOSED",
	}
	for _, a := range issue.Assignees.Nodes {
		t.Assignees = append(t.Assignees, a.Login)
	}
	for _, l := range issue.Labels.Nodes {
		t.Labels = append(t.Labels, l.Name)
	}
	for _, pr := range issue.ClosedByPullRequestsReferences.Nodes {
		if pr.State == "OPEN" {
			t.OpenLinkedPR = true
		}
	}
	return t, nil
}

func (GitHub) LinkedBranches(ref IssueRef) ([]string, error) {
	owner, repo, ok := strings.Cut(ref.Repo, "/")
	if !ok {
		return nil, fmt.Errorf("bad repo ref %q", ref.Repo)
	}
	var resp struct {
		Data struct {
			Repository struct {
				Issue struct {
					LinkedBranches struct {
						Nodes []struct {
							Ref struct {
								Name string `json:"name"`
							} `json:"ref"`
						} `json:"nodes"`
					} `json:"linkedBranches"`
				} `json:"issue"`
			} `json:"repository"`
		} `json:"data"`
	}
	query := `
query($owner: String!, $repo: String!, $num: Int!) {
  repository(owner: $owner, name: $repo) {
    issue(number: $num) { linkedBranches(first: 20) { nodes { ref { name } } } }
  }
}`
	if err := gh.JSON(&resp, "api", "graphql", "-f", "query="+query,
		"-F", "owner="+owner, "-F", "repo="+repo, "-F", fmt.Sprintf("num=%d", ref.Number)); err != nil {
		return nil, err
	}
	var names []string
	for _, n := range resp.Data.Repository.Issue.LinkedBranches.Nodes {
		if n.Ref.Name != "" {
			names = append(names, n.Ref.Name)
		}
	}
	return names, nil
}

func (GitHub) TaskBranches(repo string) ([]string, bool, error) {
	owner, name, ok := strings.Cut(repo, "/")
	if !ok {
		return nil, false, fmt.Errorf("bad repo ref %q", repo)
	}
	var resp struct {
		Data struct {
			Repository struct {
				Refs struct {
					Nodes []struct {
						Name string `json:"name"`
					} `json:"nodes"`
					PageInfo struct {
						HasNextPage bool `json:"hasNextPage"`
					} `json:"pageInfo"`
				} `json:"refs"`
			} `json:"repository"`
		} `json:"data"`
	}
	// Deliberately one page: this listing reports hasNextPage instead of
	// walking the cursor, and the sweep that reads it says "swept in part"
	// when a repo carries more task branches than a page holds (SPEC §6).
	query := `
query($owner: String!, $repo: String!, $prefix: String!) {
  repository(owner: $owner, name: $repo) {
    refs(refPrefix: $prefix, first: 100) {
      nodes { name }
      pageInfo { hasNextPage }
    }
  }
}`
	if err := gh.JSON(&resp, "api", "graphql", "-f", "query="+query,
		"-f", "owner="+owner, "-f", "repo="+name,
		"-f", "prefix=refs/heads/"+TaskBranchPrefix); err != nil {
		return nil, false, err
	}
	// GitHub returns each node's name with the queried prefix removed, so
	// the branch name is rebuilt from the same constant the query used.
	var branches []string
	for _, n := range resp.Data.Repository.Refs.Nodes {
		if n.Name != "" {
			branches = append(branches, TaskBranchPrefix+n.Name)
		}
	}
	return branches, resp.Data.Repository.Refs.PageInfo.HasNextPage, nil
}

// OpenPRsForBranch asks the pulls listing for open PRs with this head. gh
// builds the query string, so a branch name's slashes need no escaping of
// ours; the `head` filter's own grammar is `<owner>:<ref>`. The listing is
// paginated: a branch the sweep would keep for its open PR must not lose
// that PR to a page boundary (#264).
func (GitHub) OpenPRsForBranch(repo, branch string) ([]int, error) {
	owner, _, ok := strings.Cut(repo, "/")
	if !ok {
		return nil, fmt.Errorf("bad repo ref %q", repo)
	}
	var prs []struct {
		Number int `json:"number"`
	}
	if err := gh.JSON(&prs, "api", "--paginate", "repos/"+repo+"/pulls", "-X", "GET",
		"-f", "state=open", "-f", "per_page=100",
		"-f", "head="+owner+":"+branch); err != nil {
		return nil, err
	}
	numbers := make([]int, 0, len(prs))
	for _, pr := range prs {
		numbers = append(numbers, pr.Number)
	}
	return numbers, nil
}

func (g GitHub) BranchAhead(repo, branch string) (int, string, error) {
	info, err := g.RepoInfo(repo)
	if err != nil {
		return 0, "", err
	}
	var cmp struct {
		AheadBy int `json:"ahead_by"`
		Commits []struct {
			SHA string `json:"sha"`
		} `json:"commits"`
	}
	if err := gh.JSON(&cmp, "api", fmt.Sprintf("repos/%s/compare/%s...%s", repo, info.DefaultBranch, branch)); err != nil {
		return 0, "", err
	}
	var ref struct {
		Object struct {
			SHA string `json:"sha"`
		} `json:"object"`
	}
	if err := gh.JSON(&ref, "api", fmt.Sprintf("repos/%s/git/ref/heads/%s", repo, branch)); err != nil {
		return 0, "", err
	}
	return cmp.AheadBy, ref.Object.SHA, nil
}

func (GitHub) DeleteBranch(repo, branch string) error {
	_, err := gh.Run("api", "-X", "DELETE", fmt.Sprintf("repos/%s/git/refs/heads/%s", repo, branch))
	return err
}

func (GitHub) RepoInfo(repo string) (RepoInfo, error) {
	var r struct {
		DefaultBranch       string `json:"default_branch"`
		DeleteBranchOnMerge bool   `json:"delete_branch_on_merge"`
	}
	if err := gh.JSON(&r, "api", "repos/"+repo); err != nil {
		return RepoInfo{}, err
	}
	return RepoInfo{DefaultBranch: r.DefaultBranch, DeleteBranchOnMerge: r.DeleteBranchOnMerge}, nil
}
