package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// legacy1x is a protocol 1.0 pointer as v1 scaffolded and hands edited it:
// bare identities, no coordinator row, comments and blank lines that a
// human maintains and the migration must not eat.
const legacy1x = `codecrew: "1.0" # protocol version (SPEC.md §5): a different major is refused
hub: self

# Advisory role routing, read by whoever dispatches agents (see SPEC.md §5).
roles:
  implementer: { identity: myorg-coder }
  reviewer: { identity: alice }
  qa: { identity: myorg/qa-crew }
  doc-synthesizer: { identity: ~ }
`

// stubAccounts stands in for the GitHub users API: a login in the map
// resolves to that account type, anything else is a 404.
func stubAccounts(t *testing.T, accounts map[string]string) {
	t.Helper()
	prev := lookupAccount
	lookupAccount = func(login string) (string, bool, error) {
		kind, ok := accounts[login]
		return kind, ok, nil
	}
	t.Cleanup(func() { lookupAccount = prev })
}

// legacyRepo builds a committed 1.x repository: the pointer, and each
// named file under roles/ (or wherever the name says).
func legacyRepo(t *testing.T, pointer string, files map[string]string) string {
	t.Helper()
	// A temp repository has no origin, so gh could not name it anyway —
	// and asking the real one would put a network call inside every
	// migrate test. The label step prints its note and migrate carries on,
	// which is the point of the note; the two tests that exercise the step
	// stub a target of their own over this one.
	stubLabelTarget(t, nil, "", errors.New("no git remotes found"))
	dir := gitRepo(t)
	write := func(rel, content string) {
		path := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if pointer != "" {
		write(config.LegacyPointer, pointer)
	}
	for rel, content := range files {
		write(rel, content)
	}
	write("README.md", "# project\n")
	if _, err := git(dir, "add", "-A"); err != nil {
		t.Fatal(err)
	}
	if _, err := git(dir, "commit", "-q", "-m", "init"); err != nil {
		t.Fatal(err)
	}
	return dir
}

func exists(t *testing.T, dir, rel string) bool {
	t.Helper()
	_, err := os.Stat(filepath.Join(dir, filepath.FromSlash(rel)))
	return err == nil
}

func read(t *testing.T, dir, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func headSubject(t *testing.T, dir string) string {
	t.Helper()
	out, err := git(dir, "log", "-1", "--format=%s")
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func refusalCode(t *testing.T, err error) string {
	t.Helper()
	var r refusal
	if !errors.As(err, &r) {
		t.Fatalf("error %v is not a refusal", err)
	}
	return r.Code
}

// A full 1.x hub moves whole: the pointer and all ten role files, the
// version bumped, every identity typed by kind, the coordinator row added,
// and the emptied roles/ gone — in one commit that leaves the operator's
// own work alone.
func TestMigrateHub(t *testing.T) {
	files := map[string]string{}
	for _, role := range config.RoleNames {
		files["roles/"+role+".md"] = "# Role: " + role + "\n"
		files["roles/"+role+localSuffix] = "<!-- " + role + " -->\n"
	}
	dir := legacyRepo(t, legacy1x, files)
	stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})

	// The operator's own work, staged and unstaged, before the migration.
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# project\nedited\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.md"), []byte("staged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := git(dir, "add", "notes.md"); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := migrate(&out, dir, false); err != nil {
		t.Fatal(err)
	}
	for _, role := range config.RoleNames {
		for _, name := range []string{role + ".md", role + localSuffix} {
			if !exists(t, dir, config.RolesDir+"/"+name) {
				t.Errorf("%s/%s not moved", config.RolesDir, name)
			}
		}
	}
	if exists(t, dir, config.LegacyPointer) || exists(t, dir, config.LegacyRolesDir) {
		t.Errorf("1.x layout left behind:\n%s", out.String())
	}
	pointer := read(t, dir, config.Pointer)
	for _, want := range []string{`codecrew: "2.0"`, "app:myorg-coder", "user:alice", "team:myorg/qa-crew", "coordinator"} {
		if !strings.Contains(pointer, want) {
			t.Errorf("rewritten pointer missing %q:\n%s", want, pointer)
		}
	}
	// The pointer must be one this binary can now read.
	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatalf("migrated pointer does not load: %v\n%s", err, pointer)
	}
	if cfg.Roles["coordinator"].Identity.Kind != config.KindOperator {
		t.Errorf("coordinator = %v, want the operator", cfg.Roles["coordinator"].Identity)
	}
	if got := headSubject(t, dir); got != migrateSubject {
		t.Errorf("HEAD subject = %q, want %q", got, migrateSubject)
	}
	status, err := git(dir, "status", "--short")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(status, "README.md") || !strings.Contains(status, "notes.md") {
		t.Errorf("the operator's own work did not survive the commit: %q", status)
	}
	if strings.Contains(status, config.Pointer) {
		t.Errorf("the migration was not committed: %q", status)
	}
}

// keptLine is the file list from the action-needed block, or "" when the
// block is not there. Read off the line itself rather than searched for in
// the whole output: the skip and write reports above it already name every
// root file, so a Contains over the output would pass whatever the heading
// said (checky's finding on PR #278).
func keptLine(out string) string {
	for _, line := range strings.Split(out, "\n") {
		if after, ok := strings.CutPrefix(line, "Kept: "); ok {
			return after
		}
	}
	return ""
}

