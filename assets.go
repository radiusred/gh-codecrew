// Package codecrew embeds the repo's own protocol assets so the shipped
// binary can scaffold new projects (`gh codecrew init`) with the role
// contracts at their installed release — no copies to drift.
package codecrew

import "embed"

// Roles holds the five role contracts under ".codecrew/roles/" — listed by
// name, not globbed, for two reasons: .codecrew/roles/<role>.local.md files
// are this project's own extensions (SPEC §7) and must never ship inside
// the binary, where init would scaffold them into other projects and the
// drift check would treat them as contracts; and a pattern naming a
// directory would drop the whole tree anyway, embed excluding dot-prefixed
// names in that form.
//
//go:embed .codecrew/roles/implementer.md .codecrew/roles/reviewer.md .codecrew/roles/qa.md .codecrew/roles/doc-synthesizer.md .codecrew/roles/coordinator.md
var Roles embed.FS
