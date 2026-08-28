# Changelog

All notable changes to CodeCrew. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); the CLI follows
semantic versioning, and the protocol carries its own version (SPEC §5).

## [Unreleased]

### The verb takes requirements
- `milestone new --requirement TEXT` (repeatable) writes each requirement
  as a bold-ID line under `## Requirements`, numbered `M<n>-R1`, `R2`, …
  in the order given, and prints the IDs it counted. Both milestones the
  orchestrator opened had their IDs under Goal, because `--goal` was the
  only text input the verb offered — the shape #144's `NO_REQUIREMENTS`
  refuses at close (#147; #119 findings 19a, 28, 32). Text that brings
  its own ID is refused. A title whose `M<k>` prefix disagrees with the
  number the CLI derives is refused instead of doubled (numberguess #11,
  "M3: M2 — …", closed as a duplicate); one that agrees is stripped.

### The contracts say what runs
- Every command the contracts, the `AGENTS.md` scaffold and the CLI's own
  output name is written `gh codecrew …` — the installed form. The qa
  contract's first act was `codecrew milestone evidence`, and a dispatched
  qa agent ran it literally into `command not found`, twice (#146, #119
  finding 31). SPEC keeps the bare name for the protocol's verbs and says
  so. A test scans the embedded contracts, the scaffold, `init`'s next
  steps and the usage text for the bare form.
- The implementer contract carries the identity reflex the orchestrator
  run showed was missing: mint first, per session, as `GH_TOKEN` only; a
  401 means mint again, never escalate first; commit as the App's bot
  user; the env-var names platforms actually inject. Also: no task issue
  means stop and ask for one, and a decision that hands an obligation to
  another seat is not the doer's (#119, findings 6, 9, 10, 12–14, 22).
- The qa contract says why `contents: read` is deliberate; the reviewer
  contract says to post with the token on the command line and confirm
  the review's author; the doc-synthesizer delivers the milestone document
  as a task, and `DOC_MISSING`'s detail names that path instead of "merge
  its PR" (#119, findings 14 and 27).
- The dispatch recipe's identity check works under an installation token:
  the stub's App ID against `gh api /apps/<slug>`, never `GET /app`
  (#139). `identities.md` also tells platforms to give each agent its own
  `GH_CONFIG_DIR`, that a crew App's permission set cannot create a
  repository (no `Administration: write`), and that a protected default
  branch makes the scaffold the first PR — which `init`'s
  next step and the quickstart now say too (#119, findings 3, 4, 22).
- The prerequisites state the `gh` floor: 2.50.0, for `gh pr checks
  --json`, which `task finish` and the close's sweep read (#119, findings
  21 and 30; the in-CLI check is #149).

## [1.0.2] — 2026-08-28

### The close verifies something
- `milestone close` refuses `NO_REQUIREMENTS` when the milestone's
  `## Requirements` section yields no bold IDs — previously the verdict
  gate iterated zero requirements and passed vacuously (the orchestrator
  run closed a milestone over a "not satisfied" verdict this way, #119
  finding 28). `milestone new`, `status` and `milestone evidence` say so
  first; the template says IDs under Goal or Gates do not count. (#144)

## [1.0.1] — 2026-08-27

### Launch readiness
- `init` also writes `CLAUDE.md`, whose first line imports `AGENTS.md`:
  Claude Code loads `CLAUDE.md` and never `AGENTS.md`, so a fresh scaffold
  was invisible to the harness most adopters reach for — `gh codecrew init`
  then `claude` started blind. Found while assessing the README brief
  (#140). Idempotent like the rest of the scaffold; spokes are unchanged.
  (#141, PR #142)
- The README leads with the four-line start — install, `init`, start your
  agent, "Let's build this project!" — and says who does what: the agent
  runs the verbs, the operator answers the gates. (#140)

## [1.0.0] — 2026-08-27

The first stable release of the `gh` extension and of the protocol it
implements: SPEC 1.0. Everything below merged through the protocol — a
task, a plan, a PR, the reviewer seat's approval, `task finish` — and is
recorded on the linked issues.

### The protocol has a version, and the CLI checks it
- `codecrew version` prints both: `v1.0.0 (protocol 1.0)`. The pointer's
  `codecrew:` field is checked by every verb that reads it — another
  protocol major refuses `PROTOCOL_MISMATCH`; `"0.1"` (the pre-1.0 form of
  these same conventions) and a missing field proceed with a note. `init`
  scaffolds `codecrew: "1.0"`. Decided at a recorded gate. (#114, PR #127)
- SPEC is "Version 1.0" — no longer draft — and §10 states what 1.0
  promises (below).

### Contracts extend without forking
- `roles/<role>.local.md` — a project's own instructions for a role,
  loaded after the contract: hub contract, hub extension, spoke extension.
  `roles show <role>` prints the composition a dispatched session loads;
  `status`'s drift report never sees extensions. The first extension is
  this hub's editorial voice for outward-facing writing. (#122, PR #123)

### The record gathers completely
- `milestone close` gathers one record per labelled paragraph: qualified
  labels (`**Decision (…):**`) are captured verbatim, a record written
  after other text in the same comment is found, CRLF comments are
  handled, and each task's PRs are listed as summary pointers. Proven on
  the real M5 comment corpus, committed as a regression test: 3 records
  gathered before, 7 after. PR bodies are never gathered (SPEC §4).
  (#113, PR #126)

### A close leaves the repo clean
- `task finish` deletes the head branch it merged. `milestone close`
  sweeps the tasks' branches after every gate has passed — deleting only
  a branch whose PR merged and still sits at the merged commit, or one
  with no open PR and nothing beyond the default branch; never a fork's
  head or the default branch; reporting everything else. `status` notes
  when the repo does not delete branches on merge. (#129, PR #130)

### Launch readiness
- A scaffold that stands alone: every reference `init` writes resolves
  from an adopter's repo (upstream URLs; the token script installs with one
  line). `--help` exits 0 on every verb and is resolved before any verb
  runs. SPEC promises only what ships: no PR-body deviation check; github.com
  only (GitHub Enterprise Server is a non-goal for now). The quickstart
  names pull-request CI as a prerequisite. CONTRIBUTING, SECURITY and a
  teardown inventory. (#131, PR #135)

### Docs and branding
- A landing-page README in the house voice, with four illustrations
  (#111, PR #124); the former README is `docs/introduction.md` — the map,
  what exists, and the catalogue of all twenty-one refusal codes by verb
  (#112, PR #125); the CodeCrew logo and a mark adopters may use for their
  own crew Apps (#110, PR #121); the Copilot coexistence statement (#104,
  PR #105); the M5 milestone document (#106, PR #107).

### What 1.0 promises
Within a major release series of the CLI: verb names and their flags are
additive — nothing is renamed or removed; a refusal code's meaning is
stable — codes may be added in a minor, never repurposed, removed only in
a major; the `refused[CODE]: detail` line and the `version` output are
stable shapes, other human-facing text is not; pointer fields are
additive; the embedded role contracts may change in a minor, with the
drift report and `roles diff` as the mechanism. A protocol change that
invalidates existing pointers or recorded comments is a protocol major,
and the CLI that implements it refuses the old pointer.

[Unreleased]: https://github.com/radiusred/gh-codecrew/compare/v1.0.2...HEAD
[1.0.2]: https://github.com/radiusred/gh-codecrew/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/radiusred/gh-codecrew/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/radiusred/gh-codecrew/compare/v0.5.0...v1.0.0