// The root entry points are migrate's business exactly as they are init's:
// a file that is not there is CodeCrew's to write — from the same scaffold
// constant, in the migration's own commit — and a file that is there is the
// project's prose, kept byte for byte and reported only when it arrives
// nowhere. A 1.x spoke never had either file, so before #301 five of twelve
// fleet repos were migrated into a state a fresh init would never produce.
//
// The states a 1.x repo can present are one table, run twice: live, and
// again on an identical repo with --dry-run, because a preview that
// disagrees with the run it previews is not a preview (M15-R1, #301).
func TestMigrateWritesTheEntryPoints(t *testing.T) {
	// A 1.x root AGENTS.md: the instructions themselves, naming the paths
	// the migration just moved. It reaches nothing under 2.0.
	const legacyRoot = "# Agents\n\nRead `roles/implementer.md` before doing anything else.\n"
	const reachingRoot = "# Agents\n\n" + entryPointLines

	for _, c := range []struct {
		name string
		// the root files on disk before the migration
		existing map[string]string
		// the root files migrate must write, from their scaffolds
		write []string
		// the files the action-needed block must name, exactly; none
		// means the block must not appear at all
		kept []string
	}{{
		name:  "a 1.x spoke, which never had a root entry point at all",
		write: []string{"AGENTS.md", "CLAUDE.md"},
	}, {
		name:     "a 1.x hub whose root AGENTS.md holds the old instructions",
		existing: map[string]string{"AGENTS.md": legacyRoot},
		write:    []string{"CLAUDE.md"},
		kept:     []string{"AGENTS.md"},
	}, {
		name:     "both root files kept, neither reaching",
		existing: map[string]string{"AGENTS.md": legacyRoot, "CLAUDE.md": "@AGENTS.md\n"},
		kept:     []string{"AGENTS.md", "CLAUDE.md"},
	}, {
		name:     "a kept CLAUDE.md reaching through its AGENTS.md import",
		existing: map[string]string{"AGENTS.md": reachingRoot, "CLAUDE.md": "@AGENTS.md\n"},
	}, {
		name:     "a kept AGENTS.md that reaches, and no CLAUDE.md",
		existing: map[string]string{"AGENTS.md": reachingRoot},
		write:    []string{"CLAUDE.md"},
	}, {
		// The one shape the pending set exists for: the kept file's
		// import lands on a root AGENTS.md that is not there yet and is
		// about to be. Judged against disk alone it reaches nothing, and
		// the migration asks the operator to hand-edit a file that
		// reaches the instructions the moment it writes the AGENTS.md it
		// just said it would write (checky's finding 1 on PR #315).
		name:     "a kept CLAUDE.md reaching through an AGENTS.md the migration writes",
		existing: map[string]string{"CLAUDE.md": "@AGENTS.md\n"},
		write:    []string{"AGENTS.md"},
	}} {
		for _, dry := range []bool{false, true} {
			files := map[string]string{"roles/qa.md": "# Role: qa\n"}
			for name, content := range c.existing {
				files[name] = content
			}
			dir := legacyRepo(t, legacy1x, files)
			stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})

			var out bytes.Buffer
			if err := migrate(&out, dir, dry); err != nil {
				t.Fatalf("%s (dry=%v): %v", c.name, dry, err)
			}
			got := out.String()
			label := c.name
			if dry {
				label += " (dry run)"
			}

			// CodeCrew's own entry point, and every root file that was
			// missing, are written from the constants init writes them
			// from — or, in a dry run, listed and not written.
			verb := "wrote "
			if dry {
				verb = "would write "
			}
			for _, f := range append([]string{config.AgentsFile}, c.write...) {
				if !strings.Contains(got, verb+f) {
					t.Errorf("%s: output does not report %q:\n%s", label, verb+f, got)
				}
			}
			for f, want := range map[string]string{
				config.AgentsFile: agentsScaffold,
				"AGENTS.md":       rootEntryPointScaffolds["AGENTS.md"],
				"CLAUDE.md":       rootEntryPointScaffolds["CLAUDE.md"],
			} {
				if f != config.AgentsFile && !slices.Contains(c.write, f) {
					continue
				}
				if dry {
					if exists(t, dir, f) {
						t.Errorf("%s: the dry run wrote %s", label, f)
					}
					continue
				}
				if got := read(t, dir, f); got != want {
					t.Errorf("%s: %s = %q, want the scaffold init writes", label, f, got)
				}
			}
			// A file the project already had is never rewritten.
			for f, want := range c.existing {
				if got := read(t, dir, f); got != want {
					t.Errorf("%s: the project's %s was rewritten: %q", label, f, got)
				}
			}
			// The action needed block names the kept files and nothing
			// else — an absent one is written, not asked for.
			if want := strings.Join(c.kept, ", "); keptLine(got) != want {
				t.Errorf("%s: action needed reads \"Kept: %s\", want %q:\n%s", label, keptLine(got), want, got)
			}
			if len(c.kept) == 0 {
				if strings.Contains(got, "action needed") {
					t.Errorf("%s: migrate asked for an action it does not need:\n%s", label, got)
				}
				continue
			}
			if !strings.Contains(got, entryPointLines) {
				t.Errorf("%s: the block printed no lines to add:\n%s", label, got)
			}
		}
	}
}

