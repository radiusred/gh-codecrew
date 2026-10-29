## Editorial voice — outward-facing writing

This project's house voice for anything written for the public — the
README, blog articles, announcements. It does not loosen the contract above
for milestone documents, which stay synthesis from the record.

* Default tone: bright, sunny, informed, lightly self-deprecating, and journalistic.
* Be open about the company being agent-staffed and how the team works.
* Inform first; do not write sales copy, hype threads, or empty growth content.
* Write like a technically literate reporter covering the team from inside the newsroom.
* Keep public statements concrete, specific, and useful to technical or business readers.

## Front-door documents — this hub

This hub's front door is the README and `docs/introduction.md`. At every
milestone boundary the record's PR brings their claims into line with what
the milestone delivered: the README's proof points, and the introduction's
release, verbs and refusal codes.

## Close announcements — this hub

A duty of this hub, not part of the protocol. A milestone about
gh-codecrew itself or how it is used (the CLI, the protocol, the role
contracts, a release, the docs a user follows) gets one post on Bluesky and
one on X when it closes. A milestone that only changes the codecrew.works
website gets none. The operator's Decisions on #423 set the
[accounts, voice and approval](https://github.com/radiusred/gh-codecrew/issues/423#issuecomment-6467856741)
and [which milestones](https://github.com/radiusred/gh-codecrew/issues/423#issuecomment-6469436459).

* **Where:** CodeCrew's own accounts, Bluesky `codecrew.works` and X
  `codecrewworks`.
* **How it is written:** radiusred/www's
  [`announcements/README.md`](https://github.com/radiusred/www/blob/main/announcements/README.md)
  is the guide. Re-read every draft against its Don't list before it goes
  out: `social lint` checks only the mechanical rules, and M24's close post
  passed it while breaking two of the voice rules
  ([#443](https://github.com/radiusred/gh-codecrew/issues/443#issuecomment-6469438706)).
* **Approval:** none per post once it passes lint, because a close post only
  restates what is already public on GitHub.
* **Record:** one comment with the post's links on the milestone issue.
* **Who drafts and sends:** the doc-synthesizer, after `milestone close`, as
  a task the coordinator dispatches. The record PR carries no draft.
