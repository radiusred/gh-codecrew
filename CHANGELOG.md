# Changelog

All notable changes to CodeCrew. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); the CLI follows
semantic versioning, and the protocol carries its own version (SPEC §5).

## [Unreleased]

### Removed
- The ten crew badge PNGs, `assets/codecrew-{code,coord,docs,review,test}.png`
  and their `-t` variants: nothing in the repo referenced them, the App
  avatars are set by hand and the shipped copies live in codecrew-www's
  `docs/assets/images/crew/`. The logo, the mark, the social preview and
  `assets/svg/` stay. (#214)

### The M9 record
- `docs/milestones/9-the-docs-at-codecrew-works.md` — the milestone document
  for "The docs at codecrew.works": the build-time sync that gave the hub's
  documentation a web home, the introduction becoming the section index rather
  than a second README, the operator's amendment that took the milestone
  records off the marketing site and the Decision it reversed, the reviewer
  pass a sandbox stopped from posting, and the qa `not satisfied` verdict that
  produced a fix task and a strict deploy. The ROADMAP row is Done and the
  README's milestone count is refreshed with it. Docs only. (#205)

### The M8 record
- `docs/milestones/8-a-product-home-page-for-codecrew-works.md` — the
  milestone document for "A product home page for codecrew.works": the
  per-page template and the copy that stayed Markdown, the measured
  42-character rule, the images wired in and taken out again, the CSS-only
  popovers and the three defects only a render found, the drawer repair and
  why `hidden` was the mechanism, and the twenty-two changes one task issue
  absorbed. The ROADMAP row is Done and the README's milestone count is
  refreshed with it. Docs only. (#203)

### The M7 record
- `docs/milestones/7-the-coordinator-seat-and-platform-interop.md` — the
  milestone document for "The coordinator seat and platform interop": the
  M7-R4 amendment and what it rejected, the coordinator seat and the four
  verbs the run asked for, cycle 4 on `radiusred/snake` with its fold-back
  map, the announcement gate, and the clause-level statement of what cycle 4
  did not exercise of the shipped coordinator contract. The ROADMAP row is
  Done and the README's milestone count is refreshed with it. Docs only.
  (#189)

## [1.1.0] — 2026-08-31

### The scaffold is a commit
- `init` commits exactly the files it wrote — `chore: scaffold codecrew`,
  a pathspec commit, so the operator's own staged and unstaged work is
  untouched (no stash) — on the current branch, or on
  `codecrew-bootstrap` cut from the default branch when that branch
  requires pull requests (asked through `gh`; assumed when it cannot be
  asked, since a commit stranded on a protected `main` is the worse
  outcome). It never pushes, never runs `git init`, refuses a
  subdirectory (the pointer belongs at the repository root — `init` used
  to mistake one for "not a repository"), leaves a detached HEAD
  uncommitted with the command to run, and a rerun that writes nothing
  commits nothing.
  In a fresh repository the scaffold is the root commit and the last
  commit before the protocol starts; behind a ruleset the scaffold PR
  remains the one merge the operator does by hand — the pre-milestone
  gate — and delete-on-merge cleans the branch (#164 findings 51, 68;
  #172, #183)

### The ladder's last rung: hosting a crew on an orchestration platform
- `docs/platform-interop.md` — the page SPEC §9 pointed at, written from the
  orchestrator run's sixty-eight findings (#119 cycles 1–3, #164 cycle 4) and
  nothing else: the separation of concerns (the platform keeps dispatch and
  discussion, CodeCrew owns the record and routing); the coordinator seat,
  its permission set, and why it is its own agent rather than the platform's
  lead; mapping agents to roles by the routing table, with `roles show` as
  the bundle and `roles/<role>.local.md` as the platform overlay; credential
  injection and the 401 reflex; the three wake kinds and the one-wake-path
  rule; an eleven-row onboarding checklist in setup order; the per-cycle cost
  tables reproduced as recorded; the Paperclip recipe as the worked example,
  ids as placeholders; and the seams still open, named as gaps. Linked as the
  last rung from the introduction, the quickstart's ladder, the README and
  SPEC §9. (M7-R7, #54, #182)

### An App's webhook signs for its platform
- `identity webhook <slug> [--show] [--url U] [--secret S | --rotate-secret]`
  works an App's hook under its own key: prints the URL, content type,
  whether a secret is set and the subscribed events; sets the receiver's
  URL and secret (nothing stored, nothing printed); rotates the secret
  and prints it once. Two things stay on the settings page, said rather
  than pretended: an App minted without a webhook has no hook
  configuration and GitHub's API cannot create one — `NO_WEBHOOK`, the
  twenty-ninth code, names the page where it is activated by hand — and
  event subscriptions are readable but not settable after creation.
- `identity new --with-webhook` now subscribes `pull_request` and
  `pull_request_review` — the transitions a platform routes to seats —
  instead of 1.0's five (`issues`, `issue_comment`, `check_suite` were
  wakes for nothing on a platform: #119, #164 findings 46, 53); `--events`
  names others, validated against the role's permissions, and
  `--webhook-secret S` sets the receiver's secret as soon as the App
  exists — before it is installed anywhere, which is the only way
  repository events reach it, so the creation ping (signed with GitHub's
  generated secret, rejected by the receiver, harmless) is the only
  delivery that precedes it. identities.md gains "The receiver side": one App hook covers
  every repository its installation sees — no repository hooks — the
  events per seat, what a receiver does, and the Paperclip routine as
  the worked example. (M7-R3, #157, #180)

### Dry runs
- `milestone new --dry-run` prints the number the milestone would get, its
  title, the requirement IDs it would number and the ROADMAP row, and
  creates nothing — requirement prose can be written knowing the number
  (a closed duplicate still counts toward it; #119 finding 45).
  `task finish --dry-run` and `milestone close --dry-run` evaluate exactly
  the gates the live verb would, in order, print each as ok, refused with
  its code, not reached or not applicable, then the actions a clean pass
  takes — the comments it would post, the merge and the head it would
  delete; every branch the sweep would delete or keep and why, and the
  closing comment — writing nothing and exiting with the first refusal's
  code. One code path builds the plan for both modes, so the preview
  cannot disagree with the run. (M7-R5, #133, #178)

### A seat finishes only its own task
- `task finish` refuses `NOT_OWNER` — the twenty-eighth code — when the
  caller is not the seat that started the task, read from the
  `**Started by**` record `task start` now posts on every start (accepted
  only from the login it names; the assignee is the fallback for tasks
  that predate it): the same login with the `[bot]` suffix ignored, or the same
  routed seat — a team-held role is any member. Handover is `task start`
  again by the new seat (latest record wins; the path when the starter has
  left). The operator's own auth is not
  exempt; `--bypass` is the recorded override, an operator's act as it
  already was (`CREW_BYPASS` for a crew identity), and the PR comment
  names the owner overridden. A task with no start record is not gated.
  The contracts say so; SPEC §6 and §8 list the gate. Cycle 4's
  implementer merged the doc-synthesizer's document because the
  coordinator's table named a fixed seat (#164 finding 58). (#165, #175)

### The scaffold asks what is local
- `init` writes a blank `roles/<role>.local.md` beside every contract —
  one comment saying what the file is (the project's extension, loaded
  after the contract, append-only, composed by `roles show`) with two
  upstream links: the new `docs/extensions.md` examples page and SPEC §7.
  The mechanism is made visible at onboarding the way the routing table
  is; the binary ships no opinion about what goes in it, and the examples
  change without touching anyone's scaffold. A comments-only extension
  composes to nothing, so `roles show` prints the bare contract until the
  project writes something. Rerunning `init` in an existing hub adds the
  blanks and nothing else. The teardown section lists them. Decided at
  the M7-R4 amendment (#163). (M7-R4, #159, #173)
- Three wordings from the orchestrator run (#158): the milestone's ROADMAP
  row rides in the milestone's first PR (the implementer's) and the
  document PR flips it; every seat contract says *landed means done* —
  hand back the way your platform wakes the coordinator, never park
  yourself until its next verb; `<milestone number>` replaces `<n>` in
  the qa and coordinator contracts, identities.md and the usage text.
- `milestone evidence` reads a URL as ending at the first character that
  cannot be in one — a prose ellipsis, quotes, a backtick, anything
  outside ASCII — trims trailing punctuation, and keeps a closing
  parenthesis only when the URL opened it (#138).

### The mint is a verb
- `identity token [<slug>] [--installation <id>]` mints an installation
  token in Go: the App id and private key from the environment under the
  names platforms bind (`GITHUB_APP_ID`/`GITHUB_CLIENT_ID`,
  `GITHUB_PRIVATE_KEY`/`GITHUB_PEM`, PEM text or a path), else the
  `~/.config/codecrew/` key and stub for the slug; the installation is
  discovered from the App itself — a hinted id is used only when the App
  can see it, a stale one is reported and overridden (#119 finding 35).
  The token alone on stdout, a receipt on stderr, nothing written to
  `gh`'s config. Four refusal codes, the twenty-fourth to twenty-seventh:
  `NO_CREDENTIALS`, `BAD_CREDENTIALS`, `NO_INSTALLATION`,
  `INSTALLATION_AMBIGUOUS`. The contracts name it as the first act and the
  401 recovery — no seat writes an RS256 helper of its own again (#119
  findings 2, 10, 12; #164 findings 56, 67) — and `scripts/codecrew-token`
  is a one-line wrapper around it. (M7-R2, #132, #170)

### The coordinator is a seat
- `roles/coordinator.md` ships in the binary and composes like the crew's
  four: `roles show coordinator` (with `roles/coordinator.local.md`),
  `roles diff`, the drift report, and `init`, whose scaffolded routing
  table now declares `coordinator: { identity: ~ }` as its fifth row.
  Unrouted the seat is the operator, as it always was; a 1.0 table without
  the row still answers `role coordinator` with `~`. `identity new
  coordinator` mints the App with the set #119 finding 16 specified —
  contents: read, issues: write, pull requests: read, metadata — never
  contents: write. The contract states what the orchestrator run taught
  the coordination layer (#119, #164): open milestones with
  `--requirement` and tasks with plans, dispatch by the routing table, own
  the review loop in both directions with the task's owner finishing it,
  raise gates with `checkpoint` (on the scaffold PR before a milestone
  exists), one wake path per transition with a per-seat table, re-read
  state at the act, execution events one-shot, dispatch on the platform
  and cite on GitHub, never the milestone number in requirement prose.
  SPEC §5, §7 and §9 and identities.md name the seat. (M7-R1, #168)

### The M6 record
- `docs/milestones/6-polish-and-1-0.md` — the milestone document for
  "Polish and 1.0": the protocol-version gate, the four releases and the
  CHANGELOG discipline, the branding requirement amended at a gate, the
  orchestrator run and its fold-backs, and the requirement-outcomes table.
  The ROADMAP row is Done; the introduction's shipped release and the
  README's milestone count are refreshed with it. Docs only. (#161)

### The orchestrator rung has been run
- The README, the quickstart's ladder and SPEC §9 say so, with the
  receipt: a Paperclip company drove three milestones on
  `radiusred/numberguess`, the third on the App's webhook events with one
  gate and no other operator touch on the workflow; the fifty findings and
  their fold-backs are the
  log on #119, the interop doc written from them is #54.
- SPEC §5 and identities.md name the credentials the way platforms bind
  them — `GITHUB_APP_ID`/`GITHUB_CLIENT_ID`, `GITHUB_PRIVATE_KEY`/
  `GITHUB_PEM` — and say the installation is discovered from the App, a
  supplied id being a hint at most (#119 findings 12 and 35). Docs only;
  the verb that does the resolving is #132.

## [1.0.3] — 2026-08-28

The orchestrator run's fold-backs (#119): four PRs, one release.

### The floor is a refusal, not a stack trace
- Every verb that reads `.codecrew.yml` checks the installed `gh` once and
  refuses `GH_TOO_OLD` below 2.50.0 — the release that added
  `gh pr checks --json`, which `task finish` and the close's branch sweep
  read. The crew's container carried a distribution-packaged 2.46: the
  gate failed inside `gh`, the sweep silently skipped, and the agent found
  the floor by the failure (#119 findings 21 and 30; #149). The
  twenty-third refusal code; an unparseable `gh --version` proceeds with
  a note.

### The qa seat earns its keep
- `roles/qa.md` asks for judgment, not a rerun: green tests on merged
  `main` are the floor; for each requirement the verdict says what the
  shipped suite proves and what it assumes (a gap is a finding even when
  the behaviour is right), and cites at least one probe the suite does not
  enumerate — order a hairbrush. A verdict with no findings says what was
  tried that failed to break it. From the operator's read of the
  orchestrator run's QA leg, which ran the implementer's cases with
  different hands and found nothing (#119 finding 37).

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

[Unreleased]: https://github.com/radiusred/gh-codecrew/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/radiusred/gh-codecrew/compare/v1.0.3...v1.1.0
[1.0.3]: https://github.com/radiusred/gh-codecrew/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/radiusred/gh-codecrew/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/radiusred/gh-codecrew/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/radiusred/gh-codecrew/compare/v0.5.0...v1.0.0
