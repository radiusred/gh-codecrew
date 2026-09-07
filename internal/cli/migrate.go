package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/gh"
	"gopkg.in/yaml.v3"
)

// migrateSubject is the subject of the one commit migrate makes. It is a
// conventional-commit chore for the same reason init's is: the move is
// housekeeping the protocol asks for, not a change to the project.
const migrateSubject = "chore: migrate codecrew to the 2.0 layout"

// migrateCmd is the one-shot move from the protocol 1.x layout to 2.0
// (SPEC §6, M13-R2). It reads no pointer — the repos it exists for are
// exactly the ones every other verb refuses — so, like init, it is exempt
// from the protocol check, and like init it refuses a subdirectory: the
// layout is judged at the repository root.
func migrateCmd(w io.Writer, args []string) error {
	flags := flag.NewFlagSet("migrate", flag.ContinueOnError)
	dryRun := flags.Bool("dry-run", false, "print every step and write nothing")
	if err := flags.Parse(args); err != nil {
		return err
	}
	root := repoRoot(".")
	if root == "" {
		return fmt.Errorf("run migrate inside a git repository: the move is a git mv, so the history follows the files")
	}
	if cwd, err := os.Getwd(); err == nil && !sameDir(cwd, root) {
		return fmt.Errorf("run migrate at the repository root (%s), not in a subdirectory", root)
	}
	return migrate(w, root, *dryRun)
}

// move is one file the migration relocates. tracked reports whether git
// knew the source: an untracked 1.x file is renamed on the filesystem
// instead, and only its destination belongs in the commit's pathspec.
type move struct {
	from, to string
	tracked  bool
}

