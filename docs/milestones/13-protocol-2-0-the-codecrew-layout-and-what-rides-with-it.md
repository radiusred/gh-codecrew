# M13: Protocol 2.0: the .codecrew/ layout and what rides with it

Tracking issue: [#254](https://github.com/radiusred/gh-codecrew/issues/254) ·
Synthesized 2026-09-06 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #254's **nine** requirements as opened, **one struck
and none added** — the issue's `userContentEdits` query returns a count of
two, an original and a single edit at 11:31:54Z whose whole diff is the
removal of M13-R8 — leaving **eight**; its Gates section; the **two
fresh-context scan reports** attached to it as comments; its **two Decision
comments, one Deviation comment and three QA comments**; the ten delivery and
remedy task issues, all in this repository
([#255](https://github.com/radiusred/gh-codecrew/issues/255),
[#256](https://github.com/radiusred/gh-codecrew/issues/256),
[#257](https://github.com/radiusred/gh-codecrew/issues/257),
[#258](https://github.com/radiusred/gh-codecrew/issues/258),
[#259](https://github.com/radiusred/gh-codecrew/issues/259),
[#260](https://github.com/radiusred/gh-codecrew/issues/260),
[#261](https://github.com/radiusred/gh-codecrew/issues/261),
[#262](https://github.com/radiusred/gh-codecrew/issues/262),
[#285](https://github.com/radiusred/gh-codecrew/issues/285),
[#288](https://github.com/radiusred/gh-codecrew/issues/288)) plus
[#263](https://github.com/radiusred/gh-codecrew/issues/263), closed as moved;
their **ten merged PRs**; the **thirty-three Decision records and six
Deviation records** across the milestone issue and those task issues; the
**twenty-three review submissions** on the ten PRs — one approval dismissed
by a rebase, and two bodies that reached GitHub as the wrong text and were
edited after submission; the three adopted backlog captures the merges closed;
and the five captures the milestone filed and did not bundle. The house form
is [M12's document](12-v1-2-0-and-the-field-fixes-behind-it.md), with
[M11's](11-housekeeping.md) and
[M10's](10-protocol-bookkeeping-from-the-field.md) behind it; the record
standard is the one the reviewer set on their PRs,
[#232](https://github.com/radiusred/gh-codecrew/pull/232),
[#240](https://github.com/radiusred/gh-codecrew/pull/240) and
[#252](https://github.com/radiusred/gh-codecrew/pull/252) — verdict words
verbatim, no ordinal claim without every prior record linked, no judgement the
record does not carry, and a claim checked against the trail as it stands when
the record is written.

The gather from `gh codecrew milestone close 13` is again not part of the raw
material, for the reason
[M8's](8-a-product-home-page-for-codecrew-works.md),
[M9's](9-the-docs-at-codecrew-works.md),
[M10's](10-protocol-bookkeeping-from-the-field.md),
[M11's](11-housekeeping.md) and
[M12's](12-v1-2-0-and-the-field-fixes-behind-it.md) records all give:
`--dry-run` stops at the gate that counts tasks, naming the task that writes
this file. A binary built from `main` at
[c58a7df](https://github.com/radiusred/gh-codecrew/commit/c58a7df) — the merge
of this milestone's last delivery — prints:

```
gate milestone open: ok
gate no gate raised: ok
gate tasks closed: refused[OPEN_TASKS]: tasks not closed: radiusred/gh-codecrew#284 (ready)
gate requirements declared: not reached
gate QA verdicts: not reached
gate milestone document: not reached
dry run: nothing written — the live verb stops at the first refusal above
```

The released `v1.2.0` binary reaches no gate at all against this hub. Run at
the repo root, `status` and `roles show reviewer` both print
`codecrew: no .codecrew.yml found (not a CodeCrew repo?)` and exit 1 — a plain
error, with no `refused[CODE]:` in it, because a 1.x binary has no vocabulary
for a repository on the 2.0 layout; only `version` still works, printing
`v1.2.0 (protocol 1.0)` and exiting 0. It is looking for a root
`.codecrew.yml`, which M13-R1 moved. That is the milestone in one line, and it
is why the hub has been running on a local build since 13:06:47Z on
2026-09-06 — `version` on the operator's machine reports
`dev-c58a7df (protocol 2.0)`.

Every record below was read from its issue, PR, review, commit or timeline
directly.

**This PR adds the M13 ROADMAP row; it does not flip one** — the convention
M10-R1 introduced, which
[M10's record PR](https://github.com/radiusred/gh-codecrew/pull/232),
[M11's](https://github.com/radiusred/gh-codecrew/pull/240) and
[M12's](https://github.com/radiusred/gh-codecrew/pull/252) were written under.
M12-R1 shipped the verb that stopped writing the row at `milestone new`, and
#254 was the first milestone opened under it: the trail carries no Decision
discarding a row, because there was no row to discard. The three earlier
records each had one
([M10](10-protocol-bookkeeping-from-the-field.md#the-roadmap-row-belongs-to-the-doc-synthesizer-at-both-ends),
[M11](11-housekeeping.md#the-roadmap-row-discarded-again),
[M12](12-v1-2-0-and-the-field-fixes-behind-it.md#the-roadmap-row-discarded-for-what-the-operator-calls-the-last-time),
whose comment called itself "the last milestone opened under that binary").

## Goal and outcome

#254's goal starts from a collision. "CodeCrew writes its operational files to
the root of an adopter's repository — `roles/`, the pointer — where they
collide with real project layouts (Ansible's `roles/`, for one). The fix is a
layout move, which is a protocol major under SPEC §5, and a major is the one
moment a breaking change is cheap: about a dozen repos, one operator."

Two things follow from that sentence, and they shaped the whole milestone.
The first is that there is no compatibility shim anywhere: "a 2.0 binary
refuses the 1.x layout and names the one-shot migrate verb". The goal names
its reason — "the scar-tissue principle from `docs/gsd-vs-frontier-orchestration.md`"
— an essay in this repository whose argument is that most of an orchestration
framework's machinery is "portability scaffolding and scar tissue that a
frontier model mostly doesn't need, and it is genuinely expensive". A dual
read of both layouts would have been exactly that, and it was rejected in the
conversation that opened the milestone; the record carries the rejection in
the goal and in the adoption Decision, not as a quoted exchange.

The second is that if a major is the cheap moment, the question is what else
should ride it. #254 answers by commissioning two fresh-context scans of the
protocol's public surface — one Codex (`gpt-5.5`, high), one Claude (Opus) —
against one brief: what else would force a 3.0 within weeks if 2.0 shipped as
the layout move alone. Both are attached to the milestone issue in full
([Codex](https://github.com/radiusred/gh-codecrew/issues/254#issuecomment-5558889556),
[Claude](https://github.com/radiusred/gh-codecrew/issues/254#issuecomment-5558889638)),
ten findings each, and the operator's
[Decision](https://github.com/radiusred/gh-codecrew/issues/254#issuecomment-5558890102)
posted six seconds after the second says which were adopted, which were
blessed permanent, and which were filed for later. Requirements R1 through R7
are those adoptions — R1 is the layout move, and the Decision lists it among
them, because the scans widened it from `roles/` to every CodeCrew-owned file
including the pointer. R9 is two open captures riding along.

Nine requirements were present when the milestone opened at 11:22:12Z. The
nine tasks were created between 11:22:57Z and 11:23:21Z, twenty-four seconds
across the nine. Seven minutes and forty-two seconds later, at 11:31:55Z, the
operator's
[second Decision](https://github.com/radiusred/gh-codecrew/issues/254#issuecomment-5558926520)
struck M13-R8 — the v2.0.0 release — and moved it to M14; the body edit
followed one second earlier at 11:31:54Z, #263 was closed as moved at
11:31:56Z, and M14
([#269](https://github.com/radiusred/gh-codecrew/issues/269)) was opened at
11:31:59Z. **Nothing shipped in M13.** The eight remaining requirements were
delivered unreleased, on `main`, and
[#269](https://github.com/radiusred/gh-codecrew/issues/269) is where v2.0.0
and the fleet migration happen.

The eight delivery tasks and two remedy tasks, in merge order:

- **[#262](https://github.com/radiusred/gh-codecrew/issues/262) / PR
  [#274](https://github.com/radiusred/gh-codecrew/pull/274)** — M13-R9. Three
  commits, approved first round, merged 11:57:44Z
  ([b48ffb6](https://github.com/radiusred/gh-codecrew/commit/b48ffb6)). Two
  Decisions. Adopts [#253](https://github.com/radiusred/gh-codecrew/issues/253)
  and [#250](https://github.com/radiusred/gh-codecrew/issues/250), and closes
  both.
- **[#260](https://github.com/radiusred/gh-codecrew/issues/260) / PR
  [#275](https://github.com/radiusred/gh-codecrew/pull/275)** — M13-R6. Four
  commits, two review rounds, merged 12:23:09Z
  ([b49adb4](https://github.com/radiusred/gh-codecrew/commit/b49adb4)). Four
  Decisions and one Deviation, the last of which records that there was none.
- **[#258](https://github.com/radiusred/gh-codecrew/issues/258) / PR
  [#276](https://github.com/radiusred/gh-codecrew/pull/276)** — M13-R4. Five
  commits, four review rounds, merged 12:50:17Z
  ([07e4d0c](https://github.com/radiusred/gh-codecrew/commit/07e4d0c)). Two
  Decisions. The milestone's dismissed approval is here.
- **[#255](https://github.com/radiusred/gh-codecrew/issues/255) / PR
  [#277](https://github.com/radiusred/gh-codecrew/pull/277)** — M13-R1, the
  layout itself. Three commits, two review rounds, merged 13:06:47Z
  ([1c989ed](https://github.com/radiusred/gh-codecrew/commit/1c989ed)). Three
  Decisions and two Deviations.
- **[#257](https://github.com/radiusred/gh-codecrew/issues/257) / PR
  [#278](https://github.com/radiusred/gh-codecrew/pull/278)** — M13-R3. Six
  commits, three review rounds, merged 13:40:43Z
  ([c088717](https://github.com/radiusred/gh-codecrew/commit/c088717)). Four
  Decisions and one Deviation.
- **[#259](https://github.com/radiusred/gh-codecrew/issues/259) / PR
  [#279](https://github.com/radiusred/gh-codecrew/pull/279)** — M13-R5. Four
  commits, two review rounds, merged 13:56:20Z
  ([1004651](https://github.com/radiusred/gh-codecrew/commit/1004651)). Two
  Decisions and one Deviation.
- **[#256](https://github.com/radiusred/gh-codecrew/issues/256) / PR
  [#280](https://github.com/radiusred/gh-codecrew/pull/280)** — M13-R2,
  `gh codecrew migrate`. Seven commits, three review rounds, merged 14:22:48Z
  ([5fbf14f](https://github.com/radiusred/gh-codecrew/commit/5fbf14f)). Nine
  Decisions, the milestone's largest set on one task. Closes
  [#281](https://github.com/radiusred/gh-codecrew/issues/281), the capture
  filed against its own review.
- **[#261](https://github.com/radiusred/gh-codecrew/issues/261) / PR
  [#282](https://github.com/radiusred/gh-codecrew/pull/282)** — M13-R7. Six
  commits, two review rounds, merged 14:50:15Z
  ([3811ba4](https://github.com/radiusred/gh-codecrew/commit/3811ba4)). Two
  Decisions.
- **[#285](https://github.com/radiusred/gh-codecrew/issues/285) / PR
  [#286](https://github.com/radiusred/gh-codecrew/pull/286)** — the first R6
  remedy, opened from a QA finding. Two commits, two review rounds, merged
  15:25:01Z ([6c7ece5](https://github.com/radiusred/gh-codecrew/commit/6c7ece5)).
  One Decision.
- **[#288](https://github.com/radiusred/gh-codecrew/issues/288) / PR
  [#289](https://github.com/radiusred/gh-codecrew/pull/289)** — the second R6
  remedy, from the next QA finding. Two commits, two review rounds, merged
  15:58:34Z ([c58a7df](https://github.com/radiusred/gh-codecrew/commit/c58a7df)).
  Two Decisions.

From the milestone opening to the last delivery merging: four hours,
thirty-six minutes and twenty-two seconds. Eight requirements, eight
`satisfied` verdicts standing — one of them reached on the third QA round,
after two remedies.

What the hub looks like afterwards. Every CodeCrew-owned operational file is
under `.codecrew/`: the pointer at `.codecrew/config.yml` declaring
`codecrew: "2.0"` with typed identities, the five contracts and their
extensions under `.codecrew/roles/`, and the agent instructions at
`.codecrew/AGENTS.md` with a two-line root `AGENTS.md` pointing at them.
`ROADMAP.md`, `docs/milestones/`, `AGENTS.md` and `CLAUDE.md` stayed where
they were, by decision. A repo still on the 1.x layout refuses
`LAYOUT_LEGACY` naming `gh codecrew migrate`, the verb M13-R2 added.
`PROTOCOL_MISMATCH` keeps only its newer-than-this-binary meaning. A spoke
that cannot read its hub's table refuses `HUB_UNREADABLE` rather than
resolving every seat to the operator. SPEC §10 carries a table of all
**forty-two** refusal codes, held to the source by a test, and SPEC §6 states
the exit-code contract: `0` on success, `1` on every failure, the
`refused[CODE]: detail` line on stderr as the machine channel. Nothing of
this is in a release.

## The two scans, and what the operator did with them

The scans are the milestone's unusual raw material and they are worth reading
as a pair. Both were run against `main` at
[3f4de53](https://github.com/radiusred/gh-codecrew/commit/3f4de53), v1.2.0 /
protocol 1.0, on the same brief, in fresh context, by different models. Each
produced ten ranked findings with evidence down to file and line, a
"not breaking, do not bundle" section, and a verdict.

They agreed on the shape of the answer and disagreed about almost nothing.
Both said the layout move alone is not a sound 2.0 — Codex: "Layout-only 2.0
is sound only if 'layout' means every CodeCrew-owned filesystem surface";
Claude: "The layout change alone is not a sound 2.0. Two things must ride with
it, and both are consequences of the move rather than opportunistic
additions." Both independently found the pointer's own location (Codex 1,
Claude 5), the migration verb's need for its own refusal code and its own
exemption from the pointer check (Codex 3, Claude 3), and the requirement-ID
grammar being unenforced (Codex 5, Claude 6–7 in the same family). Where they
differed was in emphasis: Codex pressed hardest on namespace choices that
would be expensive to change later — the `M<n>:` title, `task/` branches, the
`cc:` label prefix — while Claude pressed hardest on two behaviours that fail
open, the missing-`codecrew:`-field rule and the unreadable hub table.

The operator's
[Decision](https://github.com/radiusred/gh-codecrew/issues/254#issuecomment-5558890102)
sorts every finding into one of three piles, and names the finding numbers as
it goes.

**Adopted into requirements.** The whole CodeCrew-owned filesystem surface
moves, pointer included (Codex 1 and 2, Claude 5) — R1. `migrate` is
conservative and has its own refusal codes, distinct from
`PROTOCOL_MISMATCH` (Codex 3, Claude 3) — R2. The entry point (Codex 10) —
R3. Typed identities (Claude 4) — R4. Routing fails closed, with a hub/spoke
skew refused "and an unreachable GitHub named by its own code so an offline
operator is never told the hub table is missing" (Claude 2 and 8) — R5. The
record grammar (Claude 6 and 7, Codex 5) — R6. The 1.0 shims deleted and the
exit-code contract and refusal-code table written into SPEC (Claude 1, 9 and
10; Codex 4) — R7.

**Blessed as permanent, no change.** The `M<n>:` title as milestone identity
(Codex 6), `task/<n>-<slug>` branches (Codex 7), the `cc:` label prefix
(Codex 8), and `ROADMAP.md` at the repository root with `docs/milestones/`
under `docs/` (Claude 5). The reason given for the last: "the records and
roadmap are the human-facing product of the protocol and stay where readers
look." Codex 9, a structured provenance stamp, is not adopted — "the stamp is
stripped by prefix and moves with the files; a schema for it can ship in any
minor." A blessing is a decision as much as a change is, and this is the only
place on the trail where these four are settled rather than merely unchanged.

**Filed, not bundled.** Five captures, opened between 11:24:17Z and 11:24:21Z
— seconds *before* the Decision that names them:
[#264](https://github.com/radiusred/gh-codecrew/issues/264) (first-page reads
drop the newest comments, verdicts and sub-issues past 100),
[#265](https://github.com/radiusred/gh-codecrew/issues/265) (routing-table
logins compared case-sensitively),
[#266](https://github.com/radiusred/gh-codecrew/issues/266) (`status` reports
contract drift but never contract absence),
[#267](https://github.com/radiusred/gh-codecrew/issues/267) (no verb creates
the `cc:` labels) and
[#268](https://github.com/radiusred/gh-codecrew/issues/268) (the offline
workflow). All five came out of the scans' "not breaking" sections and the
offline question R5 raised.

One clause of the Decision is a grammar ruling rather than an adoption, and
every later task reads it: "`~` covers the operator and any agent session
acting under the operator's own auth — the grammar types the GitHub
principal, not who is at the keyboard; a main agent with its own App is
`app:<slug>` on the coordinator row."

The Decision also fixes the milestone's sequencing, and the trail follows it
exactly: "after #255 merges the installed v1.2.0 extension cannot read this
hub; the operator reinstalls the extension from the local checkout
(`gh extension install .`) and the rest of the milestone runs on it. #256
(migrate) lands after #255 and #258. #261 (the refusal-code table) lands
last among the code tasks. #263 (release) is last." #263 was closed seven
minutes later; the other four orderings held.

## Decisions

Thirty-three Decision records across the milestone issue and the ten task
issues. Two are the operator's, on #254; the other thirty-one were written by
the seat doing the work, `radiusred-cody[bot]`, which started all ten
delivery and remedy tasks. They group by the question each task had to settle.

### The refusal that names a verb that does not exist yet

M13-R1 shipped before M13-R2, so the binary that first refuses a 1.x repo did
so by naming `gh codecrew migrate` before the verb existed. The
[Decision](https://github.com/radiusred/gh-codecrew/issues/255#issuecomment-5559035129)
takes that deliberately: "The refusal that a 1.x repo actually meets is this
one, so it is the only place the operator can be told what to do; a refusal
that names no fix is a dead end, and the alternative — shipping the refusal
mute and editing its text a task later — leaves a window in which the only
thing the binary says about a 1.x repo is that it will not read it."
**Trade-off:** "between #255 merging and #256 merging, a local build refuses a
1.x repo by naming a verb it does not yet have. The window is one task long,
inside one unreleased milestone, and the only repo affected is this hub."
**Rejected:** a placeholder `migrate` that prints "not yet implemented" — "a
verb in the surface is a promise under SPEC §10, and one that exists to be
empty is worse than a message pointing at the one that is coming."

The same comment settles what the two directions of `PROTOCOL_MISMATCH` say,
and neither says "update the pointer": "`codecrew: "1.0"` is a statement about
the repo's layout and conventions, not a preference — editing it by hand would
make the file lie about the tree beside it." So a pointer ahead of the binary
reads "upgrade the extension" and one behind it is sent to `migrate`. An
unparseable major takes the upgrade wording rather than a migration "that
would not know what to do with it".

### Where the 1.x layout is looked for, after a review found a false refusal

The reviewer's round-one finding on PR #277 was that `config.Load` checked for
the 1.x layout at every level of its upward walk, so a valid 2.0 repo refused
`LAYOUT_LEGACY` the moment a verb ran from a directory holding an ordinary
nested `roles/` — a playbook's, say. The
[Decision](https://github.com/radiusred/gh-codecrew/issues/255#issuecomment-5559157411)
that answers it is blunt about the stakes: "That is precisely the collision
protocol 2.0 exists to end, so the pointer must win wherever it sits above the
caller." The rule: walk upward for the nearest ancestor holding a `.git`
entry — "found by the same filesystem walk that finds the pointer, no `git`
subprocess" — and fall back to the starting directory outside a repository.
**Trade-off:** "a 1.x checkout that is not a git repository refuses only when
the verb runs at its top level", the position `init` is already in.
**Rejected:** special-casing the five contract names harder — "the false
positive is a *directory* collision, not a filename one, and no name list
fixes a rule applied at the wrong level" — and shelling out to
`git rev-parse --show-toplevel`, "which would put a subprocess on the path of
every verb".

### `migrate` refuses rather than half-moves, and reads `roles/` only when it is CodeCrew's

M13-R2's own text says two things about a hub's `roles/` that read
differently: move the recognised names "leaving anything else there and
reporting it", and "refuse when `roles/` holds an unrecognised entry rather
than guessing". The
[Decision](https://github.com/radiusred/gh-codecrew/issues/256#issuecomment-5559499209)
picks the refusal, and explains what makes it safe: "`roles/` is opened only
when at least one of the ten names is in it. A project's own `roles/` —
Ansible's, the collision that motivated the whole layout move — holds none, so
it is never read, never moved and never a refusal." Once one of the ten is
there the directory is CodeCrew's, and every entry must be one of the ten.
**Trade-off:** "a hub whose `roles/` genuinely mixes CodeCrew's contracts with
the project's own files cannot be migrated by the verb at all"; the refusal
names both sets "so the separation is a copy-paste rather than an
investigation". **Rejected:** `ROLES_UNRECOGNISED` as the code — "it names the
input rather than the diagnosis" — and a `--force`/`--only <paths>` escape
hatch, "a second way to be wrong about the same directory". A second Decision
in the same comment keeps the ten recognised names identical in a hub and a
spoke, so a spoke holding a stale `roles/qa.md` moves it instead of refusing
over "a file that is unmistakably CodeCrew's own".

### Four migration codes, and a rule for when a code says nothing

The
[Decision](https://github.com/radiusred/gh-codecrew/issues/256#issuecomment-5559500777)
names `BOTH_LAYOUTS`, `FOREIGN_ROLES_DIR`, `IDENTITY_UNRESOLVED` and
`MIGRATION_UNSUPPORTED`, and states the test it applied to each: "One code,
one meaning — *migrate cannot type this row, so you must* — with the detail
saying which of the four it was. The fix is the same in every case, which is
the test for whether a second code would say anything." **Trade-off:** "four
codes are four permanent promises under SPEC §10 for a verb an adopter runs
once. They are worth it because the once is the *first* contact with 2.0 and
every one of these outcomes needs a different human act." **Rejected:**
`PROTOCOL_MISMATCH` for the version case, because the scans' finding "is
precisely that a migration refusal must not share a code with 'upgrade the
extension'"; `IDENTITY_AMBIGUOUS` as a fifth code, because "a code an agent
cannot act on differently is a code that says nothing"; and codes for the
three non-protocol stops, which stay plain errors because "none of them is a
gate an agent can clear by acting on the code".

That rule was applied twice more in the same task, both times to *not* add a
code. A destination that already exists became `BOTH_LAYOUTS` rather than a
new `DESTINATION_EXISTS` — "the condition is the same fact the pointer check
already catches one level up ... and the operator's act is identical"
([Decision](https://github.com/radiusred/gh-codecrew/issues/256#issuecomment-5559594973))
— and a 1.x spoke pointer carrying a `roles:` block reuses `SPOKE_ROUTING`
from #259, "so the count stays at forty-two"
([Decision](https://github.com/radiusred/gh-codecrew/issues/256#issuecomment-5559804307)).

### The 404 that is the first half of a question

The task's rule for typing a bare identity was: probe `users/<login>`, `Bot`
means `app:`, `User` means `user:`, a 404 refuses. The
[Decision](https://github.com/radiusred/gh-codecrew/issues/256#issuecomment-5559502898)
records why that would have failed on the commonest row there is: "A 1.0 table
wrote an App as its bare slug — `identity: radiusred-cody` ... but a GitHub
App's account login *is* `radiusred-cody[bot]`, so `users/radiusred-cody` is a
404 and every App-held seat in every 1.0 hub would have refused. This hub's
own pre-#276 table is four such rows." So the 404 became "the first half of
the question": probe the value, then `<value>[bot]`. The second probe also
exposes the case the single probe would have silently mis-typed — a human
login `foo` beside an App `foo[bot]` — which is "exactly the ambiguity M13-R4
exists to end" and refuses `IDENTITY_UNRESOLVED` naming both. **Trade-off:**
two API calls per bare row, "and a migration cannot be run fully offline.
Migrate is a once-per-repo verb whose whole job here is asking GitHub a
question nothing local can answer, so the cost is the feature." **Rejected:**
guessing from the value's shape; probing `[bot]` first, which "would type a
real human whose login collides with an App slug as the App"; and accepting an
`Organization`, which "is not a principal that holds a role under SPEC §5".

The same comment settles how the pointer is rewritten: a `yaml.Node`
round-trip rather than re-marshalling a struct, because "the pointer is a file
a human maintains" and re-marshalling "would silently delete all of it". Blank
lines, which `yaml.v3` drops, are carried across as a `#codecrew-blank-line`
sentinel comment. **Trade-off:** "a round-trip is not byte-preserving" —
inner padding in flow mappings, indentation normalised — against "regex
surgery on lines", which "is fragile across the block, flow, quoted and
multi-document shapes a hand-edited pointer can take". The commit is local and
unpushed "precisely so the operator reads the diff before it goes anywhere".

### The entry point carries both forms, and a spoke gets one too

M13-R3 left the shape of the root `AGENTS.md` pointer to the implementer, and
the
[Decision](https://github.com/radiusred/gh-codecrew/issues/257#issuecomment-5559462995)
writes both a sentence and an `@` import. "The `@` import is Claude Code's
mechanism ... it means nothing to a harness that reads `AGENTS.md` as plain
markdown, which needs the path in a sentence it can follow. Neither form alone
reaches every harness, and the failure mode of picking one is invisible — an
agent that starts blind, exactly the
[#141](https://github.com/radiusred/gh-codecrew/issues/141) bug." The block is
a single constant, `entryPointLines`, which is both the scaffold and, byte for
byte, what `init` prints for a kept root file; a test asserts the two cannot
drift.

A [second Decision](https://github.com/radiusred/gh-codecrew/issues/257#issuecomment-5559463082)
settles the open question the Codex scan's finding 10 asked to have settled —
"define whether spokes get an entrypoint or remain pointer-only" — as yes, in
all three files. "The dispatch story is identical in hub and spoke ... so the
file a fresh session reads must exist in both. Writing `.codecrew/AGENTS.md`
but not `CLAUDE.md` would have reproduced #141 in every spoke." What stays
hub-only is what the hub owns: `ROADMAP.md`, the contracts and their
extensions, "a test asserts it does not" grow a second copy in a spoke.

A [third](https://github.com/radiusred/gh-codecrew/issues/257#issuecomment-5559463171)
points the block printed for a kept `CLAUDE.md` straight at
`.codecrew/AGENTS.md` rather than at `AGENTS.md`: "One block, correct in every
combination of kept files, beats two blocks whose correctness depends on which
other file exists." A
[fourth](https://github.com/radiusred/gh-codecrew/issues/257#issuecomment-5559463257)
declines to add a root `CLAUDE.md` to this hub, which the Codex scan noted the
dogfood repo lacks: "adding a root file to this repository is a change the
operator should see as its own decision rather than as a side effect of a
scaffold task."

### The entry point rides the migration, on the operator's call

The one place in the milestone where a seat raised a scope question rather
than deciding it. A migrated 1.x hub was landing on the 2.0 layout with no
`.codecrew/AGENTS.md`, "its root `AGENTS.md` still carrying the old
instructions and naming `roles/` and `.codecrew.yml`, paths that no longer
exist". The implementer flagged it on PR #280; the
[Decision](https://github.com/radiusred/gh-codecrew/issues/256#issuecomment-5559741725)
records the operator's answer and the scope reading that makes it R2's work
rather than a new requirement: "R2 moves a protocol 1.0 repo to the 2.0
layout, and R3 defines what the 2.0 layout's entry point *is* ... A migration
that produces every other 2.0 file and not that one has not finished the move
it names." **Trade-off:** "migrate now writes a file as well as moving files,
which widens what a one-shot verb touches. It is bounded to one path that
CodeCrew owns and that the repo does not already have." **Rejected:**
rewriting the root `AGENTS.md` to the 2.0 pointer scaffold — "it would be the
cleanest possible migration and it is exactly what #257 decided against for
`init`: the root entry point is the adopter's file ... and a verb that
rewrites prose is a verb nobody runs twice" — and leaving the file to a
follow-up `init` run, which "makes the migration a two-command ritual whose
second half is easy to skip, and the skip is silent".

### `IDENTITY_UNTYPED`, and a flag whose emptiness is the answer

The
[Decision](https://github.com/radiusred/gh-codecrew/issues/258#issuecomment-5559004128)
names the code for a bare routing value and puts it where the pointer is
loaded: it "names the condition (the value carries no type), not a guess at
the fix, and sits in the family with `PROTOCOL_MISMATCH`". **Trade-off:** the
condition is detected in `internal/config`, which cannot import the CLI's
refusal type, so `config.Parse` returns a typed error and `loadConfig` maps it
— "exactly the split `config.Compatible` / `PROTOCOL_MISMATCH` already uses".
**Rejected:** `IDENTITY_INVALID` ("too broad"), `BARE_IDENTITY` ("describes
the input, not the failure"), and refusing inside `UnmarshalYAML`, which "does
not know which role row it is decoding, so the detail could not name it".

Typing the value creates a second problem the
[next Decision](https://github.com/radiusred/gh-codecrew/issues/258#issuecomment-5559005094)
answers: `gh pr create --reviewer` wants a bare handle, and
`role <name>` now prints `app:radiusred-checky`. `role <name> --login` prints
the review-requestable handle — the login for `user:`, `<org>/<slug>` for
`team:` — and nothing at all for an App or the operator, so the implementer
contract's rule becomes one line whose emptiness is the branch.
**Trade-off:** "a script cannot tell 'no handle' from 'the verb failed' by
exit code alone — but it can, from `role <name>` without the flag, and the
empty string is the branch the contract actually wants to take. A nonzero exit
would make the common, correct case look like an error." **Rejected:**
printing a team's bare slug, because `gh pr create --reviewer` "splits a team
from a user on the slash and needs `org/slug`"; keeping `role` printing the
bare value, which "throws away the kind the whole requirement exists to
expose"; and a separate verb.

### Routing is resolved once, eagerly, and by topology rather than emptiness

M13-R5's whole subject is a failure that was invisible. PR #279's description
states it plainly: a spoke fetches the hub's table, "and on any failure it
kept the local pointer's table — which in a spoke is empty. An empty table
resolves every seat to `~`, and `~` is the branch where `task finish`'s
holder-review gate falls through to 'any non-author approved' and
`milestone close` counts a QA verdict from any commenter. The two gates the
protocol exists to enforce, disabled by a 404 nobody sees." The `.codecrew/`
move guarantees that 404 for a whole migration window, in both directions of
skew, which is why the fix rides 2.0 rather than a minor.

The
[first Decision](https://github.com/radiusred/gh-codecrew/issues/259#issuecomment-5559469157)
resolves the table once, eagerly, in `load()` — "a `ctx` cannot exist without
one". **Trade-off:** one extra `FileContent` call per spoke run, "on verbs
that may never resolve a role. That is cheap: `load()` already calls
`gh repo view`, so no verb that builds a `ctx` was ever usable offline."
**Rejected:** threading `(bool, error)` through every role predicate, because
"it leaves the same shape of mistake available at every one of them — a `bool`
returned alongside an error a caller may drop. The structural fix is the
smaller and the safer one."

The
[second](https://github.com/radiusred/gh-codecrew/issues/259#issuecomment-5559469243)
chooses which table governs by topology, not by the local one being empty:
`hub: self` reads the local pointer from disk, `hub: owner/repo` always takes
the hub's. The reason it matters here is stated in the comment: "without it,
'routing fails closed' would mean a hub with no routing table refuses every
verb the moment the network drops — a hub declaring no table is legitimately
`~` everywhere (M13-R5 says so explicitly), and it must stay that way with the
network down." **Rejected:** keeping the emptiness test, which "in a hub
turned a local absence into a network dependency".

### One placement rule for the label family, and one code stripper

M13-R6 reclassifies text already written on GitHub, which is why it rides a
major. The
[Decision](https://github.com/radiusred/gh-codecrew/issues/260#issuecomment-5559036839)
states the gate rule precisely: comments in order, paragraph by paragraph, on
the same split `ExtractRecords` uses; a `**Gate raised:**` paragraph opens a
gate; a `**Gate resolved:**` paragraph "closes every gate opened before it and
still open — one comment may still answer several questions, which the old
trailing-index rule got right — and never a gate opened after it"; nothing
else closes a gate. **Trade-off:** the placement rule is now uniform across
`**Decision:**`, `**Deviation:**`, `**Gate raised:**` and `**Gate
resolved:**`, "at the cost of reclassifying text already written on GitHub ...
A gate the old rule cleared with an unrelated Decision now blocks
`task finish` until a human writes the answer down; that is the failure
direction a human can fix in one comment, where the reverse (a gate silently
cleared) is the one that loses the record." **Rejected:** matching each gate
to the *next* resolution one-for-one, which "breaks the case SPEC §8 and the
old code both support", and keeping a bare `**Decision:**` as a resolution
"when it is the only comment after the gate", "a special case with no wording
in the SPEC behind it".

A
[second Decision](https://github.com/radiusred/gh-codecrew/issues/260#issuecomment-5559038681)
moves `stripCode` out of `internal/cli/evidence.go` into the tracker as
`tracker.StripCode`, so "the citation-reading rule and the verdict-reading
rule cannot drift apart again — the drift Claude's finding 7 named". It also
states the boundary it does *not* cross, so the reviewer could weigh it:
`ParseVerdicts` reads through `StripCode`; `ExtractRecords` and
`UnresolvedGates` do not, "because they are one label family read by one
placement rule, and stripping code under them would reclassify record text
this task was not asked to touch". The failure that leaves open is named and
its direction argued: a fenced block containing a line opening
`**Gate raised:**` reads as a raised gate, "which blocks `task finish` until a
human answers or edits, rather than losing a gate". **Trade-off:** "a verdict
inside a fence is content while a gate label inside a fence is a label, which
is an inconsistency in the same binary. It is the smaller of the two."

A [third](https://github.com/radiusred/gh-codecrew/issues/260#issuecomment-5559037700)
splits the new `REQUIREMENT_ID_MISMATCH` by verb: `status` prints a line and
carries on, `milestone close` and `milestone evidence` refuse. "A refusal
there would hide the tasks of every *other* milestone behind one milestone's
malformed body, and `status` gates nothing, so there is nothing for the
refusal to protect." **Trade-off:** "the same condition has two shapes in the
CLI's output, which a reader could mistake for two rules" — against the
precedent that `status` already does this for the empty-Requirements case.

### The missing `codecrew:` field is kept, and the "0.1" acceptance is not

The Claude scan's finding 1 proposed inverting the "a missing protocol field
is assumed current" rule so that it would refuse with the migration code. The
[Decision](https://github.com/radiusred/gh-codecrew/issues/261#issuecomment-5559905739)
keeps it, and the argument turns on what R1 changed: "That reasoning held
while the pointer sat at the repository root, where the same file could belong
to either major and the field was the only thing distinguishing them. Under
2.0 the layout itself answers the question: the file being read at all means
it is at `.codecrew/config.yml`, and a repo still on 1.x has no such file."
**Rejected:** inverting to a refusal, which "would refuse the one pointer
shape that cannot be lying about its layout: a hand-written two-line spoke
pointer under `.codecrew/`". The `"0.1"` acceptance beside it is deleted —
"which is the part that was genuinely scar tissue".

This is the one adopted scan finding the implementing task declined to
implement as written, and the record says so at the moment it happened rather
than after.

### One refusal-code catalogue, and a test that keeps it honest

`docs/introduction.md` carried a full second catalogue of the same codes. The
[Decision](https://github.com/radiusred/gh-codecrew/issues/261#issuecomment-5559925473)
makes SPEC §10's table the only one: "§10 already promises that a code's
meaning is stable within a major, and the list of codes that promise covers
belongs beside the promise rather than in a page a reader reaches through two
links." **Trade-off:** the introduction's per-code prose carried "what to do
about it" wording the table's one-line meanings do not — "that guidance is not
lost: each refusal's own detail string carries it, and the detail is what a
human reads at the moment it matters". **Rejected:** keeping both lists with a
test holding them to the same set: "The set was never the part that drifted —
the meanings were, and three tasks in this milestone (#256, #259 and this one)
each had to write the same code twice in two different registers."
`internal/cli/refusals_test.go` now holds the table, the stated count and the
`refuse("CODE"` sites to the same set.

### The indented code block, twice

The two remedy tasks each carry a Decision, and together they are a small
lesson in reading a specification. #285's
[Decision](https://github.com/radiusred/gh-codecrew/issues/285#issuecomment-5560134695)
states the indented-block rule measured from column 0 — four columns or more,
a tab advancing to the next multiple of four, opening "when the line above it
is blank or it is the first line of the text". **Trade-off:** "no block-structure
parser and no list-nesting state in a scanner whose whole job is to blank code
out of record text, at the cost of one known misread". **Rejected:** tracking
each list item's content column, "a Markdown block parser ... it would buy
correctness only for a shape the contracts do not ask anyone to write." The
comment names what it deliberately does not handle, which is how the next
finding came to be judged fairly.

#288's
[Decision](https://github.com/radiusred/gh-codecrew/issues/288#issuecomment-5560281160)
corrects the opener half: "an indented code block opens wherever it would not
interrupt a paragraph, which is CommonMark 4.4's only restriction on the
form." The `afterBlank` flag becomes `canOpen`, set by the start of the text,
a blank line, an ATX heading, a thematic break, and every line of a fenced
block. **Trade-off:** "the opener list is enumerated rather than derived from a
block parser, so a block-level construct nobody writes in a record — an HTML
block, a link reference definition, a blockquote's contents — does not open a
block on the line after it. Those cost a missed strip in a shape the contracts
never ask for, never a lost verdict." **Rejected:** distinguishing a setext
heading underline from a thematic break, since "both end the paragraph, so
both open a block". The same comment restates #285's column-0 boundary
unchanged, with its reason: "every form the contracts prescribe is a column-0
form ... A rule that reads column 0 therefore cannot lose one of them."

### The two smaller shapes: `status` between milestones, `evidence` after a close

M13-R9's
[first Decision](https://github.com/radiusred/gh-codecrew/issues/262#issuecomment-5559005315)
resolves a closed milestone through `MilestoneIssues`, the `state=all` listing
already in the `Tracker` interface, "rather than a new tracker method", and
learns the state from one `Task` read. **Trade-off:** one extra API call per
run "against a `Tracker` interface left untouched while #260 works in the same
package". **Rejected:** widening `Milestone`/`TitledIssue` or adding an
`AllMilestones` method — "both change the interface #260 is editing, for a
fact one existing method already answers" — and reading `OpenMilestones` first
with a fallback, "two listings to say what one says, and leaves the open path
able to drift from the closed one". The
[second](https://github.com/radiusred/gh-codecrew/issues/262#issuecomment-5559006238)
records two output calls, including why `status` with no open milestone prints
no `gates raised: none` line: "A `gates raised: none` line under an empty board
would be the verb asserting something it did not look for."

## Deviations

Six Deviation records. Five record a departure; the sixth, on
[#260](https://github.com/radiusred/gh-codecrew/issues/260#issuecomment-5559036839),
uses the label to say there was none — "none from the plan; this is the rule
the task asked to be defined precisely."

**[The reviewer seat ran on a different harness for part of the milestone](https://github.com/radiusred/gh-codecrew/issues/254#issuecomment-5559253147).**
The operator's Deviation, posted at 12:35:48Z: "the reviewer seat's declared
harness (codex, `gpt-5.5`) hit its usage limit at 13:35 BST on 2026-09-06 with
a reset three hours out." **Why:** "rather than idle two approved-or-nearly-approved
PRs for that window, the coordination layer runs the pending reviewer rounds
as fresh Claude (Opus) sessions under the same reviewer contract and the same
App identity (`radiusred-checky`) — the protocol's floor for independence is a
clean-context session holding the contract, the task and the diff; the codex
harness is the routing table's advisory preference for de-correlated
judgment." It names the rounds it expected to affect — "#276 round three (a
rebase check after codex approved the substance in round two), and #277 round
two if the quota is still out when it is rebased" — and says "Codex resumes
for every round after the reset." The reset was around 15:20 BST; PR #280's
third round, submitted at 14:21:32Z (15:21 BST), is the first round after it.

This is the distinction the protocol actually makes, stated on the record: the
routing table's `harness` and `model` fields are the coordination layer's
preference, while the gate is the contract, the identity and the clean
context. All twenty-three review submissions in the milestone are by
`radiusred-checky[bot]`, and the reviewer verified its own App ID against
`gh api /apps/radiusred-checky --jq .id` in the reviews of six of the ten PRs
— #275, #276, #277, #282, #286 and #289.

**[The docs pass touched two files the task's list did not name](https://github.com/radiusred/gh-codecrew/issues/255#issuecomment-5559086300).**
`docs/first-milestone.md`, "not in the task's doc list, but it describes what
`init` writes and names the pointer and the contracts directory throughout.
Left alone it would have told a new adopter to look for files the scaffold no
longer writes." And one path mention inside #262's own CHANGELOG entry and
SPEC row. **Why:** "the instruction was 'path mentions only, where they
describe the CURRENT layout'. Both are current-layout descriptions; neither is
a record of what happened."

**[SPEC §10's "github.com only in 1.0" became "in 1.0 and 2.0"](https://github.com/radiusred/gh-codecrew/issues/255#issuecomment-5559086300).**
**Why:** "§10 was being edited for the 2.0 paragraph anyway, and the sentence
is a live statement about what the verbs support, not a historical note."

The same comment carries a note the seat is explicit is *not* a deviation: the
`//go:embed` directive names the five contracts under `.codecrew/roles/` with
no workaround, because "Go excludes dot-prefixed names only when a pattern
names a directory; explicitly-named files embed, verified against the
toolchain before the move and by the build after it."

**[`init` prints the lines only when the kept root file does not already reach the instructions](https://github.com/radiusred/gh-codecrew/issues/257#issuecomment-5559565143).**
M13-R3 and the task's Plan both say `init` prints them "whenever a root
`AGENTS.md` or `CLAUDE.md` already exists". **Why:** "taken unconditionally,
the requirement makes `init` lie on its most common path. A rerun in a repo
`init` itself scaffolded keeps the pointer files it wrote and would print
`action needed` about them, which is false ... It also contradicts the
idempotency SPEC §6's `init` row promises in the same diff." The narrowing was
raised by the reviewer as a blocking finding in
[round one](https://github.com/radiusred/gh-codecrew/pull/278#pullrequestreview-5125429582),
fixed in
[dc64a2b](https://github.com/radiusred/gh-codecrew/commit/dc64a2b), and
re-verified in
[round two](https://github.com/radiusred/gh-codecrew/pull/278#pullrequestreview-5125452597),
"which also asked for this labelled record". The Deviation also records the
limitation accepted with it: `reachesInstructions` is a substring test, so it
reads "names the path", not "links to it", and "the failure mode is a missing
nudge, never a wrong file or an overwrite".

**[A plan step was dropped because reading the file made it vacuous](https://github.com/radiusred/gh-codecrew/issues/259#issuecomment-5559572170).**
Plan step 5 promised a correction to `docs/platform-interop.md`'s "advisory"
paragraph. **Why:** "on reading it, the plan step was vacuous ... that is a
statement about *dispatch*, not about the fetch, and it is exactly the one
sense of 'advisory' that the amended SPEC §5 preserves. Correcting it would
have made the doc wrong. Recorded rather than left silent because an
undeclared drop of a plan step is a finding whether or not the outcome is
right." The reviewer had reached the same conclusion independently and raised
it as a non-blocking finding; the record notes that too.

PRs #274, #276, #280, #282, #286 and #289 state "no deviations" in their
descriptions, and none is recorded on their task issues.

## The gates

**#254's Gates section was left as the scaffold's placeholder** — the template
line "_What "done" means beyond CI: e2e suites, manual UAT, sign-offs._",
unedited. [M12's record](12-v1-2-0-and-the-field-fixes-behind-it.md#the-gates)
says the same of #241, [M11's](11-housekeeping.md#the-gates) of #233,
[M10's](10-protocol-bookkeeping-from-the-field.md#the-gates) of #207,
[M9's](9-the-docs-at-codecrew-works.md#the-gates) of #201 and
[M8's](8-a-product-home-page-for-codecrew-works.md#the-gates) of #196. No
`**Gate raised:**` or `**Gate resolved:**` record exists anywhere in this
milestone, no `cc:needs-decision` label was applied to any issue in the
repository, and `checkpoint` was not used — on a task or on the milestone
issue.

That is worth stating twice over here, because M13-R6's subject is what the
gate grammar means and M13-R2's `migrate` is the verb an adopter meets first.
Both were verdicted against cases the milestone did not produce: the gate
rules from tests over the tracker's real M5 comment corpus, and `migrate` from
hand-made 1.x hub and spoke repositories the qa seat built.

**Every one of the ten task Plans wrote `None` into its Ask-the-human section,
and every one gave a reason.** #255's names its single judgment call — naming
`migrate` before the verb exists — as "settled by the task brief and M13-R1".
#256's, #257's, #258's, #259's, #260's, #261's, #262's, #285's and #288's each
name the calls they made and say they are recorded as Decisions rather than
raised as gates. [M12's record](12-v1-2-0-and-the-field-fixes-behind-it.md#the-gates)
observed one of its four Plans leaving the scaffold's own line unedited; none
did here.

The one scope question a seat did not answer itself did not go through
`checkpoint` either. It was flagged on PR #280 for the operator, answered, and
written up as a
[Decision](https://github.com/radiusred/gh-codecrew/issues/256#issuecomment-5559741725)
naming the operator's call — "I flagged the gap on PR #280 rather than
deciding it". The protocol has a verb for exactly that shape, and the record
shows the shape without the verb.

What gated the work: CI on every PR — `Go build and test` and
`Lint commit messages`, both required — and an independent approval on each of
the ten PRs, then the QA verdicts. No release gate existed, M13-R8 having been
struck.

## The review rounds

Ten PRs, **twenty-three review submissions**, all by the reviewer role holder
(`radiusred-checky[bot]`), and **twelve change requests**, every one resolved
in the next round. Eleven approvals were submitted; ten stand, one was
dismissed.

- **PR [#274](https://github.com/radiusred/gh-codecrew/pull/274)** —
  [approved first round](https://github.com/radiusred/gh-codecrew/pull/274#pullrequestreview-5125231119)
  at 11:56:47Z, six minutes thirty-seven seconds after opening. "I reviewed
  the diff before the PR body ... No findings."
- **PR [#275](https://github.com/radiusred/gh-codecrew/pull/275)** — two
  rounds. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/275#pullrequestreview-5125242306)
  is about test coverage of the call sites rather than the parser: "the
  changed call sites in `internal/cli/evidence.go` and `internal/cli/status.go`
  need verb/report-level tests that would fail if those call sites were
  removed." The
  [approval](https://github.com/radiusred/gh-codecrew/pull/275#pullrequestreview-5125284341)
  confirms it by mutation: "a compile-safe mutation that disables those
  branches makes `TestMilestoneEvidenceRefusesMismatchedRequirementID` and
  `TestStatusReportsRequirementIDMismatch` fail for the missing behavior."
- **PR [#276](https://github.com/radiusred/gh-codecrew/pull/276)** — four
  rounds, and the milestone's dismissed approval, described below.
- **PR [#277](https://github.com/radiusred/gh-codecrew/pull/277)** — two
  rounds. Four findings in
  [round one](https://github.com/radiusred/gh-codecrew/pull/277#pullrequestreview-5125268644),
  the first of them the walk-order bug, reproduced on the built binary: "from
  a scaffolded 2.0 repo plus `playbook/roles/qa.md`, `gh-codecrew roles show
  qa --latest` exits 1 with `refused[LAYOUT_LEGACY]` naming `playbook`". The
  other three are documentation surfaces the layout move missed — two SVG
  assets and two source comments still naming `roles/`. The
  [approval](https://github.com/radiusred/gh-codecrew/pull/277#pullrequestreview-5125392342)
  re-ran the scenario in both directions and checked the rebase clause by
  clause.
- **PR [#278](https://github.com/radiusred/gh-codecrew/pull/278)** — three
  rounds, described below.
- **PR [#279](https://github.com/radiusred/gh-codecrew/pull/279)** — two
  rounds. Round one opens by saying the change itself could not be broken —
  "the fail-closed shape is structural ... and I could not break it" — and
  then blocks on two refusal *details*. First, `HUB_UNREADABLE` "tells every
  failed fetch to run `migrate`, including a 403", when "a spoke's App that
  has not been granted the hub is the most likely 403 in a private topology,
  and the remediation it is handed is a *writing* verb against a hub that is
  fine". Second, the self-naming hub's `SPOKE_ROUTING` detail read "delete the
  routing here and declare it here. The operator has no next move." The
  [approval](https://github.com/radiusred/gh-codecrew/pull/279#pullrequestreview-5125518232)
  records that the fix went further than the finding and checked why: "`gh api
  repos/<private repo>/contents/.codecrew/config.yml` under this seat's
  installation token, on a private repo this App is not installed on, answers
  **404**, not 403. So 'absent' alone would have leaked the same misdirection
  through the other door."
- **PR [#280](https://github.com/radiusred/gh-codecrew/pull/280)** — three
  rounds, both change requests genuine bugs. Round one: `applyMoves` fell back
  from any `git mv` failure to `os.Rename`, so "a target path such as
  `.codecrew/roles/qa.md` already exists while `.codecrew/config.yml` does
  not, `git mv` fails because the destination exists, the fallback overwrites
  the existing 2.0 file, and the old tracked source deletion is left outside
  the migration commit". Round two: a 1.x spoke pointer carrying a `roles:`
  block "is accepted by `migrate`; it commits `.codecrew/config.yml` ... then
  this same binary immediately refuses `gh codecrew status` with
  `refused[SPOKE_ROUTING]`. That is not a completed one-shot move to the 2.0
  layout." The
  [approval](https://github.com/radiusred/gh-codecrew/pull/280#pullrequestreview-5125595670)
  confirms the refusal order on a scratch repo that also had a destination
  clash: "the observed refusal was `SPOKE_ROUTING`, HEAD stayed on the fixture
  commit ... and `git status` stayed clean."
- **PR [#282](https://github.com/radiusred/gh-codecrew/pull/282)** — two
  rounds, over one source comment: `internal/cli/context.go` "still says
  `"0.1"` and a missing field proceed with a note on stderr. That describes the
  shim this PR removes." Round one had already verified the substance: "42
  source codes, 42 SPEC rows, no differences."
- **PR [#286](https://github.com/radiusred/gh-codecrew/pull/286)** — two
  rounds, all three round-one findings the same class: three source-level
  contracts — the doc-synthesizer contract, `extractURLs`' doc comment and
  `ParseVerdicts`' doc comment — still described the two-form code rule the PR
  had just made a three-form rule. The
  [approval](https://github.com/radiusred/gh-codecrew/pull/286#pullrequestreview-5125787032)
  checked the shipped copies, not just the files: "Built the binary and ran
  `roles show doc-synthesizer` and `roles show qa` against it: the embedded
  copies carry the three-form wording, so an agent dispatched from the binary
  reads the same rule the code enforces."
- **PR [#289](https://github.com/radiusred/gh-codecrew/pull/289)** — two
  rounds. The
  [change request](https://github.com/radiusred/gh-codecrew/pull/289#pullrequestreview-5125846294)
  is a coverage finding proved by mutation: `thematicBreak`'s "nothing else on
  the line" clause "is the only clause with no test behind it", and the row
  that looks like its guard is not one. "I deleted that arm ... and ran
  `go test ./...`: `internal/tracker` and `internal/cli` both pass. No row
  anywhere in the repo fails." It names the direction of the risk: without the
  arm, "`--- x` and `*** NOTE ***` become openers, and the four-column line
  under either is dropped — an over-strip, a lost verdict, not the missed
  strip the Decision accepts." The
  [approval](https://github.com/radiusred/gh-codecrew/pull/289#pullrequestreview-5125864825)
  re-ran the whole mutation sweep on the new head and tabulates nine mutations
  against the rows that fail for each.

**The dismissed approval on PR #276.** The chronology from the timeline: a
[change request](https://github.com/radiusred/gh-codecrew/pull/276#pullrequestreview-5125255118)
at 12:08:41Z (the scaffolded pointer's comment still taught the 1.0 identity
grammar, so "a fresh hub initialized by this binary tells the operator to
write bare values that the same binary refuses with `IDENTITY_UNTYPED`"); a
force-push at 12:11:00Z; an
[approval](https://github.com/radiusred/gh-codecrew/pull/276#pullrequestreview-5125298621)
at 12:29:16Z; then a rebase force-push at 12:31:40Z onto #260's merge, which
dismissed that approval in the same second. Round three at 12:43:50Z is
therefore a review of a rebase, and it says so in its heading: "the rebase
claim holds; one finding that is mine to have missed at 16f5afc, not a rebase
regression." It verifies the rebase with `git range-diff` and a patch-to-patch
comparison — "Same 27 files, same `837 insertions(+), 207 deletions(-)`. **No
Go file's added or removed lines changed at all.**" — and then reports a
README count one behind the introduction's, with a candid note about where it
came from: "This is not rebase damage — it is my miss at `16f5afc` ... The
rebase carried the gap forward one notch rather than creating it, and I should
have caught it in round two." Round four is the one-line check.

**PR #278's three rounds** turn on a test that could not fail. Round one
blocked on `init` printing a false `action needed` on its own scaffolded
files, quoting the run verbatim and giving three reasons it is a finding
rather than a nit, the last of which is "It fires on this hub." Round two
verified the fix and then found that the test written to close it does not
close it: the assertion asked `strings.Contains(got, f)` of the whole output,
"and `got` already contains `kept existing CLAUDE.md` from the skip report a
dozen lines earlier. So the assertion passes whether or not the `Kept:` line
names the file — the table never checks the one thing its `stranded` column
claims to specify." The review then proves the cost by mutation: deleting the
one-hop check leaves the case green. Round two also asked for the narrowing to
be recorded as a labelled Deviation rather than as PR prose, giving the
mechanical reason: "SPEC §4 gathers the record by labelled paragraph, so the
doc-synthesizer will not see it." Round three approved after repeating both
mutations itself.

**Two review bodies reached GitHub as the wrong text, and were edited after
submission.** PR #278's round-two and round-three reviews were submitted at
13:32:29Z and 13:38:37Z with a body consisting of a single line: an `@`
followed by the path of a `verdict.md` file in the coordination layer's
dispatch directory — a file reference the dispatch never expanded. Both were
edited to the real verdict text at 13:39:34Z and 13:39:35Z, before the merge
at 13:40:43Z; the PR's `userContentEdits` history holds both versions of each.
Recorded plainly as a coordination-layer tooling fault: the reviewer's
judgment was formed and written, and the transport lost it. Both edits are
recorded under the reviewer identity, `radiusred-checky`, a minute after the
second submission; *who or what noticed the fault is not on the record*, and
an edit under an App identity says only which credential made it. No verdict
changed, and the approval that gated the merge is the reviewer's own.

**What the twelve change requests were about.** Five were real behaviour bugs
found by running the code: the walk-order false refusal (#277), the two
migration bugs (#280), and the two refusal details that sent an operator the
wrong way (#279). Four were documented surfaces that no longer matched the
behaviour the PR shipped — a scaffold's own comment (#276), two SVG assets and
two source comments (#277), a stale comment about a deleted shim (#282), three
source-level contracts describing the two-form code rule (#286). Two were
tests that did not pin what they claimed: the call sites in #275 and the
`Kept:` assertion in #278. One was a coverage gap proved by deleting the
clause and watching the suite stay green (#289). Every approval was reached
with the reviewer running the code — repo-local builds, scratch hubs and
spokes, live read-only probes, and in five of the ten PRs — #275, #278, #279,
#280 and #289 — a mutation applied to the tree and then reverted.

## QA: three rounds, and the requirement that took all three

The qa role holder (`radiusred-testy[bot]`) verdicted every requirement in a
[first comment](https://github.com/radiusred/gh-codecrew/issues/254#issuecomment-5560084992)
at 15:01:59Z, eleven minutes after the last delivery task's PR merged. Seven
came back `satisfied`. **M13-R6 came back `not satisfied`**, and the two
comments that follow are the milestone's most instructive stretch of record.

**Round one.** The verdict on R6 credits what the shipped tests prove — "gates
are read per paragraph, only `Gate resolved` resolves them, verdict tally is
per-comment with first match inside a comment, fenced/inline code verdicts are
stripped, and foreign requirement IDs refuse `REQUIREMENT_ID_MISMATCH` before
walking evidence" — and then names what the probe found: "a verdict-shaped
line in a four-space indented code block is still counted by
`tracker.ParseVerdicts` (`[{M13-R1 satisfied qa}]`), while a fenced version is
ignored". The
[finding](https://github.com/radiusred/gh-codecrew/issues/260#issuecomment-5560080971)
was filed on #260 with the probe body quoted, what happened, and what should
happen: "A Markdown indented code block is still code, so this quoted example
should not supersede a real QA verdict." Remedy task #285 opened seventy
seconds after the verdict, at 15:03:09Z; PR #286 merged twenty-two minutes
after that.

**Round two.** The
[re-verdict](https://github.com/radiusred/gh-codecrew/issues/254#issuecomment-5560246370)
at 15:30:10Z opens by scoping itself — "This supersedes my earlier R6 verdict
in comment 5560084992 for this requirement only" — the supersession rule
M13-R6 itself shipped, used on its own verdict. It confirms the remedy fixes
"the original four-space-after-blank probe, tab/CRLF indented blocks, internal
blank lines, URLs inside indented blocks, and the paragraph/list continuation
cases that must still count", and that a deeper-indented list probe correctly
did not count. Then: **`not satisfied`** again, on a narrower shape. "The
remaining hairbrush is a heading boundary ... `StripCode` only opens indented
blocks after a blank/start/fence. CommonMark 4.4 permits indented code
immediately after headings and other non-paragraph blocks; #285's Decision
deliberately excludes nested-list block parsing, but not headings." The
distinction the verdict draws is the reason it is a finding rather than a
quibble: a boundary a Decision names is a boundary; one it does not name is a
gap. #288 opened one minute later; PR #289 merged twenty-seven minutes after
that.

**Round three.** The
[final verdict](https://github.com/radiusred/gh-codecrew/issues/254#issuecomment-5560428532)
at 16:02:49Z is `satisfied`, again scoped to R6 alone. It records probes for
"a deeper-indented block inside a list item, an indented block immediately
after a heading, and CRLF line endings", each returning "only the real
`not satisfied` verdict below" the quoted sample, and it accepts the remaining
boundary by reading both Decisions: "the list-nesting boundary is acceptable
for R6 because the prescribed record forms are column-0 or plain-bullet forms,
while the #288 opener rule covers ordinary non-paragraph record structure such
as headings, thematic breaks and fence closes."

The seven other requirements were verdicted once each, in round one, against a
build of merged `main` and scratch repositories the seat made:

- **M13-R1** on a fresh scaffold (`.codecrew/config.yml` carrying
  `codecrew: "2.0"`, all five contracts and extensions, `.codecrew/AGENTS.md`,
  the root pointer, `CLAUDE.md` and `ROADMAP.md`), a hand-made 1.x repo
  refusing `LAYOUT_LEGACY`, and "a valid 2.0 repo with nested
  `playbook/roles/qa.md` ran `roles show qa` from `playbook/` without legacy
  refusal" — round one's finding on PR #277, re-run independently.
- **M13-R2** on hand-made 1.x hub and spoke repos: `--dry-run` writing
  nothing, a live run producing "a migration-only pathspec commit" that typed
  bare identities, set the version and added the coordinator row "while
  leaving unrelated staged work staged"; `FOREIGN_ROLES_DIR`, `BOTH_LAYOUTS`
  and `SPOKE_ROUTING` each raised on their own shape, and "rerun on the
  migrated hub exited 0 as a no-op".
- **M13-R3** on `init` over a pre-existing foreign root `AGENTS.md`, keeping
  the file and printing the exact lines, then "the foreign-file case stopped
  printing `action needed` after those exact lines were added" — the
  Deviation's narrowing, tested from the other end.
- **M13-R4** on `role reviewer` printing `app:radiusred-checky` and
  `--login` printing zero bytes for it, with a bare row refusing
  `IDENTITY_UNTYPED` "naming `roles.reviewer.identity` and the four allowed
  forms".
- **M13-R5** on a scratch spoke whose hub config could not be read refusing
  `HUB_UNREADABLE`, "not a gate fallback"; and with `GH_HOST` pointed at a
  loopback address, `version`, `help`, `roles show reviewer` and
  `roles diff reviewer` all exiting 0.
- **M13-R7** on the code sets: "The source refusal-code set from non-test
  `internal/cli` and SPEC §10 table both contained 42 codes with an empty
  diff", the `0.1` pointer refusing under both layouts, `role coordinator`
  erroring instead of falling back to `~`, and an assigned task with no start
  record refusing `NOT_OWNER` "with no phantom owner".
- **M13-R9** on `milestone evidence 12` resolving the closed milestone and
  printing `note: milestone M12 (radiusred/gh-codecrew#241) is closed` before
  its citation report, and on "live `status` on the installed
  `dev-3811ba4 (protocol 2.0)` build" showing no contract-drift lines.

Cross-cutting, round one: `go test ./...`, `go vet ./...` and `gofmt -l .` all
clean "with Go temp/cache dirs under the milestone workspace and outside the
clone". Rounds two and three each re-ran the suite and the focused
code-stripping tests on the merged remedy, and each verified the seat's own App
ID against `gh api /apps/radiusred-testy --jq .id` before verdicting.

## Requirement outcomes

The status column is the **standing** verdict word as the qa role holder wrote
it, verbatim and unqualified — the latest comment carrying a verdict for the
ID, which is the supersession rule M13-R6 shipped. M13-R6's earlier two
verdicts are in the Notes column and in the section above.

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M13-R1 — the 2.0 layout: every CodeCrew-owned operational file under `.codecrew/` in hub and spoke alike, `init` and every verb reading there, this hub's own files moved in the same PR, `LAYOUT_LEGACY` naming `gh codecrew migrate`, `PROTOCOL_MISMATCH` narrowed, and no dual-read of the 1.x layout anywhere | [#255](https://github.com/radiusred/gh-codecrew/issues/255) / PR [#277](https://github.com/radiusred/gh-codecrew/pull/277) | `satisfied` | Task closed; verdicted in QA round one, on a fresh scaffold, a hand-made 1.x repo and the nested-`roles/` case the review had found. Two review rounds; round one's finding was a false `LAYOUT_LEGACY` on a valid 2.0 repo. Unreleased |
| M13-R2 — `gh codecrew migrate`: a one-shot verb, exempt from the pointer check, moving a 1.0 hub or spoke to the 2.0 layout in one pathspec commit, never pushing; conservative about `roles/`, refusing rather than guessing, idempotent | [#256](https://github.com/radiusred/gh-codecrew/issues/256) / PR [#280](https://github.com/radiusred/gh-codecrew/pull/280) | `satisfied` | Task closed; verdicted in QA round one on hand-made 1.x hub and spoke repos, dry run and live. Nine Decisions, three review rounds, two of them real bugs. The 2.0 entry point was added to the verb mid-review on the operator's call. Closes [#281](https://github.com/radiusred/gh-codecrew/issues/281). Unreleased |
| M13-R3 — the entry point stands on its own: `init` always writes `.codecrew/AGENTS.md`, the root `AGENTS.md` scaffold is a short pointer that imports it, and when a root entry point already exists `init` prints the exact line to add | [#257](https://github.com/radiusred/gh-codecrew/issues/257) / PR [#278](https://github.com/radiusred/gh-codecrew/pull/278) | `satisfied` | Task closed; verdicted in QA round one on `init` over a foreign root file and on a clean rerun. Three review rounds; one recorded Deviation narrows the literal wording, at the reviewer's request and with the reviewer's agreement. Unreleased |
| M13-R4 — typed identities: `~`, `app:<slug>`, `user:<login>` or `team:<org>/<slug>`; a bare string refused with a code; `role <name>` prints the kind; `crew()` and the request-review rule read the kind instead of inferring it | [#258](https://github.com/radiusred/gh-codecrew/issues/258) / PR [#276](https://github.com/radiusred/gh-codecrew/pull/276) | `satisfied` | Task closed; verdicted in QA round one against a scratch hub with typed routing and a bare row refusing `IDENTITY_UNTYPED`. Four review rounds, one approval dismissed by a rebase. Unreleased |
| M13-R5 — routing fails closed: `HUB_UNREADABLE` instead of an empty table that resolves every seat to `~`; a hub declaring no table is still legitimately `~`; protocol skew and a spoke `roles:` block refused; GitHub unreachable named by its own code; the local verbs keep working without the network | [#259](https://github.com/radiusred/gh-codecrew/issues/259) / PR [#279](https://github.com/radiusred/gh-codecrew/pull/279) | `satisfied` | Task closed; verdicted in QA round one on a scratch spoke and on the four hub-local verbs run with the host unreachable. Two review rounds, both findings on refusal details rather than on the change. Unreleased |
| M13-R6 — the record grammar tightened: `**Gate raised:**` per paragraph, only `**Gate resolved:**` resolves and per gate, verdict supersession per comment with code stripped, and a requirement ID must carry its milestone's own number | [#260](https://github.com/radiusred/gh-codecrew/issues/260) / PR [#275](https://github.com/radiusred/gh-codecrew/pull/275); remedies [#285](https://github.com/radiusred/gh-codecrew/issues/285) / PR [#286](https://github.com/radiusred/gh-codecrew/pull/286) and [#288](https://github.com/radiusred/gh-codecrew/issues/288) / PR [#289](https://github.com/radiusred/gh-codecrew/pull/289) | `satisfied` | Task closed; **three QA rounds**. `not satisfied` at 15:01:59Z (a verdict in a four-space indented block still counted), `not satisfied` at 15:30:10Z (an indented block immediately after a heading), `satisfied` at 16:02:49Z. Each supersession is scoped "for this requirement only". Unreleased |
| M13-R7 — the 1.0 shims removed and the machine contract written down: the `"0.1"` acceptance, the coordinator special case and the first-assignee fallback deleted; SPEC §6 states the exit-code contract and SPEC §10 lists every stable refusal code in one table | [#261](https://github.com/radiusred/gh-codecrew/issues/261) / PR [#282](https://github.com/radiusred/gh-codecrew/pull/282) | `satisfied` | Task closed; verdicted in QA round one on an empty diff between the source's forty-two codes and SPEC §10's forty-two rows, plus each removal's own refusal. One adopted scan finding — inverting the missing-field rule — was declined with a recorded Decision. Unreleased |
| M13-R9 — two open captures ride the release: `status` prints the drift report and the branch note with no open milestone (adopts [#253](https://github.com/radiusred/gh-codecrew/issues/253)); `milestone evidence` resolves a closed milestone, saying it is closed, while `close` and `status` keep their open-only reads (adopts [#250](https://github.com/radiusred/gh-codecrew/issues/250)) | [#262](https://github.com/radiusred/gh-codecrew/issues/262) / PR [#274](https://github.com/radiusred/gh-codecrew/pull/274) | `satisfied` | Task closed; verdicted in QA round one on `milestone evidence 12` against the closed M12 and on the no-open-milestone tests. Approved first round. Unreleased |

**M13-R8 was struck**, not delivered and not verdicted: "v2.0.0 ships: the
`[Unreleased]` changelog becomes the 2.0.0 section leading with what broke and
the migration steps; `docs/introduction.md` and README flip; the tag is cut by
the implementer identity after the flip merges; `release.yml` builds it; the
operator verifies `gh extension upgrade codecrew` prints v2.0.0 and `status`
runs clean on this migrated hub; then the operator migrates the remaining
twelve repos, hubs before their spokes, and records the run on this issue."
The [Decision](https://github.com/radiusred/gh-codecrew/issues/254#issuecomment-5558926520)
that struck it names the consequence and accepts it: "this hub runs on the
local build from #255's merge until M14's release, a longer window than
planned and accepted." No requirement was added or amended; R9 kept its
number, "the close gate reads the IDs present".

**Captures adopted and closed by the merges — three:**
[#253](https://github.com/radiusred/gh-codecrew/issues/253) (`status` skips
the contract-drift report when no milestone is open) and
[#250](https://github.com/radiusred/gh-codecrew/issues/250)
(`milestone evidence` refuses `NOT_FOUND` on a closed milestone), both closed
by PR #274's merge at 11:57:45Z; and
[#281](https://github.com/radiusred/gh-codecrew/issues/281) (the legacy-layout
tests are flaky and environment-dependent), filed by the operator at 13:40:04Z
from two seats' reports during the milestone and closed by PR #280's merge at
14:22:50Z. #250 is the capture
[M12's record](12-v1-2-0-and-the-field-fixes-behind-it.md#requirement-outcomes)
describes being filed from the reviewer's live probe on its own PR #248 — a
capture filed by one milestone's review and closed by the next milestone's
first merge.

**Captures filed during the milestone and left open — five**, all from the
scans and the offline question, all still open:
[#264](https://github.com/radiusred/gh-codecrew/issues/264),
[#265](https://github.com/radiusred/gh-codecrew/issues/265),
[#266](https://github.com/radiusred/gh-codecrew/issues/266),
[#267](https://github.com/radiusred/gh-codecrew/issues/267) and
[#268](https://github.com/radiusred/gh-codecrew/issues/268). #267 has since
been adopted into M14 as task
[#283](https://github.com/radiusred/gh-codecrew/issues/283).

**One capture decided and deliberately left open.**
[#177](https://github.com/radiusred/gh-codecrew/issues/177) (a spoke shared by
several hubs) received an operator
[Decision](https://github.com/radiusred/gh-codecrew/issues/177#issuecomment-5558957307)
at 11:38:06Z, "while shaping protocol 2.0 on #254": the pointer keeps a single
`hub:`, blessed permanent, and the fix for a shared spoke is task-scoped
resolution through the task's parent milestone. "That is a behaviour
refinement, not a schema change, so it does not need to ride 2.0; the SPEC §3
sentence recording the rule rides #255. This capture stays open for the
resolution work." The
[scope addition](https://github.com/radiusred/gh-codecrew/issues/255#issuecomment-5558957229)
to #255 says the same in one line: "The sentence only; the task-scoped
resolution itself is a later minor."

## Protocol-discipline observations

Eight things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **A milestone commissioned a review of its own scope, and the review's
  findings became the requirements.** The pattern is new to this record and
  cheap to describe: two fresh-context sessions on the same brief against the
  same commit, ten findings each with file-and-line evidence, both attached to
  the milestone issue in full, and one Decision comment sorting every finding
  into adopted, blessed or filed by number. What it produced is auditable in a
  way a conversation is not — every requirement from R2 to R7 traces to a
  numbered finding, and the four blessings are decisions with a date rather
  than defaults nobody wrote down. What it cost is not on the record: the
  briefs are in the coordination layer's dispatch directory, and the scans'
  own token or wall-clock cost is nowhere.
- **Two independent scans agreed, and the agreement is itself evidence.** The
  findings that both models reached independently — the pointer's location,
  `migrate` needing its own code and its own exemption, the requirement-ID
  grammar — are the three the operator adopted with the least hedging. Where
  they diverged, the Decision took both sides: Codex's namespace questions
  were answered as blessings, Claude's fail-open behaviours as requirements.
  Nothing in the record says the agreement was weighed that way, but the
  adoptions and the blessings fall along that line.
- **A verdict superseded itself twice, under the rule the same milestone
  shipped.** M13-R6 taught `milestone close` that "the latest comment carrying
  a verdict for an ID wins", and the qa seat then used exactly that mechanism
  to re-verdict R6 twice, each time writing "This supersedes my earlier R6
  verdict in comment N for this requirement only". The rule was exercised on
  its own requirement, by hand, before any close ever read it — which is a
  better test than the milestone's tests, and unplanned.
- **Both R6 remedies came from probing the boundary a Decision had named.**
  #285's Decision listed what it deliberately did not handle; the qa seat
  probed the list, found the listed cases correctly excluded and one unlisted
  case wrong, and said so in those terms — "#285's Decision deliberately
  excludes nested-list block parsing, but not headings". The final verdict
  accepts the surviving boundary by citing both Decisions. A Decision that
  states its own limits gives QA something to test against; one that does not
  leaves QA guessing what was meant.
- **Three of the twelve change requests were about text that no longer matched
  behaviour, and one of them was about the contracts the binary ships.** PR
  #286's three findings were all source-level or contract-level sentences
  describing the two-form code rule after the code became three-form, and the
  approval checked the shipped copies by running `roles show` against the
  built binary rather than reading the files. [M12's
  record](12-v1-2-0-and-the-field-fixes-behind-it.md#protocol-discipline-observations)
  observed the same class on `internal/cli/cli.go`'s help, and its remedy —
  nothing in the tree pins the sentence — is still true here: no test pinned
  any of the three doc comments.
- **The one scope question a seat could not answer went to the operator, and
  did not go through `checkpoint`.** The gap PR #280 found — a migrated hub
  with no 2.0 entry point — is exactly the shape `checkpoint` and
  `cc:needs-decision` exist for: a question about what "done" means that the
  seat should not settle. It was raised in a PR comment, answered out of band,
  and written up as a Decision naming the operator's call. The outcome is
  correct and the record is complete; the mechanism the protocol provides went
  unused, in the milestone that tightened the grammar of the labels it uses.
- **The doc-synthesizer seat held no delivery task, and the routing matched
  the table.** All ten delivery and remedy tasks were started by
  `radiusred-cody[bot]`, as M12's were.
  [M10's](10-protocol-bookkeeping-from-the-field.md#protocol-discipline-observations)
  and [M11's](11-housekeeping.md#protocol-discipline-observations) records
  each observed `radiusred-wordy[bot]` starting delivery tasks under the
  implementer contract and noted that nothing in the contracts says whether a
  role's holder may be routed to another role's task. That is still unwritten;
  two milestones running have simply not needed it.
- **The refresh obligation had something to fix this time, and it is the
  mirror of M12's.**
  [M12's record](12-v1-2-0-and-the-field-fixes-behind-it.md#protocol-discipline-observations)
  found `docs/introduction.md` correct because its refusal-code count named
  `main` as its source while running ahead of the binary. Here the same page's
  `**Shipped:**` line correctly names v1.2.0 and protocol 1.0, while the rest
  of the page — the layout paths, the identity grammar, `gh codecrew migrate`
  in the install block — describes a `main` no released binary implements. The
  claim was not false, but the page could not be read straight. `README.md`
  has the same gap one page earlier in the funnel, and it does not need a
  version number to have it: its "Start now" block pairs
  `gh extension install radiusred/gh-codecrew` with
  `gh codecrew init # writes and commits .codecrew/, AGENTS.md, CLAUDE.md,
  ROADMAP.md`, and the released v1.2.0 `init`, run in a scratch repository,
  writes `.codecrew.yml` and `roles/` instead. The sweep for this document
  added the same dated paragraph to both pages, naming what is on `main` and
  unreleased and pointing at M14; it changed no version number, and M14's
  release task deletes both. *The claim a boundary refresh has to test is not
  only "does this page name a version" but "would a reader who followed this
  page get what it describes".*

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **The conversation that opened the milestone is on the record only through
  #254's goal and its Decisions.** The operator's question about `.codecrew/`,
  the dual-read proposal and its rejection against the scar-tissue essay reach
  the record as the goal's summary — "Decided in conversation with the
  operator 2026-09-06" — and as the adoption Decision's "No compatibility shim
  anywhere". No transcript, no quotes, no timestamps for the exchange. *The
  reasoning is recorded; the deliberation is not.*
- **The two scan briefs are named and not attached.** #254's goal says the
  "briefs and findings [are] in the operator's dispatch dir, 2026-09-06", and
  the findings were pasted into the issue in full. The briefs were not. Both
  scan comments paraphrase the brief in one line — "what else on the
  protocol's public surface would force a 3.0 if 2.0 shipped as the layout
  move alone" — which is enough to judge the findings against, and less than
  the findings themselves get. *Recorded as cited, not seen*, as
  [M10's record](10-protocol-bookkeeping-from-the-field.md#what-the-record-does-not-contain)
  and [M12's](12-v1-2-0-and-the-field-fixes-behind-it.md#what-the-record-does-not-contain)
  each say of material that reached them the same way.
- **Which reviewer rounds ran on which harness is only partly recorded.** The
  Deviation names two rounds it expected to affect (#276 round three, #277
  round two) and says codex resumes after the reset, which was around 15:20
  BST. Ten review submissions fall between the quota exhaustion at 12:35Z
  and 14:20Z, and nothing on the trail attributes them individually. The
  identity is constant and verified throughout — every one is
  `radiusred-checky[bot]`, and the App ID check appears in six of the ten
  PRs' reviews — so the gate the protocol enforces is intact; the harness
  behind each round is not reconstructable.
- **The `@path` fault has no artefact beyond the edits.** The two PR #278
  reviews whose bodies arrived as an unexpanded file reference are recoverable
  only from GitHub's own edit history. Nothing on the trail records the fault,
  how it was noticed, or whether it affected any other dispatch in the
  milestone; no Deviation was written for it. *Inferred from the edit history,
  not recorded:* that the coordination layer passed a path where it meant a
  file's contents, and repaired it by editing both reviews about a minute
  later.
- **`migrate` has never run against a repository that was not made for it.**
  Every exercise of the verb on the record is against a scratch 1.x repo built
  for the test — the implementer's fixtures, the reviewer's scratch hubs and
  the qa seat's hand-made hub and spoke — plus one run against "a real copy of
  this hub's own pre-#276 1.0 pointer" with the live users API. The fleet is
  about a dozen repositories and none of them has been migrated. That is M14's
  work ([#272](https://github.com/radiusred/gh-codecrew/issues/272)), and
  until it runs the verb's conservative refusals are untested against a
  repository nobody designed.
- **`HUB_UNREADABLE` and the skew refusals have never fired in the topology
  they were written for.** M13-R5 exists because the `.codecrew/` move
  guarantees a hub/spoke skew during migration; the fix landed before any
  spoke was migrated, so the window it protects has not opened yet. The
  verdict rests on a scratch spoke and on tests over nine failure shapes.
- **No gate was declared or raised, and `checkpoint` was not used.** #254's
  Gates section is the scaffold's placeholder text, unedited; no
  `cc:needs-decision` label exists on any issue in the repository; and the
  milestone that rewrote the gate grammar never exercised it on itself. The
  M5 comment corpus in `internal/tracker/testdata/` carries the project's one
  hand-raised gate, and that corpus is what the new rules were tested against.
- **Nothing measures what the milestone cost.** Ten PRs, twenty-three
  reviews, three QA rounds and two remedies inside four and a half hours, with
  no record of tokens, wall-clock per seat, or how many dispatches were made.
  The scar-tissue essay the goal cites argues about exactly that cost;
  the milestone that cites it produced no figure of its own.
- **This record cannot be re-checked after the close by the released
  binary.** `milestone evidence 13` resolves a closed milestone as of M13-R9,
  so the verb can read this document after #254 closes — but only from a build
  of `main`. The installed release refuses this hub outright, and until v2.0.0
  ships in M14 there is no released binary that can check these citations at
  all.
