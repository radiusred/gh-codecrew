# M11: Housekeeping

Tracking issue: [#233](https://github.com/radiusred/gh-codecrew/issues/233) ·
Synthesized 2026-09-04 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #233's **three** requirements — all three present when
the milestone opened, **M11-R1 amended by the operator on 2026-09-04 and its
audience refined two minutes later**, with the original clause quoted in the
issue's Amendments section — its Gates section and its **one** QA round; the
four delivery task issues, two in this repository
([#235](https://github.com/radiusred/gh-codecrew/issues/235),
[#236](https://github.com/radiusred/gh-codecrew/issues/236)) and two in the
spoke [radiusred/codecrew-www](https://github.com/radiusred/codecrew-www)
([#21](https://github.com/radiusred/codecrew-www/issues/21),
[#22](https://github.com/radiusred/codecrew-www/issues/22)); their four merged
PRs; the **thirteen Decision comments and no Deviation comment** across the
milestone issue and those task issues; the **six review submissions** on the
four PRs; the three adopted backlog captures the merges closed and the two the
milestone filed; and #233's GitHub edit history, which timestamps both
amendment edits. The house form is
[M10's document](10-protocol-bookkeeping-from-the-field.md), with
[M9's](9-the-docs-at-codecrew-works.md) behind it; the record standard is the
one the reviewer set on their PRs,
[#204](https://github.com/radiusred/gh-codecrew/pull/204),
[#206](https://github.com/radiusred/gh-codecrew/pull/206) and
[#232](https://github.com/radiusred/gh-codecrew/pull/232) — verdict words
verbatim, no ordinal claim without every prior record linked, no judgement the
record does not carry.

The gather from `gh codecrew milestone close 11` is again not part of the raw
material, for the reason [M8's](8-a-product-home-page-for-codecrew-works.md),
[M9's](9-the-docs-at-codecrew-works.md) and
[M10's](10-protocol-bookkeeping-from-the-field.md) records all give: `--dry-run`
stops at the second gate, naming the task that writes this file.

```
gate milestone open: ok
gate tasks closed: refused[OPEN_TASKS]: tasks not closed: radiusred/gh-codecrew#239 (ready)
gate requirements declared: not reached
gate QA verdicts: not reached
gate milestone document: not reached
dry run: nothing written — the live verb stops at the first refusal above
```

Every record below was read from its issue, PR or timeline directly.

**This PR adds the M11 ROADMAP row; it does not flip one** — the convention
M10-R1 introduced, which
[M10's record PR](https://github.com/radiusred/gh-codecrew/pull/232) was
written under. The convention is on `main` and not yet in a release: the latest
tagged version is `v1.1.0` (2026-08-31), and
[#208](https://github.com/radiusred/gh-codecrew/issues/208) / PR
[#216](https://github.com/radiusred/gh-codecrew/pull/216) sits under
`## [Unreleased]`. So the installed binary that opened this milestone —
`gh codecrew version` reports `v1.1.0 (protocol 1.0)` — still appended an `Open`
row to a local `ROADMAP.md` and told the implementer to carry it. The operator
[discarded it](https://github.com/radiusred/gh-codecrew/issues/233#issuecomment-5546407791),
as M10's was discarded, "because the convention it belongs to is gone from
`main`'s source and survives only in the unreleased binary gap".

## Goal and outcome

The M10 close left three backlog captures on the table, all small and all about
the hub keeping only what it needs to, and #233's goal is the list: the README
duplicates codecrew.works now that the home page (M8) and the docs section (M9)
exist, so it is cut down and sends everyone else to the site
([#224](https://github.com/radiusred/gh-codecrew/issues/224)); the docs
introduction gets a way to a routing-table example
([#227](https://github.com/radiusred/gh-codecrew/issues/227)); and the two
CodeCrew repositories stop carrying their own commitlint workflow and config now
that `radiusred/.github` publishes the shared action
([#225](https://github.com/radiusred/gh-codecrew/issues/225)). The README cut
has a spoke-side consequence — codecrew-www's sync translates README anchors to
home-page anchors — so the spoke followed in the same milestone.

Two of the three captures were filed by M10 and are named in
[its record](10-protocol-bookkeeping-from-the-field.md#requirement-outcomes) as
such: #224 and #227. The third came from elsewhere in the organisation: #225
records that it was "captured from radiusred/ops#17 (M5, 'Shared commitlint
across the org', requirement M5-R3)", where the operator had already migrated
three non-hub repositories and handed the two CodeCrew repositories to this hub.
**All three captures were adopted and all three were closed by the merges**; two
new ones were filed and remain open.

Three requirements were present when the milestone opened at 21:01:15Z on
2026-09-04. M11-R1 was amended at 22:10:08Z and its audience refined at
22:12:09Z, both by the operator, both timestamped in #233's edit history. Four
tasks delivered the work — two in the hub and two in the spoke — and this
document is the fifth. The last delivery merged at 22:36:08Z, one hour
thirty-four minutes and fifty-three seconds after the milestone opened.

- **[#236](https://github.com/radiusred/gh-codecrew/issues/236) / PR
  [#237](https://github.com/radiusred/gh-codecrew/pull/237)** — M11-R3, hub
  half. One commit, approved first round, merged 21:13:43Z
  ([66632eb](https://github.com/radiusred/gh-codecrew/commit/66632eb)). Two
  Decisions. Adopts #225.
- **[codecrew-www#21](https://github.com/radiusred/codecrew-www/issues/21) / PR
  [codecrew-www#23](https://github.com/radiusred/codecrew-www/pull/23)** —
  M11-R3, spoke half. One commit, approved first round, merged 21:18:11Z
  ([0f9eb8a](https://github.com/radiusred/codecrew-www/commit/0f9eb8a)). One
  Decision.
- **[#235](https://github.com/radiusred/gh-codecrew/issues/235) / PR
  [#238](https://github.com/radiusred/gh-codecrew/pull/238)** — M11-R1 and
  M11-R2. Four commits, two reviewer rounds either side of the operator's tone
  pass, merged 22:34:52Z
  ([cc41991](https://github.com/radiusred/gh-codecrew/commit/cc41991)). Seven
  Decisions — four by the seat and three by the operator. Adopts #224 and #227.
- **[codecrew-www#22](https://github.com/radiusred/codecrew-www/issues/22) / PR
  [codecrew-www#25](https://github.com/radiusred/codecrew-www/pull/25)** —
  M11-R1, spoke half. One commit, approved first round, merged 22:36:08Z
  ([009a7ea](https://github.com/radiusred/codecrew-www/commit/009a7ea)), seventy-six
  seconds after #238 and in that order on the seat's own recommendation. Two
  Decisions.

Delivered in full, in one QA round. M11-R1, M11-R2 and M11-R3 were all verdicted
`satisfied` in the
[qa comment](https://github.com/radiusred/gh-codecrew/issues/233#issuecomment-5547288500)
of 22:44:39Z. No verdict in this milestone came back other than `satisfied`.

What the hub looks like afterwards: `README.md` is 133 lines where it was 226,
with three headings (`#the-routing-table`, `#start-now`, `#read-next`) and two
links to codecrew.works; `docs/introduction.md` opens by calling the site the
introduction and the README the technical entry point, and carries a paragraph
of gloss on the routing table with links to the two worked examples rather than
a third copy of the block; `commitlint.config.mjs` is gone from both
repositories and both CI jobs call
`radiusred/.github/.github/actions/commitlint@main` behind a job still named
`Lint commit messages`.

## Decisions

Thirteen Decision comments. Four are the operator's — one on the milestone issue
and three on task #235; the other nine were written by the seats doing the work.
They group by the question each task had to settle.

### The receipts move to the site, and the README keeps a pointer

The README's `## The receipts` section was the repository's evidence that it is
agent-staffed, and the home page's `## CodeCrew Works` section carries four of
its five items in a richer form. The seat's
[Decision](https://github.com/radiusred/gh-codecrew/issues/235#issuecomment-5546469553)
cuts the section: "Keeping the README's bullets is precisely the '2 slightly
different versions of the same thing' the operator declined to maintain on
[codecrew-www#18](https://github.com/radiusred/codecrew-www/issues/18)."

**Trade-off, recorded in full:** "a GitHub visitor who never clicks through to
codecrew.works no longer meets the evidence that this repository is
agent-staffed — the single most persuasive thing the README had, and the reason
the section existed." Mitigated by two links, "not eliminated". **Rejected:**
keeping a trimmed receipts section, "the duplication is the thing being removed,
and a trimmed copy drifts faster than a full one"; and dropping the receipts
with no pointer at all.

The same Decision names a second consequence and welcomes it: the
doc-synthesizer contract's obligation to keep the README's proof points true at
every milestone boundary "loses its stalest instance — the hard-coded milestone
count, refreshed by hand in the M9 and M10 record PRs". The home page's
equivalent card says "Every milestone shipped this way" and needs no count.

Two of the five receipts have no counterpart on the home page — the stranger run
on [davison/numberguess](https://github.com/davison/numberguess), and the "not
yet: a per-platform onboarding script" caveat
([#54](https://github.com/radiusred/gh-codecrew/issues/54)) — and were captured
during the task as
[codecrew-www#24](https://github.com/radiusred/codecrew-www/issues/24).

### The README is the technical entry point, and its reference links stay at source

The milestone's largest change of direction was the operator's, made while
reviewing PR #238 in flight and recorded verbatim in two Decisions on #235. The
[first](https://github.com/radiusred/gh-codecrew/issues/235#issuecomment-5547013443)
opens "ok, my fault for a poor brief" and sets the shape: keep the
codecrew.works link near the top "but feel free to refer to it as the
marketing/intro site and note that this page (README) is the developer landing
page. Links in the README to docs should remain on github - link to the docs at
source. The README doesn't need to tell a story or sell anything, it can be as
dry and technical as required." The
[second](https://github.com/radiusred/gh-codecrew/issues/235#issuecomment-5547029878),
two minutes later, refines the audience: "use language that informs agentic
coders looking for more detail :) They may not be traditional developers".

Both edited M11-R1 on #233, which quotes the original clause in its Amendments
section. What the amendment cost the branch is visible in the two follow-up
commits and the seat's
[comment](https://github.com/radiusred/gh-codecrew/pull/238#issuecomment-5547046344)
on them: `f9c6a6b` moved every reference link back to source in this repository
and removed the blockquote pitch, the centred strapline and the narrative
framing, adding the refusal codes each verb raises (`NO_PLAN`, `GH_TOO_OLD`,
`NO_CHECKS`); `bbe062f` removed the phrase "developer landing page" entirely —
"the word appears nowhere in the file" — and pointed Read next at `AGENTS.md`
and `roles/` for an agent dispatched into the repository. The three headings
were unchanged across both commits, "so the anchor list on #235 still holds and
the mapping in codecrew-www#25 stays valid".

The operator's
[third Decision](https://github.com/radiusred/gh-codecrew/issues/235#issuecomment-5547148828)
answers the reviewer's second finding, which had asked for the stranger receipt
to stay until the site carried it. It does not stay: "The operator approved the
page at bbe062f without it, after amending M11-R1 to say the README tells no
story and sells nothing; the receipts' new home is the site, and codecrew-www#24
… carries the two receipts the home page lacks, this one included. The README
does not keep a receipt while the site catches up." **Trade-off:** "until www#24
lands there is no public page showing the stranger run; the record on #235 and
the old README in history keep the evidence." **Rejected:** keeping a receipts
section until the site catches up, "it would reintroduce exactly the second copy
the milestone removes".

### The routing-table anchor moves, and the spoke repairs the links behind it

`### 2. Four seats, always staffed` became `## The routing table`, so
`#2-four-seats-always-staffed` became `#the-routing-table`. The
[Decision](https://github.com/radiusred/gh-codecrew/issues/235#issuecomment-5546469707)
is about the numbering rather than the words: "a section number with no 1, 3 or
4 around it is a dangling reference, and the ordinal would have to be maintained
against a structure that no longer exists." The YAML block inside it is
byte-for-byte what is on `main` — [#229](https://github.com/radiusred/gh-codecrew/issues/229)
settled in M10 that it is a worked example and not a mirror, so nothing keeps it
in step with `.codecrew.yml` by force.

**Trade-off:** the live home page's crew gloss and `sync_docs.py`'s
`HOME_ANCHORS` both named the old anchor and went stale on merge. "GitHub does
not error on an unknown fragment — the link lands at the top of the README — so
the failure is a degraded link, not a broken one, and codecrew-www#22 is in the
same milestone to repair it." **Rejected:** keeping the heading text verbatim to
preserve the anchor, and adding an invisible HTML anchor alias — "a second name
for one section, the same objection recorded against an alias id on the home
page in the M9 record".

The seat posted the
[full anchor list](https://github.com/radiusred/gh-codecrew/issues/235#issuecomment-5546473495)
on #235 — every heading removed or renamed, with the anchor GitHub generates and
where the content now lives — and linked it from codecrew-www#22 before the hub
PR merged, saying so explicitly: "The hub PR is not merged yet; the anchors above
are what it lands with."

The spoke's
[Decision](https://github.com/radiusred/codecrew-www/issues/22#issuecomment-5546565671)
then chose what the table should say: a row per heading the README still has,
and no rows for the ones the cut removed. A hub page still linking a removed
anchor "gets the documented behaviour — anchor dropped, link lands on the home
page, a `note:` in the sync output — rather than a silent redirect to the section
the content moved to." **Rejected:** keeping the removed anchors' rows as safe
defaults, because "the table then describes a README that no longer exists — the
staleness codecrew-www#8 names — and quietly repairs upstream links that are
broken on GitHub".

That Decision also priced the merge window, having measured it: with the hub's
`main` still before #238, the recipe "builds strict-clean with exactly that one
note", the hub introduction's `README.md#the-receipts` link landing on the home
page root instead of `#codecrew-works`. "Merge this after #238, or accept one
deploy's degraded link." It was merged after #238.

### The introduction links the example rather than carrying a third copy

M11-R2 asked for "a routing-table example — carrying the README's block or
linking to it and to the home page's", and the
[Decision](https://github.com/radiusred/gh-codecrew/issues/235#issuecomment-5546469847)
takes the linking half deliberately. `docs/introduction.md` gains a paragraph
saying what the table is — a row per seat with its identity, harness and model,
`~` for a human — then links this repository's own example in the README and the
annotated generic one on the home page, with SPEC §5 as the field reference. The
reasoning is the milestone's own: "the milestone's purpose is to stop
maintaining the same content twice, and a third copy of the same eighteen lines
of YAML — with no drift test left to keep any of them honest since #229 — would
be a third thing to update when the hub reroutes a seat."

**Trade-off:** "a reader on codecrew.works/docs/ still does not *see* a routing
table on the page they are on; they see one sentence saying what it is and two
links." **Rejected:** carrying the block (three copies, no test); linking only
the README ("the site reader should not be sent to GitHub for something the site
already renders one click away").

### The unreferenced diagrams stay, against a precedent that points at deletion

Cutting the four beats left `hub-and-spokes.svg`, `milestone-lifecycle.svg` and
`gate-refusing.svg` referenced by nothing in the repository; `four-seats.svg` is
kept above the routing table. The
[Decision](https://github.com/radiusred/gh-codecrew/issues/235#issuecomment-5546469972)
names the precedent against itself: "#214 deleted the ten crew badge PNGs on
exactly this reasoning — nothing referenced them — so the precedent points at
deletion. This PR does not follow it, because those PNGs had a known shipped home
in codecrew-www while these three may be hotlinked by raw.githubusercontent URLs
from pages this repository cannot see." **Trade-off:** "three unreferenced files
stay in the tree, which is the untidiness #214 was opened to clear."
**Rejected:** deleting them here — "a link check across two other properties,
unrelated to the README's prose, inside a docs-cut PR".

### The commit lint is called, not carried

Both halves of M11-R3 kept the check context the org ruleset `require-lint`
matches and moved everything behind it into the shared action. The
[hub Decision](https://github.com/radiusred/gh-codecrew/issues/236#issuecomment-5546437853)
is about the difference between a workflow name and a check context: the
workflow keeps `name: CI` because it also carries the Go job, while "the check
context the org ruleset `require-lint` matches is the *job* name, which stays
exactly `Lint commit messages`, not the workflow name". **Trade-off:** "the file
is no longer a byte-for-byte copy of the reference caller". **Rejected:**
renaming the workflow, which "would mislabel the Go job's runs in the Actions
tab"; splitting the Go job out, as outside the task; and job-level permissions
"as a departure from the documented caller for no security gain — both grants
are read-only".

The
[spoke Decision](https://github.com/radiusred/codecrew-www/issues/21#issuecomment-5546448356)
made the mirror-image call on the same trade: the caller stays a job inside the
existing `ci.yml` rather than becoming the separate `commitlint.yml` the org
guide's copy-paste block names, and there `pull-requests: read` *is* scoped to
the job, so "the Tests job gains no permission it does not use". **Rejected:** a
second workflow triggered on the same event "for one repo of this size".

The hub's
[second Decision](https://github.com/radiusred/gh-codecrew/issues/236#issuecomment-5546437960)
leaves the two mentions of the deleted config under `docs/milestones/` as
written — the M1 and M3 records describe what those milestones shipped and "stay
true as history" — after confirming no live guidance names `wagoid` or the local
config. **Rejected:** editing the records to say "since removed", because
"records describe the milestone as closed, and rewriting them for later events is
the pattern M9/M10 corrections avoided".

Both PRs were cited on `radiusred/ops#17`, where #225 came from, at 21:07:45Z
and 21:08:22Z — before either had merged, each citation naming the run and the
check that proved the migrated context.

### codecrew-www#8 is unchanged, and stays open

The spoke task was asked to read
[codecrew-www#8](https://github.com/radiusred/codecrew-www/issues/8) — the
capture that `HOME_ANCHORS` can go semantically stale without breaking the build
— and say whether it changed anything.
[It did not](https://github.com/radiusred/codecrew-www/issues/22#issuecomment-5546565847).
The new site-build assertion checks the invariant rather than a fixed fragment:
every README anchor that actually crosses over into the synced introduction must
be a table value *and* an id the built home page has, which "turns a rename of
`#the-crew` or `#start-now` into a red test only while the hub introduction links
those anchors". The general check #8 asks for — every non-empty table value
against the built home page's ids — is still that capture's. **Trade-off,
recorded by the seat:** the assertion "passes vacuously on the fragments if a hub
checkout links no README anchor at all (it still requires at least one home-page
link)".

### The ROADMAP row, discarded again

The operator's one Decision on the milestone issue is the counterpart of M10's:
the `Open` row the installed v1.1.0 binary wrote locally is discarded, and this
document's PR adds the row already Done. **Trade-off:** "`ROADMAP.md` on `main`
shows M11 only once its record merges; `gh codecrew status` shows it now."
**Rejected:** carrying the row, since "the convention it belongs to is gone from
`main`'s source and survives only in the unreleased binary gap".

## Deviations

**No Deviation comment exists anywhere in this milestone.** The single comment
whose text matches the word is the introduction Decision on #235, which closes
by recording the absence of one: "**Deviation from the plan on this issue:**
none — this is the option the plan named, recorded here with its reasoning at
the point it was written."

PR #238 and PR codecrew-www#25 both state in their descriptions that the plan
was followed as written. PR #237 and PR codecrew-www#23 make no such statement,
and none is recorded on their task issues either.

The one change of direction the milestone did take — the operator's tone pass —
is not a deviation from a plan by a seat, and was not recorded as one. It amended
the requirement instead, which is why it appears in #233's Amendments section and
in three Decisions rather than here.

## The gates

**#233's Gates section was left as the scaffold's placeholder** — the template
line "_What "done" means beyond CI: e2e suites, manual UAT, sign-offs._",
unedited. [M10's record](10-protocol-bookkeeping-from-the-field.md#the-gates)
says the same of #207, [M9's](9-the-docs-at-codecrew-works.md#the-gates) of #201
and [M8's](8-a-product-home-page-for-codecrew-works.md#the-gates) of #196. No
`**Gate raised:**` or `**Gate resolved:**` comment exists anywhere in this
milestone, no `cc:needs-decision` label was applied, and `checkpoint` was not
used — on a task or on the milestone issue.

Both task Plans that carried an Ask-the-human section wrote `None` into it and
gave reasons: #235's says the two judgement calls its brief flagged "are recorded
as Decisions with trade-offs rather than escalated", because the operator's own
words on codecrew-www#18 "settle the direction for both"; #236's says "the check
context, permissions and config location are all fixed by radiusred/.github's
CONTRIBUTING.md". codecrew-www#21's says the task "fixes every judgment call",
and codecrew-www#22's notes the merge order "for the reviewer, not gated".

What gated the work: CI on every PR — `Lint commit messages` and `Go build and
test` in the hub, `Lint commit messages` and `Tests` in the spoke — an
independent approval on each of the four PRs, and the QA verdicts at the close.
The one course correction the milestone needed did not come through any of them:
the operator read the PR in flight, and the correction arrived as an amendment
and three Decisions.

## The review rounds

Four PRs, **six review submissions** — five by the reviewer role holder
(`radiusred-checky[bot]`) and one by the operator — and **one change request**,
resolved in the next round.

- **PR [#237](https://github.com/radiusred/gh-codecrew/pull/237)** — approved
  first round, "I found no findings". The reviewer verified the caller, the
  documented read permissions, the preserved check context, the removed local
  config and the untouched Go job, and that "the PR's own check passed under
  that exact context and its log shows the shared action ran".
- **PR [codecrew-www#23](https://github.com/radiusred/codecrew-www/pull/23)** —
  approved first round, no findings, on the same substance plus the spoke's
  suite (`64 passed`). The implementer's evidence comment on that PR quotes the
  job log resolving `radiusred/.github@main` to `f9926ef…`, staging the config
  and reporting `Lint free! 🎉`, and the repository ruleset requiring exactly one
  context.
- **PR [codecrew-www#25](https://github.com/radiusred/codecrew-www/pull/25)** —
  approved first round, "no correctness findings", after running the full deploy
  recipe with the requested environment and probing the rewrite directly against
  the #238 introduction shape. The review confirms the merge-order caveat rather
  than waiving it: "PR #238 is still open, so the merge-order caveat in this PR
  body is real; this review is approving the spoke change as sound."
- **PR [#238](https://github.com/radiusred/gh-codecrew/pull/238)** — the
  milestone's only change request and its only multi-round review, described
  below.

**PR #238's chronology, from its timeline.** Opened 21:10:31Z. Force-pushed at
21:15:10Z for the rebase onto #237's merge — the only conflict was `CHANGELOG.md`
and both `Unreleased` sections were kept. The operator's tone pass was
transcribed onto #235 at 22:10:09Z and 22:12:11Z, one second after each of the
two edits to #233 that amended M11-R1. The two revision commits followed at
22:12:12Z and 22:13:11Z, with no force-push, and the seat's summary of them at
22:14:11Z. The operator submitted an **approving review with an empty body** at
22:20:43Z — the approval the later Decision refers to as "approved the page at
bbe062f". The reviewer requested changes at 22:25:48Z. Pushing the fix commit
`b2257355` at 22:27:39Z **dismissed the operator's approval**, GitHub recording
`b2257355` as the dismissal commit and no dismissal message. The seat answered at
22:28:38Z, the reviewer approved at 22:34:03Z, and the PR merged at 22:34:52Z.

**The change request's two findings.** The first was a documentation defect the
cut created: `docs/introduction.md:13` still narrated the README's old onboarding
block — "the four lines under `[Start now](../README.md#start-now)`" and "start
your agent, 'Let's build this project!'" — while the new block "is three shell
lines … and no longer starts the agent". The seat fixed it and two neighbouring
sentences for the same reason, then grepped both files for narration of removed
content ("four lines", "Let's build", "why you'd", "four beats", "the receipts",
"developer") and reported the only surviving hit as the intended home-page
receipts pointer.

The second was the dropped stranger receipt. The reviewer had checked the live
site — "the home page receipts cover the milestone record, this repo being
agent-staffed, `radiusred/www`, and the Paperclip runs, but not the
stranger/scaffold-alone receipt" — and named the standard it was applying: "The
task brief for review says to judge removed content too: anything the site does
not carry that the README dropped is a finding." It was not actioned, by the
operator's Decision. The seat's answer says so plainly and separates the two
things: "Your reading of the evidence was right; the call on what to do about it
is the operator's, and it is recorded."

The reviewer also recorded, on PR #238, the check that the protocol's non-doer
gate depends on: "Reviewer identity verified: config App ID `4719924` matches
`gh api /apps/radiusred-checky --jq .id`, and differs from PR author
`app/radiusred-wordy`."

## QA: one round, three verdicts

The qa role holder (`radiusred-testy[bot]`) verdicted all three requirements in a
[single comment](https://github.com/radiusred/gh-codecrew/issues/233#issuecomment-5547288500)
at 22:44:39Z, against merged `main` in both repositories, having re-run the gates
with repo-local caches: the hub's `go build` and `go test ./...`, and the spoke's
full recipe (`64 passed`, ten files synced, `zensical build --clean --strict`
finding no issues).

For each requirement the seat states what the shipped tests can and cannot prove
before going past them. For M11-R1: "The shipped tests prove the spoke
rewrite/build path for current README anchors and assume the source docs corpus
they are run against; I went past that by fetching GitHub and the live site." For
M11-R3: "the hub side has no dedicated workflow-shape test, so I verified the
YAML and GitHub API directly."

**M11-R1** was verdicted against the rendered README, the source tree and the
deployed site together: "a 133-line technical entry point"; two codecrew.works
links rather than three "because the operator amendment keeps docs links at
GitHub source"; `#the-routing-table`, `#start-now` and `#read-next` present in
the rendered page; the only live inbound README links under `docs/` being
`../README.md#start-now` and `../README.md#the-routing-table`, both existing,
with the broad grep hit for `README.md#the-receipts` identified as "historical
inline code in an old milestone record, not a link"; `HOME_ANCHORS` carrying "no
stale key for a removed README heading"; and the live `/docs/` page rewriting to
`../#the-crew` and `../#start-now`, "and both target IDs exist on the live home
page". The verdict also checks the receipts Decision's own consequences:
`codecrew.works/#codecrew-works` exists, #224 is closed with back-references, and
"www#8 remains open exactly as the Decision on codecrew-www#22 says".

**M11-R2** was verdicted on content the tests do not reach: the introduction's
gloss "saying the table is one row per seat with identity, harness/model, and `~`
for a human", both example links resolving on GitHub and on the live site, and
#227 closed and adopted "with the linking-not-copying Decision recorded".

**M11-R3** was verdicted from the YAML and the GitHub API: both workflows calling
`radiusred/.github/.github/actions/commitlint@main` after a depth-0 checkout with
the job name kept exactly, "recursive file probes found no local
`commitlint.config.*` or `.commitlintrc*` in either repo", and four heads — the
two migration PRs and the latest merged PR in each repository — "all report
successful check runs named exactly `Lint commit messages`".

## Requirement outcomes

The status column is the verdict word as the qa role holder wrote it, verbatim
and unqualified. Everything else the milestone knows about a requirement is in
the columns around it.

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M11-R1 — the hub README is the developer landing page: dry and technical, written to inform agentic coders looking for more detail, with the site link near the top, the routing table, the install line and first verbs, and documentation links at source on GitHub; every inbound README anchor resolves after the cut, with the spoke updated in this milestone and its `--strict` build clean (adopts [#224](https://github.com/radiusred/gh-codecrew/issues/224); **amended** 2026-09-04, original clause quoted on #233) | [#235](https://github.com/radiusred/gh-codecrew/issues/235) / PR [#238](https://github.com/radiusred/gh-codecrew/pull/238); spoke half [codecrew-www#22](https://github.com/radiusred/codecrew-www/issues/22) / PR [codecrew-www#25](https://github.com/radiusred/codecrew-www/pull/25) | `satisfied` | Tasks closed; the one QA round, against the rendered README, the source tree, `HOME_ANCHORS` and a live fetch of the deployed site. The dropped stranger receipt is accepted by operator Decision, with [codecrew-www#24](https://github.com/radiusred/codecrew-www/issues/24) open for it |
| M11-R2 — `docs/introduction.md` gives its reader a routing-table example — carrying the README's block or linking to it and to the home page's — and says what the table means in one line; decided together with the README cut and recorded (adopts [#227](https://github.com/radiusred/gh-codecrew/issues/227)) | [#235](https://github.com/radiusred/gh-codecrew/issues/235) / PR [#238](https://github.com/radiusred/gh-codecrew/pull/238) | `satisfied` | Task closed; the one QA round, reading the gloss and both example links on GitHub and on the live site. The linking half was taken deliberately and recorded with its trade-off |
| M11-R3 — gh-codecrew and codecrew-www call the shared commitlint action published by radiusred/.github instead of carrying their own workflow steps and config; the check context stays exactly "Lint commit messages", the local config files are deleted, and the PRs are cited on `radiusred/ops#17` (adopts [#225](https://github.com/radiusred/gh-codecrew/issues/225)) | [#236](https://github.com/radiusred/gh-codecrew/issues/236) / PR [#237](https://github.com/radiusred/gh-codecrew/pull/237); spoke half [codecrew-www#21](https://github.com/radiusred/codecrew-www/issues/21) / PR [codecrew-www#23](https://github.com/radiusred/codecrew-www/pull/23) | `satisfied` | Tasks closed; the one QA round, from the YAML, recursive config probes and check-run records on four heads. Both citations were posted on `radiusred/ops#17` before either PR merged |

No requirement was dropped and none was added. One was amended after the
milestone opened, twice by the operator on 2026-09-04, with the original clause
quoted on #233.

**Captures adopted and closed by the merges — three:**
[#224](https://github.com/radiusred/gh-codecrew/issues/224) (cut the README down
and send readers to codecrew.works),
[#227](https://github.com/radiusred/gh-codecrew/issues/227)
(`docs/introduction.md` could carry the hub's routing table) and
[#225](https://github.com/radiusred/gh-codecrew/issues/225) (adopt the shared
commitlint action). **Captures filed during the milestone — two, both open:**
[#234](https://github.com/radiusred/gh-codecrew/issues/234) (`task new` refuses
`NOT_FOUND` for a milestone created seconds earlier), opened by the operator at
21:02:01Z; and
[codecrew-www#24](https://github.com/radiusred/codecrew-www/issues/24) (the two
receipts the README cut leaves without a home on the site), opened by the
doc-synthesizer seat at 21:09:18Z from the receipts Decision.
[codecrew-www#8](https://github.com/radiusred/codecrew-www/issues/8) was read by
codecrew-www#22 and left open by Decision.

## Protocol-discipline observations

Six things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **A requirement was amended while its PR was in review, and the verdict came
  well after.** [M10's record](10-protocol-bookkeeping-from-the-field.md#protocol-discipline-observations)
  observed that "a requirement can be amended between a QA run and its verdict",
  its amendment having landed two minutes and three seconds before the round-one
  verdict. Here the amendment edits landed at 22:10:08Z and 22:12:09Z, the PR
  merged at 22:34:52Z and the single verdict was posted at 22:44:39Z — thirty-two
  minutes after the second edit, against the amended clause and against merged
  `main`. Nothing in the protocol produced that ordering; the qa dispatch simply
  came after the merges, as it did not in M10.
- **The milestone's course correction arrived through no protocol mechanism at
  all.** Every task Plan declared no ask-the-human points, and each gave reasons.
  The operator read PR #238 in flight, disagreed with the tone, and the
  correction reached the record as an amendment to #233 plus three Decisions on
  #235 — no gate, no `checkpoint`, no `cc:needs-decision`. The protocol's
  escalation path is designed for a seat that knows it is stuck; this was an
  operator noticing something a seat believed it had settled.
- **The operator's review is the one intervention with nothing in it.** The
  approving review on PR #238 carries an empty body, and the next push dismissed
  it. A reader who opens PR #238 alone sees a dismissed approval, a change
  request and an approval; the reasoning that shaped the merged README lives one
  hop away on #235 and #233. That the words survive at all is because the seat
  and the operator transcribed them verbatim into Decisions and the Amendments
  section — a convention, not a mechanism.
- **The first `task new` calls of the milestone were refused, and the capture is
  numbered before three of the four tasks.** #234 records three
  `task new --milestone 11` refusals of `NOT_FOUND` seconds after `milestone new`
  created #233, with the fourth call succeeding; the issue numbers bear it out —
  the only task created before the capture is the spoke's, codecrew-www#21 at
  21:01:19Z, four seconds after the milestone issue, while #235, #236 and
  codecrew-www#22 were all created after #234 at 21:02:01Z. The cause the capture
  gives is the one M10-R2 closed for `milestone new`'s own number check: "`task
  new` still trusts a single listing read." M10 fixed the collision at one end of
  the same eventual-consistency window and left the other end open; it was hit
  while opening the very next milestone.
- **The doc-synthesizer seat held a delivery task again, and nothing recorded the
  choice.** #235 was started by `radiusred-wordy[bot]` under the implementer
  contract; #236 and both spoke tasks by `radiusred-cody[bot]`. M10's record made
  the same observation about three of its ten tasks and noted that nothing in the
  contracts says whether a role's holder may be routed to another role's task.
  Nothing recorded it this time either. The non-doer review gate is unaffected —
  `radiusred-checky[bot]` reviewed all four PRs, and verified its own distinctness
  from the author on the one it changed.
- **The refresh obligation found nothing to fix, and the tasks are why.** The
  doc-synthesizer must keep `README.md` and `docs/introduction.md` true at every
  milestone boundary. The sweep for this document found no stale claim: the
  milestone shipped no Go, no release and no new refusal code, so the README's
  "thirty-one refusal codes" and the introduction's "**Shipped:** v1.1.0" and
  "All thirty-one" still hold; the introduction's own description of the README
  was corrected inside PR #238, by the seat that changed the README, and the
  milestone count that M9's and M10's record PRs each refreshed by hand no longer
  exists to refresh. The obligation is smaller at this boundary than at the last
  two because a task deliberately removed the thing that kept making it work.

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **The operator's tone-pass words exist on the trail only as quotations.** They
  were said outside GitHub and transcribed — by the operator into two Decisions on
  #235 and into #233's Amendments section — and the operator's own review on PR
  #238 has an empty body. What is recorded is when they were transcribed
  (22:10:09Z, 22:12:11Z) and when #233 was edited (22:10:08Z, 22:12:09Z).
  *Recorded as quoted, not as independently sourced:* nothing on the trail shows
  the original channel, the exact wording before transcription, or the moment the
  operator formed the view.
- **#234's own timestamp does not match the issues it describes.** The capture
  cites "this session's shell output, 2026-09-04 ~19:55 BST", while #233 was
  created at 21:01:15Z — 22:01 BST — and #234 itself at 21:02:01Z. *Recorded as
  written; nothing on the trail reconciles the two.* No shell log is attached, so
  the three refusals rest on the capture's account, and the fourth call's success
  (codecrew-www#21, 21:01:19Z) is the only corroborating artefact in either
  direction.
- **Whether the three unreferenced SVGs are hotlinked from outside this
  repository was never checked.** The Decision keeps them on the possibility that
  raw URLs elsewhere point at them; no check was run, no capture holds the
  question, and #214's precedent points the other way. The file that resolves it
  is a link sweep across codecrew-www and the announcement blog post, and nothing
  owns it.
- **Neither commitlint task weighed the moving action ref.** Both callers
  reference `radiusred/.github/.github/actions/commitlint@main`, which is what the
  organisation's documented caller specifies. The evidence comments record the SHA
  it resolved to on the day (`f9926ef…`), but no Decision anywhere considers
  pinning it, and no capture holds the question. *Absent from the record, not
  decided against.*
- **codecrew-www#24 is filed, open, and unadopted.** The two receipts the home
  page lacks — the stranger run and the onboarding-script caveat — have a capture
  and no task, so the trade-off the operator's Decision records stands until a
  later milestone adopts it: no public page shows the stranger run.
- **No gate was declared or raised.** #233's Gates section is the scaffold's
  placeholder text, unedited. `checkpoint` was not used in this milestone, on a
  task or on the milestone issue — the same absence M10's record notes for the
  requirement that taught the verb to accept a milestone ref.
