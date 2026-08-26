package cli

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	codecrew "github.com/radiusred/gh-codecrew"
)

const hubConfigScaffold = `codecrew: "1.0" # protocol version (SPEC.md §5): a different major is refused; not the CLI release — see codecrew version
hub: self

# Role routing: who holds each role (SPEC §5). Declare all four at
# onboarding. ~ routes the role to you, the human operator; replace it with
# a GitHub App slug or a username to delegate, or an org/team-slug to
# route the role to any member of a GitHub team.
roles:
  implementer: { identity: ~ }
  reviewer: { identity: ~ }
  qa: { identity: ~ }
  doc-synthesizer: { identity: ~ }
`

const roadmapScaffold = `# Roadmap

| Milestone | Goal | Tracking issue | Status |
|-----------|------|----------------|--------|
`

const agentsScaffold = `# Agents

This repository is part of a CodeCrew project — coordination state lives in
GitHub issues and PRs, per the protocol at
https://github.com/radiusred/gh-codecrew (SPEC.md).

- ` + "`.codecrew.yml`" + ` names the hub; the hub's ` + "`roles/`" + ` holds the role
  contracts. Read the contract for the role you were dispatched as before
  doing anything else.
- ` + "`gh codecrew status`" + ` shows where the project is; ` + "`gh codecrew help`" + `
  lists the workflow verbs. Blocked gates refuse with
  ` + "`refused[CODE]: detail`" + ` — act on the code, don't work around it.
- Plans before commits, decisions recorded when made, and the verifier is
  never the doer. Reviews are model reviews: a clean-context session under
  the reviewer contract — even in pure solo, where its findings land as a
  PR comment before the operator confirms.
- **Contract drift.** ` + "`codecrew status`" + ` reports when a ` + "`roles/`" + ` contract
  differs from the one embedded in the installed CLI. When it does, the
  coordination layer compares (` + "`codecrew roles diff <role>`" + `, full upstream
  text via ` + "`codecrew roles show <role> --latest`" + `), decides what to adopt —
  contracts are this project's own fork, and local conventions are
  legitimate — and routes the reconciliation through a normal task and PR
  with the decision recorded. Never overwrite blindly.
- **Dispatch authorization.** If you are the operator's primary session —
  not dispatched as any specific role — then when a role is routed to a
  GitHub App and that role's action is needed (a review, a verdict),
  dispatching a clean-context sub-agent session as that App is authorized
  and expected; use the dispatch prompt in the hub's docs/identities.md. A
  session dispatched *as* a role never dispatches another role — that
  belongs to its coordination layer (platform, orchestrating session, or
  operator) — and never chooses or briefs its own judge.
`

// scaffold writes the greenfield files into dir. Hub mode (hub == "self")
// writes the full set; spoke mode writes only the pointer. Existing files
// are never touched — they are reported as skipped.
func scaffold(dir, hub string, contracts fs.FS) (written, skipped []string, err error) {
	files := map[string]string{
		".codecrew.yml": hubConfigScaffold,
	}
	if hub != "self" {
		files[".codecrew.yml"] = fmt.Sprintf("codecrew: \"%s\" # protocol version (SPEC.md §5): a different major is refused; not the CLI release\nhub: %s\n", protocolVersion, hub)
	} else {
		files["ROADMAP.md"] = roadmapScaffold
		files["AGENTS.md"] = agentsScaffold
		entries, err := fs.ReadDir(contracts, "roles")
		if err != nil {
			return nil, nil, err
		}
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), localSuffix) {
				continue // never scaffold one project's extensions into another
			}
			data, err := fs.ReadFile(contracts, "roles/"+e.Name())
			if err != nil {
				return nil, nil, err
			}
			// Provenance stamp: names the release these contracts shipped
			// with, so drift can be judged three-way later (the base is
			// fetchable from the upstream repo at this version).
			files[filepath.Join("roles", e.Name())] = contractStamp(e.Name()) + string(data)
		}
	}
	for rel, content := range files {
		path := filepath.Join(dir, rel)
		if _, err := os.Stat(path); err == nil {
			skipped = append(skipped, rel)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return written, skipped, err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return written, skipped, err
		}
		written = append(written, rel)
	}
	slices.Sort(written)
	slices.Sort(skipped)
	return written, skipped, nil
}

// inGitRepo reports whether dir carries a .git entry (directory in normal
// clones, file in worktrees). A conservative check: the protocol needs a
// GitHub repo, and an agent shouldn't have to discover its absence itself.
func inGitRepo(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// initCmd scaffolds a new CodeCrew repo: hub mode by default, spoke mode
// with --hub owner/repo (pointer only — contracts and roadmap live in the
// hub).
func initCmd(w io.Writer, args []string) error {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	hub := flags.String("hub", "self", "hub repo this spoke points at (default: this repo is the hub)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	written, skipped, err := scaffold(".", *hub, codecrew.Roles)
	if err != nil {
		return err
	}
	for _, f := range written {
		fmt.Fprintf(w, "wrote %s\n", f)
	}
	for _, f := range skipped {
		fmt.Fprintf(w, "kept existing %s\n", f)
	}
	if !inGitRepo(".") {
		fmt.Fprintln(w, "\nnote: this directory is not a git repository — the protocol lives in GitHub.")
		fmt.Fprintln(w, "first: git init && git add -A && git commit -m \"chore: scaffold codecrew\" &&")
		fmt.Fprintln(w, "       gh repo create <owner>/<name> --private --source=. --push")
	}
	if *hub == "self" {
		fmt.Fprintln(w, "\nnext: declare who holds each role in .codecrew.yml (~ = you),")
		fmt.Fprintln(w, "commit, then `codecrew milestone new --title \"...\" --goal \"...\"`")
	} else {
		fmt.Fprintf(w, "\nnext: commit the pointer; tasks for this spoke attach to milestones in %s\n", *hub)
	}
	return nil
}
