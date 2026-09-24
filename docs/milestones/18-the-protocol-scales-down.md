# M18: The protocol scales down

Tracking issue: [#353](https://github.com/radiusred/gh-codecrew/issues/353) ·
Synthesized 2026-09-24 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #353's **eight** requirements as opened, **none added,
none struck and none amended** — the issue's `userContentEdits` holds the
original and one revision, nine seconds later, which added the three lines
under Gates and nothing else; its **seven comments** — the operator's two
release legs, **two coordination-layer Decisions** and one correction
between them, the qa role holder's **one comment holding all eight
verdicts**, and this record's **one Deviation**; the **ten** task issues, all in this repository
([#354](https://github.com/radiusred/gh-codecrew/issues/354),
[#355](https://github.com/radiusred/gh-codecrew/issues/355),
[#358](https://github.com/radiusred/gh-codecrew/issues/358),
[#359](https://github.com/radiusred/gh-codecrew/issues/359),
[#362](https://github.com/radiusred/gh-codecrew/issues/362),
[#363](https://github.com/radiusred/gh-codecrew/issues/363),
[#364](https://github.com/radiusred/gh-codecrew/issues/364),
[#369](https://github.com/radiusred/gh-codecrew/issues/369),
[#371](https://github.com/radiusred/gh-codecrew/issues/371),
[#372](https://github.com/radiusred/gh-codecrew/issues/372)); their **ten
merged pull requests** and the **forty-five commits** on them; the one
housekeeping pull request the milestone's release plan called for,
[#375](https://github.com/radiusred/gh-codecrew/pull/375), and its **one**
commit; the **twenty-six Decision records and seven Deviation records**
across the milestone issue and its ten tasks, of which **two and one** are on
the milestone issue and **twenty-four and six** on the tasks, with **no**
`**Gate resolved:**` record anywhere; the **twenty review submissions** on
the eleven pull requests — **eleven approvals, every one of which stands**,
**eight change requests**, and **one** review posted under the operator's
login and dismissed; the **nine** adopted backlog captures the merges
closed; and the **one** capture filed in the milestone's window, which is
open. The house form is
[M17's document](17-the-frontier-builds-the-coordinator.md), with
[M16's](16-bedding-in-the-venue-seam-and-the-reference.md) behind it; the
record standard is the one the reviewer set on
[PR #322](https://github.com/radiusred/gh-codecrew/pull/322) and re-applied on
[PR #341](https://github.com/radiusred/gh-codecrew/pull/341) and
[PR #352](https://github.com/radiusred/gh-codecrew/pull/352), so every count
above and below is re-derived from the source — records by the rule
`tracker.ExtractRecords` applies (a paragraph-initial `**Decision:**`,
`**Deviation:**` or `**Gate resolved:**` label, bare or parenthetically
qualified), review states from each pull request's reviews API and timeline,
commits from each pull request's commits API, and the requirement history
from the milestone issue's `userContentEdits`.

**Counted at 2026-09-24T19:11:54Z, against a trail frozen from this
document's dispatch.** The last post on the milestone's issues before that
instant is this record's own
[Deviation](https://github.com/radiusred/gh-codecrew/issues/353#issuecomment-5820516906),
at 19:11:43Z; the one before it is the qa role holder's verdict comment at
19:05:04Z. The Deviation is inside every total above.

**This record has no task.** It is delivered as a housekeeping pull request
under M18-R7, which this milestone shipped: the doc-synthesizer contract's
"Deliver as a housekeeping PR", with the milestone issue as its charter and
its commits referencing #353.
[M17's record](17-the-frontier-builds-the-coordinator.md) and
[M16's](16-bedding-in-the-venue-seam-and-the-reference.md) each counted their
own task's Decisions; this one has no task and none to count, and its one
record of its own is the Deviation on the milestone issue, where SPEC §4 now
puts it.

**`main` moved by exactly forty-six commits**, in a straight line with no
merge commits, from M17's last record commit
([f6b769f](https://github.com/radiusred/gh-codecrew/commit/f6b769f),
committed 2026-09-23T11:41:22Z) to the housekeeping merge
([cab43cc](https://github.com/radiusred/gh-codecrew/commit/cab43cc)):
forty-five that reference one of the ten tasks and one, `chore: regenerate
the contract history for v2.1.0`, that references none. Every one of the
forty-six carries the trailer `Co-Authored-By: Claude Opus 5.5 (1M context)`
and the author `radiusred-cody[bot]`.

**The trail is checked by the verb.** The installed extension reports
`v2.1.0 (protocol 2.1)`, the release this milestone shipped. Run from this
repository at the counting instant:

```
$ gh codecrew milestone evidence 18
requirements counted: M18-R1, M18-R2, M18-R3, M18-R4, M18-R5, M18-R6, M18-R7, M18-R8 (8)
all 10 cited links resolve across 11 issues — evidence is reachable
```

Eleven issues: the milestone and its ten sub-issues. Ten links, all on
github.com: the two release runs and the two release pages, the
coordination-layer Decision that widened M18-R8 (cited from #353 and from
#371's plan, counted once), the operator's answers on #362, and four comments
on #371 — its three Decisions, cited from its own plan, and the release
evidence, cited from the operator's v2.1.0 leg on #353. QA's M18-R7 verdict
reports the same run: "live `milestone evidence 18` resolved all ten
citations across eleven issues". This record's Deviation cites no URL.

The gather from `gh codecrew milestone close 18 --dry-run` reads differently
from the one earlier records report. [M11's](11-housekeeping.md),
[M12's](12-v1-2-0-and-the-field-fixes-behind-it.md),
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md),
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md),
[M15's](15-v2-0-1-what-the-fleet-migration-taught.md),
[M16's](16-bedding-in-the-venue-seam-and-the-reference.md) and
[M17's](17-the-frontier-builds-the-coordinator.md) records each report the
dry run stopping at the gate that counts tasks, naming the task that wrote
the record. With no record task, it passes five gates and stops at the
document. As it read at 19:09:13Z, before this pull request existed:

```
gate milestone open: ok
gate no gate raised: ok
gate tasks closed: ok
gate requirements declared: ok
gate QA verdicts: ok
gate milestone document: refused[DOC_MISSING]: docs/milestones/18-*.md not on the default branch of radiusred/gh-codecrew — dispatch the doc-synthesizer: it delivers the record as a housekeeping PR with no task (SPEC §4), reviewed as the reviewer seat's routing requires (in pure solo, the operator confirms on the PR), and its author rebase-merges it; then rerun
dry run: nothing written — the live verb stops at the first refusal above
```

The detail on the last gate line is the one #369 wrote. The dry run prints no
raw material, so the Decisions and Deviations below were read from each
issue's comments directly.

Every record below was read from its issue, pull request, review, commit or
timeline directly. The prose is wrapped at this file's normal width; the
requirement-outcomes table's rows, the fenced output above and a handful of
single links are not, because a newline ends a Markdown table row and
breaking a link breaks it — the shape
[M17's record](17-the-frontier-builds-the-coordinator.md#requirement-outcomes)
carries, and the shape `ROADMAP.md` has had since M1.

**This PR adds the M18 ROADMAP row; it does not flip one.** `ROADMAP.md`
carried no M18 row before this PR. The PR also changes one word in
`docs/introduction.md`, for the reason in the section on the front door
below.

## Goal and outcome

#353 was opened at 11:42:19Z on 2026-09-23, nineteen seconds after
[`milestone close 17`](https://github.com/radiusred/gh-codecrew/issues/345#issuecomment-5794168852)
closed M17 with "all 4 tasks done, milestone document merged". Its goal opens
on the fleet: "Projects running 2.0 outside this hub — a solo blog, a browser
extension's first release, a project running stale contracts until a
milestone was free to take them — reported ceremony that was the same size
whatever the change." It lifts M16's rule for this milestone only: "The
operator lifts the bedding-in rule for this milestone (2026-09-23), and only
for the captures it adopts."

It names two kinds of work, "released separately". First, "corrections that
change no behaviour", shipped as v2.0.2 with M16's unreleased venue seam and
`CLI.md`, "so the fleet gets a clean 2.0 base before anything moves". Second,
"protocol 2.1, all additive: a requirement can be struck by a recorded
decision; housekeeping gets a reviewed light path and a `roles sync` verb;
the milestone record takes that light path instead of a task; the
doc-synthesizer contract stops assuming the hub is a product repository."
And it names what does not move: "The requirement/verdict pair, the non-doer
review and the refusal codes' meanings are unchanged." The provenance line:
"Decided with the operator 2026-09-23; the design notes are the appraisals on
#343, #348, #349 and #350." Those four appraisals were posted at 10:49Z the
same morning, before the milestone existed.

The eight requirements, in short (the table below carries each in full):

- **M18-R1** — the SPEC and the docs match the CLI on where `note:` lines
  go, which verbs read no pointer, and which verbs take `--dry-run`, "with a
  drift test pinning that list"; "no verb's output changes".
- **M18-R2** — "the CLI describes itself truly": three printed citations of
  SPEC §6 move to where the fact now lives, and `--help` shows `identity new
  --no-route`; "refusal codes and exit statuses unchanged".
- **M18-R3** — ship v2.0.2; "protocol version unchanged; lands before any of
  R4–R7 merges".
- **M18-R4** — a struck requirement: a verb posts the struck line as a
  comment, "never editing the body", and "refuses without a link to a
  recorded Decision"; `milestone close` counts `struck` as terminal.
  "Striking belongs to the coordination layer, not QA."
- **M18-R5** — the doc-synthesizer's front-door obligation is general, and
  this hub's specifics move into its `.local.md`.
- **M18-R6** — the housekeeping light path, defined mechanically ("a tool
  defines the target and the diff is the whole decision"), and `roles sync`.
- **M18-R7** — "the record takes the light path"; `DOC_MISSING` and
  `milestone evidence` unchanged; "the `NO_CHECKS` gate is out of scope".
- **M18-R8** — ship v2.1.0: "protocol version 2.1, the flip, the tag and five
  assets".

Six of the eight name the captures they adopt; R3 and R8 are releases and
adopt none. Nine captures in all, each closed by `task finish` after its
task merged.

The ten tasks, in merge order:

- **[#355](https://github.com/radiusred/gh-codecrew/issues/355) / PR
  [#356](https://github.com/radiusred/gh-codecrew/pull/356)** — M18-R2;
  adopts #339 and #334. Three commits, approved first round, merged
  2026-09-23T12:02:28Z. One Decision, one Deviation.
- **[#354](https://github.com/radiusred/gh-codecrew/issues/354) / PR
  [#357](https://github.com/radiusred/gh-codecrew/pull/357)** — M18-R1;
  adopts #337 and #338. Five commits, three review rounds, merged 12:38:38Z.
  Four Decisions, two Deviations (one the operator's).
- **[#358](https://github.com/radiusred/gh-codecrew/issues/358) / PR
  [#360](https://github.com/radiusred/gh-codecrew/pull/360)** — M18-R1, a
  follow-up "found by a blind replay of its review". Two commits, approved
  first round, merged 22:56:00Z. One Deviation (the operator's).
- **[#359](https://github.com/radiusred/gh-codecrew/issues/359) / PR
  [#361](https://github.com/radiusred/gh-codecrew/pull/361)** — M18-R3, the
  v2.0.2 flip. Two commits, approved first round, merged 23:02:53Z; tagged
  and released by 23:04:13Z. One Decision.
- **[#363](https://github.com/radiusred/gh-codecrew/issues/363) / PR
  [#366](https://github.com/radiusred/gh-codecrew/pull/366)** — M18-R5;
  adopts #350. Three commits, approved first round, merged
  2026-09-24T15:52:39Z. Two Decisions (one the operator's).
- **[#362](https://github.com/radiusred/gh-codecrew/issues/362) / PR
  [#367](https://github.com/radiusred/gh-codecrew/pull/367)** — M18-R4;
  adopts #348. Six commits, two review rounds, merged 16:04:23Z. Two
  Decisions, one Deviation.
- **[#364](https://github.com/radiusred/gh-codecrew/issues/364) / PR
  [#368](https://github.com/radiusred/gh-codecrew/pull/368)** — M18-R6;
  adopts #343 and #266. Five commits, three review rounds, merged 16:21:27Z.
  Two Decisions.
- **[#369](https://github.com/radiusred/gh-codecrew/issues/369) / PR
  [#370](https://github.com/radiusred/gh-codecrew/pull/370)** — M18-R7;
  adopts #349. Four commits, two review rounds, merged 17:19:06Z. One
  Decision.
- **[#372](https://github.com/radiusred/gh-codecrew/issues/372) / PR
  [#373](https://github.com/radiusred/gh-codecrew/pull/373)** — carried under
  M18-R8 by the coordination layer's Decision. Eleven commits, three review
  rounds, merged 18:31:04Z. Eight Decisions, one Deviation.
- **[#371](https://github.com/radiusred/gh-codecrew/issues/371) / PR
  [#374](https://github.com/radiusred/gh-codecrew/pull/374)** — M18-R8, the
  v2.1.0 flip. Four commits, approved first round, merged 18:41:56Z; tagged
  and released by 18:43:33Z. Three Decisions.

Then the housekeeping pull request #371's plan names as its Part C:
**[#375](https://github.com/radiusred/gh-codecrew/pull/375)**, `chore:
regenerate the contract history for v2.1.0` — one commit, no task, approved
by the reviewer seat's holder at 18:50:44Z and rebase-merged by its author at
18:51:12Z.

Every one of the eleven was rebase-merged by `radiusred-cody[bot]`, #375
by its author on the light path.

**The milestone ran in three sittings over thirty-one hours**: 11:42Z to
12:38Z on 2026-09-23, 22:49Z to 23:59Z the same night, and 15:43Z to 19:05Z
the next afternoon. From the
milestone opening to the v2.1.0 flip merging: **one day, six hours,
fifty-nine minutes and thirty-seven seconds**; to the QA comment, one day,
seven hours, twenty-two minutes and forty-five seconds. The corrections
(#355, #354) merged within fifty-seven minutes of the milestone opening. The
next task, #358, was created ten hours, eleven minutes and seven seconds
after #357 merged, and #358 and the v2.0.2 release were done inside the next
fifteen minutes. The three protocol tasks for R4–R6 were created at
23:54:15Z–23:54:24Z the same night and their plans written by 23:59:42Z; the
operator's answers were relayed at 15:43:47Z–15:44:02Z the next afternoon,
fifteen hours and forty-four minutes later, and all three pull requests
merged within thirty-eight minutes of the first answer. R7 and R8 followed
in the next two and a half hours. Nothing on the trail says what happened in
either gap; the last section carries it.

**Pull requests ran in parallel twice.** #356 and #357 were open together
from 11:51:16Z to 12:02:28Z on 2026-09-23. #366, #367 and #368 were all open
from 15:51:41Z to 15:52:39Z on 2026-09-24, and #367 and #368 together until
16:04:23Z. [M17's record](17-the-frontier-builds-the-coordinator.md#goal-and-outcome)
reports no two pull requests open at once;
[M16's](16-bedding-in-the-venue-seam-and-the-reference.md#goal-and-outcome)
reports four.

What exists afterwards that did not before. Two releases.
[v2.0.2](https://github.com/radiusred/gh-codecrew/releases/tag/v2.0.2),
published 2026-09-23T23:04:13Z, carries M16's venue seam and `CLI.md` and
M18-R1/R2's corrections at protocol 2.0.
[v2.1.0](https://github.com/radiusred/gh-codecrew/releases/tag/v2.1.0),
published 2026-09-24T18:43:33Z, implements protocol 2.1, which the CHANGELOG
calls "The first protocol minor. 2.1 is additive". In it: `gh codecrew
milestone strike`, and `milestone close` counting a struck requirement as
terminal; SPEC §4's Housekeeping subsection and its per-tier approval table;
`gh codecrew roles sync`, which writes the embedded contracts and
`.codecrew/AGENTS.md` where they are absent or an earlier release's text and
refuses a fork; `status` reporting a missing contract and the agents file;
the milestone record on the light path; the general front-door obligation;
and a version check at the top of the scaffolded `.codecrew/AGENTS.md`. SPEC
§10's catalogue went from forty-three codes to forty-seven —
`REQUIREMENT_UNDECLARED`, `DECISION_UNRECORDED`, `CONTRACT_FORKED`,
`AGENTS_FORKED` — "with none repurposed", as the CHANGELOG says and the
reviewer checked on #374.
This hub's pointer reads `codecrew: "2.1"`.

## Decisions

Twenty-six Decision records. Two are on the milestone issue, both the
coordination layer's, posted by the operator's account. Twenty-four are on
the tasks: one posted by the operator's account on #363, and twenty-three by
the implementer seat (`radiusred-cody[bot]`) — of which five record the
operator's answers to a task's ask-the-human points, each saying it was
"relayed by the coordination layer". Three carry a parenthetical qualifier;
none is a gate resolution. #358 carries no Decision.

### The corrections: which list, which test, which citation

#354's [first Decision](https://github.com/radiusred/gh-codecrew/issues/354#issuecomment-5794285140)
found that the three-verb `--dry-run` list its capture described was no
longer on `main` — #327 had replaced it with a pointer to `CLI.md` — so
"R1's 'lists `migrate` among the `--dry-run` verbs' means adding the list
back, now with four verbs, and pinning it". **Trade-off:** "the introduction
once again names verbs that CLI.md also covers. The drift test is what stops
that list going stale." **Rejected:** restoring the old sentence, whose
"every gate in order" did not describe `milestone new --dry-run`.

The [second](https://github.com/radiusred/gh-codecrew/issues/354#issuecomment-5794285389)
is the drift test: it reads the package's Go source with `go/ast` and takes
each function pairing `NewFlagSet("<verb>", …)` with `Bool("dry-run", …)`.
**Trade-off:** "the test depends on how the source is written".
**Rejected:** comparing against `--help`'s markers, "text about the flag and
not the flag set R1 names", and running each verb, "since most verbs load the
pointer and reach GitHub before they parse it". The reviewer's round one
found the gap in that reading — `BoolVar` was not recognised — and the fix is
in the review rounds below.

The [third](https://github.com/radiusred/gh-codecrew/issues/354#issuecomment-5794285599)
keeps "only" out of SPEC §6 and adds `identity token`'s stale-installation
note to `CLI.md`'s Output channels. **Rejected:** an exclusive list in §6,
which "would therefore be wrong again in a new way". It names the same flaw
in `CLI.md`'s "three things, and only these" and leaves it "for the reviewer
to weigh because it is outside R1's text". The reviewer weighed it and asked
for it fixed.

The [fourth](https://github.com/radiusred/gh-codecrew/issues/354#issuecomment-5794645219),
after review, changes the test for SPEC §10's exception list from verbs that
"read no pointer" to verbs that "neither require nor validate" one, because
`identity token` reads one best-effort to prefer the hub's owner.
**Rejected:** "dropping `identity token` from the exception list, which would
wrongly say it can raise the common pointer refusals", and changing
`hubOwner`, "which would change behaviour and is outside R1".

#355's [one Decision](https://github.com/radiusred/gh-codecrew/issues/355#issuecomment-5794256419)
fixes each of the three citations: `GH_UNREACHABLE` to `(CLI.md, Common
refusals)`, both `HUB_UNREADABLE` details to `(SPEC §5; CLI.md, Common
refusals)`, both operator-confirmation comments to `SPEC §5`. **Trade-off:**
a section name rather than a link, because "a full URL would roughly double
the GH_UNREACHABLE line". **Rejected:** the full URL, SPEC §5 alone for
`HUB_UNREADABLE`, and `CLI.md`'s `task finish` section for the confirmations
("it states the flag, not the tier").

### The v2.0.2 flip leaves one example behind

#359's [one Decision](https://github.com/radiusred/gh-codecrew/issues/359#issuecomment-5804336227)
leaves `CLI.md`'s `version` example at `v2.0.1 (protocol 2.0)`: the flip moves
the version "in exactly the places #320 moved it", and `CLI.md` did not exist
then. **Trade-off:** "from the tag onward the reference's example is a patch
behind what the binary prints". **Rejected:** flipping it, which "widens the
flip past the #320 precedent", and a placeholder. The reviewer read it as
"an illustration of the stable `<tag> (protocol <major.minor>)` shape, not a
claim that 2.0.1 is current". #371's plan changed it to `v2.1.0 (protocol
2.1)` in the next release, once the protocol itself had moved.

### Striking: the operator's nine answers

#362's plan put nine numbered ask-the-human points, each with options and a
recommendation, and says the operator "accepted every recommendation". The
[Decision recording them](https://github.com/radiusred/gh-codecrew/issues/362#issuecomment-5817343992),
qualified "(the operator's answers to the ask-the-human points, davison,
2026-09-24, relayed by the coordination layer)", quotes the relay: the verb is
`milestone strike <n> <ID> --decision <url>`; the Decision may be on "the
milestone or one of its sub-issue tasks"; it must be a "Decision or Gate
resolved record that also names the ID"; "a verified strike is terminal
regardless of QA verdicts, and `--reinstate` ships now"; "close counts only
the coordinator holder's struck/reinstated lines and re-verifies each link";
a still-bold struck-through ID stays a requirement; a re-run posts nothing;
`milestone close` gathers the milestone issue's own records; and the two new
codes as named. **Trade-off:** "`--reinstate` widens the grammar by one line
form in a milestone about scaling down, and it is accepted because otherwise
the only way to undo a strike is to delete a comment. Counting only the
coordinator holder's lines brings the actor back at close, while the verb
itself still checks none." **Rejected:** among others, "latest-comment-wins
across verdicts and strikes", "counting any author, and trusting a line
without re-checking its link", and "reading a body strikethrough as struck".

The gathering answer (8(a)) is why the two coordination-layer Decisions on
#353, below, reach the record's raw material at all: before #367, `milestone
close` read Decisions from the tasks and their pull requests only.

The [second](https://github.com/radiusred/gh-codecrew/issues/362#issuecomment-5817344244)
changes the coordinator contract's "never … post a verdict" to "never … post
a QA verdict", because the coordinator strikes and R8 calls the struck line
"the struck verdict"; and it touches `docs/introduction.md` a task early,
because R1's drift test would fail CI once `milestone strike` defines
`--dry-run`. **Rejected:** "leaving the introduction for R8, which would fail
CI on this PR".

### The front door: six answers, and an approval that stood

#363's [first Decision](https://github.com/radiusred/gh-codecrew/issues/363#issuecomment-5817341651)
records the operator's acceptance of all six recommendations: the general
bullet names no document and defines the front door by function, "the pages
a newcomer reads first to learn what the project is and what works"; it
points at the `.local.md` and says "where the milestone changed no such
claim, there is nothing to refresh"; `docs/extensions.md` gains a dated
example; no SPEC change; no guard test; and this hub's extension does not
send the doc-synthesizer to codecrew.works. **Trade-off:** "a contract that
names no document leaves an unextended hub's doc-synthesizer to judge what
its front door is; the pointer to the `.local.md` is where a hub removes that
judgment." **Rejected:** "naming 'README' by example (an assumption about the
hub's shape)".

The [second](https://github.com/radiusred/gh-codecrew/issues/363#issuecomment-5817476509),
posted by the operator's account at 15:52:11Z, is the coordination layer's:
"checky's approval on PR #366 stands through the commit-message reword; no
second review round." **Trade-off:** the approval "was given on head 3679da5;
the force-push (rewrapping two commit bodies that failed
`body-max-line-length`) produced head f029762 with the identical tree
`91ebaa7f0bd5a2928eb1195fb59dc24a44b327ee`, verified by the coordination
layer, and GitHub kept the approval". **Rejected:** "re-dispatching the
reviewer for a message-only change — it would re-judge a diff it has already
approved." The pull request merged twenty-eight seconds later.

### Housekeeping and `roles sync`: nine answers and one engineering call

#364's plan put nine points. Its
[first Decision](https://github.com/radiusred/gh-codecrew/issues/364#issuecomment-5817345471)
quotes the relayed answers: "1(a) text only, no verb (file the `chore finish`
backlog capture you offered …); 2(a) the PR's author rebase-merges once the
tier's approval / solo confirmation is on the PR …; 3(a) SHA-256 table of
every tagged contract text from v0.1.0 …; 4(a) all-or-nothing
CONTRACT_FORKED …; 5(a) literal — a tool must define the target (write it so
R7 can add the record's clause later); 6(a) coordinator contract changes
here; 7(a) branch-cut on the default branch; 8(a) CONTRACT_FORKED; 9(a) green
checks where any are reported." **Trade-off:** "Text-only enforcement (1a)
keeps R6 to what the requirement and the appraisal on #343 ask for, at the
cost that in a private free-plan repo only the contract stops an unreviewed
housekeeping merge; author-merges (2a) keeps the path drivable by an
orchestrator with no human in each loop." **Rejected:** "A `chore finish`
verb enforcing the review gate and checks (1b) — deferred to a backlog
capture", which is #365; operator-only or reviewer merges; full embedded
texts or a network fetch for the release history; partial writes on a fork;
"admitting tool-less corrections (5b)".

The [second](https://github.com/radiusred/gh-codecrew/issues/364#issuecomment-5817437553)
is the implementer's: `roles sync` keeps each file's form, stamped or not,
and classifies a contract "stamp-stripped and with CRLF read as LF, so a
Windows checkout of an unedited release text is recognised rather than called
a fork". **Rejected:** always stamping, never stamping, and "comparing bytes
exactly".

### The record's path: seven answers

#369's [one Decision](https://github.com/radiusred/gh-codecrew/issues/369#issuecomment-5818613347)
records the answers to its seven points, every recommendation accepted: the
record is a "named closed second class" in SPEC §4 Housekeeping, "R6's tool
sentence untouched"; a record needing its own Decision "falls back to the
task path"; `docs:` commits with the milestone issue as the reference; the
doc-synthesizer's and close-time Deviations "go on the milestone issue"; the
`DOC_MISSING` gate's looseness is "DELIBERATE — the gate checks presence";
the reviewer checks the record's citations, with "NO capture for extending
milestone evidence"; and the author deletes the branch on merge. It also
records that the operator "asked that every fix from the phase-1 'wrong or
underspecified' list be folded in", and names them. **Trade-off:**
"Deliberate looseness keeps the close gate as it is and leaves the record's
review to the contracts, as for every housekeeping PR, at the cost that a
direct commit of the record passes `milestone close` unseen; the reviewer, not
a verb, checks the record's citations." **Rejected:** widening R6's sentence
to "or the protocol", `chore:` record commits, "a provenance check in
`milestone close` (5b)", "extending `milestone evidence` to read the document
(6b)", and "a `record/` branch convention the sweep learns (7b)".

This document is delivered under that Decision. The same answer to point 5
asked for "a comment on #365 (as cody) noting the milestone record as a
second consumer of a light-path verb — no separate capture"; it was posted on
#365 at 17:06:39Z.

### v2.1.0: the one recommendation not taken

#371's plan put six points, the first with four parts. The
[Decision recording the answers](https://github.com/radiusred/gh-codecrew/issues/371#issuecomment-5819260879)
quotes the relay: `init` writes `"2.1"`; this hub's pointer moves by hand in
the release task "(#114 precedent)"; "NOTHING for existing 2.0 hubs"; and on
1(d), where the plan recommended one CHANGELOG sentence telling operators to
upgrade every seat's binary and a capture for a later binary note, "the
operator judged a CHANGELOG sentence insufficient (neither a skipping
operator nor a coordinating agent reads it) — so a NEW task, #372, carried
under M18-R8, lands BEFORE your Part A". It is the one answer across the five
protocol tasks' thirty-seven numbered points that is not the plan's
recommendation. **Trade-off:** "a 2.0 pointer is left alone everywhere except
this hub, so the fleet will show mixed `"2.0"`/`"2.1"` pointers. That is
accepted because SPEC §5 will say an earlier minor of the same major is
current." **Rejected:** "a CHANGELOG sentence alone, a `status` note for a 2.0
pointer, `migrate` rewriting a same-major minor, and a binary forward note."

The [second](https://github.com/radiusred/gh-codecrew/issues/371#issuecomment-5819261176)
keeps `migrate`'s "already on the protocol 2.1 layout" line for a pointer
that says 2.0, as "true of the layout, which 2.1 does not change".
**Rejected:** rewording it, which "would change a verb's output for no
behaviour". The [third](https://github.com/radiusred/gh-codecrew/issues/371#issuecomment-5819261383)
has four tests read `protocolVersion` while `TestVersionCmd` keeps "its
deliberate pin, which moves to `"2.1"`, so a bump still has to be made on
purpose". **Trade-off:** "the migrate tests stop catching an accidental change
to the constant."

### The widening: a version check, then a delivery path for it

The coordination layer's
[first Decision on #353](https://github.com/radiusred/gh-codecrew/issues/353#issuecomment-5819250070),
at 17:50:02Z, qualified "(coordination layer, with the operator,
2026-09-24)", widens M18-R8 "by one task, #372, landing before v2.1.0's flip
(#371)": the scaffolded `.codecrew/AGENTS.md` and the coordinator contract
gain a version check, and "'the project's floor' (named in the coordinator's
wake step, defined nowhere) becomes the pointer's protocol version".
**Trade-off:** #371's plan "verified that a binary older than its hub's
pointer fails closed silently", and "the CHANGELOG reaches neither an
operator who skips it nor a coordinating agent, which has no obligation to
read it — the check has to be in the text every dispatched agent loads."
**Rejected:** a CHANGELOG sentence alone; "`roles sync` also refreshing
`.codecrew/AGENTS.md`" (new verb behaviour, "not taken now; existing hubs
receive the new text through `init`/`migrate` or by hand"); deferring to M19.

Five minutes later a
[correction](https://github.com/radiusred/gh-codecrew/issues/353#issuecomment-5819321400),
labelled `**Correction (coordination layer, 2026-09-24)**` and so not itself a
gathered record, withdrew that premise: "#372's implementer verified on a
built binary that neither does … The option put to the operator was worded on
the wrong premise; it is being re-put." The
[second Decision](https://github.com/radiusred/gh-codecrew/issues/353#issuecomment-5819365327),
at 17:58:05Z, is the re-put answer: "`roles sync` also writes
`.codecrew/AGENTS.md` — CodeCrew's own file — under the same never-clobber
rule as the contracts", riding #372 before the flip. **Trade-off:** "new verb
behaviour inside a release's scope … in exchange, the rule reaches existing
hubs without anyone reading a CHANGELOG". **Rejected:** "hand edits across the
fleet (the thing the check exists to avoid); deferring to M19."

#372's own eight Decisions carry both phases. The version check is the
scaffold's
[first bullet](https://github.com/radiusred/gh-codecrew/issues/372#issuecomment-5819309687),
because "reading the contract is itself a verb (`roles show`), so a check
placed after it would not be 'before the first verb'". The
[scaffold's Go comment is corrected](https://github.com/radiusred/gh-codecrew/issues/372#issuecomment-5819309978):
it said "a later init or migrate rewrites it whole", and a built binary
showed both write the file only when absent — "the old comment was the claim
the #353 Decision leaned on". Phase two's four, all at 18:06:27Z–18:06:30Z:
`roles sync` delivers the agents file
[in spokes too](https://github.com/radiusred/gh-codecrew/issues/372#issuecomment-5819485840),
because "hub-only would leave the version check out of most dispatches"; a
forked agents file alone
[refuses the new `AGENTS_FORKED`](https://github.com/radiusred/gh-codecrew/issues/372#issuecomment-5819486208)
rather than stretching `CONTRACT_FORKED`, since "a widened meaning is a
repurpose in all but name"; the file is
[named by its path](https://github.com/radiusred/gh-codecrew/issues/372#issuecomment-5819486534),
with the consequence that "this hub's hand-written file is a fork, so here
bare `roles sync` refuses `AGENTS_FORKED`"; and released scaffold texts reach
the history table by
[evaluating the `agentsScaffold` constant](https://github.com/radiusred/gh-codecrew/issues/372#issuecomment-5819486863)
from each tag's source with a narrow evaluator. **Rejected:** "building each
tag and running `init`" and "hard-coding three hashes".

The last two answer the reviewer. After round one, a spoke
[reads the hub's `codecrew:` with `gh api`](https://github.com/radiusred/gh-codecrew/issues/372#issuecomment-5819619250)
under the session's own auth, falling back to the seat identity.
After round two, the check
[finishes before any verb but `version` and `identity token`](https://github.com/radiusred/gh-codecrew/issues/372#issuecomment-5819760610),
with the mint-and-retry inside it; the Decision carries a `**Supersedes:**`
paragraph naming the round-one ordering, "which let `roles show` run first".

## Deviations

Seven Deviation records: six on the tasks — two the operator's, four the
implementer seat's — and one on the milestone issue, this record's.

**[The reviewer and qa seats' model changes, and rides #354's pull request](https://github.com/radiusred/gh-codecrew/issues/354#issuecomment-5794231035).**
The operator's, at 11:46:41Z, eighteen seconds after #354 was created: both
seats' `model` moves from `gpt-5.5` to `gpt-5.6-sol` in
`.codecrew/config.yml` and the README's snippet. **Why:** "the operator chose
to carry it in the milestone's first hub PR rather than a separate `chore:`
PR, the same way #351 carried the implementer's model change in M17. The #343
light path is still only a proposal (M18-R6)." It closes: "Review rounds from
this PR on run under `gpt-5.6-sol`." The commit is
[aa9bfa3](https://github.com/radiusred/gh-codecrew/commit/aa9bfa3) on `main`,
`chore(config): the reviewer and qa seats are dispatched under gpt-5.6-sol
(#354)`; #358's is
[a24d14c](https://github.com/radiusred/gh-codecrew/commit/a24d14c).

**[The reviewer moves again, and rides #358's](https://github.com/radiusred/gh-codecrew/issues/358#issuecomment-5804247993).**
The operator's, at 22:49:55Z: the reviewer seat moves to `gpt-5.6-terra`; "The
qa seat stays on `gpt-5.6-sol`." **Why:** "a blind replay (network off, GitHub
state as it stood at the time) of two review rounds compared the models. On
PR #357, terra took 6m36s against sol's 10m42s and found the same two
substantive defects; sol also found this task's introduction defect, and both
missed a third. On PR #356, terra approved correctly in 4m57s; sol's live
review took 11m04s. … QA stays on sol as the check at the end of the
milestone. The next reviews are the trial." #358 itself is the introduction
defect the replay found.

**[`docs/working-offline.md` follows the citation](https://github.com/radiusred/gh-codecrew/issues/355#issuecomment-5794266333).**
#355's plan did not list the page, but it quotes the `GH_UNREACHABLE` line and
a test pins the quote. **Why:** "the drift test is doing its job".

**[A sweep beyond the plan](https://github.com/radiusred/gh-codecrew/issues/354#issuecomment-5794885176).**
Commit `5bb6b60` on #357 corrects three more per-verb sentences in `CLI.md` —
`init`, `migrate`, `task finish` — each of which "claimed a channel carried
only certain output, and the binary showed otherwise". **Why:** "the
coordination layer asked for this fix and the sweep before round two of
review." It was posted at 12:33:23Z, after the reviewer's round two asked for
it as a comment; the review rounds below carry that.

**[The strike receipt, and the introduction's verb list](https://github.com/radiusred/gh-codecrew/issues/362#issuecomment-5817450907).**
`milestone strike`'s receipt names the milestone issue and the posted line,
not the new comment's URL, and the introduction's verb list gains the four
milestone verbs. **Why:** "the venue's `Comment` call returns no URL, and
adding a return value to the interface for a receipt line is not worth
widening the venue in this task"; and a `--dry-run` sentence naming a verb
the verb list omits "would be a stale claim of the kind the doc-synthesizer
contract calls a defect".

**[A test and its commit](https://github.com/radiusred/gh-codecrew/issues/372#issuecomment-5819311780).**
This hub's `.codecrew/AGENTS.md` step went into #373's first commit, not the
docs commit. **Why:** "one test pins the version check in both agents files,
so splitting them would leave the first commit red. No change to scope."

**[The record takes harness facts from its dispatch brief](https://github.com/radiusred/gh-codecrew/issues/353#issuecomment-5820516906).**
This record's own, at 19:11:43Z. The brief states that every implementer
session ran under Claude Opus 5.5 and every reviewer and QA round in Codex,
with the models above; the trail carries the commit trailers, the routing
table and the operator's two Deviations, and "No review or QA comment names
its own harness or model." **Why:** "The brief is not on the trail, so this
comment puts its account there, and the record cites it rather than the
brief."

No task's pull request body reports a Deviation that is not also a comment:
#360's says "No work outside the Plan, so no Deviation of mine was needed",
#368's "No deviations from the Plan", and #370's "No deviations from the
plan".

## The gates

**#353's Gates section keeps the scaffold's placeholder line and adds three
of its own**, in the edit nine seconds after the issue was created:

- "R3 (v2.0.2) merges before any of R4–R7."
- "R6 lands before R7 (R7 needs R6's path)."
- "R8 is the last task before the record, which is the first record to take
  R7's path."

[M17's record](17-the-frontier-builds-the-coordinator.md#the-gates) says its
Gates section was left as the placeholder, as do
[M16's](16-bedding-in-the-venue-seam-and-the-reference.md#the-gates) and
[M15's](15-v2-0-1-what-the-fleet-migration-taught.md#the-gates).

**All three held, by the timestamps.** v2.0.2 was published at
2026-09-23T23:04:13Z and the first R4–R7 merge was #366 at
2026-09-24T15:52:39Z; QA's M18-R3 verdict checks the same: "Its
2026-09-23T23:04:13Z publication followed PR #361 and preceded every R4–R7
merge on 2026-09-24." #368 (R6) merged at 16:21:27Z and #370 (R7) at
17:19:06Z. The last task pull request, #374, merged at 18:41:56Z; #375
followed at 18:51:12Z, which #371's plan had placed there: "Part C is
housekeeping, not a task, so 'R8 is the last task before the record' still
holds."

**No gate was raised.** No task used `gh codecrew checkpoint`, no issue
carries `cc:needs-decision`, and the dry run reports `gate no gate raised:
ok`. The human judgment the protocol tasks needed went through their plans
instead: #362, #363, #364, #369 and #371 each wrote numbered ask-the-human
points with options and a recommendation before the first commit, and each
records the answers as a Decision before `task start` ran. The plan edits
and the starts, by timestamp: #362's plan at 23:58:30Z, answers at 15:43:56Z,
started 15:44:12Z; #363's at 23:56:59Z, 15:43:47Z, 15:43:57Z; #364's at
23:59:42Z, 15:44:02Z, 15:44:21Z; #369's at 16:52:17Z, 17:06:16Z, 17:06:27Z;
#371's at 17:36:15Z, 17:50:48Z, and a start at 18:31:43Z that waited for #372
to merge. #372's own section reads "None", because both phases were decided
on #353. The four earlier tasks each resolved theirs to none.

**The two releases' operator legs are on the milestone issue.** The
[v2.0.2 leg](https://github.com/radiusred/gh-codecrew/issues/353#issuecomment-5804411355)
at 23:05:23Z, one minute and ten seconds after the release: `upgraded from
v2.0.1 to v2.0.2`, all five `roles diff` clean, `migrate --dry-run` writing
nothing, and "No migration leg this release: the protocol is still 2.0. R4–R7 are now free
to land". The
[v2.1.0 leg](https://github.com/radiusred/gh-codecrew/issues/353#issuecomment-5820095959)
at 18:44:54Z: `upgraded from v2.0.2 to v2.1.0`, `version` printing `v2.1.0
(protocol 2.1)` "satisfying #372's version check against this hub's
`codecrew: "2.1"`", and `status` showing "the expected `agents file drift:`
line for this hub's hand-written `.codecrew/AGENTS.md`".

What gated the work: CI on every pull request — `Lint commit messages` and
`Go build and test` — the reviewer seat holder's approval on each of the
eleven, the release workflow's new contract-history guard on the v2.1.0 tag,
and the QA verdicts.

## The review rounds

Eleven pull requests, **twenty review submissions**. Nineteen are by the
reviewer role holder (`radiusred-checky[bot]`), an App distinct from the
authoring seat: **eleven approvals, one per pull request, all standing**,
and **eight change requests**, each resolved in a later round. The twentieth
was posted under the operator's login and dismissed. Six pull requests were
approved in their first round; #367 and #370 took two; #357, #368 and #373
took three.

- **PR [#356](https://github.com/radiusred/gh-codecrew/pull/356)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/356#pullrequestreview-5290669107),
  "M18-R2 is satisfied with no findings". It proves the tests bite: "removing
  CLI.md's `[--no-route]` line made `TestReferenceSynopsesMatchHelp` fail.
  Restoring the old §6 citations made `TestResolveRolesFailsClosed` fail".
- **PR [#357](https://github.com/radiusred/gh-codecrew/pull/357)** — three
  rounds. The
  [first change request](https://github.com/radiusred/gh-codecrew/pull/357#pullrequestreview-5290799220)
  has three findings: `CLI.md`'s "stderr carries three things, and only these"
  is disproved by the built binary; SPEC §10 and `CLI.md` say `identity token`
  "read[s] no pointer" when `hubOwner` reads one; and the drift test
  "recognizes only `FlagSet.Bool("dry-run", ...)`", so a `BoolVar` dry-run
  flag added to `role` "left `TestIntroductionListsTheDryRunVerbs` green". The
  [second](https://github.com/radiusred/gh-codecrew/pull/357#pullrequestreview-5290984078)
  closes all three and asks for one record: commit `5bb6b60`'s rationale
  "currently exists only in the editable PR-body section 'Review round 1',
  which does not satisfy the comment record." The implementer's Deviation
  followed forty-seven seconds later, and the
  [approval](https://github.com/radiusred/gh-codecrew/pull/357#pullrequestreview-5291039512)
  at 12:38:03Z cites it by number.
- **PR [#360](https://github.com/radiusred/gh-codecrew/pull/360)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/360#pullrequestreview-5297736120),
  checking the channel claim on the built binary and the routing snippet
  "byte-identical to the config `roles:` section".
- **PR [#361](https://github.com/radiusred/gh-codecrew/pull/361)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/361#pullrequestreview-5297778137):
  "The eleven former Unreleased entries match main byte-for-byte", and "Each
  claim in the new lead is supported by an entry below it".
- **PR [#366](https://github.com/radiusred/gh-codecrew/pull/366)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/366#pullrequestreview-5306767994)
  at 15:48:40Z: the embedded contract "now names neither a document nor this
  hub's release, verbs, or refusal codes". The head it approved failed the
  commit lint; the force-push and the Decision that the approval stood are
  above.
- **PR [#367](https://github.com/radiusred/gh-codecrew/pull/367)** — two
  rounds. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/367#pullrequestreview-5306873977)
  finds `docs/first-milestone.md` still saying "every requirement needs a QA
  verdict", "now false", and a test that claims strike-after-reinstate but
  "stops after asserting the reinstated post". The
  [approval](https://github.com/radiusred/gh-codecrew/pull/367#pullrequestreview-5306944631)
  names the new `TestStrikeReinstateStrikeEndsStruck`.
- **PR [#368](https://github.com/radiusred/gh-codecrew/pull/368)** — three
  rounds. The
  [first change request](https://github.com/radiusred/gh-codecrew/pull/368#pullrequestreview-5307029651):
  "`GH_UNREACHABLE`’s offline detail is now false" — it named `roles show/diff`
  as the hub-local verbs, and this pull request makes `roles sync` one too. The
  [second](https://github.com/radiusred/gh-codecrew/pull/368#pullrequestreview-5307088334)
  is about the pull request body alone: its "For the reviewer" section still
  described the detail as left alone, "Please update or remove that paragraph
  so the final PR description accurately describes the head." The
  [approval](https://github.com/radiusred/gh-codecrew/pull/368#pullrequestreview-5307136542)
  followed on the same head.
- **PR [#370](https://github.com/radiusred/gh-codecrew/pull/370)** — two
  rounds. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/370#pullrequestreview-5307723665)
  finds the change "otherwise coherent" and asks for one thing: a commit
  subject starting `docs: CLI.md, …`, which "Commitlint rejects … as
  sentence/start case". The
  [approval](https://github.com/radiusred/gh-codecrew/pull/370#pullrequestreview-5307777088)
  verifies the new head "has the same tree as the first-round head".
- **PR [#373](https://github.com/radiusred/gh-codecrew/pull/373)** — three
  rounds. The
  [first change request](https://github.com/radiusred/gh-codecrew/pull/373#pullrequestreview-5308354517)
  finds that the scaffold asks a spoke to compare against the hub's pointer
  "but gives no way to obtain the hub pointer before the first CodeCrew verb",
  so "An agent can therefore pass the local check without ever checking the
  project's floor"; and a stale bare `roles sync` example in `CLI.md`. It
  also reports building v2.0.2 and running `init` to check the scaffold's
  hash against the history row. The
  [second](https://github.com/radiusred/gh-codecrew/pull/373#pullrequestreview-5308440908)
  finds the fix's private-hub fallback still let `roles show` run first, and a
  117-character source line. The
  [approval](https://github.com/radiusred/gh-codecrew/pull/373#pullrequestreview-5308536908):
  "The scaffold now permits only `gh codecrew version` and `gh codecrew
  identity token` before the check".
- **PR [#374](https://github.com/radiusred/gh-codecrew/pull/374)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/374#pullrequestreview-5308645999).
  It counts the catalogue against the previous tag — "SPEC §10 contains 43
  codes at v2.0.2 and 47 here" — and accepts the window the pull request body
  flagged, between the pointer moving to 2.1 and the operator installing
  v2.1.0: "it is the operator's recorded choice … and the window ends with
  publication and the operator upgrade."
- **PR [#375](https://github.com/radiusred/gh-codecrew/pull/375)** —
  [approved](https://github.com/radiusred/gh-codecrew/pull/375#pullrequestreview-5308758223)
  as housekeeping: "after fetching tags, rerunning `scripts/contract-history`
  at `origin/main` generated a byte-identical table with only v2.1.0's five
  contract rows and its `.codecrew/AGENTS.md` scaffold row."

**The dismissed review.** At 12:31:37Z on #357 a change request was posted
under the operator's login (`davison`), with a body identical to the
reviewer App's round two, which followed at 12:32:36Z. It was dismissed at
12:32:37Z by `radiusred-checky[bot]` with the message "Dismissed because the
reviewer session accidentally posted this duplicate under the operator login;
review 5291006191 is the reviewer App record." No review with that ID exists
on the pull request — the reviews API returns 404 for it — and the App's
round-two review is
[5290984078](https://github.com/radiusred/gh-codecrew/pull/357#pullrequestreview-5290984078).
No comment elsewhere records the event.

**What the eight change requests were about.** They carry thirteen
findings. Five are a false claim in text: `CLI.md`'s exclusive stderr list,
the "read no pointer" category, the quickstart's "every requirement needs a
QA verdict", the offline detail, and `CLI.md`'s bare `roles sync` example.
Two are a test that did not cover what it claimed: the drift test's `BoolVar`
blind spot, and the missing strike-after-reinstate sequence. Two are the
version check's ordering — that it can actually run before the first verb,
which is what #372 exists for. Two are about the trail: a rationale living
only in a pull request body, and a pull request body left stale. Two are
mechanics: a commit subject's case and a 117-character source line.

## QA: one round, eight verdicts

The qa role holder (`radiusred-testy[bot]`) wrote
[one comment on #353](https://github.com/radiusred/gh-codecrew/issues/353#issuecomment-5820419251)
at 19:05:04Z, thirteen minutes and fifty-two seconds after #375 merged, with
one verdict per requirement and a closing cross-cutting probe. There is no
earlier QA comment and so no supersession.
[M17's record](17-the-frontier-builds-the-coordinator.md#qa-two-rounds-one-remedy-and-the-verdict-that-supersedes)
describes two rounds and a remedy;
[M16's](16-bedding-in-the-venue-seam-and-the-reference.md#qa-one-round-six-verdicts)
one round.

Each verdict opens "Merged-main `go test ./...` passed", then says what the
shipped tests prove and what they assume, then what the seat checked beyond
them — mostly against the released binary rather than a source build:

- **M18-R1** — "on released v2.1.0, `status 2>/dev/null` retained its report
  and advisory `note:`, `status >/dev/null` was silent, and `identity token
  >/dev/null` left its receipt on stderr"; the #354/#358 commits "changed
  documentation/tests but no production Go".
- **M18-R2** — `git grep` at v2.1.0 shows each citation where the task put
  it, and "Released `identity new --help` showed `[--no-route]`".
- **M18-R3** — "I downloaded and ran the released Linux v2.0.2 asset and saw
  `v2.0.2 (protocol 2.0)`", with the ordering check quoted under the gates.
- **M18-R4** — "against the live hub, dry runs refused an unknown ID with
  `REQUIREMENT_UNDECLARED` and a non-Decision comment for M18-R4 with
  `DECISION_UNRECORDED`, explicitly writing nothing."
- **M18-R5** — "direct inspection found the embedded front-door obligation
  general", and released `roles show doc-synthesizer` composing "the embedded
  contract first and both local sections afterward".
- **M18-R6** — "in committed scratch fixtures, released v2.1.0 named a v2.0.2
  QA contract in dry-run, cut `codecrew-roles-sync`, and added exactly one
  clean local commit, while a hand fork refused `CONTRACT_FORKED` with file
  and HEAD unchanged"; and "PR #375 was a taskless `chore:` PR, approved by
  Checky and author-merged by Cody."
- **M18-R7** — "`git diff --exit-code v2.0.2 v2.1.0 -- internal/cli/evidence.go`
  was clean, both tags use the same `HasMilestoneDoc` presence gate, the
  verdict regex remains `satisfied|not satisfied|untestable`".
- **M18-R8** — "a v2.0.2-initialized spoke was upgraded by v2.1.0 `roles
  sync` in one agents-file commit; this hub's hand-written file refused
  `AGENTS_FORKED` unchanged".

The cross-cutting probe: "all ten M18 tasks and captures #337, #338, #339,
#334, #348, #350, #343, #266, and #349 were closed; #365 remained open as the
intended capture", and `status` showed "only the expected agents-file drift
line".

**Every standing verdict is `satisfied`, and all eight are in that one
comment.** No requirement was struck, so `milestone strike` — shipped in this
milestone — was exercised on this trail only by QA's dry runs.

## Requirement outcomes

The status column is the **standing** verdict word as the qa role holder
wrote it, verbatim and unqualified. All eight stand in comment
[5820419251](https://github.com/radiusred/gh-codecrew/issues/353#issuecomment-5820419251).

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M18-R1 — the SPEC and the docs match the CLI: SPEC §6's channel paragraph says `note:` lines go to the verb's stdout beside its report, and names the notes that go to stderr (the ones raised while the pointer is read, and `identity token`'s); SPEC §10's "any verb" sentence names `identity webhook` among the verbs that read no pointer; `docs/introduction.md` lists `migrate` among the `--dry-run` verbs, with a drift test pinning that list to the verbs whose flag set defines the flag; no verb's output changes (adopts #337, #338, and #334's §10 half) | [#354](https://github.com/radiusred/gh-codecrew/issues/354) / PR [#357](https://github.com/radiusred/gh-codecrew/pull/357) and the follow-up [#358](https://github.com/radiusred/gh-codecrew/issues/358) / PR [#360](https://github.com/radiusred/gh-codecrew/pull/360) | `satisfied` | Both tasks closed; #337 and #338 closed by #354's finish. Three review rounds on #357 (two change requests, one duplicate dismissed), one on #360. Four Decisions (the four-verb list, the `go/ast` drift test, no "only" in §6, "require nor validate"); three Deviations, two of them the operator's routing changes |
| M18-R2 — the CLI describes itself truly: the `GH_UNREACHABLE` and `HUB_UNREADABLE` details and `task finish`'s two operator-confirmation comments cite the section that now carries the fact (CLI.md or SPEC §5) instead of SPEC §6, with the tests that assert them updated; `--help` shows `identity new --no-route`, and CLI.md's synopsis follows through the drift test; refusal codes and exit statuses unchanged (adopts #339, and #334's `--help` half) | [#355](https://github.com/radiusred/gh-codecrew/issues/355) / PR [#356](https://github.com/radiusred/gh-codecrew/pull/356) | `satisfied` | Task closed; #339 and #334 closed by its finish. Approved first round with no findings. One Decision (which section each citation names), one Deviation (`docs/working-offline.md` follows its pinned quote) |
| M18-R3 — ship v2.0.2: the release flip, the tag and five assets, with the CHANGELOG naming M16's venue seam and CLI.md and M18-R1/R2; protocol version unchanged; lands before any of R4–R7 merges | [#359](https://github.com/radiusred/gh-codecrew/issues/359) / PR [#361](https://github.com/radiusred/gh-codecrew/pull/361) | `satisfied` | Task closed. Tag `v2.0.2` on [9e74335](https://github.com/radiusred/gh-codecrew/commit/9e74335); [release](https://github.com/radiusred/gh-codecrew/releases/tag/v2.0.2) published 2026-09-23T23:04:13Z with five assets; operator leg on #353. One Decision (`CLI.md`'s `version` example left at v2.0.1) |
| M18-R4 — a struck requirement: a verb posts `**M<n>-R<k> — struck.** <decision link>` as a comment on the milestone issue, never editing the body, and refuses without a link to a recorded Decision on the milestone or one of its tasks; `milestone close` counts `struck` as terminal, `status` reports it, and the record's requirement table carries the word. Striking belongs to the coordination layer, not QA: whoever holds coordination owns the scope change and records it, whether that is the human, a human and an agent jointly, or an agent the human has deliberately put in charge of coordination; SPEC §7 and the qa contract say so; a still-bold struck-through ID is a test case (adopts #348) | [#362](https://github.com/radiusred/gh-codecrew/issues/362) / PR [#367](https://github.com/radiusred/gh-codecrew/pull/367) | `satisfied` | Task closed; #348 closed by its finish. Two review rounds. Two Decisions (the operator's nine answers, including `--reinstate` and coordinator-only counting; "never post a QA verdict"), one Deviation (the receipt and the verb list). No requirement of M18 was struck |
| M18-R5 — the doc-synthesizer's front-door obligation is general: the embedded contract obliges refreshing "the hub's front-door documents wherever they make claims the milestone changed", with no reference to an introduction, release, verbs or refusal codes; gh-codecrew's `.codecrew/roles/doc-synthesizer.local.md` carries this hub's specifics (the README's proof points, `docs/introduction.md`'s release, verbs and refusal codes) (adopts #350) | [#363](https://github.com/radiusred/gh-codecrew/issues/363) / PR [#366](https://github.com/radiusred/gh-codecrew/pull/366) | `satisfied` | Task closed; #350 closed by its finish. Approved first round; the approval stood through a message-only force-push by the coordination layer's Decision. Two Decisions |
| M18-R6 — the housekeeping light path and `roles sync`: the implementer and reviewer contracts and SPEC §4/§7 define housekeeping mechanically (a tool defines the target and the diff is the whole decision — no Decision to write), and it travels as a `chore:` PR with review and no task, milestone, plan, verdict or record entry; the text says what the path is in each routing tier, pure solo included, without relying on a ruleset; `roles sync` writes the embedded contracts into `.codecrew/roles/` as one local commit, never pushed, writing a contract only where the local file is absent or equals an embedded release's text and refusing otherwise with the diff named; `status` reports a missing contract as well as a drifted one and names `roles sync` (adopts #343, #266) | [#364](https://github.com/radiusred/gh-codecrew/issues/364) / PR [#368](https://github.com/radiusred/gh-codecrew/pull/368) | `satisfied` | Task closed; #343 and #266 closed by its finish. Three review rounds. Two Decisions (the operator's nine answers, text-only enforcement among them; stamp form and CRLF). The enforcing verb is deferred to capture [#365](https://github.com/radiusred/gh-codecrew/issues/365), open. The path's first use is PR [#375](https://github.com/radiusred/gh-codecrew/pull/375) |
| M18-R7 — the record takes the light path: the doc-synthesizer contract's "Deliver as a task" becomes delivery as a housekeeping PR; `milestone close`'s `DOC_MISSING` gate and `milestone evidence` are unchanged; the requirement/verdict pair is untouched; the `NO_CHECKS` gate is out of scope (adopts #349) | [#369](https://github.com/radiusred/gh-codecrew/issues/369) / PR [#370](https://github.com/radiusred/gh-codecrew/pull/370) | `satisfied` | Task closed; #349 closed by its finish. Two review rounds (a commit subject's case). One Decision (the operator's seven answers, deliberate presence-not-provenance among them). This document is the path's first record |
| M18-R8 — ship v2.1.0: protocol version 2.1, the flip, the tag and five assets; CHANGELOG and `docs/introduction.md` name the struck verdict, the light path and `roles sync` | [#371](https://github.com/radiusred/gh-codecrew/issues/371) / PR [#374](https://github.com/radiusred/gh-codecrew/pull/374), with [#372](https://github.com/radiusred/gh-codecrew/issues/372) / PR [#373](https://github.com/radiusred/gh-codecrew/pull/373) carried under it | `satisfied` | Both tasks closed. Tag `v2.1.0` on [e2b8242](https://github.com/radiusred/gh-codecrew/commit/e2b8242); [release](https://github.com/radiusred/gh-codecrew/releases/tag/v2.1.0) published 2026-09-24T18:43:33Z with five assets; operator leg on #353. #372 widened the requirement by two coordination-layer Decisions on #353. Eleven Decisions and one Deviation across the two tasks; #373 took three review rounds, #374 one |

**Captures adopted and closed by the merges — nine**, each with an "Adopted
by" comment when its task was opened and a "Closed by `gh codecrew task
finish <n>`" comment naming the task, the pull request and the merge commit:
[#337](https://github.com/radiusred/gh-codecrew/issues/337) and
[#338](https://github.com/radiusred/gh-codecrew/issues/338) (by #354),
[#339](https://github.com/radiusred/gh-codecrew/issues/339) and
[#334](https://github.com/radiusred/gh-codecrew/issues/334) (by #355),
[#350](https://github.com/radiusred/gh-codecrew/issues/350) (by #363),
[#348](https://github.com/radiusred/gh-codecrew/issues/348) (by #362),
[#343](https://github.com/radiusred/gh-codecrew/issues/343) and
[#266](https://github.com/radiusred/gh-codecrew/issues/266) (by #364), and
[#349](https://github.com/radiusred/gh-codecrew/issues/349) (by #369). #334 was
split between two tasks: #355 adopts it and #354 delivers its SPEC §10 half,
which both tasks' goals say. The oldest, #266, was filed 2026-09-06; three of
the nine (#348, #349, #350) came from projects outside this hub, as the goal
describes.

**Captures filed during the milestone and left open — one.**
[#365](https://github.com/radiusred/gh-codecrew/issues/365), filed by the
implementer seat at 15:44:03Z on 2026-09-24 under #364's answer 1(a): "the
housekeeping light path has no verb enforcing its review gate — a chore
finish". Its one comment, from #369's implementer at 17:06:39Z, names "A
second consumer: the milestone record." No issue was filed in
[radiusred/codecrew-www](https://github.com/radiusred/codecrew-www) in the
window.

## The front door

The doc-synthesizer contract, as M18-R5 left it, sends this hub's front-door
list to its `.local.md`: "the README's proof points, and the introduction's
release, verbs and refusal codes". #374 had already moved both documents to
v2.1.0. Read against what the milestone delivered, at the counting instant:

- **`README.md`** — true as it stands. "The current release, v2.1.0,
  implements protocol 2.1 … so a 2.0 repo has no migration step"; the gates
  paragraph asks "a QA verdict on every requirement not struck by a recorded
  decision", which #367 wrote; and the routing snippet, which presents itself
  as this repository's `.codecrew/config.yml` "as it stands today", is
  identical to that file's `roles:` section after both of M18's routing
  changes.
- **`docs/introduction.md`** — the release (`v2.1.0 … implementing protocol
  2.1`) and the verb list (`milestone new/evidence/close/strike`, `roles
  diff/show/sync`) match `gh codecrew help` for the installed v2.1.0, and the
  six `--dry-run` verbs are pinned by R1's drift test. The refusal-code
  section points at SPEC §10 and states no count. One claim was false: the
  sentence #374 added says `milestone close` "counts a struck requirement as
  terminal, as it does a QA `pass`". There is no `pass` verdict — the grammar
  is `satisfied|not satisfied|untestable`, as QA's M18-R7 verdict quotes it,
  and SPEC §4's gate asks for `satisfied`. This PR changes that one word to
  `satisfied`, and nothing else in either document.

## Protocol-discipline observations

Eight things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **The ask-the-human point became the main channel for human judgment, and
  it ran through the implementer.** Five tasks put thirty-seven numbered
  points to the operator before any code, with options and a recommendation
  each; no gate was raised. Every answer reached the trail as a Decision
  posted by `radiusred-cody[bot]`, "relayed by the coordination layer", and
  thirty-six of the thirty-seven take the plan's recommendation. The
  operator's own account posted three Decisions in the milestone, all
  coordination-layer calls rather than answers to a plan.
- **A coordination-layer Decision was corrected by a comment the grammar does
  not gather.** The first Decision on #353 rested on a premise #372's
  implementer disproved on a built binary, and the correction followed
  within five minutes. It is labelled `**Correction …**`, so the record
  grammar reads it as prose; the re-put answer three minutes later is a
  Decision and is gathered. That Decision opens "re-put on the corrected
  premise", so the gathered raw material points at the correction without
  carrying it.
- **Routing changes rode task pull requests as Deviations until the light
  path existed, and the pointer move still did after.** #354 and #358 each
  carried an operator Deviation and a `chore(config):` commit, because "The
  #343 light path is still only a proposal". After R6, #374 moved the hub
  pointer inside the release task by the #114 precedent; its pull request
  body notes that a separate pull request "would need its own task, because a
  hand edit is not a tool's output". #375 is the path used as SPEC §4 now
  defines it: a tool, its output as the whole diff, the holder's approval,
  the author's merge.
- **The requirement text was the brief the reviewer held the change to.**
  #373's two change requests are both about whether the version check can
  run "before the first verb", the phrase #372's goal uses; half of #367's
  is about whether the tests cover the ordering R4's grammar introduced. The
  reviews report building the branch and exercising the verbs, not reading
  the diff alone.
- **Two commit-lint failures were repaired by force-push, and the protocol
  handled them two ways.** On #366 the failing head was already approved, and
  a coordination-layer Decision let the approval stand on the identical tree.
  On #370 the reviewer's round one asked for the reword and round two
  verified the tree unchanged. A further four force-pushes — #357 at
  12:03:15Z, #367 at 15:51:33Z and 15:53:23Z, #368 at 16:05:23Z — have no
  comment; three follow a sibling pull request's merge by forty-four to
  sixty seconds, consistent with rebases onto `main`, and the fourth, on
  #367, came twenty-one seconds after the pull request opened. *Undocumented,
  inferred from the timeline.*
- **A reviewer session posted under the wrong identity, and the only record
  is a dismissal message that cites a review that does not exist.** The
  duplicate on #357 carried the operator's login and was dismissed by the
  reviewer App a second after its own round two. It was a change request,
  and the review gate counts only the holder's approval; the trail does not
  say how the session came to post under that login.
- **The light path was defined, used for a tool's output, and taken by a
  record, all within the milestone.** R6 landed at 16:21:27Z, R7 at
  17:19:06Z, #375 used the path at 18:51:12Z, and this document is the
  record the Gates call "the first record to take R7's path". The one gap
  this record met while following the contract as written is recorded as
  its Deviation above: the dispatch brief supplied facts the trail did
  not.
- **Scope widened inside a release requirement, by Decision, and the release
  still shipped the same day.** #372 added a verb behaviour and a refusal
  code under M18-R8 after the flip's plan was written; #371's `task start`
  waited for it, and the CHANGELOG lead names it. The widening is on the
  trail as two Decisions and a correction on the milestone issue, each
  qualified "coordination layer".

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **The conversation that scoped the milestone is on the record only through
  #353's goal**, "Decided with the operator 2026-09-23", and the four
  appraisals it names, which predate the issue. *The reasoning is recorded;
  the deliberation is not*, as
  [M17's record](17-the-frontier-builds-the-coordinator.md#what-the-record-does-not-contain)
  and [M16's](16-bedding-in-the-venue-seam-and-the-reference.md#what-the-record-does-not-contain)
  each say of their own openings.
- **The operator's answers are on the trail only as relays.** Each of the
  five answer Decisions quotes what the coordination layer relayed; the
  exchange in which the operator gave them is not on the trail. #369's
  Decision refers to a "phase-1 'wrong or underspecified' list" the operator
  asked to be folded in; the list itself was not posted on any issue.
- **The blind replay behind the reviewer's routing change is on the trail
  only as #358's Deviation.** Its timings and findings are quoted there; the
  replayed reviews themselves were not posted.
- **Nothing on the trail says why the milestone paused twice** — ten hours
  between #357's merge and #358's creation, and fifteen hours and
  forty-four minutes between the R4–R6 plans and their answers.
- **Which harness and model ran each review and QA round is not recorded by
  the seats.** This record's Deviation carries the coordination layer's
  account and the trail's partial support for it; no review names its own
  model, so which rounds after #358 ran under `gpt-5.6-terra` rests on that
  account.
- **The stray review under the operator's login has no comment**, and its
  dismissal message points at review 5291006191, which the reviews API does
  not return.
- **Four force-pushes have no comment.** The observation above infers
  rebases from the timing.
- **`milestone strike` was not used live.** Its behaviour on a real milestone
  is on the trail through QA's refusing dry runs and the tests; no M18
  requirement was struck.
