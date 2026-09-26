# M20: The 2.1 post

Tracking issue: [#384](https://github.com/radiusred/gh-codecrew/issues/384) ·
Synthesized 2026-09-26 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #384's **one** requirement as opened, **none added,
none struck and none amended** — the issue's `userContentEdits` query returns
a total count of zero; its Gates section, left as the scaffold's placeholder;
its **two comments**, the qa role holder's **one verdict** and this record's
**one Deviation**; the **one** task issue, in the spoke
([radiusred/codecrew-www#66](https://github.com/radiusred/codecrew-www/issues/66));
its **one merged pull request** and the **one commit** on it; the **two
Decision records and one Deviation record** across the milestone issue and
its task, the two Decisions on the task and the Deviation on the milestone
issue; the **one review submission**, an approval that stands; **no**
adopted backlog captures; and **no** captures filed in the milestone's window.
The house form is [M19's document](19-the-front-page-tells-the-crew-story.md),
with [M17's](17-the-frontier-builds-the-coordinator.md) the closest in shape;
the record standard is the one the reviewer set on
[PR #322](https://github.com/radiusred/gh-codecrew/pull/322) and re-applied
since, so every count is re-derived from the source — records by the rule
`tracker.ExtractRecords` applies (a paragraph-initial `**Decision:**`,
`**Deviation:**` or `**Gate resolved:**` label, bare or parenthetically
qualified), the review state from the pull request's reviews API and
timeline, the commit from the pull request's commits API and the spoke's
compare API, and the requirement history from the milestone issue's
`userContentEdits`.

**Counted at 2026-09-26T11:14:30Z, against a trail frozen from this
document's dispatch.** The last post on the milestone's issues before that
instant is this record's own
[Deviation](https://github.com/radiusred/gh-codecrew/issues/384#issuecomment-5845774859)
at 11:14:19Z, and it is inside every total above; the one before it is the qa
role holder's verdict at 11:11:56Z.

**This record has no task.** It is delivered as a housekeeping pull request
under SPEC §4's milestone-record clause, as
[M18's](18-the-protocol-scales-down.md) and
[M19's](19-the-front-page-tells-the-crew-story.md) were: the milestone issue
is its charter and its commit references #384.

**The whole milestone happened in the spoke.** `main` in this repository did
not move between M19's record
([13eebb4](https://github.com/radiusred/gh-codecrew/commit/13eebb4), committed
2026-09-26T10:17:59Z) and this document's pull request. QA's verdict reports
the same for the protocol's own files: "`git -C repo log --oneline
13eebb4..origin/main -- SPEC.md CLI.md internal cmd .codecrew/roles` printed
nothing". `main` in
[radiusred/codecrew-www](https://github.com/radiusred/codecrew-www) moved by
exactly **one commit** from M19's last spoke merge
([2230e98](https://github.com/radiusred/codecrew-www/commit/2230e98)):
[e01e267](https://github.com/radiusred/codecrew-www/commit/e01e267), the
post, authored by `radiusred-wordy[bot]` and carrying the trailer
`Co-Authored-By: Claude Opus 5.5 (1M context)`.

**The trail is checked by the verb.** The installed extension reports `v2.1.0
(protocol 2.1)`; M20 cut no release. Run from this repository at the counting
instant:

```
$ gh codecrew milestone evidence 20
requirements counted: M20-R1 (1)
all 0 cited links resolve across 2 issues — evidence is reachable
```

Two issues: the milestone and its one sub-issue. No links: the Decisions, the
verdict and this record's Deviation refer to issues as bare `#N` or
`owner/repo#N`, which the verb does not read as a link, and name the post's
sources without linking them.

`gh codecrew milestone close 20 --dry-run`, as it read at 11:14:11Z, before
this pull request existed, passes five gates and stops at the document, as
M19's did:

```
gate milestone open: ok
gate no gate raised: ok
gate tasks closed: ok
gate requirements declared: ok
gate QA verdicts: ok
gate milestone document: refused[DOC_MISSING]: docs/milestones/20-*.md not on the default branch of radiusred/gh-codecrew — dispatch the doc-synthesizer: it delivers the record as a housekeeping PR with no task (SPEC §4), reviewed as the reviewer seat's routing requires (in pure solo, the operator confirms on the PR), and its author rebase-merges it; then rerun
dry run: nothing written — the live verb stops at the first refusal above
```

The dry run prints no raw material, so every record below was read from its
issue, pull request, review, commit or timeline directly. The prose is wrapped
at this file's normal width; the requirement-outcomes table's row, the fenced
output above and a handful of single links are not, because a newline ends a
Markdown table row and breaking a link breaks it.

**This PR adds the M20 ROADMAP row and changes nothing else in the hub.**
`ROADMAP.md` carried no M20 row before it. The section on the front door
below says why neither front-door document changes.

## Goal and outcome

#384 was opened by the operator's account at 10:33:52Z on 2026-09-26,
fourteen minutes and two seconds after
[`milestone close 19`](https://github.com/radiusred/gh-codecrew/issues/382#issuecomment-5845426703)
closed M19. Its goal names the gap: "Protocol 2.1 shipped in v2.1.0 on
2026-09-24 (M18, #353) and has no public account yet; the 2.0 release had one
(the M15 post)." It sets the scope in one sentence: "This milestone is that
account and nothing more: a brief codecrew.works post on what 2.1 changes for
an adopter." It carries the rule forward — "The bedding-in rule stands for 2.1
— no verb changes its arguments, output, writes or refusal codes, and SPEC.md
is untouched." — and puts the experiment outside the milestone: "The
Projects-as-orchestrator experiment (#344) stays a capture, gated on beta
access." The provenance line is "Decided with the operator 2026-09-26."

The one requirement, as the body carries it:

- **M20-R1** — "the codecrew.works post on protocol 2.1: brief (the 2.0
  post's shape, well under its length), in the doc-synthesizer's voice; names
  what 2.1 adds — the struck requirement and milestone strike, the light path
  for housekeeping and the milestone record, roles sync, the version check in
  the agent entry point — and what an adopter runs to take it up; every claim
  cited by link to the hub's CHANGELOG, SPEC, CLI.md or the M18 record,
  anything in quotation marks character-exact from its source; the strict
  site build clean".

The requirement adopts no capture, and the task carries no `## Adopts`
section.

The one task:

- **[radiusred/codecrew-www#66](https://github.com/radiusred/codecrew-www/issues/66)
  / PR [radiusred/codecrew-www#67](https://github.com/radiusred/codecrew-www/pull/67)**
  — M20-R1, delivered by the doc-synthesizer seat (`radiusred-wordy[bot]`),
  which the requirement's "in the doc-synthesizer's voice" names. One commit,
  approved first round, merged 2026-09-26T11:00:48Z
  ([e01e267](https://github.com/radiusred/codecrew-www/commit/e01e267)). Two
  Decisions, no Deviations.

**The task ran in order, by the timestamps.** The operator's account created
www#66 seven seconds after the milestone. Its body carried a Goal,
`## Requirements` naming M20-R1, and the scaffold's placeholders under Plan
and Ask-the-human points; the text under both is the seat's edit at
10:35:49Z, six seconds before the `**Started by**` record at 10:35:55Z. The
two Decisions were posted at 10:38:14Z and 10:38:15Z, and the commit is
authored at 10:38:26Z. The pull request opened at 10:38:43Z; its
`Lint commit messages` and `Tests` checks were green by 10:39:37Z; the approval
came at 10:44:51Z; the rebase merge, the task's closure by PR #67 and the head
branch's deletion are all the seat's, at 11:00:48Z–11:00:50Z. The task
adopted no capture, so no closing comment was due on any. From the milestone
opening to the merge: **twenty-six minutes and fifty-six seconds**; to the QA
verdict: **thirty-eight minutes and four seconds**.

What exists afterwards that did not before. codecrew.works carries a third
blog post, beside those of 2026-09-07 and 2026-09-19: "Protocol 2.1: CodeCrew
Learns to Scale Down", at `/blog/posts/2026-09-26-protocol-2-1-scales-down/`,
dated 2026-09-26; the blog's navigation, archive and Atom feed list it. The
pull request's four files are the post and its entries in `zensical.toml`,
`docs/blog/archive.md` and `docs/blog/atom.xml` — the last with the feed's
`<updated>` changed, as the first Decision says — 44 lines added and one
removed. Nothing in this repository's code, `SPEC.md`, `CLI.md` or `docs/`
changed.

## Decisions

Two Decision records, both on www#66 and both by the doc-synthesizer seat
(`radiusred-wordy[bot]`). None is on the milestone issue, none carries a
parenthetical qualifier and none is a gate resolution.

**[Only this post's generated entries](https://github.com/radiusred/codecrew-www/issues/66#issuecomment-5845542848).**
The seat commits "only this post's own entries in the three files `main.py`
regenerates", and leaves out "the 2026-09-19 post's missing nav, archive and
feed entries that the regeneration also drags in". **Trade-off:** "the
committed generated files stay one post behind what `main.py` would write
until #41 is taken up; the scheduled build regenerates them anyway, so the
live site is unaffected, and the strict build passes on the committed state
as well as the regenerated one." **Rejected:** "committing the full
regeneration (it is #41's fix, not this task's); restoring all three files and
shipping the post with no nav, archive or feed entry of its own." The Plan had
named the Decision before the work — "leave out any older post's metadata
that is www#41's, and record that as a Decision." The reviewer's approval
cites it: "The omitted 2026-09-19 metadata is deliberately left for www#41 by
the recorded Decision".

**[GitHub citations, and the title](https://github.com/radiusred/codecrew-www/issues/66#issuecomment-5845542939).**
The seat cites "the hub's sources by their GitHub URLs (CHANGELOG, SPEC.md and
CLI.md with section anchors, the M18 record, the radiusred/gh-codecrew#353
Decisions, the captures), not the synced copies under codecrew.works/docs/",
and titles the post "Protocol 2.1: CodeCrew Learns to Scale Down", "after
M18's own title". **Trade-off:** "readers leave the site for GitHub on every
citation; in exchange each link is the exact source M20-R1 names, the
CHANGELOG and the milestone record (which are not synced) link the same way as
the rest, and a reviewer can diff every quotation against the file it cites."
**Rejected:** "site-relative links to `docs/spec.md` and `docs/cli.md` as the
2.0 post uses (the CHANGELOG, M18 record and issues would still have to go to
GitHub, so citations would be split two ways)."

Both Decisions were edited within eight seconds of being posted, before the
commit. The edit histories show what changed and nothing else: the first's
two `www#41` became `#41`, and the second's `#353` became
`radiusred/gh-codecrew#353` — each reference rewritten to resolve from the
spoke, where the comments live.

## Deviations

One Deviation record, on the milestone issue and this record's own. The task
issue and the pull request carry none.

**[The record takes harness facts from its dispatch
brief](https://github.com/radiusred/gh-codecrew/issues/384#issuecomment-5845774859).**
Posted at 11:14:19Z. The brief states that the post was written by a Claude
Opus 5.5 sub-agent, that the review ran in Codex under `gpt-5.6-terra` and
that QA ran in Codex under `gpt-5.6-sol`. The comment sets out what the trail
supports — the routing table's two Codex models, the doc-synthesizer's
missing `model` key, and the one commit's Opus 5.5 trailer — and that "No
review or QA comment names its own harness or model." **Why:** "the brief is
not on the trail, so this comment puts its account there, and the record cites
this comment rather than the brief — as the M19 record did on #382."

## The gates

**#384's Gates section was left as the scaffold's placeholder** — "_What
"done" means beyond CI: e2e suites, manual UAT, sign-offs._", unedited, as
[M17's record](17-the-frontier-builds-the-coordinator.md#the-gates) says of
#345.

**No gate was raised.** The task did not use `gh codecrew checkpoint`, no
issue carries `cc:needs-decision`, and the dry run reports `gate no gate
raised: ok`. www#66's Ask-the-human section reads "None — the scope, sources
and voice are fixed by M20-R1 and the dispatch brief."

What gated the work: the spoke's CI on the pull request — `Lint commit
messages` and `Tests`, both `SUCCESS` — the reviewer seat holder's approval,
and the QA verdict. No release gate exists because there is no release.

## The review round

One pull request, **one review submission**, by the reviewer role holder
(`radiusred-checky[bot]`), an App distinct from the authoring seat.

- **PR [radiusred/codecrew-www#67](https://github.com/radiusred/codecrew-www/pull/67)**
  — approved
  [first round](https://github.com/radiusred/codecrew-www/pull/67#pullrequestreview-5325647420)
  at 10:44:51Z, six minutes and eight seconds after the pull request opened,
  on the head it merged at
  ([0db4188](https://github.com/radiusred/codecrew-www/commit/0db4188), rebased
  to e01e267). It opens "Approved — M20-R1 is satisfied." and says it
  "Reviewed the diff before the PR body." It reports that "I diffed all nine
  quotations (including the changelog block quote) character-for-character
  against their fetched sources" and that "All 24 links resolve, including
  each GitHub anchor". It rules on the two citations the requirement's list
  does not name: "The extra AGENTS.md and M20 citations clarify the
  coordinator brief; they do not depart from R1 in substance." It checks the
  process as well as the post — "The task's original Plan predates commit
  0db4188; its two Decisions are present. The PR's closing references contain
  #66 alone." — and reproduces the build: "`uv run pytest -q` (89 passed),
  `sync_docs.py`, `main.py`, and `zensical build --clean --strict` (No issues
  found)."

No change request was submitted and no review was dismissed. No force-push
landed on the pull request.

## QA: one round, one verdict

The qa role holder (`radiusred-testy[bot]`) wrote
[one comment](https://github.com/radiusred/gh-codecrew/issues/384#issuecomment-5845760572)
on #384, at 11:11:56Z, eleven minutes and eight seconds after the merge. It
is one verdict: **M20-R1 — satisfied.**

It first bounds what the tests prove: the suite covers "generic newest-post
title/date/own-description rendering and nav/archive/feed generation, while
assuming the copy, citations, deployment, responsive layout and theme; I
checked those seams independently." It then takes the requirement clause by
clause. Brevity: "The body is 895 words versus the 2.0 post's 1,772". The four
additions and the adopter's sequence are each named. The sources: against the
CHANGELOG, SPEC §4, §7 and §10, CLI.md, the M18 record, #353's Decisions,
`.codecrew/AGENTS.md`, #384 and `gh release view v2.1.0`, "all eight quoted
strings and the block quote character-exact including punctuation". The live
site: "all 25 authored links returned 200 and both fragments landed", the
served `<main>` "matched the local main build byte-for-byte", and the deploy
run for e01e267 "was the live source". The layout, in Playwright "at 375px and
1280px in light and dark", with no horizontal overflow. And the voice:
"bright, informed, journalistic and openly agent-staffed, without hype,
unsourced superlatives, or any claim that CodeCrew dispatches, schedules or
hosts agents."

## Requirement outcomes

The status column is the verdict word as the qa role holder wrote it,
verbatim and unqualified, in comment
[5845760572](https://github.com/radiusred/gh-codecrew/issues/384#issuecomment-5845760572).

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M20-R1 — the codecrew.works post on protocol 2.1: brief (the 2.0 post's shape, well under its length), in the doc-synthesizer's voice; names what 2.1 adds — the struck requirement and milestone strike, the light path for housekeeping and the milestone record, roles sync, the version check in the agent entry point — and what an adopter runs to take it up; every claim cited by link to the hub's CHANGELOG, SPEC, CLI.md or the M18 record, anything in quotation marks character-exact from its source; the strict site build clean | [radiusred/codecrew-www#66](https://github.com/radiusred/codecrew-www/issues/66) / PR [radiusred/codecrew-www#67](https://github.com/radiusred/codecrew-www/pull/67) | `satisfied` | Task closed. One commit, approved first round. Two Decisions (the post's generated entries only, leaving www#41's; GitHub citations and the title). Verdicted on the merged build and the live post |

**Captures adopted and closed by the merge — none.** The requirement adopts
no capture and the task carries no `## Adopts` section.

**The capture the Goal names** is
[#344](https://github.com/radiusred/gh-codecrew/issues/344), the
Projects-as-orchestrator experiment, which "stays a capture, gated on beta
access". At the counting instant it is open with one comment, the same count
[M17's record](17-the-frontier-builds-the-coordinator.md#what-the-record-does-not-contain)
gives. www#41, which the first Decision leaves the older post's metadata to,
was filed in M17 and is open.

**Captures filed during the milestone — none.** From #384's creation to the
counting instant, the only issues or pull requests created in this repository
and in radiusred/codecrew-www are #384, www#66 and www#67; an organisation-wide
issue search over the same window, as this record's token sees the
organisation, returns the same three.

## The front door

The doc-synthesizer contract sends this hub's front-door list to its
`.local.md`: "the README's proof points, and the introduction's release, verbs
and refusal codes". M20 changed none of those — no release, no verb, no
refusal code, and no hub commit. `docs/introduction.md`'s `**Shipped:**` line
reads v2.1.0, which is still the shipped release. Both documents mention
codecrew.works: the README calls it "the marketing and introduction site",
which "carries the blog", and the introduction lists it as "the introduction
site". Neither lists individual posts, so a third post makes neither claim
false, and this PR changes neither document. SPEC §4 lists the record's diff
as the document, its ROADMAP row and the front-door claims the milestone
changed, "and nothing else", and says the housekeeping path carries "no
record entry", so this PR adds no CHANGELOG entry — as M18's and M19's record
pull requests added none.

## Protocol-discipline observations

Three things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **Two seats measured the 2.0 post differently, and the difference is the
  front matter.** The Plan's "well under its 1850" and the pull request body's
  "895 words of body against the 2.0 post's 1850" set against the reviewer's
  and QA's "1,772". Counted with `wc -w` for this record, the 2.0 post's file
  is 1,850 words and its body after the front matter is 1,772; the new post's
  body is 895, as all three say. The pull request body compares a body count
  with a whole-file count. The requirement's "well under its length" holds on
  either figure.
- **The seat corrected its cross-repository references before the work, and
  the edit history keeps the originals.** Both Decisions were edited within
  eight seconds of posting, and only their issue references changed. A
  reference written for one repository and read in another resolves to a
  different issue or none; nothing in the protocol checks that, and the seat
  caught it itself. *Undocumented; inferred from the two edit histories.*
- **Every quotation was checked three times, and the counts of what was
  checked differ.** The pull request body lists eight quoted strings and the
  blockquote; the reviewer counts "all nine quotations (including the
  changelog block quote)"; QA "all eight quoted strings and the block quote".
  Those agree. The link counts do not quite: the reviewer's "All 24 links"
  and QA's "all 25 authored links"; the post's Markdown carries 25 link
  targets, 22 of them distinct. Neither comment says how it counted.

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **The conversation that scoped the milestone is on the record only through
  #384's goal**, "Decided with the operator 2026-09-26". Why the post was
  wanted now, and why #344 stays a capture beyond "gated on beta access", is
  not on the trail. *The reasoning is recorded; the deliberation is not*, as
  [M19's record](19-the-front-page-tells-the-crew-story.md#what-the-record-does-not-contain)
  says of its own opening.
- **The dispatch brief the task cites is not on the trail.** www#66's
  Ask-the-human section says the scope, sources and voice "are fixed by M20-R1
  and the dispatch brief"; the brief was not posted.
- **Which harness and model ran each session is not recorded by the seats.**
  This record's Deviation carries the coordination layer's account and the
  trail's partial support for it; the one commit carries a trailer, and
  neither the review nor the QA comment names its model.
- **Nothing on the trail says why the merge waited fifteen minutes and
  fifty-seven seconds after the approval**, from 10:44:51Z to 11:00:48Z, with
  the checks green throughout. Nothing names the verb that merged, either:
  `task finish` posts no comment on a task that adopts nothing, and the
  timeline shows only that the seat merged by rebase, closed the task and
  deleted the branch within two seconds.
- **The post's content is not verified by any test.** QA's verdict says the
  suite assumes "the copy, citations, deployment, responsive layout and
  theme"; the quotation audit that stands is the seat's, the reviewer's and
  QA's readings of the sources as they were fetched on 2026-09-26.
