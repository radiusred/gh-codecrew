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
		st, err := classifyFile(dir, role, contractPath(role), string(embedded), history)
		if err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, nil
}

// classifyAgents measures .codecrew/AGENTS.md exactly as a contract is
// measured: against the scaffold this binary writes and every release's
// scaffold before it (#372). Its history rows carry its path as the role.
func classifyAgents(dir, embedded string, history []releasedContract) (contractStatus, error) {
	return classifyFile(dir, config.AgentsFile, config.AgentsFile, embedded, history)
}

// classifyFile is one file roles sync owns, measured: absent, the embedded
// text, a release's text (the latest release that shipped it), or a fork.
func classifyFile(dir, name, rel, embedded string, history []releasedContract) (contractStatus, error) {
	data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(rel)))
	if errors.Is(err, fs.ErrNotExist) {
		return contractStatus{Role: name, State: contractAbsent}, nil
	} else if err != nil {
		return contractStatus{}, err
	}
	text, stamped := normalizeContract(data)
	st := contractStatus{Role: name, Stamped: stamped}
	switch {
	case text == embedded:
		st.State = contractCurrent
	case releasedAs(history, name, text) != "":
		st.State, st.Release = contractRelease, releasedAs(history, name, text)
	default:
		st.State = contractForked
	}
	return st, nil
}

// syncSource is what roles sync writes from: the role contracts (nil in a
// spoke, which holds none) and the .codecrew/AGENTS.md scaffold ("" when
// the agents file is not among the targets a caller offers).
type syncSource struct {
	contracts fs.FS
	agents    string
}

// syncTargets resolves roles sync's arguments. None means everything the
// repo holds: in a hub every embedded contract and the agents file, in a
// spoke the agents file alone. A role name selects that contract and never
// the agents file; the path .codecrew/AGENTS.md selects the agents file.
func syncTargets(src syncSource, args []string) (roles []string, agents bool, err error) {
	if len(args) == 0 {
		if src.contracts != nil {
			if roles, err = embeddedRoles(src.contracts); err != nil {
				return nil, false, err
			}
		}
		return roles, src.agents != "", nil
	}
	for _, a := range args {
		if filepath.ToSlash(a) == config.AgentsFile {
			agents = true
			continue
		}
		roles = append(roles, a)
	}
	return roles, agents, nil
}

// rolesSyncCmd parses `roles sync [<role>|.codecrew/AGENTS.md ...]
// [--dry-run]` — targets and the flag in any order — and runs the sync: in
// a hub over the contracts and the agents file, in a spoke over the agents
// file alone.
func rolesSyncCmd(w io.Writer, args []string, hub, dir string, contracts fs.FS) error {
	var flagArgs, targets []string
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			flagArgs = append(flagArgs, a)
		} else {
			targets = append(targets, a)
		}
	}
	flags := flag.NewFlagSet("roles sync", flag.ContinueOnError)
	dryRun := flags.Bool("dry-run", false, "print what would be written and committed; write nothing")
	if err := flags.Parse(flagArgs); err != nil {
		return err
	}
	src := syncSource{contracts: contracts, agents: agentsScaffold}
	if hub != "self" {
		for _, t := range targets {
			if filepath.ToSlash(t) != config.AgentsFile {
				return fmt.Errorf("roles sync writes the hub's contracts, and this repo is a spoke of %s (spokes hold no contracts) — run it in the hub; here it syncs %s alone", hub, config.AgentsFile)
			}
		}
		src.contracts = nil
	}
	return rolesSync(w, dir, src, contractHistory, targets, *dryRun)
}

