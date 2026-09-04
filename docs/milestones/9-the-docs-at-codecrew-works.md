# M9: The docs at codecrew.works

Tracking issue: [#201](https://github.com/radiusred/gh-codecrew/issues/201) ·
Synthesized 2026-09-04 by the doc-synthesizer role
([radiusred-wordy](https://github.com/apps/radiusred-wordy)) from the
milestone's own trail: #201's five requirements — **M9-R1 and M9-R3 as amended
by the operator on 2026-09-04**, with both original clauses quoted in the
issue's Amendments section — its Gates section and its **two** QA rounds; the
three task issues in the spoke
[radiusred/codecrew-www](https://github.com/radiusred/codecrew-www) and their
merged PRs ([#6](https://github.com/radiusred/codecrew-www/issues/6) / PR
[#7](https://github.com/radiusred/codecrew-www/pull/7),
[#11](https://github.com/radiusred/codecrew-www/issues/11) / PR
[#12](https://github.com/radiusred/codecrew-www/pull/12),
[#15](https://github.com/radiusred/codecrew-www/issues/15) / PR
[#16](https://github.com/radiusred/codecrew-www/pull/16)); the **nine
Decision and one Deviation comments** across those three issues, and the
implementer's correction to the Deviation nine minutes after posting it; the
reviewer's five review submissions on the three PRs; the six backlog captures
the milestone filed ([codecrew-www#8](https://github.com/radiusred/codecrew-www/issues/8),
[#13](https://github.com/radiusred/codecrew-www/issues/13),
[#14](https://github.com/radiusred/codecrew-www/issues/14),
[#17](https://github.com/radiusred/codecrew-www/issues/17),
[gh-codecrew#202](https://github.com/radiusred/gh-codecrew/issues/202) and
[radiusred/www#61](https://github.com/radiusred/www/issues/61)); the deployed
section at <https://codecrew.works/docs/> and the deploy that first published
it strict,
[Actions run 33862409652](https://github.com/radiusred/codecrew-www/actions/runs/33862409652).
The M7 and M8 documents supplied the house form — M8 read from PR
[#204](https://github.com/radiusred/gh-codecrew/pull/204)'s branch, since it
had not yet reached `main` when this was written.

The gather from `gh codecrew milestone close 9` is again not part of the raw
material, for the reason each of the last four documents has recorded: the
live close refuses `DOC_MISSING` until this file is on `main`, and `--dry-run`
stops earlier still, at `OPEN_TASKS`, naming
[#205](https://github.com/radiusred/gh-codecrew/issues/205) — the task that
writes it. Five milestones now. Every record below was read from its issue
directly.

**This is M9's first PR in the hub**, as M8's document PR was M8's. Every task
in the milestone lived in the spoke, so the ROADMAP row went in as an operator
commit straight to `main`
([5fb0add](https://github.com/radiusred/gh-codecrew/commit/5fb0add)) — the
third occurrence of the shape capture
[#197](https://github.com/radiusred/gh-codecrew/issues/197) describes, and the
observations at the end return to it.

## Goal and outcome

M8 gave codecrew.works a product home page. The reference documentation was
still Markdown in this repository and nowhere else on the web: `www.radiusred.uk`
had stopped serving it the day before the milestone opened, when
[radiusred/www#59](https://github.com/radiusred/www/issues/59) flipped that
project's sync to `readme_only` and pointed the Projects entry's docs icon at
codecrew.works. So the hero's "Read the docs" button led to a GitHub blob and
the site had a documentation-shaped hole in it. M9's goal
([#201](https://github.com/radiusred/gh-codecrew/issues/201)) was to close the
loop: build the upstream docs into `codecrew.works/docs/` on every deploy, give
them a nav tab, and re-point the button.

**Why this was a milestone and not a task on M8.** M8's goal had already said
so: "The rest of the site (blog, later docs) keeps the standard Zensical page
layout"
([#196](https://github.com/radiusred/gh-codecrew/issues/196)) — the docs were
named as later work while M8 was being written. The point was settled again
mid-M8, in the operator's
[answers to the content review](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5518717573):
answer 5, "all the docs will be brought on to the zensical site.. this button
will link to them", which the M8 record carries as "the docs coming onto this
site as their own task; the button flips to the site path then, not now". M9
is that task, grown into a milestone because it turned out to need a sync
engine, a link rewriter, a generated nav block and a CI change rather than a
href edit.

Five requirements, M9-R1 through M9-R5, all present when the milestone opened
at 17:20:29Z on 2026-09-03; none added later, and two amended (below). Three
tasks delivered the work, all three in `radiusred/codecrew-www`, a spoke whose
`.codecrew.yml` names this repository as its hub; this document is the fourth
task and the only one in the hub. The site deploys on every push to the
spoke's `main`, so each merge below was live within minutes.

- **[#6](https://github.com/radiusred/codecrew-www/issues/6) / PR
  [#7](https://github.com/radiusred/codecrew-www/pull/7)** — the whole
  mechanism in one PR: `sync_docs.py`, the link rewriter, the generated Docs
  tab, the hero button and the CI checkout. Opened 17:21:11Z and merged
  17:48:26Z on 2026-09-03
  ([6a1bc42](https://github.com/radiusred/codecrew-www/commit/6a1bc42)), four
  commits, 46 tests. Six Decisions and the milestone's one Deviation.
- **[#11](https://github.com/radiusred/codecrew-www/issues/11) / PR
  [#12](https://github.com/radiusred/codecrew-www/pull/12)** — the milestone
  records excluded from the site, on the operator's amendment. Opened and
  merged inside eleven minutes the same evening
  ([63a4822](https://github.com/radiusred/codecrew-www/commit/63a4822)), two
  commits, 53 tests. Two Decisions.
- **[#15](https://github.com/radiusred/codecrew-www/issues/15) / PR
  [#16](https://github.com/radiusred/codecrew-www/pull/16)** — the deploy
  builds strict, and the strict test fails loudly instead of skipping.
  Adopted from capture #14 after QA found M9-R5 not satisfied; merged
  2026-09-04
  ([218c8a1](https://github.com/radiusred/codecrew-www/commit/218c8a1)),
  three commits, 62 tests. One Decision, and a premise correction.

Delivered in full, in two QA rounds. M9-R1 through M9-R4 were satisfied in the
[first round](https://github.com/radiusred/gh-codecrew/issues/201#issuecomment-5534185587)
against the deployed site; M9-R5 came back **not satisfied** in the same
comment and was satisfied in a
[superseding verdict](https://github.com/radiusred/gh-codecrew/issues/201#issuecomment-5539054178)
nine hours later, after the fix task the verdict produced. The section is live:
`/docs/` returns the introduction under Home · Docs · Blog, and
`/docs/milestones/` returns 404, which is the amendment working.

## Decisions

Nine Decision comments across the three task issues. The first six set the
mechanism; the amendment then reversed one of them.

### The introduction becomes the section index, and the README does not sync

`radiusred/www`'s sync gives every project a `README.md` → `index.md` landing
page. This one
[does not](https://github.com/radiusred/codecrew-www/issues/6#issuecomment-5529603996):
`docs/introduction.md` becomes `/docs/`, and the README is not copied at all.
The reason is M8 — the home page **is** the README's argument in page form,
so syncing the README as the docs index would have said the same thing twice
on one site, one click apart. The introduction is the page that calls itself
"the map", which makes it the honest landing.

**Trade-off:** the README's links now cross a boundary the sync has to model,
because `../README.md` from a synced page has to mean the home page.
**Rejected:** sync the README anyway and let `/docs/` open on a duplicate.

This is the only one of the four shaping calls named in #6's Ask-the-human
section that got a Decision comment; the observations at the end return to the
other three.

### README anchors cross through an explicit table, and an unknown one degrades

`docs/introduction.md` links `README.md#the-receipts`. The home page names that
section `#codecrew-works`, because it is not a copy of the README and its
heading ids are its own — and `zensical build --strict` validates anchors, so
this was a build failure rather than a cosmetic mismatch.
[The fix](https://github.com/radiusred/codecrew-www/issues/6#issuecomment-5529604198)
is `HOME_ANCHORS`, a hand-maintained map from the README's headings to the
page's. An anchor the table does not know is **dropped and reported** — the
link still lands on the home page — so a new upstream link degrades to a less
precise link instead of breaking the deploy.

**Trade-off:** a hand-maintained table that can rot silently as the home page
changes, and does so in a direction no build catches. That is exactly the gap
capture [#8](https://github.com/radiusred/codecrew-www/issues/8) was opened
for. **Rejected:** adding an alias id to the home page (M8's test asserts the
section's id, and an alias reads as a second name for one section), and
dropping every anchor unconditionally (loses the one link that has a real
counterpart).

### The guides are ordered by the reading order the introduction prescribes

`radiusred/www`'s sync sorts a project's pages by filename.
[This one does not](https://github.com/radiusred/codecrew-www/issues/6#issuecomment-5529604508):
`docs/introduction.md` has a "Read in this order" section, which is an
editorial judgement already made in the hub, and reproducing it costs one
tuple. M9-R3's coarser order still holds around it.

**Trade-off:** `PAGE_ORDER` names upstream filenames, so a rename upstream
silently drops a page to the end of the guides rather than breaking the build.
That failure mode was chosen deliberately: an unknown page is appended before
contributing and security, alphabetically, so a new upstream document appears
on the site with no change here.

### The generated nav block is committed, and the site-build test may skip

Two smaller decisions with the same reasoning behind them — the checked-in
config should describe the site that actually deploys.

[The Docs nav block is committed to `zensical.toml`](https://github.com/radiusred/codecrew-www/issues/6#issuecomment-5529604825),
between `BEGIN_DOCS_NAV` / `END_DOCS_NAV` markers, the way `main.py` already
maintains the blog's. It means the whole tab is reviewable in the PR rather
than only visible after a build; the cost is that running the sync locally
after an upstream docs change leaves `zensical.toml` dirty.

[The hero links `docs/index.md`](https://github.com/radiusred/codecrew-www/issues/6#issuecomment-5529604673)
— a real file reference Zensical resolves and validates, so a rename upstream
is caught at build time. The cost is that a `--strict` build needs the synced
tree present, so `tests/test_site_build.py` skipped with a reason naming
`SYNC_SOURCE_BASE` when no `../gh-codecrew` was on disk. **Rejected:** an
absolute `/docs/` href, which no build step validates.

That skip is the decision QA later found insufficient. It was not wrong about
the hero link — the href stayed — but the shape of the skip is what M9-R5's
not-satisfied verdict was about, and task #15 replaced it. A decision made in
the milestone's first task and revised in its last, on the record both times.

### The amendment: the milestone records are not marketing-site documentation

The milestone's largest decision, and the only one that changed a requirement.

M9-R1 and M9-R3 as written named the milestone records explicitly: R1 copied
"the upstream `docs/` tree" whole, and R3's nav order ran "the introduction
first, then the remaining guides, the spec, **the milestone records**, and
contributing and security last". Task #6 implemented exactly that, and went
one step further — the sixth Decision on #6,
[a generated `docs/milestones/index.md`](https://github.com/radiusred/codecrew-www/issues/6#issuecomment-5529604359),
because two upstream pages link the bare directory `milestones/` and with
nothing there the rewrite had to send them to a GitHub tree view from a page
that was otherwise entirely on-site. Its recorded trade-off was that one page
on the site would have no upstream source — "prose the hub never reviewed",
kept to three lines and a list of links for that reason.

Five hours after the section went live the operator
[amended both clauses](https://github.com/radiusred/gh-codecrew/issues/201)
on #201 — dated 2026-09-04, attributed, and quoting both originals: **codecrew.works is
the marketing site; the per-milestone records are an internal engineering
artefact and do not belong in a product's documentation section.**
`docs/milestones/` is excluded from the sync, and links into it resolve to
GitHub like any other unsynced path. M9-R2 was left unchanged and, the
amendment notes, now carries more weight — the two `milestones/` links in
`docs/introduction.md` become GitHub URLs by the rule that already handles
`AGENTS.md` and `CHANGELOG.md`.

What it removed, beyond the pages: the generated index and its writer, the
record ordering key, the nested nav section, and the milestones filter in the
nav's extra-page computation — all of which existed only to serve the records,
and all of which went with them in PR #12. The Docs tab went flat. The
generated-index Decision is the one decision in this milestone that was made,
shipped, and reversed inside it, under six hours apart.

The amendment also settles a question the trade-off had raised and could not
answer from inside the spoke: whether a marketing site should carry an
engineering audit trail. It should not — and the record it does not carry is
this one.

### The exclusion is a named mechanism, not a deleted copy step

Task #11 could have satisfied the instruction with a condition in the copy
loop. It
[did not](https://github.com/radiusred/codecrew-www/issues/11#issuecomment-5533278362):
an `EXCLUDE` tuple holding `"milestones/"`, checked in one place — `is_excluded`,
consulted by both the copy loop and `map_to_dest`, so exclusion and
link-rewriting cannot drift apart and anything excluded is automatically a
GitHub link. `radiusred/www`'s sync already has the same mechanism, so the two
spokes stay legible to each other.

**Trade-off, in the implementer's words:** "a tuple with one entry is more
structure than one boolean needs today. Taken because the marketing/engineering
split this encodes is a judgement that will be made again, and the comment on
the constant is where the next person will look for it."

### Directory links get `tree`, decided by asking the filesystem

Excluding the records turned `milestones/` into the docs index's most prominent
outbound link, and `github.com/…/blob/main/<dir>` is a 301 to `tree/main/<dir>`.
Shipping a redirect on the page's most-clicked external link is avoidable, so
[`github_url` learned the distinction](https://github.com/radiusred/codecrew-www/issues/11#issuecomment-5533278484).

The interesting half is how it decides. The tempting test — "no file extension,
therefore a directory" — is wrong, and the upstream proves it: `LICENSE` has no
extension and is a file. So the code asks the checkout, `(source_dir() /
repo_path).is_dir()`, which is exact and free because the checkout is already
on disk for the sync. **Trade-off:** `github_url` now depends on the source
tree and cannot be called meaningfully without one; it is only ever called
during a sync, which refuses to run without a checkout, so the coupling is real
but unreachable. A test pins both sides.

The reviewer accepted this and named its edges in the same breath: a local
checkout on a non-`main` ref can disagree with the hard-coded `BRANCH = "main"`,
and `is_dir()` follows symlinks where GitHub renders the symlink itself. Both
are captured in [#13](https://github.com/radiusred/codecrew-www/issues/13),
deliberately unfixed, with two candidate shapes written down — read the type
from the ref being linked (`git cat-file -t main:<path>`), or drop the
distinction and accept the redirect, which is what the code did before #12.

### The deploy builds strict, and a broken upstream blocks it

The milestone's last Decision, and the one with the sharpest cost.
[`site.yml`'s build becomes `zensical build --clean --strict`](https://github.com/radiusred/codecrew-www/issues/15#issuecomment-5538939035),
unconditionally — on pushes to `main`, on the daily 00:05 UTC schedule, and on
demand.

**Trade-off:** a broken upstream — a link in this repository's docs that the
sync cannot resolve — now fails the deploy outright, blocking a deploy of
unrelated blog changes and turning the daily scheduled build red on a morning
nobody pushed anything to the spoke. Accepted because *the daily build is
precisely where an upstream break arrives*: the docs are synced at build time,
so no PR in the spoke ever sees the change and CI's strict check on pull
requests never runs for it. "Publishing a docs site with broken links to keep a
blog post on schedule is the wrong way round for a site whose product is the
docs; a red Actions run is visible, and the fix is a rewrite rule here or a
link upstream, either a one-line change."

Three alternatives are recorded as rejected, and the reviewer verified that all
three remained rejected in the shipped diff:

- keeping the deploy non-strict and relying on the strict build in the test
  suite — the status quo QA found not satisfied;
- a strict build as a separate non-blocking step (`continue-on-error`) before
  the real one — "the silent skip with a different name: the deploy still
  publishes the broken artefact, and nothing in the record says it was broken";
- falling back to a non-strict build when the strict one fails, so the blog
  still ships — "it publishes exactly the artefact strict mode exists to catch,
  and rewards not fixing the break".

## Deviations

One Deviation comment in the milestone, and it is about the protocol failing
rather than the plan changing.

**[The reviewer's first pass could not post its review](https://github.com/radiusred/codecrew-www/issues/6#issuecomment-5529696277).**
`radiusred-checky` was dispatched on codex — its declared harness — at PR #7's
first push. The session ran with `--sandbox workspace-write`, whose network
access is off by default, so every `gh codecrew identity token` call failed
with `socket: operation not permitted` before returning a token. The App never
authenticated; it could not read the PR, the task issue or the milestone, and
could not submit the review. No token was exposed.

It reviewed anyway, from disk: it read the diff and both repositories, ran the
suite (42 passed), reached **request changes** on three findings, and wrote its
verdict to a file. The findings were relayed by the coordination layer and
recorded **verbatim** in the Deviation, which states plainly that "the protocol
expects that conversation on the PR. It did not happen there" and that the
re-dispatched review, not the comment, is the gate. The three findings — a
build without the upstream checkout that did not fail as the docstring and
README claimed, nav labels escaping only the double quote, and a link escaping
above the repository root being silently republished with different semantics
— were fixed in
[48b1728](https://github.com/radiusred/codecrew-www/pull/7/commits/48b1728)
with tests, including one that parses the generated `zensical.toml` back with
`tomllib` after feeding the sync a backslash, a control character and quotes.

The implementer's own reading of finding 1 is on the record rather than left
to inference: "a false claim in my own docstring and README, not a missing
edge case: I tested the suite's skip path and never ran a strict build without
an upstream."

**And then the Deviation was corrected.** Nine minutes after posting it, the
implementer
[withdrew its own closing explanation](https://github.com/radiusred/codecrew-www/issues/6#issuecomment-5529807664):
the Deviation had said the `network_access` gap "had not surfaced before", and
that was wrong — the flag was already in the operator's notes for this project,
just not in the note named for codex dispatch, which was the one that was read.
"The dispatch failure was avoidable, not novel." The same correction rides on
[gh-codecrew#202](https://github.com/radiusred/gh-codecrew/issues/202), the
capture that asks `docs/identities.md`'s dispatch section to say a dispatched
session must be able to reach `api.github.com` — a gap in this repository's
documentation, filed by a spoke, and still open.

## The gates

**#201's Gates section was left as the scaffold's placeholder** — the template
line "_What "done" means beyond CI: e2e suites, manual UAT, sign-offs._",
unedited. M8's record says the same of #196; M7's declared seven gate
conditions. No `**Gate raised:**` or `**Gate resolved:**` comment exists
anywhere in this milestone. The reviewer of PR #12 confirmed the
other half of the same mechanism from the outside: "No unresolved
`cc:needs-decision` label is present."

What actually gated the work were two things the protocol does have words for
and one it does not. Each task set itself the same self-imposed gate and met
it: `uv run pytest` green and `zensical build --strict` reporting "No issues
found", quoted with counts in every PR description (46, 52 → 53, 62). QA's
verdicts gated the close, and one of them refused. And the operator's
amendment arrived as a requirement change on the milestone issue rather than as
a gate on a task: it altered what "done" meant rather than blocking a decision
a seat could not make, and #201 records it as an amendment with no
`cc:needs-decision` label raised.

Two of #6's four shaping calls also functioned as pre-milestone gates without
being recorded as any: the Ask-the-human section says all four "were put to the
operator before the milestone was opened and are settled". That they were
settled is on the record; what was decided, for three of the four, is not.

## The review rounds

Three PRs, five review submissions by the reviewer role holder
(`radiusred-checky[bot]`), one change request — plus the unpostable pass above,
which the protocol did not see.

- **PR [#7](https://github.com/radiusred/codecrew-www/pull/7)** — two posted
  rounds. The
  [first requested changes](https://github.com/radiusred/codecrew-www/pull/7#pullrequestreview-5105031907)
  **against the record, not the code**: 48b1728 had fixed all three relayed
  defects, but the PR description still described the superseded behaviour
  ("a build without one still succeeds") and still cited 13 sync tests and 42
  total against a head with 17 and 46. The reviewer reproduced the
  contradiction in a clean temporary copy — sync exit 1, then a strict build
  exit 1 on the missing hero target — and called the description "the durable
  implementation narrative", which is the reason the finding was blocking. The
  [second round](https://github.com/radiusred/codecrew-www/pull/7#pullrequestreview-5105056003)
  verified the corrected description against head 48b1728 and approved. It also
  filed a residual-gap list that became capture #8.
- **PR [#12](https://github.com/radiusred/codecrew-www/pull/12)** — two rounds,
  neither a change request. The
  [first approved](https://github.com/radiusred/codecrew-www/pull/12#pullrequestreview-5107650141)
  with non-blocking notes, having proved the exclusion was load-bearing by
  emptying `EXCLUDE` in an isolated copy and getting exactly five failures. The
  implementer then pushed a commit acting on two of those notes — a boundary
  fixture proving `milestones-archive/` is not caught by the exclusion, and a
  fragment stripped before the directory check — which **dismissed the standing
  approval** under the repository's stale-review rule. The
  [second round](https://github.com/radiusred/codecrew-www/pull/12#pullrequestreview-5107669902)
  is the more interesting one: rather than re-reading the diff, the reviewer
  loosened `is_excluded` to a plain `startswith` in a disposable copy to check
  the new fixture was discriminating rather than cosmetic, watched it fail, and
  approved. Everything it did not fix went into capture #13.
- **PR [#16](https://github.com/radiusred/codecrew-www/pull/16)** —
  [approved first round](https://github.com/radiusred/codecrew-www/pull/16#pullrequestreview-5111721559),
  after checking the diff *before* the PR narrative, then the task's Plan and
  Decision, then the adopted capture and QA's verdict. It executed all four
  configurations the PR claimed (62 passed; two collection-error shapes at exit
  2; 32 passed and 30 skipped unset) and noted specifically that the premise
  correction was "recorded in the task Plan as required, not merely in the PR
  body".

The one round that requested changes named a defect that was then fixed, and
the defect was in the record rather than in the software: a description that
had stopped describing its own head. M8's record states the same of its own
single change request. In both, the correction was to the narrative and no code
changed.

## QA: two rounds, one refusal, and the task it produced

The qa role holder (`radiusred-testy[bot]`) verdicted against the **deployed
site**, not the source, in both rounds.

**[First round](https://github.com/radiusred/gh-codecrew/issues/201#issuecomment-5534185587)**,
01:04:09Z on 2026-09-04: M9-R1 through M9-R4 satisfied, M9-R5 **not
satisfied**. R1 and R3 were assessed against the amended clauses, with the
amendment record itself audited first and found "adequate: dated,
operator-attributed, reasoned, and quotes both originals" — QA checking the
paperwork of a requirement change before verdicting against it. The probes
worth reading: a fresh sync from the real hub copying exactly the ten intended
files and the strict build emitting exactly their ten pages; `/docs/milestones/`
and a record URL both fetched **live** and returning 404, so the exclusion was
proved as an absence from the internet rather than an absence from a nav block;
all ten live docs pages crawled with every on-site link and fragment resolved,
and every rewrite-generated GitHub target followed to a 200; and the hero's
`href="docs/"` read from the deployed home page rather than the Markdown, "so
the button and destination are demonstrably present in the same deployment".

The R5 refusal is precise, and it is about a gap between a requirement's words
and a workflow's:

> Both workflows do check out the hub, and sync precedes main/build; the
> mapping, rewrite, exact nav, button target, built-target, and exclusion tests
> are substantive… But the shipped deploy in `site.yml` runs `zensical build
> --clean` without `--strict`, while the only strict build is inside a
> module-wide `skipif`… Thus the actual deployment is not required to stay
> warning-free under `--strict`.

QA had already reproduced this independently and filed it as
[#14](https://github.com/radiusred/codecrew-www/issues/14) — half an hour
before the verdict, during an earlier verdict run that "died to an environment
failure before it could post", recorded so the observation would not be lost
with the session and explicitly labelled "a lead to verify, not a delivered
verdict". The verdict then pointed at it: "independently reproduced and
recorded at codecrew-www#14, the fix home."

**The fix task.** #14 was adopted as
[#15](https://github.com/radiusred/codecrew-www/issues/15) with a plan, one
Decision, and its own tests, and PR #16 closed both. Two earlier records
describe the same shape: M4 remedied M4-R2 and M4-R3 with
[#63](https://github.com/radiusred/gh-codecrew/issues/63) / PR #64 inside M4,
and M5 remedied M5-R1 with
[#104](https://github.com/radiusred/gh-codecrew/issues/104) / PR #105 inside
M5, both before superseding verdicts; the M1–M3 and M6–M8 records describe
none. What is not in those two: the fix here lived in a spoke, and the remedy
adopted a capture the qa seat had filed itself, before the verdict that cited
it. The verdict and the backlog item were written by the same seat, in that
order, and the implementer's task was to adopt one and satisfy the other.

**[Second round](https://github.com/radiusred/gh-codecrew/issues/201#issuecomment-5539054178)**,
10:19:23Z, on fresh merged `main` `218c8a1`: **M9-R5 satisfied**, superseding.
The suite passed 62 tests; the two misconfiguration shapes failed collection at
exit 2 as designed; unset with no sibling checkout passed 32 and skipped only
the 30 upstream-dependent tests. The proof that mattered is a mutation probe —
the workflow test failed when `--strict` was dropped from the deploy command —
and a real deployment: Actions run 33862409652 ran the strict command on
`218c8a1`, reported "No issues found", and deployed green. The verdict waited
for that run rather than asserting the workflow would work.

It also came with a hairbrush finding that did not change the verdict and was
recorded anyway: a directory merely *named* `gh-codecrew` — empty, or an
initialised repository with no `docs/` — does not trigger the collection error
#15 and PR #16 advertise; execution still fails at sync, one stage later, so
the suite is red rather than silently green. QA wrote it on the **already-closed**
task #15 and it was lifted into its own capture,
[#17](https://github.com/radiusred/codecrew-www/issues/17), the same day: "a
diagnostic-boundary gap, not a coverage gap."

## Requirement outcomes

The status column is the verdict as the qa role holder wrote it, verbatim and
unqualified — for M9-R5, the word from the
[superseding verdict](https://github.com/radiusred/gh-codecrew/issues/201#issuecomment-5539054178),
with the refused first verdict named beside it. Everything else the milestone
knows about a requirement is in the two columns around it.

| Requirement | Delivered by | QA status | Notes |
|-------------|--------------|-----------|-------|
| M9-R1 — the docs build into `codecrew.works/docs/`, regenerated every build, never committed (**amended** 2026-09-04: `docs/milestones/` excluded) | [#6](https://github.com/radiusred/codecrew-www/issues/6) / PR [#7](https://github.com/radiusred/codecrew-www/pull/7); amended by [#11](https://github.com/radiusred/codecrew-www/issues/11) / PR [#12](https://github.com/radiusred/codecrew-www/pull/12) | `satisfied` | Tasks closed; first QA round, assessed against the amended clause and a live crawl |
| M9-R2 — links survive the move: on-site between synced pages, README → the home page, everything unsynced → absolute GitHub URLs | [#6](https://github.com/radiusred/codecrew-www/issues/6) / PR [#7](https://github.com/radiusred/codecrew-www/pull/7); [#11](https://github.com/radiusred/codecrew-www/issues/11) / PR [#12](https://github.com/radiusred/codecrew-www/pull/12) | `satisfied` | Tasks closed; first QA round. The verdict names future-shape gaps captured at [#8](https://github.com/radiusred/codecrew-www/issues/8) and [#13](https://github.com/radiusred/codecrew-www/issues/13), none present in the shipped corpus |
| M9-R3 — a Docs tab between Home and Blog, entries generated from what synced (**amended** 2026-09-04: no milestone records) | [#6](https://github.com/radiusred/codecrew-www/issues/6) / PR [#7](https://github.com/radiusred/codecrew-www/pull/7); amended by [#11](https://github.com/radiusred/codecrew-www/issues/11) / PR [#12](https://github.com/radiusred/codecrew-www/pull/12) | `satisfied` | Tasks closed; first QA round, assessed against the amended clause and the live tabs |
| M9-R4 — the hero's "Read the docs" button points on-site, in the same deploy that first publishes the docs | [#6](https://github.com/radiusred/codecrew-www/issues/6) / PR [#7](https://github.com/radiusred/codecrew-www/pull/7) | `satisfied` | Task closed; first QA round, read from the deployed home page rather than the Markdown |
| M9-R5 — CI checks out the hub and syncs before the build; the build stays warning-free under `--strict`; tests prove the mapping, rewrites, nav and button | [#6](https://github.com/radiusred/codecrew-www/issues/6) / PR [#7](https://github.com/radiusred/codecrew-www/pull/7); remedy [#15](https://github.com/radiusred/codecrew-www/issues/15) / PR [#16](https://github.com/radiusred/codecrew-www/pull/16) | `satisfied` | **`not satisfied` in the first round**, superseded in the second against `218c8a1` and [Actions run 33862409652](https://github.com/radiusred/codecrew-www/actions/runs/33862409652). The second verdict records a diagnostic-boundary finding, captured at [#17](https://github.com/radiusred/codecrew-www/issues/17) |

Two requirements were amended after the milestone opened, both by the operator
on 2026-09-04, both with their originals quoted on #201. None was added or
dropped. Captures filed by the milestone: six —
[codecrew-www#14](https://github.com/radiusred/codecrew-www/issues/14) adopted
as task #15 and closed inside the milestone;
[codecrew-www#8](https://github.com/radiusred/codecrew-www/issues/8),
[#13](https://github.com/radiusred/codecrew-www/issues/13),
[#17](https://github.com/radiusred/codecrew-www/issues/17),
[gh-codecrew#202](https://github.com/radiusred/gh-codecrew/issues/202) and
[radiusred/www#61](https://github.com/radiusred/www/issues/61) still open.

## Protocol-discipline observations

- **The ROADMAP row was an operator commit for the third time.** M9's tasks
  were all in the spoke, so there was no hub PR for the row to ride, and the
  operator pushed it to `main` directly
  ([5fb0add](https://github.com/radiusred/gh-codecrew/commit/5fb0add)) under
  the org ruleset's admin bypass, as for M8
  ([ae09d01](https://github.com/radiusred/gh-codecrew/commit/ae09d01)) and
  `radiusred/ops` (facc2ec).
  [#197](https://github.com/radiusred/gh-codecrew/issues/197) has carried four
  candidate shapes since 2026-09-02 and has not been adopted, so the shape it
  describes has now recurred once per spoke-only milestone since it was filed.
  One of its options — the doc-synthesizer owning both ends of the row — would
  have made this document's PR the only place the row was ever touched; as it
  is, this PR flips a row the operator wrote, for the second milestone running.
- **Three of the milestone's four shaping calls have no recorded rationale.**
  #6's Ask-the-human section names four decisions "put to the operator before
  the milestone was opened and are settled": a new milestone rather than a task
  on M8, `/docs/` with its own tab, syncing SPEC/CONTRIBUTING/SECURITY
  alongside `docs/`, and the introduction as the section index. Only the last
  produced a `**Decision:**` comment. The first can be reconstructed from #196
  and the operator's M8 answer 5, which is why this document could write it
  down; the middle two exist only as requirement text. *Inferred, not recorded:*
  that the tab and the three root files were chosen for the reasons M9-R1 and
  M9-R3 state, because no comment anywhere gives another. Settling a question
  before the milestone opens is not the same as recording the answer, and the
  Ask-the-human section is where that difference becomes invisible.
- **A verdict refused, and the milestone absorbed it without reopening
  anything.** M9-R5's `not satisfied` produced a new task in the spoke rather
  than a reopened #6, and QA re-verdicted nine hours later against a real
  deployment. The loop worked exactly as M4's and M5's did; what is new is that
  the evidence for the refusal was already a filed capture, written by the same
  seat during a session that crashed before it could post a verdict. The
  capture pre-dated the verdict by half an hour, carried the reproduction and
  two candidate shapes, and the verdict cited it as "the fix home"; the task
  that followed adopted it by number.
- **A decision was shipped and reversed inside the milestone, and both halves
  are on the record.** The generated `milestones/index.md` was a considered
  Decision with a named trade-off on 2026-09-03, and was deleted on 2026-09-04
  by an amendment that made its problem disappear. Nothing in the protocol
  handles supersession of a Decision by a requirement change; what made it
  legible here is that the amendment quotes the original clauses and PR #12's
  description enumerates what it removed and why. A reader gathering
  `**Decision:**` comments alone would find the index decision and no
  indication it no longer holds.
- **The one change request was against a claim, not code.** The PR
  description had drifted from its own head after the fix commit — behaviour
  inverted, test counts stale — and the reviewer reproduced the description's
  claim to prove it false before blocking on it, calling the description "the
  durable implementation narrative". No code changed for that finding. M8's
  record reports the same of its single change request.
- **A seat did real work the protocol could not see.** The codex-dispatched
  reviewer formed a sound verdict on three genuine defects and could not post
  any of it, because a sandbox denied network access before the App could mint.
  From the hub's side that is indistinguishable from a reviewer that was never
  dispatched. It was recovered by relay and recorded verbatim as a Deviation,
  and the implementer then corrected its own Deviation when its explanation
  turned out to be wrong — the gap was in a note nobody read, not a gap nobody
  had found. Both the failure and the correction are on the issue;
  [#202](https://github.com/radiusred/gh-codecrew/issues/202) asks this
  repository's own `docs/identities.md` to close the doc gap, and is still open.
- **The Gates section was the scaffold placeholder again.** #201 declared no
  gate beyond CI, and the operator's amendment — the largest intervention in
  the milestone — arrived as an edit to the milestone body rather than through
  any gate mechanism. Two human interventions shaped the milestone — the
  amendment and the four pre-milestone shaping calls — and neither has a
  `**Gate raised:**` / `**Gate resolved:**` pair behind it. M8's record reports
  the same placeholder on #196.
- **The README and the introduction needed no change for this close, and that
  is a finding.** The doc-synthesizer's refresh obligation exists so a
  milestone boundary cannot leave a stale claim standing. The check here: does
  either file send a reader to GitHub for a page that now has a rendered home
  at codecrew.works/docs/? Neither does — both link repository-relative paths,
  which are correct *in* the repository and are precisely what `sync_docs.py`
  rewrites on the way to the site. So the pages that describe the docs were
  made web-readable without a word of them changing, which is what M9-R2 asked
  the rewriter to do. Only the milestone count moved.
- **The docs site does not carry this document, by design.** M9 published ten
  of the hub's Markdown pages to codecrew.works and deliberately excluded one
  kind of page from them. The milestone whose subject was publishing the
  documentation is the one that decided the engineering record stays on
  GitHub; this record is the first written after that amendment.

## What the record does not contain

Gathered here rather than left implicit, because each one is a gap in the
trail rather than a gap in the work.

- **Three of the four shaping calls have no recorded rationale.** #6's
  Ask-the-human section states that all four "were put to the operator before
  the milestone was opened and are settled", and only the
  introduction-as-section-index call produced a `**Decision:**` comment. The
  new-milestone-rather-than-an-M8-task call is reconstructible from
  [#196](https://github.com/radiusred/gh-codecrew/issues/196)'s goal and the
  operator's
  [answer 5](https://github.com/radiusred/codecrew-www/issues/4#issuecomment-5518717573)
  on codecrew-www#4. The other two — `/docs/` with its own tab, and syncing
  `SPEC.md`, `CONTRIBUTING.md` and `SECURITY.md` alongside `docs/` — exist
  only as the text of M9-R1 and M9-R3. *Inferred, not recorded:* that they
  were chosen for the reasons those clauses state, because no comment
  anywhere gives another.
- **The adoption of capture #14 as task #15 has no recorded Decision.** The
  alternatives available after a `not satisfied` verdict — amend M9-R5,
  reopen #6, close the milestone with the requirement unmet — are not weighed
  anywhere. #15's Goal states the choice as made, and the coordination-layer
  act that made it left no comment. *Inferred:* the adoption followed the
  verdict's own pointer, "recorded at codecrew-www#14, the fix home".
- **No gate was declared or raised.** #201's Gates section is the scaffold's
  placeholder text, unedited, so the operator's amendment and the four
  pre-milestone calls have no `**Gate raised:**` / `**Gate resolved:**` pair
  behind them.
- **What prompted the amendment is not on the record.** #201's Amendments
  section gives the reasoning — a marketing site is not the place for an
  engineering artefact — but not the occasion: whether the operator reached it
  from seeing the records live at codecrew.works after PR #7 deployed, or had
  intended it before. The two events are five hours apart and nothing connects
  them in writing.
- **The relayed review exists in the record only as three quoted findings.**
  The unpostable reviewer wrote its verdict to a file; the Deviation on #6
  quotes the three blocking findings verbatim and summarises the rest. The
  file itself is not in either repository, so the fix-verification detail and
  the residual-gap list of that pass survive only through the second, posted
  review.
- **The amendment is dated 2026-09-04 while the task citing it opened at
  23:03Z on 2026-09-03.** #11 was created and merged before midnight UTC and
  states that the amendment happened; the amendment's own date is a day later.
  Nothing in the trail states which clock either is on, so the ordering below
  a few hours cannot be read off the record.
