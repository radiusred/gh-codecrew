**Reach snapshot (GitHub side), 2026-10-24 14:00Z**, read by `scripts/reach`.

| Measure | example/hub | example/www | example/empty |
|---|---|---|---|
| Stars / forks / watchers | 7 / 1 / 2 | 0 / 0 / 0 | 0 / 0 / 0 |
| Views, 14 days (unique) | 10 (6) | 0 (0) | 0 (0) |
| Clones, 14 days (unique) | 383 (129) | 22 (9) | 0 (0) |
| Referrers | codecrew.works 4 (3), github.com 1 (1) | none | none |

- **Non-crew authors of issues and PRs, all time:**
  - example/hub: **2**: alice (2), dependabot[bot] (5). The crew authors are example-cody[bot] (209), op (176).
  - example/www: none. The crew authors are op (41).
  - example/empty: **1**: example-cody (1).
- **Non-crew authors of discussions, all time:**
  - example/hub: **1**: bob (1). The crew authors are example-checky[bot] (1).
  - example/www: Discussions are not enabled.
  - example/empty: none, and no discussions yet.
- **Release downloads, example/hub**, raw counts. `gh extension install` and `upgrade` each fetch one asset and the count is not unique, so subtract the crew's own known installs and upgrades by hand for an upper bound on outside installs:
  - v2.1.1, 3 in all: darwin-arm64 1, linux-amd64 2; 0 on darwin-amd64, windows-amd64
  - v2.0.0, 2 in all: linux-amd64 1, windows-amd64 1
  - v0.1.0: no assets
- **Clones** include the crew's own dispatch clones and CI, so they are not a measure of outside use.
- **Search Console:** not read by this script, because only the operator can see it. Operator: add the indexed and not-indexed counts, and impressions and clicks for the previous 28 days.
