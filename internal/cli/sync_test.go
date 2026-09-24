package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	codecrew "github.com/radiusred/gh-codecrew"
	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/gosrc"
)

// The texts the sync tests reason about: fakeContracts embeds
// "# Role: implementer\n" and "# Role: qa\n"; oldImplementer is what an
// earlier release shipped, and fakeHistory is that release's table.
const oldImplementer = "# Role: implementer\nOld text.\n"

func sha(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// oldAgents and newAgents stand in for a released .codecrew/AGENTS.md
// scaffold and the one the binary writes; fakeHistory carries the first.
const (
	oldAgents = "# Agents\nOld scaffold.\n"
	newAgents = "# Agents\nNew scaffold.\n"
)

var fakeHistory = []releasedContract{
	{Release: "v1.0.0", Role: "implementer", SHA256: sha("# Role: implementer\nOlder text.\n")},
	{Release: "v1.1.0", Role: "implementer", SHA256: sha(oldImplementer)},
	{Release: "v1.2.0", Role: "implementer", SHA256: sha(oldImplementer)},
	{Release: "v2.0.0", Role: config.AgentsFile, SHA256: sha(oldAgents)},
}

func writeAgents(t *testing.T, dir, content string) {
	t.Helper()
	p := filepath.Join(dir, filepath.FromSlash(config.AgentsFile))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readAgents(t *testing.T, dir string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(config.AgentsFile)))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func writeContract(t *testing.T, dir, role, content string) {
	t.Helper()
	p := filepath.Join(dir, rolesPath(role+".md"))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readContract(t *testing.T, dir, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, rolesPath(name)))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// writeEmbeddedContracts gives dir the binary's own contracts, so a hub in
// a test is current and status has nothing to say about it.
func writeEmbeddedContracts(t *testing.T, dir string) {
	t.Helper()
	roles, err := embeddedRoles(codecrew.Roles)
	if err != nil {
		t.Fatal(err)
	}
	for _, role := range roles {
		data, err := fs.ReadFile(codecrew.Roles, contractPath(role))
		if err != nil {
			t.Fatal(err)
		}
		writeContract(t, dir, role, string(data))
	}
}

// syncRepo is a git repo holding one commit, so a sync has a branch and a
// history to add to.
func syncRepo(t *testing.T) string {
	t.Helper()
	dir := gitRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hub\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mustGitIn(t, dir, "add", "README.md")
	mustGitIn(t, dir, "commit", "-q", "-m", "chore: start")
	return dir
}

func mustGitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := git(dir, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func stateOf(t *testing.T, statuses []contractStatus, role string) contractStatus {
	t.Helper()
	for _, s := range statuses {
		if s.Role == role {
			return s
		}
	}
	t.Fatalf("no status for %s in %v", role, statuses)
	return contractStatus{}
}

// A local contract is one of four things, and roles sync and status both
// read it the same way: the stamp and the line endings are not the text.
func TestClassifyContracts(t *testing.T) {
	cases := []struct {
		name, local string
		absent      bool
		want        contractState
		release     string
		stamped     bool
	}{
		{name: "absent", absent: true, want: contractAbsent},
		{name: "current", local: "# Role: implementer\n", want: contractCurrent},
		{name: "current, stamped", local: contractStamp("implementer.md") + "# Role: implementer\n", want: contractCurrent, stamped: true},
		{name: "a release's text, named by the latest release that shipped it", local: oldImplementer, want: contractRelease, release: "v1.2.0"},
		{name: "a release's text, stamped", local: contractStamp("implementer.md") + oldImplementer, want: contractRelease, release: "v1.2.0", stamped: true},
		{name: "a release's text, CRLF", local: strings.ReplaceAll(contractStamp("implementer.md")+oldImplementer, "\n", "\r\n"), want: contractRelease, release: "v1.2.0", stamped: true},
		{name: "a fork", local: oldImplementer + "- A local convention.\n", want: contractForked},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if !tc.absent {
				writeContract(t, dir, "implementer", tc.local)
			}
			statuses, err := classifyContracts(dir, fakeContracts, fakeHistory, []string{"implementer"})
			if err != nil {
				t.Fatal(err)
			}
			got := stateOf(t, statuses, "implementer")
			if got.State != tc.want || got.Release != tc.release || got.Stamped != tc.stamped {
				t.Errorf("got %+v, want state %d release %q stamped %v", got, tc.want, tc.release, tc.stamped)
			}
		})
	}
	if _, err := classifyContracts(t.TempDir(), fakeContracts, fakeHistory, []string{"navigator"}); err == nil {
		t.Error("a role the binary does not embed was classified")
	}
}

