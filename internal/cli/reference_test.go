package cli

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// CLI.md is the reference for the verbs, and two things in it are copies of
// facts that live elsewhere: the synopsis of each verb, which `--help`
// prints, and the refusal codes each verb can exit with, which SPEC §10
// catalogues. These two tests are what keeps the copies true (M16-R2).

// TestReferenceSynopsesMatchHelp reads both surfaces and requires the same
// verbs with the same synopsis lines: a verb added to one alone, a flag
// added to one alone, and a reworded synopsis all fail here.
func TestReferenceSynopsesMatchHelp(t *testing.T) {
	help := helpSynopses(t)
	ref := referenceSynopses(t)

	for _, verb := range sortedKeys(help) {
		lines, ok := ref[verb]
		if !ok {
			t.Errorf("--help lists %q and CLI.md has no section for it", verb)
			continue
		}
		if !equalLines(help[verb], lines) {
			t.Errorf("%q: synopsis differs\n--help:\n  %s\nCLI.md:\n  %s",
				verb, strings.Join(help[verb], "\n  "), strings.Join(lines, "\n  "))
		}
	}
	for _, verb := range sortedKeys(ref) {
		if _, ok := help[verb]; !ok {
			t.Errorf("CLI.md documents %q and --help does not list it", verb)
		}
	}
}

// TestReferenceRefusalsMatchSpecTable reads CLI.md's per-verb refusal lists
// and SPEC §10's "Raised by" column and requires them to agree in both
// directions: a code a verb's section names must have a §10 row naming that
// verb, and every verb a §10 row names must name that code. A row raised by
// "any verb" is carried by the reference's Common refusals section, which
// stands for all of them.
func TestReferenceRefusalsMatchSpecTable(t *testing.T) {
	perVerb, common := referenceRefusals(t)
	table := specRaisedBy(t)

	for _, code := range sortedKeys(table) {
		for _, by := range table[code] {
			if strings.HasPrefix(by, "any verb") {
				if !common[code] {
					t.Errorf("SPEC §10 says %s is raised by %q and CLI.md's Common refusals do not list it", code, by)
				}
				continue
			}
			if _, ok := perVerb[by]; !ok {
				t.Errorf("SPEC §10 says %s is raised by %q and CLI.md has no section for that verb", code, by)
				continue
			}
			if !perVerb[by][code] {
				t.Errorf("SPEC §10 says %s is raised by %q and CLI.md's %q section does not list it", code, by, by)
			}
		}
	}

	for _, verb := range sortedKeys(perVerb) {
		for _, code := range sortedKeys(perVerb[verb]) {
			if !raisedBy(table[code], verb) {
				t.Errorf("CLI.md's %q section lists %s, which SPEC §10 does not name it for", verb, code)
			}
		}
	}
	for _, code := range sortedKeys(common) {
		if !raisedByAnyVerb(table[code]) {
			t.Errorf("CLI.md's Common refusals list %s, which SPEC §10 does not raise from any verb", code)
		}
	}
}

func raisedBy(list []string, verb string) bool {
	for _, by := range list {
		if by == verb {
			return true
		}
	}
	return false
}

func raisedByAnyVerb(list []string) bool {
	for _, by := range list {
		if strings.HasPrefix(by, "any verb") {
			return true
		}
	}
	return false
}

// helpSynopses reads the usage screen as two columns: the synopsis on the
// left, the description on the right. A verb's entry starts at indent 2 and
// its further synopsis lines are indented further; a line indented to the
// description column continues the description and carries no synopsis.
func helpSynopses(t *testing.T) map[string][]string {
	t.Helper()
	lines := strings.Split(usage, "\n")
	start := -1
	for i, l := range lines {
		if strings.TrimSpace(l) == "verbs:" {
			start = i + 1
			break
		}
	}
	if start < 0 {
		t.Fatal("the usage screen has no \"verbs:\" list; the parser is broken, not the reference")
	}
	out := map[string][]string{}
	descCol, verb := -1, ""
	for _, l := range lines[start:] {
		if strings.TrimSpace(l) == "" {
			break
		}
		indent := len(l) - len(strings.TrimLeft(l, " "))
		if descCol > 0 && indent >= descCol {
			continue // a wrapped description, not a synopsis
		}
		left, col := splitUsageColumns(l)
		if col > 0 && descCol < 0 {
			descCol = col
		}
		if indent <= 2 {
			verb = verbKey(left)
			if verb == "" {
				t.Fatalf("no verb name at the head of usage line %q", l)
			}
			out[verb] = []string{left}
			continue
		}
		if verb == "" {
			t.Fatalf("usage line %q continues no verb", l)
		}
		out[verb] = append(out[verb], left)
	}
	if len(out) == 0 {
		t.Fatal("no verbs parsed out of the usage screen")
	}
	return out
}

var columnGap = regexp.MustCompile(`  +`)

// splitUsageColumns returns a usage line's synopsis column and the index the
// description column starts at (0 when the line has none).
func splitUsageColumns(line string) (string, int) {
	indent := len(line) - len(strings.TrimLeft(line, " "))
	trimmed := strings.TrimRight(line[indent:], " ")
	if gap := columnGap.FindStringIndex(trimmed); gap != nil {
		return trimmed[:gap[0]], indent + gap[1]
	}
	return trimmed, 0
}

