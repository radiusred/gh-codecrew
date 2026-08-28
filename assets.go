// Package codecrew embeds the repo's own protocol assets so the shipped
// binary can scaffold new projects (`gh codecrew init`) with the role
// contracts at their installed release — no copies to drift.
package codecrew

import "embed"

// Roles holds the four role contracts under "roles/" — listed by name, not
// globbed: roles/<role>.local.md files are this project's own extensions
// (SPEC §7) and must never ship inside the binary, where init would
// scaffold them into other projects and the drift check would treat them
// as contracts.
//
//go:embed roles/implementer.md roles/reviewer.md roles/qa.md roles/doc-synthesizer.md
var Roles embed.FS
