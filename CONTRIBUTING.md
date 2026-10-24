# Contributing

CodeCrew is developed with CodeCrew: the protocol in [SPEC.md](SPEC.md) is
the contribution process.

- **Ideas and bugs** — open an issue. Leave it unlabelled: that is a backlog
  capture, not protocol traffic, until it is adopted into a milestone — at
  which point the task that delivers it names it (`task new --adopts`), and
  `task finish` closes it with the merge.
- **Work** — every change is a task under the open milestone
  (`gh codecrew status` names it; [ROADMAP.md](ROADMAP.md) lists the
  finished ones): a task issue with a plan written
  before the first commit, atomic conventional commits referencing the task
  (`(#123)`), a PR that closes it, and a review from the reviewer seat — an
  App identity here, whose review arrives by dispatch, not by request. The
  operator merges through `task finish`; nothing merges around it.
- **Housekeeping** — a change where a tool states the target and the diff is
  the whole decision takes the light path instead (SPEC §4): a `chore:`
  commit whose body names the tool, a PR reviewed by the reviewer seat, and
  no task. Regenerating the release table `roles sync` reads — every
  release's contracts and `.codecrew/AGENTS.md` scaffold — is one: after a
  release is tagged, `scripts/contract-history` in a PR of its own — the
  release workflow will not build the next release until the table covers
  every tag. The milestone record is the one other kind: the
  doc-synthesizer's `docs:` PR, referencing the milestone issue, reviewed
  the same way and merged by its author.
- **Reach** — a milestone that measures who finds and uses CodeCrew takes a
  snapshot at its open and again at its close with `scripts/reach`, so both
  are read the same way. It prints the Markdown to post on the milestone
  issue; `--json` keeps the raw figures and `--from` renders them again.
  GitHub's traffic figures need push access, so it runs under the
  operator's own `gh` auth. It only reads; what only the operator knows
  (Search Console, which downloads were the crew's own) is left as a
  placeholder to fill in before posting.
- **Code ships with tests** in the same PR (the #46 convention); documented
  commands are executed by the reviewer, verbatim.
- **Decisions and deviations** are recorded as comments when they happen
  (SPEC §4), so the milestone document can be synthesized from the record.

If you are an agent, start at [.codecrew/AGENTS.md](.codecrew/AGENTS.md).
