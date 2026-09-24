package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/config"
)

// syncBranch is where roles sync commits when it is run on the default
// branch: the light path is always a pull request (SPEC §4), and a commit
// on the default branch invites the direct push the contracts forbid.
const syncBranch = "codecrew-roles-sync"

// releasedContract is one row of contractHistory: the SHA-256 of a role
// contract's text as a tagged release shipped it.
type releasedContract struct {
	Release, Role, SHA256 string
}

// contractState is what a local contract is, measured against the
// contracts this binary embeds and every release's text before them.
type contractState int

const (
	contractCurrent contractState = iota // the embedded text: nothing to do
	contractAbsent                       // no local file: roles sync writes it
	contractRelease                      // an earlier release's text: roles sync replaces it
	contractForked                       // anything else: the project's fork, never touched
)

// contractStatus is one role's classification. Release names the latest
// release that shipped the local text, for contractRelease; Stamped says
// the local file carried init's provenance stamp, which a rewrite keeps.
type contractStatus struct {
	Role    string
	State   contractState
	Release string
	Stamped bool
}

// normalizeContract reads a local contract the way every comparison does:
// line endings are not part of the text (a CRLF checkout is the same
// contract), and init's provenance stamp is not either.
func normalizeContract(data []byte) (text string, stamped bool) {
	s := strings.ReplaceAll(string(data), "\r\n", "\n")
	return stripStamp(s), strings.HasPrefix(s, stampPrefix)
}

// releasedAs names the latest release whose role contract was exactly
// text, or "" when no release shipped it.
func releasedAs(history []releasedContract, role, text string) string {
	sum := sha256.Sum256([]byte(text))
	want := hex.EncodeToString(sum[:])
	found := ""
	for _, h := range history {
		if h.Role == role && h.SHA256 == want {
			found = h.Release // history is in release order: the last match is the latest
		}
	}
	return found
}

// embeddedRoles lists the role contracts the binary embeds, in the embed's
// (alphabetical) order; an extension is never a contract.
func embeddedRoles(contracts fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(contracts, config.RolesDir)
	if err != nil {
		return nil, err
	}
	var roles []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), localSuffix) || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		roles = append(roles, strings.TrimSuffix(e.Name(), ".md"))
	}
	return roles, nil
}

// classifyContracts measures each named role's local contract in dir (all
// embedded roles when roles is empty) against the embedded contract and the
// release history. It reads only local files and the binary. A role the
// binary does not embed is an error, as it is for roles diff.
func classifyContracts(dir string, contracts fs.FS, history []releasedContract, roles []string) ([]contractStatus, error) {
	all, err := embeddedRoles(contracts)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		roles = all
	}
	var out []contractStatus
	for _, role := range roles {
		if !slices.Contains(all, role) {
			return nil, fmt.Errorf("no embedded contract for role %q", role)
		}
		embedded, err := fs.ReadFile(contracts, contractPath(role))
		if err != nil {
			return nil, err
		}
		data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(contractPath(role))))
		if errors.Is(err, fs.ErrNotExist) {
			out = append(out, contractStatus{Role: role, State: contractAbsent})
			continue
		} else if err != nil {
			return nil, err
		}
		text, stamped := normalizeContract(data)
		st := contractStatus{Role: role, Stamped: stamped}
		switch {
		case text == string(embedded):
			st.State = contractCurrent
		case releasedAs(history, role, text) != "":
			st.State, st.Release = contractRelease, releasedAs(history, role, text)
		default:
			st.State = contractForked
		}
		out = append(out, st)
	}
	return out, nil
}

// rolesSyncCmd parses `roles sync [<role>...] [--dry-run]` — roles and the
// flag in any order — and runs the sync in the hub.
func rolesSyncCmd(w io.Writer, args []string, hub, dir string, contracts fs.FS) error {
	var flagArgs, roles []string
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			flagArgs = append(flagArgs, a)
		} else {
			roles = append(roles, a)
		}
	}
	flags := flag.NewFlagSet("roles sync", flag.ContinueOnError)
	dryRun := flags.Bool("dry-run", false, "print what would be written and committed; write nothing")
	if err := flags.Parse(flagArgs); err != nil {
		return err
	}
	if hub != "self" {
		return fmt.Errorf("roles sync writes the hub's contracts, and this repo is a spoke of %s (spokes hold no contracts) — run it in the hub", hub)
	}
	return rolesSync(w, dir, contracts, contractHistory, roles, *dryRun)
}

