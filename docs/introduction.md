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

**On `main` and not yet released:** protocol 2.0 — the `.codecrew/` layout
this page describes throughout, the typed identity grammar, routing that
refuses rather than degrading, and `gh codecrew migrate`, which moves a repo
still on the 1.x layout. A 2.0 binary refuses a 1.x repo and names that verb;
there is no dual read. The released v1.2.0 extension implements none of it.
The release is [M14](https://github.com/radiusred/gh-codecrew/issues/269);
what 2.0 broke and how to migrate is in
[CHANGELOG.md](../CHANGELOG.md) and
[the M13 record](milestones/13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md).

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

A blocked gate exits 1 — every failure does, and there is no exit-code
taxonomy — with one line on stderr:

```
codecrew: refused[CODE]: detail
```

The **code** is for the agent: a fixed vocabulary whose meanings are stable
within a protocol major, so an orchestrator branches on it. The **detail**
is for the human, and is free to be reworded in any release; nothing should
parse it. A `note:` line on the same stream is advisory and not a failure.

**The catalogue is [SPEC §10](../SPEC.md#10-the-cli)** — every code, the
verbs that raise it and what it means, in one table. It is the single list:
this page carried a second one until protocol 2.0, and keeping two in step
across a milestone's worth of new codes proved to be exactly the drift the
table exists to prevent. The source of record behind both is `refuse("CODE"`
in `internal/cli/`, and a test fails the build when the table and the source
disagree.

The fix for a refusal is in the detail line the CLI prints, which names the
condition met and the way out — read it before reaching for a code table.
[Your first milestone](first-milestone.md) walks through the ones a new
project meets first: `NO_PLAN`, `NO_CHECKS`, `VERDICT_MISSING`,
`DOC_MISSING`.

Licensed under [Apache 2.0](../LICENSE).
