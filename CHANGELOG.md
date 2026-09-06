# Changelog

All notable changes to CodeCrew. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); the CLI follows
semantic versioning, and the protocol carries its own version (SPEC §5).

## [Unreleased]

### The `cc:` labels are created, and restyled by `migrate`, from the crew palette
- **Nothing defined the protocol's labels.** `cc:milestone`, `cc:task` and
  `cc:needs-decision` were created implicitly by the first
  `gh issue create --label` or `POST /issues/N/labels` that mentioned
  them, so they wore whatever colour GitHub generated and carried no
  description at all — and the first gate in a repository was an untested
  path. `init` now ensures all three exist, hub and spoke alike, and
  `checkpoint` defines `cc:needs-decision` before applying it. An existing
  label is never touched, colour and description included: a project may
  have restyled one deliberately, and nothing can tell that from a GitHub
  default.
- **The colours come from the crew images.** `cc:milestone` takes the
  mark's cyan `#01d4ff`, `cc:task` the test seat's `#92edff` — the same
  hue, lightened, a task being part of a milestone — and
  `cc:needs-decision` the review seat's `#f0aeff`, because a gate is a
  question for a human and review is the seat whose job is asking one. The
  pairing is chosen for the two that co-occur on a gated task issue. SPEC
  §4 carries the table as the protocol's defaults; only the names are
  protocol.
- **`migrate` restyles; `init` does not.** The one-shot move to the 2.0
  layout now brings the three labels to the defaults whatever they were
  wearing, creating the missing ones and reporting each — the one place
  the protocol overwrites a label's styling. A 1.x repository's `cc:`
  labels were all created implicitly, so leaving them is leaving the
  migration half done, and the migration's promise is a repository
  indistinguishable from a fresh 2.0 `init`. `--dry-run` lists the
  creations and restyles beside the file steps and writes neither. The
  asymmetry is deliberate: `migrate` runs once, on a repository whose
  labels nobody chose; `init` reruns, on one whose labels somebody may
  have.