// verbKey is the verb a synopsis line names: the leading bare words, before
// the first argument placeholder or flag. "milestone close <n>" is
// "milestone close"; "role <name> [--login]" is "role".
func verbKey(synopsis string) string {
	var words []string
	for _, w := range strings.Fields(synopsis) {
		if !isBareWord(w) {
			break
		}
		words = append(words, w)
	}
	return strings.Join(words, " ")
}

// isBareWord reports whether a synopsis word is part of the verb's name: a
// lowercase word, so a flag (`--milestone`) and a placeholder (`<ref>`) both
// end the name.
func isBareWord(w string) bool {
	if w == "" || w[0] < 'a' || w[0] > 'z' {
		return false
	}
	for _, r := range w {
		if (r < 'a' || r > 'z') && r != '-' {
			return false
		}
	}
	return true
}

const (
	referenceHeading = "## "
	commonHeading    = "Common refusals"
)

// referenceSynopses reads the first fenced block of every verb section of
// CLI.md. The block's first line carries the `gh codecrew ` prefix a reader
// types; the rest align under it, and both sides are compared trimmed.
func referenceSynopses(t *testing.T) map[string][]string {
	t.Helper()
	out := map[string][]string{}
	for heading, body := range referenceSections(t) {
		if heading == commonHeading {
			continue
		}
		block := firstFencedBlock(body)
		if len(block) == 0 {
			continue
		}
		block[0] = strings.TrimPrefix(block[0], "gh codecrew ")
		verb := verbKey(block[0])
		if verb == "" {
			continue
		}
		if _, dup := out[verb]; dup {
			t.Errorf("CLI.md documents %q twice", verb)
		}
		out[verb] = block
	}
	if len(out) == 0 {
		t.Fatal("no verb sections parsed out of CLI.md")
	}
	return out
}

var codeSpan = regexp.MustCompile("`([A-Z][A-Z_]+)`")

// referenceRefusals reads the codes each verb section's "**Refusals.**"
// paragraph names, and the codes the Common refusals table names.
func referenceRefusals(t *testing.T) (map[string]map[string]bool, map[string]bool) {
	t.Helper()
	perVerb := map[string]map[string]bool{}
	common := map[string]bool{}
	for heading, body := range referenceSections(t) {
		if heading == commonHeading {
			for _, l := range strings.Split(body, "\n") {
				if !strings.HasPrefix(l, "| [`") {
					continue
				}
				if m := codeSpan.FindStringSubmatch(l); m != nil {
					common[m[1]] = true
				}
			}
			continue
		}
		block := firstFencedBlock(body)
		if len(block) == 0 {
			continue
		}
		verb := verbKey(strings.TrimPrefix(block[0], "gh codecrew "))
		if verb == "" {
			continue
		}
		codes := map[string]bool{}
		for _, m := range codeSpan.FindAllStringSubmatch(refusalsParagraph(body), -1) {
			codes[m[1]] = true
		}
		perVerb[verb] = codes
	}
	if len(common) == 0 {
		t.Fatalf("CLI.md's %q section names no refusal codes", commonHeading)
	}
	return perVerb, common
}

// refusalsParagraph is the section's "**Refusals.**" paragraph — from that
// marker to the next blank line. Codes are named there and nowhere else in
// a verb's section.
func refusalsParagraph(body string) string {
	const marker = "**Refusals.**"
	i := strings.Index(body, marker)
	if i < 0 {
		return ""
	}
	rest := body[i:]
	if end := strings.Index(rest, "\n\n"); end >= 0 {
		return rest[:end]
	}
	return rest
}

// referenceSections splits CLI.md on its "## " headings, keyed by the
// heading text with any code span markers removed.
func referenceSections(t *testing.T) map[string]string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "CLI.md"))
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	heading, body := "", &strings.Builder{}
	flush := func() {
		if heading != "" {
			out[heading] = body.String()
		}
		body = &strings.Builder{}
	}
	for _, l := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(l, referenceHeading) {
			flush()
			heading = strings.ReplaceAll(strings.TrimPrefix(l, referenceHeading), "`", "")
			continue
		}
		body.WriteString(l)
		body.WriteString("\n")
	}
	flush()
	return out
}

// firstFencedBlock returns the lines of a section's first ``` block,
// trimmed. That block is the synopsis.
func firstFencedBlock(body string) []string {
	var out []string
	in := false
	for _, l := range strings.Split(body, "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), "```") {
			if in {
				return out
			}
			in = true
			continue
		}
		if in {
			out = append(out, strings.TrimSpace(l))
		}
	}
	return nil
}

var specRow = regexp.MustCompile("(?m)^\\| `([A-Z_]+)` \\| ([^|]*?) \\|")

// specRaisedBy reads SPEC §10's table as code -> the verbs its "Raised by"
// column names, code spans stripped.
func specRaisedBy(t *testing.T) map[string][]string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "SPEC.md"))
	if err != nil {
		t.Fatal(err)
	}
	out := map[string][]string{}
	for _, m := range specRow.FindAllStringSubmatch(string(data), -1) {
		var by []string
		for _, part := range strings.Split(m[2], ",") {
			by = append(by, strings.ReplaceAll(strings.TrimSpace(part), "`", ""))
		}
		out[m[1]] = by
	}
	if len(out) == 0 {
		t.Fatal("no rows parsed out of SPEC §10's table")
	}
	return out
}

func equalLines(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if strings.TrimSpace(a[i]) != strings.TrimSpace(b[i]) {
			return false
		}
	}
	return true
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
