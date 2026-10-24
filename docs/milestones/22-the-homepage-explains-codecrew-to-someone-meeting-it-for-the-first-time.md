# M22: The homepage explains CodeCrew to someone meeting it for the first time

Tracking issue: [#420](https://github.com/radiusred/gh-codecrew/issues/420) ·
Synthesized 2026-10-24 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #420's **three** requirements as opened, **none added,
none struck and none amended** — the issue's `userContentEdits` holds the
original and one revision, fourteen seconds later, which added the two lines
under Gates and nothing else; its **five comments** — the coordination layer's
**one Deviation** and the operator's **one gate resolution**, both posted by
the operator's account, the qa role holder's **two** verdict comments, and
this record's **one Deviation**; the **three** task issues, all in the spoke
[radiusred/codecrew-www](https://github.com/radiusred/codecrew-www)
([#69](https://github.com/radiusred/codecrew-www/issues/69),
[#71](https://github.com/radiusred/codecrew-www/issues/71),
[#73](https://github.com/radiusred/codecrew-www/issues/73)); their **three
merged pull requests** and the **seven commits** on them, two of which carry
forward the **two** commits pushed to `task/69` after its pull request merged;
the **eight Decision records, one `**Gate resolved:**` record and two
Deviation records** across the milestone issue and its three tasks — all eight
Decisions on the tasks, and the gate resolution and both Deviations on the
milestone issue; the **seven review submissions** on the three pull requests —
**three approvals, every one of which stands**, **three change requests** from
the reviewer seat's App, and **one change request submitted under the
operator's account and dismissed**; the **one** adopted backlog capture the
first task's finish closed; and the **one** capture filed in the milestone's
window, open. The house form is [M19's
document](19-the-front-page-tells-the-crew-story.md), whose milestone also
changed codecrew.works alone, with
[M21's](21-v2-1-1-what-the-fleets-first-fortnight-on-2-1-taught.md) behind it;
the record standard is the one the reviewer set on [PR
#322](https://github.com/radiusred/gh-codecrew/pull/322) and re-applied since,
on [PR #383](https://github.com/radiusred/gh-codecrew/pull/383) and [PR
#413](https://github.com/radiusred/gh-codecrew/pull/413) among others, so
every count is re-derived from the source — records by the rule
`tracker.ExtractRecords` applies (a paragraph-initial `**Decision:**`,
`**Deviation:**` or `**Gate resolved:**` label, bare or parenthetically
qualified), review states from each pull request's reviews API and timeline,
commits from each pull request's commits API and from the spoke's compare API,
and the requirement history from the milestone issue's `userContentEdits`.

**Counted at 2026-10-24T02:46:50Z, against a trail frozen from this document's
dispatch.** The last post on the milestone's issues before that instant is
this record's own
[Deviation](https://github.com/radiusred/gh-codecrew/issues/420#issuecomment-6378919539)
at 02:46:23Z, and it is inside every total above; the one before it is the qa
role holder's re-verification of M22-R3 at 02:40:24Z.

**This record has no task.** It is delivered as a housekeeping pull request
under SPEC §4's milestone-record clause, as
[M18's](18-the-protocol-scales-down.md),
[M19's](19-the-front-page-tells-the-crew-story.md),
[M20's](20-the-2-1-post.md) and
[M21's](21-v2-1-1-what-the-fleets-first-fortnight-on-2-1-taught.md) were: the
milestone issue is its charter and its commit references #420. Its record of
its own is the one Deviation on the milestone issue.

**The whole milestone happened in the spoke.** `main` in this repository did
not move between M21's record
([0d16fff](https://github.com/radiusred/gh-codecrew/commit/0d16fff), committed
2026-10-11T00:07:54Z) and this document's pull request. QA's M22-R2 verdict
reports the same for the protocol's own files: "As a cross-cutting probe, `git
log 0d16fff..origin/main -- SPEC.md CLI.md internal cmd .codecrew/roles`
produced no commits." `main` in radiusred/codecrew-www moved by exactly
**seven commits**, in a straight line with no merge commits, from M20's merge
([e01e267](https://github.com/radiusred/codecrew-www/commit/e01e267),
committed 2026-09-26T11:00:48Z) to #73's merge
([53df5ff](https://github.com/radiusred/codecrew-www/commit/53df5ff)). Every
one of the seven has the author `radiusred-cody[bot]` and carries the trailer
`Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`. Between the two
commits three files changed: `docs/index.md` (101 lines added, 83 removed),
`docs/stylesheets/extra.css` (91 and 17) and `tests/test_site_build.py` (243
and 32).

**The trail is checked by the verb.** The installed extension reports `v2.1.1
(protocol 2.1)`; M22 cut no release. Run from this repository at the counting
instant:

```
$ gh codecrew milestone evidence 22
requirements counted: M22-R1, M22-R2, M22-R3 (3)
all 7 cited links resolve across 4 issues — evidence is reachable
```

Four issues: the milestone and its three sub-issues. Seven links, counted once
each: the live homepage, which both of QA's verdict comments and its finding
on #71 cite; that finding and the operator's gate resolution, both cited by
QA's first verdict comment; the commit `53df5ff` and the deploy run for it,
both cited by the re-verification; #69's fourth Decision, cited by #71's
first; and #73's first Decision, cited in the qualifier of #73's second.
Neither Deviation on #420 cites a URL.

`gh codecrew milestone close 22 --dry-run`, as it read at the counting
instant, before this pull request existed, passes five gates and stops at the
document, as M18's to M21's did:

```
gate milestone open: ok
gate no gate raised: ok
gate tasks closed: ok
gate requirements declared: ok
gate QA verdicts: ok
gate milestone document: refused[DOC_MISSING]: docs/milestones/22-*.md not on the default branch of radiusred/gh-codecrew — dispatch the doc-synthesizer: it delivers the record as a housekeeping PR with no task (SPEC §4), reviewed as the reviewer seat's routing requires (in pure solo, the operator confirms on the PR), and its author rebase-merges it; then rerun
dry run: nothing written — the live verb stops at the first refusal above
```

The dry run prints no raw material, so every record below was read from its
issue, pull request, review, commit or timeline directly. The prose is wrapped
at this file's normal width; the requirement-outcomes table's rows, the fenced
output above and a handful of single links are not, because a newline ends a
Markdown table row and breaking a link breaks it.

**This PR adds the M22 ROADMAP row and changes nothing else in the hub.**
`ROADMAP.md` carried no M22 row before it. The section on the front door below
says why neither front-door document changes.

## Goal and outcome

#420 was opened by the operator's account at 01:26:59Z on 2026-10-24, thirteen
days, one hour and eighteen minutes after [`milestone close
21`](https://github.com/radiusred/gh-codecrew/issues/393#issuecomment-6103609099)
closed M21. Its goal starts from a comparison: "An outside appraisal of
CodeCrew (ChatGPT, on v2.1.1) summed up what the framework does well more
clearly than codecrew.works does." It names what the appraisal's section 2
covers — "deterministic gates, GitHub as the source of truth, separate
identities and decisions recorded as they happen" — and says "a reader meeting
CodeCrew for the first time would understand the framework faster from it." It
sets the work: "This milestone brings that material onto the homepage and
removes the copy it makes redundant (radiusred/codecrew-www#68)." It carries
the operator's three rulings, then bounds the scope: "It changes
codecrew.works only. The protocol, the CLI and the synced docs are untouched,
so the bedding-in rule stands."

The capture it names,
[radiusred/codecrew-www#68](https://github.com/radiusred/codecrew-www/issues/68),
was filed by the operator's account at 01:18:00Z, eight minutes and fifty-nine
seconds before the milestone. It says: "The operator asked ChatGPT for an
appraisal of CodeCrew at v2.1.1. Its section 2, "Where I think the design is
genuinely strong", is a better synopsis of what CodeCrew does well than the
homepage currently gives." It quotes that section verbatim, maps each of its
five points to the homepage copy it overlaps in a table, says the "A second
model, with other blind spots" panel "has no counterpart in section 2 and
should probably stay", and raises three ask-the-human points: "Quote or
adopt?", "Accuracy at the current release" and "Keep the appraisal's balance
in view". The section 2 quoted there is all of the appraisal the trail holds.

The operator's answers are on #69, in a
[comment](https://github.com/radiusred/codecrew-www/issues/69#issuecomment-6378158446)
posted by the operator's account at 01:27:25Z, "in the operator's words": (1)
"rewrite as copy but with substantially similar content"; (2) "make it
accurate and drop the version/date"; (3) on the appraisal's reservations,
"Leave them both out. The point of adding this from the appraisal is not to
inflate or misrepresent the framework, it's to explain it better to a human
encountering it for the first time." #420's goal restates the three as rulings
of 2026-10-24.

The three requirements, in short (the table below carries each in full):

- **M22-R1** — the homepage explains, "in the site's own voice and with
  substantially the appraisal's content", five things: the gates in code and
  what each checks; GitHub as the only state store; one App per role "with its
  own permissions, rather than every agent sharing one privileged token, with
  the required-review caveat"; decisions recorded as they happen; and "the
  project is working software with tests, releases and CI, not a set of
  prompts". "No first person, no attribution to the appraisal, and none of its
  ratings or reservations".
- **M22-R2** — "every claim in the new copy is accurate at the current
  release, checked against SPEC.md and CLI.md", with the gate lists matching
  the refusal table, the `--with-approval-permission` qualifier, and "no
  version number or release date".
- **M22-R3** — "nothing on the homepage says the same thing twice", the
  second-model panel kept, the page reading cleanly at 375px and 1280px, in
  light and dark, with keyboard focus and the mobile drawer working, docs and
  blog pages unchanged, and "the site-build tests pin the new structure".

No requirement adopts a capture. One task does: #69 adopts #68.

The three tasks, in merge order:

- **[#69](https://github.com/radiusred/codecrew-www/issues/69) / PR
  [#70](https://github.com/radiusred/codecrew-www/pull/70)** — M22-R1, M22-R2
  and M22-R3; adopts #68. One commit, approved in its one round, merged
  01:42:36Z, two minutes and eighteen seconds after the pull request opened.
  Four Decisions. The review and the merge were not dispatched by this hub's
  coordination layer; the Deviations below carry the account.
- **[#71](https://github.com/radiusred/codecrew-www/issues/71) / PR
  [#72](https://github.com/radiusred/codecrew-www/pull/72)** — M22-R1 and
  M22-R3, the operator's revision of #70's page. Four commits, the first two
  carried from `task/69`; two review rounds; merged 02:08:18Z. Two Decisions.
- **[#73](https://github.com/radiusred/codecrew-www/issues/73) / PR
  [#74](https://github.com/radiusred/codecrew-www/pull/74)** — M22-R3, after
  QA's `not satisfied`. Two commits, three review rounds, merged 02:34:53Z.
  Two Decisions.

Every one of the three was rebase-merged by `radiusred-cody[bot]`, and every
merge was followed by a successful `Build and deploy site` run on the spoke's
`main` — three runs, the last,
[40054765515](https://github.com/radiusred/codecrew-www/actions/runs/40054765515),
finishing at 02:35:29Z. Each pull request's head passed `Lint commit messages`
and `Tests`.

**The work was serial.** No two of the three pull requests were open at the
same time, and each task was opened after the previous pull request merged:
#71 twelve minutes and thirteen seconds after #70, #73 eight minutes and four
seconds after #72. Each task's plan was written by the implementer seat
between five and seventeen seconds before that task's "Started by" comment.

**The milestone ran for one hour, thirteen minutes and twenty-five seconds**,
from the issue opening to QA's last verdict. The first version of the page
merged fifteen minutes and thirty-seven seconds after the milestone opened.
The coordination layer's Deviation followed twelve minutes and eleven seconds
after that merge, and #71 two seconds after the Deviation. QA's first verdicts
came seven minutes and thirty-eight seconds after #72 merged; #73 was opened
twenty-six seconds after them, and QA's re-verification five minutes and
thirty-one seconds after #74 merged.

What exists afterwards that did not before, read on the live page at 02:45:20Z
(generator `zensical-0.0.58`, title "CodeCrew — Your coding agents, working as
a crew on GitHub"). The hero's second paragraph is the commodity line alone:
"The coordinator running your agents is a commodity now; separation of duties
is not, and distinct identities are what make it real." A Why section follows
it directly, in five panels: "One App per role, with that role's permissions",
"A second model, with other blind spots", "Gates in code, not instructions",
"GitHub is the only state store" and "Decisions, recorded as they are made".
The crew comes next, laid out as an org chart with the coordinator above the
four other roles, then the worked example, whose `task finish` note now reads
"The gates run here: if one fails, it refuses with its reason." The proof
section's first receipt is "It is working software, not a set of prompts.",
and its agent-staffed receipt says that "`gh codecrew task finish` merges a
task here only once the reviewer App has approved". The required-review trade
— write access, GitHub's required review, branch protection and the paid plan
— is explained under Start now and nowhere else. The page names no version and
no release date, and neither ChatGPT nor an appraisal. The spoke's test suite
went from 94 passing tests on PR #70 to 98 on PR #74. The hub, its protocol,
its CLI and its synced docs are unchanged.

## Decisions

Eight Decision records, all on the tasks and all posted by the implementer
seat (`radiusred-cody[bot]`). One carries a parenthetical qualifier. Two
record calls that the comments attribute to the operator — #69's first, the
rulings, and #69's fourth, the revision after the operator's review of the
rendered page — and #71's two carry the second of those and the operator's
org-chart layout. The gate resolution is under the gates below.

### The rulings, and where the five ideas go

#69's [first
Decision](https://github.com/radiusred/codecrew-www/issues/69#issuecomment-6378240657)
opens "The coordination layer's rulings bind this task, and I do not re-open
them", and restates the three: the copy "is the site's own" and "does not name
ChatGPT or "an appraisal" anywhere on the page"; every claim accurate "and the
copy names no version and no date"; the ratings and reservations left out.
**Trade-off:** "Rewriting loses the outside voice that an attributed quote
would carry. It also avoids framing an AI model's opinion as an endorsement,
and it lets each claim be corrected against SPEC." **Rejected:** "Quoting
section 2 with attribution, its version and its date. The operator rejected
this in ruling 1."

The
[second](https://github.com/radiusred/codecrew-www/issues/69#issuecomment-6378240852)
places the ideas: four of them become Why panels "alongside the kept
second-model panel"; the identity panel "absorbs "See who built it, and who
checked it""; the gates and decisions panels "split "The record is the work"
between them"; working software "goes into the proof section ("CodeCrew
Works"), merged into the first receipt, which is the receipt that duplicates
it"; and the hero's clause about the reviewer App's approval "moves to the
gates panel". **Trade-off:** "Five panels do not fill a three-column grid
evenly", so the grid becomes six tracks; and working software "is evidence
rather than a reason to adopt, so it sits with the other evidence instead of
under "Why"." **Rejected:** "A sixth "working software" panel under "Why"",
and "A single prose passage replacing the panels".

### The gate lists, against the refusal table

#69's [third
Decision](https://github.com/radiusred/codecrew-www/issues/69#issuecomment-6378283774)
checks "each gate claim against the hub's SPEC §10 refusal table" and names
"only what the table backs": for `task finish`, `NO_PR`, `NOT_OWNER`,
`NO_CHECKS`, `CHECKS_FAILING`, `NO_HOLDER_REVIEW` and `GATE_UNRECORDED`; for
`milestone close`, `OPEN_TASKS`, `VERDICT_UNSATISFIED` and `DOC_MISSING`. It
corrects the source: "The source said `task finish` checks the plan. It does
not: `task start` is the verb that refuses `NO_PLAN`". "A new test reads the
SPEC table and pins each claim to its code." **Trade-off:** the identity panel
carries the required-review caveat "only as the qualifier", linking to Start
now for the rest; "The alternative was a third full copy of the caveat, next
to Start now's note and the agent-staffed receipt, and the brief leaves that
receipt alone." **Rejected:** "Stating the whole caveat in the panel, because
that says the same thing again and M22-R3 rules out repetition."

### The operator's review: Why first, the crew before the example, the copy cut

#69's [fourth
Decision](https://github.com/radiusred/codecrew-www/issues/69#issuecomment-6378348092)
was posted at 01:45:51Z, three minutes and fifteen seconds after PR #70
merged, and edited once, at 01:47:32Z; the edit recast it from a list into
three numbered changes and added the continuity fixes, and kept its trade-off
and rejection in substance. As it stands: "Following the operator's review of
the rendered page on PR #70, I made three changes. All are the operator's
call." Why "now comes straight after the hero, so a first-time reader meets
the ideas first"; "The page order is now: hero → Why → the crew → the worked
example → CodeCrew Works → Start now. Readers know the roles before the
example tells its story through them"; and "Each Why panel is one idea and its
reason, in about 36–47 words. The section is about 210 words, down from about
500." Its continuity fixes gloss the roles where Why first names them, credit
the reviewer's reason to "the reviewer's instructions", "not to a contract the
page has not introduced yet", turn "Human gate" into "every question raised
for a human answered", and point the example's `task finish` note back at the
gates. **Trade-off:** "The gates panel no longer lists every check. It names
green CI, the reviewer's approval, answered human questions and QA
verification of every requirement"; "For length, I dropped "No server, no
dashboard" and the "a bot name proves the App, not the model" sentence."
**Rejected:** "The panel order from the first commit, with the gates first.
The hero ends on "distinct identities are what make it real", so the
identities panel picks that up, and the gates panel would have opened cold."

Two of the phrases that trade-off names did not ship. The reviewer's first
round on PR #72 found both inaccurate, and the panel now reads "the approval
your routing requires" and "every requirement's latest QA verdict is
`satisfied`"; the implementer's [answer to that
round](https://github.com/radiusred/codecrew-www/pull/72#issuecomment-6378514714)
records the change, and no Decision on #69 or #71 does.

### The carried revision and the org chart

#71's [first
Decision](https://github.com/radiusred/codecrew-www/issues/71#issuecomment-6378451829)
carries the fourth Decision's two commits onto a fresh branch: "PR #70 merged
early, before the revision landed, so it never reached main. The two commits
are 8228b63 and 7042d67 from task/69, cherry-picked onto main at e4107cc, with
each reference reworded to #71. The trees are identical to task/69's head."
**Trade-off:** "The reasoning stays recorded once, on #69, and this comment
links to it." **Rejected:** "Squashing the two commits into one. Each one
matches one of the operator's two review rounds."

The
[second](https://github.com/radiusred/codecrew-www/issues/71#issuecomment-6378466207)
makes the crew "an org chart, as the operator asked": the coordinator "comes
first in the markup and is centred on top", a CSS grid, and from 45em
connector lines that "are pseudo-elements, so they add no markup and nothing
to the accessibility tree"; on a phone, "the four sit beneath it in a 2×2
grid, with no lines". The coordinator's summary is unchanged — "It already
reads as the head of the chart" — and the routing example's YAML "keeps its
order, with the coordinator row last". **Trade-off:** "A 2×2 grid on a phone
is half the height of a stack, and keeps the four reading as peers under one
head." **Rejected:** "A stack of four on a phone", "Connector elements in the
markup with `aria-hidden`", and "Reordering the YAML so the coordinator comes
first. That table shows the routing format, not the chart".

### The receipt: from the repository's history to the merge gate

#73's [first
Decision](https://github.com/radiusred/codecrew-www/issues/73#issuecomment-6378663046)
has the agent-staffed receipt state "this project's own fact instead of the
general rule: "each merge here waits on the reviewer App's approval, not a
person's"", and calls it "true at every tier of this project's routing".
**Trade-off:** "The receipt no longer says how the project's merges also
satisfy GitHub's own required review. That explanation is Start now's alone,
as M22-R3 requires." **Rejected:** a wording that "still carries half of the
rule's reasoning, which is the repeat QA flagged", and "cutting the sentence".

The
[second](https://github.com/radiusred/codecrew-www/issues/73#issuecomment-6378796342),
labelled "**Decision (superseding the receipt wording in …):**", replaces that
claim: "The earlier claim, "each merge here waits on the reviewer App's
approval, not a person's", is false for the history, because bootstrap PR #1
merged without review. radiusred-checky[bot] found this in round 1 of the
review of PR #74". The receipt "now describes the merge gate, not the
repository's history". **Trade-off:** the claim rests on the reviewer seat's
routing and `NO_HOLDER_REVIEW`; "`--bypass` does not lift this"; "It no longer
covers merges made outside `task finish`, which is the honest bound."
**Rejected:** "A claim about every merge since a given date or PR. It would
need a date or a count that goes stale", and "A claim about how the seat is
routed alone". The first Decision is left as written.

## Deviations

Two Deviation records, both on the milestone issue. No task issue or pull
request carries one, and no pull request body mentions one.

**[PR #70 was reviewed and merged before the operator's review of the rendered
page had
finished](https://github.com/radiusred/gh-codecrew/issues/420#issuecomment-6378437386).**
The coordination layer's, posted by the operator's account at 01:54:47Z,
twelve minutes and eleven seconds after the merge. "As a result, M22 ships its
homepage change in two tasks rather than one." It gives the times — approved
"by `radiusred-checky[bot]` at 01:41:51Z and merged through `task finish` by
`radiusred-cody[bot]` at 01:42:36Z" — and says: "This hub's coordination layer
dispatched neither step. Another Paperclip company's seat routines did
(RRT-1372, RRT-1373), triggered by the App-level webhooks of the shared checky
and cody Apps. So the review did not run under the hub's routing (Codex,
`gpt-5.6-terra`)." The merged version "carries M22-R1 and M22-R2's content,
but not the operator's changes made after viewing the rendered page", which
"are on the task branch at `7042d67`". **Why:** "The cody and checky Apps are
installed org-wide and shared with another coordination layer whose routines
do not filter by repository. Capture #421 records the gap. The operator has
paused those routines for the rest of M22. A second task under M22 carries the
revision, and the reviewer will be dispatched under the routing table." It
gives the pull request's opening as 01:40:20Z, the time GitHub's timeline
connects it to the task; the pull request API's creation time is 01:40:18Z.

**[The record takes harness facts from its dispatch
brief](https://github.com/radiusred/gh-codecrew/issues/420#issuecomment-6378919539).**
This record's own, at 02:46:23Z. The brief states that the implementer ran
under Claude Opus 5.5 in Claude Code, "as one session across all three tasks,
resumed between them"; that every review round ran in Codex under
`gpt-5.6-terra` and both QA rounds under `gpt-5.6-sol`, "except PR
radiusred/codecrew-www#70's review and merge"; that the coordination layer was
the operator's own Claude Code session; and that the operator iterated through
a local dev server mirroring the implementer's working tree. The comment sets
out what the trail supports — the commit trailers, the "Generated with Claude
Code" line that ends every task pull request body, the routing table, the
coordination layer's Deviation and #421, #71's Plan and the gate resolution —
and that "Nothing on the trail says whether the three tasks shared one
implementer session. No review or QA comment names its own harness or model."
It also puts on the record the change request PR #72 received under the
operator's account, which the brief does not mention. **Why:** "the brief is
not on the trail, so this comment puts its account there, and the record cites
this comment rather than the brief, as the M18 to M21 records did on #353,
#382, #384 and #393."

## The gates

**#420's Gates section keeps the scaffold's placeholder line and adds two of
its own**, in the edit fourteen seconds after the issue was created:

- "QA checks the deployed page on codecrew.works, not only the build (the M19
  lesson: the live site found what CI and review missed)."
- "The operator checks the rendered page before the record is written."

**The first was discharged on the live site, twice.** QA's first comment reads
the deployed homepage "top to bottom in Chromium" and ends: "Gates: QA's
“checks the deployed page” gate is discharged by the live probes above." The
re-verification cites the [deploy
run](https://github.com/radiusred/codecrew-www/actions/runs/40054765515) for
the merged head and reports that "the live homepage matched the local `main`
build byte-for-byte" after "normalizing only Cloudflare's email obfuscation".

**The second was resolved on the page as it stood at PR #72's third commit.**
The operator's [**Gate
resolved:**](https://github.com/radiusred/gh-codecrew/issues/420#issuecomment-6378478542)
comment, at 01:59:03Z, fifty-two seconds after PR #72 opened, reads: "The
operator checked the rendered page before the record was written. They viewed
it on a local dev server built from PR radiusred/codecrew-www#72's head
(`e170874`), and approved it on 2026-10-24 in these words: "It's good - the
story builds in a more natural way now and explains more detailed concepts
after the basic and important ones. Let's ship that"." It names the three
revisions and says they "came from the same review of the rendered page, first
of PR #70's version and then of the branch as it developed." Three commits
changed the page after `e170874`: `d6b7d03`, the gates panel's wording, three
minutes and fifteen seconds after the gate resolution; and PR #74's two, the
identity panel and the agent-staffed receipt. No later comment of the
operator's is on the trail. QA's first comment cites the gate resolution as
the gate's answer: "The operator's rendered-page gate is resolved in the
milestone's Gate resolved record."

**The operator's review was in progress when PR #70 merged.** The coordination
layer's Deviation says so, and #421 says the merge "pre-empted an operator
review of the rendered page that was in progress". The approval on PR #70
ends: "Rendered interaction checks remain the author's reported evidence;
deployed QA and operator review remain the milestone's gates."

**No gate was raised.** No task used `gh codecrew checkpoint`, no issue
carries `cc:needs-decision`, and the dry run reports `gate no gate raised:
ok`. Each task's Ask-the-human section records none and says where the answer
came from: #69's "The operator answered the three questions from www#68 in a
comment on this task", #71's "The operator set the layout in the coordinator's
dispatch. The phone layout is an ordinary call, recorded as a Decision", and
#73's "The coordinator's dispatch sets out the remedy."

What gated the work: CI on every pull request — `Lint commit messages` and
`Tests`, both green on the head each of the three merged at — the reviewer
seat holder's approval on each of the three, the operator's gate, and the QA
verdicts.

## The review rounds

Three pull requests, **seven review submissions**: six by the reviewer role
holder (`radiusred-checky[bot]`), an App distinct from the authoring seat —
**three approvals, one per pull request, all standing**, and **three change
requests**, two resolved by a later commit and one by a superseding Decision
and a pull request body edit with no commit — and one by the operator's
account, a change request dismissed thirty-three seconds after it was
submitted. No pull request carries an inline review comment, and no pull
request's timeline carries a force-push event. PR #70 was approved in its
first round; PR #72 took two; PR #74 took three.

- **PR [#70](https://github.com/radiusred/codecrew-www/pull/70)** — one round,
  outside the hub's routing (the coordination layer's Deviation). The
  [approval](https://github.com/radiusred/codecrew-www/pull/70#pullrequestreview-5636222432),
  at 01:41:51Z, one minute and thirty-three seconds after the pull request
  opened, checked "the gate lists, default App permissions and approval opt-in
  against upstream SPEC.md and docs/identities.md" at `0d16fff` and reports
  "No blocking findings" and 94 tests passing. Its validation command names
  the working path `/tmp/checky-sources-1372`.
- **PR [#72](https://github.com/radiusred/codecrew-www/pull/72)** — two rounds
  by the reviewer seat. At 02:00:51Z a change request was submitted under the
  operator's account; at 02:01:24Z `radiusred-checky[bot]` dismissed it with
  the message "Superseded by the reviewer App review.", and two seconds later
  submitted the [same
  text](https://github.com/radiusred/codecrew-www/pull/72#pullrequestreview-5636321437)
  as its own change request: "M22-R2 needs two accuracy corrections in the
  shortened **Gates in code, not instructions** panel." The panel's "the
  reviewer's approval" "is not true for an operator-held reviewer role", and
  "QA can verify an unsatisfied requirement" so "has verified every
  requirement" overstates `VERDICT_UNSATISFIED`. The
  [approval](https://github.com/radiusred/codecrew-www/pull/72#pullrequestreview-5636343690)
  at `d6b7d03` found the new wording matches "`NO_HOLDER_REVIEW`,
  `NO_NONDOER_APPROVAL`, `VERDICT_MISSING`, and `VERDICT_UNSATISFIED`", and
  that "restoring either old phrase made the gates-vs-SPEC test fail".
- **PR [#74](https://github.com/radiusred/codecrew-www/pull/74)** — three
  rounds. The [first change
  request](https://github.com/radiusred/codecrew-www/pull/74#pullrequestreview-5636413619):
  "the replacement receipt overstates the project's history." The reviewer
  "audited all 29 merged PRs, requiring an `APPROVED` review by
  `radiusred-checky` before the merge", and found that PR #1 "was merged by
  `davison` with an empty reviews list". The
  [second](https://github.com/radiusred/codecrew-www/pull/74#pullrequestreview-5636432709):
  "the source and its test are correct, but the task record is not final" —
  the pull request description and "its linked receipt Decision still say
  "each merge here waits on the reviewer App's approval"" — and it asked to
  "post a superseding structured Decision on #73". The implementer
  [answered](https://github.com/radiusred/codecrew-www/pull/74#issuecomment-6378798199)
  "with no new commit"; the pull request's edit history shows its body changed
  at 02:32:53Z. The
  [approval](https://github.com/radiusred/codecrew-www/pull/74#pullrequestreview-5636438618):
  "the PR description quotes the current merge-gate receipt throughout and
  contains no "each merge" wording."

**What the three change requests were about.** They carry four findings. Three
are claims on the page broader than what CodeCrew or this project does: an
approval stated as universal, QA "verified" where the gate is a satisfied
latest verdict, and "each merge" where one merge had no review. The fourth is
the record: a pull request body and a Decision that still quoted the claim the
code had dropped.

## QA: two rounds, one remedy

The qa role holder (`radiusred-testy[bot]`) wrote two verdict comments on #420
and one finding on #71.

**The
[first](https://github.com/radiusred/gh-codecrew/issues/420#issuecomment-6378639442)**,
at 02:15:56Z, seven minutes and thirty-eight seconds after PR #72 merged,
holds one verdict per requirement. Each opens with the merged-main suite —
"The merged-main suite passed 97 tests" — and then reads the deployed page:

- **M22-R1** — `satisfied`: "all five source ideas in the site's own voice
  with their reasons", with the failed-session claim's "loss bounded to
  unpushed work"; the repository, Actions, releases and milestone-document
  links "each returned HTTP 200"; and "The live copy has no first person,
  ChatGPT/appraisal attribution, rating or reservation."
- **M22-R2** — `satisfied`: every live claim matched "against hub `SPEC.md`
  §§5/10 and `CLI.md`"; the identity panel "bounds read-only reviewer/QA
  permissions with “by default” and names `--with-approval-permission`"; the
  quotation "character-exact"; "no version or release date".
- **M22-R3** — `not satisfied`: "the required-review caveat appears three
  times: the Why identity panel says the opt-in gives the reviewer write
  access so approval can count; the agent-staffed receipt says a write-enabled
  reviewer App satisfies required review; and Start now repeats the opt-in,
  write-access and required-review explanation." The [finding on
  #71](https://github.com/radiusred/codecrew-www/issues/71#issuecomment-6378633702),
  thirty-eight seconds earlier, carries the three quotations and names the
  suite's gap: "`test_each_of_the_five_ideas_is_on_the_page_once` counts
  selected literal phrases and does not reject the repeated semantic claim."
  The rest of R3 "passed live": the order, both widths in both schemes, "Tab
  traversed 37 mobile and 41 desktop visible focus stops without a trap", the
  drawer, and docs and blog pages matching "a pre-M22 build".

**The
[second](https://github.com/radiusred/gh-codecrew/issues/420#issuecomment-6378865032)**,
at 02:40:24Z, five minutes and thirty-one seconds after PR #74 merged, is
`satisfied` for M22-R3 and says "This supersedes my earlier M22-R3 verdict."
It reports 98 tests on merged `main` at `53df5ff`; "the required-review rule,
its write-access reason, branch protection, and the paid-private-plan
qualifier occur only in Start now"; the receipt "accurate against the spoke's
hub pointer, the hub routing of reviewer to `app:radiusred-checky`, and SPEC
§10's `NO_HOLDER_REVIEW`"; "38 mobile and 42 desktop focus stops"; and a
scratch-copy mutation in which the new test "failed exactly on `cc-why` with
`approval can count` and `required review`, proving the new regression test
rejects the defect rather than merely passing the shipped page."

**Every standing verdict is `satisfied`.** M22-R1 and M22-R2 stand in the
first comment; M22-R3 stands in the second. No requirement was struck.

## Requirement outcomes

The status column is the **standing** verdict word as the qa role holder wrote
it, verbatim and unqualified: M22-R1 and M22-R2 in comment
[6378639442](https://github.com/radiusred/gh-codecrew/issues/420#issuecomment-6378639442),
and M22-R3 in comment
[6378865032](https://github.com/radiusred/gh-codecrew/issues/420#issuecomment-6378865032),
which supersedes the first comment's `not satisfied` for that requirement
alone.

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M22-R1 — the homepage explains, in the site's own voice and with substantially the appraisal's content: (a) task finish and milestone close are gates in code, not prompt instructions, with what each one checks; (b) GitHub is the only state store, so there is no workflow database to drift from the repository and a failed session loses no task context; (c) each role can act through its own GitHub App with its own permissions, rather than every agent sharing one privileged token, with the required-review caveat; (d) decisions and deviations are recorded as structured comments when they happen and synthesized into the milestone's document; (e) the project is working software with tests, releases and CI, not a set of prompts. No first person, no attribution to the appraisal, and none of its ratings or reservations | [#69](https://github.com/radiusred/codecrew-www/issues/69) / PR [#70](https://github.com/radiusred/codecrew-www/pull/70), then [#71](https://github.com/radiusred/codecrew-www/issues/71) / PR [#72](https://github.com/radiusred/codecrew-www/pull/72) | `satisfied` | Both tasks closed; #68 closed by #69's finish. #70 approved in its one round, outside the hub's routing; #72 two rounds. Six Decisions (the rulings, the placement, the gate lists, the operator's revision on #69; the carried revision and the org chart on #71). Verdicted after #72, before #74 |
| M22-R2 — every claim in the new copy is accurate at the current release, checked against SPEC.md and CLI.md: the gate lists match the refusal table; the claim that review and QA need not be able to change code carries the --with-approval-permission qualifier; no version number or release date appears in the new copy | [#69](https://github.com/radiusred/codecrew-www/issues/69) / PR [#70](https://github.com/radiusred/codecrew-www/pull/70) | `satisfied` | Task closed. The gate lists checked against SPEC §10 by #69's third Decision; #72's round one then corrected two of the shortened panel's phrases, recorded in a PR comment. Verdicted after #72, before #74 |
| M22-R3 — nothing on the homepage says the same thing twice: copy that the new material duplicates (the Why panels, the hero's second paragraph, the receipts) is removed or merged into it, and the second-model panel is kept unless it is folded in whole; the page reads cleanly at 375px and 1280px, in light and dark, with keyboard focus and the mobile drawer working; docs and blog pages are unchanged; the site-build tests pin the new structure | [#69](https://github.com/radiusred/codecrew-www/issues/69) / PR [#70](https://github.com/radiusred/codecrew-www/pull/70) and [#71](https://github.com/radiusred/codecrew-www/issues/71) / PR [#72](https://github.com/radiusred/codecrew-www/pull/72), then [#73](https://github.com/radiusred/codecrew-www/issues/73) / PR [#74](https://github.com/radiusred/codecrew-www/pull/74) | `satisfied` | All three tasks closed. `not satisfied` (the required-review caveat three times on the live page), then `satisfied` after #74. #74 took three review rounds. Two Decisions on #73, the second superseding the first |

**Captures adopted and closed by the merges — one.**
[radiusred/codecrew-www#68](https://github.com/radiusred/codecrew-www/issues/68),
filed by the operator's account at 01:18:00Z, the milestone's source. It
carries an "Adopted by" comment naming #69 at 01:27:17Z, from the operator's
account, and a "Closed by `gh codecrew task finish 69`" comment from
`radiusred-cody[bot]` at 01:42:41Z naming PR #70 and the merge commit — the
finish the coordination layer's Deviation says another company's routine ran.

**Captures filed during the milestone and left open — one.**
[#421](https://github.com/radiusred/gh-codecrew/issues/421), in this
repository, by the operator's account at 01:54:32Z, fifteen seconds before the
Deviation that cites it: "App-level webhooks on shared crew Apps let another
hub's routines review and merge this hub's PRs". It says "The routines'
instructions say "on the project repo", but nothing compares the event's
repository with the project", and that "Nothing on the trail says which
coordination layer acted. The review and the merge read exactly as if the hub
had dispatched them." It weighs three shapes — docs, a record naming the
dispatcher ("This connects to #378"), and a CLI check it calls "probably out
of reach" — and relates itself to #177 and #347.

From #420's creation to the counting instant, the only issues or pull requests
created in this repository and the spoke are #420, #421, the three tasks and
their three pull requests; an organisation-wide issue search over the same
window, as this record's token sees the organisation, returns the same eight.
#68 names [#416](https://github.com/radiusred/gh-codecrew/issues/416) to
[#419](https://github.com/radiusred/gh-codecrew/issues/419), filed by the
operator's account at 01:17:19Z–01:17:21Z, as "the project's answer to those
reservations"; all four are open, unlabelled, and outside M22.

## The front door

The doc-synthesizer contract sends this hub's front-door list to its
`.local.md`: "the README's proof points, and the introduction's release, verbs
and refusal codes". M22 changed none of those — no release, no verb, no
refusal code, and no hub commit. What M22 changed is codecrew.works, so the
check here, as in [M19's
record](19-the-front-page-tells-the-crew-story.md#the-front-door), is whether
either front-door document makes a claim about the site that M22 made false.
Read against the live page at 02:45:20Z:

- **`README.md`** calls codecrew.works "the marketing and introduction site"
  that "carries the blog". True as it stands; the blog answers 200.
- **`docs/introduction.md`** makes three claims about the site. It links "an
  annotated generic" routing table at `#the-crew`: the anchor is on the live
  page and the illustrative `roles:` table is under it, in the order #71's
  second Decision kept. It says the receipts are on the home page at
  `#codecrew-works`: the anchor is live and the receipts are under it. And it
  lists the site as "the introduction site: what the problem is, how the
  protocol works, and what it has delivered": the page has the Why section,
  the crew and the worked example, and the proof section. The spoke's
  `sync_docs.py` maps the README's two homepage anchors to `#the-crew` and
  `#start-now`, and both are on the live page. None of the claims is false, so
  this PR changes neither document.

SPEC §4 lists the record's diff as the document, its ROADMAP row and the
front-door claims the milestone changed, "and nothing else", and says the
housekeeping path carries "no record entry", so this PR adds no CHANGELOG
entry — as the M19 to M21 record pull requests added none.

## Protocol-discipline observations

Six things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **A seat's identity did not say who was driving it.** The approval and the
  merge on PR #70 came from the reviewer seat's App and the seat that started
  the task, and #421 calls the merge "protocol-compliant". Neither the review
  nor the finish carries anything that distinguishes a dispatch by another
  company's routine from one by this hub; the trail's account is the
  coordination layer's Deviation and #421, both written afterwards. The
  approval itself named the operator's review as an open milestone gate, and
  `task finish` reads the task issue, not the milestone's Gates section.
- **A merged pull request's summary describes work it did not merge.** SPEC §4
  makes the merged PR the task's summary. PR #70's body was edited at
  01:48:00Z, five minutes and twenty-four seconds after the merge, to describe
  the section order, the cut panels and the example's note — the content of
  the two commits pushed to `task/69` after the merge, which reached `main`
  through PR #72. The edit links the fourth Decision; no comment records the
  edit, and the merged commit `e4107cc` holds the earlier page.
- **The operator's gate was answered on an earlier head than the one that
  shipped.** The gate resolution names `e170874`; the reviewer's round on PR
  #72 and QA's R3 remedy then changed three things on the page. Each change
  was reviewed and QA verified the result, and no later operator comment is on
  the trail. [M19's
  record](19-the-front-page-tells-the-crew-story.md#the-gates) reports that
  milestone's gate resolution approving the page "as it stands" before its
  last fix landed, on a condition the trail then met.
- **A review appeared under the wrong identity and was corrected on the
  trail.** PR #72's first-round change request was submitted under the
  operator's account, dismissed by the reviewer App and resubmitted by it word
  for word. The dismissal message records the correction; no comment says how
  the first submission came to be made, and the dismissed review remains in
  the pull request's history.
- **Two review rounds ended differently on the record.** On PR #72 the
  reviewer's two findings changed wording a Decision had recorded, and the
  change is recorded in a pull request comment only, which SPEC §4 does not
  gather. On PR #74 the reviewer's second round asked for a superseding
  Decision because the task's record still quoted the dropped claim, and one
  was posted.
- **QA's live read found a repetition the task's own records had weighed.**
  #69's third Decision named a third copy of the required-review caveat as the
  alternative it avoided, and PR #70's body says the agent-staffed receipt
  "still carries its required-review sentence, because the brief allows
  trimming receipts only where they duplicate point (e)". Two approvals
  followed. QA counted three statements of the caveat on the live page; #73's
  dispatch-set remedy moved it to Start now alone and added a test that
  rejects the semantic repeat.

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **The appraisal is on the trail as one section.** #68 quotes section 2
  verbatim and quotes two of the appraisal's ratings in its third
  ask-the-human point; the rest of the document, and the conversation in which
  the operator commissioned it, are not on GitHub. *The reasoning is recorded;
  the deliberation is not*, as [M21's
  record](21-v2-1-1-what-the-fleets-first-fortnight-on-2-1-taught.md#what-the-record-does-not-contain)
  says of its own opening.
- **The operator's review of the rendered page is recorded through the
  implementer.** The operator's own words on the trail are the three rulings
  and the gate resolution's approval; the revisions themselves — what the
  operator asked for, in what order, on which server — are known from #69's
  fourth Decision, #71's goal and Decisions, and the gate resolution's
  summary.
- **The other company's routines are not on the trail.** RRT-1372 and RRT-1373
  are Paperclip tasks, named by the Deviation and #421; the routines'
  instructions are quoted only in #421, which session or model ran PR #70's
  review is not recorded, and when the operator paused the routines is not
  stated beyond "for the rest of M22".
- **The coordination layer's dispatch briefs are not on the trail.** #71's
  "The operator set the layout in the coordinator's dispatch", #73's "The
  coordinator's dispatch sets out the remedy" and PR #70's "the brief allows
  trimming receipts only where they duplicate point (e)" cite them.
- **Which harness and model ran each session is not recorded by the seats.**
  This record's Deviation carries the coordination layer's account and the
  trail's partial support for it; no review or QA comment names its model.
- **Nothing on the trail explains PR #72's review under the operator's
  account**, beyond the dismissal message.
- **`task/69` was pushed again after its pull request deleted it**, and it is
  in the spoke at `7042d67` at the counting instant; no comment mentions the
  re-push or the branch.
- **The screenshots behind the pull requests' rendered checks were not
  committed or posted.** Each of the three pull request bodies says so; the
  observations are quoted in the bodies and in QA's verdicts.
