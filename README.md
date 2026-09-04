<p align="center">
  <img src="assets/codecrew-logo.webp" alt="CodeCrew" width="320">
</p>

<p align="center"><strong>Agent-driven software delivery, with the receipts kept in GitHub.</strong></p>

CodeCrew is a small protocol for running coding agents on real work and
keeping an audit trail worth reading: **the record is the work.** Milestones
are GitHub issues. Tasks are issues with a plan in them. Decisions and
deviations are comments written at the moment they happen, in a fixed shape a
machine can find later. The gates — CI green, an independent approval, a human
sign-off wherever one was asked for — are enforced by a CLI that *refuses*
rather than reminds, and at milestone close one role compiles the comments
into a document explaining why the system is the way it is. There is no
server, no dashboard, no new place to look: it is `gh`, issues, PRs, and five
short role contracts any agent harness can read.

> **The full story is at [codecrew.works](https://codecrew.works)** — how it
> works and what it has done on the [home page](https://codecrew.works), the
> guides and the per-milestone records under
> [codecrew.works/docs](https://codecrew.works/docs/), and the writing at
> [codecrew.works/blog](https://codecrew.works/blog/). The protocol itself is
> [SPEC.md](SPEC.md), here in this repo.

## The routing table

<img src="assets/svg/four-seats.svg" alt="The four seats — implementer, reviewer, qa, doc-synthesizer — each a contract file, each held by exactly one of: you, a colleague by username, a GitHub team, or an App identity." width="720">

Implementer, reviewer, qa, doc-synthesizer, and the coordinator that
dispatches them. Each is a contract — a short markdown file — not an account,
and each is held by you (`~`), a colleague by username, a GitHub team, or a
GitHub App identity minted for the job. Who holds which seat is one table in
`.codecrew.yml`.

Here is a worked example: the `roles:` section of this repository's own
`.codecrew.yml`, as it stands today. Each row is a seat — the identity that
holds it, and the harness and model it is dispatched under, which can differ
from row to row. There are five rows because the coordinator, the seat that
dispatches the other four, is routed too; `~` means a human holds it, here the
operator who coordinates this hub by hand. Your table will look different, and
every seat pointing at `~` is a complete one.

```yaml
roles:
  implementer:
    harness: claude-code
    model: claude-fable-5
    identity: radiusred-cody
  reviewer:
    harness: codex
    model: gpt-5.5
    identity: radiusred-checky
  qa:
    harness: codex
    model: gpt-5.5
    identity: radiusred-testy
  doc-synthesizer:
    harness: claude-code
    identity: radiusred-wordy
  coordinator:
    identity: ~   # the operator: this hub is coordinated by hand (SPEC §7)
```

Solo is not a degraded mode; it is this table with every seat pointing at you.
When a project outgrows that, `identity new <role>` mints a dedicated App
through GitHub's manifest flow and reroutes the seat. The protocol does not
change — only the table does.

## Start now

Four lines, then one sentence to your agent:

```sh
gh extension install radiusred/gh-codecrew
cd my-project            # any repo on GitHub, brand new or years old
gh codecrew init         # writes and commits .codecrew.yml, roles/, AGENTS.md, CLAUDE.md, ROADMAP.md
claude                   # or codex, or whichever coding agent you run
```

> Let's build this project!

That is the whole onboarding. **You do not run the verbs. Your agent does.**
It opens a milestone and hangs tasks off it —

```sh
gh codecrew milestone new --title "Walking skeleton" \
  --goal "A deployed hello-world proving the delivery pipeline end to end." \
  --requirement "visiting the app's URL returns a greeting"
gh codecrew task new --milestone 1 --title "Serve the greeting" --requirements M1-R1
```

— then plans, works and merges under gates that refuse rather than remind.
You are needed at three moments: when a gate asks you a question, when a PR
wants your review, and when a milestone wants your verdict.

Two prerequisites worth knowing before the first PR. The repo needs
pull-request CI of some kind, because `task finish` refuses a PR that reports
no checks at all (absence never satisfies a gate); ten lines of workflow do,
and the [quickstart](docs/first-milestone.md#5-finish-the-task) shows them.
And `gh` must be 2.50.0 or later — `gh --version` says. From there, the
[quickstart](https://codecrew.works/docs/first-milestone/) walks one milestone
end to end: open it, plan a task, do the work, verdict it, close it, and read
the document the close produced.

## Read next

- [codecrew.works](https://codecrew.works) — what CodeCrew is, how it works,
  and [what it has done](https://codecrew.works/#codecrew-works)
- [SPEC.md](SPEC.md) — the protocol itself: topology, state model,
  configuration, verbs, roles, gates
- [docs/](docs/) — the guides, rendered at
  [codecrew.works/docs](https://codecrew.works/docs/); the per-milestone
  records are in [docs/milestones/](docs/milestones/), one per milestone
  delivered by the protocol it documents
- [ROADMAP.md](ROADMAP.md) says what each finished milestone was for
  (`gh codecrew status` names the open one), and
  [CHANGELOG.md](CHANGELOG.md) what each release shipped
- [CONTRIBUTING.md](CONTRIBUTING.md) — the contribution process, which is the
  protocol

Licensed under [Apache 2.0](LICENSE).
