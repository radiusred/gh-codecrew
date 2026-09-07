package cli

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// The venue is one interface in one package (M16-R1, #194): internal/gh is
// wrapped by internal/tracker and by nothing else, and no code outside that
// package runs gh, speaks to api.github.com, or builds a github.com URL.
// This is the test that keeps it that way — the leak it guards against is
// not a bug on the day it lands, it is the drift that put nine of these
// calls in internal/cli in the first place. It reads source, so a new leak
// is named by file and line and reads as "move it behind a venue method".

// venuePkgs are the packages allowed everything below: the venue itself,
// and the gh wrapper it owns.
var venuePkgs = []string{
	filepath.Join("internal", "tracker"),
	filepath.Join("internal", "gh"),
}

// venueExceptions are the string literals outside the venue that are
// allowed, each with the reason it is not a venue call. Every entry is a
// recorded Decision on the task, and the list is meant to stay this short.
var venueExceptions = map[string]string{
	// init.go's U is the documentation base for CodeCrew's own repository,
	// interpolated into the links the scaffold writes into AGENTS.md and
	// the role files. It addresses this project's docs for a human to
	// read; nothing fetches it, and no venue call is made through it.
	"https://github.com/radiusred/gh-codecrew/blob/main": "internal/cli/init.go: the documentation base URL init writes into the scaffold",
}

// venueLeaks reports every venue call in one parsed file, as the messages a
// failure prints. It is the whole rule set, in one place, so the guard and
// its own self-test cannot drift apart.
func venueLeaks(fset *token.FileSet, f *ast.File) []string {
	var found []string
	at := func(p token.Pos, format string, args ...any) {
		found = append(found, fset.Position(p).String()+": "+fmt.Sprintf(format, args...))
	}
	for _, imp := range f.Imports {
		if p, _ := strconv.Unquote(imp.Path.Value); strings.HasSuffix(p, "/internal/gh") {
			at(imp.Pos(), "imports internal/gh — the gh wrapper has exactly one importer, internal/tracker; reach it through a tracker.Tracker method")
		}
	}
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			if isExecCommandGH(node) {
				at(node.Pos(), "runs gh directly — every invocation of gh belongs behind a tracker.Tracker method")
			}
		case *ast.BasicLit:
			if node.Kind != token.STRING {
				return true
			}
			s, err := strconv.Unquote(node.Value)
			if err != nil {
				return true
			}
			if _, allowed := venueExceptions[s]; allowed {
				return true
			}
			switch {
			case strings.Contains(s, "api.github.com"):
				at(node.Pos(), "names api.github.com — the REST base is tracker.AppAPIBase and the calls through it are venue methods: %q", s)
			case strings.Contains(s, "https://github.com/"), strings.Contains(s, "http://github.com/"):
				at(node.Pos(), "builds a github.com URL — a venue knows its own URL grammar; use the tracker package's: %q", s)
			}
		}
		return true
	})
	return found
}

// isExecCommandGH matches exec.Command("gh", …) — the raw invocation that
// bypasses even internal/gh's own seam.
func isExecCommandGH(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Command" {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok || pkg.Name != "exec" || len(call.Args) == 0 {
		return false
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	name, err := strconv.Unquote(lit.Value)
	return err == nil && name == "gh"
}

func TestNoVenueCallsOutsideTheVenue(t *testing.T) {
	root := filepath.Join("..", "..")
	fset := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		if d.IsDir() {
			if name := d.Name(); name == ".git" || name == "testdata" {
				return fs.SkipDir
			}
			for _, p := range venuePkgs {
				if rel == p {
					return fs.SkipDir
				}
			}
			return nil
		}
		// Test files carry github.com links as fixtures — a record's
		// citations, a manifest's URL — and are not the code that runs.
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, parseErr := parser.ParseFile(fset, path, nil, 0)
		if parseErr != nil {
			return parseErr
		}
		for _, leak := range venueLeaks(fset, f) {
			t.Errorf("%s (M16-R1)", leak)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// The guard is only worth having if it fires: every rule is exercised here
// against source parsed in memory, through the same venueLeaks the walk
// uses, so a rule that stops matching fails rather than silently passing.
func TestVenueGuardCatchesEachLeak(t *testing.T) {
	for _, tc := range []struct {
		name, src string
		want      bool
	}{
		{"an internal/gh import", "package p\nimport \"github.com/radiusred/gh-codecrew/internal/gh\"\nvar _ = gh.Run\n", true},
		{"a raw gh exec", "package p\nimport \"os/exec\"\nvar _ = exec.Command(\"gh\", \"repo\", \"view\")\n", true},
		{"an api.github.com literal", "package p\nvar base = \"https://api.github.com\"\n", true},
		{"a github.com URL", "package p\nvar u = \"https://github.com/o/r/issues/1\"\n", true},
		{"an insecure github.com URL", "package p\nvar u = \"http://github.com/o/r\"\n", true},
		{"a github.com URL built by concatenation", "package p\nfunc f(o string) string { return \"https://github.com/organizations/\" + o }\n", true},
		{"the allowed documentation base", "package p\nconst U = \"https://github.com/radiusred/gh-codecrew/blob/main\"\n", false},
		{"a git exec", "package p\nimport \"os/exec\"\nvar _ = exec.Command(\"git\", \"status\")\n", false},
		{"the module's own import path", "package p\nimport \"github.com/radiusred/gh-codecrew/internal/tracker\"\nvar _ = tracker.GitHub{}\n", false},
		{"help-text prose naming the host", "package p\nvar help = \"a dead github.com link refuses\"\n", false},
	} {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "leak.go", tc.src, 0)
		if err != nil {
			t.Fatalf("%s: parsing the case: %v", tc.name, err)
		}
		leaks := venueLeaks(fset, f)
		if got := len(leaks) > 0; got != tc.want {
			t.Errorf("%s: guard fired = %v (%v), want %v", tc.name, got, leaks, tc.want)
		}
	}
}