// The whole verb on a hub with one contract an earlier release's text and
// one absent: both written in one commit of exactly those paths, the body
// naming the tool and the version, the extension and the operator's other
// work untouched.
func TestRolesSyncWritesAbsentAndReleaseTextsInOneCommit(t *testing.T) {
	dir := syncRepo(t)
	writeContract(t, dir, "implementer", oldImplementer)
	ext := filepath.Join(dir, rolesPath("implementer"+localSuffix))
	if err := os.WriteFile(ext, []byte("- House style.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mustGitIn(t, dir, "add", ".")
	mustGitIn(t, dir, "commit", "-q", "-m", "chore: contracts")
	// The operator's work in progress: an edited extension and a staged file.
	if err := os.WriteFile(ext, []byte("- House style.\n- More.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mustGitIn(t, dir, "add", "notes.txt")

	var out bytes.Buffer
	if err := rolesSync(&out, dir, syncSource{contracts: fakeContracts}, fakeHistory, nil, false); err != nil {
		t.Fatal(err)
	}
	// Unstamped stays unstamped; absent is written as init writes it.
	if got := readContract(t, dir, "implementer.md"); got != "# Role: implementer\n" {
		t.Errorf("implementer.md = %q", got)
	}
	if got := readContract(t, dir, "qa.md"); got != contractStamp("qa.md")+"# Role: qa\n" {
		t.Errorf("qa.md = %q", got)
	}
	files := mustGitIn(t, dir, "show", "--name-only", "--format=", "HEAD")
	if files != contractPath("implementer")+"\n"+contractPath("qa") {
		t.Errorf("the commit holds %q, want exactly the two contracts", files)
	}
	msg := mustGitIn(t, dir, "log", "-1", "--format=%B")
	for _, want := range []string{
		"chore: sync codecrew role contracts to " + version,
		contractPath("implementer") + ": written (was the v1.2.0 text)",
		contractPath("qa") + ": written (absent)",
		"`gh codecrew roles sync`",
		"Housekeeping (SPEC §4)",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("commit message lacks %q:\n%s", want, msg)
		}
	}
	status := mustGitIn(t, dir, "status", "--porcelain")
	for _, want := range []string{"M " + filepath.ToSlash(rolesPath("implementer"+localSuffix)), "A  notes.txt"} {
		if !strings.Contains(status, want) {
			t.Errorf("the operator's work was disturbed; status lacks %q:\n%s", want, status)
		}
	}
	for _, want := range []string{"wrote " + contractPath("implementer") + " (was the v1.2.0 text)", "committed ", "never pushes"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, out.String())
		}
	}
}

// A fork anywhere among the named roles refuses before anything is
// written; naming the other roles syncs them past it.
func TestRolesSyncRefusesAForkBeforeWriting(t *testing.T) {
	dir := syncRepo(t)
	fork := oldImplementer + "- A local convention.\n"
	writeContract(t, dir, "implementer", fork)
	head := mustGitIn(t, dir, "rev-parse", "HEAD")

	err := rolesSync(&bytes.Buffer{}, dir, syncSource{contracts: fakeContracts}, fakeHistory, nil, false)
	var r refusal
	if !errors.As(err, &r) || r.Code != "CONTRACT_FORKED" {
		t.Fatalf("err = %v, want CONTRACT_FORKED", err)
	}
	for _, want := range []string{"gh codecrew roles diff implementer", localSuffix, "gh codecrew roles sync <role>"} {
		if !strings.Contains(r.Detail, want) {
			t.Errorf("detail lacks %q: %s", want, r.Detail)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, rolesPath("qa.md"))); err == nil {
		t.Error("qa.md was written by a refused sync")
	}
	if got := readContract(t, dir, "implementer.md"); got != fork {
		t.Error("the fork was touched")
	}
	if got := mustGitIn(t, dir, "rev-parse", "HEAD"); got != head {
		t.Error("a refused sync committed")
	}

	if err := rolesSync(&bytes.Buffer{}, dir, syncSource{contracts: fakeContracts}, fakeHistory, []string{"qa"}, false); err != nil {
		t.Fatalf("syncing the other role past the fork: %v", err)
	}
	if got := readContract(t, dir, "implementer.md"); got != fork {
		t.Error("the fork was touched by a sync naming another role")
	}
	if files := mustGitIn(t, dir, "show", "--name-only", "--format=", "HEAD"); files != contractPath("qa") {
		t.Errorf("the commit holds %q", files)
	}
}

func TestRolesSyncDryRunWritesNothing(t *testing.T) {
	dir := syncRepo(t)
	writeContract(t, dir, "implementer", oldImplementer)
	head := mustGitIn(t, dir, "rev-parse", "HEAD")
	var out bytes.Buffer
	if err := rolesSync(&out, dir, syncSource{contracts: fakeContracts}, fakeHistory, nil, true); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"would write " + contractPath("implementer") + " (was the v1.2.0 text)", "would write " + contractPath("qa") + " (absent)", "would commit 2 paths on main", "dry run: nothing written"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("dry run lacks %q:\n%s", want, out.String())
		}
	}
	if readContract(t, dir, "implementer.md") != oldImplementer {
		t.Error("the dry run wrote implementer.md")
	}
	if mustGitIn(t, dir, "rev-parse", "HEAD") != head {
		t.Error("the dry run committed")
	}
	// A fork refuses the dry run too, with the same code.
	writeContract(t, dir, "qa", "# Role: qa\nForked.\n")
	var r refusal
	if err := rolesSync(&bytes.Buffer{}, dir, syncSource{contracts: fakeContracts}, fakeHistory, nil, true); !errors.As(err, &r) || r.Code != "CONTRACT_FORKED" {
		t.Errorf("dry run over a fork: %v", err)
	}
}

