# M16: Bedding in: the venue seam and the reference

Tracking issue: [#325](https://github.com/radiusred/gh-codecrew/issues/325) ·
Synthesized 2026-09-07 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #325's **six** requirements as the body stands, **seven
as it was opened** — one withdrawn and one renumbered by the operator's
Decision, the whole of the difference; its Gates section, left as the
scaffold's placeholder; its **four comments**, carrying **one Decision, two
Deviations and one QA comment holding all six verdicts**; the **eight**
delivery task issues, five in this repository
([#326](https://github.com/radiusred/gh-codecrew/issues/326),
[#327](https://github.com/radiusred/gh-codecrew/issues/327),
[#328](https://github.com/radiusred/gh-codecrew/issues/328),
[#329](https://github.com/radiusred/gh-codecrew/issues/329),
[#330](https://github.com/radiusred/gh-codecrew/issues/330)) and three in the
spoke ([radiusred/codecrew-www#30](https://github.com/radiusred/codecrew-www/issues/30),
[radiusred/codecrew-www#33](https://github.com/radiusred/codecrew-www/issues/33),
[radiusred/codecrew-www#35](https://github.com/radiusred/codecrew-www/issues/35)),
plus [#340](https://github.com/radiusred/gh-codecrew/issues/340), the task this
document is; their **eight merged pull requests** and the **twenty-seven
commits** on them — twenty-four here and three in the spoke; the **twenty-five
Decision records and seven Deviation records** across the milestone issue and
all nine of its task issues, of which **twenty-three and seven** are on the
milestone issue and the eight delivery tasks and the remaining **two** on this
document's own; the **twelve review submissions** on the eight pull requests —
**four change requests** and **eight approvals, every one of which stands**,
none dismissed; the **six** adopted backlog captures the merges closed; and the
**six** captures the milestone filed and left open — five in this hub and one in
the spoke. The house form is
[M15's document](15-v2-0-1-what-the-fleet-migration-taught.md), with
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md) and
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md) behind
it; the record standard is the one the reviewer set on
[PR #304](https://github.com/radiusred/gh-codecrew/pull/304) and re-applied on
[PR #322](https://github.com/radiusred/gh-codecrew/pull/322), so every count
above and below is re-derived from the source — records by the rule
`tracker.ExtractRecords` itself applies (a paragraph-initial `**Decision:**`,
`**Deviation:**` or `**Gate resolved:**` label, bare or parenthetically
qualified), review states from each pull request's reviews API and timeline,
commits from each pull request's commits API, and the requirement history from
the milestone issue's `userContentEdits`.

**One of the twenty-five Decision records is a gate resolution.** The
`**Gate resolved:**` label on #327 is the milestone's one gate.
`ExtractRecords` gathers a gate resolution as a Decision (SPEC §8), so it is
counted as one here; it is also counted as a gate
in the gates section below, and it is the same record both times.

**Counted at 2026-09-07T14:09:08Z, against a trail frozen from this document's
dispatch.** The coordination layer posted nothing on the milestone after the QA
comment at 14:01:10Z and posts nothing further before this pull request merges.
This document's own task's **two Decisions** were written before the count was
taken and are inside every total above, so a reader running the same queries
gets the same numbers.
[M15's record](15-v2-0-1-what-the-fleet-migration-taught.md) had its counts
expire between its first commit and its review and gave both figures;
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md) left its own task's records
out. This one has one figure because the records were written first.

**The requirements were amended once, and the amendment is the milestone's
first Decision.** #325's `userContentEdits` query returns a total count of
**two** — the body as it was created at 12:05:54Z and one edit at 12:09:53Z,
two seconds before the Decision that explains it. Nothing else in the body
moved. The seven requirements a reader would have seen at 12:05:54Z are not the
six the issue carries now, and the section below on the withdrawn requirement is
the whole of that difference.

**The trail this record is compiled from is checked by the verb the milestone
did not change.** `gh codecrew milestone evidence` walks the milestone issue and
its sub-issues — bodies and comments, which is where every citation below comes
from — and refuses when a github.com citation does not resolve. The installed
extension reports `v2.0.1 (protocol 2.0)`, which is also what a binary built
from `main` at [42ab6f8](https://github.com/radiusred/gh-codecrew/commit/42ab6f8)
reports, because M16 cut no release. Run from this repository at the same
instant as the counts above:

```
$ gh codecrew milestone evidence 16
requirements counted: M16-R1, M16-R2, M16-R3, M16-R4, M16-R5, M16-R6 (6)
all 4 cited links resolve across 10 issues — evidence is reachable
```

Ten issues: the milestone and its nine sub-issues — the eight delivery tasks,
and this document's own. Four links, because this milestone's prose cites
almost entirely by bare `#N` reference: the only github.com URLs in the whole
trail outside code are one in the QA comment and three in this document's own
Decision comments. The QA seat ran the same verb before those three existed and
recorded `all 0 cited links resolve across 9 issues`; both readings are correct
at their own instant, and the difference is this task.

The gather from `gh codecrew milestone close 16 --dry-run` is again not part of
the raw material, for the reason
[M11's](11-housekeeping.md),
[M12's](12-v1-2-0-and-the-field-fixes-behind-it.md),
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md),
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md) and
[M15's](15-v2-0-1-what-the-fleet-migration-taught.md) records all give:
`--dry-run` stops at the gate that counts tasks, naming the task that writes
this file. It is quoted as it read before this document's own pull request
opened; with that pull request open the same gate names #340 `(in review)`.

```
gate milestone open: ok
gate no gate raised: ok
gate tasks closed: refused[OPEN_TASKS]: tasks not closed: radiusred/gh-codecrew#340 (ready)
gate requirements declared: not reached
gate QA verdicts: not reached
gate milestone document: not reached
dry run: nothing written — the live verb stops at the first refusal above
```

Every record below was read from its issue, pull request, review, commit or
timeline directly. The prose is wrapped at this file's normal width; the
requirement-outcomes table's rows and a handful of single links are not, because
a newline ends a Markdown table row and breaking a link breaks it — the same
shape
[M15's record](15-v2-0-1-what-the-fleet-migration-taught.md#requirement-outcomes)
and [M14's](14-adoption-tidy-and-the-v2-0-0-release.md#requirement-outcomes)
carry, and the shape `ROADMAP.md` has had since M1.

**This PR adds the M16 ROADMAP row; it does not flip one** — the convention
M10-R1 introduced and M12-R1 shipped, which
[M15's record PR](https://github.com/radiusred/gh-codecrew/pull/322),
[M14's](https://github.com/radiusred/gh-codecrew/pull/304) and
[M13's](https://github.com/radiusred/gh-codecrew/pull/290) were written under.
`ROADMAP.md` carried no M16 row before this PR, so there was nothing to discard
and the trail carries no Decision discarding one.

## Goal and outcome

#325 was opened at 12:05:54Z on 2026-09-07, nine hours and twenty-four minutes
after
[`milestone close 15`](https://github.com/radiusred/gh-codecrew/issues/305#issuecomment-5564293567)
closed M15 with the line "all 8 tasks done, milestone document merged". Its goal
opens on a rule rather than on a list, and the rule is the milestone's whole
shape: "Protocol 2.0 beds in. Nothing in this milestone changes the protocol or
what an operator or a seat sees: no verb gains or loses an argument, changes its
output, its writes or its refusal codes; the pointer gains no key; SPEC §4, §5,
§8 and §10 keep their meaning."

Under that rule the goal admits exactly two kinds of work. "The first is
architecture with no surface: the venue seam
([#194](https://github.com/radiusred/gh-codecrew/issues/194)) — every use of gh
and of the GitHub API behind one interface, GitHub the only implementation, no
second venue and nothing that names one; the GitLab notes on #194 were a
feasibility check and stay notes; the gain now is a fake venue that makes the
tests simpler and honest." The second is documentation: "the CLI reference that
replaces SPEC §6 and sits at the root
([#323](https://github.com/radiusred/gh-codecrew/issues/323), the operator's
decisions there), and the docs captures with no behaviour in them — the
self-gating workflow for docs-only PRs
([#191](https://github.com/radiusred/gh-codecrew/issues/191)), working offline
([#268](https://github.com/radiusred/gh-codecrew/issues/268)'s docs half; its
behaviour half is [#324](https://github.com/radiusred/gh-codecrew/issues/324)),
external contributors
([#31](https://github.com/radiusred/gh-codecrew/issues/31)) — plus one
site-mechanics fix in the spoke
([radiusred/codecrew-www#29](https://github.com/radiusred/codecrew-www/issues/29))."
The provenance line is one sentence: "Decided with the operator 2026-09-07."

Each requirement adopts exactly one capture, and the mapping is one-to-one with
no requirement taking a pair and none adopting nothing: M16-R1 adopts #194,
M16-R2 adopts #323, M16-R3 adopts #191, M16-R4 adopts #268, M16-R5 adopts #31,
and M16-R6 adopts
[radiusred/codecrew-www#29](https://github.com/radiusred/codecrew-www/issues/29).
That is a different shape from
[M15's](15-v2-0-1-what-the-fleet-migration-taught.md#goal-and-outcome), where
two requirements took a pair and one — the release — adopted nothing.

The eight delivery tasks, in merge order:

- **[#330](https://github.com/radiusred/gh-codecrew/issues/330) / PR
  [#331](https://github.com/radiusred/gh-codecrew/pull/331)** — M16-R5. One
  commit, approved first round, merged 12:22:09Z
  ([0ab7586](https://github.com/radiusred/gh-codecrew/commit/0ab7586)). Two
  Decisions. Adopts
  [#31](https://github.com/radiusred/gh-codecrew/issues/31).
- **[radiusred/codecrew-www#30](https://github.com/radiusred/codecrew-www/issues/30)
  / PR [radiusred/codecrew-www#31](https://github.com/radiusred/codecrew-www/pull/31)**
  — M16-R6. One commit, approved first round, merged 12:29:14Z
  ([037f7c7](https://github.com/radiusred/codecrew-www/commit/037f7c7)). Two
  Decisions. Adopts
  [radiusred/codecrew-www#29](https://github.com/radiusred/codecrew-www/issues/29).
- **[#328](https://github.com/radiusred/gh-codecrew/issues/328) / PR
  [#332](https://github.com/radiusred/gh-codecrew/pull/332)** — M16-R3. Two
  commits, approved first round, merged 12:37:26Z
  ([00b254b](https://github.com/radiusred/gh-codecrew/commit/00b254b)). Two
  Decisions and three Deviations, more Deviations than any other task in the
  milestone. Adopts
  [#191](https://github.com/radiusred/gh-codecrew/issues/191).
- **[#327](https://github.com/radiusred/gh-codecrew/issues/327) / PR
  [#335](https://github.com/radiusred/gh-codecrew/pull/335)** — M16-R2. Seven
  commits, two review rounds, merged 13:09:16Z
  ([67256a3](https://github.com/radiusred/gh-codecrew/commit/67256a3)). Three
  Decisions and the milestone's one gate, raised and resolved on the issue.
  Adopts [#323](https://github.com/radiusred/gh-codecrew/issues/323).
- **[radiusred/codecrew-www#33](https://github.com/radiusred/codecrew-www/issues/33)
  / PR [radiusred/codecrew-www#34](https://github.com/radiusred/codecrew-www/pull/34)**
  — M16-R2's spoke half. One commit, approved first round, merged 13:22:31Z
  ([9132d54](https://github.com/radiusred/codecrew-www/commit/9132d54)). No
  Decisions: the task's own plan says the two rows and their position "were both
  settled by the milestone requirement and by the hub PR #335's follow-up
  paragraph". Adopts nothing — it is the spoke half of a hub requirement, opened
  after the hub half merged.
- **[#329](https://github.com/radiusred/gh-codecrew/issues/329) / PR
  [#333](https://github.com/radiusred/gh-codecrew/pull/333)** — M16-R4. Eight
  commits, three review rounds, merged 13:22:37Z
  ([0981557](https://github.com/radiusred/gh-codecrew/commit/0981557)). Three
  Decisions and one Deviation, and the two review rounds that produced the
  Deviation. Adopts
  [#268](https://github.com/radiusred/gh-codecrew/issues/268).
- **[radiusred/codecrew-www#35](https://github.com/radiusred/codecrew-www/issues/35)
  / PR [radiusred/codecrew-www#36](https://github.com/radiusred/codecrew-www/pull/36)**
  — M16-R4's spoke half. One commit, approved first round, merged 13:32:38Z
  ([1060594](https://github.com/radiusred/codecrew-www/commit/1060594)). One
  Decision. Adopts nothing, for the same reason as www#33.
- **[#326](https://github.com/radiusred/gh-codecrew/issues/326) / PR
  [#336](https://github.com/radiusred/gh-codecrew/pull/336)** — M16-R1, and the
  only code in the milestone. Six commits, two review rounds, merged 13:44:35Z
  ([42ab6f8](https://github.com/radiusred/gh-codecrew/commit/42ab6f8)). Eight
  Decisions and one Deviation, more Decisions than any other task in the
  milestone. Adopts
  [#194](https://github.com/radiusred/gh-codecrew/issues/194).

From the milestone opening to the last delivery merging: **one hour,
thirty-eight minutes and forty-one seconds**. From the first `task start` —
#328's, at 12:11:04Z — to that merge: one hour, thirty-three minutes and
thirty-one seconds. The QA verdicts landed at 14:01:10Z, sixteen minutes and
thirty-five seconds after the last merge.

**Two of the eight tasks did not exist when the milestone opened.** The five hub
tasks and the first spoke task were created between 12:06:02Z and 12:06:26Z, in
one batch, eight to thirty-two seconds after the milestone issue. The other two
spoke tasks were opened later and each one after the hub half it follows:
www#33 at 13:10:17Z, sixty-one seconds after #335 merged, and www#35 at
13:24:14Z, ninety-seven seconds after #333 merged. Both were named as follow-up
in the hub pull request that made them necessary — #335's body specifies the two
`sync_docs.py` rows by name, and #333's page had to be placed in the site's
navigation once it existed — so the site's reading order matches the
introduction's at every point where both were true at once.

**Four hub pull requests were open at once**, between 12:31:07Z and 12:37:26Z:
#332, #333, #335 and #336. #331 opened at 12:15:14Z and merged before any of
them, and the three spoke pull requests never overlapped each other.
[M15's record](15-v2-0-1-what-the-fleet-migration-taught.md#goal-and-outcome)
reports five hub lanes open at once and
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md#goal-and-outcome) four.

What exists afterwards that did not before. Every invocation of `gh` and every
GitHub REST or GraphQL call the CLI makes goes through `tracker.Tracker`;
`internal/gh` has exactly one importer, and a test that reads the whole module's
source fails on an `internal/gh` import, a `gh` exec, an `api.github.com`
literal or a `github.com` URL anywhere outside the two venue packages, with one
recorded exception. `internal/tracker/faketracker` implements the whole
interface — forty-three methods — and records every call. `CLI.md` is a
reference at the repository root: twenty-one `##` sections in 925 lines, four
shared and one for each of the seventeen verbs `--help` lists, and SPEC §6 is a
stub that keeps its number so §7–§13 keep theirs. `docs/first-milestone.md` §5
carries a self-gating workflow for a repository whose gate is expensive, and the
`NO_CHECKS` refusal names the `skipping` path that satisfies it.
`docs/working-offline.md` says which verbs run with no network, which wait and
exactly what each prints, with six of its quotations rebuilt from the code by a
test. `docs/introduction.md` carries the conventions for a project whose issue
list takes reports from strangers. The site publishes the reference and the
offline page top-level, in the introduction's own reading order, and the blog's
nav label carries no HTML.

**M16 shipped no release, and nothing in the trail says otherwise.** No tag was
cut, `docs/introduction.md`'s `**Shipped:**` line still reads v2.0.1 and the
installed extension still prints `v2.0.1 (protocol 2.0)`. SPEC §10's table still
counts **forty-three** refusal codes: this milestone added none, which is what
its Goal's rule means when it is measured rather than asserted. What is on
`main` and unreleased is in `CHANGELOG.md`'s `## [Unreleased]` section, whose
six entries this pull request makes seven — five from this milestone's tasks,
one from M15's own record, and this document's.

## The requirement that was withdrawn

The milestone opened with seven requirements. Three minutes and fifty-nine
seconds later it had six, and the whole of the change is on the record twice —
once as the body edit, once as the Decision that explains it.

The requirement first declared as M16-R6 was "migrating an OpenGSD project to
CodeCrew": a `docs/migrating-from-opengsd.md` with a cutover half and an
audit-trail half, "with a recorded Decision for each of the four mechanisms #36
lists (a retroactive `docs/milestones/0-pre-codecrew.md` synthesized from the
GSD archives and citing them by SHA; a ROADMAP tombstone row for the OpenGSD
releases; which load-bearing GSD decisions are re-recorded as comments on the
first milestone; `.planning/` archived in place, never deleted)". It adopted
[#36](https://github.com/radiusred/gh-codecrew/issues/36), and the spoke's blog
requirement was M16-R7.

The [operator's Decision](https://github.com/radiusred/gh-codecrew/issues/325#issuecomment-5570412243)
at 12:09:55Z withdraws it and renumbers M16-R7 → M16-R6. Its reason is a
refusal: "`task new --adopts 36` refused `ADOPT_NOT_OPEN`: the operator closed
#36 as not planned on 2026-09-05, and the coordination layer scoped this
milestone from the capture's text without reading its state. The 2026-09-05
decision stands; nothing about #36 is reopened here." **Trade-off:** "the body
is edited after creation, which the M15 record noted as a thing this hub had not
done; the edit is this Decision's subject and is recorded before any task
starts." **Rejected:** "leaving R6 declared with no task — `milestone close`
would demand a verdict on work nobody did."

Three things about that are worth keeping. The first is that a verb caught it:
`task new --adopts` reads the capture's state, and `ADOPT_NOT_OPEN` is what
stopped a milestone from carrying a requirement whose whole subject had been
declined two days earlier. Nothing else in the trail would have noticed —
`milestone new` writes the issue and does not read what its requirements name.
The second is the timing: the edit landed at 12:09:53Z and the first `task
start` at 12:11:04Z, so no task ever ran against the seven-requirement body.
The third is what the withdrawal cost the record — #36 is still closed as not
planned, and `gh codecrew milestone evidence 16` counts six requirements
because six is what the body declares.

## Decisions

Twenty-five Decision records: one on the milestone issue, twenty-two on the
eight delivery tasks and two on this document's own. Two of the twenty-five are
the operator's — the withdrawal above, and the gate resolution on #327 — and the
rest were written by the seat doing the work: `radiusred-cody[bot]` on all eight
delivery tasks, `radiusred-wordy[bot]` on this document's.

One carries a parenthetical qualifier rather than the bare label —
`**Decision (operator, 2026-09-07):**` on #325 — and one is a gate resolution
rather than a decision label. Both are records under `ExtractRecords`' rule and
both are invisible to a search for `**Decision:**` alone, exactly as
[M15's](15-v2-0-1-what-the-fleet-migration-taught.md#decisions) one qualified
label and [M14's](14-adoption-tidy-and-the-v2-0-0-release.md#decisions) two
were. The other twenty-three are bare. They group by the question each task had
to settle.

### The seam that kept its name

M16-R1's subject is a boundary, and #326's eight Decisions run from what to call
it out to what the guard is allowed to forgive.

The [first](https://github.com/radiusred/gh-codecrew/issues/326#issuecomment-5570608772)
keeps the seam's name: `tracker.Tracker` in `internal/tracker`, extended in
place. **Trade-off:** "The seam's name says 'tracker' where the milestone says
'venue' — against a rename touching every file in `internal/cli` and
`internal/tracker` for no behavioural gain, in a milestone whose rule is that
nothing an operator or a seat sees changes." **Rejected:** `venue.Venue` with
`Tracker` as an alias, on the ground that "an alias leaves two names for one
interface — which is the drift this task exists to stop, in a different form."
The name a reader arrives at is documented instead, in the interface's doc
comment, the package's, and `docs/founding-decisions.md`; "the identifier they
have to type is the one already in forty files."

The [second](https://github.com/radiusred/gh-codecrew/issues/326#issuecomment-5570609002)
draws the line between the interface and the package around it: I/O became nine
interface methods, and the venue's pure knowledge became package functions
beside them. **Rejected:** making `RecordAPIPath`, `IsRecordLink`, the URL
builders, `Unreachable` and `CompareVersions` methods — "They are knowledge, not
calls: there is nothing to script, nothing to record, and a fake that could
answer them differently would let a CLI test pass against a grammar the venue
does not have." The rule it states is the one a reader can apply to the next
addition: "if a fake would want to lie about it, it is a method; if lying about
it would only weaken a test, it is a function."

The [third](https://github.com/radiusred/gh-codecrew/issues/326#issuecomment-5570609241)
keeps every refusal in `internal/cli`. The App-JWT surface is one low-level
method, `AppRequest`, handing back a status and a body. **Rejected:** typed
`AppInstallations` and `MintInstallationToken` methods, because each "would have
had to carry GitHub's 401 and 404 bodies back out for `internal/cli` to word its
refusals from, or take the refusal wording with it — the first is `AppRequest`
with extra steps, the second moves protocol text out of the package that owns
the refusal table (SPEC §10)."

The [fourth](https://github.com/radiusred/gh-codecrew/issues/326#issuecomment-5570609448)
makes the fake an ordinary package rather than a test file, for a reason that is
purely Go's: "a `fake_test.go` inside `internal/tracker` is invisible to every
other package's tests, which is where the CLI verb tests live." Its contract is
stated in the same record — every call recorded with its arguments in order,
every answer from a `<Method>Fn` field, "nil meaning 'the zero value, no
error', so a test scripts only what it exercises and a verb never panics half
way through on a corner it did not think about" — and its own tests walk the
interface by reflection "so an interface that grows fails with a count rather
than a type error."

The [fifth](https://github.com/radiusred/gh-codecrew/issues/326#issuecomment-5570615700)
is the one the requirement asked for by name: which ad-hoc fakes the new one
replaces, and why the others stay. Two are replaced, where the venue is now
reached through the ctx. The rest stay because "Every one of them embeds a nil
`tracker.Tracker`, so a verb that calls a method the test did not think about
**panics and names it** — that is the assertion, and `faketracker.Venue`'s
'answer with the zero value' contract deliberately does the opposite." The
record also supplies the evidence that they are not the leak: "Because they
embed the interface, none of them needed a line changed when it grew by nine
methods."

The [sixth](https://github.com/radiusred/gh-codecrew/issues/326#issuecomment-5570615879)
is the guard's one exception: `init.go`'s `U`, the documentation base URL `init`
writes into the scaffold. **Trade-off:** "An exception list is a place for a
second entry to appear, against forbidding CodeCrew from linking its own
documentation." **Rejected:** routing `U` through `tracker.RepoURL` — "nothing
fetches it, it addresses this project's own docs, and it is interpolated into
the links `init` writes into `AGENTS.md` and the role files for a human to
click. Constructing it through the venue would say the opposite of what it is."
The same record explains why the guard reads the AST rather than the text: an
import path passes, and `cli.go`'s help-text prose "a dead github.com link
refuses" is not a URL and passes.

The [seventh](https://github.com/radiusred/gh-codecrew/issues/326#issuecomment-5570616068)
is the one place in the whole milestone where a request to GitHub changed shape:
`AccountType` path-escapes the login for every caller, where `identity new`'s
owner lookup did not. **Rejected:** two methods, or an unescaped parameter the
caller escapes. The argument for it being inside the Goal's rule is measured
rather than asserted — "For every login GitHub will actually issue, the two are
byte-identical, and for anything else the unescaped form was a malformed path,
not a different answer" — and the argument for it belonging to the venue is the
seam's own: "a caller that has to remember to escape a path segment is a caller
that is building the request."

The [eighth](https://github.com/radiusred/gh-codecrew/issues/326#issuecomment-5571351623)
was written at 13:27:50Z, after another lane's merge had put a `github.com`
literal into a file this branch's guard would fail on: `--help`'s last line
composes its `CLI.md` link from `U` rather than earning a second entry in the
exception list. **Trade-off:** "This branch edits one line of code #327 landed
hours earlier, against an exception list that doubles in length on the first PR
after the guard exists." (#327 merged at 13:09:16Z, eighteen minutes and
thirty-four seconds before the Decision, not hours; the substance holds and the
interval does not.)
**Rejected:** adding the literal to `venueExceptions` — "It would have
been the correct kind of exception — a documentation link, not a venue call —
and that is the problem: the same address was already in the tree as `U`,
allowed for the same reason, so the list would have carried the same URL twice
and taught the next reader that adding an entry is the way past the guard."
The record closes on the check that makes it safe: `--help` byte-identical to
`origin/main`'s, "checked against a binary built from `0981557` in a scratch
worktree — stdout and stderr — so #327's drift test reads exactly what it read
before, and the exception list stays a list of one."

### The reference, and the question that went to a human

#327's three Decisions and the milestone's one gate all answer the same
question from different sides: what a reference is allowed to be.

The gate is
[raised](https://github.com/radiusred/gh-codecrew/issues/327#issuecomment-5570463643)
at 12:13:55Z, twelve seconds before the task started and before any commit:
does §6 become a stub, or is it deleted and §7–§13 renumbered? It carries a
rider on the same answer — §6 also holds the exit-code contract and the
machine-channel paragraph, "protocol promises, not verb detail" — and it says
where the plan already leans. The operator's
[resolution](https://github.com/radiusred/gh-codecrew/issues/327#issuecomment-5570638590)
at 12:28:39Z takes (a) and keeps the rider: "§6 becomes a one-paragraph stub
naming `CLI.md`; §7–§13 keep their numbers, and nothing outside the SPEC is
renumbered." **Why:** "the stub delivers both without touching every `§N`
citation in contracts, docs, code comments, the site and fifteen milestone
records — a renumbering that changes what every seat reads is the kind of
surface change M16's Goal rules out. Renumbering can be done in a later major if
the stub ever reads as clutter."

The [first Decision](https://github.com/radiusred/gh-codecrew/issues/327#issuecomment-5570470037)
names the file `CLI.md`. **Trade-off:** it "names its subject … next to a
`SPEC.md` that names the protocol", and "gives the site a stable slug when the
spoke's sync gains its row (`CLI.md` → `cli.md` → `/docs/cli/`)"; the cost is
"that the name says 'CLI' where the protocol elsewhere says 'verbs'".
**Rejected:** `REFERENCE.md` — "it names the genre, not the subject, and a
repository whose root already carries a specification, a changelog and a roadmap
has more than one thing that is a reference" — and `docs/cli.md`, which #323's
operator Decision had already ruled out by putting the file at the root.

The [second](https://github.com/radiusred/gh-codecrew/issues/327#issuecomment-5570623955)
answers the question #323 left open: hand-written and held true by two drift
tests, not generated. **Trade-off:** "What the CLI holds in code is the synopsis
line and the flag set; what the reference has to say beyond that … is not in a
form any generator could reach without inventing a second schema to annotate
every verb with. Generating half a page per verb and hand-writing the other half
would give a file no one could edit confidently." **Rejected:** generating the
whole reference — "the annotations would be the reference, moved into Go and out
of a reader's reach" — and shipping no test at all, "on the grounds that a
reference is documentation — the M15 discussion on #306 is what one unguarded
copy already costs."

The [third](https://github.com/radiusred/gh-codecrew/issues/327#issuecomment-5570888002)
is the one a review forced, and it is the milestone's clearest case of the
Goal's rule biting. The reviewer found that `CLI.md` said `status` writes its
`note:` lines to stderr and the CLI writes them to stdout. The Decision states
the CLI's behaviour in the reference and leaves SPEC §6's contrary sentence
exactly as it is. **Trade-off:** "the cost is that §6 and `CLI.md` disagree on
one sentence until §6 is corrected — but §6's paragraph is one of the two the
gate on this task deliberately kept as protocol, M16's Goal forbids changing
what a verb writes to make prose true, and correcting the SPEC's sentence is a
decision with its own trade-off (fix the prose, or move every note to stderr and
break callers reading them out of stdout)." **Rejected:** "softening the
reference to 'notes are advisory and the stream is unspecified' — a reference
whose whole job is to say what a verb writes cannot decline to name the
stream." The disagreement is
[capture #337](https://github.com/radiusred/gh-codecrew/issues/337), filed by
the implementer seat at 12:49:39Z, thirteen seconds before the Decision that
names it.

### A sentence in the docs, and the tests that were missing under it

#328's two Decisions are both about where a claim belongs.

The [first](https://github.com/radiusred/gh-codecrew/issues/328#issuecomment-5570475713)
puts the `skipping` sentence in `docs/introduction.md`'s refusal section as a
short paragraph rather than in a per-code catalogue, and corrects the
requirement's own premise while doing it: "the task and M16-R3 both say 'the
introduction's refusal-code catalogue', but the introduction has not carried one
since protocol 2.0 — it deliberately defers to SPEC §10's table, because keeping
two lists in step was the drift that table exists to prevent." **Rejected:**
adding detail to SPEC §10's `NO_CHECKS` row, "M16's Goal says §10 keeps its
meaning, the table is one line per code by construction, and the refusal's own
detail line — which the CLI prints and the reader actually meets — now carries
the clause."

The [second](https://github.com/radiusred/gh-codecrew/issues/328#issuecomment-5570475877)
is the engineering call inside the documented workflow: classify changed files
by allowlist, not denylist. **Trade-off:** "a denylist reads more naturally …
but it fails open: a directory nobody thought about — a new `migrations/`, a
`Makefile` moved to `build/` — silently stops running the gate, and the failure
is invisible because the check still reports green." **Rejected:** "a `paths:`
filter on the `on: pull_request` trigger, which is the shape most projects try.
A workflow that never triggers reports nothing at all, so a required check
filtered that way sits pending forever and `task finish` sees exactly the
absence `[skip ci]` produces." The documentation states the rule as "gate the
job, not the workflow".

### A page, a guard, and a network with no route out

#329's three Decisions are about how much a documentation task is allowed to
build, and how it is allowed to know things.

The [first](https://github.com/radiusred/gh-codecrew/issues/329#issuecomment-5570528308)
makes it its own page rather than a section of the introduction. **Trade-off:**
"What the requirement asks for is a per-verb table plus a procedure, and
`docs/introduction.md` says of itself that it 'is the map'; a map that carries a
reference table stops being one." It also names the parallel-lane cost it
avoids: "A separate page also keeps this task's edit to `docs/introduction.md`
down to a single list entry, which matters while #327, #328 and #330 are editing
the same file." **Rejected:** appending the section — "the page is already 173
lines of map and this section is 190 of reference — the tail would have
outweighed the map" — and appending the reading-order entry at the end to avoid
renumbering, because "The list is a reading order; a practical operating page
belongs before the protocol text and the history, not after them."

The [second](https://github.com/radiusred/gh-codecrew/issues/329#issuecomment-5570530108)
ships a Go file in a documentation task and argues for it against the brief:
"the dispatch brief says 'No Go change', and this is a Go file. It changes no
behaviour … so it is inside M16's Goal; and the value of a page whose whole
promise is 'this is what the CLI prints' collapses the first time a detail line
is reworded, which the refusal detail explicitly reserves the right to be."
**Rejected:** shipping the page unguarded, and "quoting the CLI loosely
(paraphrasing what it prints) so that nothing could drift, which would have
removed the reason to read the page."

The [third](https://github.com/radiusred/gh-codecrew/issues/329#issuecomment-5570534093)
chooses a network namespace over a fake `gh` on `PATH`. **Trade-off:** "a shim
proves what the CLI does with a message I wrote, and the requirement is what it
does when the machine is offline. The namespace is the real thing — real `gh`,
real DNS failure — at the cost of one Linux-only step that a reviewer on another
platform reproduces with an unroutable `GH_HOST` instead." **Rejected:**
"trusting the code read alone. Two of the findings below are not visible from
the call graph" — and the record then lists them: `init` and `migrate` complete
offline, and `init` without its branch-protection probe takes the cautious path
and commits on `codecrew-bootstrap` rather than the default branch, which it
does not do online in an unprotected repository.

### Where a convention lives, and the link that was not added

#330's two Decisions are both about blast radius.

The [first](https://github.com/radiusred/gh-codecrew/issues/330#issuecomment-5570450362)
puts the conventions section in `docs/introduction.md` rather than
`docs/extensions.md`. **Rejected:** `extensions.md`, which "is scoped strictly
to `.codecrew/roles/<role>.local.md` … How a project handles issues filed by
people outside the crew is not a role extension, and putting it there would turn
a single-subject page into a two-subject one, with the new section unreachable
from the page's own title."

The [second](https://github.com/radiusred/gh-codecrew/issues/330#issuecomment-5570450552)
declines a cross-link the dispatch brief had allowed, and the reason is the
sharpest statement in the milestone of what "nothing an operator or a seat sees
changes" is worth: "`assets.go` embeds the five contracts from those exact
paths, so `roles diff` and `status` in this repository would stay clean — the
local file and the embedded copy are the same bytes. The cost is downstream:
`contractDrift` compares a project's scaffolded `.codecrew/roles/<role>.md`
against the copy embedded in the binary it is running, so a link added to
`coordinator.md` here ships in the next release and makes `status` report
contract drift on every project that scaffolded its coordinator contract from
v2.0.1, for a change none of them made." **Trade-off:** "a coordinator reading
only its contract does not learn from it that a community report is adopted
rather than relabelled." The cross-link runs one way instead.

### The nav label, and the version the tests do not run

radiusred/codecrew-www#30's two Decisions are both about a bug that was
invisible to the suite that should have caught it.

The [first](https://github.com/radiusred/codecrew-www/issues/30#issuecomment-5570453788)
makes the blog nav label the plain post title, with no date in any form.
**Trade-off:** "The nav key *is* the page title as far as the theme is
concerned — that is the whole bug — so any date in the key is a date in
`<title>` and `<h1>` under the theme version that ships. Encoding it as markup
is what we are removing; encoding it as plain text (`Title (2026-09-07)`) leaks
the same date to the same four places, just without the tags." The cost is
priced: the archive lists every post by date, the Atom feed carries `<updated>`,
and the post's own meta block prints the date in full. **Rejected:** teaching
the theme to take the page title from front matter — "it fixes the symptom one
template at a time — the nav sidebar would still render
`<small class="muted">(2026-09-07)</small>` beside the entry."

The [second](https://github.com/radiusred/codecrew-www/issues/30#issuecomment-5570456428)
is the one that makes the test worth having. The requirement asks for an
assertion that `<title>` contains no `<small`; the Decision asserts that the
built page carries no `<small` anywhere, because "Built against `uv.lock`'s
zensical 0.0.58 the post's `<title>` and `<h1>` are already the plain
front-matter title and the nav key only reaches the sidebar; built against
0.0.59 — what `.github/workflows/site.yml`'s unpinned `pip install zensical`
fetches, and what is serving codecrew.works right now — the nav key wins and the
markup lands in `<title>`, `<h1>`, `og:title` and `twitter:title`." So "the
assertion the requirement names literally … passes today on the pinned version
and would have gone green with the bug still in `main.py`." **Rejected:**
pinning the deploy or bumping the lock — "Both are the right idea and neither
belongs in this task: M16's goal binds R6 to 'nothing else on the page changes',
and a theme-version move changes the built page by definition." The split is
[capture radiusred/codecrew-www#32](https://github.com/radiusred/codecrew-www/issues/32),
filed by the operator at 12:21:44Z from the seat's own report.

### One row, and the guard that was not added with it

radiusred/codecrew-www#35's
[one Decision](https://github.com/radiusred/codecrew-www/issues/35#issuecomment-5571332895)
adds the `PAGE_ORDER` row and pins its position, and declines the general guard
that would have caught the drift earlier. **Trade-off:** "A guard would have
caught this drift when the hub merged `docs/working-offline.md` instead of after
it shipped to the live nav in the wrong place, which is the failure this task
exists to correct. Against that: the guard turns every new upstream page into a
red spoke build until the spoke catches up, and the spoke's CI checks the hub's
`main` out on every run and on a daily schedule — so the hub merging a page
would break the spoke's `main`, not just a PR." The softer contract already in
place is named as the reason it is survivable: "a page it does not know is
appended alphabetically rather than dropped, so the page still publishes and the
only cost is its position." **Rejected:** the guard — and the record adds the
cheaper alternative it did not take, "a warning line from the sync naming pages
that fell through to the alphabetical tail", with the reason it is out of scope
and the note "so the reviewer can raise one if they disagree." The reviewer did
not.

### This document's own two

The [first](https://github.com/radiusred/gh-codecrew/issues/340#issuecomment-5571842536)
fixes the counting instant and puts this task's own records inside the totals,
against
[M15's record](15-v2-0-1-what-the-fleet-migration-taught.md) giving both figures
and [M14's](14-adoption-tidy-and-the-v2-0-0-release.md) giving neither.
**Rejected:** counting before the Decisions are written, which "makes the
header's numbers unreproducible for anyone who runs the same query afterwards,
which is the one thing a synthesized count has to survive."

The [second](https://github.com/radiusred/gh-codecrew/issues/340#issuecomment-5571844826)
discharges the refresh obligation by reading `README.md` and
`docs/introduction.md` against what the milestone delivered and editing neither,
because both were made true by the milestone's own merges. It is in the
protocol-discipline section below, with what was checked.

## Deviations

Seven Deviation records: two on the milestone issue, both the operator's and
both about a harness, and five on the task issues — three on #328, one on #326
and one on #329, all five by the implementer seat. This document's own task
carries none.

**[Three test moves the plan did not name](https://github.com/radiusred/gh-codecrew/issues/326#issuecomment-5570681911).**
Four unit tests left `internal/cli` for `internal/tracker/venue_test.go` when
the URL grammar did. **Why:** "a test cannot follow a function only half way.
Leaving delegating one-line wrappers in `internal/cli` purely so their tests
could stay put would have left two names for each piece of grammar, which is the
shape this task exists to remove." The record is specific about what did not
change — "The assertions themselves are unchanged, case for case" — and about
what improved: `internal/tracker/venue_test.go` "adds `RepoURL` and
`AppInstallURL`, which had no test before." The same comment records a second
widening in the safe direction: "the leak guard walks the whole module outside
the venue packages rather than the four directories the plan listed, so a
package added later is covered without anyone remembering to add it."

**[The test the task asked to update did not exist](https://github.com/radiusred/gh-codecrew/issues/328#issuecomment-5570475559).**
#328's plan said to update the expected string for the `NO_CHECKS` refusal
detail; no test asserted it, and nothing pinned the `skipping` bucket as
satisfying either. **Why:** "the sentence the task adds to the docs is a claim
about behaviour — that `skipping` satisfies the gate — and it was undefended by
any test. Adding the two is the same size of change as editing an assertion
would have been, and it stops the docs and the classification drifting apart."
The record also names the check that makes the new test worth anything:
"Confirmed the first test fails against the old detail string before the
change."

**[The skip marker in a commit message, discovered on this branch](https://github.com/radiusred/gh-codecrew/issues/328#issuecomment-5570642284).**
This is the milestone's one deviation that is also a finding, and its evidence
is a chain of timestamps: "PR #332 opened at 12:15:20Z on a base that was still
clean (#331 did not merge until 12:22:09Z) and GitHub created no workflow run at
all for head `06fda46`, whose commit message body quoted the marker while
explaining why it cannot pass the gate. An amend to `33b4bd6`, marker still in
the body, produced no run either; a close/reopen produced no run. The push that
removed the token from both messages, rebased onto the new main, ran CI within a
minute." **Why:** "The evidence is exactly the claim the section makes, so it
belongs in the section — and it is load-bearing here: `task finish`
rebase-merges, so a commit message carrying the token would carry the skip onto
`main`." Both commit messages were reworded to say "a skip-ci commit marker"
instead of the literal token, which is why neither
[1ac14b4](https://github.com/radiusred/gh-codecrew/commit/1ac14b4) nor
[663b822](https://github.com/radiusred/gh-codecrew/commit/663b822) carries one.

**[The changelog conflict, resolved the way the contract says](https://github.com/radiusred/gh-codecrew/issues/328#issuecomment-5570642447).**
Rebased onto `main` after #330 merged, keeping both `## [Unreleased]` sections.
**Why:** "The two entries were written against the same insertion point. The
brief's instruction for a changelog conflict is to keep both." It is the first
of four passes over that hunk in this milestone; the review-rounds section
below has the rest.

**[The guard covers five texts, not two — and then six](https://github.com/radiusred/gh-codecrew/issues/329#issuecomment-5570892831).**
#329's plan and its own Decision comment both said the page quotes two texts
that live in the source. **Why:** "checky counted them on PR #333 and was right.
The argument the Decision made for guarding two applies unchanged to the other
three, and the third is load-bearing — the whole 'what step 6 does to a branch
you already have' section turns on `locally: git fetch && git switch <branch>`
being what the verb prints, since following that line is the step that silently
does nothing." The Deviation is explicit about the record it corrects rather
than editing it: "The earlier comments stand as they were written; this one
corrects them." It also carries two more corrections from the same review: the
`--dry-run` rows claimed the flag for six verbs where three define it, and the
branch-protection paragraph gained the two offline cases that do not reach the
bootstrap branch. By the merge the count had moved once more — the guard covers
six texts, `migrate --dry-run`'s footer having joined them in
[b384b7d](https://github.com/radiusred/gh-codecrew/commit/b384b7d) — and that
last move has no comment of its own; the reviewer's round-three approval names
the gap.

**[The reviewer seat's harness ran out](https://github.com/radiusred/gh-codecrew/issues/325#issuecomment-5571465834).**
The operator's Deviation at 13:37:24Z: "the reviewer seat's declared harness
(codex) could not complete the rebase check on PR #336 — its usage limit was
reached mid-review (resets 18:16 BST) and nothing was posted. The check runs in
a Claude Code session under the same identity, radiusred-checky, as M14 and M15
did when the harness failed." **Why:** "the identity gates the merge, the
harness does not". The record's own list of the rounds that stand is discussed
in the review-rounds section below, because its count and its list disagree.

**[The qa seat's harness, the same limit](https://github.com/radiusred/gh-codecrew/issues/325#issuecomment-5571563800).**
Eight minutes and seventeen seconds later, at 13:45:41Z: "the qa seat's declared
harness (codex) is at its usage limit until 18:16 BST; the M16 verification runs
in a Claude Code session under the qa identity, radiusred-testy, as the reviewer
seat's last round did." **Why:** "the identity gates the verdict, the harness
does not; waiting nine hours to close a milestone whose last merge landed at
13:44Z buys nothing the record needs." The QA comment says the same in its own
closing line.

## The gates

**#325's Gates section was left as the scaffold's placeholder** — the template
line "_What "done" means beyond CI: e2e suites, manual UAT, sign-offs._",
unedited. [M15's record](15-v2-0-1-what-the-fleet-migration-taught.md#the-gates)
says the same of #305,
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md#the-gates) of #269,
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#the-gates)
of #254, [M12's](12-v1-2-0-and-the-field-fixes-behind-it.md#the-gates) of #241
and [M11's](11-housekeeping.md#the-gates) of #233.

**One gate was raised, through the verb the protocol provides for it, and it was
answered in fourteen minutes and forty-four seconds.** #327's plan declared its
ask-the-human point in the body — "Stub or renumber (raised, blocking)" — and
`gh codecrew checkpoint 327` posted it as a `**Gate raised:**` comment at
12:13:55Z with the `cc:needs-decision` label; the operator's
`**Gate resolved:**` came at 12:28:39Z. The task's first commit
([9adb872](https://github.com/radiusred/gh-codecrew/commit/9adb872)) is authored
at 12:26:49Z, before the resolution, which the plan accounts for: the work below
the gate "is drafted on (a)", the answer the operator's instinct on #323 already
pointed at, and the SPEC edit
([1993aba](https://github.com/radiusred/gh-codecrew/commit/1993aba), 12:27:02Z)
carries only what (a) allows. The reviewer checked exactly that ordering: "the
stub-vs-renumber gate was raised with `checkpoint` before the commits and
resolved by the operator before the SPEC edit."

**That is the verb being used for what it is for.**
[M13's record](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#protocol-discipline-observations)
reports one question reaching the operator without it,
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md#the-gates) three, and
[M15's](15-v2-0-1-what-the-fleet-migration-taught.md#protocol-discipline-observations)
one, and M15's closes on the count: "Five questions across three milestones have
now taken this route and none has used the verb the protocol provides for it."
This one used it, and the `cc:needs-decision` label it applies is what put the
task in `status`'s `gates raised:` block while the question was open — which is
where PR #336's evidence block, captured at 12:31Z, shows it.

**The other seven delivery tasks wrote a resolved Ask-the-human section, and all
seven resolved it to none.** Each names the calls the seat made and says where
they went rather than writing the bare word: #326's "Every choice above is an
ordinary engineering call, recorded as a Decision comment on this issue";
#328's, which points at the judgment #191 had already made and the milestone had
adopted; #329's and #330's, each naming the one judgment the requirement
delegates to the task and saying it is recorded as a Decision;
radiusred/codecrew-www#30's the same of the nav date; www#33's, which says the
change "and its position were both settled by the milestone requirement and by
the hub PR #335's follow-up paragraph"; and www#35's, "The position is fixed by
the hub's own reading order."

What gated the work: CI on every pull request — `Go build and test` and `Lint
commit messages`, both required in this hub — an independent approval on each of
the eight, then the QA verdicts. No requirement in this milestone added a gate
of its own, and no release gate exists because there is no release.

## The review rounds

Eight pull requests, **twelve review submissions**, every one by the reviewer
role holder (`radiusred-checky[bot]`), an App distinct from the authoring seat.
**Four change requests**, each resolved in the next round. **Eight approvals
were submitted and all eight stand** — no approval in this milestone was
dismissed, where
[M15's](15-v2-0-1-what-the-fleet-migration-taught.md#the-review-rounds) reports
three dismissals from five lanes and
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md#the-review-rounds) two from
four. The reason is visible in the timestamps rather than in any rule: in every
lane that was force-pushed, the last push preceded the approval. Ten
force-pushes landed on four hub branches — three each on #332, #333 and #335,
one on #336 — and none on the three spoke branches or on #331.

- **PR [#331](https://github.com/radiusred/gh-codecrew/pull/331)** — approved
  [first round](https://github.com/radiusred/gh-codecrew/pull/331#pullrequestreview-5131915305)
  at 12:20:52Z, the milestone's first. It checks a documentation claim against
  the platform it describes — "Markdown issue-template `labels` front matter and
  issue-form `labels` automatically add labels to issues created from those
  templates/forms, while GitHub's repository role table gives Read users no
  apply/dismiss-label permission and Triage+ that permission" — and then runs
  the verbs the section names: "`task start 31` refuses `NOT_A_TASK` for the
  unlabelled adopted report, and `task finish 330 --dry-run` refuses
  `NOT_OWNER`". Its pull request is the only one in the milestone whose
  diff contains no source file at all — two files, `docs/introduction.md` and
  `CHANGELOG.md`.
- **PR [radiusred/codecrew-www#31](https://github.com/radiusred/codecrew-www/pull/31)**
  — approved
  [first round](https://github.com/radiusred/codecrew-www/pull/31#pullrequestreview-5131983754)
  at 12:28:08Z. It reproduces the Decision's own version argument rather than
  taking it: "With the old label format restored in a scratch review worktree,
  the new focused guards fail under both zensical 0.0.58 and zensical 0.0.59",
  and it diffs the built HTML against `origin/main` under both — "under 0.0.58
  only the two sidebar label occurrences change; under 0.0.59 the six hunks are
  exactly the PR body's title, social title, sidebar, and `h1` label changes".
  It also records what it deliberately did not ask for: "`uv.lock` and
  `.github/workflows/site.yml` are untouched; the version split is recorded but
  correctly left out of this PR."
- **PR [#332](https://github.com/radiusred/gh-codecrew/pull/332)** — approved
  [first round](https://github.com/radiusred/gh-codecrew/pull/332#pullrequestreview-5132061933)
  at 12:36:20Z. It verifies the documented workflow against Actions semantics —
  "workflow-level skips/path filters produce no usable reported fact for this
  gate, while a job skipped by `if:` is still reported by the workflow run" —
  and it verifies the Deviation rather than reading it: "The skip-marker
  deviation is supported by the old head commits and Actions run history, and
  the current commit messages contain no literal skip marker."
- **PR [#335](https://github.com/radiusred/gh-codecrew/pull/335)** — two rounds.
  The [change request](https://github.com/radiusred/gh-codecrew/pull/335#pullrequestreview-5132142372)
  at 12:45:56Z is one finding, and it is the one that produced the reference's
  third Decision and capture #337: "`CLI.md:237` says `status` writes 'The
  report goes to stdout, `note:` lines to stderr.' That is not the CLI's current
  behavior", with the four `status.go` line numbers where the notes are written
  to the stdout writer and a run of the built binary to confirm it. Its ask is
  precise about which side must move: "The reference is required to specify each
  verb's output/writes against the existing CLI surface, and M16 forbids
  changing that surface to make the prose true." The
  [round-two approval](https://github.com/radiusred/gh-codecrew/pull/335#pullrequestreview-5132338389)
  at 13:08:22Z re-measures the streams three ways — `status 2>/dev/null`,
  `status >/dev/null`, a refusing `milestone close --dry-run` and an
  `identity token` with no credentials — and checks two more verbs' Writes lines
  that the finding had not touched.
- **PR [radiusred/codecrew-www#34](https://github.com/radiusred/codecrew-www/pull/34)**
  — approved
  [first round](https://github.com/radiusred/codecrew-www/pull/34#pullrequestreview-5132451282)
  at 13:21:33Z. It reads the built site rather than the Markdown: "In
  `site/docs/cli/index.html`, the CLI page is top-level in the Docs nav and
  immediately precedes the SPEC … The built article has 177 hrefs, 89 distinct;
  all local links and fragments resolve, and the only absolute href in the
  article is `https://github.com/radiusred/gh-codecrew/blob/main/LICENSE`."
- **PR [#333](https://github.com/radiusred/gh-codecrew/pull/333)** — three
  rounds, and the only lane in the milestone that took more than two. The
  [first change request](https://github.com/radiusred/gh-codecrew/pull/333#pullrequestreview-5132120821)
  at 12:43:31Z reproduces the whole page in its own network namespace before it
  says anything — "from this hub clone, a scaffolded spoke, a fresh repository,
  a repository with no `origin`, and two protocol 1.x repositories" — and finds
  two things: the `--dry-run` rows promised a flag three of six verbs do not
  have ("This is the one claim on the page that a reader can act on and be wrong
  about"), and the guard covered two of the five source-owned strings the page
  quotes, with the three unguarded ones named by file and line. It also clears
  the Go file the brief had ruled out, on the merits: "It is inside M16's Goal
  and it is in the plan. Not a finding."
  The [second change request](https://github.com/radiusred/gh-codecrew/pull/333#pullrequestreview-5132281689)
  at 13:02:35Z re-verifies every item of round one by running it — "I reworded
  each of the five in `docs/working-offline.md` in turn and watched its own test
  — and only its own test — fail, then restored the file" — and then finds one
  new thing in the corrected text: `migrate` takes `--dry-run` too, so the
  replacement sentence's count of three is wrong, and "`migrate --dry-run` is
  the only dry run that completes with no network — which is exactly the sort of
  thing this page exists to say". It prices the error in the reader's terms: "an
  operator about to migrate offline is told no preview exists and runs the real
  thing — the CLI's one destructive local operation". The
  [round-three approval](https://github.com/radiusred/gh-codecrew/pull/333#pullrequestreview-5132438508)
  at 13:20:01Z executes the fix rather than reading it, in a scratch 1.x repo
  under `unshare -rn`, both with typed identities and with a bare login, and
  reports six nits including two the record itself needs: that the Deviation
  correcting "two texts" now says five where the guard covers six, and that
  "round two's correction … has no comment of its own".
- **PR [radiusred/codecrew-www#36](https://github.com/radiusred/codecrew-www/pull/36)**
  — approved
  [first round](https://github.com/radiusred/codecrew-www/pull/36#pullrequestreview-5132536636)
  at 13:31:51Z, against a hub checkout at `0981557` because the page it places
  had merged eight minutes earlier. It reads the built navigation and its
  neighbours: "`CodeCrew CLI reference` -> `Working offline` -> `CodeCrew
  Protocol Specification` as top-level docs entries, with Working offline's
  previous/next links pointing back to CLI and forward to SPEC." It also
  endorses the Decision that declined the general guard: "the no-general-guard
  choice is recorded as a Decision and is reasonable for this narrow spoke fix".
- **PR [#336](https://github.com/radiusred/gh-codecrew/pull/336)** — two rounds,
  the second a rebase check, and the only change request in the milestone that
  is not about the diff's own subject. Round one's
  [change request](https://github.com/radiusred/gh-codecrew/pull/336#pullrequestreview-5132211201)
  at 12:54:18Z opens "Requesting changes for one rebase/record issue" and the
  issue is one changelog hunk: "In a direct comparison with current main, this
  branch drops that whole `#328` changelog entry. The review brief explicitly
  asks that every changelog entry from the docs PRs that landed before this
  survives." The substance was already clear, and the round says so by listing
  what it checked — the branch binary compared byte-for-byte against the
  installed v2.0.1 for `status`, `milestone close 16 --dry-run`, every `--help`,
  `milestone evidence 15`, `roles show implementer` and `help`, and the guard
  driven by putting an `internal/gh` import and a `gh.Run` back into
  `internal/cli` and watching it name the file.
  The [round-two approval](https://github.com/radiusred/gh-codecrew/pull/336#pullrequestreview-5132632354)
  at 13:43:28Z is the milestone's most mechanical round and the one the harness
  Deviation covers. It compares the pre- and post-rebase patches at the blob
  level and reports every seam file byte-identical, names the single non-rebase
  change (`cli.go`'s help line composing from `U`) and proves it forced — "I put
  #327's exact literal into a throwaway file on a scratch worktree and the guard
  fired on it — `venueExceptions` is an exact-match map" — then proves it
  invisible by building `main` at `0981557` in a scratch worktree and diffing
  `--help` against a build of `b6165ae`: "**Identical, sha256 `d6f018b1…`,
  `cmp` clean.**" It also builds and tests each of the six commits in turn, "so
  the guard never lands on a tree it would fail."

**What the four change requests were about.** #333's two rounds are the only
ones that found something wrong with what a pull request set out to do: a
documented claim a reader could act on and be wrong about, twice, both about
`--dry-run`. #335's was a claim about which stream a verb writes to, which the
Goal's rule then made a documentation change rather than a code one. And #336's
was a changelog hunk. Every approval in this milestone was reached with the
reviewer running something: binaries built from the branch and from the
installed release compared byte for byte, guards fed deliberate leaks and
watched to name the file and line, drift tests broken from both sides, pages
reworded quotation by quotation to watch the right test fail, a site built
strictly under two theme versions, and a verb executed inside a network
namespace with no route out.

**One changelog hunk, four passes, and one of them a change request.** The
`## [Unreleased]` section of `CHANGELOG.md` is where four of the five hub lanes
met. #328's Deviation records the first reconciliation, after #330 merged;
#335's pull request body records the second, "rebased onto `main` at 0ab7586,
keeping both `## [Unreleased]` entries with this one on top"; #333's round-three
approval records the third by checking that the diff against `main` "deletes no
line"; and #336's round-one change request is the fourth, the one time the
reconciliation was got wrong and a review caught it. The implementer contract's
instruction for that hunk is one sentence — keep both entries — and this
milestone exercised it four times across the eighty-two minutes between its
first merge and its last.
[M15's record](15-v2-0-1-what-the-fleet-migration-taught.md#the-review-rounds)
reports the same hunk reconciled four times in eighty minutes across five lanes.

## QA: one round, six verdicts

The qa role holder (`radiusred-testy[bot]`) verdicted every requirement in
[one comment](https://github.com/radiusred/gh-codecrew/issues/325#issuecomment-5571749294)
at 14:01:10Z, sixteen minutes and thirty-five seconds after the last delivery
merged. All six came back `satisfied`, so there is no supersession to read and
no remedy loop to follow — the standing verdict for each requirement is its only
verdict.
[M15's record](15-v2-0-1-what-the-fleet-migration-taught.md#qa-one-round-seven-verdicts)
describes one round and seven verdicts;
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md#qa-two-rounds-and-the-requirement-that-was-answered-rather-than-remedied)
two rounds; and
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#qa-three-rounds-and-the-requirement-that-took-all-three)
three rounds and a remedy task.

**The comment verifies the milestone's rule before it verifies any
requirement**, under its own heading, "The milestone's own rule, checked before
the requirements" — which follows from a Goal whose first claim is a negative
one. What it ran: `--help` from the built binary "differs from installed
v2.0.1 by exactly the two lines naming `CLI.md` and nothing else"; `status`,
`milestone close 16 --dry-run` and `milestone evidence 15` are "**byte-identical**
between the two against this hub (`cmp` clean, all three)"; no flag definition
changed since v2.0.1, checked as "`git diff v2.0.1 HEAD` over every
`fs.Bool/String/Var` line is empty"; the refusal-code vocabulary is unchanged,
"the only new `[A-Z_]` literal in the whole package is `"refused[NO_CHECKS]"`
inside a new test's expectation"; `.codecrew/config.yml` and `internal/config`
untouched, so "the pointer gains no key"; and `SPEC.md`'s diff against v2.0.1 is
"ten hunks and nothing else: §6 replaced by its stub, §8's one sentence, and
eight citation repoints from `(§6)` to `[CLI.md](CLI.md)` with the paragraphs
reflowed around them; §§1–13 keep their numbers."

Each verdict then names what it ran for that requirement, and each carries a
"What I tried that failed to break it" paragraph. Four are worth quoting for
what they establish beyond the tests.

**M16-R1's verdict counts the fake against the interface, method by method.**
"I counted 43 interface methods against 43 `<Method>Fn` fields and 43 `*Venue`
methods … and checked by script that every one of the 43 both consults its `Fn`
and records the call, so the zero-value contract has no silent hole". It then
attacks the guard from four directions in a scratch copy — a `gh.Run` in
`internal/cli`, a raw `exec.Command("gh", …)` and a `github.com` literal in
`internal/config`, and a concatenated host in `cmd/codecrew/main.go` — and
reports each named by file, line and reason. It also states the guard's two
blind spots without calling either a defect: it is literal-based, so
`exec.Command(ghBin, …)` and a host split across concatenations pass, and it
skips `_test.go` files by design.

**M16-R1's verdict is also the one that corrects a Decision's arithmetic.** "the
record's Decision on kept fakes says 'eleven other per-test fakes stay' and then
names twelve; the tree has exactly those twelve … so every kept fake is recorded
as the requirement asks and only the count in the prose is off by one." The
requirement asked for each kept fake to be recorded; each is.

**M16-R2's verdict reads the reference against the code by hand, because the
drift tests do not.** It names the boundary itself: "Neither drift test compares
`CLI.md` against the Go source for options, defaults, reads/writes or exit
status — that half is prose held by review, which is why I read three verbs
against the flag sets by hand." What it found on those three: `task finish`
documents "all fourteen refusal codes — exactly the fourteen `refuse("…")` calls
in its code path, no more and no fewer", and `identity new` documents all eight
flags "including `--no-route`, which the synopsis does not carry (matching
`--help`, and the subject of capture #334)". It also hunts the requirement's
"no narrative" clause as a text property — "no 'because', no
'historically/previously/used to', no `(#N)` issue reference, no milestone or
requirement ID, no first person" — and settles the stream disagreement in the
reference's favour by running it: "`status 2>/dev/null` prints the report **and**
its `note:` line, `status >/dev/null` prints nothing at all".

**M16-R4's verdict reproduces the page rather than reading it, and then reports
the one case it does not cover.** Everything the page claims was re-run under
`unshare -rn` — "more than the four asked for" — including `init` committing on
`codecrew-bootstrap` and `migrate --dry-run` completing offline with the working
tree untouched. The exception is stated and then bounded: "in a repo with **no
git remote at all**, an offline verb fails with a bare `gh repo: no git remotes
found` and no refusal code rather than `GH_UNREACHABLE` — the page's 'what
waits' table describes an offline machine with a normal clone, which is the case
it was written for, and adding a remote to my scratch repo produced the
documented refusal on every one of the six verbs. Not a finding; noted so the
next reader knows it was checked."

The other two are verdicted the same way. M16-R3's checks the documented
workflow against Actions semantics and then breaks both sides of the behaviour —
"Removing `"skipping"` from the bucket switch fails
`TestPRInfoSkippingChecksSatisfyTheGate` … Deleting the new clause from the
refusal format string fails `TestPlanFinishNoChecksNamesTheSkippingPath` on all
three tokens" — and confirms the gate's position and exit status are unchanged.
M16-R5's drives `task start`'s refusal through the real verb with a fake venue
"rather than reading it": "a task-shaped issue with its labels stripped refuses
`NOT_A_TASK`, and the fake recorded zero comments, zero assignments and zero
branches — the refusal is genuinely before the writes, which is what 'invisible
to every verb' has to mean." M16-R6's fetches the live post, reports "**zero**
occurrences of `<small` anywhere in the served HTML", and reproduces the
version-dependence the Decision argued from by putting the markup back in both
places and rebuilding.

**The QA round filed a capture rather than a finding, and says why.**
[#339](https://github.com/radiusred/gh-codecrew/issues/339) — three shipped CLI
strings still cite `SPEC §6` for text this milestone moved to `CLI.md` — is
"Out of M16's scope by the Goal's own rule — #326's Decision keeps every printed
string byte-identical on purpose, and the byte-identical comparison against
v2.0.1 above depends on it — so it is a capture and not a finding against a
requirement." The comment also names the withdrawn requirement in its own
closing lines, "Nothing here is silent", and records the harness it ran under.

## Requirement outcomes

The status column is the **standing** verdict word as the qa role holder wrote
it, verbatim and unqualified. Each requirement was verdicted once, so the
standing verdict is the only one.

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M16-R1 — the venue seam: every invocation of gh and every GitHub REST or GraphQL call sits behind one interface in one package; a test asserts internal/cli invokes no gh, imports no internal/gh and builds no github.com API path; the GitHub implementation is the only one and nothing names a second; one fake venue in the test tree implements the whole interface and the CLI tests use it where the ad-hoc per-test fakes it replaces are not simpler (each kept fake recorded); no verb changes its arguments, output, writes or refusal codes and SPEC.md is untouched; docs/founding-decisions.md names the seam (adopts [#194](https://github.com/radiusred/gh-codecrew/issues/194)) | [#326](https://github.com/radiusred/gh-codecrew/issues/326) / PR [#336](https://github.com/radiusred/gh-codecrew/pull/336) | `satisfied` | Task closed; verdicted by counting 43 interface methods against 43 fake methods and fields, by attacking the guard from four directions in a scratch copy, and by comparing the branch binary with the installed v2.0.1 byte for byte. Two review rounds, the first over a changelog hunk the rebase had dropped and not over the seam. Eight Decisions and one Deviation; the verdict corrects one of them, which names twelve kept fakes and counts eleven |
| M16-R2 — the CLI reference replaces SPEC §6: a root file beside SPEC.md (its name a recorded Decision), pure technical detail and no narrative, per verb: synopsis, purpose, every option with argument, default and effect, what it reads and writes, --dry-run behaviour, the refusal codes it can exit with each linked to §10, exit status, an example; §6's verb table and gate sentences move into it and §6 is left as a stub or the sections renumbered (the task's first ask-the-human point); named by --help's last line and the README's opening lines, placed in docs/introduction.md's reading order before the SPEC, top-level in the codecrew.works navigation; two drift tests (adopts [#323](https://github.com/radiusred/gh-codecrew/issues/323)) | [#327](https://github.com/radiusred/gh-codecrew/issues/327) / PR [#335](https://github.com/radiusred/gh-codecrew/pull/335) and [radiusred/codecrew-www#33](https://github.com/radiusred/codecrew-www/issues/33) / PR [radiusred/codecrew-www#34](https://github.com/radiusred/codecrew-www/pull/34) | `satisfied` | Both halves closed; verdicted by reading three verbs' sections against their flag sets by hand, breaking both drift tests in both directions, and fetching the live `/docs/cli/` page. The milestone's one gate — stub or renumber — was raised with `checkpoint` and resolved by the operator in fourteen minutes; the rider keeping the exit-code and channel paragraphs in §6 is honoured. Three Decisions, the third forced by a review and captured as [#337](https://github.com/radiusred/gh-codecrew/issues/337) |
| M16-R3 — docs-only pull requests: docs/first-milestone.md §5 gains the self-gating workflow shape, a cheap job that always reports and a heavy job behind needs/if on a changed-files test computed by the committed workflow, never from the commit message; docs/introduction.md's refusal section and the NO_CHECKS refusal's detail text say that a check reporting skipping satisfies the gate and [skip ci] produces no fact and is refused (the code, the gate and the exit status unchanged); SPEC §8's absence-never-satisfies paragraph gains one sentence (adopts [#191](https://github.com/radiusred/gh-codecrew/issues/191)) | [#328](https://github.com/radiusred/gh-codecrew/issues/328) / PR [#332](https://github.com/radiusred/gh-codecrew/pull/332) | `satisfied` | Task closed; verdicted against Actions semantics, by deleting the `skipping` bucket and the new refusal clause in turn and watching the two new tests fail, and by confirming the allowlist fails closed on an empty diff, a new top-level directory and a moved base. Approved first round. Two Decisions and three Deviations, one of them the skip-marker finding the branch made on itself |
| M16-R4 — working offline: a documented section that says which verbs run without GitHub, which wait and how each refuses when it cannot reach GitHub, and the honest recipe for work begun offline as the CLI behaves TODAY — verified against the code, not the wish; task start's reconciliation of a pre-existing branch is [#324](https://github.com/radiusred/gh-codecrew/issues/324) and not part of this (adopts [#268](https://github.com/radiusred/gh-codecrew/issues/268)) | [#329](https://github.com/radiusred/gh-codecrew/issues/329) / PR [#333](https://github.com/radiusred/gh-codecrew/pull/333) and [radiusred/codecrew-www#35](https://github.com/radiusred/codecrew-www/issues/35) / PR [radiusred/codecrew-www#36](https://github.com/radiusred/codecrew-www/pull/36) | `satisfied` | Both halves closed; verdicted by re-running every claim under `unshare -rn` from a hub, a spoke, a fresh repository and a 1.x one, and by mutating each of the six guarded quotations. Three review rounds, the only lane past two; both change requests were about `--dry-run` claims a reader could act on. Three Decisions and one Deviation, the Deviation correcting the guard's own count |
| M16-R5 — external contributors: a documented conventions section stating that no issue template may auto-apply a cc: label, that a community report becomes work by task new --adopts <report> — never by relabelling or rewriting the report — so task finish closes it with the task, and that a labeled-event workflow reverting cc: labels applied below maintain is deferred until a project needs it; no SPEC text changes, no CLI change (adopts [#31](https://github.com/radiusred/gh-codecrew/issues/31)) | [#330](https://github.com/radiusred/gh-codecrew/issues/330) / PR [#331](https://github.com/radiusred/gh-codecrew/pull/331) | `satisfied` | Task closed; verdicted against GitHub's own template and permission behaviour and by driving `task start`'s `NOT_A_TASK` refusal through the real verb with a fake venue, which recorded zero writes. Approved first round, the milestone's first, on the only diff in it with no source file. Two Decisions, the second declining a cross-link that would have reported contract drift on every project scaffolded from v2.0.1 |
| M16-R6 — the codecrew.works blog: the nav label main.py generates carries no HTML — the page <title> and <h1> are the plain post title and the date is shown by the theme's meta block alone — with a test that renders one post and asserts <title> contains no <small; the strict build stays clean; nothing else on the page changes (adopts [radiusred/codecrew-www#29](https://github.com/radiusred/codecrew-www/issues/29)) | [radiusred/codecrew-www#30](https://github.com/radiusred/codecrew-www/issues/30) / PR [radiusred/codecrew-www#31](https://github.com/radiusred/codecrew-www/pull/31) | `satisfied` | Task closed; verdicted on the live post returning 200 with a plain `<title>` and zero `<small` anywhere, and by restoring the markup in both `main.py` and `zensical.toml` to watch each guard fail. Approved first round. Two Decisions; the second widens the assertion the requirement names literally, because that one would have passed with the bug still in place under the pinned theme version. The version split is [radiusred/codecrew-www#32](https://github.com/radiusred/codecrew-www/issues/32) and stays open |

**Captures adopted and closed by the merges — six**, and none of them by a
`Closes` line in a pull request body:
[#31](https://github.com/radiusred/gh-codecrew/issues/31) by `task finish 330`,
[radiusred/codecrew-www#29](https://github.com/radiusred/codecrew-www/issues/29)
by `task finish 30`,
[#191](https://github.com/radiusred/gh-codecrew/issues/191) by
`task finish 328`,
[#323](https://github.com/radiusred/gh-codecrew/issues/323) by
`task finish 327`,
[#268](https://github.com/radiusred/gh-codecrew/issues/268) by
`task finish 329`, and
[#194](https://github.com/radiusred/gh-codecrew/issues/194) by
`task finish 326`. Every one was closed one or two seconds after the task that
adopted it, which is `task finish` doing both in one run. The two spoke tasks that
carry no `## Adopts` section — www#33 and www#35 — closed nothing, because they
are the spoke halves of hub requirements rather than adoptions of their own.

**Captures filed during the milestone and left open — six**, five in this hub
and one in the spoke, and every one of them traceable to a specific moment in
the work. [#324](https://github.com/radiusred/gh-codecrew/issues/324) was filed
by the operator two seconds before the milestone issue itself, splitting
#268's behaviour half out so M16-R4 could adopt the documentation half without
touching a verb: "Not in M16 by the operator's rule: no behaviour change while
2.0 beds in."
[radiusred/codecrew-www#32](https://github.com/radiusred/codecrew-www/issues/32)
— CI tests build with the locked zensical while the deploy installs it unpinned
— comes from the implementer's report on the spoke's first pull request, and is
explicit that "A test cannot guard this; the workflow shape is the guard."
[#334](https://github.com/radiusred/gh-codecrew/issues/334) — `--help` omits
`identity new --no-route`, and SPEC §10 omits `identity webhook` from the
no-pointer verbs — is what writing every option down found, and it is filed
rather than fixed because "M16 forbids changing a verb's surface and the
synopsis drift test pins the reference to `--help` as it stands."
[#337](https://github.com/radiusred/gh-codecrew/issues/337) is the stream
disagreement, filed by the implementer seat from the reviewer's finding and
carrying the two ways out for whoever takes it up.
[#338](https://github.com/radiusred/gh-codecrew/issues/338) is discussed in the
observations below. [#339](https://github.com/radiusred/gh-codecrew/issues/339)
is the QA seat's, filed at 13:59:33Z, ninety-seven seconds before the verdicts.

## Protocol-discipline observations

Nine things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **A verb refused a requirement into existence.** M16's seventh requirement
  did not survive its own first `task new`, because `--adopts` reads the
  capture's state and `#36` was closed as not planned two days earlier. Nothing
  else in the protocol reads what a requirement names: `milestone new` writes
  the issue and stops, `milestone evidence` walks the links a body already
  carries, and `status` counts tasks. The check that caught a scoping error was
  a side effect of a flag that exists to keep captures and tasks linked. It cost
  four minutes and one Decision, and it is the strongest argument in this
  milestone for adopting by verb rather than by prose.
- **A milestone body was edited after creation for the first time on this hub,
  and the edit is a Decision's own subject.**
  [M15's record](15-v2-0-1-what-the-fleet-migration-taught.md) reports #305's
  `userContentEdits` count as zero and says "the requirements a reader sees are
  the requirements that were written". #325's count is two — the body as created
  and one edit — and the operator's Decision names the edit as its own trade-off
  before any task started. The protocol has no verb that amends a requirement
  and no rule about when a body may be edited; what stood in for one here was a
  Decision written two seconds after the edit, and the `userContentEdits` query
  that makes the difference recoverable by anyone.
- **The gate verb was used, and the label it applies did the reporting.** #327
  raised its ask-the-human point with `gh codecrew checkpoint` before its first
  commit, and the `cc:needs-decision` label put the task in `status`'s
  `gates raised:` block for the fourteen minutes the question was open — which
  is visible in PR #336's evidence block, captured while it was. Five questions
  across
  [M13](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#protocol-discipline-observations),
  [M14](14-adoption-tidy-and-the-v2-0-0-release.md#the-gates) and
  [M15](15-v2-0-1-what-the-fleet-migration-taught.md#protocol-discipline-observations)
  reached the operator without it; this one did not.
- **The rule "nothing an operator or a seat sees changes" was measured four
  separate times, by three different parties, and once it forced a design.**
  The implementer measured it on #336 (`status`, `milestone close --dry-run`,
  `milestone evidence 15`, `roles show` and `help` byte-identical against the
  installed v2.0.1) and again on #326's eighth Decision (`--help` byte-identical
  against a build of `main`); the reviewer measured it on #336's two rounds; and
  QA measured it across the whole milestone before verdicting anything. The one
  time it forced a design was #327's third Decision: a reference could not be
  made true by changing what a verb writes, so the reference states the CLI and
  the SPEC keeps its sentence, with a capture holding the difference. A rule
  that is only asserted costs nothing; this one cost a Decision, a review round
  and a capture.
- **A documentation task shipped Go, twice, and both times a Decision argued for
  it against its own brief.** #329's guard rebuilds six CLI strings from the
  code that prints them; #328 wrote two tests that did not exist under a claim
  it was about to document. Neither changes behaviour, and both were checked
  against the Goal by the reviewer rather than waved through — "It is inside
  M16's Goal and it is in the plan. Not a finding." The pattern the repository
  now has three instances of (M15's #299, and these two) is that a documentation
  page which quotes the CLI ships with a test that rebuilds the quotation.
- **The parallel-lane cost fell entirely on one changelog hunk, and it was paid
  in rebases rather than in dismissed approvals.** Ten force-pushes across four
  hub branches, four reconciliations of `## [Unreleased]`, and one change
  request for getting one of them wrong — but no approval was dismissed, because
  every lane's last push preceded its approval.
  [M15's record](15-v2-0-1-what-the-fleet-migration-taught.md#protocol-discipline-observations)
  reports three dismissals from five lanes and
  [M14's](14-adoption-tidy-and-the-v2-0-0-release.md#protocol-discipline-observations)
  two from four; four lanes here produced none. Nothing in the protocol
  sequences pull requests, and nothing here needed it to.
- **A capture was answered by a merge fourteen minutes before it was filed.**
  [#338](https://github.com/radiusred/gh-codecrew/issues/338) asks for a clause
  naming `migrate` as the fourth `--dry-run` verb in `docs/introduction.md`,
  quoting the sentence "`milestone new`, `task finish` and `milestone close`
  take `--dry-run`". That sentence is in the file at
  [00b254b](https://github.com/radiusred/gh-codecrew/commit/00b254b) and is not
  in it at [67256a3](https://github.com/radiusred/gh-codecrew/commit/67256a3),
  PR #335's merge at 13:09:16Z; the capture was filed at 13:23:31Z from a review
  note written at 13:02:35Z, when the sentence was still there. Two lanes editing
  the same page cannot see each other's merges from inside a review, and the
  protocol offers nothing that would have caught it — `milestone evidence` checks
  that citations resolve, not that they are still true. The capture is left open
  and this record states the finding; correcting the page here would have
  invented a defect to fix.
- **Two seats ran on a harness their routing table does not name, under the
  identity it does.** The reviewer's round two on #336 and the entire QA round
  ran in Claude Code sessions because codex was at its usage limit, each
  recorded as an operator Deviation before the work. The argument is the same
  one [M14's](14-adoption-tidy-and-the-v2-0-0-release.md#the-gates) and
  [M15's](15-v2-0-1-what-the-fleet-migration-taught.md#deviations) records give
  — the identity gates, the harness routes — and this is the third consecutive
  milestone to need it. What is new is that it reached the qa seat as well as
  the reviewer's, so both non-doer gates in one milestone were held by a harness
  the table does not name.
- **The refresh obligation found nothing to fix, because the milestone's own
  tasks had already done it.** `README.md`'s opening block and its "Read next"
  list carry `CLI.md`, added by #335; its reference list carries
  `docs/working-offline.md`, added by #333 in answer to a non-blocking review
  nit. `docs/introduction.md`'s reading order runs 6 `CLI.md`, 7 Working
  offline, 8 `SPEC.md`, and its "What exists" prose now defers per-verb detail to
  the reference rather than enumerating it. `**Shipped:** v2.0.1` is still true —
  M16 cut no tag — and SPEC §10 still catalogues forty-three refusal codes. That
  is a different outcome from
  [M15's](15-v2-0-1-what-the-fleet-migration-taught.md#protocol-discipline-observations),
  where the obligation "found two incomplete sentences and nothing false", and
  the reason is structural rather than lucky: four of M16's five hub
  requirements name `docs/introduction.md` or `README.md` in their own text, so
  the pages were refreshed by the work instead of after it.

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **The conversation that scoped the milestone is on the record only through
  #325's goal.** "Decided with the operator 2026-09-07" is the whole of it. Each
  of the six adopted captures carries its own history, but the judgement that
  these six belong together — and that #36 belonged with them until a verb
  said otherwise — has no transcript. *The reasoning is recorded; the
  deliberation is not*, as
  [M15's record](15-v2-0-1-what-the-fleet-migration-taught.md#what-the-record-does-not-contain),
  [M14's](14-adoption-tidy-and-the-v2-0-0-release.md#what-the-record-does-not-contain)
  and [M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#what-the-record-does-not-contain)
  each say of their own openings.
- **Which harness ran eleven of the twelve review rounds is not recorded, and
  the one Deviation that attributes any of them miscounts its own list.** The
  [reviewer Deviation](https://github.com/radiusred/gh-codecrew/issues/325#issuecomment-5571465834)
  says "the six earlier M16 rounds by codex" and then names eight — "#331, #332,
  #335 r1 and r2, www#31, www#34, www#36 and #336 r1" — and #333's three rounds
  appear in neither the count nor the list. Eleven submissions preceded #336's
  round two. The identity is constant across all twelve, so the gate the
  protocol enforces is intact; the harness behind each round is not
  reconstructable, exactly as
  [M15's](15-v2-0-1-what-the-fleet-migration-taught.md#what-the-record-does-not-contain)
  and [M14's](14-adoption-tidy-and-the-v2-0-0-release.md#what-the-record-does-not-contain)
  records report of their own.
- **The sixth guarded quotation on #329 has no record of its own.** The
  [Deviation](https://github.com/radiusred/gh-codecrew/issues/329#issuecomment-5570892831)
  corrects "two texts" to five and the merged guard covers six, `migrate
  --dry-run`'s footer having been added in the round-two fix. The reviewer's
  round-three approval names the gap — "round two's correction … has no comment
  of its own" — and the task closed without one. The commits
  ([8d0a6f4](https://github.com/radiusred/gh-codecrew/commit/8d0a6f4),
  [b384b7d](https://github.com/radiusred/gh-codecrew/commit/b384b7d)) and the
  pull request body, which says six and names them, are where the sixth is
  recorded.
- **Nothing on the trail says what the venue seam is for beyond the tests.** The
  Goal rules out a second implementation — "GitHub the only implementation, no
  second venue and nothing that names one" — and says "the gain now is a fake
  venue that makes the tests simpler and honest". #194's own GitLab notes "were a
  feasibility check and stay notes". `docs/founding-decisions.md` records that
  another venue waits on community demand. Whether the seam is worth its
  forty-three methods to anything other than the test suite is not argued
  anywhere, and cannot be until something asks for a second venue.
- **The two blind spots QA found in the leak guard are recorded only in the QA
  comment.** `exec.Command(ghBin, …)` behind a variable, a `CommandContext`
  form, and a host name split across concatenated literals all pass; so does a
  leak introduced in a `_test.go` file, which the walk skips by design. No
  capture was filed for either, and the guard's own source does not say so. QA
  called them "Two assumptions the suite makes, neither a blocker", which is the
  whole of the record.
- **Nothing measures whether SPEC §6's stub reads as clutter.** The gate
  resolution says "Renumbering can be done in a later major if the stub ever
  reads as clutter", which puts a judgment in the future with nothing that will
  prompt anyone to make it. There is no capture, no roadmap row and no test.
- **Whether the stale branches that motivated M14-R4 still exist is still not
  known, and nothing in this milestone looked.**
  [M15's record](15-v2-0-1-what-the-fleet-migration-taught.md#what-the-record-does-not-contain)
  reports the same, for the same reason: `milestone close`'s candidate set is
  the hub plus the repositories its own tasks name, and this milestone's tasks
  name this hub and the spoke. Nothing here changed either mechanism.
- **Nothing measures what the milestone cost.** Eight pull requests, twelve
  reviews and one QA round, with no record of tokens, wall-clock per seat, or
  how many dispatches were made.
  [M15's record](15-v2-0-1-what-the-fleet-migration-taught.md#what-the-record-does-not-contain),
  [M14's](14-adoption-tidy-and-the-v2-0-0-release.md#what-the-record-does-not-contain)
  and [M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#what-the-record-does-not-contain)
  say the same of their own, and the answer has not changed. What the trail does
  carry is the wall-clock of the work itself — one hour, thirty-eight minutes
  and forty-one seconds from the milestone opening to the last merge — which is
  a fact about the lane rather than about its cost.
