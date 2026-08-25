# Identities: running solo, staffing a crew

Roles are contracts, not accounts. Nothing in CodeCrew requires a GitHub App
to exist before a role can be staffed — a role with no identity acts as the
human operator, and the whole protocol runs on `gh auth login` alone. App
identities are infrastructure you add when you want them (attribution) or
need them (enforced independent review). This document covers both ends:
the zero-setup solo path, and minting a proper App identity for a role.

The conceptual model — identity tiers, credential resolution order — is in
[SPEC §5](../SPEC.md); this is the operational companion.

## Running solo

Prerequisites: a GitHub account, `gh` authenticated (`gh auth login`), the
extension installed (`gh extension install radiusred/gh-codecrew`), and a
`.codecrew.yml` in your repo (`hub: self` for a single-repo project).

Everything works, because solo is a routing configuration, not a reduced
protocol: every role is always staffed, and in pure solo *you* hold each
unrouted one. You (or agents running under your auth) create milestones and
tasks, write plans, commit, open PRs, raise and resolve gates. You also
perform the qa contract yourself: post per-requirement verdicts on the
milestone issue in the standard `**M1-R1 — satisfied.**` form — the close
gate counts the qa role holder's verdicts, and unrouted, that is you.
Declare the routing table in your hub's `.codecrew.yml` at onboarding (all
four roles, `~` for the ones you embody) — `codecrew init` scaffolds exactly
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
where a logo lands.
The App is owned by the hub's account unless `--owner` says otherwise, and
the verb refuses an App named after a role — crew-member names only. The
alternatives weighed at the gate — a published manifest template with no
verb, docs-only tier-sharpening, fine-grained PATs from a second account —
are recorded there; a shared App published by radiusred stays off the table
by prior decision, because the framework must never be key custodian for an
adopter's org.

Webhooks stay off by default — a crew App acts, it never listens — with an
opt-in `--with-webhook --webhook-url <receiver>` that subscribes the App to
the protocol-traffic events for platform users, whose orchestrators watch
through webhooks on the identity Apps they already own (the watch seam,
[#54](https://github.com/radiusred/gh-codecrew/issues/54)). The webhook
secret is printed once at creation and never written to disk. The alternatives weighed — a published manifest template with no
verb, docs-only tier-sharpening, fine-grained PATs from a second account —
are recorded there; a shared App published by radiusred stays off the table
by prior decision, because the framework must never be key custodian for an
adopter's org.

**Who should climb this way, and who shouldn't:** guided App creation is
the natural progression for a solo operator whose project is growing in
sophistication — agents earn attributable seats one at a time — and close
to a requirement for orchestration platforms (a Paperclip-style platform
drives every seat through App identities it owns). It is usually the wrong
first move for a *human* team adopting CodeCrew for isolated milestones or
tasks: humans are already distinct, attributable principals whose approvals
GitHub accepts, so route those roles to usernames in `.codecrew.yml` and
mint Apps only for the seats agents actually fill. Agents reading this:
before recommending `identity new`, check whether the seat in question is
held by a human — if it is, routing beats minting.

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

Mint a short-lived installation token per invocation — either via the
env-var path from SPEC §5 (`GITHUB_CLIENT_ID` / `GITHUB_PRIVATE_KEY` /
`GITHUB_INSTALLATION_ID`, for orchestrators), or with the bootstrap script
in this repo: `scripts/codecrew-token <app-slug>` — the **full** App slug
(`myorg-coder`, not `coder`). Then:

- `GH_TOKEN=$tok gh codecrew …` / `GH_TOKEN=$tok gh …` for API actions.
- Push over HTTPS as `https://x-access-token:$TOKEN@github.com/owner/repo.git`.
- Commit author is the App's bot user:
  `<slug>[bot] <UID+<slug>[bot]@users.noreply.github.com>`, where UID comes
  from `gh api 'users/<slug>%5Bbot%5D' --jq .id`.

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

- **Fresh context.** The role session starts clean — no implementation
  conversation, no shared scratch state. Independence is contextual as well
  as credential: an approval means little if the approver watched the code
  being written. For a sole rider, the floor is the same model in a separate
  session holding nothing but the role contract, the task, and the diff;
  better, when available, is a different model or harness — `harness:` in
  the routing table is the declared intent, and de-correlated judgment is
  the point of splitting the seats.
- **QA dispatches start with reachable evidence.** Before dispatching the
  qa role, the coordination layer runs `codecrew milestone evidence <n>` —
  citations that 404 cost a requirement its verdict once and the check is
  deterministic, so it runs as code, not hope. The qa contract refuses past
  an unreachable record from its side too.
- **Credentials.** The env-var path from SPEC §5, or
  `scripts/codecrew-token <slug>`. Beside the private key, `identity new`
  writes a credential stub (`~/.config/codecrew/<slug>.json` — the App ID
  and client ID, nothing secret): with it, the helper mints the App's JWT
  and discovers the installation itself, touching no user credential. For
  an App minted before the stub existed, write it by hand —
  `{"slug":"<slug>","app_id":<numeric App ID>}` — the account-lookup
  fallbacks depend on what the operator's token happens to see, and
  personal-account installations have twice needed hand-fed IDs without
  it. The session exports `GH_TOKEN` and is the App.
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
> `<app-slug>` (mint a token; never print it — use it only as GH_TOKEN).
> Confirm your identity differs from the PR author. Review correctness,
> plan, premises, tests, and consequences. Submit an ordinary GitHub PR
> review as the App: approve only if sound, otherwise request changes with
> concrete findings. Report your verdict and evidence back. Do not merge or
> run task finish.

### Known quirks

- **Bot identities are not assignable to issues.** `task start` handles this:
  it records a `**Started by**` comment instead. Expected, not an error.
- **The viewer login carries a `[bot]` suffix** (`myorg-coder[bot]`) while
  the routing table names the bare slug; the CLI normalises this everywhere
  it resolves roles.
- **Approvals from Apps and required-review rules:** verified both ways —
  approvals count only from principals with **write access**. An App with
  `Contents: write` satisfies a required-review rule by its own approval
  (proven live in the org); one with the reviewer table's read-only set
  does not (proven live on numberguess) — the superseding Decision on
  [#73](https://github.com/radiusred/gh-codecrew/issues/73). Where the
  approval is uncounted, `task finish` refuses with `REVIEW_NOT_COUNTED`;
  the paths are a non-author human approving on the reviewer's
  recommendation, granting the reviewer App write access, or
  `task finish --bypass` where the ruleset lists the operator as a bypass
  actor (the bypass is recorded on the PR).