// migrate does the whole move for the repository rooted at root. Every
// refusal is raised before anything is written, so a refused migration
// leaves the repo exactly as it found it — and so --dry-run, which stops
// at the same point, is a true preview.
func migrate(w io.Writer, root string, dryRun bool) error {
	legacy := config.LegacyLayout(root)
	_, err := os.Stat(filepath.Join(root, filepath.FromSlash(config.Pointer)))
	current := err == nil
	switch {
	case len(legacy) > 0 && current:
		return refuse("BOTH_LAYOUTS", "%s holds both layouts — the protocol 1.x files (%s) and %s — and migrate will not choose between them: keep whichever the project uses, remove the other, then rerun",
			root, strings.Join(legacy, ", "), config.Pointer)
	case current:
		// The files are where they belong; the labels may not be. A rerun
		// is the documented recovery from a label step that could not
		// reach GitHub, and it only works if this path does it (checky's
		// finding 1 on PR #291), so migrate is idempotent in what it
		// moves rather than in what it does.
		fmt.Fprintf(w, "already on the protocol %s layout (%s) — nothing to move\n", protocolVersion, config.Pointer)
		migrateLabelsRerun(w, dryRun)
		if dryRun {
			fmt.Fprintln(w, "dry run: nothing written")
		}
		return nil
	case len(legacy) == 0:
		return fmt.Errorf("%s holds no %s and no protocol 1.x layout (not a CodeCrew repo?)", root, config.Pointer)
	case !slices.Contains(legacy, config.LegacyPointer):
		return fmt.Errorf("%s holds 1.x role contracts (%s) but no %s: nothing there says whether the repo is a hub or a spoke, and migrate does not guess — restore the pointer, or move the files by hand",
			root, strings.Join(legacy, ", "), config.LegacyPointer)
	}

	src, err := os.ReadFile(filepath.Join(root, config.LegacyPointer))
	if err != nil {
		return err
	}
	doc, err := parsePointer(src)
	if err != nil {
		return err
	}
	// Cheapest and most local first, so a refusal costs no API call: the
	// protocol major, then the directory, then the identities.
	if err := checkMigratable(doc); err != nil {
		return err
	}
	if err := checkSpokeRouting(doc); err != nil {
		return err
	}
	roleMoves, err := rolesMoves(root)
	if err != nil {
		return err
	}
	moves := append([]move{{from: config.LegacyPointer, to: config.Pointer}}, roleMoves...)
	if err := checkDestinations(root, moves); err != nil {
		return err
	}
	// The 2.0 entry point is a file of CodeCrew's that the 1.x layout had
	// nowhere to put: its instructions lived in the root AGENTS.md, which
	// belongs to the project (M13-R3). Migrating the tree without it would
	// leave the repo on the 2.0 layout with no 2.0 entry point, so the
	// scaffold init writes goes in when it is absent — and so do the root
	// entry points, from the same constants init writes them from. Only a
	// file that is already there is the project's; that one is reported,
	// never rewritten.
	written, kept := entryPointPlan(root)
	changes, err := rewritePointer(doc)
	if err != nil {
		return err
	}
	pointer, err := encodePointer(doc)
	if err != nil {
		return err
	}

	emptied := len(roleMoves) > 0

	moved, removed, wrote, rewrote := "moved", "removed", "wrote", "rewrote"
	if dryRun {
		moved, removed, wrote, rewrote = "would move", "would remove", "would write", "would rewrite"
	}
	for _, m := range moves {
		fmt.Fprintf(w, "%s %s -> %s\n", moved, m.from, m.to)
	}
	if emptied {
		fmt.Fprintf(w, "%s the emptied %s/\n", removed, config.LegacyRolesDir)
	}
	for _, f := range written {
		fmt.Fprintf(w, "%s %s\n", wrote, f.path)
	}
	fmt.Fprintf(w, "%s %s:\n", rewrote, config.Pointer)
	for _, c := range changes {
		fmt.Fprintf(w, "  %s\n", c)
	}

	if dryRun {
		fmt.Fprintf(w, "would commit %d paths on %s: %q\n", commitPathCount(moves)+len(written), currentBranch(root), migrateSubject)
		migrateLabels(w, true)
		fmt.Fprintln(w, "dry run: nothing written")
		reportEntryPoint(w, kept)
		return nil
	}

	branch, err := git(root, "symbolic-ref", "--short", "HEAD")
	if err != nil {
		// A commit here would strand the migration on no branch. The files
		// still move — the operator asked for that — and the commit is
		// theirs to make once they are on a branch.
		branch = ""
	}
	if err := applyMoves(root, moves); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(config.Pointer)), pointer, 0o644); err != nil {
		return err
	}
	for _, f := range written {
		if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(f.path)), []byte(f.content), 0o644); err != nil {
			return err
		}
	}
	if emptied {
		// os.Remove refuses a directory that is not empty, which is exactly
		// the guard wanted: anything left behind keeps the directory.
		if err := os.Remove(filepath.Join(root, config.LegacyRolesDir)); err != nil {
			fmt.Fprintf(w, "note: %s/ is not empty and was kept (%v)\n", config.LegacyRolesDir, err)
		}
	}

	stage, paths := commitPaths(moves, written)
	if branch == "" {
		// The moves are already on disk, so the labels are done here too —
		// one rule, "once the move reached disk the labels run", rather
		// than a path on which the operator is told to finish the commit
		// by hand and has no reason to suspect the labels were skipped
		// (checky's finding 2). The note stays last: it is the act asked
		// of a human.
		migrateLabels(w, false)
		fmt.Fprintf(w, "note: HEAD is detached — the files are moved but not committed; switch to a branch and run: %s\n", commitByHand(migrateSubject, paths))
		return nil
	}
	sha, ok := pathspecCommit(w, root, migrateSubject, "the migration", stage, paths)
	if !ok {
		// Same rule: the move is on disk, so the labels run — after
		// pathspecCommit's own note, which named the command to run.
		migrateLabels(w, false)
		return nil
	}
	fmt.Fprintf(w, "committed %s on %s: %q — the migration only; your other changes are as they were\n", sha, branch, migrateSubject)
	// Then the one part of the migration that is not a file: the three
	// labels, brought to the protocol's colour and description whatever a
	// 1.x repo had them wearing (§4). After the commit, so a GitHub
	// failure is a note and the move stands (M14-R5, #267).
	migrateLabels(w, false)
	fmt.Fprintf(w, "next: read it (git show %s), then push and open a pull request — migrate never pushes\n", sha)
	// Last, so the one thing needing a human is the last thing on screen.
	reportEntryPoint(w, kept)
	return nil
}

