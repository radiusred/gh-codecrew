# CodeCrew

A lightweight framework for agent-driven software delivery: the auditability
and reproducible discipline of heavyweight frameworks, without the ceremony.
Project state lives in the tools teams already use — GitHub issues, PRs, and
CI — and the only documents the framework produces are per-milestone records
of the decisions that shaped the system and why.

CodeCrew is three things:

1. **A protocol** ([SPEC.md](SPEC.md)) — conventions for representing
   milestones, tasks, plans, decisions, deviations, and gates in GitHub, and
   how agents and humans transact over them.
2. **Role contracts** ([roles/](roles/)) — harness-neutral prompt files for
   the implementer, reviewer, qa, and doc-synthesizer roles, loadable by any
   agent (Claude Code, Codex, Gemini CLI, or an orchestrator's company).
3. **A CLI** — `codecrew`, a single static Go binary wrapping `gh`, providing
   the workflow verbs with gates enforced as code.

## Status

The protocol is specified and the CLI works — released as v0.2.3: `init`,
`status`, `milestone new/close`, `task new/start/finish`, `checkpoint`,
`role`, and `version` are all implemented, with machine-readable refusals
(`refused[CODE]: detail`) when a gate blocks — `NO_CHECKS` when a PR reports
no CI checks at all (the deterministic gate cannot be satisfied by absence,
and there is no override), `GATE_UNRECORDED` when a raised gate lacks a
recorded resolution, and `VERDICT_MISSING`/`VERDICT_UNSATISFIED` when a
milestone tries to close without a satisfied QA verdict on every
requirement. `task start` is role-aware: roles whose contracts forbid
commits (qa, reviewer) get no linked development branch.

Solo is a routing configuration, not a degraded tier: every role is always
staffed, by a GitHub App, a named human, or — `~` — the operator themselves;
the qa holder's verdicts are the ones that count at close, and the reviewer
holder is who PRs ask to review. A solo operator therefore needs nothing but
`gh auth login`; `codecrew init` scaffolds a new project with every role
routed to `~`, and the [quickstart](docs/first-milestone.md) walks the first
milestone end to end.

Four milestones have been delivered *with* the protocol — agent-authored
PRs under GitHub App identities, non-doer review, deterministic CI gates, QA
verdicts enforced at close, and synthesized closing documents. The first
spoke is live: [radiusred/www](https://github.com/radiusred/www) is driven
from this hub through the installed extension, and its first delivery was
[a blog post introducing CodeCrew, published by the protocol it
describes](https://www.radiusred.uk/blog/posts/2026-08-20-this-post-was-delivered-by-the-framework-it-introduces/).
The first stranger's project is public too:
[davison/numberguess](https://github.com/davison/numberguess) was taken from
`gh codecrew init` to a closed milestone by one human and a Codex session
working from the scaffold alone, transcript included. See
[docs/milestones/](docs/milestones/) for the per-milestone records.

Not yet here: any backend other than GitHub.

## Install and use

```sh
gh extension install radiusred/gh-codecrew   # precompiled, all platforms
gh codecrew version        # confirm what you installed (gh never auto-updates extensions)
gh codecrew init           # scaffold a new project (see docs/first-milestone.md)
gh codecrew status         # open milestones, inferred task states, raised gates
gh codecrew role reviewer  # who holds a role: an App, a username, or ~ (you)
gh codecrew help           # the full verb list

# or build from source (single static binary; requires gh on PATH):
go build ./cmd/codecrew
```

Every repo in a CodeCrew project carries a `.codecrew.yml` pointing at the
hub — a spoke's is a two-line pointer (`init --hub owner/repo`); this repo
is its own hub (`hub: self`; see SPEC §3 for choosing yours). The hub's
config also routes the four roles, and `init` writes that table for you.
Agents dispatched into a CodeCrew repo start at [AGENTS.md](AGENTS.md).
Agent identities are GitHub Apps; short-lived tokens come from
`scripts/codecrew-token` (see SPEC §5 and
[docs/identities.md](docs/identities.md) for the identity tiers — a solo
operator needs nothing but `gh auth login`, self-confirmation recorded).

## Documents

- [SPEC.md](SPEC.md) — the protocol
- [docs/first-milestone.md](docs/first-milestone.md) — the quickstart: a
  stranger's first milestone, end to end
- [docs/identities.md](docs/identities.md) — running solo, and minting
  per-role GitHub App identities
- [docs/founding-decisions.md](docs/founding-decisions.md) — design decisions
  and their trade-offs
- [docs/milestones/](docs/milestones/) — the per-milestone "why" records
- [docs/gsd-vs-frontier-orchestration.md](docs/gsd-vs-frontier-orchestration.md)
  — the analysis that motivated the project

Licensed under [Apache 2.0](LICENSE).