// Everything migrate writes rides the migration's own commit, not the
// operator's next one, and the commit is still limited to the migration:
// work the operator had staged or unstaged is untouched.
func TestMigrateCommitsTheEntryPointsItWrites(t *testing.T) {
	dir := legacyRepo(t, legacy1x, map[string]string{"roles/qa.md": "# Role: qa\n"})
	stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})

	var out bytes.Buffer
	if err := migrate(&out, dir, false); err != nil {
		t.Fatal(err)
	}
	files := committedFiles(t, dir, "HEAD")
	for _, f := range append([]string{config.AgentsFile}, rootEntryPoints...) {
		if !strings.Contains(files, f) {
			t.Errorf("%s is not in the migration commit:\n%s", f, files)
		}
	}
	status, err := git(dir, "status", "--short")
	if err != nil {
		t.Fatal(err)
	}
	if status != "" {
		t.Errorf("the migration left work uncommitted: %q", status)
	}
}

// An older scaffold — four contracts, no extensions — moves what is there
// and invents nothing.
func TestMigrateHubWithoutExtensions(t *testing.T) {
	files := map[string]string{}
	for _, role := range []string{"implementer", "reviewer", "qa", "doc-synthesizer"} {
		files["roles/"+role+".md"] = "# Role: " + role + "\n"
	}
	dir := legacyRepo(t, legacy1x, files)
	stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})

	var out bytes.Buffer
	if err := migrate(&out, dir, false); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(filepath.Join(dir, filepath.FromSlash(config.RolesDir)))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, e := range entries {
		got = append(got, e.Name())
	}
	want := []string{"doc-synthesizer.md", "implementer.md", "qa.md", "reviewer.md"}
	if !slices.Equal(got, want) {
		t.Errorf("%s holds %v, want exactly %v", config.RolesDir, got, want)
	}
}

// An unrecognised entry beside CodeCrew's own is refused, by name, with
// nothing written: migrate moves ten known names and does not guess about
// an eleventh.
func TestMigrateForeignRolesEntry(t *testing.T) {
	dir := legacyRepo(t, legacy1x, map[string]string{
		"roles/qa.md":        "# Role: qa\n",
		"roles/webserver.md": "# our own\n",
	})
	stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})

	var out bytes.Buffer
	err := migrate(&out, dir, false)
	if code := refusalCode(t, err); code != "FOREIGN_ROLES_DIR" {
		t.Fatalf("code = %s, want FOREIGN_ROLES_DIR (%v)", code, err)
	}
	if !strings.Contains(err.Error(), "webserver.md") {
		t.Errorf("the refusal does not name the entry: %v", err)
	}
	if !exists(t, dir, config.LegacyPointer) || exists(t, dir, config.Pointer) {
		t.Error("a refused migration wrote something")
	}
	if headSubject(t, dir) == migrateSubject {
		t.Error("a refused migration committed")
	}
}

// A project's own roles/ — no CodeCrew file in it — is never read and
// never moved. That collision is the whole reason for the 2.0 layout.
func TestMigrateLeavesAProjectsOwnRolesDir(t *testing.T) {
	dir := legacyRepo(t, legacy1x, map[string]string{
		"roles/webserver/tasks/main.yml": "- name: install\n",
	})
	stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})

	var out bytes.Buffer
	if err := migrate(&out, dir, false); err != nil {
		t.Fatal(err)
	}
	if !exists(t, dir, "roles/webserver/tasks/main.yml") {
		t.Error("the project's own roles/ was touched")
	}
	if !exists(t, dir, config.Pointer) {
		t.Error("the pointer did not move")
	}
	if exists(t, dir, config.RolesDir) {
		t.Errorf("%s was created for a repo with no contracts", config.RolesDir)
	}
}

// Both layouts at once is the one case migrate must not resolve on its own.
func TestMigrateBothLayouts(t *testing.T) {
	dir := legacyRepo(t, legacy1x, map[string]string{
		config.Pointer: "codecrew: \"2.0\"\nhub: self\n",
	})
	var out bytes.Buffer
	err := migrate(&out, dir, false)
	if code := refusalCode(t, err); code != "BOTH_LAYOUTS" {
		t.Fatalf("code = %s, want BOTH_LAYOUTS (%v)", code, err)
	}
	for _, want := range []string{config.LegacyPointer, config.Pointer} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the refusal does not name %s: %v", want, err)
		}
	}
}