// scaffolded is a file migrate writes outright rather than moves: the path
// it goes to, and the bytes init would write there.
type scaffolded struct{ path, content string }

// entryPointPlan decides what the migration does about the entry points
// before it writes anything, so the dry run previews the state the live run
// creates rather than the one both start from. write holds every file that
// is absent — .codecrew/AGENTS.md, which is CodeCrew's own, and each root
// entry point, from the constants init writes them from. kept holds the
// root files that are already there and would still not reach
// .codecrew/AGENTS.md once write is on disk: those are the project's prose,
// and the only thing left for a human (#301). A kept CLAUDE.md counts as
// reaching through an AGENTS.md this migration is about to write, exactly
// as it would on a rerun once that file exists.
func entryPointPlan(root string) (write []scaffolded, kept []string) {
	absent := func(rel string) bool {
		_, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel)))
		return err != nil
	}
	if absent(config.AgentsFile) {
		write = append(write, scaffolded{config.AgentsFile, agentsScaffold})
	}
	pending := map[string]bool{}
	for _, f := range rootEntryPoints {
		if absent(f) {
			write = append(write, scaffolded{f, rootEntryPointScaffolds[f]})
			pending[f] = true
		}
	}
	for _, f := range rootEntryPoints {
		if !pending[f] && !reachesInstructions(root, f, pending) {
			kept = append(kept, f)
		}
	}
	return write, kept
}

// reportEntryPoint names the kept root entry points that do not reach
// .codecrew/AGENTS.md. A 1.x repo's root AGENTS.md holds the instructions
// themselves and names the paths that just moved, so after a migration it
// almost always needs the two lines — but it is the project's file, and
// migrate rewrites nobody's prose. An absent one is not the operator's
// work at all: migrate writes it, as init does, and it never reaches here.
func reportEntryPoint(w io.Writer, kept []string) {
	if len(kept) == 0 {
		return
	}
	entryPointAction(w, kept)
}

// currentBranch is the branch a commit would land on, or "HEAD" when it
// cannot be named — the dry run reports where, without needing the commit
// to be possible.
func currentBranch(dir string) string {
	if branch, err := git(dir, "symbolic-ref", "--short", "HEAD"); err == nil {
		return branch
	}
	return "HEAD"
}

// commitPaths splits the migration into what has to be staged (every
// destination and every file written outright) and the whole pathspec the
// commit is limited to (the same, plus the sources git already knows, so
// their removal is recorded in the same commit).
func commitPaths(moves []move, written []scaffolded) (stage, paths []string) {
	for _, m := range moves {
		stage = append(stage, m.to)
		paths = append(paths, m.to)
		if m.tracked {
			paths = append(paths, m.from)
		}
	}
	for _, f := range written {
		stage = append(stage, f.path)
		paths = append(paths, f.path)
	}
	return stage, paths
}

// commitPathCount is what the pathspec will hold once the moves are made.
// The dry run has not run git mv, so it counts what a tracked repo gives:
// both ends of every move.
func commitPathCount(moves []move) int { return 2 * len(moves) }

