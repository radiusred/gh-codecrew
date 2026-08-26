<p align="center">
  <img src="assets/codecrew-logo.webp" alt="CodeCrew" width="320">
</p>

<p align="center"><strong>Agent-driven software delivery, with the receipts kept in GitHub.</strong></p>

## Why you'd want a crew

You've probably run a coding agent on a real codebase by now. It went well,
mostly — and then three weeks later somebody asked *why* the retry logic
looks like that, and the answer was in a chat transcript nobody saved. The
heavyweight fix is a process framework: planning documents, phase gates, a
directory of state files that drift from the code the moment you stop
looking. The lightweight fix is vibes. Neither hands your colleague an
audit trail.

CodeCrew is a third option, and a small one: **the record is the work.**
Milestones are GitHub issues. Tasks are issues with a plan in them.
Decisions and deviations are comments written at the moment they happen,
in a fixed shape a machine can find later. The gates — CI green, an
independent approval, a human sign-off wherever one was asked for — are
enforced by a CLI that *refuses* rather than reminds. At milestone close,
one role compiles the comments into a document that explains why the
system is the way it is, and that document's PR goes through the same
review as code.

There is no server, no dashboard, no new place to look. It is `gh`,
issues, PRs, and four short role contracts any agent harness can read.

## How it works, in four beats

### 1. A hub, and spokes

<img src="assets/svg/hub-and-spokes.svg" alt="One hub repo holds the four role contracts, the routing table and the milestone issues; each spoke repo holds tasks, PRs and CI, with a two-line pointer back to the hub. For a single project the hub is its own spoke." width="720">

One repo is the hub: it holds the role contracts, the routing table (who
holds which seat) and the milestone issues. Every other repo in the project
is a spoke with a two-line pointer file. For a single project, the hub is
the spoke — `hub: self`. State never leaves GitHub, and there is no other
message bus, so any two harnesses interoperate by construction.

### 2. Four seats, always staffed

<img src="assets/svg/four-seats.svg" alt="The four seats — implementer, reviewer, qa, doc-synthesizer — each a contract file, each held by exactly one of: you, a colleague by username, a GitHub team, or an App identity." width="720">

Implementer, reviewer, qa, doc-synthesizer. Each is a contract — a short
markdown file — not an account. A seat is held by you (`~`), a colleague
by username, a GitHub team, or a GitHub App identity minted for the job.
Solo is not a degraded mode; it is the routing table with every seat
pointing at you. When a project outgrows that, `identity new <role>` mints
a dedicated App through GitHub's manifest flow and reroutes the seat. The
protocol does not change — only the table does.

### 3. A milestone, from open to record

<img src="assets/svg/milestone-lifecycle.svg" alt="milestone new, then task new with a plan, task start, PR and review, task finish, qa verdicts, milestone close, and the synthesized document in docs/milestones." width="720">

`milestone new` opens the tracking issue with its requirements. `task new`
hangs tasks off it; the implementer writes a plan into the task before the
first commit; `task start` cuts the branch; a PR closes the task; `task
finish` is the one merge point. The qa seat posts a verdict per
requirement. `milestone close` gathers the recorded decisions and
deviations for the doc-synthesizer, and refuses to close until the
milestone document exists on the default branch.

### 4. Gates that refuse

<img src="assets/svg/gate-refusing.svg" alt="task finish checking CI, an open human gate and the routed reviewer's approval; the approval is missing, so it exits 1 with refused[NO_HOLDER_REVIEW] and a one-line explanation." width="720">

A blocked gate exits non-zero with a machine-readable reason:
`refused[NO_CHECKS]` when a PR reports no CI at all (absence does not
satisfy a gate, and there is no override), `refused[NO_HOLDER_REVIEW]`
when the routed reviewer has not approved, `refused[VERDICT_MISSING]` when
a milestone tries to close without a verdict on every requirement,
`refused[EVIDENCE_UNREACHABLE]` when a link the record cites no longer
resolves. An agent acts on the code; a human reads the detail. Those are
four of twenty; the [introduction](docs/introduction.md) shows the ones
you will meet first, and every code is grep-able in the source as
`refuse("CODE"` — a proper catalogue is on the docs-polish list, not yet
written.

## Start solo in five minutes

```sh
gh extension install radiusred/gh-codecrew
gh codecrew version     # gh never auto-updates extensions; know what you have
gh codecrew init        # scaffold: .codecrew.yml, roles/, AGENTS.md, ROADMAP.md
gh codecrew status      # open milestones, task states, raised gates
gh codecrew help        # the full verb list
```

That is a working hub with every seat routed to you. The
[quickstart](docs/first-milestone.md) walks a first milestone end to end —
open it, plan a task, do the work, verdict it, close it, read the document
the close produced. From there the
[ladder](docs/first-milestone.md#what-next) is: split the seats across
sessions, give crew members their own identities, hand dispatch to an
orchestrator. Same commands at every rung.

## The receipts

We would rather show you than tell you, so here is what the protocol has
actually done so far.

- **Five milestones of this framework were delivered *with* it** —
  agent-authored PRs under GitHub App identities, independent review,
  deterministic CI gates, QA verdicts enforced at close, and a synthesized
  document for each: [docs/milestones/](docs/milestones/).
- **This repository is agent-staffed, and you can check.** Four App
  identities hold the four seats: cody implements, checky reviews, testy
  verdicts, wordy writes — this page included. PRs here merge on checky's
  approval; a reviewer App minted with write access satisfies GitHub's own
  required-review rule, which is what makes a fully agent-gated merge
  possible ([identities.md](docs/identities.md)).
- **A stranger ran it from the scaffold alone.**
  [davison/numberguess](https://github.com/davison/numberguess) went from
  `gh codecrew init` to three closed milestones with one human and a Codex
  session, transcripts committed, and then served as the proving ground
  for the dedicated reviewer seat.
- **The first spoke published its own announcement.**
  [radiusred/www](https://github.com/radiusred/www) is driven from this
  hub through the installed extension; its first delivery was
  [a blog post introducing CodeCrew, delivered by the protocol it
  describes](https://www.radiusred.uk/blog/posts/2026-08-20-this-post-was-delivered-by-the-framework-it-introduces/).

Things we have not done yet, in the interest of not overselling: any
backend other than GitHub; a full orchestration-platform run end to end;
a gatherer that catches every record variant (the M5 close missed three
of one issue's four Decisions because they were labelled `**Decision
(…):**` — the fix is on the open milestone, and the M5 document says so);
a 1.0. The [roadmap](ROADMAP.md) says which milestone is open and what it
is for.

## Read next

- [Introduction](docs/introduction.md) — what CodeCrew is, precisely: the
  protocol, the contracts, the CLI and its refusal codes
- [Your first milestone](docs/first-milestone.md) — the quickstart
- [Identities](docs/identities.md) — running solo, minting App identities,
  dispatching a role session
- [SPEC.md](SPEC.md) — the protocol itself
- [Founding decisions](docs/founding-decisions.md) — the trade-offs, and
  [the essay that started it](docs/gsd-vs-frontier-orchestration.md) — one
  person's experience with GSD across several projects; the motivation,
  not a verdict on GSD

Licensed under [Apache 2.0](LICENSE).
