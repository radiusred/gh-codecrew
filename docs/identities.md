# Identities: running solo, staffing a crew

Roles are contracts, not accounts. A seat is held by exactly one of four
kinds of principal: **you**, the operator (`~` in the routing table — the
whole protocol runs on `gh auth login` alone); **a named human** (a
username); **a GitHub team** (`identity: org/team-slug` — any member holds
the seat); or **a GitHub App identity** minted for a crew member. Nothing
requires an App to exist before a role can be staffed — App identities are
infrastructure you add when you want them (attribution) or need them
(enforced independent review). This document covers the ends of that range:
the zero-setup solo path, and minting a proper App identity for a role —
with the advice on which to choose in between.

The conceptual model — identity tiers, credential resolution order — is in
[SPEC §5](../SPEC.md); this is the operational companion.

## Running solo

Prerequisites: a GitHub account, `gh` 2.50.0 or later (`task finish` and
the close's branch sweep read `gh pr checks --json`, which older `gh`
lacks; `gh --version`) authenticated (`gh auth login`), the extension
installed (`gh extension install radiusred/gh-codecrew`), and a
`.codecrew.yml` in your repo (`hub: self` for a single-repo project).

Everything works, because solo is a routing configuration, not a reduced
protocol: every role is always staffed, and in pure solo *you* hold each
unrouted one. You (or agents running under your auth) create milestones and
tasks, write plans, commit, open PRs, raise and resolve gates. You also
perform the qa contract yourself: post per-requirement verdicts on the
milestone issue in the standard `**M1-R1 — satisfied.**` form — the close
gate counts the qa role holder's verdicts, and unrouted, that is you.
Declare the routing table in your hub's `.codecrew.yml` at onboarding (all
five roles — the four crew seats and the coordinator that dispatches them —
`~` for the ones you embody) — `gh codecrew init` scaffolds exactly
this, along with the roadmap seed and role contracts; an absent table works
but the CLI will nag.

Exactly one thing degrades: the non-doer review gate. GitHub forbids
approving your own PR, and in pure solo there is no second principal to
approve — so `task finish` accepts `--operator-confirm` instead, and records
an explicit confirmation comment on the PR. When the confirming operator is
also the PR author, the recorded comment says so in as many words: *no
independent principal exists in this project*.

Two honesty notes:

- The confirmation must come from a **human** identity. An identity carrying
  `[bot]`, or one routed to a role in `.codecrew.yml`, is refused
  (`refused[SELF_CONFIRM]`) — agents can never waive review, in any tier.
- Pure solo means no independent verifier, and the record shows it. That is
  the honest floor, not a failure mode; climb a tier when the review gate
  needs teeth.

## Minting a crew member

**Who should climb this way, and who shouldn't:** guided App creation is
the natural progression for a solo operator whose project is growing in
sophistication — agents earn attributable seats one at a time — and close
to a requirement for orchestration platforms (a Paperclip-style platform
drives every seat through App identities it owns). It is usually the wrong
first move for a *human* team adopting CodeCrew for isolated milestones or
tasks: humans are already distinct, attributable principals whose approvals
GitHub accepts, so route those roles to usernames — or to a GitHub team
(`identity: org/team-slug`, any member holds the role; SPEC §5 has the
semantics and the footprint) — and mint Apps only for the seats agents
actually fill. Agents reading this:
before recommending `identity new`, check whether the seat in question is
held by a human — if it is, routing beats minting.

When a role should act as itself — attributable work, and an approver
distinct from the author — it gets its own **GitHub App**. One App per role,
named for the crew member, not the role (`myorg-coder`, not
`myorg-implementer`: identities outlive role reassignments).