// checkDestinations refuses before anything is written when a 2.0 file
// already sits where a 1.x one would move. It is the same condition the
// pointer check catches one level up — the repo carries both layouts — and
// the answer is the same: migrate will not overwrite the newer file to
// reach the older one, and it will not leave the tracked source's deletion
// outside its own commit (checky's finding on PR #280).
func checkDestinations(root string, moves []move) error {
	var clashes []string
	for _, m := range moves {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(m.to))); err == nil {
			clashes = append(clashes, fmt.Sprintf("%s, where %s would go", m.to, m.from))
		}
	}
	if len(clashes) == 0 {
		return nil
	}
	return refuse("BOTH_LAYOUTS", "%s carries both layouts: it already holds %s — keep whichever the project uses, remove the other, then rerun",
		root, strings.Join(clashes, "; "))
}

// applyMoves performs each rename with git mv, so the history follows the
// file. The one case that falls back to the filesystem is a source git does
// not know — an uncommitted 1.x scaffold — because naming an untracked path
// in the commit's pathspec would fail. Every other git mv failure is
// returned: git refusing to move a file it tracks is the repository's
// answer, and renaming underneath it would strand the source's deletion
// outside the migration commit.
func applyMoves(root string, moves []move) error {
	for i := range moves {
		to := filepath.Join(root, filepath.FromSlash(moves[i].to))
		if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
			return err
		}
		_, mvErr := git(root, "mv", "--", moves[i].from, moves[i].to)
		if mvErr == nil {
			moves[i].tracked = true
			continue
		}
		if isTracked(root, moves[i].from) {
			return fmt.Errorf("moving %s to %s: %w", moves[i].from, moves[i].to, mvErr)
		}
		if err := os.Rename(filepath.Join(root, filepath.FromSlash(moves[i].from)), to); err != nil {
			return fmt.Errorf("moving %s to %s: %w", moves[i].from, moves[i].to, err)
		}
	}
	return nil
}

// isTracked reports whether git holds an index entry for path.
func isTracked(root, path string) bool {
	_, err := git(root, "ls-files", "--error-unmatch", "--", path)
	return err == nil
}

// migratableNames are the only file names migrate takes out of a 1.x
// roles/ directory: the five role contracts and their local extensions
// (SPEC §7). Everything else in there is somebody's — the project's, or a
// leftover nobody has claimed — and migrate refuses rather than deciding.
func migratableNames() []string {
	names := make([]string, 0, 2*len(config.RoleNames))
	for _, role := range config.RoleNames {
		names = append(names, role+".md", role+localSuffix)
	}
	return names
}

// rolesMoves lists what comes out of root/roles/. The directory is read at
// all only when it holds at least one recognised name: a project's own
// roles/ — Ansible's, say — is not CodeCrew's to touch, and that collision
// is the whole reason for the 2.0 layout. Once the directory is CodeCrew's,
// every entry in it must be recognised, or the verb refuses
// FOREIGN_ROLES_DIR rather than guessing which of them it owns (M13-R2).
func rolesMoves(root string) ([]move, error) {
	entries, err := os.ReadDir(filepath.Join(root, config.LegacyRolesDir))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	known := migratableNames()
	var ours, foreign []string
	for _, e := range entries {
		if e.Type().IsRegular() && slices.Contains(known, e.Name()) {
			ours = append(ours, e.Name())
			continue
		}
		foreign = append(foreign, e.Name())
	}
	if len(ours) == 0 {
		return nil, nil // not CodeCrew's directory; leave it where it is
	}
	if len(foreign) > 0 {
		return nil, refuse("FOREIGN_ROLES_DIR", "%s/ holds CodeCrew's files (%s) and entries it does not recognise (%s): migrate moves the five role contracts and their %s extensions and nothing else, and will not guess about the rest — move or remove them, then rerun",
			config.LegacyRolesDir, strings.Join(ours, ", "), strings.Join(foreign, ", "), localSuffix)
	}
	moves := make([]move, 0, len(ours))
	for _, name := range ours {
		moves = append(moves, move{
			from: config.LegacyRolesDir + "/" + name,
			to:   path.Join(config.RolesDir, name),
		})
	}
	return moves, nil
}

