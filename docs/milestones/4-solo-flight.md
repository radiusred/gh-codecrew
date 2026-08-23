# M4: Solo flight

Tracking issue: [#41](https://github.com/radiusred/gh-codecrew/issues/41) ·
Synthesized 2026-08-23 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the records
gathered by `codecrew milestone close 4` (six: two gate resolutions on #42,
Decisions on #45 and #49, a Deviation on #51, a gate resolution on #59), the
task plans and merged PR descriptions, the flight's findings log on #57, the
public test-drive record in
[davison/numberguess](https://github.com/davison/numberguess), and QA's two
rounds of verdicts on #41. The close itself lands when `milestone close 4`
runs after this document merges.

## Goal and outcome

Make CodeCrew adoptable by a stranger: a solo operator with no prior knowledge
takes a brand-new project from empty repo to first closed milestone following
only the documentation — fixing, on the way, what the docs audit found: a
verdict gate no pure-solo project could pass, no quickstart, and an
undocumented greenfield bootstrap. Delivered in full, and proven by walking
it: [davison/numberguess](https://github.com/davison/numberguess) went from
`gh codecrew init` to a closed first milestone driven by one human and a
Codex session that had never seen the protocol
([findings on #57](https://github.com/radiusred/gh-codecrew/issues/57#issuecomment-5382665394);
[SESSION_TRANSCRIPT.md](https://github.com/davison/numberguess/blob/main/SESSION_TRANSCRIPT.md)).

Along the way the milestone reframed what "solo" means — a routing
configuration, not a degraded tier
([#42](https://github.com/radiusred/gh-codecrew/issues/42) / PR
[#43](https://github.com/radiusred/gh-codecrew/pull/43)) — and grew the CLI by
three verbs: `init` ([#49](https://github.com/radiusred/gh-codecrew/issues/49)
/ PR [#50](https://github.com/radiusred/gh-codecrew/pull/50)), `role`
([#45](https://github.com/radiusred/gh-codecrew/issues/45) / PR
[#46](https://github.com/radiusred/gh-codecrew/pull/46)), and `version`
([#55](https://github.com/radiusred/gh-codecrew/issues/55) / PR
[#56](https://github.com/radiusred/gh-codecrew/pull/56)). The quickstart
([docs/first-milestone.md](../first-milestone.md),
[#48](https://github.com/radiusred/gh-codecrew/issues/48) / PR
[#52](https://github.com/radiusred/gh-codecrew/pull/52)) walks the stranger
into every refusal deliberately. Four releases shipped — v0.2.0 through
v0.2.3 — two of them mid-flight remediation cuts, on the principle that the
test drive and QA must exercise the binary the docs tell a stranger to
install.

## Decisions

- **Solo is a routing configuration, not a degraded tier — every role is
  always staffed.** M4-R1 was framed as "how does the verdict gate get
  satisfied when no `qa` identity is routed?", and the gate raised on #42
  offered three mechanisms: implicit human verdicts, a per-requirement waiver
  marker (the implementer's recommendation, matching the #38 self-confirm
  precedent), or a close-time flag. The operator rejected all three as
  framed: each treated solo as role-absence. Instead, the qa role's routed
  identity defines whose verdicts count, and routing supports an App slug, a
  human's GitHub username, or — `~`/absent — the human operator themselves,
  per SPEC §5's existing "acts as the human operator" rule. A solo operator's
  verdicts count because they hold the role, not because the gate was
  waived; the same model covers teams (roles routed to different humans or
  agents). Trade-off: loses the explicit per-requirement waiver
  acknowledgement — independence becomes a configuration choice, visible in
  the routing table rather than restated in each verdict.
  ([Gate resolved on #42](https://github.com/radiusred/gh-codecrew/issues/42#issuecomment-5379578178);
  delivered by PR [#43](https://github.com/radiusred/gh-codecrew/pull/43) as
  `config.HoldsRole`, with crew identities never able to hold a role they are
  not routed to.)
- **Role routing is required at project onboarding.** Follow-on to the
  above: the hub's `.codecrew.yml` should declare all four roles explicitly
  (App, username, or `~`), and an orchestrator finding no table should prompt
  for one rather than assume defaults; the CLI stays tolerant of an absent
  table but says so. Trade-off: one config block in the greenfield path,
  versus silent defaults that hide who holds which role. The prompt's
  materialization was deferred to M4-R2/R3 — and became the `~`-routed table
  `init` writes.
  ([Decision on #42](https://github.com/radiusred/gh-codecrew/issues/42#issuecomment-5379587723).)
- **`codecrew init` earns its place; docs-only bootstrap rejected.** M4-R3
  asked the question explicitly. The operator decided for the verb: hub mode
  scaffolds the pointer file with the full `~`-routed roles table, the
  ROADMAP seed, the role contracts embedded in the binary at build time, and
  an `AGENTS.md` entry point; spoke mode writes the two-line pointer;
  idempotent. Trade-off: a verb to maintain and embedded contracts to ship,
  versus three paste blocks before the first verb works — accepted for
  first-run experience and so onboarding changes have one home instead of
  smearing across quickstart, SPEC, and identities docs.
  ([Decision on #49](https://github.com/radiusred/gh-codecrew/issues/49#issuecomment-5380794790);
  PR [#50](https://github.com/radiusred/gh-codecrew/pull/50). The quickstart
  plan on [#48](https://github.com/radiusred/gh-codecrew/issues/48) records
  the consequence: its paste-blocks goal was superseded by the verb.)
- **Holder resolution is a CLI verb (`codecrew role <name>`), not a
  parse-the-YAML convention.** From a pointer-only spoke the YAML is not
  local; the CLI already owns hub-fallback resolution (memoized since #42);
  and a deterministic, script-consumable query is the orchestrator-friendly
  seam (`--reviewer $(codecrew role reviewer)`). Rejected: convention-only
  (breaks in spokes, every harness reimplements it) and folding it into
  `status` output (not script-consumable). The same PR routed the hub's
  reviewer role to `davison` by name — expressible only since #42 let
  identities be usernames — and added the implementer's obligation to
  request review from the role holder at PR creation, coexisting with
  CODEOWNERS.
  ([Decision on #45](https://github.com/radiusred/gh-codecrew/issues/45#issuecomment-5379993748);
  PR [#46](https://github.com/radiusred/gh-codecrew/pull/46), the first PR
  to request its reviewer this way.)
- **A checkless PR refuses at `task finish` (`refused[NO_CHECKS]`), with no
  override.** Flight finding 5 reproduced on a throwaway repo
  ([evidence on #59](https://github.com/radiusred/gh-codecrew/issues/59#issuecomment-5382684061)):
  `gh pr checks` on a checkless PR exits 1 with nothing on stdout, so `task
  finish` was dying with a raw gh error — an accident that happened to
  block. The "vacuous pass on zero checks" path M3 worried about
  (radiusred/www#41 finding 2) was unreachable. The gate raised on #59
  offered refuse-always, warn-and-proceed, or refuse-with-recorded-override;
  the operator chose Option A: a repo with no checks has no deterministic
  gate, so the gate must not be satisfiable by absence (SPEC §2, "gates run
  code"), the fix is a ten-line workflow, and the one observed refusal in
  the wild produced exactly that workflow unprompted. Trade-off: docs-only
  and tiny repos must carry at least one `pull_request` workflow to use the
  verb — a repo that truly never wants CI merges outside it, visibly, rather
  than being handed a waiver for the only deterministic gate.
  ([Gate resolved on #59](https://github.com/radiusred/gh-codecrew/issues/59#issuecomment-5383079436);
  PR [#61](https://github.com/radiusred/gh-codecrew/pull/61). Settles the M3
  finding in the strict direction.)
- **Code touched in a task ships with tests in the same PR; no backfill
  tasks.** Set by the operator in review on PR #46 ("code is building up but
  tests aren't…"), folded into `roles/implementer.md` in the same PR, and
  cited as "the #46 convention" by every code PR after it. No Decision
  comment exists — this is a reviewer-set rule, visible only in the PR's
  review thread and the contract diff.
  ([Changes requested on PR #46](https://github.com/radiusred/gh-codecrew/pull/46#pullrequestreview-5000040680);
  commit `ef7d980`.)
- **Release parity before QA dispatch.** Not a gathered record but a
  principle stated in the #51 plan and applied three more times: the test
  drive and QA must exercise the released binary the quickstart tells a
  stranger to install, so the quickstart merged only after v0.2.0
  ([#51](https://github.com/radiusred/gh-codecrew/issues/51)), v0.2.1 cut
  finding 0 mid-flight ([#55](https://github.com/radiusred/gh-codecrew/issues/55)),
  v0.2.2 preceded QA dispatch
  ([#62](https://github.com/radiusred/gh-codecrew/issues/62)), and v0.2.3
  preceded the re-dispatch
  ([#65](https://github.com/radiusred/gh-codecrew/issues/65)). Recorded in
  the plans and the `**Release verified.**` evidence comments, not as a
  Decision.

## The flight: davison/numberguess

M4-R4's evidence anchor ([#57](https://github.com/radiusred/gh-codecrew/issues/57)):
a brand-new project, its own hub, taken from `gh codecrew init` through a
full milestone lifecycle by one human and a Codex session running the
released v0.2.0 binary, following only the scaffolded artifacts. The
operator's curated transcript —
[SESSION_TRANSCRIPT.md](https://github.com/davison/numberguess/blob/main/SESSION_TRANSCRIPT.md)
— preserves the prompts, responses, and commands verbatim.

**What the flight proved**
([findings comment](https://github.com/radiusred/gh-codecrew/issues/57#issuecomment-5382665394)):
a foreign harness with zero prior context drove the whole protocol from the
scaffold alone. It read `AGENTS.md`, `.codecrew.yml`, and `ROADMAP.md`
unprompted straight after `init`, explained the four roles and the `~`
routing when the operator said the roles were new to them (and recommended
keeping every role operator-held), planned before code, recorded a
well-formed [Decision](https://github.com/davison/numberguess/issues/2#issuecomment-5382529302)
and [Deviation](https://github.com/davison/numberguess/issues/2#issuecomment-5382550066)
at the moment each occurred, posted six verdicts in the mandatory format,
created a doc task when `DOC_MISSING` refused, and closed through every gate
— [milestone #1](https://github.com/davison/numberguess/issues/1), tasks
[#2](https://github.com/davison/numberguess/issues/2) and
[#4](https://github.com/davison/numberguess/issues/4), PRs
[#3](https://github.com/davison/numberguess/pull/3) and
[#5](https://github.com/davison/numberguess/pull/5), and the project's own
[milestone document](https://github.com/davison/numberguess/blob/main/docs/milestones/1-playable-number-guessing-game.md).
It was also the first live fire, by an independent harness, of the pure-solo
confirm (viewer == author under #42's bot-guard) and of the solo verdict
path (qa unrouted → operator-held, six verdicts counted at the close gate).
No protocol concept needed explaining; every finding was edge ergonomics.

**Friction findings and their fold-back**
([close comment on #57](https://github.com/radiusred/gh-codecrew/issues/57#issuecomment-5385201768)):

0. **The binary carried no version.** Found before the first verb: the
   release tag reached `scripts/build-extension` and was never embedded, so
   a stranger could not tell what they had installed. Fixed mid-flight —
   `codecrew version`, stamped via ldflags; the quickstart's install step
   gained verification and the manual-update truth (`gh extension upgrade`
   is always the user's act). ([#55](https://github.com/radiusred/gh-codecrew/issues/55)
   / PR [#56](https://github.com/radiusred/gh-codecrew/pull/56); v0.2.1.)
1. **`init` ran happily in a non-repo.** The agent discovered the missing
   git repo by its own inspection and back-filled `git init` + `gh repo
   create`. Now `init` detects the absence and leads its next-steps with the
   bootstrap command — guide, not refuse, since the scaffold is wanted either
   way. ([#58](https://github.com/radiusred/gh-codecrew/issues/58) / PR
   [#60](https://github.com/radiusred/gh-codecrew/pull/60).)
2. **`codecrew: "0.1"` read as a tool version.** It is the protocol version
   (SPEC 0.1), independent of the binary (v0.2.x); the operator had to ask.
   The scaffold now says so in a comment, and SPEC §5 in prose. (#58 / PR #60.)
3. **No supported way to edit milestone requirements → a raw `gh api` PATCH
   with an empty shell variable momentarily wiped the issue body.** A
   near-miss data loss, visible in the
   [transcript](https://github.com/davison/numberguess/blob/main/SESSION_TRANSCRIPT.md).
   The doc half landed in the quickstart (#58 / PR #60) — then had to land
   again (see QA below); a requirements verb remains a separate question.
4. **`task start` created the linked branch remotely and said nothing about
   the local checkout.** The implementation commit landed on local `main`
   and needed surgery. `task start` now prints `locally: git fetch && git
   switch <branch>`. (#58 / PR #60.)
5. **Checkless-repo `task finish` behaviour was accidental.** The agent's
   response to the accidental refusal — add a real CI workflow, recorded as
   its Deviation — was exactly what a designed gate should nudge toward. Made
   deliberate as `refused[NO_CHECKS]` (the #59 decision above; PR #61).

## QA: two rounds under the enforced gate

Before dispatch, two preparations were recorded on #41
([prep notes](https://github.com/radiusred/gh-codecrew/issues/41#issuecomment-5385271973)):
v0.2.2 was released so QA would verify the shipped binary, and
davison/numberguess was made public — because **App installations are
per-account**: the org-installed testy could not have read a private repo
under the operator's personal account without making the App
public-installable, a second installation, and a `codecrew-token` change to
resolve non-org installations. Recorded as input to
[#47](https://github.com/radiusred/gh-codecrew/issues/47) (the no-App-fleet
onboarding decision): the App fleet does not cross account boundaries
without ceremony.

**First round**
([verdicts](https://github.com/radiusred/gh-codecrew/issues/41#issuecomment-5385304745)):
`radiusred-testy[bot]` returned one satisfied, two not satisfied, one
untestable — and for the first time the remedy loop ran under the verdict
gate as code rather than on trust (M2's loop predated it; M3 passed clean).

- **M4-R1 satisfied** — `TestHoldsRole` from a clean build covers
  routed-App, routed-human, unrouted-human, routed-other-role, and `[bot]`
  exclusion; #42 records the pre-implementation gate; numberguess#1 closed
  on six verdicts from its unrouted human qa holder.
- **M4-R2 not satisfied** — QA executed the quickstart verbatim against
  v0.2.2 and found the requirements-editing command was not runnable as
  written: `gh issue edit 1 --body-file <(...)` invokes a literal `...` and
  preserves no body. The web-UI alternative worked, but M4-R2 promises one
  complete command path.
- **M4-R3 not satisfied** — v0.2.2's `init` scaffolded correctly (seven
  files, idempotent rerun, pointer-only spoke mode), but its non-git guidance
  said `git init && gh repo create … --source=. --push` before any commit
  existed, so `--push` could not complete as written.
- **M4-R4 untestable** — the lifecycle was independently confirmed (closed
  milestone and tasks, merged green PRs with solo confirmations, six
  verdicts, five findings, fold-back PRs merged, `NO_CHECKS` resolved), but
  the cited `SESSION_TRANSCRIPT.md` returned 404: it had only ever existed in
  the operator's working tree, so the "following only the docs" execution
  path could not be verified.

**Remedies:** the two doc defects became
[#63](https://github.com/radiusred/gh-codecrew/issues/63) / PR
[#64](https://github.com/radiusred/gh-codecrew/pull/64) — a fetch-edit-put
sequence for requirements (`gh issue view … > milestone.md`, edit, `gh
issue edit --body-file milestone.md`), and `init`'s guidance reordered to
commit before `gh repo create --push`; the transcript was committed to
numberguess@main (`2ab727a4`); and v0.2.3 was cut
([#65](https://github.com/radiusred/gh-codecrew/issues/65)) so the
superseding verdicts would test the shipped string, not `main`.

**Second round**
([superseding verdicts](https://github.com/radiusred/gh-codecrew/issues/41#issuecomment-5385339708)):
all three superseded to satisfied against released v0.2.3 — QA executed the
repaired fetch step and obtained the full milestone body; ran the corrected
`init` guidance through init, stage, commit, and the documented create step;
and reconciled the committed transcript's command-and-dialogue sequence with
the live numberguess record end to end.

QA changed no files, commits, or branches, per its contract.

## Deviations

One recorded, on [#51](https://github.com/radiusred/gh-codecrew/issues/51#issuecomment-5381243676):
the v0.2.0 release was not tag-and-verify as planned. The release workflow
failed because `cli/gh-extension-precompile@v2` sets up its own Go (1.18,
`GOTOOLCHAIN: local`) and ignores the workflow's `setup-go` step; the `init`
verb was the first shipped code to use `slices` (Go 1.21+), so v0.2.0 was
the first release build to hit the pin — invisible to PR CI, which used
`setup-go` correctly. Fix: `go_version_file: go.mod` on the precompile
action (PR [#53](https://github.com/radiusred/gh-codecrew/pull/53)); the
failed tag was deleted and re-pushed, and the
[release verified](https://github.com/radiusred/gh-codecrew/issues/51#issuecomment-5381277519)
from the installed extension. The numberguess flight's own Deviation (the CI
workflow added when `task finish` refused) is recorded in that project, not
this one.

The three operational release tasks (#51, #62, #65) had no PR and were closed
directly with evidence comments, as their plans said they would — `task
finish` cannot close a no-PR task (`NO_PR` by design). The #51 plan named
the unused linked branch and the no-PR close as protocol friction for ops
tasks; it remains unaddressed.

## Requirement outcomes

| Requirement | Delivered by | Outcome |
|-------------|--------------|---------|
| M4-R1 — a pure-solo project passes the verdict gate | [#42](https://github.com/radiusred/gh-codecrew/issues/42) / PR #43 | Done; QA satisfied, first round |
| M4-R2 — the first-milestone quickstart | [#48](https://github.com/radiusred/gh-codecrew/issues/48) / PR #52; remedy [#63](https://github.com/radiusred/gh-codecrew/issues/63) / PR #64 | Done; not satisfied → satisfied (v0.2.3) |
| M4-R3 — greenfield bootstrap documented or automated | [#49](https://github.com/radiusred/gh-codecrew/issues/49) / PR #50; remedy #63 / PR #64 | Done; not satisfied → satisfied (v0.2.3) |
| M4-R4 — the test drive, findings folded back | [#57](https://github.com/radiusred/gh-codecrew/issues/57); fold-back #55 / PR #56, #58 / PR #60, #59 / PR #61 | Done; untestable → satisfied (transcript committed) |

Supporting work with no requirement ID: the `role` verb and reviewer-request
obligation (#45 / PR #46), the v0.2.0 release and toolchain fix (#51 / PR
#53), and the v0.2.2 and v0.2.3 parity releases (#62, #65).

## Protocol-discipline gaps observed

- **Runnability auditing caught what review didn't.** Both not-satisfied
  verdicts were doc commands that read correctly and failed when executed —
  one of them (the `--body-file` placeholder) had passed review on PR #60,
  itself a fold-back of a flight finding about that very step. QA's method
  of executing the docs verbatim is the audit M4-R2 promised, and it found
  defects that a reading review and the flight both missed. Candidate: the
  reviewer contract could ask for doc commands to be executed, not read.
- **The record can exist only in a working tree.** M4-R4's untestable
  verdict was correct: the transcript the findings comment cited was never
  committed. Nothing in the protocol checks that cited evidence is reachable.
  Candidate: a QA-prep step that resolves every evidence link before
  dispatch.
- **Ops habit outran protocol on #62.** The v0.2.2 tag was pushed moments
  before the task's plan was written
  ([close comment](https://github.com/radiusred/gh-codecrew/issues/62#issuecomment-5385273759))
  — harmless, self-reported, and corrected on #65 ("in protocol order this
  time"). Worth recording because release tasks have no PR and therefore no
  gate: nothing but habit enforces plan-before-act on them.
- **The App fleet stops at the account boundary.** The finding lives in a
  prep comment on #41 and names #47 and `docs/identities.md` as its
  destinations; #47 carries no comment and `docs/identities.md` does not yet
  mention it. Recorded here so it is not lost a second time.
- **Two decisions of record have no Decision comment.** The tests-ride-along
  convention (set in review on PR #46) and release parity (stated in plans
  and evidence comments) both shaped the milestone and were both invisible
  to the gatherer — M3's "judgment calls live in plans and PR bodies" gap,
  with a new variant: review threads.
- **Gates recorded as designed, for the first time under enforcement.** The
  positive observation: both checkpoint gates this milestone (#42, #59) were
  raised before implementation, resolved with trade-off and rejected
  alternatives in the prescribed form, and gathered at close — and the
  remedy loop resolved through the gate's own rule — a later verdict from
  the qa holder supersedes an earlier one — with the first round's verdicts
  standing in the record as `VERDICT_UNSATISFIED` would have read them,
  rather than being edited away. The discipline M2 ran on trust and M3 built
  as code held under its first real load.