// A 2.0 file already sitting where a 1.x one would go is both layouts one
// level down: migrate refuses before writing rather than renaming over the
// newer file and stranding the tracked source's deletion outside its own
// commit (checky's finding on PR #280).
func TestMigrateRefusesAnExistingDestination(t *testing.T) {
	dir := legacyRepo(t, legacy1x, map[string]string{
		"roles/qa.md":              "# Role: qa (1.x)\n",
		config.RolesDir + "/qa.md": "# Role: qa (2.0)\n",
	})
	stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})

	var out bytes.Buffer
	err := migrate(&out, dir, false)
	if code := refusalCode(t, err); code != "BOTH_LAYOUTS" {
		t.Fatalf("code = %s, want BOTH_LAYOUTS (%v)", code, err)
	}
	for _, want := range []string{config.RolesDir + "/qa.md", "roles/qa.md"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the refusal does not name %s: %v", want, err)
		}
	}
	if got := read(t, dir, config.RolesDir+"/qa.md"); got != "# Role: qa (2.0)\n" {
		t.Errorf("the existing 2.0 file was overwritten: %q", got)
	}
	if !exists(t, dir, "roles/qa.md") || !exists(t, dir, config.LegacyPointer) {
		t.Error("a refused migration moved files")
	}
	status, err := git(dir, "status", "--short")
	if err != nil {
		t.Fatal(err)
	}
	if status != "" {
		t.Errorf("a refused migration touched the index or the tree: %q", status)
	}
}

// A source git does not track is renamed on the filesystem and left out of
// the commit's pathspec, which is the only case that may skip git mv.
func TestMigrateMovesAnUntrackedPointer(t *testing.T) {
	dir := legacyRepo(t, "", nil)
	if err := os.WriteFile(filepath.Join(dir, config.LegacyPointer), []byte("codecrew: \"1.0\"\nhub: myorg/hub\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := migrate(&out, dir, false); err != nil {
		t.Fatal(err)
	}
	if !exists(t, dir, config.Pointer) || exists(t, dir, config.LegacyPointer) {
		t.Error("the untracked pointer did not move")
	}
	if got := headSubject(t, dir); got != migrateSubject {
		t.Errorf("HEAD subject = %q, want %q", got, migrateSubject)
	}
}

// A spoke holds the pointer and, at most, its own extensions.
func TestMigrateSpoke(t *testing.T) {
	dir := legacyRepo(t, "codecrew: \"1.0\"\nhub: myorg/hub\n", map[string]string{
		"roles/implementer" + localSuffix: "<!-- ours -->\n",
	})
	var out bytes.Buffer
	if err := migrate(&out, dir, false); err != nil {
		t.Fatal(err)
	}
	pointer := read(t, dir, config.Pointer)
	if !strings.Contains(pointer, `codecrew: "2.0"`) || !strings.Contains(pointer, "hub: myorg/hub") {
		t.Errorf("spoke pointer = %q", pointer)
	}
	if strings.Contains(pointer, "coordinator") {
		t.Errorf("a coordinator row was invented for a table-less spoke: %q", pointer)
	}
	if !exists(t, dir, config.RolesDir+"/implementer"+localSuffix) {
		t.Error("the spoke's extension did not move")
	}
	// The end state is a pointer this binary reads: a one-shot move that
	// leaves the next verb refusing has not landed.
	if _, err := config.Load(dir); err != nil {
		t.Errorf("the migrated spoke does not load: %v", err)
	}
}

// Protocol 1.0 let a spoke carry a routing table and 2.0 does not, so
// migrating one forward would write a pointer this same binary refuses on
// the next verb. The block is the operator's routing, not migrate's to
// drop, so the move stops before it starts (checky's second-round finding
// on PR #280).
func TestMigrateRefusesASpokeCarryingRouting(t *testing.T) {
	dir := legacyRepo(t, "codecrew: \"1.0\"\nhub: myorg/hub\nroles:\n  reviewer: { identity: myorg-checky }\n  qa: { identity: alice }\n", map[string]string{
		"roles/implementer" + localSuffix: "<!-- ours -->\n",
	})
	stubAccounts(t, map[string]string{"myorg-checky[bot]": "Bot", "alice": "User"})
	before := headSubject(t, dir)

	var out bytes.Buffer
	err := migrate(&out, dir, false)
	if code := refusalCode(t, err); code != "SPOKE_ROUTING" {
		t.Fatalf("code = %s, want SPOKE_ROUTING (%v)", code, err)
	}
	for _, want := range []string{"myorg/hub", "qa, reviewer"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the refusal does not name %q: %v", want, err)
		}
	}
	if exists(t, dir, config.Pointer) || exists(t, dir, config.AgentsFile) || !exists(t, dir, config.LegacyPointer) {
		t.Error("a refused migration wrote something")
	}
	if headSubject(t, dir) != before {
		t.Error("a refused migration committed")
	}
	status, err2 := git(dir, "status", "--short")
	if err2 != nil {
		t.Fatal(err2)
	}
	if status != "" {
		t.Errorf("a refused migration touched the tree: %q", status)
	}
}

// A spoke with no table migrates as before, and one whose `roles:` key is
// empty is not carrying a table at all — the reading config.Parse takes.
func TestMigrateSpokeWithAnEmptyRolesKey(t *testing.T) {
	dir := legacyRepo(t, "codecrew: \"1.0\"\nhub: myorg/hub\nroles:\n", nil)
	var out bytes.Buffer
	if err := migrate(&out, dir, false); err != nil {
		t.Fatal(err)
	}
	if _, err := config.Load(dir); err != nil {
		t.Errorf("the migrated spoke does not load: %v", err)
	}
}

