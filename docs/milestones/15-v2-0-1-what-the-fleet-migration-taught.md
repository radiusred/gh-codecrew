# M15: v2.0.1: what the fleet migration taught

Tracking issue: [#305](https://github.com/radiusred/gh-codecrew/issues/305) ·
Synthesized 2026-09-07 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #305's **seven** requirements as opened, **none
added, none struck and none amended** — the issue's `userContentEdits`
query returns a total count of zero, so the requirements a reader sees are
the requirements that were written; its Gates section; its **four
comments**, carrying **two Deviation records, the operator's verification of
the release, and one QA comment holding all seven verdicts**; the **seven**
delivery task issues, six in this repository
([#306](https://github.com/radiusred/gh-codecrew/issues/306),
[#307](https://github.com/radiusred/gh-codecrew/issues/307),
[#308](https://github.com/radiusred/gh-codecrew/issues/308),
[#309](https://github.com/radiusred/gh-codecrew/issues/309),
[#310](https://github.com/radiusred/gh-codecrew/issues/310),
[#311](https://github.com/radiusred/gh-codecrew/issues/311)) and one in the
spoke ([radiusred/codecrew-www#27](https://github.com/radiusred/codecrew-www/issues/27)),
plus [#321](https://github.com/radiusred/gh-codecrew/issues/321), the task
this document is; their **seven merged pull requests** and the **twenty-four
commits** on them; the **twenty-nine Decision records and eight Deviation
records** across the milestone issue and all eight of its task issues, of
which **twenty-six and seven** are on the milestone issue and the seven
delivery tasks and the rest on this document's own; the **fourteen
review submissions** on the seven pull requests — **four change requests**,
**ten approvals submitted of which seven stand and three were dismissed by a
rebase**; the **eight** adopted backlog captures the merges closed; and the
**four** captures the milestone filed and left open — three in this hub and
one in the spoke. The house form is
[M14's document](14-adoption-tidy-and-the-v2-0-0-release.md), with
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md) and
[M12's](12-v1-2-0-and-the-field-fixes-behind-it.md) behind it; the record
standard is the one the reviewer set on
[PR #304](https://github.com/radiusred/gh-codecrew/pull/304), whose two
rounds found first a check the M14 record's task claimed and its branch did
not carry, and then two counts in the correction itself that were five
short — so every count above and below is re-derived from the source, records
by the rule `tracker.ExtractRecords` itself applies (a paragraph-initial
`**Decision:**`, `**Deviation:**` or `**Gate resolved:**` label, bare or
parenthetically qualified), review states from each pull request's reviews
API and timeline, and the release from the git, releases and actions APIs
rather than from the evidence comment that reports it.

**Counted at 2026-09-07T02:26:50Z, against a trail that is frozen from
there.** Two records landed on the milestone after this document's first
commit and before its review: the operator's
[Decision](https://github.com/radiusred/gh-codecrew/issues/306#issuecomment-5564002686)
on #306 at 02:09:01Z, and a second
[Deviation](https://github.com/radiusred/gh-codecrew/issues/305#issuecomment-5564026934)
on #305 at 02:11:36Z, about this document's own review round. The counts
above are re-derived as of that instant and include this document's own
task, whose three Decisions and whose
[Deviation](https://github.com/radiusred/gh-codecrew/issues/321#issuecomment-5564196823)
recording the recount are the last records the milestone takes — the
coordination layer posts nothing further on it before this pull request
merges.
[M14's record](14-adoption-tidy-and-the-v2-0-0-release.md) left its own
task's records out of its totals; this one gives both figures rather than
choose one in silence.

**The trail this record is compiled from is checked by the release this
milestone cut.** `gh codecrew milestone evidence` walks the milestone issue
and its sub-issues — bodies and comments, which is where every citation
below comes from — and refuses when a github.com citation does not resolve.
The installed extension reports `v2.0.1 (protocol 2.0)`, and run from this
repository at the same instant as the counts above it printed:

```
$ gh codecrew milestone evidence 15
requirements counted: M15-R1, M15-R2, M15-R3, M15-R4, M15-R5, M15-R6, M15-R7 (7)
all 20 cited links resolve across 9 issues — evidence is reachable
```

Nine issues: the milestone and its eight sub-issues — the seven delivery
tasks, and this document's own.

The gather from `gh codecrew milestone close 15` is again not part of the raw
material, for the reason
[M10's](10-protocol-bookkeeping-from-the-field.md),
[M11's](11-housekeeping.md),
[M12's](12-v1-2-0-and-the-field-fixes-behind-it.md),
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md) and
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md) records all give:
`--dry-run` stops at the gate that counts tasks, naming the task that writes
this file. It is quoted as it read before this document's own pull request
opened; with that pull request open the same gate names #321 `(in review)`.

```
gate milestone open: ok
gate no gate raised: ok
gate tasks closed: refused[OPEN_TASKS]: tasks not closed: radiusred/gh-codecrew#321 (ready)
gate requirements declared: not reached
gate QA verdicts: not reached
gate milestone document: not reached
dry run: nothing written — the live verb stops at the first refusal above
```

Every record below was read from its issue, pull request, review, commit,
tag, release, workflow run or timeline directly. The prose is wrapped at this
file's normal width; the requirement-outcomes table's rows and a handful of
single links are not, because a newline ends a Markdown table row and
breaking a link breaks it — the same shape
[M14's record](14-adoption-tidy-and-the-v2-0-0-release.md#requirement-outcomes)
and [M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#requirement-outcomes)
carry, and the shape `ROADMAP.md` has had since M1.

**This PR adds the M15 ROADMAP row; it does not flip one** — the convention
M10-R1 introduced and M12-R1 shipped, which
[M14's record PR](https://github.com/radiusred/gh-codecrew/pull/304),
[M13's](https://github.com/radiusred/gh-codecrew/pull/290),
[M12's](https://github.com/radiusred/gh-codecrew/pull/252) and
[M11's](https://github.com/radiusred/gh-codecrew/pull/240) were written
under. `ROADMAP.md` carried no M15 row before this PR, so there was nothing
to discard and the trail carries no Decision discarding one — the same as
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md) and
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md).

## Goal and outcome

#305 was opened at 23:51:44Z on 2026-09-06, sixteen minutes and forty-eight
seconds after
[`milestone close 14`](https://github.com/radiusred/gh-codecrew/issues/269#issuecomment-5563049997)
closed M14 with the line "all 7 tasks done, milestone document merged". Its
goal is one long sentence of provenance followed by one of scope, and the
provenance is the milestone's whole argument: "v2.0.0 shipped on 2026-09-06
and the operator migrated the fleet the same evening; running the new verbs
on their own milestone and moving ten repos surfaced seven captures, all of
them implementation corrections and none a protocol change."

Then it names the seven, in the words of the captures themselves: `migrate`
leaves a spoke without the root entry points `init` would have written
([#301](https://github.com/radiusred/gh-codecrew/issues/301)); `task start`
reports a 403 on every App-held seat and `status` never names who holds an
App-run task ([#287](https://github.com/radiusred/gh-codecrew/issues/287));
every listing read stops at the first hundred, so a long milestone issue
loses its newest verdicts
([#264](https://github.com/radiusred/gh-codecrew/issues/264)); a body saved
with CRLF defeats the tracker's line-anchored scans
([#296](https://github.com/radiusred/gh-codecrew/issues/296)) and nothing
asserts that `task new` applies `cc:task`
([#297](https://github.com/radiusred/gh-codecrew/issues/297)); a skipped
stale-branch sweep is invisible until the next close
([#295](https://github.com/radiusred/gh-codecrew/issues/295)) and pull
request prose can hand GitHub an unintended closing reference
([#303](https://github.com/radiusred/gh-codecrew/issues/303)). The scope
follows from that list: "This milestone fixes them, ships v2.0.1, and then —
with released behaviour to describe in the present tense — publishes the
site's first blog post on protocol 2.0 and the migration."

All seven captures were filed on 2026-09-06, between 11:24Z and 22:30Z,
while M13 and M14 were running — the last two of them after v2.0.0 shipped,
#301 from the operator's own migration run and #303 from the QA seat's
verification of it. Every one is answered by a requirement that names it,
and the mapping is
one-to-one but for two requirements that take a pair: M15-R1 adopts #301,
M15-R2 adopts #287, M15-R3 adopts #264, M15-R4 adopts #296 and #297, M15-R5
adopts #295 and #303, M15-R6 adopts nothing — it is the release — and M15-R7
adopts [radiusred/codecrew-www#26](https://github.com/radiusred/codecrew-www/issues/26),
the blog post the site has had a section for and no posts in.

The seven delivery tasks, in merge order:

- **[#307](https://github.com/radiusred/gh-codecrew/issues/307) / PR
  [#313](https://github.com/radiusred/gh-codecrew/pull/313)** — M15-R2.
  Three commits, approved first round, merged 00:11:45Z
  ([8f9d9e9](https://github.com/radiusred/gh-codecrew/commit/8f9d9e9)).
  Three Decisions. Adopts
  [#287](https://github.com/radiusred/gh-codecrew/issues/287).
- **[#308](https://github.com/radiusred/gh-codecrew/issues/308) / PR
  [#312](https://github.com/radiusred/gh-codecrew/pull/312)** — M15-R3. Two
  commits, two review rounds, the second a rebase check, merged 00:20:46Z
  ([55a04ee](https://github.com/radiusred/gh-codecrew/commit/55a04ee)).
  Three Decisions. Adopts
  [#264](https://github.com/radiusred/gh-codecrew/issues/264).
- **[#306](https://github.com/radiusred/gh-codecrew/issues/306) / PR
  [#315](https://github.com/radiusred/gh-codecrew/pull/315)** — M15-R1. Four
  commits, two review rounds, merged 00:29:53Z
  ([9e225d9](https://github.com/radiusred/gh-codecrew/commit/9e225d9)). Five
  Decisions, the fifth the operator's and recorded after the merge. Adopts
  [#301](https://github.com/radiusred/gh-codecrew/issues/301).
- **[#309](https://github.com/radiusred/gh-codecrew/issues/309) / PR
  [#314](https://github.com/radiusred/gh-codecrew/pull/314)** — M15-R4. Six
  commits, three review rounds, the third a rebase check, merged 00:40:31Z
  ([7610047](https://github.com/radiusred/gh-codecrew/commit/7610047)). Four
  Decisions and three Deviations, one of which records no deviation at all.
  Adopts [#296](https://github.com/radiusred/gh-codecrew/issues/296) and
  [#297](https://github.com/radiusred/gh-codecrew/issues/297).
- **[#310](https://github.com/radiusred/gh-codecrew/issues/310) / PR
  [#317](https://github.com/radiusred/gh-codecrew/pull/317)** — M15-R5. Six
  commits, three review rounds, the third a rebase check, merged 00:50:13Z
  ([9c322e0](https://github.com/radiusred/gh-codecrew/commit/9c322e0)). Four
  Decisions and one Deviation, and two of the milestone's open captures come
  out of its review. Adopts
  [#295](https://github.com/radiusred/gh-codecrew/issues/295) and
  [#303](https://github.com/radiusred/gh-codecrew/issues/303).
- **[#311](https://github.com/radiusred/gh-codecrew/issues/311) / PR
  [#320](https://github.com/radiusred/gh-codecrew/pull/320)** — M15-R6, Part
  A. Two commits, two review rounds, merged 01:08:18Z
  ([2dc147a](https://github.com/radiusred/gh-codecrew/commit/2dc147a)). Four
  Decisions, one of them the operator's, and one Deviation. Part B — the tag,
  the release run and the assets — and the operator's own verification ran on
  the closed issue.
- **[radiusred/codecrew-www#27](https://github.com/radiusred/codecrew-www/issues/27)
  / PR [radiusred/codecrew-www#28](https://github.com/radiusred/codecrew-www/pull/28)**
  — M15-R7, and the milestone's only task outside this hub. One commit,
  approved first round, merged 01:25:44Z
  ([9931d40](https://github.com/radiusred/codecrew-www/commit/9931d40)).
  Three Decisions. Adopts
  [radiusred/codecrew-www#26](https://github.com/radiusred/codecrew-www/issues/26).

From the milestone opening to the last delivery merging: **one hour,
thirty-four minutes**, to the second. From the first `task start` — #308's,
at 23:59:09Z — to that merge: one hour, twenty-six minutes and thirty-five
seconds; to the last merge in this hub, #320's, one hour, nine minutes and
nine seconds.

Five of the six hub pull requests were open at once. #312 opened at
00:03:50Z, #313 at 00:04:24Z, #314 at 00:05:53Z, #315 at 00:05:56Z and #317
at 00:09:02Z — five minutes and twelve seconds across all five — and they
merged in the order #313, #312, #315, #314, #317, each later branch rebased
onto the merge before it. That parallelism is where all three dismissed
approvals and all three rebase-check rounds come from, and it is the section
below on the review rounds.

Five fixes, then a release, then a post: **seven `satisfied` verdicts, all
reached in one QA round**, where
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#qa-three-rounds-and-the-requirement-that-took-all-three)
took three rounds and
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md#qa-two-rounds-and-the-requirement-that-was-answered-rather-than-remedied)
two.

What exists afterwards that did not before. `migrate` writes an absent root
`AGENTS.md` or `CLAUDE.md` from the same constants `init` writes them from,
in the same pathspec commit, listed by `--dry-run` — so a migrated spoke is
indistinguishable from a fresh 2.0 scaffold, and the `action needed` block
has one condition left, a file the project already had that reaches nowhere.
`task start` does not offer a GitHub App as an assignee and says nothing
about it; `status` names the holder of a task in progress or in review from
its latest `**Started by**` record. Four listing reads in the tracker walk
the whole listing, and the two that deliberately do not each say so in a code
comment, one of them held there by a test asserting it does not paginate. Every issue body and comment body the tracker
returns has its CRLF normalised, and so does every exported scanner's input;
SPEC §4 says line endings are not part of the record grammar. `status`
reports the stale task branches of the repository it runs in with the verdict
the next `milestone close` would give them, computed by that close's own
function; `task finish` names any issue its merge would close besides the
task, and both role contracts now carry the rule that produced it. v2.0.1 is
released and installed. And [codecrew.works](https://codecrew.works) has its
first blog post. SPEC §10's table still counts **forty-three** refusal codes:
this milestone added none, which is what "implementation corrections, none a
protocol change" means when it is measured rather than asserted.

## The release, and the blog post

M15-R6 is the one requirement whose delivery is split between a seat and the
operator, and M15-R7 is the one that waits on it. The trail follows both
splits exactly.

**Part A**, PR [#320](https://github.com/radiusred/gh-codecrew/pull/320),
merged 01:08:18Z at
[2dc147a](https://github.com/radiusred/gh-codecrew/commit/2dc147a):
`## [Unreleased]` became `## [2.0.1] — 2026-09-07` with a fresh
`## [Unreleased]` carrying `Nothing yet.` above it, seven entries moving down
unchanged; the compare links updated; `docs/introduction.md` and `README.md`
flipped to v2.0.1. Where 2.0.0's section opens with a
`### What broke, and what to do` block, this one opens with a single
paragraph, and the pull request body says why: "what a reader of a `.1` needs
first is not a break but the absence of one." The paragraph names where the
release came from — ten repositories migrated onto the `.codecrew/` layout,
then the new verbs worked on a milestone of their own — and ends on the two
facts an adopter is looking for: the protocol stays at 2.0, and a repository
already on the 2.0 layout has no migration step to take.

**Part B**, run by the implementer seat on the closed issue and evidenced in
one [comment](https://github.com/radiusred/gh-codecrew/issues/311#issuecomment-5563624254)
at 01:10:40Z. Every claim in it was re-checked for this record against the
API: the [annotated tag `v2.0.1`](https://github.com/radiusred/gh-codecrew/releases/tag/v2.0.1)
is tag object `fabdc17eba781a98af792e5542d730fab0285be9`, pointing at commit
`2dc147ac7fa5539e7f1b171f72aa8d984840b2a2`, message `v2.0.1 (protocol 2.0)`,
tagger
`radiusred-cody[bot] <270330637+radiusred-cody[bot]@users.noreply.github.com>`
at 01:08:33Z — fifteen seconds after the merge it names; `release.yml`
[run 34072047784](https://github.com/radiusred/gh-codecrew/actions/runs/34072047784)
was created at 01:08:40Z on that same commit and completed with conclusion
`success`; the
[release](https://github.com/radiusred/gh-codecrew/releases/tag/v2.0.1) was
published at 01:09:55Z, is neither draft nor prerelease, is named
`gh-codecrew 2.0.1`, and carries **five** assets — darwin-amd64,
darwin-arm64, linux-amd64, linux-arm64 and windows-amd64.exe. The evidence
comment adds one detail the API does not: the run's only annotation is the
runner's standing Node 20 deprecation notice, "identical to the v2.0.0
run's".

**The operator's legs**, posted on #311 at 01:11:08Z as a
[three-step checklist](https://github.com/radiusred/gh-codecrew/issues/311#issuecomment-5563627340)
and answered on #305 at 01:21:41Z by a
[record](https://github.com/radiusred/gh-codecrew/issues/305#issuecomment-5563697433)
that reads back the four commands it ran. `gh extension list` shows
`gh codecrew radiusred/gh-codecrew v2.0.1`; `gh codecrew version` prints
`v2.0.1 (protocol 2.0)`; `gh codecrew roles diff <role>` for all five roles
reports "matches the embedded v2.0.1 contract"; and `gh codecrew migrate
--dry-run` reports the hub already on the 2.0 layout with its labels at the
defaults and nothing written. The `status` line is the one worth quoting,
because it is two of this milestone's five fixes observed at once in a
single command: "the M15 board with six tasks done and radiusred/codecrew-www#27
in review, naming its holder `@radiusred-wordy[bot]` from the start record
(M15-R2 live); no gates; no contract-drift line; no stale-branch line
(M15-R5's report, silent because clean — the hub carries no `task/` branch
whose task is closed); the delete-on-merge note prints as before."

**There is no migration leg this time, and the record says so in one line**:
"No migration leg this release: the protocol is still 2.0." M14-R3's
equivalent step was a fleet migration — a checklist of twelve repositories
and a migration record of eleven; M15-R6's is one `gh extension upgrade`.
That difference between a protocol major and a patch is what the operator's
one line records, and it is the only place on the trail where it is said.

**The blog post waited for the release, and the trail shows the wait.**
M15-R7's own wording — "merged only after M15-R6" — is the sequencing gate,
and it holds by timestamps rather than by anyone enforcing it: the release
was published at 01:09:55Z, the task
[radiusred/codecrew-www#27](https://github.com/radiusred/codecrew-www/issues/27)
started at 01:14:35Z, its pull request opened at 01:19:34Z and merged at
01:25:44Z. The post is
`docs/blog/posts/2026-09-07-protocol-2-0-and-the-migration.md`, titled
"Protocol 2.0: CodeCrew Moves Out of Your Repository's Root", an unheaded
opening and five sections that follow the adopted capture's shape: why a
protocol major, what changed, what an adopter does, what deliberately did not
change, and what the migration taught, the last of them ending on the seven
captures v2.0.1 answers. Its generated
artefacts — the nav block between `BEGIN_BLOG_POSTS` and `END_BLOG_POSTS` in
`zensical.toml`, `docs/blog/archive.md` and `docs/blog/atom.xml` — ride in
the same single commit, because all three are tracked. The QA seat fetched
the built page and recorded a 200.

The post is the site's first, and the section it lands in was empty before
it: the capture that asked for it opens on the observation that "the site
has a Blog section ... and no posts", and that "protocol 2.0 is the first
change that asks every adopter to do something".

## Decisions

Twenty-nine Decision records: twenty-three on the six task issues in this
hub, three on the spoke's task, and three on this document's own — the
milestone issue itself carries none. Two of the twenty-nine are the
operator's, both on the trail after the work they judge: the linked-branch
Decision on #311, and the
[declined marker](https://github.com/radiusred/gh-codecrew/issues/306#issuecomment-5564002686)
on #306, written after that task had closed. The rest were written by the
seat doing the work — `radiusred-cody[bot]` for the six code and release
tasks, `radiusred-wordy[bot]` for the post and for this document.

One carries a parenthetical qualifier rather than the bare label, and it
matters to the count: `**Decision (operator, after the merge):**` on #306 is
a record under `ExtractRecords`' rule and is invisible to a search for
`**Decision:**` alone, exactly as
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md#decisions) two qualified
labels were. The other twenty-eight are bare. They group by the question each
task had to settle.

### A migrated spoke, and a released changelog section edited to describe it

M15-R1's whole subject is what `migrate` leaves behind, and #306's four
Decisions run from a message string out to the record itself.

The [first](https://github.com/radiusred/gh-codecrew/issues/306#issuecomment-5563218004)
deletes the `(kept)` suffix and, with it, `entryPointAction`'s headline and
list parameters, so `init` and `migrate` print one identical `action needed`
block. **Trade-off:** an operator who has seen the 2.0.0 output sees
different words, and "the two verbs are now coupled through one message
rather than two callers of a shared payload" — against which "the suffix
existed only to distinguish `(kept)` from `(absent)`, and `(absent)` is
gone". **Rejected:** keeping the suffix against a third state arriving,
because there is none: "a root entry point is either the project's file or
one CodeCrew wrote, and the second never needs a human."

The [second](https://github.com/radiusred/gh-codecrew/issues/306#issuecomment-5563218057)
is the one the review then went after, and it is this milestone's clearest
statement of what a dry run is for: "the dry run is judged against the state
it previews, not the state it starts from." `reachesInstructions` takes
the set of files the caller is about to write, so an absent root `AGENTS.md`
beside a kept `CLAUDE.md` holding `@AGENTS.md` is judged against the
`AGENTS.md` the migration is about to create. **Trade-off:** one more
parameter on a helper both verbs call, for one case — but "`migrate`'s own
doc comment promises that `--dry-run` 'is a true preview', so the divergence
is a bug rather than a rough edge." **Rejected:** computing the block after
the writes and staying silent in a dry run: "a dry run that is silent about
the one act it cannot perform for you is worse than one that is wrong about
it."

The [third](https://github.com/radiusred/gh-codecrew/issues/306#issuecomment-5563218136)
edits a released section of `CHANGELOG.md`, which is the milestone's most
deliberate departure from a convention it names. Two passages in the 2.0.0
section describe `migrate` naming an absent root entry point under
`action needed`; both now say what the verb does, each marked
`(2.0.1, #306)`. **Trade-off:** "Keep a Changelog treats a released section
as a record, and this edits one. But that section is also this project's only
migration guide — it is what an operator with a repo still on 1.x reads
before running the verb, and they will run v2.0.1 — so leaving it would hand
every remaining 1.x adopter an instruction to create by hand a file the CLI
now writes." **Rejected:** stating the correction only under `[Unreleased]`:
"the guide is read top-down from the 2.0.0 heading by someone who has just
upgraded; a correction eight hundred lines above it is not where they are
looking." The reviewer judged the departure explicitly rather than passing
over it — "an amendment, not a silent rewrite" — and the QA verdict for R6
records the 2.0.0 section as preserved "apart from recorded #306 correction
notes".

The [fourth](https://github.com/radiusred/gh-codecrew/issues/306#issuecomment-5563218236)
draws the line the third one needed: the milestone records under
`docs/milestones/` are left alone, "the record of what those milestones
decided and shipped, not documentation of the current verb", even though
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md) is the record that names
capture #301 in the first place. **Trade-off:** "a reader who greps `docs/` for `action needed`
finds prose describing a behaviour the CLI no longer has. That is what a
record is."

A [fifth Decision](https://github.com/radiusred/gh-codecrew/issues/306#issuecomment-5564002686),
the operator's, arrived long after the task closed — three minutes after this
document's first commit, which had recorded its absence as a gap. The
reviewer's round-two note on PR #315 had observed that a third passage of the
2.0.0 section, the `**What it writes.**` bullet, states the same behaviour
without a `(2.0.1, #306)` marker, and left it to the author. The Decision
declines it rather than leaving it, on a ground the note had not weighed:
"That passage describes what the 2.0.0 binary wrote, which is true of 2.0.0;
the two marked passages are the adopter's migration steps, which had to
change because the steps a reader follows today are 2.0.1's. The 2.0.1
section's own entry says what changed." **Trade-off:** "a reader of the 2.0.0
section alone sees behaviour the current release has improved on; the marked
steps beside it say so." **Rejected:** marking the third passage, which
"would turn a description of a released version into a running commentary."

### An App is not an assignee, and the board reads the record instead

M15-R2's three Decisions on #307 are all about where a fact should come from.

The [first](https://github.com/radiusred/gh-codecrew/issues/307#issuecomment-5563206012)
makes the skip wider than the requirement's word: an unrouted login carrying
the `[bot]` suffix is treated as an App too, because the predicate is the
existing `crewIdentity`, "true for a `[bot]` suffix and for a login the
routing table resolves to an `app:`-typed row". **Trade-off:** the skip is
wider than "the routing table says `app:`" — and the argument is that the two
sets are the same one: "a `[bot]` login this project's table does not name is
still a GitHub App by construction, and GitHub refuses it as an assignee
exactly as it refuses a routed one, so the call skipped is a call that could
only have produced the 403 #287 is about." **Rejected:** a narrower test on
the routed row, which "would leave an unrouted App seat ... printing the
permanent 403 note that this task exists to remove, and it would put a
second, subtly different definition of 'is an App' next to `crewIdentity`."

The [second](https://github.com/radiusred/gh-codecrew/issues/307#issuecomment-5563206107)
keeps the first assignee as a display fallback in `status`, and answers the
objection before anyone raises it: M13-R7 deleted exactly such a fallback
from `tracker.StartedBy`. "It is not the same thing: `StartedBy` answers 'who
owns this task', a gate `task finish` refuses on, and it still has no
fallback ... `status` answers 'what is GitHub showing about this task', and
there the fallback is forced by the state machine above it: SPEC §4 infers
'in progress' *from* the assignee being set, so a task in that state with no
start record would print a state with no name capable of explaining it."
**Rejected:** dropping it for consistency, which "prints a self-contradictory
line — `[in progress]` with nobody named — for the one case that produced the
state."

The [third](https://github.com/radiusred/gh-codecrew/issues/307#issuecomment-5563206199)
is the milestone's clearest refusal to widen scope, and it is the one that
became a capture. `tracker.InferState` is left alone, so an App-run task that
has started but has no pull request open still reads `[ready]` and never
reaches a state with a holder to name. The Decision says the consequence out
loud — "R2's own words are 'in progress or in review', so the change
satisfies the requirement as written, but the App case is thinner than it
first reads" — and gives two reasons for stopping: teaching the state machine
the start record "is an edit to SPEC §4's lifecycle table ... and so a
protocol change; M15's Goal states that all seven captures it adopts are
implementation corrections", and it "would also make the state machine need a
comments read for every task in every state, where the requirement bounds the
cost at one read per task in two states." The reviewer confirmed the reading
in round one — "Declining to decide it here was right. It reads to me like a
backlog capture for the operator" — and
[#316](https://github.com/radiusred/gh-codecrew/issues/316) is that capture.

### One paging mechanism, and the reads that are not listings

M15-R3's three Decisions on #308 are a scope boundary drawn twice and a test
seam built to prove the change.

The [first](https://github.com/radiusred/gh-codecrew/issues/308#issuecomment-5563208102)
takes `gh api --paginate` over a cursor loop, and the reasoning is a
founding-decision argument: gh "follows the `Link` header and joins a REST
array endpoint's pages into a single JSON array, in request order, so
`gh.JSON` unmarshals the whole listing exactly as it already does for
`Labels`". It records the check rather than the belief — "I checked the join
against this repo's own `labels` and `pulls` endpoints on gh 2.100 before
relying on it, including the `-X GET -f …` form `OpenPRsForBranch` uses" —
and names the cost: "the call is unbounded: a pathological listing is now
several requests instead of one truncated read, which is the correct trade
for a gate that must not miss a verdict." **Rejected:** hand-rolled cursor
loops, "in a package whose founding decision is that gh owns the transport".

The [second](https://github.com/radiusred/gh-codecrew/issues/308#issuecomment-5563208160)
says what M15-R3's "every listing read" does not reach: the GraphQL
`first: N` connections — `ClosingPRs`, `LinkedBranches` and `Task`'s nested
assignees, labels and closing references — "are per-issue *relations*, not
repository listings: an issue with twenty-one closing PRs or eleven assignees
is not a shape the protocol produces, and every one of them is read for a
yes/no ... rather than for a set that must be complete." **Rejected:**
cursor-walking them for symmetry, which "would be untested paging code on
paths no fixture can drive honestly."

The [third](https://github.com/radiusred/gh-codecrew/issues/308#issuecomment-5563208243)
is about the test rather than the change, and it is why the change is
believable. Neither existing `gh` fake could express pagination — "a fixture
of 101 comments served unconditionally would pass with or without
`--paginate`, proving nothing" — so a third one, `pagingGH`, serves the whole
listing only when the call carried the flag and truncates to a hundred
otherwise. The verification is stated as a mutation: "removing `--paginate`
from all four readers fails all four tests with 'returned 100 of …', and
restoring it passes." The same comment ends with a paragraph headed
**Deviation from nothing else:**, recording that `fakeGH`'s `/pulls` matcher
now finds the path anywhere in the argument list — a paragraph the record
grammar does not read as a Deviation, which is discussed below.

### Two layers for one rule, and a claim made good rather than withdrawn

M15-R4 asks for normalisation "once at its boundary", and #309's
[first Decision](https://github.com/radiusred/gh-codecrew/issues/309#issuecomment-5563218322)
does it twice, in the same comment that says why the requirement's word
cannot be met as written: "`Tracker` is an interface, so the reader is the
boundary a real body crosses but not the only way a body reaches a scanner: a
`string` crossing that seam carries no promise about its line endings."
Normalising only in `github.go` leaves a second backend or a fake tracker
"able to feed CRLF straight into a `(?m)…$` scan"; normalising only in the
scanners leaves `Comment.Body` carrying CRLF into `internal/cli`. So both,
"and the second pass is free where there is nothing to replace". **Rejected:**
pushing it down into the unexported helpers, "the same work spread over more
sites ... below the level a reader looks for it", and keeping the ad-hoc
`strings.ReplaceAll` inside `paragraphs` as a third layer.

The same comment carries a second Decision — the table test drives each
scanner *through* `GitHub{}.IssueBody` and `GitHub{}.Comments` rather than
with hand-made strings, so it "proves the path a GitHub body actually takes
instead of the path the test invented" — and one sentence in it turned out to
be false. That sentence is the subject of the milestone's most instructive
exchange, and it is recorded under the Deviations below rather than edited
away.

Two further Decisions land in a
[second comment](https://github.com/radiusred/gh-codecrew/issues/309#issuecomment-5563219531).
SPEC §4 gains one sentence, and the argument for adding it is about who reads
the spec: "a human writing a record does so in GitHub's web editor, which is
exactly the thing that produces CRLF: the reader who needs the guarantee is
the reader §4 is written for, and until now nothing told them whether a
Decision typed in the browser would be found." **Rejected:** describing the
normalisation itself, because "SPEC §4 states what a record *is*, not how the
tracker reads bytes." The other holds `UnresolvedGates` to returning the
`Comment` values it was handed rather than normalised copies: "Scanning is
the scanner's job; rewriting its input is not."

### What silence is allowed to mean, and what an advisory read may cost

M15-R5's four Decisions on #310 are two pairs: what a report may leave
unsaid, and what a note may cost when it fails.

The [first](https://github.com/radiusred/gh-codecrew/issues/310#issuecomment-5563197623)
scopes the stale-branch report to `c.current`, the repository the verb runs
in, and the reasoning is a cost argument rather than a coverage one:
`status` "has no repo set to walk that is both complete and cheap. The hub's
routing table names roles, not repos ... the open milestones' tasks name
repos, but the report has to work between milestones too — that quiet period
is exactly when #295 says a skipped sweep goes unseen, and it is when the
board is empty." **Trade-off:** a spoke accumulating stale branches stays
invisible until someone runs `status` there. **Rejected:** walking the
milestones' repos as the close does — "the same set, at a much higher rate,
with the wrong coverage between milestones."

The [second](https://github.com/radiusred/gh-codecrew/issues/310#issuecomment-5563197714)
prints nothing when there is nothing stale, and defends the silence by naming
what makes it safe: "the unreadable case is covered instead: it prints its
own `note:` line, so silence means clean and a failure is never silent." The
reasoning is about attention rather than about correctness: "a permanent
clean line on every run is the kind of output an operator stops reading,
which is precisely how the line that does matter gets missed." Round one of
the review then found the third unreadable case the Decision had not counted
— a failing `RepoInfo` — and the fix is that a failed read prints a `note:`
where the report would have been, so the property the Decision asserted is
now true of every path.

The [third](https://github.com/radiusred/gh-codecrew/issues/310#issuecomment-5563343962)
exists because the reviewer asked for it to exist. `ClosingReferences` reads
one page, `closingIssuesReferences(first: 50)`, while M15-R3 — open in
another lane at the time — required every listing read in the tracker to
paginate. The Decision draws the boundary: R3 "exists for listings that grow
without bound while the thing they describe stays one thing", whereas "a pull
request's closing references are bounded by the body that declares them and
by nothing else — they are written once, by one author, in one text."
**Trade-off:** a body with more than fifty closing references would truncate
silently, "judged acceptable against the alternative — every existing
closing-reference set on this hub is one to three — and recorded here rather
than only in the Go comment so QA verifying M15-R3 can see the boundary in
the record."

The [fourth](https://github.com/radiusred/gh-codecrew/issues/310#issuecomment-5563344058)
turns a hard error into a note, in answer to a reviewer's explicitly
non-blocking judgment rather than to a request. "The read is advisory by
construction — it feeds one sentence, not a
gate and not an action the merge depends on — and this codebase already fixes
what an advisory failure costs": `status` carries on with a `note:` when a
listing fails, `deleteHead` says a failure there "is a note, never an error",
and SPEC §6 says nothing after the merge refuses. "Aborting a finish whose
every gate has passed, after nothing has been written, over a GraphQL hiccup
on an informational line, is the odd one out — and it is the exact asymmetry
checky named on PR #317, in the same change that has `status` do the
opposite." **Rejected:** consistency with `PRInfo` and `planAdoptions`, on
the ground that "those two are consistent with each other because they are
load-bearing; consistency with them is not a property this read shares."

### What the flip does not touch

Three Decisions on #311 are about restraint, and each carries a precedent.
The [first](https://github.com/radiusred/gh-codecrew/issues/311#issuecomment-5563529562)
adds no changelog entry for the flip itself, "as #186 (v1.1.0), #242 (v1.2.0)
and #272 (v2.0.0) each made the same call": the entry "would be pure
self-reference, and it would have to be written under the very heading it
announces". The
[second](https://github.com/radiusred/gh-codecrew/issues/311#issuecomment-5563529646)
leaves SPEC §5's `v2.0.0 (protocol 2.0)` illustration standing one patch
behind, and is unusual in stating a trade-off in both directions: the stale
CLI half "is a real if small wrongness", against which "flipping it makes
SPEC — the protocol document, versioned by the protocol — a file every patch
release must touch, which is the shape of a claim that rots." It ends by
inviting the overrule: "Cheap to overrule if the reviewer reads the line as a
claim rather than an illustration." The reviewer did not.

The [third](https://github.com/radiusred/gh-codecrew/issues/311#issuecomment-5563529747)
is the one edit in the flip that is not a version bump. The README sent a 1.x
adopter to "the top of `CHANGELOG.md`" for the migration steps, and this pull
request is what makes that false — "from now on the top of that file is
whatever shipped last, while what a 1.x adopter needs stays where it is" — so
the sentence names the 2.0.0 section instead. **Rejected:** a deep link to
the section anchor, because "GitHub's anchor for `## [2.0.0] — 2026-09-06` is
generated from the heading text, date included, so it would break the next
time that heading is touched — a worse rot than the one being fixed."

The fourth Decision on #311 is the operator's, and it answers a Deviation;
both are below.

### The post that had to be true in the present tense

The three Decisions on
[radiusred/codecrew-www#27](https://github.com/radiusred/codecrew-www/issues/27)
are the doc-synthesizer seat's, and all three are about what a public page
may say.

The [first](https://github.com/radiusred/codecrew-www/issues/27#issuecomment-5563678709)
sets a link policy: relative links into the site's own synced `/docs/` pages
for anything the docs sync publishes — the spec and the GSD essay — and
absolute github.com URLs for everything it excludes, "the CHANGELOG, the
milestone records ..., the issues, the issue-comment permalinks and the two
releases. Every one of those was resolved through the API before the post was
committed." **Trade-off:** "the relative links only resolve once
`sync_docs.py` has run, so the post is unreadable as a bare file in the
repository and correct only as a built page ... and the strict build is what
catches it." **Rejected:** absolute GitHub URLs for the two synced pages,
because "sending a reader to GitHub for a page the site publishes is the
drift the sync exists to prevent."

The [second](https://github.com/radiusred/codecrew-www/issues/27#issuecomment-5563678809)
records the single departure from a verbatim quotation: the adopter's five
steps are the CHANGELOG's own, with the bare `#306` at the end of step 5
rendered as a link, "off GitHub it is inert text that names nothing".
**Rejected:** paraphrasing the steps, "which would put a second, drifting
copy of the migration instructions on a marketing site while the CHANGELOG
stays the source of truth."

The [third](https://github.com/radiusred/codecrew-www/issues/27#issuecomment-5563678916)
is the one that could only be made by someone reading the site's own rules:
the fleet migration is described by counts and outcomes and never by
repository name, because "the record behind it names twelve repositories
across two GitHub accounts, and some of them are private", and the site's
README forbids naming private repositories in public copy. **Trade-off:** "a
reader cannot check the fleet from the post alone; they follow the two links
to do it."

## Deviations

Eight Deviation records. Seven are on the milestone issue and the delivery
tasks — two by the operator on #305, three on #309 and one each on #310 and
#311, those five by the implementer seat — and the eighth is on this
document's own task. None is recorded on #306, #307, #308 or the spoke's
task, and the bodies of PR #312 and PR #315 each say "No deviations from the
plan" in as many words.

**[The reviewer seat ran on a different harness again — for a different
reason](https://github.com/radiusred/gh-codecrew/issues/305#issuecomment-5563478569).**
The operator's Deviation at 00:46:06Z: "the reviewer seat's declared harness
(codex) could not run the round-three rebase check on PR #317: two `codex
exec` dispatches in a row were stopped by the operator's machine for low
memory before posting anything (nothing reached GitHub either time)." The
check ran "in a Claude Code session under the same identity,
radiusred-checky, as M14 did for the same reason". **Why:** "the identity is
what gates the merge, not the harness behind it; the machine, not the seat,
failed." The last clause is the difference from
[M14's version](14-adoption-tidy-and-the-v2-0-0-release.md#deviations) and
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#deviations),
which were both about a usage limit: this one is a resource failure on
the operator's own hardware. The Deviation also names which rounds are
unaffected — "#312 r1 and r2, #314 r3" — which is the only harness
attribution anywhere in the milestone's delivery.

**[And again, on this record's own review round](https://github.com/radiusred/gh-codecrew/issues/305#issuecomment-5564026934).**
The operator's second Deviation, at 02:11:36Z, reports the same failure on
PR #322's first round — `codex exec` "stopped by the operator's machine for
low memory before posting anything, this time with 28 GiB free, so the
earlier explanation (page cache read as pressure) is incomplete" — and the
round ran in a Claude Code session under the same identity. **Why:**
"unchanged — the identity gates, the harness does not; a third codex attempt
in the same session state is not worth the round." Two things follow it. The
reviewer's declared harness has now been stood down twice in this milestone
alone, both times for the operator's machine rather than for a quota, and the
first of the two explanations is withdrawn on the record rather than left
standing. And this is a record
about the review of the document you are reading, which is why the counts in
its opening paragraph name the instant they were taken.

**[The two layers are now pinned by tests, and a claim is made good rather
than withdrawn](https://github.com/radiusred/gh-codecrew/issues/309#issuecomment-5563303825).**
The Deviation that round one on PR #314 produced, and the one place in this
milestone where a recorded claim was measured and found false. The Decision
on #309 argued for normalising in two places and said, of its table test, "It also means the
test would still catch a regression if the boundary layer were removed and
only the scanner layer left, and the other way round." The reviewer measured
it: deleting all seven scanner-entry calls left the suite green but for a
test that pre-dated the task, so "six of the seven call sites were held by
nothing". The Deviation quotes the false half back and says what it costs:
"**'The other way round' was false**, for every scanner but `ExtractRecords`,
and it was the argument for the second layer — so it had to become true or
stop being made. It is made good rather than withdrawn; this comment corrects
it on the record and the earlier one is left as written." Each table row
gained a third reading — the CRLF fixture handed straight to the scanner, no
reader in the path — and the comment then states which assertion pins which
layer, measured one layer at a time. It also draws a distinction the test
file's header now carries: three of the eight rows are pins and five are
guards, because "only three of the scans can be defeated by a CR at all" and
the other five "assert the property the second layer makes unconditional".

**[The `UnresolvedGates` row projects the comment's author and URL rather
than its body](https://github.com/radiusred/gh-codecrew/issues/309#issuecomment-5563303825).**
A one-paragraph Deviation in the same comment, and it is the consequence of a
Decision two comments earlier: the scanner deliberately returns the `Comment`
values it was handed, "so a body reaching it directly comes back carrying the
CR it arrived with, and projecting the body would have asserted a rewrite the
scanner deliberately does not do."

**[No deviation from the plan](https://github.com/radiusred/gh-codecrew/issues/309#issuecomment-5563218322).**
The third record on #309 is a paragraph reading "**Deviation:** none from the
plan; this is the plan as written, recorded here as the contract asks." Under
`ExtractRecords`' rule it is a Deviation record, and this record counts it as
one; what it records is the absence of one. It is counted here because the
alternative is a number no reader can reproduce from the source — see the
observations below.

**[The closing-reference note names a qualified
ref](https://github.com/radiusred/gh-codecrew/issues/310#issuecomment-5563248817).**
The task's brief wrote the note's shape as a bare `#<n>`; it ships as
`radiusred/gh-codecrew#42`. **Why:** "GitHub's closing keywords accept
`owner/repo#n`, so a pull request can genuinely carry a closing reference to
another repository's issue, and `ClosingReferences` reads each node's own
`nameWithOwner` rather than assuming the PR's. A bare number would name the
wrong issue in exactly the case the note exists to catch." The reviewer
picked it out in round one as correct on those grounds.

**[The pull request body cannot un-link its own
task](https://github.com/radiusred/gh-codecrew/issues/311#issuecomment-5563538019).**
The milestone's one escalation, and it is the same one M14 had, recognised
immediately because M14 recorded it. PR #320's body "carries no closing
keyword anywhere — every `#ref` in it was checked before opening — but
`gh pr view 320 --json closingIssuesReferences` reports `[311]`". **Why:**
"the link is not body-derived. `gh codecrew task start` creates the working
branch with `gh issue develop`, a *linked branch* on the issue's development
panel, and GitHub converts that link into a closing reference the moment a
pull request is opened from the branch. Editing the body cannot clear it."
The Deviation cites the precedent by permalink — PR #298, the v2.0.0 flip,
"reports `[272]` to this day for the same reason" — and states the check that
is actually available: "exactly one reference, this task's own number, and
nothing else."

The operator's
[Decision](https://github.com/radiusred/gh-codecrew/issues/311#issuecomment-5563544774)
one minute and six seconds later accepts it: "the merge will therefore close
this task before Part B (the tag) and the operator's verification. Accepted:
Part B's evidence and the operator's record post on the closed issue, exactly
as #272 did for v2.0.0." **Trade-off:** "the task's state reads done while the
release is still being cut; the milestone's R6 verdict, not the issue state,
is what says the release shipped." **Rejected:** "unlinking the branch or
re-opening the issue after the merge — churn in the record for no gain the QA
verdict does not already give."

That exchange took a minute where M14's took an investigation. M14's
Deviation on #272 lists what was tried against the API before concluding the
link had no public handle; #311's cites it and moves on. The cost of the
second occurrence of a known protocol wrinkle is one comment and one answer.

**[The counts moved because the trail moved](https://github.com/radiusred/gh-codecrew/issues/321#issuecomment-5564196823).**
The eighth Deviation is this document's own, and it is the reason the numbers
above are what they are. The first commit of this record counted twenty-five
Decisions and six Deviations; the operator's #306 Decision landed three
minutes later and the second #305 Deviation five minutes after that, and the
review found both — "the record's whole claim on a reader is that its counts
are re-derivable from the trail, and two of them stopped being so at
02:09:01Z". **Why:** a count taken at synthesis and left alone is true only
until the next comment, and this milestone's trail kept moving after
synthesis where
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md) did not. **Trade-off:**
the totals now cover this document's own task, where M14's record left its
own out — so both figures are stated. **Rejected:** dating the snapshot and
shipping it, "a first paragraph that a reader re-deriving it on merge day
finds wrong by two, with a note explaining why that is acceptable — in a
document whose subject includes a seat correcting its own recorded claim
rather than leaving it standing".

## The gates

**#305's Gates section was left as the scaffold's placeholder** — the
template line "_What "done" means beyond CI: e2e suites, manual UAT,
sign-offs._", unedited.
[M14's record](14-adoption-tidy-and-the-v2-0-0-release.md#the-gates) says the
same of #269,
[M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#the-gates)
of #254,
[M12's](12-v1-2-0-and-the-field-fixes-behind-it.md#the-gates) of #241 and
[M11's](11-housekeeping.md#the-gates) of #233. No `**Gate raised:**` or
`**Gate resolved:**` record exists anywhere in this milestone, neither
repository carries an issue labelled `cc:needs-decision`, and `checkpoint`
was not used — on a task or on the milestone issue. `gh codecrew status` run
against this hub reports `gates raised: none`.

**All seven delivery task Plans wrote a resolved Ask-the-human section, and
all seven resolved it to none.** #309's is the bare word. The other six name
the calls the seat made and say where they went: #306's "the `(kept)` suffix
and the CHANGELOG 2.0.0 prose are ordinary engineering calls; both are
recorded as Decisions on this issue"; #307's the same of two; #308's of the
paging mechanism; #310's of "what `status` says when it has nothing to
report"; the spoke task's that "the judgment calls this post rests on ...
were all made and recorded by the operator before it was written; this task
compiles them". #311's is the only one that names a condition rather than a
call — "The date in the `## [2.0.1]` heading is 2026-09-07; if the tag slips
past that day it is amended before tagging" — and the evidence comment closes
it: "The commit's date, 2026-09-07 UTC, matches the `## [2.0.1] — 2026-09-07`
heading it ships, so the heading needed no amendment."

That is a milestone in which no question reached a human through a Plan.
[M14's record](14-adoption-tidy-and-the-v2-0-0-release.md#the-gates) reports
three questions reaching the operator, one of them raised in a Plan. This
milestone had one, and it went by Deviation: PR #320's closing reference,
answered by the operator's Decision above. It did not go through `checkpoint`
either.

What gated the work: CI on every pull request — `Go build and test` and
`Lint commit messages`, both required — an independent approval on each of
the seven, then the QA verdicts. M15-R6 added a release gate of its own that
no verb enforces, and M15-R7 a sequencing one: the operator's installation
checks, and a post that may not merge before the release it describes.

## The review rounds

Seven pull requests, **fourteen review submissions**, every one by the
reviewer role holder (`radiusred-checky[bot]`), an App distinct from both
authoring seats. **Four change requests**, each resolved in the next round.
**Ten approvals were submitted; seven stand, and three were dismissed by a
rebase force-push.**

- **PR [#313](https://github.com/radiusred/gh-codecrew/pull/313)** — approved
  [first round](https://github.com/radiusred/gh-codecrew/pull/313#pullrequestreview-5127113957)
  at 00:10:17Z. It checks the behaviour by mutation rather than by test name — "I reverted each behaviour
  in a scratch tree (unwrapped the `Assign` call; made `taskHolder` return
  the assignee first) and re-ran" — and then runs the built binary against
  this hub to see the capture's own symptom gone: "the four in-review M15
  tasks now name `@radiusred-cody[bot]`, which is exactly the line #287 said
  was missing". Its five findings are all non-blocking. Two look past the diff —
  that `status`'s new comments read will walk a long task's whole history
  once #308 paginates `Comments`, and that `crewIdentity` now answers two
  questions that are the same set today and are not the same concept — and a
  third reaches the same reading as capture #316, which the operator had
  already filed four minutes earlier from the implementer's own report.
- **PR [#312](https://github.com/radiusred/gh-codecrew/pull/312)** — two
  rounds, the second a rebase check. Round one
  [approved](https://github.com/radiusred/gh-codecrew/pull/312#pullrequestreview-5127116465)
  at 00:11:18Z — twenty-seven seconds before #313 merged — having verified
  `--paginate` against "a forced multi-page REST array endpoint" rather than
  taking the mechanism on trust, and was dismissed at 00:13:31Z by the rebase
  onto that merge. The
  [round-two approval](https://github.com/radiusred/gh-codecrew/pull/312#pullrequestreview-5127137666)
  compares the patch before and after "excluding `CHANGELOG.md`", finds the
  stable patch IDs match, and names the only content differences as #307's
  mainline files.
- **PR [#315](https://github.com/radiusred/gh-codecrew/pull/315)** — two
  rounds. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/315#pullrequestreview-5127121195)
  holds the pull request on the one mechanism it invented: "Neutering it
  locally ... leaves `go test ./...`
  entirely green. Every branch this PR added to that function is dead to the
  suite." It then observes that the missing case is the one the PR body, the
  plan and the Decision all describe in the same words — "an absent root
  `AGENTS.md` beside a kept `CLAUDE.md` holding `@AGENTS.md`" — pastes the
  output of the built binary on a synthetic spoke both with and without the
  mechanism, and supplies the four-line table row that closes it. The
  [round-two approval](https://github.com/radiusred/gh-codecrew/pull/315#pullrequestreview-5127160910)
  re-measures the fix by mutation, quotes the four assertions that then fail,
  and checks the rebase by diffing the round-one patch against the new one.
- **PR [#314](https://github.com/radiusred/gh-codecrew/pull/314)** — three
  rounds, the third a rebase check. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/314#pullrequestreview-5127121996)
  is the finding that produced the Deviation above, and it is precise about
  what it is asking for: not a behaviour change but a claim made true — "The
  claim is the argument for the second layer, so it should be true or it
  should not be made" — with the remedy spelled out, including "a short
  follow-up comment on the task saying which layer each assertion actually
  pins". The
  [round-two approval](https://github.com/radiusred/gh-codecrew/pull/314#pullrequestreview-5127162948)
  re-measures both layers itself, notes that the Decision comment is
  unedited (`created_at` equals `updated_at`), and leaves a further
  non-blocking note that one bullet of the new test-file header had become
  false in the fixing — which the third commit then corrected. It was
  dismissed at 00:31:42Z by the rebase onto #315's merge; the
  [round-three approval](https://github.com/radiusred/gh-codecrew/pull/314#pullrequestreview-5127188399)
  compares the approved patch with the rebased one and names the one
  difference, that header comment.
- **PR [#317](https://github.com/radiusred/gh-codecrew/pull/317)** — three
  rounds, the third a rebase check, and the review that produced two of the
  milestone's four open captures. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/317#pullrequestreview-5127132501)
  carries two findings, one record ask and two judgments, and opens by saying
  what is right before what is wrong. Finding 1 is the third unreadable case
  the silent-when-clean Decision had not counted, quoted against the
  Decision's own words. Finding 2 is a test that "does not test what it
  says": "deleting the `if !task.Closed` early return from
  `staleBranchAction` and re-running leaves this test **PASS** ... the output
  is empty for the wrong reason", with the remedy being to assert cost rather
  than output. Finding 3 asks for no code at all — that the one-page
  `closingIssuesReferences` read be recorded as a Decision naming M15-R3, "so
  the boundary is in the record rather than in the source", because "QA
  verifying R3 will not go looking" in a Go comment. The
  [round-two approval](https://github.com/radiusred/gh-codecrew/pull/317#pullrequestreview-5127172850)
  verifies each fix by mutation and quotes the failure messages, and was
  dismissed at 00:42:12Z by the rebase onto #314's merge. The
  [round-three approval](https://github.com/radiusred/gh-codecrew/pull/317#pullrequestreview-5127213266)
  is the rebase check: "the added and removed lines are byte-identical (608
  lines each, diff empty)", with both incoming changes from the two merges in
  between confirmed present.
- **PR [#320](https://github.com/radiusred/gh-codecrew/pull/320)** — two
  rounds, and the closing-reference discipline this milestone shipped
  applied to its own last pull request. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/320#pullrequestreview-5127252184)
  is one finding and it is about a single word: the body's opening paragraph
  read "this body carries no auto-close reference", and "still contains the
  close-family word `close`". The
  [round-two approval](https://github.com/radiusred/gh-codecrew/pull/320#pullrequestreview-5127263637)
  confirms the head is unchanged, the body carries no close-family keyword
  before a ref, and GitHub reports #311 alone — "matching the accepted
  linked-branch decision on task 311".
- **PR [radiusred/codecrew-www#28](https://github.com/radiusred/codecrew-www/pull/28)**
  — [approved first round](https://github.com/radiusred/codecrew-www/pull/28#pullrequestreview-5127314740)
  at 01:24:39Z. It reviews a document the way the other rounds review code:
  the site is built (`zensical build --clean --strict`), the generated
  artefacts are regenerated to prove they are already current
  ("Regenerating left `zensical.toml`, `docs/blog/archive.md`, and
  `docs/blog/atom.xml` unchanged"), the built HTML is opened, and "Every
  GitHub source URL in the post returned 200". Its verdict names the four
  things a post from this seat can get wrong: "I found no factual overreach,
  hype, private-repository naming, or closing-reference problem."

**What the four change requests were about.** #315's was a mechanism the
pull request invented that no test held, proved by neutering it and watching
the suite stay green. #317's carried two findings: a failure path that
deleted the whole stale-branch report in silence, against the Decision that
had just promised silence means clean, and a test that passed with its own
subject deleted. #314's was not about the code at all but about a recorded
claim the tests did not support. And #320's was one close-family word in a
pull request body. Every approval in this milestone was reached with the
reviewer running something: built binaries on synthetic 1.x repositories, a
forced multi-page REST endpoint, normalisation calls stripped one layer at a
time, guards turned into deleters, patch-to-patch comparisons across each
rebase, a GraphQL query run verbatim against PR #294, and a full site
build.

**All three dismissed approvals are the same shape**, and it is the shape
five parallel lanes make inevitable: an approval, then a rebase onto the
merge that landed underneath it, then a round that reviews the rebase alone
and says so in its own heading. #312's round one stood two minutes and
thirteen seconds; #314's round two, two minutes; #317's round two, eight
minutes and twenty-nine seconds. In all three cases round three — or two —
verified the rebase mechanically rather than re-reading the diff.
[M14's record](14-adoption-tidy-and-the-v2-0-0-release.md#the-review-rounds)
describes two such dismissals from four lanes; this milestone had three from
five.

**The one hunk every parallel lane touched was the same one.** #313 merged
first and was never force-pushed; #320 opened after everything else had
merged and was never force-pushed either. The other four were force-pushed —
once on #312, twice on #315, three times each on #314 and #317 — and every
round that followed one of those pushes checks the same file by name: that `CHANGELOG.md`'s
`[Unreleased]` section still carries the entries that were already on `main`
as well as the branch's own. "`CHANGELOG.md` keeps both the #308 and #307
Unreleased entries" (#312); "carries all four entries that were on `main` ...
plus this task's new one. Nothing was dropped or reordered" (#315); "the
CHANGELOG entry sits above #308's ... and #308's entry is intact" (#314);
"the Unreleased section carries main's six entry groups untouched ... plus
this PR's" (#317). The implementer contract's instruction for that hunk is
one sentence — keep both entries — and it was exercised four times in eighty
minutes.

## QA: one round, seven verdicts

The qa role holder (`radiusred-testy[bot]`) verdicted every requirement in
[one comment](https://github.com/radiusred/gh-codecrew/issues/305#issuecomment-5563791131)
at 01:36:41Z, eleven minutes after the last delivery merged. All seven came
back `satisfied`, so there is no supersession to read and no remedy loop to
follow — the standing verdict for each requirement is its only verdict.
[M13's record](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#qa-three-rounds-and-the-requirement-that-took-all-three)
describes three rounds and a remedy task;
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md#qa-two-rounds-and-the-requirement-that-was-answered-rather-than-remedied)
describes two and a Decision that cleared a failure with nothing built.

Every verdict opens on the same cross-cutting check — "Main
build/test/vet/gofmt are clean" — and then names what was run for that
requirement specifically. Three are worth quoting for what they establish
beyond the tests.

**M15-R1 was verdicted against a scratch repository, not against the suite.**
The seat built a 1.x hub and ran the verb on it: `migrate --dry-run` listing
`.codecrew/AGENTS.md`, `AGENTS.md` and `CLAUDE.md`; the live run "writing
root `AGENTS.md`/`CLAUDE.md` byte-equal to a fresh `init --hub` scaffold in
the same `chore: migrate codecrew to the 2.0 layout` commit"; and a kept
non-reaching `AGENTS.md` producing a dry run that names "only `Kept:
AGENTS.md` under action needed". That is the requirement's own sentence —
"a migrated spoke is indistinguishable from a fresh 2.0 scaffold" — measured
rather than inferred.

**M15-R3's verdict says what it could not check, and why.** The requirement
is about a listing past a hundred rows, and the seat looked for one: "the
requested live >100-comment hub issue does not exist under the App token
today (highest issue comment count is #119 with 69), so the live boundary is
covered by code/tests and the M15 dry-run close path rather than a natural
>100-comment record." The tests it credits instead are the ones built for
exactly this — "the paging tests drive 101 comments and 150-row listings so
dropping `--paginate` truncates".

**M15-R5's verdict includes the capture bookkeeping.** Besides the report
being silent on a hub whose `git ls-remote origin 'refs/heads/task/*'` is
empty, and the contracts carrying the new checks, it records that "backlog
captures #316, #318 and #319 remain open and outside the adopted set" — the
distinction between what a milestone fixed and what it merely found, checked
rather than assumed. It also exercises a verb against a closed task:
`task finish 306 --dry-run` "refuses `CLOSED` and prints `dry run: nothing
written`".

The other four are verdicted against the merged code and the shipped
artefacts: M15-R2 on the `task start` and `status` tests, SPEC §6's rows, and
the operator's own record capturing live `status` "naming an App-held
in-review task from its start record before all M15 tasks closed"; M15-R4 on
the CRLF table "through GitHub readers and directly at scanner entry", SPEC
§4's sentence, and `TestTaskNewAppliesTheTaskLabel` asserting exactly
`cc:task`; M15-R6 on the installed extension, the release's five assets, the
tag object, the run, the changelog and the operator record; and M15-R7 on the
merged spoke pull request, the closed capture, the live page returning HTTP
200, the nav block, and the post's present tense against v2.0.1. The R7
verdict also names what it is deliberately not counting:
"radiusred/codecrew-www#29 remains an open backlog capture for the nav label
HTML leak outside this requirement."

## Requirement outcomes

The status column is the **standing** verdict word as the qa role holder
wrote it, verbatim and unqualified. Each requirement was verdicted once, so
the standing verdict is the only one.

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M15-R1 — `migrate` writes an absent root `AGENTS.md` or `CLAUDE.md` from the same scaffold constants `init` writes, in the same pathspec commit as the rest of the move, and `--dry-run` lists each file it would write; the `action needed` block names only a kept root file that does not reach `.codecrew/AGENTS.md`, so a migrated spoke is indistinguishable from a fresh 2.0 scaffold; a kept root file is still never rewritten; SPEC §6's `migrate` row follows (adopts [#301](https://github.com/radiusred/gh-codecrew/issues/301)) | [#306](https://github.com/radiusred/gh-codecrew/issues/306) / PR [#315](https://github.com/radiusred/gh-codecrew/pull/315) | `satisfied` | Task closed; verdicted on a scratch 1.x hub probe — dry run, live run and a kept non-reaching file — against a fresh `init --hub` scaffold. Two review rounds; round one blocked on the one mechanism the PR invented having no test that failed without it. Five Decisions, including the edit to the released 2.0.0 changelog section — marked `(2.0.1, #306)` in the two passages it touched — and the operator's, after this record's first commit, declining the marker for the third |
| M15-R2 — `task start` does not try to assign an `app:`-typed caller and prints no note about it — the routing table says the kind — while a `user:`-typed caller is still assigned; `status` shows who holds an in-progress or in-review task from its latest **Started by** record rather than from the assignee list, so an App-run task names its holder; SPEC §6's `task start` and `status` rows follow (adopts [#287](https://github.com/radiusred/gh-codecrew/issues/287)) | [#307](https://github.com/radiusred/gh-codecrew/issues/307) / PR [#313](https://github.com/radiusred/gh-codecrew/pull/313) | `satisfied` | Task closed; verdicted on the targeted tests, SPEC §6's rows and the operator's live `status` naming an App-held in-review task from its start record. Approved first round, with the review running the built binary against this hub to see the capture's symptom gone. Three Decisions; the third declines to teach `InferState` the start record because that is SPEC §4's lifecycle table, and is captured as [#316](https://github.com/radiusred/gh-codecrew/issues/316) |
| M15-R3 — every listing read in the tracker is paginated — `Comments`, `SubIssues`, `listIssues` and any other first-page read — so a milestone issue past a hundred comments keeps its newest QA verdicts and `milestone close` cannot report `VERDICT_MISSING` for a requirement whose satisfied verdict sits on a later page; a test drives each paginated reader with a synthetic listing longer than one page (adopts [#264](https://github.com/radiusred/gh-codecrew/issues/264)) | [#308](https://github.com/radiusred/gh-codecrew/issues/308) / PR [#312](https://github.com/radiusred/gh-codecrew/pull/312) | `satisfied` | Task closed; verdicted on the source reads and the paging tests — 101 comments and 150-row listings — with the verdict stating that no live issue past a hundred comments exists to check against. Two review rounds, the second a rebase check; the first approval was dismissed two minutes and thirteen seconds after it was submitted. Three Decisions, two of them scope boundaries: `--paginate` over cursor loops, and the GraphQL relations left unpaginated |
| M15-R4 — the tracker normalises CRLF line endings to LF once at its boundary — every issue body and comment body it reads — so each line-anchored scan reads a body GitHub's web editor saved with CRLF exactly as it reads the LF one, with a table test per scanner on a CRLF fixture and no regexp changes; and the `task new` tests assert the created issue carries exactly the `cc:task` label (adopts [#296](https://github.com/radiusred/gh-codecrew/issues/296), [#297](https://github.com/radiusred/gh-codecrew/issues/297)) | [#309](https://github.com/radiusred/gh-codecrew/issues/309) / PR [#314](https://github.com/radiusred/gh-codecrew/pull/314) | `satisfied` | Task closed; verdicted on rows read both through the GitHub readers and directly at scanner entry, SPEC §4's new sentence, and the label assertion. Three review rounds, the third a rebase check; round one found the second layer held by no test and a Decision's claim about it false. Four Decisions and three Deviations — the normalisation is applied twice where the requirement says once, and the record says why |
| M15-R5 — `status` reports stale task branches — `task/<n>-…` branches on the remote whose task issue is closed — with the same delete-or-keep verdict and reason `milestone close`'s two sweeps would give, computed by the same code so the report and the sweep can never disagree, one prefix-filtered listing per repo and a note when that listing was partial; and the record guards against accidental closing references: the implementer contract says no closing keyword may precede any ref in a PR body but the task's own, and `task finish` (and its `--dry-run`) prints a note naming every issue the PR's closing references would close besides the task itself, before the merge (adopts [#295](https://github.com/radiusred/gh-codecrew/issues/295), [#303](https://github.com/radiusred/gh-codecrew/issues/303)) | [#310](https://github.com/radiusred/gh-codecrew/issues/310) / PR [#317](https://github.com/radiusred/gh-codecrew/pull/317) | `satisfied` | Task closed; verdicted on a hub whose remote carries no `task/` branch, on `task finish 306 --dry-run` refusing `CLOSED` and writing nothing, on both contracts through the installed release's `roles diff`, and on the three filed captures being open and outside the adopted set. Three review rounds, the third a rebase check; round one found an unreported failure path, a test that passed with its behaviour deleted, and asked for a Decision rather than a code change. Four Decisions and one Deviation |
| M15-R6 — v2.0.1 ships: the `[Unreleased]` changelog becomes the 2.0.1 section, a fresh empty `[Unreleased]` above it and the compare links updated, and `docs/introduction.md` and `README.md` name v2.0.1 where they name the release; the annotated tag is cut by the implementer identity on the flip merge commit; `release.yml` builds it with its five assets; the evidence (tag SHA, run, assets) is posted on the task; the operator verifies `gh extension upgrade codecrew` prints v2.0.1 and `status` runs clean on this hub, and records that on this issue | [#311](https://github.com/radiusred/gh-codecrew/issues/311) / PR [#320](https://github.com/radiusred/gh-codecrew/pull/320) | `satisfied` | Task closed; verdicted on the installed extension, the release, the tag and the run, and on the operator's record. Tag `v2.0.1` (tag object `fabdc17` → commit [2dc147a](https://github.com/radiusred/gh-codecrew/commit/2dc147a)) by `radiusred-cody[bot]` at 01:08:33Z, run 34072047784 on that commit, five assets, release published 01:09:55Z. Two review rounds, the first over a close-family word in the body's prose. No migration leg: the protocol stays at 2.0 |
| M15-R7 — codecrew.works publishes its first blog post, by the doc-synthesizer seat: protocol 2.0 and the migration — why a protocol major, what changed (one plain paragraph each, no code tour), what an adopter does (the exact steps from the CHANGELOG 2.0 entry), and what deliberately did not change — in the site voice, in the present tense against the shipped release, listed in the nav `BEGIN_BLOG_POSTS` block, and merged only after M15-R6 (adopts [radiusred/codecrew-www#26](https://github.com/radiusred/codecrew-www/issues/26)) | [radiusred/codecrew-www#27](https://github.com/radiusred/codecrew-www/issues/27) / PR [radiusred/codecrew-www#28](https://github.com/radiusred/codecrew-www/pull/28) | `satisfied` | Task closed; verdicted on the merged pull request, the closed capture, the live page returning HTTP 200, the nav block and the post's present tense against v2.0.1. Approved first round, by a review that rebuilt the site strictly and re-resolved every GitHub URL in the post. Three Decisions, all about what a public page may say. The sequencing held by timestamps: the release published at 01:09:55Z, the task started at 01:14:35Z |

**Captures adopted and closed by the merges — eight**, and none of them by a
`Closes` line in a pull request body:
[#287](https://github.com/radiusred/gh-codecrew/issues/287) by
`task finish 307`,
[#264](https://github.com/radiusred/gh-codecrew/issues/264) by
`task finish 308`,
[#301](https://github.com/radiusred/gh-codecrew/issues/301) by
`task finish 306`,
[#296](https://github.com/radiusred/gh-codecrew/issues/296) and
[#297](https://github.com/radiusred/gh-codecrew/issues/297) by
`task finish 309`,
[#295](https://github.com/radiusred/gh-codecrew/issues/295) and
[#303](https://github.com/radiusred/gh-codecrew/issues/303) by
`task finish 310`, and
[radiusred/codecrew-www#26](https://github.com/radiusred/codecrew-www/issues/26)
by `task finish 27`. Each carries two comments and nothing else: an
**Adopted by** note written when the task was created, and a closure naming
the task, the pull request and the merge commit.

**Captures filed during the milestone and left open — four**, every one
filed by the operator from a seat's own report. Three are in this hub, and
all three come out of M15-R2's and M15-R5's tasks.
[#316](https://github.com/radiusred/gh-codecrew/issues/316) — an App-run task
without a pull request reads `[ready]`, because SPEC §4's lifecycle table
infers "in progress" from the assignee that an App can never be — comes from
the implementer's third Decision on #307, filed at 00:06:27Z from the seat's
own report and four minutes before the review that independently reached the
same reading; the capture is explicit that the fix "is a protocol minor
(2.1), or folded into the next major".
[#318](https://github.com/radiusred/gh-codecrew/issues/318) — `task finish`
does not say when the pull request's closing references *omit* the task —
comes from the reviewer's round-one finding 5 on PR #317, which observed that
the new read holds exactly the data for the inverse defect and does nothing
with it. [#319](https://github.com/radiusred/gh-codecrew/issues/319) — a
failed branch comparison drops a branch from the sweep and from the new
report silently — comes from the same review's round two, against the SPEC
sentence that same pull request added, "that the report's silence may only
ever mean clean". The fourth is in the spoke:
[radiusred/codecrew-www#29](https://github.com/radiusred/codecrew-www/issues/29),
the blog nav label's HTML leaking into the page title, filed from the
doc-synthesizer's report on the merged post. Nobody could have found it
before, because "`docs/blog/posts/` was empty until #27".

## Protocol-discipline observations

Nine things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **The adoption link M14 shipped closed every capture this milestone
  adopted, and no pull request body closed one.**
  [M14's record](14-adoption-tidy-and-the-v2-0-0-release.md#protocol-discipline-observations)
  says all four of its own adopted captures were closed by the convention
  `task new --adopts` exists to replace, "because all four tasks were opened
  before it existed", and that the two tasks opened after the flag landed
  adopted nothing. Every one of this milestone's eight adoptions ran on the
  flag instead: each capture carries an **Adopted by** comment written when
  the task was created and a closure written by `task finish` naming the
  task, the pull request and the merge commit, and every pull request body's
  `closingIssuesReferences` names its task alone. A requirement delivered in
  one milestone did its whole job for the first time in the next.
- **M14-R4's second branch pass ran for the first time sixteen minutes before
  this milestone opened, and reported nothing.** M14's record closes on the
  observation that "M14-R4 has never run" and that "the first close that will
  exercise it is this milestone's, after this document merges". It ran
  earlier than that: `milestone close 14`'s
  [comment](https://github.com/radiusred/gh-codecrew/issues/269#issuecomment-5563049997)
  at 23:34:56Z reads, in full, "all 7 tasks done, milestone document merged"
  — no branch named, in a verb whose contract is to report each one it
  considers. Its candidate set is the hub plus every repository the
  milestone's tasks name, which for M14 was this hub alone. M15's close will
  be the first with a second repository in that set, because one of its tasks
  lives in the spoke.
- **Five parallel lanes cost three dismissed approvals, three rebase rounds
  and four passes over one changelog hunk — and the cost was paid entirely by
  the reviewer.** Nothing in the protocol sequences pull requests and nothing
  needs to: every rebase was clean and every rebase check passed. But three
  of the milestone's fourteen review submissions exist only because of the
  ordering, three approvals were thrown away by force-pushes that changed no
  behaviour, and the `[Unreleased]` section of `CHANGELOG.md` was reconciled
  four times in eighty minutes.
  [M14's record](14-adoption-tidy-and-the-v2-0-0-release.md#protocol-discipline-observations)
  reports two dismissals from four lanes; one more lane here produced one
  more dismissal.
- **A review asked for a record rather than a change, and got one.** Round
  one on PR #317 agreed with the code — "**I agree with the call**" — and
  still blocked on where the reasoning lived: the one-page
  `closingIssuesReferences` read was justified in a Go comment while M15-R3
  was open in another lane requiring every listing read to paginate, and "QA
  verifying R3 will not go looking" in the source. The remedy was a Decision
  comment on the task naming M15-R3 and arguing the boundary. Nothing
  compiled differently; a boundary that existed only in code now exists where
  the next reader will look for it.
- **The record grammar read two honest paragraphs oppositely, in the same
  milestone.** `**Deviation:** none from the plan` on #309 opens a paragraph
  with a bare label, so `ExtractRecords` counts a Deviation record whose
  content is that there was no deviation. `**Deviation from nothing else:**`
  on #308 opens a paragraph that describes a real one — a test helper's
  matcher changed to accommodate the new flag — and the scanner does not see
  it, because the label is neither bare nor parenthetically qualified. Both
  are the seat writing honestly; the grammar reads one in and one out. Both are
  counted here as the scanner counts them, because a count no reader can
  reproduce from the source is worse than one that needs explaining.
- **One question reached the operator, and it was the one M14 had already
  answered.** #311's Deviation on the linked-branch closing reference went up
  at 00:56:23Z and came back as a Decision at 00:57:29Z — one minute and six
  seconds, because
  [M14's Deviation on #272](https://github.com/radiusred/gh-codecrew/issues/272#issuecomment-5561314143)
  had already established there is no public handle for the link and
  [its Decision](https://github.com/radiusred/gh-codecrew/issues/272#issuecomment-5561318885)
  had already chosen the answer. It did not raise a gate, apply
  `cc:needs-decision` or call `checkpoint`, which
  [M13's record](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#protocol-discipline-observations)
  observed of one question and
  [M14's](14-adoption-tidy-and-the-v2-0-0-release.md#protocol-discipline-observations)
  of three. Five questions across three milestones have now taken this route
  and none has used the verb the protocol provides for it.
- **The guard this milestone shipped was applied to its own last pull
  request, and the finding was in prose that created no reference.** PR
  #320's body opened "this body carries no auto-close reference", and round
  one requested a change over the word `close` inside `auto-close`. The
  contract rule M15-R5 added forbids a closing word *before a ref*, and that
  word preceded none; `closingIssuesReferences` already named #311 alone.
  The review cites "the task-specific review instruction", which is stricter
  than the contract — no closing keyword in the body at all. Worth recording
  precisely, because the gap between what the shipped rule says and what the
  dispatch asked for is the gap a future reader will otherwise have to
  reconstruct.
- **A milestone record's counts have a shelf life, and this one's expired
  between its commit and its review.** The document was committed at
  02:06:04Z with twenty-five Decisions and six Deviations counted; the
  operator's #306 Decision landed at 02:09:01Z and a second #305 Deviation at
  02:11:36Z, and the review that found them opened on the observation that
  "the record's whole claim on a reader is that its counts are re-derivable
  from the trail". Nothing comparable happened to
  [M14's record](14-adoption-tidy-and-the-v2-0-0-release.md): the only
  comment its milestone issue gained after synthesis was `milestone close`'s
  own, and that came after the record had merged. The protocol has no verb
  that freezes a trail, and no rule about when a record is counted; this one
  was settled by the coordination layer holding the milestone still and by
  the record naming the instant it counted. Both are conventions this
  milestone invented on the spot.
- **The refresh obligation found two incomplete sentences and nothing false.**
  Both version claims — `docs/introduction.md`'s `**Shipped:**` line and its
  `version` example, and `README.md`'s "Already running CodeCrew?" paragraph
  — were flipped by M15-R6's own pull request and read v2.0.1 before this
  task started; the installed extension prints `v2.0.1 (protocol 2.0)`; SPEC
  §10's catalogue still counts forty-three refusal codes, because the
  milestone added none. What was out of date was the introduction's sentence
  enumerating what the verbs do: it named three things `task finish` does and
  there are now four, and it did not mention `status`'s holder line or its
  stale-branch report at all. Both are added here. `README.md` is unchanged:
  its one sentence about `status` is a landing page's summary rather than an
  enumeration and is as true after this milestone as before.

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **The conversation that opened the milestone is on the record only through
  #305's goal.** "Decided with the operator 2026-09-07" is the whole of it.
  The seven captures each carry their own history, but the judgement that
  they belong together — and that the release and the blog post belong with
  them — has no transcript. *The reasoning is recorded; the deliberation is
  not*, as
  [M14's record](14-adoption-tidy-and-the-v2-0-0-release.md#what-the-record-does-not-contain)
  and [M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#what-the-record-does-not-contain)
  each say of their own openings. That sentence also dates itself a day
  ahead of the API: GitHub timestamps #305 at 2026-09-06T23:51:44Z, and the
  spoke's task, created fifty-five seconds later, plans a post dated
  2026-09-07. Everything written after midnight UTC agrees with its own
  prose, so the offset is a local clock rather than a mistake — but it is
  the same class of discrepancy
  [M14's record](14-adoption-tidy-and-the-v2-0-0-release.md#what-the-record-does-not-contain)
  reports of its own trail, and every time in this record is taken from the
  API for that reason.
- **Which harness ran ten of the fourteen review rounds is not recorded.**
  The operator's Deviation names three that codex ran — #312's two rounds and
  #314's third — and one it could not, #317's third. The other ten,
  including every round on #313, #315 and #320 and the spoke's, are not
  attributed anywhere. The identity is constant across all fourteen, so the
  gate the protocol enforces is intact; the harness behind each round is not
  reconstructable, exactly as
  [M14's record](14-adoption-tidy-and-the-v2-0-0-release.md#what-the-record-does-not-contain)
  and [M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#what-the-record-does-not-contain)
  report of their own.
- **A third passage of the released 2.0.0 changelog section still describes
  the old behaviour — and this is the one gap in this list that was closed
  after the list was written.** PR #315's round-two approval named it as a
  non-blocking note: the `**What it writes.**` bullet under
  `### gh codecrew migrate: the one-shot move to the 2.0 layout` says an
  absent root entry point "is named under an `action needed` heading", it
  sits *between* the two passages that were corrected, and it is the only one
  of the three without a `(2.0.1, #306)` marker. It is still there, and this
  document's first commit recorded that no Decision said whether that was
  deliberate. The operator's
  [Decision](https://github.com/radiusred/gh-codecrew/issues/306#issuecomment-5564002686)
  three minutes later is that Decision: the marker is declined, because the
  unmarked passage describes what the 2.0.0 binary wrote and is still true of
  2.0.0, while the two that carry markers are the adopter's steps. The bullet
  stays to say so, because the paragraph it replaces asserted the opposite.
- **M15-R3's live boundary was never crossed.** The QA verdict says so in its
  own words: no issue in this hub is past a hundred comments — "highest issue
  comment count is #119 with 69" — so the requirement is verdicted from the
  paging tests and the source rather than from a listing that actually needed
  a second page. The condition #264 describes has still never been observed
  in this project, only prevented.
- **Whether the stale branches that motivated M14-R4 still exist is not
  known, and nothing in this milestone looked.** They are in
  `radiusred/numberguess`, per
  [M14's record](14-adoption-tidy-and-the-v2-0-0-release.md#decisions), and
  neither mechanism reaches it: `milestone close`'s candidate set is the hub
  plus the repositories its own tasks name, and M15-R5's new `status` report
  covers "the repo you are standing in" by an explicit Decision. The report
  makes a skipped sweep visible where someone is already working; it does
  nothing for a repository nobody visits, which is the case #167 was filed
  from.
- **What `migrate` printed on the fleet's ten successful runs is still not on
  the record.**
  [M14's record](14-adoption-tidy-and-the-v2-0-0-release.md#what-the-record-does-not-contain)
  reports that the migration record is eleven `git log` excerpts and nothing
  else, and that [#301](https://github.com/radiusred/gh-codecrew/issues/301)
  exists because the operator noticed one behaviour while running it and
  filed it by hand. This milestone fixed that behaviour; it did not recover
  what the other runs printed, and nothing can now.
- **The companion post named in the adopted capture is not part of this
  milestone.** [radiusred/codecrew-www#26](https://github.com/radiusred/codecrew-www/issues/26)
  ends on "A short companion post on radiusred/www points here (radiusred/www
  capture filed the same day)". M15-R7 does not name it, no task delivered
  it, and nothing on this milestone's trail says whether it exists. Nothing
  records where the published post was announced or who has read it, either.
- **Nothing measures what the milestone cost.** Seven pull requests, fourteen
  reviews, one QA round and a release, with no record of tokens, wall-clock
  per seat, or how many dispatches were made.
  [M14's record](14-adoption-tidy-and-the-v2-0-0-release.md#what-the-record-does-not-contain)
  and [M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#what-the-record-does-not-contain)
  say the same of their own, and the answer has not changed. What the trail
  does carry is the wall-clock of the work itself — one hour and thirty-four
  minutes from the milestone opening to the last merge — which is a fact
  about the lane rather than about its cost.
