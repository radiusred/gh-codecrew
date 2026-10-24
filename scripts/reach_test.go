// Package scripts holds no Go: this file tests the repo's shell scripts
// from `go test ./...`, the only test CI runs.
package scripts

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func needTools(t *testing.T) {
	t.Helper()
	for _, tool := range []string{"bash", "jq"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s not installed", tool)
		}
	}
}

func reach(t *testing.T, args ...string) string {
	t.Helper()
	out, err := exec.Command("bash", append([]string{"reach"}, args...)...).Output()
	if err != nil {
		t.Fatalf("reach %s: %v", strings.Join(args, " "), err)
	}
	return string(out)
}

// TestReachRendersSnapshot renders a saved snapshot with `reach --from`,
// which reads nothing from GitHub, and compares it with the expected
// Markdown. The fixture covers what the live repos do not yet show:
// non-crew authors, a non-crew bot, a login that matches a crew App's slug
// without being its bot user, discussion authors, a spoke with Discussions
// off, an earlier release fetched on every platform, one with no downloads
// and one with no assets.
func TestReachRendersSnapshot(t *testing.T) {
	needTools(t)
	got := reach(t, "--from", "testdata/reach.json")
	want, err := os.ReadFile("testdata/reach.md")
	if err != nil {
		t.Fatal(err)
	}
	if got != string(want) {
		t.Errorf("reach --from testdata/reach.json differs from testdata/reach.md\n got:\n%s\nwant:\n%s", got, want)
	}
}

// operatorSays is the one sentence of the release bullet that only the
// operator can write; the script leaves placeholders in it.
var operatorSays = regexp.MustCompile(`The operator confirms that .*? outside installs\*\*\.`)

// TestReachMatchesOpeningSnapshot holds the script to the shape of M23's
// opening snapshot (testdata/opening-snapshot.md, comment 6384122881 on
// #423, taken by hand): fed the figures read at open
// (testdata/opening.json), the script must print the same table and the
// same four bullets, line for line. The only freedoms are the first line's
// prose after the timestamp, and the operator's own words where the script
// prints an _(operator: ...)_ placeholder (the crew-install count and
// Search Console). Anything the script adds goes after the snapshot's last
// line, in the appendix.
func TestReachMatchesOpeningSnapshot(t *testing.T) {
	needTools(t)
	raw, err := os.ReadFile("testdata/opening-snapshot.md")
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Split(strings.TrimRight(string(raw), "\n"), "\n")
	got := strings.Split(reach(t, "--label", "M23-R1 opening snapshot", "--from", "testdata/opening.json"), "\n")
	if len(got) < len(want) {
		t.Fatalf("reach printed %d lines, the opening snapshot has %d", len(got), len(want))
	}
	if got[len(want)] != "" || !strings.HasPrefix(got[len(want)+1], "<details>") {
		t.Errorf("after the snapshot's %d lines, want a blank line and the collapsed appendix, got %q, %q", len(want), got[len(want)], got[len(want)+1])
	}
	for i, w := range want {
		g := got[i]
		switch {
		case i == 0:
			const head = "**M23-R1 opening snapshot (GitHub side), 2026-10-24 "
			if !strings.HasPrefix(w, head) || !strings.HasPrefix(g, head) {
				t.Errorf("line 1: want both to open %q\n got: %s\nwant: %s", head, g, w)
			}
		case strings.HasPrefix(w, "- **Search Console"):
			const lead = "- **Search Console (operator's report, "
			if !strings.HasPrefix(g, lead) || !strings.Contains(g, "_(operator:") {
				t.Errorf("line %d: want %q with operator placeholders, got: %s", i+1, lead, g)
			}
		case strings.HasPrefix(w, "- **Release downloads:**"):
			if !strings.Contains(g, "_(operator: how many of these 3)_") {
				t.Errorf("line %d: no placeholder for the crew's own installs: %s", i+1, g)
			}
			g = operatorSays.ReplaceAllString(g, "<operator>")
			w = operatorSays.ReplaceAllString(w, "<operator>")
			fallthrough
		default:
			if g != w {
				t.Errorf("line %d differs from the opening snapshot\n got: %s\nwant: %s", i+1, g, w)
			}
		}
	}
}

// TestReachOffline: --from must not need gh or the network, so the
// snapshot a record quotes can be re-rendered from its saved JSON anywhere.
func TestReachOffline(t *testing.T) {
	needTools(t)
	jq, _ := exec.LookPath("jq")
	bin := t.TempDir()
	if err := os.Symlink(jq, bin+"/jq"); err != nil {
		t.Fatal(err)
	}
	bash, _ := exec.LookPath("bash")
	cmd := exec.Command(bash, "reach", "--from", "testdata/reach.json")
	cmd.Env = append(os.Environ(), "PATH="+bin)
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("reach --from with only jq on PATH: %v\n%s", err, out.String())
	}
}