// A repo already on 2.0 moves nothing: it says so, commits nothing and
// exits 0, so a rerun is safe.
func TestMigrateAlreadyCurrent(t *testing.T) {
	dir := legacyRepo(t, "", map[string]string{
		config.Pointer: "codecrew: \"2.0\"\nhub: self\n",
	})
	before := headSubject(t, dir)
	var out bytes.Buffer
	if err := migrate(&out, dir, false); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "already on the protocol 2.0 layout") {
		t.Errorf("output = %q", out.String())
	}
	if headSubject(t, dir) != before {
		t.Error("an idempotent rerun committed")
	}
}

// A pointer below protocol 1.0 is out of scope, and the refusal names the
// version rather than moving files it cannot vouch for.
func TestMigrateProtocolBelowOne(t *testing.T) {
	dir := legacyRepo(t, "codecrew: \"0.1\"\nhub: self\n", nil)
	var out bytes.Buffer
	err := migrate(&out, dir, false)
	if code := refusalCode(t, err); code != "MIGRATION_UNSUPPORTED" {
		t.Fatalf("code = %s, want MIGRATION_UNSUPPORTED (%v)", code, err)
	}
	if !strings.Contains(err.Error(), "0.1") {
		t.Errorf("the refusal does not name the version: %v", err)
	}
	if exists(t, dir, config.Pointer) {
		t.Error("a refused migration wrote something")
	}
}

// A repo with neither layout is not a CodeCrew repo, and says so.
func TestMigrateNotACodecrewRepo(t *testing.T) {
	dir := legacyRepo(t, "", nil)
	var out bytes.Buffer
	err := migrate(&out, dir, false)
	if err == nil || !strings.Contains(err.Error(), "not a CodeCrew repo") {
		t.Fatalf("err = %v", err)
	}
}

// Contracts with no pointer say nothing about hub or spoke, so migrate
// stops rather than guessing.
func TestMigrateContractsWithoutPointer(t *testing.T) {
	dir := legacyRepo(t, "", map[string]string{"roles/qa.md": "# Role: qa\n"})
	var out bytes.Buffer
	err := migrate(&out, dir, false)
	if err == nil || !strings.Contains(err.Error(), "hub or a spoke") {
		t.Fatalf("err = %v", err)
	}
}

// --dry-run prints every step the live verb takes and writes none of them.
func TestMigrateDryRun(t *testing.T) {
	dir := legacyRepo(t, legacy1x, map[string]string{"roles/qa.md": "# Role: qa\n"})
	stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})
	before := headSubject(t, dir)

	var out bytes.Buffer
	if err := migrate(&out, dir, true); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"would move " + config.LegacyPointer + " -> " + config.Pointer,
		"would move roles/qa.md -> " + config.RolesDir + "/qa.md",
		"would remove the emptied roles/",
		"would write " + config.AgentsFile,
		"would write AGENTS.md",
		"would write CLAUDE.md",
		"would rewrite " + config.Pointer,
		"roles.implementer.identity: myorg-coder -> app:myorg-coder",
		"roles.coordinator: added",
		"dry run: nothing written",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("dry run did not print %q:\n%s", want, out.String())
		}
	}
	if exists(t, dir, config.Pointer) || !exists(t, dir, config.LegacyPointer) || !exists(t, dir, "roles/qa.md") || exists(t, dir, config.AgentsFile) {
		t.Error("the dry run wrote to disk")
	}
	for _, f := range rootEntryPoints {
		if exists(t, dir, f) {
			t.Errorf("the dry run wrote %s to disk", f)
		}
	}
	if headSubject(t, dir) != before {
		t.Error("the dry run committed")
	}
}

// The rewrite is a yaml round-trip, and a pointer is a file a human keeps:
// its comments, key order and blank lines come through the move.
func TestRewritePointerKeepsTheFilesShape(t *testing.T) {
	stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})
	doc, err := parsePointer([]byte(legacy1x))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := rewritePointer(doc); err != nil {
		t.Fatal(err)
	}
	out, err := encodePointer(doc)
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	for _, want := range []string{
		"# protocol version (SPEC.md §5)",
		"# Advisory role routing, read by whoever dispatches agents (see SPEC.md §5).",
		"hub: self\n\n#",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the rewrite lost %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, blankLineMarker) {
		t.Errorf("the blank-line marker leaked into the file:\n%s", got)
	}
	if !strings.HasSuffix(got, "\n") || strings.HasSuffix(got, "\n\n") {
		t.Errorf("the file does not end in exactly one newline: %q", got[len(got)-4:])
	}
	// The added row takes the style the table already uses — flow here,
	// block below — so the file still reads as the one its operator wrote.
	if !strings.Contains(got, "coordinator: {identity: ~}") {
		t.Errorf("the coordinator row does not match the table's flow style:\n%s", got)
	}
}

