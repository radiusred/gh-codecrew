# M17: The frontier builds the coordinator

Tracking issue: [#345](https://github.com/radiusred/gh-codecrew/issues/345) ·
Synthesized 2026-09-23 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #345's **two** requirements as opened, **none added,
none struck and none amended** — the issue's `userContentEdits` query returns
a total count of zero; its Gates section, left as the scaffold's placeholder;
its **two comments**, both the qa role holder's and together holding **three
verdicts**; the **three** delivery task issues, all three in the spoke
([radiusred/codecrew-www#37](https://github.com/radiusred/codecrew-www/issues/37),
[radiusred/codecrew-www#39](https://github.com/radiusred/codecrew-www/issues/39),
[radiusred/codecrew-www#42](https://github.com/radiusred/codecrew-www/issues/42)),
plus [#351](https://github.com/radiusred/gh-codecrew/issues/351), the task
this document is; their **three merged pull requests** and the **four
commits** on them; the **six Decision records and two Deviation records**
across the milestone issue and all four of its task issues, of which **three
and one** are on the three delivery tasks, **three and one** on this
document's own, and **none** on the milestone issue; the **four review
submissions** on the three pull requests — **one change request** and **three
approvals, every one of which stands**, none dismissed; **no** adopted backlog
captures; and the **six** captures filed in the milestone's window and left
open, of which three are the milestone's own and three came from other
projects' use of the protocol. The house form is
[M16's document](16-bedding-in-the-venue-seam-and-the-reference.md), with
[M15's](15-v2-0-1-what-the-fleet-migration-taught.md) behind it; the record
standard is the one the reviewer set on
[PR #322](https://github.com/radiusred/gh-codecrew/pull/322) and re-applied
on [PR #341](https://github.com/radiusred/gh-codecrew/pull/341), so every
count above and below is re-derived from the source — records by the rule
`tracker.ExtractRecords` itself applies (a paragraph-initial `**Decision:**`,
`**Deviation:**` or `**Gate resolved:**` label, bare or parenthetically
qualified), review states from each pull request's reviews API and timeline,
commits from each pull request's commits API, and the requirement history
from the milestone issue's `userContentEdits`.

**Counted at 2026-09-23T11:22:52Z, against a trail frozen from this document's
dispatch.** The coordination layer's last posts on the milestone before that
instant were the opening of #351 at 11:20:08Z and its Deviation seven seconds
later, and it posts nothing further on the milestone's issues before this
pull request merges. This document's own task's **three Decisions** were
written at 11:22:45Z–11:22:47Z, before the count, and are inside every total
above, as
[its first Decision](https://github.com/radiusred/gh-codecrew/issues/351#issuecomment-5793927587)
says they would be — the method
[M16's record](16-bedding-in-the-venue-seam-and-the-reference.md) used.

**The whole milestone happened in the spoke.** `main` in this repository did
not move between M16's record
([7a2febd](https://github.com/radiusred/gh-codecrew/commit/7a2febd), committed
2026-09-07T14:52:09Z) and this document's pull request, which carries the only
two hub commits the milestone makes. `main` in
[radiusred/codecrew-www](https://github.com/radiusred/codecrew-www) moved by
exactly four commits, in a straight line from M16's last spoke merge
([1060594](https://github.com/radiusred/codecrew-www/commit/1060594)):
[b8fadc4](https://github.com/radiusred/codecrew-www/commit/b8fadc4) and
[066b97b](https://github.com/radiusred/codecrew-www/commit/066b97b) (the post),
[5c2a9e6](https://github.com/radiusred/codecrew-www/commit/5c2a9e6) (the front
page) and [199db8c](https://github.com/radiusred/codecrew-www/commit/199db8c)
(the remedy).

**The trail is checked by the verb, and it cites nothing on github.com.** The
installed extension reports `v2.0.1 (protocol 2.0)`, which is still the shipped
version because M17 cut no release. Run from this repository at the counting
instant:

```
$ gh codecrew milestone evidence 17
requirements counted: M17-R1, M17-R2 (2)
all 3 cited links resolve across 5 issues — evidence is reachable
```

Five issues: the milestone and its four sub-issues. Three links, all of them
outside github.com and all three in the QA seat's comments and the remedy
task: the announcement at `claude.com`, the live post, and the site's front
page. The verb counts distinct URLs, so the announcement's four citations
across #345, www#37 and www#42 count once. Every other reference on the trail
is a bare `#N` or `owner/repo#N`, which the verb does not read as a link.

The gather from `gh codecrew milestone close 17 --dry-run` is again not part of
the raw material, for the reason
[M11's](11-housekeeping.md) through
[M16's](16-bedding-in-the-venue-seam-and-the-reference.md) records all give:
`--dry-run` stops at the gate that counts tasks, naming the task that writes
this file. As it read at the counting instant, before this document's pull
request opened:

```
gate milestone open: ok
gate no gate raised: ok
gate tasks closed: refused[OPEN_TASKS]: tasks not closed: radiusred/gh-codecrew#351 (ready)
gate requirements declared: not reached
gate QA verdicts: not reached
gate milestone document: not reached
dry run: nothing written — the live verb stops at the first refusal above
```

Every record below was read from its issue, pull request, review, commit or
timeline directly. The prose is wrapped at this file's normal width; the
requirement-outcomes table's rows and a handful of single links are not,
because a newline ends a Markdown table row and breaking a link breaks it —
the shape [M16's record](16-bedding-in-the-venue-seam-and-the-reference.md#requirement-outcomes)
carries, and the shape `ROADMAP.md` has had since M1.

**This PR adds the M17 ROADMAP row; it does not flip one.** `ROADMAP.md`
carried no M17 row before this PR, so there was nothing to discard.

## Goal and outcome

#345 was opened at 09:48:01Z on 2026-09-19, eleven days, eighteen hours and
fifty-three minutes after
[`milestone close 16`](https://github.com/radiusred/gh-codecrew/issues/325#issuecomment-5572420267)
closed M16 with the line "all 9 tasks done, milestone document merged". Its
goal names its occasion in the first sentence: "Anthropic's redesigned
Projects (2026-09-17) puts a coordinator, branch-per-thread cloud sessions and
project memory into Claude Code itself — the coordinator shape CodeCrew
reached in M7 and under Paperclip." The second sentence sets the scope: "This
milestone is the public response and nothing more: a codecrew.works post that
reports what overlaps and what does not, and the front-page line it promises."

The third carries M16's rule forward — "The bedding-in rule stands —
no verb changes its arguments, output, writes or refusal codes, and SPEC.md is
untouched" — and the fourth puts the experiment outside the milestone: "The
Projects-as-orchestrator experiment is #344 and stays a capture: it is gated
on beta access and is findings, not a requirement." The provenance line is one
sentence: "Decided with the operator 2026-09-19."

The two requirements, as the body carries them:

- **M17-R1** — the post: "written in the doc-synthesizer's voice, quoting the
  announcement verbatim where it is quoted, naming what Projects does that
  CodeCrew never owned (the running) and the four protocol claims it leaves
  untouched (the open record, no self-review, any model in a seat, the human
  gate as a protocol event), hedging the identity claim as unverified, linking
  #344 as the experiment; the essay framed as one operator's experience and
  not positioned against GSD; the strict site build clean".
- **M17-R2** — the front page "carries the post's one-line position — the
  coordinator is a commodity, the record and the separation of duties are not
  — as a sentence in the existing intro, not a new section; the sync table and
  everything else on the page unchanged".

Neither requirement adopts a capture, and no task carries a `## Adopts`
section. That is a different shape from
[M16's](16-bedding-in-the-venue-seam-and-the-reference.md#goal-and-outcome),
where each of six requirements adopted exactly one.

The three delivery tasks, in merge order:

- **[radiusred/codecrew-www#37](https://github.com/radiusred/codecrew-www/issues/37)
  / PR [radiusred/codecrew-www#38](https://github.com/radiusred/codecrew-www/pull/38)**
  — M17-R1, delivered by the doc-synthesizer seat (`radiusred-wordy[bot]`),
  which the requirement's "written in the doc-synthesizer's voice" names. Two
  commits, two review rounds, merged 2026-09-19T11:25:30Z
  ([066b97b](https://github.com/radiusred/codecrew-www/commit/066b97b)). One
  Deviation, no Decisions.
- **[radiusred/codecrew-www#39](https://github.com/radiusred/codecrew-www/issues/39)
  / PR [radiusred/codecrew-www#40](https://github.com/radiusred/codecrew-www/pull/40)**
  — M17-R2, by the implementer seat (`radiusred-cody[bot]`). One commit,
  approved first round, merged 2026-09-23T11:01:36Z
  ([5c2a9e6](https://github.com/radiusred/codecrew-www/commit/5c2a9e6)). Two
  Decisions.
- **[radiusred/codecrew-www#42](https://github.com/radiusred/codecrew-www/issues/42)
  / PR [radiusred/codecrew-www#43](https://github.com/radiusred/codecrew-www/pull/43)**
  — M17-R1's remedy, opened after QA's first verdict, by the implementer seat.
  One commit, approved first round, merged 2026-09-23T11:14:49Z
  ([199db8c](https://github.com/radiusred/codecrew-www/commit/199db8c)). One
  Decision.

**The milestone ran in two sittings four days apart.** www#37 was created
thirteen seconds after the milestone issue, and its pull request merged at
11:25:30Z, one hour, thirty-seven minutes and twenty-nine seconds after the
milestone opened. www#39, the second requirement's task, was not created until
2026-09-23T10:55:03Z — three days, twenty-three hours and twenty-nine minutes
after the first merge — and everything after it, the front page, both QA
rounds, the remedy and this document's task, happened inside the next
twenty-five minutes. From the milestone opening to the last delivery merging:
**four days, one hour, twenty-six minutes and forty-eight seconds**. Nothing on
the trail says why M17-R2 waited; the section on what the record does not
contain carries it.

**No two pull requests were open at once.** www#38 merged four days before
www#40 opened, and www#43 opened seven minutes and forty-five seconds after
www#40 merged.
[M16's record](16-bedding-in-the-venue-seam-and-the-reference.md#goal-and-outcome)
reports four hub pull requests open at once.

What exists afterwards that did not before. codecrew.works carries a second
blog post, "The Coordinator Is a Commodity Now. The Record Is Not.", at
`/blog/posts/2026-09-19-the-coordinator-is-a-commodity-now/`, dated
2026-09-19. The front page's hero paragraph carries one new sentence, as PR
www#40 quotes it: "The coordinator running your agents is a commodity now; the
record and the separation of duties are not." The words "a commodity now" link
the post on-site. The spoke's hero test asserts the
sentence and the link render. Nothing in this repository's code, `SPEC.md`,
`CLI.md` or `docs/` changed.

**M17 shipped no release, and nothing in the trail says otherwise.** No tag was
cut, `docs/introduction.md`'s `**Shipped:**` line still reads v2.0.1 and the
installed extension still prints `v2.0.1 (protocol 2.0)`. The bedding-in rule
was measured rather than asserted: QA's first comment reports that "the hub
main has no `SPEC.md`, `internal/`, or `cmd/` changes since M16". This pull
request's second commit touches `.codecrew/config.yml`, the advisory routing
table; the Deviations section below carries why it rides here.

## Decisions

Six Decision records: three on the delivery tasks, all by the implementer seat
(`radiusred-cody[bot]`) — two on www#39 and one on www#42 — and three on this
document's own, by `radiusred-wordy[bot]`. None is on the milestone issue, none
carries a parenthetical qualifier and none is a gate resolution. www#37
carries no Decision.

### A link in the hero, and a test the brief did not ask for

www#39's
[first Decision](https://github.com/radiusred/codecrew-www/issues/39#issuecomment-5793588158)
links the new sentence to the post, "on the words 'a commodity now', with an
on-site relative link". **Trade-off:** "the hero gains a third link beside its
two buttons, and a link in the lead can draw a first-time reader away from
'Start now'. The post is the argument behind the claim, though; an unlinked
assertion in the hero reads as a slogan, and the linked one shows its working.
The link stays on-site, so the hero test's 'no github.com in the hero' rule
holds." **Rejected:** "no link (a bare claim a reader cannot check, and the
post itself says it is the source of this line); linking the whole sentence
(too much pink in a paragraph that already carries the bold punchline)."

The [second](https://github.com/radiusred/codecrew-www/issues/39#issuecomment-5793588380)
adds one assertion to the existing hero test, so the diff is `docs/index.md`
plus `tests/test_site_build.py`. **Trade-off:** "the dispatch brief asked for
a `git diff --stat main` showing only `docs/index.md`; the implementer contract
and this repo's precedent (#19, #22: every home-page copy change carried a
hero/section assertion) say the change ships with a test." **Rejected:**
"copy-only PR with no test (leaves the requirement unprotected and would need
a backfill later, which the crew does not do)." The reviewer's approval
addresses it: "The test-file addition matches the recorded Decision on
radiusred/codecrew-www#39, and I agree with it". The merged diff is the one the
pull request body shows: two files, three insertions and one deletion.

### Which quotation moves, and which were checked

www#42's [one Decision](https://github.com/radiusred/codecrew-www/issues/42#issuecomment-5793748856)
is the remedy's scope. Only the "reviews the outputs" quotation changes, and
every other quotation in the post was checked against the announcement,
"fetched 2026-09-23, tags stripped, whitespace collapsed", and left as it is.
The record lists them — the whole coordinator sentence, "learns more about the
project details", "why the export was dropped", "also remembers your working
and communication style", "come after that" — and names the three quoted
strings it did not hold to the source because they are not quotations of it:
"the operator decided", "can Projects run a CodeCrew crew" and "the operator's
seats, yes; the crew's, not yet." **Trade-off:** "none in the fix itself; the
check relied on a text extraction of the live page, so a later edit to the
announcement would not be caught." **Rejected:** "dropping the quotation marks
from the fragment (QA's alternative) — the fragment is genuinely verbatim once
the punctuation is right, and the sentence is about it being a quote."

### This document's own three

The [first](https://github.com/radiusred/gh-codecrew/issues/351#issuecomment-5793927587)
fixes the counting instant after this task's records are written, the method
[M16's record](16-bedding-in-the-venue-seam-and-the-reference.md) used.
**Rejected:** counting before these records exist, "which makes the header
unreproducible for anyone who runs the same query afterwards".

The [second](https://github.com/radiusred/gh-codecrew/issues/351#issuecomment-5793927820)
draws the line for which captures count as the milestone's at the instant #345
was created, so #344 — filed six minutes and thirty-three seconds earlier and
cited by the Goal as existing — is the Goal's capture and not one filed during
M17. **Trade-off:** "the M16 record counted #324, filed two seconds before its
milestone, among that milestone's captures; this record draws the line at the
milestone's creation instant instead, so the two records do not count the same
way." The requirement-outcomes section applies it.

The [third](https://github.com/radiusred/gh-codecrew/issues/351#issuecomment-5793928043)
discharges the refresh obligation by reading `README.md` and
`docs/introduction.md` against what M17 delivered and editing neither: the
milestone's work is a blog post and one sentence in the spoke, and "no verb,
refusal code, release or hub document changed". **Rejected:** "adding a line
about the post to the README or the introduction — that is news, not a claim
about what exists and works, and those pages carry the latter." The README
already links the blog; `**Shipped:** v2.0.1` is still true.

## Deviations

Two Deviation records: one on www#37, by the doc-synthesizer seat, and one on
this document's own task, by the operator.

**[The post was drafted before its task existed](https://github.com/radiusred/codecrew-www/issues/37#issuecomment-5740858343).**
"the post was drafted on the operator's machine before this task existed — the
operator asked for an appraisal of the announcement, then for it as a post,
then chose the vehicle (M17 and this task). The draft is committed unchanged
on the linked branch; the plan above describes it after the fact." **Why:**
"the protocol wants the plan on the issue before the first commit (SPEC §4).
No commit had been made — the draft was an untracked file — so the letter of
the rule holds, but the deciding-before-doing it protects did not: the plan
was written knowing the draft. Recorded rather than hidden." The timestamps
agree with it. The milestone issue was created at 09:48:01Z, the task thirteen
seconds later; the seat's plan edit is timestamped 09:49:03Z, the Deviation
09:49:04Z, `task start`'s comment 09:49:07Z, and the first commit's author date
09:49:41Z — one minute and forty seconds from the milestone to a 46-line post.

**[The implementer seat's model changes, and rides on this record's pull request](https://github.com/radiusred/gh-codecrew/issues/351#issuecomment-5793895499).**
The operator's Deviation at 11:20:15Z, seven seconds after #351 was created:
"this record task's PR also carries a one-line change to `.codecrew/config.yml`
outside the milestone's scope: the implementer seat's `model` moves from
`claude-fable-5` to `claude-opus-5-5`, by the operator's decision on
2026-09-23." **Why:** "the operator chose this vehicle over a separate `chore:`
PR (a light path for such changes is only proposed so far, in #343). M17's only
other open work was in the spoke `radiusred/codecrew-www`, so this is the
milestone's one remaining hub PR. The change is in its own commit so the
reviewer can judge it apart from the record. The two M17 spoke tasks after the
change (www#39 and www#42) were already dispatched under Opus 5.5." Both
tasks' commits carry the trailer `Co-Authored-By: Claude Opus 5.5 (1M context)`,
so the routing table named `claude-fable-5` for the implementer seat while the
seat ran under Opus 5.5 for both of its M17 tasks, and this pull request's
second commit brings the table in line. The commit changes that one value and
nothing else in the file.
[#343](https://github.com/radiusred/gh-codecrew/issues/343), the light-path
capture, was filed on 2026-09-11 and is open.

## The gates

**#345's Gates section was left as the scaffold's placeholder** — the template
line "_What "done" means beyond CI: e2e suites, manual UAT, sign-offs._",
unedited. [M16's record](16-bedding-in-the-venue-seam-and-the-reference.md#the-gates)
says the same of #325, and
[M15's](15-v2-0-1-what-the-fleet-migration-taught.md#the-gates) of #305.

**No gate was raised.** No task used `gh codecrew checkpoint`, no issue carries
a `cc:needs-decision` label, and `milestone close 17 --dry-run` reports `gate
no gate raised: ok`. All three delivery tasks wrote a resolved Ask-the-human
section, each resolved to none, and each says why: www#37's "None: the
operator chose the vehicle (M17, a task, a PR) on 2026-09-19 before this task
was opened"; www#39's "The wording is house-voice copy the reviewer judges; the
operator's M17 plan already fixed the position and its placement"; and www#42's
bare "None."

What gated the work: CI on every pull request — the spoke's `Lint commit
messages` and `Tests` — an independent approval on each of the three, then the
QA verdicts. No requirement added a gate of its own, and no release gate exists
because there is no release.

## The review rounds

Three pull requests, **four review submissions**, every one by the reviewer
role holder (`radiusred-checky[bot]`), an App distinct from both authoring
seats. **One change request**, resolved in the next round. **Three approvals
were submitted and all three stand.** One force-push landed in the milestone,
on www#38 at 09:51:53Z, before any review; the observations below carry it.

- **PR [radiusred/codecrew-www#38](https://github.com/radiusred/codecrew-www/pull/38)**
  — two rounds. The
  [change request](https://github.com/radiusred/codecrew-www/pull/38#pullrequestreview-5255571103)
  at 2026-09-19T11:14:43Z opens on the requirement's own words: "M17-R1 and the
  task plan both require announcement quotes to be verbatim, and two quoted
  phrases are paraphrases." The first is the memory sentence, quoted as "remembers
  project details, decisions, and communication preferences," — "that exact
  phrase does not appear in Anthropic's post"; the second is the rollout,
  quoted as "later", where the source says "come after that". The same round
  lists what it checked and passed, and one item in that list is the fragment
  QA later failed: "the coordinator sentence and "reviews the outputs" are
  verbatim". It also verified the post's other factual claims against the hub —
  "the hub first commit is 2026-08-20; reviewer and QA are Codex-routed in the
  hub config; M7 is the coordinator-seat record" — and "every hub link in the
  post returned HTTP 200". The fix commit
  ([c338cf1](https://github.com/radiusred/codecrew-www/commit/c338cf1)) is
  authored one minute and fifty-one seconds after the review. The
  [round-two approval](https://github.com/radiusred/codecrew-www/pull/38#pullrequestreview-5255587274)
  at 11:20:25Z says it "Verified only the requested follow-up": the four
  repaired fragments are verbatim, and `git diff 0fc7d25..c338cf1` touches only
  the post, "with exactly the memory sentence and the Team/Enterprise plan-tier
  sentence changed".
- **PR [radiusred/codecrew-www#40](https://github.com/radiusred/codecrew-www/pull/40)**
  — approved
  [first round](https://github.com/radiusred/codecrew-www/pull/40#pullrequestreview-5290085878)
  at 2026-09-23T11:00:55Z, three minutes and forty-nine seconds after the pull
  request opened, opening "No findings." It reads the built site rather than
  the Markdown — "`site/index.html` renders the hero link as
  `blog/posts/2026-09-19-the-coordinator-is-a-commodity-now/`, and that target
  exists" — and judges the voice against the requirement: the sentence "does not
  position CodeCrew against a named tool, and stays bright, informed, and
  inform-first rather than hype or sales copy." It also rules on the regenerated
  blog metadata the pull request body flagged: it "does not belong in this PR;
  a capture is the right vehicle for that state."
- **PR [radiusred/codecrew-www#43](https://github.com/radiusred/codecrew-www/pull/43)**
  — approved
  [first round](https://github.com/radiusred/codecrew-www/pull/43#pullrequestreview-5290221273)
  at 11:14:18Z. It checks the diff before the body — "one file … with one
  changed line; `git diff --word-diff` shows only the full stop moving" — and
  then repeats the Decision's audit rather than taking it: every quotation
  attributed to the announcement is present in the source "exactly as the post
  now quotes them", "The source does not contain `"reviews the outputs."`", and
  "The implementer's classification also holds" for the three quoted strings
  that are the post's own prose. It accepts the absence of a test on the
  merits: "the repo tests cover generation/layout and do not assert individual
  blog prose, and the required check for this one-character correction is the
  independent source quotation audit plus strict build."

**What the one change request was about.** Two phrases inside quotation marks
that paraphrased the source, on a requirement whose text is "quoting the
announcement verbatim where it is quoted". The same round checked a third
fragment and passed it with a full stop inside the marks that the source does
not have there; the QA section below is what happened next. Every round ran
the spoke's deploy-equivalent recipe — `uv run pytest -q` (70 passed each
time), `sync_docs.py`, `main.py` and a strict Zensical build.

## QA: two rounds, one remedy, and the verdict that supersedes

The qa role holder (`radiusred-testy[bot]`) wrote two comments on #345 and
three verdicts. The first round found M17-R1 not satisfied; a remedy task was
opened, reviewed and merged; the second round re-verified M17-R1 and says in
its own words which comment it supersedes.
[M16's record](16-bedding-in-the-venue-seam-and-the-reference.md#qa-one-round-six-verdicts)
and [M15's](15-v2-0-1-what-the-fleet-migration-taught.md#qa-one-round-seven-verdicts)
each describe one round with no supersession;
[M14's](14-adoption-tidy-and-the-v2-0-0-release.md#qa-two-rounds-and-the-requirement-that-was-answered-rather-than-remedied)
two rounds and a requirement answered by a Decision rather than a remedy task;
and [M13's](13-protocol-2-0-the-codecrew-layout-and-what-rides-with-it.md#qa-three-rounds-and-the-requirement-that-took-all-three)
three rounds and a remedy task.

**Round one, at 2026-09-23T11:07:46Z** —
[the comment](https://github.com/radiusred/gh-codecrew/issues/345#issuecomment-5793731346).
Its M17-R2 verdict is `satisfied`, on the live served HTML at the front page
containing the sentence "inside the existing `cc-hero__sub` paragraph with the
post link resolving HTTP 200", and on a diff probe of the merge commit that
"changes only that paragraph, the before/after section structure is
identical". Its M17-R1 verdict is `not satisfied`. It first says what the clean
build does and does not prove — "that proves the site builds but assumes the
post's factual and quote checks are correct" — then passes every other clause
of the requirement and names the one it fails: "the post quotes `reviews the
outputs.` with a period, while the source only has that phrase inside `scopes
the request, delegates the work, coordinates parallel threads, reviews the
outputs, and assembles the finished result.`; finding filed at
radiusred/codecrew-www#37." The finding is
[a comment on the closed task](https://github.com/radiusred/codecrew-www/issues/37#issuecomment-5793721573)
at 11:07:01Z, forty-five seconds before the verdict, with the source fragment
quoted and the two ways out stated: "either preserve the source punctuation
exactly, or stop presenting the changed-punctuation fragment as a quotation."

**The remedy.** The operator opened
[radiusred/codecrew-www#42](https://github.com/radiusred/codecrew-www/issues/42)
at 11:08:15Z, twenty-nine seconds after the verdict, with the remedy in its
Goal: "Move the full stop outside the closing quotation mark so the quoted text
is character-exact; nothing else in the post changes". Its plan widens the
check to every quotation in the post. The pull request merged at 11:14:49Z,
seven minutes and three seconds after the verdict.

**Round two, at 11:19:10Z** —
[the comment](https://github.com/radiusred/gh-codecrew/issues/345#issuecomment-5793882104),
eleven minutes and twenty-four seconds after round one. It is M17-R1 alone, and
it names what it replaces: "This re-verification supersedes comment 5793731346
for R1." It checks every announcement-attributed quotation "character-exact"
against the source — the same six strings the reviewer checked on www#43 — and
reports that "the old bad fragment `reviews the outputs.` is not in the source
and the live post now renders `"reviews the outputs".` with punctuation outside
the quote". It then re-checks the rest of the requirement against the live post
rather than carrying it over, quoting the post's own phrases for each clause:
"`the running of it`", the four claims by their headings, the hedge "`As far as
we can tell from the announcement`", the link to #344, and the GSD essay framed
"as one operator's origin story rather than a rival product claim". It closes
on the capture it filed: "The known generated metadata drift is real but
already filed outside this requirement as radiusred/codecrew-www#41."

**The standing verdicts are round two's for M17-R1 and round one's for
M17-R2.** Both are `satisfied`.

## Requirement outcomes

The status column is the **standing** verdict word as the qa role holder wrote
it, verbatim and unqualified.

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M17-R1 — the codecrew.works post on the Projects announcement: written in the doc-synthesizer's voice, quoting the announcement verbatim where it is quoted, naming what Projects does that CodeCrew never owned (the running) and the four protocol claims it leaves untouched (the open record, no self-review, any model in a seat, the human gate as a protocol event), hedging the identity claim as unverified, linking #344 as the experiment; the essay framed as one operator's experience and not positioned against GSD; the strict site build clean | [radiusred/codecrew-www#37](https://github.com/radiusred/codecrew-www/issues/37) / PR [radiusred/codecrew-www#38](https://github.com/radiusred/codecrew-www/pull/38) and the remedy [radiusred/codecrew-www#42](https://github.com/radiusred/codecrew-www/issues/42) / PR [radiusred/codecrew-www#43](https://github.com/radiusred/codecrew-www/pull/43) | `satisfied` | Standing verdict [5793882104](https://github.com/radiusred/gh-codecrew/issues/345#issuecomment-5793882104), which supersedes the `not satisfied` of [5793731346](https://github.com/radiusred/gh-codecrew/issues/345#issuecomment-5793731346): a full stop inside the marks on "reviews the outputs" that the source does not have there. Both tasks closed. The post took two review rounds over two paraphrased quotations; the remedy one. One Deviation (the draft preceded the task) and one Decision (the remedy's scope and the quotations it checked) |
| M17-R2 — the codecrew.works front page carries the post's one-line position — the coordinator is a commodity, the record and the separation of duties are not — as a sentence in the existing intro, not a new section; the sync table and everything else on the page unchanged | [radiusred/codecrew-www#39](https://github.com/radiusred/codecrew-www/issues/39) / PR [radiusred/codecrew-www#40](https://github.com/radiusred/codecrew-www/pull/40) | `satisfied` | Task closed; verdicted on the live front page and on a diff of the merge commit that changes only the `cc-hero__sub` paragraph. Approved first round with no findings. Two Decisions: the on-site link to the post, and the hero-test assertion the dispatch brief had not asked for |

**Captures adopted and closed by the merges — none.** No requirement adopts a
capture and no task carries a `## Adopts` section.

**The capture the Goal was written around** is
[#344](https://github.com/radiusred/gh-codecrew/issues/344), the
Projects-as-orchestrator experiment, filed by the operator at 09:41:28Z on
2026-09-19, six minutes and thirty-three seconds before the milestone issue.
It is open, and it is findings, not a requirement: "the deliverable is
findings. Any behaviour change it suggests is a separate capture, parked with
the rest until the operator lifts the bedding-in rule." Under this document's
[second Decision](https://github.com/radiusred/gh-codecrew/issues/351#issuecomment-5793927820)
it is not counted below.

**Captures filed during the milestone and left open — three**, two in this hub
and one in the spoke.
[#346](https://github.com/radiusred/gh-codecrew/issues/346), filed at 11:15:52Z
on 2026-09-19, is the operator's, "reading the page on codecrew.works":
`docs/platform-interop.md` still pins itself to "the installed release"
v1.0.3; an addendum the same day adds that three of its sections describe the
Paperclip recipe as it stood before Paperclip v2026.916.0.
[#347](https://github.com/radiusred/gh-codecrew/issues/347), filed at 11:59:32Z,
is the sibling experiment for the other orchestrator — whether a crew still
holds typed seats when Paperclip agents act as "delegates of a responsible
person" — and the operator's comment on #344 links the two: "Findings on each
stay on its own issue."
[radiusred/codecrew-www#41](https://github.com/radiusred/codecrew-www/issues/41)
is the QA seat's, filed at 2026-09-23T11:07:11Z: running the deploy-equivalent
generator on the spoke's `main` regenerates `docs/blog/archive.md`,
`docs/blog/atom.xml` and `zensical.toml` with the M17-R1 post, which `main`
does not carry committed.

**Captures filed in the window from other projects — three**, all the
operator's, none of them M17 work.
[#348](https://github.com/radiusred/gh-codecrew/issues/348) (2026-09-20) comes
from closing `davison/md-notes` M8: a requirement struck by decision has no
verdict word, so `milestone close` refuses on it.
[#349](https://github.com/radiusred/gh-codecrew/issues/349) and
[#350](https://github.com/radiusred/gh-codecrew/issues/350) (2026-09-23, four
seconds apart) come from closing M1 in `davison/davison.github.io`: the
milestone-close ceremony costs the same for a three-line milestone as for a
large one, and the doc-synthesizer's README-and-introduction obligation assumes
the hub is a product repository. Each carries a coordination-layer appraisal
dated 2026-09-23 that adopts it "into the next protocol milestone". None is
adopted by any M17 task, and all three are open.

## Protocol-discipline observations

Six things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **The verdict loop ran end to end, and supersession was written down by the
  seat.** A `not satisfied` verdict, a finding on the delivering task, a remedy
  task opened twenty-nine seconds later, an independent review, a merge, and a
  re-verification that names the comment it supersedes by number — in eleven
  minutes and twenty-four seconds. The protocol's rule already makes the latest
  QA comment carrying a verdict for a requirement the standing one (M13-R6, in
  `tracker.ParseVerdicts`), so round two would have superseded round one without
  saying so; it says so anyway, by comment number. The record a reader needs to
  follow the loop is all on #345, www#37 and www#42.
- **A review and a QA round checked the same fragment against the same source
  and came to different answers.** The reviewer's round one on www#38 lists
  "reviews the outputs" among the quotations it found verbatim; QA found a full
  stop inside the marks that the source does not have there. The difference is
  one character, and the requirement's word is "verbatim". What the remedy then
  did is visible in three places: the task's Decision, the reviewer's approval
  of www#43 and QA's round two each list the same six strings and each check
  them character for character, punctuation inside the marks included.
- **A seat's commit was rejected by the commit lint and repaired by a
  force-push, and no comment records it.** The first commit on www#38
  ([e1a688d](https://github.com/radiusred/codecrew-www/commit/e1a688d)) carried
  a commit-message body written as single long lines, the longest 731
  characters; its `Lint commit messages` check concluded `failure` at 09:50:03Z.
  The force-push at 09:51:53Z replaced it with
  [0fc7d25](https://github.com/radiusred/codecrew-www/commit/0fc7d25), whose
  tree is identical and whose body is wrapped under 100 columns. The pull
  request opened on the failing commit; the review ran on the repaired one.
  *Undocumented, inferred from the two commits and the check runs.*
- **The routing table and the practice disagreed for the whole of M17's
  implementer work, and the table is advisory by design.**
  `.codecrew/config.yml` names `claude-fable-5` for the implementer seat; both
  M17 implementer tasks ran under Claude Opus 5.5, which the operator's
  Deviation records and the commit trailers confirm. The file's own header
  calls the table "Advisory role routing, read by whoever dispatches agents",
  and nothing refused. This pull request's second
  commit brings the table in line with what was dispatched. The doc-synthesizer
  seat has no `model` key; the post's two commits carry the trailer
  `Co-Authored-By: Claude Fable 5.1`.
- **Three seats saw the unregenerated blog metadata, and the third filed it.**
  PR www#40's body tells the reviewer that running `main.py` "regenerates
  `zensical.toml`, `docs/blog/archive.md` and `docs/blog/atom.xml` with the
  M17-R1 post, which main does not have committed"; the reviewer's approval
  says "a capture is the right vehicle for that state"; neither filed one, and
  QA filed www#41 five minutes and thirty-five seconds after the merge. The deploy workflow regenerates
  the files, so the live site is unaffected; what the capture records is that
  the spoke's `main` and its deployed output differ.
- **The trail cites nothing on github.com, so the evidence verb checked only
  what lies outside the protocol's reach.** Every reference between the
  milestone's issues is a bare `#N` or `owner/repo#N`, and the three links
  `milestone evidence 17` resolves are the announcement, the live post and the
  front page. [M16's record](16-bedding-in-the-venue-seam-and-the-reference.md)
  reports six links, five of them in its own task's records. A citation the
  verb cannot see is one it cannot check.

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **The conversation that scoped the milestone is on the record only through
  #345's goal**, "Decided with the operator 2026-09-19", and through www#37's
  Deviation, which says the operator "asked for an appraisal of the
  announcement, then for it as a post, then chose the vehicle". The appraisal
  itself is not on the trail; #344 names it as its source. *The reasoning is
  recorded; the deliberation is not*, as
  [M16's record](16-bedding-in-the-venue-seam-and-the-reference.md#what-the-record-does-not-contain)
  and [M15's](15-v2-0-1-what-the-fleet-migration-taught.md#what-the-record-does-not-contain)
  each say of their own openings.
- **Nothing on the trail says why M17-R2's task was opened four days after
  M17-R1 merged.** The Goal names the front-page line as something the post
  "promises", and the post itself says the week's changes include "a line on
  this site's front page". www#39 was created on 2026-09-23, and no comment
  between 2026-09-19 and then mentions it.
- **Which harness ran the reviewer and QA rounds is not on the trail.** No
  review, QA comment or Deviation names one. The routing table routes both
  seats to `codex` with `gpt-5.5`, and the coordination layer's account for
  this record is that all four review submissions and both QA rounds ran
  there. Unlike
  [M16](16-bedding-in-the-venue-seam-and-the-reference.md#deviations), no
  Deviation says any round ran elsewhere. The implementer's harness is on the
  trail through the operator's Deviation and the commit trailers; the
  doc-synthesizer's through the post's commit trailers only.
- **The force-push on www#38 has no comment.** The observation above
  reconstructs it from the check runs and the two commits; the seat recorded
  neither the failure nor the repair.
- **The post's content is not verified by any test.** QA's round one says it
  plainly — the clean build "proves the site builds but assumes the post's
  factual and quote checks are correct" — and the reviewer on www#43 says the
  spoke's tests "do not assert individual blog prose". The quotation audit that
  now stands is three seats' readings of the announcement as it was served on
  2026-09-23, and www#42's Decision names the consequence: "a later edit to the
  announcement would not be caught."
- **#344 and #347 have no findings yet.** Both experiments are open and both
  are gated on something outside the project — #344 on "a Projects beta seat",
  #347 on a Paperclip on v2026.916.0 or later. The post says of #344 "the
  finding will go on an issue like everything else"; at the counting instant it
  has one comment, the operator's link to its sibling.
