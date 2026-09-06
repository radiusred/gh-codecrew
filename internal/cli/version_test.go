package cli

import (
	"strings"
	"testing"
)

func TestVersionCmd(t *testing.T) {
	var b strings.Builder
	if err := versionCmd(&b); err != nil {
		t.Fatal(err)
	}
	want := "dev (protocol " + protocolVersion + ")"
	if got := strings.TrimSpace(b.String()); got != want {
		t.Errorf("source build version = %q, want %q", got, want)
	}
	if protocolVersion != "2.0" {
		t.Errorf("protocolVersion = %q; the .codecrew/ layout is protocol 2.0 (SPEC §5, §10)", protocolVersion)
	}
}