func TestRewritePointerAddsTheCoordinatorInTheTablesStyle(t *testing.T) {
	stubAccounts(t, map[string]string{"alice": "User"})
	doc, err := parsePointer([]byte("codecrew: \"1.0\"\nhub: self\nroles:\n  implementer:\n    identity: alice\n"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := rewritePointer(doc); err != nil {
		t.Fatal(err)
	}
	out, err := encodePointer(doc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "coordinator:\n    identity: ~") {
		t.Errorf("the coordinator row does not match the table's block style:\n%s", out)
	}
}

// Typing is by kind, and the two shapes a 1.0 table could not tell apart —
// an App slug and a username — are separated by asking GitHub.
func TestTypeIdentity(t *testing.T) {
	stubAccounts(t, map[string]string{
		"alice":             "User",
		"myorg-coder[bot]":  "Bot",
		"myorg-checky":      "User",
		"myorg-checky[bot]": "Bot",
		"acme":              "Organization",
	})
	for _, tc := range []struct {
		value, want, code string
	}{
		{value: "~", want: "~"},
		{value: "", want: "~"},
		{value: "app:myorg-coder", want: "app:myorg-coder"},
		{value: "user:alice", want: "user:alice"},
		{value: "team:myorg/crew", want: "team:myorg/crew"},
		{value: "myorg/crew", want: "team:myorg/crew"},
		{value: "alice", want: "user:alice"},
		{value: "myorg-coder", want: "app:myorg-coder"},
		{value: "myorg-coder[bot]", want: "app:myorg-coder"},
		{value: "myorg-checky", code: "IDENTITY_UNRESOLVED"}, // a user and an App answer to it
		{value: "nobody", code: "IDENTITY_UNRESOLVED"},
		{value: "acme", code: "IDENTITY_UNRESOLVED"}, // an organization holds no seat
		{value: "team:crew", code: "IDENTITY_UNRESOLVED"},
	} {
		got, err := typeIdentity("reviewer", tc.value)
		switch {
		case tc.code != "":
			if err == nil {
				t.Errorf("%q typed to %q, want refused[%s]", tc.value, got, tc.code)
				continue
			}
			if code := refusalCode(t, err); code != tc.code {
				t.Errorf("%q refused %s, want %s", tc.value, code, tc.code)
			}
			if !strings.Contains(err.Error(), "reviewer") {
				t.Errorf("%q: the refusal does not name the row: %v", tc.value, err)
			}
		case err != nil:
			t.Errorf("%q: %v", tc.value, err)
		case got != tc.want:
			t.Errorf("%q typed to %q, want %q", tc.value, got, tc.want)
		}
	}
}

// GitHub out of reach is its own condition after #259, not a value that
// could not be typed: migrate names it with the shared code and still
// writes nothing.
func TestMigrateNamesAnUnreachableGitHub(t *testing.T) {
	dir := legacyRepo(t, legacy1x, nil)
	prev := lookupAccount
	lookupAccount = func(string) (string, bool, error) {
		return "", false, errors.New("gh api: dial tcp 140.82.121.6:443: connect: network is unreachable")
	}
	t.Cleanup(func() { lookupAccount = prev })

	var out bytes.Buffer
	err := migrate(&out, dir, false)
	if code := refusalCode(t, err); code != "GH_UNREACHABLE" {
		t.Fatalf("code = %s, want GH_UNREACHABLE (%v)", code, err)
	}
	if exists(t, dir, config.Pointer) || !exists(t, dir, config.LegacyPointer) {
		t.Error("a refused migration moved files")
	}
}

// An identity that cannot be typed stops the whole migration before it
// writes: half a migration is worse than none.
func TestMigrateUnresolvedIdentityWritesNothing(t *testing.T) {
	dir := legacyRepo(t, legacy1x, map[string]string{"roles/qa.md": "# Role: qa\n"})
	stubAccounts(t, map[string]string{"alice": "User"}) // myorg-coder is nobody
	var out bytes.Buffer
	err := migrate(&out, dir, false)
	if code := refusalCode(t, err); code != "IDENTITY_UNRESOLVED" {
		t.Fatalf("code = %s, want IDENTITY_UNRESOLVED (%v)", code, err)
	}
	if exists(t, dir, config.Pointer) || !exists(t, dir, "roles/qa.md") {
		t.Error("a refused migration moved files")
	}
}

// A 1.x repository's cc: labels were all created implicitly — a colour
// GitHub generated, no description — so migrate restyles them to the
// protocol's defaults and creates the ones that never came up, which is
// the difference between it and init: migrate's promise is a repository
// indistinguishable from a fresh 2.0 scaffold, init reruns on a repository
// the project may have restyled on purpose (the operator's Decision on
// #283). The step runs after the commit, so nothing GitHub says can reach
// the move.
func TestMigrateRestylesTheLabels(t *testing.T) {
	dir := legacyRepo(t, legacy1x, map[string]string{"roles/qa.md": "# Role: qa\n"})
	stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})
	f := &labelFake{defined: []tracker.Label{implicit("bug"), implicit(tracker.LabelTask)}}
	stubLabelTarget(t, f, "o/r", nil)

	var out bytes.Buffer
	if err := migrate(&out, dir, false); err != nil {
		t.Fatal(err)
	}
	if got := names(f.restyled); !slices.Equal(got, []string{tracker.LabelTask}) {
		t.Errorf("restyled %v, want the one label the 1.x repo had", got)
	}
	if got := names(f.created); !slices.Equal(got, []string{tracker.LabelMilestone, tracker.LabelNeedsDecision}) {
		t.Errorf("created %v", got)
	}
	if _, restyled := tracker.FindLabel(f.restyled, "bug"); restyled {
		t.Error("migrate restyled a label that is not the protocol's")
	}
	for _, want := range []string{
		"restyled label cc:task (#ededed -> #92edff)",
		"created label cc:milestone (#01d4ff)",
		"created label cc:needs-decision (#f0aeff)",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("migrate did not report %q:\n%s", want, out.String())
		}
	}
	// After the commit, and before the guidance that follows it.
	commit := strings.Index(out.String(), "committed ")
	label := strings.Index(out.String(), "restyled label")
	next := strings.Index(out.String(), "next: read it")
	if commit < 0 || !(commit < label && label < next) {
		t.Errorf("the label step is not between the commit and the guidance:\n%s", out.String())
	}
}

