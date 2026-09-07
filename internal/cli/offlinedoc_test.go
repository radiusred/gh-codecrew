package cli

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	codecrew "github.com/radiusred/gh-codecrew"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// docs/working-offline.md tells an operator what the CLI prints when it
// cannot reach GitHub, and a page that quotes output is only worth reading
// while the quotation is true. Each test below produces one of the five
// texts the page quotes verbatim from the code that prints them, and fails
// when the page has stopped carrying it — the shape
// TestRefusalCodesMatchTheSpecTable uses for SPEC §10's catalogue. They
// assert nothing about GitHub and run offline themselves.
//
// Six, because checky's reviews of PR #333 counted them: the first pass
// guarded two of five, and round two added migrate's dry-run footer.
const offlineDoc = "working-offline.md"

func readOfflineDoc(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "docs", offlineDoc))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// quoted checks the page carries a line the CLI prints. Everything the
// caller names as variable — the gh message interpolated into a refusal, a
// branch name, a login — is cut out first, leaving the static text the CLI
// itself writes; whitespace is flattened on both sides, since the page
// wraps its prose at 80 columns and a quotation may straddle a line break.
func quoted(t *testing.T, doc, printed string, variable ...string) {
	t.Helper()
	parts := []string{printed}
	for _, v := range variable {
		var next []string
		for _, p := range parts {
			next = append(next, strings.Split(p, v)...)
		}
		parts = next
	}
	flat := strings.Join(strings.Fields(doc), " ")
	for _, part := range parts {
		part = strings.Join(strings.Fields(part), " ")
		if part == "" {
			continue
		}
		if !strings.Contains(flat, part) {
			t.Errorf("docs/%s does not quote what the CLI prints:\n%q", offlineDoc, part)
		}
	}
}

// ghSaid stands in for whatever gh prints when the network is gone. The
// page says in prose that this part varies, so no test pins it.
const ghSaid = "gh repo: error connecting to api.github.com"

// TestOfflineDocQuotesTheUnreachableRefusal guards the refusal every verb
// that transacts with the record stops at when it is offline — the one the
// page's "What waits" section is built around.
func TestOfflineDocQuotesTheUnreachableRefusal(t *testing.T) {
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
	stubLabelTarget(t, nil, "", errors.New(ghSaid))

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

// TestOfflineDocQuotesTheBootstrapNote guards the second thing `init` says
// offline: the branch-protection probe went unanswered, so the scaffold was
// committed on the bootstrap branch rather than the one the operator was
// on. The page's claim that an offline `init` behaves differently from an
// online one in an unprotected repo rests on this line.
func TestOfflineDocQuotesTheBootstrapNote(t *testing.T) {
	stubProtection(t, false, false) // asked, and GitHub did not answer
	dir := gitRepo(t)
	// A born branch: an unborn HEAD cannot be branched from, so the
	// bootstrap path — and this note with it — needs a commit to exist.
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("one\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := git(dir, "add", "README.md"); err != nil {
		t.Fatal(err)
	}
	if _, err := git(dir, "commit", "-q", "-m", "initial"); err != nil {
		t.Fatal(err)
	}
	written, _, err := scaffold(dir, "self", fakeContracts)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	commitScaffold(&out, dir, written)

	note := ""
	for _, line := range strings.Split(out.String(), "\n") {
		if strings.HasPrefix(line, "note:") {
			note = line
		}
	}
	if note == "" {
		t.Fatalf("an unanswered branch-protection probe produced no note:\n%s", out.String())
	}
	quoted(t, readOfflineDoc(t), note)
}

// TestOfflineDocQuotesTheSpokeDiffMessage guards the one line in the
// "what waits" section that is not a network condition at all: `roles diff`
// in a repo holding no contracts fails the same way online and offline, and
// the page says so to keep an operator from reading it as being offline.
func TestOfflineDocQuotesTheSpokeDiffMessage(t *testing.T) {
	err := rolesDiff(io.Discard, t.TempDir(), codecrew.Roles, "implementer")
	if err == nil {
		t.Fatal("roles diff found a contract in a directory holding none")
	}
	quoted(t, readOfflineDoc(t), err.Error())
}

// TestOfflineDocQuotesTaskStartsReceipt guards the load-bearing quotation:
// the whole "what step 6 does to a branch you already have" section turns
// on `locally: git fetch && git switch <branch>` being what the verb
// prints, since following that line is the step that silently does nothing.
func TestOfflineDocQuotesTaskStartsReceipt(t *testing.T) {
	f := &startFake{task: startingTask(), viewer: "someone"}
	var out bytes.Buffer
	if err := runTaskStart(startCtx(f, crewRoles), &out, f.task.Ref); err != nil {
		t.Fatal(err)
	}
	if len(f.branches) != 1 {
		t.Fatalf("task start created %v, not one linked branch", f.branches)
	}
	quoted(t, readOfflineDoc(t), out.String(), f.branches[0], f.task.Ref.String(), f.viewer)
}

// TestOfflineDocQuotesTheMigrateDryRunFooter guards the one preview the
// page tells an operator to reach for offline. `migrate --dry-run` is
// local — it reads the repo, not GitHub — so it completes with no network
// and says so on its last line, which is what makes it worth naming beside
// three dry runs that do not.
func TestOfflineDocQuotesTheMigrateDryRunFooter(t *testing.T) {
	// Typed identities: an offline migrate that had to ask GitHub what a
	// bare 1.0 login is refuses instead, which is the other half of what
	// the page says about this verb.
	dir := legacyRepo(t, "codecrew: \"1.0\"\nhub: self\nroles:\n  implementer: {identity: user:alice}\n",
		map[string]string{"roles/qa.md": "# Role: qa\n"})
	stubLabelTarget(t, nil, "", errors.New(ghSaid))

	var out bytes.Buffer
	if err := migrate(&out, dir, true); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	footer := lines[len(lines)-1]
	if !strings.HasPrefix(footer, "dry run:") {
		t.Fatalf("the dry run did not end with its footer:\n%s", out.String())
	}
	quoted(t, readOfflineDoc(t), footer)
}
