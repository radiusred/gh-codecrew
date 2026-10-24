// Package scripts holds no Go: this file tests the repo's shell scripts
// from `go test ./...`, the only test CI runs.
package scripts

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

// TestReachRendersSnapshot renders a saved snapshot with `reach --from`,
// which reads nothing from GitHub, and compares it with the Markdown the
// milestone comment carries (M23-R1, #424). The fixture covers what the
// live repos do not yet show: non-crew authors, a non-crew bot, a login
// that matches a crew App's slug without being its bot user, discussion
// authors, a spoke with Discussions off, and a release with no assets.
func TestReachRendersSnapshot(t *testing.T) {
	for _, tool := range []string{"bash", "jq"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s not installed", tool)
		}
	}
	got, err := exec.Command("bash", "reach", "--from", "testdata/reach.json").Output()
	if err != nil {
		t.Fatalf("reach --from: %v", err)
	}
	want, err := os.ReadFile("testdata/reach.md")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("reach --from testdata/reach.json differs from testdata/reach.md\n got:\n%s\nwant:\n%s", got, want)
	}
}

// TestReachOffline: --from must not need gh or the network, so the
// snapshot a record quotes can be re-rendered from its saved JSON anywhere.
func TestReachOffline(t *testing.T) {
	for _, tool := range []string{"bash", "jq"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s not installed", tool)
		}
	}
	jq, _ := exec.LookPath("jq")
	bin := t.TempDir()
	if err := os.Symlink(jq, bin+"/jq"); err != nil {
		t.Fatal(err)
	}
	bash, _ := exec.LookPath("bash")
	cmd := exec.Command(bash, "reach", "--from", "testdata/reach.json")
	cmd.Env = append(os.Environ(), "PATH="+bin)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("reach --from with only jq on PATH: %v\n%s", err, out)
	}
}
