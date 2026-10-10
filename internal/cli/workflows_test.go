package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// repoWorkflows is the workflow half of a fake repository, embedded in the
// task new and task start fakes: what the Actions listing answers, and the
// files behind it. Its zero value is a repository with one workflow on
// pull_request, so a test about something else meets no workflow note;
// a test about the note sets flows (an empty, non-nil slice for none at
// all), listErr, or a file.
type repoWorkflows struct {
	flows    []tracker.Workflow
	listErr  error
	files    map[string]string // path → content
	fileErr  error             // what every file read fails with, when set
	listedAt []string          // the repos the listing was asked about
}

const reportingWorkflow = "name: CI\non: [pull_request]\njobs:\n  test:\n    runs-on: ubuntu-24.04\n    steps:\n      - run: true\n"

func (r *repoWorkflows) workflows(repo string) ([]tracker.Workflow, error) {
	r.listedAt = append(r.listedAt, repo)
	if r.flows == nil && r.listErr == nil {
		return []tracker.Workflow{{Path: ".github/workflows/ci.yml", State: "active"}}, nil
	}
	return r.flows, r.listErr
}

func (r *repoWorkflows) file(path string) ([]byte, error) {
	if r.fileErr != nil {
		return nil, r.fileErr
	}
	if r.files == nil {
		return []byte(reportingWorkflow), nil
	}
	content, ok := r.files[path]
	if !ok {
		return nil, errors.New("gh api: Not Found (HTTP 404)")
	}
	return []byte(content), nil
}

func (f *startFake) Workflows(repo string) ([]tracker.Workflow, error) { return f.workflows(repo) }
func (f *startFake) FileContent(_, path string) ([]byte, error)        { return f.file(path) }
func (f *taskNewFake) Workflows(repo string) ([]tracker.Workflow, error) {
	return f.repoWorkflows.workflows(repo)
}
func (f *taskNewFake) FileContent(_, path string) ([]byte, error) { return f.file(path) }

// Every shape Actions accepts for `on:` is read, and only pull_request
// counts: pull_request_target runs the base branch's workflow, a different
// event.
func TestOnPullRequest(t *testing.T) {
	for _, tc := range []struct {
		name, yaml string
		want       bool
	}{
		{"one event", "on: pull_request\n", true},
		{"a list", "on: [push, pull_request]\n", true},
		{"a block list", "on:\n  - push\n  - pull_request\n", true},
		{"a map", "on:\n  pull_request:\n  push:\n    branches: [main]\n", true},
		{"a map with filters", "on:\n  pull_request:\n    branches: [main]\n    paths-ignore: ['docs/**']\n", true},
		{"a flow map", "on: {push: {}, pull_request: {}}\n", true},
		{"a quoted key", "\"on\": pull_request\n", true},
		{"push only", "on: push\n", false},
		{"push and a schedule", "on:\n  push:\n  schedule:\n    - cron: '0 0 * * *'\n", false},
		{"pull_request_target alone", "on: pull_request_target\n", false},
		{"pull_request_target in a list", "on: [push, pull_request_target]\n", false},
		{"pull_request_target in a map", "on:\n  pull_request_target:\n    types: [opened]\n", false},
		{"the event named only in a job", "on: push\njobs:\n  pull_request:\n    runs-on: x\n", false},
		{"no on at all", "name: nothing\n", false},
		{"empty", "", false},
	} {
		got, err := onPullRequest([]byte(tc.yaml))
		if err != nil {
			t.Errorf("%s: %v", tc.name, err)
			continue
		}
		if got != tc.want {
			t.Errorf("%s: onPullRequest = %v, want %v", tc.name, got, tc.want)
		}
	}
	if _, err := onPullRequest([]byte("on: [unclosed\n")); err == nil {
		t.Error("a file that is not YAML must be an error: the caller cannot say")
	}
}