// blankLineMarker stands in for a blank line across the yaml round-trip
// that rewrites the pointer. The round-trip keeps comments, which is why
// it is used at all, but drops blank lines — and in a file a human
// maintains those carry as much of its shape as the comments do.
const blankLineMarker = "#codecrew-blank-line"

func hideBlankLines(s string) string {
	lines := strings.Split(s, "\n")
	for i := 0; i < len(lines)-1; i++ { // the empty tail after a final newline is not a line
		if strings.TrimSpace(lines[i]) == "" {
			lines[i] = blankLineMarker
		}
	}
	return strings.Join(lines, "\n")
}

func showBlankLines(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		if strings.TrimSpace(l) == blankLineMarker {
			lines[i] = ""
		}
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
}

// parsePointer decodes a 1.x pointer into the node tree the rewrite edits
// in place — comments, key order and flow styles included, so the file the
// operator reads after the move is the file they wrote before it.
func parsePointer(src []byte) (*yaml.Node, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(hideBlankLines(string(src))), &doc); err != nil {
		return nil, fmt.Errorf("%s: %w", config.LegacyPointer, err)
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%s: not a mapping — migrate rewrites a pointer file, and this is not one", config.LegacyPointer)
	}
	return doc.Content[0], nil
}

func encodePointer(root *yaml.Node) ([]byte, error) {
	var b strings.Builder
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	if err := enc.Encode(root); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	return []byte(showBlankLines(b.String())), nil
}

// mapValue is the value node a mapping holds for key, or nil.
func mapValue(m *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}
	return nil
}

// checkMigratable refuses a pointer migrate has nothing to say about.
// Protocol 1.0 is the only major it moves forward: 0.x predates the
// conventions this move assumes, and anything above 1 is not a 1.x repo at
// all, whatever the files beside it look like.
func checkMigratable(root *yaml.Node) error {
	v := mapValue(root, "codecrew")
	if v == nil || v.Tag == "!!null" || v.Value == "" {
		return nil // no version stated: the layout says 1.x, and that is enough
	}
	major, _, _ := strings.Cut(v.Value, ".")
	if n, err := strconv.Atoi(major); err == nil && n == 1 {
		return nil
	}
	return refuse("MIGRATION_UNSUPPORTED", "%s says protocol %q; migrate moves a protocol 1.x repo to %s and nothing else — a pointer below 1.0 predates the conventions it assumes (SPEC §5)",
		config.LegacyPointer, v.Value, protocolVersion)
}

// checkSpokeRouting refuses a 1.x spoke pointer that carries a routing
// table. Protocol 1.0 allowed the shape and 2.0 does not (M13-R5, #259),
// so migrating it faithfully would write a .codecrew/config.yml that this
// same binary refuses SPOKE_ROUTING on the very next verb — a one-shot
// move that has not landed. The block is not dropped on the operator's
// behalf: it is routing they wrote, and where it belongs is theirs to say
// (checky's second-round finding on PR #280). The rows are named in sorted
// order, as config.Parse names them, so the refusal does not move between
// runs.
func checkSpokeRouting(root *yaml.Node) error {
	hub := mapValue(root, "hub")
	if hub == nil || hub.Value == "" || hub.Value == "self" {
		return nil
	}
	roles := mapValue(root, "roles")
	if roles == nil || roles.Kind != yaml.MappingNode || len(roles.Content) == 0 {
		return nil // an empty table is not a table, exactly as Parse reads it
	}
	names := make([]string, 0, len(roles.Content)/2)
	for i := 0; i+1 < len(roles.Content); i += 2 {
		names = append(names, roles.Content[i].Value)
	}
	slices.Sort(names)
	return refuse("SPOKE_ROUTING", "%s names the hub %s and carries a roles: block (%s), which protocol 1.0 allowed and 2.0 does not — the hub carries the one routing table for the project (SPEC §5), so migrating this table forward would write a %s that every verb then refuses. Migrate will not drop routing you wrote: move these rows into %s's %s, or delete them here, then rerun",
		config.LegacyPointer, hub.Value, strings.Join(names, ", "), config.Pointer, hub.Value, config.Pointer)
}

