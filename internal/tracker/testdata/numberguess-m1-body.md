## Goal
Cycle 1 of the engineering-loop demo: a small, reviewable, well-tested number-guessing game, delivered as pure game-logic first (engine), then tests, then docs. This milestone charters the requirements the engine PR (#2, approved) must satisfy and that QA verdicts at the end.

Requirements:
- **M1-R1** — `newGame(min=1, max=100, seed?)` returns a fresh state `{min,max,secret,attempts:0,status:'playing'}` with `secret` an integer in `[min,max]` inclusive; an integer `seed` makes `secret` deterministic (mulberry32), omitting it uses `Math.random`.
- **M1-R2** — `guess(state,n)` returns a NEW state (input never mutated) with `result` one of `too_low|too_high|correct` (frozen `RESULTS` map), `attempts` incremented; a correct guess sets `status:'won'`.
- **M1-R3** — Input validation: non-integer guess → `TypeError`; out-of-range guess → `RangeError`; guessing a finished game → `Error`; `newGame` rejects non-integer bounds/seed and `min >= max`.
- **M1-R4** — Pure module: no I/O, no UI, no runtime deps, no build step; Node ESM at `src/engine.js`, exported as the package root.
- **M1-R5** — Automated tests (**implementer**) cover R1–R3 including seed determinism, immutability, attempt counting, and all error paths. The tests are the implementer's deliverable, landed via an implementer PR and gated by the reviewer; the **qa seat verdicts** them and never writes code. _(Re-chartered 2026-08-28: qa never authors tests — Bossy.)_
- **M1-R6** — Milestone record synthesized to `docs/milestones/1-core-engine.md` (Wordy) from the issue/PR record.

Gates charter: implementer=radiusred-cody, reviewer=radiusred-checky, qa=radiusred-testy, doc-synthesizer=radiusred-wordy (per .codecrew.yml).

## Requirements
_One line per requirement, its ID in bold: M1-R1, M1-R2, … — only
bolded IDs count as requirements, so this placeholder does not._

## Gates
_What "done" means beyond CI: e2e suites, manual UAT, sign-offs._


