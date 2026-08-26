# M5: Crew expansion

Tracking issue: [#67](https://github.com/radiusred/gh-codecrew/issues/67) ·
Synthesized 2026-08-26 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the records
gathered by `codecrew milestone close 5` (three: the gate resolution on #68,
the release-parity Decision on #73, a Deviation on #76 — three further
Decisions on #73 escaped the gatherer; see the discipline section), the
thirteen findings and four Decisions on
[#73](https://github.com/radiusred/gh-codecrew/issues/73) — the milestone's
spine — the task plans and merged PR descriptions (#68 through #105), the
proving-ground record in
[davison/numberguess](https://github.com/davison/numberguess)
([M2_SESSION_TRANSCRIPT.md](https://github.com/davison/numberguess/blob/main/M2_SESSION_TRANSCRIPT.md),
[M3_SESSION_TRANSCRIPT.md](https://github.com/davison/numberguess/blob/main/M3_SESSION_TRANSCRIPT.md),
PRs #8, #11, #13, #15), the cc-approval-probe experiment as preserved in the
superseding Decision, and QA's two verdict rounds on #67. The close itself
lands when `milestone close 5` runs after this document merges — it has
already refused `DOC_MISSING` once, exactly as designed.

## Goal and outcome

Move from the solo pair to a crew: staff the last unstaffed seat with a
dedicated reviewer App driven by a second harness, decide how an operator
without a pre-built App fleet gets attributable role identities, let a role
be held by a GitHub team, and fold back the two M4 discipline gaps that were
cheap to fix and bit anyway
([#67](https://github.com/radiusred/gh-codecrew/issues/67)).

Delivered in full, after one remedy round. At open, the reviewer seat was
the only one still human-routed (`reviewer: davison`, the M4 arrangement);
at close, all four seats are App-staffed —
[radiusred-cody](https://github.com/apps/radiusred-cody) (implementer),
[radiusred-checky](https://github.com/apps/radiusred-checky) (reviewer,
minted with the milestone's own verb),
[radiusred-testy](https://github.com/apps/radiusred-testy) (qa),
[radiusred-wordy](https://github.com/apps/radiusred-wordy) (doc-synthesizer)
— and the flip itself was agent-gated: PR
[#100](https://github.com/radiusred/gh-codecrew/pull/100) routed checky as
reviewer, checky's first dispatch reviewed that very PR, and its counted
approval satisfied both CodeCrew's holder gate and GitHub's required-review
rule. The seat that reviewed its own seating.

Getting there grew the CLI by `identity new` (guided App minting via the
manifest flow, with `--with-webhook`, `--with-approval-permission`, and
auto-routing; [#70](https://github.com/radiusred/gh-codecrew/issues/70) /
PR [#71](https://github.com/radiusred/gh-codecrew/pull/71) plus four
live-fire fold-backs), `roles diff` / `roles show --latest` with drift
detection in `status` ([#85](https://github.com/radiusred/gh-codecrew/issues/85)
/ PR [#93](https://github.com/radiusred/gh-codecrew/pull/93)),
`milestone evidence` ([#98](https://github.com/radiusred/gh-codecrew/issues/98)
/ PR [#102](https://github.com/radiusred/gh-codecrew/pull/102)), team
routing for `identity:`
([#97](https://github.com/radiusred/gh-codecrew/issues/97) / PR
[#101](https://github.com/radiusred/gh-codecrew/pull/101)), and
`task finish`'s `REVIEW_NOT_COUNTED`/`--bypass` and `NO_HOLDER_REVIEW`
gates ([#82](https://github.com/radiusred/gh-codecrew/issues/82) / PR
[#89](https://github.com/radiusred/gh-codecrew/pull/89),
[#91](https://github.com/radiusred/gh-codecrew/issues/91) / PR
[#95](https://github.com/radiusred/gh-codecrew/pull/95)). Four releases
shipped — v0.3.0, v0.3.1, v0.4.0, v0.5.0
([#72](https://github.com/radiusred/gh-codecrew/issues/72),
[#80](https://github.com/radiusred/gh-codecrew/issues/80),
[#99](https://github.com/radiusred/gh-codecrew/issues/99),
[#103](https://github.com/radiusred/gh-codecrew/issues/103)) — under a new,
scoped release-parity rule the milestone's first live fire forced into
being.

## Decisions

- **Guided App creation via the manifest flow (`identity new <role>`), not a
  template or docs-only tiering.** The M5-R2 gate, raised on #68 before any
  commit with four candidate paths honestly enumerated: (a) a guided-creation
  verb, (b) a published manifest template with no verb, (c) docs-only
  tier-sharpening so adopters mint Apps only when a gate demands independent
  review, (d) fine-grained PATs from a second account (weak attribution —
  enumerate, likely reject). A radiusred-published shared App was off the
  table by prior decision: the framework is never key custodian for an
  adopter's org. The operator chose (a): the logical progression for solo
  operators growing with their project, and near-required for orchestration
  platforms like Paperclip. Trade-off, recorded in the resolution: it may be
  a poor option for human teams using CodeCrew on isolated milestones, so
  the docs must carry the recommendation split for both human and agent
  consumption — which
  [docs/identities.md](../identities.md) now does.
  ([Gate raised](https://github.com/radiusred/gh-codecrew/issues/68#issuecomment-5395550079)
  and [resolved on #68](https://github.com/radiusred/gh-codecrew/issues/68#issuecomment-5395599850);
  decision recorded by PR [#69](https://github.com/radiusred/gh-codecrew/pull/69),
  implemented by PR [#71](https://github.com/radiusred/gh-codecrew/pull/71).)
- **Release parity binds proof runs, not iteration.** The M4 principle — QA
  and stranger-path evidence must cite the binary an adopter installs — met
  its first cost when the shipped v0.3.0 verb failed live (finding 1, below)
  and a release per fix would have meant version churn per bug. The operator
  scoped the principle: discovery runs in a live-fire loop use local source
  builds (`codecrew version` prints `dev`, keeping the tier honest in any
  transcript), fold-back fixes ride normal tasks with no release each, and a
  single release precedes the recorded proof run, which re-runs the critical
  path end to end against the shipped binary — so release-build-only defects
  (the #51/#55 class) are still caught where the record depends on them.
  Trade-off: intermediate iterations lose release-pipeline coverage.
  Rejected: a public patch release per fold-back (an ops task per bug);
  prerelease tags with `--pin` (keeps per-iteration coverage but retains the
  ceremony without the proof-run benefit).
  ([Decision on #73](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5398340528);
  applied as v0.3.1 before the numberguess proof run
  ([#80](https://github.com/radiusred/gh-codecrew/issues/80)) and v0.4.0/v0.5.0
  before the hub flip and QA dispatch
  ([#99](https://github.com/radiusred/gh-codecrew/issues/99),
  [#103](https://github.com/radiusred/gh-codecrew/issues/103)).)
- **The counting question, answered twice — and the second answer supersedes
  the first.** The requirement demanded an empirical answer to "does a
  GitHub App's approving review count toward a required-review rule?" before
  any `task finish` expectation was wired to it. The first Decision, from
  the numberguess M3 proof run: **never** — on
  [numberguess#11](https://github.com/davison/numberguess/pull/11) CodeCrew's
  gate accepted `davison-review-bot`'s approval while GitHub held
  `reviewDecision: REVIEW_REQUIRED`, `--admin` was rejected
  (`bypass_actors: []`), and the merge needed a bypass actor; the predicted
  fallback — operator approves on the agent's recommendation — turned out
  unavailable exactly in pure solo, because the operator authored the PR and
  GitHub forbids self-approval. Supported configurations were recorded with
  their platform footprint, and the rejected alternatives with them: wiring
  `task finish` to expect App approvals to count (the platform said no), and
  telling adopters to drop required-review rules (loses GitHub-native
  enforcement where humans participate).
  ([First Decision on #73](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5412720900).)
  Then an operator hypothesis, proven live on a throwaway repo
  (radiusred/cc-approval-probe#1, since deleted), superseded the claim:
  under the org ruleset, `reviewDecision` flipped to APPROVED on
  `radiusred-testy[bot]`'s approval alone. Testy holds `contents: write`;
  davison-review-bot holds the deliberate `contents: read`. **Corrected
  rule: a write-access App's approving review counts; a read-only App's
  does not** (org-vs-personal not fully excluded as a co-factor — the
  discriminating experiment is recorded for whoever wants the last word).
  Consequences: no gate-logic change (`mergeGate` keys on `reviewDecision`,
  mechanism-agnostic), but every "never counted" text had to be corrected
  ([#90](https://github.com/radiusred/gh-codecrew/issues/90) / PR
  [#92](https://github.com/radiusred/gh-codecrew/pull/92)), and the hub-flip
  calculus changed: a write-access reviewer App can fully satisfy the org
  gate. The trade-off that opens — privilege beyond contract need, held by
  an auditable identity, versus least-privilege plus human-approval paths —
  was left as a per-project operator's choice, surfaced at minting as
  `identity new reviewer --with-approval-permission`: reviewer-only, never
  the default, the trade stated by the verb at minting
  ([#94](https://github.com/radiusred/gh-codecrew/issues/94) / PR
  [#96](https://github.com/radiusred/gh-codecrew/pull/96) — operator
  direction recorded in the plan, no Decision comment).
  ([Superseding Decision on #73](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5413715331).)
- **A model review is the protocol's expected review; the holder's approval
  is the one that counts.** Operator direction: a clean-context model review
  under the reviewer contract — optionally a different harness — is the
  expected review whether or not a human gate also applies. Mechanized where
  a distinct principal holds the seat: `task finish`'s review gate tightened
  from "any non-doer approval" to *an approving review from the reviewer
  role's holder* (`refused[NO_HOLDER_REVIEW]` otherwise; other approvals,
  CODEOWNERS included, coexist but do not satisfy the gate, and
  `--operator-confirm` cannot stand in when a holder exists). In pure solo
  the contract and scaffold strongly encourage a dispatched clean-context
  reviewer session whose findings land as a PR comment before confirmation.
  Evidence base cited by the Decision: the numberguess M2/M3 reviews and the
  hub's own review cycle, where clean-context App review caught
  implementer-blind defects in every round it ran. Trade-off: teams lose
  "anyone senior approves" unless they route the seat accordingly — routing
  is the declared intent, so that is a feature.
  ([Decision on #73](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5413733166);
  PR [#95](https://github.com/radiusred/gh-codecrew/pull/95).)
- **The current contract rides the binary; `_latest` is virtual.** The
  contract-propagation design: `init` stamps each scaffolded contract with a
  provenance header naming the release and upstream path, `status` reports
  drift between local `roles/` and the embedded copies, `roles diff` and
  `roles show --latest` give the comparison and full upstream text on
  demand, and the AGENTS.md scaffold carries a standing reconciliation
  instruction — contracts are the project's own fork, reconciled through a
  task and PR, never overwritten. Rejected: a physical `_latest` shadow file
  (authority confusion, side-effect writes, messenger staleness). No
  Decision comment exists — the design and its rejected alternative are
  recorded only in the [#85 plan](https://github.com/radiusred/gh-codecrew/issues/85)
  and [PR #93](https://github.com/radiusred/gh-codecrew/pull/93), the M4
  "judgment lives in plans and PR bodies" gap recurring.

## The proving ground: davison/numberguess

M5-R1 forbade changing the hub's own seat until the pattern was proven
elsewhere: the operator minted a reviewer App in
[davison/numberguess](https://github.com/davison/numberguess) (the M4
test-drive project) with the new verb, and a Codex harness drove it through
two feature milestones there. Three live fires reshaped the middle of the
milestone; every finding was recorded on #73 as it happened and folded back
as a hub task.

**Live fires 1–2: the verb itself.** The shipped v0.3.0 minting leg failed
at GitHub's manifest validation with a misleading `"url" wasn't supplied` —
the real defect was that the webhook-off default sent
`"hook_attributes": {"active": false}`, and GitHub requires
`hook_attributes.url` whenever the object is supplied at all
([finding 1](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5398340272);
fixed by omitting the object entirely,
[#74](https://github.com/radiusred/gh-codecrew/issues/74) / PR
[#75](https://github.com/radiusred/gh-codecrew/pull/75)). The second run
minted end to end and surfaced ergonomics: the verb should write the role's
routing into a local `.codecrew.yml` itself rather than print an
instruction ([finding 2](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5398479264);
[#76](https://github.com/radiusred/gh-codecrew/issues/76) / PR
[#77](https://github.com/radiusred/gh-codecrew/pull/77)).

**Live fire 3: the M2 review dispatch**
([findings 3–6](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5402963241),
[7](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5402974446),
[8](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5403029004);
[M2 transcript](https://github.com/davison/numberguess/blob/main/M2_SESSION_TRANSCRIPT.md)).
GitHub cannot resolve an App as a review-requestable user, so the #45
request-review obligation failed — the session raised a gate rather than
working around it, and the resolution became framework behaviour: human and
team holders are *requested*, App-held seats are *dispatched*
([#78](https://github.com/radiusred/gh-codecrew/issues/78) / PR
[#79](https://github.com/radiusred/gh-codecrew/pull/79)). The token
helper's org-only installation discovery could not see a personal-account
App, forcing hand-fed IDs — patched with a `/user/installations` fallback
(PR #79) that then failed in practice too, and was properly fixed by
persisting a credential stub at minting and discovering installations under
the App's own JWT ([finding 11](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5412721150);
[#83](https://github.com/radiusred/gh-codecrew/issues/83) / PR
[#88](https://github.com/radiusred/gh-codecrew/pull/88)). The transcript
also carries the first hard evidence for the milestone's premise: the
delegated clean-context reviewer caught a real bug the implementer's tests
missed — styled-mode detection ignored stream encoding, and an ASCII
pseudo-TTY crashed with `UnicodeEncodeError` — refused approval, and
re-reproduced the exact failing case before approving
([numberguess#8](https://github.com/davison/numberguess/pull/8)). It also
held the contract under pressure: it refused to approve while
`cc:needs-decision` was unresolved, and declined to fall back to the
operator's credential when the App token failed. The blemish, recorded
honestly: the request-changes half of that first review never reached
GitHub as an App-posted review (the token helper 404'd at exactly that
moment), so the R1 promise — both halves App-posted — still needed the
proof run.

**The proof run: M3 under a real ruleset**
([findings 9–12](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5412721150);
[M3 transcript](https://github.com/davison/numberguess/blob/main/M3_SESSION_TRANSCRIPT.md)).
With branch protection requiring one approving review, the counting
experiment ran (the first counting Decision, above), and `task finish`'s
inability to complete in that configuration — gates passed, GitHub refused
the merge, three PRs merged around the protocol's one merge point with
`gh pr merge --admin` — became
`refused[REVIEW_NOT_COUNTED]` and the recorded `--bypass` path
([#82](https://github.com/radiusred/gh-codecrew/issues/82) / PR
[#89](https://github.com/radiusred/gh-codecrew/pull/89)). Codex's runtime
guardrail declined to auto-dispatch the reviewer because nothing in the
scaffold authorized sub-agents — folded back as the conditional standing
authorization in the AGENTS.md scaffold, addressed to the operator's
primary session only, so a platform-dispatched role agent reads a
prohibition, not a licence
([#81](https://github.com/radiusred/gh-codecrew/issues/81) / PR
[#87](https://github.com/radiusred/gh-codecrew/pull/87)). The milestone
template's bold placeholder ID matched the requirement parser and refused
M2's close on a phantom requirement
([#84](https://github.com/radiusred/gh-codecrew/issues/84) / PR
[#86](https://github.com/radiusred/gh-codecrew/pull/86)). And the R1
promise landed whole: on
[numberguess#15](https://github.com/davison/numberguess/pull/15) the App
posted [changes requested](https://github.com/davison/numberguess/pull/15#pullrequestreview-5019676389)
[twice](https://github.com/davison/numberguess/pull/15#pullrequestreview-5019691367)
— audit-trail chronology findings on the milestone document itself, held
until fixed — then
[approved](https://github.com/davison/numberguess/pull/15#pullrequestreview-5019699910),
under the active one-approval ruleset.

## The other proving ground: the hub's own review cycle

The milestone's premise — de-correlated review earns its seat — was
demonstrated on the hub while the hub was being changed.

**Testy's eight-PR cycle.** With the dedicated reviewer App still pending,
the operator directed testy (the qa App, `contents: write`) to stand in as
reviewer for the fold-back batch — PRs
[#86](https://github.com/radiusred/gh-codecrew/pull/86),
[#87](https://github.com/radiusred/gh-codecrew/pull/87),
[#88](https://github.com/radiusred/gh-codecrew/pull/88),
[#89](https://github.com/radiusred/gh-codecrew/pull/89),
[#92](https://github.com/radiusred/gh-codecrew/pull/92),
[#93](https://github.com/radiusred/gh-codecrew/pull/93),
[#95](https://github.com/radiusred/gh-codecrew/pull/95),
[#96](https://github.com/radiusred/gh-codecrew/pull/96) — running every
round by experiment, not reading. Each finding round caught an
implementer-blind defect: the #86 regression test never rendered the real
template ([review](https://github.com/radiusred/gh-codecrew/pull/86#pullrequestreview-5021285539));
the #88 token script died silently under `set -euo pipefail` on a bad stub,
still depended on user credentials, and — found only by exercising the fix
— garbled gh's 403 body into a garbage JWT
([three rounds](https://github.com/radiusred/gh-codecrew/pull/88#pullrequestreview-5021290056));
the #89 bypass comment recorded a merge that might not have happened
([review](https://github.com/radiusred/gh-codecrew/pull/89#pullrequestreview-5021293034));
the #92 "never counted" sweep left two stale code comments asserting the
superseded rule ([review](https://github.com/radiusred/gh-codecrew/pull/92#pullrequestreview-5021685283));
and #95 minted the same stale-prose class fresh in SPEC §8 — the third
round running to catch it
([review](https://github.com/radiusred/gh-codecrew/pull/95#pullrequestreview-5021806132)).
Out of that came a process adoption: the full-branch grep sweep for
falsified phrases, recorded for the gatherer in the
[#92 approval](https://github.com/radiusred/gh-codecrew/pull/92#pullrequestreview-5021703941).
PRs #92 and #93 are also the hub's first fully agent-gated merges:
`task finish` merged both on testy's counted approval alone — no human
approval, no bypass, no gate waived
([finding 13](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5413945476)).

**Checky's first four.** radiusred-checky was minted by the operator with
released v0.4.0 `identity new reviewer --with-approval-permission` — the
flag's first live use; App ID 4719924, `contents: write` verified, the
routing line auto-written by the verb — and its first dispatch reviewed
[PR #100](https://github.com/radiusred/gh-codecrew/pull/100), the PR that
seats it. Its rounds kept the pattern: on
[#101](https://github.com/radiusred/gh-codecrew/pull/101#pullrequestreview-5025370306)
it found the unrouted-role privilege bug — a team member held every
*unrouted* role, cross-lane leakage in exactly the mixed configuration team
routing enables — with a runnable probe demonstrating it; on
[#102](https://github.com/radiusred/gh-codecrew/pull/102#pullrequestreview-5025370430)
it applied the executed-docs obligation *that very diff introduced* (built
the branch, ran `milestone evidence 5` verbatim) and blocked on an
undeclared deviation: the one untested function was precisely the
network-touching `checkHTTP` the plan had singled out for httptest
coverage. Both cleared in one round each;
[#105](https://github.com/radiusred/gh-codecrew/pull/105#pullrequestreview-5025566029)
approved the remedy PR clean.

## QA: two rounds under the verdict gate

Before dispatch, the QA-prep recorded on #67
([prep comment](https://github.com/radiusred/gh-codecrew/issues/67#issuecomment-5418707754)):
v0.5.0 was released and installed so QA would exercise the shipped binary,
and — M5-R4 applied to its own close — the pre-dispatch evidence check ran
as code: `codecrew milestone evidence 5` reported all cited links reachable
across the milestone's 21 issues.

**First round**
([verdicts](https://github.com/radiusred/gh-codecrew/issues/67#issuecomment-5418766885)):
`radiusred-testy[bot]`, as qa role holder, returned three satisfied, one
not satisfied.

- **M5-R1 not satisfied** — everything the seat promises verified live
  (checky routed and counted-approving on #100–#102; the App-posted
  changes-requested-then-approve pair on numberguess#15 under the active
  ruleset; both counting Decisions preceding the wiring merges; the hub
  flip only after the proof stood) — but the requirement's final clause
  failed: no statement of Copilot code review coexistence existed anywhere
  in the docs or history, despite the #73 plan asserting the fold-backs
  already carried it
  ([finding on #73](https://github.com/radiusred/gh-codecrew/issues/73#issuecomment-5418764382)).
- **M5-R2 satisfied** — the gate demonstrably preceded implementation
  (resolution 13:05, first code PR 14:09 on 2026-08-24), the rejected
  alternatives and the per-account installation boundary are in
  docs/identities.md, and checky's minting is the released-binary proof.
  A doc-duplication defect (a passage spliced twice into identities.md) was
  [filed on #70](https://github.com/radiusred/gh-codecrew/issues/70#issuecomment-5418764507)
  as a quality finding, verdict unaffected.
- **M5-R3 satisfied** — the released binary resolves team routing live and
  SPEC §5 states cross-member supersession, the zero-framework-code
  composition with GitHub's counts, and the platform footprint (org-only;
  private-repo counts plan-gated). Honestly scoped: live member-set
  resolution was not exercisable end to end — the org has no team and QA
  may not create one — so that sliver rests on the checky-reviewed #101
  implementation and its tests.
- **M5-R4 satisfied** — both halves verified on the shipped binary against
  this milestone's own close: the executed-not-read obligation is embedded
  in v0.5.0 and matches main, and `milestone evidence 5` ran as the
  dispatch's first act, refusal path probed.

**Remedy:** the R1 clause became
[#104](https://github.com/radiusred/gh-codecrew/issues/104) / PR
[#105](https://github.com/radiusred/gh-codecrew/pull/105) — a "Coexistence
with Copilot code review" section in docs/identities.md (advisory signal
alongside, never the reviewer role: not role-routable, not crew-attributed,
reviews uncounted toward required approvals, paid-plan-gated, and no
substitute for SPEC §8's independent-review layer) — plus the unsplicing of
the #70 duplicate. Docs-only, so no release rode it: the parity Decision
binds what the binary ships, and repo prose is not in it.

**Superseding verdict**
([2026-08-26, against main](https://github.com/radiusred/gh-codecrew/issues/67#issuecomment-5418807399)):
M5-R1 superseded to satisfied — the coexistence section verified clause by
clause against the requirement, the duplicate confirmed gone, all other R1
evidence standing. R2, R3, R4 unchanged and current.

## Deviations

One gathered, on
[#76](https://github.com/radiusred/gh-codecrew/issues/76#issuecomment-5398503549):
`identity new` gained a next-step line beyond the plan, handing the
operator the App's settings URL for the coming CodeCrew logo — the manifest
has no avatar field and no API uploads one, so the verb's best move is the
right URL. Prompted by the operator mid-task.

Two more were recorded where they occurred rather than on task issues, both
surfaced by review: PR #86's test had silently transcribed the template
instead of rendering it (recorded in the review-response comment once
caught), and PR #95's plan-promised refusal-path tests were declared as a
Deviation on the PR with rationale (no Tracker fake in the cli package) —
the honest-record form the contract asks for, as
[testy's re-review noted](https://github.com/radiusred/gh-codecrew/pull/95#pullrequestreview-5021824979).

The four release tasks (#72, #80, #99, #103) had no PR and were closed
directly with `**Release verified.**` evidence comments, per the
established ops-task precedent.

## Requirement outcomes

| Requirement | Delivered by | Outcome |
|-------------|--------------|---------|
| M5-R1 — the reviewer role staffed, proven in numberguess, counting question decided | [#73](https://github.com/radiusred/gh-codecrew/issues/73) / PR #100; fold-backs #78/PR #79, #81/PR #87, #82/PR #89, #90/PR #92, #91/PR #95, #94/PR #96; remedy [#104](https://github.com/radiusred/gh-codecrew/issues/104) / PR #105 | Done; not satisfied → satisfied (superseding verdict) |
| M5-R2 — the no-App-fleet onboarding story decided at a gate | [#68](https://github.com/radiusred/gh-codecrew/issues/68) / PR #69; [#70](https://github.com/radiusred/gh-codecrew/issues/70) / PR #71; fold-backs #74/PR #75, #76/PR #77, #83/PR #88 | Done; satisfied, first round |
| M5-R3 — roles routed to GitHub teams | [#97](https://github.com/radiusred/gh-codecrew/issues/97) / PR #101 | Done; satisfied, first round |
| M5-R4 — the M4 discipline gaps closed and proven on this close | [#98](https://github.com/radiusred/gh-codecrew/issues/98) / PR #102 | Done; satisfied, first round |

Supporting work with no requirement ID: the four parity releases (#72,
#80, #99, #103), the milestone-template placeholder fix (#84 / PR #86),
and the contract-drift tooling (#85 / PR #93) — which caught its own first
live drift: v0.4.0's embedded contracts lagged the amendments merged after
it, `status` reported the gap, and the report cleared at v0.5.0, exactly
the loop the design intends
([#103](https://github.com/radiusred/gh-codecrew/issues/103)).

## Protocol-discipline observations

- **The gatherer missed most of the milestone's spine.** `milestone close 5`
  gathered three records; #73 alone carries four Decisions. The three it
  missed are labeled `**Decision (…):**` — the parenthetical qualifier
  defeats the gatherer's `**Decision:**` matcher — and the PR-recorded
  deviations (#86, #95) live outside issue comments entirely. M3 and M4
  observed "judgment lives in plans and PR bodies"; this is the same gap in
  mechanical form, now with a precise reproduction. This document works from
  the full #73 record, not the gather.
- **A plan asserted delivery that hadn't happened.** #73's wrap-up plan
  claimed "fold-backs already merged carry the docs (dispatch recipe,
  counting rules, coexistence)" — two of the three were verifiably present;
  the Copilot coexistence statement was not, and QA's grep over the docs
  and full history proved it had never existed. The stale-prose class the
  review rounds kept catching in diffs, appearing at plan level, caught by
  the verdict gate instead.
- **The record miscounts itself in one place.** The QA-prep comment cites
  "#73 (thirteen findings, five Decisions)"; #73 carries thirteen numbered
  findings and four Decision comments — the milestone's fifth Decision is
  the gate resolution on #68. Trivial, but a synthesizer counting from the
  prep line alone would look for a Decision that is not there.
- **Experiment evidence can be deleted.** The cc-approval-probe repo that
  settled the counting question was a throwaway and is gone; the
  superseding Decision preserves the result and method, but the experiment
  is not re-runnable from the record, and the recorded discriminating
  experiment (org-vs-personal as co-factor) remains unrun. `milestone
  evidence` would have flagged a URL; a shorthand citation to a deleted
  repo is the class it cannot.
- **Declared harness and actual harness have drifted.** The routing table
  declares `harness: codex` for the reviewer seat, and the numberguess
  proving-ground dispatches ran under Codex — but this milestone's hub
  dispatches of checky ran on claude sub-agents (operator-reported at
  synthesis time; no record of the dispatching harness exists anywhere in
  the milestone's issues or PRs). The table is documented as declared
  intent, so this is an unreconciled intent/practice drift — either the
  declaration or the practice should move, and nothing currently records
  which harness actually ran a dispatch.
- **Supersession worked as designed, at every level it exists.** A Decision
  superseded a Decision (the counting pair — the first stands in the record
  with its reasoning; the second names exactly what it corrects); a verdict
  superseded a verdict (M5-R1, the remedy loop's second run under the
  enforced gate); reviews superseded reviews (every changes-requested round
  closed by a later approval naming what it verified). Latest-wins with the
  history intact, nothing edited away.
- **The gates held under delegation and against their own subjects.** The
  #68 gate was raised before any commit and resolved before the first code
  PR merged. The M2 reviewer sub-agent refused to approve while a gate was
  unresolved and refused to smear identities when its token failed. R4 was
  proven on the close that shipped it: the evidence verb ran against its
  own milestone, and the executed-docs obligation was first applied to the
  diff that introduced it. And `milestone close 5` refused `DOC_MISSING`
  until this document exists on the default branch — the refusal is in the
  gather output this document was synthesized from.
