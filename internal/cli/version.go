package cli

import (
	"fmt"
	"io"
)

// version is stamped by the release build (scripts/build-extension, via
// -ldflags -X); source builds report "dev".
var version = "dev"

// protocolVersion is the SPEC version this binary implements — independent
// of the release tag (SPEC §5, §10). Pointers of another protocol major are
// refused; see config.Compatible.
const protocolVersion = "1.0"

func versionCmd(w io.Writer) error {
	fmt.Fprintf(w, "%s (protocol %s)\n", version, protocolVersion)
	return nil
}
