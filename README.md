<p align="center">
  <img src="assets/codecrew-logo.webp" alt="CodeCrew" width="320">
</p>

CodeCrew is a protocol for agent-driven software delivery and a `gh` extension
that enforces it. Project state lives in GitHub and nowhere else: a milestone
is an issue, a task is an issue with a plan in its body, decisions and
deviations are comments in a fixed shape, and the work of a task is a branch
and a PR. The gates — CI green, an approving review from whoever holds the
reviewer seat, a QA verdict on every requirement, a human sign-off wherever one
was raised — are checked by the CLI, which exits non-zero with
`refused[CODE]: detail` when one is unmet. At milestone close the
doc-synthesizer seat compiles the recorded comments into
`docs/milestones/<n>-<slug>.md`. There is no server, no database and no state
files: the dependencies are `gh` 2.50.0 or later, a GitHub repo, and CI on its
pull requests.

[codecrew.works](https://codecrew.works) is the marketing and introduction
site, and carries the [blog](https://codecrew.works/blog/). This page is the
technical entry point: it is written for whoever wants more detail than the
site gives — someone running a coding harness, or the agent reading on their
behalf — and every reference it links is at source in this repository rather
than on the site.

## The routing table

<img src="assets/svg/four-seats.svg" alt="The four seats — implementer, reviewer, qa, doc-synthesizer — each a contract file, each held by exactly one of: you, a colleague by username, a GitHub team, or an App identity." width="720">

Four seats — implementer, reviewer, qa, doc-synthesizer — and a coordinator
that dispatches them. Each is a contract file under [roles/](roles/), not an
account, and each is held by one of four kinds of principal: you (`~`), a
username, a GitHub team (`owner/team-slug`), or a GitHub App identity. Who
holds which seat is the `roles:` table in the hub's `.codecrew.yml`.

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

Solo is a routing configuration, not a reduced protocol: `gh codecrew init`
writes the table with every seat routed to `~`, and the same gates apply.
`gh codecrew identity new <role>` mints a dedicated App for a seat through
GitHub's manifest flow and rewrites its row; nothing else changes
([docs/identities.md](docs/identities.md)). The field reference is
[SPEC.md](SPEC.md) §5.

## Start now

```sh
gh extension install radiusred/gh-codecrew
cd my-project            # any repo on GitHub, brand new or years old
gh codecrew init         # writes and commits .codecrew.yml, roles/, AGENTS.md, CLAUDE.md, ROADMAP.md
```

`init` scaffolds the project with every seat routed to `~`. After it, the verbs
are run by the coding agent rather than by the person driving it:
[AGENTS.md](AGENTS.md) tells an agent dispatched into the repo where the hub
and the contracts are, and `gh codecrew roles show <role>` prints the contract
it works to, local extensions included.

```sh
gh codecrew milestone new --title "Walking skeleton" \
  --goal "A deployed hello-world proving the delivery pipeline end to end." \
  --requirement "visiting the app's URL returns a greeting"
gh codecrew task new --milestone 1 --title "Serve the greeting" --requirements M1-R1
gh codecrew task start 3     # refuses NO_PLAN until the task body carries a plan
```

`gh codecrew help` lists the rest; `gh codecrew status` reports the open
milestone, the inferred task states and any raised gate. A human is needed at
three points, whatever harness is driving: a gate raised for a decision
(`cc:needs-decision`), a PR that needs the reviewer seat's approval, and the QA
verdicts a milestone cannot close without.

Two requirements the CLI checks. `gh` must be 2.50.0 or later — `task finish`
reads `gh pr checks --json`, and an older `gh` refuses `GH_TOO_OLD` before any
verb runs. And the repo needs pull-request CI: `task finish` refuses `NO_CHECKS`
on a PR that reports no checks at all, with no override. `init` does not write a
workflow; [docs/first-milestone.md](docs/first-milestone.md#5-finish-the-task)
carries a ten-line one.

## Read next

Reference documentation, at source in this repository — an agent dispatched
into a CodeCrew repo starts at [AGENTS.md](AGENTS.md) and its role contract
under [roles/](roles/):

- [docs/introduction.md](docs/introduction.md) — what CodeCrew is, precisely:
  the three parts, what is shipped, and all thirty-one refusal codes by the
  verb that raises each
- [docs/first-milestone.md](docs/first-milestone.md) — one milestone end to
  end, solo: open it, plan a task, do the work, verdict it, close it
- [docs/identities.md](docs/identities.md) — routing seats to humans, teams
  and GitHub App identities, and dispatching a role session
- [docs/platform-interop.md](docs/platform-interop.md) — hosting the crew on
  an orchestration platform: the coordinator seat, wake paths, onboarding
- [docs/extensions.md](docs/extensions.md) — extending a role contract in
  `roles/<role>.local.md` without forking it
- [docs/founding-decisions.md](docs/founding-decisions.md) — the trade-offs
  the design was chosen against
- [docs/milestones/](docs/milestones/) — one record per delivered milestone,
  synthesized from the decisions recorded while it was built
- [SPEC.md](SPEC.md) — the protocol: topology, state model, configuration,
  verbs, roles, gates
- [CONTRIBUTING.md](CONTRIBUTING.md) and [SECURITY.md](SECURITY.md)
- [ROADMAP.md](ROADMAP.md) — the finished milestones; [CHANGELOG.md](CHANGELOG.md)
  — what each release shipped

Licensed under [Apache 2.0](LICENSE).
