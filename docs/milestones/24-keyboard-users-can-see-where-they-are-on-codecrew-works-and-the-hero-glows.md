# M24: Keyboard users can see where they are on codecrew.works, and the hero glows

Tracking issue: [#443](https://github.com/radiusred/gh-codecrew/issues/443) ·
Synthesized 2026-10-28 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #443's **five** requirements as opened, **none added,
none struck and none amended** — the issue's `userContentEdits` is empty, so
the body has never been edited; its **three comments** — the qa role holder's
**one** verdict comment and this record's **two Deviations**; the **two** task
issues, both in the spoke
[radiusred/codecrew-www](https://github.com/radiusred/codecrew-www)
([#86](https://github.com/radiusred/codecrew-www/issues/86),
[#87](https://github.com/radiusred/codecrew-www/issues/87)); their **two
merged pull requests** and the **four commits** on them; the **five Decision
records and two Deviation records** across the milestone issue and its two
tasks — four Decisions on #86, one on #87, and both Deviations on the
milestone issue, with no `**Gate resolved:**` record anywhere; the **two
review submissions**, both approvals and both standing; the **four** adopted
backlog captures the two finishes closed; and **no** capture filed in the
milestone's window. The house form is [M22's
document](22-the-homepage-explains-codecrew-to-someone-meeting-it-for-the-first-time.md),
cut down for a smaller milestone. Counts follow the record standard the
reviewer set on [PR #322](https://github.com/radiusred/gh-codecrew/pull/322)
and has re-applied since: records by the rule `tracker.ExtractRecords` applies
(a paragraph-initial `**Decision:**`, `**Deviation:**` or `**Gate resolved:**`
label), review states from each pull request's reviews API and timeline,
commits from each pull request's commits API and the spoke's compare API.

**Counted at 2026-10-28T01:36:41Z, against a trail frozen from this document's
dispatch.** The last post on the milestone's issues at that instant is this
record's second
[Deviation](https://github.com/radiusred/gh-codecrew/issues/443#issuecomment-6458305735),
at 01:36:04Z, and it is inside every total above; before it come this record's
first
[Deviation](https://github.com/radiusred/gh-codecrew/issues/443#issuecomment-6458213308)
at 01:29:22Z and the qa role holder's verdict comment at 01:23:51Z.

**This record has no task.** It is a housekeeping pull request under SPEC §4's
milestone-record clause, as the M18 to M22 records were: the milestone issue
is its charter and its commit references #443.

**The whole milestone happened in the spoke.** `main` in this repository stood
at [8ad5637](https://github.com/radiusred/gh-codecrew/commit/8ad5637),
committed at 00:27:13Z, five minutes and fifty-one seconds before #443 opened,
and had not moved at the counting instant. `main` in radiusred/codecrew-www
moved by exactly **four commits**, in a straight line with no merge commits,
from [f74c898](https://github.com/radiusred/codecrew-www/commit/f74c898) (task
#83's logo, the commit #85's mock was made against) to
[fc548c4](https://github.com/radiusred/codecrew-www/commit/fc548c4). Every one
of the four has the author `radiusred-cody[bot]`. Six files changed:
`docs/stylesheets/extra.css` (125 lines added, 26 removed),
`overrides/home.html` (5 and 129), `overrides/main.html` (9 and 0), the new
`overrides/partials/keyboard.html` (148), `tests/test_site_build.py` (224 and
90) and `zensical.toml` (4 and 2).

**The trail is checked by the verb.** The installed extension reports `v2.1.1
(protocol 2.1)`; M24 cut no release. Run from this repository at the counting
instant:

```
$ gh codecrew milestone evidence 24
requirements counted: M24-R1, M24-R2, M24-R3, M24-R4, M24-R5 (5)
all 1 cited link resolve across 3 issues — evidence is reachable
```

Three issues: the milestone and its two sub-issues. The one link is the live
homepage, which QA's verdict comment cites. Neither Deviation cites a URL.

`gh codecrew milestone close 24 --dry-run`, run at the counting instant,
before this pull request existed, passes four gates and stops at the fifth, a
gate earlier than M18's to M22's did:

```
gate milestone open: ok
gate no gate raised: ok
gate tasks closed: ok
gate requirements declared: ok
gate QA verdicts: refused[VERDICT_MISSING]: no QA verdict on radiusred/gh-codecrew#443 for: M24-R5 — dispatch QA (.codecrew/roles/qa.md)
gate milestone document: not reached
dry run: nothing written — the live verb stops at the first refusal above
```

**This PR adds the M24 ROADMAP row and changes nothing else in the hub.** The
front-door section below says why neither front-door document changes.

## Goal and outcome

#443 was opened by the operator's account at 00:33:04Z on 2026-10-28. Its
goal: "Fix the three keyboard-focus faults found in M19's end-to-end check, so
a keyboard user's focus is always somewhere visible on codecrew.works, and
give the homepage hero the slow multi-colour glow the operator mocked and
approved." It bounds the scope: "Spoke work in radiusred/codecrew-www only; no
protocol or CLI change (the bedding-in rule stands)." Its Gates section is the
scaffold's placeholder line and nothing else.

The four captures it draws on:

- [radiusred/codecrew-www#57](https://github.com/radiusred/codecrew-www/issues/57)
  and [#58](https://github.com/radiusred/codecrew-www/issues/58), filed by
  `radiusred-cody[bot]` at 02:26:30Z and 02:26:31Z on 2026-09-26, each "Found
  during the M19-R7 end-to-end check (#53)": every fenced code block's empty
  line-anchor links, "each one an invisible keyboard Tab stop", and the closed
  phone drawer, whose links "take keyboard focus off-screen, on every page".
- [#64](https://github.com/radiusred/codecrew-www/issues/64), filed by the
  same App at 06:36:55Z that day, "Found while fixing the same thing on the
  homepage (#63, M19-R7)": on the docs and blog pages "the palette toggle can
  take keyboard focus with no visible ring".
- [#85](https://github.com/radiusred/codecrew-www/issues/85), filed by the
  operator's account at 19:21:37Z on 2026-10-25: "Give the homepage hero
  (`.cc-hero`) a slow, multi-colour glow along its top edge". The operator
  "mocked it locally against `main` at f74c898 (the transparent logo, #83),
  checked it under both palette settings, and approved the result", and the
  capture carries that mock's CSS.

M19 had fixed all three on the homepage alone — #57 and #64 by its tasks #53
and #63, #58 by PR #60, as the captures and #58's comment record — and left
the docs and blog pages as they were.

The five requirements, in short (the table below carries each in full):
**M24-R1**, no invisible Tab stop in a code block; **M24-R2**, the closed
phone drawer takes no focus on any page; **M24-R3**, the palette toggle shows
its focus ring on the visible control on every page; **M24-R4**, the hero glow
as #85 specifies it; **M24-R5**, the close-post duty.

The two tasks, both opened by the operator's account nine and nineteen seconds
after the milestone, in merge order:

- **[#87](https://github.com/radiusred/codecrew-www/issues/87) / PR
  [#88](https://github.com/radiusred/codecrew-www/pull/88)** — M24-R4; adopts
  #85. One commit, approved in its one round, merged as
  [972874e](https://github.com/radiusred/codecrew-www/commit/972874e) at
  00:49:37Z. One Decision.
- **[#86](https://github.com/radiusred/codecrew-www/issues/86) / PR
  [#89](https://github.com/radiusred/codecrew-www/pull/89)** — M24-R1, M24-R2
  and M24-R3; adopts #57, #58 and #64. Three commits, rebased onto 972874e
  before the pull request opened, approved in its one round, merged as
  [fc548c4](https://github.com/radiusred/codecrew-www/commit/fc548c4) at
  01:00:34Z. Four Decisions.

Both were rebase-merged by `radiusred-cody[bot]`, each head passed `Lint
commit messages` and `Tests`, and each merge was followed by a successful
`Build and deploy site` run on the spoke's `main`, the last,
[40612391691](https://github.com/radiusred/codecrew-www/actions/runs/40612391691),
finishing at 01:01:10Z.

**The work ran in parallel.** #87 was started at 00:35:55Z and #86 at
00:38:18Z; #86's first commit is dated 00:39:10Z, before PR #88 opened, and
its last 00:51:10Z, after PR #88 merged. #87's plan named the overlap in
advance: "If sibling #86 merges first, I rebase, keep both changes to
`extra.css`, and regenerate the hash." #87 merged first, and PR #89's body
says "The branch is rebased onto #87's hero glow (972874e), and both sets of
`extra.css` changes are kept."

**The milestone ran for fifty minutes and forty-seven seconds**, from the
issue opening to QA's verdict comment. The first merge came sixteen minutes
and thirty-three seconds after the milestone opened, the second ten minutes
and fifty-seven seconds after that, and QA's verdicts twenty-three minutes and
seventeen seconds after the second.

What exists afterwards that did not before, as the merged pull requests
describe it and QA's verdicts found it live: no page carries a `__codelineno`
anchor; below the drawer breakpoint every page's shut drawer and shut "On this
page" pop-up are hidden, and a script loaded on every page makes the burger
and the pop-up's button keyboard controls; the palette toggle's ring falls on
the label on show on every page; and the homepage hero carries #85's glow,
drifting on a 12-second cycle over a plain ink ground. This record's own read
of the live site at 01:26:51Z (generator `zensical-0.0.58`, stylesheet
`extra.css?h=61e934bba8`) found no `__codelineno` anchor on the homepage or
`/docs/cli/`, and the `cc-glow-drift` animation in the served stylesheet. The
spoke's test suite went from 150 passing tests on PR #88 to 153 on PR #89 and
on QA's merged-`main` rerun.

## Decisions

Five Decision records, all on the tasks, all posted by the implementer seat
(`radiusred-cody[bot]`), each with a trade-off and a rejected alternative.
None is on the milestone issue or a pull request.

**The line anchors go at the source** ([#86,
first](https://github.com/radiusred/codecrew-www/issues/86#issuecomment-6457493841),
00:38:36Z). "`anchor_linenums = false` in `zensical.toml`", and the homepage's
hiding rule goes because it "then matches nothing". The trade-off: "a deep
link to `#__codelineno-N-M` would stop resolving. None exists". Rejected: the
homepage's `visibility: hidden` rule without its scope, which "leaves 93 empty
links in `/docs/cli/` alone, and M24-R1 asks for them to be gone." #57 had
offered both options.

**The keyboard script loads on every page** ([#86,
second](https://github.com/radiusred/codecrew-www/issues/86#issuecomment-6457543100),
00:42:08Z). With the shut drawer hidden site-wide, "a burger with no Tab stop
would leave a keyboard user on a phone no way into the site's navigation", so
the homepage's script moves to a partial that every page loads, and the shut
search panel's controls leave the tab order everywhere. The trade-off: "the
docs and blog pages now run a small script of ours, where before they ran only
the theme's." Rejected: hiding the drawer and leaving the burger a label, and
keeping the search fix homepage-only, since "The probe counts 3 invisible
stops in the shut search panel on every docs and blog page (4 at 375px)".
#58's own comment had asked for both parts: "keep the shut drawer out of the
tab order, and give the burger a way in from the keyboard."

**The phone's "On this page" pop-up is fixed too** ([#86,
third](https://github.com/radiusred/codecrew-www/issues/86#issuecomment-6457544933),
00:42:16Z), the one change beyond the three captures. Below 60em the theme
leaves the shut pop-up's links in the tab order: "23 on `/docs/cli/` and 5 on
the newest post." The trade-off: "this goes beyond the three adopted captures,
which never measured the pop-up. It is the same fault as the drawer's ... and
M24-R2 says Tab never lands on a link drawn off-screen on any page, so it
rides here rather than in a new capture." Rejected: hiding the pop-up without
giving its button a Tab stop, and leaving it for a later capture.

**The browser proof lives in the pull request body** ([#86,
fourth](https://github.com/radiusred/codecrew-www/issues/86#issuecomment-6457675232),
00:51:17Z). The build tests "pin what that proof rests on" — the rules, the
theme rules and markup they key off, and the script on every page. The
trade-off: "CI does not re-run the Tab walk". Rejected: "a Playwright test in
`tests/`", because CI installs no browser, "so it would add a new CI
dependency, which the brief rules out." #86's plan had said this Decision
would be recorded.

**The hero keeps a plain ink ground** ([#87, the
one](https://github.com/radiusred/codecrew-www/issues/87#issuecomment-6457474431),
00:37:14Z). The purple-to-ink gradient goes as #85's note 2 asks, and
"`.cc-hero` keeps a plain `background: var(--cc-ink)` in its place rather than
no background at all." The trade-off: the glow is translucent, and with no
ground "the top of the hero rendered as a pale sky blue behind the headline"
in the light scheme; with ink, "At the very top edge this sits a little darker
than the mock, which had the gradient's `--cc-purple` there, so the type's
contrast there is a little higher, not lower." Rejected: deleting the
declaration outright, and keeping the old gradient. #87's plan had named this
fallback in advance and said it would be recorded.

## Deviations

Two Deviation records, both on the milestone issue and both this record's own.

**[The record takes harness facts from its dispatch
brief](https://github.com/radiusred/gh-codecrew/issues/443#issuecomment-6458213308)**,
at 01:29:22Z. The brief states that the two implementer sessions "were Claude
Opus 5.5 under Claude Code and ran in parallel"; that the reviewer and QA both
ran in Codex under `gpt-5.6-terra`; and that the coordination layer was the
operator's own Claude Code session. The comment sets out what the trail
supports — PR #89's three commits carry the Claude Opus 5.5 trailer and PR
#88's one commit carries none, both pull request bodies end with the
"Generated with Claude Code" line, and the two tasks overlapped — and that
"The hub's routing table routes the reviewer to Codex under `gpt-5.6-terra`
and QA to Codex under `gpt-5.6-sol`, so the brief's account of QA's model
differs from the routing table." **Why:** "the brief is not on the trail, so
this comment puts its account there".

**[QA ran under the reviewer's model, by a dispatch
slip](https://github.com/radiusred/gh-codecrew/issues/443#issuecomment-6458305735)**,
at 01:36:04Z, following up the first. "QA ran in Codex under `gpt-5.6-terra`,
not under the `gpt-5.6-sol` that the hub's routing table names for qa". It
carries the coordination layer's account, given after the first Deviation was
posted: "the run reused the reviewer's command line, so it did not follow the
routing table. It was a dispatch slip, not a decision." The coordination layer
says the M24-R5 verdict "will run under `gpt-5.6-sol`." **Why:** the first
Deviation "left the difference unexplained"; this one records the explanation
"and leaves the earlier comment as it was posted."

No task issue or pull request carries a Deviation. PR #88's body says "No
deviations from the plan."

## Gates and review

**No gate beyond CI was declared, and none was raised.** #443's Gates section
holds only the scaffold's placeholder; no task used `gh codecrew checkpoint`;
no issue carries `cc:needs-decision`. #86's plan records no ask-the-human
point. #87's plan records one, conditional: below WCAG AA contrast "the mock's
strength and reach are the operator's call, not mine. I will raise a
checkpoint". PR #88's body reports every figure above 4.5:1, "so the plan's
ask-the-human point was not triggered."

**Two pull requests, two review submissions, both approvals by
`radiusred-checky[bot]` at the head that merged, both standing.** No change
request, no inline review comment, no comment on either pull request, and no
force-push event on either timeline.

- **PR [#88](https://github.com/radiusred/codecrew-www/pull/88)** —
  [approved](https://github.com/radiusred/codecrew-www/pull/88#pullrequestreview-5676436083)
  at 00:49:09Z, four minutes and fifty-two seconds after it opened: "the
  appended glow body in `docs/stylesheets/extra.css` is byte-for-byte the
  adopted #85 mock after its header comment"; the plain-ink ground "the task's
  documented, sound response to light-scheme wash-out"; "150 tests pass". It
  spot-checked "The pink-link fallback contrast ... 6.13:1 on ink" and calls
  the recorded 4.8:1 glow minimum "plausible".
- **PR [#89](https://github.com/radiusred/codecrew-www/pull/89)** —
  [approved](https://github.com/radiusred/codecrew-www/pull/89#pullrequestreview-5676503139)
  at 01:00:10Z, four minutes and thirty-nine seconds after it opened: "153
  tests"; the built site "retains the hero-glow stylesheet changes from
  972874e"; the reviewer "exercised the built homepage, CLI page and newest
  post in Chromium", and found that both phone controls "still work as native
  labels with JavaScript disabled". "The PR closes #86 only; #57, #58 and #64
  remain adopted captures."

## QA: one comment, four verdicts

The qa role holder (`radiusred-testy[bot]`) wrote [one
comment](https://github.com/radiusred/gh-codecrew/issues/443#issuecomment-6458137655)
on #443, at 01:23:51Z. Each verdict opens with the merged-`main` rerun ("153
tests"), names what that suite assumes, and then checks that assumption on the
live site:

- **M24-R1** — `satisfied`: the live stylesheet "is byte-identical to the
  rebuilt `main` asset", and `a[href^="#__codelineno"]` "= 0" across six
  pages, "while fenced code blocks still rendered" on five of them, 43 on
  `/docs/cli/`.
- **M24-R2** — `satisfied`: full Tab walks at 375px, 1000px and 1280px on the
  homepage, CLI page and blog, where "Every stop was in the viewport,
  non-hidden, and non-zero-sized"; at 375px "Enter and Space opened the burger
  (and the docs/blog “On this page” control) ... and Escape returned it to the
  visible control", including "after live instant navigation".
- **M24-R3** — `satisfied`: three pages, two widths, both OS schemes, fresh
  and stored palette choices: "all 96 cases drew the `rgb(0, 208, 248) auto
  3px` outline on a visible, in-viewport label."
- **M24-R4** — `satisfied`: the live hero "computes the adopted 0.55-strength,
  400%-sized oklab gradient with the lifted cyan plateau, 12s drift, ink mask,
  and ink ground; it visibly glows"; no animation under reduced motion;
  readable with the gradient removed. Its own raster sample measured
  "worst-case 8.14:1 headline and 4.67:1 pink-link contrast (the CSS comment
  records 7.7:1 and 4.8:1 respectively; both live samples remain above AA)".
  And a trace at "375px/4× CPU recorded ... 60.0 fps", which answers PR #88's
  note that the repaint cost had not been profiled and that "The live-site QA
  after deploy ... is the place to see it."

The comment ends: "M24-R5 is not QA-judged yet: the close-post links occur
after this record."

## Requirement outcomes

The status column is the qa role holder's verdict word, verbatim, from comment
[6458137655](https://github.com/radiusred/gh-codecrew/issues/443#issuecomment-6458137655),
the only verdict comment. M24-R5 has none at the counting instant.

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M24-R1 — No fenced code block on the site carries a Tab stop a keyboard user cannot see: the empty line-anchor links are gone (radiusred/codecrew-www#57) | [#86](https://github.com/radiusred/codecrew-www/issues/86) / PR [#89](https://github.com/radiusred/codecrew-www/pull/89) | `satisfied` | Task closed; #57 closed by #86's finish. One review round. The line anchors off at the source (#86's first Decision) |
| M24-R2 — On a phone-width screen the closed navigation drawer takes no keyboard focus; Tab never lands on a link drawn off-screen, on any page (radiusred/codecrew-www#58) | [#86](https://github.com/radiusred/codecrew-www/issues/86) / PR [#89](https://github.com/radiusred/codecrew-www/pull/89) | `satisfied` | Task closed; #58 closed by #86's finish. The script on every page and the "On this page" pop-up, beyond the captures (#86's second and third Decisions); the browser proof in the PR body (fourth) |
| M24-R3 — On every page, keyboard focus on the light/dark palette toggle shows a visible focus ring on the visible control (radiusred/codecrew-www#64) | [#86](https://github.com/radiusred/codecrew-www/issues/86) / PR [#89](https://github.com/radiusred/codecrew-www/pull/89) | `satisfied` | Task closed; #64 closed by #86's finish. 96 live cases |
| M24-R4 — The homepage hero carries the slow multi-colour glow along its top edge as specified in radiusred/codecrew-www#85: the operator's mock with the lifted cyan, the stylesheet hash re-stamped, the hero's type keeping its measured contrast at the glow's brightest point, no animation under prefers-reduced-motion, and an unglowing hero that still looks right in browsers without oklab/color-mix | [#87](https://github.com/radiusred/codecrew-www/issues/87) / PR [#88](https://github.com/radiusred/codecrew-www/pull/88) | `satisfied` | Task closed; #85 closed by #87's finish. One review round. The plain ink ground (#87's one Decision). QA's live pink-link sample, 4.67:1, is below the 4.8:1 the CSS comment records and above AA |
| M24-R5 — The close comes with a Bluesky and an X post under M23-R9's duty as amended on #423: drafted in the record PR, approved by the operator, links commented here | this record's pull request (the draft); the operator (approval, posting, the links comment) | none | Deferred by QA: "M24-R5 is not QA-judged yet: the close-post links occur after this record." The close's QA gate refuses `VERDICT_MISSING` for it at the counting instant |

**Captures adopted and closed by the merges — four.** Each carries an "Adopted
by" comment from the operator's account at 00:33:16Z–00:33:25Z and a "Closed
by `gh codecrew task finish`" comment from `radiusred-cody[bot]`: #85 at
00:49:42Z, naming PR #88 and 972874e; #57, #58 and #64 at 01:00:37Z–01:00:42Z,
naming PR #89 and fc548c4.

**Captures filed during the milestone — none.** From #443's creation to the
counting instant, the only issues or pull requests created in this repository
and the spoke are #443, the two tasks, their two pull requests and this
record's own pull request,
[#444](https://github.com/radiusred/gh-codecrew/pull/444); an
organisation-wide issue search over the same window, as this record's token
sees the organisation, returns the same six.

## The front door

The hub's front door is the README and `docs/introduction.md`, checked for
"the README's proof points, and the introduction's release, verbs and refusal
codes". M24 changed none of those: no release, no verb, no refusal code and no
hub commit. As in [M22's
record](22-the-homepage-explains-codecrew-to-someone-meeting-it-for-the-first-time.md#the-front-door),
the remaining check is whether either document says something about
codecrew.works that M24 made false. The README calls it "the marketing and
introduction site" that "carries the blog"; the introduction links the
homepage's `#the-crew` and `#codecrew-works` anchors. Both anchors are on the
live page at 01:26:51Z, and the blog answers 200. Nothing is false, so this PR
changes neither document, and, as SPEC §4 gives the record's diff, adds no
CHANGELOG entry.

## Protocol-discipline observations

- **The close-post requirement and the close gate sit in an order the duty's
  text does not.** The hub's `doc-synthesizer.local.md` says that "Soon after
  the close the operator approves the exact text for each network, the post
  goes out ... and one comment with its links goes on the milestone issue."
  M24-R5 makes that comment a requirement, QA deferred its verdict until the
  links exist, and `milestone close` refuses `VERDICT_MISSING` while any
  requirement lacks one. As the trail stands, M24 can close only after the
  post is out and QA has verdicted R5.
- **One task went beyond its captures, on the record.** The "On this page"
  pop-up was in no capture and in no requirement's text by name; #86's third
  Decision tied it to M24-R2's "on any page" before the pull request opened,
  the pull request body flagged it for the reviewer, and QA's R2 verdict
  tested it.
- **QA checked each suite's assumption rather than the suite.** Every verdict
  names what the merged-`main` tests take on trust — the deployed asset being
  current, the browser's focus behaviour, the rendered control, the compositor
  — and probes that on the live site, which is the gap #86's fourth Decision
  left open by keeping the Tab walk out of CI.
- **A seat ran off its routing, and nothing a seat wrote showed it.** QA's
  verdict comment names no model, and CodeCrew does not dispatch agents, so
  the model a seat runs under is the dispatcher's to honour (SPEC §5); the
  difference surfaced because the dispatch brief stated the model QA ran
  under, and its cause is on the trail only in this record's second Deviation.
- **One of the four commits has no co-author trailer.** 972874e, PR #88's one
  commit, carries no `Co-Authored-By` line; PR #89's three do. Both pull
  request bodies end with the "Generated with Claude Code" line.

## What the record does not contain

- **The close post itself.** M24-R5 is unverdicted: the draft is in this pull
  request's description, and the operator's approval, the posts and the links
  comment come after the record.
- **The operator's mock session.** #85 says the operator mocked the glow
  locally, compared three cyans and "approved the result"; the mock's CSS is
  on the trail, the session is not.
- **The coordination layer's dispatch briefs.** #86's fourth Decision cites
  "the brief", which is not on GitHub.
- **Which harness and model ran each session**, as the seats record it. This
  record's first Deviation carries the coordination layer's account and the
  trail's partial support for it; the second gives the coordination layer's
  explanation for QA's model differing from the routing table. Neither is
  something a seat recorded.
- **The screenshots and probe scripts behind the rendered checks.** PR #89's
  body says "Screenshots stay in the dispatch directory", and QA's R4 verdict
  says "screenshots retained locally"; the figures are quoted in the bodies
  and the verdicts.
