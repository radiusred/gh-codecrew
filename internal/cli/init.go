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
	"github.com/radiusred/gh-codecrew/internal/config"
)

// U is the upstream repository every scaffolded reference resolves to —
// an adopter's hub holds no copy of the protocol or its docs.
const U = "https://github.com/radiusred/gh-codecrew/blob/main"

const hubConfigScaffold = `codecrew: "%s" # protocol version (SPEC.md §5): a different major is refused; not the CLI release — see gh codecrew version
hub: self

# Role routing: who holds each role (SPEC §5). Declare all five at
# onboarding. The identity is typed — one of exactly four forms:
#   ~                    you, the human operator (and any session under
#                        your own auth)
#   app:<slug>           a GitHub App, an agent acting as itself
#   user:<login>         one named human
#   team:<org>/<slug>    any member of the team
# A value carrying no type prefix is refused (IDENTITY_UNTYPED): an App
# slug and a username are the same string, and the protocol treats them
# differently. The coordinator is the seat that dispatches the other four
# (SPEC §7); unrouted, it is you.
roles:
  implementer: { identity: ~ }
  reviewer: { identity: ~ }
  qa: { identity: ~ }
  doc-synthesizer: { identity: ~ }
  coordinator: { identity: ~ }
`

const roadmapScaffold = `# Roadmap

| Milestone | Goal | Tracking issue | Status |
|-----------|------|----------------|--------|
`

// agentsScaffold is the CodeCrew instructions themselves, written to
// .codecrew/AGENTS.md in hub and spoke alike: CodeCrew's file, under
// CodeCrew's directory, so a later init or migrate rewrites it whole
// without touching a line the project wrote. The root AGENTS.md an adopter
// owns only points at it (entryPointLines).
const agentsScaffold = `# Agents

This repository is part of a CodeCrew project — coordination state lives in
GitHub issues and PRs, per the protocol at
` + U + `/SPEC.md.

- ` + "`.codecrew/config.yml`" + ` names the hub; the hub's ` + "`.codecrew/roles/`" + `
  holds the role contracts. Read the contract for the role you were
  dispatched as before doing anything else — ` + "`gh codecrew roles show <role>`" + `
  prints it with this project's ` + "`.codecrew/roles/<role>.local.md`" + ` extension
  appended (blank until the project writes one; in a hub ` + "`init`" + ` scaffolds the
  file with a comment saying what belongs there).
- ` + "`gh codecrew status`" + ` shows where the project is; ` + "`gh codecrew help`" + `
  lists the workflow verbs. Blocked gates refuse with
  ` + "`refused[CODE]: detail`" + ` — act on the code, don't work around it.
- Plans before commits, decisions recorded when made, and the verifier is
  never the doer. Reviews are model reviews: a clean-context session under
  the reviewer contract — even in pure solo, where its findings land as a
  PR comment before the operator confirms.
- **Contract drift.** ` + "`gh codecrew status`" + ` reports when a ` + "`.codecrew/roles/`" + ` contract
  differs from the one embedded in the installed CLI. When it does, the
  coordination layer compares (` + "`gh codecrew roles diff <role>`" + `, full upstream
  text via ` + "`gh codecrew roles show <role> --latest`" + `), decides what to adopt —
  contracts are this project's own fork, and local conventions are
  legitimate — and routes the reconciliation through a normal task and PR
  with the decision recorded. Never overwrite blindly.
- **Dispatch authorization.** If you are the operator's primary session —
  not dispatched as any specific role — then when a role is routed to a
  GitHub App and that role's action is needed (a review, a verdict),
  dispatching a clean-context sub-agent session as that App is authorized
  and expected; use the dispatch prompt in
  ` + U + `/docs/identities.md. A
  session dispatched *as* a role never dispatches another role — that
  belongs to its coordination layer (platform, orchestrating session, or
  operator) — and never chooses or briefs its own judge.
`

// entryPointLines connects a repository's root entry point to the CodeCrew
// instructions in .codecrew/AGENTS.md, and carries both forms because
// neither alone reaches every harness: the sentence is what a harness
// reading plain markdown follows, and the bare @-import is what Claude Code
// resolves — transitively, through CLAUDE.md's @AGENTS.md — which is why it
// must not be wrapped in backticks, read there as literal text. This block
// is the whole of the root AGENTS.md scaffold and, verbatim, what init
// prints when a root entry point already exists and is kept.
const entryPointLines = "This repository is a CodeCrew project — the instructions for an agent\n" +
	"dispatched here are in `" + config.AgentsFile + "`; read that file first.\n" +
	"\n" +
	"@" + config.AgentsFile + "\n"

// agentsPointerScaffold is the root AGENTS.md init writes: a heading and the
// pointer, nothing else. Everything that can change with a release lives in
// .codecrew/AGENTS.md, so the file an adopter edits stays two lines long.
const agentsPointerScaffold = `# Agents

` + entryPointLines

