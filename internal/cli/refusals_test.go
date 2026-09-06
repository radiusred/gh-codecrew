package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// SPEC §10's table is the catalogue of refusal codes the protocol promises
// stable, and the source is what actually raises them. This test is the
// thing that keeps the two the same set: a verb that adds a code cannot
// merge without its row, a row cannot outlive the code it names, and the
// count SPEC states once is the count of the table it states it about
// (M13-R7).
func TestRefusalCodesMatchTheSpecTable(t *testing.T) {
	source := refusalCodesInSource(t)
	spec, err := os.ReadFile(filepath.Join("..", "..", "SPEC.md"))
	if err != nil {
		t.Fatal(err)
	}
	table := refusalCodesInSpec(string(spec))

	for _, code := range sorted(source) {
		if !table[code] {
			t.Errorf("%s is raised in internal/cli/ and has no row in SPEC §10's table", code)
		}
	}
	for _, code := range sorted(table) {
		if !source[code] {
			t.Errorf("SPEC §10's table has a row for %s, which nothing raises", code)
		}
	}

	// The count is stated once, in SPEC, and the sentence stating it opens
	// the table.
	want := fmt.Sprintf("**The refusal codes.** %s,", spellOut(len(source)))
	if !strings.Contains(string(spec), want) {
		t.Errorf("SPEC does not say %q; there are %d refusal codes", want, len(source))
	}
}

// refusalCodesInSource reads the catalogue of record — every `refuse("CODE"`
// in the package's own non-test sources, which is where every refusal is
// raised.
func refusalCodesInSource(t *testing.T) map[string]bool {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	call := regexp.MustCompile(`refuse\("([A-Z_]+)"`)
	codes := map[string]bool{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range call.FindAllStringSubmatch(string(data), -1) {
			codes[m[1]] = true
		}
	}
	if len(codes) == 0 {
		t.Fatal("no refuse(\"CODE\" calls found; the scan is broken, not the catalogue")
	}
	return codes
}

// refusalCodesInSpec reads the leading cell of every row of the §10 table —
// the only place in SPEC where a line begins with a code in a code span
// followed by a column.
func refusalCodesInSpec(spec string) map[string]bool {
	row := regexp.MustCompile("(?m)^\\| `([A-Z_]+)` \\|")
	codes := map[string]bool{}
	for _, m := range row.FindAllStringSubmatch(spec, -1) {
		codes[m[1]] = true
	}
	return codes
}

func sorted(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// spellOut writes a count the way SPEC's prose does. Two digits is all the
// catalogue will ever need; anything else is a bug in the test, not a
// number to render.
func spellOut(n int) string {
	ones := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	tens := []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
	switch {
	case n < 0 || n > 99:
		return fmt.Sprint(n)
	case n < 20:
		return strings.ToUpper(ones[n][:1]) + ones[n][1:]
	}
	word := tens[n/10]
	if n%10 != 0 {
		word += "-" + ones[n%10]
	}
	return strings.ToUpper(word[:1]) + word[1:]
}
