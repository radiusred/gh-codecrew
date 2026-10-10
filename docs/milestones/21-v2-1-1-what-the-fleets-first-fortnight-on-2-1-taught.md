# M21: v2.1.1: what the fleet's first fortnight on 2.1 taught

Tracking issue: [#393](https://github.com/radiusred/gh-codecrew/issues/393) ·
Synthesized 2026-10-10 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #393's **eight** requirements as opened, **none
added, none struck and none amended** — the issue's `userContentEdits`
returns a total count of zero; its Gates section, left as the scaffold's
placeholder; its **four comments** — the operator's release leg, the qa role
holder's **two verdict comments**, and this record's **one Deviation**; the
**eight** task issues, all in this repository
([#394](https://github.com/radiusred/gh-codecrew/issues/394),
[#395](https://github.com/radiusred/gh-codecrew/issues/395),
[#396](https://github.com/radiusred/gh-codecrew/issues/396),
[#397](https://github.com/radiusred/gh-codecrew/issues/397),
[#398](https://github.com/radiusred/gh-codecrew/issues/398),
[#399](https://github.com/radiusred/gh-codecrew/issues/399),
[#400](https://github.com/radiusred/gh-codecrew/issues/400),
[#401](https://github.com/radiusred/gh-codecrew/issues/401)); their **eight
merged pull requests** and the **twenty-three commits** on them; the one
housekeeping pull request #401's plan calls its Part C,
[#412](https://github.com/radiusred/gh-codecrew/pull/412), and its **one**
commit; the **nineteen Decision records and four Deviation records** across
the milestone issue and its eight tasks — all nineteen Decisions and three of
the Deviations on the tasks, and one Deviation, this record's, on the
milestone issue — with **no** `**Gate resolved:**` record anywhere; the
**twelve review submissions** on the nine pull requests — **nine approvals,
every one of which stands**, and **three change requests**; the **seven**
adopted backlog captures the tasks' finishes closed; and the **two** captures
filed in the milestone's window, both open. The house form is
[M18's document](18-the-protocol-scales-down.md), whose milestone also
shipped a release, with
[M20's](20-the-2-1-post.md) behind it; the record standard is the one the
reviewer set on [PR #322](https://github.com/radiusred/gh-codecrew/pull/322)
and re-applied since, so every count is re-derived from the source — records
by the rule `tracker.ExtractRecords` applies (a paragraph-initial
`**Decision:**`, `**Deviation:**` or `**Gate resolved:**` label, bare or
parenthetically qualified), review states from each pull request's reviews
API and timeline, commits from each pull request's commits API, and the
requirement history from the milestone issue's `userContentEdits`.

**Counted at 2026-10-10T23:44:00Z, against a trail frozen from this
document's dispatch.** The last post on the milestone's issues before that
instant is this record's own
[Deviation](https://github.com/radiusred/gh-codecrew/issues/393#issuecomment-6103422065)
at 23:42:57Z, and it is inside every total above; the one before it is the
qa role holder's second verdict comment at 23:35:21Z.

**This record has no task.** It is delivered as a housekeeping pull request
under SPEC §4's milestone-record clause, as
[M18's](18-the-protocol-scales-down.md),
[M19's](19-the-front-page-tells-the-crew-story.md) and
[M20's](20-the-2-1-post.md) were: the milestone issue is its charter and its
commit references #393.

**`main` moved by exactly twenty-four commits**, in a straight line with no
merge commits, from M20's record
([164fd42](https://github.com/radiusred/gh-codecrew/commit/164fd42),
committed 2026-09-26T11:21:45Z) to the housekeeping merge
([2aa1bb4](https://github.com/radiusred/gh-codecrew/commit/2aa1bb4)):
twenty-three that reference one of the eight tasks and one, `chore:
regenerate the contract history for v2.1.1`, that references none. Every one
of the twenty-four has the author `radiusred-cody[bot]` and carries the
trailer `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`. Between
the two commits, 34 files changed, with 2,998 lines added and 216 removed;
`SPEC.md` is not among them.

**The trail is checked by the verb.** The installed extension reports
`v2.1.1 (protocol 2.1)`, the release this milestone shipped. Run from this
repository after the counting instant, with nothing posted in between:

```
$ gh codecrew milestone evidence 21
requirements counted: M21-R1, M21-R2, M21-R3, M21-R4, M21-R5, M21-R6, M21-R7, M21-R8 (8)
all 16 cited links resolve across 9 issues — evidence is reachable
```

Nine issues: the milestone and its eight sub-issues. Sixteen links, all on
github.com: the release run, cited from the operator's leg on #393 and from
#401's release evidence and counted once; the release page; #359's Decision,
cited by #401's first Decision; twelve from the qa role holder's second
comment — its own first-round comment, #393, the eight task pull requests,
the start record on #395 and the Deviation on #395; and #395 itself, cited by
that Deviation. Two of the sixteen, the Actions run and the release tag, are
shapes M21-R3 moved onto the API. QA's first round reports the same verb
resolving "all 3 across 9"; the R6 remedy and the verdict that followed it
added the other thirteen. This record's Deviation cites no URL.

`gh codecrew milestone close 21 --dry-run`, as it read at 23:40:21Z, before
this pull request existed, passes five gates and stops at the document, as
M18's, M19's and M20's did:

```
gate milestone open: ok
gate no gate raised: ok
gate tasks closed: ok
gate requirements declared: ok
gate QA verdicts: ok
gate milestone document: refused[DOC_MISSING]: docs/milestones/21-*.md not on the default branch of radiusred/gh-codecrew — dispatch the doc-synthesizer: it delivers the record as a housekeeping PR with no task (SPEC §4), reviewed as the reviewer seat's routing requires (in pure solo, the operator confirms on the PR), and its author rebase-merges it; then rerun
dry run: nothing written — the live verb stops at the first refusal above
```

The dry run prints no raw material, so every record below was read from its
issue, pull request, review, commit or timeline directly. The prose is
wrapped at this file's normal width; the requirement-outcomes table's rows,
the fenced output above and a handful of single links are not, because a
newline ends a Markdown table row and breaking a link breaks it.

**This PR adds the M21 ROADMAP row and changes nothing else in the hub.**
`ROADMAP.md` carried no M21 row before it. #411 had already moved both
front-door documents to v2.1.1; the section on the front door below says why
neither changes again.

## Goal and outcome

#393 was opened by the operator's account at 20:03:58Z on 2026-10-10,
fourteen days, eight hours, forty-one minutes and forty-eight seconds after
[`milestone close 20`](https://github.com/radiusred/gh-codecrew/issues/384#issuecomment-5845849474)
closed M20. Its goal opens on the release it asks for: "The first patch after
protocol 2.1. Running v2.1.0 across md-notes, a new hub with four spokes, the
tradedesk repositories and a private hub surfaced seven places where the CLI
falls short of what SPEC and GitHub already say". It lists them: "a merge the
CLI refuses that GitHub would allow, evidence links a private repository
cannot resolve, an init commit that breaks SPEC §4, a NO_CHECKS refusal
nobody is warned about, a reviewer contract that contradicts CLI.md, a task
left open after a successful finish, and runners that move image on
2026-10-19." It sets the bound — "All CLI patch fixes and contract
clarifications; no protocol change — the 2.1 bedding-in rule stands
(operator, 2026-10-10: CLI patch releases do not affect it)." — and names
what waits: "Deferred to protocol 2.2 as captures: #387 #365 #381 (the light
path's author fence and an operator-authored small change), the self-close
half of #386, #377–#379, and #391's scaffolded workflow."

**The seven places are seven backlog captures, all filed by the operator's
account**, one by the coordination layer from this repository's own release
runs and six from projects running v2.1.0 outside it:

- [#380](https://github.com/radiusred/gh-codecrew/issues/380), created
  2026-09-24T23:01:22Z, which its body attributes to the "Coordination layer,
  2026-09-25, from the v2.0.2 and v2.1.0 release runs": both runs carried the
  annotation that `ubuntu-latest` "will migrate to Ubuntu 26 beginning
  October 19, 2026".
- [#386](https://github.com/radiusred/gh-codecrew/issues/386), 2026-09-30,
  "from running `davison/md-notes` M12 on protocol 2.1 / CLI v2.1.0": on
  davison/md-notes#231 `task finish` printed "closes via its closing keyword"
  and "#231 stayed open".
- [#388](https://github.com/radiusred/gh-codecrew/issues/388), 2026-10-07,
  "in a private hub": Actions job URLs reported unreachable although `gh api`
  resolves the run; a comment of 2026-10-09 adds a second shape, a bare
  repository link on radiusred/trading-hub.
- [#389](https://github.com/radiusred/gh-codecrew/issues/389), 2026-10-09:
  `init`'s scaffold commit has no body, "Observed 2026-10-09 onboarding four
  spokes to a new hub: the reviewer seat … requested changes on all four PRs
  for exactly this".
- [#390](https://github.com/radiusred/gh-codecrew/issues/390), 2026-10-09:
  "Two reviewer sessions today … requested changes on otherwise-clean task
  PRs because the task issue had **no assignee**"; "one cost ~60k tokens, the
  other ~80k, plus a round each."
- [#391](https://github.com/radiusred/gh-codecrew/issues/391), 2026-10-09: a
  hub `init` scaffolded has no workflow, so its first task was refused
  `NO_CHECKS` after approval — "Observed 2026-10-09 on
  `radiusred/trading-hub` (#3 → #4 refused; #5 → #6 adds the workflows)."
- [#392](https://github.com/radiusred/gh-codecrew/issues/392), 2026-10-09:
  `gh pr merge` refused client-side on `BLOCKED` although the caller was the
  ruleset's bypass actor — "Observed 2026-10-09 on
  `radiusred/tradedesk-marketdata#88`".

One further report, [#387](https://github.com/radiusred/gh-codecrew/issues/387)
(radiusred/tradedesk-miner, 2026-10-06), was not adopted. Its one comment,
the same afternoon, reads "**Assessed 2026-10-06; held as a capture under the
2.1 bedding-in rule.**", and the goal defers it to 2.2 with #365 and #381.

The eight requirements, in short (the table below carries each in full):

- **M21-R1** — the release workflows on "a pinned ubuntu-24.04 image before
  the ubuntu-latest label moves on 2026-10-19" (adopts #380).
- **M21-R2** — `task finish` "merges in a way GitHub decides, so a ruleset
  bypass the caller legitimately holds is honoured instead of refused
  client-side"; `task start` and `task finish --dry-run` note a required
  check that has never reported (adopts #392).
- **M21-R3** — `milestone evidence` resolves "every github.com URL shape it
  meets … through the authenticated API" (adopts #388).
- **M21-R4** — `init`, `roles sync` and `migrate` write a commit body "naming
  the tool, its release and the protocol, as SPEC §4 requires" (adopts
  #389).
- **M21-R5** — "A fresh repository learns before its first task that task
  finish requires a reporting check" (adopts #391; "the scaffolded workflow
  is deferred to 2.2").
- **M21-R6** — the reviewer contract's task-start evidence for an App
  implementer, its housekeeping paragraph citing `init`'s body, and the
  implementer contract asking "the PR body to link the start record" (adopts
  #390).
- **M21-R7** — `task finish` reads the task and the captures it closes back
  from GitHub after the merge, and "does not close the task itself (that half
  of #386 is re-filed as a 2.2 capture)" (adopts #386).
- **M21-R8** — "v2.1.1 ships: version flip, tag on the merge commit,
  release.yml assets, evidence on the milestone"; the operator upgrades, and
  "the fleet's hubs pick up the contract changes through roles sync at their
  own pace."

Seven of the eight adopt one capture each; R8 is the release and adopts none.

The eight tasks, in merge order. Every one was opened by the operator's
account and carried out by the implementer seat (`radiusred-cody[bot]`),
which wrote each Plan, ran `task start`, opened the pull request and
rebase-merged it:

- **[#399](https://github.com/radiusred/gh-codecrew/issues/399) / PR
  [#403](https://github.com/radiusred/gh-codecrew/pull/403)** — M21-R6;
  adopts #390. Three commits, approved first round, merged
  2026-10-10T22:12:45Z. Two Decisions, one Deviation.
- **[#394](https://github.com/radiusred/gh-codecrew/issues/394) / PR
  [#402](https://github.com/radiusred/gh-codecrew/pull/402)** — M21-R1;
  adopts #380. Two commits, two review rounds, merged 22:18:48Z. One
  Decision.
- **[#395](https://github.com/radiusred/gh-codecrew/issues/395) / PR
  [#405](https://github.com/radiusred/gh-codecrew/pull/405)** — M21-R2;
  adopts #392. Three commits, approved first round, merged 22:32:31Z. Four
  Decisions; one Deviation, posted after QA's first round.
- **[#397](https://github.com/radiusred/gh-codecrew/issues/397) / PR
  [#404](https://github.com/radiusred/gh-codecrew/pull/404)** — M21-R4;
  adopts #389. One commit, two review rounds, merged 22:40:05Z. Three
  Decisions, one Deviation.
- **[#398](https://github.com/radiusred/gh-codecrew/issues/398) / PR
  [#410](https://github.com/radiusred/gh-codecrew/pull/410)** — M21-R5;
  adopts #391. Five commits, approved first round, merged 22:55:07Z. Two
  Decisions.
- **[#396](https://github.com/radiusred/gh-codecrew/issues/396) / PR
  [#406](https://github.com/radiusred/gh-codecrew/pull/406)** — M21-R3;
  adopts #388. Four commits, two review rounds, merged 23:01:21Z. Three
  Decisions.
- **[#400](https://github.com/radiusred/gh-codecrew/issues/400) / PR
  [#408](https://github.com/radiusred/gh-codecrew/pull/408)** — M21-R7;
  adopts #386. Three commits, approved first round, merged 23:05:03Z. Two
  Decisions.
- **[#401](https://github.com/radiusred/gh-codecrew/issues/401) / PR
  [#411](https://github.com/radiusred/gh-codecrew/pull/411)** — M21-R8, the
  v2.1.1 flip. Two commits, approved first round, merged 23:15:54Z; tagged
  at 23:16:25Z and released by 23:17:40Z. Two Decisions.

Then the housekeeping pull request #401's plan names as its Part C:
**[#412](https://github.com/radiusred/gh-codecrew/pull/412)**, `chore:
regenerate the contract history for v2.1.1` — one commit, no task, approved
by the reviewer seat's holder at 23:22:59Z and rebase-merged by its author at
23:23:28Z.

Seven of the eight merges were followed within seconds by a "Closed by `gh
codecrew task finish <n>`" comment on the adopted capture, naming the task,
the pull request and the merge commit; #401's
[release evidence](https://github.com/radiusred/gh-codecrew/issues/401#issuecomment-6103242243)
says PR #411 "was rebase-merged by `task finish 401`".

**The milestone ran in three and a half hours, by the timestamps.** The
operator's account created the eight tasks between 20:04:37Z and 20:05:09Z,
the first thirty-nine seconds after the milestone; each was labelled `cc:task`
and attached to #393 within three seconds, and each capture got its "Adopted
by" comment two seconds after its task was created. The trail then carries
nothing until #394's Plan edit at 22:02:56Z. Five tasks were planned and
started between then and 22:05:31Z — #394, #399, #396, #397 and #395, in that
order — #400 at 22:34:06Z, #398 at 22:42:43Z and #401 at 23:06:42Z. In every
task the implementer's Plan edit is the body's one revision after the
operator's original, and it precedes the `**Started by**` record by between
three and twenty-eight seconds. From the milestone opening to the v2.1.1 flip
merging: **three hours, eleven minutes and fifty-six seconds**; to the
release's publication, three hours, thirteen minutes and forty-two seconds; to
QA's second-round verdict, three hours, thirty-one minutes and twenty-three
seconds. From the first Plan edit to the release: one hour, fourteen minutes
and forty-four seconds.

**Pull requests ran in parallel.** #402, #403, #404, #405 and
#406 were all open from 22:12:05Z, when #406 opened, to 22:12:45Z, when #403
merged; three were open at once again from 22:39:18Z to 22:40:05Z (#404,
#406, #408) and from 22:47:09Z to 22:55:07Z (#406, #408, #410).

What exists afterwards that did not before.
[v2.1.1](https://github.com/radiusred/gh-codecrew/releases/tag/v2.1.1),
published 2026-10-10T23:17:40Z with five platform assets, the tag an annotated
object (`3a36f9a`) on the flip's merge commit
([9757001](https://github.com/radiusred/gh-codecrew/commit/9757001)), message
`v2.1.1 (protocol 2.1)`. The CHANGELOG's lead calls it "The first patch under
protocol 2.1, carrying what the fleet's first fortnight on it surfaced." In
it: `task finish` merges through GitHub's REST merge endpoint, pinned to the
head the gates read, and notes when GitHub reports the pull request `BLOCKED`
and when a required check has never reported — the second in `task start` too;
after the merge it reads the task and each capture it closed back from GitHub
and notes one still open; `milestone evidence` resolves github.com links
through a table of URL shapes and warns, without refusing, on a private
repository's unmapped shape; `init`, `migrate` and `roles sync` write the SPEC
§4 body on their commits; `init`'s `next:` lines and a `task new`/`task start`
note say a repository needs a check reporting on every pull request; the
reviewer contract names the start record as an App task's evidence and the
implementer contract asks the pull request body to link it; and every workflow
job runs on `ubuntu-24.04`. "No verb, flag or refusal code is added, removed
or repurposed, and the protocol stays at 2.1", as the lead puts it and as QA's
M21-R8 verdict checks. Three embedded contracts changed — implementer,
reviewer and qa — and #412 added their rows, with the unchanged
doc-synthesizer, coordinator and agents-file rows, to the contract history.
This hub's pointer still reads `codecrew: "2.1"`.

## Decisions

Nineteen Decision records, all on the tasks and all posted by the implementer
seat (`radiusred-cody[bot]`). Eight record a call the coordination layer made
— one on each of #394, #395, #396, #397, #398 and #400, two on #399 — and each
says so in its text: #394's and #395's, for instance, "recorded here by the
implementer, not re-opened", and #399's two "(coordination layer, M21 brief)".
One answers a reviewer's finding. None is on the milestone issue, none carries
a parenthetical qualifier on its label, and none is a gate resolution. The
operator's account posted no Decision in the milestone.

### The runners: pin now, move later

#394's [one Decision](https://github.com/radiusred/gh-codecrew/issues/394#issuecomment-6102651066)
pins every job "to `ubuntu-24.04` now, and move to Ubuntu 26 deliberately
later — the first of the two shapes #380 proposed. This is the coordination
layer's decision, made in the task Goal and M21-R1". **Trade-off:** "The
repository stops picking up image updates automatically; moving off 24.04
becomes a change someone has to make". **Rejected:** "Running the workflows
once on `ubuntu-26.04` and staying on `ubuntu-latest` — it would leave the
release build moving image on GitHub's schedule, not ours."

### The merge: GitHub decides, and the notes stay advisory

#395 carries four. The
[first](https://github.com/radiusred/gh-codecrew/issues/395#issuecomment-6102671225)
records two points the operator settled "on #392 and in M21-R2": `task
finish` merges through GitHub's REST merge endpoint "so GitHub evaluates the
request", and "The `--bypass` path, which went through `gh pr merge
--admin`, moves to the same endpoint with the same recorded comment and
receipt"; and the never-reported check is "a note, not a refusal — no new
code, no exit-status change — and nothing at all when the rules endpoint is
unavailable." **Trade-off:** "`gh pr merge`'s client-side check refused early
and locally on `BLOCKED`; the merge now asks GitHub and reports its answer,
so a refusal arrives from GitHub after the protocol's comments have posted".
**Rejected:** "`gh pr merge` only when `mergeStateStatus` is `CLEAN` (two
merge paths to test, and `--admin` would still be the only way past
`BLOCKED`, which the protocol forbids a crew identity); `--auto` (merges
after the protocol's one merge point)."

The [second](https://github.com/radiusred/gh-codecrew/issues/395#issuecomment-6102709390)
moves both venue methods to the endpoint and keeps them separate, and treats
"A 200 answering `merged: false`" as a failure. **Trade-off:** "pinning `sha`
is new — `gh pr merge` sent no expected head — so a push between the gates'
read and the merge is now refused by GitHub (409) instead of merging commits
nobody gated. That is a stricter merge, not a different record, so no
checkpoint." **Rejected:** "collapsing the two methods into one".

The [third](https://github.com/radiusred/gh-codecrew/issues/395#issuecomment-6102709550)
answers what `task start` compares against when no pull request exists yet:
"the heads of the **three newest pull requests into the default branch**
(open or closed), unioned". **Trade-off:** "a check that reports only on
some pull requests (a path-filtered workflow) can read as missing when the
last three skipped it". **Rejected:** "the default branch's head commit
(rebase merges leave pull_request-only checks off it, so every such check
would read as missing); parsing the workflow files".

The [fourth](https://github.com/radiusred/gh-codecrew/issues/395#issuecomment-6102709684)
makes the notes advisory "end to end": "**any** failed read … prints nothing,
because a context missing from a listing that could not be read is not a
context nothing reports", and "Only rulesets are read
(`rules/branches/{base}`), not classic branch protection, whose
required-check read needs administration access." **Trade-off:** "the
rule-type list can include types that never block a merge (`deletion`,
`non_fast_forward`)". **Rejected:** "a note for an unreadable policy (noise
on every Free-plan private repository)."

### Evidence: a table of shapes, and a warning for what it cannot map

#396's [first Decision](https://github.com/radiusred/gh-codecrew/issues/396#issuecomment-6102662279),
"(coordination layer, settled in the M21-R3 dispatch for #388)", has an
unmapped github.com path in a private repository warn rather than refuse: "A
refusal must mean "this evidence is dead"; an unmappable private path is
"this evidence is unverifiable", a different fact". A public repository's
404 still refuses `EVIDENCE_UNREACHABLE`. **Trade-off:** "a dead private
link of an unmapped shape passes with a warning; QA weighs it, as it weighs a
dead external link." **Rejected:** "fetching unmapped shapes as an
authenticated page (the Goal's last sentence): `gh` holds an API token, not a
browser session, and github.com does not serve private HTML pages to a bearer
token, so the fetch would 404 exactly as the anonymous one does." The
rejected option is the one #396's own goal had named: "a shape with no API
path goes through an authenticated fetch."

The [second](https://github.com/radiusred/gh-codecrew/issues/396#issuecomment-6102720802)
puts the table in the venue as data, held by a test that "fails on a row with
none", and lists six engineering calls — among them that numeric captures
must be digits, that github.com's own top-level pages "are never read as an
owner", and that "a comment anchor now resolves the comment itself …, so a
deleted comment now refuses where it used to pass on its issue."
**Trade-off:** "a branch name with a slash in `tree/feat/x` or
`blob/feat/x/file` is still split at the first slash (pre-existing,
unchanged)". **Rejected:** "routing every github.com link under a readable
repository through `repos/{o}/{r}` alone — it proves the repository, not the
cited run or comment, and would turn a dead job link into a pass."

The [third](https://github.com/radiusred/gh-codecrew/issues/396#issuecomment-6102739085)
records that "The qa contract's `warning:` sentence (`.codecrew/roles/qa.md`)
was updated in step with the code in PR #406", and that "no other contract,
doc or SPEC passage made the "warning means external" claim." It carries no
Trade-off or Rejected line. It was posted at 22:13:46Z, three seconds after
the commit it describes (e12e86b, authored 22:13:43Z) and a minute and
forty-one seconds after PR #406 opened.

### Commit bodies: one shape, and a hub name split at its slash

#397's [first Decision](https://github.com/radiusred/gh-codecrew/issues/397#issuecomment-6102667153),
"(coordination layer, settled on #389 and M21-R4)", fixes the shape: "every
commit `init`, `roles sync` and `migrate` make carries a body of one shape"
— the tool line with release and protocol, the files, then `Housekeeping
(SPEC §4); no task.` **Trade-off:** "the verbs' commit messages get longer,
and `roles sync`'s existing body changes shape; in return a reviewer checks
one shape against SPEC §4 for all three". **Rejected:** "leaving the body to
the operator (four spoke PRs were sent back for it on 2026-10-09); a body
naming only the tool and release without the protocol."

The [second](https://github.com/radiusred/gh-codecrew/issues/397#issuecomment-6102667304)
settles the details: the invocation as the verb read it, a tool line that
wraps at a space, "no trailing full stop" after a path, a count past 20
files, and a printed by-hand command that carries the body "as a second
`-m`, single-quoted so the newlines survive a paste into sh, bash, zsh or
fish". **Trade-off:** "a multi-line quoted argument in a printed command is
less tidy than a one-liner, but it is the only form that makes the commit
SPEC §4 asks for in one paste". **Rejected:** "`-F <file>` with the text
written by `init` (the file would be swept into the scaffold by `git add
-A`); a separate `-m` per line (git makes each its own paragraph, breaking
the shape)."

The [third](https://github.com/radiusred/gh-codecrew/issues/397#issuecomment-6102864803),
"(answering checky's finding on PR #404)", splits a word too long for any
line "just after a slash where one fits", so the longest `--hub` GitHub
accepts, "a 39-character owner plus a 100-character repository, 140
characters in all", spans two lines; "A word with no usable slash is cut at
the 100-character limit." **Trade-off:** "in that one case the hub's name
spans two lines of the body, which a reader has to put back together".
**Rejected:** naming the hub generically when it does not fit, and "refusing
an overlong hub, because it would change `init`'s refusal behaviour."

### The first task's check: say it, do not scaffold it

#398's [first Decision](https://github.com/radiusred/gh-codecrew/issues/398#issuecomment-6102966881),
"the coordination layer's, settled on #391 and in M21-R5", ships #391's
options 1 and 3 and defers option 2, `init` scaffolding an always-reporting
workflow, to protocol 2.2 as the capture #409. **Trade-off:** "a fresh
repository still cannot finish its first task without adding a workflow by
hand; it is now told so before the task is chartered and again when it is
started, instead of after review." **Rejected:** "scaffolding the workflow in
v2.1.1 — it writes files every adopter inherits and adds an argument to
`init`, a protocol-level choice that the 2.1 bedding-in rule keeps out of a
CLI patch."

The [second](https://github.com/radiusred/gh-codecrew/issues/398#issuecomment-6102968027)
reads `GET /repos/{o}/{r}/actions/workflows`, keeps active workflows under
`.github/workflows/`, and parses each one's `on:`; "`pull_request` counts,
`pull_request_target` does not." **Trade-off:** GitHub-managed workflows have
no file to read, "so a repository whose only reporting check is CodeQL
default setup would see the note". **Rejected:** "the contents listing of
`.github/workflows` — it answers 404 both for a repository with no workflow
directory (the case the note exists for) and for one the token cannot see",
and "Counting any listed workflow without reading its `on:` — a push-only
workflow reports nothing on a pull request".

### The start record: contract text only, history later

#399's two Decisions are both qualified in their text "(coordination layer,
M21 brief)". The
[first](https://github.com/radiusred/gh-codecrew/issues/399#issuecomment-6102661665)
leaves the contract history alone until R8: "`scripts/contract-history` runs
after the v2.1.1 tag". **Trade-off:** "Between this merge and R8, `roles
sync` from a v2.1.0 binary still offers the v2.1.0 text; the fleet picks up
the change with v2.1.1." **Rejected:** "Regenerating the table here — it
reads release tags, and there is no v2.1.1 tag to read." #412 is where that
landed.

The [second](https://github.com/radiusred/gh-codecrew/issues/399#issuecomment-6102661822)
has the reviewer contract's housekeeping paragraph cite "the body `init`
writes in the shape task #397 (M21-R4) is adding in parallel".
**Trade-off:** "If #397 settles on different wording, this citation must
follow it before v2.1.1 ships; both land in the same release." **Rejected:**
citing `roles sync`'s existing body, because "the operator named init's, and
#389 asks for exactly that." #404's change request reports the two in
agreement: "The reviewer-contract quotation and `init`'s normal body agree".

### The read-back: report, do not close

#400's [first Decision](https://github.com/radiusred/gh-codecrew/issues/400#issuecomment-6102904484),
"the coordination layer's, settled on #386 and in M21-R7", has `task finish`
read the task back and report a still-open one as a `note:`, and re-files the
self-close half of #386 as #407. **Trade-off:** "A task GitHub leaves open
still needs a human to close it by hand, but the verb now says so at once".
**Rejected:** "Closing the task from the verb now — it changes what the verb
writes and what "the only merge point" means (SPEC §4 has the task close
through its closing keyword), so it needs SPEC text, which the 2.1
bedding-in rule keeps out of a CLI patch."

The [second](https://github.com/radiusred/gh-codecrew/issues/400#issuecomment-6102934484)
sets the lines and the bound: the merge line "stops asserting the outcome",
and the read-back is "bounded at 4 reads 2 seconds apart (about 6 seconds of
waiting at most)", the same for each adopted capture. **Trade-off:** "A task
GitHub never closes costs the finish about 6 extra seconds before the note".
**Rejected:** "Keeping "closes via its closing keyword" on the live merge
line and adding the read-back after it — the line would still assert a close
the next line might contradict, which is the confusion #386 reports", and an
exponential back-off.

### v2.1.1: the flip moves the release where #361 and #374 did

#401's [first Decision](https://github.com/radiusred/gh-codecrew/issues/401#issuecomment-6103160893)
moves the shipped version "in exactly the places #361 and #374 moved it" —
the introduction's `**Shipped:**` line and `version` example, and the
README's current release — and leaves `CLI.md`'s `version` example and SPEC
§5's illustration at v2.1.0, as
[#359's Decision](https://github.com/radiusred/gh-codecrew/issues/359#issuecomment-5804336227)
kept `CLI.md`'s example on the v2.0.2 patch. **Trade-off:** "from the tag
onward `CLI.md`'s example is a patch behind what the binary prints; a
reviewer who reads it as a claim has a one-word change." **Rejected:**
"flipping `CLI.md` and SPEC here — it widens the flip past the precedent the
brief names, and touching SPEC's text in a patch release under the bedding-in
rule invites the question of whether the protocol moved."

The [second](https://github.com/radiusred/gh-codecrew/issues/401#issuecomment-6103161031)
keeps the README pointing a 2.0 repo at the CHANGELOG's 2.1.0 section, since
that "is where the steps for that move are". **Trade-off:** "a reader on
2.1.0 is not sent to the 2.1.1 section by name; it is the newest release
section, the first one they meet." **Rejected:** "the 2.1.0 and 2.1.1
sections", because "the 2.1.1 section has no migration content a 2.0 repo
needs that the 2.1.0 one lacks".

## Deviations

Four Deviation records: three on the tasks, all the implementer seat's, and
one on the milestone issue, this record's.

**[The reviewer contract had no task-start check to amend](https://github.com/radiusred/gh-codecrew/issues/399#issuecomment-6102661534).**
Posted on #399 at 22:04:18Z, eighteen seconds after its start record and
before its first commit. "The v2.1.0 reviewer contract has no task-start
check to amend — it never tells the reviewer to verify an assignee, or that
`task start` assigns the task", so the sentence goes in "as a short bullet of
its own under "Review against"". **Why:** "#390's premise … does not hold
against the embedded text"; the two reviews' reading "most likely came from
SPEC §4's lifecycle table, where "In progress" is signalled by "Assignee set"
— a row the 2.1 bedding-in rule keeps out of this patch." The SPEC row "is
flagged for the coordination layer as a possible 2.2 capture, not changed
here."

**[The by-hand command's staging, beyond the plan](https://github.com/radiusred/gh-codecrew/issues/397#issuecomment-6102693696).**
Posted on #397 at 22:08:11Z, before its commit was authored. The plan changed
only the message of the printed by-hand commit; running `migrate`'s
detached-HEAD command as printed "failed before any message was written:
`fatal: pathspec '.codecrew.yml' did not match any files`", because it staged
sources `git mv` had already removed. **Why:** "for `migrate` that
instruction could not be followed at all, so the body would never have
reached a commit. The fix is confined to the printed command and is covered
by the same test."

**[PR #405's body corrected after QA](https://github.com/radiusred/gh-codecrew/issues/395#issuecomment-6103345982).**
Posted on #395 at 23:32:22Z, fifty-nine minutes and fifty-one seconds after
the pull request merged and one minute and eighteen seconds after QA's first
round. The body "named
https://github.com/radiusred/gh-codecrew/issues/395 as the start record
instead of linking the **Started by** comment itself", and "was written
before PR #403 merged the implementer-contract text (#399) that asks the PR
body to link that exact comment." **Why:** "M21-R6's evidence is that every
M21 PR body links its task's start record, so the reviewer never has to look
for it. #405 was the one of eight that did not. The PR has already merged, so
the fix touches only the record." The pull request's edit history shows the
correction at 23:32:19Z, three seconds before the comment.

**[The record takes coordination-layer facts from its dispatch brief](https://github.com/radiusred/gh-codecrew/issues/393#issuecomment-6103422065).**
This record's own, at 23:42:57Z. The brief states the harness and model of
every seat, that the qa-contract sentence on PR #406 was added at the
coordination layer's request, and that each pull request's rebase before its
gating review met a CHANGELOG conflict. The comment sets out what the trail
supports: the commit trailers, the "Generated with Claude Code" line that
ends every task pull request body, the routing table, the Decision on #396,
and the force-push times — and that "No review or QA comment names its own
harness or model." **Why:** "the brief is not on the trail, so this comment
puts its account there, and the record cites this comment rather than the
brief, as the M18, M19 and M20 records did on #353, #382 and #384."

Two task pull request bodies state that none was needed — #402's "No
deviations." and #410's "No deviations." — and no other pull request body
reports a Deviation that is not also a comment.

## The gates

**#393's Gates section was left as the scaffold's placeholder** — "_What
"done" means beyond CI: e2e suites, manual UAT, sign-offs._", unedited, as
[M20's record](20-the-2-1-post.md#the-gates) says of #384. The ordering the
milestone needed is in the task goals instead: #394's "Smallest task in the
milestone; land it first", and #401's "Runs last, after R1–R7 are merged".
The second held: #411 merged after all seven others. The first did not:
#403 merged first, and #394's pull request six minutes and three seconds
later.

**No gate was raised.** No task used `gh codecrew checkpoint`, no issue
carries `cc:needs-decision`, and the dry run reports `gate no gate raised:
ok`. Four tasks left their Ask-the-human section as the placeholder (#394,
#398, #399, #401) and four wrote that there were none, each saying where the
answer was settled — #395's "both design points are settled by the operator
on #392 and in M21-R2", #396's "settled by the coordination layer in the
dispatch for M21-R3", #397's "the body's shape is settled on #389 and
M21-R4", and #400's "the Goal and the operator's settlement on #386 / M21-R7
decide the open points".

**The release's operator leg is on the milestone issue.** The
[leg](https://github.com/radiusred/gh-codecrew/issues/393#issuecomment-6103248238),
at 23:19:12Z, one minute and thirty-two seconds after the release was
published: `upgraded from v2.1.0 to v2.1.1`; `gh codecrew version` printing
`v2.1.1 (protocol 2.1)`; `roles diff` for all five roles each "`matches the
embedded v2.1.1 contract`", so that "the implementer/reviewer drift `status`
had reported since #403 merged (the #399 text against the installed v2.1.0)
is gone"; and `status` showing "every M21 task done, no gates raised", with
"the expected agents-file drift" the only line left. It ends: "The
contract-history table for v2.1.1 is Part C of #401, in flight as a
housekeeping PR." #412 was opened six seconds later.

What gated the work: CI on every pull request — `Lint commit messages` and
`Go build and test`, which ran on `ubuntu-24.04` for every merged head except
#403's, which merged before #402 and ran on `ubuntu-latest` — the reviewer
seat holder's approval on each of the nine, the release workflow's
contract-history guard on the v2.1.1 tag, and the QA verdicts.

## The review rounds

Nine pull requests, **twelve review submissions**, every one by the reviewer
role holder (`radiusred-checky[bot]`), an App distinct from the authoring
seat: **nine approvals, one per pull request, all standing**, and **three
change requests**, each resolved in a later round. Six pull requests were
approved in their first round; #402, #404 and #406 took two. No review was
dismissed, and no inline review comment was posted. Each approval was given
on the head that merged.

- **PR [#403](https://github.com/radiusred/gh-codecrew/pull/403)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/403#pullrequestreview-5480988781),
  "Approved — PR #403 correctly delivers M21-R6 / task #399." It quotes
  CLI.md's `task start` rule against the new bullet, and reports that
  "Grepping the reviewer, implementer, coordinator, QA, and doc-synthesizer
  contracts found no conflicting assignment claim."
- **PR [#402](https://github.com/radiusred/gh-codecrew/pull/402)** — two
  rounds. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/402#pullrequestreview-5480973170)
  has one finding: the CHANGELOG entry "ends `(#394).`, rather than ending
  `(#394)`, and it never identifies adopted capture `#380`." The
  implementer's
  [reply](https://github.com/radiusred/gh-codecrew/pull/402#issuecomment-6102692324)
  followed forty-six seconds later. The
  [approval](https://github.com/radiusred/gh-codecrew/pull/402#pullrequestreview-5481004897)
  checks the rebase — "Rebasing on `main` introduced no unrelated change" —
  and accepts the body's pre-rebase run links: "the PR's attached checks are
  the authoritative record for the current head."
- **PR [#405](https://github.com/radiusred/gh-codecrew/pull/405)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/405#pullrequestreview-5481045422),
  "Approved — M21-R2 is sound." It reads the REST call's shape and the `sha`
  pin — "A moved head therefore returns GitHub's actionable 409 response" —
  and reports a live read-only probe: the reviewer identity's `task finish
  395 --dry-run` "correctly stops at `NOT_OWNER` before policy reads." It
  states that "the start record is linked from the PR"; the body then linked
  the task issue, which QA's first round found.
- **PR [#404](https://github.com/radiusred/gh-codecrew/pull/404)** — two
  rounds. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/404#pullrequestreview-5481029923):
  "`wrapLine` intentionally leaves an overlong word whole", so a built-source
  `init --hub` with a valid 39-character owner and 100-character repository
  "emits the `owner/repo` token as a 140-character line in the commit body",
  breaking "M21-R4's and the recorded Decision's hard “every line ≤100”
  condition". It also reports running all three verbs from a built binary,
  and the printed `init` and `migrate` commands through `sh -c`. The
  [approval](https://github.com/radiusred/gh-codecrew/pull/404#pullrequestreview-5481064399):
  "The slash-split continuation reconstructs the exact hub name and complete
  tool line", and "The new regression coverage would fail under the prior
  overlong-word behavior."
- **PR [#410](https://github.com/radiusred/gh-codecrew/pull/410)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/410#pullrequestreview-5481099837).
  It checks the trigger parser — "yaml.v3 node parsing retains the literal
  `on` key, so the YAML 1.1 boolean-key trap is not present" — notes that
  "`task new` and `task start` have no dry-run surface, so I did not invoke
  their writing paths", and rules on a commit type: "`feat(tracker)` is
  acceptable here: it is the internal additive read supporting the additive
  CLI notice, not a semver or protocol claim."
- **PR [#406](https://github.com/radiusred/gh-codecrew/pull/406)** — two
  rounds. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/406#pullrequestreview-5481082570)
  has two findings. The body's "For the reviewer" section said the qa
  contract was left alone, which the diff contradicts: "The PR body is part of
  the task's audit record and must match the final diff." And a test listed
  an App URL without asserting where it went: "The test therefore passes if
  an accidental new API mapping routes the App page through `APIReachable`".
  The
  [approval](https://github.com/radiusred/gh-codecrew/pull/406#pullrequestreview-5481118468)
  finds both fixed — "Removing `apps` from the reserved owners would therefore
  fail the test" — and reports "the rebase base is 8d9c907 with no unrelated
  PR changes."
- **PR [#408](https://github.com/radiusred/gh-codecrew/pull/408)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/408#pullrequestreview-5481127411),
  "Approved — M21-R7 is satisfied." It rules on the changed output line: "SPEC
  §4 still correctly describes the mechanism as "issue closed via closing
  keyword"", and "this is an M21 CLI patch rather than a protocol or
  downstream-parser contract change."
- **PR [#411](https://github.com/radiusred/gh-codecrew/pull/411)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/411#pullrequestreview-5481153361).
  It checks the patch's bound: "The v2.1.0 and head unique source
  refusal-code sets are identical (49 codes), and no flag declaration
  changed", and runs the release guard under a temporary local tag: "Part C
  should therefore regenerate history after tagging for the next release."
- **PR [#412](https://github.com/radiusred/gh-codecrew/pull/412)** —
  [approved](https://github.com/radiusred/gh-codecrew/pull/412#pullrequestreview-5481171647)
  as housekeeping: "Approved — this satisfies the SPEC §4 housekeeping
  fence." Rerunning `scripts/contract-history` "produced exactly the six
  existing `v2.1.1` rows and no other diff."

**What the three change requests were about.** They carry four findings. One
is record mechanics: a CHANGELOG entry's ending and its missing capture. One
is a hard bound broken at the edge of the input space: the 140-character
line. One is a pull request body left stale by a change made after it was
written. One is a test that did not cover what it claimed. Each change
request was answered with a new head, and by an implementer comment naming
it.

**Ten force-pushes; eight have no comment.** Beyond the two that answer
#402's and #404's change requests, eight follow a sibling pull request's
merge, by thirty-four seconds to ten minutes and three seconds: #402 at
22:13:24Z, #404 at 22:19:35Z and 22:33:07Z, #405 at 22:28:51Z, #406 at
22:40:53Z and 22:55:45Z, and #408 at 22:48:17Z and 23:01:55Z. None has a
comment; two approvals mention a rebase, as quoted above. #406's answer to
its change request was an ordinary push of a new `test:` commit at
22:48:22Z.

## QA: two rounds, one remedy, and the verdict that supersedes

The qa role holder (`radiusred-testy[bot]`) wrote two comments on #393, the
shape [M17's record](17-the-frontier-builds-the-coordinator.md#qa-two-rounds-one-remedy-and-the-verdict-that-supersedes)
describes for its own milestone.

**The first round.**
[One comment](https://github.com/radiusred/gh-codecrew/issues/393#issuecomment-6103336491)
at 23:31:04Z, seven minutes and thirty-six seconds after #412 merged, with one
verdict per requirement: seven `satisfied` and **M21-R6 — not satisfied.**
Each verdict opens with the build line — "On merged main, `go test ./...`,
`go vet ./...`, gofmt, and `go build ./cmd/codecrew` passed.", to which R8's
adds the release's contract-history guard — then says what the seat checked
beyond the tests, and where it could not:

- **M21-R1** — "exactly three `runs-on` declarations, all `ubuntu-24.04` with
  the #380 migration comment"; through the jobs API, every M21 pull request
  run from #402 on reported `ubuntu-24.04`, while "the earlier PR #403 run
  reported `ubuntu-latest` because it ran before #402 merged", and the
  release run's job "reported `ubuntu-24.04` and succeeded."
- **M21-R2** — "the executable merge path is `PUT
  repos/{o}/{r}/pulls/{n}/merge -f merge_method=rebase -f sha=<gated head>`,
  with no `gh pr merge` call"; a live `task finish 401 --dry-run` "printed
  `refused[CLOSED]` at the first gate and no policy note because later gates
  were not reached"; and "`git diff v2.1.0 v2.1.1 -- SPEC.md` was empty and
  both releases expose the same 49 refusal codes."
- **M21-R3** — with the released binary, "`milestone evidence 18` resolved
  all 10 citations across 11 issues and `milestone evidence 21` resolved all
  3 across 9", and: "I had no private hub for a live probe, so that path is
  test-proven".
- **M21-R4** — released-binary `init` runs in scratch repositories, hub and
  spoke, "made `chore: scaffold codecrew` commits whose bodies named the
  exact invocation, v2.1.1, protocol 2.1, every written file and
  `Housekeeping (SPEC §4); no task.`, with no line over 100 columns".
- **M21-R5** — both scratch `init` runs printed the requirement; "I did not
  run `task new` or `task start` against a workflow-less live repository
  because neither has a read-only mode", so that note is tested, not
  observed.
- **M21-R6** — the contracts carry the text and `roles diff` matches, but
  "comparing every task's actual `Started by` comment URL with all eight M21
  PR bodies found only 7/8 exact links: PR #405 calls
  `https://github.com/radiusred/gh-codecrew/issues/395` the start record but
  does not link the actual comment". "The required PR-body evidence is
  therefore incomplete."
- **M21-R7** — six focused tests read and run; "`CloseIssue` only for adopted
  captures—not the task"; and "Part C was housekeeping and therefore supplied
  no live `task finish` observation".
- **M21-R8** — the tag "peels to `9757001`", the release has "exactly five
  platform assets", and "The bedding-in probes found protocol constants still
  2.1, no SPEC or verb-dispatch diff, no flag-declaration changes, and
  identical 49-code refusal sets versus v2.1.0." It closes on the captures:
  all seven adopted ones closed "with `task finish` comments naming their
  tasks; #407 and #409 remain open captures".

**The remedy.** Seventy-five seconds after the verdict, at 23:32:19Z, the
implementer seat edited PR #405's body to link the start comment; three
seconds later it posted the Deviation on #395 described above. Nothing on the
trail records who asked for either; *undocumented, inferred from the
timestamps*.

**The second round.**
[One comment](https://github.com/radiusred/gh-codecrew/issues/393#issuecomment-6103367665)
at 23:35:21Z: **M21-R6 — satisfied.** It supersedes "my first-round
finding" — "for this requirement alone; the other seven verdicts in that
comment stand." It reports re-running the build and `milestone evidence
21`, rereading M21-R6 "verbatim", and checking PR #405 through "Its REST
record, timeline, and two GitHub user-content snapshots": the body "was
written before #403 merged and was corrected after #405 merged; the only
prose/link change is the start URL from the task issue to that comment. The
edit also added one trailing blank line, but changed no other content." Its
exhaustive probe derives each expected URL from the task's actual start
comment and compares all eight bodies: "8/8 are exact matches."

**Every standing verdict is `satisfied`**: M21-R1 to R5, R7 and R8 in the
first-round comment, M21-R6 in the second. No requirement was struck.

## Requirement outcomes

The status column is the **standing** verdict word as the qa role holder
wrote it, verbatim and unqualified: M21-R1 to R5, R7 and R8 in comment
[6103336491](https://github.com/radiusred/gh-codecrew/issues/393#issuecomment-6103336491),
and M21-R6 in comment
[6103367665](https://github.com/radiusred/gh-codecrew/issues/393#issuecomment-6103367665),
which supersedes the first round's `not satisfied` for that requirement
alone.

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M21-R1 — The release workflows run on a pinned ubuntu-24.04 image before the ubuntu-latest label moves on 2026-10-19 (adopts #380). | [#394](https://github.com/radiusred/gh-codecrew/issues/394) / PR [#402](https://github.com/radiusred/gh-codecrew/pull/402) | `satisfied` | Task closed; #380 closed by its finish. Two review rounds (the CHANGELOG entry's ending and its capture). One Decision (pin now, the coordination layer's). The release job's pin was first exercised by the v2.1.1 release run |
| M21-R2 — task finish merges in a way GitHub decides, so a ruleset bypass the caller legitimately holds is honoured instead of refused client-side, and the dry run reports when the merge will rely on one; task start and task finish --dry-run say as a note: when a check required by a branch policy has never reported on the PR (adopts #392). | [#395](https://github.com/radiusred/gh-codecrew/issues/395) / PR [#405](https://github.com/radiusred/gh-codecrew/pull/405) | `satisfied` | Task closed; #392 closed by its finish. Approved first round. Four Decisions (the REST endpoint and both notes, the coordination layer's; the `sha` pin; three newest PR heads for `task start`; advisory reads, rulesets only). One Deviation (its body's start link, corrected after merge) |
| M21-R3 — milestone evidence resolves every github.com URL shape it meets — Actions run and job URLs, a repository root — through the authenticated API the way it resolves issues and PRs, so a private hub's record is not refused EVIDENCE_UNREACHABLE on reachable evidence (adopts #388). | [#396](https://github.com/radiusred/gh-codecrew/issues/396) / PR [#406](https://github.com/radiusred/gh-codecrew/pull/406) | `satisfied` | Task closed; #388 closed by its finish. Two review rounds (a stale PR body; a test that did not assert the App fetch). Three Decisions (warn on an unmapped private shape, the coordination layer's; the shape table; the qa contract's `warning:` sentence). Private path test-proven, not probed live |
| M21-R4 — init, roles sync and migrate write a body on the commit they make naming the tool, its release and the protocol, as SPEC §4 requires of a housekeeping change, and CLI.md states the body each writes (adopts #389). | [#397](https://github.com/radiusred/gh-codecrew/issues/397) / PR [#404](https://github.com/radiusred/gh-codecrew/pull/404) | `satisfied` | Task closed; #389 closed by its finish. Two review rounds (a 140-character line at GitHub's maximum name lengths). Three Decisions (the shape, the coordination layer's; the rendering; the slash split). One Deviation (`migrate`'s by-hand command staged removed sources) |
| M21-R5 — A fresh repository learns before its first task that task finish requires a reporting check: init's next: lines say so, task new and task start emit a note: when the target repository has no workflow on pull_request, and CLI.md cross-references NO_CHECKS from init and task finish (adopts #391; the scaffolded workflow is deferred to 2.2). | [#398](https://github.com/radiusred/gh-codecrew/issues/398) / PR [#410](https://github.com/radiusred/gh-codecrew/pull/410) | `satisfied` | Task closed; #391 closed by its finish. Approved first round. Two Decisions (option 2 deferred to 2.2 as [#409](https://github.com/radiusred/gh-codecrew/issues/409), the coordination layer's; the workflows listing and the `on:` parse). The `task new`/`task start` note is test-proven, not observed live |
| M21-R6 — The reviewer contract's task-start check says what the evidence is when the implementer is an App (no assignee: the start record and the linked branch), its housekeeping paragraph cites the commit body init writes, and the implementer contract asks the PR body to link the start record (adopts #390). | [#399](https://github.com/radiusred/gh-codecrew/issues/399) / PR [#403](https://github.com/radiusred/gh-codecrew/pull/403) | `satisfied` | Task closed; #390 closed by its finish. Approved first round. Two Decisions (history regenerated at R8; the `init` body cited in #397's shape), one Deviation (no task-start check to amend). Verdicted `not satisfied` in QA's first round for PR #405's body; satisfied in the second after the correction and #395's Deviation |
| M21-R7 — After the merge, task finish reads the task — and the captures it closes by --adopts — back from GitHub, retrying briefly, and says plainly when one is still open instead of asserting that it closes; it does not close the task itself (that half of #386 is re-filed as a 2.2 capture) (adopts #386). | [#400](https://github.com/radiusred/gh-codecrew/issues/400) / PR [#408](https://github.com/radiusred/gh-codecrew/pull/408) | `satisfied` | Task closed; #386 closed by its finish. Approved first round. Two Decisions (read back, do not close, the coordination layer's, re-filed as [#407](https://github.com/radiusred/gh-codecrew/issues/407); the output lines and the four-read bound). No live read-back observed in the milestone |
| M21-R8 — v2.1.1 ships: version flip, tag on the merge commit, release.yml assets, evidence on the milestone; the operator upgrades the installed extension, and the fleet's hubs pick up the contract changes through roles sync at their own pace. | [#401](https://github.com/radiusred/gh-codecrew/issues/401) / PR [#411](https://github.com/radiusred/gh-codecrew/pull/411), with the housekeeping PR [#412](https://github.com/radiusred/gh-codecrew/pull/412) | `satisfied` | Task closed. Tag `v2.1.1` on [9757001](https://github.com/radiusred/gh-codecrew/commit/9757001); [release](https://github.com/radiusred/gh-codecrew/releases/tag/v2.1.1) published 2026-10-10T23:17:40Z with five assets from [run 38094448333](https://github.com/radiusred/gh-codecrew/actions/runs/38094448333) on `ubuntu-24.04`; operator leg on #393. Two Decisions (where the flip moves the release; the README's 2.1.0 pointer) |

**Captures adopted and closed by the tasks' finishes — seven**, each with an
"Adopted by" comment from the operator's account when its task was created
and a "Closed by `gh codecrew task finish <n>`" comment naming the task, the
pull request and the merge commit:
[#390](https://github.com/radiusred/gh-codecrew/issues/390) (by #399),
[#380](https://github.com/radiusred/gh-codecrew/issues/380) (by #394),
[#392](https://github.com/radiusred/gh-codecrew/issues/392) (by #395),
[#389](https://github.com/radiusred/gh-codecrew/issues/389) (by #397),
[#391](https://github.com/radiusred/gh-codecrew/issues/391) (by #398),
[#388](https://github.com/radiusred/gh-codecrew/issues/388) (by #396) and
[#386](https://github.com/radiusred/gh-codecrew/issues/386) (by #400). The
oldest, #380, was open fifteen days, twenty-three hours, seventeen minutes
and twenty-nine seconds; the annotation it quotes dates the label's
migration from October 19, 2026.

**Captures filed during the milestone and left open — two**, both by the
implementer seat, unlabelled, each re-filing the half of an adopted capture
its task's Decision deferred to protocol 2.2:
[#407](https://github.com/radiusred/gh-codecrew/issues/407) at 22:34:21Z,
"task finish closes the task itself when GitHub does not act on the closing
reference", from #386 by #400; and
[#409](https://github.com/radiusred/gh-codecrew/issues/409) at 22:42:23Z,
"init scaffolds an always-reporting workflow so a new repository can finish
its first task", from #391 by #398. From #393's creation to the counting
instant, the only issues or pull requests created in this repository are
#393, its eight tasks, the nine pull requests and these two; an
organisation-wide issue search over the same window, as this record's token
sees the organisation, returns the same twenty.

**The captures the goal defers to 2.2** —
[#387](https://github.com/radiusred/gh-codecrew/issues/387),
[#365](https://github.com/radiusred/gh-codecrew/issues/365),
[#381](https://github.com/radiusred/gh-codecrew/issues/381) and
[#377](https://github.com/radiusred/gh-codecrew/issues/377)–[#379](https://github.com/radiusred/gh-codecrew/issues/379)
— are all open at the counting instant, as are #407 and #409.

## The front door

The doc-synthesizer contract sends this hub's front-door list to its
`.local.md`: "the README's proof points, and the introduction's release,
verbs and refusal codes". #411 moved both documents' release to v2.1.1 under
#401's first Decision. Read against what the milestone delivered, at the
counting instant:

- **`README.md`** — true as it stands. "The current release, v2.1.1,
  implements protocol 2.1", with the pointer to the CHANGELOG's 2.1.0 section
  #401's second Decision kept; the routing snippet, which presents itself as
  this repository's `.codecrew/config.yml` "as it stands today", is identical
  to that file's `roles:` section; the `gh` floor sentence — "`task finish`
  reads `gh pr checks --json`" — still describes the source, which reads the
  checks that way and refuses `GH_TOO_OLD` below 2.50.0; and "`init` does not
  write a workflow", which M21-R5 left true by deferring the scaffolded
  workflow to #409.
- **`docs/introduction.md`** — the release ("**Shipped:** v2.1.1 of the `gh`
  extension, implementing protocol 2.1", and the `version` example) is the
  installed extension's; the verb list matches `gh codecrew help` from v2.1.1;
  the six `--dry-run` verbs are unchanged, since no flag declaration changed,
  as the reviewer on #411 and QA's M21-R8 verdict report; and the refusal-code
  section points at SPEC §10 and states no count, the catalogue being the same
  49 codes as v2.1.0's. Its sentence that `task finish` closes an adopted
  capture "after the merge with a comment naming the task, the pull request
  and the commit the merge left" is still true; M21-R7 added a read-back after
  that close.

Neither document makes a claim the milestone made false, so this PR changes
neither. SPEC §4 lists the record's diff as the document, its ROADMAP row and
the front-door claims the milestone changed, "and nothing else", and says
the housekeeping path carries "no record entry", so this PR adds no CHANGELOG
entry — as M18's, M19's and M20's record pull requests added none.

## Protocol-discipline observations

Seven things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **A contract change and the pull requests it governs were in flight
  together, and the one it caught was fixed on the record, not the code.**
  #403 merged the implementer-contract text asking a pull request body to
  link the start record at 22:12:45Z, while #402, #404, #405 and #406 were
  open with bodies already written. Three of the four linked the start
  comment anyway; #405 linked the task issue. The reviewer's approval on
  #405 says "the start record is linked from the PR"; QA's first round,
  comparing exact URLs, did not accept it. The remedy was an edit to a merged
  pull request's body, and QA's second round verified it through GitHub's
  edit history rather than the body alone.
- **The coordination layer's calls reach the trail through the implementer.**
  Eight of the nineteen Decisions record a call the coordination layer made,
  and every one was posted by `radiusred-cody[bot]`; the operator's account
  posted the milestone, the tasks, the "Adopted by" comments and the release
  leg, and no record. Three Decisions point at a source the trail does not
  hold: #396's first, "settled in the M21-R3 dispatch", and #399's two, "M21
  brief". One of the eight rejects the option its own task goal named: #396's
  goal asked for an authenticated fetch of unmapped shapes, and the
  coordination layer's Decision gives the reason it cannot work.
- **#401 closed when its flip merged, though its pull request said it would
  not.** PR #411's body: "This PR does not complete the task. #401 stays open
  until Part B is done". #401 closed at 23:15:55Z, a second after the merge,
  and #401's evidence comment says so: "#401 closed through the linked branch
  at 23:15:55Z." The tag, the release run, the evidence and Part C all
  followed on a closed task. QA's live `task finish 401 --dry-run` was
  refused `CLOSED` at its first gate, so the dry run of the milestone's last
  task showed none of M21-R2's notes.
- **Three contracts changed; the release's account names one task's.** #412's
  body says the implementer, reviewer and qa hashes are new "because M21
  changed those three contracts" — the qa contract by #396. The CHANGELOG
  lead tells a hub that wants "#399's contract clarifications" to run `roles
  sync`, the operator leg reports "the implementer/reviewer drift" gone, and
  #401's evidence says "#399 changed its contracts". None of the three is
  false; none mentions the qa contract's new sentence.
- **A Decision withdrew a phrase no Decision had used.** #397's third
  Decision says "My earlier Decision said an overlong word "stays whole";
  that is withdrawn". Neither earlier Decision on #397 contains the phrase,
  and neither was edited; the reviewer's change request found the behaviour
  in the code — "`wrapLine` intentionally leaves an overlong word whole".
  *Undocumented; inferred from the comments and their edit histories.*
- **A placeholder in a plan notified a stranger.** #395's Plan says the note
  is "asked of GitHub as `@caller`", written as a mention; the issue's
  timeline records a `mentioned` and a `subscribed` event for the GitHub
  account `caller`, which is not a member of the crew, at 22:05:39Z, twenty
  seconds after the Plan edit. Nothing on the trail remarks on it.
- **A SPEC row was named as the likely root of #390, and it was already
  captured before M21.** #399's Deviation and #403's body both say SPEC §4's
  lifecycle table, where "In progress" is "Assignee set", probably produced
  the reviews' reading, and flag it "for the coordination layer as a possible
  2.2 capture". Capture
  [#316](https://github.com/radiusred/gh-codecrew/issues/316), filed
  2026-09-07 and open at the counting instant, already records it: "SPEC §4's
  lifecycle table infers `in progress` from the task having an assignee", and
  "GitHub refuses an App as an assignee". Neither the Deviation nor the body
  cites #316; no new capture was filed in the milestone's window, and the
  goal's list of what 2.2 inherits does not name #316.

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **The conversation that scoped the milestone is on the record only through
  #393's goal**, and its one line of provenance, "(operator, 2026-10-10: CLI
  patch releases do not affect it)". Why the seven captures were taken
  together now, and why #387 waits while #389 and #390 do not, is not on the
  trail beyond #387's "held as a capture under the 2.1 bedding-in rule".
  *The reasoning is recorded; the deliberation is not*, as
  [M20's record](20-the-2-1-post.md#what-the-record-does-not-contain) says of
  its own opening.
- **The dispatch briefs three Decisions cite are not on the trail** — #396's
  "M21-R3 dispatch", #399's "M21 brief" — nor is "the precedent the brief
  names" in #401's first Decision.
- **Which harness and model ran each session is not recorded by the seats.**
  This record's Deviation carries the coordination layer's account and the
  trail's partial support for it; the commits carry a trailer and the pull
  request bodies a line, and no review or QA comment names its model.
- **Nothing on the trail says why the milestone paused** for the hour and
  fifty-seven minutes between the tasks' creation and the first Plan edit,
  or who asked for the qa-contract sentence on PR #406 and for PR #405's
  correction; the Deviation carries the brief's account of the first of
  those two.
- **The eight force-pushes that follow a sibling merge have no comment.** The
  approvals on #402 and #406 mention a rebase; no comment mentions a
  conflict, which the brief reports on every one.
- **Two behaviours shipped in v2.1.1 were not observed live in the
  milestone.** QA's M21-R7 verdict says "Part C was housekeeping and
  therefore supplied no live `task finish` observation" of the read-back, and
  its M21-R3 verdict says the private-hub path "is test-proven". Which binary
  ran each `task finish` — the installed v2.1.0 or a build of the branch — is
  not recorded, and none of the finishes posted its output.
- **The first private hub behind #388 is not named.** The capture says "in a
  private hub"; only its second shape, two days later, names
  radiusred/trading-hub.