- **A GitHub failure never reaches the local work, and a rerun finishes
  it.** An unreadable label listing, a refused write, a `gh` that cannot
  name the repository, or a directory that is not a git repository yet is a
  `note:` line and nothing more — the files are written, the scaffold or
  the move is committed, and no verb ever refuses over a label.
  `checkpoint` degrades the same way: a creation that fails still raises
  the gate, since applying an unknown label creates it implicitly as
  before. What that leaves behind is recoverable rather than permanent:
  `migrate` on a repo already on 2.0 moves nothing and still runs the label
  step, and it runs on the two paths where the commit never happens — a
  detached HEAD, and a commit `git` refused — so "the files moved" and "the
  labels were done" are never two different answers. SPEC §6 says the
  *move* runs once and a rerun repeats the label step alone, and §4 accepts
  what follows: a `cc:` label restyled deliberately after the migration is
  set back to the defaults by the next `migrate`, previewed by `--dry-run`
  and named in the receipts, with `init` — which never touches an existing
  label — as the verb a project reruns instead.
  (#283, #267)

### A close sweeps the branches earlier closes left behind
- **`milestone close` no longer walks past a stale task branch.** The sweep
  visited only the closing milestone's own tasks, so a branch whose task
  shipped under a milestone that closed before the sweep worked — or whose
  delete failed once — was invisible to every later verb and stood forever:
  `radiusred/numberguess` still carries two of them from closes on 2026-08-28
  (#167). A second pass now follows the milestone's own: one prefix-filtered
  listing of each repo's `task/<n>-…` branches — the hub the milestone issue
  lives in, and every repo its tasks name — and every branch whose task issue
  is closed judged by the same two delete conditions the milestone's own meet
  (its PR merged and the tip still at the merged commit, or no open PR and
  nothing beyond the default branch), which is also the only test that
  catches a rebase-merged branch. What goes is named in the closing comment
  under its own sentence, `Swept from earlier closes: …`; a branch with
  unmerged commits, or one whose task is still open, is named and left with
  the reason. Bounded: one listing per repo, and an open task's branch costs
  a single issue read. `--dry-run` lists them beside the milestone's own and
  writes nothing, and the milestone's own branches are never revisited by the
  second pass. Two guards make the wider candidate set safe: a branch about
  to go because no PR is open is checked once more against the pull requests
  whose head it is — whatever they close, which is the relation a task's own
  closing PRs cannot see — and a candidate the verb could not read at all is
  reported and left standing, never deleted. A repo with more task branches
  than one listing holds is swept in part and says so. (#273)

### A task adopts a backlog capture, and the merge closes it
- **`task new --adopts <ref>[,<ref>]`** records the backlog issues a task
  takes up: repeatable and comma-separated, a bare number resolving against
  the task's own repo and `owner/repo#n` naming one anywhere. The refs go
  into an `## Adopts` section of the task body and each capture gets a
  comment naming the task. Every ref must be an open issue, checked before
  anything is created — `refused[ADOPT_NOT_OPEN]` for one that cannot be
  read as much as for one already closed — so a refusal leaves no
  half-adopted task; duplicates collapse, and a comment that fails once the
  task exists is a `note:`, the body carrying the link either way.
- **`task finish` closes them after the merge**, each with a comment naming
  the task, the pull request and the commit the merge left on the default
  branch. Nothing there refuses, because the merge has happened: a capture
  already closed is reported, and one that cannot be closed is a `note:`
  naming it. `--dry-run` lists each capture it would close, and each it
  would skip, beside the branch it would delete.
- **The section is read strictly**, since what it lists is closed after a
  merge where nothing can refuse: the heading counts only as a heading, so
  prose quoting `## Adopts` cannot shadow the real section (a task about
  this feature quotes it); only the ref at the head of a list line is an
  adoption, the prose after it being the capture's title; and the body is
  read through `StripCode`, so a ref inside a code span, a fenced block or
  an indented one is content, exactly as it is for the verdict scan and the
  citation walk (M13-R6).
- The protocol now does the bookkeeping the link already describes: an
  adopted capture stays open exactly as long as the task carrying it, and no
  PR body has to remember a `Closes` line for one. Every hub milestone since
  M3 closed with a manual sweep of the captures it had adopted, or forgot
  one. SPEC §4 says what adoption means, §6's `task new` and `task finish`
  rows say what the verbs do, and the coordinator and implementer contracts
  follow (#270, closing #193).

### `task finish` tidies the clone it ran in
- **The local task branch goes with the remote one.** `task finish` merged
  the PR and deleted the head branch on GitHub, and left the operator
  standing on a local branch whose upstream had just vanished — eighteen of
  them after one milestone close, all safely merged, none deletable with
  `git branch -d` because a rebase-merge rewrites the commits (#192). The
  verb now finishes the job in the clone it runs in: it fetches with
  `--prune`, switches off the task branch to the default branch,
  fast-forwards that branch to the merge, and deletes the local task
  branch, printing every step. The force-delete is allowed on two grounds
  only — the branch sits at the commit GitHub merged, or it is contained in
  the fetched default branch — so a branch with an unpushed commit is named
  and kept, and a local default branch that has diverged is named and left
  alone. Outside a repository, in a clone of another repo, or with no local
  branch of that name, nothing happens and nothing is printed; `--dry-run`
  names the local steps beside the remote gates without fetching or moving
  a ref. The CLI can only tidy the clone it runs in — a multi-clone setup
  still sweeps the others, and `docs/identities.md`'s cleanup guidance now
  tells the remaining remote refs (`git branch -r`) from the local branches
  left in each clone (`git branch`), the ones kept for their commits
  included (#271).

### The M13 record
- `docs/milestones/13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md`
  — the milestone document for "Protocol 2.0: the .codecrew/ layout and what
  rides with it": the layout move and the two fresh-context scans that decided
  what rode with it, the operator's Decision adopting seven findings into
  requirements and blessing four surfaces permanent, and the release struck out
  to M14 (#269) seven minutes after the milestone opened, so nothing shipped
  here. Then the eight requirements as delivered — `.codecrew/` everywhere and
  `LAYOUT_LEGACY`, `gh codecrew migrate`, the entry point standing on its own,
  typed identities, routing that fails closed, the tightened record grammar,
  the 1.0 shims removed with the exit-code contract and the forty-two-code
  table in SPEC, and the two captures that rode along — with the twelve change
  requests behind them, the approval a rebase dismissed, the reviewer seat's
  harness deviation, the two review bodies that reached GitHub as an
  unexpanded file reference, and the three QA rounds in which M13-R6 was
  verdicted `not satisfied` twice and remedied twice before it stood. The
  ROADMAP row is added Done, and the boundary refresh names what is on `main`
  and unreleased in both `README.md` and `docs/introduction.md`, while the
  introduction's shipped-release claim stays v1.2.0. Docs only. (#284)

### An indented code block opens wherever it would not interrupt a paragraph
- **A verdict quoted under a heading is content too.** `StripCode` opened an
  indented code block only after a blank line or at the start of the text,
  so `### Example` followed straight by a four-space verdict line still
  counted that line as a verdict and hid the real one below it. CommonMark
  restricts the form in one way only — indented code may not interrupt a
  paragraph — so a block now opens after a blank line, the start of the
  text, an ATX heading, a thematic break, or any line of a fenced block,
  and a prose line still prevents one, which keeps an indented paragraph
  continuation and a list item's continuation as prose. Headings and
  thematic breaks count as openers only below four columns of indentation,
  where CommonMark still reads them as such (#288).

### Indented code blocks are code, for verdicts and for citations
- **A verdict quoted in a four-space indented block no longer counts as a
  verdict.** `StripCode` blanked Markdown's inline spans and fenced blocks
  but not its third code form, so a QA comment that quoted an earlier
  verdict as an indented block — the shape a paste picks up when nobody
  reaches for backticks — superseded the verdict written below it. It now
  strips a run of lines indented four columns or more (a tab counting to
  the next multiple of four) that ends at the first non-blank line indented
  less than four; blank lines inside the run belong to it. One rule, one
  implementation, so the citation walk stops reading a URL in an indented
  block as evidence in the same commit. A line that continues a paragraph or
  a list item is not a block, because a paragraph is already open;
  indentation is measured from column 0, not from a list item's own content
  column (#285). Where such a block may open is the entry above (#288).

### The 1.0 shims deleted, and the machine contract written down
- **Breaking (protocol 2.0).** Three pieces of 1.0 scar tissue are gone,
  each of them a behaviour an adopter could have depended on, which is why
  they go at a major and not later (#254, the Claude scan's findings 1, 9
  and 10; the Codex scan's finding 4).
- **`codecrew: "0.1"` is no longer accepted.** A 1.x binary took the
  pre-1.0 form of the same conventions with a note to update the field. It
  is two majors back now, and its repo is on the 1.x layout, so it meets
  the same `gh codecrew migrate` refusal every older pointer meets. A
  pointer with no `codecrew:` field is still assumed current, with a note —
  kept deliberately, not by omission: under 2.0 the pointer lives at
  `.codecrew/config.yml`, so the file's own path is proof of the layout it
  speaks, and a 1.x repo has no such file to reach the check with.
- **A declared routing table names every seat.** `role coordinator`
  resolved to `~` when a declared table had no such row, because the row
  arrived after 1.0 hubs had scaffolded their tables. `init` scaffolds it
  and `migrate` writes it into a 1.x table, so the special case's only
  remaining effect was to infer a holder for a seat the table does not
  name; it now errors as any other missing role does. A table that declares
  nothing at all is unaffected — every seat, coordinator included, is the
  operator's.
- **An assignee is not a start record.** `StartedBy` fell back to a task's
  first assignee "for tasks started before the record existed", giving
  every assigned-but-never-started task an implicit owner for the ownership
  gate. Deleted — and with it the gate's old reading of an empty owner as
  "nobody to hold anyone to": `task finish` on a task with no
  `**Started by**` record now refuses `refused[NOT_OWNER]`, its detail
  naming `gh codecrew task start`. `--bypass` still overrides, recorded
  with wording that says there was no start rather than naming a starter
  that does not exist.
- **The exit-code contract, in SPEC §6.** Every failure exits 1 — a
  refused gate, a bad flag, an unreachable GitHub — and no exit-code
  taxonomy is coming within this major, because a status finer than
  "this did not happen" would break every caller already asserting on 1 the
  day it arrived. The machine channel is the `refused[CODE]: detail` line
  on stderr: the code is the branch point, the detail is prose for a human
  that nothing should parse, and a `note:` line is advisory. Written down
  because silence is what makes a later change breaking.
- **Every refusal code in one table, in SPEC §10.** All forty-two — the
  code, the verbs that raise it, one line of meaning — under the stability
  promise that already lived there. `docs/introduction.md` carried a second
  catalogue of the same codes and now points at the table instead: three
  tasks in this milestone alone had to edit both, which is the drift a
  single list exists to prevent. A test holds the table, the count SPEC
  states, and the `refuse("CODE"` sites in `internal/cli/` to the same set,
  so a verb cannot add a code without its row.
- SPEC §6's `task finish` row also gains `--operator-confirm`, which was in
  the usage string and load-bearing for the solo tier but missing from its
  signature. (#261)

### `gh codecrew migrate`: the one-shot move to the 2.0 layout
- The verb that moves a protocol 1.0 repo — hub or spoke — to the 2.0
  layout, once, in one local commit it never pushes. It reads no pointer,
  so it is exempt from the protocol check the way `init` is: the repos it
  exists for are exactly the ones every other verb now refuses. It runs at
  the repository root, and `--dry-run` prints every step and writes
  nothing. The steps an adopter runs are under *Protocol 2.0: the
  .codecrew/ layout* below.
- **What it moves.** `.codecrew.yml` to `.codecrew/config.yml`, and out of
  a root `roles/` only the ten names CodeCrew owns — the five role
  contracts and their `<role>.local.md` extensions — into
  `.codecrew/roles/`, by `git mv`, so the history follows the files; the
  emptied directory is removed. A root `roles/` is read at all only when it
  already holds one of the ten, so a project's own `roles/` (Ansible's, the
  collision the layout move exists to end) is never touched. Once the
  directory is CodeCrew's, an entry outside the ten refuses
  `refused[FOREIGN_ROLES_DIR]` naming it: migrate does not guess which
  files it owns.
- **What it writes.** `.codecrew/AGENTS.md`, the 2.0 entry point, when the
  repo has none — under 1.x those instructions lived in the root
  `AGENTS.md`, which belongs to the project, so a tree migrated without it
  would be on the 2.0 layout with no 2.0 entry point. The root `AGENTS.md`
  and `CLAUDE.md` are then treated exactly as `init` treats a kept one: one
  that already reaches `.codecrew/AGENTS.md` asks for nothing, one that
  does not is named under an `action needed` heading with the exact lines
  to paste, and neither is ever rewritten.
- **What it rewrites.** The pointer, in place: `codecrew: "2.0"`, a
  `coordinator` row when the table declares none, and every identity typed
  per the grammar above. A bare 1.0 value is typed by asking GitHub what it
  is — `users/<login>`, then `users/<login>[bot]`, since a 1.0 table wrote
  an App as its bare slug while the App's account carries the suffix — and
  a value that answers to nothing, to both a user and an App, or to an
  organization refuses `refused[IDENTITY_UNRESOLVED]` rather than encoding
  a guess in the routing table. The rewrite keeps the file's comments,
  blank lines and key order: the pointer is a file its project maintains.
- **What it refuses.** `BOTH_LAYOUTS` when the two layouts overlap — a 2.0
  pointer beside the 1.x one, or a 2.0 file already sitting where a 1.x one
  would move — naming what it found; `MIGRATION_UNSUPPORTED` when the
  pointer's protocol major is not 1, naming the version; `SPOKE_ROUTING`
  when a 1.x spoke's pointer carries a `roles:` block — 1.0 allowed the
  shape, 2.0 does not, and the rows are the operator's routing to move or
  delete rather than migrate's to drop; the two above; and
  `GH_UNREACHABLE` when the identity lookup could not reach GitHub at all,
  never folded into the value's own refusal. Every refusal is raised before
  anything is written, and the one rename that may skip `git mv` is a
  source git does not track. A repo already on
  2.0 says so, moves nothing, commits nothing and exits 0, so a rerun is
  safe — from #283 it still runs the label step there, so the rerun is
  also the recovery when that step could not reach GitHub.
- SPEC §6 carries the verb's row and §10 names it; `docs/introduction.md`
  gains the four codes (thirty-eight → forty-two, with the README's count).
  (#256)

### Routing fails closed
- **Breaking (protocol 2.0).** A spoke that cannot read its hub's
  `.codecrew/config.yml` now refuses `refused[HUB_UNREADABLE]`, naming the hub
  and the path, instead of degrading to its own empty table. An empty table
  resolves every seat to `~`, so the old fallback silently turned `task
  finish`'s holder-review gate into "any non-author approved" and `milestone
  close`'s verdict count into "anyone commented" — the two gates the protocol
  exists to enforce, failing open on a 404. The `.codecrew/` move guarantees
  that 404 for a whole migration window, in both directions of skew (#254,
  the Claude scan's finding 2). A hub that reads fine and declares no table is
  a different condition and is still legitimately `~` everywhere.
- The routing table is resolved once, at load, so no verb can consult an
  unresolved one: a hub reads its own pointer from disk — no fetch, so a hub
  resolves roles with the network down — and a spoke always reads the hub's.
- **Breaking (protocol 2.0).** The protocol-version check becomes
  topology-wide: a spoke reads the hub pointer's `codecrew:` major on that
  same fetch, and a major differing from the one the binary implements refuses
  `PROTOCOL_MISMATCH` naming both sides. One project speaks one protocol
  major.
- **Breaking (protocol 2.0).** A spoke's pointer carrying a `roles:` block is
  refused at load with `refused[SPOKE_ROUTING]`, naming the hub and the rows.
  The hub carries the one table; a copy in a spoke either silently outranks it
  or goes stale, and the protocol will not pick a winner (#254, finding 8).
- GitHub being unreachable — no route, no DNS, no credentials — is now its own
  `refused[GH_UNREACHABLE]`, never reported as a missing hub table and never a
  bare `gh` error. The classification was read off the installed `gh` rather
  than guessed, and is narrow on purpose: an HTTP 403 or 404 means GitHub
  answered. `codecrew version`, `codecrew help`, and `roles show`/`roles diff`
  in a hub read only local and embedded files and keep working with the
  network cut; from a spoke `roles show` needs the hub and raises the new code
  cleanly.
- SPEC §5 states that the hub carries the one table and §6 loses "A hub's
  routing table fetched from a spoke is advisory and is not checked", stating
  the checks in its place; the coordinator contract says what to do with the
  new codes; the introduction's refusal-code list gains three (thirty-five →
  thirty-eight, with the README's count). (#259)

### The entry point stands on its own
- **Breaking.** The CodeCrew instructions move out of the root `AGENTS.md`
  and into `.codecrew/AGENTS.md` — CodeCrew's file, in CodeCrew's directory,
  which a later `init` or `migrate` rewrites whole without touching a line
  the project wrote. The root `AGENTS.md` `init` scaffolds is now a two-part
  pointer at it: one sentence naming the path, for a harness that reads
  plain markdown, and a bare `@.codecrew/AGENTS.md` import, which Claude Code
  resolves transitively through `CLAUDE.md`'s `@AGENTS.md`. Neither form
  alone reaches every harness, so the scaffold carries both. `CLAUDE.md`
  imports `AGENTS.md` as before.
- **A spoke gets the entry point too** — `.codecrew/AGENTS.md`, the root
  pointer and `CLAUDE.md` — because an agent is dispatched into a spoke
  exactly as into a hub and needs the same file to land on. What stays
  hub-only is what the hub owns: the roadmap, the contracts and their
  extensions.
- **A kept `AGENTS.md` or `CLAUDE.md` is no longer only reported.** `init`
  keeps an existing root entry point untouched, as it keeps every existing
  file, and — when that file does not already reach `.codecrew/AGENTS.md` —
  prints the exact lines to paste into it, byte for byte the ones its own
  pointer carries, under an `action needed` heading naming each stranded
  file, last in its output. Instructions on disk that nothing arrives at
  were the one skip that left a project incomplete. A kept file that already
  arrives — the pointer `init` wrote on an earlier run — asks for nothing,
  so a rerun stays idempotent in what it says as well as in what it writes.
- This hub's own entry point follows the same shape, and SPEC §3, §6, §7 and
  §10, the README, `CONTRIBUTING.md`, `docs/first-milestone.md`,
  `docs/introduction.md`, `docs/identities.md` and `docs/extensions.md`
  follow. (#257)


### Protocol 2.0: the .codecrew/ layout
- **Breaking.** Every CodeCrew-owned operational file moves under
  `.codecrew/`: the pointer from `.codecrew.yml` to `.codecrew/config.yml`,
  the role contracts and their local extensions from `roles/` to
  `.codecrew/roles/`, in hub and spoke alike. The framework had been writing
  into the root of repositories it does not own, where `roles/` collides with
  real project layouts (Ansible's, for one). `init` writes the new tree,
  `config.Load` walks upward for the new pointer — reporting the directory
  that *contains* `.codecrew/`, so every path is still resolved from the repo
  root — and drift, `roles diff`, `roles show`, the spoke's hub fetch and the
  contract provenance stamp all read and name it. The human-facing record
  does not move: `ROADMAP.md`, `docs/milestones/`, `AGENTS.md` and
  `CLAUDE.md` stay at the root. The protocol version is `2.0`, and `init`
  writes it.
- **A 2.0 binary refuses a 1.x repo; it does not read one.** A root
  `.codecrew.yml`, or a root `roles/` holding one of the five contracts, with
  no `.codecrew/config.yml` above it, refuses `refused[LAYOUT_LEGACY]` naming
  what was found and `gh codecrew migrate` — the one-shot verb that moves a
  repo forward, arriving in this release. There is no dual-read and no
  compatibility shim anywhere. `init` reads no pointer, so it is exempt from
  the protocol check, but it raises the same refusal rather than writing a
  second layout beside the first.
- `PROTOCOL_MISMATCH` stops being symmetric: a pointer ahead of the binary
  asks for an extension upgrade, one behind it is told the repo predates this
  protocol and is moved with `migrate`. Neither wording asks anyone to edit
  the version field by hand — it describes the repo rather than choosing for
  it.
- This hub's own files moved in the same change, the embedded contracts being
  built from them. SPEC §3 (paths, and the blessing of the single `hub:`
  field: a spoke belongs to one hub, and a task created in a spoke by another
  hub's milestone resolves its hub through its parent milestone, not the
  pointer — #177), §5, §6, §7 and §10 follow, as do the five contracts,
  `AGENTS.md`, the README and the docs; the introduction's refusal-code list
  gains `LAYOUT_LEGACY` (thirty-four → thirty-five, with the README's count).
  (#255)
- **Moving a repo forward.** Upgrade the extension, then migrate each repo —
  hub first, since its `.codecrew/roles/` is what spokes read:

  ```
  gh extension upgrade codecrew     # or: gh extension install radiusred/gh-codecrew
  gh codecrew migrate --dry-run     # every step, nothing written
  gh codecrew migrate               # the move, in one local commit
  git show                          # read it: the pointer's rewrite is in there
  git push -u origin HEAD           # migrate never pushes; open the PR yourself
  ```

  `migrate` ends with an `action needed` block whenever the repo's root
  `AGENTS.md` or `CLAUDE.md` does not yet reach `.codecrew/AGENTS.md` — a
  1.x root entry point holds the old instructions, so it usually does not.
  Paste the two lines it prints into each file it names; migrate does not
  edit them itself.

  The move also brings the repo's `cc:` labels to the protocol's defaults —
  created where missing, restyled where a 1.x repo had them from implicit
  creation — reported after the commit and listed by `--dry-run`. It needs
  GitHub, and a repo it cannot reach gets a `note:` and the migration
  stands. Rerunning `migrate` is then the recovery: a repo already on 2.0
  moves nothing and still does the labels, printing `labels already at the
  protocol defaults` when there is nothing left to do (#283).

  `migrate` is described in its own entry above. (#256)

### Typed identities in the routing table
- A routing row's `identity` now names the kind of GitHub principal that
  holds the seat: `~` (the operator, and any session acting under the
  operator's own auth), `app:<slug>` (a GitHub App), `user:<login>` (one
  named human) or `team:<org>/<slug>` (any member of the team). **Breaking
  (protocol 2.0):** a bare value — a 1.0 table's `my-org-coder` — is
  refused at load with `refused[IDENTITY_UNTYPED]`, naming the row and the
  four forms; `gh codecrew migrate` rewrites an existing table (#256).
  An App slug and a username are the same string, so a 1.0 table could not
  say which a seat held, and the protocol treats them differently (#254,
  M13-R4).
- What the kind decides, instead of being guessed from the value's shape:
  a crew identity — refused `--operator-confirm` and `--bypass` — is now
  the `app:`-typed holders and `[bot]` logins only, so a human holding a
  seat (`user:`, or a member of a `team:`) keeps the operator's acts where
  before every routed login was refused; role holding matches an App with
  or without its `[bot]` suffix and a user exactly; a team is a `team:`
  row rather than any value containing a slash.
- `gh codecrew role <name>` prints the typed value (`app:radiusred-checky`),
  and the new `--login` prints the handle a review request can name — the
  login for `user:`, `<org>/<slug>` for `team:` — and nothing at all for an
  App or the operator, which is the implementer contract's whole branch at
  PR creation. `identity new` routes and prints `identity: app:<slug>`.
- SPEC §5 carries the grammar as a table and §6's every-verb row the new
  code; `docs/identities.md`, `docs/introduction.md`,
  `docs/first-milestone.md`, `docs/platform-interop.md`, the README's
  worked example, this hub's own routing table and the five role contracts
  follow. `init`'s scaffolded routing table teaches the four forms too,
  so a fresh hub is never told to write a value its own binary refuses.
  (#258)

### The record grammar, tightened
- Four rules about what recorded text *means*, each a protocol-2.0 break:
  text already written on GitHub is reclassified, which is why they ride a
  major.
- **A gate is read per paragraph.** `**Gate raised:**` counts anywhere in a
  comment, exactly as `**Decision:**` does. A gate raised as a comment's
  second paragraph used to be invisible to `task finish`, so
  `GATE_UNRECORDED` never fired for it.
- **Only a `**Gate resolved:**` record resolves a gate, and per gate.** A
  resolution answers every gate raised before it and still open — one
  comment may still answer several questions — and never one raised after
  it; a bare `**Decision:**` on any subject no longer clears every gate on
  the issue, which is what SPEC §8 always said. The `cc:needs-decision`
  label stays the hard block.
- **Verdict supersession is per comment.** The latest comment carrying a
  verdict for a requirement ID wins, and the first verdict for that ID
  inside it counts, so a QA comment may quote the verdict it supersedes
  without superseding itself. Code spans and fenced blocks are stripped
  before the scan, so the record-reading rule is now the citation-reading
  rule — `stripCode` moved from `internal/cli` to the tracker package as
  `StripCode`, one implementation serving both.
- **A requirement ID under `## Requirements` carries its milestone's own
  number.** `milestone close` and `milestone evidence` refuse the new
  `refused[REQUIREMENT_ID_MISMATCH]`, naming every offending ID, before a
  verdict can be counted against a requirement belonging to another
  milestone; `status` prints the same condition as a line rather than
  dying, because it reports the board rather than gating it.
- SPEC §4, §6 and §8, the qa and implementer contracts, and the
  introduction's refusal-code list (thirty-two → thirty-three, with the
  README's count) follow. (M13-R6, #260)

### Two captures from the field: status between milestones, evidence after a close
- `status` no longer stops at `no open milestones in <hub>`: that line replaces
  the board and the gates section, and the two advisory checks below it still
  run — the delete-branch-on-merge note and the contract-drift report. Drift is
  purely local (the hub's `.codecrew/roles/` against the contracts embedded
  in the binary) and has nothing to do with milestone state, and the quiet
  period between milestones is exactly when an operator reconciles a fork
  against a
  new release; a hub in that state had shown no drift line while
  `roles diff` showed the divergence (#253).
- `milestone evidence <n>` resolves a closed milestone too, reading the hub's
  milestone listing regardless of state under the same `M<n>:` title rule, and
  prints a note naming it as closed before the citation report. A reviewer
  checking M11's shipped record had been refused `NOT_FOUND`; link rot in a
  finished record is what the verb is most useful for. `refused[NOT_FOUND]` now
  means no milestone carries that number in either state, and `milestone close`
  and `status` keep their open-only reads (#250).
- SPEC §6's `status` and `milestone evidence` rows and the introduction's
  refusal-code list follow. (#262)

### The M12 record
- `docs/milestones/12-v1-2-0-and-the-field-fixes-behind-it.md` — the
  milestone document for "v1.2.0 and the field fixes behind it": the changelog
  flip and the introduction at the boundary, the annotated tag cut by the
  implementer identity fourteen seconds after the flip merged, the release run
  and its five assets, and the operator's verification of the installed binary;
  then the three fixes that merged behind the tag and ride the next release —
  `task new` reading the newest issues and retrying before `NOT_FOUND`,
  `milestone close` refusing `MILESTONE_GATED` while the milestone issue
  carries a gate, and `milestone evidence` checking citations rather than every
  URL — with the four review rounds on the close help line, the round-one
  verdict relayed after a machine crash and recorded as the milestone's one
  Deviation, the three captures adopted and closed and the one filed. The
  ROADMAP row is added Done, and the boundary refresh of `README.md` and
  `docs/introduction.md` found no stale claim to fix. Docs only. (#251)

### milestone close refuses while the milestone issue carries a gate
- `milestone close` reads the milestone issue's own labels and refuses
  `MILESTONE_GATED` while `cc:needs-decision` is on it — a requirement-level
  question raised there by `checkpoint` (#200) could until now stay open
  while the milestone closed under it (#219). The gate is the close's
  second, after "milestone open" and before "tasks closed", named "no gate
  raised" as `task finish`'s is; `--dry-run` prints the gate line and marks
  the rest not reached. `checkpoint`'s milestone-issue comment and receipt
  say the close refuses, beside `status` listing the gate. SPEC's
  `milestone close` row and §8, the introduction's refusal-code list
  (thirty-one → thirty-two, with the README's count) and the coordinator
  contract follow. (#244)

### A link is a citation; code is content
- `milestone evidence` checks the links a record cites — a URL in prose or
  in a Markdown link — and not every URL in it: a URL inside an inline code
  span or a fenced code block is content, a probe target that is unreachable
  by design or a hostname in a verbatim command or error string, and is not
  scanned. Two NXDOMAIN-by-design hosts quoted as curl targets in a survey
  comment had refused a complete record on radiusred/ops. An unreachable
  citation now splits two ways: a github.com link that does not resolve is
  still `EVIDENCE_UNREACHABLE`; any other host prints a `warning:` line and
  the verb passes, so a dead external link is visible to the qa seat without
  blocking it. The SPEC row, the qa and doc-synthesizer contracts and the
  introduction state the rule — write a probe target in code, write evidence
  as a link. (M12-R4, #222, #245)

### task new waits for a milestone the listing has not caught up with
- `task new --milestone <n>` no longer trusts a single read of the
  open-milestone listing. When the listing lacks `M<n>` it reads the hub's
  newest issues regardless of label and accepts an open `M<n>:` issue carrying
  `cc:milestone` there; failing that it waits two seconds and reads both
  again, three reads in all, before `refused[NOT_FOUND]` — whose detail now
  says the listing can lag a milestone created seconds ago. A milestone found
  by either fallback is noted in the output. Three `task new` calls straight
  after `milestone new` were refused in the field and the fourth, seconds
  later, succeeded: the label-filtered listing lags a just-created issue, the
  same window #209 closed for `milestone new`'s own number check. The
  coordinator contract says a `task new` right after `milestone new` may take
  a moment and retries itself; SPEC's `task new` row and the introduction's
  refusal list describe the resolution. Adopts #234. (#243)

## [1.2.0] — 2026-09-05

### The M11 record
- `docs/milestones/11-housekeeping.md` — the milestone document for
  "Housekeeping": the three backlog captures adopted and closed, the README cut
  from 226 lines to 133 with the receipts moved to the home page and the
  reference links kept at source in this repository, the routing-table anchor
  renamed and the spoke's sync table and home-page gloss repaired behind it, the
  introduction linking the example rather than carrying a third copy, and both
  repositories' commit lint moved onto the organisation's shared action with the
  `Lint commit messages` check context unchanged — with the operator's tone pass
  that amended M11-R1 mid-review, and the reviewer finding the introduction still
  narrating the README it no longer describes. The ROADMAP row is added Done, and
  the boundary refresh of `README.md` and `docs/introduction.md` found no stale
  claim to fix. Docs only. (#239)

### The README is the technical entry point
- `README.md` goes from 226 lines to 133, dry and technical, written for
  whoever wants more detail than the site gives — someone running a coding
  harness, or the agent reading on their behalf: what CodeCrew is and what it
  depends on in a paragraph, one line near the top naming
  [codecrew.works](https://codecrew.works) as the marketing and introduction
  site and this page as the technical entry point, the routing-table example
  with its gloss, the install line and the first verbs with the refusal codes
  they raise, and a Read next that is a plain list of the reference
  documentation **at source in this repository** — `docs/*.md`, `SPEC.md`,
  `CONTRIBUTING.md`, `SECURITY.md`, `ROADMAP.md`, `CHANGELOG.md`, and
  `AGENTS.md` with `roles/` for an agent dispatched into the repo. "Why you'd
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

[Unreleased]: https://github.com/radiusred/gh-codecrew/compare/v1.2.0...HEAD
[1.2.0]: https://github.com/radiusred/gh-codecrew/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/radiusred/gh-codecrew/compare/v1.0.3...v1.1.0
[1.0.3]: https://github.com/radiusred/gh-codecrew/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/radiusred/gh-codecrew/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/radiusred/gh-codecrew/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/radiusred/gh-codecrew/compare/v0.5.0...v1.0.0