// rewritePointer edits the pointer in place into its 2.0 form and returns
// what it changed, for the receipt: the protocol version, every identity
// typed per M13-R4, and the coordinator row a 1.0 table could leave out
// (its absence used to be a special case in the code; the row makes it a
// statement).
func rewritePointer(root *yaml.Node) ([]string, error) {
	var changes []string

	if v := mapValue(root, "codecrew"); v == nil {
		root.Content = append([]*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: "!!str", Value: "codecrew"},
			{Kind: yaml.ScalarNode, Tag: "!!str", Style: yaml.DoubleQuotedStyle, Value: protocolVersion},
		}, root.Content...)
		changes = append(changes, fmt.Sprintf("codecrew: added, %q (the 1.x layout stated no version)", protocolVersion))
	} else {
		was := v.Value
		v.Kind, v.Tag, v.Style, v.Value = yaml.ScalarNode, "!!str", yaml.DoubleQuotedStyle, protocolVersion
		if was == "" {
			changes = append(changes, fmt.Sprintf("codecrew: %q (it stated no version)", protocolVersion))
		} else {
			changes = append(changes, fmt.Sprintf("codecrew: %q -> %q", was, protocolVersion))
		}
	}

	roles := mapValue(root, "roles")
	if roles == nil || roles.Kind != yaml.MappingNode {
		return changes, nil // a spoke carries no table; there is nothing to type
	}
	for i := 0; i+1 < len(roles.Content); i += 2 {
		name, row := roles.Content[i].Value, roles.Content[i+1]
		if row.Kind != yaml.MappingNode {
			continue
		}
		v := mapValue(row, "identity")
		if v == nil || v.Tag == "!!null" {
			continue // absent or ~ — the operator, and already in the grammar
		}
		if v.Kind != yaml.ScalarNode {
			return nil, refuse("IDENTITY_UNRESOLVED", "roles.%s.identity is not a single value, so migrate cannot type it — write one of `~`, `app:<slug>`, `user:<login>` or `team:<org>/<slug>` by hand, then rerun", name)
		}
		typed, err := typeIdentity(name, v.Value)
		if err != nil {
			return nil, err
		}
		if typed == v.Value {
			continue
		}
		changes = append(changes, fmt.Sprintf("roles.%s.identity: %s -> %s", name, v.Value, typed))
		v.Value, v.Style = typed, 0
	}
	if mapValue(roles, "coordinator") == nil {
		// Written in the style the table already uses, block or flow: the
		// row appears in a file its operator maintains, and a migration
		// that reformats what it touches is a migration nobody reads.
		identity := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "~"}
		style := yaml.FlowStyle
		if len(roles.Content) > 1 && roles.Content[1].Style != yaml.FlowStyle {
			style = 0
			identity.LineComment = "the operator, until this seat is routed (SPEC §7)"
		}
		roles.Content = append(roles.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "coordinator"},
			&yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Style: style, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: "!!str", Value: "identity"},
				identity,
			}},
		)
		changes = append(changes, "roles.coordinator: added, identity ~ (the seat that dispatches the other four, SPEC §7)")
	}
	return changes, nil
}