// rootEntryPoints are the two root files a harness discovers on its own,
// in the order both verbs report them. Neither init nor migrate ever
// overwrites one; when either keeps one that does not already reach the
// instructions it prints the lines to add (reachesInstructions).
var rootEntryPoints = []string{"AGENTS.md", "CLAUDE.md"}

// rootEntryPointScaffolds is what each root entry point is written from
// when it is absent — the one statement of that mapping, so init and
// migrate write the same bytes rather than each carrying a copy. An absent
// file is CodeCrew's to write; a file already there is the project's, and
// is only ever reported (#301).
var rootEntryPointScaffolds = map[string]string{
	"AGENTS.md": agentsPointerScaffold,
	"CLAUDE.md": claudeScaffold,
}

// reachesInstructions reports whether a root entry point already leads to
// .codecrew/AGENTS.md: by naming the path — the sentence and the @-import
// both do — or, for CLAUDE.md, through its @AGENTS.md import into a root
// AGENTS.md that does. A file that already arrives needs no line pasted
// into it, so a rerun in a repo init itself scaffolded reports a plain
// skip and asks for nothing (SPEC §6, idempotency).
//
// pending names the root files the caller is about to write and has not
// written yet, so a verb can judge the state it is creating rather than
// the one it found: every one of them is written from
// rootEntryPointScaffolds, which reaches by construction. init passes nil
// — it judges only what it kept — and migrate passes what its dry run
// would write, which is how the preview and the live run agree.
func reachesInstructions(dir, name string, pending map[string]bool) bool {
	// A pending file is not on disk yet, so judge the bytes that are about
	// to be there: the question is about the repository the verb is
	// making, not the one it found.
	content := rootEntryPointScaffolds[name]
	if !pending[name] {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return false
		}
		content = string(data)
	}
	if strings.Contains(content, config.AgentsFile) {
		return true
	}
	// One hop, and only this one: CLAUDE.md's own scaffold reaches the
	// instructions through the root pointer rather than naming them.
	return name == "CLAUDE.md" && strings.Contains(content, "@AGENTS.md") && reachesInstructions(dir, "AGENTS.md", pending)
}

// claudeScaffold bridges Claude Code to the harness-neutral entry point:
// Claude Code loads CLAUDE.md, never AGENTS.md (code.claude.com/docs/en/memory),
// so without this file a fresh scaffold is invisible to it. The import must
// be the first line; anything Claude-specific goes below it.
const claudeScaffold = `@AGENTS.md

<!-- CodeCrew: Claude Code reads CLAUDE.md, not AGENTS.md; the import above
     loads the entry point every harness shares. Add Claude-specific
     instructions below it. -->
`

// extensionScaffold is the blank .codecrew/roles/<role>.local.md init writes beside
// every contract: the mechanism made visible at onboarding the way the
// routing table is, so "how do we customise this?" is answered where the
// question arises. Stable and brief — the invariant (SPEC §7) and two
// upstream links; the examples live on the page they point at and change
// without touching anyone's scaffold. Comments only, so it composes to
// nothing until someone writes something (M7-R4, #173).
const extensionScaffold = `<!--
.codecrew/roles/%[1]s.local.md — this project's extension to the %[1]s contract.
Loaded after .codecrew/roles/%[1]s.md, append-only, never a replacement: an extension
that contradicts its contract is a review finding. gh codecrew roles show %[1]s
prints the composition. Comments only, it adds nothing.

What belongs here, with worked examples (house style, repo conventions,
a platform's wake syntax, ids and tooling):
` + U + `/docs/extensions.md
Protocol: ` + U + `/SPEC.md (section 7)
-->
`

// scaffold writes the greenfield files into dir. Both modes get the pointer
// and the entry point — .codecrew/AGENTS.md and the two root files that
// point at it — because an agent is dispatched into a spoke exactly as into
// a hub; hub mode adds the roadmap, the contracts and their extensions.
// Existing files are never touched — they are reported as skipped.
func scaffold(dir, hub string, contracts fs.FS) (written, skipped []string, err error) {
	pointer := filepath.FromSlash(config.Pointer)
	rolesDir := filepath.FromSlash(config.RolesDir)
	files := map[string]string{
		pointer:                               fmt.Sprintf(hubConfigScaffold, protocolVersion),
		filepath.FromSlash(config.AgentsFile): agentsScaffold,
	}
	for _, name := range rootEntryPoints {
		files[name] = rootEntryPointScaffolds[name]
	}
	if hub != "self" {
		files[pointer] = fmt.Sprintf("codecrew: \"%s\" # protocol version (SPEC.md §5): a different major is refused; not the CLI release\nhub: %s\n", protocolVersion, hub)
	} else {
		files["ROADMAP.md"] = roadmapScaffold
		entries, err := fs.ReadDir(contracts, config.RolesDir)
		if err != nil {
			return nil, nil, err
		}
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), localSuffix) {
				continue // never scaffold one project's extensions into another
			}
			data, err := fs.ReadFile(contracts, config.RolesDir+"/"+e.Name())
			if err != nil {
				return nil, nil, err
			}
			// Provenance stamp: names the release these contracts shipped
			// with, so drift can be judged three-way later (the base is
			// fetchable from the upstream repo at this version).
			files[filepath.Join(rolesDir, e.Name())] = contractStamp(e.Name()) + string(data)
			role := strings.TrimSuffix(e.Name(), ".md")
			files[filepath.Join(rolesDir, role+localSuffix)] = fmt.Sprintf(extensionScaffold, role)
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

