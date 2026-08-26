# Security

- **Nothing secret lives in a CodeCrew repo.** App private keys and
  credential stubs live in `~/.config/codecrew/` on the operator's machine
  or in an orchestrator's secret store; the pointer file, contracts and
  scaffolds carry no credentials.
- **Tokens are short-lived.** Crew identities act with GitHub App
  installation tokens minted per invocation (one hour); nothing long-lived
  is written.
- **Least privilege by role.** Each App carries its contract's minimal
  permission set ([docs/identities.md](docs/identities.md)); a reviewer
  App's write access, when granted, is an explicit operator choice that is
  recorded.
- **Reporting.** Email security@radiusred.uk, or open a private
  vulnerability report on this repository. Please do not open a public
  issue for a credential or permission flaw.
