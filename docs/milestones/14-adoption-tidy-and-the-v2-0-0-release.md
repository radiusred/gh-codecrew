# M14: Adoption, tidy, and the v2.0.0 release

Tracking issue: [#269](https://github.com/radiusred/gh-codecrew/issues/269) ·
Synthesized 2026-09-06 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #269's **three** requirements as opened and **two
added, none struck and none amended** — the issue's `userContentEdits` query
returns a count of three, an original and two edits whose whole diffs are the
addition of M14-R4 at 11:38:01Z and M14-R5 at 14:49:05Z — leaving **five**;
its Gates section; its **eight comments**, carrying **three Decision records,
one Deviation record, one bare operator instruction, the operator's fleet
migration record and two QA comments**; the **six** delivery task issues, all
in this repository
([#270](https://github.com/radiusred/gh-codecrew/issues/270),
[#271](https://github.com/radiusred/gh-codecrew/issues/271),
[#272](https://github.com/radiusred/gh-codecrew/issues/272),
[#273](https://github.com/radiusred/gh-codecrew/issues/273),
[#283](https://github.com/radiusred/gh-codecrew/issues/283),
[#299](https://github.com/radiusred/gh-codecrew/issues/299)) plus
[#302](https://github.com/radiusred/gh-codecrew/issues/302), the task this
document is; their **six merged PRs** and the **twenty-five commits** on
them; the **twenty-six Decision records and four Deviation records** across
the milestone issue and those task issues; the **thirteen review
submissions** on the six PRs — **five change requests**, **eight approvals
submitted of which six stand and two were dismissed by a rebase**; the
**four** adopted backlog captures the merges closed; and the **five**
captures the milestone filed and left open. The house form is
[M13's document](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md),
with [M12's](12-v1-2-0-and-the-field-fixes-behind-it.md) and
[M11's](11-housekeeping.md) behind it; the record standard is the one the
reviewer set on
[PR #290](https://github.com/radiusred/gh-codecrew/pull/290), whose two
rounds found seven claims in the M13 record that did not survive a check,
three of them counts — so every count above and below is re-derived
from the source rather than carried across, records by the rule
`tracker.ExtractRecords` itself applies (a paragraph-initial `**Decision:**`,
`**Deviation:**` or `**Gate resolved:**` label, bare or parenthetically
qualified), review states from each PR's reviews API and timeline, and the
fleet from the operator's migration record rather than from the checklist
that preceded it.

**This record's citations can be checked by a released binary; M13's could
not.** The v1.2.0 extension could not read this hub at all after M13-R1
moved the pointer, which is why
[M13's record](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md)
closes on the observation that "until v2.0.0 ships in M14 there is no
released binary that can check these citations at all". It shipped here. The
installed extension now reports `v2.0.0 (protocol 2.0)` — not
`dev-<sha> (protocol 2.0)` — and against this hub it prints:

```
$ gh codecrew milestone evidence 13
note: milestone M13 (radiusred/gh-codecrew#254) is closed — checking the record as it stands
requirements counted: M13-R1, M13-R2, M13-R3, M13-R4, M13-R5, M13-R6, M13-R7, M13-R9 (8)
all 9 cited links resolve across 13 issues — evidence is reachable
```

The gather from `gh codecrew milestone close 14` is again not part of the raw
material, for the reason
[M9's](9-the-docs-at-codecrew-works.md),
[M10's](10-protocol-bookkeeping-from-the-field.md),
[M11's](11-housekeeping.md),
[M12's](12-v1-2-0-and-the-field-fixes-behind-it.md) and
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md)
records all give: `--dry-run` stops at the gate that counts tasks, naming the
task that writes this file. What is new is which binary says so — the
released one, run from the repository root:

```
gate milestone open: ok
gate no gate raised: ok
gate tasks closed: refused[OPEN_TASKS]: tasks not closed: radiusred/gh-codecrew#302 (ready)
gate requirements declared: not reached
gate QA verdicts: not reached
gate milestone document: not reached
dry run: nothing written — the live verb stops at the first refusal above
```

Every record below was read from its issue, PR, review, commit, tag, release,
workflow run or timeline directly. The prose is wrapped at this file's normal
width; the requirement-outcomes table's rows and a handful of single links
are not, because a newline ends a Markdown table row and breaking a link
breaks it — the same shape [M13's record](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#requirement-outcomes)
and [M12's](12-v1-2-0-and-the-field-fixes-behind-it.md#requirement-outcomes)
carry, and the shape `ROADMAP.md` has had since M1.

**This PR adds the M14 ROADMAP row; it does not flip one** — the convention
M10-R1 introduced and M12-R1 shipped, which
[M13's record PR](https://github.com/radiusred/gh-codecrew/pull/290),
[M12's](https://github.com/radiusred/gh-codecrew/pull/252),
[M11's](https://github.com/radiusred/gh-codecrew/pull/240) and
[M10's](https://github.com/radiusred/gh-codecrew/pull/232) were written
under. `ROADMAP.md` carried no M14 row before this PR, so there was nothing
to discard and the trail carries no Decision discarding one — the same as
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md),
and unlike
[M10](10-protocol-bookkeeping-from-the-field.md#the-roadmap-row-belongs-to-the-doc-synthesizer-at-both-ends),
[M11](11-housekeeping.md#the-roadmap-row-discarded-again) and
[M12](12-v1-2-0-and-the-field-fixes-behind-it.md#the-roadmap-row-discarded-for-what-the-operator-calls-the-last-time),
which each had one.

## Goal and outcome

#269 was opened at 11:31:59Z on 2026-09-06, four seconds after the Decision
on #254 struck M13-R8 and moved the release here. Its goal names three
things and joins them with one sentence of reasoning: "M13 ships protocol
2.0's breaking changes; this milestone finishes the bookkeeping the protocol
already knows enough to do itself, then cuts the release that carries all of
it."

The bookkeeping is two captures from the field, and the goal states what they
have in common rather than describing them separately: "a backlog issue a
milestone adopts stays open after its task ships unless a brief remembers a
`Closes` line ([#193](https://github.com/radiusred/gh-codecrew/issues/193),
from davison/topos M3 and every hub milestone since), and local task branches
whose upstream the protocol deleted accumulate until someone sweeps by hand
([#192](https://github.com/radiusred/gh-codecrew/issues/192)). Both are 'the
protocol knows the link, so the protocol should do the bookkeeping'." That
sentence is the milestone's whole thesis, and the two requirements added
after it were added because they are the same shape: a link or a name the
protocol already holds and was not acting on.

Three requirements were present at 11:31:59Z: R1 (adoption as a first-class
link), R2 (`task finish`'s local cleanup) and R3 (the release and the fleet
migration). **M14-R4** was added at 11:38:01Z, with the
[operator Decision](https://github.com/radiusred/gh-codecrew/issues/269#issuecomment-5558957391)
recording it seven seconds later, adopting
[#167](https://github.com/radiusred/gh-codecrew/issues/167): "the
stale-branch sweep is the same file family as #192's local cleanup and
belongs in the bookkeeping milestone." The same Decision names six captures
read and deliberately not bundled — #108, #191, #134, #194, #36 and #31 —
which is the only place on the trail where those six are weighed at all.
**M14-R5** was added at 14:49:05Z, from a bare operator
[comment](https://github.com/radiusred/gh-codecrew/issues/269#issuecomment-5559826513)
at 14:18:50Z ("Adopt #267 and use colours from the crew image palette") that
the coordination layer wrote up as a
[Decision](https://github.com/radiusred/gh-codecrew/issues/269#issuecomment-5560010556)
thirty minutes later, adopting
[#267](https://github.com/radiusred/gh-codecrew/issues/267). Each addition is
purely additive: the two edits touch nothing that was already there, and no
requirement was struck or amended anywhere in the milestone.

The six delivery tasks, in merge order:

- **[#271](https://github.com/radiusred/gh-codecrew/issues/271) / PR
  [#292](https://github.com/radiusred/gh-codecrew/pull/292)** — M14-R2.
  Three commits, two review rounds, merged 17:20:48Z
  ([8a1e5ce](https://github.com/radiusred/gh-codecrew/commit/8a1e5ce)).
  Three Decisions. Adopts and closes
  [#192](https://github.com/radiusred/gh-codecrew/issues/192).
- **[#270](https://github.com/radiusred/gh-codecrew/issues/270) / PR
  [#294](https://github.com/radiusred/gh-codecrew/pull/294)** — M14-R1. Six
  commits, two review rounds, merged 17:33:29Z
  ([490a410](https://github.com/radiusred/gh-codecrew/commit/490a410)). Four
  Decisions and one Deviation. Adopts and closes
  [#193](https://github.com/radiusred/gh-codecrew/issues/193) — the last
  capture closed by the PR-body convention this task replaced.
- **[#273](https://github.com/radiusred/gh-codecrew/issues/273) / PR
  [#293](https://github.com/radiusred/gh-codecrew/pull/293)** — M14-R4. Six
  commits, three review rounds, the third a rebase check, merged 17:43:56Z
  ([c452fcf](https://github.com/radiusred/gh-codecrew/commit/c452fcf)). Five
  Decisions. Adopts and closes
  [#167](https://github.com/radiusred/gh-codecrew/issues/167); its second
  half is split out as [#295](https://github.com/radiusred/gh-codecrew/issues/295).
- **[#283](https://github.com/radiusred/gh-codecrew/issues/283) / PR
  [#291](https://github.com/radiusred/gh-codecrew/pull/291)** — M14-R5. Five
  commits, three review rounds, two of them after a rebase, merged 17:56:30Z
  ([dd5913d](https://github.com/radiusred/gh-codecrew/commit/dd5913d)). Five
  Decisions and one Deviation, which a later Decision withdraws. Adopts and
  closes [#267](https://github.com/radiusred/gh-codecrew/issues/267).
- **[#272](https://github.com/radiusred/gh-codecrew/issues/272) / PR
  [#298](https://github.com/radiusred/gh-codecrew/pull/298)** — M14-R3,
  Part A. Three commits, two review rounds, merged 18:46:22Z
  ([2b75fdc](https://github.com/radiusred/gh-codecrew/commit/2b75fdc)). Four
  Decisions and one Deviation. Part B — the tag, the release run and the
  assets — and the operator's verification and fleet migration all ran on
  the closed issue.
- **[#299](https://github.com/radiusred/gh-codecrew/issues/299) / PR
  [#300](https://github.com/radiusred/gh-codecrew/pull/300)** — M14-R3's
  spirit rather than its text, and the task says so. Two commits, approved
  first round, merged 18:57:41Z
  ([a407c8c](https://github.com/radiusred/gh-codecrew/commit/a407c8c)). Two
  Decisions, the second reversing the first.

From the milestone opening to the last delivery merging: seven hours,
twenty-five minutes and forty-two seconds. From the first `task start` to
that merge: **one hour, fifty-seven minutes and five seconds**. The gap is
M13: #254 closed at 16:56:50Z and #271 started at 17:00:36Z, three minutes
and forty-six seconds later, because the two milestones share this hub and
M14's own tasks were created at its opening and left to wait.

Four of the six PRs were open at once. #291 opened at 17:07:13Z, #292 at
17:08:31Z, #293 at 17:08:44Z and #294 at 17:09:18Z — two minutes and five
seconds across all four — and they merged in the order #292, #294, #293,
#291, each later branch rebased onto the merge before it. That parallelism is
where both dismissed approvals and both rebase-check rounds come from, and
[#293's PR body](https://github.com/radiusred/gh-codecrew/pull/293) records
the sequencing as a review finding taken ("**Merge order** (finding 6): #292
(#271) has since merged").

Five requirements, five `satisfied` verdicts standing — one of them reached
on the second QA round, after a Decision rather than a remedy task.

What exists afterwards that did not before. `task new --adopts <ref>[,<ref>]`
records adopted captures under an `## Adopts` section and comments on each
one; `task finish` closes them after the merge with a comment naming the
task, the PR and the merge commit, and refuses `ADOPT_NOT_OPEN` at creation
time for a ref it cannot read or that is already closed — SPEC §10's table
counts **forty-three** codes now, from forty-two. `task finish` also tidies
the clone it ran in: fetch, switch off the task branch, fast-forward the
default branch, delete the local task branch, and name-and-keep anything
carrying commits the merge did not. `milestone close` runs a second branch
pass over every `task/<n>-<slug>` branch in the hub and in each task's repo,
not only its own milestone's. `init`, `checkpoint` and `migrate` bring
`cc:milestone`, `cc:task` and `cc:needs-decision` into existence with fixed
descriptions and colours sampled from the project's own crew images. This hub
carries the `CLAUDE.md` its own `init` writes, with a test that neither root
entry point can drift from its scaffold. And v2.0.0 is released and installed
over the local build, with ten repositories moved to the 2.0 layout, one
refused by design and one retired.

## The release, and the fleet

M14-R3 is the requirement M13 struck and moved here, unchanged in wording,
and it is the only one whose delivery is split between a seat and the
operator. The trail follows that split exactly.

**Part A**, PR [#298](https://github.com/radiusred/gh-codecrew/pull/298),
merged 18:46:22Z at
[2b75fdc](https://github.com/radiusred/gh-codecrew/commit/2b75fdc):
`## [Unreleased]` became `## [2.0.0] — 2026-09-06` with a fresh empty
`## [Unreleased]` above it, nineteen entries moving down unchanged but for
one tidy; a new `### What broke, and what to do` block leading the section;
`docs/introduction.md` and `README.md` flipped to v2.0.0 and protocol 2.0.

**Part B**, run by the implementer seat on the closed issue and evidenced in
one [comment](https://github.com/radiusred/gh-codecrew/issues/272#issuecomment-5561368504)
at 18:48:23Z. Every claim in it was re-checked for this record against the
API: the [annotated tag `v2.0.0`](https://github.com/radiusred/gh-codecrew/releases/tag/v2.0.0)
is tag object `93bf2f37818e1ccc7f2f3cff497628aae846d805`, pointing at commit
`2b75fdc6fdf5f6c8bb44a7ce8068fd0e3ebfad09`, message
`v2.0.0 (protocol 2.0)`, tagger
`radiusred-cody[bot] <270330637+radiusred-cody[bot]@users.noreply.github.com>`
at 18:46:35Z; `release.yml`
[run 34052790620](https://github.com/radiusred/gh-codecrew/actions/runs/34052790620)
completed with conclusion `success`; the
[release](https://github.com/radiusred/gh-codecrew/releases/tag/v2.0.0) is
published, neither draft nor prerelease, named `gh-codecrew 2.0.0`, and
carries **five** assets — darwin-amd64, darwin-arm64, linux-amd64,
linux-arm64 and windows-amd64.exe.

**The operator's legs**, posted on the closed #272 at 18:49:08Z as a
[four-step checklist](https://github.com/radiusred/gh-codecrew/issues/272#issuecomment-5561372695):
install the release over the local build, verify on this hub, migrate the
fleet hubs before their spokes, and record the run on #269. The checklist
names twelve repositories — six hubs (radiusred/ops, radiusred/numberguess,
radiusred/snake, davison/davisononline-www, davison/md-notes, davison/topos),
five spokes (radiusred/.github, radiusred/infrastructure, radiusred/www,
radiusred/codecrew-www, davison/topos-plugins) and davison/numberguess,
"excluded by decision (protocol 0.1 sandbox, never touched again)".

**The migration record**, one
[comment](https://github.com/radiusred/gh-codecrew/issues/269#issuecomment-5562585598)
at 22:20:39Z, is a fenced block of `git log` output per repository. It names
**eleven**: ten with two commits each — `chore: migrate codecrew to the 2.0
layout` followed by `chore: update agent entry points` (or, on
radiusred/numberguess, `chore: update agent instructions for codecrew 2.0`) —
and `davison/numberguess`, which carries the refusal verbatim instead of
commits:

```
codecrew: refused[MIGRATION_UNSUPPORTED]: .codecrew.yml says protocol "0.1"; migrate moves a protocol 1.x repo to 2.0 and nothing else — a pointer below 1.0 predates the conventions it assumes (SPEC §5)
```

That is `migrate` refusing by design, which is the outcome the checklist
predicted for it. It is also the first exercise of any of M13-R2's
conservative refusals against a repository nobody built for the test —
[M13's record](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#what-the-record-does-not-contain)
closes on the observation that "`migrate` has never run against a repository
that was not made for it".

**The one repository missing from the record was retired rather than
migrated, and the QA seat found the gap.** The record lists eleven, the
checklist named twelve, and `davison/md-notes` appears in neither the record
nor as a repository the QA App's token can reach; M14-R3 came back
`not satisfied` on that alone. The operator's
[Decision](https://github.com/radiusred/gh-codecrew/issues/269#issuecomment-5562689587)
nine minutes later settles it: "davison/md-notes is retired from the CodeCrew
fleet, not migrated ... it was an empty project, and the working directory
was removed and re-created rather than moved to the 2.0 layout." It states
the fleet as it now stands — "the eleven repos in the migration record above
... ten migrated ... plus davison/numberguess, refused
`MIGRATION_UNSUPPORTED` by design as a protocol-0.1 sandbox" — and says what
happens if the repository ever comes back: "it starts fresh with
`gh codecrew init` on a clean tree; its remote may still carry the 1.x files
and a 2.0 binary run there refuses `LAYOUT_LEGACY`, which is the intended
outcome for a retired repo." **Rejected:** "migrating an empty project for
the sake of the count."

So the fleet of record is eleven repositories, ten of them migrated. Twelve
is what M14-R3's text and the checklist say, and it was right when they were
written; this record states the fleet as decided rather than as first listed,
which is the whole reason the Decision exists.

## Decisions

Twenty-six Decision records across the milestone issue and the six delivery
task issues. Four are the operator's — three on #269 and one on #272; the
other twenty-two were written by the seat doing the work,
`radiusred-cody[bot]`, which started all six delivery tasks. Two of the
twenty-six carry a parenthetical qualifier rather than the bare label, and
both matter to the count: `**Decision (superseding the scope Decision
above):**` on #283 and `**Decision (follow-up):**` on #299 are records under
`ExtractRecords`' rule and are invisible to a search for `**Decision:**`
alone. They group by the question each task had to settle.

### What an adoption is, in a body a human also reads

M14-R1 had to choose what counts as an adopted ref in a task body, and the
[Decision](https://github.com/radiusred/gh-codecrew/issues/270#issuecomment-5560800434)
picks the narrowest possible answer: "only the ref at the head of a list line
is an adoption; the prose after it is the capture's title, which `task new`
writes for the reader and nothing parses back." **Why:** the title is free —
`task new --adopts` already reads each capture to check it is open — and it
makes the section readable, "but a title is arbitrary text: a capture called
'fix #42 in the parser' would give a body-wide ref scan a second adoption,
and `task finish` closes what it finds." **Rejected:** scanning the whole
section for refs, and writing bare refs with no titles ("safe, but the
section stops being readable").

A [second Decision](https://github.com/radiusred/gh-codecrew/issues/270#issuecomment-5560800510)
splits the work in two so `--dry-run` can be honest: `planAdoptions` decides
each capture's fate *before* the merge, "the way `planSweep`/`executeSweep`
already split the milestone close's branch sweep". A capture whose state
cannot be read is planned as a close. **Why:** reading each capture's state
after the merge "would make the dry run say 'would close' about an issue it
knows is closed", and failing the verb there "would refuse a finish whose
gates have all passed, over bookkeeping". **Trade-off:** a capture that the
PR's own body closes at merge time is planned as an open one; the attempted
close then no-ops and prints a note, "which is the right shape for a step
that cannot refuse".

A [third](https://github.com/radiusred/gh-codecrew/issues/270#issuecomment-5560800367)
adds `MergeCommit(repo, number)` to the `Tracker` interface so the closing
comment names a commit that exists: "`PR.HeadSHA` is the head as GitHub last
saw it, and a rebase merge rewrites that commit, so quoting it would point a
reader at a commit that is on no branch." **Trade-off:** one more method and
one more API call, only on a finish that adopted something. **Rejected:**
widening `PRInfo`'s query, "a heavier call ... and it widens the struct three
other gates read".

### The record-reading rule gets its third reader, after a review found two holes

The
[fourth Decision on #270](https://github.com/radiusred/gh-codecrew/issues/270#issuecomment-5560881164)
and the Deviation in the same comment are the milestone's clearest instance
of a review finding a hole a Decision had already described. The reviewer
reproduced two ways past the head-of-line rule: a `- #42` inside a fenced or
four-column code block *inside* the `## Adopts` section was adopted, and any
earlier prose carrying the literal string `## Adopts` — #270's own Goal does
— shadowed the real section and returned an issue nobody named.

The Decision answers the reviewer's proposed general fix rather than the
specific one: `section()` is left alone for `PlanPresent` and
`RequirementIDs`, and the `## Adopts` section gets its own heading-anchored
reader. "checky's finding 2 is right that all three callers share the
substring weakness, but not that they share the consequence. A shadowed
`## Plan` costs a presence check that `task start` re-asks on the next run; a
shadowed `## Requirements` costs an ID list a human reads in the refusal. A
shadowed `## Adopts` closes an issue nobody adopted, after a merge, in a step
that by design cannot refuse — the one output in the package that acts on
other people's issues." **Rejected:** changing `section()` for all three
callers now, "the right change if the other two ever gain teeth, and a task
of its own with its own test sweep, not a rider here."

The Deviation beside it is unusually candid about its own scope: "`AdoptedRefs`
reads the body through `tracker.StripCode` and finds its heading as a
heading; the plan described neither, and my own Decision on the head-of-line
rule described only half the hole it was closing." M13-R6's one-rule
requirement — one code-stripping implementation shared by the verdict scan
and the citation walk — "now has its third reader".

### Two grounds for a force-delete, and a dry run that does not fetch

M14-R2 automates `git branch -D`, which is the thing #192 opens by observing
that operators hesitate to script. The
[Decision](https://github.com/radiusred/gh-codecrew/issues/271#issuecomment-5560797163)
allows it on two grounds only: the branch tip is the exact commit GitHub
reported as merged (`PR.HeadSHA`, frozen at merge), or the tip is contained
in the fetched default branch. **Trade-off:** a clone merely *behind* the
pushed head is kept, "so the tidy is occasionally a no-op where a delete
would have been safe. That is the direction to be wrong in." **Rejected:**
the ordinary ancestry test, "useless here: the hub rebase-merges, so a
cleanly merged task branch's commits are never ancestors of `main` (this is
why `git branch -d` refuses them, the observation #192 opens with)."

The [second](https://github.com/radiusred/gh-codecrew/issues/271#issuecomment-5560797256)
holds `--dry-run` to writing nothing at all, the clone included: "It reads
the clone exactly as it stands and names the steps the live verb will take,
including the fetch itself." **Trade-off:** the fast-forward line becomes a
prediction rather than a measurement — with the argument for why the shape of
the answer cannot change across the merge: "`origin/<default>` only ever
moves forward, so a local default branch that is an ancestor now is still one
afterwards, and a diverged one is still diverged. The deletion decision is
measured, not predicted." **Rejected:** fetching during the dry run, since "a
fetch moves remote-tracking refs".

The [third](https://github.com/radiusred/gh-codecrew/issues/271#issuecomment-5560797350)
takes the one place the requirement's wording and the useful behaviour part
company. M14-R2 words the exception as being about the deletion alone; the
Decision switches the operator off the task branch even when the branch is
then kept. **Trade-off:** "it moves the operator's checkout out from under
work they have not pushed ... the commits are untouched and the line naming
the kept branch is the last thing printed, so the way back is on screen."
**Rejected:** staying put whenever the branch is kept, which "would leave the
common case ... with a merged, deleted-upstream branch checked out, which is
the state #192 is about."

### A wider candidate set, and the two guards a review added to it

M14-R4's whole change is which branches a close looks at, and three Decisions
fix the boundaries. The
[first](https://github.com/radiusred/gh-codecrew/issues/273#issuecomment-5560794368)
lists "the hub the milestone issue lives in, plus every repo the milestone's
tasks name, deduped and hub first", because "the hub is where every milestone
issue lives and where the record is read, so it is the one repo a close can
always be certain about, and it is the repo a solo project's stale branches
accumulate in — #167's own evidence is `radiusred/numberguess`, a solo hub."
**Rejected:** every repo in the org or in the spoke list, which would let "a
close ... reach repos this milestone never touched".

The [second](https://github.com/radiusred/gh-codecrew/issues/273#issuecomment-5560794454)
keeps the delete conditions identical to the first pass's, at a real API
cost, and gives the reason the cheap test cannot work: "this project
rebase-merges, so a merged branch's commits are rewritten and it is never
`ahead == 0` — only 'its PR merged and the tip still at the merged commit'
catches it, and that needs the PR." Reusing `branchAction` also "means the
two passes can never drift apart". **Rejected:** reading PRs out of the
branch listing itself, which "would give the `Tracker` seam a second, richer
PR type for one caller's benefit".

The [third](https://github.com/radiusred/gh-codecrew/issues/273#issuecomment-5560794527)
makes a branch a candidate on its name alone, without requiring its issue to
carry `cc:task`, and states the safety property in one clause: "An issue
number that resolves to nothing at all is a `note:` and a skip, never a
delete." That clause is what the reviewer's round-one finding then went after,
because nothing tested it.

Two more Decisions answer that review.
[The open-PR guard](https://github.com/radiusred/gh-codecrew/issues/273#issuecomment-5560890882)
is taken "as required by the coordination layer" and its reasoning is the
sharpest scoping argument in the milestone: `branchAction`'s "no PR" arm "is
only as good as the PR set it is handed, and `ClosingPRs` finds the PRs that
close `<n>` — so a PR opened on a task branch without a `Closes #<n>` line is
invisible to it, and deleting the branch would close that PR. Pass one's
candidates come from the task's own PRs and one conventional name, so its
exposure is narrow; pass two makes every `task/<n>-…` branch ... a candidate,
which is wider than the arm was written for." **Rejected:** widening the PR
set for every candidate, and "leaving it to the Decision text, which was the
reviewer's non-blocking suggestion — a guard that only exists in prose is the
shape finding 1 was about."
[The truncation flag](https://github.com/radiusred/gh-codecrew/issues/273#issuecomment-5560890979)
reports rather than pages: "one page is a hundred *task* branches, which no
repo in this project approaches ... paging would buy correctness nobody can
currently reach, at the cost of an unbounded walk in a verb whose whole
design claim is a bounded one."

### The palette, measured rather than remembered

M14-R5 asked for colours "drawn from the crew image palette" and named one
hex value in its own text — the mark's cyan `#1ad1ff` — with the three seat
accents named by colour only. The
[Decision](https://github.com/radiusred/gh-codecrew/issues/283#issuecomment-5560760544)
replaces four remembered values with measured ones, and says so: "The operator's rough
sample was close and is superseded by the measurement: the mark's cyan is
#01d4ff rather than #1ad1ff, the review pink #f0aeff rather than #e8a8f8, the
test teal #92edff rather than #88d8e8, the docs magenta #e162ff rather than
#c858e8." The method is stated well enough to repeat — every accent pixel in
`codecrew-www`'s `docs/assets/images/crew/*.png` at 512×512 with alpha ≥ 200,
saturation > 0.35 and lightness between 0.20 and 0.90, counted by frequency,
each image carrying its accent at a pure tone and a shaded one, and the pure
tone taken.

The mapping is argued from where the colours are actually read: "`cc:needs-decision`
takes the review pink because a gate is a question for a human, and review is
the seat whose whole job is asking one. `cc:milestone` takes the mark's cyan
... `cc:task` takes the test teal — the same hue as the milestone's cyan,
lightened — because a task is a part of a milestone and the family should
read as one, and because it is the pairing that keeps the two labels which
actually co-occur on one issue far apart." **Trade-off:** `cc:milestone` and
`cc:task` are told apart by lightness alone, "deliberate and costs nothing in
practice — the two never appear on the same issue". **Rejected:** `cc:task`
in the docs magenta, which "puts `cc:task` (h289 l0.69) beside
`cc:needs-decision` (h289 l0.84) on every gated task issue, which is the one
place the label colours have to do work. The magenta is left unassigned
rather than spent badly."

Those three colours and their descriptions are live on this repository's
label list today, which is the cheapest possible check of the whole
requirement: `cc:milestone` `#01d4ff`, `cc:task` `#92edff`,
`cc:needs-decision` `#f0aeff`.

### The one verb the requirement did not name, and the Decision that reversed itself

The most instructive Decision pair in the milestone is on #283, and it is a
scope question answered twice.

The [first](https://github.com/radiusred/gh-codecrew/issues/283#issuecomment-5560793417)
declines to touch `migrate`: "M14-R5 does not name `migrate`, and the rule R5
does state — an existing label is left untouched — is the opposite
instruction ... Doing both in one PR would ship a verb that contradicts the
SPEC sentence in the same diff." It names the cost plainly — "the twelve
repos M14-R3 migrates all have their `cc:` labels from implicit creation, so
they keep their random colours and no description, and `init` will not fix
them" — and says whose call that is: "it is the operator's call to pay or
not, which is why it is on the record as an ask-the-human point in the Plan
rather than resolved here." A **Deviation** in the same comment records the
consequence: the PR's `Closes #267` "closes the capture with that comment
unaddressed", flagged so the operator can open a follow-up.

Nine minutes later the
[superseding Decision](https://github.com/radiusred/gh-codecrew/issues/283#issuecomment-5560845742)
takes the other branch on the operator's answer, delivered through the
coordination layer from the operator's own second comment on #267. `migrate`
restyles all three labels; `init` and `checkpoint` still do not. The
asymmetry is defended rather than smoothed over: "`migrate` is the one-shot
move to the 2.0 layout and its whole promise is a repository
indistinguishable from a fresh 2.0 `init` — and a 1.x repository's `cc:`
labels were, without exception, created implicitly by the first verb that
mentioned one: a colour GitHub generated, no description. There is no project
intent there to preserve. `init` is the verb that reruns, on a repository
whose labels somebody may have chosen deliberately, so it keeps its
never-restyle rule." **Trade-off:** two rules where there was one, "bought
with the twelve repositories M14-R3 migrates: without it every one of them
lands on 2.0 still wearing GitHub's grey, and no verb would ever fix them."
**Rejected:** a `migrate --restyle-labels` opt-in flag, which "would be off
by default and so unused in exactly the twelve runs it was added for"; and
renaming a recased label while restyling it, since "GitHub matches label
names case-insensitively ... and a rename is a visible change to issue
history that the instruction did not ask for."

The same comment withdraws the Deviation explicitly — "no longer holds and is
withdrawn: the capture is answered in full" — so the trail carries the
Deviation and its retraction side by side rather than an edited comment.

A [third Decision](https://github.com/radiusred/gh-codecrew/issues/283#issuecomment-5560935066)
settles where the label step runs after the reviewer offered two options and
the coordinator left the choice with the seat: "on every path where the move
reached disk — after the commit, after a detached-HEAD note, and after a
commit `git` refused — rather than being announced as skipped. ... The rule
it buys is one sentence: once the migration is on disk, the labels are done."
**Trade-off:** "the labels are written to GitHub for a migration the operator
may still abandon ... an abandoned migration leaves a repository with three
well-styled labels and nothing else — which is exactly what `init` would have
left." **Rejected:** running the step before the commit attempt, which "would
put a remote write ahead of the local work on the one path that succeeds".
A fourth, `**Decision (checky's nit 3):**` in the same comment, replaces an
unstated premise with a structural one: the restyle now passes "the
repository's own spelling of the label name to `UpdateLabel`", because
"checky was right that `labels.go` passed the canonical name and that the
outcome held only through GitHub's case-insensitive path lookup — an untested
and unstated premise."

### What a release flip does and does not write down

Four Decisions on #272, three of them editorial and each with a precedent.
The [first](https://github.com/radiusred/gh-codecrew/issues/272#issuecomment-5561101979)
adds no changelog entry for the flip itself, "as it was on #186's and #242's
flips": "An entry describing the 2.0.0 release, sitting under the Unreleased
heading that release just created, would be a changelog line about the
changelog." The
[second](https://github.com/radiusred/gh-codecrew/issues/272#issuecomment-5561102069)
puts a `### What broke, and what to do` block at the top of the section and
accepts that the migration steps now appear twice in the file: "the block a
1.x adopter needs must be the first thing under the version heading, not four
hundred lines down, and the older block carries detail ... the lead
deliberately compresses." **Rejected:** deleting the older block, because
"editing an entry to remove content it recorded rewrites the history the file
exists to keep." The
[third](https://github.com/radiusred/gh-codecrew/issues/272#issuecomment-5561102158)
adds `migrate` to `docs/introduction.md`'s shipped-verb list — "The verb
shipped in #256 while 2.0 was on `main` and unreleased, so omitting it from a
paragraph headed **Shipped:** was correct until this flip and is wrong the
moment the claim reads v2.0.0" — with the list checked against
`internal/cli/cli.go`'s dispatch rather than against memory.

The
[fourth](https://github.com/radiusred/gh-codecrew/issues/272#issuecomment-5561318885)
is the operator's, and it answers the one question a seat declined to settle.
It is described with the Deviation that raised it, below.

### One guard or two, and a Decision taken back before review

#299 was opened and merged inside thirteen minutes and carries two Decisions
that contradict each other, which is exactly what makes it worth recording.
The
[first](https://github.com/radiusred/gh-codecrew/issues/299#issuecomment-5561376868)
guards `CLAUDE.md` only and argues scope discipline: widening the test to the
root `AGENTS.md` too "would cost one line and close that hole in the same
PR", but "this task's brief is one file plus one test ... and a test that
fails on a doc this PR was not asked to touch would drag an unrelated edit
into scope if the two ever turned out not to match." **Rejected:** covering
both here — "the narrower guard ships now; the `AGENTS.md` half is flagged in
PR #300's body for the reviewer, and is worth a backlog capture."

Two minutes and sixteen seconds later the
[follow-up](https://github.com/radiusred/gh-codecrew/issues/299#issuecomment-5561389674)
reverses it on the coordinator's instruction, before any review: the guard
becomes a table over `rootEntryPoints` checked in both directions, "so a
third root file cannot arrive unguarded". Two details in it are the reason it
reads as a record rather than a correction. "`AGENTS.md` proved already
byte-identical, so no file changed, only its ability to stop being identical
unnoticed." And it was carried "as a second commit rather than an amend, so
the record shows the flag being taken rather than silently absorbing it — and
so the pushed branch a reviewer may already hold is not rewritten."

## Deviations

Four Deviation records: one by the operator on #269, and three by the
implementer seat. One of the three was withdrawn by a later Decision on the
same issue, above. PRs [#293](https://github.com/radiusred/gh-codecrew/pull/293)
and [#300](https://github.com/radiusred/gh-codecrew/pull/300) state "no
deviations" in their descriptions, and none is recorded on #271 or #273.

**[The reviewer seat ran on a different harness again](https://github.com/radiusred/gh-codecrew/issues/269#issuecomment-5560870108).**
The operator's Deviation at 17:20:03Z: "the reviewer seat's declared harness
(codex) hit its usage limit again during M14's first review rounds
(2026-09-07, reset 19:29 local)." **Why:** "as in M13 (Deviation on #254),
the coordination layer runs the affected rounds as fresh Claude (Opus)
sessions under the reviewer contract and the same App identity
(`radiusred-checky`); the identity is what gates, the harness is the routing
table's advisory preference." It names four rounds — "#293 and #294 round one
(already Claude by choice, to parallelise), #292 round two, and #291 round
one" — and adds "codex resumes when the window returns". The parenthetical
is new relative to M13's version of the same Deviation: two of the four
rounds were Claude by choice rather than by quota, to run the four parallel
lanes at once.

**[`AdoptedRefs` reads through `StripCode` and anchors its heading](https://github.com/radiusred/gh-codecrew/issues/270#issuecomment-5560881164).**
The plan described neither, and the Deviation says the seat's own earlier
Decision "described only half the hole it was closing". **Why:** the reviewer
reproduced both holes against the branch, "and #270's own Goal carries the
literal string `## Adopts`, which shadowed the real section in the
reproduction and returned an issue nobody named."

**[The PR body cannot un-close its own task](https://github.com/radiusred/gh-codecrew/issues/272#issuecomment-5561314143).**
The milestone's most technically detailed record, and the one place a seat
refused to decide. The reviewer's change request on PR #298 asked for two things — the
literal closing keyword out of the PR body, and `closingIssuesReferences`
empty afterwards. "The first is done. **The second cannot be reached by
editing the body**, and the PR still reports `closingIssuesReferences:
[#272]`." **Why:** "the link is not body-derived. `gh codecrew task start`
creates the working branch with `gh issue develop`, which is a *linked
branch* on the issue's development panel; when a pull request is opened from
that branch GitHub converts the branch link into a linked pull request that
closes the issue on merge. That link has no public API handle" — and the
Deviation lists what was tried: `gh issue develop 272 --list` and the
`linkedBranches` GraphQL connection both empty, no public mutation for a PR
link, the surviving `ConnectedEvent` not deletable, and removing the body's
`Closes #272` producing no `DisconnectedEvent`. It then states the general
fact this exposes: "This is how every task in this hub has closed —
`task finish` does not close the task issue itself ...; it relies on exactly
this link."

The Deviation ends by refusing to decide — "**Not mine to decide,** so it
goes to the coordination layer rather than into a commit" — and lays out
three options with their costs: accept the close and run Part B on the closed
issue as M12 did with #242; reopen and re-close around the operator's
verification; or close PR #298 and open a replacement from the same branch,
"the one act that would clear the link ... at the price of the PR number, its
URL and checky's review". "I have done none of these and made no commit and
no push."

The operator's
[Decision](https://github.com/radiusred/gh-codecrew/issues/272#issuecomment-5561318885)
fifty-one seconds later takes the first: "Accepted, as in M12 (#242): Part B
— the annotated tag, the release run and the evidence — and the operator's
verification and fleet migration are posted on this issue after it closes; no
reopen." **Rejected:** the replacement PR, "it would discard the reviewed PR
and its number for a cosmetic difference in the timeline."

**[Closing #267 would have dropped its second comment](https://github.com/radiusred/gh-codecrew/issues/283#issuecomment-5560793417)**
— withdrawn. **Why:** "the dispatch brief fixes the `Closes` line, and #270's
`--adopts` link is not landed yet, so there is no other mechanism to record
the adoption." The withdrawal is itself a fact about the milestone's
sequencing: M14-R1 shipped the mechanism that would have made this Deviation
unnecessary, and M14-R5's PR merged twenty-three minutes after M14-R1's,
still under the old convention.

## The gates

**#269's Gates section was left as the scaffold's placeholder** — the
template line "_What "done" means beyond CI: e2e suites, manual UAT,
sign-offs._", unedited.
[M13's record](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#the-gates)
says the same of #254,
[M12's](12-v1-2-0-and-the-field-fixes-behind-it.md#the-gates) of #241,
[M11's](11-housekeeping.md#the-gates) of #233,
[M10's](10-protocol-bookkeeping-from-the-field.md#the-gates) of #207,
[M9's](9-the-docs-at-codecrew-works.md#the-gates) of #201 and
[M8's](8-a-product-home-page-for-codecrew-works.md#the-gates) of #196. No
`**Gate raised:**` or `**Gate resolved:**` record exists anywhere in this
milestone, no issue in this repository carries the `cc:needs-decision` label,
and `checkpoint` was not used — on a task or on the milestone issue.

That is worth saying once more here, because M14-R5's own subject is the
label `checkpoint` applies. The verb that creates `cc:needs-decision` shipped
in this milestone, the label now exists on this repository at the sampled
pink, and nothing in the milestone put it on an issue. The requirement was
verdicted from tests and from a `migrate --dry-run` reporting the labels at
their defaults, not from a raised gate.

**Every one of the six task Plans wrote a resolved Ask-the-human section, and
five of them wrote `None`.** #271's and #273's are the bare word; #270's,
#272's and #299's name the calls the seat made and say they are recorded as
Decisions. **#283's is the exception, and it is the only question raised in a Plan's
Ask-the-human section.** Its Plan, written at 17:00:42Z, opened with "**Closing
#267 drops its second comment**" and argued the case both ways; at 17:15:59Z
the same section was rewritten to "**Settled.** ... answered by the operator
through the coordination layer: `migrate` restyles, `init` and `checkpoint`
do not ... Nothing is left open."

**That question did not go through `checkpoint` either**, and neither did the
release-link escalation on #272 or the coordinator's instruction reversing
#299's first Decision. Three questions in one milestone reached the operator
and came back answered; all three were raised in a Plan or a comment, and
none used the verb the protocol provides for exactly that shape. M13's record
[made the same observation about one question](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#protocol-discipline-observations);
this milestone had three.

What gated the work: CI on every PR — `Go build and test` and `Lint commit
messages`, both required — an independent approval on each of the six PRs,
then the QA verdicts. M14-R3 added a release gate of its own that no verb
enforces: Part B and the operator's own installation checks, which the
implementer's evidence comment is explicit it cannot run — "`gh` extensions
are user-global, so a seat cannot install one in its dispatch dir".

## The review rounds

Six PRs, **thirteen review submissions**, all by the reviewer role holder
(`radiusred-checky[bot]`), and **five change requests**, every one resolved
in the next round. **Eight approvals were submitted; six stand, and two were
dismissed by a rebase force-push.**

- **PR [#292](https://github.com/radiusred/gh-codecrew/pull/292)** — two
  rounds. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/292#pullrequestreview-5126040911)
  is a single documentation finding, and it is about a sentence the PR made
  false rather than one it left stale: `docs/identities.md`'s teardown
  guidance said "anything left is listed by `git branch -r`", and after this
  PR "those branches are not shown by `git branch -r`". The
  [approval](https://github.com/radiusred/gh-codecrew/pull/292#pullrequestreview-5126052111)
  reads the new bullet back against `internal/cli/clone.go` line by line and
  then greps SPEC, docs, the contracts, the README and the help text for any
  other sentence still describing remote-only cleanup, finding none.
- **PR [#294](https://github.com/radiusred/gh-codecrew/pull/294)** — two
  rounds. Both round-one findings are the same failure the implementer's own
  Decision had named, reached by different doors, and the review says so:
  "Both reproduce against the merged branch, both are closed by one line each
  in `internal/tracker/tracker.go`, and neither has a test." It quotes the
  reproduction output for each. The
  [approval](https://github.com/radiusred/gh-codecrew/pull/294#pullrequestreview-5126079592)
  confirms both by mutation rather than by reading — reverting each fix in
  turn and naming the test that then fails — and corrects a premise stated
  earlier in the same round: running `AdoptedRefs` against #270's live body
  returns nothing, "and nothing is what it should return. **#270 has no
  `## Adopts` section** — it was opened before the flag existed."
- **PR [#293](https://github.com/radiusred/gh-codecrew/pull/293)** — three
  rounds, the third a rebase check. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/293#pullrequestreview-5126049579)
  opens by saying the change could not be broken and then blocks on the one
  new behaviour with no test: "the safety property the third Decision on #273
  states in its own words — 'an issue number that resolves to nothing at all
  is a `note:` and a skip, never a delete' — is the only new behaviour in
  this PR with no test behind it, and the suite stays green when I turn it
  into a deleter." It quotes the mutated `skip` function and the green test
  run. The
  [round-two approval](https://github.com/radiusred/gh-codecrew/pull/293#pullrequestreview-5126086714)
  tabulates ten mutations against the assertions that catch each, then was
  dismissed by the rebase onto #294's merge at 17:37:47Z. The
  [round-three approval](https://github.com/radiusred/gh-codecrew/pull/293#pullrequestreview-5126106551)
  verifies the rebase with `git range-diff` and a patch-to-patch comparison
  and names the only three real differences, all conflict resolutions.
- **PR [#291](https://github.com/radiusred/gh-codecrew/pull/291)** — three
  rounds, and the milestone's second dismissed approval. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/291#pullrequestreview-5126068576)
  is the milestone's one finding where documentation and code disagreed and
  the *code* was wrong: the CHANGELOG told an operator to rerun `migrate`
  after a failed label step, and "the rerun cannot do it ...
  `internal/cli/migrate.go:68-70` returns at the idempotence check, 120 lines
  before the label step". It runs both commands and pastes the output, then
  names the stakes: "that repository is now on 2.0 wearing GitHub's grey with
  nothing in the tool that will ever fix it — which is the exact state M14-R5
  exists to end, and the twelve repositories M14-R3 migrates are the runs
  this was added for. One flaky `gh` in the middle of those twelve produces
  it silently." The
  [round-two approval](https://github.com/radiusred/gh-codecrew/pull/291#pullrequestreview-5126118975)
  exercises the fix on a live repository whose three labels were in exactly
  the described state — one recased and mis-coloured, two wearing `#ededed`
  with no description — and checks the dry run left the listing "byte-for-byte
  unchanged". It was dismissed by the rebase onto #293's merge at 17:49:52Z;
  the
  [round-three approval](https://github.com/radiusred/gh-codecrew/pull/291#pullrequestreview-5126136448)
  is the rebase check and confirms the extra commit is "the two notes I left"
  and no production code.
- **PR [#298](https://github.com/radiusred/gh-codecrew/pull/298)** — two
  rounds. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/298#pullrequestreview-5126239552)
  is four lines long and it caught something a reading would not: the body
  claimed to be "deliberately *not* `Closes #272`" while still spelling the
  literal keyword, and `gh pr view 298 --json closingIssuesReferences`
  reported #272. The
  [approval](https://github.com/radiusred/gh-codecrew/pull/298#pullrequestreview-5126265846)
  accepts the surviving reference on the strength of the Decision recorded in
  between — "GitHub still reports #272 through the linked branch, and the
  Decision recorded on #272 accepts that flow ... so I am not holding the PR
  on that reference" — and re-checks the flip against M14-R3.
- **PR [#300](https://github.com/radiusred/gh-codecrew/pull/300)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/300#pullrequestreview-5126292938)
  at 18:56:52Z, eight minutes after opening. It checks the PR's claim "by
  construction rather than by reading": builds the binary, runs `init` into
  an empty repository and diffs the scaffold against the hub's files — "same
  189 bytes, same md5 ... and the same holds for `AGENTS.md` (172 bytes,
  identical)". It also leaves the doc-synthesizer a note, quoted in the
  observations below, about how to read this task's requirement claim.

**What the five change requests were about.** One was a real behaviour bug
found by running the code and reading the documentation against it (#291's
unreachable recovery path). Two were parser holes in a new reader, each
reproduced and each closed by one line (#294). One was a coverage gap proved
by turning a guard into a deleter and watching the suite stay green (#293).
One was a documented surface the PR had made false (#292). And one was a PR
body whose metadata contradicted its own first sentence (#298). Every
approval was reached with the reviewer running something: built binaries,
scratch repositories, a live label listing restyled and re-read, `git
range-diff` across a rebase, and in four of the six PRs — #291, #293, #294
and #300 — a mutation applied to the tree and then reverted.

**Both dismissed approvals are the same shape**, and it is a shape the
parallel lanes made inevitable. #293's round-two approval at 17:35:45Z was
dismissed at 17:37:47Z by the rebase onto #294's merge; #291's round-two
approval at 17:47:05Z was dismissed at 17:49:52Z by the rebase onto #293's
merge. In both cases round three reviewed the rebase alone and said so in its
own heading, and in both cases the reviewer verified it mechanically —
`git range-diff` plus a patch-to-patch comparison — rather than re-reading
the diff. #291 was rebased twice in all, the first time onto #294's merge
before round two.

**The App-ID check appears in the reviews of one of the six PRs.** #294's
round-two approval records "minted as `radiusred-checky`, App 4719924,
matching `GET /apps/radiusred-checky`; distinct from the PR's author
`radiusred-cody[bot]`". The other twelve submissions do not carry it in any
wording. All thirteen are by `radiusred-checky[bot]`, so the gate the
protocol enforces — a non-author holding the reviewer contract — is intact on
every PR; the seat's own identity check is not part of the record for five of
the six.

## QA: two rounds, and the requirement that was answered rather than remedied

The qa role holder (`radiusred-testy[bot]`) verdicted every requirement in a
[first comment](https://github.com/radiusred/gh-codecrew/issues/269#issuecomment-5562642135)
at 22:31:38Z. Four came back `satisfied`. **M14-R3 came back `not
satisfied`**, and what happened next is different in kind from M13's remedy
loop: no code changed, no remedy task was opened, and the fix was a Decision.

**Round one.** The R3 verdict credits the whole release side first — "installed
`gh codecrew version` prints `v2.0.0 (protocol 2.0)`, tag `v2.0.0` is
annotated by `radiusred-cody[bot]` on `2b75fdc`, release run 34052790620
succeeded with five assets, README/docs/introduction name v2.0.0, the
changelog leads with 'What broke, and what to do', and v1.2.0-and-older
changelog content is intact apart from compare links" — and then names
exactly one thing: "the #269 migration record lists ten successful migration
entries plus `davison/numberguess` refused, omits `davison/md-notes` from the
operator checklist entirely, and this QA App gets 404 for `davison/md-notes`,
so the required fleet migration record is incomplete." The verdict is against
the *record*, not against the migrations — which is the correct reading of a
requirement whose last clause is "records the run on this issue".

**The answer was a Decision, not a task.** The operator's retirement
[Decision](https://github.com/radiusred/gh-codecrew/issues/269#issuecomment-5562689587)
landed at 22:40:53Z, nine minutes and fifteen seconds after the verdict.
Nothing was migrated, nothing was built, and no PR was opened.

**Round two**, two minutes and fifty seconds later at 22:43:43Z, is a
[re-verdict scoped to R3 alone](https://github.com/radiusred/gh-codecrew/issues/269#issuecomment-5562703297):
"Re-verifying M14-R3 only; this is the standing verdict and supersedes my
earlier R3-only failure" — the supersession rule M13-R6 shipped, used for the
second milestone running. It re-runs the release checks, then says what
changed: "For the only previously open ground, I read the latest #269
Decision: `davison/md-notes` is retired, so its remote carrying 1.x files is
accounted for and not a migration gap; the R3 fleet is the eleven repos in
the migration record (ten migrated plus `davison/numberguess` refused
`MIGRATION_UNSUPPORTED` as protocol 0.1). That record plus the Decision
accounts for the #272 checklist." It closes with live spot-checks of
`radiusred/ops` and `davison/topos` default branches carrying
`.codecrew/config.yml` with `codecrew: "2.0"` and no root `.codecrew.yml`.

The other four requirements were verdicted once each, in round one, against a
build of merged `main` and against the installed release:

- **M14-R1** on the built binary refusing `task new --adopts #193` with
  `refused[ADOPT_NOT_OPEN]` and a unique-title issue search returning `[]`
  (so the refusal really did create nothing), `task finish 299 --dry-run`
  listing no adopted captures, and "a scratch probe confirmed a trailing ref
  in `## Adopts` prose is ignored while CRLF fails closed as #296 records".
- **M14-R2** on the targeted clone tests "including the temp-repo case where
  a local task branch with an unpushed commit is switched off, named, kept,
  and left at the extra commit", plus the live evidence from #272's own
  finish, which "fast-forward[ed] local `main` and delet[ed] the local task
  branch".
- **M14-R4** on the guard tests "for stale delete/keep decisions, open-task
  skip, open-PR guard, rebase-merge handling, unreadable candidates, and
  truncated listings", with `#167` cited for the live stale-branch evidence.
  The verdict is explicit about what it did *not* do: "I did not run a live
  close; `milestone close 14 --dry-run` currently stops at
  `refused[OPEN_TASKS]` for #302 before branch planning, which is the current
  post-QA doc task state rather than stale-sweep behavior."
- **M14-R5** on the hub's own label list matching the recorded palette and
  descriptions, targeted tests for `init` in hub and spoke and for `migrate`'s
  restyle and rerun, and `gh codecrew migrate --dry-run` from the installed
  release reporting "labels already at the protocol defaults".

Cross-cutting, both rounds: `go build`, `go test ./...`, `go vet ./...` and
`gofmt -l .` clean on merged `main`, the second round noting the caches were
kept under the QA workspace. Round one also closes with the capture
bookkeeping — "captures #193/#192/#167/#267 are closed and #295/#296/#297
remain open" — and files
[#303](https://github.com/radiusred/gh-codecrew/issues/303) for what it found
on PR #294's body, described under the observations below.

## Requirement outcomes

The status column is the **standing** verdict word as the qa role holder
wrote it, verbatim and unqualified — the latest comment carrying a verdict
for the ID, which is the supersession rule M13-R6 shipped. M14-R3's earlier
verdict is in the Notes column and in the section above.

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M14-R1 — backlog adoption is a first-class link: `task new --adopts <ref>[,<ref>]` records the adopted issues under an `## Adopts` line and comments on each capture; `task finish` closes each after the merge with a comment pointing at the task and its PR; `--dry-run` lists them; an already-closed capture is reported, not an error; SPEC §6, the coordinator and implementer contracts follow (adopts [#193](https://github.com/radiusred/gh-codecrew/issues/193)) | [#270](https://github.com/radiusred/gh-codecrew/issues/270) / PR [#294](https://github.com/radiusred/gh-codecrew/pull/294) | `satisfied` | Task closed; verdicted in QA round one on the built binary's `ADOPT_NOT_OPEN` refusal, a dry-run finish and a scratch parser probe. Two review rounds; both round-one findings were parser holes reproduced against the branch, and the fix put `AdoptedRefs` on M13-R6's shared code-stripping rule. `ADOPT_NOT_OPEN` takes SPEC §10's table from forty-two codes to forty-three |
| M14-R2 — `task finish` cleans up locally too: from a clone whose current or local branch is the task's, it switches to the default branch, fast-forwards it and deletes the local task branch after the merge, reporting each step; a branch with unpushed commits is never deleted and is named instead; nothing happens from elsewhere; `--dry-run` prints the local steps beside the remote ones (adopts [#192](https://github.com/radiusred/gh-codecrew/issues/192)) | [#271](https://github.com/radiusred/gh-codecrew/issues/271) / PR [#292](https://github.com/radiusred/gh-codecrew/pull/292) | `satisfied` | Task closed; verdicted in QA round one on temp-repo clone tests and on the live finish that merged PR #298. Two review rounds; the one finding was a teardown sentence in `docs/identities.md` the change had made false. Three Decisions, one of which departs from the requirement's literal wording and says why |
| M14-R3 — v2.0.0 ships: the `[Unreleased]` changelog becomes the 2.0.0 section leading with what broke and the migration steps; `docs/introduction.md` and README flip; the tag is cut by the implementer identity; `release.yml` builds it; the operator verifies the installed release and migrates the remaining repos, hubs before their spokes, recording the run on #269 | [#272](https://github.com/radiusred/gh-codecrew/issues/272) / PR [#298](https://github.com/radiusred/gh-codecrew/pull/298), and [#299](https://github.com/radiusred/gh-codecrew/issues/299) / PR [#300](https://github.com/radiusred/gh-codecrew/pull/300) under its spirit | `satisfied` | Task closed; **two QA rounds**. `not satisfied` at 22:31:38Z — the migration record accounted for eleven of the checklist's twelve repositories — then `satisfied` at 22:43:43Z after the operator's retirement Decision for `davison/md-notes`, scoped "Re-verifying M14-R3 only". Tag `v2.0.0` (tag object `93bf2f3` → commit [2b75fdc](https://github.com/radiusred/gh-codecrew/commit/2b75fdc)) by `radiusred-cody[bot]`, run 34052790620, five assets. Ten repositories migrated, `davison/numberguess` refused `MIGRATION_UNSUPPORTED` by design |
| M14-R4 — `milestone close` sweeps stale task branches from earlier closes too: after its own tasks' branches it deletes any `task/<n>-<slug>` branch whose task issue is closed and whose branch is merged or empty, reporting each in the closing comment, and never touches a branch with unmerged commits or one whose task is open (adopts [#167](https://github.com/radiusred/gh-codecrew/issues/167)) | [#273](https://github.com/radiusred/gh-codecrew/issues/273) / PR [#293](https://github.com/radiusred/gh-codecrew/pull/293) | `satisfied` | Task closed; verdicted in QA round one on the guard tests and #167's recorded evidence, with the verdict stating that no live close was run. Three review rounds, the third a rebase check; round one's finding was an untested safety guard that stayed green when turned into a deleter. Five Decisions, two of them added by the review. #167's second half is captured as [#295](https://github.com/radiusred/gh-codecrew/issues/295) |
| M14-R5 — the protocol's labels exist before they are needed: `init` (hub and spoke) creates `cc:milestone`, `cc:task` and `cc:needs-decision` when missing, each with a fixed description and a colour drawn from the crew image palette, reporting what it created and leaving an existing label untouched; `checkpoint` creates `cc:needs-decision` rather than failing on a repo's first gate (adopts [#267](https://github.com/radiusred/gh-codecrew/issues/267)) | [#283](https://github.com/radiusred/gh-codecrew/issues/283) / PR [#291](https://github.com/radiusred/gh-codecrew/pull/291) | `satisfied` | Task closed; verdicted in QA round one against this hub's own label list, the targeted tests and the installed release's `migrate --dry-run`. Three review rounds, two of them after a rebase; round one found the documented recovery path unreachable in code. Five Decisions, including the sampled palette and a scope Decision the operator's answer superseded — `migrate` restyles, which R5's text does not ask for |

**Captures adopted and closed by the merges — four**, each by the PR
delivering the requirement that adopted it:
[#192](https://github.com/radiusred/gh-codecrew/issues/192) by PR #292,
[#193](https://github.com/radiusred/gh-codecrew/issues/193) by PR #294,
[#167](https://github.com/radiusred/gh-codecrew/issues/167) by PR #293 and
[#267](https://github.com/radiusred/gh-codecrew/issues/267) by PR #291. All
four were closed by a `Closes` line in the PR body, including #193 — the
capture whose delivery *is* the mechanism that replaces that convention. Its
PR body says so: "the last capture closed by the PR-body convention. From the
next task that adopts one, `task new --adopts` records it and `task finish`
closes it."

**Captures filed during the milestone and left open — five.**
[#295](https://github.com/radiusred/gh-codecrew/issues/295) (`status` reports
stale task branches, split out of #167 by the implementer at the reviewer's
request, before the merge closed it), [#296](https://github.com/radiusred/gh-codecrew/issues/296)
(the tracker's line-anchored regexps miss CRLF bodies, from the reviewer's
round two on PR #294),
[#297](https://github.com/radiusred/gh-codecrew/issues/297) (nothing asserts
that `task new` applies `cc:task`, from the reviewer's round two on PR #291),
[#301](https://github.com/radiusred/gh-codecrew/issues/301) (`migrate` prints
`action needed` for an *absent* root entry point instead of writing it as
`init` does, from the operator's own question while migrating the fleet) and
[#303](https://github.com/radiusred/gh-codecrew/issues/303) (a PR body's
prose can create unintended closing references, from the QA seat's M14
verification). Three of the five come from review rounds, one from the
operator running the release's own verb on real repositories, and one from
QA reading a PR's metadata rather than its diff.

## Protocol-discipline observations

Seven things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **A milestone shipped the mechanism that would have made its own
  bookkeeping unnecessary, and could not use it.** M14-R1 landed `task new
  --adopts` at 17:33:29Z; M14-R5's PR merged twenty-three minutes later
  carrying `Closes #267` in its body, and #283's Deviation says why — "#270's
  `--adopts` link is not landed yet, so there is no other mechanism to record
  the adoption". All four of this milestone's adopted captures were closed by
  the convention M14-R1 exists to replace, because all four tasks were opened
  before it existed. The two tasks opened after the flag landed — #299 and
  this record's own #302 — adopt nothing, so this record is the boundary
  where the old convention was still the only one in use.
- **A `not satisfied` verdict was cleared by a Decision, with nothing built.**
  M13's two `not satisfied` verdicts each opened a remedy task and merged a
  PR. M14's one was against the *record* rather than the code — a fleet entry
  missing from an operator comment — and the fix was the operator writing down
  what had happened to the repository the record left out. Nine minutes from
  verdict to
  Decision, three more to the re-verdict. That is the protocol working on a
  class of failure it has no verb for: the requirement said "records the run
  on this issue", QA checked the issue, and the remedy was prose.
- **The record-grammar rule M13 shipped is now load-bearing in two
  directions.** M13-R6's supersession rule was exercised on its own
  requirement by hand; here it carried a re-verdict for a different
  milestone's requirement, written by the same seat in the same form ("this
  is the standing verdict and supersedes my earlier R3-only failure"). And
  M13-R6's other half — one code-stripping implementation shared by every
  record reader — gained its third reader under review pressure rather than
  by design: the reviewer's finding on PR #294 was, in as many words, that
  the new reader "is the third record-reading scan in the package and the
  only one that does not use it — while being the only one whose output
  closes other people's issues".
- **Four parallel lanes produced two dismissed approvals and two rebase-check
  rounds, and the cost was paid entirely by the reviewer.** #291, #292, #293
  and #294 opened within two minutes of each other and merged in a different
  order than they opened. Nothing in the protocol sequences PRs, and nothing
  needs to — the merges were clean and both rebase checks passed — but two
  of the milestone's thirteen review submissions are rebase checks that exist
  only because of the ordering, and two approvals were thrown away by
  force-pushes that changed no behaviour.
  [M13's record](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#the-review-rounds)
  describes one such dismissal; this milestone had two, twelve minutes
  apart.
- **Three questions went to the operator; none went through `checkpoint`.**
  #283's Plan raised one and was rewritten when it was answered; #272's
  Deviation raised one and refused to decide it; #299's first Decision was
  reversed by a coordinator instruction before review. All three are
  well-recorded — a Plan edit, a Deviation with three costed options, and a
  follow-up Decision carried as a second commit "so the record shows the flag
  being taken". None of them raised a gate, applied `cc:needs-decision` or
  called the verb. M13's record observed the same of its single escalation:
  four questions across two milestones have now taken this route, which
  suggests the gap is in the dispatch habit rather than in any one seat's
  judgement.
- **The refresh obligation found nothing false this time, and the reason is
  the release.**
  [M13's record](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#protocol-discipline-observations)
  reported that `README.md`'s "Start now" block paired
  `gh extension install radiusred/gh-codecrew` with a comment saying `init`
  "writes and commits `.codecrew/`, `AGENTS.md`, `CLAUDE.md`, `ROADMAP.md`",
  while the released v1.2.0 `init` wrote `.codecrew.yml` and `roles/`. Run
  again for this document from the installed v2.0.0, in a scratch repository,
  `init` writes exactly what the block says. The page did not change; the
  binary caught up with it. *The check that matters is still the one PR #290
  named — would a reader who followed this page get what it describes — and
  it is only answerable by running the release.*
- **The reviewer left the doc-synthesizer a reading, in the review, for the
  one task whose requirement claim is a stretch.** PR #300's approval records
  it as a non-blocking note: "M14-R3 as worded on #269 is the release flip,
  the tag and the fleet migration; this is none of those. The task's Goal and
  Plan and the PR body all say so in the same words ('its spirit … no new
  requirement') ... Nothing to fix — recorded so the milestone document has
  the reading in front of it: if the synthesizer later wants R3's evidence to
  be only the release, this PR is the one task under it that is not." This
  record takes that reading: the R3 row names both tasks and marks which one
  is under the requirement's spirit rather than its text.

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **The conversation that opened the milestone is on the record only through
  #269's goal.** "Decided with the operator 2026-09-06" is the whole of it.
  The two captures the goal names carry their own history, but the judgement
  that they belong together in one milestone with the release — the thesis
  the goal states in one sentence — has no transcript. *The reasoning is
  recorded; the deliberation is not*, as
  [M13's record](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#what-the-record-does-not-contain)
  says of its own opening.
- **Which reviewer rounds ran on which harness is mostly not recorded.** The
  Deviation names four of the thirteen submissions and says codex resumes
  when the window returns; PR #294's round-two approval self-attributes as a
  Claude-run review. The remaining eight — including every round on #298 and
  #300 — are not attributed anywhere, and the Deviation's own reset time
  cannot settle them: it reads "2026-09-07, reset 19:29 local" on a comment
  GitHub timestamps 2026-09-06T17:20:03Z, and two of the milestone's captures
  ([#296](https://github.com/radiusred/gh-codecrew/issues/296),
  [#297](https://github.com/radiusred/gh-codecrew/issues/297)) likewise date
  themselves 2026-09-07 in prose while being created on 2026-09-06. Every
  time in this record is taken from the API rather than from the prose for
  that reason. The identity is constant throughout, so the gate the
  protocol enforces is intact; the harness behind each round is not
  reconstructable, exactly as
  [M13's record](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#what-the-record-does-not-contain)
  reports of its own.
- **The six captures the M14-R4 Decision read and did not bundle are named by
  number and nothing else.** "Reviewed and not bundled: #108 (additive verb,
  any minor), #191 (docs), #134, #194, #36, #31." Four of the six carry no
  reason at all. A blessing and a deferral are both decisions, and these four
  are recorded as neither.
- **The fleet migration is eleven `git log` excerpts, and nothing else.** The
  migration record shows two commits per repository and, for
  `davison/numberguess`, one refusal line. It does not show what `migrate`
  printed on any successful run — no `action needed` blocks, no label
  receipts, no notes — and the operator's checklist step 3 asked for "anything
  migrate refused or noted". [#301](https://github.com/radiusred/gh-codecrew/issues/301)
  exists because the operator noticed one such behaviour while running it and
  filed it by hand; whether the other nine runs printed anything worth
  keeping is not on the record. *This is the first fleet-wide exercise of
  M13-R2 against repositories nobody built for the test, and what it
  reports is the commits rather than the runs.*
- **`davison/md-notes`'s current state is asserted, not shown.** The
  retirement Decision says the working directory "was removed and re-created"
  and that the remote "may still carry the 1.x files". The QA seat's probe
  returned a 404 under its own token, which distinguishes nothing between
  private, renamed and deleted. Nobody looked, and the Decision does not
  claim anyone did.
- **M14-R4 has never run.** The verdict says so in its own words: "I did not
  run a live close." The second branch pass is verdicted from unit tests and
  from #167's recorded evidence that two stale branches exist in
  `radiusred/numberguess`; the first close that will exercise it is this
  milestone's, after this document merges. Its `--dry-run` cannot preview it
  either, because the task-count gate refuses before branch planning is
  reached — which is a property of the verb, not of the test.
- **Nothing measures what the milestone cost.** Six PRs, thirteen reviews,
  two QA rounds and a fleet migration, with no record of tokens, wall-clock
  per seat, or how many dispatches were made.
  [M13's record](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#what-the-record-does-not-contain)
  says the same of its own, and the answer has not changed.
- **The accidental closing reference on PR #294 is captured but not
  explained.** The PR closed
  [#42](https://github.com/radiusred/gh-codecrew/issues/42) as well as #193
  and #270, from prose about the example title "fix #42 in the parser" — the
  very hazard M14-R1's own parser was hardened against, arriving through
  GitHub's separate PR-body parser instead. #42 had been closed since
  2026-08-22 by PR #43, so nothing changed state.
  [#303](https://github.com/radiusred/gh-codecrew/issues/303) records the
  observation and proposes a check; which exact wording GitHub matched is not
  established there, and nobody has reproduced it.
