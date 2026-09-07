# Working offline

CodeCrew keeps its state in GitHub on purpose: the plan, the decisions, the
gates and the verdicts are issues and pull requests, not files (SPEC §1). So
the verbs that transact over that state need GitHub, and there is no local
write queue that would replay them later — a queue is scar tissue, and under
latest-wins it is a consistency hazard. What is left is still most of the
day's work: reading the contract, planning, branching, implementing, running
the tests. This page says which verbs run with no network, which wait and
exactly what each one prints when it cannot reach GitHub, and the recipe for
work that was begun offline.

Everything below was read off the CLI, not off its intentions: every quoted
line comes from running the binary in a network namespace with no route out
(`unshare -rn`), from a hub, from a spoke, from a fresh repository and from a
protocol 1.x one.

## What runs with no network

| Verb | Offline behaviour |
| --- | --- |
| `codecrew version` | Prints the release and protocol version. Both are stamped into the binary at build time; nothing is fetched. |
| `codecrew help`, and `--help` on any verb | Prints usage and exits 0. `--help` is read before the verb runs, so it never reaches a gate or a fetch. |
| `codecrew roles show <role>` *(in a hub)* | Composes the contract from disk: the hub's `.codecrew/roles/<role>.md`, then its `.local.md` extension. A hub reads its own pointer and its own contracts, so nothing is fetched (SPEC §6). |
| `codecrew roles show <role> --latest` | Prints the contract embedded in the binary. Works anywhere, hub or spoke. |
| `codecrew roles diff <role>` *(in a hub)* | Compares the local contract with the embedded one. Both sides are on this machine. |
| `codecrew init` | Scaffolds the files and makes the one local commit. See the note below: it completes, with two `note:` lines. |
| `codecrew migrate` | Moves the repo to the 2.0 layout and commits, locally. Completes unless a 1.x identity has to be typed; see below. |

And everything git-local, which is where most of a task's time goes: reading
`.codecrew/AGENTS.md` and the role contract, drafting the plan in a file,
creating a branch, committing, running the tests, rebasing.

The `gh` floor check is local too. Every verb that reads the pointer asks
`gh --version` first and refuses `GH_TOO_OLD` below 2.50.0 — that call runs
no HTTP, so an offline machine is never told its `gh` is unreadable.

**`init` offline.** Its work is files plus one commit, and both are local; the
two things it asks GitHub are not, so both degrade to notes and the scaffold
stands:

```
note: could not ask GitHub whether the default branch requires pull requests, so the commit is on codecrew-bootstrap — git push -u origin codecrew-bootstrap, then open a PR (or merge it locally if the branch is unprotected)
note: could not ask GitHub which repository this is (…) — the cc: labels are created on first use instead, with GitHub's own colour
```

The first is the branch-protection probe: unanswered, `init` takes the
cautious branch and commits on `codecrew-bootstrap` rather than on the
default branch. (A repository with no `origin` at all is a different case —
that is *known* not to require pull requests, and the commit lands on the
current branch.) The second is the label step, which is the one thing the
scaffold cannot write to disk; the `cc:` labels then get created implicitly,
with GitHub's own colour, the first time a verb applies one. Running
`gh codecrew init` again once you are online creates them properly.

**`migrate` offline.** Same shape, and for the same reason: the migration is
a local move and one local commit, so no remote failure stands in its way.
It reports the same label note and exits 0. There is one exception, and it
refuses rather than guesses: a 1.x routing table that names a bare login
(`identity: someuser`) cannot be moved to the 2.0 grammar without asking
GitHub whether that login is a human or an App, so `migrate` refuses
`GH_UNREACHABLE` **before writing anything** and the whole migration waits.
A 1.x table whose identities already carry `user:`, `app:` or a `team`'s
slash needs no lookup and migrates with no network.

## What waits, and what it prints

Every verb that transacts with the record goes through one load path: read
the pointer, check the `gh` floor, then ask GitHub which repository this is.
That last call is where an offline run stops — before any write, on every one
of these verbs, with the same line on stderr and exit status 1:

```
codecrew: refused[GH_UNREACHABLE]: GitHub could not be reached (gh repo: error connecting to api.github.com
check your internet connection or https://githubstatus.com) — check the network and that gh is authenticated (gh auth status), or mint the seat's token with gh codecrew identity token <slug>; codecrew version, help, and roles show/diff in a hub need no network (SPEC §6)
```

The parenthesis is `gh`'s own message, so the wording there varies with how
the network is failing — no route, no DNS, no credentials all land here. The
code does not: `GH_UNREACHABLE` is raised for exactly this condition and is
never folded into another one, which is what lets an orchestrator tell
"offline" apart from "that issue does not exist" (SPEC §6, §10).

| Verb | Where it stops |
| --- | --- |
| `codecrew status` | The load. Nothing of the board is printed. |
| `codecrew role <name>` | The load — the routing answer is local in a hub, but the verb still asks GitHub to name the current repository first. |
| `codecrew task new/start/finish` | The load, `--dry-run` included. |
| `codecrew milestone new/close/evidence` | The load, `--dry-run` included. |
| `codecrew checkpoint` | The load. The gate is not raised; nothing local records it either. |
| `codecrew roles show <role>` *(from a spoke)* | The hub's contract has to be fetched, so it refuses `GH_UNREACHABLE` rather than reporting a contract that is merely elsewhere. |

