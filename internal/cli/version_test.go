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
	if got := strings.TrimSpace(b.String()); got != "dev" {
		t.Errorf("source build version = %q, want dev", got)
	}
}
