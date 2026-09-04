package codecrew

import (
	"os"
	"strings"
	"testing"
)

// The README shows this hub's own routing table under "Four seats, always
// staffed". GitHub renders the README from the repo and cannot transclude
// another file, so the block is a copy — and this test is what keeps the
// copy honest: reroute a seat, change a harness or pin a model in
// .codecrew.yml, and the build fails until the README block says the same.
func TestREADMECarriesTheRoutingTable(t *testing.T) {
	cfg, err := os.ReadFile(".codecrew.yml")
	if err != nil {
		t.Fatal(err)
	}
	i := strings.Index(string(cfg), "\nroles:\n")
	if i < 0 {
		t.Fatal(".codecrew.yml has no roles: section")
	}
	table := strings.TrimRight(string(cfg)[i+1:], "\n") + "\n"

	readme, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(readme), "```yaml\n"+table+"```\n") {
		t.Errorf("the README's yaml block is not the roles: section of .codecrew.yml, verbatim.\n"+
			"Replace the block under \"Four seats, always staffed\" with:\n\n%s", table)
	}
}
