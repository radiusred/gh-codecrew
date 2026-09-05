# M12: v1.2.0 and the field fixes behind it

Tracking issue: [#241](https://github.com/radiusred/gh-codecrew/issues/241) ·
Synthesized 2026-09-05 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #241's **four** requirements — all four present when the
milestone opened, **none amended and none added** (the issue's
`userContentEdits` query returns a count of zero and no `lastEditedAt`) — its
Gates section and its **one** QA round; the four delivery task issues, all in
this repository
([#242](https://github.com/radiusred/gh-codecrew/issues/242),
[#243](https://github.com/radiusred/gh-codecrew/issues/243),
[#244](https://github.com/radiusred/gh-codecrew/issues/244),
[#245](https://github.com/radiusred/gh-codecrew/issues/245)); their four merged
PRs; the annotated tag `v1.2.0`, the release run and the release it published;
the **six Decision comments and one Deviation comment** across the milestone
issue and those task issues; the **nine review submissions** on the four PRs,
one of them relayed; the three adopted backlog captures the merges closed and
the one the milestone filed; and #241's GitHub edit history, which is empty.
The house form is [M11's document](11-housekeeping.md), with
[M10's](10-protocol-bookkeeping-from-the-field.md) behind it; the record
standard is the one the reviewer set on their PRs,
[#204](https://github.com/radiusred/gh-codecrew/pull/204),
[#206](https://github.com/radiusred/gh-codecrew/pull/206),
[#232](https://github.com/radiusred/gh-codecrew/pull/232) and
[#240](https://github.com/radiusred/gh-codecrew/pull/240) — verdict words
verbatim, no ordinal claim without every prior record linked, no judgement the
record does not carry, and a claim checked against the trail as it stands when
the record is written, Corrections included.

The gather from `gh codecrew milestone close 12` is again not part of the raw
material, for the reason [M8's](8-a-product-home-page-for-codecrew-works.md),
[M9's](9-the-docs-at-codecrew-works.md),
[M10's](10-protocol-bookkeeping-from-the-field.md) and
[M11's](11-housekeeping.md) records all give: `--dry-run` stops at the gate
that counts tasks, naming the task that writes this file. The installed
`v1.2.0` binary prints five gates:

```
gate milestone open: ok
gate tasks closed: refused[OPEN_TASKS]: tasks not closed: radiusred/gh-codecrew#251 (ready)
gate requirements declared: not reached
gate QA verdicts: not reached
gate milestone document: not reached
dry run: nothing written — the live verb stops at the first refusal above
```

A binary built from `main` at
[da630c5](https://github.com/radiusred/gh-codecrew/commit/da630c5) — the
merge of this milestone's last delivery — prints six, with the gate M12-R3
added between the first and the second, `gate no gate raised: ok`, and stops at
the same refusal. The difference between the two outputs is this milestone.

Every record below was read from its issue, PR, tag, run or timeline directly.

**This PR adds the M12 ROADMAP row; it does not flip one** — the convention
M10-R1 introduced, which
[M10's record PR](https://github.com/radiusred/gh-codecrew/pull/232) and
[M11's](https://github.com/radiusred/gh-codecrew/pull/240) were written under.
Until this milestone the convention was on `main` and in no release, so the
installed binary that opened M10 and M11 still appended an `Open` row locally.
It opened M12 the same way: #241 was created at 23:41:58Z on 2026-09-04 under
the installed `v1.1.0`, seventeen minutes and fifty-five seconds before this
milestone's own first task tagged `v1.2.0`, and the operator
[discarded the row](https://github.com/radiusred/gh-codecrew/issues/241#issuecomment-5547713635)
twenty-six seconds after creating the issue, "as for M10 (#207) and M11
(#233)", adding: "This is the last milestone opened under that binary: M12-R1
ships the verb that no longer writes the row." Under the binary M12-R1 shipped,
`milestone new --dry-run` prints no row — the operator's
[verification](https://github.com/radiusred/gh-codecrew/issues/242#issuecomment-5547856508)
and the
[qa comment](https://github.com/radiusred/gh-codecrew/issues/241#issuecomment-5549049198)
both say so.

## Goal and outcome

#241's goal is in two halves, and the order between them is the milestone's
shape. First, ship what M10 and M11 built: "The installed CLI is still v1.1.0,
so every run in the field gets the old behaviour: milestone new writes a
ROADMAP row nobody carries, milestone numbers can collide, status hides
milestone-issue gates and reports contract drift on three roles, task finish
dies on the raw GraphQL error on private repos. v1.2.0 is the boundary that
makes the docs true of the installed binary." Then, "behind it", three
bookkeeping gaps the field surfaced while M10 and M11 ran: `task new` refusing
`NOT_FOUND` seconds after `milestone new`
([#234](https://github.com/radiusred/gh-codecrew/issues/234)), `milestone
close` ignoring a gate raised on the milestone issue
([#219](https://github.com/radiusred/gh-codecrew/issues/219)), and `milestone
evidence` treating every URL in a record as a citation
([#222](https://github.com/radiusred/gh-codecrew/issues/222)).

All three gaps were backlog captures already written down when the milestone
opened, each filed from a field run an earlier record describes: #219 from the
ask-the-human point in #210's Plan during M10
([M10's record](10-protocol-bookkeeping-from-the-field.md#status-and-checkpoint-learn-what-a-milestone-issue-is)
quotes the operator's answer, "yes in principle, not in this task"); #222 from
`gh codecrew milestone evidence 4` on the `radiusred/ops` hub, "seen with
codecrew v1.1.0 on 2026-09-04"; and #234 from the three refused `task new`
calls while M11 was being opened, which
[M11's record](11-housekeeping.md#protocol-discipline-observations) names among
its observations. **All three captures were adopted and all three were closed
by the merges**; one new capture was filed and remains open.

**The release went first, and the fixes ride the next one.** The three
adoption comments, posted by the operator within three seconds of each other at
23:43:24Z–23:43:26Z, each end the same way: "The fix lands after v1.2.0
(M12-R1) is tagged, so it rides the next release"
([#234](https://github.com/radiusred/gh-codecrew/issues/234#issuecomment-5547719853),
[#219](https://github.com/radiusred/gh-codecrew/issues/219#issuecomment-5547719974),
[#222](https://github.com/radiusred/gh-codecrew/issues/222#issuecomment-5547720086)).
That is what happened: PR #246 merged at 23:59:39Z, the tag was cut fourteen
seconds later, and the three fix PRs merged after it. The installed `v1.2.0`
therefore carries none of M12-R2, M12-R3 or M12-R4; `main` carries all three.

Four requirements were present when the milestone opened at 23:41:58Z on
2026-09-04. None was amended. The four tasks were created between 23:42:26Z
and 23:42:35Z, nine seconds apart across the four, and this document is the
fifth. The last delivery merged at 01:29:42Z on 2026-09-05, one hour
forty-seven minutes and forty-four seconds after the milestone opened.

- **[#242](https://github.com/radiusred/gh-codecrew/issues/242) / PR
  [#246](https://github.com/radiusred/gh-codecrew/pull/246)** — M12-R1, Part
  A. One commit, amended once between two review rounds, merged 23:59:39Z
  ([f64d4c9](https://github.com/radiusred/gh-codecrew/commit/f64d4c9)). One
  Decision. Then Part B, outside any PR: the annotated tag `v1.2.0` on that
  commit at 23:59:53Z, tagger `radiusred-cody[bot]`, message
  `v1.2.0 (protocol 1.0)`; the release run
  [33931480565](https://github.com/radiusred/gh-codecrew/actions/runs/33931480565),
  created 23:59:56Z and `completed` / `success` at 00:01:21Z; and the
  [release](https://github.com/radiusred/gh-codecrew/releases/tag/v1.2.0),
  "gh-codecrew 1.2.0", published 00:01:07Z with five assets — darwin-amd64,
  darwin-arm64, linux-amd64, linux-arm64, windows-amd64.exe. The implementer's
  [evidence comment](https://github.com/radiusred/gh-codecrew/issues/242#issuecomment-5547849331)
  at 00:01:43Z and the operator's
  [verification](https://github.com/radiusred/gh-codecrew/issues/242#issuecomment-5547856508)
  at 00:02:35Z close the requirement end to end.
- **[#243](https://github.com/radiusred/gh-codecrew/issues/243) / PR
  [#247](https://github.com/radiusred/gh-codecrew/pull/247)** — M12-R2. Two
  commits, approved first round, merged 00:10:41Z
  ([cbf32fe](https://github.com/radiusred/gh-codecrew/commit/cbf32fe)). One
  Decision. Adopts #234 and cites its Correction.
- **[#245](https://github.com/radiusred/gh-codecrew/issues/245) / PR
  [#248](https://github.com/radiusred/gh-codecrew/pull/248)** — M12-R4. Three
  commits, two review rounds, merged 00:24:21Z
  ([3828a00](https://github.com/radiusred/gh-codecrew/commit/3828a00)). One
  Decision. Adopts #222. The reviewer's live probe during round one produced
  the milestone's one filed capture,
  [#250](https://github.com/radiusred/gh-codecrew/issues/250).
- **[#244](https://github.com/radiusred/gh-codecrew/issues/244) / PR
  [#249](https://github.com/radiusred/gh-codecrew/pull/249)** — M12-R3. Five
  commits, four review rounds, merged 01:29:42Z
  ([da630c5](https://github.com/radiusred/gh-codecrew/commit/da630c5)). Two
  Decisions and the milestone's one Deviation. Adopts #219.

Delivered in full, in one QA round. M12-R1, M12-R2, M12-R3 and M12-R4 were all
verdicted `satisfied` in the
[qa comment](https://github.com/radiusred/gh-codecrew/issues/241#issuecomment-5549049198)
of 03:29:07Z, one hour fifty-nine minutes after the last merge. No verdict in
this milestone came back other than `satisfied`.

What the hub looks like afterwards. `CHANGELOG.md` has a `## [1.2.0] —
2026-09-05` block of sixteen `###` sections — M10's and M11's work and the M7
to M11 records — under a `## [Unreleased]` that already holds the three fix
entries. `docs/introduction.md` says `**Shipped:** v1.2.0` and catalogues
thirty-two refusal codes. `task new` reads the hub's newest issues and retries
before `NOT_FOUND`; `milestone close` has six gates, the second of them
`no gate raised`, refusing `MILESTONE_GATED`; `milestone evidence` reads a URL
inside code as content, refuses on a dead `github.com` citation and warns on
any other. `gh codecrew status` run against the installed `v1.2.0` prints
`contract drift` for `roles/coordinator.md`, `roles/doc-synthesizer.md` and
`roles/qa.md` — the three contracts the fixes edited after the tag.

## Decisions

Six Decision comments. One is the operator's, on the milestone issue; the other
five were written by the seat doing the work, `radiusred-cody[bot]`, which
started all four delivery tasks. They group by the question each task had to
settle.

### The ROADMAP row, discarded for what the operator calls the last time

The operator's one Decision, posted at 23:42:24Z before any task existed, is
the counterpart of
[M10's](10-protocol-bookkeeping-from-the-field.md#the-roadmap-row-belongs-to-the-doc-synthesizer-at-both-ends)
and
[M11's](11-housekeeping.md#the-roadmap-row-discarded-again): the `Open` row
the installed `v1.1.0` binary appended locally is discarded, and this
document's PR adds the row already Done. Unlike those two, the
[comment](https://github.com/radiusred/gh-codecrew/issues/241#issuecomment-5547713635)
carries no Trade-off or Rejected line of its own; it refers to the earlier
two by number — "discarded, as for M10 (#207) and M11 (#233)" — and adds the
one thing that is new: "This is the last milestone opened under that binary:
M12-R1 ships the verb that no longer writes the row."

### The version-history lines keep their version

Flipping the shipped-release claim meant finding every `1.1.0` outside
`CHANGELOG.md` history and `docs/milestones/`. The
[Decision](https://github.com/radiusred/gh-codecrew/issues/242#issuecomment-5547737114)
leaves the five that the grep found: `docs/extensions.md`'s dated section
headings and `docs/platform-interop.md`'s "From v1.1.0 …" sentences. "They
record *when* a behaviour arrived, not what is installed; flipping them to
v1.2.0 would make them false." **Trade-off:** "a reader grepping for the
shipped version still hits five history lines; the one claim that describes
the installed release (`docs/introduction.md` "What exists") is the one that
moves." **Rejected:** rewording the history lines to avoid the version number,
because "the version is the point of those sentences (which release a
platform needs for the behaviour)".

The PR body records two more things that deliberately did not change, without
a Decision comment: the contract drift the installed `v1.1.0` reported needed
no edit, because `status` and `roles diff` compare the hub's `roles/*.md`
against the copies embedded in the running binary and the drift "clears when
the v1.2.0 binary is installed"; and `version` is stamped from the tag by
ldflags, so nothing in Go was bumped. The task's Plan named `a7a459d`, the
v1.1.0 flip from [#186](https://github.com/radiusred/gh-codecrew/issues/186),
as the exact form to follow, and the PR says it did.

### `task new` tries two routes, in order, three times

#234 proposed either a retry with back-off or a fallback search by title
prefix. The
[Decision](https://github.com/radiusred/gh-codecrew/issues/243#issuecomment-5547873151)
takes both, in a fixed order per read: the open-milestone listing first; on a
miss, the hub's newest issues regardless of label (`RecentIssues`, "in the
tracker since #209"), where an `M<n>:` title "counts only when the issue
itself carries `cc:milestone` and is open — confirmed through the existing
`Task(ref)` fetch, since that listing is unfiltered and shows closed issues";
on a miss there too, a two-second wait and both reads again, "three reads in
all before `refused[NOT_FOUND]`". The sleeper is a package variable
(`sleep = time.Sleep`, "the `ghVersion` pattern") so tests record the waits
rather than take them, and a milestone found by either fallback is printed as
a note "so the run's output shows the listing lagged".

**Trade-off:** "the worst case for a milestone that truly does not exist is
~4s slower (two waits) and up to six listing reads; the common case — the
number is there — is unchanged at one read." The fallback runs before the
first wait, "so in the #234 shape the verb most likely succeeds with no wait
at all (the unfiltered listing caught the just-created issue in #195's
analysis)". **Rejected:** a new tracker method to fetch one issue's labels and
state — "`Task(ref)` already returns both, and the seam stays as it is";
retrying alone without the fallback — "it would have made the #234 shell
script wait out the window every time instead of finding the issue at once";
and `milestone new` printing a `--milestone-ref` for `task new` to take, which
"#234 itself sets that aside as a larger change".

The second listing is the one
[M10's #209](10-protocol-bookkeeping-from-the-field.md#milestone-new-repairs-a-collision-instead-of-refusing-on-it)
added as "a second source for the floor and nothing more", because whether it
indexes a just-created issue faster "is not something the test suite can
measure". #243 relies on it for the same eventual-consistency window from the
other end, and records the same uncertainty as a likelihood rather than a
guarantee.

### The close gate is second, and reads labels the way `checkpoint` does

#244's
[first Decision](https://github.com/radiusred/gh-codecrew/issues/244#issuecomment-5547877612)
places the new gate: `milestone close`'s second, "after "milestone open" and
before "tasks closed", named "no gate raised" as `task finish`'s is", reading
the milestone issue's labels through `Tracker.IssueLabels` — "the REST issues
endpoint, added on PR #218" — because the `Milestone` value `OpenMilestones`
returns carries no labels. **Trade-off:** "one more REST read per close (a hub
closes a milestone rarely), and a dry run on a gated milestone shows every
bookkeeping gate as not reached — the operator learns about open tasks or
missing verdicts only once the question is answered." **Rejected:** placing it
last, after the document gate — "a requirement-level question may change what
"done" means, so surfacing it only after the tasks, verdicts and document have
all been counted is the wrong order"; and reading the labels through
`Tracker.Task` — "it would work for an issue, but `IssueLabels` is the route
the codebase already uses to read a gate's label off a ref, and one route is
easier to keep honest."

`IssueLabels` has a history in the record. #210's Decision in M10 rejected
exactly that method for `status` — "a second issue query for a field `Task`
already returns" — and
[M10's one seat Deviation](10-protocol-bookkeeping-from-the-field.md#deviations)
then added it for `checkpoint`, after the reviewer showed `Task` answers
`NOT_FOUND` for a pull request. #244 takes it up for the third verb and, this
time, rejects `Task` on the ground that the REST route now exists. `status`
still reads through `Task`.

PR #249's note for the reviewer records a count that does not match: #219's
capture said the close had "four gates", while `closeGates` had five with
`milestone open` counted, "so the docs here name the new gate by position
rather than by count". The four review rounds below turn on that list.

### `checkpoint` says the close refuses, because "nothing blocks" stopped being true

#244's
[second Decision](https://github.com/radiusred/gh-codecrew/issues/244#issuecomment-5547877773)
changes output M12-R3 does not name. `checkpoint`'s milestone-issue wording —
the `**Gate raised:**` comment and the stdout receipt — was set by PR #218 in
M10 to say `status` lists the gate, because nothing mechanical then blocked on
the label; the code comment on `raiseGate` said the same. After this PR both
would be false "on the record the moment it is written", so the comment and
receipt now say `status` lists the gate *and* `milestone close` refuses while
the label is present. **Trade-off:** "a change to `checkpoint`'s output that
M12-R3 does not name, with the wording test updated alongside." **Rejected:**
leaving the wording as PR #218 set it — "the comment is what the person
resolving the gate reads".

### A citation is a link outside code; a dead one refuses or warns by host

#222 offered three shapes: skip URLs inside code; treat only GitHub links as
citations; or report non-GitHub failures as warnings. #245's
[Decision](https://github.com/radiusred/gh-codecrew/issues/245#issuecomment-5547876907)
is "#222's first shape combined with its third". A citation is "a URL in prose
or in a Markdown link outside code"; a URL inside an inline code span ("any
backtick-run length, CommonMark matching") or a fenced block ("``` or ~~~,
closed by a fence of the same character at least as long") is content and is
not scanned. An unreachable citation then splits by host: `github.com` refuses
`EVIDENCE_UNREACHABLE` as before; any other host "is reported on a `warning:`
line and the verb passes, with the summary line counting both".

**Trade-off:** ""GitHub" is the literal `github.com` host, so
`docs.github.com`, `raw.githubusercontent.com` and the hub's own site
(codecrew.works) are external — a dead link to one of them warns rather than
refuses." The Decision calls that "the conservative reading of #222", whose
examples are "issues, PRs, comments, commits, runs", and notes that "a warning
is still printed, so nothing dies silently"; the refusal's numerator counts
`github.com` failures only while its denominator is every citation, "so `1 of
4` reads as before". **Rejected:** the second shape alone — "it would drop
external links from the report entirely, and a dead link to a design doc or a
release page is worth a line"; and a `note:` prefix to match the CLI's other
advisory lines — "this line is a finding for QA to weigh, not a bookkeeping
aside, and the brief names it a warning."

This document is written under the rule that Decision produced. The
doc-synthesizer contract now says a URL that is not evidence goes inside a
code span, "and nobody edits a comment to hide it from the scanner" — the
remedy #222 records the coordinator having applied on `radiusred/www#64`. The
two NXDOMAIN-by-design hosts that capture quotes, `hooks.radiusred.uk` and
`zoo.radiusred.uk`, appear here in code for that reason.

## Deviations

One Deviation comment, on #244, recorded by `radiusred-cody[bot]` at
01:29:35Z — twenty-three seconds after the approval it refers to and seven
seconds before the merge.

**[The reviewer's round-one verdict on PR #249 was relayed after a machine crash](https://github.com/radiusred/gh-codecrew/issues/244#issuecomment-5548420342).**
"The reviewer's round-one verdict on PR #249 (at 372ac24) was written by the
reviewer run but the machine crashed before the review was submitted; the
coordination layer posted that verdict verbatim as `radiusred-checky[bot]`
with a relay note, and rounds two to four were fresh reviewer runs. The
approval that gates this merge (review 5119172717 at 081de3f) is the
reviewer's own." The relayed
[review](https://github.com/radiusred/gh-codecrew/pull/249#pullrequestreview-5119146783)
carries the note at its foot, set off by a rule: "_Relayed by the coordination
layer: this verdict was written by the reviewer run at 372ac24 but the machine
crashed before the review was submitted; the text above is that run's
verdict.md verbatim. Round two will be a fresh reviewer run._" The verdict
above the note is a complete review — the finding, the checks run, a live
dry-run result, and the reviewer's identity verification against
`gh api /apps/radiusred-checky`.

The Deviation records the fact the non-doer gate depends on, and the trail
bears it out: the four review submissions on #249 are all by
`radiusred-checky[bot]`, and the merge was gated by the fourth, submitted at
01:29:12Z by a run the Deviation calls fresh. What the record does not show of
the crash is in the last section.

PRs #246, #247 and #248 make no Deviation of their own. #247's and #248's
descriptions say "no deviations from the plan"; #246's makes no such
statement, and what its Plan's step 7 said would follow the merge — the tag,
the run, the evidence comment — is what the trail shows. None is recorded on
their task issues.

## The gates

**#241's Gates section was left as the scaffold's placeholder** — the template
line "_What "done" means beyond CI: e2e suites, manual UAT, sign-offs._",
unedited. [M11's record](11-housekeeping.md#the-gates) says the same of #233,
[M10's](10-protocol-bookkeeping-from-the-field.md#the-gates) of #207,
[M9's](9-the-docs-at-codecrew-works.md#the-gates) of #201 and
[M8's](8-a-product-home-page-for-codecrew-works.md#the-gates) of #196. No
`**Gate raised:**` or `**Gate resolved:**` comment exists anywhere in this
milestone, no `cc:needs-decision` label was applied, and `checkpoint` was not
used — on a task or on the milestone issue.

M12-R3's subject is what `milestone close` does when such a label *is* on the
milestone issue, so — as M10-R3 before it, whose subject was what `status`
does in the same case — the requirement was verdicted against the case the
milestone did not produce. The qa seat ran the built verb's dry run live
against #241 and read the negative: `gate milestone open: ok`,
`gate no gate raised: ok`, then `OPEN_TASKS` on this task, "this proves #241
has no milestone-issue gate and the live close was not run". The positive path
was verified from `TestPlanCloseRefusesWhileMilestoneIssueIsGated`, which the
verdict describes as exercising "the hairbrush ordering: with `cc:needs-decision`
present and an open task also present, dry-run reports
`gate no gate raised: refused[MILESTONE_GATED]` and leaves `tasks closed` not
reached; removing the label reaches `OPEN_TASKS`."

Three of the four task Plans wrote `None` into their Ask-the-human section and
gave reasons: #242's, "the release content is what M10 and M11's merged tasks
already decided; the version and date follow from them"; #243's, "the fix is
the shape #234 and the milestone requirement describe"; #244's, "the shape
(code, gate family, dry-run line) is fixed by M12-R3 and #219; the ordering and
the `checkpoint` wording are engineering calls recorded as Decisions". #245's
Plan left the scaffold's own line in place — "_None identified yet — the
implementer adds any before starting._" — and added none.

What gated the work: CI on every PR — `Lint commit messages`, which was failing
on PR #246's first head and is where the milestone's first change request came
from, and `Go build and test` — an independent approval on each of the four
PRs, and the QA verdicts at the close. For M12-R1 the requirement itself named
a further gate outside any verb: the operator's own machine, where
`gh extension upgrade codecrew` had to print `v1.2.0`, `status` had to stop
reporting drift, and `milestone new --dry-run` had to print no row. The
[verification comment](https://github.com/radiusred/gh-codecrew/issues/242#issuecomment-5547856508)
records all four, with the upgrade line quoted: `upgraded from v1.1.0 to
v1.2.0`.

## The review rounds

Four PRs, **nine review submissions**, all by the reviewer role holder
(`radiusred-checky[bot]`) — one of them posted by the coordination layer under
that identity, with a relay note — and **five change requests**, every one
resolved in the next round.

- **PR [#246](https://github.com/radiusred/gh-codecrew/pull/246)** — two
  rounds. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/246#pullrequestreview-5118762931)
  at 23:51:09Z was not about the diff: "The PR is mechanically correct, but
  the required commitlint check is failing." Commit `aa8a042`'s body was one
  line of 293 characters against the 100-character limit. The seat rewrapped
  the body and force-pushed `aca564c` at 23:51:45Z — "subject unchanged, diff
  unchanged — the tree is identical to aa8a042" — and the
  [approval](https://github.com/radiusred/gh-codecrew/pull/246#pullrequestreview-5118805030)
  at 23:58:37Z confirmed both: "the tree is identical to aa8a042, the amended
  commit body has no line over 100 characters, and both required checks are
  green".
- **PR [#247](https://github.com/radiusred/gh-codecrew/pull/247)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/247#pullrequestreview-5118872479)
  at 00:09:51Z, four and a half minutes after opening. The review states its
  reading order — "I reviewed the diff before the PR body" — and its sources,
  including "capture #234 including its Correction comment", and records that
  "task plan and Decision precede the commits".
- **PR [#248](https://github.com/radiusred/gh-codecrew/pull/248)** — two
  rounds. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/248#pullrequestreview-5118908845)
  at 00:17:16Z had three findings: `roles/qa.md` "still says "Every link the
  record cites must resolve before you test against it"" though the PR lets
  external citations pass as warnings; `internal/cli/cli.go`'s help "still says "verify every cited link
  in the milestone's record resolves""; and the tests did not cover "the
  explicit "fence inside a list item" shape from the review brief". The
  [fix commit](https://github.com/radiusred/gh-codecrew/commit/c5b702e)
  answered all three, with the seat
  [recording](https://github.com/radiusred/gh-codecrew/pull/248#issuecomment-5547991122)
  that "no test pinned the old line" of help, and that the fence opener now
  skips a leading list marker so "a fence opening on the bullet's own line is
  a fence". The
  [approval](https://github.com/radiusred/gh-codecrew/pull/248#pullrequestreview-5118949993)
  at 00:23:26Z verified "marker-line fences that fail without stripping the
  list marker". The same review's verification list records the probe that
  became #250: `./gh-codecrew milestone evidence 11` "refused `NOT_FOUND`
  because M11 issue #233 is closed and the verb only searches open
  milestones". The operator filed the capture thirty-four seconds after the
  review was submitted.
- **PR [#249](https://github.com/radiusred/gh-codecrew/pull/249)** — the
  milestone's only multi-round review beyond two, described below.

**PR #249's chronology, from its timeline.** Opened 00:06:33Z with two commits.
Force-pushed at 00:25:40Z for the rebase onto #248's merge — the seat's
[comment](https://github.com/radiusred/gh-codecrew/pull/249#issuecomment-5548046419)
records a `CHANGELOG.md` conflict resolved by keeping this entry at the top
above #245's and #243's, and a `SPEC.md` conflict resolved by keeping #245's
new `milestone evidence` row and reapplying the `milestone close` edit. Then
nothing on the PR for fifty-two minutes. The
[round-one review](https://github.com/radiusred/gh-codecrew/pull/249#pullrequestreview-5119146783)
at 01:18:48Z is the relayed one. Its finding: `internal/cli/cli.go:19` "still
says `milestone close` is `close a milestone (gates: tasks closed, doc
merged)`", which "is now stale" and is "an operator-facing documented surface".
The fix
[9686ebf](https://github.com/radiusred/gh-codecrew/commit/9686ebf) landed
twenty-seven seconds later, listing "no gate raised, tasks closed,
requirements, QA verdicts, doc merged" — five gates.
[Round two](https://github.com/radiusred/gh-codecrew/pull/249#pullrequestreview-5119156835)
at 01:22:53Z: `closeGates` "starts with `milestone open`, and the live dry-run
prints `gate milestone open` first"; the help should too. The fix
[f4662de](https://github.com/radiusred/gh-codecrew/commit/f4662de) added it —
six gates.
[Round three](https://github.com/radiusred/gh-codecrew/pull/249#pullrequestreview-5119164892)
at 01:26:16Z: the help "still does not name all six `closeGates` entries
exactly" — "requirements" and "doc merged" are not the strings `requirements
declared` and `milestone document`. The fix
[081de3f](https://github.com/radiusred/gh-codecrew/commit/081de3f) used the
exact names, and
[round four](https://github.com/radiusred/gh-codecrew/pull/249#pullrequestreview-5119172717)
at 01:29:12Z approved: "it now lists the six closeGates strings verbatim in
order". The seat's three answering comments each say "no test pinned the old
line" or that gofmt, vet and tests are green; the reviewer's each say the
follow-up diff touches only `internal/cli/cli.go`. From the relayed review to
the merge: ten minutes fifty-four seconds; from the first fix commit to the
approval, three commits and three reviews inside ten minutes.

**What the five change requests were about.** One was a required check the
diff itself did not touch — a commit body line over the limit, on the one PR
in the milestone that changed no Go. Four were documented surfaces that no
longer matched the behaviour the PR shipped: a contract (`roles/qa.md`), two
help lines in `internal/cli/cli.go` (the `milestone evidence` entry once, the
`milestone close` entry three times), and a test that did not lock a shape the
review brief asked for. On both PRs the seat recorded that no test pinned the
help line. [M10's record](10-protocol-bookkeeping-from-the-field.md#the-review-rounds)
describes the change request on PR #216 as the same surface for the same
reason — the top-level help advertised a removed behaviour, and "the new
`runMilestoneNew` test bypasses `cli.go` and therefore cannot catch this
regression". Every approval was reached with the reviewer running the code:
repo-local builds, live dry-runs of `milestone close 12` and probes of
`milestone evidence` against the hub, a range check on each follow-up diff.

The reviewer recorded its identity check on two of the four PRs — in the round-one
review of #248 and inside the relayed verdict on #249: the config App ID
`4719924` matched `gh api /apps/radiusred-checky --jq .id` and differed from
the PR author `radiusred-cody[bot]`. That the check appears in the relayed text
means it was run by the reviewer before the crash, not by the relay.

## QA: one round, four verdicts

The qa role holder (`radiusred-testy[bot]`) verdicted all four requirements in
a
[single comment](https://github.com/radiusred/gh-codecrew/issues/241#issuecomment-5549049198)
at 03:29:07Z, against merged `main` and the installed `v1.2.0` both, having
run `go test ./...` and `go vet ./...` with repo-local caches. It records one
wrinkle of that setup: a literal `gofmt -l .` "lists only repo-local
`gomodcache/gopkg.in/yaml.v3@v3.0.1` files created by the required cache
setting, while `gofmt -l $(git ls-files '*.go')` is clean".

**M12-R1** was verdicted from the artefacts: `milestone evidence 12` passing
"installed-first" and built alike (`all 2 cited links resolve across 6
issues`); `gh release view v1.2.0` reporting the tag and the five assets by
name; `git cat-file -p v1.2.0` showing "an annotated tag by
`radiusred-cody[bot]` pointing at PR #246's merge commit"; the tag's
`CHANGELOG.md` with "a fresh empty Unreleased and a 16-section `1.2.0` block";
the introduction's two `v1.2.0` claims; the installed extension printing
`v1.2.0 (protocol 1.0)`; and `milestone new --dry-run` printing "no ROADMAP
row". On drift the verdict is explicit about which binary and which moment:
"The operator's #242 release verification recorded no drift at the v1.2.0
boundary; my current-main `status` run prints only the expected post-tag drift
lines for `roles/coordinator.md`, `roles/doc-synthesizer.md`, and
`roles/qa.md`, because #243/#245/#244 changed those contracts after the
release tag."

**M12-R2** was verdicted without a live probe, and says so first: "I did not
probe `task new` live because that would create an issue." The verdict reads
`resolveMilestone`'s order — `OpenMilestones`, then `RecentIssues` confirmed
through `Task(ref)` "as open and `cc:milestone`", then "a bounded 2s wait for
three reads before `refused[NOT_FOUND]`" — and the suite's cases, including
"rejection of a closed prefixed issue and a non-milestone prefixed issue, and
bounded NOT_FOUND with no create/link side effects". It reads the coordinator
contract's new line back — "operators should not script their own pause or
retry loop" — and confirms #234 closed by PR #247 "with the adoption
back-reference".

**M12-R3** was verdicted on the live dry run against #241 described under
[The gates](#the-gates), the test's hairbrush ordering, and the surfaces the
four review rounds had fought over: "Top-level help now lists all six
`closeGates` strings exactly, `docs/introduction.md` and README say thirty-two
refusal codes, and non-test source has 32 distinct `refuse("CODE"` sites." It
names the Deviation as part of what closed #219: "including the recorded
crash-relay deviation for review round one and the final approval on the
corrected help text".

**M12-R4** was verdicted on a case the seat built, because the milestone's
own record gave the fix nothing to do: "there was no behavioral difference on this record
because it contains no quoted probe URLs." The seat built the case instead —
"an uncommitted scratch record" at a path under its dispatch directory, with
"an inline-code URL, a fenced-block URL inside a list item, a dead external
citation, and a dead github.com citation", and "a temporary Go test in the
clone" that "proved only the external and github.com citations were checked,
with the external classified as a warning and the github.com link as
`EVIDENCE_UNREACHABLE`". The committed suite "also covers list-item fences and
warning/refusal classification"; both contracts "state the rule"; #222 is
closed by PR #248, and "#250 remains open as the unrelated backlog capture for
running evidence on closed milestones".

## Requirement outcomes

The status column is the verdict word as the qa role holder wrote it, verbatim
and unqualified. Everything else the milestone knows about a requirement is in
the columns around it.

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M12-R1 — v1.2.0 ships: the `[Unreleased]` changelog becomes the 1.2.0 section with a fresh empty Unreleased above it and the compare links updated; the introduction's shipped-release claim flips; the tag is cut by the implementer identity after the flip merges; `release.yml` builds it; and the release is verified by `gh extension upgrade codecrew` printing v1.2.0 and `status` no longer reporting contract drift | [#242](https://github.com/radiusred/gh-codecrew/issues/242) / PR [#246](https://github.com/radiusred/gh-codecrew/pull/246); tag `v1.2.0` on [f64d4c9](https://github.com/radiusred/gh-codecrew/commit/f64d4c9); run [33931480565](https://github.com/radiusred/gh-codecrew/actions/runs/33931480565); [release](https://github.com/radiusred/gh-codecrew/releases/tag/v1.2.0) | `satisfied` | Task closed; the one QA round, against the tag object, the release assets, the installed extension and the built tree. The operator's own-machine checks are recorded in a [comment](https://github.com/radiusred/gh-codecrew/issues/242#issuecomment-5547856508); the no-drift state it records lasted until PR #247 merged eight minutes later |
| M12-R2 — `task new` resolves a milestone created seconds earlier: when the label-filtered listing has no open `M<n>` it retries with a short back-off, or reads the milestone by another route, before refusing `NOT_FOUND`; the coordinator contract says what to expect right after `milestone new` (adopts [#234](https://github.com/radiusred/gh-codecrew/issues/234)) | [#243](https://github.com/radiusred/gh-codecrew/issues/243) / PR [#247](https://github.com/radiusred/gh-codecrew/pull/247) | `satisfied` | Task closed; the one QA round, from the built source and its suite, with no live probe — "that would create an issue". Not in the installed `v1.2.0` |
| M12-R3 — `milestone close` refuses with a code while the milestone issue itself carries `cc:needs-decision`, in the same family as `task finish`'s no-gate-raised gate, and `--dry-run` prints that gate line (adopts [#219](https://github.com/radiusred/gh-codecrew/issues/219)) | [#244](https://github.com/radiusred/gh-codecrew/issues/244) / PR [#249](https://github.com/radiusred/gh-codecrew/pull/249) | `satisfied` | Task closed; the one QA round, on a live dry run of the built verb against #241 in the negative and the test in the positive, plus the help line, both docs' count and a count of `refuse("` sites. Four review rounds, the first relayed. Not in the installed `v1.2.0` |
| M12-R4 — `milestone evidence` checks only the links a record cites as evidence: a URL quoted as a probe target meant to be unreachable, or as a verbatim command or error string, is not a citation and does not fail the gate; the doc-synthesizer contract states the rule (adopts [#222](https://github.com/radiusred/gh-codecrew/issues/222)) | [#245](https://github.com/radiusred/gh-codecrew/issues/245) / PR [#248](https://github.com/radiusred/gh-codecrew/pull/248) | `satisfied` | Task closed; the one QA round, on a scratch record and a temporary test the seat built because M12's own record "contains no quoted probe URLs". Not in the installed `v1.2.0` |

No requirement was dropped, added or amended.

**Captures adopted and closed by the merges — three:**
[#234](https://github.com/radiusred/gh-codecrew/issues/234) (`task new`
refuses `NOT_FOUND` for a milestone created seconds earlier; its
[Correction](https://github.com/radiusred/gh-codecrew/issues/234#issuecomment-5547397489)
was read with it, and the record uses the issue timestamps it names),
[#219](https://github.com/radiusred/gh-codecrew/issues/219) (`milestone close`
ignores a gate raised on the milestone issue) and
[#222](https://github.com/radiusred/gh-codecrew/issues/222) (`milestone
evidence` treats every URL in the record as a citation). Each was closed by
its PR's merge, at 00:10:42Z, 01:29:44Z and 00:24:23Z. **Capture filed during
the milestone — one, open:**
[#250](https://github.com/radiusred/gh-codecrew/issues/250) (`milestone
evidence` refuses `NOT_FOUND` on a closed milestone), opened by the operator at
00:17:50Z from the reviewer's live probe on PR #248, with its source named:
"checky's review of PR #248 (task #245), 2026-09-05".

## Protocol-discipline observations

Seven things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **The release went first by design, and the trail shows what that costs
  within the hour.** The operator verified no contract drift on the installed
  `v1.2.0` at 00:02:35Z; PR #247 merged at 00:10:41Z and changed
  `roles/coordinator.md`; by 01:29:42Z three contracts differed from the copies
  embedded in the binary, and the qa comment's `status` run — like the one made
  for this document — prints three `contract drift` lines again. The line is the
  same one `v1.1.0` printed when the binary was behind the hub; now it prints
  because the hub is ahead of the binary. `status` does not distinguish the two,
  and nothing in the record asks it to.
- **All three adopted captures were named by an earlier record before this
  milestone adopted them.** M10's record quoted the operator's "yes in
  principle, not in this task" that produced #219; M11's record described the
  three refused `task new` calls behind #234 and said M10 "left the other end
  open; it was hit while opening the very next milestone"; #222 was filed on
  the day M10 ran, from another hub on the released binary. M10's record made
  the same observation of #197: "a milestone document named a protocol gap, and
  the next milestone adopted the capture it pointed at." That is now the shape
  of three consecutive milestones' backlogs.
- **A review verdict reached the record by relay, and the record says so in
  three places.** The relayed review carries its own note; the Deviation on
  #244 names the crash, the relay and which approval gated the merge; and the
  qa verdict on M12-R3 names "the recorded crash-relay deviation". What the
  protocol did not produce is a timestamp for the event itself, and the
  Deviation was written after round four's approval rather than when the relay
  happened — SPEC §4's "at the moment they occur" is fifty-two minutes of silence
  on the PR followed by a review, and the comment that explains the review came
  ten minutes forty-seven seconds after it.
- **Three of the five change requests were about one help line, and no test
  pins it.** `internal/cli/cli.go`'s top-level help was the surface on PR #216
  in M10, on the `milestone evidence` entry in PR #248, and three rounds
  running on the `milestone close` entry in PR #249, where the finding narrowed
  from "stale" to "five of six" to "not the exact names". The seat wrote "no
  test pinned the old line" on both PRs. The reviewer read the help against
  `closeGates` and the dry run's output each time; nothing in the tree does.
- **The hub lane serialised on `CHANGELOG.md` again.** The three fix PRs
  opened within seventy-four seconds — 00:05:19Z, 00:05:51Z, 00:06:33Z — and
  each of the later two force-pushed a rebase onto the previous merge
  (00:11:45Z, 00:25:40Z), recording a `CHANGELOG.md` conflict resolved by
  keeping both entries, and on #249 a `SPEC.md` conflict as well.
  [M10's record](10-protocol-bookkeeping-from-the-field.md#protocol-discipline-observations)
  observed the same ordering under the same cause with four PRs.
- **The doc-synthesizer seat held no delivery task, and nothing recorded that
  either.** M10's and M11's records each observed `radiusred-wordy[bot]`
  starting delivery tasks under the implementer contract and noted that nothing
  in the contracts says whether a role's holder may be routed to another role's
  task. Here all four were started by `radiusred-cody[bot]`. The routing is as
  unrecorded when it follows the table as when it did not.
- **The refresh obligation found nothing stale, and the reason is a number
  that runs ahead of the binary.** `docs/introduction.md` says `**Shipped:**
  v1.2.0` and, with `README.md`, "thirty-two" refusal codes; the shipped
  `v1.2.0` binary has thirty-one, because `MILESTONE_GATED` merged after the
  tag. The introduction names its source as "the catalogue of record —
  `refuse("CODE"` in `internal/cli/`", which is `main`, and the same
  relationship held under `v1.1.0` after
  [M10's #209 and #211](10-protocol-bookkeeping-from-the-field.md#protocol-discipline-observations)
  moved the count to thirty-one on `main` with no release behind it. The
  sweep for this document changed neither file: the claims are true of the
  source they name, and the record states the gap rather than "fixing" a count
  the next release will make true.

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **The crash has no artefact of its own.** The relay note and the Deviation
  record that a reviewer run at `372ac24` wrote a verdict, that "the machine
  crashed before the review was submitted", and that rounds two to four were
  fresh runs. Neither gives a time, a machine, or the run; the only trace on
  the trail is the fifty-two-minute interval between the seat's rebase comment
  at 00:26:34Z and the relayed review at 01:18:48Z. The implementer side shows
  a split too: the two commits PR #249 opened with carry one `Claude-Session`
  trailer and the three follow-up commits carry another. *Inferred, not
  recorded:* that the same crash ended the implementer's first session; nothing
  on the trail says so.
- **The operator's M12-R1 verification is a comment, not a transcript.** The
  four installed-binary checks the requirement names were run on the operator's
  machine and their outcomes reported; no shell output is attached. The qa seat
  later reproduced two of them on the installed extension — `version` and
  `milestone new --dry-run` — but the third, `status` printing no drift, was
  true for eight minutes and cannot be reproduced from the trail since.
  *Recorded as reported, corroborated in part.*
- **`task new`'s retry has never run against the lag it was written for.** The
  qa seat did not probe it live, by its own account, and no run in the record
  shows the note a fallback prints. The #234 shape — three refusals and a fourth
  success — has not recurred on the trail since the capture, and cannot be met
  by the fix from an installed binary until the next release.
- **`MILESTONE_GATED` has never been raised live.** As with M10-R3, the
  milestone's own issue carried no gate, so the verdict rests on the dry run's
  negative and the test's positive. No hub in the record has yet had
  `milestone close` refuse on a milestone-issue gate.
- **The qa seat's M12-R4 proof is not in the tree, and the original failure
  was not re-run.** The scratch record at
  `/home/darren/.cache/codecrew-dispatch/testy-m12/evidence-scratch.md` and the
  temporary test `TestScratchEvidenceFixtureClassifiesCitationsOnly` are
  described in the verdict and exist nowhere else. `milestone evidence 4` on
  `radiusred/ops`, the run #222 was filed from, was not repeated with the fixed
  verb, and nothing records whether the comment on `radiusred/www#64` that the
  coordinator edited "to hide hostnames from the link scanner" was put back.
- **The review brief for PR #248 is not on the trail.** The reviewer's third
  finding cites "the explicit "fence inside a list item" shape from the review
  brief", a shape neither #245's Plan nor M12-R4 names; it reached the code
  through that finding and commit `c5b702e`. *Recorded as cited, not seen:* the
  brief lives in the coordination layer's dispatch directory, as
  [M10's record](10-protocol-bookkeeping-from-the-field.md#what-the-record-does-not-contain)
  says of two events that reached it the same way.
- **#250 is filed, open and unadopted — and it applies to this record.**
  `milestone evidence` resolves a milestone through the open-milestone listing
  only, so once #241 closes the verb cannot check this document's citations.
  The four it counts today — two on #242's evidence comment, two in this
  task's Plan — are checked only while the milestone is open.
- **No gate was declared or raised.** #241's Gates section is the scaffold's
  placeholder text, unedited; `checkpoint` was not used on a task or on the
  milestone issue, including for the verb M12-R3 taught to refuse on its
  label; and #245's Plan carries the template's Ask-the-human line unedited.
