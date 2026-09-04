# Changelog

All notable changes to CodeCrew. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); the CLI follows
semantic versioning, and the protocol carries its own version (SPEC §5).

## [Unreleased]

### The README is the developer landing page
- `README.md` goes from 226 lines to 128, dry and technical: what CodeCrew is
  and what it depends on in a paragraph, one line near the top naming
  [codecrew.works](https://codecrew.works) as the marketing and introduction
  site and this page as the developer landing page, the routing-table example
  with its gloss, the install line and the first verbs with the refusal codes
  they raise, and a Read next that is a plain list of the reference
  documentation **at source in this repository** — `docs/*.md`, `SPEC.md`,
  `CONTRIBUTING.md`, `SECURITY.md`, `ROADMAP.md`, `CHANGELOG.md`. "Why you'd
  want a crew", the four beats, the ladder and "The receipts" are gone: the
  home page (M8) and the docs section (M9) carry that argument, and the README
  no longer makes it. Headings: the routing-table section is renamed, so its
  anchor moves from `#2-four-seats-always-staffed` to `#the-routing-table`;
  `#start-now` and `#read-next` are unchanged; the YAML block is byte-for-byte
  what it was. `docs/introduction.md` gains a one-paragraph gloss on what the
  routing table is, with links to the README's example and the home page's, and
  its `README.md#the-receipts` link now points at the home page's receipts
  section. Docs only. Under M11-R1 (amended by the operator on #233,
  2026-09-04) and M11-R2. (#235)

### The commit lint is the org's shared action
- `.github/workflows/commitlint.yml`'s `commitlint` job calls
  `radiusred/.github/.github/actions/commitlint@main` — the composite action
  radiusred/.github publishes, carrying the organisation's commitlint config
  and its pinned `wagoid/commitlint-github-action` — after an
  `actions/checkout@v4` with `fetch-depth: 0`, the thin caller its
  CONTRIBUTING.md documents; the workflow grants `pull-requests: read`
  alongside `contents: read` as that caller does. The job id and its name,
  `Lint commit messages`, are unchanged, so the check context the org ruleset
  `require-lint` requires still matches; the `test` job is untouched.
  `commitlint.config.mjs` is deleted — the config lives inside the action.
  Under M11-R3 (#233); adopts #225 (hub half). (#236)

### The M10 record
- `docs/milestones/10-protocol-bookkeeping-from-the-field.md` — the milestone
  document for "Protocol bookkeeping from the field": the seven backlog captures
  adopted and closed, the ROADMAP row moving to the doc-synthesizer at both ends,
  `milestone new` repairing a number collision instead of refusing on it, the
  refusal that names the missing App permission, the dispatch guidance verified
  against the installed Codex CLI, the codex seats re-pinned mid-milestone, and
  the README routing table shipped as a mirror and amended into a worked example
  — with the home page's placeholder-identity variant that diverges from it by
  the operator's decision. The ROADMAP row is added Done — the first record PR to
  add one rather than flip one — and the README's milestone count is refreshed
  with it. Docs only. (#231)

### The README's routing table is an example, not a mirror
- The block under "Four seats, always staffed" stays, and its gloss now says
  what it is: this hub's table *as it stands today*, showing what a routing
  table can do — a different harness and model per seat, a human on the
  coordinator row — with the note that another project's will look different.
  `routing_table_test.go`, which failed the build unless the block matched
  `.codecrew.yml` byte for byte, is deleted; the README documents what is
  possible, not what this repo happens to run. Amends M10-R6 (#207) and
  supersedes the second decision recorded on #213. (#229)

### The codex seats are pinned to gpt-5.5
- `.codecrew.yml` pins `model: gpt-5.5` on the `reviewer` and `qa` rows,
  the model the two codex-harnessed seats run under since the operator's
  Decision on #207 (2026-09-04) moved them off `gpt-5.6-sol` for cost; the
  README's mirrored block, and the routing-table examples in SPEC §5 and
  docs/platform-interop.md, say the same. The dispatch guidance's
  reasoning-effort note now gives gpt-5.5's own default (`medium`) and the
  levels it accepts (`low`, `medium`, `high`, `xhigh`), read from the
  Codex CLI's model catalog. Under M10-R5. (#226)

### The README shows the hub's routing table
- The README's "Four seats, always staffed" beat now carries this repo's own
  `roles:` table from `.codecrew.yml`, with a gloss: each row is a seat, the
  identity that holds it, and the harness and model it is dispatched under;
  `~` is a human. GitHub cannot transclude, so the block is a copy — a test
  in the root package fails the build if it is not the file's roles section
  verbatim. Hub half of codecrew-www#18. (#213)

### task finish names the missing permission
- `task finish` refuses `NO_CHECKS_PERMISSION` — naming the App, the PR and
  the permission the installation token lacks — when GitHub answers the
  status check rollup with `Resource not accessible by integration`, instead
  of dying on the raw GraphQL error before the gate list; `--dry-run` prints
  it as the CI checks gate's line. A private repo needs `checks: read` for
  the rollup and `actions: read` for the workflow run behind each check
  suite, and the refusal names whichever is missing; an unrelated failure
  still surfaces raw. `identity new`'s permission table grants
  `actions: read` to the four seats that read checks, and docs/identities.md
  says how an App minted before the table carried either permission gains
  it — on the App's settings page, then accepted on the installation, neither
  through the API. Adopts #198. (#211)


### milestone new never reuses a milestone number
- The number is derived from the max over two listings — the label-filtered
  milestone listing and the hub's newest issues regardless of label — and
  verified after the issue is created: both listings are read again and,
  when another issue already carries the `M<n>:` prefix, the new issue is
  renumbered to the next free number, title and every `M<n>-R<k>` ID, with
  a `renumbered:` line saying so (bounded to three rounds). A repair that
  fails, or a number still taken after that, is `refused[MILESTONE_NUMBER_TAKEN]`
  naming both issues and the hand fix. Three calls seconds apart came back
  as M2 in the field because the label-filtered listing lagged the issue
  the first call had just created. The tracker seam gains `MilestoneIssues`
  (replacing `AllMilestoneTitles`), `RecentIssues` and `EditIssue`; SPEC's
  CLI table, the refusal-code list and the coordinator contract say
  milestones can be opened back to back. Adopts #195. (#209)
### Gates on milestone issues show on the board
- `status` reads the milestone issue's own labels: a `cc:needs-decision`
  raised there — a question about a requirement that no task carries — is
  marked on the milestone's line and listed under `gates raised:` beside
  the task gates as `<ref> — <title> (milestone)`, where it was hidden
  before (#200). `checkpoint` keeps accepting a milestone ref, and its
  comment and receipt now say that `status` lists the gate rather than
  that `task finish` refuses; a task keeps the task wording. The
  coordinator contract says a requirement-level gate goes on the milestone
  issue, and SPEC's `status` and `checkpoint` rows and §8 say the same.
  `status` and `checkpoint` move onto the shared context so a fake tracker
  can drive them, with tests for both. (#210)
### Dispatch guidance covers network reach, and the codex seats pin their model
- `docs/identities.md`'s "Dispatching a role session" gains a **Reachability**
  bullet: the dispatched session must reach `api.github.com`, which every
  credential step assumes; a sandboxed harness may deny network by default and
  the seat then does real work and can post none of it (#202, the reviewer pass
  on codecrew-www#7). Codex CLI is the worked example — `--sandbox
  workspace-write` denies network unless `-c sandbox_workspace_write.network_access=true`
  is passed — and the bullet says how a codex seat's model and reasoning effort
  are set at dispatch (`-m`, `-c model_reasoning_effort=<level>`), both verified
  against codex-cli 0.152.1. The hub's `.codecrew.yml` pins
  `model: gpt-5.6-sol` on the `reviewer` and `qa` rows as it already did for
  the implementer, and the routing-table examples in SPEC §5 and
  `docs/platform-interop.md` do the same; the config test asserts a `model` on
  a codex row loads. Docs, config and a test fixture only. (#212)

### The ROADMAP row belongs to the doc-synthesizer at both ends
- `milestone new` creates the tracking issue and nothing else: the local
  append to `ROADMAP.md`, the "rides in this milestone's first PR" line and
  `--dry-run`'s row are gone. The row had no PR to ride in when a milestone's
  tasks all lived in spokes — hit three times in the field — so the
  doc-synthesizer now adds it, already Done, in the record PR; the roadmap
  lists finished milestones and `status` reports the open one. The
  implementer contract drops "The ROADMAP row is yours", the doc-synthesizer's
  "Flip the ROADMAP row" becomes "Add the ROADMAP row", and SPEC §4, the CLI
  table, the first-milestone guide and the interop page say so; the README
  and CONTRIBUTING stop claiming the roadmap names the open milestone.
  Adopts #197, shape 3. (#208)

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