Two conditions in this table are not about the network at all, and reading
the code rather than the symptom saves the confusion: `roles diff` from a
spoke fails with `no local .codecrew/roles/implementer.md — run from the hub
(spokes hold no contracts)` whether or not you are online, and a `--dry-run`
is a dry run against GitHub's *current* state, so it is a read and needs the
network as much as the write does.

**The identity verbs are the exception worth knowing.** `identity token` and
`identity webhook` talk to `api.github.com` directly rather than through
`gh`, so they never reach the refusal above. Offline they exit 1 with the
transport error alone and no refusal code:

```
codecrew: Get "https://api.github.com/app/installations?per_page=100": dial tcp: lookup api.github.com: Temporary failure in name resolution
```

Their *local* refusals are unaffected: a missing key or an unusable one still
refuses `NO_CREDENTIALS` or `BAD_CREDENTIALS` with no network, because
resolving the credential and signing the App JWT happen on this machine. Only
the two calls after that — discovering the installation and exchanging the
JWT for the one-hour installation token — need GitHub. A token you already
hold keeps working until it expires, and it cannot be renewed offline.

## The recipe: work begun offline

The boundaries of a task are GitHub's; the middle is yours. So:

**Offline.**

1. Read the contract — `gh codecrew roles show implementer` in a hub, or the
   copy in `.codecrew/roles/` — and `.codecrew/AGENTS.md`.
2. Draft the plan in a file. Be honest with yourself about what this costs:
   SPEC §4 wants the plan on the issue *before* the first commit, and
   offline you cannot put it there. Writing it first and posting it
   unchanged when you reconnect keeps the substance of the rule — deciding
   before doing — and posting a plan you have already implemented against,
   edited to match what you did, does not. `task start` refuses `NO_PLAN`
   until the section is on the issue either way.
3. Branch locally, with the name `task start` would create:
   `task/<issue number>-<slug of the title>`. The slug is the title
   lower-cased with every run of non-alphanumeric characters collapsed to a
   single `-`, trimmed to 40 characters.
4. Implement, commit, run the tests. Commit messages reference the task
   (`(#123)`) exactly as they would online.

**Back online.**

5. Put the plan on the task issue.
6. `gh codecrew task start <ref>` — it assigns the seat, verifies the plan,
   posts the start record that makes you the owner `task finish` will hold to,
   and creates the linked branch.
7. Reconcile your local branch with the one it just created, then push, open
   the pull request, and carry on.

**What step 6 actually does to a branch you already have.** `task start`
creates the linked branch through GitHub (`gh issue develop`), which knows
nothing about your machine. A local branch of that name is not consulted, and
the branch GitHub creates is cut from the default branch's current head and
carries none of your commits. The verb reports success and prints the two
lines it always prints:

```
linked branch task/329-working-offline-what-runs-what-waits-and created
locally: git fetch && git switch task/329-working-offline-what-runs-what-waits-and
```

Follow that second line literally and nothing happens, which is the part to
watch for: `git switch` finds the local branch you already have and stays on
it — "Already on …" — rather than checking out the branch that was just
created, and because a local branch of the name existed, git set up no
tracking either. You are on your own commits, with no upstream, next to a
remote branch of the same name that is linked to the issue and empty.

Reconciling the two is one rebase, and the push that follows sets the
upstream that is missing:

```sh
git fetch origin
git rebase origin/task/329-…            # your commits, replayed onto the linked branch
git push -u origin task/329-…
```

If you branched from the same commit the linked branch was cut from — the
common case, since both come from the default branch — the rebase is a no-op
and the push is a fast-forward. If the default branch moved on while you were
offline, the rebase is where you find that out, and it is the rebase you
would have done before opening the pull request anyway. Push under the same
name: the link between the issue and the branch is the name, so a branch
pushed as something else is not the linked one.

Two things not to do. Do not `git switch` to the fetched branch under a new
name and cherry-pick onto it — that is the same commits twice, and a tangle
to unpick at review. And run `task start` before you push, not after: the
linked branch is the one `gh issue develop` creates, and what it does with a
branch already sitting on the remote under that name, put there by something
else, is not a question this page has tested. Taking the steps in the order
above never asks it.

Re-running `task start` after the linked branch exists is safe on GitHub's
side — the call is idempotent, it answers with the existing branch, creates
no second one and moves nothing — but the verb posts another start record
every time it runs, since the latest record is what names the task's owner
(SPEC §4). Run it once.

None of this is `task start` being clever, because today it is not:
[#324](https://github.com/radiusred/gh-codecrew/issues/324) is the capture
for making it reconcile a branch begun offline — assign, verify the plan and
link the branch that is already there rather than create a second — and until
that ships, the rebase above is the step you do yourself. This page describes
what the CLI does today, and promises nothing of #324.

## Why there is no offline mode

Two things follow from the protocol's own shape, and both are deliberate
(SPEC §1, §4).

The record is the source of truth, and it is shared. A decision comment
written to a local queue is a decision nobody can see, react to or gate on
until it syncs; a gate raised offline is not raised. Under latest-wins,
replaying a queue after a conflicting change on GitHub silently overwrites
somebody's answer. So the verbs refuse loudly instead, with a code an
orchestrator can branch on, and the operator decides what to do about it.

And there is no local mirror to fall back to. A hub reads its own pointer and
its own contracts from disk — which is why `roles show` and `roles diff` work
there — but the issues, the labels, the reviews and the checks have no
on-disk representation to consult. What CodeCrew keeps on your machine is
what git already keeps: the code, the branch, and the contracts.
