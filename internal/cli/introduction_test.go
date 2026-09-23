package cli

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// TestIntroductionListsTheDryRunVerbs pins docs/introduction.md's list of
// the verbs that take `--dry-run` to the verbs whose flag set defines the
// flag. The list went stale once already: it named three verbs after
// `migrate` gained the flag as the fourth (#338), because nothing tied the
// prose to the code. The code side is read from this package's source —
// every function that builds a flag set named for its verb and defines a
// `dry-run` bool on it — so a fifth verb gaining the flag fails here until
// the introduction names it too.
func TestIntroductionListsTheDryRunVerbs(t *testing.T) {
	code := dryRunVerbs(t)
	if len(code) == 0 {
		t.Fatal("found no flag set defining dry-run; the source scan is broken")
	}
	doc := introductionDryRunList(t)

	if strings.Join(code, ", ") != strings.Join(doc, ", ") {
		t.Errorf("docs/introduction.md names the --dry-run verbs\n  %s\nand the flag sets that define it belong to\n  %s",
			strings.Join(doc, ", "), strings.Join(code, ", "))
	}
}

// dryRunVerbs returns, sorted, the name of every flag set in this package
// that defines a `dry-run` flag, by `Bool` or `BoolVar`. A flag set is
// named for its verb (`flag.NewFlagSet("task finish", …)`), and the flag
// is defined in the same function that builds the set, so the pairing is
// per function.
func dryRunVerbs(t *testing.T) []string {
	t.Helper()
	fset := token.NewFileSet()
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var verbs []string
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, name, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			var sets []string
			dryRun := false
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				switch sel.Sel.Name {
				case "NewFlagSet":
					if arg, ok := stringArg(call, 0); ok {
						sets = append(sets, arg)
					}
				case "Bool":
					if arg, ok := stringArg(call, 0); ok && arg == "dry-run" {
						dryRun = true
					}
				case "BoolVar":
					if arg, ok := stringArg(call, 1); ok && arg == "dry-run" {
						dryRun = true
					}
				}
				return true
			})
			if !dryRun {
				continue
			}
			if len(sets) != 1 {
				t.Fatalf("%s: %s defines dry-run beside %d flag sets; the scan pairs one set per function",
					name, fn.Name.Name, len(sets))
			}
			verbs = append(verbs, sets[0])
		}
	}
	sort.Strings(verbs)
	return verbs
}

// stringArg is call's i-th argument when it is a string literal: the flag
// name is the first argument of Bool and NewFlagSet, and the second of
// BoolVar, whose first is the variable it binds.
func stringArg(call *ast.CallExpr, i int) (string, bool) {
	if len(call.Args) <= i {
		return "", false
	}
	lit, ok := call.Args[i].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	s, err := strconv.Unquote(lit.Value)
	return s, err == nil
}

// introductionDryRunList returns, sorted, the verbs the introduction's
// sentence names as taking `--dry-run`: the code spans in the sentence that
// ends "take `--dry-run`…", read from the start of that sentence.
func introductionDryRunList(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "docs", "introduction.md"))
	if err != nil {
		t.Fatal(err)
	}
	flat := strings.Join(strings.Fields(string(data)), " ")
	const marker = " take `--dry-run`"
	at := strings.Index(flat, marker)
	if at < 0 {
		t.Fatalf("docs/introduction.md has no sentence naming the verbs that%s", marker)
	}
	sentence := flat[:at]
	if dot := strings.LastIndex(sentence, ". "); dot >= 0 {
		sentence = sentence[dot+2:]
	}
	var verbs []string
	for _, m := range regexp.MustCompile("`([^`]+)`").FindAllStringSubmatch(sentence, -1) {
		verbs = append(verbs, m[1])
	}
	sort.Strings(verbs)
	return verbs
}
