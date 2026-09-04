# M10: Protocol bookkeeping from the field

Tracking issue: [#207](https://github.com/radiusred/gh-codecrew/issues/207) ·
Synthesized 2026-09-04 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #207's **seven** requirements — six present when the
milestone opened, **M10-R6 amended and M10-R7 added by the operator on
2026-09-04**, with the original R6 clause quoted in the issue's Amendments
section — its Gates section and its **two** QA rounds; the ten delivery task
issues, nine in this repository
([#208](https://github.com/radiusred/gh-codecrew/issues/208),
[#209](https://github.com/radiusred/gh-codecrew/issues/209),
[#210](https://github.com/radiusred/gh-codecrew/issues/210),
[#211](https://github.com/radiusred/gh-codecrew/issues/211),
[#212](https://github.com/radiusred/gh-codecrew/issues/212),
[#213](https://github.com/radiusred/gh-codecrew/issues/213),
[#214](https://github.com/radiusred/gh-codecrew/issues/214),
[#226](https://github.com/radiusred/gh-codecrew/issues/226),
[#229](https://github.com/radiusred/gh-codecrew/issues/229)) and one in the
spoke [radiusred/codecrew-www](https://github.com/radiusred/codecrew-www)
([#19](https://github.com/radiusred/codecrew-www/issues/19)); their ten merged
PRs; the **twenty-eight Decision and two Deviation comments** across the
milestone issue and those task issues; the **sixteen review submissions** on
the ten PRs; the seven adopted backlog captures the merges closed and the three
the milestone filed; and #207's GitHub edit history, which timestamps both
amendments. The house form is
[M9's document](9-the-docs-at-codecrew-works.md) and
[M8's](8-a-product-home-page-for-codecrew-works.md); the record standard is the
one the reviewer set on their PRs,
[#204](https://github.com/radiusred/gh-codecrew/pull/204) and
[#206](https://github.com/radiusred/gh-codecrew/pull/206) — verdict words
verbatim, no ordinal claim without every prior record linked, no judgement the
record does not carry.

The gather from `gh codecrew milestone close 10` is again not part of the raw
material, for the reason [M8's](8-a-product-home-page-for-codecrew-works.md)
and [M9's](9-the-docs-at-codecrew-works.md) records both give: `--dry-run`
stops at the second gate, naming the task that writes this file.

```
gate milestone open: ok
gate tasks closed: refused[OPEN_TASKS]: tasks not closed: radiusred/gh-codecrew#231 (ready)
gate requirements declared: not reached
gate QA verdicts: not reached
gate milestone document: not reached
dry run: nothing written — the live verb stops at the first refusal above
```

Every record below was read from its issue directly.

**This PR adds the M10 ROADMAP row; it does not flip one.** That is the
convention M10-R1 introduced and
[#208](https://github.com/radiusred/gh-codecrew/issues/208) / PR
[#216](https://github.com/radiusred/gh-codecrew/pull/216) shipped; the
convention did not exist until that PR merged at 15:21:58Z on 2026-09-04, so no
earlier record PR could have been written under it. `milestone new` had
appended an `Open` row to a local `ROADMAP.md` when #207 was created; the
operator
[discarded it](https://github.com/radiusred/gh-codecrew/issues/207#issuecomment-5540526036)
rather than carry it, so that the milestone which retires the old convention
does not also practise it.

## Goal and outcome

Three field runs on v1.1.0 — the `radiusred/ops` company hub, the
`davison/topos` solo hub and the `codecrew-www` spoke — turned up a cluster of
bookkeeping gaps, and #207's goal is the list: the ROADMAP row has nowhere to
ride when a milestone's work is all in spokes (hit three times); `milestone new`
hands out a taken number when calls come seconds apart; `status` hides a gate
raised on a milestone issue; `task finish` dies on a raw GraphQL error when the
App lacks `checks: read` on a private repo; and the dispatch guidance never says
the dispatched session must reach GitHub. Two pieces of hub hygiene rode
alongside: the routing table shown in the README, and the unreferenced crew
badge PNGs removed.

Each of the five gaps was a backlog capture already written down when the
milestone opened — #197, #195, #200, #198 and #202 — and so were both hygiene
items, #199 and codecrew-www#18. **Seven captures were adopted and all seven
were closed by the merges**; three new ones were filed and remain open.

Six requirements were present when the milestone opened at 12:36:50Z on
2026-09-04. M10-R6 was amended at 17:16:24Z and M10-R7 added at 17:19:00Z, both
by the operator, both timestamped in #207's edit history. Ten tasks delivered
the work — nine in the hub and one in the spoke — and this document is the
eleventh. The last delivery merged at 18:00:53Z, five hours and twenty-four
minutes after the milestone opened.

- **[#214](https://github.com/radiusred/gh-codecrew/issues/214) / PR
  [#215](https://github.com/radiusred/gh-codecrew/pull/215)** — the ten
  unreferenced crew badge PNGs deleted from `assets/`, one `chore:` commit,
  approved first round and merged at 12:45:22Z
  ([d125950](https://github.com/radiusred/gh-codecrew/commit/d125950)). Adopts
  #199. The PR records "No Decisions or Deviations recorded: the plan was
  followed as written", and none exist.
- **[#208](https://github.com/radiusred/gh-codecrew/issues/208) / PR
  [#216](https://github.com/radiusred/gh-codecrew/pull/216)** — M10-R1. Four
  commits, two review rounds, merged 15:21:58Z
  ([a6dcd76](https://github.com/radiusred/gh-codecrew/commit/a6dcd76)). Two
  Decisions. Adopts #197 (shape 3).
- **[#212](https://github.com/radiusred/gh-codecrew/issues/212) / PR
  [#217](https://github.com/radiusred/gh-codecrew/pull/217)** — M10-R5. Four
  commits, three review rounds, merged 15:46:18Z
  ([15ed1c7](https://github.com/radiusred/gh-codecrew/commit/15ed1c7)). Three
  Decisions. Adopts #202.
- **[#210](https://github.com/radiusred/gh-codecrew/issues/210) / PR
  [#218](https://github.com/radiusred/gh-codecrew/pull/218)** — M10-R3. Four
  commits, two review rounds, merged 15:56:36Z
  ([5467568](https://github.com/radiusred/gh-codecrew/commit/5467568)). Two
  Decision comments and the milestone's first Deviation. Adopts #200.
- **[#209](https://github.com/radiusred/gh-codecrew/issues/209) / PR
  [#221](https://github.com/radiusred/gh-codecrew/pull/221)** — M10-R2. Three
  commits, two review rounds, merged 16:34:40Z
  ([4049261](https://github.com/radiusred/gh-codecrew/commit/4049261)). One
  Decision. Adopts #195.
- **[#211](https://github.com/radiusred/gh-codecrew/issues/211) / PR
  [#220](https://github.com/radiusred/gh-codecrew/pull/220)** — M10-R4. Five
  commits, two review rounds, merged 16:43:13Z
  ([01f5285](https://github.com/radiusred/gh-codecrew/commit/01f5285)). Four
  Decisions. Adopts #198.
- **[#213](https://github.com/radiusred/gh-codecrew/issues/213) / PR
  [#223](https://github.com/radiusred/gh-codecrew/pull/223)** — M10-R6 as
  written. One commit, approved first round, merged 16:51:55Z
  ([eb2dc59](https://github.com/radiusred/gh-codecrew/commit/eb2dc59)). Four
  Decisions by the seat and one by the operator.
- **[#226](https://github.com/radiusred/gh-codecrew/issues/226) / PR
  [#228](https://github.com/radiusred/gh-codecrew/pull/228)** — the codex seats
  re-pinned to `gpt-5.5` under M10-R5. Three commits, approved first round,
  merged 17:09:40Z
  ([6e8ffa6](https://github.com/radiusred/gh-codecrew/commit/6e8ffa6)). Two
  Decisions and the milestone's second Deviation.
- **[#229](https://github.com/radiusred/gh-codecrew/issues/229) / PR
  [#230](https://github.com/radiusred/gh-codecrew/pull/230)** — the remedy for
  the M10-R6 amendment: the drift test deleted, the gloss rewritten. One commit,
  approved first round, merged 17:24:27Z
  ([00284f2](https://github.com/radiusred/gh-codecrew/commit/00284f2)). One
  Decision.
- **[codecrew-www#19](https://github.com/radiusred/codecrew-www/issues/19) / PR
  [codecrew-www#20](https://github.com/radiusred/codecrew-www/pull/20)** —
  M10-R7, the home page's example table. Two commits, approved first round,
  merged 18:00:53Z
  ([630ae87](https://github.com/radiusred/codecrew-www/commit/630ae87)). Three
  Decisions by the seat, one by the operator, and one comment marking a Decision
  superseded. Adopts codecrew-www#18.

Delivered in full, in two QA rounds. M10-R1 through M10-R6 were verdicted
`satisfied` in the
[first round](https://github.com/radiusred/gh-codecrew/issues/207#issuecomment-5544060296);
M10-R6 as amended and M10-R7 were verdicted `satisfied` in the
[second](https://github.com/radiusred/gh-codecrew/issues/207#issuecomment-5544580039),
which states that its R6 verdict supersedes the first. No verdict in this
milestone came back other than `satisfied`.

## Decisions

Twenty-eight Decision comments. Six are the operator's — three on the milestone
issue and three on task issues; the rest were written by the seats doing the
work. They group by the gap each task closed.

### The ROADMAP row belongs to the doc-synthesizer at both ends

[#197](https://github.com/radiusred/gh-codecrew/issues/197) had carried the gap
since 2026-09-02, and
[M9's record](9-the-docs-at-codecrew-works.md#protocol-discipline-observations)
named it first among its observations: a milestone whose tasks all live in spokes has
no hub PR for its ROADMAP row to ride, so the row went in as an operator commit
straight to `main` for both M8 and M9. The operator chose shape 3 — the
doc-synthesizer owns the row at both ends — over shapes 1 and 2, and #208
implemented it: `milestone new` stops touching `ROADMAP.md`, and the record PR
appends the row already `Done`.

[The verb prints nothing about the row at all](https://github.com/radiusred/gh-codecrew/issues/208#issuecomment-5540596686)
— not the eventual `Done` row as a hint, not a pointer to who adds it.
**Trade-off:** an operator used to v1.1.0's "ROADMAP.md updated locally" line
sees it vanish with nothing in its place; the CHANGELOG and the two contracts
are where the new ownership is stated. **Rejected:** printing the row template
as a copyable hint — the slug is not known when the milestone opens, and
`roles/doc-synthesizer.md` already carries the template verbatim for a seat
that reads it weeks later; and an informational line naming the doc-synthesizer
— "a verb that creates an issue should not narrate what another seat does at
the milestone's end".

[The repo-wide sweep pulled in two files the plan had not listed](https://github.com/radiusred/gh-codecrew/issues/208#issuecomment-5540615170).
`README.md` and `CONTRIBUTING.md` both claimed the roadmap names the open
milestone, which stops being true when rows are added only at a milestone's
end; both now point at `gh codecrew status` for the open one and at the roadmap
for the finished ones. **Rejected:** leaving them to this record PR — "a claim
the CLI change itself falsifies belongs with the change". `docs/identities.md`
was judged and left, being a file inventory rather than a description of the
row's lifecycle.

The operator's own Decision on #207 completes the pair: M10's `Open` row was
**discarded** rather than carried, because "carrying an `Open` row in the first
PR and then rewriting the verb that wrote it would put the old convention and
the new one in the same milestone's history". **Trade-off, recorded:** until
this PR merges, `ROADMAP.md` on `main` does not show M10 — `gh codecrew status`
does, "which is the state SPEC §3 calls live".

### `milestone new` repairs a collision instead of refusing on it

Three `milestone new` calls seconds apart in `radiusred/ops` all came back as
M2, because the label-filtered `cc:milestone` listing lagged the issue the first
call had just created
([#195](https://github.com/radiusred/gh-codecrew/issues/195)). #209's
[single Decision](https://github.com/radiusred/gh-codecrew/issues/209#issuecomment-5542759419)
chose repair over refusal: after `CreateIssue`, both listings are read again,
and when another issue already carries the `M<n>:` prefix the new issue is
**renumbered** — title and every generated `M<n>-R<k>` ID — through a new
`EditIssue` tracker method, with a `renumbered:` line printed and the check
repeated, bounded to three rounds. `MILESTONE_NUMBER_TAKEN` is kept for what
repair cannot cover: the edit fails, or the number is still taken after the
bound.

**Trade-off:** the repair moves the number after `created milestone …` has
already printed, so a coordinator reading only the first line can carry a stale
one; `roles/coordinator.md` now says to read the `renumbered:` line.
**Rejected:** refusing on detection, which "would have left the record with two
issues sharing a prefix — the exact state #195 had to clean up by hand".

Also rejected, and the reason is a dependency between two of this milestone's
own requirements: #195's second layer was a number floor read from
`ROADMAP.md`, and M10-R1 had just removed the row that made the file know the
open milestone. The floor became the maximum over two listings instead — the
label-filtered one and the hub's newest issues with no label filter — with the
second treated as "a second source for the floor and nothing more", because
whether the unfiltered listing indexes a just-created issue faster "is not
something the test suite can measure". The post-create check is the guarantee.

### `status` and `checkpoint` learn what a milestone issue is

[#200](https://github.com/radiusred/gh-codecrew/issues/200) reported a gate
raised by `checkpoint` on a milestone issue that `status` then printed as
`gates raised: none`. The capture offered two ways out — `status` lists such
gates, or `checkpoint` refuses a milestone ref and says where the gate belongs.
#210's Plan records the operator taking the first: "a gate on a milestone issue
is legitimate (a question about a requirement has no task to carry it), so
`checkpoint` keeps accepting a milestone ref and `status` shows the gate."
That choice has no Decision comment of its own; the Plan is where it survives.

[`status` reads the milestone issue's labels through the existing `Tracker.Task(ref)`](https://github.com/radiusred/gh-codecrew/issues/210#issuecomment-5540686227)
rather than a new tracker method: `Task` is a plain issue query that already
returns labels, so it serves a milestone issue as it serves a task.
**Trade-off:** one more API call per open milestone, and a `Task`-typed value
standing for a milestone issue at the call site, named `issue` there to say so.
**Rejected:** an `IssueLabels(ref)` method — "a second issue query for a field
`Task` already returns, plus a fake to carry it in every test that builds one".
The Deviation below is that rejection being taken up again the same afternoon,
for the other verb.

The second Decision in the same comment: `checkpoint` reads its target's labels
**before** it writes, so its comment and receipt can say what actually holds.
On a milestone issue they say `status` lists the gate beside the tasks' gates
until the label is removed, rather than that `task finish` refuses — which is
false there. **Rejected:** one neutral wording for both targets, which "would
have to leave out what blocks and where the gate is seen, which is the
information the person resolving it needs".

The task's Plan raised one ask-the-human point — should a milestone-issue gate
block `milestone close`? — and
[the operator answered it on the issue](https://github.com/radiusred/gh-codecrew/issues/210#issuecomment-5542315770):
"yes in principle, not in this task", captured as
[#219](https://github.com/radiusred/gh-codecrew/issues/219). **Trade-off:**
until then a milestone can close under an open requirement-level question;
`status` now shows it. **Rejected:** folding the close gate into PR #218 — "it
would ship a change no requirement names".

### `task finish` names the permission instead of the error

[#198](https://github.com/radiusred/gh-codecrew/issues/198) recorded
`task finish` dying on `GraphQL: Resource not accessible by integration
(…statusCheckRollup)` before its gate list. #211 carries four Decisions.

[The refusal names whichever permission is missing](https://github.com/radiusred/gh-codecrew/issues/211#issuecomment-5542391247),
by reading the GraphQL path out of gh's error: `statusCheckRollup` is
`checks: read`, and `checkSuite.workflowRun` is `actions: read` — the second
shape #198 recorded after `checks: read` had been granted. Both match only when
the message says `Resource not accessible by integration` *and* the path is
under `statusCheckRollup`. **Trade-off:** "parsing an error string is brittle by
nature — a future `gh` may reword it, in which case the raw error returns (the
pre-#211 behaviour), never a wrong refusal." **Rejected:** probing the App's
permissions up front, because `GET /apps/{slug}` reports the App's *default*
set, which the owner may have edited without the installation accepting it —
"exactly the case in the field".

[`actions: read` joins the manifest table](https://github.com/radiusred/gh-codecrew/issues/211#issuecomment-5542391477)
for the implementer, reviewer, qa and doc-synthesizer, because `gh pr checks`
reads `checkSuite.workflowRun` unconditionally and a freshly minted crew would
otherwise hit the new refusal on its first private hub. The coordinator set is
unchanged. **Rejected:** #198's other option — stop reading `workflowRun` —
which means replacing `gh pr checks` with a hand-written GraphQL query, against
the founding decision to wrap `gh` rather than speak the API.

[The `status` warning for a routed App missing a permission is out of the task](https://github.com/radiusred/gh-codecrew/issues/211#issuecomment-5542391662),
for the same reason the up-front probe was rejected: the warning "would go quiet
at exactly the moment the operator needs it", and a correct check needs the
installation's granted set, which requires an App JWT or org-admin scope.

The fourth came after review.
[The adapter-level test the reviewer asked for runs through a new seam in `internal/gh`](https://github.com/radiusred/gh-codecrew/issues/211#issuecomment-5543417926):
`var Command = exec.Command`, reassigned by tests only, with the fake gh being
the test binary re-entered. **Trade-off:** one exported variable in a package
that had none. **Rejected:** a PATH shim ("a file and a `/bin/sh` dependency for
what a re-entered binary does with neither") and an injectable runner interface
threaded through `tracker.GitHub` ("a constructor change for every caller to
cover four lines of wiring").

### The dispatch guidance states verified facts, not remembered ones

[#202](https://github.com/radiusred/gh-codecrew/issues/202) was filed by the
spoke during M9, after a codex-harnessed reviewer did a full pass it could not
post because its sandbox denied network. #212's
[first Decision](https://github.com/radiusred/gh-codecrew/issues/212#issuecomment-5540634202)
is the verification itself, not the wording: neither `codex --help` nor
`codex exec --help` names the two config keys, so they were checked with
`--strict-config`, which refuses unrecognised `-c` overrides — a bogus key
fails, `model_reasoning_effort=high` and
`sandbox_workspace_write.network_access=true` pass validation and fail only on a
deliberately bogus provider, and a non-boolean network value is rejected by
type. The seats' own rollout logs recorded the network flag in use and
`"reasoning_effort":null`. **Trade-off:** "the verification is against one CLI
version on one machine; the doc names the version so a reader with a different
one knows to re-check." **Rejected:** stating the key names from the operator's
notes alone — "#202's own comment is about a note that gave the wrong recipe".

[`.codecrew.yml` pins `model:` on the reviewer and qa rows](https://github.com/radiusred/gh-codecrew/issues/212#issuecomment-5540634466),
read from the CLI banner and the seats' rollout logs rather than assumed.
**Trade-off:** "a pinned default is a snapshot: when the CLI's default moves,
the file still names what the seats *are meant* to run under, which is the point
of the pin." The table gains no reasoning-effort key — the schema has none,
adding one is a protocol change, and effort is a dispatch-time property.
**Rejected:** a model slug with the level baked in, "the CLI models them as
separate settings, and the table should mirror the CLI".

[The routing-table examples in SPEC §5 and `docs/platform-interop.md` pin the codex model too](https://github.com/radiusred/gh-codecrew/issues/212#issuecomment-5540634679).
**Rejected:** leaving them generic, which "kept teaching the asymmetry the hub's
file just corrected: `model:` shown on one row and absent from the other two
reads as *codex seats do not take a model*, which is how the hub's table came to
be missing it"; and removing `model:` from every example row, since SPEC's block
is the example decided at the M6 gate and showing the field is how the doc says
it exists.

### The model the de-correlated seats run under, changed twice in one afternoon

Both changes are the operator's, on #207, and the second supersedes part of the
first.

[The first](https://github.com/radiusred/gh-codecrew/issues/207#issuecomment-5542305318),
at 14:57:37Z: reviewer and qa dispatches under codex pass
`-m gpt-5.6-sol -c model_reasoning_effort=high`. The occasion is recorded — the
implementer's work on #212 found the codex seats had been running at the model's
default effort, `low` for `gpt-5.6-sol`, "because nothing at dispatch set it",
and the reviewer of PR #216 that day had been dispatched before this and ran at
the default. **Trade-off:** "slower and costlier reviews for de-correlated
judgment that is actually exercised." **Rejected:** adding an effort key to the
`.codecrew.yml` schema — "the table declares which model a seat runs under; how
hard it thinks is a dispatch-time property of the harness".

[The second](https://github.com/radiusred/gh-codecrew/issues/207#issuecomment-5543428246),
at 16:23:22Z: the codex seats run `gpt-5.5`, not `gpt-5.6-sol`. The stated
reason is cost — "the reviews of PRs #215–#220 today ran on gpt-5.6-sol … and
consumed the operator's plan faster than the milestone can afford".
**Trade-off:** "a smaller model on the de-correlated seats for the rest of the
milestone, in exchange for finishing it." **Rejected:** dropping effort to
`medium` on `gpt-5.6-sol` — "the same family at lower effort is the more
expensive way to save." The Decision names its own follow-through, and #226 is
that task: the routing table "must say what the seats actually run under".

#226 recorded two Decisions of its own.
[The effort note now names both models' levels](https://github.com/radiusred/gh-codecrew/issues/226#issuecomment-5543798022)
— `gpt-5.5` defaults to `medium` and accepts `low`, `medium`, `high`, `xhigh`;
`gpt-5.6-sol` defaults to `low` and adds `max` and `ultra` — read from the CLI's
cached model catalog with its fetch timestamp recorded. **Rejected:** naming
only `gpt-5.5`'s four levels, since "a level list that was silently one model's
would mislead a dispatcher who passes `max` to gpt-5.5".
[The parser fixture in `internal/config/config_test.go` keeps its `gpt-5.6-sol` string](https://github.com/radiusred/gh-codecrew/issues/226#issuecomment-5543798273):
"a test that asserts a literal should not churn every time the hub's routing
table moves", and changing it "would put a Go touch in a PR that has no
behaviour change to test". The cost is one grep hit that is not a pin; the qa
verdict on M10-R5 names that exception explicitly.

### The README's routing table, shipped as a mirror and amended into an example

#213 put the hub's own `roles:` table in the README and recorded four
Decisions.
[The block goes in beat 2](https://github.com/radiusred/gh-codecrew/issues/213#issuecomment-5542985212),
"Four seats, always staffed", after the paragraph that explains what a seat is
— **rejected:** "Where it goes from here", because "a reader arriving there has
already been told about seats without ever being shown one".
[The coordinator row keeps its comment](https://github.com/radiusred/gh-codecrew/issues/213#issuecomment-5543000494)
— it "says in the file's own words the one thing a row cannot say for itself".
[`docs/introduction.md` does not get the block](https://github.com/radiusred/gh-codecrew/issues/213#issuecomment-5543000678),
because the introduction "explains the protocol in the abstract" and SPEC §5
already carries a generic table; adding it "wants its own capture", which is
[#227](https://github.com/radiusred/gh-codecrew/issues/227).

[The second Decision is the one the milestone reversed](https://github.com/radiusred/gh-codecrew/issues/213#issuecomment-5543000282):
the block is a curated copy "kept honest by a test", `routing_table_test.go`,
failing the build unless the fenced block was `.codecrew.yml`'s `roles:` section
byte for byte. Its recorded reasoning was "the same enforce-don't-remind stance
the CLI's refusals take", and its rejected alternative was transclusion, which
GitHub does not offer. The reviewer approved PR #223 with that test verified
both ways, and the operator
[replied on the task issue](https://github.com/radiusred/gh-codecrew/issues/213#issuecomment-5544027158)
at 17:15:26Z:

> This really isn't necessary. The important part for the docs is that it shows
> what is *possible*, not necessarily what we have in this repo. Other users
> will want very different routing tables. A table that shows an example of
> setting a different harness and a different model is fine.

M10-R6 was amended on #207 fifty-eight seconds later, quoting the original
clause; the task issue carries
[the adoption note](https://github.com/radiusred/gh-codecrew/issues/213#issuecomment-5544039631)
saying Decision 2 "is superseded by that amendment; the record will say so".
This is that.

#229 is the remedy, and
[its Decision](https://github.com/radiusred/gh-codecrew/issues/229#issuecomment-5544058429)
is the replacement gloss: "Here is a worked example: the `roles:` section of
this repository's own `.codecrew.yml`, as it stands today. … Your table will
look different, and every seat pointing at `~` is a complete one."
**Trade-off, in the seat's words:** "'as it stands today' is an explicit expiry
date on the block — a reader is told the copy may lag the file, which is weaker
than the byte-equality the deleted test enforced… In exchange the paragraph says
the true thing: the table is illustrative, and a reader whose project looks
nothing like this hub is not being shown a requirement." **Rejected:** dropping
the dating clause and letting the block float free of the file, and inventing a
fictional table — "the README's argument is 'this repository is agent-staffed,
and you can check', which a made-up example undercuts".

### The home page's table, and the rule from M8 that met a requirement from M10

M10-R7 put the same example in the codecrew.works home page's "The crew"
section, in the spoke.
[The block is copied from the hub README, not from `.codecrew.yml`](https://github.com/radiusred/codecrew-www/issues/19#issuecomment-5544239376),
"so the two pages tell one story", and nothing in the spoke checks it against
either. **Trade-off:** the page can go stale — the price the M10-R6 amendment
already accepted. **Rejected:** a byte-equality test against the hub, "the shape
#213 shipped and #229 deleted", which "would put a network fetch, or a second
repo, inside this repo's test suite".

[One line departs from the README byte for byte](https://github.com/radiusred/codecrew-www/issues/19#issuecomment-5544239585):
the coordinator row's comment, refitted from
`# the operator: this hub is coordinated by hand (SPEC §7)` (75 characters) to
`# a human: the operator` (41). The reason is a measured rule from the earlier
milestone —
[M8's 42-character ceiling for every code line on the page](8-a-product-home-page-for-codecrew-works.md#the-42-character-rule-measured-once-and-then-enforced-everywhere),
taken in headless Chromium at 420px — which the README has no equivalent of.
**Rejected:** keeping the line and letting the block scroll, and moving the
comment onto its own line, "which fits but invents a shape neither the README
nor `.codecrew.yml` has". Verified by rendering: at 420px the block's
`scrollWidth` equals its `clientWidth` in both colour schemes.

The third Decision recorded a collision between two rules, and was superseded.
[The seat narrowed the home page's no-crew-names test](https://github.com/radiusred/codecrew-www/issues/19#issuecomment-5544239830)
so the example block could carry the four real bot logins, recording the tension
plainly: M8's operator rule is that "the names are not part of the framework,
they are the radiusred crew", while M10-R7 and the capture it adopts "ask for
that table by name, so the rule and the requirement genuinely collide". The seat
took the later, more specific requirement, flagged the call for the operator on
the task issue and in PR #20, and named the remedy if the answer went the other
way — genericise the four identities and change nothing else.

[The operator decided against it](https://github.com/radiusred/codecrew-www/issues/19#issuecomment-5544432453):
"the M8 rule holds — crew members are not named on the product page, example
block included." The home-page table uses `coder-bot`, `review-bot`, `qa-bot`
and `doc-bot`; harnesses, models, the `~` coordinator row and its comment stay;
the test returns to its full form. **Trade-off, the operator's own:** "the
home-page block and the hub README's block no longer match line for line — the
README shows this hub's real table, the home page shows what a table looks
like." **Rejected:** amending the M8 rule for code blocks.

The seat then
[marked its own Decision superseded](https://github.com/radiusred/codecrew-www/issues/19#issuecomment-5544452216)
and recorded what changed in `ed97564`: the placeholder identities, a gloss that
no longer claims the block is this project's own file and instead points at the
hub README for the real table, the restored test, and one added structural
assertion — the five `identity:` values are exactly the four placeholders and
`~`, "the same rule stated from the other side, so the block cannot quietly
acquire a real login later". **The divergence is deliberate and recorded at both
ends:** the README carries this hub's real identities, the product page carries
placeholders.

## Deviations

Two Deviation comments, one from a seat and one from the operator.

**[`checkpoint` reads labels through REST, not through `Task`](https://github.com/radiusred/gh-codecrew/issues/210#issuecomment-5542909968).**
The reviewer's blocking finding on PR #218 was that `Tracker.Task` queries
GraphQL's issue-only `repository.issue(number:)`, which answers `NOT_FOUND` for
a pull request — so the label read the Decision had just introduced rejected PR
refs before the comment and label were written, and `roles/coordinator.md`
requires the pre-milestone gate to be recorded on the scaffold PR. The reviewer
demonstrated it against PR #218 itself: `issue: null` / `NOT_FOUND`, while
`issueOrPullRequest` identified the same number as a pull request. The fix takes
up the alternative the Decision had rejected, but only for the verb that needs
it: `checkpoint` gains `Tracker.IssueLabels(ref)` over the REST issues endpoint,
"which is what the comment and label writes already use"; `status` still reads
through `Task`, "where the target is always an issue". The regression fake
leaves `Task` undefined "so a regression to the GraphQL read would panic", and a
pull-request case is in the test. Fix commit
[5467568](https://github.com/radiusred/gh-codecrew/commit/5467568).

**[The reviewer dispatch on PR #228 stopped at the Codex usage limit](https://github.com/radiusred/gh-codecrew/issues/226#issuecomment-5543841707)**,
recorded by the operator: the session stopped after `go test ./...` had passed
and before any review was posted, on `You've hit your usage limit`. The stated
cause is the milestone's own consumption — "the eight reviews earlier today
(five on gpt-5.6-sol at high effort) consumed the plan; the operator moved the
seats to gpt-5.5 for exactly this reason". The Deviation records what did *not*
happen as clearly as what did: "No review exists on #228; the dispatch is
re-run when the limit resets. No other harness stands in: the routing table says
the reviewer runs under codex." The re-run approved PR #228 at 17:08:21Z.

## The gates

**#207's Gates section was left as the scaffold's placeholder** — the template
line "_What "done" means beyond CI: e2e suites, manual UAT, sign-offs._",
unedited. [M9's record](9-the-docs-at-codecrew-works.md#the-gates) says the same
of #201 and [M8's](8-a-product-home-page-for-codecrew-works.md#the-gates) of
#196. No `**Gate raised:**` or `**Gate resolved:**` comment exists anywhere in
this milestone, and no `cc:needs-decision` label was applied.

M10-R3's subject is what `status` does when such a label *is* applied to a
milestone issue, so the requirement was verdicted against the case the milestone
did not produce. The qa seat verdicted the requirement by reading `status`
live against #207 and finding the negative case: "an unmarked
`M10: Protocol bookkeeping from the field (radiusred/gh-codecrew#207)` header
and `gates raised: none`, matching the fact that no gate label is present",
adding "I did not raise a live gate." The positive path was verified from the
code and the tests instead.

What did gate the work: CI and the seats' own suites on every PR
(`go test ./...`, `go vet ./...`, `gofmt -l .` in the hub; `uv run pytest` and a
strict Zensical build in the spoke), an independent approval on each of the ten
PRs, and the QA verdicts at the close. Three questions a seat could not answer
were put to the operator and answered as Decisions on the issue that raised them
— the `milestone close` gate question on #210, the README cut on #213, and the
crew-names collision on codecrew-www#19 — none of them through `checkpoint`,
and each producing either a capture or a change to the work in flight.

## The review rounds

Ten PRs, **sixteen review submissions** by the reviewer role holder
(`radiusred-checky[bot]`), **five change requests**, all five resolved in the
next round.

- **PR [#215](https://github.com/radiusred/gh-codecrew/pull/215)** — approved
  first round. The reviewer checked that "no individual deleted filename is
  referenced at HEAD" and that `assets.go` embeds no `assets/` content.
- **PR [#216](https://github.com/radiusred/gh-codecrew/pull/216)** — one change
  request on a surface the new tests could not reach: `internal/cli/cli.go`'s
  top-level help still advertised `[--dry-run] (print the number and row it
  would get; create nothing)`, "the old behavior M10-R1 removes", and "the new
  `runMilestoneNew` test bypasses `cli.go` and therefore cannot catch this
  regression". Round two verified the fix with a test that "fails against the
  old ffd400c wording".
- **PR [#217](https://github.com/radiusred/gh-codecrew/pull/217)** — three
  rounds. The change request was a factual correction to the new dispatch
  guidance: the doc said an unset effort means the model's default, and the
  reviewer disproved it on the installed CLI by writing
  `model_reasoning_effort = "medium"` into a temporary config and watching the
  banner report `medium` with no `-c` at dispatch. "As written, the end-to-end
  dispatch guidance can tell an operator the seat ran at the model default when
  it actually inherited a local setting." Round two approved the corrected
  precedence; GitHub then dismissed that approval when the branch was rebased,
  and round three re-approved after checking a range-diff — "ff34e0f differs
  from the previously approved edfb29e only by main changes through a6dcd76".
- **PR [#218](https://github.com/radiusred/gh-codecrew/pull/218)** — the change
  request is the Deviation above. Round two approved at `12e4105`, noting the
  fake "omits `Task`" so the old read would panic.
- **PR [#220](https://github.com/radiusred/gh-codecrew/pull/220)** — a change
  request with two findings, both about tests that could not fail: reverting the four production lines wiring
  `MissingChecksPermission` into `PRInfo` "leaves every added test green and
  restores #198's raw GraphQL failure", and `TestBuildManifestPermissionsPerRole`
  "uses `rolePermissions` itself as its expected value", so removing the new
  grants "changes the oracle and the test still passes". The review also
  recorded a non-blocking permission-scope judgment with three external
  citations, concluding `checks: read` alone is insufficient for current
  `gh pr checks`. Round two approved both fixes.
- **PR [#221](https://github.com/radiusred/gh-codecrew/pull/221)** — the change
  request was a stale count in current documentation: `README.md` still said
  "four of twenty-nine" while the branch's introduction said thirty. The review
  named the cross-PR consequence without treating it as the defect — "whichever
  of #220 or #221 lands second will need to rebase and make the final count
  thirty-one" — which is what happened.
- **PR [#223](https://github.com/radiusred/gh-codecrew/pull/223)** — approved
  first round, having reproduced the drift test's failure in a scratch worktree
  with a one-character README change, and having checked the operator's
  later `gpt-5.5` Decision against the PR: "#226 explicitly carries that
  follow-up after #223 lands … so it is not a blocker for this PR."
- **PR [#228](https://github.com/radiusred/gh-codecrew/pull/228)** — approved
  first round, on the re-run after the usage-limit Deviation.
- **PR [#230](https://github.com/radiusred/gh-codecrew/pull/230)** — approved
  first round, "no findings", verifying that the block itself is unchanged and
  the drift test gone.
- **PR [codecrew-www#20](https://github.com/radiusred/codecrew-www/pull/20)** —
  approved first round against "M10-R7 as amended, task #19, adopted capture
  #18, and the M8 home-page rules", with the deploy-equivalent sequence run
  locally.

What the five change requests were about: two of them (#216, and both findings
on #220) were changes **not protected by a test that fails without them**; one
was a compatibility regression (#218); and two were defects in documentation
rather than in code — a false statement about Codex's effort precedence (#217)
and a stale refusal-code count (#221). Every approval in the milestone was
reached with the reviewer running the code: scratch worktrees, mutated copies,
live dry-runs against the hub, and a temporary Codex config.

## QA: two rounds, and a supersession recorded by the seat

The qa role holder (`radiusred-testy[bot]`) verdicted against merged `main` in
both rounds, running `go test ./...`, `go vet ./...` and `gofmt -l .` for each
requirement and then a "hairbrush" probe beyond the suite.

**[First round](https://github.com/radiusred/gh-codecrew/issues/207#issuecomment-5544060296)**,
17:18:27Z: **M10-R1 through M10-R6 all `satisfied`.** Three of the probes ran
live against this repository: `milestone new --dry-run` run
from the clone, printing "only the would-be issue/title/requirement IDs plus
`dry run: nothing written`" with tracked `git status` clean and no ROADMAP text;
two live dry-runs for M10-R2 that "created nothing, with newest issue unchanged
across the probes", one of them using a matching `M11:` title prefix and no
`--requirement` to check the prefix is stripped rather than doubled; and `status`
run live on #207 for M10-R3. For M10-R4 the verdict counted the refusal
catalogue by hand — "`docs/introduction.md` has 31 actual refusal-code bullets,
matching 31 unique production refusal codes" — and checked the docs permission
table against `internal/cli/identity.go` row by row. For M10-R5 it read the
installed CLI (`codex-cli 0.152.1`), re-ran the strict-config probes, and swept
`git grep -n 'gpt-5'` to confirm "only those current pins plus the recorded
`gpt-5.6-sol` parser-fixture exception". For M10-R6 it drove the drift test to
failure in a scratch worktree with a one-character README change.

**[Second round](https://github.com/radiusred/gh-codecrew/issues/207#issuecomment-5544580039)**,
18:06:55Z: **M10-R6 `satisfied`** as amended, and **M10-R7 `satisfied`**. The
seat states the supersession in its own words — "This R6 verdict supersedes my
round-one R6 verdict, which cited that now-deleted drift test" — and records
what the amendment cost: "I compared the README block to the
current `.codecrew.yml`; as a fact today it still reflects the hub table … but
no test enforces byte equality anymore." M10-R7 was verdicted against the
deployed site, "a live curl of `https://codecrew.works/`", checking the table's
position in the crew section, the surviving badge popovers and stylesheet rules
from M8, and that the live HTML "contains `coder-bot`, `review-bot`, `qa-bot`,
and `doc-bot`, and does not contain the radiusred crew App slugs in that page
table".

**The first round's chronology, stated exactly**, because #207's edit history
carries it. From that history and the comment timestamps: the operator's comment on #213 at 17:15:26Z; the R6 amendment edit
to #207 at 17:16:24Z; the remedy task #229 opened at 17:16:26Z; the round-one qa
comment posted at 17:18:27Z, two minutes and three seconds after the amendment,
with an R6 paragraph that assesses the drift test — the original clause; PR #230
opened at 17:18:36Z; the R7 amendment edit at 17:19:00Z; PR #230 merged at
17:24:27Z. The verdict was posted after the amendment landed and describes the
state of `main` before the remedy merged; nothing on the trail records when the
seat read #207.

## Requirement outcomes

The status column is the verdict word as the qa role holder wrote it, verbatim
and unqualified — for M10-R6, the word from the second round, which the seat
states supersedes its first. Everything else the milestone knows about a
requirement is in the columns around it.

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M10-R1 — the ROADMAP row belongs to the doc-synthesizer at both ends: `milestone new` no longer writes it, the record PR adds it Done, and the contracts, SPEC and guides say so (adopts [#197](https://github.com/radiusred/gh-codecrew/issues/197)) | [#208](https://github.com/radiusred/gh-codecrew/issues/208) / PR [#216](https://github.com/radiusred/gh-codecrew/pull/216) | `satisfied` | Task closed; first QA round, with a live `milestone new --dry-run` from the clone and a repo-wide ROADMAP sweep. This document's PR adds the M10 row rather than flipping one, which is the convention the requirement introduced |
| M10-R2 — `milestone new` never hands out a taken number: calls seconds apart get distinct numbers, and a collision is repaired after creation or refused with a code (adopts [#195](https://github.com/radiusred/gh-codecrew/issues/195)) | [#209](https://github.com/radiusred/gh-codecrew/issues/209) / PR [#221](https://github.com/radiusred/gh-codecrew/pull/221) | `satisfied` | Task closed; first QA round, with two live dry-runs against the hub that created nothing and a read of the repair path |
| M10-R3 — `status` lists a gate raised on a milestone issue beside task gates, `checkpoint` accepts a milestone ref, and the coordinator contract says where a requirement-level gate goes (adopts [#200](https://github.com/radiusred/gh-codecrew/issues/200)) | [#210](https://github.com/radiusred/gh-codecrew/issues/210) / PR [#218](https://github.com/radiusred/gh-codecrew/pull/218) | `satisfied` | Task closed; first QA round. Verified live on #207 in the negative — unmarked header, `gates raised: none` — with the positive path read from the code and the tests; the seat records "I did not raise a live gate" |
| M10-R4 — `task finish` refuses with a code naming the App and the missing permission instead of the raw GraphQL error, and `docs/identities.md` records the private-repo requirement (adopts [#198](https://github.com/radiusred/gh-codecrew/issues/198)) | [#211](https://github.com/radiusred/gh-codecrew/issues/211) / PR [#220](https://github.com/radiusred/gh-codecrew/pull/220) | `satisfied` | Task closed; first QA round, against #198's two verbatim error strings, a hand count of 31 refusal codes and a row-by-row check of the docs permission table against `internal/cli/identity.go`. The private-repo path itself was not exercised live — see [What the record does not contain](#what-the-record-does-not-contain) |
| M10-R5 — the dispatch guidance covers network reach and reasoning effort, and `.codecrew.yml` pins `model:` for the two codex-harnessed seats (adopts [#202](https://github.com/radiusred/gh-codecrew/issues/202)) | [#212](https://github.com/radiusred/gh-codecrew/issues/212) / PR [#217](https://github.com/radiusred/gh-codecrew/pull/217); re-pinned by [#226](https://github.com/radiusred/gh-codecrew/issues/226) / PR [#228](https://github.com/radiusred/gh-codecrew/pull/228) | `satisfied` | Tasks closed; first QA round, against the installed `codex-cli 0.152.1`, its model catalog and a `git grep -n 'gpt-5'` sweep naming the one recorded fixture exception |
| M10-R6 — the README shows a routing table where it explains the seats, as a worked example rather than a mirror (**amended** 2026-09-04; original clause quoted on #207) | [#213](https://github.com/radiusred/gh-codecrew/issues/213) / PR [#223](https://github.com/radiusred/gh-codecrew/pull/223); remedy [#229](https://github.com/radiusred/gh-codecrew/issues/229) / PR [#230](https://github.com/radiusred/gh-codecrew/pull/230) | `satisfied` | Second QA round, against the amended clause on merged `main` with `routing_table_test.go` gone. The first round's R6 verdict — also `satisfied`, but against the original clause and citing the drift test — is superseded, in the seat's own words |
| M10-R7 — the codecrew.works home page shows the same example routing table in "The crew" (adopts [codecrew-www#18](https://github.com/radiusred/codecrew-www/issues/18); **added** by amendment 2026-09-04) | [codecrew-www#19](https://github.com/radiusred/codecrew-www/issues/19) / PR [codecrew-www#20](https://github.com/radiusred/codecrew-www/pull/20) | `satisfied` | Task closed; second QA round, against the spoke's gates on `main` and a live fetch of <https://codecrew.works/>. The deployed block carries placeholder identities by the operator's Decision, and the gloss points at the README for the hub's real table |

One requirement was amended and one added after the milestone opened, both by
the operator on 2026-09-04, with the original R6 clause quoted on #207. None was
dropped.

**Captures adopted and closed by the merges — seven:**
[#197](https://github.com/radiusred/gh-codecrew/issues/197),
[#195](https://github.com/radiusred/gh-codecrew/issues/195),
[#200](https://github.com/radiusred/gh-codecrew/issues/200),
[#198](https://github.com/radiusred/gh-codecrew/issues/198),
[#202](https://github.com/radiusred/gh-codecrew/issues/202),
[#199](https://github.com/radiusred/gh-codecrew/issues/199) and
[codecrew-www#18](https://github.com/radiusred/codecrew-www/issues/18).
**Captures filed by the milestone — three, all open:**
[#219](https://github.com/radiusred/gh-codecrew/issues/219) (`milestone close`
ignores a gate raised on the milestone issue),
[#224](https://github.com/radiusred/gh-codecrew/issues/224) (cut the README down
and send readers to codecrew.works) and
[#227](https://github.com/radiusred/gh-codecrew/issues/227)
(`docs/introduction.md` could carry the hub's routing table). All three were
opened by the operator: #219 and #224 one second before the Decision comment
answering the question that produced each, and #227 thirty-two seconds after PR
#223 merged, following the #213 Decision that said the introduction's block
"wants its own capture".

## Protocol-discipline observations

Six things the milestone showed about the protocol itself, none of them a
restatement of the sections above.

- **The gap M9's record named first among its observations is closed, and this
  PR is written under the replacement.**
  [M9's first observation](9-the-docs-at-codecrew-works.md#protocol-discipline-observations)
  was that a spoke-only milestone has no way to put its ROADMAP row in a PR, and
  that #197 "has carried four candidate shapes since 2026-09-02 and has not been
  adopted"; it named the shape that would fix it — "the doc-synthesizer owning
  both ends of the row" — and noted that M9's own record PR "flips a row the
  operator wrote". M10-R1 adopted that shape, and this PR adds a row nobody
  wrote earlier: a milestone document named a protocol gap, and the next
  milestone adopted the capture it pointed at.
- **Two Decisions were superseded inside the milestone, and both supersessions
  were written by hand.** M9's record observed that "nothing in the protocol
  marks a Decision as superseded" and that a reversed Decision "still reads as
  standing". Here both reversals were followed by a comment saying so — on
  [#213](https://github.com/radiusred/gh-codecrew/issues/213#issuecomment-5544039631)
  after the amendment, and on
  [codecrew-www#19](https://github.com/radiusred/codecrew-www/issues/19#issuecomment-5544452216)
  under the heading `**Decision (superseded):**`. Neither is a protocol
  mechanism: the gathering verb still reads `**Decision:**` headings, and both
  notes are conventions the seats chose. The gap M9 recorded is unchanged; what
  changed is the practice around it.
- **A requirement can be amended between a QA run and its verdict.** The
  amendment edit landed at 17:16:24Z and the round-one verdict at 17:18:27Z,
  assessing the clause as it read before. Nothing in the protocol tells a
  running qa seat that the requirement it is verdicting has moved, and nothing
  invalidates a verdict written against a superseded clause — the seat itself
  handled it, in the next round, by naming which verdict supersedes which.
- **The hub lane ran one PR at a time because of `CHANGELOG.md`.** Four PRs
  were open at once at the peak, each adding a section at the top of
  `[Unreleased]` — the same lines — so each had to rebase onto whichever merged
  before it, and the three opened within seven minutes of each other at
  12:45:47Z, 12:49:03Z and 12:52:37Z merged in that same order at 15:21:58Z,
  15:46:18Z and 15:56:36Z. The costs are
  visible in the record rather than inferred: the reviewer of #217 had a
  standing approval dismissed by a rebase and spent a third round proving via
  range-diff that only `main` changes had arrived; the reviewer of #218 verified
  that "the latest rebase changed only main's CHANGELOG context"; and the
  reviewer of #221 found a stale refusal-code count and correctly identified it
  as a cross-PR ordering consequence rather than a defect in the branch, since
  #220 and #221 were both bumping the same number.
- **The doc-synthesizer seat held three of the ten delivery tasks.** #213, #229
  and codecrew-www#19 were started by `radiusred-wordy[bot]` under the
  implementer contract, the other seven by `radiusred-cody[bot]`. The seat that
  shipped the drift test also deleted it, and the seat writing this record wrote
  three of the tasks it describes. The protocol's non-doer review gate is
  unaffected — the reviewer is a different App on every one of the ten PRs — but
  nothing in the contracts says a role's holder may or may not be routed to
  another role's task, and nothing recorded the choice.
- **The refresh obligation found only a number, because the tasks fixed the
  README as they went.** The doc-synthesizer must keep the README and
  `docs/introduction.md` true at every milestone boundary. Here #216 corrected
  the two sentences claiming the roadmap names the open milestone, #209 and #211
  moved the refusal-code count to thirty-one in both files, and #213 and #229
  added and then reframed the routing-table block — so the only claim left stale
  at the boundary was the milestone count in the receipts. That is the
  obligation working from the other end: a claim a change falsifies was fixed by
  the change, on the recorded reasoning that "a claim the CLI change itself
  falsifies belongs with the change".

## What the record does not contain

Gathered here rather than left implicit, because each is a gap in the trail
rather than a gap in the work.

- **Two process events reached this document from the coordination layer and
  appear nowhere on the trail.** The dispatch brief for this task states that
  the implementer seat for #211 was stopped and re-dispatched before it recorded
  anything, because of permission prompts, and that one codex review run died
  with its shell wrapper and was re-run. Neither has a comment, a Deviation or a
  timestamp on any issue or PR. What the trail shows of #211 is one Plan edit,
  at 15:03:19Z by `radiusred-cody`, with no earlier comment on the issue — the
  same shape as #209 and #213, which were simply dispatched later in the lane.
  *Recorded as reported, not as verified:* nothing on the record distinguishes a
  first dispatch that wrote nothing from a task that started late, and nothing
  identifies which review run died, on which PR, or when.
- **The re-run of the #228 review is not explained against the limit message it
  quotes.** The Deviation quotes `try again at 9:13 PM` and was posted at
  16:58:28Z; the approval on PR #228 arrived at 17:08:21Z, well before that
  time. *Inferred, not recorded:* nothing says whether the limit reset early,
  whether a different account or plan was used, or whether the quoted time was
  approximate.
- **`NO_CHECKS_PERMISSION` was never exercised against a private repository.**
  PR #220's description says so in its note to the reviewer — "the detail's
  wording was not exercised live against a private repo in this session (the hub
  is public)" — and the reviewer agreed, treating it as out of reach rather than
  as a finding. The two error strings in the tests are the ones #198 recorded
  verbatim from the field. The end-to-end proof the PR names, "a live run on
  `radiusred/ops` as a seat still lacking `actions: read`", was not run, and no
  task or capture holds it.
- **The choice between #200's two options has no Decision comment.** The
  capture offered `status` listing milestone gates or `checkpoint` refusing a
  milestone ref, and the operator took the first; the only record of it is the
  sentence in #210's Plan that states it as already decided. *Inferred, not
  recorded:* that the reasoning is the one the Plan gives — a requirement-level
  question has no task to carry it — because no comment anywhere gives another.
- **The alternatives behind the milestone's scope are named but not weighed
  here.** The operator's first Decision on #207 records that the scope was
  chosen from three options and that #197's shape 3 was chosen over shapes 1 and
  2, and it names the two rejected scopes in a clause each; the reasoning that
  separated them lives on #197 and in the operator's judgement, not in a comment
  this document can quote.
- **No gate was declared or raised.** #207's Gates section is the scaffold's
  placeholder text, unedited, so the three operator answers this milestone
  needed — the `milestone close` question, the README cut and the crew-names
  collision — have no `**Gate raised:**` / `**Gate resolved:**` pair behind
  them. Each was asked in a task Plan's Ask-the-human section and answered as a
  Decision comment. `checkpoint` was not used in this milestone, on a task or a milestone issue,
  including for the requirement M10-R3 added it to.
- **The effort flag's application per dispatch is not evidenced.** The
  operator's Decision of 14:57:37Z says "every later dispatch in this milestone
  carries the flags", and #212's verification shows the seats' earlier rollout
  logs recording `"reasoning_effort":null`. No later rollout log or dispatch
  record is quoted anywhere, so the flags' use after that point rests on the
  Decision's own statement.
