# CodeCrew CLI reference

Every `gh codecrew` verb: what it takes, what it reads, what it writes, what
it prints, what it refuses, and what it exits with. The protocol these verbs
implement is [SPEC.md](SPEC.md); the refusal-code catalogue is
[SPEC §10](SPEC.md#10-the-cli), which every code below links to.

Every verb is safe to run by any role from any repo in the project: the
pointer file resolves the hub.

- [Invocation](#invocation)
- [Exit status](#exit-status)
- [Output channels](#output-channels)
- [Common refusals](#common-refusals)
- Verbs: [`init`](#gh-codecrew-init) · [`migrate`](#gh-codecrew-migrate) ·
  [`status`](#gh-codecrew-status) ·
  [`milestone new`](#gh-codecrew-milestone-new) ·
  [`milestone close`](#gh-codecrew-milestone-close) ·
  [`milestone evidence`](#gh-codecrew-milestone-evidence) ·
  [`task new`](#gh-codecrew-task-new) · [`task start`](#gh-codecrew-task-start) ·
  [`task finish`](#gh-codecrew-task-finish) ·
  [`checkpoint`](#gh-codecrew-checkpoint) · [`role`](#gh-codecrew-role) ·
  [`roles diff`](#gh-codecrew-roles-diff) ·
  [`roles show`](#gh-codecrew-roles-show) ·
  [`identity new`](#gh-codecrew-identity-new) ·
  [`identity webhook`](#gh-codecrew-identity-webhook) ·
  [`identity token`](#gh-codecrew-identity-token) ·
  [`version`](#gh-codecrew-version)

## Invocation

```
gh codecrew <verb> [<argument>…] [<option>…]
```

Installed as a `gh` extension the command is `gh codecrew`; the binary built
from source takes the same arguments without the `gh` prefix. `gh` 2.50.0 or
later must be on `PATH`.

`--help` or `-h` anywhere in the arguments prints the verb list on stderr and
exits 0, before the verb runs. `gh codecrew help` prints the same list on
stdout. `--version` and `-v` are `version`. An unknown verb or subcommand
prints the verb list on stderr and exits 1.

An issue reference — the `<ref>` argument of `task start`, `task finish` and
`checkpoint` — is `N`, `#N` or `owner/repo#N`; a bare number resolves against
the repository the verb runs in. A milestone number is `<n>`, and
`milestone close` also accepts `M<n>`. Where a verb takes options as well as
a leading argument, the argument may be given before or after them;
`roles diff` and `roles show` want the role name first, and `task start` and
`milestone evidence` take their argument and nothing else.

## Exit status

`0` when the verb did what it was asked, `1` on every failure — a refused
gate, an unusable flag, an unreachable GitHub, a `gh` that failed. There is
no exit-code taxonomy: `1` says only "this did not happen". `--dry-run`
follows the same rule and exits `1` when it reports a gate that would refuse.

## Output channels

**stdout** carries what a caller consumes: a minted token, `role <name>`'s
identity, a report and its `warning:` lines.

**stderr** carries the machine channel, one line:

```
codecrew: refused[CODE]: detail
```

The code is a fixed vocabulary, catalogued in [SPEC §10](SPEC.md#10-the-cli);
the detail is prose for a human and may be reworded in any release. A `note:`
line is the other thing stderr carries: advisory, never a failure, printed
alongside a verb that goes on to succeed.

## Common refusals

Every verb that loads the working repo's `.codecrew/config.yml` can raise
these, before its own work begins. The verbs that load no pointer are
`version`, `help`, `init`, `migrate`, `identity token` and
`identity webhook`; `init` and `migrate` raise the layout codes themselves.

| Code | Condition |
|------|-----------|
| [`LAYOUT_LEGACY`](SPEC.md#10-the-cli) | The repo is still on the protocol 1.x layout: a root `.codecrew.yml`, or a root `roles/` holding one of the five contracts, with no `.codecrew/config.yml`. The detail names what was found and `gh codecrew migrate`. |
| [`PROTOCOL_MISMATCH`](SPEC.md#10-the-cli) | The pointer's protocol major differs from the one this binary implements — the local pointer, or the hub's on a spoke's fetch. |
| [`IDENTITY_UNTYPED`](SPEC.md#10-the-cli) | A routing row's `identity` carries no type prefix; the detail names the row and the four forms. |
| [`SPOKE_ROUTING`](SPEC.md#10-the-cli) | A spoke's pointer carries a `roles:` block; the hub carries the one routing table. |
| [`HUB_UNREADABLE`](SPEC.md#10-the-cli) | In a spoke: the hub's pointer could not be fetched or parsed, so no role resolves. Routing fails closed. |
| [`GH_TOO_OLD`](SPEC.md#10-the-cli) | The installed `gh` is below 2.50.0, checked once, up front. |
| [`GH_UNREACHABLE`](SPEC.md#10-the-cli) | `gh` never reached GitHub — no route, no DNS, no credentials. Never folded into another condition. |

A hub reads its own pointer from disk and resolves roles with no network.
`version`, `help`, and `roles show` / `roles diff` in a hub need no network
at all.

## `gh codecrew init`

```
gh codecrew init [--hub owner/repo]
```

Scaffolds a repository into a CodeCrew project and creates the protocol's
labels. Hub mode — the default — writes the pointer with the full `~`-routed
roles table, a `ROADMAP.md` seed, the role contracts embedded in the
installed release under `.codecrew/roles/`, and a blank
`.codecrew/roles/<role>.local.md` extension beside each. Spoke mode writes
the pointer alone of those. Both modes write the entry point,
`.codecrew/AGENTS.md`, and the two root files that reach it: an `AGENTS.md`
naming the path and importing it, and a `CLAUDE.md` importing `AGENTS.md`.
Existing files are kept and reported. It then commits exactly the files it
wrote and ensures the three `cc:` labels exist. Rerunning it writes no file
it already wrote, commits nothing when it wrote nothing, and creates no label
it already created.

**Options**

| Option | Argument | Default | Effect |
|--------|----------|---------|--------|
| `--hub` | `owner/repo`, or `self` | `self` | Spoke mode: the pointer names that hub and nothing else is scaffolded. `self` is hub mode. |

**Reads.** The working directory, which must be the repository root; the
repository's default branch and its pull-request requirement, through `gh`;
the repository's existing labels.

**Writes.** `.codecrew/config.yml`, `.codecrew/AGENTS.md`, `AGENTS.md`,
`CLAUDE.md`, and in hub mode `ROADMAP.md`, `.codecrew/roles/<role>.md` and
`.codecrew/roles/<role>.local.md` — each only when it is absent. One commit
of exactly those paths, on the current branch, or on `codecrew-bootstrap` cut
from the default branch when the default branch requires pull requests; it
never pushes. In GitHub: the missing `cc:task`, `cc:milestone` and
`cc:needs-decision` labels, with the protocol's colour and description; an
existing label is left exactly as it is. On stdout: a line per file written
or kept, the commit, a line per label, and last an `action needed` block
naming any kept root entry point that does not reach `.codecrew/AGENTS.md`,
with the lines to add to it.

**Refusals.** [`LAYOUT_LEGACY`](SPEC.md#10-the-cli) — the repo is on the 1.x
layout; `gh codecrew migrate` moves it. A directory that is not the
repository root, and a scaffold that cannot be written, exit 1 without a
code. A label step that cannot reach GitHub is a `note:`, not a refusal, and
the scaffold stands; so is a directory that is not yet a git repository, and
the labels arrive on the rerun once it is one.

**Exit status.** 0 when the scaffold is complete, including a rerun that
writes nothing; 1 on a refusal or an unusable argument.

```
gh codecrew init                              # this repo is the hub
gh codecrew init --hub myorg/platform         # a spoke pointing at that hub
```

## `gh codecrew migrate`

```
gh codecrew migrate [--dry-run]
```

Moves a repository from the protocol 1.x layout to 2.0. `git mv` takes
`.codecrew.yml` to `.codecrew/config.yml` and, out of a root `roles/`, the
five role contracts and their `<role>.local.md` extensions into
`.codecrew/roles/`, removing the emptied directory. The pointer is then
rewritten in place: `codecrew: "2.0"`, a `coordinator` row when the table
declares none, and every identity typed by asking GitHub what each bare 1.0
login is. Comments, blank lines and key order are kept. It writes
`.codecrew/AGENTS.md` when the repo has none, and does for the root
`AGENTS.md` and `CLAUDE.md` what `init` does. Every refusal is raised before
anything is written. A repo already on the 2.0 layout says so, moves nothing
and commits nothing — the label step still runs.

**Options**

| Option | Argument | Default | Effect |
|--------|----------|---------|--------|
| `--dry-run` | none | off | Prints every file step and every label creation and restyle; writes nothing. |

**Reads.** The working directory, which must be the repository root of a git
repository; `.codecrew.yml`; a root `roles/`, and only when it already holds
one of the ten CodeCrew names; GitHub, for the account behind each bare 1.0
identity (`users/<login>`, then `users/<login>[bot]`).

**Writes.** The moved and rewritten files above, in one commit of exactly
those paths on the current branch; it never pushes. In GitHub: the missing
`cc:` labels created, and the existing ones restyled to the protocol's colour
and description, each reported; every other label untouched. On stdout: a
line per step, and last an `action needed` block for a kept root entry point
that does not reach the instructions.

**`--dry-run`.** The same steps in the same order, the same refusals, and the
same label lines; nothing written.

**Refusals.** [`BOTH_LAYOUTS`](SPEC.md#10-the-cli) — the 1.x and 2.0 layouts
overlap. [`FOREIGN_ROLES_DIR`](SPEC.md#10-the-cli) — a CodeCrew-owned root
`roles/` holds an entry outside the ten names.
[`MIGRATION_UNSUPPORTED`](SPEC.md#10-the-cli) — the pointer's protocol major
is not 1. [`SPOKE_ROUTING`](SPEC.md#10-the-cli) — a 1.x spoke's pointer
carries a `roles:` block. [`IDENTITY_UNRESOLVED`](SPEC.md#10-the-cli) — a
bare identity answers to nothing, to both a user and an App, or to an
organization. [`GH_UNREACHABLE`](SPEC.md#10-the-cli) — the lookup could not
reach GitHub at all. A directory that is not a git repository, or not its
root, exits 1 without a code. A label step that cannot reach GitHub is a
`note:`; the rerun completes it.

**Exit status.** 0 when the repo is on the 2.0 layout, including a rerun that
moves nothing; 1 on a refusal. `--dry-run` exits with the first refusal it
reports.

```
gh codecrew migrate --dry-run    # every step, nothing written
gh codecrew migrate              # the move, in one local commit
```

## `gh codecrew status`

```
gh codecrew status
```

Reports where the project is: the open milestones and their tasks, each
task's inferred state, and every raised gate — on tasks and on milestone
issues alike, the latter marked `(milestone)`. A task in progress or in
review names its holder, taken from the login in its latest `**Started by**`
record and falling back to the first assignee only when nothing records a
start. It then reports the stale task branches of the repository it runs in —
every `task/<n>-…` branch on the remote whose task issue is closed — with the
delete-or-keep verdict the next `milestone close` would give it. With no open
milestone it says so in place of the board and the gates.

**Options.** None.

**Reads.** The pointer and, from a spoke, the hub's; the hub's open
milestones, their sub-issues, labels and comments; the running repository's
remote `task/<n>-…` branches and one issue per branch; the repository's
delete-branch-on-merge setting; the local `.codecrew/roles/` contracts, for
the drift note.

**Writes.** Nothing. The report goes to stdout, `note:` lines to stderr.

**Refusals.** The [common](#common-refusals) ones. Nothing in the report
refuses: an ID that is not the milestone's own prints as a line and the
report carries on, a branch whose task cannot be read is a `note:` and is
left, a listing that held more than one page says it is partial, and a
repository that could not be read at all replaces the branch report with a
`note:` — the report's silence means clean and nothing else.

**Exit status.** 0 when the report printed; 1 on a common refusal.

```
gh codecrew status
```

## `gh codecrew milestone new`

```
gh codecrew milestone new --title T [--goal G]
                          [--requirement R]...
                          [--dry-run]
```

Creates a milestone tracking issue in the hub from the template, labelled
`cc:milestone` and titled `M<n>: <title>`. Each `--requirement` becomes a
bold-ID line under `## Requirements`, numbered `M<n>-R1`, `R2`, … in the
order given — the section the close gate reads. The number `n` is derived
twice: before creating, as one past the highest `M<k>:` title across the
hub's label-filtered milestone listing and its newest unfiltered issues; and
after creating, when both listings are read again and the number must be the
new issue's alone. Another issue already carrying the prefix has the new
issue renumbered to the next free number, title and `M<n>-R<k>` IDs, printed
as a `renumbered:` line.

**Options**

| Option | Argument | Default | Effect |
|--------|----------|---------|--------|
| `--title` | text | none; required | The milestone's title. A title carrying an `M<k>` prefix that disagrees with the derived number is refused; one that agrees is stripped. |
| `--goal` | text | `_To be written._` | The `## Goal` paragraph. |
| `--requirement` | text | none | One requirement, repeatable. The CLI numbers them; text that brings its own ID is refused. |
| `--dry-run` | none | off | Prints the number, title and requirement IDs the milestone would get; creates nothing. |

**Reads.** The pointer and, from a spoke, the hub's; the hub's milestone
listing and its newest issues.

**Writes.** One issue in the hub, labelled `cc:milestone`. No file: the
milestone's `ROADMAP.md` row is added, Done, by its document PR. On stdout:
the issue reference, the title, the requirement IDs counted, and any
`renumbered:` line.

**`--dry-run`.** The number, title and requirement IDs; no issue.

**Refusals.** [`MILESTONE_NUMBER_TAKEN`](SPEC.md#10-the-cli) — the issue was
created and another milestone holds its `M<n>:` prefix, and the renumbering
did not settle it; the detail names both issues and the hand fix. A missing
`--title`, a title whose own `M<k>` prefix disagrees, and a `--requirement`
carrying an ID exit 1 without a code. Plus the [common](#common-refusals)
refusals.

**Exit status.** 0 when the issue exists as described; 1 otherwise.

```
gh codecrew milestone new --title "Greeting service" --dry-run
gh codecrew milestone new --title "Greeting service" --goal "…" \
  --requirement "visiting the app's URL returns a greeting" \
  --requirement "the greeting is configurable"
```

## `gh codecrew milestone close`

```
gh codecrew milestone close <milestone number>
                            [--dry-run]
```

Closes a milestone once every gate passes, then sweeps the task branches it
leaves behind. The gates, in order: the milestone issue carries no
`cc:needs-decision`; every task under it is closed; its `## Requirements`
section declares at least one bold requirement ID; every ID it declares is
the milestone's own; every requirement's latest QA verdict is `satisfied` —
only verdicts from the qa role's holder count, and the latest comment
carrying a verdict for an ID wins; and the milestone document is on the
default branch. It then deletes each task branch whose PR merged and which
still sits at the merged commit, or which has no open PR and carries nothing
beyond the default branch, reporting every branch it keeps and why. A second
pass reaches the branches earlier closes left behind, across the hub and
every repo the milestone's tasks name, and names them in the closing comment
under `Swept from earlier closes:`. It gathers the milestone's
Decision and Deviation comments as raw material for the doc-synthesizer.

**Options**

| Option | Argument | Default | Effect |
|--------|----------|---------|--------|
| `--dry-run` | none | off | Prints every gate in order with its outcome, then every branch either sweep would delete or keep and why, then the closing comment; writes nothing; exits with the first refusal's code. |

**Reads.** The pointer and, from a spoke, the hub's; the milestone issue, its
sub-issues, their labels and comments; `docs/milestones/<n>-*.md` on the
hub's default branch; the `task/<n>-…` branches of the hub and of every repo
the tasks name, and the pull requests whose head each candidate branch is.

**Writes.** In GitHub: the closing comment on the milestone issue, the issue
closed, and the branches the two sweeps delete. Nothing on disk. The report
goes to stdout.

**`--dry-run`.** The same gates in the same order, the same sweep verdicts
and the same closing comment; nothing written.

**Refusals.** [`NOT_FOUND`](SPEC.md#10-the-cli) — no open milestone with that
number. [`MILESTONE_GATED`](SPEC.md#10-the-cli) — the milestone issue carries
`cc:needs-decision`. [`OPEN_TASKS`](SPEC.md#10-the-cli) — tasks under it are
still open. [`NO_REQUIREMENTS`](SPEC.md#10-the-cli) — the `## Requirements`
section declares no bold ID. [`REQUIREMENT_ID_MISMATCH`](SPEC.md#10-the-cli)
— an ID declared there is not the milestone's own; the grammar is
`M<milestone>-R<k>`. [`VERDICT_MISSING`](SPEC.md#10-the-cli) — a requirement
has no QA verdict from the qa seat's holder.
[`VERDICT_UNSATISFIED`](SPEC.md#10-the-cli) — the latest verdict on a
requirement is not `satisfied`. [`DOC_MISSING`](SPEC.md#10-the-cli) — no
`docs/milestones/<n>-*.md` on the default branch. Plus the
[common](#common-refusals) refusals. Nothing in the sweep refuses: a branch
the verb could not read is a `note:` and is left standing, and a repo
carrying more task branches than one listing holds is swept in part and says
so.

**Exit status.** 0 when the milestone is closed; 1 on the first refusal,
whose code it exits with under `--dry-run` too.

```
gh codecrew milestone close 4 --dry-run
gh codecrew milestone close 4
```

## `gh codecrew milestone evidence`

```
gh codecrew milestone evidence <milestone number>
```

Walks the milestone's record — the tracking issue and every sub-issue, bodies
and comments — and verifies that every citation resolves: github.com
references through the API under the caller's auth, everything else by HTTP.
A citation is a URL in prose or in a Markdown link outside code; a URL inside
an inline code span, a fenced block, or a block indented four columns where
it does not continue a paragraph is content and is not checked. It checks the
milestone's own requirement-ID grammar before the walk. The milestone
resolves whatever its state: a closed one resolves too and is reported as
closed before the citation report.

**Options.** None.

**Reads.** The pointer and, from a spoke, the hub's; the hub's milestone
listing in every state; the milestone issue and its sub-issues, bodies and
comments; each cited URL.

**Writes.** Nothing. The report goes to stdout, including a `warning:` line
per unresolved external citation.

**Refusals.** [`NOT_FOUND`](SPEC.md#10-the-cli) — no milestone carries that
number, open or closed. [`REQUIREMENT_ID_MISMATCH`](SPEC.md#10-the-cli) — an
ID under `## Requirements` is not the milestone's own, checked before the
walk. [`EVIDENCE_UNREACHABLE`](SPEC.md#10-the-cli) — a github.com citation
does not resolve. An external citation that does not resolve prints a
`warning:` and does not block. Plus the [common](#common-refusals) refusals.

**Exit status.** 0 when every github.com citation resolved; 1 when one did
not, or on any other refusal.

```
gh codecrew milestone evidence 4
```

## `gh codecrew task new`

```
gh codecrew task new --milestone N --title T
                     [--repo owner/repo] [--goal G] [--requirements IDs]
                     [--adopts <ref>[,<ref>]]
```

Creates a task issue from the template, labelled `cc:task`, in the spoke
named by `--repo` or in the current repository, and attaches it to the
milestone as a sub-issue. `--adopts` names the backlog captures the task
takes up: they are written under an `## Adopts` section of the body, each
capture is commented on to say this task carries it, and `task finish` closes
them after the merge. Every adopted ref must be an open issue and is checked
before anything is created, so a refusal leaves no half-adopted task;
duplicates collapse. The milestone resolves by number from the hub's
open-milestone listing, then from the hub's newest issues regardless of
label, then again after a short wait — three reads in all, and a milestone
found by either fallback is named in the output.

**Options**

| Option | Argument | Default | Effect |
|--------|----------|---------|--------|
| `--milestone` | number | none; required | The milestone the task links into. |
| `--title` | text | none; required | The task's title. |
| `--repo` | `owner/repo` | the current repository | The spoke the task issue is created in. |
| `--goal` | text | `_To be written._` | The `## Goal` paragraph. |
| `--requirements` | comma-separated IDs | `None directly.` | The `## Requirements` line: the milestone requirement IDs this task serves. |
| `--adopts` | `<ref>[,<ref>]` | none | A backlog issue this task adopts and `task finish` closes; repeatable, and each value may be a comma-separated list. A bare number resolves against the task's own repo. |

**Reads.** The pointer and, from a spoke, the hub's; the hub's milestone
listing and newest issues; each adopted issue's state.

**Writes.** In GitHub: one issue in the target repo labelled `cc:task`,
linked to the milestone as a sub-issue, and one comment on each adopted
capture. Nothing on disk. On stdout: the issue reference and the adoptions
recorded.

**Refusals.** [`NOT_FOUND`](SPEC.md#10-the-cli) — no open milestone with that
number, after all three reads. [`ADOPT_NOT_OPEN`](SPEC.md#10-the-cli) — an
adopted ref could not be read, or is already closed; checked before the task
is created. [`GH_UNREACHABLE`](SPEC.md#10-the-cli) — GitHub could not be
reached. Plus the [common](#common-refusals) refusals. A comment on a capture
that fails once the task exists is a `note:`; the body carries the link
either way.

**Exit status.** 0 when the task issue exists and is linked; 1 otherwise.

```
gh codecrew task new --milestone 4 --title "Serve the greeting" \
  --requirements M4-R1 --adopts 31
```

## `gh codecrew task start`

```
gh codecrew task start <ref>
```

Starts a task: it verifies that the issue is an open `cc:task` whose `## Plan`
section is not empty, posts the `**Started by** @<login>.` record that makes
the caller the task's owner — the record `task finish` holds to — assigns the
caller when the routing table types them as a human, and creates the linked
development branch `task/<n>-<slug>`. An `app:`-typed caller, and any login
carrying the `[bot]` suffix, is not assigned: GitHub does not accept a GitHub
App as an issue assignee, and the record is the fact. A caller whose role
routing resolves to a contract that forbids commits — `qa`, `reviewer` — gets
no branch.

**Options.** None.

**Reads.** The pointer and, from a spoke, the hub's routing table; the task
issue, its labels, its body's `## Plan` section; the caller's login.

**Writes.** In GitHub: the `**Started by**` comment on the task issue, the
assignment where the caller is human, and the linked branch. Nothing on disk
— the branch is created on the remote, and the printed line names the
`git fetch && git switch` that brings it local. On stdout: the branch
created, the local command, and the start record.

**Refusals.** [`CLOSED`](SPEC.md#10-the-cli) — the task issue is already
closed. [`NOT_A_TASK`](SPEC.md#10-the-cli) — the issue is not labelled
`cc:task`. [`NO_PLAN`](SPEC.md#10-the-cli) — the task's Plan section is
empty. Plus the [common](#common-refusals) refusals. A failed assignment is a
`note:`.

**Exit status.** 0 when the record is posted; 1 on a refusal.

```
gh codecrew task start 12
gh codecrew task start myorg/service#12
```

## `gh codecrew task finish`

```
gh codecrew task finish <ref> [--operator-confirm]
                        [--bypass]
                        [--dry-run]
```

The gatekeeper: it checks every gate, merges the closing pull request by
rebase, closes the task and the captures it adopted, deletes the head branch
and tidies the clone it runs in. The gates, in order: the task is open; it
carries no `cc:needs-decision`; every gate raised on it has a
`**Gate resolved:**` answer; the caller is the seat that started it — the
login in the latest `**Started by**` record, the same login with the `[bot]`
suffix ignored, or the same routed seat, a team-held role being any member; a
pull request closes the task; its CI checks exist and are green; the reviewer
role's holder has approved when that role routes to a distinct principal, and
otherwise any non-doer has; and GitHub's own required-review rule counts the
approval. Before the merge it prints a `note:` naming every issue the pull
request's closing references would close besides the task itself. After the
merge it closes each adopted capture with a comment naming the task, the pull
request and the merge commit, deletes the head branch, and — in a clone of
the task's own repo that holds the merged branch — fetches with prune,
switches to the default branch if that branch is checked out, fast-forwards
it, and deletes the local task branch. Nothing after the merge refuses.

**Options**

| Option | Argument | Default | Effect |
|--------|----------|---------|--------|
| `--operator-confirm` | none | off | Solo tier: records an explicit operator confirmation on the pull request in place of a non-doer approval. A crew identity is refused. |
| `--bypass` | none | off | Merges with the ruleset's administrator bypass when GitHub does not count the recorded approval, recorded as a comment on the pull request. For a human operator the ruleset lists as a bypass actor; a crew identity is refused. |
| `--dry-run` | none | off | Prints every gate in order — ok, refused with its code, not reached, not applicable — then the comments, merge, adopted captures, head deletion and local cleanup it would perform; writes nothing, the clone included; exits with the first refusal's code. |

**Reads.** The pointer and, from a spoke, the hub's routing table; the task
issue, its labels and comments; the pull requests that close it, their
checks, reviews and review decision; the closing references GitHub parsed
from the pull request body; the local clone's branches and remotes.

**Writes.** In GitHub: the operator-confirmation or bypass comment where one
applies, the rebase merge, the task closed, a comment on each adopted capture
and that capture closed, and the head branch deleted. On disk, and only in a
clone of the task's own repo: a fetch with prune, a switch to the default
branch, a fast-forward of it, and the deletion of the local task branch —
forced, and allowed only when the branch sits at the merge commit or is
contained in the fetched default branch. A branch carrying anything else is
named and kept. Run anywhere else, the local half does nothing and prints
nothing.

**`--dry-run`.** Every gate and every write above, printed and not performed;
the same `note:` about the pull request's other closing references; the same
exit code as the real run's first refusal.

**Refusals.** [`CLOSED`](SPEC.md#10-the-cli) — the task issue is already
closed. [`GATED`](SPEC.md#10-the-cli) — the task carries `cc:needs-decision`.
[`GATE_UNRECORDED`](SPEC.md#10-the-cli) — a gate was raised and the label
removed with no `**Gate resolved:**` comment answering it.
[`NOT_OWNER`](SPEC.md#10-the-cli) — the caller is not the seat that started
the task, or nothing records a start. [`NO_PR`](SPEC.md#10-the-cli) — no open
pull request closes the task.
[`NO_CHECKS_PERMISSION`](SPEC.md#10-the-cli) — the caller's installation
token cannot read the pull request's checks at all.
[`NO_CHECKS`](SPEC.md#10-the-cli) — the pull request reports no checks; there
is no override. [`CHECKS_PENDING`](SPEC.md#10-the-cli) — its checks are still
running. [`CHECKS_FAILING`](SPEC.md#10-the-cli) — a check failed.
[`NO_HOLDER_REVIEW`](SPEC.md#10-the-cli) — the reviewer seat's holder has not
approved. [`NO_NONDOER_APPROVAL`](SPEC.md#10-the-cli) — the reviewer seat is
operator-held and no non-author has approved.
[`REVIEW_NOT_COUNTED`](SPEC.md#10-the-cli) — the protocol's review gate
passed and GitHub's own required-review rule has not; the detail names the
supported paths. [`SELF_CONFIRM`](SPEC.md#10-the-cli) — `--operator-confirm`
was given by a crew identity. [`CREW_BYPASS`](SPEC.md#10-the-cli) —
`--bypass` was given by a crew identity. Plus the
[common](#common-refusals) refusals. After the merge nothing refuses: a
capture already closed is reported, one that cannot be closed is a `note:`,
and so is a failed deletion.

**Exit status.** 0 when the merge happened and the task is closed; 1 on the
first refusal, whose code `--dry-run` exits with too.

```
gh codecrew task finish 12 --dry-run
gh codecrew task finish 12
gh codecrew task finish 12 --operator-confirm   # solo tier, human operator
```

## `gh codecrew checkpoint`

```
gh codecrew checkpoint <ref> --question "..."
```

Raises a human gate: it posts the question as a comment and applies
`cc:needs-decision`, creating the label with the protocol's colour and
description if the repository does not define it yet. The ref is a task, or
the milestone issue when the question is about a requirement and no task
carries it; the comment and the receipt say which. The gate blocks
`task finish` until a human removes the label and a `**Gate resolved:**`
comment records the answer.

**Options**

| Option | Argument | Default | Effect |
|--------|----------|---------|--------|
| `--question` | text | none; required | The judgment call a human must make. Posted verbatim as the comment. |

**Reads.** The pointer and, from a spoke, the hub's; the referenced issue and
the repository's labels.

**Writes.** In GitHub: the question as a comment on the referenced issue, the
`cc:needs-decision` label created if absent, and the label applied. Nothing
on disk. On stdout: the receipt naming the issue and that the gate is raised.

**Refusals.** The [common](#common-refusals) ones. A missing `--question` or
`<ref>` exits 1 without a code. A label creation that fails is a `note:` and
the gate is still raised: applying an unknown label creates it implicitly,
and only the colour is lost.

**Exit status.** 0 when the gate is raised; 1 otherwise.

```
gh codecrew checkpoint 12 --question "which of the two schemas do we keep?"
```

## `gh codecrew role`

```
gh codecrew role <name> [--login]
```

Prints the typed identity holding a role — `app:<slug>`, `user:<login>`,
`team:<org>/<slug>`, or `~` for the operator — on one line. Script-consumable,
and correct from a pointer-only spoke, where the hub's routing table is
fetched at load. With `--login` it prints instead the handle GitHub will
accept a review request for, and nothing at all for an `app:` or `~` seat,
neither of which can be requested: the emptiness is the caller's whole
decision.

**Options**

| Option | Argument | Default | Effect |
|--------|----------|---------|--------|
| `--login` | none | off | Prints the review-requestable handle: the login of a `user:` seat, `<org>/<slug>` for a `team:` seat, and nothing for an `app:` or `~` seat. |

**Reads.** The pointer and, from a spoke, the hub's routing table.

**Writes.** Nothing. One line on stdout, or no line at all under `--login` for
an unrequestable seat.

**Refusals.** The [common](#common-refusals) ones — a hub that cannot be read
refuses rather than answering `~`. A missing or unknown role name exits 1
without a code.

**Exit status.** 0 when the role resolved, including the empty `--login`
output; 1 otherwise.

```
gh codecrew role reviewer            # app:myorg-checker
reviewer=$(gh codecrew role reviewer --login)
[ -n "$reviewer" ] && gh pr create --reviewer "$reviewer"
```

## `gh codecrew roles diff`

```
gh codecrew roles diff <role>
```

Shows how the project's `.codecrew/roles/<role>.md` differs from the contract
embedded in the installed CLI. A `.codecrew/roles/<role>.local.md` extension
is never drift. Contracts are the project's own fork: reconciliation is a
judgment routed through a task and a pull request, never an overwrite.

**Options.** None.

**Reads.** The pointer; the local `.codecrew/roles/<role>.md`; the contract
embedded in the binary.

**Writes.** Nothing. The diff goes to stdout.

**Refusals.** The [common](#common-refusals) ones. A role name the CLI does
not embed exits 1 without a code.

**Exit status.** 0 whether or not the contracts differ; 1 on a refusal.

```
gh codecrew roles diff implementer
```

## `gh codecrew roles show`

```
gh codecrew roles show <role> [--latest]
```

Prints the contract a dispatched session loads: the hub's
`.codecrew/roles/<role>.md` with its local extensions appended in order — the
hub's, then the spoke's. With `--latest` it prints the contract embedded in
the installed binary, whole and without extensions.

**Options**

| Option | Argument | Default | Effect |
|--------|----------|---------|--------|
| `--latest` | none | off | Prints the contract embedded in the installed binary instead of the composed one. |

**Reads.** The pointer; in a hub, `.codecrew/roles/<role>.md` and
`.codecrew/roles/<role>.local.md` from disk; from a spoke, the hub's contract
and extension from its default branch through GitHub, and the spoke's own
extension from disk.

**Writes.** Nothing. The contract goes to stdout.

**Refusals.** [`GH_UNREACHABLE`](SPEC.md#10-the-cli) — from a spoke, the
hub's contract could not be fetched. Plus the [common](#common-refusals)
refusals. In a hub the verb reads only local files and needs no network.

**Exit status.** 0 when the contract printed; 1 on a refusal.

```
gh codecrew roles show implementer
gh codecrew roles show implementer --latest
```

## `gh codecrew identity new`

```
gh codecrew identity new <role> --name N
                         [--owner O] [--with-webhook --webhook-url U]
                         [--with-approval-permission]
                         [--events E,E] [--webhook-secret S]
```

Mints the role's App identity through the GitHub App manifest flow: it builds
a manifest carrying the role's minimal permission set, serves a one-click
loopback URL for the operator to open, stores the private key GitHub returns,
writes the role's routing into the hub's `.codecrew/config.yml` in typed form
(`identity: app:<slug>`), and prints the steps that remain — installing the
App, which is per account, and the optional display polish. Webhooks are off
by default.

**Options**

| Option | Argument | Default | Effect |
|--------|----------|---------|--------|
| `--name` | App name | none; required | The App's name — a crew member (`myorg-coder`), not a role. |
| `--owner` | account | the hub's owner | The account that owns the App. |
| `--no-route` | none | off | Prints the routing step instead of writing it into the hub's pointer. |
| `--with-webhook` | none | off | Subscribes the App to protocol traffic; requires `--webhook-url`. |
| `--webhook-url` | URL | none | The receiver for `--with-webhook` deliveries. |
| `--with-approval-permission` | none | off | Reviewer only: grants `contents: write`, so the App's approvals satisfy GitHub's required-review rules. |
| `--events` | comma-separated events | `pull_request,pull_request_review` | With `--with-webhook`: the events to subscribe, validated against the role's permissions. |
| `--webhook-secret` | secret | none | With `--with-webhook`: sets the hook secret to the receiver's as soon as the App exists. |

**Reads.** The pointer and the hub it names; the owner account's type,
through GitHub; the manifest conversion GitHub returns at the end of the
browser flow.

**Writes.** On disk: `~/.config/codecrew/<slug>.<date>.private-key.pem` and
the credential stub `~/.config/codecrew/<slug>.json`; the hub's
`.codecrew/config.yml` routing row, uncommitted, unless `--no-route`. In
GitHub: the App itself, and its hook secret when `--webhook-secret` is given.
On stdout: the loopback URL to open, the stored paths, and the numbered steps
that remain.

**Refusals.** The [common](#common-refusals) ones. A missing `<role>` or
`--name`, a manifest the flags do not make valid, `--webhook-secret` without
`--with-webhook`, and a browser flow that does not complete within the hour
exit 1 without a code.

**Exit status.** 0 when the App exists and its key is stored; 1 otherwise.

```
gh codecrew identity new implementer --name myorg-coder
gh codecrew identity new reviewer --name myorg-checker --with-approval-permission
```

## `gh codecrew identity webhook`

```
gh codecrew identity webhook <slug> [--show]
                             [--url U] [--secret S|--rotate-secret]
```

Reads and sets an existing App's own hook, under the App's own key: `--show`
prints the URL, the content type, whether a secret is set and the App's
subscribed events; `--url` and `--secret` set them; `--rotate-secret` mints a
new secret, sets it and prints it once. An App hook covers every repository
its installation sees. Event subscriptions cannot be set after creation —
GitHub offers no endpoint — and the verb prints the settings page where they
are ticked.

**Options**

| Option | Argument | Default | Effect |
|--------|----------|---------|--------|
| `--show` | none | off | Prints the hook's URL, content type, whether a secret is set, and the subscribed events. |
| `--url` | URL | none | Sets the receiver URL. |
| `--secret` | secret | none | Sets the hook secret to the receiver's. Never stored, never printed. Exclusive with `--rotate-secret`. |
| `--rotate-secret` | none | off | Generates a new secret, sets it, and prints it once. |
| `--events` | comma-separated events | none | Not settable after creation: the verb prints where they change, and what is subscribed now. |

**Reads.** The App's private key and credential stub — the environment
bindings first (`GITHUB_APP_ID` or `GITHUB_CLIENT_ID`, and `GITHUB_PRIVATE_KEY`
or `GITHUB_PEM` as PEM text or a path), else `~/.config/codecrew/<slug>`; the
App's own record and hook configuration from GitHub. It reads no pointer and
runs from anywhere.

**Writes.** In GitHub: the App's hook URL and secret, where those flags are
given. Nothing on disk. On stdout: the hook's configuration, and a rotated
secret once.

**Refusals.** [`NO_CREDENTIALS`](SPEC.md#10-the-cli) — nothing to sign with:
no App id and key in the environment, and no key and stub under
`~/.config/codecrew/`. [`BAD_CREDENTIALS`](SPEC.md#10-the-cli) — GitHub
rejected the App JWT, or knows no App by the id it was signed as.
[`NO_WEBHOOK`](SPEC.md#10-the-cli) — the App was minted without a webhook and
GitHub's API cannot create one; the detail names the settings page. `--secret`
with `--rotate-secret`, and `--events`, exit 1 without a code.

**Exit status.** 0 when the hook was printed or set; 1 otherwise.

```
gh codecrew identity webhook myorg-coder --show
gh codecrew identity webhook myorg-coder --url https://example.test/hook --rotate-secret
```

## `gh codecrew identity token`

```
gh codecrew identity token [<slug>]
                           [--installation ID]
```

Mints a short-lived installation token as the App and prints it alone on
stdout, with a receipt on stderr naming the slug, App id, installation and
account it minted for. Credentials come from the environment first, under the
names platforms bind, and otherwise from the key and stub under
`~/.config/codecrew/` for the slug. It signs the App JWT, discovers the
installation from the App itself — one installation is taken, several narrow
to the hub's owner — and never writes `gh`'s config. The `<slug>` argument
selects the key and stub under `~/.config/codecrew/` and is optional when the
environment binds an App id and key. It reads no pointer and runs from
anywhere.

**Options**

| Option | Argument | Default | Effect |
|--------|----------|---------|--------|
| `--installation` | id | `GITHUB_INSTALLATION_ID`, else none | A hint, used only when the App can see that installation. |

**Reads.** `GITHUB_APP_ID` or `GITHUB_CLIENT_ID`, and `GITHUB_PRIVATE_KEY` or
`GITHUB_PEM` (PEM text or a path), then `GITHUB_INSTALLATION_ID`; failing
those, `~/.config/codecrew/<slug>.json` and its private key. GitHub, for the
App's installations and the token. The pointer, only to learn the hub's owner
for narrowing, and never as a requirement.

**Writes.** Nothing on disk, and nothing in GitHub beyond the token GitHub
mints. The token alone on stdout; the receipt on stderr.

**Refusals.** [`NO_CREDENTIALS`](SPEC.md#10-the-cli) — no App id and key in
the environment, and no key and stub under `~/.config/codecrew/`.
[`BAD_CREDENTIALS`](SPEC.md#10-the-cli) — GitHub rejected the App JWT, or
knows no App by the id it was signed as.
[`NO_INSTALLATION`](SPEC.md#10-the-cli) — the App is installed on no account
its key can see. [`INSTALLATION_AMBIGUOUS`](SPEC.md#10-the-cli) — the App is
installed on several accounts and nothing selects one.

**Exit status.** 0 when a token was printed; 1 otherwise.

```
export GH_TOKEN=$(gh codecrew identity token myorg-coder)
gh codecrew identity token myorg-coder --installation 12345678
```

## `gh codecrew version`

```
gh codecrew version
```

Prints the installed release tag and the protocol version it implements, as
`<tag> (protocol <major.minor>)`. A build from source reports `dev`. The
shape is stable within a major.

**Options.** None. `--version` and `-v` are the same verb.

**Reads.** Nothing: no pointer, no network, no files.

**Writes.** Nothing. One line on stdout.

**Refusals.** None.

**Exit status.** 0.

```
gh codecrew version     # v2.0.1 (protocol 2.0)
```

Licensed under [Apache 2.0](LICENSE).