// sameDir compares two paths as directories, symlinks resolved.
func sameDir(a, b string) bool {
	ra, err1 := filepath.EvalSymlinks(a)
	rb, err2 := filepath.EvalSymlinks(b)
	if err1 != nil || err2 != nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	return ra == rb
}

// initCmd scaffolds a new CodeCrew repo: hub mode by default, spoke mode
// with --hub owner/repo (pointer and entry point only — the contracts and
// the roadmap live in the hub).
func initCmd(w io.Writer, args []string) error {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	hub := flags.String("hub", "self", "hub repo this spoke points at (default: this repo is the hub)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	// The pointer belongs at the repository root — config.Load walks
	// upward and would find a nested one first — so init refuses to
	// scaffold a subdirectory before writing anything.
	if root := repoRoot("."); root != "" {
		if cwd, err := os.Getwd(); err == nil && !sameDir(cwd, root) {
			return fmt.Errorf("run init at the repository root (%s), not in a subdirectory", root)
		}
	}
	// init scaffolds rather than loading a pointer, so it is exempt from
	// the pointer check — but not from the layout. A repo still on 1.x is
	// migrated, never given a second layout beside the first.
	if legacy := config.LegacyLayout("."); len(legacy) > 0 {
		dir, err := filepath.Abs(".")
		if err != nil {
			dir = "."
		}
		return refuseLegacyLayout(dir, legacy, "run gh codecrew migrate to move it, rather than scaffolding a second layout beside it")
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
	// A kept AGENTS.md or CLAUDE.md that does not reach the instructions is
	// the one skip that leaves the project disconnected: the instructions
	// are on disk and nothing arrives at them. Reporting the skip is not
	// enough — print the lines to paste. A kept file that already reaches
	// them (a rerun on what init wrote) is a plain skip and asks nothing.
	var stranded []string
	for _, f := range rootEntryPoints {
		if slices.Contains(skipped, f) && !reachesInstructions(".", f, nil) {
			stranded = append(stranded, f)
		}
	}
	if repoRoot(".") == "" {
		fmt.Fprintln(w, "\nnote: this directory is not a git repository — the protocol lives in GitHub.")
		fmt.Fprintln(w, "first: git init && git add -A && git commit -m \"chore: scaffold codecrew\" &&")
		fmt.Fprintln(w, "       gh repo create <owner>/<name> --private --source=. --push")
	} else {
		// The scaffold is the last commit before the protocol starts: made
		// here, from exactly the files written, never pushed (#172).
		commitScaffold(w, ".", written)
	}
	// Then the one thing the scaffold cannot write to disk: the protocol's
	// labels, in the repository GitHub knows. After the commit, never
	// before it — a GitHub failure here is a note, and the scaffold stands
	// whatever it says (M14-R5, #267).
	initLabels(w, ".")
	if *hub == "self" {
		fmt.Fprintf(w, "\nnext: every seat is routed to you (~ in %s); to hand one to a\n", config.Pointer)
		fmt.Fprintln(w, "colleague, a team or an App later, see "+U+"/docs/identities.md")
		fmt.Fprintln(w, "then `gh codecrew milestone new --title \"...\" --goal \"...\"`")
	} else {
		fmt.Fprintf(w, "\nnext: tasks for this spoke attach to milestones in %s\n", *hub)
	}
	// Last, so the one thing needing a human is the last thing on screen.
	if len(stranded) > 0 {
		entryPointAction(w, stranded)
	}
	return nil
}

// entryPointAction prints the one thing these verbs cannot do for the
// operator: the root entry points they KEPT that do not reach
// .codecrew/AGENTS.md, and the exact lines to paste into each —
// entryPointLines verbatim, so what is pasted is what the scaffold would
// have written. init and migrate both end with it, because instructions on
// disk that nothing arrives at are the same incomplete project either way,
// and both reach it in the same state: an absent root file is written from
// the scaffold by whichever verb found it missing, so the only file either
// can list here is one whose prose belongs to the project (#301). One
// condition, one message — the block takes no wording from its caller.
func entryPointAction(w io.Writer, kept []string) {
	fmt.Fprint(w, "\naction needed — a kept entry point does not reach CodeCrew's instructions.\n")
	fmt.Fprintf(w, "Kept: %s\n", strings.Join(kept, ", "))
	fmt.Fprintf(w, "Add these lines to each, so an agent dispatched here finds %s:\n\n", config.AgentsFile)
	fmt.Fprint(w, entryPointLines)
}
