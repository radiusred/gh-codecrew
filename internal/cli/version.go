package cli

import (
	"fmt"
	"io"
)

// version is stamped by the release build (scripts/build-extension, via
// -ldflags -X); source builds report "dev".
var version = "dev"

func versionCmd(w io.Writer) error {
	fmt.Fprintln(w, version)
	return nil
}
