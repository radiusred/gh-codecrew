package cli

import (
	"fmt"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// What the branch policy says about a merge before task finish asks GitHub
// for one (#392). Everything here is advisory: it reads, it prints `note:`
// lines, and a read that fails says nothing — a plan without rulesets
// answers the rules endpoint 404 or 403, and absence of policy is not a
// finding. No gate reads it and no exit status depends on it.

// recentPRSample is how many of the newest pull requests into the default
// branch task start reads for the checks that have reported there: no PR
// exists for the task yet, so the repository's recent ones stand in for
// "has anything ever reported this check". Three heads cost seven calls in
// a repository whose policy requires checks, and none in one whose does not.
const recentPRSample = 3

// requiredChecks lists the status-check contexts the rules require, once
// each, in the order the rules name them.
func requiredChecks(rules []tracker.BranchRule) []string {
	var out []string
	seen := map[string]bool{}
	for _, r := range rules {
		for _, c := range r.Checks {
			if !seen[c] {
				seen[c] = true
				out = append(out, c)
			}
		}
	}
	return out
}

// ruleTypes lists the rule types in force, once each, in the order read.
func ruleTypes(rules []tracker.BranchRule) []string {
	var out []string
	seen := map[string]bool{}
	for _, r := range rules {
		if !seen[r.Type] {
			seen[r.Type] = true
			out = append(out, r.Type)
		}
	}
	return out
}

// unreportedChecks returns the required contexts that no check run and no
// commit status has reported under on any of heads. Nil when there is
// nothing to compare — no requirement, no head — or when a head's checks
// cannot be read: a context missing from a listing that could not be read
// is not a context nothing reports.
func unreportedChecks(t tracker.Tracker, repo string, required, heads []string) []string {
	if len(required) == 0 || len(heads) == 0 {
		return nil
	}
	reported := map[string]bool{}
	for _, sha := range heads {
		names, err := t.ReportedChecks(repo, sha)
		if err != nil {
			return nil
		}
		for _, n := range names {
			reported[n] = true
		}
	}
	var out []string
	for _, c := range required {
		if !reported[c] {
			out = append(out, c)
		}
	}
	return out
}

// unreportedNote is the line for one required context nothing reports;
// where says what was looked at.
func unreportedNote(base, context, where string) string {
	return fmt.Sprintf("note: the policy on %s requires the status check %q, and nothing has reported it on %s — no workflow in this repository reports it, so the merge will block or rely on a bypass until one does", base, context, where)
}

// blockedNote is said when GitHub reports the PR BLOCKED although every
// protocol gate passed: the merge is a request GitHub decides, and the
// line names who asks and, when the rules could be read, what is in force.
// It reports the record; it does not predict — whether the caller holds a
// bypass is readable only by an administrator, and is not read.
func blockedNote(pr tracker.PR, caller string, rules []tracker.BranchRule, rulesRead bool) string {
	line := fmt.Sprintf("note: GitHub reports PR #%d as BLOCKED although every protocol gate passed — the merge is asked of GitHub as @%s, and GitHub honours a ruleset bypass that caller holds and refuses otherwise", pr.Number, caller)
	if types := ruleTypes(rules); rulesRead && len(types) > 0 {
		line += fmt.Sprintf(" (rules in force on %s: %s)", pr.BaseRef, strings.Join(types, ", "))
	}
	return line
}

// startPolicyNotes is task start's half: the default branch's required
// checks that none of the newest pull requests into it has reported —
// the policy the task's pull request will meet, said before its first
// commit, when adding the missing workflow is cheapest.
func startPolicyNotes(t tracker.Tracker, repo string) []string {
	info, err := t.RepoInfo(repo)
	if err != nil || info.DefaultBranch == "" {
		return nil
	}
	base := info.DefaultBranch
	rules, err := t.BranchRules(repo, base)
	if err != nil {
		return nil
	}
	required := requiredChecks(rules)
	if len(required) == 0 {
		return nil
	}
	heads, err := t.RecentPRHeads(repo, base, recentPRSample)
	if err != nil {
		return nil
	}
	var notes []string
	for _, c := range unreportedChecks(t, repo, required, heads) {
		notes = append(notes, unreportedNote(base, c, fmt.Sprintf("the newest pull requests into %s", base)))
	}
	return notes
}
