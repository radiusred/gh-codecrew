# Your first milestone

This is the whole path, start to finish, for the most common setup: **one
human, one coding agent to hand** (a Claude Code, Codex, or similar session),
and a brand-new project. Install CodeCrew, scaffold a repo, and drive one
milestone from opening to audited close. No GitHub Apps, no configuration
beyond what the tool writes for you, no prior knowledge assumed — and no
paid GitHub plan: everything here works on the free tier.

Two ideas carry everything else. First, **project state lives in GitHub** —
issues, sub-issues, labels, and PRs — so there is nothing to keep in sync
and nothing to lose. Second, **the verbs refuse**: when a gate isn't met,
the CLI exits with `refused[CODE]: detail` telling you exactly what's
missing. A refusal is the guide-rail working, not an error. This guide walks
you into each one you'll meet, on purpose.

## You and your agent

The expected division of labour: **your agent runs the loop, you make the
calls**. The agent acts as orchestrator and implementer in one session — it
runs the verbs, writes plans into task issues, does the work, opens the
PRs, and stops at every gate; you review, resolve raised gates, confirm
merges, and give verdicts. Start your agent in the scaffolded repo and
say "Let's build this project!" — `init` writes the entry point it reads
(`AGENTS.md`, and from v1.0.1 a `CLAUDE.md` that imports it), so it knows
where it is; point it at this page too if you want it to follow the long
form. The refusals keep it honest even when you're not watching.

Every command below is equally pasteable by a human — the protocol is fully
human-operable end to end, and running one milestone by hand is a fine way
to learn what your agent is doing. But the framework earns its keep when
the agent does the typing and the record proves what happened.

## 0. Prerequisites