// The dry run lists the label steps beside the file ones and writes
// neither. A GitHub that will not answer is a note in both modes and never
// a refusal — the migration is a local move, and no remote failure may
// stand in its way.
func TestMigrateDryRunListsTheLabelsAndWritesNothing(t *testing.T) {
	dir := legacyRepo(t, legacy1x, map[string]string{"roles/qa.md": "# Role: qa\n"})
	stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})
	f := &labelFake{defined: []tracker.Label{implicit(tracker.LabelTask)}}
	stubLabelTarget(t, f, "o/r", nil)

	var out bytes.Buffer
	if err := migrate(&out, dir, true); err != nil {
		t.Fatal(err)
	}
	if len(f.created) != 0 || len(f.restyled) != 0 {
		t.Errorf("the dry run wrote: created %v restyled %v", names(f.created), names(f.restyled))
	}
	for _, want := range []string{
		"would restyle label cc:task (#ededed -> #92edff)",
		"would create label cc:milestone (#01d4ff)",
		"would create label cc:needs-decision (#f0aeff)",
		"dry run: nothing written",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("the dry run did not print %q:\n%s", want, out.String())
		}
	}

	// And with GitHub refusing, in both modes: a note, no refusal, and the
	// migration still committed.
	for _, dryRun := range []bool{true, false} {
		dir := legacyRepo(t, legacy1x, map[string]string{"roles/qa.md": "# Role: qa\n"})
		stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})
		stubLabelTarget(t, &labelFake{readErr: errors.New("403")}, "o/r", nil)
		var out bytes.Buffer
		if err := migrate(&out, dir, dryRun); err != nil {
			t.Fatalf("dryRun=%v: %v", dryRun, err)
		}
		if want := "note: could not read o/r's labels (403)"; !strings.Contains(out.String(), want) {
			t.Errorf("dryRun=%v: output missing %q:\n%s", dryRun, want, out.String())
		}
		if committed := strings.Contains(out.String(), "committed "); committed == dryRun {
			t.Errorf("dryRun=%v: committed = %v", dryRun, committed)
		}
	}
}

// A rerun on a repository already on 2.0 is the documented recovery from a
// label step that could not reach GitHub, so it has to actually do the
// labels: migrate is idempotent in what it moves, not in what it does.
// Without this the first run's transient `gh` failure left the repository
// on 2.0 wearing GitHub's grey with nothing in the tool that would ever
// fix it (checky's finding 1 on PR #291).
func TestMigrateRerunDoesTheLabelsAndNothingElse(t *testing.T) {
	for _, tc := range []struct {
		name    string
		fake    *labelFake
		dryRun  bool
		created []string
		lines   []string
		absent  []string
	}{
		{
			// The state finding 1 describes: the move landed, the labels
			// did not.
			name:    "the first run could not reach GitHub",
			fake:    &labelFake{},
			created: []string{tracker.LabelMilestone, tracker.LabelTask, tracker.LabelNeedsDecision},
			lines:   []string{"already on the protocol 2.0 layout", "created label cc:milestone (#01d4ff)", "created label cc:needs-decision (#f0aeff)"},
			absent:  []string{"labels already at the protocol defaults", "would create"},
		},
		{
			// Nothing to do, and it says so: on this path silence is
			// indistinguishable from the step never having run.
			name:   "a second rerun, with the labels already right",
			fake:   &labelFake{defined: defaults(tracker.LabelMilestone, tracker.LabelTask, tracker.LabelNeedsDecision)},
			lines:  []string{"labels already at the protocol defaults"},
			absent: []string{"created label", "restyled label"},
		},
		{
			// --dry-run previews the recovery and writes none of it.
			name:   "the dry run previews it",
			fake:   &labelFake{defined: []tracker.Label{implicit(tracker.LabelTask)}},
			dryRun: true,
			lines:  []string{"would restyle label cc:task (#ededed -> #92edff)", "would create label cc:milestone", "dry run: nothing written"},
			absent: []string{"labels already at the protocol defaults"},
		},
		{
			// And a GitHub that still will not answer is still a note.
			name:   "GitHub is still unreachable",
			fake:   &labelFake{readErr: errors.New("403")},
			lines:  []string{"note: could not read o/r's labels (403)"},
			absent: []string{"labels already at the protocol defaults", "created label"},
		},
	} {
		dir := legacyRepo(t, "", map[string]string{config.Pointer: "codecrew: \"2.0\"\nhub: self\n"})
		stubLabelTarget(t, tc.fake, "o/r", nil)
		before := headSubject(t, dir)

		var out bytes.Buffer
		if err := migrate(&out, dir, tc.dryRun); err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if got := names(tc.fake.created); !slices.Equal(got, tc.created) {
			t.Errorf("%s: created %v, want %v", tc.name, got, tc.created)
		}
		for _, want := range tc.lines {
			if !strings.Contains(out.String(), want) {
				t.Errorf("%s: output missing %q:\n%s", tc.name, want, out.String())
			}
		}
		for _, never := range tc.absent {
			if strings.Contains(out.String(), never) {
				t.Errorf("%s: output must not contain %q:\n%s", tc.name, never, out.String())
			}
		}
		// Nothing else happens on this path: no commit, no file moved or
		// written, and the 1.x layout was never there to remove.
		if headSubject(t, dir) != before {
			t.Errorf("%s: the rerun committed", tc.name)
		}
		if st, _ := git(dir, "status", "--porcelain"); st != "" {
			t.Errorf("%s: the rerun wrote to disk: %q", tc.name, st)
		}
		for _, never := range []string{"moved ", "would move ", "wrote ", "rewrote ", "committed "} {
			if strings.Contains(out.String(), never) {
				t.Errorf("%s: the rerun reported %q:\n%s", tc.name, never, out.String())
			}
		}
	}
}

