# CodeCrew, precisely

The [landing page](../README.md) says why you would want a crew. This page
says what CodeCrew is, what exists today, and where everything else lives —
it is the map.

## Read in this order

Two things first, because the pages below can give the wrong impression.
The verbs are documented so you *can* run them; in practice your coding
agent runs them and you answer the gates. And the whole of onboarding is
the four lines on the landing page — install, `init`, start your agent,
"Let's build this project!" — the quickstart is the long form of what
happens next, not a prerequisite.

1. [Why you'd want a crew](../README.md) — the landing page: the problem,
   the four beats, the receipts.
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
2. **Role contracts** ([roles/](../roles/)) — harness-neutral prompt files
   for the implementer, reviewer, qa, and doc-synthesizer roles, and for the
   coordinator that dispatches them, loadable by any agent (Claude Code,
   Codex, Gemini CLI, or an orchestrator's company).
   A project extends a contract without forking it in
   `roles/<role>.local.md` (SPEC §7).
3. **A CLI** — `codecrew`, a single static Go binary wrapping `gh`, providing
   the workflow verbs with gates enforced as code.

## What exists

**Shipped:** v1.1.0 of the `gh` extension, implementing protocol 1.0
(`version` prints both: `v1.1.0 (protocol 1.0)`; the pointer's protocol
major is checked, another major refuses; `gh` itself must be 2.50.0 or
later, or the CLI refuses `GH_TOO_OLD` before any verb runs). Verbs:
`init`, `status`, `milestone new/evidence/close`,
`task new/start/finish`, `checkpoint`,
`role`, `roles diff/show`, `identity new/token/webhook`, and `version` — all implemented,
with machine-readable refusals (`refused[CODE]: detail`, catalogued below)
when a gate blocks. `task start` is role-aware: roles whose contracts forbid
commits (qa, reviewer) get no linked development branch; `roles show <role>`
prints a contract with its `roles/<role>.local.md` extensions appended;
`task finish` deletes the branch it merged and `milestone close` sweeps what
its tasks left; `milestone new`, `task finish` and `milestone close` take
`--dry-run` — every gate in order with its outcome, then what the verb would
do, nothing written, the same refusal code. What changed and when: [CHANGELOG.md](../CHANGELOG.md). Not
yet here: any backend other than GitHub, and GitHub Enterprise Server —
github.com only.

**Who holds a seat.** Every role is always staffed, by exactly one of four
kinds of principal: the operator themselves (`~` in the routing table), a
named human (a username), a GitHub team (`identity: org/team-slug` — any
member holds the role), or a GitHub App identity. Solo is a routing
configuration, not a degraded tier: the qa holder's verdicts are the ones
that count at close, and the reviewer holder's approving review is the one
`task finish` requires. A solo operator therefore needs nothing but `gh auth
login`; `gh codecrew init` scaffolds a project with every role routed to `~`.
When a project outgrows solo, `identity new <role>` mints a dedicated App
identity through the manifest flow and routes it for you — minted with
`--with-approval-permission`, a reviewer App's approvals count toward
GitHub's own required-review rules, which makes fully agent-gated merges
possible ([identities.md](identities.md)).

**How a repo joins.** Every repo in a CodeCrew project carries a
`.codecrew.yml` pointing at the hub — a spoke's is a two-line pointer
(`init --hub owner/repo`); this repo is its own hub (`hub: self`; SPEC §3
on choosing yours). The hub's config also routes the five roles — the four
crew seats and the coordinator, which unrouted is you — and `init` writes
that table for you. Agents dispatched into a CodeCrew repo start at
[AGENTS.md](../AGENTS.md). Agent identities are GitHub Apps; a seat's
first act is `export GH_TOKEN=$(gh codecrew identity token <slug>)`, which
mints a short-lived installation token from the platform's env bindings
or the local key and stub and discovers the installation from the App
itself — see [identities.md](identities.md) (SPEC §5 for the tiers).

**What it has done** — the receipts are on the [landing
page](../README.md#the-receipts); the per-milestone records are in
[milestones/](milestones/). The largest environment it has run in — an
orchestration platform dispatching all five seats — has its own page:
[Platform interop](platform-interop.md).

## Install and use

```sh
gh extension install radiusred/gh-codecrew   # precompiled, all platforms
gh codecrew version        # confirm what you installed (gh never auto-updates extensions)
gh codecrew init           # scaffold a new project (see first-milestone.md)
gh codecrew status         # open milestones, inferred task states, raised gates, notes
gh codecrew role reviewer  # who holds a role: an App, a username, or ~ (you)
gh codecrew help           # the full verb list

# or build from source (single static binary; requires gh on PATH):
go build -o gh-codecrew ./cmd/codecrew
```

## Refusal codes

A blocked gate exits non-zero with `refused[CODE]: detail`. The code is for
the agent; the detail is for the human. All thirty, by the verb that raises
them (the source is the catalogue of record — `refuse("CODE"` in
`internal/cli/`):

**any verb that loads `.codecrew.yml`**

- `PROTOCOL_MISMATCH` — the pointer's protocol major differs from the one
  this binary implements (SPEC §5); `"0.1"` and a missing field proceed
  with a note.
- `GH_TOO_OLD` — the installed `gh` is older than 2.50.0, the floor
  `task finish` and the close's branch sweep need (`gh pr checks --json`);
  the detail names both versions. A `gh --version` banner the CLI cannot
  parse proceeds with a note.

**`task new`**

- `NOT_FOUND` — no open milestone with that number in the hub.

**`task start`**

- `NOT_A_TASK` — the issue is not labelled `cc:task`.
- `CLOSED` — the task is already closed.
- `NO_PLAN` — the Plan section is empty; plans come before work (SPEC §4).

**`task finish`** (in the order the gates are checked)

- `CLOSED` — the task is already closed.
- `GATED` — a `cc:needs-decision` gate is raised; a human resolves it and
  removes the label.
- `GATE_UNRECORDED` — a gate was raised and the label removed, but no
  `**Gate resolved:**` comment records the decision (SPEC §8).
- `NOT_OWNER` — the task was started by another seat (the `**Started by**`
  record `task start` posts on every start, accepted only from the login
  it names; the assignee for tasks that predate it; the same seat is the
  same login or the same routed role — any member of a team-held role). The seat that
  started a task finishes it — dispatch it; hand it over by running
  `task start` as the new seat (the latest record wins — the path when the
  starter has left); or a human operator overrides on the record with
  `--bypass`. A task with no start record is not gated.
- `NO_PR` — no open PR closes the task.
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
- `SELF_CONFIRM` — `--operator-confirm` was given by a crew identity; only a
  human operator can waive review.
- `CREW_BYPASS` — `--bypass` was given by a crew identity; a bypass is an
  operator's act.

**`milestone new`**

- `MILESTONE_NUMBER_TAKEN` — the issue was created, but another milestone
  already carries its `M<n>:` prefix (the listing the number came from
  lagged) and the verb's own renumbering failed or found the next number
  taken too; the detail names both issues and the fix: retitle the new one
  to the next free number and rewrite its `M<n>-R<k>` IDs. A repair that
  succeeds is not a refusal — it prints a `renumbered:` line.

**`milestone close`**

- `NOT_FOUND` — no open milestone with that number.
- `OPEN_TASKS` — tasks are still open, listed with their inferred state.
- `NO_REQUIREMENTS` — the milestone's `## Requirements` section yields no
  bold IDs, so there is nothing to verdict; IDs written elsewhere in the
  body do not count (`new`, `status` and `evidence` note this first).
- `VERDICT_MISSING` — a requirement has no QA verdict from the qa holder.
- `VERDICT_UNSATISFIED` — the latest verdict on a requirement is not
  `satisfied`.
- `DOC_MISSING` — `docs/milestones/<n>-*.md` is not on the default branch;
  the gathered records are printed above the refusal for the
  doc-synthesizer, which delivers the document as a task (plan, `task
  start`, a PR with `Closes #<task>`, `task finish`).

**`milestone evidence`**

- `NOT_FOUND` — no open milestone with that number.
- `EVIDENCE_UNREACHABLE` — cited links in the milestone's record do not
  resolve; repair them before dispatching QA.

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
