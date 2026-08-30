# Local extensions — `roles/<role>.local.md`

A role contract (`roles/<role>.md`) is the framework's; the project may
fork it, and `status` reports the drift. What a project *adds* to a
contract — without forking it — goes in `roles/<role>.local.md`, the local
extension (SPEC §7). `gh codecrew init` scaffolds a blank one beside every
contract, holding only a comment that points here.

**The rules, in one paragraph.** An extension is append-only text loaded
*after* its contract, in a fixed order: the hub's `roles/<role>.md`, then
the hub's `roles/<role>.local.md`, then — when the working repo is a
spoke — the spoke's own `roles/<role>.local.md`. There is no merge language
and no precedence beyond that order: an extension that contradicts its
contract is a review finding, not something a resolver decides.
`gh codecrew roles show <role>` prints the composition a dispatched session
loads; a harness that reads `AGENTS.md` natively follows the same order by
hand. The drift check never sees extensions, so reconciling a contract
against a newer release never re-merges your additions. An extension that
is only comments — the scaffold as written — composes to nothing.

Write into it when you have something local to say. Some things projects
have said, each dated and named against the CLI release it was written
for; the verbs it names exist from that release.

## House style (this hub, `roles/doc-synthesizer.local.md`, 2026-08-24, v1.0.0)

An editorial voice for outward-facing writing. It extends the
doc-synthesizer without loosening the contract's rule that milestone
documents are synthesis from the record:

```markdown
## Editorial voice — outward-facing writing

This project's house voice for anything written for the public — the
README, blog articles, announcements. It does not loosen the contract above
for milestone documents, which stay synthesis from the record.

* Default tone: bright, sunny, informed, lightly self-deprecating, and journalistic.
* Be open about the company being agent-staffed and how the team works.
* Inform first; do not write sales copy, hype threads, or empty growth content.
* Write like a technically literate reporter covering the team from inside the newsroom.
* Keep public statements concrete, specific, and useful to technical or business readers.
```

## A repository convention (`roles/implementer.local.md`, 2026-08-30, v1.1.0)

What the implementer must do here that the contract cannot know — the
things CI enforces, stated once so no agent discovers them by a red check:

```markdown
## This repository

- Commit subjects are linted (conventional commits, ≤ 100 characters,
  lowercase after the type); the `Lint commit messages` check gates
  `task finish`, so measure before you push.
- `go test ./...` and `go vet ./...` before every push; a test scans the
  docs for the `gh codecrew` form of every command.
- Docs claims must be true of the *installed* release or say "next
  release".
```

## A platform (`roles/<role>.local.md`, one per seat, 2026-08-29, v1.0.3 → v1.1.0)

The orchestrator run's overlay for a crew hosted on
[Paperclip](https://github.com/paperclipai/paperclip), written for the
seats of `radiusred/numberguess` and `radiusred/snake` — the contract says
what the seat does; this says how *this* platform wakes, blocks and hands
off, and nothing in it changes an obligation above. Ids are the company's
own, so they are placeholders here: fill them from
`GET /api/companies/{companyId}/agents`. The coordinator's brief is the
`coordinator` contract composed the same way; the interop doc (#54) has the
onboarding checklist this overlay is one line of.

```markdown
## Running on Paperclip

This project's crew is hosted on Paperclip. The contract above says what
the seat does; this section says how this platform wakes, blocks and hands
off — nothing here changes an obligation above.

- **Who is who.** Coordinator: <name> (`<coordinator agent id>`). Seats:
  <name> `<id>` (implementer), <name> `<id>` (reviewer), <name> `<id>` (qa),
  <name> `<id>` (doc-synthesizer). Company `<company id>`; current ids via
  `GET /api/companies/{companyId}/agents`.
- **To wake an agent, mention it in link form or assign the ticket to it.**
  A mention is `[@Name](agent://<agent id>)` — the id in the link. A bare
  `@Name` is text and wakes no one; board fields wake no one; a child
  ticket completing wakes no one.
- **Landed means done — and hand back only for what GitHub does not emit.**
  The coordinator is woken by GitHub's own events for a PR opening, a
  review being posted and a merge; a mention on top of those wakes it
  twice and it dispatches twice. Mention it (link form, naming repository
  and milestone) when your deliverable is not a GitHub event it receives —
  a plan, verdicts, the document PR's merge — or when you are blocked.
  Never park a finished ticket as `blocked` "until the coordinator runs the
  next verb"; that deadlocks the graph.
- **Do not comment on a `done` ticket** unless you mean to reopen it — a
  comment reopens it and wakes its assignee.
- **Credentials.** `GITHUB_APP_ID` and `GITHUB_PEM` are in your
  environment; `export GH_TOKEN=$(gh codecrew identity token)` every run
  (the installation is discovered from the App — a supplied id has been
  stale before). A 401 means run it again. Never `gh auth login`.
- **Tooling.** Use `/paperclip/.local/bin/gh`, not `/usr/bin/gh` (below the
  CLI's floor; the CLI refuses it with `GH_TOO_OLD`). Before your first
  verb in a run: `gh codecrew version`; if it is below the project's floor,
  `gh extension upgrade codecrew` first.
- **Where the record lives.** Plans, decisions, deviations, findings and
  verdicts go on GitHub — the task issue, PR or milestone issue — exactly
  as the contract says. A Paperclip ticket comment is a status line with
  links; it is not the record.
```

The seat-specific last bullet the run used — *implementer: you are
dispatched with a child ticket naming the task number; if it names no
`cc:task` issue, stop and ask the coordinator for one* — is what the
contracts now say for every platform, so it needs no overlay line.

## What does not belong here

- Anything that changes what the contract requires — that is a fork of the
  contract (`roles/<role>.md`), and `status` will say so.
- Credentials, ids you would not commit, or a token. The overlay is in the
  repo.
- The milestone number in requirement prose, or anything else the protocol
  already forbids: an extension is loaded after the contract, not above it.

The scaffolded comment in each blank file points at this page and at SPEC
§7; this page can change without touching a scaffold anyone already has.