The framework's path here is guided creation via the App manifest flow
(decided at the gate on
[#68](https://github.com/radiusred/gh-codecrew/issues/68)):

```
gh codecrew identity new reviewer --name myorg-reviewy
```

builds a manifest carrying the role's minimal permission set (the table
below), serves it as a one-click local URL, and — once you confirm the
creation on GitHub — stores the returned private key under the
`~/.config/codecrew/` convention, routes the role in the hub's
`.codecrew.yml` for you (run it in the hub; `--no-route` skips, and on a
pointer-only spoke it prints the routing line instead), and prints what
stays manual: installing the App (per-account — see step 4 below) and,
optionally, giving it the crew logo — the manifest has no avatar field and
no API uploads one, so Display information on the App's settings page is
where a logo lands. If you want the CodeCrew mark for your own crew, it is
in the hub's [`assets/`](../assets/) — `codecrew-mark.png` is avatar-sized.
The App is owned by the hub's account unless `--owner` says otherwise, and
the verb refuses an App named after a role — crew-member names only. The
alternatives weighed at the gate — a published manifest template with no
verb, docs-only tier-sharpening, fine-grained PATs from a second account —
are recorded there; a shared App published by radiusred stays off the table
by prior decision, because the framework must never be key custodian for an
adopter's org.

Webhooks stay off by default — a crew App acts, it never listens — with an
opt-in `--with-webhook --webhook-url <receiver>` for platform users, whose
orchestrators watch through webhooks on the identity Apps they already own
(the watch seam, [#54](https://github.com/radiusred/gh-codecrew/issues/54)).
It subscribes the App to the two transitions a platform routes to seats —
`pull_request` and `pull_request_review` — and `--events` names others
(`issues` for a receiver that wants gates; an event the role's permissions
cannot receive is refused before the browser opens). GitHub generates the
hook secret; it is printed once and never written to disk — or, with
`--webhook-secret <the receiver's>`, replaced by the receiver's as soon as
the App exists — before its first protocol delivery; the creation `ping`
GitHub sends is signed with the generated secret and fails the receiver's
check, harmlessly, since a ping fires nothing (finding 60 below). Mint
with the hook if the App will ever listen: an App minted without one has
no hook configuration GitHub's API can create — see "The receiver side".

When the guided flow can't serve (no browser at hand, an enterprise
quirk), the manual ritual it automates — one App per role:

1. **Create the App, owned by the org** (Org → Settings → Developer settings
   → GitHub Apps → New GitHub App). Homepage URL can be anything; deactivate
   webhooks — CodeCrew's Apps only ever act, they never listen.
2. **Grant the minimum repository permissions** for the role's contract.
   `Metadata: read` is always required; beyond that:

   | Role            | Contents      | Issues        | Pull requests  | Checks |
   |-----------------|---------------|---------------|----------------|--------|
   | implementer     | Read & write  | Read & write  | Read & write   | Read   |
   | reviewer        | Read          | Read & write  | Read & write   | Read   |
   | qa              | Read          | Read & write  | Read & write   | Read   |
   | doc-synthesizer | Read & write  | Read & write  | Read & write   | Read   |
   | coordinator     | Read          | Read & write  | Read           | —      |

   Add `Workflows: read & write` to the implementer if it will ever touch
   `.github/workflows/` files — `Contents` alone cannot push those.

   The reviewer row's `Contents: read` is least-privilege and deliberate —
   but GitHub counts approvals only from write-access principals, so a
   read-only reviewer App's approval never satisfies a required-review
   rule. Granting it `Contents: write` makes its approvals count, held by
   an identity whose contract forbids editing code: an auditable trade of
   privilege for an agent-gated merge, and plainly trusted-agent territory
   — which is why it is never the default. Per-project operator's choice,
   made explicitly: at minting with
   `identity new reviewer --with-approval-permission`, or mid-life by
   bumping the App's Contents permission in its settings — an escalation
   GitHub itself gates behind the installation's approval of the change.

   The qa row's `Contents: read` is just as deliberate: the seat files
   what it finds and never fixes it, so it needs no branch and no push —
   and the permission enforced that mechanically when an orchestrator
   mapped "write the tests" onto the qa agent, whose push failed 403
   ([#119](https://github.com/radiusred/gh-codecrew/issues/119), finding 14). Tests ride the
   implementer's PR. Do not grant qa write access to make a mis-mapping
   work; re-map the work.

   The coordinator row is the smallest of all: the seat opens milestones
   and tasks, comments, labels and raises gates — `Issues: read & write`
   is the whole of its writing — and reads PRs only to learn where the
   review loop stands. It never pushes, reviews or merges, so `Contents`
   and `Pull requests` stay read, and it reads gate results through the
   verbs, not the checks API. A platform binds this App's credentials to
   the agent that runs `roles/coordinator.md`
   (`identity new coordinator --name <crew-member>` mints it); solo, the
   seat is unrouted and the operator's own `gh` auth is the identity. The
   orchestrator run's coordination layer had no identity at all and read
   every seat's credentials through its own 401
   ([#119](https://github.com/radiusred/gh-codecrew/issues/119), finding 16).
3. **Generate a private key** and store it outside any repo. Convention:
   `~/.config/codecrew/<app-slug>.<date>.private-key.pem`.
4. **Install the App on the org**, scoped to all repositories or at least to
   the hub and every spoke the role will touch. A spoke the App cannot see
   is a spoke the role cannot serve. Installations are **per-account**: an
   App installed on your org sees nothing under any other account — a
   personal-account repo included — until the App is made public-installable,
   installed on that account too, and tokens are minted against the right
   installation. The fleet does not cross account boundaries without
   ceremony (found the hard way in M4's QA prep,
   [#41](https://github.com/radiusred/gh-codecrew/issues/41)).
5. **Route the role** in the hub's `.codecrew.yml`:
   `roles.<role>.identity: <app-slug>`.

### Acting as the App

Mint a short-lived installation token per invocation:

```
export GH_TOKEN=$(gh codecrew identity token <app-slug>)
```

The verb reads the env vars an orchestrator binds first (SPEC §5: the App
id and private key under the platform's names — `GITHUB_APP_ID`/
`GITHUB_CLIENT_ID`, `GITHUB_PRIVATE_KEY`/`GITHUB_PEM`, PEM text or a file
path — with the slug then optional), else the key and credential stub
`identity new` wrote under `~/.config/codecrew/` for the **full** App slug
(`myorg-coder`, not `coder`). It discovers the installation from the App
itself — a supplied `GITHUB_INSTALLATION_ID` or `--installation <id>` is
used only when the App can see it, and a stale one is reported and
overridden (the platform handed an agent a stale id once; #119 finding 35).
It prints the token alone on stdout, a receipt on stderr (`minted for
<slug> (App <id>, from <source>): installation <id> on <account>`), never
touches `gh`'s config, and refuses with a code when it cannot mint:
`NO_CREDENTIALS` (nothing to sign with — the detail names what it looked
for and how to write the stub by hand), `BAD_CREDENTIALS` (GitHub rejected
the JWT: key and id disagree, retrying will not help), `NO_INSTALLATION`
(install the App — step 4), `INSTALLATION_AMBIGUOUS` (several accounts;
pass the id). It runs from anywhere — no `.codecrew.yml` needed — and the
`scripts/codecrew-token` the 1.0 contracts pointed at is now a one-line
wrapper around it. Then:

- `GH_TOKEN=$tok gh codecrew …` / `GH_TOKEN=$tok gh …` for API actions.
- Push over HTTPS as `https://x-access-token:$TOKEN@github.com/owner/repo.git`.
- Commit author is the App's bot user:
  `<slug>[bot] <UID+<slug>[bot]@users.noreply.github.com>`, where UID comes
  from `gh api 'users/<slug>%5Bbot%5D' --jq .id`.
- Per session, and only as `GH_TOKEN`: never `gh auth login` with an
  installation token and never persist it — a one-hour token in a shared
  `hosts.yml` is the next agent's 401. A 401 means the token expired; mint
  again before anything else. Platforms: give each dispatched agent its own
  `GH_CONFIG_DIR` (or `HOME`), so no run starts from another identity's
  dead token ([#119](https://github.com/radiusred/gh-codecrew/issues/119), findings 10 and 22).
- A crew App cannot create a repository: creating one under an org takes
  `Administration: write`, which no role's permission set carries and no
  contract asks for. The operator creates the hub and each spoke — with
  any branch rule — before the crew's first act; the scaffold then lands
  as the operator's commit, or as the project's first PR when the default
  branch is protected (#119, findings 3 and 4).

### Dispatching a role session

A crew App acts when dispatched; what varies by tier is who notices there
is work and starts the session. An App minted without `--with-webhook`
never watches — nothing polls on its behalf, and the operator (solo) or an
orchestrating session dispatches by hand, as the numberguess flight did.
An App minted `--with-webhook` is watched *by its platform*: the receiver
named at minting gets the protocol-traffic events (a PR opening among
them) and the platform dispatches the role session itself — no human act,
and nothing for the implementer to do in any tier: dispatch always belongs
to the coordination layer, never to the doer.

The tiers differ only in where that layer lives. Solo, it is the same
session changing hats: implementer while building, coordination layer when
the review is due — and the fixed dispatch prompt below is what keeps the
hat-switch honest, because the briefing channel is a standing template,
not implementer-authored content. On a platform, the hats are separate
processes with platform-assigned identities: every session is dispatched
*as* a role, none is the operator's primary session, so the scaffold's
dispatch authorization never applies to them — the platform watches the
webhooks and does all dispatching itself. What matters for the dispatch,
however it is triggered:

- **The contract, composed.** What the session loads is the hub's
  contract plus the project's local extensions (`roles/<role>.local.md`,
  hub then spoke — SPEC §7): `gh codecrew roles show <role>` prints exactly
  that, so a dispatch prompt can point at one command instead of a file
  list.
- **Fresh context.** The role session starts clean — no implementation
  conversation, no shared scratch state. Independence is contextual as well
  as credential: an approval means little if the approver watched the code
  being written. For a sole rider, the floor is the same model in a separate
  session holding nothing but the role contract, the task, and the diff;
  better, when available, is a different model or harness — `harness:` in
  the routing table is the declared intent, and de-correlated judgment is
  the point of splitting the seats.
- **QA dispatches start with reachable evidence.** Before dispatching the
  qa role, the coordination layer runs `gh codecrew milestone evidence <milestone number>` —
  citations that 404 cost a requirement its verdict once and the check is
  deterministic, so it runs as code, not hope. The qa contract refuses past
  an unreachable record from its side too. Brief the seat for what the
  implementer could not have produced: an assessment of the shipped tests
  against each requirement and a probe the suite does not enumerate — not
  a rerun of `npm test` (the contract's hairbrush).
- **Credentials.** `gh codecrew identity token <slug>` — the env-var path
  from SPEC §5 first, the local key and stub second, a refusal third (the
  verb above). Beside the private key, `identity new`
  writes a credential stub (`~/.config/codecrew/<slug>.json` — the App ID
  and client ID, nothing secret): with it, the helper mints the App's JWT
  and discovers the installation itself, touching no user credential. For
  an App minted before the stub existed, write it by hand —
  `{"slug":"<slug>","app_id":<numeric App ID>}` — the account-lookup
  verb refuses `NO_CREDENTIALS` without it and says exactly that. The
  session exports `GH_TOKEN` and is the App. To confirm *which*
  App, read the verb's stderr receipt (the slug, App id, installation and
  account it minted for) and compare the App ID with
  `gh api /apps/<slug> --jq .id` — the public App record, no token needed. `GET /app` does not work here:
  it wants the App's own JWT and answers an installation token with a 401
  ([#139](https://github.com/radiusred/gh-codecrew/issues/139)). After posting, check the review's
  `user.login` is `<slug>[bot]`.
- **The reviewer specifically** is never requested through GitHub's
  reviewers field — Apps are not review-requestable, and the implementer
  contract says to stand down rather than gate on it. The dispatched session
  reads the PR cold and posts its review directly:
  `GH_TOKEN=$tok gh pr review <n> --request-changes --body …` and, when
  satisfied, `--approve`. Whether that approval also *satisfies a
  required-review count* turns on the App's access: write-access Apps
  count, read-only Apps do not (see Known quirks below).

A dispatch prompt proven in the field (the numberguess flight,
davison/numberguess#8 — where the isolated reviewer caught a real encoding
crash the implementer's tests had missed, refused approval until it was
fixed, then verified the exact reproduction before approving):

> Act as the CodeCrew reviewer for `<repo>` PR #N. Read AGENTS.md and the
> hub's roles/reviewer.md first and follow them exactly. You are not the
> implementer and must not edit code. Inspect the PR diff BEFORE its
> description, then the task and milestone issues. Authenticate as
> `<app-slug>` (mint a token; never print it — use it only as GH_TOKEN, on
> the same command line as each gh call). Confirm the identity: the App ID
> in ~/.config/codecrew/<app-slug>.json equals `gh api /apps/<app-slug>
> --jq .id` (not GET /app, which 401s under an installation token), and it
> differs from the PR author. Review correctness, plan, premises, tests,
> and consequences. Submit an ordinary GitHub PR review as the App: approve
> only if sound, otherwise request changes with concrete findings. Then
> confirm the review's author is `<app-slug>[bot]`. Report your verdict and
> evidence back. Do not merge or run task finish.

### The receiver side

An App webhook is one per App, and GitHub delivers it for **every
repository the App's installation covers** — an installation on "All
repositories" includes repositories created later. Installations are per
account (step 4 above): an org's installation delivers the org's
repositories, a personal account's its own, and the same App installed on
both delivers both. So a platform needs no repository webhooks at all; the
two hand-pasted repository hooks of the orchestrator run's fourth cycle
([#164](https://github.com/radiusred/gh-codecrew/issues/164)) were the
workaround for seats whose Apps had no receiver, not the design.

Each seat's App points at the receiver that dispatches that seat, and
subscribes only to the transition that is its wake:

| App | events | the receiver dispatches |
|-----|--------|-------------------------|
| the reviewer's | `pull_request` | a review of the PR opened or updated |
| the implementer's | `pull_request_review` | the fix, the re-review request, or — approved — `task finish` by the task's owner |
| the coordinator's | what the seats do not take: `issues` for `cc:needs-decision` gates, nothing by default | the off-table events, and the milestone verbs |

`gh codecrew identity webhook <slug>` works an *active* hook under the
App's own key (credentials as `identity token` resolves them): no flags or
`--show` prints the URL, content type, whether a secret is set, and the
subscribed events; `--url <receiver>` and `--secret <the receiver's>` set
them (nothing stored, nothing printed); `--rotate-secret` mints a new one
and prints it once. Two things GitHub keeps for the settings page, and the
verb says so rather than pretend: an App minted **without** a webhook has
no hook configuration at all — `GET` and `PATCH` both answer 404, and no
endpoint creates one — so the verb refuses `NO_WEBHOOK` naming the page
where Webhook → Active and the URL are set by hand, after which `--secret`
and `--show` work; and an existing App's **event subscriptions** have no
endpoint either — set by the manifest at creation (`--events`) or ticked on
that page.

A receiver verifies `X-Hub-Signature-256` (HMAC-SHA256 of the raw body
with the secret), answers 2xx, and reads `action` from the payload. Two
things the run taught: a `ping` cannot fire a receiver that requires
payload fields — the first real `pull_request` is the test, and a 401
there means the secret, a 422 the payload
([#164](https://github.com/radiusred/gh-codecrew/issues/164), finding 60);
and the receiver's own last-fired timestamp is the health check for the
seam, because a redundant wake path hides a dead one for hours (finding 59).

**Worked example — a Paperclip routine.** A routine with a webhook trigger
(`POST /api/routines/<id>/triggers`, kind `webhook`, `signingMode:
github_hmac`) mints its own secret and a public fire URL
(`/api/routine-triggers/public/<publicId>/fire`). Mint the seat's App with
`identity new reviewer --name myorg-reviewy --with-webhook --webhook-url
<fire URL> --webhook-secret <the trigger's secret>` — or, for an App whose
hook is already active, `identity webhook myorg-reviewy --url <fire URL>
--secret <the trigger's secret>` — install it on the account, and the routine's
`lastFiredAt` moves on the first PR (the creation `ping` arrives signed
with GitHub's generated secret and is rejected — expected). The routine's body should say where
the payload lives (finding 61); the interop doc (#54) carries the whole
onboarding checklist.

### Coexistence with Copilot code review

Some teams have Copilot's PR auto-review enabled. Keep it if you like it —
but it is an **advisory signal alongside** the reviewer role, never the
role itself: it cannot be role-routed, its comments are not crew-attributed,
its reviews do not count toward required approvals (GitHub documents this),
and it needs a paid Copilot plan — platform-tier intentionality applies.
Nothing about it satisfies SPEC §8's independent-review layer; the
dispatched reviewer seat does that, whatever else also comments on the PR.

### Known quirks

- **Bot identities are not assignable to issues.** `task start` handles this:
  it records a `**Started by**` comment instead. Expected, not an error.
- **The viewer login carries a `[bot]` suffix** (`myorg-coder[bot]`) while
  the routing table names the bare slug; the CLI normalises this everywhere
  it resolves roles.
- **Approvals from Apps and required-review rules:** count only from
  principals with **write access** — the full rule, the trade it implies
  and the `--with-approval-permission` path are in step 2 of the manual
  ritual above. Verified both ways: an App with `Contents: write`
  satisfies a required-review rule by its own approval (proven live in the
  org); one with the read-only set does not (proven live on numberguess) —
  the superseding Decision on
  [#73](https://github.com/radiusred/gh-codecrew/issues/73). Where the
  approval is uncounted, `task finish` refuses with `REVIEW_NOT_COUNTED`;
  the paths are a non-author human approving on the reviewer's
  recommendation, granting the reviewer App write access, or
  `task finish --bypass` where the ruleset lists the operator as a bypass
  actor (the bypass is recorded on the PR).

## Teardown

What CodeCrew created, and what to do with it when a project ends or you
stop using the framework:

- **Record — keep.** Milestone and task issues, their comments, the PRs and
  the merged milestone documents under `docs/milestones/` are the audit
  trail; nothing needs deleting for the framework to be gone.
- **The pointer, contracts and extensions** — `.codecrew.yml`, `roles/`
  (the contracts and, from 1.1, the scaffolded `roles/<role>.local.md`
  extensions beside them — blank unless the project wrote into them),
  `AGENTS.md`, `CLAUDE.md` (hub only; it imports `AGENTS.md` for Claude
  Code — keep it if you had one of your own), `ROADMAP.md` in each repo.
  Delete or keep; they are plain files with no hooks.
- **Labels** — `cc:milestone`, `cc:task`, `cc:needs-decision` on each repo,
  created on first use. Remove in the repo's label settings if you like.
- **Task branches** — `task finish` deletes a merged head and `milestone
  close` sweeps; anything left is listed by `git branch -r`.
- **App identities** — each crew App under the owning account's Developer
  settings: uninstall it from the org or account, then delete the App. Its
  private key and credential stub live only in `~/.config/codecrew/`;
  delete the files. Tokens minted from them expire within the hour.
- **The extension** — `gh extension remove codecrew`.