A GitHub account, [`gh`](https://cli.github.com/) 2.50.0 or later
installed and signed in (`gh auth login`; `gh --version` to check —
`task finish` reads `gh pr checks --json`, and a distribution-packaged
`gh` can be older), `git`, and (recommended) your coding agent of choice.
One thing to know before your first PR: `task finish` refuses a PR that
reports no CI check at all, with no override — `init` does not write a
workflow, so step 5 shows the ten-line one you add yourself.
You and your agent are the whole crew today (SPEC §5 calls this solo
routing; the [What next](#what-next) ladder covers expanding it).

## 1. Install and scaffold

```sh
gh extension install radiusred/gh-codecrew
gh codecrew version        # confirm what you got (gh extension list agrees)

gh repo create my-project --private --clone   # or use an existing repo
cd my-project
gh codecrew init
```

One thing worth knowing up front: `gh` never auto-updates extensions — not
for patches, not for majors. Updating is always your act:
`gh extension upgrade codecrew`, then `gh codecrew version` to confirm.

`init` writes your `.codecrew.yml` (the project's pointer file, with a
routing table declaring all five roles — `~` means "held by you"), a
`ROADMAP.md` seed, the five role contracts under `roles/`, and an
`AGENTS.md` entry point for any agent you later dispatch — plus, from
v1.0.1, a `CLAUDE.md` that imports it, because Claude Code loads
`CLAUDE.md` and never `AGENTS.md`. Commit and push:

```sh
git add -A && git commit -m "chore: scaffold codecrew" && git push
```

If your default branch is protected (pull requests required), the
scaffold is your project's first PR — and the one merge you do yourself:
no milestone or task exists yet, so there is nothing for `task finish` to
merge it through. Everything after it goes through a task.

## 2. Open a milestone

```sh
gh codecrew milestone new --title "Walking skeleton" \
  --goal "A deployed hello-world proving the delivery pipeline end to end." \
  --requirement "visiting the app's URL returns a greeting"
```

This creates the milestone tracking issue — the canonical milestone object —
and adds a row to `ROADMAP.md` (commit it with your next PR). Each
`--requirement` (repeat it for more) lands under the issue's
**Requirements** section with a bold ID the tooling checks "done" against
— numbered by the CLI, `M1-R1`, `M1-R2`, … in the order you gave them, and
printed back so you see what the close gate will count:

```markdown
## Requirements
- **M1-R1** — visiting the app's URL returns a greeting
```

Only bold IDs in that section are requirements; ones written under Goal
are prose, and a milestone with none refuses to close. To add or reword
one later, edit it in the web UI, or safely from the shell with
fetch-edit-put — pull the current body down, change only the Requirements
section, put it back:

```sh
gh issue view 1 --json body -q .body > milestone.md
# edit milestone.md — fill in the Requirements section, leave the rest
gh issue edit 1 --body-file milestone.md && rm milestone.md
```

Never a raw `gh api` PATCH built from shell variables: one empty variable
silently replaces the whole body.

## 3. Open a task, plan it, start it

```sh
gh codecrew task new --milestone 1 --title "Serve the greeting" --requirements M1-R1
```

The task issue is attached to the milestone as a sub-issue — `gh codecrew
status` now shows your project. Try to start it:

```sh
gh codecrew task start 3   # your task's issue number
```

It refuses: `refused[NO_PLAN]`. Plans come before work — edit the task
issue's **Plan** section with your intended changes (a trivial task may have
a trivial plan; it may not have an absent one), then run `task start` again.
It assigns you and creates a linked branch.

## 4. Do the work

```sh
git fetch && git switch task/3-serve-the-greeting
# ...write the code, in commits that reference the task...
git commit -m "feat: serve the greeting (#3)"
git push
gh pr create --fill --body "Closes #3"
```

The PR description is the task's summary artifact — say what was done and
which requirement it satisfies. If a decision came up mid-task (a trade-off,
a rejected alternative), record it as an issue or PR comment starting
`**Decision:**` *at the moment you make it*; if the work departs from the
plan, a comment starting `**Deviation:**` with a `**Why:**` line. Both
shapes are collected automatically at milestone close — one record per
comment, so give each its own.

## 5. Finish the task

```sh
gh codecrew task finish 3
```

Refused again — no approving review from a non-author, and you can't approve
your own PR. Solo has an honest degradation:

```sh
gh codecrew task finish 3 --operator-confirm
```

This records an explicit confirmation comment on the PR — stating that
author and operator are the same principal — then rebase-merges and deletes
the task branch. The task issue closes via its `Closes #3` keyword.

If your repo has no CI yet, you'll meet one more refusal first:
`refused[NO_CHECKS]`. The deterministic gate can't be satisfied by a repo
that runs nothing, and there's no override — the fix is a ten-line workflow:

```yaml
# .github/workflows/test.yml
name: Test
on: [pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: make test   # or your language's test command
```

Push it to the task branch, let the check report, and `task finish` (with
failing or pending checks it refuses too — that's the gate having teeth).

## 6. Verdict — you hold the qa role

Try to close the milestone:

```sh
gh codecrew milestone close 1
```

`refused[VERDICT_MISSING]` — every requirement needs a QA verdict, and the
qa role is unrouted, which means it's yours. Change hats: build what was
merged, exercise the requirement's *intent* (not just what the code's own
tests prove), and post a comment on the **milestone issue**, one line per
requirement, in exactly this form:

```markdown
**M1-R1 — satisfied.** Deployed from main, visited the URL, got the greeting.
```

The evidence is yours, from your own execution — `roles/qa.md` is the
contract you just performed.

## 7. The milestone document

```sh
gh codecrew milestone close 1
```

One refusal left: `refused[DOC_MISSING]` — and above it, the CLI prints
every `**Decision:**` and `**Deviation:**` comment recorded during the
milestone. That's the raw material. Write `docs/milestones/1-<slug>.md` from
it — the "why" record: what was decided, what was traded off, what deviated
from plan — flip the ROADMAP row to Done, and merge it like any other change
(a PR, `--operator-confirm` if you like the discipline).

## 8. Close

```sh
gh codecrew milestone close 1
```

All gates pass: tasks closed, verdicts satisfied, document merged. Any task
branch still lying around is swept first (only merged or empty ones — the
output says what went and what stayed), then the milestone issue closes
with the record comment, and everything above — the
plan, the decisions, the confirmation, the verdicts, the document — is
permanent, linked, and public to whoever you choose.

## What next

The loop above is the same shape at every size — more requirements per
milestone, more tasks per requirement, CI required checks to give `task
finish`'s gate real teeth ([SPEC §8](../SPEC.md)). The interesting growth is
in the crew, and it's a ladder:

1. **Split the roles across sessions.** Your agent doing the work and then
   verdicting its own milestone shares one context's blind spots. The cheap
   fix: have your harness launch a sub-agent per role, each briefed with the
   role's contract from `roles/` (plus any `roles/<role>.local.md` beside
   it — from the next release, `gh codecrew roles show <role>` prints the
   two composed) — so a fresh-context QA reading `roles/qa.md` probes what
   the implementer's context wouldn't. Or cross model families:
   dispatch another LLM through its own CLI for the reviewer or qa seat.
   Same identities, same commands — just different eyes.
2. **Give crew members their own identities.** When you want the record to
   show *which* agent did what — and GitHub itself to enforce that the
   approver isn't the author — route roles to other humans by username, to
   a GitHub team (`identity: org/team-slug`, any member holds the seat), or
   to GitHub App identities: [identities.md](identities.md). The protocol
   doesn't change; only the routing table does.
3. **Full orchestration platforms** — an orchestrator dispatching the whole
   crew against the routing table, webhooks instead of polling. Run end to
   end on Paperclip: three milestones on a proving-ground repo, the third
   with one gate and no operator touch on the workflow besides it — the
   findings and what they changed are on
   [#119](https://github.com/radiusred/gh-codecrew/issues/119);
   the interop doc written from them is
   [#54](https://github.com/radiusred/gh-codecrew/issues/54).
