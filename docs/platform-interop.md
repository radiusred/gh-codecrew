# Platform interop: hosting a crew on an orchestration platform

This is the last rung of the ladder the [quickstart](first-milestone.md)
teaches: one human and one agent, then roles split across sessions, then
crew members with their own [App identities](identities.md), and then this —
a full orchestration platform dispatching the whole crew, with a human only
at the gates. [SPEC §9](../SPEC.md) is the anchor; this page is the
operational form of it, and it was deliberately held back
([#54](https://github.com/radiusred/gh-codecrew/issues/54)) until the rungs
below it had been proven.

Everything below was learned by running it. An orchestrator run drove four
cycles of real milestones on two proving grounds: cycles 1–3 on
[radiusred/numberguess](https://github.com/radiusred/numberguess), logged as
findings 1–50 on
[#119](https://github.com/radiusred/gh-codecrew/issues/119), and cycle 4 on
[radiusred/snake](https://github.com/radiusred/snake) — a fresh repo driven
by a dedicated coordinator agent from the first event — as findings 51–68 on
[#164](https://github.com/radiusred/gh-codecrew/issues/164). Every claim on
this page cites one of those findings, a shipped contract, or the PR that
changed the CLI because of it. Where the record names something that is still not solved, the last
section names it rather than rounding it up.

**A note on releases.** This page is written against the 1.1 line. The
installed release is **v1.0.3**, and six things named below ship in
**v1.1.0**: the coordinator seat — `roles/coordinator.md`, the `coordinator`
row in the routing table, and `gh codecrew identity new coordinator`
([#169](https://github.com/radiusred/gh-codecrew/pull/169)); `gh codecrew identity token`
([#171](https://github.com/radiusred/gh-codecrew/pull/171)); the blank `roles/<role>.local.md` that
`gh codecrew init` scaffolds beside each contract
([#174](https://github.com/radiusred/gh-codecrew/pull/174)); `task finish` refusing `NOT_OWNER`
when a seat that did not start a task tries to finish it
([#176](https://github.com/radiusred/gh-codecrew/pull/176)); `gh codecrew identity webhook`,
the `--webhook-secret` flag and the narrowed default event set
([#181](https://github.com/radiusred/gh-codecrew/pull/181)); and `gh codecrew init` committing
the scaffold it wrote ([#184](https://github.com/radiusred/gh-codecrew/pull/184)). On v1.0.3
the coordination layer is briefed by hand, a token is minted by a script or
in-line, an App's hook is worked on its settings page, and the scaffold is
committed by whoever ran `init` — which is precisely what the run did, at
the cost its findings record. Everything else here — the routing table,
`gh codecrew roles show`, the workflow verbs, the wake rules and the rest of
the onboarding checklist — is true of v1.0.3 today.

The crew that wrote and reviewed this page is itself agent-staffed, on the
protocol it describes.

## The separation of concerns

A platform and CodeCrew are not competing for the same job, and the run's
expensive moments were all one trying to do the other's.

| the platform owns | CodeCrew owns |
|---|---|
| waking agents, and what a wake means | the record: milestones, tasks, plans, decisions, gates, verdicts |
| discussion between agents, boards, tickets, queues | the routing table — who holds which seat, under which identity |
| the run loop, concurrency, cost | the gates: what may merge, what may close, and what refuses |

The run's central finding is what happens when the line moves. In cycle 1 the
platform kept its own build → review → test → docs loop *around* the repo,
so numberguess PR #2 existed with no milestone, no task issue, no plan, no CI and no
tests, and its Decision sat in the PR body where the record's gatherer never
reads: *"CodeCrew's gates bind agents that use its verbs; an orchestrator
that keeps its own task system around the repo bypasses all of them without
noticing"*
([#119 finding 7](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5443566637)).
The gates are not ambient. They bind at the moment a verb is run
([finding 36](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5453727497)).

The rule that falls out of the other direction is the one the coordinator
contract now states outright: **dispatch on the platform, cite on GitHub.**
A dispatch is a platform object — a fresh ticket, assigned to the seat,
woken the platform's way. The record is GitHub. Cycle 4 stalled twice on the
confusion: a coordinator posted the platform's own mention syntax as a
comment on a GitHub issue and nothing woke
([#164 findings 64 and 65](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465353819)).
The closing read of the whole run puts it plainly — *"CodeCrew owns the
record next to the code, the platform owns dispatch and discussion, and every
stall in this run was one trying to own the other's half"*
([#119 entry 18](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462998256)).

## The coordinator is a seat, and its own agent

Something has to open milestones, charter tasks, dispatch seats and drive the
review loop. For three cycles nothing owned that job by contract, and the
findings pile up in one direction: the coordination layer had no artifact
addressed to it
([finding 1](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5443566637)),
kept its own loop outside the protocol
([finding 7](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5443566637)),
inherited job descriptions written from a platform template that named a
different protocol
([finding 15](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5445671135)),
had no identity of its own and read every seat's 401 as its own
([finding 16](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5445972756)),
left the review loop's second round to start itself — twice, for 36 minutes
and a night
([findings 20](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5446623977)
and [26](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5447587101)),
could not write the ROADMAP row the CLI told it to commit
([finding 29](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5451814672)),
and learned the gates well enough to route around ever invoking them
([finding 36](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5453727497)).

So the coordination layer became the fifth seat: `roles/coordinator.md`
([the contract](../roles/coordinator.md), added in
[#169](https://github.com/radiusred/gh-codecrew/pull/169)), a
`coordinator` row in the routing table, and an identity of its own.
Unrouted (`~`) it is the operator — a solo project has a coordinator too, and
it is you ([SPEC §7](../SPEC.md)).

**Its identity.** `gh codecrew identity new coordinator` mints an App with
contents: read, issues: write, pull requests: read and metadata — never
contents: write, never pull requests: write. That set is
[finding 16](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5445972756)'s,
and it is deliberately not enough to merge, review or push: everything the
seat does is an issue, a comment or a label, and a 403 on anything else is
the contract enforced rather than a misconfiguration. The consequence is real
and worth planning for — the coordinator cannot commit the ROADMAP row
`milestone new` writes locally, so the implementer carries it in the
milestone's first PR
([finding 29](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5451814672);
the contracts now say so, [#174](https://github.com/radiusred/gh-codecrew/pull/174)).

**Why it is its own agent, and not the platform's CEO.** On a platform with a
lead agent, the obvious move is to hand that agent the coordination job. The
run tried it, and the operator's decision on
[#54](https://github.com/radiusred/gh-codecrew/issues/54#issuecomment-5462401509)
went the other way: the platform's lead agent spent four to eight minutes per
wake walking its stock CEO checklist before it read the event, which under
single-flight became a seven-deep queue of five-minute no-ops. The reasoning
generalises past any one platform — *a platform instance is unlikely to drive
a CodeCrew project as its only job, even when its product is software; the
lead agent should keep correct lead-agent instructions, and a dedicated
coordinator agent should carry explicit, specific instructions and role
files.*

So: a distinct agent, whose instruction bundle is
`gh codecrew roles show coordinator` (the contract plus the project's
`roles/coordinator.local.md` overlay), whose run loop is a lean router —
read the wake, `gh codecrew version` / `gh codecrew status`, act by the
contract, close the execution event, exit — which holds the coordination App
and is the receiver's assignee. Cycle 4 ran exactly that shape and the
numbers moved: the coordinator's share of the bill fell from three quarters
to a quarter (see [What it costs](#what-it-costs)). Its remaining defects
were all in the brief, not the shape
([#164 findings 62, 63, 65](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465148858)),
and each is now a line in the contract.

One more obligation the contract carries because a platform is not the only
kind of coordinator: **the watch.** What a coordinator watches, in the record
itself, is each open PR's reviews, each task issue's plan and labels
(`cc:needs-decision` raised or removed, `**Gate resolved:**`), the milestone
issue's comments, the checks on each PR, and the default branch. On a
platform the events replace the polling — one signal per transition, never
both; a coordinator with no platform is the operator running
`gh codecrew status`, which is why `status` reports every transition the
event table names
([operator, #54](https://github.com/radiusred/gh-codecrew/issues/54#issuecomment-5462792878)).

## Mapping agents to roles

The routing table in the hub's `.codecrew.yml` is the map, and it is
advisory: CodeCrew does not dispatch anything, so the table is a contract for
whatever does ([SPEC §5](../SPEC.md)).

```yaml
roles:
  implementer:     { harness: claude-code, model: claude-fable-5, identity: my-org-coder }
  reviewer:        { harness: codex, identity: my-org-reviewer }
  qa:              { harness: codex, identity: my-org-qa }
  doc-synthesizer: { harness: claude-code, identity: my-org-docs }
  coordinator:     { identity: my-org-coordinator }
```

One platform agent per row. The platform's job is to make each of its agents
*be* that row: the App identity for credentials, the composed contract for
instructions, and nothing else briefing it.

**The bundle is one command.** `gh codecrew roles show <role>` prints the
composition a dispatched session should load — the hub's `roles/<role>.md`,
then the hub's `roles/<role>.local.md`, then the spoke's if the working repo
is one. An onboarding script installs that output as the agent's instruction
file. Never choose a seat's model, harness or identity by hand when the table
says otherwise, and never brief a seat past its contract; both are in the
coordinator's Never list.

**The platform's paragraph is an overlay, never an edit.** What *this*
platform's wake syntax, agent ids and tool paths are goes in
`roles/<role>.local.md` — append-only text loaded after the contract, which
`gh codecrew init` scaffolds blank beside every contract it writes
([#174](https://github.com/radiusred/gh-codecrew/pull/174)). The worked
Paperclip overlay, with the ids as placeholders, is in
[docs/extensions.md](extensions.md); this page's checklist is the rest of
what that one file is a line of. An extension that contradicts its contract
is a review finding, not something a resolver decides
([SPEC §7](../SPEC.md)).

**Bundles belong to the company; the project is a parameter of the wake.**
Cycle 4's seats ran the whole cycle on the *previous* project's composed
bundles and it worked — the overlay's ids are the company's, not the
repository's — which is an accident the run recorded rather than a design
([#164 finding 57](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218)).
The coordinator showed the other half: its brief named one repository, so a
hand-back that named none sent it to `status` in the wrong clone, where it
read "no open milestones" and correctly concluded nothing was due
([finding 62](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465148858)).
Hence the contract's rule — a wake names its repository, and a nameless one
means `gh codecrew status` in every project you coordinate before concluding
anything — and hence the checklist's first row: give the platform a project
object, and the question cannot arise.

## Credentials

Each seat acts as its own GitHub App, and the mint is one verb on every
harness:

```sh
export GH_TOKEN=$(gh codecrew identity token <slug>)
```

`gh codecrew identity token` reads the App's id and private key from the
environment under the names platforms actually bind, then from the local key
and credential stub, then refuses with a code
([#171](https://github.com/radiusred/gh-codecrew/pull/171);
[SPEC §5](../SPEC.md)):

| what the platform binds | names the verb accepts |
|---|---|
| the App's id | `GITHUB_APP_ID`, `GITHUB_CLIENT_ID` |
| the private key | `GITHUB_PRIVATE_KEY`, `GITHUB_PEM` — PEM text or a path to it |
| the installation | discovered from the App; `GITHUB_INSTALLATION_ID` is a hint at most |

Three namings for two secrets is not an accident of documentation — it is
what the run met. Cycle 1's agents were handed `GITHUB_PRIVATE_KEY` +
`GITHUB_APP_ID`, the contracts named `GITHUB_CLIENT_ID` /
`GITHUB_PRIVATE_KEY` / `GITHUB_INSTALLATION_ID`, and a third agent's
environment carried `GITHUB_APP_ID` + `GITHUB_PEM`
([findings 2](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5443566637)
and [12](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5444149238)).
A supplied installation id was stale, and the agent that ignored it and
discovered its own was right
([finding 35](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5453727497)).
From v1.1.0 the verb does all of that, so the platform binds whatever it
binds; before it, every seat wrote its own JWT helper — which is what the run
watched them do, three times, in three dialects.

**The 401 reflex: mint again.** This is the single most expensive habit the
run found, and it is a habit, not a knowledge gap. Three of three agents
escalated a 401 to a human rather than re-minting. One had a working PEM in
its environment the whole time; one operator question — *"are you unable to
mint tokens from the PEM file in your environment?"* — and 67 seconds later
it had built the JWT, discovered the installation, minted, and posted its
review as itself
([finding 12](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5444149238)).
The tokens last an hour; a 401 means the hour is up. Every seat contract in
this hub now opens with the mint and says a 401 means run it again.

Three more rules the run paid for:

- **Per session, never `gh auth login`, never a shared config.** A reviewer's
  approval never reached GitHub because the `gh` config it inherited was
  shared across agents and held another App's expired token
  ([finding 10](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5443700318)).
  Give each agent its own `GH_CONFIG_DIR`; that is the platform's half of it.
- **A tool nobody's instructions name does not get used.** A skill that mints
  installation tokens sat materialised in one seat's bundle for a whole cycle
  and appears in none of its runs — every run derived the JWT by hand instead,
  because the contract's credentials bullet had already told it *how*
  ([#164 finding 67](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465449989)).
  Put the command in the contract, not the tool in the workspace.
- **A sketch in prose is not a verb.** A coordinator brief that read the
  token response with the wrong field name made every mint raise, and the
  agent diagnosed it as "flaky" and proposed a retry loop
  ([#164 finding 56](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218)).
  That is the argument for `gh codecrew identity token` existing at all.

## Wake paths, and the one-wake-path rule

A wake is how a platform tells an agent there is work. Three kinds appear in
the run, and they are not interchangeable:

| wake | what it is | what the run found |
|---|---|---|
| **mention** | naming the agent in a comment, in whatever link form the platform requires | the reliable one, and the only one the coordinator learned by example ([#119 finding 39](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5455026917)) |
| **assignment** | assigning a ticket to the agent | equally a wake; the receiver prior art used it exclusively ([#54](https://github.com/radiusred/gh-codecrew/issues/54#issuecomment-5461604709)) |
| **webhook** | a GitHub event delivered to a receiver that spawns the agent's work item | the fact rather than a proxy for it; the case for it is [finding 38](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5454810691) |

Board fields wake no one, and a seat that parks a finished ticket "until the
coordinator runs the next verb" deadlocks the graph — that shape cost the run
its two longest stalls, once between implementer and reviewer and once
between doc-synthesizer and coordinator
([#119 findings 20](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5446623977)
and [40](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5455071029)).
A child ticket completing wakes no one — or rather,
it did once and then did not, indistinguishably, from the API
([finding 38](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5454810691)):
*"a wake that is not observable in the board's structured state is a wake the
operator cannot debug."* That finding is why the design points at GitHub's
events: the review landing on the PR is the fact; the ticket completing is a
proxy, and the proxy failed.

**The one-wake-path rule: one signal per transition, never two.** This is the
rule the run cost the most to learn, and it is counter-intuitive because
redundancy feels safe.

- A coordinator woken by both the webhook and the ticket graph runs
  concurrently with itself. One PR opening produced three coordinator runs in
  44 seconds, each dispatching a reviewer; every duplicate re-entered the
  webhook as a new event, so *the duplication compounds one transition later*
  ([finding 46](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462194623)).
- Single-flight fixes the race and creates a queue: a run that read the PR's
  state before it queued acts on a snapshot minutes old
  ([finding 48](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462312601)).
  The answer is to **re-read at the act, not at the wake** — the PR's reviews
  *and your own open dispatches* — and skip a dispatch whose trigger the
  record already reflects.
- A dispatch template that ends *"mark this ticket done and mention me"*
  re-creates the second path the overlay just removed. Cycle 4's coordinator
  did this to itself: two dispatches per transition, three of its fourteen
  runs bought nothing at all
  ([#164 findings 53 and 54](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218)).
- And redundancy hides outages. Cycle 4's webhook was dead for four hours and
  the milestone completed anyway on the "duplicate" mention path; nobody
  noticed until the log was read
  ([finding 59](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5464097267)).
  *A single wake path would have stalled visibly — the stall is the alarm.*

The rule has a symmetric failure the run also paid for, so it is stated per
seat and per deliverable rather than as a slogan. Told *"mention me only if
you are blocked"*, a qa seat posted its verdicts and closed its ticket —
verdicts are not a GitHub event, so nothing woke and the milestone stopped
([finding 63](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465324165)).
The table lives in [`roles/coordinator.md`](../roles/coordinator.md): PR
opened, a fix pushed and a review posted travel by GitHub's event with no
hand-back; a plan, a set of verdicts and a merged document PR are handed back
by the platform's own wake path, naming repository and milestone.

**Execution events are one-shot.** Where a receiver spawns a work item per
event, that item is handled by the table and closed *in the same run* —
including when the event is not in the table. Leaving one open cost four
recovery runs for a single non-event, because the platform kept asking for a
disposition
([finding 44](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5461988272)).
The next event arrives on its own.

**One receiver shape or another.** A receiver can be platform-native — a
routine or trigger the platform exposes, which is the Paperclip recipe below
— or a small stateless service that maps events to wakes. The previous
venture's `paperclip-webhook-relay` was the latter: PR opened → the reviewer
seat, review submitted → the PR author's seat, merge → author, each as a
ticket *assigned* to the seat, and it is prior art rather than the design
because its routing was hard-coded rather than read from `.codecrew.yml`
([#54](https://github.com/radiusred/gh-codecrew/issues/54#issuecomment-5461604709)).
Either way, one App hook per seat is the delivery mechanism — see
["The receiver side"](identities.md) in identities.md for the events per seat
and what a receiver must do with a delivery.

## The onboarding checklist

In order. Each row is a step the run performed by hand, and the citation is
where it cost something.

| # | step | why, from the record |
|---|---|---|
| 1 | **The Project first.** Declare the hub as the platform's own unit of scope — name, the GitHub remote, a workspace, lead = the coordinator — and create the routines and the kickoff item under it. | Every run of cycle 4 opened with *"No project or prior session workspace was available"*, so each wake re-cloned or reused whatever checkout it last had; the coordinator's was another project's. The Project makes hub ↔ project one to one and finding 62's "which repo?" cannot arise ([#164](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465157957)). |
| 2 | **App hooks, not repository hooks.** Point each seat's App webhook at the receiver that dispatches that seat: `gh codecrew identity new <role> --with-webhook --webhook-url <receiver> --webhook-secret <the receiver's>`, or `gh codecrew identity webhook <slug> --url … --secret …` for an App whose hook is already **active** — an App minted without a webhook has no hook configuration at all, GitHub's API cannot create one, and the verb refuses `NO_WEBHOOK` naming the settings page where Webhook → Active is ticked by hand first. This is opt-in by design — *a crew App acts, it never listens* until an operator says otherwise ([#54](https://github.com/radiusred/gh-codecrew/issues/54#issuecomment-5454775046)). | One App hook delivers for every repository its installation covers, so a platform needs no repository hooks at all; cycle 4's two hand-pasted repository hooks were a workaround for seats whose Apps had no receiver ([identities.md, "The receiver side"](identities.md); [#181](https://github.com/radiusred/gh-codecrew/pull/181)). |
| 3 | **Instructions layered, never replaced.** Keep or supply the platform's own managed files (its run-loop, its wake semantics), *append* `gh codecrew roles show <role>`, and add the platform overlay. | Seats created by a lead agent's hiring plan carried exactly one file each and were never told how the platform wakes, blocks or hands off; the free-text parking that deadlocked the graph is what an agent invents without it ([#119 finding 41](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5455092737); [#54](https://github.com/radiusred/gh-codecrew/issues/54#issuecomment-5455092966)). |
| 4 | **Reinstall bundles from `roles show` whenever a contract or overlay changes** — a scripted step, run for every seat, not a coordinator's chore. The sequence is: contracts on `main` → `gh codecrew roles show <role>` → bundles. | A fix committed to the overlay is inert until the seats carry it; fifteen minutes after the change landed the seats were still emitting the old hand-back, and the reinstall had to be done through the API by hand ([#119, change point 3](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462283121)). |
| 5 | **Single-flight: one run at a time per agent** (`maxConcurrentRuns: 1` on Paperclip), set on the coordinator and, in a one-task-at-a-time crew, on the seats — before the first event. | The platform default assumes independent tickets, and CodeCrew's transitions are not independent; a coordinator running three copies of itself dispatched three reviewers for one PR ([#119, change point 2](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462261574); [finding 46](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462194623)). Pair it with the act-time re-read, or the queue goes stale ([finding 48](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462312601)). |
| 6 | **Execution events are one-shot** — say so in the brief: handle by the table, close in the same run, including for events not in the table. | Four recovery runs for one non-event, and the receiver marked its own run failed ([#119 finding 44](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5461988272)). |
| 7 | **Provisioning routines is a user act, not a coordinator act.** Create the receiver's routines and triggers with a user-level key; the onboarding script or plugin runs as a user. | An agent creating a routine assigned to *another* agent is refused — *"Agents can only manage routines assigned to themselves"* — and there is no permission key to grant. The coordinator built both payloads and correctly stopped ([#164 finding 55](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218)). |
| 8 | **State the payload path in the receiver's template.** Say in the spawned work item where the event body can be read. | A spawned item shows the templated text with its variables resolved and nothing else; a seat's first wake spent its opening tool calls hunting for the payload before giving up and reading the PR from GitHub. *A seat should not have to discover the platform's data model on its first wake* ([#164 finding 61](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465148858)). |
| 9 | **Watch the receiver's last-fired timestamp** as the seam's health check, and put a watchdog on the ingress if it is a tunnel. | Four hours dead, hidden by a redundant path; and the same ingress dropped again the next milestone ([#164 finding 59](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5464097267), and its second occurrence [at the start of M2](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465148858)). |
| 10 | **Do not test with `ping`.** A receiver whose template variables are required answers a ping `422 Missing routine variables` — correct behaviour, shown by GitHub as a failed delivery. The first real `pull_request` is the test; a **401 means the secret, a 422 the payload**. | [#164 finding 60](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5464097267). Secrets can be checked without GitHub by signing a body with the stored secret. |
| 11 | **Hand over coordinators at a milestone boundary, with a checklist:** the receiver's routines reassigned, the overlays repointed at the new agent's id, every seat's bundle reinstalled, single-flight set on the new agent, and **every open ticket the old coordinator holds for the project reassigned or closed**. | A mid-milestone handover left both coordinators live and each dispatched a reviewer for the same PR within a minute. *Two coordinators for one loop is the cost of handing over mid-milestone; the ticket graph, not the brief, decides who wakes* ([#119 finding 50](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462878896); [change point 5](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462791038)). |

## What it costs

Both logs recorded per-cycle numbers. They are reproduced here exactly as
recorded, in two tables rather than one, because the two logs count human
touches differently: #119 counts *non-gate touches*, and #164 splits
*onboarding* touches (bringing the repo and the platform up) from *workflow*
touches (standing in for a wake that should have fired).

**Cycles 1–3, on `radiusred/numberguess`**
([#119 entry 18](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462998256)):

| cycle | milestone | PRs | wall-clock | runs | $ reported | runs / PR | non-gate touches | gates |
|---|---|---|---|---|---|---|---|---|
| 1 (08-27/28) | M1 #3 | 3 (+ scaffold) | ~19.5 h | 85 | 68.49 | 28 | 5 | 0 |
| 2 (08-28) | M2 #9 | 2 | 4 h 01 | 30 | 26.62 | 15 | 6 | 0 |
| 3 (08-29) | M4 #15 | 4 | 3 h 08 | 97 | 82.76 | 24 (13 after the fixes) | 0 | 1 |
| **run** | 3 closed | 9 | ~27 h | **212** | **177.87** | | **11** | **1** |

The coordinator's share is recorded for two of those three cycles: **54% in
cycle 2, 75% in cycle 3** — $62.35 of $82.76 on 51 runs of a large model
against $6.09 for twelve of a smaller one.

**Cycle 4, on `radiusred/snake`, with a dedicated coordinator agent**
(*Loopy* in the table below)
([#164 entry 2](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465449989)):

| | runs | $ | coordinator share | PRs | runs / PR | wall-clock | onboarding touches | workflow touches | gates |
|---|---|---|---|---|---|---|---|---|---|
| M1, coordinator shape | 71 | 36.90 | Loopy 35 runs, 52% | 3 + scaffold | 24 | 1 h 43 | 4 | **0** | 0 (2 platform confirmations) |
| M2, seat routines | 33 | 17.40 | Loopy 8 runs, 23% | 4 | **8** | 1 h 40 (incl. ~30 min of stalls) | 4 | 3 | 0 |
| cycle 3 (numberguess M4), for scale | 97 | 82.76 | 75% | 4 | 24 | 3 h 08 | — | 0 | 1 |

Three things to read off these, all stated in the logs rather than derived
here:

- **Coordination was more than half the bill, and most of it bought nothing.**
  In cycle 2 the coordination layer took 11 of 30 runs and $14.34 of $26.62
  for three dispatches and one close; seven of its eleven runs answered an
  operator or found the state unchanged
  ([#119 finding 42](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5455234167)).
- **Routing transitions straight to the seats is the lever.** Cycle 4's second
  milestone put PR events on per-seat receivers and kept the coordinator for
  the milestone verbs, gates and off-table events: *"the seat routines cut
  runs per PR by two thirds against both M1 and cycle 3, and the
  coordinator's share of the bill from three quarters to a quarter"*
  ([#164 entry 2](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465449989)).
  One PR ran nine runs, three review rounds, coordinator idle throughout.
- **Human touches went to zero once the wiring was right, and came back only
  for brief defects.** Cycle 3 closed a milestone with zero touches that were
  not the one gate; cycle 4's first milestone had zero workflow touches; its
  second had three, *"all wakes after coordinator-brief defects, not platform
  faults, and each became a change point"* (same entry).

The honest framing: this is more expensive than one competent session doing
the same milestone — an order of magnitude even at cycle 2's improved rate
([finding 42](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5455234167)).
What you buy is that every step is attributable, independently reviewed and
recorded, with the human only at the gates. The trend across four cycles is
the argument, not any single figure.

## The worked example: Paperclip

[Paperclip](https://github.com/paperclipai/paperclip) is the platform all
four cycles ran on, so it is the recipe with evidence behind it. Ids,
hostnames and secrets are the company's own: every one below is a
placeholder, and none of them belongs in a repository.

**1. The Project.** Create a Project for the hub — `repoUrl` pointing at the
GitHub remote, a workspace with a checkout, `leadAgentId` the coordinator —
and set `projectId` on the routines and the kickoff ticket. Everything else
inherits it
([#164](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465157957)).

**2. The agents.** One per routing-table row, plus whatever else the company
does. The coordinator is a *distinct* agent from the platform's lead: its
`AGENTS.md` is the output of `gh codecrew roles show coordinator`, its
`HEARTBEAT.md` a lean router loop, `maxConcurrentRuns: 1`
([#54](https://github.com/radiusred/gh-codecrew/issues/54#issuecomment-5462401509)).
Bundles are installed per file through the agents API
(`PUT /api/agents/<agent id>/instructions-bundle/file`), idempotently, by the
onboarding script — the agent ids it writes into the overlay come from
`GET /api/companies/<company id>/agents`.

**3. The identities.** `gh codecrew identity new <role>` per seat and
`gh codecrew identity new coordinator` for the fifth, installed on the
account that owns the repositories; bind `GITHUB_APP_ID` and `GITHUB_PEM`
per agent, and give each agent its own `GH_CONFIG_DIR`.

**4. The routines and their triggers.** A routine per wake, assigned to the
seat it wakes, with a webhook trigger in GitHub's own signing mode:

```sh
paperclipai routine trigger:create <routine id> \
  --payload-json '{"kind":"webhook","signingMode":"github_hmac"}'
```

That returns the public fire URL and a secret, shown once
(`trigger:rotate-secret` mints another) — the seam the run found for a
platform whose tasks have no webhook-invokable wake
([#54](https://github.com/radiusred/gh-codecrew/issues/54#issuecomment-5454545834)). Cycle 4 ran three: a review dispatch
on `pull_request` → the reviewer, an implementer dispatch on
`pull_request_review` → the implementer, and a router for what the seats do
not take. **Create them with a user-level key** — an agent cannot create a
routine assigned to another agent (checklist row 7).

The routine body is a template; its variables interpolate from the payload
(`{{action}}`, `{{pull_request.number}}`, `{{review.state}}`). Two things to
get right in it: write the variables in plain text — a markdown editor that
escapes the underscore in `{{pull_request.number}}` leaves a variable that
can never resolve
([#119 finding 44](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5461988272))
— and state where the full payload lives, since the spawned issue does not
carry it: it is `triggerPayload` on the routine *run*
(`GET /api/routines/<originId>/runs`, the run whose `linkedIssue.id` is the
issue)
([#164 finding 61](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5465148858)).

**5. The App hooks.** Point each seat's App at its routine's fire URL and set
the App's secret to the trigger's, so GitHub signs what the trigger verifies:

```sh
gh codecrew identity new reviewer --name my-org-reviewer \
  --with-webhook --webhook-url <fire URL> --webhook-secret <the trigger's secret>

# or, for an App whose hook is already active:
gh codecrew identity webhook my-org-reviewer --url <fire URL> --secret <the trigger's secret>
```

From v1.1.0, `gh codecrew identity new --with-webhook` subscribes
`pull_request` and `pull_request_review` by default — the transitions a
platform routes to seats; `--events` names others
([#181](https://github.com/radiusred/gh-codecrew/pull/181)). On v1.0.3 it
subscribes five, three of which wake a platform for nothing. The creation
ping is signed with GitHub's generated secret and rejected; that is expected,
and it is the only delivery that can precede the App's install. No repository
webhooks are needed (checklist row 2).

**6. Ingress and the allowlist.** The fire URL must be reachable from GitHub.
Cycle 3 and 4 used a Tailscale Funnel on the platform host, `serve reset`
first and **path-scoped** to the trigger endpoint — a stray serve handler
went public on the first try — plus the platform's own hostname allowlist
(`npx paperclipai allowed-hostname <host>`)
([#54](https://github.com/radiusred/gh-codecrew/issues/54#issuecomment-5461972581)).
It dropped twice, silently, in cycle 4; the watchdog in checklist row 9 is
not optional.

**7. The overlay.** Write `roles/<role>.local.md` per seat — the coordinator's
name and id, the seats' ids, the mention form, the credential bindings, the
`gh` path, and the one-wake-path line. The full text as the run used it is
in [docs/extensions.md](extensions.md).

**8. Kick off.** One line to the coordinator, naming the repository. From
there the loop is the protocol's: `gh codecrew milestone new --requirement`,
`gh codecrew task new`, plans, PRs, reviews, `gh codecrew task finish` by the
seat that started the task, `gh codecrew milestone evidence`, verdicts, the
milestone document, `gh codecrew milestone close`.

## What is not solved yet

The record names these and nothing has closed them; they are here so this
page is not read as a finished story.

- **No onboarding script or plugin ships.** Steps 1–7 above are still done by
  hand or by a script you write. The
  [#54 inventory](https://github.com/radiusred/gh-codecrew/issues/54#issuecomment-5461972581)
  says where each piece should live — the CLI owns the GitHub side, the
  contracts own the neutral obligations and the overlay slots, a per-platform
  script owns bundles, ids, routines, triggers and hostnames, and this page
  owns the ingress recipe and the templates. The middle piece is the one that
  does not exist. The `init --platform <name>` scaffold that was proposed for
  it was deliberately dropped: `gh codecrew init` writes a *blank* extension
  pointing at [docs/extensions.md](extensions.md) instead, so the platform
  recipe can change without touching anyone's scaffold
  ([#174](https://github.com/radiusred/gh-codecrew/pull/174)).
- **The coordinator contract has not been run as a composed contract on a
  platform.** Cycle 4 ran on a hand-written brief; `roles/coordinator.md` was
  written *from* that brief and its findings and shipped afterwards
  ([#169](https://github.com/radiusred/gh-codecrew/pull/169)). Every
  obligation in it traces to a finding, but the contract itself awaits its
  first cycle.
- **The scaffold still costs the operator one merge, where a ruleset makes
  it a PR.** This one has narrowed since the run. From v1.1.0
  `gh codecrew init` commits exactly the files it wrote — on the current
  branch, or on
  `codecrew-bootstrap` cut from the default branch when that branch requires
  pull requests — so the scaffold is the last commit before the protocol
  starts rather than a PR with no task behind it, and delete-on-merge sweeps
  the branch the run found stranded
  ([SPEC §6](../SPEC.md)'s `init` row;
  [#183](https://github.com/radiusred/gh-codecrew/issues/183),
  [PR #184](https://github.com/radiusred/gh-codecrew/pull/184), which closed
  the capture the run's findings opened). What remains: behind a ruleset the
  scaffold arrives as a pull request, and it is the one merge the operator
  does by hand — recorded, because no milestone exists yet to `checkpoint`
  on, as a `**Gate raised:**` / `**Gate resolved:**` pair on the scaffold PR
  itself ([#164 finding 52](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218);
  [`roles/coordinator.md`](../roles/coordinator.md)). On a fresh repository
  whose org also requires a check that cannot report yet, that hand merge is
  an administrator merge — three repositories have met it
  ([#164 findings 51 and 68](https://github.com/radiusred/gh-codecrew/issues/164#issuecomment-5463692218)) —
  because `init` writes no CI workflow: the protocol reads check results and
  never defines CI ([SPEC §5](../SPEC.md)).
- **Wake coalescing is the platform's half, and is missing.** Single-flight
  plus one wake per signal is a queue, and a queue of stale wakes costs about
  what concurrency did, minus the collisions
  ([#119 finding 48](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462312601)).
  The act-time re-read is the framework's mitigation, not a fix.
- **A platform quirk went unexplained.** One agent's API key answered
  *"Routine not found"* for a routine the observer key could read; recorded
  in the [#54 inventory](https://github.com/radiusred/gh-codecrew/issues/54#issuecomment-5461972581)
  and never diagnosed.
- **The platform's tickets are not committed beside the record.** Decided,
  twice, on cost: *"more noise than they're worth"*. The findings logs quote
  what the fold-back needed, and the ticket ids are cited throughout, but the
  transcripts themselves are not in any repository
  ([#164](https://github.com/radiusred/gh-codecrew/issues/164),
  restated at [#119 entry 18](https://github.com/radiusred/gh-codecrew/issues/119#issuecomment-5462998256)).
- **Cycle 1's coordinator share was not measured.** The share is on the record
  for cycles 2, 3 and 4 only; cycle 1's split between the coordination layer
  and the seats is not recoverable from the log.

## Where to go next

- [Identities](identities.md) — minting the Apps, dispatching a role session,
  and "The receiver side": the events per seat and what a receiver does.
- [Local extensions](extensions.md) — `roles/<role>.local.md`, with the
  Paperclip seat overlay as a worked example.
- [`roles/coordinator.md`](../roles/coordinator.md) — the seat's contract,
  which is also the coordinator agent's instruction bundle.
- [SPEC §9](../SPEC.md) — the environments the protocol supports, of which
  this is the largest.
- [#119](https://github.com/radiusred/gh-codecrew/issues/119) and
  [#164](https://github.com/radiusred/gh-codecrew/issues/164) — the findings
  logs themselves, entry by entry, with the numbers and the ticket ids.
