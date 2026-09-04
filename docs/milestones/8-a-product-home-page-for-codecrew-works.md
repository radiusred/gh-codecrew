# M8: A product home page for codecrew.works

Tracking issue: [#196](https://github.com/radiusred/gh-codecrew/issues/196) ·
Synthesized 2026-09-04 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #196's four requirements, its Gates section and its
single QA round; the three task issues in the spoke
[radiusred/codecrew-www](https://github.com/radiusred/codecrew-www) and their
merged PRs ([#2](https://github.com/radiusred/codecrew-www/issues/2) / PR
[#3](https://github.com/radiusred/codecrew-www/pull/3),
[#4](https://github.com/radiusred/codecrew-www/issues/4) / PR
[#5](https://github.com/radiusred/codecrew-www/pull/5),
[#9](https://github.com/radiusred/codecrew-www/issues/9) / PR
[#10](https://github.com/radiusred/codecrew-www/pull/10)); the **25 Decision
and 13 Deviation comments** across those three issues, and the twenty-six
coordinator notes they answer; the reviewer's four review submissions on the
three PRs; the two backlog captures the milestone filed in the hub
([#197](https://github.com/radiusred/gh-codecrew/issues/197),
[#199](https://github.com/radiusred/gh-codecrew/issues/199)); and the deployed
page at <https://codecrew.works/>. The M6 and M7 documents supplied the house
form.

The gather from `gh codecrew milestone close 8` is again not part of the raw
material, for the reason the last three documents have each recorded: the live
close refuses `DOC_MISSING` until this file is on `main`, and `--dry-run`
stops earlier still, at `OPEN_TASKS`, naming
[#203](https://github.com/radiusred/gh-codecrew/issues/203) — the task that
writes it. Four milestones now. Every record below was read from its issue
directly.

**This is M8's first PR in the hub.** Every task in the milestone lived in the
spoke, so the hub had no branch for the milestone's work; the ROADMAP row went
in as an operator commit straight to `main`
([ae09d01](https://github.com/radiusred/gh-codecrew/commit/ae09d01), under the
org ruleset's admin bypass), which is what capture
[#197](https://github.com/radiusred/gh-codecrew/issues/197) is about. The
observations at the end return to it.

## Goal and outcome

`codecrew.works` served a placeholder home page built, like every other page
on the site, from the standard Zensical page layout. M8's goal
([#196](https://github.com/radiusred/gh-codecrew/issues/196)) was to turn that
page into a product page — "in the flow shared by opengsd.net, openrouter.ai
and paperclip.ing but simpler: hero, how it works, why, proof, install" —
while every other page on the site kept the reading column, the sidebars and
the footer navigation it already had.

Four requirements, M8-R1 through M8-R4, all present when the milestone opened
at 22:11:01Z on 2026-09-02; none added later. Three tasks delivered the work,
all three in `radiusred/codecrew-www`, a spoke whose `.codecrew.yml` names
this repository as its hub; this document is the fourth task and the only one
in the hub. The site deploys on every push to the spoke's `main`
(`.github/workflows/site.yml`), so each merge below was live within minutes.

- **[#2](https://github.com/radiusred/codecrew-www/issues/2) / PR
  [#3](https://github.com/radiusred/codecrew-www/pull/3)** — the first cut:
  the per-page template, the front-matter hide rules and the hero-to-install
  flow. Opened and merged inside eighteen minutes on 2026-09-02
  ([76907cd](https://github.com/radiusred/codecrew-www/commit/76907cd)), three
  commits, 8 tests.
- **[#4](https://github.com/radiusred/codecrew-www/issues/4) / PR
  [#5](https://github.com/radiusred/codecrew-www/pull/5)** — the styling pass,
  and the bulk of the record. Opened as five changes; finished as
  twenty-two, seventeen of them added by coordinator note while the PR was
  open. Merged 2026-09-03
  ([1ef031e](https://github.com/radiusred/codecrew-www/commit/1ef031e)), thirty
  commits, 24 tests.
- **[#9](https://github.com/radiusred/codecrew-www/issues/9) / PR
  [#10](https://github.com/radiusred/codecrew-www/pull/10)** — a defect the
  operator found on a phone after the page was live: the burger opened an
  empty drawer. Merged the same evening
  ([bd9c935](https://github.com/radiusred/codecrew-www/commit/bd9c935)), one
  commit, 49 tests.

Delivered in full. All four requirements carry a `satisfied` verdict from the
qa role holder in a
[single round](https://github.com/radiusred/gh-codecrew/issues/196#issuecomment-5534140553)
on 2026-09-04, against merged `main` (53 tests by then) and against the
deployed page rather than the source.

Two things the milestone shipped that no requirement asked for: the operator's
own artwork and screen-grabs, which arrived mid-PR and were recorded as they
landed; and a content review by this seat, dispatched mid-milestone, which
became the page's copy. Both are below.

## Decisions

Twenty-five Decision comments, grouped. Each group links every comment it
covers.

### The mechanism: a per-page template, and copy that stays Markdown

Three decisions on [#2](https://github.com/radiusred/codecrew-www/issues/2)
set the shape the rest of the milestone worked inside.

`docs/overrides/home.html`
[extends `main.html`, not `base.html`](https://github.com/radiusred/codecrew-www/issues/2#issuecomment-5517256711)
as the brief's wording said. `main.html` is the repository's own override,
carrying the Open Graph, Twitter Card and JSON-LD head; extending `base.html`
would have silently dropped the social-preview metadata from the one page most
likely to be shared, or forced a copy of that block into the new template.
The block chain is otherwise identical.

[The copy lives in `docs/index.md` as Markdown](https://github.com/radiusred/codecrew-www/issues/2#issuecomment-5517256854),
with `md_in_html` `<section markdown>` wrappers carrying the `cc-*` classes;
the template supplies only the frame. **Rejected:** hard-coding the sections
in the template and passing copy through front matter — MiniJinja has no
Markdown filter, so every paragraph would have become raw HTML in front
matter and the page would have stopped being a document. That choice is why
seventeen later changes could be made in the Markdown rather than in the
template, and it held: `home.html` changed once more in the whole milestone
(one class on the footer), until #9 reopened it.

[The hero carries its own colour scope](https://github.com/radiusred/codecrew-www/issues/2#issuecomment-5517256983)
— `data-md-color-scheme="slate"` on its own `<section>` — so it is the logo's
dark ground in *both* schemes and does not follow the palette toggle. That is
deliberate: the logo artwork sits on a dark ground and would read as a dark
rectangle on the light scheme's off-white. The same comment carries the
full-bleed decision: `.md-main__inner:has(> .cc-home) { max-width: none }`,
because the `width: 100vw; margin-left: calc(50% - 50vw)` trick scrolls
sideways whenever a vertical scrollbar is present, and JavaScript was ruled
out by the brief. Where `:has()` is unsupported the page degrades to the
standard reading column, which still reads.

### Colour, contrast, and a stylesheet finding

[Change 8](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5517955836)
is the milestone's densest colour decision, and it contains a finding worth
keeping: Zensical's *modern* stylesheet paints `.md-footer` from
`--md-default-bg-color` and the meta strip from
`--md-default-fg-color--lightest`, so the `--md-footer-bg-color` values sitting
in `extra.css` **had never applied** — which is why the footer had always read
as part of the band above it. The fix paints the home footer through the
variables the modern stylesheet actually reads.

The rest of the group is contrast arithmetic, all of it measured:
[the hero buttons](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5517657753)
take one new variable, `--cc-cyan-tint` (the brand cyan at 20% alpha), as a
translucent fill over the gradient rather than a solid block — a solid
`--cc-cyan-deep` with white text measures about 4.2:1, under the 4.5:1 that
0.8rem semibold text wants — and lose the underline Zensical applies to every
`.md-typeset a`. A second accent, `--cc-pink`, sampled from the docs icon,
becomes the home page's link colour, with `--cc-pink-deep` where the light
scheme drops it to 3.1:1; hover is a `color-mix()` towards the scheme's text
colour, so no third variable, and it is applied through `--md-typeset-a-color`
on `.cc-home` alone so the blog's links are untouched.

### The 42-character rule, measured once and then enforced everywhere

[The install block's line ceiling is 42 characters, not the 60 the coordinator
note suggested](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5517657870).
The number is measured, not chosen: in headless Chromium at 420px the code
column holds 43 characters before it scrolls, and the longest command,
`gh extension install radiusred/gh-codecrew`, is 42. A 60-character ceiling
would have passed lines that scroll on the very phone width the change was
for. The test pins 42 and a comment in it says where the number came from.

The rule then propagated. Step snippets were first held at
[44](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5517686406),
the measured inline-code capacity, which forced the four How-it-works steps
from a row into a two-by-two grid: four across a 56rem column gives each step
256px and a 44-character snippet needs 356px, so the first build ran the
fourth step off the right edge at 1440px. Then the
[superseding decision](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5518733462)
cut the snippets to bare verbs (`gh codecrew task start`, longest 27
characters) on the operator's answer that "it's marketing not a technical
manual", and brought the step ceiling down to the install block's 42 — one
rule for every code line on the page. That comment is explicitly marked
`**Decision (superseding):**` and names the decision it replaces, which is the
protocol working as intended.

### The images that came, and went

The most expensive sequence in the milestone, and the one a reader three
months from now is most likely to be puzzled by, because the page carries
almost none of it.

Change 10 wired the operator's four screen-grabs in as dimmed backgrounds
behind the How-it-works steps. Its
[Decision](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5518874387)
anchors each grab by a share of its own height (`--cc-step-bg-y`) rather than
by `background-position` percentage, because a percentage depends on how far
the image overflows the card and so lands on a different feature at 1440px
than at 420px; it also records that Chromium resolves a `url()` held in a
custom property against the *stylesheet* that substitutes it, which is why the
first build showed plain ink. A
[third refinement](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5518936581)
eased the filter and re-anchored one grab on a retake.

[Change 12](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5519198173)
then deleted all of it — images, rules, filter, fade and the tests that pinned
them — and put a backdrop hook on the proof section instead, shipped
deliberately unused with a stylesheet comment naming the file it expected and a
test covering both the unwired and the wired case. That hook was itself
superseded by
[change 14](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5524980794),
which brought the capture into the foreground as an angled card. Change 14 is
also where the seat records a specificity trap it only found by rendering:
`.md-typeset .cc-capture { transform: none }` matched the textual test but lost
in the browser to the `:nth-child` rotations, so the phone rule is written as
`.md-typeset .cc-capture:nth-child(n)` and the test pins the form that wins.

[Change 15](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5525048919)
reshaped the single card into two — the same page captured upper and lower,
faded into each other with CSS `mask-image` on the whole card so the frames
fade with the pictures, which is what makes two framed pictures read as one
document; a
[refinement](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5525101287)
enlarged them and overlapped the second above the first by a measured 73px.

Net effect on the shipped page: the step grabs are gone, the proof backdrop is
gone, and the proof pair is what survived. Each Decision states plainly that
the dead rules were removed rather than left "in case", which is why the
stylesheet carries none of this history.

### The terminal that copies clean

[Change 11](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5519178317)
turns the Start now block into a terminal window: each command a
`<span class="cc-term__line">` inside one `<pre><code>`, with the `$` prompt as
a `::before` and the dim output line as an `::after` from a `data-out`
attribute. Generated content is in neither the DOM text nor a selection, so
Zensical's copy button — which reads `getSelection().toString()` — yields
exactly the five commands, verified in headless Chromium in both schemes at
both widths. **Rejected:** shell comments as the output lines (the copy would
carry them); a fenced block (pygments spans give no per-line hook); and
`user-select: none`, because a selection string is not guaranteed to omit it
across browsers and generated content is guaranteed to. The same comment
records the one place the 42-character rule bends: the CSS prompt is fitted
into the measured slack rather than counted against the rule, since counting
it would have put the verbatim install command over the limit.

### Popovers, and three defects found by rendering

The four Decisions from
[change 16](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5525240314)
through
[change 19](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5528220243)
build a CSS-only popover — a Markdown panel inside a focusable trigger, hidden
at rest, shown on `:hover` and `:focus-within`, a fixed bottom sheet below
45em so it cannot leave the screen — and then fix three defects that only a
render could have found:

- **The phone lift pinned its own sheet.** A transformed element becomes the
  containing block for `position: fixed` descendants, so `translateY(-2px)` on
  a badge trigger pinned the badge's bottom sheet inside the badge's 55px box.
  On phones the lift is `top: -2px` on the relatively positioned trigger
  instead
  ([change 16](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5525240314)).
- **A bright flash on hover, root-caused rather than tuned away.** At rest the
  trigger had `outline-style: none` and `outline-color` computing to
  `currentColor`, so the 125ms colour transition began from the *text* colour
  and the first frames painted a near-white ring in slate and a near-black one
  in the light scheme. The fix rests the trigger at
  `outline: 0.08rem solid transparent` so only the colour transitions, from
  nothing, and deletes the separate `:focus-visible` rule that was brighter
  than the sustained state. Sampled frame by frame afterwards: one hue, only
  the alpha moving
  ([change 17](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5528061612)).
- **Row-1 panels painted under row 2.** Confirmed with `elementFromPoint`
  before the fix. The lifted trigger's transform makes it a stacking context
  at `z-index: auto`, so the panel's `z-index: 4` ranked only *inside* that
  context and the later row-2 receipts painted over the whole of it. The open
  trigger now takes `z-index: 5`; a z-index on the panel could never have
  escaped, and an unconditional equal z-index on every trigger would still
  have lost to DOM order
  ([change 18](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5528113188)).

The same group carries the crew-badge treatment
([change 9 visuals](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5518771641),
[change 20](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5528897130)):
the badges are white marks on transparency and vanish on the light scheme's
off-white, so each sits on a rounded purple tile — the ground the App avatars
carry on GitHub. When the operator redrew all five, change 20 re-measured the
new artwork (50–66% near-white pixels), found the tile still necessary, and
said so rather than dropping it silently, as the coordinator note had asked.

### The two decisions on the drawer

[#9](https://github.com/radiusred/codecrew-www/issues/9) is small and both its
Decisions are load-bearing.

[Emit the sidebar `hidden` rather than writing a media query of our own.](https://github.com/radiusred/codecrew-www/issues/9#issuecomment-5531858105)
Zensical emits the primary sidebar with `hidden` for a page whose front matter
says `hide: navigation` — which `docs/index.md` already said — and its
stylesheet then gives `.md-sidebar--primary` `display: block` inside
`@media screen and (max-width: 76.234375em)`, the same breakpoint at which the
nav tabs collapse. So upstream, `hidden` means *absent on desktop, drawer on
mobile*: exactly what M8-R1 asks for on either side of that breakpoint. The
override had removed the element instead of hiding it. **Rejected:** our own
`@media` rule, which would have duplicated 76.234375em in a repository with
nothing keeping it in step with the theme's, whose failure mode is a window
width showing neither tabs nor drawer. The cost of the chosen route is a
dependency on a rule this repository does not own, so a test asserts that rule
still exists at that exact query — a theme upgrade that drops it fails in CI
rather than silently on a phone.

[The drawer keeps Material's nested "On this page" list.](https://github.com/radiusred/codecrew-www/issues/9#issuecomment-5531857882)
Rendering the primary sidebar brings `partials/nav.html` with it, and that
nests a secondary list under the *active* item, so the home page's drawer now
lists its own section headings under "Home". `partials/nav.html` does not
consult `hide:`, so this happens for any page with `hide: toc`. Left in place:
it is the theme's standard mobile behaviour, every other page already does it,
and on a phone it is the only way to reach "Start now" without scrolling the
whole product page. The comment names the objection itself — "a reader who
reads M8-R1 strictly could call this a partial regression" — and QA
independently agreed it is not; see below.

## Deviations

Thirteen Deviation comments, every one of them on
[#4](https://github.com/radiusred/codecrew-www/issues/4), every one recording
the same shape of event: **a change added to an open PR by coordinator note,
after the plan was posted.** Each links the note it came from.

| # | The change | Deviation |
|---|-----------|-----------|
| 7 | The "you do not run the verbs" lead moves to How it works; each step gains a verb snippet | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5517641219) |
| 8 | Darker bands, a footer of its own, a second accent for links | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5517884546) |
| 9 | The page copy is replaced by the content review, with the operator's seven answers applied | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5518733337) |
| 11 | Start now becomes a terminal window; the Why panels lose their figure for glyphs | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5519129472) |
| 12 | The step grabs are removed again; a backdrop hook goes on the proof section | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5519185925) |
| 13 | Each step's two lines become a two-bubble conversation (carries a second **Deviation (copy)** for the lead's clause) | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5519346121) |
| 14 | The proof capture comes to the foreground as an angled card, superseding 12's backdrop | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5524954824) |
| 15 | The capture is retaken as two halves that fade into each other | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5525017269) |
| 16 | Caption, one-sentence receipts, CSS-only popovers on the badges and the fourth receipt | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5525199398) |
| 17 | "It has shipped" becomes "CodeCrew Works"; four uniform receipt cards; the flash fixed | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5527997838) |
| 18 | Row-1 popovers painting under row 2, fixed at the source; panels raised off the page ground | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5528097133) |
| 19 | The ring goes from the receipt cards; every popover larger, clearer, rising out of its card | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5528203492) |
| 21 | The payoff line moves beside the terminal, larger, in quotes | [link](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5529040319) |

Two are worth reading beyond the table.
[Change 14](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5524954824)
records that the superseded backdrop had already been wired in a local commit
minutes earlier and was folded into the replacement before pushing, so the
branch shows the image arriving as a card rather than as a wire-then-unwire
pair — a deliberate choice about what the history should say, disclosed rather
than hidden. And
[change 13](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5519346121)
carries a second `**Deviation (copy):**` heading for a one-clause edit to copy
that had been approved verbatim: the smallest deviation in the milestone,
recorded because the coordinator note asked for it to be.

**Three changes produced no Deviation comment.** Changes 10 (the step grabs)
and 20 (doubling the badges) arrived by coordinator note exactly as the
thirteen above did and were recorded only as Decisions; change 22 (committing
the operator's redrawn favicon and logo) produced neither, and its only trace
is the commit
[`ef2188ba`](https://github.com/radiusred/codecrew-www/pull/5/commits/ef2188ba)
and the PR description. Nothing was concealed — the coordinator notes are all
on the issue and the PR body enumerates all twenty-two changes — but a reader
gathering `**Deviation:**` comments alone would count thirteen of sixteen
plan departures.

## The content review, and a seat used out of turn

Between task #4's plan and its ninth change, this seat — the doc-synthesizer —
was dispatched to read the page as copy rather than to record it. The result
is a
[22,000-character content review](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5518094157)
on #4, posted with an unusual first line: "proposal, not applied. Nothing in
this comment has been committed anywhere; it is copy for the operator to
accept, amend or reject." Its sharpest finding was a page arguing with itself:
a How-it-works lead saying "you do not run the verbs" above four steps that
formatted those verbs exactly as the install block formats commands the reader
*is* told to type. It proposed replacement copy and closed with seven open
questions.

The operator
[answered all seven in eight lines](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5518717573)
and the
[change-9 note](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5518726720)
turned the answers into instructions. Three of the answers shaped the page
that shipped:

- **Crew members are not named on the product page.** "The names are not part
  of the framework, they are the radiusred crew. Role names should be used" —
  applied to the copy and to the agent-staffed receipt, with the bot logins
  allowed only inside a link target or a screen-grab. The rule was scoped
  expressly to the product page, not to this repository's own README, which
  still names them.
- **No standing refresh obligation for the home page.** Answer 7 declined to
  add codecrew.works to the doc-synthesizer's per-milestone refresh duty
  (which covers the README and `docs/introduction.md`). The consequence is a
  copy rule rather than a process one: the page carries no counts that age —
  "every milestone of the framework itself", not "seven milestones" — so a
  milestone boundary cannot leave it stale. This document is the first test of
  that: nothing on the home page needed changing for M8's close.
- **"Read the docs" stays off-site for now.** The hero's second button kept
  its GitHub README target, with the docs "coming onto this site as their own
  task; the button flips to the site path then, not now." It was the second of
  the two ask-the-human points raised on
  [#2](https://github.com/radiusred/codecrew-www/issues/2#issuecomment-5517257115),
  and M9 closed it — the QA verdict below reaches an on-site Docs CTA.

The first ask-the-human point from #2 — that the hero's "a third option" had
lost the antecedent that earns it — was answered inside the milestone: the
operator's [change-2 amendment](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5517583561)
first replaced the phrase with "an engineering process framework", and the
content review's answer 1 then restored an antecedent in one sentence.

## The gates

**#196's Gates section was left as the scaffold's placeholder** — the template
line "_What "done" means beyond CI: e2e suites, manual UAT, sign-offs._",
unedited. So the milestone declared no gate beyond CI, and none was raised or
resolved through the protocol.

What actually gated the work was the operator at a dev server, repeatedly and
informally. Several of the twenty-two changes on #4 arrived from a local
render rather than from a plan — change 8 "from the operator's local review of
PR #5", and the
[change-10](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5518920915),
[change-11](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5519152414)
and [change-15](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5525076286)
refinements, each captioned "from the dev server". The
[change-22 note](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5529090074)
is the closest thing to a sign-off the record holds: "This is genuinely the
last change: the operator has approved the page and is ready for review."

Alongside that, two gates the tasks set themselves and met every time: a
`zensical build --clean --strict` that must report "No issues found", and a
headless-Chromium render at 1440px and 420px in both colour schemes, with the
findings quoted as measurements. The PR #5 description lists forty of those
measurements. M8-R4 is satisfied on the strength of them.

## The review rounds

Three PRs, four review submissions by the reviewer role holder
(`radiusred-checky[bot]`), one change request.

- **PR [#3](https://github.com/radiusred/codecrew-www/pull/3)** —
  [approved first round](https://github.com/radiusred/codecrew-www/pull/3#pullrequestreview-5095817689).
  The reviewer ran the build, the suite and a strict build, and compared the
  built blog page against a strict build from `origin/main` — byte-for-byte
  identical, which is how M8-R3 was actually proved rather than asserted.
- **PR [#5](https://github.com/radiusred/codecrew-www/pull/5)** —
  [approved first round](https://github.com/radiusred/codecrew-www/pull/5#pullrequestreview-5104626456),
  after thirty commits and twenty-two changes. The reviewer re-ran the
  headless-Chromium checks independently, tabbed through all nine popover
  triggers for focus behaviour, and repeated the byte-identical blog
  comparison. It also filed the one reservation in the milestone's review
  record: "the tests are unusually implementation-specific in places", accepted
  because their structural assertions accurately describe the generated output
  and are backed by strict-build and browser verification.
- **PR [#10](https://github.com/radiusred/codecrew-www/pull/10)** — two
  rounds. The
  [first requested changes](https://github.com/radiusred/codecrew-www/pull/10#pullrequestreview-5106642388)
  **against the record, not the code**: the PR claimed all three new tests fail
  against the previous template, and the reviewer reverted only the `site_nav`
  block in a scratch copy and got 2 failed, 1 passed — the third test is a
  forward guard against a theme upgrade, not a regression test. "No code change
  is required for this finding." The implementer
  [agreed in the terms the reviewer used](https://github.com/radiusred/codecrew-www/pull/10#issuecomment-5531903682)
  — "I read a count and asserted a set" — corrected the PR body, amended the
  commit message (because a rebase merge would have put the false sentence in
  `main`'s history permanently) and posted a
  [correction on the task issue](https://github.com/radiusred/codecrew-www/issues/9#issuecomment-5531903899).
  The
  [second round](https://github.com/radiusred/codecrew-www/pull/10#pullrequestreview-5106680852)
  verified that the force-push changed only commit metadata by resolving both
  commits to the same tree object and diffing them to empty, then approved.

Every round that requested changes named a defect that was then fixed, which
is the eighth milestone in a row that can say so — and this one is the first
where the defect was in the verification narrative rather than in the
software. The reviewer read a claim about tests, ran the experiment the claim
described, and got a different number.

## QA: one round, four satisfied

The qa role holder (`radiusred-testy[bot]`) verdicted all four requirements in
a [single comment](https://github.com/radiusred/gh-codecrew/issues/196#issuecomment-5534140553)
at 00:58:01Z on 2026-09-04, against merged `main` (53 tests) and — for every
verdict — against the deployed page rather than the source. The probes worth
reading: the live home page curled and read for the product narrative in order
rather than asserted from the Markdown; every local image, stylesheet, script
and internal destination on it resolved to 200, including both proof captures;
`/blog/` and `/docs/` fetched live and confirmed to keep the reading column,
both sidebars and prev/next navigation; the palette controls checked to map a
light OS preference to `default` and a dark one to `slate`, with representative
foreground/background pairs measured in both (lowest about 6:1).

The M8-R1 verdict is the one to read beside the #9 Decision above. QA followed
the header's collapsed-navigation contract as a "hairbrush probe" — below the
tabs breakpoint the hidden element becomes a populated drawer reaching Home,
Docs and Blog — and concluded in its own words that the drawer's "nested
heading links are content inside the mobile replacement for the tabs, not the
absent persistent left navigation or right-hand panel, so the #9 repair does
not regress R1." The seat that made the change had named the same objection
and reached the same answer independently; the verdict is what settles it.

## Requirement outcomes

| Requirement | Delivered by | Outcome |
|-------------|--------------|---------|
| M8-R1 — the home page's own template: header retained; left nav, "on this page" panel and prev/next footer nav absent; site footer retained | [#2](https://github.com/radiusred/codecrew-www/issues/2) / PR [#3](https://github.com/radiusred/codecrew-www/pull/3); repaired by [#9](https://github.com/radiusred/codecrew-www/issues/9) / PR [#10](https://github.com/radiusred/codecrew-www/pull/10) | Done; **satisfied**, first round — the mobile drawer explicitly held not to regress it |
| M8-R2 — the page reads as a product page: hero, how it works, why, proof, closing install, full-bleed sections | [#2](https://github.com/radiusred/codecrew-www/issues/2) / PR [#3](https://github.com/radiusred/codecrew-www/pull/3); [#4](https://github.com/radiusred/codecrew-www/issues/4) / PR [#5](https://github.com/radiusred/codecrew-www/pull/5) | Done; **satisfied**, first round — verdicted against the deployed page |
| M8-R3 — blog and every other page unaffected | [#2](https://github.com/radiusred/codecrew-www/issues/2) / PR [#3](https://github.com/radiusred/codecrew-www/pull/3) | Done; **satisfied**, first round — proved by byte-identical built output in both code reviews |
| M8-R4 — both colour schemes render correctly, the build stays warning-free, tests prove template selection and hide rules | all three tasks | Done; **satisfied**, first round |

No requirement was added, amended or dropped after the milestone opened.
Captures filed by the milestone and left open by choice:
[#197](https://github.com/radiusred/gh-codecrew/issues/197) and
[#199](https://github.com/radiusred/gh-codecrew/issues/199).

## Protocol-discipline observations

- **A milestone whose work is entirely in a spoke has no first PR in its
  hub.** `milestone new` appends the ROADMAP row locally and says it "rides in
  this milestone's first PR (the implementer's)"; the implementer contract
  repeats it. M8's implementer had no branch in this repository to ride, so
  the operator pushed the row to `main` directly under the org ruleset's admin
  bypass ([ae09d01](https://github.com/radiusred/gh-codecrew/commit/ae09d01)),
  and GitHub recorded a bypassed rule violation for it. The coordination
  layer cannot do it either: its App holds `contents: read` by design (#119
  finding 29). [#197](https://github.com/radiusred/gh-codecrew/issues/197)
  captures the shape with four options, one of which is that the
  doc-synthesizer should own *both* ends of the row. This document's PR is the
  first hub PR the milestone has had, and it flips the row the operator wrote.
- **Twenty-two changes on one task issue, and the Deviation mechanism carried
  sixteen of them.** Task #4 opened with five changes in its goal and
  finished with twenty-two, seventeen added by coordinator note while its PR
  was open. Thirteen produced a `**Deviation:**` comment naming the note it
  came from; changes 10 and 20 produced Decisions but no Deviation; change 22
  produced neither. The mechanism held under a load nothing in the protocol
  anticipated — the issue is a 61-comment design conversation — but the
  gathering verb sees `**Decision:**` and `**Deviation:**` headings, and three
  plan departures are outside them. An alternative reading of the same record
  is that #4 should have been five tasks; it was not, and the cost is one
  issue nobody will read end to end against the benefit of a design
  conversation that sits in one place with the PR answering it. Both halves
  are on the record; this document does not adjudicate it.
- **The Gates section was left as the scaffold's placeholder.** #196 declared
  no gate beyond CI, and the operator's approval — the thing that actually
  gated the merge of PR #5 — arrived as a sentence at the end of a coordinator
  note about the favicon. M6 and M7 both raised and resolved gates in the open
  with `**Gate raised:**` / `**Gate resolved:**`; M8 had a real human sign-off
  and no gate to record it in.
- **A defect reached production and came back as a task in eleven minutes.**
  The empty drawer was live from PR #3's merge on 2026-09-02 until PR #10 on
  the evening of 2026-09-03 — through the whole styling pass, three review
  submissions and every headless-Chromium check, none of which opened the
  burger. Every verification in the milestone rendered the page at 420px and
  read what was *on* it; none operated a control. The fix arrived as a
  properly-formed task with a plan, two Decisions and its own regression
  tests, one of which now asserts the theme rule the fix depends on so a
  Zensical upgrade fails in CI rather than silently on a phone.
- **The reviewer's one change request was against a claim, not code.** It
  required the PR body, the commit message and the task record to be
  corrected, and nothing else — the review gate acting on the *record* as a
  deliverable, which is what this project claims the record is.
- **The doc-synthesizer was dispatched as an editor mid-milestone, and
  labelled its own output.** The contract has no editorial-review duty; the
  seat was given one, and it opened its comment by stating that nothing had
  been committed and the copy was the operator's to accept or reject. The
  operator then answered seven open questions in one comment and a coordinator
  note turned the answers into a change. Nothing about this is in the protocol
  — it worked because the seat marked the boundary between proposing and
  recording, and because the answers went on the issue rather than into a
  chat.
- **Two captures, both still open, both about things nobody is looking at.**
  [#197](https://github.com/radiusred/gh-codecrew/issues/197) is above.
  [#199](https://github.com/radiusred/gh-codecrew/issues/199) observes that
  this repository's ten crew badge PNGs are referenced by nothing in it: their
  only consumers are the five App avatar settings pages and
  `codecrew-www`'s own committed copies of the five transparent variants. M8
  made the second of those two, so the milestone that created the duplicate
  also filed the capture about it.
- **M9's work landed in the same spoke while M8 was still running.** The suite
  moved 8 → 24 → 49 → 53 across the milestone, and the jump from 24 to 49
  happened between M8's second and third tasks because the docs section
  arrived in between — the reviewer of PR #10 ran the suite with
  `SYNC_SOURCE_BASE` pointing at a hub checkout, and QA's M8-R3 verdict names
  "the later-added Docs section". Two open milestones sharing one spoke is not
  a shape the protocol forbids, and nothing collided; it is worth noting only
  because M8's test counts cannot be read as M8's own without it.