// rolesSync writes the embedded contract for every named role whose local
// file is absent or an earlier release's text, as one local commit that is
// never pushed — housekeeping's light path (SPEC §4): the tool defines the
// target, and the commit body names the tool. A fork anywhere among the
// named roles refuses CONTRACT_FORKED before anything is written, so exit 1
// still means nothing happened (SPEC §6).
func rolesSync(w io.Writer, dir string, contracts fs.FS, history []releasedContract, roles []string, dryRun bool) error {
	statuses, err := classifyContracts(dir, contracts, history, roles)
	if err != nil {
		return err
	}
	var forked []string
	for _, s := range statuses {
		if s.State == contractForked {
			forked = append(forked, fmt.Sprintf("%s (gh codecrew roles diff %s)", contractPath(s.Role), s.Role))
		}
	}
	if len(forked) > 0 {
		verb := "differ"
		if len(forked) == 1 {
			verb = "differs"
		}
		return refuse("CONTRACT_FORKED", "%s %s from the embedded %s contract and from every release's text: a fork is the project's (SPEC §7), and roles sync never overwrites one — reconcile it in a task, moving the project's additions to .codecrew/roles/<role>%s, or sync the other roles by naming them (gh codecrew roles sync <role>...)",
			strings.Join(forked, ", "), verb, version, localSuffix)
	}

	if repoRoot(dir) == "" {
		return fmt.Errorf("run roles sync inside a git repository: its result is one commit for a pull request")
	}

	write, wrote, kept := "wrote", "", "kept"
	if dryRun {
		write, kept = "would write", "would keep"
	}
	var lines, paths []string
	for _, s := range statuses {
		p := contractPath(s.Role)
		switch s.State {
		case contractCurrent:
			fmt.Fprintf(w, "%s %s — already the embedded %s contract\n", kept, p, version)
			continue
		case contractAbsent:
			wrote = "absent"
		case contractRelease:
			wrote = "was the " + s.Release + " text"
		}
		line := fmt.Sprintf("%s: written (%s)", p, wrote)
		lines = append(lines, line)
		paths = append(paths, p)
		fmt.Fprintf(w, "%s %s (%s)\n", write, p, wrote)
	}
	if len(paths) == 0 {
		fmt.Fprintf(w, "every contract is the embedded %s text — nothing to write\n", version)
		return nil
	}

	subject := "chore: sync codecrew role contracts to " + version
	body := strings.Join(lines, "\n") + "\n\n" +
		fmt.Sprintf("Target: the role contracts embedded in gh codecrew %s, written by\n`gh codecrew roles sync` (SPEC §7). Housekeeping (SPEC §4): reviewed, no task.", version)

	branch, detached := currentBranchOf(dir)
	onDefault := !detached && branch == defaultBranchOffline(dir)
	if onDefault {
		if _, err := git(dir, "rev-parse", "--verify", "--quiet", "refs/heads/"+syncBranch); err == nil {
			return fmt.Errorf("roles sync commits on %s when run on the default branch, and that branch already exists — merge or delete it, or switch to it, then rerun", syncBranch)
		}
	}

	if dryRun {
		target := branch
		switch {
		case detached:
			target = "no branch (HEAD is detached: the files would be written, not committed)"
		case onDefault:
			target = syncBranch + ", cut from " + branch
		}
		fmt.Fprintf(w, "would commit %d paths on %s: %q\n", len(paths), target, subject)
		fmt.Fprintln(w, "dry run: nothing written")
		return nil
	}

	for _, s := range statuses {
		if s.State == contractCurrent {
			continue
		}
		data, err := fs.ReadFile(contracts, contractPath(s.Role))
		if err != nil {
			return err
		}
		content := string(data)
		// The file keeps its form: init's stamp where the file had one or
		// there was no file (as init would write it), none where the
		// project keeps its contracts unstamped, as a hub's own are.
		if s.State == contractAbsent || s.Stamped {
			content = contractStamp(s.Role+".md") + content
		}
		full := filepath.Join(dir, filepath.FromSlash(contractPath(s.Role)))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			return err
		}
	}

	if detached {
		fmt.Fprintf(w, "note: HEAD is detached — the contracts are written but not committed; switch to a branch and run: %s\n", commitByHand(subject, paths))
		return nil
	}
	if onDefault {
		if _, err := git(dir, "switch", "-q", "-c", syncBranch); err != nil {
			fmt.Fprintf(w, "note: could not create %s (%v) — the contracts are written but not committed; on a branch of your own run: %s\n", syncBranch, err, commitByHand(subject, paths))
			return nil
		}
		branch = syncBranch
	}
	sha, ok := pathspecCommit(w, dir, subject, body, "the contracts", paths, paths)
	if !ok {
		return nil
	}
	fmt.Fprintf(w, "committed %s on %s: %q — the contracts only; your other changes are as they were\n", sha, branch, subject)
	fmt.Fprintf(w, "next: read it (git show %s), then push and open a pull request for review — roles sync never pushes; it is housekeeping, with no task (SPEC §4)\n", sha)
	return nil
}

// currentBranchOf names the checked-out branch, or reports a detached HEAD.
func currentBranchOf(dir string) (branch string, detached bool) {
	b, err := git(dir, "symbolic-ref", "--short", "HEAD")
	if err != nil {
		return "", true
	}
	return b, false
}

// defaultBranchOffline is the default branch as the clone records it —
// origin/HEAD — with no network call; "" when the clone does not know.
func defaultBranchOffline(dir string) string {
	ref, err := git(dir, "symbolic-ref", "--short", "refs/remotes/origin/HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(ref, "origin/")
}
