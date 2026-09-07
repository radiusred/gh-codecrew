package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// docs/working-offline.md tells an operator what the CLI prints when it
// cannot reach GitHub, and a page that quotes output is only worth reading
// while the quotation is true. These tests build the two texts it quotes
// from the code that prints them and fail when the page has stopped
// carrying them — the shape TestRefusalCodesMatchTheSpecTable uses for
// SPEC §10's catalogue. They assert nothing about GitHub and run offline
// themselves.
const offlineDoc = "working-offline.md"

func readOfflineDoc(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "docs", offlineDoc))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// quoted checks the page carries a line the CLI prints. The `gh` message
// interpolated into a refusal varies with how the network is failing, so
// the check is on the static halves the CLI itself writes, either side of
// whatever gh said.
func quoted(t *testing.T, doc, printed, variable string) {
	t.Helper()
	for _, part := range strings.Split(printed, variable) {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if !strings.Contains(doc, part) {
			t.Errorf("docs/%s does not quote what the CLI prints:\n%q", offlineDoc, part)
		}
	}
}

// TestOfflineDocQuotesTheUnreachableRefusal guards the refusal every verb
// that transacts with the record stops at when it is offline — the one the
// page's "What waits" section is built around.
func TestOfflineDocQuotesTheUnreachableRefusal(t *testing.T) {
	const ghSaid = "gh repo: error connecting to api.github.com"
	err := unreachable(errors.New(ghSaid))
	if err == nil {
		t.Fatal("unreachable() did not classify a gh connection failure")
	}
	var r refusal
	if !errors.As(err, &r) || r.Code != "GH_UNREACHABLE" {
		t.Fatalf("offline gh failure is %v, not a GH_UNREACHABLE refusal", err)
	}
	quoted(t, readOfflineDoc(t), r.Detail, ghSaid)
}

// TestOfflineDocQuotesTheLabelNote guards the other half of the page: the
// verbs whose work is local carry on offline and say so in a note. The
// note is produced here the way an offline run produces it — the label
// step's target lookup failing — rather than copied into the test.
func TestOfflineDocQuotesTheLabelNote(t *testing.T) {
	const ghSaid = "gh repo: error connecting to api.github.com"
	restore := labelTarget
	labelTarget = func() (tracker.Tracker, string, error) {
		return nil, "", errors.New(ghSaid)
	}
	defer func() { labelTarget = restore }()

	var out bytes.Buffer
	withLabelTarget(&out, func(tracker.Tracker, string) {
		t.Error("the label step ran although GitHub could not name the repository")
	})
	note := strings.TrimSpace(out.String())
	if note == "" {
		t.Fatal("an unreachable GitHub produced no note from the label step")
	}
	quoted(t, readOfflineDoc(t), note, ghSaid)
}