// On the default branch the light path's commit goes on a branch of its
// own; one left over from an earlier sync stops the verb before it writes.
func TestRolesSyncCutsABranchOnTheDefaultBranch(t *testing.T) {
	dir := syncRepo(t)
	mustGitIn(t, dir, "update-ref", "refs/remotes/origin/main", "HEAD")
	mustGitIn(t, dir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	var out bytes.Buffer
	if err := rolesSync(&out, dir, syncSource{contracts: fakeContracts}, fakeHistory, []string{"qa"}, true); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "would commit 1 paths on "+syncBranch+", cut from main") {
		t.Errorf("dry run does not name the branch:\n%s", out.String())
	}
	if err := rolesSync(&bytes.Buffer{}, dir, syncSource{contracts: fakeContracts}, fakeHistory, []string{"qa"}, false); err != nil {
		t.Fatal(err)
	}
	if b := mustGitIn(t, dir, "symbolic-ref", "--short", "HEAD"); b != syncBranch {
		t.Errorf("committed on %s, want %s", b, syncBranch)
	}
	if mustGitIn(t, dir, "rev-parse", "main") == mustGitIn(t, dir, "rev-parse", "HEAD") {
		t.Error("the default branch moved")
	}

	mustGitIn(t, dir, "switch", "-q", "main")
	err := rolesSync(&bytes.Buffer{}, dir, syncSource{contracts: fakeContracts}, fakeHistory, []string{"qa"}, false)
	if err == nil || !strings.Contains(err.Error(), syncBranch) {
		t.Errorf("a leftover %s was not reported: %v", syncBranch, err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, rolesPath("qa.md"))); statErr == nil {
		t.Error("qa.md was written although the verb stopped")
	}

	// Off the default branch the commit lands where the operator is.
	mustGitIn(t, dir, "switch", "-q", "-c", "chores")
	if err := rolesSync(&bytes.Buffer{}, dir, syncSource{contracts: fakeContracts}, fakeHistory, []string{"qa"}, false); err != nil {
		t.Fatal(err)
	}
	if b := mustGitIn(t, dir, "symbolic-ref", "--short", "HEAD"); b != "chores" {
		t.Errorf("committed on %s, want chores", b)
	}
}

