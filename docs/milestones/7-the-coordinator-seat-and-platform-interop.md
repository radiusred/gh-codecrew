# M7: The coordinator seat and platform interop

Tracking issue: [#163](https://github.com/radiusred/gh-codecrew/issues/163) ·
Synthesized 2026-08-31 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #163's requirements — including **M7-R4 as amended by
the operator on 2026-08-30**, with the original clause quoted inside the
amended text — its Gates section and its single QA round; the ten hub task
issues and their merged PRs (#164/PR #188, #168/PR #169, #170/PR #171,
#173/PR #174, #175/PR #176, #178/PR #179, #180/PR #181, #182/PR #185,
#183/PR #184, #186/PR #187) and the three in
[radiusred/www](https://github.com/radiusred/www) (#51/PR #52, #53/PR #54,
#56/PR #57); the eight Decision and two Deviation comments across them; the
cycle-4 findings log on
[#164](https://github.com/radiusred/gh-codecrew/issues/164) — seven log
comments, two of them numbered entries, carrying findings 51–68, two milestone
tables and a fold-back map, this milestone's largest single record;
[#119](https://github.com/radiusred/gh-codecrew/issues/119) for the three
cycles before it; the v1.1.0 release and
[CHANGELOG.md](../../CHANGELOG.md); and `docs/platform-interop.md` as it
shipped. The M1–M6 documents supplied the house form.

The gather from `gh codecrew milestone close 7` is again not part of the raw
material, for the reason M6's document predicted and one more: the live close
refuses `DOC_MISSING` until this file is on `main`, and the `--dry-run` added
by this milestone stops earlier still — at `OPEN_TASKS`, naming
[#189](https://github.com/radiusred/gh-codecrew/issues/189), the task that
writes this document. Every record below was read from its issue directly.

## Goal and outcome

M6 ended with an experiment and a gap. The orchestrator run
([#119](https://github.com/radiusred/gh-codecrew/issues/119)) had proved a
platform can drive the protocol, and had itemised what the framework lacked
for it: the coordination layer had no contract, no identity verb, and no
artifact a platform could install; credentials and webhooks were hand-wired;
the interop document ([#54](https://github.com/radiusred/gh-codecrew/issues/54))
was still unwritten. M7 closes that list — a coordinator contract and identity,
"the verbs the run's findings asked for", a fourth cycle on a fresh proving
ground with its findings folded back, the ladder's last rung written from them,
and v1.1.0 ([#163](https://github.com/radiusred/gh-codecrew/issues/163)).

Nine requirements, M7-R1 through M7-R9, all present when the milestone opened
at 16:53Z on 2026-08-29 — none added later, which is the difference from M6.
Thirteen tasks delivered the work — ten in the hub, three in the www spoke —
and this document is the fourteenth. Five captures the run had filed were
adopted and closed inside it (#133, #138, #157, #158, #159); a sixth, #132,
was adopted as M7-R2 and shipped but is still open, which the observations
below name. Two more captures were opened by cycle 4 and closed in the same
milestone (#165, #172), and
[#54](https://github.com/radiusred/gh-codecrew/issues/54) — open since
2026-08-22 — was finally written.

Delivered in full. All nine requirements carry a `satisfied` verdict from the
qa role holder in a
[single round](https://github.com/radiusred/gh-codecrew/issues/163#issuecomment-5471956921)
on 2026-08-30, against merged `main` and the shipped v1.1.0 binary. One of
them, M7-R4, was satisfied against text the operator **rewrote mid-milestone**;
that decision is the first section below. And one of #163's gates does not
close by evidence at all — it closes by this document saying which parts of
the coordinator contract cycle 4 did not exercise. That section is
["What cycle 4 did not exercise"](#what-cycle-4-did-not-exercise-of-the-coordinator-contract).

## Decisions

### The platform overlay was dropped, and a blank file took its place

The milestone's largest decision, and the only one that changed a requirement.

M7-R4 as written asked for **platform onboarding from the CLI**: `init
--platform <name>` and a `roles overlay <name>` verb, scaffolding
`roles/<role>.local.md` "from a shipped template carrying the platform
paragraph (wake syntax, landed-means-done, one wake path, no comments on done
tickets, credentials, tooling)". Cycle 4 had just paid for the absence of it —
the seats ran snake on numberguess's composed bundles because snake had no
extensions at all
([#164 finding 57](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218)).

The operator amended it on 2026-08-30
([Decision](https://github.com/radiusred/gh-codecrew/issues/163#issuecomment-5468874006)):
**no `init --platform`, no `roles overlay` verb, no shipped platform
template.** `init` scaffolds a *blank* `roles/<role>.local.md` beside every
contract it writes, holding one stable comment — what the file is (the
project's extension, loaded after the contract, append-only, composed by
`roles show`), a pointer to SPEC §7, and a pointer to an upstream examples
page, `docs/extensions.md`, "where uses — house style, CI conventions, a
platform's wake syntax, ids and tooling — live and can change without touching
anyone's scaffold". A comment-only extension composes to nothing.

**Why**, in the operator's words: the mechanism is purpose-neutral — "this
hub's own `doc-synthesizer.local.md` is an editorial voice, not a platform" —
and the original clause conflated it with one body of content and with
shipping a third-party platform's quirks inside the binary. The blank file is
the routing-table move again: onboarding starts explicit, and the pointer
arrives at the moment the question does. Rerunning `init` in an existing hub
adds the files, so there is no upgrade verb either.

**Rejected**, quoted: "a generic `roles extend <role> [--from]` verb (a verb
where a file suffices); embedding a Paperclip template (opinionated, versioned
with the CLI, and the ids are per-installation anyway — it was never a
template)."

**Trade-off.** The amended shape ships nothing that gets a platform running:
the recipe lives in `docs/extensions.md` and the interop page, which the
operator or their agent must find and apply by hand. That cost is named on the
interop page's own gap list — no onboarding script or plugin ships — and it is
the deliberate price of not versioning another product's quirks inside this
one. #163's body carries the amended clause with the original quoted inside
it; PR [#174](https://github.com/radiusred/gh-codecrew/pull/174) implements it,
closing #159, #158 and #138 with it.

One thing left the requirement rather than entering it:
[finding 68](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465449989)
(the scaffold branch nothing sweeps) was split out to become
[#172](https://github.com/radiusred/gh-codecrew/issues/172), which is the next
decision below.

### The coordination layer becomes the routing table's fifth row

The seat had to exist somewhere before it could be dispatched. The question
the implementer raised on
[#168](https://github.com/radiusred/gh-codecrew/issues/168) was whether `init`
should write it for *every* adopter, making the routing table five rows for a
solo operator who is already the coordinator.

[Decision](https://github.com/radiusred/gh-codecrew/issues/168#issuecomment-5465531521):
yes — `init` scaffolds `coordinator: { identity: ~ }` as the fifth row, and
the docs count five roles for the table. **Why:** M7-R1 asks for the entry
with unrouted meaning the operator, and SPEC §5's rule is that onboarding
starts explicit; "a row the operator embodies costs nothing, while an absent
row makes every platform onboarding add it by hand (what cycle 3 and 4 did)".
**Rejected:** leaving the row out of `init` and documenting it as
platform-only — "the seat exists in the solo shape too (the operator *is* the
coordinator), so the table should say so."

The compatibility consequence is on the record too: `role coordinator` against
a 1.0 table without the row answers `~` rather than refusing, because the seat
predates the row (PR [#169](https://github.com/radiusred/gh-codecrew/pull/169)).
The docs distinguish the five-row routing table from the passages that count
four *crew* seats, which is why the README's "Four seats, always staffed" and
its four-App receipt were left standing.

`identity new coordinator` mints the permission set
[#119 finding 16](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5445972756)
specified — contents: read, issues: write, pull requests: read, metadata
implicit — and never contents: write; `--with-approval-permission` stays
reviewer-only. The contract states the same set in prose, and adds the reading
that makes it operable: "a 403 on a push or a review is the contract enforced"
([`roles/coordinator.md`](../../roles/coordinator.md)).

### The mint becomes a verb, and refuses rather than falling back

`identity token <slug>` (M7-R2, adopted from #132) exists because the run
watched three different agents write three different RS256 helpers, and
because one of them was wrong in a way nobody diagnosed: the coordinator's
brief read the installation-token response with `["access_token"]` where
GitHub returns `token`, every mint raised `KeyError`, and the agent "called
the failure flaky, and proposed a retry loop"
([#164 finding 56](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218)).
The companion finding is sharper still: a `gh-cli` skill that mints exactly
this token had been installed in the implementer's bundle since 2026-08-28 and
appears in none of its runs, because "the contract's Credentials bullet
already told them *how*, so the model never shopped for a tool" — and the
skill was Claude-only, useless to the two Codex seats
([finding 67](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465449989)).

Two questions were resolved on
[#170](https://github.com/radiusred/gh-codecrew/issues/170) before the first
commit. **Fall back to the operator's `gh` auth when nothing mints?** No —
M7-R2 says the verb refuses with a code, and "a silent fallback to another
principal's token is the misattribution SPEC §5 exists to prevent". **Retire
`scripts/codecrew-token` or wrap it?** Wrap — "the 1.0 contracts told adopters
to install it from a fixed upstream URL, so the file must keep resolving; the
wrapper is one line."

Four refusal codes came with it — `NO_CREDENTIALS`, `BAD_CREDENTIALS`,
`NO_INSTALLATION`, `INSTALLATION_AMBIGUOUS` — taking the catalogue from
twenty-three to twenty-seven. A hinted installation id is a preference and
never trusted blind: the verb checks it against `GET /app/installations` and
reports a stale one on stderr,
[#119 finding 35](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5453727497)
made executable. The stderr receipt line — slug, App id, installation, account
— is the identity check the dispatch recipe wanted (#139), and stdout carries
the token and nothing else.

### Ownership is the seat, not the login

`task finish` refusing `NOT_OWNER` (#175, from the operator's capture #165) is
the milestone's one enforcement change, and it came from a merge that should
not have happened: snake PR #8 was Wordy's document, and `radiusred-cody[bot]`
merged it, "because Loopy's table sends approved → Cody, `task finish`,
regardless of whose PR it is, and the verb checks the reviewer gate and CI but
not that the finisher started the task"
([#164 finding 58](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5464097267)).
The operator's direction was four words: *"don't merge a task you don't own."*

The design question that followed was team-held seats, and it was decided in
the open
([Decision](https://github.com/radiusred/gh-codecrew/issues/175#issuecomment-5468991234)):
**ownership is the seat, not the login.** `task finish` accepts the same login
(`[bot]` suffix ignored) *or* the same routed seat, since a role held by a
GitHub team is held by any member. A starter who has since left the team
resolves to no seat and no longer matches; the task is handed over by the new
seat running `task start` again — which records a fresh `**Started by**`, and
the latest record wins — or an operator finishes it with `--bypass`, which
posts the override on the PR naming the owner it overrode. **Why:** "a
login-level check would have refused a teammate, which is not what the routing
table means by a team-held role." No new flag; `--operator-confirm` does not
waive ownership, and the operator's own `gh` auth is not exempt.

The twenty-eighth code. The consequence the PR states plainly: on this hub
nothing changes, because cody has started every task cody finishes; on a
crew driven by a platform, a wrongly dispatched implementer now gets a refusal
naming the owner instead of a merge
(PR [#176](https://github.com/radiusred/gh-codecrew/pull/176)).

### Plan first, then act — so a preview cannot disagree with the verb

M7-R5's `--dry-run` (from #133) could have been a second code path that prints
what the verb would do. It is not. Each of the three writing verbs was split
into a read-only *plan* that evaluates the gates in order and decides the
actions, and an *execute* step that performs them; the live verb runs both,
`--dry-run` runs the plan and prints it
([#178](https://github.com/radiusred/gh-codecrew/issues/178), PR
[#179](https://github.com/radiusred/gh-codecrew/pull/179)). "Gate logic is not
duplicated: the same functions produce the same refusals, so the dry run can
never disagree with the real run." A gate the live verb would not reach prints
as *not reached*; a path not taken as *not applicable*.

One question was open in #133 and answered on #178: **should a refused dry run
exit non-zero?** Yes — "M7-R5 says 'the same refusal codes'; a script asks a
dry run the same question it asks the verb." A clean run exits 0 and reports
`dry run: nothing written`. `milestone new --dry-run` prints the number the
milestone would get, which is the point: requirement prose can be written
knowing it, the M4 lesson where a closed duplicate shifted M3 to M4
([#119 finding 45](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462078800)).

The sequencing question #133 also raised — sweep before close, or close before
sweep — was answered by not changing it: the sweep only ever deletes a merged
or empty branch, so a close that fails after it loses nothing. That reasoning
is now a comment in the code rather than a change.

### An App's events are set at creation, and the verb says so

M7-R3 asked for `identity webhook` to set an App's hook URL, secret **and
event subscriptions**. The implementer checked what GitHub actually allows,
live, against this org's own Apps: `GET`/`PATCH /app/hook/config` read and
update the URL, secret and content type; there is **no endpoint for an App's
event subscriptions** — `PATCH /app` is 404, and they are set by the manifest
at creation or by hand on the App's settings page
([#180](https://github.com/radiusred/gh-codecrew/issues/180)).

So the verbs were shaped by the platform's real surface: events at creation
(`identity new --with-webhook --events`), URL and secret at any time, and
`--events` on an existing App **refused** with the settings-page URL and the
current list "rather than pretend". A second live check narrowed it further —
an App minted without a webhook has no hook configuration at all, `GET` and
`PATCH` both 404, and GitHub cannot create one — so every path on such an App
refuses `NO_WEBHOOK` naming the settings page where Webhook → Active is ticked
by hand
([comment on PR #181](https://github.com/radiusred/gh-codecrew/pull/181#issuecomment-5469241608)).
The twenty-ninth code.

The default event set changed in a minor release, and the ask-the-human point
on #180 says why that is allowed: "SPEC §10 promises flag names and refusal
codes stable, not the event set a manifest subscribes; the CHANGELOG says what
changed and `--events` restores any subscription." `--with-webhook` now
subscribes `pull_request` and `pull_request_review` — the two transitions the
seat routines wake on — where 1.0 subscribed five, and the three it dropped
(`issues`, `issue_comment`, `check_suite`) are the ones the run recorded
waking a coordinator for nothing
([#164 findings 46 and 53](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218)).

### `init` commits what it wrote — no `git init`, no stash

The operator proposed four steps on
[#172](https://github.com/radiusred/gh-codecrew/issues/172): check for a git
repository and `git init` if absent, stash uncommitted work, commit the
scaffold, unstash. The implementer's
[assessment](https://github.com/radiusred/gh-codecrew/issues/172#issuecomment-5468961697)
took two and rejected two, and the operator
[adopted it on that shape](https://github.com/radiusred/gh-codecrew/issues/172#issuecomment-5469484860).

**Rejected — `git init`:** "the case it 'fixes' is the one where the operator
ran `init` in the wrong directory (a home directory, a parent) — today that
mistake costs a few stray files and a note; with `git init` it costs a stray
repository around everything beneath it, and the scaffold is then committed
into it. Git being mandatory is an argument for refusing clearly, not for
creating." **Rejected — stash and unstash:** `init` creates files and never
edits one, so a pathspec commit (`git commit --only -- <the files init wrote>`)
already leaves the operator's staged and unstaged work exactly as it was;
"stash is the step that can fail in the middle (a pop that conflicts, a stash
left behind when the commit errors); not doing it is the safe version."
**Rejected — `task finish --bootstrap`:** "a milestone-less merge path the
protocol does not need once the commit is on the right branch."

**Trade-off**, stated in the adoption: on a fresh repository this is the root
commit and findings 51 and 68 both disappear; where a ruleset requires pull
requests the commit lands on `codecrew-bootstrap` and the one human merge
remains — recorded as the pre-milestone gate on that PR — with the branch
cleaned by delete-on-merge, "which `status` already nags about". The verb never
pushes. A commit that cannot be made is a note with the command to run by
hand, "never an error that hides the scaffold" (PR
[#184](https://github.com/radiusred/gh-codecrew/pull/184)).

### Two cost tables, and never a third

Twice — once for the interop page, once for the article — the doc-synthesizer
was asked implicitly to merge four cycles of cost into one table, and twice
declined on the record
([Decision on #182](https://github.com/radiusred/gh-codecrew/issues/182#issuecomment-5469539530);
[Decision on www#51](https://github.com/radiusred/www/issues/51#issuecomment-5470461265)).

**Why:** the two logs do not count the same thing. #119's column is *non-gate
touches*; #164 splits *onboarding* touches (bringing the repo, App, hooks,
ingress and Project up) from *workflow* touches (a human standing in for a
wake that should have fired). "Merging them would either silently equate two
different metrics or require me to re-derive numbers, and M7-R7's gate is that
every figure is reproduced exactly as recorded." **Trade-off:** a reader must
compare across two tables — mitigated by the "cycle 3, for scale" row #164's
own table carries, which is where the logs already join. **Rejected:** "a
third, synthesised table of 'all four cycles' with a normalised touch column.
It would be the only number on the page not traceable to a comment." The
article's version adds a second reason: cycle 1's coordinator share was never
measured, and "filling those with dashes reads as missing data; filling them
with estimates would be fabrication." The article states it as unmeasured.

### The announcement rules become standing guidelines

The operator's three findings on the announcement PR — plain speech for a
reader who knows neither product, no inline URLs in a LinkedIn post, and a
jargon purge — arrived with *"These are probably general guidelines that
should be maintained if possible."* They were written once, into
`announcements/README.md`, rather than re-derived per announcement
([Decision on www#53](https://github.com/radiusred/www/issues/53#issuecomment-5471799214)).

**Trade-off:** "a second place to keep current — the guidelines can drift from
what the announcements actually do. Accepted because the alternative is worse:
three announcements in, the rules were living in the operator's review
comments, which do not travel." **Rejected:** leaving them in this
announcement's own file, "invisible to the next one", and putting them in the
repo README, "which is about the tooling, not the editorial rules".

Decided in passing, and worth recording because it is the pattern: the
LinkedIn first comment became a verb — `social comment --urn … --text-file …`,
with tests, in the same PR — "because a guideline that depends on someone
remembering to paste a comment seconds after posting is a guideline that will
be missed."

## The gates

#163 declared seven gate conditions. What actually ran is worth separating
from what the section asked for, because only one of them travelled through
the protocol's own gate mechanism.

**The announcement gate (www#53) — the one `cc:needs-decision` of the
milestone.** M7-R8's hard clause: no social post without an operator
resolution. wordy drafted both texts, measured the tags, captured the
`check` and `--dry-run` output into the announcement file, and
[raised the gate at 19:22:53Z](https://github.com/radiusred/www/issues/53#issuecomment-5470754543).
The operator reviewed the PR at 22:47Z with three findings and *changes
requested*, the texts were rewritten, and the gate was
[resolved at 23:05:32Z](https://github.com/radiusred/www/issues/53#issuecomment-5471818828)
— *"post as approved in the PR"*. The posts went out at 23:06:41Z (Bluesky),
23:06:47Z (LinkedIn) and 23:06:52Z (the first comment), each verified at the
source rather than from the tool's own output, and the PR merged after the
Result section recorded the URLs. Three hours and forty-two minutes in which
nothing was published and nothing worked around the gate; the resolution
changed the content, which is what a gate is for.

**The article's tone pass (www#52), taken as a PR review rather than a
checkpoint.** wordy
[requested it explicitly](https://github.com/radiusred/www/pull/52#issuecomment-5470459335)
— "@davison as reviewer on the article's voice and framing (origin story, not
a takedown), and on the culture-clash section, which quotes him" — and stated
that the PR was not finished until both the operator and the reviewer seat had
passed. Four findings came back and all four were addressed at `573e02e`: the
opening no longer reads as a confession; "the platform's lead agent" became
"Paperclip's CEO agent" in all three places, because "'lead agent' could mean
anything in a Paperclip company"; every observed behaviour now names Paperclip
by name, while the claims that are about the separation of concerns rather
than one product stayed generic; and the prose was hard-wrapped at 80 columns.
The operator approved at 18:59:32Z; the reviewer seat approved four minutes
later. The same tone pass was requested on the interop page and given as an
approval on PR #185 at 15:57:31Z, which the fix push then dismissed
(`b36685a`) — the record shows it now only as a dismissed approval with an
empty body, recoverable from the dismissal event rather than from the review.

**The evidence gate.** `gh codecrew milestone evidence 7` was run live before
the QA dispatch and resolved all sixteen citations, "including
punctuation-bearing prose" — which is the M7-R4 extractor fix
([#138](https://github.com/radiusred/gh-codecrew/issues/138)) proving itself
on the milestone that shipped it. Recorded in the M7-R4 verdict.

**Release parity.** Verified twice: by the implementer, who downloaded and ran
the shipped binary after the tag
([release evidence](https://github.com/radiusred/gh-codecrew/issues/186#issuecomment-5471890289)),
and independently by qa, which fetched
`gh-codecrew_v1.1.0_linux-amd64` (SHA-256
`d98a2ede630a16581996d7d7f166d3d747afc2664a493fed3fb53a9c62f9f45c`), executed
it, and confirmed the release API lists all five platform assets.

**The M7-R4 amendment did not travel as a gate**, and that is a gap in the
record rather than a decision about gates. M6's equivalent — the M6-R1 avatar
clause — went up as a `**Gate raised:**` with two options enumerated, was
resolved in the same form, and the requirement was amended with the original
quoted inside it. Here the amendment carries the last of those three: #163's
body quotes the original clause under *"Original clause:"*, and the
[Decision comment](https://github.com/radiusred/gh-codecrew/issues/163#issuecomment-5468874006)
records the reasoning and the rejected alternatives in full. But no
`checkpoint` was raised, no `cc:needs-decision` label was ever applied to #163
(the label timeline shows only `cc:milestone`), and the question the operator
answered is not on the record — only the answer. The decision is well
documented; the deliberation is not.

**The coordinator-contract gate** is the one this document closes, below.

## Cycle 4: two milestones on radiusred/snake

M7-R6's deliverable was a findings log, not code, and it is the milestone's
largest record: six entries on
[#164](https://github.com/radiusred/gh-codecrew/issues/164), findings 51
through 68 continuing #119's numbering so cross-references stay unique.

The shape: the operator created `radiusred/snake` at 15:57Z on 2026-08-29 with
a LICENSE and the org rulesets `protect-main` and `require-lint`, and kicked
off a coordinator agent — Loopy, on Sonnet 5 — with one line at 16:08:33Z. The
four crew seats ran on the bundles installed for numberguess cycle 3. **M1**
tested the coordinator shape: the coordinator's App hook → a webhook router
routine → dispatch by table. **M2** tested seat-routed routines: the two table
transitions wired straight to the seats (`pull_request` → the reviewer,
`pull_request_review` → the implementer), the coordinator unsubscribed from
both and kept for milestone verbs, gates and off-table events. Every operator
act to bring the repository up was counted as an **onboarding touch**,
separately from **workflow touches** — a human standing in for a wake that
should have fired.

Both milestones closed. `radiusred/snake` went from an empty repository to two
closed milestones in five hours of wall-clock, "with a game the operator
reports is fun to play"
([entry 2](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465449989)).
The numbers, reproduced as recorded:

| | runs | $ | coordinator share | PRs | runs / PR | wall-clock | onboarding touches | workflow touches | gates |
|---|---|---|---|---|---|---|---|---|---|
| M1, coordinator shape | 71 | 36.90 | Loopy 35 runs, 52% | 3 + scaffold | 24 | 1 h 43 | 4 | **0** | 0 (2 platform confirmations) |
| M2, seat routines | 33 | 17.40 | Loopy 8 runs, 23% | 4 | **8** | 1 h 40 (incl. ~30 min of stalls) | 4 | 3 | 0 |
| cycle 3 (numberguess M4), for scale | 97 | 82.76 | 75% | 4 | 24 | 3 h 08 | — | 0 | 1 |

The log's own reading: "the seat routines cut runs per PR by two thirds
against both M1 and cycle 3, and the coordinator's share of the bill from
three quarters (cycle 3) to a quarter; M2's three workflow touches were all
wakes after coordinator-brief defects (62, 63, 65), not platform faults, and
each became a change point." The cost column is the instance's reported figure
for the Claude-family seats; the Codex seats report tokens and no cost.

M1's zero workflow touches is the headline the M2 column then qualifies. M1
ran the full loop with no hand on it — including cycle 1's entry-9 sequence,
the one that once deadlocked overnight: PR approved → `task finish` →
`refused[NO_CHECKS]` → CI added → the push dismissed the approval →
`refused[NO_HOLDER_REVIEW]` → hand-back → re-review → approved, "run end to
end with no hand on it"
([entry 1](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218)).
It also cost 24 runs per PR to do it, and 52% of the bill was the
coordinator's.

Three findings are worth their own space because they changed the contract's
shape rather than a line of it.

**The dispatch template re-created the second wake path it was meant to
remove.** Every dispatch the coordinator wrote ended *"mark this ticket done
and mention me"*, so each transition woke it twice — once by webhook, once by
the mention — producing two dispatch tickets for one PR, twice. Single-flight
held and the seats absorbed the duplicates at a run each; three further runs
produced nothing but the platform's "Agent did not post a summary comment this
run". The act-time re-read did not catch it because it checked the PR's
reviews and not the coordinator's own open dispatches
([findings 53 and 54](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218)).

**Redundancy hid a four-hour outage.** The webhook router's trigger had been
dead since 14:23:45Z and M1 never noticed, because the "duplicate" mention
path was carrying every transition on its own. The cause was the ingress
Funnel's public leg closing TLS with no bytes while the node reported healthy.
"Two lessons: redundancy hid an outage for four hours (a single wake path
would have stalled visibly — the stall *is* the alarm), and the routine's
`lastFiredAt` is the one observable that says whether the seam is alive"
([finding 59](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5464097267)).
It dropped a second time at the start of M2; the operator re-registered it and
installed a watchdog.

**The coordinator did not know which project it was in.** Every run in cycle 4
opened with the platform's *"No project or prior session workspace was
available"*, so each wake reused whatever checkout it last had. After a merge,
the implementer mentioned the coordinator without naming a repository; the
coordinator woke in its numberguess clone, read "no open milestones", no-opped,
and reported to the operator that the implementer had "pinged the wrong
coordinator". The operator's
diagnosis was structural rather than textual: the platform has the abstraction
that was missing — a **Project**, with a codebase pointing at the GitHub
remote — and a CodeCrew hub maps one-to-one onto one. Declaring it changed the
coordinator's next run from a re-clone to `status` in fifteen seconds
([findings 62](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465148858)
and [its fold-back](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465157957)).

Two more stalls closed the cycle and both were the coordinator obeying its own
rules in the wrong place: it told the qa seat *"mention me only if you are
blocked"* — and verdicts are not a GitHub event, so the seat did as told and
nothing moved
([finding 63](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465324165));
then it "dispatched" the doc-synthesizer by posting the platform's mention
syntax as a **GitHub comment**, so no agent woke at all
([finding 65](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465353819)).
Three change points were applied through the platform's API mid-cycle: the
brief gains "which project" and a nameless wake means `status` everywhere
(6); the dispatch template is written **per seat** rather than as a rule the
coordinator interprets (7); and, under Never, "dispatch a seat on GitHub" —
"on GitHub you cite, you never dispatch" (8).

### What the run folded back

Cycle 4's fold-back map is on the record and every line of it landed inside
this milestone:

| findings | became |
|---|---|
| 52, 53, 54, 62, 63, 64, 65 | the coordinator contract's obligations — the per-seat wake table, the project as a wake parameter, dispatch on the platform / cite on GitHub, a fresh child ticket per dispatch, pre-milestone gates on the scaffold PR ([#168](https://github.com/radiusred/gh-codecrew/issues/168) / PR [#169](https://github.com/radiusred/gh-codecrew/pull/169)) |
| 56, 67 | `identity token` ([#170](https://github.com/radiusred/gh-codecrew/issues/170) / PR [#171](https://github.com/radiusred/gh-codecrew/pull/171)) |
| 57 | the blank local extensions, in M7-R4's amended form ([#173](https://github.com/radiusred/gh-codecrew/issues/173) / PR [#174](https://github.com/radiusred/gh-codecrew/pull/174)) |
| 58, 66 | capture [#165](https://github.com/radiusred/gh-codecrew/issues/165) → `NOT_OWNER` ([#175](https://github.com/radiusred/gh-codecrew/issues/175) / PR [#176](https://github.com/radiusred/gh-codecrew/pull/176)) |
| 51, 68 | capture [#172](https://github.com/radiusred/gh-codecrew/issues/172) → `init` commits its scaffold ([#183](https://github.com/radiusred/gh-codecrew/issues/183) / PR [#184](https://github.com/radiusred/gh-codecrew/pull/184)) |
| 46, 53 | the narrowed default event set ([#180](https://github.com/radiusred/gh-codecrew/issues/180) / PR [#181](https://github.com/radiusred/gh-codecrew/pull/181)) |
| 55, 59, 60, 61, 62 | the interop page's onboarding checklist ([#182](https://github.com/radiusred/gh-codecrew/issues/182) / PR [#185](https://github.com/radiusred/gh-codecrew/pull/185)) |

The #119 open item was decided rather than deferred: **the platform's tickets
are not committed as curated transcripts** — the operator's read, twice, is
"more noise than they're worth", and the findings logs quote what the
fold-backs needed. The consequence is stated on the interop page's gap list:
the run is not fully reconstructible from any repository, the same class of
gap M5 recorded when a probe repository was deleted.

Three captures this milestone filed rather than fixed are still open: a README
strapline the operator asked for during cycle 4
([#166](https://github.com/radiusred/gh-codecrew/issues/166)); stale branches
left by closes that predate the current sweep, found on numberguess
([#167](https://github.com/radiusred/gh-codecrew/issues/167)); and a spoke
shared by several hubs, which SPEC §3's one-hub-per-spoke rule does not yet
answer ([#177](https://github.com/radiusred/gh-codecrew/issues/177)).

## What cycle 4 did not exercise of the coordinator contract

#163's fourth gate reads: *"The coordinator contract (M7-R1) is exercised by
cycle 4's second milestone before it ships, or the milestone document says
which parts were not."* Cycle 4 did not exercise it, and this section is the
alternative the gate names.

The timing settles the general case. `milestone close 2` ran on snake at
23:10:23Z on 2026-08-29; PR
[#169](https://github.com/radiusred/gh-codecrew/pull/169), which added
`roles/coordinator.md`, merged at 23:45:00Z the same evening — **thirty-five
minutes after the cycle ended**. Cycle 4 ran on Loopy's hand-written
`AGENTS.md`/`HEARTBEAT.md` bundle, carried over from cycle 3 and patched
mid-cycle through the platform's API. The contract was written *from* that
bundle and from the findings the cycle produced; the interop page records the
same limitation in its own words — "the contract itself awaits its first
cycle" ([platform-interop.md](../platform-interop.md)). What follows is the
clause-level detail that statement stands for.

**Not exercised at all.**

- **The Identity section, all but one clause.** snake's `.codecrew.yml` declares four
  roles — implementer, reviewer, qa, doc-synthesizer — and **no `coordinator`
  row**; it was scaffolded by v1.0.3, which had none to write. The coordinator
  agent therefore held no seat in the routing table it was coordinating, and
  `roles.coordinator.identity` was never read by anything. Its App was not
  minted by `gh codecrew identity new coordinator`, which shipped with the
  contract. Of the permission set the contract states, only `contents: read`
  is on the record: the coordinator "concluded correctly that it cannot run
  `init` itself (`contents: read`)" and delegated the scaffold
  ([entry 1](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218)),
  and the implementer carried the ROADMAP row because the coordinator cannot
  commit. Whether that App holds `issues: write`, `pull requests: read`, or —
  the clause that matters — *not* `contents: write` and *not* `pull requests:
  write`, is not recorded anywhere in the cycle-4 log. **Undocumented, and not
  inferable from the run beyond the one permission its behaviour demonstrated.**
- **The mint.** `export GH_TOKEN=$(gh codecrew identity token <slug>)` did not
  exist during cycle 4 — PR
  [#171](https://github.com/radiusred/gh-codecrew/pull/171) merged at 00:04:50Z
  on 2026-08-30. The cycle ran on the brief's hand-written JWT sketch, and
  finding 56 is the record of that sketch being wrong.
- **The contract as a composed bundle.** Nothing in cycle 4 ran
  `gh codecrew roles show coordinator`; there was no `roles/coordinator.md` to
  compose and no `roles/coordinator.local.md` to compose with it. Nor was the
  composition path exercised for any other seat: the crew ran on bundles
  composed for a different repository, and snake carried no
  `roles/<role>.local.md` at all
  ([finding 57](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218)).
- **`checkpoint`, and the pre-milestone gate.** Both snake milestones recorded
  **zero** protocol gates; M1's two human decisions travelled on the
  platform's own confirmation mechanism. The scaffold-merge authorization is
  precisely the pre-milestone gate the contract now requires to be recorded on
  the scaffold PR as a `**Gate raised:**` / `**Gate resolved:**` pair — and in
  cycle 4 it left no such pair, which is
  [finding 52](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218),
  the finding that wrote the clause. The clause has never been run.
- **`NOT_OWNER` and the `--dry-run` previews**, both cited by name in the
  contract's review-loop and milestone-end bullets. Both shipped after the
  cycle (PRs [#176](https://github.com/radiusred/gh-codecrew/pull/176) and
  [#179](https://github.com/radiusred/gh-codecrew/pull/179), 2026-08-30).
  Cycle 4 is the evidence *against* the first: the verb merged another seat's
  task (finding 58). The rule held for the rest of M2 by prose alone
  ([finding 66](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465449989)).
- **"Never write the milestone number into requirement prose."** The clause
  comes from numberguess and M4 (#119); the cycle-4 log records nothing for or
  against it.
- **The wake table as a table.** Change point 7 rewrote the dispatch templates
  per seat for the remainder of M2, so two of its six rows — the qa seat's
  verdicts and the doc-synthesizer's hand-back — were exercised in that form,
  and the "task's owner finishes it" row was exercised in prose. The two rows
  that say a GitHub-emitted transition is *never* also hand-mentioned were
  exercised only after change points 6–8 removed the mention tail, and the
  six-row table itself is #169's rendering of those lessons. No cycle has read
  it.

**Not recorded either way.** The contract's milestone-end order opens with
`gh codecrew milestone evidence <milestone number>`. #164's entries record the
verdicts, both documents and both closes, but no `evidence` run in either
milestone. The record does not say it did not happen; it does not say it did.
**Undocumented; not inferred.**

**Exercised in the brief's form, before or after the mid-cycle change points.**
One wake, one unit of work, and execution events one-shot (findings 53 and 54,
as failures, then held); milestones opened with `--requirement` and tasks
opened with plans; dispatch by the four-row routing table; the review loop
owned in both directions, twice through changes-requested rounds; the record
re-read at the act, though against the PR only and not the coordinator's own
open dispatches — which is exactly what the contract's wording adds; a wake
naming its repository (change point 6); dispatch on the platform and citation
on GitHub (change point 8, with the doc-synthesizer dispatch that followed it);
`**Decision:**` comments on the GitHub record; and the 403 discipline — the
seat delegated the scaffold and the ROADMAP row rather than committing them.

**The honest summary.** Every obligation in `roles/coordinator.md` traces to a
finding, and most of them trace to a finding cycle 4 produced by violating the
obligation before it existed. What has not happened is a cycle driven by the
contract itself — installed from `roles show coordinator`, under an identity
minted by `identity new coordinator`, routed by a `coordinator` row, raising
gates with `checkpoint` and previewing verbs with `--dry-run`. That is cycle
5's job, and until it runs the contract is a well-sourced hypothesis rather
than a proven one.

## The ladder's last rung, and the spoke

**`docs/platform-interop.md`** (M7-R7, from #54) is nine sections and 610
lines: the separation of concerns; the coordinator as a seat and its own
agent; mapping agents to roles; credentials; wake paths and the one-wake-path
rule; an eleven-row onboarding checklist, each row carrying the finding that
priced it; the cost tables; the Paperclip recipe with placeholders for every
id, hostname and secret; and seven named gaps. The reviewer's audit of it is
itself a receipt — "177 Markdown link occurrences (87 unique), no missing
relative targets; 58 GitHub URL targets, including comment anchors, resolved
through the API"
([review](https://github.com/radiusred/gh-codecrew/pull/185#pullrequestreview-5061284480)).

Its first review round caught something the milestone then had to fix twice
over: the page described v1.1.0 behaviour as current while the installed
release was v1.0.3. The answer was one **"A note on releases"** paragraph
naming the six things that ship in v1.1.0, each with its PR, and what the
v1.0.3 path is for each — "which is precisely what the run did, at the cost
its findings record" — plus five qualified sentences in the body rather than
"next release" sprinkled throughout. The same comment flagged that
`docs/introduction.md` had the identical problem and that fixing it belonged
to the release task, not this one, "so it does not fall between the two"
([comment](https://github.com/radiusred/gh-codecrew/pull/185#issuecomment-5469772791)).
It did not: #186 flipped it.

**The www spoke carried three tasks.** The
[article](https://www.radiusred.uk/blog/posts/2026-08-30-four-cycles-on-a-real-orchestration-platform/)
(www#51 / PR #52), ~2,100 words, both cost tables reproduced verbatim, v1.1.0
named as forthcoming and never as installed, and the visibility rules checked
rather than assumed — no ticket ids, no agent ids, no instance host, and both
proving grounds confirmed public before linking. The
[announcements](https://github.com/radiusred/www/pull/54) (www#53 / PR #54)
behind the gate above, with Bluesky measured at 293 of 300 graphemes by the
module rather than by eye, and tags chosen from measured reach — the finding
worth keeping is that "`#Orchestration`, `#MultiAgent` and `#AutonomousAgents`
are the tags instinct reaches for on this subject and are, measurably, empty
rooms." And the pointer flip (www#56 / PR #57): `codecrew: "0.1"` → `"1.0"`,
one line, proved against the installed v1.1.0 binary.

## The release

v1.1.0 was cut on 2026-08-30 at 23:19:30Z, from the merge commit of PR
[#187](https://github.com/radiusred/gh-codecrew/pull/187) (`a7a459d`), by the
implementer identity. The `[1.1.0]` heading is dated 2026-08-31, which is the
merge commit's own date in the repository's local timezone
(`2026-08-31 00:18:11 +0100`) rather than a discrepancy with the tag — the
class of drift M6's document caught on 1.0.0, checked and absent here.

The flip itself is small and deliberately so: the whole `[Unreleased]` section — nine M7 headings plus the two M6 document entries —
became `## [1.1.0]` verbatim, "nothing dropped, nothing reworded", a fresh
empty `[Unreleased]` above it, and the compare links extended. Only the
introduction's Shipped paragraph moved with it; the verb list and the
twenty-nine-code catalogue had been kept current by each task as it landed,
which is the discipline M6 established and this milestone did not have to
repeat.

The release carries every M7 task: the coordinator seat, `identity token`,
blank local extensions from the scaffold, `NOT_OWNER`, `identity webhook` with
`--webhook-secret`, `init` committing its scaffold, `--dry-run` on three
verbs, the interop page and the M6 record. Every verb
and flag is additive and every code is added rather than repurposed, which is
what SPEC §10's 1.0 promise obliges. The catalogue went 23 → 27 → 28 → 29
across #170, #175 and #180, each task cataloguing its own.

## Deviations

Two were recorded across thirteen tasks, and both were review-driven — one
from a reviewer reproducing a failure the plan had not anticipated, one from
`main` moving under a branch that was already in review.

**The bootstrap branch is cut from the default branch, not HEAD**
([#183](https://github.com/radiusred/gh-codecrew/issues/183#issuecomment-5469566151)).
The plan said `codecrew-bootstrap` is cut from HEAD; as merged it is cut from
`refs/remotes/origin/<default>`, else the local branch, HEAD only when neither
exists. **Why:** the reviewer reproduced the failure — from a feature branch
the scaffold PR would have carried the feature commits, "and the scaffold PR
must be scaffold-only for the one hand merge it needs to be safe". A second
deviation from the same review: the quickstart pushes `HEAD`, not `main`,
because `init` commits on whatever branch the operator is on and prints the
exact command. The reviewer had reproduced that too, in a repository whose
branch was named `trunk`.

**The interop page's gap entry was rewritten mid-review**
([#182](https://github.com/radiusred/gh-codecrew/issues/182#issuecomment-5469598774)).
PR #184 merged while PR #185 was open, which closed #172 — so the page's
"What is not solved yet" entry claiming the scaffold PR lands outside the
protocol "was true when written and false by the time it would be read, and a
page whose whole claim is that every statement traces to the record cannot
carry a stale one". The rewritten entry separates what narrowed (the scaffold
is a commit; delete-on-merge sweeps the branch finding 68 found stranded) from
what did not (behind a ruleset the scaffold still arrives as a PR and is the
one merge the operator does by hand; on a fresh repository whose org also
requires a check that cannot report yet, that merge is an administrator merge
— three repositories have met it). The rebase produced one conflict, in
`CHANGELOG.md`, resolved by keeping both entries intact in merge order.

Two deviations across thirteen tasks, against M6's eleven across eighteen.
The record offers no explanation for the difference and none is inferred here;
the count is stated because it is the sort of thing a reader comparing two
milestones will notice.

## The review rounds

Ten hub PRs and three www PRs, every one approved by the reviewer seat before
merge and every one merged on green CI. Six of the thirteen needed a second
round from the reviewer seat or more; a seventh, www PR #54, needed a second
round from the operator.

- **PR #171** (`identity token`): the documented mismatched-App probe did not
  produce the promised refusal. With a real key and `GITHUB_APP_ID=1`, GitHub
  answers 404 `Integration not found`, not 401, so the CLI emitted an
  uncatalogued raw error — "this makes the SPEC/docs/CHANGELOG claim false for
  a required real command". Both shapes now refuse `BAD_CREDENTIALS`, with the
  real 404 modelled in the fake server.
- **PR #176** (`NOT_OWNER`): the ownership record was neither authenticated
  nor ordered. Any comment beginning `**Started by** @…` changed the effective
  owner — "and can let that identity pass `NOT_OWNER`, despite `task start`
  never having run as it" — and human restarts left multiple assignees with no
  latest-wins. Fixed by making `task start` post the canonical record on every
  start and accepting a record only when its author is the login it names.
- **PR #179** (dry runs): the plan/execute split silently changed live output
  on a gated path — an unrouted-qa note that `main` prints before a
  `DOC_MISSING` refusal was deferred into a closure that never ran. The
  refactor's whole premise is that the preview cannot disagree with the verb,
  and the finding is the same premise applied to the verb against itself.
- **PR #181** (`identity webhook`): four rounds, the milestone's longest. Two
  rounds on untested paths and an unsupported case; then two on one sentence —
  the claim that the receiver's secret is set "before its first delivery",
  which the reviewer argued was a race the implementation could not prevent.
  It was resolved not by weakening the claim but by tightening it: repository
  events reach an App only through an installation, and at manifest-conversion
  time the App has none, so the creation `ping` is the only delivery that can
  precede the PATCH. The reviewer accepted the argument on the fourth round.
- **PR #184** (`init` commits): four findings in round one, all reproduced
  live — a subdirectory run that wrote files and made no commit, a detached
  HEAD that stranded exactly the commit the change exists not to strand, the
  feature-branch bootstrap above, and a push command that fails on any branch
  not named `main`.
- **PR #185** (the interop page): the release-boundary finding above.

Every round that requested changes named a defect that was then fixed, which
is the seventh milestone in a row that can say so; and three of the six
findings above were found by *running* the documented command in a real
repository rather than by reading the diff.

## QA: one round, nine satisfied

The qa role holder (`radiusred-testy[bot]`) verdicted all nine requirements in
a single comment at 23:34:45Z on 2026-08-30, fifteen minutes after the release
was published, against merged `main` and the installed binary. Every verdict
opens with green tests and then cites a probe the suite does not enumerate,
which is the contract M6 rewrote.

The probes worth reading: an installed v1.1.0 composed a disposable
`roles/coordinator.local.md` in the www spoke *after* the embedded contract
and read the permission set out of the live output; the mint verb was run
against a real hidden installation, a space-bearing slug, a nonexistent slug
and a genuine orphan key — the last reaching the distinct missing-stub detail
and exit 1 **without printing a token**; `identity webhook --show` was pointed
at a real App with no hook and correctly refused `NO_WEBHOOK` while changing
nothing; `milestone new --dry-run` predicted M8, its requirement ID and its
ROADMAP row without creating them, and `task finish 168 --dry-run` against a
closed task printed the ordered gate plan and stopped at `CLOSED`; the
cycle-4 log was read from GitHub rather than from fixtures, and found to
expose both milestones, the touch split, the table, findings 51–68, the
transcript decision and the fold-back map independently; the interop page's
live citations were followed; the www record was checked for the tone pass,
the gate resolution order and the posted URLs; and the release asset was
downloaded, hashed and executed.

The M7-R1 verdict is the one to read beside the gate section above. It
records, in its own words, that "cycle 4 used the hand-written bundle, and
`docs/platform-interop.md` explicitly records the composed contract as not yet
platform-exercised, satisfying the document-side alternative of the gate" —
so the requirement was verdicted satisfied on the understanding that this
document would carry the detail. It does.

## Requirement outcomes

| Requirement | Delivered by | Outcome |
|-------------|--------------|---------|
| M7-R1 — a coordinator role contract, routing row and least-privilege identity | [#168](https://github.com/radiusred/gh-codecrew/issues/168) / PR [#169](https://github.com/radiusred/gh-codecrew/pull/169); enforcement [#175](https://github.com/radiusred/gh-codecrew/issues/175) / PR [#176](https://github.com/radiusred/gh-codecrew/pull/176) | Done; **satisfied**, first round — with the gate's document-side alternative discharged above |
| M7-R2 — `identity token <slug>` (#132) | [#170](https://github.com/radiusred/gh-codecrew/issues/170) / PR [#171](https://github.com/radiusred/gh-codecrew/pull/171) | Done; **satisfied**, first round |
| M7-R3 — an App's webhook signs for a platform receiver (#157) | [#180](https://github.com/radiusred/gh-codecrew/issues/180) / PR [#181](https://github.com/radiusred/gh-codecrew/pull/181) | Done; **satisfied**, first round |
| M7-R4 — local extensions from the scaffold (#159, #158, #138), **amended 2026-08-30** | [#173](https://github.com/radiusred/gh-codecrew/issues/173) / PR [#174](https://github.com/radiusred/gh-codecrew/pull/174); [amendment](https://github.com/radiusred/gh-codecrew/issues/163#issuecomment-5468874006) on #163 | Done; **satisfied**, first round — against the amended text |
| M7-R5 — dry runs for `milestone new`, `task finish`, `milestone close` (#133) | [#178](https://github.com/radiusred/gh-codecrew/issues/178) / PR [#179](https://github.com/radiusred/gh-codecrew/pull/179) | Done; **satisfied**, first round |
| M7-R6 — cycle 4 recorded and folded back | [#164](https://github.com/radiusred/gh-codecrew/issues/164) / PR [#188](https://github.com/radiusred/gh-codecrew/pull/188) | Done; **satisfied**, first round |
| M7-R7 — the platform-interop doc (#54) | [#182](https://github.com/radiusred/gh-codecrew/issues/182) / PR [#185](https://github.com/radiusred/gh-codecrew/pull/185) | Done; **satisfied**, first round |
| M7-R8 — the Paperclip-experiment article, tone pass and gated announcement | www [#51](https://github.com/radiusred/www/issues/51) / PR [#52](https://github.com/radiusred/www/pull/52); www [#53](https://github.com/radiusred/www/issues/53) / PR [#54](https://github.com/radiusred/www/pull/54) | Done; **satisfied**, first round |
| M7-R9 — v1.1.0 ships and the docs are true of it | [#186](https://github.com/radiusred/gh-codecrew/issues/186) / PR [#187](https://github.com/radiusred/gh-codecrew/pull/187); www [#56](https://github.com/radiusred/www/issues/56) / PR [#57](https://github.com/radiusred/www/pull/57) | Done; **satisfied**, first round — `milestone close 7` runs after this document merges |

Supporting work with no requirement ID of its own:
[#183](https://github.com/radiusred/gh-codecrew/issues/183) / PR
[#184](https://github.com/radiusred/gh-codecrew/pull/184), `init` committing
its scaffold — a cycle-4 fold-back that serves M7-R9's clean repository at the
close. Captures adopted and closed inside the milestone: #133, #138, #157,
#158, #159 (filed during M6's run), #54 (open since 2026-08-22), and #165 and
#172 (opened by cycle 4). Adopted, shipped and **still open**: #132. Captures
left open by choice: #166, #167, #177.

## Protocol-discipline observations

- **A requirement was amended without a gate, and the difference shows.**
  M6-R1's amendment and M7-R4's amendment reached the same place — a
  requirement rewritten in the open with the original quoted inside it — by
  different routes. M6's went up as `**Gate raised:**` with the options
  enumerated and came back as `**Gate resolved:**`; M7's arrived as a Decision
  comment recording an answer whose question is nowhere on the record. The
  Decision is unusually thorough and nothing was smoothed over, so the loss is
  small — but it is a loss, and it is the sort a synthesizer only notices by
  looking for the gate that should be beside the decision.
- **The gather has still never been the raw material for the document it
  gathers for.** M6-R3 fixed the gatherer and M6's own document could not use
  it, because `DOC_MISSING` stands between the close and the record. M7 added
  `--dry-run` and the preview stops earlier still — at `OPEN_TASKS`, naming
  this document's own task. Three milestones now. The ordering is inherent, not
  a defect, but the fix has yet to pay off once.
- **Operator review findings arrived where they do not travel.** On www#52 the
  tone pass came back as four commit comments on `bdac473`, which "do not
  travel with the PR", so the seat quoted each one into a PR comment alongside
  what changed for it. On #185 the operator's approval was dismissed by the
  fix push and now shows as a dismissed review with an empty body — the tone
  pass is recoverable from the dismissal event, not from the review. Both are
  small; both mean a reader reconstructing the milestone from the PR alone
  would miss a gate that actually ran.
- **The milestone shipped the enforcement for a rule it had just broken.**
  The wrong-seat merge happened on snake at 17:50:23Z on 2026-08-29; PR
  [#176](https://github.com/radiusred/gh-codecrew/pull/176), which makes
  `task finish` refuse it, merged at 13:39:33Z the next day — under twenty
  hours from observed defect to enforced refusal, through a capture, a task
  and two review rounds. The reviewer then proved it from the other side: it
  ran `task finish 175` as `radiusred-checky` on the implementer's own task
  and got `refused[NOT_OWNER]` naming `radiusred-cody[bot]`, with the PR left
  unmerged.
- **A capture shipped without closing.**
  [#132](https://github.com/radiusred/gh-codecrew/issues/132), the capture
  M7-R2 was adopted from, is still open: task #170 and PR #171 both cite it in
  prose but neither carries a `Closes #132`. Every other capture this
  milestone adopted was closed by the PR that discharged it — #133 and #138,
  #157, #158, #159, #165, #172 and #54. The verb #132 asked for is shipped in
  v1.1.0 and catalogued; only the capture's state disagrees. Nothing depends
  on it, and it is exactly the kind of drift the close's sweep does not look
  for.
- **The reviewer changed the design four times, not just the code.** The
  bootstrap branch's parent, the ownership record's authentication, the
  webhook secret's stated boundary and the interop page's release boundary
  were all design decisions revised in review, and three of them were found by
  running the documented command in a real repository rather than by reading
  the diff. Only the first became a Deviation; the other three were absorbed
  into their PRs, which means the record of *why* those three shapes are what
  they are lives in review bodies rather than in a Decision or a Deviation
  anyone would gather.
- **The coordinator seat is shipped and unproven, and both halves are on the
  record.** The interop page says so, the QA verdict says so, and this
  document itemises it. That is the honest position at the boundary: M7 built
  the seat from four cycles of evidence about what a coordination layer gets
  wrong, and the fifth cycle is what will say whether writing it down was
  enough.
