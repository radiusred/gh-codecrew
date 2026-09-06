# CodeCrew, precisely

[codecrew.works](https://codecrew.works) is the introduction site; the
[repository's README](../README.md) is the technical entry point, with the
install line and the routing-table example. This page says what CodeCrew is,
what exists today, and where everything else lives — it is the map.

## Read in this order

Two things first, because the pages below can give the wrong impression.
The verbs are documented so they *can* be run by hand; in practice the coding
agent runs them and a human answers the gates. And onboarding is the three
commands under [Start now](../README.md#start-now) — install the extension,
`cd` into the repo, `gh codecrew init` — after which an agent dispatched into
the repo reads [.codecrew/AGENTS.md](../.codecrew/AGENTS.md) (the root
[AGENTS.md](../AGENTS.md) points at it) and its role contract and does the
rest; the quickstart is the long form of what happens next, not a
prerequisite.

1. [codecrew.works](https://codecrew.works) — the introduction site: what
   the problem is, how the protocol works, and what it has delivered.
2. **This page** — the three parts, what is shipped, the refusal codes.
3. [Your first milestone](first-milestone.md) — the quickstart: one human,
   one agent, a milestone from opening to audited close.
4. [Identities](identities.md) — running solo, staffing seats with humans,
   teams and App identities, dispatching a role session.
5. [Platform interop](platform-interop.md) — the ladder's last rung: how an
   orchestration platform hosts the whole crew, written from four cycles of
   doing it.
6. [SPEC.md](../SPEC.md) — the protocol itself: topology, state model,
   configuration, verbs, roles, gates.
7. [Founding decisions](founding-decisions.md) and the per-milestone
   records in [milestones/](milestones/) — the "why", as it was recorded.
8. [GSD vs. "just let the model orchestrate"](gsd-vs-frontier-orchestration.md)
   — the essay that started the project: one person's experience with GSD
   across several projects, kept as the motivation, not a verdict on GSD.

## Three parts

CodeCrew is a lightweight framework for agent-driven software delivery: the
auditability and reproducible discipline of heavyweight frameworks, without
the ceremony. Project state lives in the tools teams already use — GitHub
issues, PRs, and CI — and the only documents the framework produces are
per-milestone records of the decisions that shaped the system and why.

1. **A protocol** ([SPEC.md](../SPEC.md)) — conventions for representing
   milestones, tasks, plans, decisions, deviations, and gates in GitHub, and
   how agents and humans transact over them.
2. **Role contracts** ([.codecrew/roles/](../.codecrew/roles/)) —
   harness-neutral prompt files for the implementer, reviewer, qa, and
   doc-synthesizer roles, and for the coordinator that dispatches them,
   loadable by any agent (Claude Code, Codex, Gemini CLI, or an
   orchestrator's company).
   A project extends a contract without forking it in
   `.codecrew/roles/<role>.local.md` (SPEC §7).
3. **A CLI** — `codecrew`, a single static Go binary wrapping `gh`, providing
   the workflow verbs with gates enforced as code.

## What exists

**Shipped:** v1.2.0 of the `gh` extension, implementing protocol 1.0
(`version` prints both: `v1.2.0 (protocol 1.0)`; the pointer's protocol
major is checked, another major refuses; `gh` itself must be 2.50.0 or
later, or the CLI refuses `GH_TOO_OLD` before any verb runs). Verbs:
`init`, `status`, `milestone new/evidence/close`,
`task new/start/finish`, `checkpoint`,
`role`, `roles diff/show`, `identity new/token/webhook`, and `version` — all implemented,
with machine-readable refusals (`refused[CODE]: detail`, catalogued below)
when a gate blocks. `task start` is role-aware: roles whose contracts forbid
commits (qa, reviewer) get no linked development branch; `roles show <role>`
prints a contract with its `.codecrew/roles/<role>.local.md` extensions
appended; `task finish` deletes the branch it merged and `milestone close`
sweeps what its tasks left; `milestone new`, `task finish` and `milestone
close` take `--dry-run` — every gate in order with its outcome, then what the
verb would do, nothing written, the same refusal code. What changed and when:
[CHANGELOG.md](../CHANGELOG.md). Not yet here: any backend other than GitHub,
and GitHub Enterprise Server — github.com only.

**Who holds a seat.** Every role is always staffed, by exactly one of four
kinds of principal, named by the routing table's type prefix: the operator
themselves (`identity: ~`), a named human (`user:<login>`), a GitHub team
(`team:<org>/<slug>` — any member holds the role), or a GitHub App identity
(`app:<slug>`). Solo is a routing
configuration, not a degraded tier: the qa holder's verdicts are the ones
that count at close, and the reviewer holder's approving review is the one
`task finish` requires. A solo operator therefore needs nothing but `gh auth
login`; `gh codecrew init` scaffolds a project with every role routed to `~`.
When a project outgrows solo, `identity new <role>` mints a dedicated App
identity through the manifest flow and routes it for you — minted with
`--with-approval-permission`, a reviewer App's approvals count toward
GitHub's own required-review rules, which makes fully agent-gated merges
possible ([identities.md](identities.md)).

**The routing table.** Who holds which seat is one `roles:` table in the hub's
`.codecrew/config.yml`: a row per seat, naming the typed identity that holds
it and the harness and model it is dispatched under, with `~` where the
operator holds it. Two
worked examples are one click away — this repository's own table, as it stands
today, in the [README](../README.md#the-routing-table), and an annotated
generic one on [the home page](https://codecrew.works/#the-crew). SPEC §5 is
the field-by-field reference.

**How a repo joins.** Every repo in a CodeCrew project carries a
`.codecrew/config.yml` pointing at the hub — a spoke's is a two-line pointer
(`init --hub owner/repo`); this repo is its own hub (`hub: self`; SPEC §3
on choosing yours). The hub's config also routes the five roles — the four
crew seats and the coordinator, which unrouted is you — and `init` writes
that table for you. Every repo also carries the entry point — the
instructions in [.codecrew/AGENTS.md](../.codecrew/AGENTS.md), a root
[AGENTS.md](../AGENTS.md) pointing at them and a `CLAUDE.md` importing that
— so an agent dispatched into a spoke lands the same way it does in the hub. Agent identities are GitHub Apps; a seat's
first act is `export GH_TOKEN=$(gh codecrew identity token <slug>)`, which
mints a short-lived installation token from the platform's env bindings
or the local key and stub and discovers the installation from the App
itself — see [identities.md](identities.md) (SPEC §5 for the tiers).

**What it has done** — the receipts are on the [home
page](https://codecrew.works/#codecrew-works); the per-milestone records are in
[milestones/](milestones/). The largest environment it has run in — an
orchestration platform dispatching all five seats — has its own page:
[Platform interop](platform-interop.md).

## Install and use

```sh
gh extension install radiusred/gh-codecrew   # precompiled, all platforms
gh codecrew version        # confirm what you installed (gh never auto-updates extensions)
gh codecrew init           # scaffold a new project (see first-milestone.md)
gh codecrew migrate        # a repo still on the 1.x layout: move it to 2.0 (--dry-run first)
gh codecrew status         # open milestones, inferred task states, raised gates, notes
gh codecrew role reviewer  # who holds a role: app:<slug>, user:<login>, team:<org>/<slug>, or ~ (you)
gh codecrew help           # the full verb list

# or build from source (single static binary; requires gh on PATH):
go build -o gh-codecrew ./cmd/codecrew
```

## Refusal codes

A blocked gate exits non-zero with `refused[CODE]: detail`. The code is for
the agent; the detail is for the human. All forty-two, by the verb that
raises them (the source is the catalogue of record — `refuse("CODE"` in
`internal/cli/`):

**any verb that loads `.codecrew/config.yml`**

- `LAYOUT_LEGACY` — the repo is still on the protocol 1.x layout: a root
  `.codecrew.yml`, or a root `roles/` holding one of the five contracts,
  with no `.codecrew/config.yml`. The detail names what was found and
  `gh codecrew migrate`, the one-shot verb that moves it; nothing reads the
  old layout. `init` is exempt from the pointer check and raises this one
  too, rather than writing a second layout beside the first.
- `PROTOCOL_MISMATCH` — a protocol major differs from the one this binary
  implements (SPEC §5); a missing field proceeds with a note. A pointer
  ahead of the binary asks for an extension upgrade, one behind it is told
  the repo predates this protocol and is moved with `migrate`. The check is
  topology-wide: a spoke reads the hub's pointer to resolve roles and
  applies it there too, naming both sides — one project speaks one protocol
  major.
- `IDENTITY_UNTYPED` — a routing row's `identity` carries no type prefix,
  so it names no kind of principal; the detail names the row and the four
  forms (`~`, `app:<slug>`, `user:<login>`, `team:<org>/<slug>`).
- `SPOKE_ROUTING` — a spoke's pointer carries a `roles:` block. The hub
  carries the one routing table for the project (SPEC §5); a copy in a
  spoke either silently outranks it or goes stale, and the protocol will
  not pick a winner. The detail names the hub and the rows found.
- `HUB_UNREADABLE` — this repo is a spoke and the hub's
  `.codecrew/config.yml` could not be fetched or parsed, so no role can be
  resolved. Routing fails closed: before 2.0 the verb fell back to the
  spoke's own empty table, which resolves every seat to `~` — turning
  `task finish`'s holder-review gate into "any non-author approved" and
  `milestone close`'s verdict count into "anyone commented". A hub that
  reads fine and declares no table is a different thing entirely and is
  legitimately `~` everywhere. A hub still on the 1.x layout has no such
  file, and the detail names `gh codecrew migrate`.
- `GH_UNREACHABLE` — `gh` never reached GitHub: no route, no DNS, or no
  credentials at all. Never reported as a missing hub table, and never a
  bare `gh` error. `codecrew version`, `codecrew help`, and
  `roles show`/`roles diff` in a hub need no network, so they keep working;
  from a spoke `roles show` needs the hub and raises this.
- `GH_TOO_OLD` — the installed `gh` is older than 2.50.0, the floor
  `task finish` and the close's branch sweep need (`gh pr checks --json`);
  the detail names both versions. A `gh --version` banner the CLI cannot
  parse proceeds with a note.

**`migrate`**

- `BOTH_LAYOUTS` — the repo carries `.codecrew/config.yml` *and* a protocol
  1.x pointer or contracts; the detail names both, and migrate will not
  choose between them. Keep whichever the project uses, remove the other,
  and rerun.
- `FOREIGN_ROLES_DIR` — a root `roles/` that holds CodeCrew's own files
  also holds entries it does not recognise; the detail names both sets.
  Migrate moves the five role contracts and their `<role>.local.md`
  extensions and nothing else, so it stops rather than guessing about an
  eleventh name. Move or remove them, then rerun. A `roles/` with no
  CodeCrew file in it is a project's own and is never read.
- `MIGRATION_UNSUPPORTED` — the pointer's protocol major is not 1: below
  1.0 predates the conventions the move assumes, and above it is not a 1.x
  repo whatever the files beside it look like. The detail names the version
  read.
- `IDENTITY_UNRESOLVED` — a bare 1.0 identity could not be typed to exactly
  one GitHub principal: nothing answers to it, both a user and an App do, it
  is an organization, or GitHub could not be asked. Write the row as `~`,
  `app:<slug>`, `user:<login>` or `team:<org>/<slug>` by hand, then rerun.
  (A bare value that is already an App slug is found at `<slug>[bot]`, so
  the common 1.0 table types itself.)

**any verb that reads a milestone's `## Requirements`**

- `REQUIREMENT_ID_MISMATCH` — an ID declared under `## Requirements` is not
  the milestone's own: the grammar is `M<milestone>-R<k>` (SPEC §4), so
  M12-R3 under M13 is refused before any verdict is counted against it. The
  detail names every offending ID and the milestone read; the fix is a hand
  edit of the issue body. `milestone close` and `milestone evidence` refuse;
  `status` prints the same condition as a line and carries on, because it
  reports the board rather than gating it.

**`task new`**

- `NOT_FOUND` — no open milestone with that number in the hub, after the
  verb has also read the hub's newest issues and retried for a few seconds:
  the listing can lag a milestone created seconds earlier.

**`task start`**

- `NOT_A_TASK` — the issue is not labelled `cc:task`.
- `CLOSED` — the task is already closed.
- `NO_PLAN` — the Plan section is empty; plans come before work (SPEC §4).

**`task finish`** (in the order the gates are checked)

- `CLOSED` — the task is already closed.
- `GATED` — a `cc:needs-decision` gate is raised; a human resolves it and
  removes the label.
- `GATE_UNRECORDED` — a gate was raised and the label removed, but no
  `**Gate resolved:**` comment records the decision (SPEC §8). Both labels
  are read per paragraph, anywhere in a comment; only `**Gate resolved:**`
  answers, and it answers every gate still open before it, never one raised
  after it.
- `NOT_OWNER` — the task was started by another seat (the `**Started by**`
  record `task start` posts on every start, accepted only from the login
  it names; the assignee for tasks that predate it; the same seat is the
  same login or the same routed role — any member of a team-held role). The seat that
  started a task finishes it — dispatch it; hand it over by running
  `task start` as the new seat (the latest record wins — the path when the
  starter has left); or a human operator overrides on the record with
  `--bypass`. A task with no start record is not gated.
- `NO_PR` — no open PR closes the task.
- `NO_CHECKS_PERMISSION` — the installation token cannot read the PR's
  checks at all: on a private repo the App needs `checks: read` (the status
  check rollup) and `actions: read` (the workflow run behind each suite);
  the detail names the App and the missing one, and the fix is the App's
  settings page followed by the installation's acceptance
  (docs/identities.md).
- `NO_CHECKS` — the PR reports no CI checks at all; absence cannot satisfy
  the deterministic gate, and there is no override.
- `CHECKS_PENDING` — checks are still running.
- `CHECKS_FAILING` — a check failed.
- `NO_HOLDER_REVIEW` — the reviewer role is routed to someone, and that
  holder has not approved; the role defines whose review counts.
- `NO_NONDOER_APPROVAL` — the reviewer role is unrouted and no non-author
  has approved (solo: rerun with `--operator-confirm`).
- `REVIEW_NOT_COUNTED` — the protocol's review gate passed, but GitHub's
  own required-review rule is still unmet (an App's approval without write
  access); a non-author human approves on the reviewer's recommendation,
  the App gets write access, or `--bypass` where the ruleset allows it.
- `SELF_CONFIRM` — `--operator-confirm` was given by a crew identity (a
  `[bot]` login or an `app:`-typed seat holder); only a human operator can
  waive review, and a `user:`- or `team:`-typed holder is one.
- `CREW_BYPASS` — `--bypass` was given by a crew identity; a bypass is an
  operator's act.

**`milestone new`**

- `MILESTONE_NUMBER_TAKEN` — the issue was created, but another milestone
  already carries its `M<n>:` prefix (the listing the number came from
  lagged) and the verb's own renumbering failed or found the next number
  taken too; the detail names both issues and the fix: retitle the new one
  to the next free number and rewrite its `M<n>-R<k>` IDs. A repair that
  succeeds is not a refusal — it prints a `renumbered:` line.

**`milestone close`** (in the order the gates are checked)

- `NOT_FOUND` — no open milestone with that number.
- `MILESTONE_GATED` — the milestone issue itself carries `cc:needs-decision`:
  a requirement-level question raised there by `checkpoint`; a human
  resolves it with a `**Gate resolved:**` comment and removes the label.
- `OPEN_TASKS` — tasks are still open, listed with their inferred state.
- `NO_REQUIREMENTS` — the milestone's `## Requirements` section yields no
  bold IDs, so there is nothing to verdict; IDs written elsewhere in the
  body do not count (`new`, `status` and `evidence` note this first).
- `REQUIREMENT_ID_MISMATCH` — the section declares an ID belonging to
  another milestone (above); the same gate, checked next.
- `VERDICT_MISSING` — a requirement has no QA verdict from the qa holder.
  A verdict written inside a code span or a fenced block is content, not a
  verdict.
- `VERDICT_UNSATISFIED` — the latest verdict on a requirement is not
  `satisfied` — the latest *comment* carrying one for that ID, and the
  first verdict for the ID inside it.
- `DOC_MISSING` — `docs/milestones/<n>-*.md` is not on the default branch;
  the gathered records are printed above the refusal for the
  doc-synthesizer, which delivers the document as a task (plan, `task
  start`, a PR with `Closes #<task>`, `task finish`).

**`milestone evidence`**

- `NOT_FOUND` — no milestone with that number, open or closed. A closed
  milestone resolves and is reported as closed before the report: a record's
  citations are worth checking after the close as before it.
- `REQUIREMENT_ID_MISMATCH` — the milestone declares a requirement ID that
  is not its own (above); checked before the walk.
- `EVIDENCE_UNREACHABLE` — github.com links the milestone's record cites
  do not resolve; repair them before dispatching QA. A URL inside a code
  span or a fenced block is content, not a citation, and is not checked;
  an external citation that does not resolve is a `warning:` line, not a
  refusal.

**`identity token`**

- `NO_CREDENTIALS` — nothing to sign with: no App id and key bound in the
  environment (`GITHUB_APP_ID`/`GITHUB_CLIENT_ID`, `GITHUB_PRIVATE_KEY`/
  `GITHUB_PEM`), and no key and stub under `~/.config/codecrew/` for the
  slug; the detail says what was looked for and how to write the stub.
- `BAD_CREDENTIALS` — GitHub rejected the App JWT (401), or knows no App
  by the id it was signed as (404 "Integration not found"): the key and
  the id do not belong to the same App, or the key was revoked. Retrying
  will not help; check the id against `gh api /apps/<slug> --jq .id`.
- `NO_INSTALLATION` — the App is installed on no account the key can see;
  install it (identities.md, step 4).
- `INSTALLATION_AMBIGUOUS` — the App is installed on several accounts and
  neither a hint nor the hub's owner selects one; the detail lists them,
  and `--installation <id>` (or `GITHUB_INSTALLATION_ID`) chooses.

**`identity webhook`** (and `NO_CREDENTIALS`/`BAD_CREDENTIALS` as above)

- `NO_WEBHOOK` — the App was minted without a webhook; GitHub keeps no
  hook configuration for it and its API cannot create one. Activate the
  webhook on the App's settings page (the detail names it) with the
  receiver's URL, then `--secret` sets the receiver's secret.

Licensed under [Apache 2.0](../LICENSE).