func TestRolesSyncWithNothingToDoCommitsNothing(t *testing.T) {
	dir := syncRepo(t)
	writeContract(t, dir, "implementer", "# Role: implementer\n")
	writeContract(t, dir, "qa", contractStamp("qa.md")+"# Role: qa\n")
	head := mustGitIn(t, dir, "rev-parse", "HEAD")
	var out bytes.Buffer
	if err := rolesSync(&out, dir, syncSource{contracts: fakeContracts}, fakeHistory, nil, false); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "nothing to write") {
		t.Errorf("output:\n%s", out.String())
	}
	if mustGitIn(t, dir, "rev-parse", "HEAD") != head {
		t.Error("a sync with nothing to write committed")
	}
}

// A spoke holds no contracts, so naming a role there is still an error
// without a code, as roles diff's is; the agents file is the one thing a
// spoke's roles sync writes (#372).
func TestRolesSyncInASpokeWritesOnlyTheAgentsFile(t *testing.T) {
	err := rolesSyncCmd(&bytes.Buffer{}, []string{"qa", "--dry-run"}, "o/hub", t.TempDir(), fakeContracts)
	if err == nil || !strings.Contains(err.Error(), "spoke of o/hub") {
		t.Errorf("err = %v", err)
	}
	var r refusal
	if errors.As(err, &r) {
		t.Errorf("a spoke exits without a code, as roles diff does: %v", err)
	}

	dir := syncRepo(t)
	writeAgents(t, dir, oldAgents)
	mustGitIn(t, dir, "add", "-A")
	mustGitIn(t, dir, "commit", "-q", "-m", "chore: agents")
	var out bytes.Buffer
	if err := rolesSync(&out, dir, syncSource{agents: newAgents}, fakeHistory, nil, false); err != nil {
		t.Fatal(err)
	}
	if got := readAgents(t, dir); got != newAgents {
		t.Errorf("agents file = %q", got)
	}
	if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(config.RolesDir))); err == nil {
		t.Error("a spoke's sync wrote contracts")
	}
	if subj := mustGitIn(t, dir, "log", "-1", "--format=%s"); subj != "chore: sync the codecrew agents file to "+version {
		t.Errorf("subject = %q", subj)
	}
	if !strings.Contains(out.String(), config.AgentsFile+" (was the v2.0.0 text)") {
		t.Errorf("output:\n%s", out.String())
	}
}

// In a hub a bare sync writes the contracts and the agents file as one
// commit; a role list leaves the agents file alone, and its path names it
// alone (#372).
func TestRolesSyncCarriesTheAgentsFile(t *testing.T) {
	dir := syncRepo(t)
	src := syncSource{contracts: fakeContracts, agents: newAgents}
	var out bytes.Buffer
	if err := rolesSync(&out, dir, src, fakeHistory, nil, false); err != nil {
		t.Fatal(err)
	}
	if got := readAgents(t, dir); got != newAgents {
		t.Errorf("agents file = %q", got)
	}
	if subj := mustGitIn(t, dir, "log", "-1", "--format=%s"); subj != "chore: sync codecrew role contracts and agents file to "+version {
		t.Errorf("subject = %q", subj)
	}
	files := mustGitIn(t, dir, "show", "--name-only", "--format=", "HEAD")
	for _, want := range []string{config.AgentsFile, contractPath("implementer"), contractPath("qa")} {
		if !strings.Contains(files, want) {
			t.Errorf("the commit lacks %s:\n%s", want, files)
		}
	}
	if body := mustGitIn(t, dir, "log", "-1", "--format=%b"); !strings.Contains(body, "Target: the role contracts and "+config.AgentsFile) {
		t.Errorf("body = %q", body)
	}

	// A role list never touches the agents file; the path alone does.
	dir = syncRepo(t)
	writeAgents(t, dir, oldAgents)
	if err := rolesSync(&bytes.Buffer{}, dir, src, fakeHistory, []string{"qa"}, false); err != nil {
		t.Fatal(err)
	}
	if got := readAgents(t, dir); got != oldAgents {
		t.Errorf("a role list rewrote the agents file: %q", got)
	}
	if err := rolesSync(&bytes.Buffer{}, dir, src, fakeHistory, []string{config.AgentsFile}, false); err != nil {
		t.Fatal(err)
	}
	if got := readAgents(t, dir); got != newAgents {
		t.Errorf("naming the path did not sync it: %q", got)
	}
	if files := mustGitIn(t, dir, "show", "--name-only", "--format=", "HEAD"); strings.TrimSpace(files) != config.AgentsFile {
		t.Errorf("naming the path committed more than the agents file:\n%s", files)
	}
}