// The repositories the note is for, and the ones it is not: no workflow at
// all, workflows none of which runs on pull_request (a disabled one and a
// GitHub-managed one included), against one that does, and every read that
// fails — which says nothing rather than guess.
func workflowCases() []struct {
	name string
	repo repoWorkflows
	note bool
} {
	push := "on: push\n"
	return []struct {
		name string
		repo repoWorkflows
		note bool
	}{
		{"no workflow at all", repoWorkflows{flows: []tracker.Workflow{}}, true},
		{"a workflow lacking pull_request", repoWorkflows{
			flows: []tracker.Workflow{{Path: ".github/workflows/release.yml", State: "active"}},
			files: map[string]string{".github/workflows/release.yml": push}}, true},
		{"only pull_request_target", repoWorkflows{
			flows: []tracker.Workflow{{Path: ".github/workflows/label.yml", State: "active"}},
			files: map[string]string{".github/workflows/label.yml": "on: pull_request_target\n"}}, true},
		{"a disabled pull_request workflow", repoWorkflows{
			flows: []tracker.Workflow{{Path: ".github/workflows/ci.yml", State: "disabled_manually"}},
			files: map[string]string{".github/workflows/ci.yml": reportingWorkflow}}, true},
		{"only GitHub-managed workflows", repoWorkflows{
			flows: []tracker.Workflow{{Path: "dynamic/pages/pages-build-deployment", State: "active"}},
			files: map[string]string{}}, true},
		{"one workflow on pull_request among others", repoWorkflows{
			flows: []tracker.Workflow{{Path: ".github/workflows/release.yml", State: "active"}, {Path: ".github/workflows/ci.yml", State: "active"}},
			files: map[string]string{".github/workflows/release.yml": push, ".github/workflows/ci.yml": reportingWorkflow}}, false},
		{"the listing 404", repoWorkflows{listErr: errors.New("gh api: Not Found (HTTP 404)")}, false},
		{"the listing 403", repoWorkflows{listErr: errors.New("gh api: Resource not accessible by integration (HTTP 403)")}, false},
		{"a workflow file unreadable", repoWorkflows{
			flows:   []tracker.Workflow{{Path: ".github/workflows/ci.yml", State: "active"}},
			fileErr: errors.New("gh api: Not Found (HTTP 404)")}, false},
		{"a workflow file that is not YAML", repoWorkflows{
			flows: []tracker.Workflow{{Path: ".github/workflows/ci.yml", State: "active"}},
			files: map[string]string{".github/workflows/ci.yml": "on: [unclosed\n"}}, false},
	}
}

// #391 at task start: the note, once, beside the branch — and nothing
// when a workflow will report or the read cannot say. The start completes
// either way.
func TestTaskStartNotesARepositoryWithNoPullRequestWorkflow(t *testing.T) {
	for _, tc := range workflowCases() {
		f := &startFake{task: startingTask(), viewer: "myorg-coder[bot]", repoWorkflows: tc.repo}
		var out bytes.Buffer
		if err := runTaskStart(startCtx(f, crewRoles), &out, f.task.Ref); err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		got := out.String()
		if !strings.Contains(got, "started o/r#7 as @myorg-coder[bot]") {
			t.Errorf("%s: the start did not complete:\n%s", tc.name, got)
		}
		if len(f.listedAt) != 1 || f.listedAt[0] != "o/r" {
			t.Errorf("%s: the workflows read asked about %q, want the task's repository", tc.name, f.listedAt)
		}
		assertWorkflowNote(t, tc.name, got, "o/r", tc.note)
	}
}

// And at task new, about the repository the task is created in — the
// --repo spoke, where its pull request will be opened, not the hub.
func TestTaskNewNotesARepositoryWithNoPullRequestWorkflow(t *testing.T) {
	recordSleeps(t)
	for _, tc := range workflowCases() {
		f := &taskNewFake{open: [][]tracker.Milestone{{openMilestone(233, "M11: Housekeeping")}}, repoWorkflows: tc.repo}
		var out bytes.Buffer
		if err := runTaskNew(taskNewCtx(f), &out, 11, "o/spoke", "Cut the README", "g", "M11-R1", nil); err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		got := out.String()
		if !strings.HasPrefix(got, "created task o/spoke#21 as a sub-issue of o/hub#233\n") {
			t.Errorf("%s: the task was not created first:\n%s", tc.name, got)
		}
		if len(f.listedAt) != 1 || f.listedAt[0] != "o/spoke" {
			t.Errorf("%s: the workflows read asked about %q, want the task's repository", tc.name, f.listedAt)
		}
		assertWorkflowNote(t, tc.name, got, "o/spoke", tc.note)
	}
}

func assertWorkflowNote(t *testing.T, name, got, repo string, want bool) {
	t.Helper()
	note := noChecksNote(repo)
	if !want {
		if strings.Contains(got, "note:") {
			t.Errorf("%s: said something:\n%s", name, got)
		}
		return
	}
	if strings.Count(got, "note:") != 1 || !strings.Contains(got, note+"\n") {
		t.Errorf("%s: want exactly the workflow note:\n%s", name, got)
	}
	for _, w := range []string{"note: no workflow in " + repo + " runs on pull_request", "task finish will refuse NO_CHECKS", "SPEC §8", "add a workflow that runs on pull_request", "let it report", "docs/first-milestone.md"} {
		if !strings.Contains(note, w) {
			t.Errorf("%s: the note does not say %q:\n%s", name, w, note)
		}
	}
}
