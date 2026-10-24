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
(`.codecrew/AGENTS.md`, with a root `AGENTS.md` pointing at it and a
`CLAUDE.md` importing that), so it knows where it is; point it at this page
too if you want it to follow the long form. The refusals keep it honest even when you're not watching.

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

`init` writes your `.codecrew/config.yml` (the project's pointer file, with a
routing table declaring all five roles — `~` means "held by you"), a
`ROADMAP.md` seed, the five role contracts under `.codecrew/roles/` — each
with a blank `.codecrew/roles/<role>.local.md` beside it, the place for this
project's own additions to a contract (the file says what belongs there; it
adds nothing until you write into it) — and the entry point for any agent
you later dispatch: `.codecrew/AGENTS.md`, which holds the instructions and
is CodeCrew's to rewrite, a root `AGENTS.md` of two lines pointing at it
(a sentence naming the path and an `@.codecrew/AGENTS.md` import), and a
`CLAUDE.md` importing that root, because Claude Code loads `CLAUDE.md` and
never `AGENTS.md`. If you already have an `AGENTS.md` or a `CLAUDE.md`,
`init` keeps it untouched and prints the exact lines to paste into it —
unless it already points at the instructions, so rerunning `init` on a repo
it scaffolded asks you for nothing.
It commits exactly the files it wrote
(`chore: scaffold codecrew`, with a body naming `init`, its release and
protocol, and each file, as SPEC §4 asks of a housekeeping commit), leaving
anything else you had staged or modified as it was. It never pushes; that is yours (run `init` at the
repository root — it refuses a subdirectory, since the pointer belongs at
the root):

```sh
git push -u origin HEAD
```

(`init` prints the exact command for the branch it committed on.)

If your default branch requires pull requests, `init` will have put the
commit on a branch it made from the default branch, `codecrew-bootstrap`,
and said so: `git push -u origin codecrew-bootstrap` and open the scaffold
PR — the one merge you do yourself, since no milestone or task exists yet
for `task finish` to merge it through. Everything after it goes through a
task.

## 2. Open a milestone

```sh
gh codecrew milestone new --title "Walking skeleton" \
  --goal "A deployed hello-world proving the delivery pipeline end to end." \
  --requirement "visiting the app's URL returns a greeting"
```

