# M19: The front page tells the crew story

Tracking issue: [#382](https://github.com/radiusred/gh-codecrew/issues/382) ·
Synthesized 2026-09-26 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #382's **seven** requirements as opened, **none added,
none struck and none amended** — the issue's `userContentEdits` holds the
original and one revision, twelve seconds later, which added the two lines
under Gates and nothing else; its **seven comments** — the qa role holder's
**three**, the first holding all seven verdicts and the other two re-verifying
M19-R7, the coordination layer's **one Decision**, the operator's **one gate
resolution**, and this record's **two Deviations**; the **nine** task issues,
all in the spoke
[radiusred/codecrew-www](https://github.com/radiusred/codecrew-www)
([#44](https://github.com/radiusred/codecrew-www/issues/44),
[#46](https://github.com/radiusred/codecrew-www/issues/46),
[#48](https://github.com/radiusred/codecrew-www/issues/48),
[#49](https://github.com/radiusred/codecrew-www/issues/49),
[#51](https://github.com/radiusred/codecrew-www/issues/51),
[#52](https://github.com/radiusred/codecrew-www/issues/52),
[#53](https://github.com/radiusred/codecrew-www/issues/53),
[#61](https://github.com/radiusred/codecrew-www/issues/61),
[#63](https://github.com/radiusred/codecrew-www/issues/63)); their **nine
merged pull requests** and the **twenty-six commits** on them; the
**thirty-three Decision records, one `**Gate resolved:**` record and two
Deviation records** across the milestone issue and its nine tasks, of which
**one, one and two** are on the milestone issue and **thirty-two, none and
none** on the tasks; the **fifteen review submissions** on the nine pull
requests — **nine approvals, every one of which stands**, and **six change
requests**, one of them withdrawn by a coordination-layer Decision and none
dismissed; the **one** adopted backlog capture the merges closed; and the
**four** captures filed in the milestone's window, all four open. The house
form is [M18's document](18-the-protocol-scales-down.md), with
[M17's](17-the-frontier-builds-the-coordinator.md) behind it; the record
standard is the one the reviewer set on [PR
#322](https://github.com/radiusred/gh-codecrew/pull/322) and re-applied on [PR
#341](https://github.com/radiusred/gh-codecrew/pull/341), [PR
#352](https://github.com/radiusred/gh-codecrew/pull/352) and [PR
#376](https://github.com/radiusred/gh-codecrew/pull/376), so every count above
and below is re-derived from the source — records by the rule
`tracker.ExtractRecords` applies (a paragraph-initial `**Decision:**`,
`**Deviation:**` or `**Gate resolved:**` label, bare or parenthetically
qualified), review states from each pull request's reviews API and timeline,
commits from each pull request's commits API and from the spoke's compare API,
and the requirement history from the milestone issue's `userContentEdits`.

**Counted at 2026-09-26T10:10:22Z, against a trail frozen from this document's
dispatch.** The last two posts on the milestone's issues before that instant
are this record's own Deviations, [the
first](https://github.com/radiusred/gh-codecrew/issues/382#issuecomment-5845307364)
at 10:01:58Z and [the
second](https://github.com/radiusred/gh-codecrew/issues/382#issuecomment-5845363545)
at 10:09:55Z; the one before them is the qa role holder's final M19-R7 verdict
at 09:56:56Z. Both Deviations are inside every total above. The instant moved
once: the first count, at 10:08:42Z, was taken before the second Deviation,
which records a slip in this document's own pull request, and differed from
this one by that comment alone.

**This record has no task.** It is delivered as a housekeeping pull request
under SPEC §4's milestone-record clause, as
[M18's](18-the-protocol-scales-down.md) was: the milestone issue is its
charter and its commits reference #382. Its records of its own are the two
Deviations on the milestone issue.

**The whole milestone happened in the spoke.** `main` in this repository did
not move between M18's last record commit
([c6153c0](https://github.com/radiusred/gh-codecrew/commit/c6153c0), committed
2026-09-24T19:30:27Z) and this document's pull request. QA's cross-cutting
check reports the same for the protocol's own files: "`git -C repo log
--oneline c6153c0..origin/main -- SPEC.md CLI.md internal cmd .codecrew/roles`
returned no lines." `main` in radiusred/codecrew-www moved by exactly
**twenty-six commits**, in a straight line with no merge commits, from M17's
remedy ([199db8c](https://github.com/radiusred/codecrew-www/commit/199db8c))
to #63's merge
([2230e98](https://github.com/radiusred/codecrew-www/commit/2230e98)). Every
one of the twenty-six has the author `radiusred-cody[bot]`. Eighteen carry the
trailer `Co-Authored-By: Claude Opus 5.5 (1M context)`; the other eight —
#44's two, #49's four and #61's two — carry no trailer.

**The trail is checked by the verb.** The installed extension reports `v2.1.0
(protocol 2.1)`; M19 cut no release. Run from this repository at the counting
instant:

```
$ gh codecrew milestone evidence 19
requirements counted: M19-R1, M19-R2, M19-R3, M19-R4, M19-R5, M19-R6, M19-R7 (7)
all 10 cited links resolve across 10 issues — evidence is reachable
```

Ten issues: the milestone and its nine sub-issues. Ten links, counted once
each: five on codecrew.works (the page and four of its anchors), all from QA's
verdicts; the deploy run QA's last M19-R7 verdict cites; QA's two findings on
#53 and #61; [radiusred/snake#6](https://github.com/radiusred/snake/pull/6),
from #48's Decision; and [M8's
record](8-a-product-home-page-for-codecrew-works.md), cited by #46's and #48's
Decisions. Neither of this record's Deviations cites a URL.

`gh codecrew milestone close 19 --dry-run`, as it read at 10:00:28Z, before
this pull request existed, passes five gates and stops at the document, as
M18's did:

```
gate milestone open: ok
gate no gate raised: ok
gate tasks closed: ok
gate requirements declared: ok
gate QA verdicts: ok
gate milestone document: refused[DOC_MISSING]: docs/milestones/19-*.md not on the default branch of radiusred/gh-codecrew — dispatch the doc-synthesizer: it delivers the record as a housekeeping PR with no task (SPEC §4), reviewed as the reviewer seat's routing requires (in pure solo, the operator confirms on the PR), and its author rebase-merges it; then rerun
dry run: nothing written — the live verb stops at the first refusal above
```

The dry run prints no raw material, so the Decisions and Deviations below were
read from each issue's comments directly. Every record below was read from its
issue, pull request, review, commit or timeline directly. The prose is wrapped
at this file's normal width; the requirement-outcomes table's rows, the fenced
output above and a handful of single links are not, because a newline ends a
Markdown table row and breaking a link breaks it.

**This PR adds the M19 ROADMAP row and changes nothing else in the hub.**
`ROADMAP.md` carried no M19 row before it. The section on the front door below
says why neither front-door document changes.

## Goal and outcome

#382 was opened by the operator's account at 00:02:04Z on 2026-09-26. Its goal
starts from the site as it stood: "codecrew.works undersells the thing that
sets CodeCrew apart. The page leads with the record, its one worked example
shows a single agent with a human reviewing, and the App identities that let
different harnesses and models build and review each other's work only turn up
in the middle of the page." It sets the direction: "This milestone rebuilds
the homepage around that story: a crew of agents, each acting on GitHub as its
own App, one building and another checking. The record and the gates come in
as the payoff of that collaboration." It bounds the scope: "It changes
codecrew.works only. The protocol, the CLI and the synced docs are untouched,
so the bedding-in rule stands." And it names its source: "It works from a
Codex review of the homepage (2026-09-26), appraised with the operator."

The page it rebuilds is the one
[M8](8-a-product-home-page-for-codecrew-works.md) created and
[M17](17-the-frontier-builds-the-coordinator.md) last changed, with the
"commodity" line M19-R2 keeps.

The seven requirements, in short (the table below carries each in full):

- **M19-R1** — the homepage's configuration is valid under protocol 2.1: typed
  App identities, `--name` on every `identity new`, the snippet marked
  illustrative, the link to CodeCrew's own routing table removed, and "a
  site-build test fails if an untyped identity or a nameless `identity new`
  returns".
- **M19-R2** — "the hero leads with the crew", keeping the "coordinator is a
  commodity" line "aimed at separation of duties", with a primary call to
  action to setup and a secondary one to the example.
- **M19-R3** — "the worked example shows two agents", "told through CodeCrew's
  own crew (Cody on Claude Code writes the code, Checky on Codex reviews it,
  each under its App), labelled as an example".
- **M19-R4** — "the proof can be checked without hovering", naming the PR, its
  author App and reviewer App, with the orchestrator receipt "visible, not
  hidden in a pop-over".
- **M19-R5** — "the middle of the page speaks to an engineer evaluating the
  product": visible role summaries, three benefits, "hub/spoke and principal
  types are left to the docs".
- **M19-R6** — "Start now" goes "install → init → solo → add an agent
  reviewer", with the required-review claim's "permission qualifier and its
  plan-tier qualifier".
- **M19-R7** — homepage-only metadata for the new pitch, and the page reading
  cleanly "at 375px and 1280px, in light and dark, with keyboard focus and the
  mobile drawer working, and docs/blog pages are unchanged".

No requirement adopts a capture. One task does: #61, carried under M19-R7,
adopts
[radiusred/codecrew-www#32](https://github.com/radiusred/codecrew-www/issues/32).

The nine tasks, in merge order — R1, R3, R4, R2, R5, R6, R7, then two
follow-ups under R7:

- **[#44](https://github.com/radiusred/codecrew-www/issues/44) / PR
  [#45](https://github.com/radiusred/codecrew-www/pull/45)** — M19-R1. Two
  commits, two review rounds, merged 00:45:23Z. Four Decisions, one of them
  the coordination layer's.
- **[#46](https://github.com/radiusred/codecrew-www/issues/46) / PR
  [#47](https://github.com/radiusred/codecrew-www/pull/47)** — M19-R3. Two
  commits, approved first round, merged 01:02:33Z. Three Decisions.
- **[#48](https://github.com/radiusred/codecrew-www/issues/48) / PR
  [#50](https://github.com/radiusred/codecrew-www/pull/50)** — M19-R4. Two
  commits, approved first round, merged 01:16:45Z. Four Decisions.
- **[#49](https://github.com/radiusred/codecrew-www/issues/49) / PR
  [#54](https://github.com/radiusred/codecrew-www/pull/54)** — M19-R2. Four
  commits, three review rounds, merged 01:39:51Z. Five Decisions.
- **[#51](https://github.com/radiusred/codecrew-www/issues/51) / PR
  [#55](https://github.com/radiusred/codecrew-www/pull/55)** — M19-R5. Two
  commits, approved first round, merged 01:54:42Z. Three Decisions.
- **[#52](https://github.com/radiusred/codecrew-www/issues/52) / PR
  [#56](https://github.com/radiusred/codecrew-www/pull/56)** — M19-R6. Three
  commits, two review rounds, merged 02:11:42Z. Three Decisions.
- **[#53](https://github.com/radiusred/codecrew-www/issues/53) / PR
  [#60](https://github.com/radiusred/codecrew-www/pull/60)** — M19-R7. Eight
  commits, three review rounds, merged 05:30:55Z. Seven Decisions.
- **[#61](https://github.com/radiusred/codecrew-www/issues/61) / PR
  [#62](https://github.com/radiusred/codecrew-www/pull/62)** — M19-R7, after
  QA's first `not satisfied`; adopts #32. Two commits, approved first round,
  merged 06:15:50Z. Two Decisions.
- **[#63](https://github.com/radiusred/codecrew-www/issues/63) / PR
  [#65](https://github.com/radiusred/codecrew-www/pull/65)** — M19-R7, after
  QA's second `not satisfied`. One commit, approved first round, merged
  09:45:50Z. One Decision.

Every one of the nine was rebase-merged by `radiusred-cody[bot]`, and every
merge was followed by a successful `Build and deploy site` run on the spoke's
`main` — nine runs, the last,
[36233759002](https://github.com/radiusred/codecrew-www/actions/runs/36233759002),
finishing at 09:46:25Z.

**The work was serial.** No two of the nine pull requests were open at the
same time; each opened after the previous one merged. Several tasks were
opened ahead of their turn — #48 and #49 at 01:02:22Z and 01:02:55Z, #51, #52
and #53 at 01:16:31Z–01:17:09Z — and each task's plan was written, and `task
start` run, only when its predecessor had merged: every plan edit on the nine
issues falls between two and eleven seconds before that task's "Started by"
comment. None of the nine task issues states why the order was serial; every
one of the nine pull requests changes the homepage's markup, template or
stylesheet, and its tests.

**The milestone ran for nine hours, fifty-four minutes and fifty-two
seconds**, from the issue opening to QA's last verdict. R1 merged forty-three
minutes after the milestone opened. R3, R4, R2, R5 and R6 followed within the
next eighty-seven minutes, the last at 02:11:42Z. R7's pull request opened at
02:28:07Z and waited two hours, twenty-four minutes and fifty seconds for its
first review; it merged at 05:30:55Z, and QA's first verdicts came thirteen
minutes and fifty-two seconds later. The two follow-ups took the rest: #61's
pull request merged thirty-one minutes after QA's first `not satisfied`, and
#63's, which waited three hours, seven minutes and forty-eight seconds for its
review, three hours and fifteen minutes after the second. This record's
Deviation carries the one account there is of both waits; the last section
says what the trail does not.

What exists afterwards that did not before. The homepage leads with the
headline "Your coding agents, working as a crew on GitHub." and names GitHub
App identities and different harnesses and models in its first paragraph; its
worked example is a six-turn review loop between Cody and Checky under a
stable `#the-example` anchor; its proof names
[radiusred/snake#6](https://github.com/radiusred/snake/pull/6), both Apps, the
finding and the fix in visible text; no pop-over remains on the page; "Start
now" has two steps, solo and then an agent reviewer; the homepage alone
carries the new title and description; and the homepage's closed phone drawer,
closed search panel and palette toggle behave for a keyboard. The spoke's
deploy workflow now builds from `uv.lock`, so the live site is built with the
theme version CI tests (Zensical 0.0.58), where it had been installing the
theme unpinned. The spoke's test suite went from 71 passing tests on PR #45 to
89 on PR #65. The hub, its protocol, its CLI and its synced docs are
unchanged.

## Decisions

Thirty-three Decision records. One is on the milestone issue, the coordination
layer's, posted by the operator's account. Thirty-two are on the tasks: one
posted by the operator's account on #44, the coordination layer's, and
thirty-one by the implementer seat (`radiusred-cody[bot]`). None carries a
parenthetical qualifier. The one gate resolution is under the gates below.

### The configuration: `myorg-*`, a long line, a rendered test

#44's [first
Decision](https://github.com/radiusred/codecrew-www/issues/44#issuecomment-5841568926)
types the example's Apps as `app:myorg-coder`, `app:myorg-checker`,
`app:myorg-tester` and `app:myorg-writer`, keeps the `~` coordinator row, and
gives the command `--name myorg-checker` "so the command and the YAML agree".
**Trade-off:** "`myorg-*` is the naming CLI.md already uses for its own
examples", and the two new names "avoid Radius Red's crew names, which the
existing site-build test keeps off the page". **Rejected:** keeping the
`*-bot` stems behind an `app:` prefix, "valid, but it reads as a role name",
and "using Radius Red's real Apps: the M8 rule keeps crew members off this
page, and M19-R1 wants the example not read against the real table".

The
[second](https://github.com/radiusred/codecrew-www/issues/44#issuecomment-5841569019)
lets the 54-character `identity new … --name myorg-checker` stand above the
page's 42-character ceiling for code lines "as a recorded exception": "no
shorter form exists", and the line is inline code that "wraps at a space on a
phone rather than scrolling". #52 later split the same command at a shell
continuation in "Start now", below.

The
[third](https://github.com/radiusred/codecrew-www/issues/44#issuecomment-5841569115)
makes the new test, `test_home_configuration_is_valid_under_protocol_2_1`,
read the rendered homepage, not the markdown, and search the whole page for
`identity new`, "so a future `identity new` in any section (R6's 'Start now',
say) is held to the same rule". It "landed in its own commit first and failed
on the current page on both counts". **Rejected:** "a check on `docs/index.md`
source: it would miss whatever the build does to the text".

The
[fourth](https://github.com/radiusred/codecrew-www/issues/44#issuecomment-5841625900),
posted by the operator's account at 00:42:50Z, is the coordination layer's: it
"withdraws the reviewer's round-1 finding on PR #45 (added file lines over 100
characters) and re-dispatches the review against the unchanged head d8f8a95".
**Trade-off:** "the 100-character cap is the commitlint rule for
commit-message lines; codecrew-www has no line-length rule for files … The
round-1 brief stated the cap ambiguously, so the finding is the dispatcher's
error." **Rejected:** "asking the implementer to reflow the lines to satisfy
the finding." The approval followed at 00:44:47Z.

### The worked example: the M8 rule narrowed, a made-up change

#46's [first
Decision](https://github.com/radiusred/codecrew-www/issues/46#issuecomment-5841717716)
names Cody and Checky and their Apps in visible text, and narrows
`test_crew_section_names_the_seats_and_no_crew_member` so it "exempts those
two names inside the `cc-example` section only; Testy and Wordy stay off the
whole page, and all four stay off every other section". **Trade-off:** M8's
rule ("Crew members are not named on the product page", bot logins "only in
link targets or screen-grabs") against M19-R3, "approved by the operator on
2026-09-26": "That later, explicit approval supersedes M8 for the worked
example only, so the test narrows rather than drops the rule." **Rejected:**
"deleting the name check", and "naming the Apps only in link targets".

The
[second](https://github.com/radiusred/codecrew-www/issues/46#issuecomment-5841728808)
makes the example's change — "a rate limit whose counter a restart resets" —
made up, and says so on the page: "The change is made up; the crew and the
verbs are real." The section's heading carries `#the-example` for R2's
secondary call to action. **Trade-off:** "M19-R4 is the requirement that names
a real PR"; and "PR #45 was considered: its round-1 change request was
withdrawn by the coordination layer rather than fixed, so it is not a build,
change request, fix, approval loop." **Rejected:** "retelling PR #45 or
another real PR here".

The
[third](https://github.com/radiusred/codecrew-www/issues/46#issuecomment-5841728913)
lays the loop out as "a single-column thread, one turn per speaker", with the
crew artwork for Cody and Checky, and removes the old tour's rules.
**Trade-off:** "No Simple Icons glyph exists for Codex in this zensical build,
and one Claude glyph for every agent is what M19-R3 retires, so the crew
artwork identifies the agents and the harness is named in text." **Rejected:**
"keeping the four numbered cards", and "harness logos as avatars (none for
Codex, and a harness is not the agent)".

### The proof: which PR, whose names, what stays visible

#48's [first
Decision](https://github.com/radiusred/codecrew-www/issues/48#issuecomment-5841783593)
identifies the PR in the proof's screenshots as
[radiusred/snake#6](https://github.com/radiusred/snake/pull/6), "established
by reading, not inference": title, head branch, the `radiusred-checky[bot]`
change request (review 5058697880), the `radiusred-cody[bot]` answer and the
approval "all match the grabs, read through the API". **Rejected:** "choosing
another PR and re-shooting the grabs", and "describing the models behind
either App on that PR (not recorded on it, so not inferred)".

The
[second](https://github.com/radiusred/codecrew-www/issues/48#issuecomment-5841783752)
exempts the proof's case block from the name check as well as the worked
example; "the four receipts, and every other section, stay free of crew
names". **Trade-off:** M19-R4 requires the author and reviewer Apps "in
visible text, and the operator's 2026-09-26 approval of M19-R3, which lifted
M8 for Cody and Checky in the worked example, is extended by the dispatch to
R4's proof where it names the Apps". **Rejected:** "exempting the whole proof
section (would let names drift into the receipts)".

The
[third](https://github.com/radiusred/codecrew-www/issues/48#issuecomment-5841800279)
makes all four receipts visible text and removes the receipts' pop-over CSS:
"M19-R4's title is 'checked without hovering', and a receipt whose link sits
in a pop-over is not checkable without one". The required-review sentence,
"now visible, gains the plan-tier qualifier SPEC §5 states". The
[fourth](https://github.com/radiusred/codecrew-www/issues/48#issuecomment-5841800399)
keeps the two screenshots "as supporting evidence; the text under them carries
the proof", because "At 375px the grabs are not legible (GitHub's text renders
at roughly 3px)".

### The hero: a headline, a bounded line, two buttons

#49's five Decisions were posted at 01:20:10Z–01:20:13Z, before the pull
request opened. The
[first](https://github.com/radiusred/codecrew-www/issues/49#issuecomment-5841891001)
sets the headline, "Your coding agents, working as a crew on GitHub.", and a
first sub paragraph. **Trade-off:** "the review's 'Give your coding agents a
crew.' was the starting point; 'on GitHub' puts the place the crew works into
the headline"; the sub "says the harnesses and models *can* differ, not that
the names attest them". **Rejected:** "the review's headline verbatim (the
brief rules out pasting it)".

The
[second](https://github.com/radiusred/codecrew-www/issues/49#issuecomment-5841891145)
keeps the commodity line in the hero, re-aimed at separation of duties, its
opening clause "the M17-R2 sentence word for word". The
[third](https://github.com/radiusred/codecrew-www/issues/49#issuecomment-5841891263)
takes the lost-transcript story and "the record is the work" out of the hero;
the Why panel keeps the line. **Rejected:** "a third, supporting hero line for
the record (the hero would lead with two stories again, the finding HP-01 was
raised against)". The
[fourth](https://github.com/radiusred/codecrew-www/issues/49#issuecomment-5841891433)
gives the hero two buttons, "Set up a crew" and "See one at work", and moves
"Read the docs" out, since "the header's Docs tab (desktop) and the drawer
behind the burger (phones) already reach the docs". **Rejected:** "the
review's labels 'Build your crew' / 'See a crew in action' (the first implies
CodeCrew builds or hosts the agents, which it does not)". The
[fifth](https://github.com/radiusred/codecrew-www/issues/49#issuecomment-5841891597)
leaves the logo's size at 375px to R7.

Two of those Decisions quote wording that the review rounds then changed. The
shipped attribution sentence and commodity clause are the ones in the
implementer's [answer to round
one](https://github.com/radiusred/codecrew-www/pull/54#issuecomment-5841967208)
and the round-two commit, below; no Decision on #49 records them.

### The middle of the page: visible roles, three bounded benefits

#51's [first
Decision](https://github.com/radiusred/codecrew-www/issues/51#issuecomment-5842095234)
gives each of the five crew badges a visible role name and one-line summary
"addressed to the reader", and removes every pop-over: "no pop-over remains on
the page". **Rejected:** "keeping the contract quotes in pop-overs as a second
layer (substance behind hover again, in the agent's voice)".

The
[second](https://github.com/radiusred/codecrew-www/issues/51#issuecomment-5842095396)
keeps the illustrative `roles:` YAML under `#the-crew`. **Trade-off:** "the
hub's docs/introduction.md links `https://codecrew.works/#the-crew` as 'an
annotated generic' routing table, and `sync_docs.py`'s `HOME_ANCHORS` maps the
README's `#the-routing-table` there; moving the YAML would leave both pointing
at a section without it, and the hub is out of scope." **Rejected:** "moving
the YAML and the command to 'Start now' (breaks the docs' anchor and pre-empts
R6's setup path)".

The
[third](https://github.com/radiusred/codecrew-www/issues/51#issuecomment-5842095530)
makes Why three benefits: "See who built it, and who checked it"; "A second
model, with other blind spots", quoting the reviewer contract's
"self-evaluation shares the blind spots of the work itself" and bounded — "it
has blind spots of its own and nobody promises it catches every bug"; and "The
record is the work". The hub/spoke description becomes a link to SPEC §3.
**Rejected:** "merging the two sections into one", and "the review's 'choose
different tools for different work' as a benefit".

### Start now: a continuation line, two qualifiers, two steps

#52's [first
Decision](https://github.com/radiusred/codecrew-www/issues/52#issuecomment-5842186879)
puts `identity new` in a second terminal, split at a shell continuation so
both lines keep to the 42-character ceiling; R1's page-wide `--name` check
"now joins a continuation before it reads a command, so the split cannot hide
a nameless one". **Rejected:** "adding `--with-approval-permission` to the
shown command (docs/identities.md: it is 'never the default'…)".

The
[second](https://github.com/radiusred/codecrew-www/issues/52#issuecomment-5842186974)
gives the required-review claim both qualifiers —
`--with-approval-permission`, "a trade of privilege for an agent-gated merge",
and "on a private repo, branch protection needs a paid GitHub plan", "the
exact wording of R4's receipt", held to one constant by the test — and ends on
what holds without either: `gh codecrew task finish` "refuses until the
reviewer App has approved". The
[third](https://github.com/radiusred/codecrew-www/issues/52#issuecomment-5842187084)
sets the two steps, "First, work solo" and "Next, add an agent reviewer", with
three follow-ups: install per account, commit the route, and start the
reviewer in a new session, where "CodeCrew does not start it; you do, or an
orchestrator does". **Rejected:** "naming Codex as the reviewer's harness here
(the worked example already pairs harnesses, and the start path should not
prescribe one)".

### R7: the metadata, four planned fixes, two from review

#53's plan named the metadata and four homepage fixes found in a first pass,
and said that anything real outside the homepage "becomes an unlabeled
'Backlog capture:' issue, not a fix". Its [first
Decision](https://github.com/radiusred/codecrew-www/issues/53#issuecomment-5842323529)
sets the homepage's description and its title — "CodeCrew — Your coding
agents, working as a crew on GitHub", "the hero's headline without its full
stop" — from the page's front matter, through `htmltitle` and `extrahead`
block overrides in `home.html`. **Trade-off:** the old title "was the other
half of the old pitch … and a link preview shows the title first".
**Rejected:** "editing main.html's tagline or `site_description` (site-wide,
and the requirement says homepage only)".

The next four are the planned fixes, each with before-and-after measurements
in headless Chromium: code blocks stepping down to 11.2px below 420px so
"nothing on the homepage scrolls sideways at 375px"
([Decision](https://github.com/radiusred/codecrew-www/issues/53#issuecomment-5842332904));
the routing example's 18 empty line anchors hidden, taking the Tab walk "58 →
40 stops at 1280px and 62 → 42 at 375px"
([Decision](https://github.com/radiusred/codecrew-www/issues/53#issuecomment-5842336526));
the hero logo at 7rem below 45em, bringing "the first hero paragraph … inside
the first screen"
([Decision](https://github.com/radiusred/codecrew-www/issues/53#issuecomment-5842340045));
and each App login in the example kept whole
([Decision](https://github.com/radiusred/codecrew-www/issues/53#issuecomment-5842349724)).

The last two answer the reviewer. The
[first](https://github.com/radiusred/codecrew-www/issues/53#issuecomment-5843378028)
hides the homepage's shut phone drawer from the keyboard with CSS keyed off
the theme's toggle; it "answers checky's change request on PR #60 (review
5324647805), which the coordination layer upheld: #58 is not out of scope on
the homepage". **Rejected:** "a script toggling `inert` or `tabindex`", and
"the rule without a homepage scope". The
[second](https://github.com/radiusred/codecrew-www/issues/53#issuecomment-5843469594)
makes the burger a keyboard control with a small inline script, answering the
second change request, "which the coordination layer upheld". **Trade-off:**
"A script is the only way to give a `<label>` key handling." **Rejected:**
"Replacing the label with a `<button>` in a header override", and "A focus
trap inside the open drawer".

### The follow-ups: deploy from the lock, then two focus fixes

The coordination layer's [Decision on
#382](https://github.com/radiusred/gh-codecrew/issues/382#issuecomment-5843627028),
posted by the operator's account at 05:45:30Z, forty-three seconds after QA's
first M19-R7 verdict, answers it with "a follow-up task under the same
requirement" that "(a) makes the deploy build with the locked theme … adopting
#32; and (b) takes the closed search panel's invisible controls out of the
homepage's Tab order". **Trade-off:** "the deploy installs zensical unpinned
(0.0.65) while CI tests the lock (0.0.58), exactly the split #32 predicted.
Deploying from the lock makes the live site the build CI verified … #32 left
the direction to the operator; the operator is unavailable overnight, the
change is one workflow file and reverts cleanly, and it is flagged for the
operator's review with the milestone's rendered-page gate." **Rejected:**
"bumping the lock to 0.0.65 (every M19 rendering check would need redoing
against an untested theme); making the title template version-proof only
(leaves the split that caused it, and #32 open)." #61 was opened three seconds
later.

#61's [first
Decision](https://github.com/radiusred/codecrew-www/issues/61#issuecomment-5843650453)
mirrors the test job in the deploy — `astral-sh/setup-uv@v6`, `uv sync
--frozen`, then `uv run` for the build — and pins the workflow's shape with a
test. **Trade-off:** "capture #32 said a test cannot guard the version split.
That is true of the versions themselves … The workflow's shape can be checked
as text, though, and that is what regressed." The
[second](https://github.com/radiusred/codecrew-www/issues/61#issuecomment-5843677941)
puts a small stylesheet inside the theme's search shadow root, because "the
theme's search does not use `#__search` for its state … and page CSS cannot
select into a shadow root", and guards the theme's minified class names with a
test that reads the theme bundle.

#63's [one
Decision](https://github.com/radiusred/codecrew-www/issues/63#issuecomment-5843955027)
gives the palette toggle's visible label the theme's own focus ring whenever
either palette radio has keyboard focus, in one homepage-scoped CSS rule. It
explains the miss: the theme "never checks a radio on load", so its ring can
land on the hidden label; measured before the change, "6 of 20 steps had
`outline: none`", including "two states QA's walk did not reach"; after, "20
of 20 show the ring". **Rejected:** "a script that checks the chosen radio on
load (it would fight the theme's palette state machine…)", and "fixing it
site-wide (docs and blog stay byte-identical; a capture follows for them)".

## Deviations

Two Deviation records, both on the milestone issue and both this record's own.
No task issue or pull request carries one. Four pull request bodies say so in
words — #45's "No deviations from the plan.", #47's, #54's and #56's the same.

**[The record takes harness facts from its dispatch
brief](https://github.com/radiusred/gh-codecrew/issues/382#issuecomment-5845307364).**
Posted at 10:01:58Z. The brief states that every implementer session ran as a
Claude Opus 5.5 sub-agent, every review round in Codex under `gpt-5.6-terra`
and every QA round in Codex under `gpt-5.6-sol`, and "that the Codex usage
limit stopped two review dispatches overnight, on radiusred/codecrew-www PR 60
and PR 65, with nothing posted either time". The comment sets out what the
trail supports — the routing table's three models, the trailers on eighteen of
twenty-six commits and none on eight, and the two long waits for a first
review — and that "No review or QA comment names its own harness or model, and
nothing on the trail says why either review waited." **Why:** "the brief is
not on the trail, so this comment puts its account there, and the record cites
this comment rather than the brief."

**[The record's pull request briefly carried a closing
reference](https://github.com/radiusred/gh-codecrew/issues/382#issuecomment-5845363545).**
Posted at 10:09:55Z. [PR
#383](https://github.com/radiusred/gh-codecrew/pull/383) opened at 10:09:28Z
with a first paragraph saying it "does not close" the milestone issue,
followed by its number; GitHub read that as a closing reference, and
`closingIssuesReferences` listed #382. The body was edited at 10:09:33Z, five
seconds later by the pull request's edit history, and the list then read
empty. **Why:** "a sentence written to say the PR does not close the milestone
put the word 'close' directly before its number, which the doc-synthesizer
contract forbids and GitHub parses without regard to the negation."

## The gates

**#382's Gates section keeps the scaffold's placeholder line and adds two of
its own**, in the edit twelve seconds after the issue was created:

- "R1 merges before any other task."
- "The operator checks the rendered page before the record is written."

**The first held, by the timestamps.** PR #45 (R1) merged at 00:45:23Z; the
next task, #46, was created at 00:50:14Z.

**The second was resolved by the operator before the last fix landed, on a
condition the trail then met.** The operator's [**Gate
resolved:**](https://github.com/radiusred/gh-codecrew/issues/382#issuecomment-5844813820)
comment, at 09:01:41Z, reads: "the operator checked the live homepage on
2026-09-26 and is happy with it as it stands: 'it tells the story coherently
and brings the reader on the journey.' The record may be written once M19-R7's
remaining follow-up (radiusred/codecrew-www#63) lands and QA re-verifies R7."
At that instant PR #65 was open and awaiting review, and the live site was the
build of #61's merge. #65 merged at 09:45:50Z and QA's re-verification,
`satisfied`, followed at 09:56:56Z. All three of QA's comments, the last
included, end with the same line: "Gate note: 'The operator checks the
rendered page before the record is written' remains outstanding. It is the
operator's gate, so QA does not answer it." The trail's answer is the
operator's comment, fifty-five minutes before the last of those lines.

**No gate was raised.** No task used `gh codecrew checkpoint`, no issue
carries `cc:needs-decision`, and the dry run reports `gate no gate raised:
ok`. Each of the nine task plans records its ask-the-human points as none.

What gated the work: CI on every pull request — `Lint commit messages` and
`Tests`, both green on the head each of the nine merged at — the reviewer seat
holder's approval on each of the nine, the operator's gate, and the QA
verdicts.

## The review rounds

Nine pull requests, **fifteen review submissions**, all by the reviewer role
holder (`radiusred-checky[bot]`), an App distinct from the authoring seat:
**nine approvals, one per pull request, all standing**, and **six change
requests**, five resolved by a later commit and one withdrawn by the
coordination layer's Decision on #44. None was dismissed, and no pull request
carries an inline review comment. Five pull requests were approved in their
first round; #45 and #56 took two; #54 and #60 took three. No pull request's
timeline carries a force-push event.

- **PR [#45](https://github.com/radiusred/codecrew-www/pull/45)** — two rounds
  on one head. The [change
  request](https://github.com/radiusred/codecrew-www/pull/45#pullrequestreview-5323835396)
  found "M19-R1 is substantively sound, but the PR misses the stated
  100-character line limit" and listed six lines. The coordination layer's
  Decision withdrew it; the
  [approval](https://github.com/radiusred/codecrew-www/pull/45#pullrequestreview-5323846116)
  at the unchanged head `d8f8a95` says: "The round-1 line-length finding is
  withdrawn as a misapplied rule: the 100-character cap applies to
  commit-message lines only."
- **PR [#47](https://github.com/radiusred/codecrew-www/pull/47)** — [approved
  first
  round](https://github.com/radiusred/codecrew-www/pull/47#pullrequestreview-5323908372).
  It checked that "The recorded M19-R3 exception is narrowly implemented: Cody
  and Checky occur only in the example; Testy and Wordy occur nowhere on the
  page", and that "The wording does not turn bot names into evidence of model
  identity."
- **PR [#50](https://github.com/radiusred/codecrew-www/pull/50)** — [approved
  first
  round](https://github.com/radiusred/codecrew-www/pull/50#pullrequestreview-5323959647),
  after checking snake#6 "through GitHub's API": "The new quote is
  character-exact; the finding/fix summaries are faithful; and every changed
  external link returned 200."
- **PR [#54](https://github.com/radiusred/codecrew-www/pull/54)** — three
  rounds. The [first change
  request](https://github.com/radiusred/codecrew-www/pull/54#pullrequestreview-5324005886)
  found "Two factual overclaims in the hero": that the attribution sentence
  "reads as unconditional", when "Apps are optional", and that the approval
  clause "is false for the supported operator-held reviewer tier". The
  [second](https://github.com/radiusred/codecrew-www/pull/54#pullrequestreview-5324031286)
  accepted both fixes and asked that the attribution list be set off: "The
  reader has to recover the subject and predicate across four
  indistinguishable comma breaks." The
  [approval](https://github.com/radiusred/codecrew-www/pull/54#pullrequestreview-5324043578)
  found "The spaced em dashes match the page's established convention".
- **PR [#55](https://github.com/radiusred/codecrew-www/pull/55)** — [approved
  first
  round](https://github.com/radiusred/codecrew-www/pull/55#pullrequestreview-5324100557):
  "The reviewer quotation is character-exact against the hub contract", and
  the claims "all match the hub’s SPEC, identities guide, CLI, and role
  contracts".
- **PR [#56](https://github.com/radiusred/codecrew-www/pull/56)** — two
  rounds. The [change
  request](https://github.com/radiusred/codecrew-www/pull/56#pullrequestreview-5324138508)
  found that the solo lead's "takes your recorded confirmation" left out the
  flag: "a pure-solo operator must invoke `gh codecrew task finish <task>
  --operator-confirm`". The
  [approval](https://github.com/radiusred/codecrew-www/pull/56#pullrequestreview-5324153741)
  followed the one-sentence fix.
- **PR [#60](https://github.com/radiusred/codecrew-www/pull/60)** — three
  rounds. The [first change
  request](https://github.com/radiusred/codecrew-www/pull/60#pullrequestreview-5324647805):
  "`#58` is not out of scope for this task" — on the homepage at 375px the
  closed drawer's links took focus so "their focus rings are off-screen";
  "being inherited from the theme and present on docs/blog does not make it
  outside a requirement framed around the homepage." The
  [second](https://github.com/radiusred/codecrew-www/pull/60#pullrequestreview-5324696554):
  the burger "has no Tab stop or keyboard activation", so "a keyboard-only
  user on a phone viewport loses the page's primary navigation". The
  [approval](https://github.com/radiusred/codecrew-www/pull/60#pullrequestreview-5324751977)
  reports the keyboard walk at both widths in both schemes and 84 tests
  passing.
- **PR [#62](https://github.com/radiusred/codecrew-www/pull/62)** — [approved
  first
  round](https://github.com/radiusred/codecrew-www/pull/62#pullrequestreview-5324879647):
  the workflow guard "fails against `origin/main` and passes here", and
  "Independent Playwright walks at 375px and 1280px in light and dark found no
  closed-search or other invisible focus stops".
- **PR [#65](https://github.com/radiusred/codecrew-www/pull/65)** — [approved
  first
  round](https://github.com/radiusred/codecrew-www/pull/65#pullrequestreview-5325488784):
  "Fresh Tab, Enter and Space toggles, Shift+Tab back from the hero, and
  reloads with a stored choice all gave the shown label `auto 3px rgb(0, 208,
  248)` plus a header-pixel change".

**What the six change requests were about.** They carry seven findings. Three
are claims on the page broader than what CodeCrew does: the unconditional
attribution, the approval clause that did not hold for an operator-held
reviewer, and a solo lead without `--operator-confirm`. Two are keyboard
failures on the homepage that the theme causes on every page: the shut
drawer's off-screen focus and the burger with no keyboard access. One is
punctuation. One is the line-length rule the coordination layer withdrew as
the dispatcher's error, which the approval records as "a misapplied rule".

## QA: three rounds, one requirement re-verified twice

The qa role holder (`radiusred-testy[bot]`) wrote three comments on #382.

**The
[first](https://github.com/radiusred/gh-codecrew/issues/382#issuecomment-5843622793)**,
at 05:44:47Z, thirteen minutes and fifty-two seconds after PR #60 merged,
holds one verdict per requirement and a cross-cutting paragraph. Each verdict
opens with the full suite passing on merged `main` — "`uv run pytest -q
--basetemp=$PWD/../pytest-tmp` passed all 84 tests" — then checks the live
site against the requirement:

- **M19-R1** — `satisfied`: every `identity:` value on the live routing
  example typed or `~`, and "I made two independent scratch copies, changed
  one identity to bare `myorg-checker` and removed `--name` from the other,
  and `test_home_configuration_is_valid_under_protocol_2_1` failed each copy
  on the intended assertion."
- **M19-R2** — `satisfied`: the hero's two claims visible; "Set up a crew" and
  "See one at work" resolving to their targets "in that order".
- **M19-R3** — `satisfied`: the six turns, and "'the operator or an
  orchestrator starts every session' and 'review runs in a fresh session' each
  occur exactly once."
- **M19-R4** — `satisfied`: "I followed all 14 substantive proof links to 200
  responses … and queried GitHub's API: the quoted sentence is
  character-exact".
- **M19-R5** — `satisfied`: five visible role summaries, "exactly three
  benefits", and the different-model benefit "promises neither a guarantee nor
  every bug".
- **M19-R6** — `satisfied`: the path and both qualifiers, with the commands
  checked "against CLI.md and the installed CLI's `--help` without executing a
  write verb".
- **M19-R7** — `not satisfied`: "Two live failures remain": "the unpinned
  deploy installed Zensical 0.0.65 and emits `Home - CodeCrew` for `<title>`,
  `og:title` and `twitter:title` instead of the crew pitch", and "Tab enters
  three invisible controls in the closed search shadow panel at 1280px and
  four at 375px". The [finding on
  #53](https://github.com/radiusred/codecrew-www/issues/53#issuecomment-5843615924)
  carries the reproduction.

Its cross-cutting paragraph checks that nothing on the page "says CodeCrew
dispatches, schedules or hosts agents", that there are "no speed, quality or
cost promises and no counts that age", and that "The M8 naming rule holds:
Cody and Checky occur only in the worked example and proof case block under
the operator's 2026-09-26 approval, while Testy and Wordy occur nowhere."

**The
[second](https://github.com/radiusred/gh-codecrew/issues/382#issuecomment-5843911002)**,
at 06:30:11Z, after #61, is `not satisfied` again and says it supersedes the
first. The lock repair is live — the live homepage, docs index and a blog post
"are byte-for-byte equal to their live HTML after removing only Cloudflare's
email-obfuscation injection" — and the search stops are gone, but "in a clean
OS-dark context with no stored palette, Tab focuses the 0×0, opacity-zero
`input#__palette_0` while both palette radios are unchecked", with the header
"pixel-identical before and during `:focus-visible`". The [finding on
#61](https://github.com/radiusred/codecrew-www/issues/61#issuecomment-5843906652)
carries it.

**The
[third](https://github.com/radiusred/gh-codecrew/issues/382#issuecomment-5845274888)**,
at 09:56:56Z, after #63, is `satisfied` and supersedes the second by comment
number. It cites the [deploy of
`2230e98`](https://github.com/radiusred/codecrew-www/actions/runs/36233759002),
repeats the byte-for-byte comparison and the metadata check, and reports: "All
16 fresh/stored × forward/backward cases put `outline: auto 3px rgb(0, 208,
248)` on the visible sun/moon label … fresh OS-dark Tab, the original failure,
passed at both widths … I found no remaining R7 failure."

**Every standing verdict is `satisfied`.** M19-R1 to M19-R6 stand in the first
comment; M19-R7 stands in the third. No requirement was struck.

## Requirement outcomes

The status column is the **standing** verdict word as the qa role holder wrote
it, verbatim and unqualified: M19-R1 to M19-R6 in comment
[5843622793](https://github.com/radiusred/gh-codecrew/issues/382#issuecomment-5843622793),
M19-R7 in comment
[5845274888](https://github.com/radiusred/gh-codecrew/issues/382#issuecomment-5845274888).

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M19-R1 — the homepage's configuration is valid under protocol 2.1: every App identity in the routing example is typed (`app:`), every `identity new` command carries `--name`, the snippet is marked as illustrative, and the link to CodeCrew's own routing table is removed; a site-build test fails if an untyped identity or a nameless `identity new` returns | [#44](https://github.com/radiusred/codecrew-www/issues/44) / PR [#45](https://github.com/radiusred/codecrew-www/pull/45) | `satisfied` | Task closed. Two review rounds on one head; the round-one finding withdrawn by the coordination layer's Decision on #44. Four Decisions (the `myorg-*` names, the 54-character line, the rendered-page test, the withdrawal). Merged first, as the Gates require |
| M19-R2 — the hero leads with the crew: visible text names GitHub App identities and says that different harnesses and models build and review; the "coordinator is a commodity" line is kept, aimed at separation of duties; a primary CTA to setup and a secondary one to the example | [#49](https://github.com/radiusred/codecrew-www/issues/49) / PR [#54](https://github.com/radiusred/codecrew-www/pull/54) | `satisfied` | Task closed. Three review rounds (two overclaims, then punctuation). Five Decisions, two of whose quoted wording the review rounds changed |
| M19-R3 — the worked example shows two agents: a build → change request → fix → approval loop, told through CodeCrew's own crew (Cody on Claude Code writes the code, Checky on Codex reviews it, each under its App), labelled as an example; it states once that the operator or orchestrator starts the sessions and that review runs in a fresh session; it replaces the single-agent "How it works" tour | [#46](https://github.com/radiusred/codecrew-www/issues/46) / PR [#47](https://github.com/radiusred/codecrew-www/pull/47) | `satisfied` | Task closed. Approved first round. Three Decisions (M8's rule narrowed to the example, the made-up change and `#the-example`, the thread layout) |
| M19-R4 — the proof can be checked without hovering: it names the PR, its author App and reviewer App, summarises the finding and the fix in visible text, and links straight to the review; the orchestrator receipt (numberguess/snake) is visible, not hidden in a pop-over | [#48](https://github.com/radiusred/codecrew-www/issues/48) / PR [#50](https://github.com/radiusred/codecrew-www/pull/50) | `satisfied` | Task closed. Approved first round. Four Decisions (snake#6 identified, the Apps named in the case block, all four receipts visible, the screenshots kept as support) |
| M19-R5 — the middle of the page speaks to an engineer evaluating the product: role summaries are visible and written for the reader rather than the agent; three benefits, including a bounded claim that a reviewer on a different model doesn't share the author's blind spots; hub/spoke and principal types are left to the docs | [#51](https://github.com/radiusred/codecrew-www/issues/51) / PR [#55](https://github.com/radiusred/codecrew-www/pull/55) | `satisfied` | Task closed. Approved first round. Three Decisions (visible summaries and no pop-overs, the YAML kept under `#the-crew`, three bounded benefits) |
| M19-R6 — "Start now" goes install → init → solo → add an agent reviewer; it covers installing each App per account, and the required-review claim carries its permission qualifier and its plan-tier qualifier (branch protection on private repos) | [#52](https://github.com/radiusred/codecrew-www/issues/52) / PR [#56](https://github.com/radiusred/codecrew-www/pull/56) | `satisfied` | Task closed. Two review rounds (`--operator-confirm`). Three Decisions (the continuation line, both qualifiers, the two steps) |
| M19-R7 — the homepage metadata matches the new pitch, set for the homepage only; the page reads cleanly at 375px and 1280px, in light and dark, with keyboard focus and the mobile drawer working, and docs/blog pages are unchanged | [#53](https://github.com/radiusred/codecrew-www/issues/53) / PR [#60](https://github.com/radiusred/codecrew-www/pull/60), then [#61](https://github.com/radiusred/codecrew-www/issues/61) / PR [#62](https://github.com/radiusred/codecrew-www/pull/62) and [#63](https://github.com/radiusred/codecrew-www/issues/63) / PR [#65](https://github.com/radiusred/codecrew-www/pull/65) | `satisfied` | All three tasks closed. `not satisfied` twice, then `satisfied`; the follow-up direction set by the coordination layer's Decision on #382. #60 took three review rounds, #62 and #65 one. Ten Decisions across the three tasks. #61 adopted and closed capture #32 |

**Captures adopted and closed by the merges — one.**
[radiusred/codecrew-www#32](https://github.com/radiusred/codecrew-www/issues/32),
filed by the operator's account on 2026-09-07 — "CI tests build with the
locked zensical while the deploy installs it unpinned" — from an implementer
report on PR #31. It said then that "the next theme change will pass CI and
break the site the same way", and left the choice between deploying from the
lock and bumping it to the operator, "since either changes the built page". It
carries an "Adopted by" comment naming radiusred/codecrew-www#61 at 05:45:34Z
and a "Closed by `gh codecrew task finish 61`" comment at 06:15:52Z naming PR
#62 and the merge commit.

**Captures filed during the milestone and left open — four**, all by the
implementer seat in the spoke, and all about the docs and blog pages that
M19-R7 kept unchanged:

- [radiusred/codecrew-www#57](https://github.com/radiusred/codecrew-www/issues/57),
  02:26:30Z, from #53: "every fenced code block on the site carries empty
  line-anchor links, each one an invisible keyboard Tab stop".
- [radiusred/codecrew-www#58](https://github.com/radiusred/codecrew-www/issues/58),
  02:26:31Z, from #53: "on a phone, the closed navigation drawer's links take
  keyboard focus off-screen, on every page". Its one comment, updated after PR
  #60's second round, records that the homepage now has both the shut-drawer
  and burger fixes and "this capture remains for every other page".
- [radiusred/codecrew-www#59](https://github.com/radiusred/codecrew-www/issues/59),
  02:26:31Z, from #53: "the docs pages and the 404 still describe CodeCrew
  with the receipts-led site description"; "Whether the site-wide description
  and tagline should follow the homepage is a call for the operator."
- [radiusred/codecrew-www#64](https://github.com/radiusred/codecrew-www/issues/64),
  06:36:55Z, from #63: "on the docs and blog pages, the palette toggle can
  take keyboard focus with no visible ring".

No issue other than #382 was filed in this repository in the window.

## The front door

The doc-synthesizer contract sends this hub's front-door list to its
`.local.md`: "the README's proof points, and the introduction's release, verbs
and refusal codes". M19 changed none of those — no release, no verb, no
refusal code, and no hub commit. [#363's
Decision](https://github.com/radiusred/gh-codecrew/issues/363#issuecomment-5817341651)
in M18 records that this hub's extension "does not send the doc-synthesizer to
codecrew.works", as [M8's record](8-a-product-home-page-for-codecrew-works.md)
records the operator declining to add the home page to the refresh duty. What
M19 changed is the site, so the check here is narrower: whether either
front-door document makes a claim about codecrew.works that M19 made false.
Read at the counting instant, against the live page (generator
`zensical-0.0.58`, title "CodeCrew — Your coding agents, working as a crew on
GitHub"):

- **`README.md`** calls codecrew.works "the marketing and introduction site"
  that carries the blog. True as it stands.
- **`docs/introduction.md`** makes three claims about the site. It links "an
  annotated generic" routing table at `#the-crew`: the anchor is on the live
  page and the illustrative `roles:` table is under it, which #51's Decision
  kept there for this link. It says "the receipts are on the home page" at
  `#codecrew-works`: the anchor is live and the four receipts are under it,
  visible since #48. And it lists the site as "the introduction site: what the
  problem is, how the protocol works, and what it has delivered": the page has
  a Why section, a worked example of the protocol's review loop and the proof
  section. None of the three is false, so this PR changes neither document.

## Protocol-discipline observations

Seven things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **The human judgment in the tasks came from outside the task issues.** All
  nine plans record their ask-the-human points as none, and no gate was
  raised. The judgment is on the trail as the milestone's requirement text
  (M19-R3 names Cody and Checky), the coordination layer's two Decisions, and
  the operator's gate resolution; and, cited by the tasks without being on the
  trail, as dispatch briefs — "The operator has already ruled that the
  example's model IDs need not match anything (dispatch brief)" on #44, "the
  operator's steer on the commodity line is in the dispatch" on #49, the
  operator's approval "extended by the dispatch to R4's proof" on #48, and "as
  the brief preferred" on #61.
- **The coordination layer corrected its own dispatch on the trail.** The
  reviewer's round-one finding on PR #45 applied a rule the review brief had
  stated ambiguously, and the coordination layer's Decision withdrew it as
  "the dispatcher's error" within a minute and seven seconds of the change
  request, with the approval following on the same head. By contrast, #53's
  two review-answering Decisions each say the coordination layer "upheld" a
  change request, and no comment of the coordination layer's records either
  upholding.
- **QA found what the pull requests' checks did not see.** CI tested the
  lock's theme; the deploy installed a newer one; the live page's title
  differed from the tested build. #32 had described that split on 2026-09-07
  and said "A test cannot guard this; the workflow shape is the guard". #61
  added a test of the workflow's shape. Both of QA's `not satisfied` verdicts
  say what the suite did not exercise — "the deployment environment or a
  browser Tab walk", then "a clean OS-dark Tab walk" — and both failures were
  found on the live site with a browser.
- **One requirement took three tasks, and the direction for the second was a
  coordination-layer Decision taken while the operator was unavailable.** The
  Decision flagged the choice "for the operator's review with the milestone's
  rendered-page gate". The gate resolution, three hours and sixteen minutes
  later, approves the page "as it stands" and does not mention #32's
  direction.
- **Decisions recorded before review were not updated after it.** #49's first
  two Decisions quote hero wording the review rounds on PR #54 replaced; the
  shipped wording is in the pull request's comments and commits. Pull request
  comments are not gathered as records, and no later Decision on #49
  supersedes the two.
- **A line in QA's comments outlived the gate it describes.** Each of QA's
  comments ends by calling the operator's gate outstanding, the last
  fifty-five minutes after the operator's `**Gate resolved:**` comment, and
  says QA does not answer it. The gate's state is carried by the operator's
  comment; `milestone close` does not read the Gates section's lines.
- **A homepage-only requirement was held against the theme's behaviour on the
  homepage.** The reviewer's first change request on PR #60 brought the
  theme's drawer inside M19-R7 ("being inherited from the theme and present on
  docs/blog does not make it outside a requirement framed around the
  homepage"), and QA's two later findings, the search panel and the palette
  toggle, were theme behaviour too. Every fix stayed homepage-scoped so the
  docs and blog stayed byte-identical in the locked build, and each theme-wide
  remainder became a capture.

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **The Codex review of the homepage and its appraisal are not on GitHub.**
  #382's goal says the milestone "works from a Codex review of the homepage
  (2026-09-26), appraised with the operator". The review's findings are not
  posted on any issue; one of their identifiers, HP-01, appears once, in #49's
  third Decision, and several Decisions quote "the review's" suggested
  wording. *The reasoning is recorded; the deliberation is not*, as [M18's
  record](18-the-protocol-scales-down.md#what-the-record-does-not-contain)
  says of its own opening.
- **The operator's approval of naming Cody and Checky is on the trail only as
  a citation.** #46's and #48's Decisions and plans and QA's cross-cutting
  paragraph each cite an approval by the operator on 2026-09-26; M19-R3's own
  text, written by the operator's account, names both. No separate comment
  records the approval, and its extension to R4's proof is attributed to "the
  dispatch".
- **The dispatch briefs the tasks cite are not on the trail.** The rulings
  #44, #48, #49 and #61 attribute to their briefs are quoted in the
  observations above; the briefs themselves, and the round-one review brief on
  PR #45 that the coordination layer called ambiguous, were not posted.
- **The coordination layer's upholding of PR #60's two change requests has no
  comment.** #53's Decisions state it; nothing else records it.
- **Which harness and model ran each session is not recorded by the seats.**
  This record's first Deviation carries the coordination layer's account and
  the trail's partial support for it; eight of the twenty-six commits carry no
  trailer, and no review or QA comment names its model.
- **Nothing on the trail says why PR #60 and PR #65 waited hours for review.**
  The first Deviation carries the brief's account — a Codex usage limit — and
  nothing was posted either time.
- **The operator's view of the deploy-from-lock direction is not recorded.**
  The coordination layer's Decision flagged it for review with the gate; the
  gate resolution does not mention it.
- **The screenshots and browser captures behind the pull requests' rendered
  checks were not committed or posted.** PR #60's body says they "are in the
  dispatch directory, not committed"; the measurements are quoted in the
  Decisions and review bodies.
