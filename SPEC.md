# CodeCrew Protocol Specification

Version 2.0 — 2026-09-06

CodeCrew is a lightweight framework for agent-driven software delivery. It keeps
the auditability and reproducible discipline of heavyweight frameworks like GSD
while discarding their ceremony, and it uses the tools teams already have —
GitHub issues, PRs, and CI — as its state store and message bus instead of
maintaining a parallel documentation corpus.

The design rationale is recorded in [docs/founding-decisions.md](docs/founding-decisions.md)
and the analysis that motivated the project in
[docs/gsd-vs-frontier-orchestration.md](docs/gsd-vs-frontier-orchestration.md).

## 1. What CodeCrew is

Three parts, in order of importance:

1. **A protocol** — conventions for how project state is represented in GitHub:
   what a milestone, task, plan, decision, deviation, gate, and summary look
   like, and how agents and humans transact over them. This document is the
   protocol. Any agent that can read and write GitHub issues can participate.
2. **Role contracts** — short, harness-neutral prompt files defining what each
   role (implementer, reviewer, QA, doc-synthesizer) reads, writes, and is
   accountable for. Any frontier-model harness (Claude Code, Codex, Gemini CLI,
   …) can load them; orchestrators such as Paperclip assign them to agents.
3. **A thin CLI** — `codecrew`, a single static Go binary wrapping `gh`. It
   implements the workflow verbs so agents don't each reimplement the
   conventions, and it runs the deterministic gates as code, not judgment. It
   is the only real software in the framework, and it is deliberately small.
   Installed as a `gh` extension the command is `gh codecrew`; this document
   writes the bare `codecrew` when naming the protocol's verbs, and the
   contracts, the scaffold and the CLI's own output write the runnable form,
   because they are loaded verbatim as agent instructions.

CodeCrew does **not** dispatch or schedule agents. A human operator, a cron
job, or an orchestrator does that; CodeCrew defines what a dispatched agent
reads and writes.

## 2. Principles

The kernel of discipline CodeCrew enforces — everything else is left to the
model's judgment:

- **Externalized state.** All coordination state lives in GitHub, sharded by
  task. There is no global state file and nothing that two agents can merge-
  conflict over; GitHub serializes concurrent writes server-side.
- **Plan before nontrivial work.** A task's plan — intended changes,
  requirement IDs covered, pre-identified ask-the-human points — is written
  into the task issue before the first commit.
- **Atomic commits with greppable task references.** History is bisectable and
  progress is inspectable from reality (`git log --grep`), not from reports.
- **Gates run code, not judgment.** Build and full test suite (CI required
  checks) must be green before anything is called done. CodeCrew reads check
  status; it never defines CI.
- **The verifier is never the doer.** Every task is reviewed by a role other
  than the one that implemented it — ideally a different model on a different
  harness — because self-evaluation shares the blind spots of the work itself.
- **Human gates at designed points.** Ask-the-human points are pre-marked in
  plans, and any participant can raise one mid-task. A raised gate blocks
  progression until a human resolves it.
- **Decisions are captured at decision time.** Trade-offs and rejected
  alternatives are recorded in a structured comment the moment they are made,
  not reconstructed later.
- **Requirements have IDs**, so "done" is checkable against something.
- **Documentation is synthesis, not maintenance.** The only documents the
  framework produces are per-milestone "why" documents, compiled at milestone
  close from the decisions and deviations recorded during the work.
- **Intentional platform footprint.** The protocol knows which GitHub
  features it depends on and what plan tier they need (§5, Platform
  requirements); a dependency that would exclude free-plan adopters is taken
  knowingly or not at all.

## 3. Topology: hub and spokes

A CodeCrew project is a **hub** repository and one or more **spoke**
repositories.

**There is no single-repo mode.** A single-repo project is the degenerate case
where N=1 and the hub *is* the spoke. Representations are identical in both
cases, so growing from one repo to many requires no migration.

**The hub** holds: the roadmap, milestone tracking issues, synthesized
milestone documents (`docs/milestones/`), role contracts
(`.codecrew/roles/`), and the CodeCrew configuration.

**Spokes** hold: task issues, the PRs that implement them, and their own CI
gates. Task issues live in the repo whose code they change, so GitHub's native
traceability (closing keywords, PR linkbacks, CODEOWNERS, repo-scoped
permissions) works without convention.