// A hand-written agents file is the project's own: alone it refuses
// AGENTS_FORKED, beside a forked contract CONTRACT_FORKED names both, and
// either way nothing is written — dry run included (#372).
func TestRolesSyncNeverOverwritesAnAgentsFileOfTheProjects(t *testing.T) {
	src := syncSource{contracts: fakeContracts, agents: newAgents}
	for _, c := range []struct {
		name, implementer, code string
		dryRun                  bool
	}{
		{"agents alone", "# Role: implementer\n", "AGENTS_FORKED", false},
		{"agents alone, dry run", "# Role: implementer\n", "AGENTS_FORKED", true},
		{"beside a forked contract", "# Role: implementer\nOurs.\n", "CONTRACT_FORKED", false},
	} {
		t.Run(c.name, func(t *testing.T) {
			dir := syncRepo(t)
			writeContract(t, dir, "implementer", c.implementer)
			writeAgents(t, dir, "# Our own agents file\n")
			head := mustGitIn(t, dir, "rev-parse", "HEAD")
			var r refusal
			err := rolesSync(&bytes.Buffer{}, dir, src, fakeHistory, nil, c.dryRun)
			if !errors.As(err, &r) || r.Code != c.code {
				t.Fatalf("err = %v, want %s", err, c.code)
			}
			if !strings.Contains(err.Error(), "gh codecrew roles diff "+config.AgentsFile) {
				t.Errorf("the refusal does not name the agents file's diff: %v", err)
			}
			if readAgents(t, dir) != "# Our own agents file\n" || readContract(t, dir, "implementer.md") != c.implementer {
				t.Error("a refused sync wrote a file")
			}
			if _, err := os.Stat(filepath.Join(dir, rolesPath("qa.md"))); err == nil {
				t.Error("a refused sync wrote the absent qa contract")
			}
			if mustGitIn(t, dir, "rev-parse", "HEAD") != head {
				t.Error("a refused sync committed")
			}
		})
	}
	// In a spoke the way out names no roles: there are none to name.
	dir := syncRepo(t)
	writeAgents(t, dir, "# Our own agents file\n")
	err := rolesSync(&bytes.Buffer{}, dir, syncSource{agents: newAgents}, fakeHistory, nil, false)
	var r refusal
	if !errors.As(err, &r) || r.Code != "AGENTS_FORKED" || strings.Contains(err.Error(), "naming the roles") {
		t.Errorf("err = %v", err)
	}
}

// The agents file is measured as a contract is: CRLF is not the text, and a
// release's scaffold is recognised by the history row keyed on its path.
func TestClassifyAgents(t *testing.T) {
	for _, c := range []struct {
		name, local string
		absent      bool
		want        contractState
		release     string
	}{
		{"absent", "", true, contractAbsent, ""},
		{"current", newAgents, false, contractCurrent, ""},
		{"current, CRLF", strings.ReplaceAll(newAgents, "\n", "\r\n"), false, contractCurrent, ""},
		{"a release's", oldAgents, false, contractRelease, "v2.0.0"},
		{"the project's own", "# Ours\n", false, contractForked, ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			if !c.absent {
				writeAgents(t, dir, c.local)
			}
			st, err := classifyAgents(dir, newAgents, fakeHistory)
			if err != nil {
				t.Fatal(err)
			}
			if st.State != c.want || st.Release != c.release {
				t.Errorf("got %+v", st)
			}
		})
	}
}