// rolesSync writes the embedded text for every named target whose local
// file is absent or an earlier release's text — role contracts, and
// .codecrew/AGENTS.md (#372) — as one local commit that is never pushed:
// housekeeping's light path (SPEC §4), where the tool defines the target
// and the commit body names the tool. A fork anywhere among the targets
// refuses before anything is written — CONTRACT_FORKED for a contract,
// AGENTS_FORKED for the agents file alone — so exit 1 still means nothing
// happened (SPEC §6).
func rolesSync(w io.Writer, dir string, src syncSource, history []releasedContract, targets []string, dryRun bool) error {
	roles, withAgents, err := syncTargets(src, targets)
	if err != nil {
		return err
	}
	var statuses []contractStatus
	if len(roles) > 0 {
		if src.contracts == nil {
			return fmt.Errorf("this repo holds no role contracts to sync")
		}
		if statuses, err = classifyContracts(dir, src.contracts, history, roles); err != nil {
			return err
		}
	}
	if withAgents {
		st, err := classifyAgents(dir, src.agents, history)
		if err != nil {
			return err
		}
		statuses = append(statuses, st)
	}

	var forked []string
	agentsForked := false
	for _, s := range statuses {
		if s.State != contractForked {
			continue
		}
		if s.Role == config.AgentsFile {
			agentsForked = true
			continue
		}
		forked = append(forked, fmt.Sprintf("%s (gh codecrew roles diff %s)", contractPath(s.Role), s.Role))
	}
	agentsNote := fmt.Sprintf("%s (gh codecrew roles diff %s)", config.AgentsFile, config.AgentsFile)
	if len(forked) > 0 {
		verb := "differ"
		if len(forked) == 1 {
			verb = "differs"
		}
		also := ""
		if agentsForked {
			also = fmt.Sprintf("; %s differs from every release's scaffold too", agentsNote)
		}
		return refuse("CONTRACT_FORKED", "%s %s from the embedded %s contract and from every release's text: a fork is the project's (SPEC §7), and roles sync never overwrites one — reconcile it in a task, moving the project's additions to .codecrew/roles/<role>%s, or sync the other roles by naming them (gh codecrew roles sync <role>...)%s",
			strings.Join(forked, ", "), verb, version, localSuffix, also)
	}
	if agentsForked {
		way := "sync the contracts alone by naming the roles (gh codecrew roles sync <role>...)"
		if src.contracts == nil {
			way = "a spoke has nothing else for roles sync to write"
		}
		return refuse("AGENTS_FORKED", "%s differs from the embedded %s scaffold and from every release's: it is the project's own (SPEC §7), and roles sync never overwrites it — bring in what the new scaffold adds by hand, in a task; or %s",
			agentsNote, version, way)
	}

	if repoRoot(dir) == "" {
		return fmt.Errorf("run roles sync inside a git repository: its result is one commit for a pull request")
	}

	write, wrote, kept := "wrote", "", "kept"
	if dryRun {
		write, kept = "would write", "would keep"
	}
	var lines, paths []string
	contractsWritten, agentsWritten := false, false
	for _, s := range statuses {
		p, what := contractPath(s.Role), "contract"
		if s.Role == config.AgentsFile {
			p, what = config.AgentsFile, "scaffold"
		}
		switch s.State {
		case contractCurrent:
			fmt.Fprintf(w, "%s %s — already the embedded %s %s\n", kept, p, version, what)
			continue
		case contractAbsent:
			wrote = "absent"
		case contractRelease:
			wrote = "was the " + s.Release + " text"
		}
		if s.Role == config.AgentsFile {
			agentsWritten = true
		} else {
			contractsWritten = true
		}
		line := fmt.Sprintf("%s: written (%s)", p, wrote)
		lines = append(lines, line)
		paths = append(paths, p)
		fmt.Fprintf(w, "%s %s (%s)\n", write, p, wrote)
	}
	if len(paths) == 0 {
		fmt.Fprintf(w, "every file roles sync writes here is the embedded %s text — nothing to write\n", version)
		return nil
	}

	subject, target, noun := "chore: sync codecrew role contracts to "+version, "the role contracts", "the contracts"
	switch {
	case agentsWritten && contractsWritten:
		subject = "chore: sync codecrew role contracts and agents file to " + version
		target, noun = "the role contracts and "+config.AgentsFile, "the contracts and the agents file"
	case agentsWritten:
		subject = "chore: sync the codecrew agents file to " + version
		target, noun = config.AgentsFile, "the agents file"
	}
	body := strings.Join(lines, "\n") + "\n\n" +
		fmt.Sprintf("Target: %s embedded in gh codecrew %s, written by\n`gh codecrew roles sync` (SPEC §7). Housekeeping (SPEC §4): reviewed, no task.", target, version)

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
		rel, content := config.AgentsFile, src.agents
		if s.Role != config.AgentsFile {
			data, err := fs.ReadFile(src.contracts, contractPath(s.Role))
			if err != nil {
				return err
			}
			rel, content = contractPath(s.Role), string(data)
			// The file keeps its form: init's stamp where the file had one or
			// there was no file (as init would write it), none where the
			// project keeps its contracts unstamped, as a hub's own are.
			// The agents file is never stamped, as init writes it.
			if s.State == contractAbsent || s.Stamped {
				content = contractStamp(s.Role+".md") + content
			}
		}
		full := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			return err
		}
	}

	are := "are"
	if !contractsWritten {
		are = "is"
	}
	if detached {
		fmt.Fprintf(w, "note: HEAD is detached — %s %s written but not committed; switch to a branch and run: %s\n", noun, are, commitByHand(subject, paths))
		return nil
	}
	if onDefault {
		if _, err := git(dir, "switch", "-q", "-c", syncBranch); err != nil {
			fmt.Fprintf(w, "note: could not create %s (%v) — %s %s written but not committed; on a branch of your own run: %s\n", syncBranch, err, noun, are, commitByHand(subject, paths))
			return nil
		}
		branch = syncBranch
	}
	sha, ok := pathspecCommit(w, dir, subject, body, noun, paths, paths)
	if !ok {
		return nil
	}
	fmt.Fprintf(w, "committed %s on %s: %q — %s only; your other changes are as they were\n", sha, branch, subject, noun)
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
