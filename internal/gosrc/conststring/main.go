// Command conststring prints the value of a package-level string constant
// read from the Go source on stdin — scripts/contract-history runs it over
// each release's internal/cli/init.go to recover the .codecrew/AGENTS.md
// that release wrote. Usage: conststring <name> < file.go
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/radiusred/gh-codecrew/internal/gosrc"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: conststring <name> < file.go")
		os.Exit(2)
	}
	src, err := io.ReadAll(os.Stdin)
	if err == nil {
		var text string
		if text, err = gosrc.StringConst(src, os.Args[1]); err == nil {
			_, err = io.WriteString(os.Stdout, text)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "conststring:", err)
		os.Exit(1)
	}
}