// The move is on disk before the commit is attempted, so the labels run on
// the two paths that never reach it: a detached HEAD, and a commit git
// refused. Otherwise the operator is told to finish the commit by hand and
// has no reason to suspect the labels were skipped (checky's finding 2).
func TestMigrateDoesTheLabelsWhenTheCommitDoesNotHappen(t *testing.T) {
	dir := legacyRepo(t, legacy1x, map[string]string{"roles/qa.md": "# Role: qa\n"})
	stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})
	f := &labelFake{defined: []tracker.Label{implicit(tracker.LabelTask)}}
	stubLabelTarget(t, f, "o/r", nil)
	// Detach: the moves still apply, the commit cannot.
	if _, err := git(dir, "checkout", "-q", "--detach"); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := migrate(&out, dir, false); err != nil {
		t.Fatal(err)
	}
	if got := names(f.restyled); !slices.Equal(got, []string{tracker.LabelTask}) {
		t.Errorf("restyled %v", got)
	}
	if got := names(f.created); !slices.Equal(got, []string{tracker.LabelMilestone, tracker.LabelNeedsDecision}) {
		t.Errorf("created %v", got)
	}
	if !exists(t, dir, config.Pointer) {
		t.Error("the moves did not apply")
	}
	// The act asked of a human is still the last thing on screen.
	label := strings.Index(out.String(), "restyled label")
	note := strings.Index(out.String(), "note: HEAD is detached")
	if label < 0 || note < 0 || label > note {
		t.Errorf("the labels must be reported before the detached-HEAD note:\n%s", out.String())
	}
}

// The other commit-less path: `git commit` refused. The moves are on disk
// exactly as in the detached case, so the labels run there too — and until
// this test the claim rested on the comment above and on a live run, with
// the call itself unguarded (checky's note 2 on PR #291). The identity is
// unset after legacyRepo has made its commit, the way
// TestCommitScaffoldFailureIsANote induces the same failure.
func TestMigrateDoesTheLabelsWhenTheCommitIsRefused(t *testing.T) {
	dir := legacyRepo(t, legacy1x, map[string]string{"roles/qa.md": "# Role: qa\n"})
	stubAccounts(t, map[string]string{"myorg-coder[bot]": "Bot", "alice": "User"})
	f := &labelFake{defined: []tracker.Label{implicit(tracker.LabelTask)}}
	stubLabelTarget(t, f, "o/r", nil)
	before := headSubject(t, dir)
	git(dir, "config", "--unset", "user.email")
	git(dir, "config", "--unset", "user.name")
	for _, v := range []string{"GIT_AUTHOR_NAME", "GIT_COMMITTER_NAME", "GIT_AUTHOR_EMAIL", "GIT_COMMITTER_EMAIL", "EMAIL"} {
		t.Setenv(v, "")
	}
	t.Setenv("HOME", t.TempDir()) // no global identity either

	var out bytes.Buffer
	if err := migrate(&out, dir, false); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "note: could not commit the migration") {
		t.Fatalf("the commit did not fail:\n%s", out.String())
	}
	if headSubject(t, dir) != before {
		t.Error("a commit was made without an identity")
	}
	// The moves landed, so the labels did too.
	if !exists(t, dir, config.Pointer) {
		t.Error("the moves did not apply")
	}
	if got := names(f.restyled); !slices.Equal(got, []string{tracker.LabelTask}) {
		t.Errorf("restyled %v", got)
	}
	if got := names(f.created); !slices.Equal(got, []string{tracker.LabelMilestone, tracker.LabelNeedsDecision}) {
		t.Errorf("created %v", got)
	}
	// Here the receipts follow pathspecCommit's own note, which named the
	// command to run — the asymmetry with the detached path that the
	// Decision on #283 records.
	note := strings.Index(out.String(), "note: could not commit the migration")
	label := strings.Index(out.String(), "restyled label")
	if label < note {
		t.Errorf("the receipts must follow the commit note:\n%s", out.String())
	}
}
