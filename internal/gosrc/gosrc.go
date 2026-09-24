// Package gosrc reads values out of Go source without building it. Its one
// job is the contract history (scripts/contract-history and the guard
// that checks it): the .codecrew/AGENTS.md a release wrote is a string
// constant in that release's internal/cli/init.go, not a file at the tag,
// and evaluating the constant from the tagged source is how its text is
// recovered without checking out and building every release.
package gosrc

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

// StringConst evaluates the package-level string constant name declared in
// src. It understands exactly what the scaffolds are made of — string
// literals, other string constants declared in the same file, parentheses
// and + — and refuses anything else with an error, so a scaffold this
// cannot read fails loudly rather than yielding a wrong text.
func StringConst(src []byte, name string) (string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), "", src, parser.SkipObjectResolution)
	if err != nil {
		return "", err
	}
	consts := map[string]ast.Expr{}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			vs := spec.(*ast.ValueSpec)
			for i, n := range vs.Names {
				if i < len(vs.Values) {
					consts[n.Name] = vs.Values[i]
				}
			}
		}
	}
	e := evaluator{consts: consts, seen: map[string]bool{}}
	return e.constant(name)
}

type evaluator struct {
	consts map[string]ast.Expr
	seen   map[string]bool
}

func (e evaluator) constant(name string) (string, error) {
	expr, ok := e.consts[name]
	if !ok {
		return "", fmt.Errorf("no constant %s in this file", name)
	}
	if e.seen[name] {
		return "", fmt.Errorf("constant %s refers to itself", name)
	}
	e.seen[name] = true
	defer delete(e.seen, name)
	return e.eval(expr)
}

func (e evaluator) eval(expr ast.Expr) (string, error) {
	switch x := expr.(type) {
	case *ast.BasicLit:
		if x.Kind != token.STRING {
			return "", fmt.Errorf("a %s literal is not a string", x.Kind)
		}
		return strconv.Unquote(x.Value)
	case *ast.ParenExpr:
		return e.eval(x.X)
	case *ast.Ident:
		return e.constant(x.Name)
	case *ast.BinaryExpr:
		if x.Op != token.ADD {
			return "", fmt.Errorf("operator %s is not string concatenation", x.Op)
		}
		l, err := e.eval(x.X)
		if err != nil {
			return "", err
		}
		r, err := e.eval(x.Y)
		if err != nil {
			return "", err
		}
		return l + r, nil
	default:
		return "", fmt.Errorf("cannot evaluate a %T", expr)
	}
}
