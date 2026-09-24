# Working in a CodeCrew project

This repository follows the CodeCrew protocol. If you are an agent dispatched
to work here:

1. **Find the hub.** Read `.codecrew/config.yml`; `hub: self` means this repo
   is the hub, otherwise follow `hub: owner/repo`.
2. **Check your CLI against the hub's protocol**, before the first verb:
   `gh codecrew version` prints the protocol the binary implements, which
   must be the same major as the hub pointer's `codecrew:` field and a
   minor at least the pointer's — the CLI checks the major alone. Here that
   is the local pointer; from a spoke, read the hub's with
   `gh api repos/<hub>/contents/.codecrew/config.yml -H "Accept: application/vnd.github.raw"`
   under whatever `gh` auth the session has, or straight after step 4 if
   that read is refused. If the binary falls short and you install the
   tools, `gh extension upgrade codecrew`; otherwise raise it with
   `gh codecrew checkpoint` and stop. Never upgrade mid-task.
3. **Load your role contract** from the hub's `.codecrew/roles/` directory —
   you were dispatched as one of:
   [implementer](roles/implementer.md),
   [reviewer](roles/reviewer.md), [qa](roles/qa.md),
   [doc-synthesizer](roles/doc-synthesizer.md). If no role was
   named, you are the implementer. Then load the project's extensions to it,
   in order: the hub's `.codecrew/roles/<role>.local.md`, then the same file
   in your working repo if it is a spoke — `gh codecrew roles show <role>`
   prints the whole composition (SPEC §7).
4. **Resolve your identity** per the contract:
   `export GH_TOKEN=$(gh codecrew identity token <slug>)` — the verb reads
   orchestrator-injected env vars (`GITHUB_APP_ID`/`GITHUB_CLIENT_ID`,
   `GITHUB_PRIVATE_KEY`/`GITHUB_PEM`) first, then the role's key and stub in
   `~/.config/codecrew/`, and refuses with a code otherwise; the operator's
   own `gh` auth is the identity only for an unrouted role.
5. **The protocol is [SPEC.md](../SPEC.md)** in the hub. The short version:
   state lives in GitHub issues and PRs, not in files; plan in the task issue
   before the first commit; atomic conventional commits referencing the task
   (`(#123)`); record decisions and deviations as structured comments the
   moment they happen; never verify or approve your own work; stop for a
   human at ask-the-human points and whenever `cc:needs-decision` is raised.

Commit messages follow conventional commits and are linted in CI.
