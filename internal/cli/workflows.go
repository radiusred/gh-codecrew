package cli

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// Whether anything will report a check on the task's pull request, said
// before there is one (#391). task finish refuses NO_CHECKS on a pull
// request that reports no check (SPEC §8), and a repository with no
// workflow on pull_request produces none — so a fresh hub's first task was
// chartered, built and approved before the refusal named the gap. task new
// and task start now say it as a note: advisory, read-only, and silent on
// any read that fails, since a workflow the read could not see is not a
// workflow that is missing.

// workflowsDir is where the workflows a repository writes live; anything
// else the listing names is GitHub-managed (Pages, Dependabot), with no
// file to read.
const workflowsDir = ".github/workflows/"

// onPullRequest reports whether a workflow file's top-level `on:` includes
// the pull_request event, in any of the three shapes Actions accepts: one
// event (`on: pull_request`), a list (`on: [push, pull_request]`), or a map
// keyed by event (`on:` then `pull_request:` beneath it). pull_request_target
// is a different event — it runs the base branch's workflow — and does not
// count. A file that is not YAML is an error: the caller cannot say.
func onPullRequest(data []byte) (bool, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return false, err
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return false, nil
	}
	top := doc.Content[0].Content
	for i := 0; i+1 < len(top); i += 2 {
		if top[i].Value != "on" {
			continue
		}
		on := top[i+1]
		switch on.Kind {
		case yaml.ScalarNode:
			return on.Value == "pull_request", nil
		case yaml.SequenceNode:
			for _, e := range on.Content {
				if e.Kind == yaml.ScalarNode && e.Value == "pull_request" {
					return true, nil
				}
			}
		case yaml.MappingNode:
			for j := 0; j+1 < len(on.Content); j += 2 {
				if on.Content[j].Value == "pull_request" {
					return true, nil
				}
			}
		}
		return false, nil
	}
	return false, nil
}

// lacksPullRequestWorkflow reports whether repo has no active workflow
// whose `on:` includes pull_request — none listed at all, or none of those
// listed triggered by it. Each active workflow under .github/workflows is
// read from the default branch. Any read that fails — the listing (403,
// 404), or one file (a workflow that has run only from a branch so far) —
// answers false: the read could not say, and the note says nothing.
// A disabled workflow does not run, so it does not count.
func lacksPullRequestWorkflow(t tracker.Tracker, repo string) bool {
	flows, err := t.Workflows(repo)
	if err != nil {
		return false
	}
	for _, f := range flows {
		if f.State != "active" || !strings.HasPrefix(f.Path, workflowsDir) {
			continue
		}
		data, err := t.FileContent(repo, f.Path)
		if err != nil {
			return false
		}
		on, err := onPullRequest(data)
		if err != nil {
			return false
		}
		if on {
			return false
		}
	}
	return true
}

// noChecksNote is the line task new and task start print for a repository
// with no workflow on pull_request: the refusal task finish will meet, and
// #391's remedy — a workflow on pull_request, in place before the pull
// request is reviewed rather than after.
func noChecksNote(repo string) string {
	return fmt.Sprintf("note: no workflow in %s runs on pull_request — task finish will refuse NO_CHECKS until a check reports on the task's pull request (SPEC §8: the gate cannot be satisfied by absence); add a workflow that runs on pull_request, on the task's branch or before it, and let it report ("+U+"/docs/first-milestone.md shows a ten-line one)", repo)
}

// workflowNotes is the note, or nothing, for the repository a task's pull
// request will be opened in.
func workflowNotes(t tracker.Tracker, repo string) []string {
	if lacksPullRequestWorkflow(t, repo) {
		return []string{noChecksNote(repo)}
	}
	return nil
}