**Every repo carries a pointer file, `.codecrew/config.yml`**, so an agent
dropped into any repo can find the coordination point. In the hub it declares
`hub: self`; in a spoke it names the hub (`hub: owner/repo`). Everything
CodeCrew owns operationally lives under `.codecrew/` — the pointer and the
role contracts — so the framework never competes for a name in a project's
own tree (`roles/` is Ansible's before it is ours). The human-facing record
does not: `ROADMAP.md`, `docs/milestones/`, `AGENTS.md` and `CLAUDE.md` stay
where readers and harnesses look for them.

**A spoke belongs to one hub.** The single `hub:` field is permanent — a
repo is in one delivery stream at a time, and the pointer is the answer to
"where is my coordination point", not a list. A task created *in* a spoke by
another hub's milestone is not the exception it looks like: such a task
resolves its hub through its parent milestone, not through the repo's
pointer, so the pointer keeps naming the stream the repo belongs to
([#177](https://github.com/radiusred/gh-codecrew/issues/177) carries the
resolution itself).

### Choosing a hub

The hub is wherever the pointer files say it is — any repo the project owns
can serve. The decision ladder:

1. **`hub: self`** for a single-repo project. The degenerate case above;
   nothing to decide.
2. **The org's `.github` repo** for a multi-repo org whose delivery record
   can be public. The hub's cargo — roadmap, milestone issues, milestone
   documents, role contracts — is org-level coordination state, which is what
   `.github` exists to hold. The caveat is visibility: `.github` repos are
   public by convention, and milestone issues carry goals, decisions, and
   gate discussions, so this choice publishes the delivery narrative even
   when spoke repos are private.
3. **A dedicated hub repo** when the tracking must stay private, or when one
   org runs several projects that should not share a delivery stream.

Three constraints make every rung workable:

- **One hub per spoke.** The pointer file has a single `hub:` field; a repo
  belongs to one delivery stream at a time. The field is blessed as
  permanent — see the pointer paragraph above for the cross-hub task.
- **One milestone number-line per hub.** Milestone numbering scans the hub's
  milestone issues, so a hub is one serialized stream of work. An org running
  several concurrent projects wants a hub per project, not interleaved
  numbering in a shared one.
- **Multiple hubs per org coexist.** All protocol state — milestone issues,
  numbering, labels, the roadmap — is hub-scoped. The only org-global pieces
  are the role App identities, and those are meant to be shared: the same
  implementer identity can serve every hub in the org.

**The hub is not the framework's repo.** `radiusred/gh-codecrew` is both the
distribution point for the CLI (`gh extension install radiusred/gh-codecrew`)
and this project's own hub — a dogfooding coincidence, not a pattern.
Adopters install the extension from it and point their `.codecrew/config.yml`
at a repo they own; the protocol requires write access to the hub (labels,
sub-issue attachment, milestone close), so a hub you don't control is not a
hub.

### Growth and restructuring

- **Adding a spoke:** create the repo, add a `.codecrew/config.yml` naming
  the hub. Nothing moves, nothing converts.
- **Cross-repo references** use the qualified form `owner/repo#123`. Short
  `#123` references remain valid forever within their own repo because issues
  never leave it.
- **Hub changes happen at milestone boundaries.** The protocol's only live
  state is the open milestone and its open sub-issues; everything else is
  inert record that nothing re-reads. So a hub move (splitting a project onto
  its own hub, extracting the hub from a repo that outgrew hub-is-spoke) is:
  close the current milestone in the old hub as normal, repoint each spoke's
  `.codecrew/config.yml` in a one-line PR, and start the next milestone in the new
  hub — a fresh number-line. History stays where it happened: closed
  milestone issues, milestone documents, and decision trails are
  point-in-time snapshots, and moving them rewrites the audit trail. Leave a
  tombstone row in the old hub's ROADMAP naming where the stream continued.
  Worst case this closes a milestone slightly early with requirements carried
  forward — cheap, and it keeps every milestone's record gathered from one
  hub and synthesized into one document.
- **Mid-flight transfer is the discouraged exception.** GitHub issue transfer
  preserves comments and leaves redirects, but drops label associations and
  is unverified for sub-issue links — moving an *open* milestone risks the
  record the boundary rule keeps intact. Reach for it only when a boundary
  close is genuinely impossible.
- **Repo splits:** open task issues are transferred with GitHub's issue
  transfer. Transfer drops label associations, so conventions must be cheap to
  re-apply (re-label, re-attach the sub-issue link to the milestone).

## 4. State model

### Roadmap

`ROADMAP.md`, committed in the hub. A short ordered list of the milestones
delivered so far, with a one-line goal each and links to their tracking
issues and milestone documents. Each row is added, already Done, by the
milestone's document PR (the doc-synthesizer's, §7): nothing writes a row
while a milestone is open — `milestone new` creates the tracking issue and
nothing else — so the roadmap lists finished milestones and `codecrew status`
reports the open one. The alternative, a row committed at open time, had no
PR to ride in when a milestone's tasks all lived in spokes
([#197](https://github.com/radiusred/gh-codecrew/issues/197)).

### Milestone

A **tracking issue in the hub**, labeled `cc:milestone`. This is the canonical
milestone object — GitHub's repo-scoped Milestone feature is never canonical,
even in single-repo projects. Body structure:

```markdown
## Goal
One paragraph.

## Requirements
- **M1-R1** — <requirement>
- **M1-R2** — <requirement>

## Gates
What "done" means beyond CI: e2e suites, manual UAT, sign-offs.
```

The milestone's tasks are attached as **GitHub sub-issues**, not listed in
the body: the tracking issue shows live per-task states and progress
natively, and nothing is hand-maintained — a checkbox list rots the moment
someone forgets to tick it (state is inferred, not bookkept).

Requirement IDs are `M<milestone>-R<n>` and are referenced by task plans and
PR descriptions. A requirement may span multiple tasks; it is only marked
complete when every task covering it is closed.

The grammar is enforced, not conventional: every ID under `## Requirements`
carries its own milestone's number. `milestone close` and `milestone
evidence` refuse `REQUIREMENT_ID_MISMATCH` on an ID that does not, naming
it, so no verdict is ever counted against a requirement belonging to
another milestone; `status` prints the same condition as a line, because it
reports the board rather than gating it. IDs written outside the section
are not requirements at all.

### Task

An **issue in the spoke whose code it changes**, labeled `cc:task`, attached
to its milestone as a sub-issue. Body structure:

```markdown
## Goal
What this task delivers.

## Requirements
M1-R2, M1-R4

## Plan
Intended changes, in enough detail that a deviation is detectable.

## Ask-the-human points
Pre-identified judgment calls that must stop for a human. "None" is valid.
```

The plan is written or updated by the implementer **before the first commit**.
Trivial tasks may have trivial plans; they may not have absent ones.

### Task lifecycle

State is inferred, not bookkept:

| State       | Signal                                        |
|-------------|-----------------------------------------------|
| Ready       | Issue open, unassigned                        |
| In progress | Assignee set                                  |
| Gated       | `cc:needs-decision` label present             |
| In review   | Linked PR open                                |
| Done        | PR merged, issue closed via closing keyword   |

The label set is deliberately tiny: `cc:milestone`, `cc:task`,
`cc:needs-decision`. Anything inferable from GitHub state is inferred.

### Decisions and deviations

Recorded as issue or PR comments **at the moment they occur**, using a
structured prefix so synthesis can find them:

```markdown
**Decision:** <what was decided>
**Trade-off:** <what was weighed>
**Rejected:** <the alternative and why>
```

```markdown
**Deviation:** <what differs from the plan>
**Why:** <rationale>
```

A decision made at a human gate uses the resolution form and is gathered as
a Decision:

```markdown
**Gate resolved:** <what was decided>
**Trade-off:** <what was weighed, if anything>
```

A deviation that changes what a requirement means, or an ask-the-human point
being reached, raises a human gate (§8) instead of being silently recorded.

The label may carry a parenthetical qualifier — `**Decision (superseding the
R1 counting claim):**` — which is gathered verbatim and carries no
semantics. Gathering is per labelled *paragraph*: a comment that is one
record is one record, and a record written after other text in the same
comment is still found. A label that does not open a paragraph — mid-line,
or at the start of a line inside a paragraph — is not a record, and any
other bold label opening a later paragraph (`**Finding 3:**`) ends the
record before it. One record per comment remains the practice to prefer.

The placement rule is the whole label family's, `**Gate raised:**` and
`**Gate resolved:**` included: a gate raised in a comment's second
paragraph is a gate, and a gate answered in a comment's second paragraph is
answered.

A QA verdict (§7) is read by the same reading. Supersession is per
*comment*: for each requirement ID the latest comment carrying a verdict
for it wins, and within that comment the first verdict for the ID counts —
so a comment may quote the verdict it supersedes without superseding
itself. Code is content, in records as in citations: a verdict inside an
inline code span or a fenced block — the form quoted in a contract, an
earlier verdict shown verbatim — is not a verdict, exactly as a URL in code
is not a citation (§6, `milestone evidence`).

### Summary

The PR description, finalized at merge time: what was done, which requirements
it satisfies, and any deviations (linking the deviation comments). The merged
PR *is* the task's summary artifact; no separate file is written. It is read
by the doc-synthesizer directly and never gathered as a record: records live
in comments, and a decision that exists only in a PR body is unrecorded.

### Commits

Atomic, one logical change each, every commit referencing the task issue
(`(#123)` suffix — the local number suffices because commits live in the same
repo as their task issue). Linear history (rebase merging) is recommended.

### Milestone document

`docs/milestones/<id>-<slug>.md`, committed to the hub when the milestone
closes. It is the "why" document: the architectural, pattern, and technology
choices made during the milestone, their trade-offs and rejected alternatives —
synthesized from the Decision and Deviation comments recorded during the work,
never reconstructed from raw history. It lands via a normal PR and passes the
same review gate as code.

## 5. Configuration

`.codecrew/config.yml` in every repo. Spokes need only the pointer; the hub
carries the full configuration. The `codecrew` field is the **protocol
version** — this document's version, naming the conventions the file speaks
— and is independent of the CLI release: `codecrew version` prints both, as
`v2.0.0 (protocol 2.0)`. The CLI implements one protocol major and checks
the pointer's on every verb that loads it: a different major is refused
(`refused[PROTOCOL_MISMATCH]`), and the two directions differ — a pointer
ahead of the binary asks for an extension upgrade, one behind it is told the
repo predates this protocol and is moved with `codecrew migrate`; neither
asks anyone to edit the version field, which describes the repo rather than
choosing for it. A missing field is assumed current, with a note. Decided at
the M6 gate on
[#114](https://github.com/radiusred/gh-codecrew/issues/114):

```yaml
codecrew: "2.0"
hub: self                # spokes: owner/repo

# Advisory role routing, read by whoever dispatches agents.
roles:
  implementer: { harness: claude-code, model: claude-fable-5, identity: app:my-org-coder }
  reviewer:    { harness: codex, model: gpt-5.5, identity: "team:my-org/review-crew" }
  qa:          { harness: codex, model: gpt-5.5, identity: app:my-org-qa }
  doc-synthesizer: { harness: claude-code, identity: user:alice }
  coordinator: { identity: ~ }   # the seat that dispatches the other four; ~ = the operator
```

Role routing is **advisory**: CodeCrew does not dispatch agents, so the config
is a contract for the orchestrator (or human) that does.

**Every role is always staffed — solo is a routing configuration, not a
reduced protocol.** The routing table says *who* holds each role: a GitHub
App (an agent acting as itself), a specific human, a team, or — with no
identity (`~`) — the human operator. An orchestrator dispatches
sub-agents or delegates to other harnesses per this table; a solo operator
embodies whatever routes to `~`. The hub's config should declare all five
roles at project onboarding — the four crew seats and the coordinator that
dispatches them (§7) — each routed explicitly; an orchestrator finding
no routing table should prompt for one rather than assume. `codecrew init`
(§6) scaffolds exactly this — the table with every role routed `~` — so
onboarding starts explicit. The CLI tolerates an absent table (every role is
then operator-held) but says so in its output; a table written before the
coordinator row existed still has a coordinator — the operator — and
`codecrew role coordinator` says so.

`identity` is **typed**: its value names the kind of GitHub principal that
holds the seat, and there are exactly four forms.

| Value | Who holds the seat |
|---|---|
| `~`, or no `identity` key | **The operator** — and any agent session acting under the operator's own auth. The grammar types the GitHub principal, not who is at the keyboard: a main agent with an App of its own is `app:<slug>` on the coordinator row, not `~`. |
| `app:<slug>` | **A GitHub App** — an agent acting as itself. Its viewer login is `<slug>[bot]`, which the protocol matches with or without the suffix. |
| `user:<login>` | **One named human**, by GitHub username. |
| `team:<org>/<slug>` | **Any member of the team** (child-team members included). |

The type prefix is not decoration: it is the one fact a routing table
alone can carry about a seat, and every rule below reads it rather than
guessing from the value's shape. An identity carrying no prefix — a 1.0
table's bare `my-org-coder` — is refused at load
(`refused[IDENTITY_UNTYPED]`, naming the row and the four forms), because
an App slug and a username are the same string and the protocol treats
them differently; `gh codecrew migrate` rewrites a 1.0 table by resolving
each login's kind through the API.

A **team-held** role is held by any member for every purpose an identity
serves — the holder review gate, verdict counting, the crew-identity
refusals, and branch rules
([#44](https://github.com/radiusred/gh-codecrew/issues/44)).
Latest-wins verdict supersession runs **across members**: one teammate can
supersede another's verdict — correct for a pool, and stated here so it is
never merely implied. A team-held reviewer seat composes with GitHub's
native required-review counts with zero framework code: repo rules demand
N approving reviews from the nominated team, GitHub enforces the count,
`task finish` merely reads the result. Platform footprint (per this
section's tier-intentionality obligation): teams are org-only — a
personal-account project cannot route to them and does not need to — and
required-review *counts* on private repos are plan-gated, so the
composition win is public-repo or paid-private territory. Each agent-staffed role gets its own app identity
so agent-authored work is attributable and distinct from the operator —
without this, every PR is authored by the operator's account and GitHub's
prohibition on self-approval makes the non-doer review gate unsatisfiable.
App private keys live outside the repo (by convention
`~/.config/codecrew/<slug>.*.private-key.pem`); short-lived installation
tokens are minted per invocation.

Where a verb gates on a role — `milestone close` counts only the qa role
holder's verdicts — "holder" means the routed identity, or, for an unrouted
role, any human identity not routed to a different role. Crew identities can
never stand in for a role they are not routed to.

**Crew identity** means `app:`-typed, and nothing else. It is the test the
verbs apply where an act belongs to a human — `--operator-confirm`,
`--bypass` — and a human who holds a seat (`user:`, or a member of a
`team:`) is not one: they hold the role *and* keep the operator's acts. A
1.0 table could not distinguish the two, so every routed login was refused
as crew.

### Identity tiers

App identities are optional infrastructure, required only where a
GitHub-native approval must come from a party with no human account:

1. **Solo** — one human, agents act under the operator's own auth. The whole
   protocol works on `gh auth login` alone; the review gate degrades to
   explicit operator confirmation (§6) because GitHub forbids self-approval.
   The confirmation must come from a human identity — a `[bot]` login or an
   `app:`-typed role holder is refused — and when the confirmer is also the PR author,
   the recorded comment states that no independent principal exists. See
   [docs/identities.md](docs/identities.md) for the operational path.
2. **Multi-human team** — a colleague satisfies the review gate natively; app
   identities add only attribution (distinguishing an agent's work from its
   dispatcher's).
3. **Enforced review with one human, or agent-reviews-agent** — app
   identities are required: the author and approver must be distinct
   principals, and each automated role needs its own.

Credential resolution is uniform across tiers, and `codecrew identity token
<slug>` (§6) is the act: orchestrator-injected env vars
— the App's id and private key under whatever names the platform binds
(`GITHUB_APP_ID` or `GITHUB_CLIENT_ID`; `GITHUB_PRIVATE_KEY` or `GITHUB_PEM`,
as PEM text or a file path), with the installation discovered from the App
itself (a supplied `GITHUB_INSTALLATION_ID` is a hint at most: the run on
Paperclip was handed a stale one — #119, findings 12 and 35) — then the
locally-held private key and credential stub `identity new` wrote. The verb
refuses with a code past those two; the operator's `gh` auth is the
identity only of an unrouted role, never a fallback for a routed one.

### Platform requirements

Everything the protocol *requires* is available on every GitHub plan,
public or private: issues and sub-issues, labels, comments, pull requests,
rebase merging, GitHub Apps, and CI check reading. A free-plan solo operator
loses nothing.

Features CodeCrew benefits from but deliberately does **not** require,
because they are plan- or visibility-gated on private repos:

- **Branch rulesets / required status checks and reviews** — free on public
  repos; private repos need a paid plan. Without them, `task finish` *is*
  the enforcement: it refuses on red checks and missing approval even where
  GitHub wouldn't block the merge button.
- **Auto-merge** — a convenience some hubs enable; `task finish` merges
  through the API and never depends on it.
- **Actions minutes** — private repos have quotas; the protocol reads check
  results but never defines CI, so a spoke's CI budget is its own affair.

The standing obligation: any future protocol dependency on a GitHub feature
states its plan availability — public/private differences included — at the
moment it is introduced.

## 6. Workflow verbs

The CLI's surface, and therefore the backend interface. Every verb is safe to
run by any role from any repo in the project (the pointer file resolves the
hub).

| Verb | What it does |
|------|--------------|
| `codecrew status` | Where the project is: open milestones, task states, raised gates — on tasks and on milestone issues alike, the latter marked `(milestone)` and on the milestone's own line; notes contract drift and a repo that does not delete branches on merge. With no open milestone it says so in place of the board and the gates section, and the two notes still print: both are local and have nothing to do with milestone state, and the quiet period between milestones is when a `.codecrew/roles/` fork gets reconciled against a new release. |
| `codecrew init [--hub owner/repo]` | Scaffolds a new repo: hub mode writes `.codecrew/config.yml` with the full `~`-routed roles table, the ROADMAP.md seed, the role contracts (embedded at the installed release) under `.codecrew/roles/`, each with a blank `.codecrew/roles/<role>.local.md` extension beside it (a comment saying what the file is for, pointing at §7 and the upstream examples page; comments-only composes to nothing), and an AGENTS.md entry point; spoke mode writes the two-line pointer. Then it commits exactly the files it wrote — a pathspec commit, so the operator's own staged and unstaged work is untouched — on the current branch, or on `codecrew-bootstrap` cut from the default branch when that branch requires pull requests (asked through `gh`; assumed when it cannot be asked), never pushing; it refuses a subdirectory (the pointer belongs at the root) and leaves a detached HEAD uncommitted with the command to run: the scaffold is the last commit before the protocol starts, and where a ruleset requires it, the scaffold PR is the one merge the operator does by hand, recorded as the pre-milestone gate (§8; #172). Idempotent — existing files are kept and reported, and a rerun that writes nothing commits nothing. Scaffolded contracts carry a provenance stamp naming the release that wrote them. `init` reads no pointer, so it is exempt from the protocol check — but not from the layout: a repo still on 1.x refuses `LAYOUT_LEGACY` naming `codecrew migrate`, rather than writing a second layout beside the first. |
| `codecrew milestone new` | Creates a milestone tracking issue in the hub from the template (`--dry-run` prints the number it would assign, the title and the requirement IDs, and creates nothing — so requirement prose can be written knowing the number); each `--requirement` (repeatable) becomes a bold-ID line under `## Requirements`, numbered M<n>-R1, R2, … in the order given — the section the close gate reads — and the IDs counted are printed; text that brings its own ID is refused. The CLI derives n, twice: before creating, as one past the highest `M<k>:` title across the hub's label-filtered milestone listing and its newest unfiltered issues — either listing alone can lag an issue created seconds earlier ([#195](https://github.com/radiusred/gh-codecrew/issues/195)) — and after creating, when both listings are read again and the number must be the new issue's alone; another issue already carrying the prefix has the new issue renumbered to the next free number, title and `M<n>-R<k>` IDs, printed as a `renumbered:` line (bounded; `refused[MILESTONE_NUMBER_TAKEN]` naming both issues and the hand fix when the repair fails or the number is still taken). A title carrying an `M<k>` prefix that disagrees is refused, one that agrees is stripped. Touches no file: the milestone's ROADMAP.md row is added, Done, by its document PR (§4). |
| `codecrew task new --milestone <id> --repo <spoke>` | Creates a task issue in the spoke from the template; attaches it to the milestone as a sub-issue. The milestone is resolved by number from the hub's open-milestone listing — and, when that listing lacks it, from the hub's newest issues regardless of label (an open issue titled `M<n>:` carrying `cc:milestone`), then again after a short wait, three reads in all: the label-filtered listing can lag a milestone created seconds earlier ([#234](https://github.com/radiusred/gh-codecrew/issues/234)), and a milestone found by either fallback is noted in the output. `refused[NOT_FOUND]` only after that. |
| `codecrew task start <ref>` | Verifies a plan is present, posts the `**Started by** @<login>.` record (and assigns the caller where GitHub allows — humans; App identities are not assignable) (refuses to start a planless nontrivial task), creates the working branch — unless the caller's role routing resolves to a role whose contract forbids commits (`qa`, `reviewer`), which get no branch. |
| `codecrew checkpoint <ref> --question "…"` | Raises a human gate: posts the question as a comment, applies `cc:needs-decision`. The ref is a task, or the milestone issue when the question is about a requirement and no task carries it (§8) — the comment and the receipt say which. |
| `codecrew identity new <role> --name <app>` | Mints the role's App identity via the GitHub App manifest flow: generates a manifest with the role's minimal permission set, hands the operator a one-click loopback URL, stores the returned private key locally, writes the role's routing into the hub's `.codecrew/config.yml` in the typed form (`identity: app:<slug>`; `--no-route` opts out, and the instruction printed when the table is not local carries the same form), and prints the remaining manual steps (install — per-account — and optional display polish). Webhooks off by default; `--with-webhook --webhook-url U` opts in for platform receivers (§9), subscribing `pull_request` and `pull_request_review` (`--events` names others, validated against the role's permissions) and, with `--webhook-secret S`, setting the receiver's secret as soon as the App exists — before it is installed anywhere, and repository events reach an App only through an installation, so the creation ping (signed with GitHub's generated secret, rejected by the receiver, harmless) is the only delivery that precedes it. |
| `codecrew identity webhook <slug> [--show] [--url U] [--secret S \| --rotate-secret]` | An active App hook under the App's own key: prints the URL, content type, whether a secret is set and the subscribed events; sets the URL and secret; rotates the secret and prints it once. An App minted without a webhook has no hook configuration and GitHub's API cannot create one — `refused[NO_WEBHOOK]` names the settings page where it is activated by hand; event subscriptions are not settable after creation either (no endpoint) — the verb prints the page. An App hook covers every repository its installation sees, so a platform needs no repository hooks. |
| `codecrew identity token [<slug>] [--installation <id>]` | Mints a short-lived installation token as the App: credentials from the environment under the names platforms bind (`GITHUB_APP_ID`/`GITHUB_CLIENT_ID`, `GITHUB_PRIVATE_KEY`/`GITHUB_PEM` as PEM text or a path), else the `~/.config/codecrew/` key and stub for the slug; signs the App JWT, discovers the installation from the App (a hinted id — the flag or `GITHUB_INSTALLATION_ID` — is used only when the App can see it; one installation is taken, several narrow to the hub's owner), and prints the token alone on stdout with a receipt on stderr. Never writes `gh`'s config. Refuses `NO_CREDENTIALS`, `BAD_CREDENTIALS`, `NO_INSTALLATION`, `INSTALLATION_AMBIGUOUS`. Runs from anywhere — it reads no pointer. |
| `codecrew roles diff <role>` / `codecrew roles show <role> [--latest]` | Contract tooling: `show` prints the contract a dispatched session loads — the hub's `.codecrew/roles/<role>.md` with its local extensions appended in §7 order (hub, then spoke); `show --latest` prints the contract embedded in the installed CLI whole. Drift: `status` reports when a local `.codecrew/roles/` contract differs from the embedded copy (scaffolded contracts carry a provenance stamp naming their release) and `diff` shows the divergence; `.codecrew/roles/<role>.local.md` files are never drift. Contracts are the project's own fork — reconciliation is a judgment routed through a task and PR, never an overwrite. |
| `codecrew role <name> [--login]` | Prints the typed identity holding a role — `app:<slug>`, `user:<login>`, `team:<org>/<slug>`, or `~` for the operator (§5). Script-consumable; resolves from the hub's routing table when run in a spoke. `--login` prints instead the handle GitHub will accept a review request for — the login of a `user:` seat, `<org>/<slug>` for a `team:` seat — and **nothing at all** for an `app:` or `~` seat, neither of which can be requested: the emptiness is the caller's whole decision. The implementer uses it that way at PR creation; App-held seats are dispatched instead, Apps not being requestable (CODEOWNERS-driven requests coexist — requested reviewers union). |
| *(every verb that reads the working repo's `.codecrew/config.yml`)* | Refuses `LAYOUT_LEGACY` when the repo is still on the protocol 1.x layout — a root `.codecrew.yml`, or a root `roles/` holding one of the five contracts, with no `.codecrew/config.yml` above it: the detail names what was found and `codecrew migrate`, and nothing reads the old layout. Refuses `PROTOCOL_MISMATCH` when the pointer's protocol major differs from the one the binary implements (§5), `IDENTITY_UNTYPED` when a routing row's identity carries no type prefix — the detail names the row and the four forms — and `GH_TOO_OLD` when the installed `gh` is below the floor the verbs need (2.50.0, for `gh pr checks --json`) — checked once, up front, so the gate never fails inside `gh`. A hub's routing table fetched from a spoke is advisory and is not checked. |
| *(every verb that reads a milestone's `## Requirements`)* | Refuses `REQUIREMENT_ID_MISMATCH` when an ID declared there is not the milestone's own — `M<milestone>-R<k>` is the grammar (§4) — naming every offending ID and the milestone read; `status` prints it as a line and carries on, reporting rather than gating. |
| `codecrew task finish <ref> [--dry-run]` | The gatekeeper (`--dry-run` evaluates every gate below in order and prints each — ok, refused with its code, not reached, not applicable — then the comments, merge and head deletion it would perform, writing nothing and exiting with the first refusal's code): refuses while the task carries `cc:needs-decision` (`refused[GATED]`) and, once the label is gone, while a gate raised on it has no answer (`refused[GATE_UNRECORDED]` — both labels are read per paragraph and only a later `**Gate resolved:**` answers, §4 and §8); verifies the caller is the seat that started the task (the `**Started by**` record `task start` posts on every start — accepted only from the login it names — else, for tasks started before the record, the assignee: the same login with the `[bot]` suffix ignored, or the same routed seat — a team-held role is any member; `refused[NOT_OWNER]` otherwise — the operator's own auth is not exempt; handover is `task start` again by the new seat, latest record wins, and `--bypass` is the recorded override for a human operator; a task with no start record is not gated), that a PR exists, CI checks exist and are green (`refused[NO_CHECKS]` when a PR reports zero checks — the deterministic gate cannot be satisfied by absence, and there is no override; `refused[NO_CHECKS_PERMISSION]`, naming the App and the permission, when the caller's installation token cannot read the checks at all — a private repo requires `checks: read` and `actions: read`, granted on the App's settings page and accepted on the installation), an approving review exists from the reviewer role's holder when the role routes to a distinct principal (`refused[NO_HOLDER_REVIEW]` otherwise — other approvals coexist but do not satisfy the gate; any non-doer approval suffices only when the role is operator-held) — then merges (rebase) and closes. When GitHub's own required-review rule is still unmet at that point (`reviewDecision: REVIEW_REQUIRED` — approvals count only from principals with write access: a write-access App's approval counts, a read-only App's and an operator confirmation do not), it refuses with `refused[REVIEW_NOT_COUNTED]` naming the supported paths; `--bypass` performs the ruleset's administrator merge instead, recorded as a PR comment, and only for a human operator the ruleset lists as a bypass actor. Refuses otherwise, with the specific unmet condition. In a solo-tier project (§5) where author and operator are the same principal, the non-doer approval degrades to an explicit operator confirmation, recorded as a PR comment; the confirming identity must be human — crew identities (a `[bot]` suffix, or an `app:`-typed routed role) are refused with `refused[SELF_CONFIRM]`; a `user:`- or `team:`-typed holder is a human and is not. After the merge it deletes the head branch — the counterpart of `task start` creating it; a deletion failure is a note, the merge stands. |
| `codecrew milestone evidence <n>` | Walks the milestone's record — tracking issue and every sub-issue, bodies and comments — and verifies every citation resolves (github.com references via the API under the caller's auth, everything else by HTTP). A citation is a URL in prose or in a Markdown link outside code; a URL inside an inline code span or a fenced code block is content — a probe target meant to be unreachable, a verbatim command or error string — and is not checked. A github.com citation that does not resolve is `refused[EVIDENCE_UNREACHABLE]`; an external one prints a `warning:` line and does not block, for QA to weigh. It also checks the milestone's own ID grammar before the walk (`refused[REQUIREMENT_ID_MISMATCH]`, §4): the record is read so QA can be dispatched against it, and a requirement belonging to another milestone is not this record's to verdict. Run by the coordination layer before dispatching QA, and by QA as its first act: uncommitted evidence cost M4-R4 its verdict, and the check is deterministic, so it runs as code. The milestone is resolved from the hub's milestone listing regardless of state — a closed one resolves too and is reported as closed before the citation report, link rot in a shipped record being what a maintainer reads the verb for; `refused[NOT_FOUND]` means no milestone carries that number, open or closed. `milestone close` and `status` keep their open-only reads. |
| `codecrew milestone close <id> [--dry-run]` | (`--dry-run`: the same gates in order, then every branch the sweep would delete or keep and why, and the closing comment; nothing written, the first refusal's code.) Verifies that the milestone issue itself carries no `cc:needs-decision` (`refused[MILESTONE_GATED]` otherwise — a requirement-level gate raised there by `checkpoint` is answered before anything is counted, §8), that all tasks are closed (`refused[OPEN_TASKS]`), that the tracking issue's `## Requirements` section declares at least one bold requirement ID (`refused[NO_REQUIREMENTS]` otherwise — IDs written elsewhere in the body are not requirements, so a close can never verify nothing) and that every ID it declares is the milestone's own (`refused[REQUIREMENT_ID_MISMATCH]` otherwise, naming them — §4's grammar, checked before anything is counted against a foreign requirement), and every requirement's latest QA verdict is `satisfied` (`refused[VERDICT_MISSING]` / `refused[VERDICT_UNSATISFIED]` otherwise; only verdicts from the qa role's holder count — its routed identity, or the human operator when the role is unrouted (§5) — and supersession is per comment: the latest comment carrying a verdict for an ID wins, the first verdict for that ID inside it counts, and a verdict written inside a code span or a fenced block is content, not a verdict, §4); once every gate has passed, sweeps the tasks' branches (the heads of their PRs, the `task/<n>-…` names `task start` cut) — deleting one only when its PR merged and the branch still sits at the merged commit, or when no PR is open and it carries nothing beyond the default branch; never a fork's branch or the default branch itself; reporting every other one — so the successful close and its closing comment record what was removed and why; gathers every Decision/Deviation comment across the milestone's tasks into raw material for the doc-synthesizer; refuses to close until the milestone document PR is merged. |

Verbs exit nonzero with a machine-readable reason when a gate blocks them, so
agents can act on the refusal rather than parse prose.

## 7. Roles

Role contracts live in the hub under `.codecrew/roles/`, one short markdown file each,
loadable by any harness (and referenced from `AGENTS.md` for harnesses that
read it natively). Roles are contracts, not accounts: no GitHub App needs to
exist for a role to be staffed — every role can act as the human operator
(§5, and [docs/identities.md](docs/identities.md) for both the solo path and
App creation). v1 roles:

- **implementer** — writes the plan into the task issue, does the work in
  atomic commits, records decisions/deviations as they happen, opens the PR,
  finalizes the summary. Never approves their own PR.
- **reviewer** — reviews the PR against the plan and requirement IDs; the
  conversation is ordinary PR review comments and requested changes. Ideally a
  different model/harness from the implementer, per role routing.
- **qa** — exercises the built thing against the milestone's gates and the
  requirements' intent (not just the tests the implementer wrote); reports
  findings as issue/PR comments.
- **doc-synthesizer** — at milestone close, compiles the recorded decisions
  and deviations into the milestone document and opens its PR.
- **coordinator** — the coordination layer as a seat: opens milestones
  (`--requirement`) and tasks, dispatches the four crew seats by the routing
  table, owns the review loop in both directions (reviewer on a PR,
  implementer on changes requested, the task's owner on approval), raises
  the gates only a human can answer, and drives `milestone evidence` and
  `milestone close`. It never writes code, reviews, verdicts or merges; its
  App holds contents: read, issues: write, pull requests: read and
  metadata, never more. Unrouted it is the operator — every project has a
  coordinator, solo included. The contract states what the orchestrator
  run taught the seat (#119, #164): one wake path per transition, state
  re-read at the act, execution events one-shot, dispatch on the platform
  and cite on GitHub, never the milestone number in requirement prose.

**Local extensions.** A project's own instructions for a role — house
style, local conventions, what its orchestrator injects — go in
`.codecrew/roles/<role>.local.md`, never into the contract. The contract is the
project's fork of the framework's (§6, `roles diff`); an extension is
append-only text loaded *after* it, so reconciling the contract against a
newer release never has to re-merge project additions, and `status`'s
drift check never sees them. Load order is fixed: the hub's
`.codecrew/roles/<role>.md`, then the hub's `.codecrew/roles/<role>.local.md`,
then the working repo's `.codecrew/roles/<role>.local.md` when it is a spoke. There is no merge
language and no precedence beyond that order — an extension that
contradicts its contract is a review finding, not a resolver's job.
`codecrew roles show <role>` prints the composition a dispatched session
should load; a harness that reads `AGENTS.md` natively follows the same
order by hand.

The inter-agent protocol is **GitHub itself** — issue comments, PR reviews,
labels. There is no other message bus, so any two harnesses interoperate by
construction.

## 8. Verification and gates

Three independent layers, attacking different failure modes:

1. **Deterministic gates** — each spoke's CI required checks. Owned by the
   repo, read by `task finish`. Catch what code can catch. A PR with zero
   reported checks refuses (`refused[NO_CHECKS]`, no override): this layer
   cannot be satisfied by absence, so every repo using `task finish`
   carries at least one `pull_request` workflow.
2. **Independent review** — the reviewer role holder's approval required to
   merge when the seat routes to a distinct principal; any non-doer approval
   only when the role is operator-held. The doer is fixed too: the seat that
   started a task is the seat that finishes it (`refused[NOT_OWNER]`
   otherwise; an operator overrides on the record with `--bypass`) — a
   merge by another seat is a misattribution even when the review gate
   passed (#165). The norm is a model review: a
   clean-context session under the reviewer contract, optionally a different
   harness — even in pure solo, where its findings land as a PR comment
   before the operator confirms. Catches correlated self-evaluation failure:
   the model grading its own work shares the blind spots of the model that
   did the work, and a briefed reviewer shares the briefing's.
3. **Human gates** — pre-marked ask-the-human points plus ad-hoc
   `checkpoint`s. Catch the boundary cases where capability doesn't help:
   "should this task be done as specified at all." A `cc:needs-decision` label
   blocks `task finish` until a human removes it, and the resolution must be
   recorded as a `**Gate resolved:**` comment (§4) so it is gathered into the
   milestone record as a Decision — `task finish` refuses
   (`refused[GATE_UNRECORDED]`) while a raised gate has no resolution comment,
   even after the label is removed. Both labels are read per paragraph (§4),
   and only `**Gate resolved:**` resolves: a `**Decision:**` on another
   subject answers nothing. Resolution is per gate — a `**Gate resolved:**`
   record resolves every gate raised before it and still open, so one
   comment may answer several questions, and never a gate raised after it. A question about a requirement has no
   task to carry it: it is raised on the milestone issue, where `status`
   lists the gate beside the tasks' gates, marked (#200), and
   `milestone close` refuses (`refused[MILESTONE_GATED]`) while the label
   is present (#219).

## 9. Environments

CodeCrew is harness-neutral by construction: participation requires only the
ability to run a CLI and read/write GitHub. Supported shapes:

- **Solo operator** — one human, one harness, hub-is-spoke.
- **Mixed-model team** — role routing assigns different models to different
  roles (e.g. Claude implements, Codex reviews); they converse through PRs.
- **Orchestrated company** (e.g. Paperclip) — the orchestrator maps its agents
  to CodeCrew roles via the routing config and dispatches them; CodeCrew
  defines what each one reads and writes. Exercised end to end on Paperclip
  (#119: three milestones on a proving-ground repo, the third driven by the
  App's webhook events with one gate and no other operator touch on the
  workflow; #164: a fourth cycle on a fresh repo, driven by a coordinator
  agent from the first event). The platform installs the coordinator as a
  seat like the others — `.codecrew/roles/coordinator.md` (§7) composed by
  `roles show coordinator`, the platform's wake syntax, ids and tooling in
  the project's `.codecrew/roles/coordinator.local.md`, its identity minted by
  `identity new coordinator` — instead of a hand-written brief, and wires
  each seat's App to the receiver that dispatches it (`identity new
  --with-webhook --webhook-secret`, `identity webhook`; one App hook covers
  every repository the installation sees). The whole shape — the separation
  of concerns, the seat mapping, credentials, wake paths, the onboarding
  checklist, what four cycles cost and the Paperclip recipe — is
  [docs/platform-interop.md](docs/platform-interop.md), which also names the
  seams that remain open.

## 10. The CLI

**What 2.0 broke.** 2.0 is a protocol major, and the break is the layout:
every CodeCrew-owned operational file moved under `.codecrew/` — the pointer
from `.codecrew.yml` to `.codecrew/config.yml`, the contracts and their
extensions from `roles/` to `.codecrew/roles/` — so the framework stops
competing for names in the root of a repo it does not own. There is no
compatibility shim and no dual-read: a 2.0 binary meeting a 1.x repo refuses
`LAYOUT_LEGACY` and names `codecrew migrate`, the one-shot verb that moves
the files and rewrites the pointer. Nothing else about a 1.x project
changes — the issues, labels, branches, records and roadmap are untouched,
and `ROADMAP.md`, `docs/milestones/`, `AGENTS.md` and `CLAUDE.md` stay at
the root where readers and harnesses expect them.

**What 1.0 promises** (decided at the M6 gate, #114). Within a major release
series of the CLI: verb names and their flags are additive — nothing is
renamed or removed; a refusal code's meaning is stable — codes may be added
in a minor, never repurposed, and removed only in a major; the
`refused[CODE]: detail` line and the `version` output are stable shapes,
other human-facing text is not; pointer fields are additive; the embedded
role contracts may change in a minor — `status`'s drift report and `roles
diff` are the mechanism, reconciliation the project's judgment. A change to
this document that invalidates existing pointers or recorded comments is a
protocol major, and the CLI that implements it refuses the old pointer.

- **Go**, single static binary (`CGO_ENABLED=0`), cross-compiled for
  linux/mac/windows. Also installable as a `gh` extension (`gh codecrew …`),
  since a gh extension is just a binary named `gh-codecrew`.
- **Wraps `gh`** in v1 — authentication comes for free. github.com only in
  1.0 and 2.0: the URLs the verbs build and recognise are github.com's; GitHub
  Enterprise Server is a non-goal (§12) until it is proven.
- **Backend interface** is shaped by the workflow verbs (§6), not by GitHub's
  feature set. The GitHub adapter is the only v1 implementation; the interface
  keeps two ports conceptually separate even though GitHub fuses them — the
  **tracker** (issues, labels, comments) and the **review surface** (PRs,
  reviews, checks) — so a future backend pairing a different tracker with a
  git host stays possible.
- Clear errors over mystery: if the hub declares a spoke the caller's token
  cannot see, the CLI says exactly that.

## 11. Auditability and mutability

Live tracker state is mutable, and CodeCrew does not police edits — teams that
rewrite the record either have good reason or pay the price. The durable audit
trail is what lands in git: atomic commits with task references, merged PR
descriptions, and the committed milestone documents, which are point-in-time
snapshots of the decisions that mattered. The immutable record is in git; the
working record is in the tracker.

## 12. Non-goals

CodeCrew deliberately does not have:

- wave orchestration, worktree manifests, dispatch scripts, or heartbeats —
  execution scaffolding ages badly; frontier models isolate and merge
  correctly from intent
- a global state file — state is sharded by task and mediated by GitHub
- agent dispatching or scheduling — that is the operator's or orchestrator's
  job
- a maintained documentation corpus — documents are synthesized at milestone
  boundaries, from records captured at decision time
- a Jira adapter — GitHub only, with the backend interface as the seam if
  that ever changes
- GitHub Enterprise Server, in 1.0 — github.com only; the verbs hardcode its
  URLs, and no GHES run has been made
- configuration for behavior the model can decide well — the framework
  manages trust (specification, verification, audit), not capability

## 13. v1 scope and open questions

**v1 delivers:** this protocol; the role contracts (four crew seats at 1.0, the coordinator's from 1.1); the CLI with the §6
verbs and the GitHub adapter; CodeCrew's own hub bootstrapped with it.

**Open questions, deferred until the bootstrap surfaces evidence:**

- Whether `task finish` should also verify requirement-coverage claims in the
  PR body against the milestone's requirement list, or leave that to review.
- An optional read-only GitHub Projects mirror as a human dashboard, derived
  from issue state.
- Issue/PR templates shipped as `.github/` templates versus created by the
  CLI at `new` time.