// Dry run names the agents file among what it would write and writes none.
func TestRolesSyncDryRunCoversTheAgentsFile(t *testing.T) {
	dir := syncRepo(t)
	head := mustGitIn(t, dir, "rev-parse", "HEAD")
	var out bytes.Buffer
	if err := rolesSync(&out, dir, syncSource{contracts: fakeContracts, agents: newAgents}, fakeHistory, nil, true); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "would write "+config.AgentsFile+" (absent)") {
		t.Errorf("output:\n%s", out.String())
	}
	if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(config.AgentsFile))); err == nil || mustGitIn(t, dir, "rev-parse", "HEAD") != head {
		t.Error("a dry run wrote or committed")
	}
}

// contractHistory must hold every tagged release's contract texts, or roles
// sync calls a project that is merely behind a fork. The table is
// regenerated by scripts/contract-history; this test checks it against the
// tags in this clone. A clone without tags — a pull request's shallow
// checkout — skips, unless CODECREW_REQUIRE_TAGS is set, as the release
// workflow sets it: a release cannot build while the table misses one. A
// release's text that is the embedded contract needs no row, because the
// embed is recognised on its own.
func TestContractHistoryCoversEveryTag(t *testing.T) {
	root := filepath.Join("..", "..")
	out, err := exec.Command("git", "-C", root, "tag", "-l", "v[0-9]*").Output()
	tags := strings.Fields(string(out))
	if err != nil || len(tags) == 0 {
		if os.Getenv("CODECREW_REQUIRE_TAGS") != "" {
			t.Fatalf("CODECREW_REQUIRE_TAGS is set and this clone has no release tags (%v): fetch them (fetch-depth: 0)", err)
		}
		t.Skip("no release tags in this clone")
	}
	roles, err := embeddedRoles(codecrew.Roles)
	if err != nil {
		t.Fatal(err)
	}
	have := map[string]bool{}
	for _, h := range contractHistory {
		have[h.Release+" "+h.Role+" "+h.SHA256] = true
	}
	for _, tag := range tags {
		for _, role := range roles {
			var text []byte
			for _, p := range []string{contractPath(role), config.LegacyRolesDir + "/" + role + ".md"} {
				if b, err := exec.Command("git", "-C", root, "show", tag+":"+p).Output(); err == nil {
					text = b
					break
				}
			}
			if text == nil {
				continue // the role did not exist yet
			}
			embedded, _ := fs.ReadFile(codecrew.Roles, contractPath(role))
			if string(text) == string(embedded) {
				continue
			}
			if !have[tag+" "+role+" "+sha(string(text))] {
				t.Errorf("contractHistory lacks %s's %s contract: run scripts/contract-history", tag, role)
			}
		}
		// A release on the 2.0 layout wrote .codecrew/AGENTS.md from the
		// scaffold constant in its init.go, read here as the script reads
		// it (#372).
		if exec.Command("git", "-C", root, "cat-file", "-e", tag+":"+config.Pointer).Run() != nil {
			continue
		}
		src, err := exec.Command("git", "-C", root, "show", tag+":internal/cli/init.go").Output()
		if err != nil {
			t.Fatalf("%s: %v", tag, err)
		}
		text, err := gosrc.StringConst(src, "agentsScaffold")
		if err != nil {
			t.Fatalf("%s's agentsScaffold cannot be read from source: %v", tag, err)
		}
		if text != agentsScaffold && !have[tag+" "+config.AgentsFile+" "+sha(text)] {
			t.Errorf("contractHistory lacks %s's %s scaffold: run scripts/contract-history", tag, config.AgentsFile)
		}
	}
}

// The script and the guard recover a release's agents file by evaluating
// the scaffold constant from source; that reading of this tree's init.go
// is the constant the binary carries, byte for byte (#372).
func TestSourceReadingOfTheScaffoldIsTheScaffold(t *testing.T) {
	src, err := os.ReadFile("init.go")
	if err != nil {
		t.Fatal(err)
	}
	got, err := gosrc.StringConst(src, "agentsScaffold")
	if err != nil {
		t.Fatal(err)
	}
	if got != agentsScaffold {
		t.Errorf("source reading differs from the constant:\n%s", unifiedDiff(got, agentsScaffold))
	}
}
