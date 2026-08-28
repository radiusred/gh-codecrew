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
5. [SPEC.md](../SPEC.md) — the protocol itself: topology, state model,
   configuration, verbs, roles, gates.
6. [Founding decisions](founding-decisions.md) and the per-milestone
   records in [milestones/](milestones/) — the "why", as it was recorded.
7. [GSD vs. "just let the model orchestrate"](gsd-vs-frontier-orchestration.md)
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
   for the implementer, reviewer, qa, and doc-synthesizer roles, loadable by
   any agent (Claude Code, Codex, Gemini CLI, or an orchestrator's company).
   A project extends a contract without forking it in
   `roles/<role>.local.md` (SPEC §7).
3. **A CLI** — `codecrew`, a single static Go binary wrapping `gh`, providing
   the workflow verbs with gates enforced as code.

## What exists

**Shipped:** v1.0.0 of the `gh` extension, implementing protocol 1.0
(`version` prints both: `v1.0.0 (protocol 1.0)`; the pointer's protocol
major is checked, another major refuses). Verbs: `init`, `status`,
`milestone new/evidence/close`, `task new/start/finish`, `checkpoint`,
`role`, `roles diff/show`, `identity new`, and `version` — all implemented,
with machine-readable refusals (`refused[CODE]: detail`, catalogued below)
when a gate blocks. `task start` is role-aware: roles whose contracts forbid
commits (qa, reviewer) get no linked development branch; `roles show <role>`
prints a contract with its `roles/<role>.local.md` extensions appended;
`task finish` deletes the branch it merged and `milestone close` sweeps what
its tasks left. What changed and when: [CHANGELOG.md](../CHANGELOG.md). Not
yet here: any backend other than GitHub, and GitHub Enterprise Server —
github.com only.

**Who holds a seat.** Every role is always staffed, by exactly one of four
kinds of principal: the operator themselves (`~` in the routing table), a
named human (a username), a GitHub team (`identity: org/team-slug` — any
member holds the role), or a GitHub App identity. Solo is a routing
configuration, not a degraded tier: the qa holder's verdicts are the ones
that count at close, and the reviewer holder's approving review is the one
`task finish` requires. A solo operator therefore needs nothing but `gh auth
login`; `codecrew init` scaffolds a project with every role routed to `~`.
When a project outgrows solo, `identity new <role>` mints a dedicated App
identity through the manifest flow and routes it for you — minted with
`--with-approval-permission`, a reviewer App's approvals count toward
GitHub's own required-review rules, which makes fully agent-gated merges
possible ([identities.md](identities.md)).

**How a repo joins.** Every repo in a CodeCrew project carries a
`.codecrew.yml` pointing at the hub — a spoke's is a two-line pointer
(`init --hub owner/repo`); this repo is its own hub (`hub: self`; SPEC §3
on choosing yours). The hub's config also routes the four roles, and `init`
writes that table for you. Agents dispatched into a CodeCrew repo start at
[AGENTS.md](../AGENTS.md). Agent identities are GitHub Apps; short-lived
tokens come from the `codecrew-token` script this repo ships — your project
installs it from upstream, see [identities.md](identities.md) (SPEC §5 for
the tiers).

**What it has done** — the receipts are on the [landing
page](../README.md#the-receipts); the per-milestone records are in
[milestones/](milestones/).

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
the agent; the detail is for the human. All twenty-two, by the verb that raises
them (the source is the catalogue of record — `refuse("CODE"` in
`internal/cli/`):

**any verb that loads `.codecrew.yml`**

- `PROTOCOL_MISMATCH` — the pointer's protocol major differs from the one
  this binary implements (SPEC §5); `"0.1"` and a missing field proceed
  with a note.

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
  doc-synthesizer.

**`milestone evidence`**

- `NOT_FOUND` — no open milestone with that number.
- `EVIDENCE_UNREACHABLE` — cited links in the milestone's record do not
  resolve; repair them before dispatching QA.

Licensed under [Apache 2.0](../LICENSE).
