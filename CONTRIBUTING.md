# Contributing

CodeCrew is developed with CodeCrew: the protocol in [SPEC.md](SPEC.md) is
the contribution process.

- **Ideas and bugs** — open an issue. Leave it unlabelled: that is a backlog
  capture, not protocol traffic, until it is adopted into a milestone.
- **Work** — every change is a task under the open milestone
  (`gh codecrew status` names it; [ROADMAP.md](ROADMAP.md) lists the
  finished ones): a task issue with a plan written
  before the first commit, atomic conventional commits referencing the task
  (`(#123)`), a PR that closes it, and a review from the reviewer seat — an
  App identity here, whose review arrives by dispatch, not by request. The
  operator merges through `task finish`; nothing merges around it.
- **Code ships with tests** in the same PR (the #46 convention); documented
  commands are executed by the reviewer, verbatim.
- **Decisions and deviations** are recorded as comments when they happen
  (SPEC §4), so the milestone document can be synthesized from the record.

If you are an agent, start at [AGENTS.md](AGENTS.md).
