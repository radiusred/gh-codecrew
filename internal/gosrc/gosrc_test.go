package gosrc

import (
	"os"
	"strings"
	"testing"
)

const sample = `package x

const U = "https://example.com"

const (
	a, b = "one", "two"
)

const joined = ` + "`raw ` + U + \"/x\"" + ` + (a + b)

const notString = 3

const call = strings.Repeat("x", 2)

const loop = loop2
const loop2 = loop
`

func TestStringConst(t *testing.T) {
	got, err := StringConst([]byte(sample), "joined")
	if err != nil {
		t.Fatal(err)
	}
	if want := "raw https://example.com/xonetwo"; got != want {
		t.Errorf("joined = %q, want %q", got, want)
	}
	for _, name := range []string{"missing", "notString", "call", "loop"} {
		if _, err := StringConst([]byte(sample), name); err == nil {
			t.Errorf("%s: evaluated what it should refuse", name)
		}
	}
}

// The evaluator reads the scaffold the way the compiler does: the value it
// recovers from this tree's init.go is the constant the binary carries.
func TestStringConstReadsTheScaffold(t *testing.T) {
	src, err := os.ReadFile("../cli/init.go")
	if err != nil {
		t.Fatal(err)
	}
	got, err := StringConst(src, "agentsScaffold")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got, "# Agents\n") || !strings.Contains(got, "https://github.com/radiusred/gh-codecrew/blob/main/SPEC.md") {
		t.Errorf("agentsScaffold evaluated to:\n%s", got)
	}
}