// typeIdentity turns a 1.0 table's identity value into the 2.0 grammar. A
// value that is already in the grammar passes through; `~` and an empty
// value are the operator and stay as they are; a value carrying a slash is
// the only 1.0 shape that named its own kind, a team. Anything else is a
// bare login, and only GitHub knows whether it is an App or a human — the
// ambiguity M13-R4 exists to end, and the reason this is resolved rather
// than guessed.
func typeIdentity(role, value string) (string, error) {
	value = strings.TrimSpace(value)
	if id := config.ParseIdentity(value); id.Kind != config.KindUntyped {
		return id.String(), nil
	}
	if kind, _, ok := strings.Cut(value, ":"); ok && (kind == "app" || kind == "user" || kind == "team") {
		return "", refuse("IDENTITY_UNRESOLVED", "roles.%s.identity: %q carries the %q prefix but is not a value of that kind, and migrate will not repair it — write `app:<slug>`, `user:<login>` or `team:<org>/<slug>` by hand, then rerun", role, value, kind)
	}
	if strings.Contains(value, "/") {
		return "team:" + value, nil
	}
	return resolveIdentity(role, value)
}

// lookupAccount reports what GitHub holds at users/<login> — "User" for a
// human, "Bot" for an App's account, "Organization" — and whether anything
// is there at all. A func var so tests stand in for the API, the
// teamMembers pattern.
var lookupAccount = func(login string) (accountType string, found bool, err error) {
	var acct struct {
		Type string `json:"type"`
	}
	if err := gh.JSON(&acct, "api", "users/"+url.PathEscape(login)); err != nil {
		if strings.Contains(err.Error(), "HTTP 404") || strings.Contains(err.Error(), "Not Found") {
			return "", false, nil
		}
		return "", false, err
	}
	return acct.Type, true, nil
}

// resolveIdentity asks GitHub what a bare 1.0 login is. A GitHub App's
// account login carries a `[bot]` suffix the routing table never did, so a
// bare App slug is a 404 at users/<slug> and is found one probe later at
// users/<slug>[bot] — and when both a human and an App answer to the same
// name the value means two things, which is not migrate's to settle.
func resolveIdentity(role, value string) (string, error) {
	ask := func(login string) (string, bool, error) {
		kind, found, err := lookupAccount(login)
		if err == nil {
			return kind, found, nil
		}
		// GitHub out of reach is its own condition, named by its own code
		// (SPEC §6) — never folded into "this value cannot be typed",
		// which says the API answered and the answer did not settle it.
		if unreachable := unreachable(err); unreachable != nil {
			return "", false, unreachable
		}
		return "", false, refuse("IDENTITY_UNRESOLVED", "roles.%s.identity: GitHub would not say what %q is (%v) — the value has to be typed before the table can be read, so fix the access or write `app:<slug>`, `user:<login>` or `team:<org>/<slug>` by hand, then rerun", role, value, err)
	}
	slug := strings.TrimSuffix(value, "[bot]")
	kind, found, err := ask(value)
	if err != nil {
		return "", err
	}
	if found && kind == "Bot" {
		return "app:" + slug, nil
	}
	botKind, botFound := "", false
	if value == slug { // not already suffixed: the App's account may still be there
		if botKind, botFound, err = ask(value + "[bot]"); err != nil {
			return "", err
		}
	}
	switch {
	case found && kind == "User" && botFound && botKind == "Bot":
		return "", refuse("IDENTITY_UNRESOLVED", "roles.%s.identity: %q is both a GitHub user and the App %s[bot], and migrate will not choose — write `user:%s` or `app:%s` by hand, then rerun", role, value, value, value, value)
	case found && kind == "User":
		return "user:" + value, nil
	case botFound && botKind == "Bot":
		return "app:" + value, nil
	case found:
		return "", refuse("IDENTITY_UNRESOLVED", "roles.%s.identity: %q is a GitHub %s, which holds no seat — a routing identity is `~`, `app:<slug>`, `user:<login>` or `team:<org>/<slug>` (SPEC §5); write one by hand, then rerun", role, value, strings.ToLower(kind))
	default:
		return "", refuse("IDENTITY_UNRESOLVED", "roles.%s.identity: GitHub knows no user %q and no App %s[bot], so migrate cannot say what kind of principal it names — write `~`, `app:<slug>`, `user:<login>` or `team:<org>/<slug>` by hand, then rerun", role, value, value)
	}
}