This creates the milestone tracking issue — the canonical milestone object
— and nothing else; `ROADMAP.md` gets its row when the milestone is
finished (step 7). Each
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
status` now shows your project. In a brand-new repository the milestone is
issue #1 and this task is #2, so the examples from here on use 2; if your
numbers differ, use your own. Try to start it:

```sh
gh codecrew task start 2   # your task's issue number
```

It refuses: `refused[NO_PLAN]`. Plans come before work — edit the task
issue's **Plan** section with your intended changes (a trivial task may have
a trivial plan; it may not have an absent one), in the web UI or with step
2's fetch-edit-put on issue 2, then run `task start` again.
It assigns you and creates a linked branch.

## 4. Do the work

```sh
git fetch && git switch task/2-serve-the-greeting   # task start prints this line
# ...write the code, in commits that reference the task...
git commit -m "feat: serve the greeting (#2)"
git push
gh pr create --fill --body "Closes #2"
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
gh codecrew task finish 2
```

If your repo has no CI yet, the first refusal is `refused[NO_CHECKS]`. The
deterministic gate can't be satisfied by a repo that runs nothing, and
there's no override — the fix is a ten-line workflow:

```yaml
# .github/workflows/test.yml
name: Test
on: [pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: echo "No tests yet"   # swap in your test command once you have one
```

The last line passes on a brand-new project that has nothing to test yet.
Once you have tests, replace it with the command that runs them
(`make test`, `npm test`, `go test ./...`, …) — that's when the check
starts meaning something.

Commit the file to the task branch and push. If git refuses the push with a
message about the `workflow` scope, your `gh` login can't write workflow
files yet: `gh auth login` doesn't ask for that scope by default, and GitHub
needs it to add or change anything under `.github/workflows/` over HTTPS
([GitHub's list of scopes](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/scopes-for-oauth-apps#available-scopes)).
Add it, then push again:

```sh
gh auth refresh -s workflow
```

A check takes a little while to show up after a push, usually under a
minute. Run `task finish` straight away and it may still say `NO_CHECKS`,
even though the workflow is there: it just hasn't reported yet. Wait for it
(`gh pr checks` shows when it has) and run `task finish` again. With failing
or pending checks it refuses too — that's the gate having teeth.

The next refusal: no approving review from a non-author, and you can't
approve your own PR. Solo still gets a review; it just isn't a GitHub
approval. The norm, even with one human, is a fresh pair of eyes before you
confirm: ask your agent to start a clean session with no memory of the
work — a sub-agent, or a different coding tool — give it the reviewer's
instructions (`gh codecrew roles show reviewer` prints them), and have it
post its findings as a comment on the PR. Those instructions also serve a
reviewer with its own GitHub App identity, so skip their identity steps
when solo: minting a token, posting with `gh pr review` as `<slug>[bot]`,
and stopping if it shares the PR author's identity. Your reviewer works
under your own `gh` login and leaves a comment instead.
Read its findings and the diff, fix what needs fixing, and when you're
happy, confirm:

```sh
gh codecrew task finish 2 --operator-confirm
```

This posts a confirmation comment on the PR — that you reviewed and accepted
the work as both author and operator, since no one independent exists to
approve it — then rebase-merges and deletes the task branch. The task issue
closes via its `Closes #2` keyword.

### When the gate is expensive: let the workflow gate itself

The ten lines above run everything on every pull request, which is right
until the gate takes minutes and a run of docs-only PRs pays for it twice
over. The temptation is `[skip ci]` in the commit message. Don't: GitHub
creates no workflow run and no check run for a skipped push, so there is no
fact for `task finish` to read, and it refuses `NO_CHECKS` — correctly, and
with no override. The commit message is author-controlled; a reported check
is the platform's own fact, and the gate only ever trusts the second one.
(GitHub reads the marker out of the message wherever it appears, so a commit
that merely *mentions* it skips the run too — which is how this section
found out.)

What *does* satisfy the gate is a check that reports `skipping`. So let the
committed workflow decide, from the diff, whether the heavy job needs to
run — a cheap job that always reports, and the expensive one behind
`needs:` and `if:` on its output:

```yaml
# .github/workflows/test.yml
name: Test
on: [pull_request]

jobs:
  changes:
    name: Changed files
    runs-on: ubuntu-latest
    outputs:
      code: ${{ steps.classify.outputs.code }}
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0          # the base commit has to be in the clone
      - id: classify
        run: |
          git diff --name-only \
            ${{ github.event.pull_request.base.sha }}...HEAD > changed.txt
          # Docs are an allowlist: anything not matching is code. A new
          # top-level directory is code until you say otherwise.
          if grep -qvE '^(docs/.*|[^/]*\.md)$' changed.txt; then
            echo "code=true" >> "$GITHUB_OUTPUT"
          else
            echo "code=false" >> "$GITHUB_OUTPUT"
          fi
      - run: make docs-check   # your cheap docs gate: seconds, not minutes

  test:
    name: Go build and test
    needs: changes
    if: needs.changes.outputs.code == 'true'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: make test
```

A docs-only PR now costs one short job; everything else runs the full gate.
Three things make it hold up:

- **Gate the job, not the workflow.** `on: pull_request: paths:` looks like
  the same idea, but a workflow that never triggers reports nothing at all —
  a required check filtered that way stays pending forever and `task finish`
  sees the same absence `[skip ci]` produces. A job skipped by `if:` inside
  a run that did happen reports `Skipped`, and that is a fact.
- **Classify by allowlist.** `.github/**`, the `Makefile`, `scripts/**` and
  every source path are code. Invert the test — list what counts as docs,
  treat the rest as code — so a directory nobody thought about defaults to
  running the gate.
- **It is no weaker than what you had.** A pull request could always neuter
  the workflow; both the change and the classification are in the diff, and
  the reviewer reads that diff either way.

## 6. Verdict — you hold the qa role

Try to close the milestone:

```sh
gh codecrew milestone close 1
```

`refused[VERDICT_MISSING]` — every requirement that has not been struck
needs a QA verdict, and the qa role is yours: `init` set it to `~` ("you")
in `.codecrew/config.yml`, which the CLI and the reference docs also call
*unrouted*. There's nothing to change until you want someone else to hold
it ([What next](#what-next)).
Change hats: build what was merged, exercise the requirement's *intent*
(not just what the code's own tests prove), and post a comment on the
**milestone issue**, one line per requirement, in exactly this form:

```markdown
**M1-R1 — satisfied.** Deployed from main, visited the URL, got the greeting.
```

The evidence is yours, from your own execution — `.codecrew/roles/qa.md` is the
contract you just performed.

The one exception is a requirement struck by a recorded decision: one the
coordination layer withdrew from scope. It needs no verdict, because nothing
was built. The Decision comes first, then `gh codecrew milestone strike`
posts the struck line on the milestone issue, and close counts that line as
terminal ([SPEC §4](../SPEC.md#milestone)).

## 7. The milestone document

```sh
gh codecrew milestone close 1
```

One refusal left: `refused[DOC_MISSING]` — and above it, the CLI prints
every `**Decision:**` and `**Deviation:**` comment recorded during the
milestone. That's the raw material. (Above that is a `note: qa is unrouted`
line, now that your verdicts count: it means the `~` from step 6, and
nothing needs changing.) The refusal says to dispatch the doc-synthesizer;
solo, that's you, and this step is the job. Write
`docs/milestones/1-<slug>.md` from it — the "why" record: what was
decided, what was traded off, what deviated from plan. The `1` is the
milestone's number (the M1), not its issue number. The slug is up to you, since the close looks for any
`docs/milestones/1-*.md`; the title in lower case with hyphens is the usual
choice, as in `docs/milestones/1-walking-skeleton.md`.

Then add the milestone's row to `ROADMAP.md`, already Done and linking the
document. `init` wrote the header; the milestone's title goes in its Goal
column:

```markdown
| Milestone | Goal | Tracking issue | Status |
|-----------|------|----------------|--------|
| M1 | Walking skeleton | [#1](https://github.com/<you>/my-project/issues/1) | [Done](docs/milestones/1-walking-skeleton.md) |
```

Land both on the light path the protocol gives the record ([SPEC
§4](../SPEC.md#housekeeping)): no task this time — a `docs:` commit
referencing the milestone issue, a PR whose body names that issue without a
closing keyword (the close closes it, not the merge), a clean-context review
as in step 5, then your operator confirmation comment on the PR — you are
the author and there is no one independent to approve it — and a rebase
merge. `task finish --operator-confirm` is not involved: there is no task
for it to finish, so you post the comment yourself, in the form it posts:

```markdown
**Operator confirmation:** reviewed and accepted by @<your-login> as both author and operator (pure solo tier, SPEC §5) — no independent principal exists in this project.
```

Nothing checks that comment's wording; it is the record that you looked
before merging. Then merge it yourself:

```sh
gh pr merge 4 --rebase --delete-branch   # your record PR's number
```

The roadmap lists finished milestones; `gh codecrew status` reports the
open one.

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
   role's contract from `.codecrew/roles/` (plus any
   `.codecrew/roles/<role>.local.md` beside
   it — `gh codecrew roles show <role>` prints the two composed) — so a fresh-context QA reading `.codecrew/roles/qa.md`
   probes what
   the implementer's context wouldn't. Or cross model families:
   dispatch another LLM through its own CLI for the reviewer or qa seat.
   Same identities, same commands — just different eyes.
2. **Give crew members their own identities.** When you want the record to
   show *which* agent did what — and GitHub itself to enforce that the
   approver isn't the author — route roles to other humans
   (`identity: user:<login>`), to a GitHub team
   (`identity: team:<org>/<slug>`, any member holds the seat), or to GitHub
   App identities (`identity: app:<slug>`): [identities.md](identities.md).
   The protocol doesn't change; only the routing table does.
3. **Full orchestration platforms** — an orchestrator dispatching the whole
   crew against the routing table, webhooks instead of polling. Run end to
   end on Paperclip: three milestones on a proving-ground repo, the third
   with one gate and no operator touch on the workflow besides it, then a
   fourth cycle on a fresh repo driven by a dedicated coordinator agent —
   the findings and what they changed are on
   [#119](https://github.com/radiusred/gh-codecrew/issues/119) and
   [#164](https://github.com/radiusred/gh-codecrew/issues/164).
   [Platform interop](platform-interop.md) is the page written from them:
   the coordinator seat, credentials, wake paths, the onboarding checklist,
   what four cycles cost, and the recipe.

## Tell us how it went

If you've worked through this guide, or part of it, we'd love to hear how it
went: what you tried, how far you got, and anything that stopped you or
surprised you. Getting stuck is just as useful to hear about as finishing.
There's a short form under
[I tried it](https://github.com/radiusred/gh-codecrew/discussions/categories/i-tried-it)
in the project's Discussions, and every report is read.
