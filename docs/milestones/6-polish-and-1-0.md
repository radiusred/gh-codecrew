# M6: Polish and 1.0

Tracking issue: [#109](https://github.com/radiusred/gh-codecrew/issues/109) ·
Synthesized 2026-08-29 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #109's amended requirements, its Gates section and its
two QA rounds; the three gate trails; the eighteen hub task issues and their
merged PRs (#110 through #155) and the three in
[radiusred/www](https://github.com/radiusred/www) (#45, #46, #47);
the eighteen findings-log entries and fifty numbered findings on
[#119](https://github.com/radiusred/gh-codecrew/issues/119) — the orchestrator
run, this milestone's largest single record; the four releases v1.0.0–v1.0.3
and [CHANGELOG.md](../../CHANGELOG.md), which narrates each. The gather from
`gh codecrew milestone close 6` is not part of the raw material this time: the
close runs *after* this document merges, and `DOC_MISSING` is what stands
between it and the record. Every Decision and Deviation cited below was read
from its issue directly, which is also the first thing this milestone fixed
(M6-R3).

## Goal and outcome

Take CodeCrew from a proven protocol to a presentable 1.0: branding lands, the
docs get an adoption-facing polish pass, the gatherer bug M5's synthesis
exposed is fixed, and v1.0.0 ships with public announcements — drafted and
posted by wordy behind an operator gate, with the posting credentials
recovered from the archived ig_trader-era repos
([#109](https://github.com/radiusred/gh-codecrew/issues/109)). Nine
requirements, M6-R1 through M6-R9; the last two (launch readiness, a clean
close) were added from a pre-launch scan and a backlog capture after the
milestone opened.

Delivered in full. All nine requirements carry a `satisfied` verdict from the
qa role holder; one of them (M6-R1) is satisfied only against text that was
**amended at a gate** on the day of the verdict, and the document says so
below rather than smoothing it over.

At open, the extension was v0.5.0, the protocol had a version nobody read, the
README was a specification, `milestone close 5` had gathered three of seven
records, and no milestone anywhere had been driven end to end by an
orchestration platform. At close: **four releases** — v1.0.0 (2026-08-26),
v1.0.1, v1.0.2 and v1.0.3 (2026-08-28) — a versioned protocol the CLI checks
and refuses on, a landing-page README and a documentation map, a gatherer
proven against M5's real comment corpus, contracts that extend without
forking, branches that clean themselves up, and **three milestones closed on a
proving-ground repo by a Paperclip company of six agents**, the third with
zero operator touches that were not a gate. The 1.0 announcements are live on
Bluesky and LinkedIn, and the blog article they link to is on radiusred.uk.

The milestone also grew mid-flight. Five of its eighteen tasks (#144, #148,
#151, #153, #155) did not exist when it opened: they are fold-backs from the
orchestrator run, and they are why there are four releases instead of one.

## Decisions

### The protocol has a version, and the CLI refuses on it

The 1.0 question raised before any release work, as a gate on
[#114](https://github.com/radiusred/gh-codecrew/issues/114): does SPEC go 1.0
with the binary, or does the pointer's `codecrew: "0.1"` keep its own cadence?
The measured fact that framed it: `Config.Codecrew` was parsed and never read
by any verb — the field had no job, which is half of why M4's operator read it
as a tool version. Three candidates were enumerated at the gate: **(a)**
lockstep, one number; **(b)** decoupled, stated and *used*; **(c)** drop the
field and let the contract-drift report be the only compatibility signal.

The operator chose (b)
([gate raised](https://github.com/radiusred/gh-codecrew/issues/114#issuecomment-5429682327),
[resolved](https://github.com/radiusred/gh-codecrew/issues/114#issuecomment-5429716312):
"option (b) - Decoupled with refusals"). SPEC 1.0 and v1.0.0 ship together and
then run on separate cadences; the binary declares the protocol it implements
and `version` prints both (`v1.0.0 (protocol 1.0)`); `init` scaffolds
`codecrew: "1.0"`; and every verb that loads the pointer checks its protocol
*major* — a different major refuses `PROTOCOL_MISMATCH`, `"0.1"` is accepted
as the pre-1.0 form of 1.0 with a note, a missing field gets a note.
**Trade-off:** two version numbers for a reader to hold, against one that
would advance with releases the protocol did not change in. **Rejected:** (a)
lockstep, for that reason; (c) dropping the field, which leaves no seam by
which a v2 protocol could refuse a v1 project. Implemented in PR
[#127](https://github.com/radiusred/gh-codecrew/pull/127).

What 1.0 promises came with it, into SPEC §10: within a CLI major, verb names
and flags are additive; a refusal code's meaning is stable — codes may be
added in a minor, never repurposed, removed only in a major; the
`refused[CODE]: detail` line and the `version` output are stable shapes and
other human-facing text is not; pointer fields are additive; embedded role
contracts may change in a minor, with the drift report and `roles diff` as the
mechanism. The commitments were proposed at the gate and are part of what it
resolved — the first time this project has written down what a version number
obliges it to.

### The gatherer's unit is the labelled paragraph, not the comment

M5's synthesis found three of four Decisions on #73 invisible to
`milestone close`: the matcher took only the bare prefixes `**Decision:**`,
`**Deviation:**` and `**Gate resolved:**`, and the parenthetical form used
throughout that milestone defeated it. Baseline measured with the real
gatherer over M5's 22 tasks before anything changed: **3 records**; the
comments held **7**
([plan on #113](https://github.com/radiusred/gh-codecrew/issues/113)).

The [Decision on #113](https://github.com/radiusred/gh-codecrew/issues/113#issuecomment-5429025204):
a paragraph opening `**Decision|Deviation|Gate resolved[ (qualifier)]:**`
starts a record; unlabelled and `**Why:**` / `**Trade-off:**` / `**Rejected:**`
paragraphs after it attach; any other bold label ends it; a label mid-line is
never a record. The qualifier is kept verbatim in a new `Record.Label` and
carries no semantics. **Trade-off:** a comment holding two records now yields
two, and a Deviation buried under a review round-up is gathered without the
round-up — against "one record per comment" as the simpler mental model, which
the quickstart keeps teaching as the practice to prefer. **Rejected:** PR
bodies as a source (SPEC §4 stands — the body is the summary and links to the
records, and the doc-synthesizer contract has the seat read it directly; §4
now says so explicitly, and `milestone close` prints each task's PRs as
pointers); recognising supersession mechanically, which is a record-model
question the text already answers well enough; scanning PR *review* bodies —
checked across all eighteen M5 PRs, no record lives in one.

The regression test is the milestone's most durable artifact:
`internal/tracker/testdata/m5-corpus.json`, M5's real comments (22 issues,
43 + 8 comments, snapshotted 2026-08-26), with `TestExtractRecordsM5Corpus`
asserting all seven records by source, kind and label — including the four the
old matcher missed. After the change, the same live gatherer over the same 22
tasks returns **7** (PR [#126](https://github.com/radiusred/gh-codecrew/pull/126)).

### Contracts extend without forking

SPEC §6 called role contracts "the project's own fork", which conflated two
kinds of drift: the framework contract moving upstream (the drift report's
job) and a project's own additions, which under the old story flagged drift
permanently and so drowned the first signal
([#122](https://github.com/radiusred/gh-codecrew/issues/122)). The seam is
append-only and deliberately dumb: `roles/<role>.local.md`, loaded after the
contract in a fixed order — hub contract, hub extension, spoke extension —
with no merge language, no front-matter and no precedence rules beyond order.
A local file that contradicts its contract is a review finding, not a
resolver's job. `roles show <role>` prints the composition a dispatched agent
actually sees; `--latest` still prints the embedded contract alone; the drift
check keeps comparing only `roles/<role>.md`.

The [Decision on #122](https://github.com/radiusred/gh-codecrew/issues/122#issuecomment-5426799250)
is about the embed, and it came out of review: `assets.go` globbed
`roles/*.md`, so the first extension shipped *inside the binary* — `init`
scaffolded this company's editorial voice into strangers' projects, and drift
counted the extension as a contract (checky's reproduction on PR
[#123](https://github.com/radiusred/gh-codecrew/pull/123)). The embed now
lists the four contracts by name, and `scaffold` and `contractDrift` also skip
any `.local.md` they meet. **Trade-off:** a fifth contract now means editing
`assets.go` — one line, with a test that fails loudly — which is the price of
never shipping one project's extensions inside the binary. **Rejected:**
filtering only in `contractDrift` (leaves `init` scaffolding the voice into
strangers' projects); a separate `roles/local/` directory (moves the file away
from the contract it extends and changes nothing about the embed).

The first extension shipped in the same PR: `roles/doc-synthesizer.local.md`,
this hub's editorial voice for outward-facing writing, recovered from wordy's
Paperclip-era instructions. Every doc-synthesizer dispatch after it — the
landing page, the docs polish, the blog article, the announcements, this
document — loaded it through `roles show`. The routing-table form
(`roles.<role>.instructions: <path>`) was deliberately *not* built, pending
evidence from #119 that an orchestrator needs more than a convention file; the
run instead used exactly this mechanism, putting the platform's own rules into
`roles/<role>.local.md` on the proving-ground repo.

### The docs become a landing page and a map

Two tasks and one operator brief. [#111](https://github.com/radiusred/gh-codecrew/issues/111)
moved the old README whole to `docs/introduction.md` (a byte-identical commit,
so history follows the file) and wrote a landing page in its place: why a
developer would want a crew, four beats with an illustration each, the
receipts, each linked into the record. [#112](https://github.com/radiusred/gh-codecrew/issues/112)
then made the introduction the map — a "Read in this order" index and, closing
a hand-off #111 left open, **the refusal-code catalogue**, every code grouped
by the verb that raises it
([Decision on #112](https://github.com/radiusred/gh-codecrew/issues/112#issuecomment-5428066073)).
**Trade-off:** one page carries three jobs (map, definition, reference) and
runs to ~150 lines, against a two-file split that would leave the reader
deciding which to open first. **Rejected:** a `docs/README.md` as the index —
GitHub renders it when browsing `docs/`, but the www sync would publish it as
a second overview beside the landing page and every later refresh would have
two files to keep in step; a catalogue in SPEC §6 — the protocol names the
gates, not the CLI's wording, and it would date the SPEC with each new
refusal.

The illustrations are four hand-authored SVGs, and the
[Decision on #111](https://github.com/radiusred/gh-codecrew/issues/111#issuecomment-5427538423)
made them **theme-neutral single files** rather than light/dark pairs behind
`<picture prefers-color-scheme>`. **Trade-off:** a mid-tone reads on both
GitHub themes without being ideal on either — measured at ~5:1 on white and
~4:1 on `#0d1117`, above the 3:1 floor for large text and diagrams, below what
a dedicated pair could hit. **Rejected:** `<picture>` pairs, because www's
`sync_docs.py` rewrites `<img src>` but not `<source srcset>`, so on
radiusred.uk the dark source would 404 and only the fallback would render —
and two files per drawing doubles what every later refresh must keep in step;
embedded `@media (prefers-color-scheme)` inside one SVG, which works in `<img>`
but is invisible to the www build's class-driven theme toggle.

Then the operator wrote a brief of his own
([#140](https://github.com/radiusred/gh-codecrew/issues/140)): the docs made
adoption look like the operator's job, when the whole of onboarding is four
lines and one sentence to your agent. The README now leads with
`gh extension install` / `cd` / `gh codecrew init` / `claude` and "Let's build
this project!", and says the thing the old page only implied — *you do not run
the verbs; your agent does*. That brief immediately produced
[#141](https://github.com/radiusred/gh-codecrew/issues/141): Claude Code reads
`CLAUDE.md` and never `AGENTS.md`, so the four-line start began blind on the
harness most adopters reach for. `init` now writes `CLAUDE.md` whose first line
is `@AGENTS.md`, and v1.0.1 was cut so the README described a release that
existed (PRs [#142](https://github.com/radiusred/gh-codecrew/pull/142),
[#143](https://github.com/radiusred/gh-codecrew/pull/143)).

### A close leaves the repo clean

Fifty stale task branches sat on the hub when M6 opened: `task start` cuts a
linked branch and nothing ever removed it
([#129](https://github.com/radiusred/gh-codecrew/issues/129), adopting backlog
capture #128). Three layers, all Free-tier: `task finish` deletes the head it
merged; `milestone close` sweeps the milestone's tasks' branches after every
gate has passed; `status` notes when the repo's `delete_branch_on_merge` is
off, so adopters learn the setting exists.

The load-bearing detail is what counts as safe to delete. Rebase-merge
rewrites commits, so "the tip is an ancestor of `main`" is the wrong test — 43
of the 50 branches swept by hand had non-ancestor tips. The rule is: the
branch's PR merged **and the branch still sits at the PR's frozen
`headRefOid`**, or no PR is open and the branch is zero commits ahead of the
default branch. Where several PRs share a head, the one that forbids deletion
wins (open over merged over closed); forks' heads and the default branch are
never candidates. Everything else is reported and kept, with the reason.

Three Deviations on #129 record how that shape was arrived at, and two of them
correct each other — see *Deviations* below.

## The three gates

The protocol's own gate mechanism ran three times in this milestone, on three
different kinds of question.

**The protocol version (#114).** Raised before any commit, with three
candidates enumerated honestly and a recommendation stated; resolved by the
operator in three minutes. Nothing was committed before the resolution — the
gate did what SPEC §8 says it is for.

**The announcements (www#47).** M6-R5's hard clause: *nothing posts publicly
without an explicit operator gate resolution*. wordy drafted both texts,
attached the dry-run output, and
[raised the gate at 01:56Z](https://github.com/radiusred/www/issues/47#issuecomment-5433381414);
the operator resolved it at
[08:54Z](https://github.com/radiusred/www/issues/47#issuecomment-5436692680)
("**Gate resolved:** new tags"); the posts went out at 08:56Z. Seven hours
where nothing was published and no one worked around it. The resolution
changed the content — the tag lines — which is what a gate is for.

**M6-R1's avatar clause (#109), the one that changed a requirement.** QA's
first round returned M6-R1 *not satisfied*: the assets and the social preview
were delivered, but the requirement also asked for the exported mark to be
uploaded to each crew App's Display information, and no evidence of that
existed — a manual surface GitHub exposes through no API, so QA could not
reconstruct it either
([finding on #110](https://github.com/radiusred/gh-codecrew/issues/110#issuecomment-5463451199)).
The implementer then checked from outside — a GitHub App's avatar is public at
`https://avatars.githubusercontent.com/in/<app-id>` — and confirmed the
finding on the facts: all five crew Apps carry the Radius Red "RR" mark with a
per-role badge, not the CodeCrew mark
([evidence](https://github.com/radiusred/gh-codecrew/issues/110#issuecomment-5463467023)).
Whether it *should* be done was a requirement question, so it went up as a
[gate](https://github.com/radiusred/gh-codecrew/issues/109#issuecomment-5463467756)
with two real options.

The operator chose to amend
([Decision](https://github.com/radiusred/gh-codecrew/issues/109#issuecomment-5463477172),
[Gate resolved](https://github.com/radiusred/gh-codecrew/issues/109#issuecomment-5463478052):
"keep the RR avatars for now"): **the five crew Apps keep the Radius Red mark
with their role badges — they are Radius Red's crew, not the framework's — and
the CodeCrew mark lives on the repository.** #109's body now carries the
amended text with the original clause quoted inside it for the record.
**Rejected:** uploading the CodeCrew mark to each App's Display information —
five manual uploads for branding that would misattribute the crew. QA
re-verdicted against the amended text four minutes later
([superseding verdict](https://github.com/radiusred/gh-codecrew/issues/109#issuecomment-5463494352)),
having fetched all five avatars itself.

This is the first time in six milestones that a requirement has been changed
rather than met, and the mechanism handled it in the open: a verdict of *not
satisfied*, a gate naming the options, an operator decision with its rejected
alternative, an amended body that quotes what it replaced, and a superseding
verdict. Nothing was edited away.

## The releases, and what the CHANGELOG had to do

Four releases in three days, each cut from a merge commit by the implementer
identity, each verified from the *shipped* binary rather than a dev build —
the M4 release-parity discipline, unchanged.

| Release | Tagged on | What it carried |
|---|---|---|
| [v1.0.0](https://github.com/radiusred/gh-codecrew/releases/tag/v1.0.0) | `5cb4581` (PR [#137](https://github.com/radiusred/gh-codecrew/pull/137)) | protocol 1.0 declared, printed and checked; local contract extensions; the whole record gathered; branch cleanup; launch readiness; the landing page and the introduction |
| [v1.0.1](https://github.com/radiusred/gh-codecrew/releases/tag/v1.0.1) | `8e0d019` (PR [#143](https://github.com/radiusred/gh-codecrew/pull/143)) | `init` writes `CLAUDE.md`; the README's four-line start |
| [v1.0.2](https://github.com/radiusred/gh-codecrew/releases/tag/v1.0.2) | `d26ea35` (PR [#145](https://github.com/radiusred/gh-codecrew/pull/145)) | `milestone close` refuses `NO_REQUIREMENTS` instead of passing vacuously |
| [v1.0.3](https://github.com/radiusred/gh-codecrew/releases/tag/v1.0.3) | `39fd293` (PR [#154](https://github.com/radiusred/gh-codecrew/pull/154)) | `GH_TOO_OLD`; `milestone new --requirement`; `gh codecrew` in the contracts; the qa seat's hairbrush |

v1.0.0's [evidence comment](https://github.com/radiusred/gh-codecrew/issues/115#issuecomment-5432315375)
is the shape the release tasks have taken since M5 — no PR, closed directly,
every claim re-proved from the installed extension. It is worth reading for
one thing in particular: `milestone evidence 6` **refused first**. The initial
run returned `EVIDENCE_UNREACHABLE` for a `…/blob/main/…` citation — a literal
ellipsis in #131's plan that the extractor read as part of a URL. The record
was repaired by rewording and the extractor's greed captured as #138. The verb
was right about the record, and wrong about the URL, and both facts are in the
evidence comment.

The repo change the release needed rode its own task
([#136](https://github.com/radiusred/gh-codecrew/issues/136) / PR #137), which
is where **CHANGELOG.md** was born: Keep-a-Changelog shape, entries grouped by
what each change means to an adopter rather than by commit, every merged PR
since v0.5.0 appearing exactly once, and "What 1.0 promises" restated from
SPEC §10. It is deliberately *not* an upgrade guide — an operator direction
recorded on #136 and #131: there are no 0.x adopters to migrate. The 1.0.0
entry doubled as the GitHub release's notes, and each release since has been
cut from a PR that flipped `[Unreleased]` to the new version in the same
diff — so the changelog is never behind the tag.

## The proving ground: an orchestrator drives three milestones

[#119](https://github.com/radiusred/gh-codecrew/issues/119) is M6-R6 and the
milestone's largest record: eighteen findings-log entries, fifty numbered
findings, three cycles, and five hub releases' worth of fold-back. SPEC §1 and
§9 had claimed since M3 that "orchestrators such as Paperclip assign the roles
to agents"; nobody had run it. The operator's Paperclip instance hired four
fresh agents under a CEO agent, told them only to read the CodeCrew docs and
build something on a brand-new `radiusred/numberguess`, and the run was
observed from both sides — the protocol's record on GitHub, the platform's
tickets read over its API.

**Cycle 1 began with the framework not reaching the agents at all.** The
platform's managed job descriptions load first in every run, and Bossy's
hiring plan had written them from a template that predates CodeCrew: "a
four-agent engineering team proving that a real build → review → test → docs
loop runs cleanly", with boundaries that contradicted the seats one by one —
Testy told to "write and run the tests that gate merges", Cody told not to
"own the test suite", Wordy scoped to player-facing copy
([finding 15](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5445767722)).
The repo's `AGENTS.md` and `roles/` are reached only by agents that open the
repo, and the coordination layer never opens it. Everything earlier — a CEO
that delegated `init` to the implementer, an implementer that then routed its
own judge, a PR with no milestone, no task, no plan and its Decision in the PR
body — is that one finding seen from different angles. The run's central
finding is
[finding 7](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5443566637):
**CodeCrew's gates bind agents that use its verbs; an orchestrator that keeps
its own task system around the repo bypasses all of them without noticing.**

The intervention was to install the composed contracts as the agents' managed
instructions
([Decision, operator](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5445852159)),
with the trade-off recorded in the operator's own words: *"these instructions
essentially cause CC to take over the entire company ethos. Appropriate for
this test, likely not the right shape for a real Paperclip setup."* The design
note attached to that Decision is what the fold-back actually shipped: the
onboarding artifact should be an **overlay**, not a replacement — and cycle 3
used `roles/<role>.local.md`, M6-R7's mechanism, to do exactly that.

The contracts started working on contact. Within one heartbeat the coordinator
read the routing table, concluded on its own that "Testy → qa is
`contents: read` by design", and corrected the seat mis-mapping unprompted.
The qa seat ran its contract's mandatory first act, met
`refused[NOT_FOUND]: no open milestone M1`, and **refused to work without a
milestone, telling the CEO to open one**
([finding 17](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5446038074))
— the protocol's pressure coming from inside the crew for the first time. And
the permission model did the contract's job mechanically: the qa App's
`contents: read` refused a push at the first attempt, which is the seat
separation enforced by the App's permission set rather than by an agent's
judgement
([finding 14](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5445671135)).

Three themes account for nearly every stall, and each has a fold-back:

- **Identity and credentials** (findings 2, 6, 10, 12, 13, 16, 18, 22, 35).
  Three different namings for the same two secrets across two agents; a
  one-hour installation token persisted into a shared `gh` config with
  `gh auth login` and inherited, dead, by the next agent; and **three of three
  agents escalated a 401 they could have re-minted through** — one of them
  recovering 67 seconds after the operator asked "are you unable to mint
  tokens from the PEM file in your environment?". Knowledge was never the gap;
  the recovery reflex was. The coordination layer, meanwhile, had no identity
  at all and the framework had never said it needed one — the operator minted
  a fifth App to fit
  ([finding 16](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5445972756)).
- **The review loop's second round** (findings 20, 26, 38–41, 46, 48, 50).
  The protocol has exactly one cycle — changes requested → fix → re-review —
  and an orchestrator must own it explicitly or the second round never starts.
  Every operator touch in cycle 2 after the first was a substitute for a wake
  the platform's ticket graph does not provide: board fields wake no one, a
  child ticket completing wakes no one, and a bare `@Bossy` is text. The wake
  is a link-form mention or an assignment
  ([findings 39](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5455026917)
  and [41](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5455092737)).
- **Gates that bind only when a verb is run** (findings 28, 36). The
  coordinator learned what the close would refuse and skipped to the remedy
  without ever invoking it, so the record shows no refusal for a gate that was
  in fact respected. And once, worse: `milestone close 1` **passed vacuously**
  over a `not satisfied` verdict, because the requirements had been written
  under `## Goal` and `RequirementIDs` counted zero — nothing missing, nothing
  unsatisfied, pass
  ([finding 28](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5451814672)).
  Every agent counted six requirements; the CLI counted none and said nothing.
  That is the run's most important finding and the first thing fixed.

**Cycle 3 is the one M6-R6 is met by.** With the contracts layered rather than
replacing the platform's files, the coordinator's App webhook wired to a
Paperclip routine, and three fixes applied mid-cycle, the loop ran on GitHub's
own events: every PR opening, review and merge reached the coordinator within
three seconds. The feature task went from PR to merge in twenty-three minutes
with a real changes-requested round inside it and **no hand on it at all**
([entry 16](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462379153)).
The cycle's one human moment was a *gate*: the coordinator met
`refused[DOC_MISSING]` over a document that existed under the wrong number,
raised `gh codecrew checkpoint` with three options and a recommendation, and
the operator resolved it — SPEC §8's mechanism exercised end to end on a
platform for the first time
([the close](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462686479),
[resolved](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462791038)).

Three findings from cycle 3 are worth carrying into any future run:

- **Requirement prose is executed, not interpreted.** The coordinator wrote
  M4's requirements assuming the milestone would be M3; the verb assigned 4
  after the fact, because a milestone closed as a duplicate still consumed the
  ordinal. Four seats then built faithfully to text that named the wrong
  number — including the doc-synthesizer, which named the file
  `3-selectable-difficulty.md` and recorded the trade-off as a Decision, and
  the qa seat, which credited the wrong ROADMAP row twice
  ([findings 45](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462078800)
  and [47](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462305708)).
  Nothing could refuse it: the IDs were right and only the prose disagreed.
- **A coordinator woken by two paths races itself; a single-flight one goes
  stale.** One PR opening produced three concurrent coordinator runs and three
  reviewer tickets; setting `maxConcurrentRuns: 1` turned the race into a
  queue of wakes that were stale by the time they ran
  ([findings 46](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462194623)
  and [48](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462312601)).
  The fix is one wake path per transition, a re-read at the act rather than at
  the wake, and short runs.
- **The reviewer seat caught the doc-synthesizer inventing the record.** On
  the cycle-3 milestone document, checky found that two Decisions "invent
  rejected alternatives and rationale that do not exist in the cited record",
  and named the remedy the contract names — remove them, or label them as
  undocumented and inferred
  ([finding 49](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462506456)).
  The invented text was *plausible*, which is precisely why the record and not
  plausibility is the standard. This document is written under that finding.

The run's own numbers, from
[entry 18](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462998256)
(reported cost covers the Claude-hosted seats; the Codex seats report tokens
and no cost):

| cycle | milestone | PRs | wall-clock | runs | $ reported | runs / PR | non-gate touches | gates |
|---|---|---|---|---|---|---|---|---|
| 1 (08-27/28) | M1 #3 | 3 (+ scaffold) | ~19.5 h | 85 | 68.49 | 28 | 5 | 0 |
| 2 (08-28) | M2 #9 | 2 | 4 h 01 | 30 | 26.62 | 15 | 6 | 0 |
| 3 (08-29) | M4 #15 | 4 | 3 h 08 | 97 | 82.76 | 24 (13 after the fixes) | 0 | 1 |
| **run** | 3 closed | 9 | ~27 h | **212** | **177.87** | | **11** | **1** |

The operator's macro read, recorded verbatim at the end of cycle 1 and still
the frame the run ended on: *"CC was designed to solve an overly burdened GSD
framework — paperclip is the GSD equivalent in orchestration… It feels like a
culture clash right now."*
([operator's read](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5451965158)).
By cycle 3 the number that carried it was measurable: **75% of the cost was
the coordinator**, on a stock CEO checklist, most of it in wakes that found
nothing. The conclusion — recorded as an
[operator Decision mid-cycle](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462401564)
— is that the coordination layer should be **its own agent with its own lean
loop**: the CEO keeps the CEO's job, and a dedicated coordinator holds the
seat. Loopy was hired at 13:50Z on the last day and closed the milestone
half an hour later, at about $0.50 a run against the CEO's $1.22. The interop document
that inventory feeds is #54; it is not written here.

### What the run folded back

Five hub PRs, two releases, and a set of captures — all of it recorded as it
happened rather than at the end:

| Finding | Task | PR | Release |
|---|---|---|---|
| 28 — the close passed vacuously | [#144](https://github.com/radiusred/gh-codecrew/issues/144) | [#145](https://github.com/radiusred/gh-codecrew/pull/145) | v1.0.2 |
| 3–6, 9, 10, 12–14, 21, 22, 27, 30, 31 — wording across contracts, scaffold and docs | [#148](https://github.com/radiusred/gh-codecrew/issues/148) | [#150](https://github.com/radiusred/gh-codecrew/pull/150) | v1.0.3 |
| 19a, 28, 32 — requirements land where the parser reads them | [#151](https://github.com/radiusred/gh-codecrew/issues/151) | [#152](https://github.com/radiusred/gh-codecrew/pull/152) | v1.0.3 |
| 21, 30 — the `gh` floor is a refusal, not a stack trace | [#153](https://github.com/radiusred/gh-codecrew/issues/153) | [#154](https://github.com/radiusred/gh-codecrew/pull/154) | v1.0.3 |
| 37 — the qa seat judges the suite | [#155](https://github.com/radiusred/gh-codecrew/issues/155) | [#156](https://github.com/radiusred/gh-codecrew/pull/156) | v1.0.3 |
| 12, 35, and the rung itself | [#119](https://github.com/radiusred/gh-codecrew/issues/119) | [#160](https://github.com/radiusred/gh-codecrew/pull/160) | unreleased (docs) |

Two of those deserve their own line. **#155** rewrote `roles/qa.md` from the
operator's read of the run's QA leg — *"Testy simply ran the tests and
verified passing. There's no real value in that… Testy should be earning his
keep by validating that the tests are appropriate"*
([recorded verbatim](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5453862489)).
The contract now says the rerun is the floor, not the product: per requirement,
what the shipped suite proves and what it assumes, gaps being findings even
when the behaviour is right, and **every verdict cites at least one probe the
suite does not enumerate** — the operator's rule of thumb, *a QA engineer
walks into a bar and orders a drink; then −1 drinks; then 999,999 drinks; then
a hairbrush*. M6's own QA round ran under that contract, and every one of its
nine verdicts names its hairbrush.

**#148**'s scoping is itself a recorded Decision: the in-CLI `gh` version check
was *not* folded into it, because a refusal code moves the catalogue, the
README count and the SPEC table together — a code change with its own tests,
not a wording fold-back
([Decision on #148](https://github.com/radiusred/gh-codecrew/issues/148#issuecomment-5453252300)).
**Trade-off, stated:** the crew that hit the floor stayed one parse error away
from the same stall until #149 landed; the prerequisites at least documented
it in the meantime. It landed four hours later as #153.

## The announcements: archaeology, an article, and a gate

M6-R5 ran across three tasks in the www spoke, and it began with a dig. The
[archaeology log on www#45](https://github.com/radiusred/www/issues/45#issuecomment-5432998890)
records locations and names only, never values: four archived repos searched
in full history for posting code — **none was ever committed to any of them**;
posting had been done ad hoc from wordy's Paperclip runtime. The one surviving
trace was a secrets inventory naming the credentials and their homes, and the
values themselves were still in the platform's local-encrypted vault, company-
scoped, having outlived the agents that used them. **Decision: recover, don't
re-provision.** **Trade-off:** decrypting the vault cost one query and a
thirty-line script; re-provisioning would have cost a browser consent, a new
app password, and the feed scopes a May re-consent had added. **Rejected:**
treating the tokens as dead on age alone — introspection was cheap, and the
LinkedIn refresh token had nine months left.

What replaced the ad-hoc runtime is a `social/` package in the www repo:
stdlib-only, credentials from the environment or `~/.config/codecrew/social.env`
(never the tree), with `check`, `post --dry-run` and an `auth linkedin`
re-consent playbook so the browser leg never depends on a private note again.
Three further Decisions on www#45 record the credential home and its
write-back rule (LinkedIn rotates the refresh token on every refresh, so
without write-back the second refresh fails), the no-dependencies choice, and
the pinned API version — `202608`, so a stale pin fails loudly with
`426 NONEXISTENT_VERSION` rather than silently. A fifth, from the operator
during review, added faceted `#hashtags` to Bluesky: a tag without its facet
is plain text, and for a small account faceted tags are the only route to an
audience.

[www#46](https://github.com/radiusred/www/issues/46) is the article the
announcements link to, delivered through the spoke's own task flow and
[live on radiusred.uk](https://www.radiusred.uk/blog/posts/2026-08-27-codecrew-1-0-the-crew-is-staffed-and-the-gates-have-teeth/).
Its worked example was chosen deliberately
([Decision](https://github.com/radiusred/www/issues/46#issuecomment-5433250557)):
www#45's own review round, where the reviewer seat caught a token write-back
bug **in the very tool that would post the article's announcements**.
**Rejected:** an example from the hub's own milestones, already covered by the
introducing post.

[www#47](https://github.com/radiusred/www/issues/47) posted them, behind the
gate described above. Two Decisions shaped the texts: both posts lead with the
**project home** and carry the article second, on the operator's instruction
that the repository is the product and the article the explanation; and the
tag lines were chosen from **measured** Bluesky search volume and LinkedIn's
published community sizes rather than instinct, with `#CodeCrew` dropped from
both — a brand tag with no community spends one of five slots for nothing
today. The texts as posted, the tag research and the resulting URLs are
committed in `announcements/2026-08-27-codecrew-1-0.md` on
[PR #50](https://github.com/radiusred/www/pull/50), with the
[Bluesky](https://bsky.app/profile/radiusred.bsky.social/post/3mu2i2animl2a)
and
[LinkedIn](https://www.linkedin.com/feed/update/urn:li:share:7498664546795126784/)
posts recorded as evidence.

The fourth www task, [#51](https://github.com/radiusred/www/issues/51) — the
article about the orchestrator run — was **detached from this milestone** by
[operator Decision](https://github.com/radiusred/gh-codecrew/issues/109#issuecomment-5463407685)
on the day of the close. Cycle 4 is planned on a fresh proving-ground repo, and
it would change the article's story and numbers; writing it from cycles 1–3
alone would need a follow-up post, and holding M6 open for cycle 4 would tie a
release milestone to an experiment. **Rejected:** both of those. M6 closes on
its nine requirements; the article follows cycle 4.

## Launch readiness, from a stranger's scan

[#131](https://github.com/radiusred/gh-codecrew/issues/131) is M6-R9, and it
exists because the operator dispatched a fresh-context Codex session
(gpt-5.6-sol) with a "what have we missed?" brief before the release. Its
findings are
[recorded verbatim on the issue](https://github.com/radiusred/gh-codecrew/issues/131#issuecomment-5431924370)
— verdict: *ship after fixing 4* — and reading them is the cheapest tour of
what a 1.0 owes its first adopter. The two that mattered most:

- **`init` produced a self-incomplete project.** The scaffolded `AGENTS.md`
  pointed an agent at "the hub's `docs/identities.md`" and the implementer
  contract at `scripts/codecrew-token` — files the adopter's repo does not
  have. Every such reference is now a resolvable upstream URL, with
  `TestScaffoldReferencesResolveUpstream` scaffolding into a temp dir with the
  real embed and failing on any hub-relative reference.
- **SPEC promised a gate the binary does not run** — `task finish` "verifies…
  deviations referenced in the PR body have recorded comments". Nothing reads
  the PR body. The promise was removed rather than implemented.

Alongside those: `--help` exits 0 on every verb (it had been printing usage
and then failing); `init`'s next-steps line says what it actually wrote;
github.com-only is stated in SPEC §12 rather than implied to be enterprise-
ready; the quickstart names working PR CI as the one prerequisite, since
`NO_CHECKS` has no override and `init` writes no workflow; and the root gained
CONTRIBUTING.md, SECURITY.md and a teardown inventory. One of the scan's
findings — an upgrade story for 0.x adopters — was excluded by operator
direction on the ground that no 0.x adopters exist; the others became captures
(#132, #133, #134). All of it landed in PR
[#135](https://github.com/radiusred/gh-codecrew/pull/135), **before** the
release, which is the only order in which a pre-launch scan is worth running.

## The close that has not run yet

`gh codecrew milestone close 6` is M6-R8's first live proof, and it runs
*after* this document merges — `DOC_MISSING` refuses until the file is on the
default branch, which is the gate working as designed. What it is expected to
do, from the state QA verified on 2026-08-29:

- **Sweep two branches.** At verdict time the repository carried `main` and
  `task/114-decide-the-protocol-version-story-before`, whose tip exactly equals
  merged PR #127's frozen `headRefOid` — so the sweep should delete it with
  `deleted (PR merged)` and never touch `main`. This document's own branch,
  `task/161-…`, is deleted earlier and separately, by `task finish` at the
  moment it merges.
- **Print what it removed** as part of the close's record, with the count in
  the closing comment; anything it cannot safely delete is reported with its
  reason rather than removed.
- **Gather the record** — and this is the gather M6-R3 was for: it should now
  find the qualified `**Decision (operator, …):**` comments on #109, #144 and
  #148, and the mid-comment Deviation shapes, all of which the M5-era matcher
  would have dropped.

`status` already notes that the hub does not delete branches on merge, which
is the third layer: adopters learn the setting exists rather than being
silently depended on for it.

## Deviations

Eleven were recorded, on five task issues. Most are small and honest; three
are worth the space.

**The merge that beat the merge point** ([#110](https://github.com/radiusred/gh-codecrew/issues/110#issuecomment-5424816286)).
PR #121 was merged by GitHub auto-merge on the reviewer App's counted approval
and green CI — *before* `task finish 110` ran and before the plan's last step
(the social-preview upload) had landed. `Closes #110` then closed the task with
a requirement step outstanding. **Why:** the implementer armed
`gh pr merge --auto --rebase` at PR-open, a habit from before a reviewer App's
approvals counted; with a counted approval, auto-merge *is* a merge gate, and
it fired first. This is the mirror of M5's finding 10 — a merge that happened
around the protocol's one merge point, this time ahead of it rather than after
its refusal. The task was reopened until the evidence landed, and the lesson —
*do not arm auto-merge on task PRs; `task finish` is the merge* — held for the
remaining seventeen tasks.

**A Deviation that superseded its own predecessor** ([#129](https://github.com/radiusred/gh-codecrew/issues/129#issuecomment-5431479212)).
The branch sweep was planned around GitHub's linked-branch relation. The first
Deviation reported, from live checks, that the relation is dropped when the
issue closes — and closed tasks are the whole target of the sweep. The second
corrected it: the relation is dropped once a PR *attaches* to the issue, not at
close; the earlier observation had been taken minutes before the PR opened. The
candidate set became a three-source union (the relation, the head refs of the
task's PRs, and the conventional `task/<n>-<slug>` name), and the same review
round closed three deletion paths that could have removed unmerged work. Both
records stand; the second names exactly what it corrects.

**A record about how records are written** ([#122](https://github.com/radiusred/gh-codecrew/issues/122#issuecomment-5426860230)).
A Deviation about the extension marker's wording was first written *inside* a
Decision comment, then reposted as its own comment "because the gatherer files
one record per comment by its opening label". It was true at the time and #113
made it untrue three hours later — the deviation is now gathered either way.
A small thing, but it is the record noticing its own plumbing, which is what
M6-R3 was about.

The remaining eight: the transparent mark is for dark grounds only, not "light
or custom" as the plan said (white ink over alpha — on white, only the cyan C
survives, found by compositing during review); the introduction's nine broken
relative links were repaired in #111 rather than #112, because a landing page
whose first "read next" leads to 404s is a defect of *that* page; the
doc-synthesizer contract's freshness obligation was widened to name
`docs/introduction.md`, since the perishable claims had just moved there and no
contract covered them; the SVG monospace labels are 13px rather than the
promised 14px floor, recorded rather than bumped because 14px overruns the
boxes across the fonts GitHub may pick; README.md was touched by #112 though
its plan did not list it, because a landing-page claim that PR falsified could
not outlive the PR; the refusal catalogue's first draft filed `NOT_FOUND` under
`task start`, a grep-by-line read that missed the enclosing function, corrected
after the reviewer mapped every call site to its function; a `PRInfo` decode
test was dropped for want of a seam that would have meant restructuring
`internal/tracker/github.go` on the side, with the fields verified live
instead; and on www#45, the first review round was dispatched under an Opus
session although the routing table declares `harness: codex` for the reviewer
seat — the round's findings stand, the session was stopped mid re-review, and
the approval round ran under the declared harness, because de-correlated
judgment from the implementer's model family is the point of the seat.

## QA: one round, one supersession

The qa role holder (`radiusred-testy[bot]`) verdicted all nine requirements on
2026-08-29 against merged `main`, under the contract this milestone had
rewritten eighteen hours earlier: green tests are the first sentence, never the
evidence, and every verdict cites a probe the suite does not enumerate. Eight
came back **satisfied** in the first round; M6-R1 came back **not satisfied**,
went to the gate described above, and was superseded to satisfied against the
amended requirement.

The hairbrushes are what make the round worth reading. QA fetched the four
committed assets and the social preview over HTTP and compared them byte for
byte with the repository, then tried the analogous external check on the App
avatars and reported that the surface was not observable — the finding that
raised the gate. It fed the gatherer
`**Decision (operator (release 🧹), 2026-08-29):**` — nested parentheses and an
emoji — and got exactly one record with its label preserved. It stamped the
CLI with an unenumerated `v1.0.3+qa.hairbrush` build-metadata value and
confirmed the version line printed it verbatim. It ran an `init` in a fresh
repository **with no remote** and fetched every upstream reference the scaffold
wrote, at HTTP 200. It gave `roles show` an empty local extension (skipped) and
a populated one for a role that does not exist (no phantom role). It queried
the proving-ground repository directly to confirm #119's account of the
orchestrator run. And it got HTTP 200 from the live Bluesky post, the LinkedIn
post and the linked article **without posting anything**.

One QA finding did not affect a verdict and is recorded here because the
exchange it produced is useful: the prescribed source-build command in the QA
dispatch brief was `go build -o gh-codecrew .`, and the repository root is the
embedded-assets library, not a `main` package, so it produced an archive
([finding](https://github.com/radiusred/gh-codecrew/issues/115#issuecomment-5463451382)).
The documented command — `go build -o gh-codecrew ./cmd/codecrew`, the same
target `scripts/build-extension` uses — is correct, and QA re-ran against it
successfully once pointed at it. The
[response](https://github.com/radiusred/gh-codecrew/issues/115#issuecomment-5463466963)
records the actual lesson: **dispatch briefs should cite the documented
command rather than restate it.** No task was raised.

## Requirement outcomes

| Requirement | Delivered by | Outcome |
|-------------|--------------|---------|
| M6-R1 — branding lands: exports, README, social preview (avatar clause amended at the gate) | [#110](https://github.com/radiusred/gh-codecrew/issues/110) / PR [#121](https://github.com/radiusred/gh-codecrew/pull/121); [gate](https://github.com/radiusred/gh-codecrew/issues/109#issuecomment-5463467756) and [amendment](https://github.com/radiusred/gh-codecrew/issues/109#issuecomment-5463477172) on #109 | Done; not satisfied → **satisfied** ([superseding verdict](https://github.com/radiusred/gh-codecrew/issues/109#issuecomment-5463494352), against the amended text) |
| M6-R2 — docs polish for adopters, one navigable story | [#111](https://github.com/radiusred/gh-codecrew/issues/111) / PR [#124](https://github.com/radiusred/gh-codecrew/pull/124); [#112](https://github.com/radiusred/gh-codecrew/issues/112) / PR [#125](https://github.com/radiusred/gh-codecrew/pull/125); [#140](https://github.com/radiusred/gh-codecrew/issues/140) / PR [#143](https://github.com/radiusred/gh-codecrew/pull/143) | Done; satisfied, first round |
| M6-R3 — the record gathers completely, with a corpus regression test | [#113](https://github.com/radiusred/gh-codecrew/issues/113) / PR [#126](https://github.com/radiusred/gh-codecrew/pull/126) | Done; satisfied, first round |
| M6-R4 — v1.0.0 ships, protocol-version story decided first | [#114](https://github.com/radiusred/gh-codecrew/issues/114) / PR [#127](https://github.com/radiusred/gh-codecrew/pull/127); [#136](https://github.com/radiusred/gh-codecrew/issues/136) / PR [#137](https://github.com/radiusred/gh-codecrew/pull/137); [#115](https://github.com/radiusred/gh-codecrew/issues/115) (release, no PR) | Done; satisfied, first round |
| M6-R5 — the announcements, behind an operator gate | www [#45](https://github.com/radiusred/www/issues/45) / PR #48; www [#46](https://github.com/radiusred/www/issues/46) / PR #49; www [#47](https://github.com/radiusred/www/issues/47) / PR #50 | Done; satisfied, first round |
| M6-R6 — driven end to end by an orchestrator | [#119](https://github.com/radiusred/gh-codecrew/issues/119) / PR [#160](https://github.com/radiusred/gh-codecrew/pull/160); fold-backs #144/PR #145, #148/PR #150, #151/PR #152, #153/PR #154, #155/PR #156 | Done; satisfied, first round |
| M6-R7 — contracts extend without forking | [#122](https://github.com/radiusred/gh-codecrew/issues/122) / PR [#123](https://github.com/radiusred/gh-codecrew/pull/123) | Done; satisfied, first round |
| M6-R8 — a close leaves the repo clean | [#129](https://github.com/radiusred/gh-codecrew/issues/129) / PR [#130](https://github.com/radiusred/gh-codecrew/pull/130) | Done; satisfied, first round — the live `milestone close 6` sweep runs after this document merges |
| M6-R9 — launch readiness | [#131](https://github.com/radiusred/gh-codecrew/issues/131) / PR [#135](https://github.com/radiusred/gh-codecrew/pull/135); [#141](https://github.com/radiusred/gh-codecrew/issues/141) / PR [#142](https://github.com/radiusred/gh-codecrew/pull/142) | Done; satisfied, first round |

Supporting work with no requirement ID of its own: the four releases (#115
directly; v1.0.1, v1.0.2 and v1.0.3 recorded as release evidence on #141, #144
and #119), and the captures the milestone filed rather than fixed inline —
#132 (`identity token`), #133 (`milestone new --dry-run`), #134 (GHES), #138
(the evidence extractor's URL greed), #139, #146, #147, #149, #157, #158, #159,
and #54, the interop document the orchestrator run is raw material for.

## Protocol-discipline observations

- **The gatherer fix is proven, and this document still did not use it.**
  M6-R3's whole point was that `milestone close` should hand the synthesizer
  the milestone's records. It now does — 3 records to 7 on the M5 corpus, with
  a regression test. But `milestone close 6` cannot run until this file is on
  `main`, so the M6 document, like the M5 one, was synthesized by reading every
  issue directly. The fix will first pay off for M7. The ordering is inherent
  to `DOC_MISSING`, not a defect, but it is worth stating that the gather has
  never yet been the raw material for the document it is gathered for.
- **A requirement was amended to match reality, in the open.** M6-R1's avatar
  clause asked for something the platform has no API for and the operator did
  not want done by hand. The honest paths were to do five manual uploads or to
  change the requirement; the second was chosen, at a gate, with the original
  clause quoted inside the amended text and the qa verdict re-run against the
  new wording. Contrast M5's stale-plan finding: there, a plan asserted
  delivery that had not happened and QA caught it. Here the requirement itself
  was wrong, and the mechanism that caught it was the same one.
- **Auto-merge is a merge gate, and it beat the merge point.** The #110
  Deviation is the milestone's clearest process failure and it happened in its
  first hour. Nothing in the CLI prevents it — `task finish` is the protocol's
  merge point, not GitHub's — so the guard is a habit, and habits are exactly
  what the contracts are for. It has not recurred, but nothing mechanical
  stops it recurring.
- **Half this milestone was written by its own experiment.** Five of eighteen
  tasks and two of four releases came out of #119, a task whose deliverable was
  supposed to be *findings*. That is the fold-back loop working, and it is also
  a warning about scope: the milestone stayed open three days longer than its
  first nine tasks needed, and the decision to detach www#51 rather than wait
  for cycle 4 is the operator drawing that line explicitly.
- **The reviewer seat caught something in almost every round.** On the hub:
  the extension shipping inside the binary (#123), three deletion paths that
  could have removed unmerged work (#130), a catalogue entry filed under the
  wrong verb (#125), a start block promising agent behaviour the CLI cannot
  enforce (#143), a requirement tag that was not exact (#145), a new check that
  bypassed `status` and `roles` (#154). On the proving ground: a quick-start
  that never terminates, a CI detector that could not detect, a command that
  fails in a fresh checkout, and a milestone document that invented its own
  rationale. Six milestones in, no round of clean-context review has yet been
  wasted.
- **Two small inconsistencies in the record itself.** PR #127's description
  says `codecrew version` prints `v1.2.0 (protocol 1.0)` where it means
  `v1.0.0` — a typo in a merged summary, harmless but the sort a synthesizer
  reading only PR bodies would propagate. And CHANGELOG.md dates 1.0.0 as
  2026-08-27 while the release was published at 23:27Z on 2026-08-26 — the tag
  went out half an hour before the date its entry carries. Neither
  changes anything; both are the class of drift that only a reader comparing
  two sources ever finds.
- **The record of the orchestrator run is honest about what it did not do.**
  #119's plan promised the platform's tickets committed as curated transcripts
  beside the protocol's record on the proving ground, as the numberguess
  session transcripts were in M4 and M5. Entry 18 states plainly that this has
  not happened: the export is prepared, where it lands is the operator's call,
  and the entry records that it has not landed. The findings quote what the
  fold-backs needed, and the instance is still live — but the run is not fully
  reconstructible from the repository alone, which is the same class of gap M5
  recorded when a throwaway probe repository was deleted.
