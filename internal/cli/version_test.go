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
	if protocolVersion != "2.1" {
		t.Errorf("protocolVersion = %q; this binary implements protocol 2.1 (SPEC §5, §10)", protocolVersion)
	}
}
