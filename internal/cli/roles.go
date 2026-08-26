package cli

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	codecrew "github.com/radiusred/gh-codecrew"
	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// stampPrefix opens every provenance stamp contractStamp writes; stripStamp
// keys on it, so a reworded stamp still strips as long as this prefix opens
// the comment.
const stampPrefix = "<!-- scaffolded by codecrew "

// contractStamp is the provenance header init writes above a scaffolded
// contract: which release wrote it, and where upstream lives — enough for a
// later three-way judgment (the base is fetchable from the upstream repo at
// the stamped version).
func contractStamp(name string) string {
	return fmt.Sprintf("%s%s; upstream: radiusred/gh-codecrew roles/%s -->\n\n", stampPrefix, version, name)
}

// stripStamp removes the provenance stamp (and its trailing blank line) so
// comparisons see the contract alone. Content without a stamp — the hub's
// own roles/, or a hand-created file — passes through unchanged.
func stripStamp(content string) string {
	if !strings.HasPrefix(content, stampPrefix) {
		return content
	}
	if _, rest, ok := strings.Cut(content, "-->\n"); ok {
		return strings.TrimPrefix(rest, "\n")
	}
	return content
}

// contractDrift compares each local roles/ contract (stamp-stripped)
// against the copy embedded in this binary, returning the names that
// differ. Local files with no embedded counterpart, and embedded contracts
// with no local file (a pointer-only spoke), are not drift.
func contractDrift(dir string, contracts fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(contracts, "roles")
	if err != nil {
		return nil, err
	}
	var drifted []string
	for _, e := range entries {
		local, err := os.ReadFile(filepath.Join(dir, "roles", e.Name()))
		if err != nil {
			continue // no local copy — nothing to drift
		}
		embedded, err := fs.ReadFile(contracts, "roles/"+e.Name())
		if err != nil {
			return nil, err
		}
		if stripStamp(string(local)) != string(embedded) {
			drifted = append(drifted, strings.TrimSuffix(e.Name(), ".md"))
		}
	}
	return drifted, nil
}

// unifiedDiff is a minimal line diff (LCS) — enough to judge a contract
// change without shelling out to external tools.
func unifiedDiff(a, b string) string {
	al := strings.Split(strings.TrimSuffix(a, "\n"), "\n")
	bl := strings.Split(strings.TrimSuffix(b, "\n"), "\n")
	// lcs[i][j] = length of LCS of al[i:], bl[j:]
	lcs := make([][]int, len(al)+1)
	for i := range lcs {
		lcs[i] = make([]int, len(bl)+1)
	}
	for i := len(al) - 1; i >= 0; i-- {
		for j := len(bl) - 1; j >= 0; j-- {
			if al[i] == bl[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else {
				lcs[i][j] = max(lcs[i+1][j], lcs[i][j+1])
			}
		}
	}
	var out strings.Builder
	i, j := 0, 0
	for i < len(al) && j < len(bl) {
		switch {
		case al[i] == bl[j]:
			fmt.Fprintf(&out, "  %s\n", al[i])
			i++
			j++
		case lcs[i+1][j] >= lcs[i][j+1]:
			fmt.Fprintf(&out, "- %s\n", al[i])
			i++
		default:
			fmt.Fprintf(&out, "+ %s\n", bl[j])
			j++
		}
	}
	for ; i < len(al); i++ {
		fmt.Fprintf(&out, "- %s\n", al[i])
	}
	for ; j < len(bl); j++ {
		fmt.Fprintf(&out, "+ %s\n", bl[j])
	}
	return out.String()
}

// rolesDiff prints the local contract against the embedded one: "-" lines
// are local (this project's fork), "+" lines are the embedded contract at
// the installed release.
func rolesDiff(w io.Writer, dir string, contracts fs.FS, role string) error {
	embedded, err := fs.ReadFile(contracts, "roles/"+role+".md")
	if err != nil {
		return fmt.Errorf("no embedded contract for role %q", role)
	}
	local, err := os.ReadFile(filepath.Join(dir, "roles", role+".md"))
	if err != nil {
		return fmt.Errorf("no local roles/%s.md — run from the hub (spokes hold no contracts)", role)
	}
	stripped := stripStamp(string(local))
	if stripped == string(embedded) {
		fmt.Fprintf(w, "roles/%s.md matches the embedded %s contract\n", role, version)
		return nil
	}
	fmt.Fprintf(w, "roles/%s.md (local, -) vs embedded %s contract (+):\n", role, version)
	fmt.Fprint(w, unifiedDiff(stripped, string(embedded)))
	fmt.Fprintf(w, "\ncontracts are this project's fork — reconcile through a task and PR, never a blind overwrite\n")
	return nil
}

// localSuffix names a project's append-only extension to a role contract:
// roles/<role>.local.md, loaded after the contract itself (SPEC §7). It has
// no embedded counterpart, so contractDrift never sees it — extensions are
// not drift, and reconciling the contract never has to re-merge them.
const localSuffix = ".local.md"

// localPart is one extension in load order, labelled by where it came
// from so the composed text says which repo added each part.
type localPart struct {
	Source string // e.g. "roles/qa.local.md (hub)"
	Body   string
}

// composeContract appends each non-empty extension to the contract after
// a blank line and an HTML-comment marker naming its source. Append-only:
// neither side is parsed, and an extension that contradicts its contract
// is a review finding, not something a resolver decides.
func composeContract(base string, locals []localPart) string {
	var out strings.Builder
	out.WriteString(base)
	for _, l := range locals {
		if strings.TrimSpace(l.Body) == "" {
			continue
		}
		if !strings.HasSuffix(out.String(), "\n") {
			out.WriteString("\n")
		}
		fmt.Fprintf(&out, "\n<!-- extension: %s -->\n\n%s", l.Source, l.Body)
		if !strings.HasSuffix(l.Body, "\n") {
			out.WriteString("\n")
		}
	}
	return out.String()
}

// composedContract assembles what a dispatched session loads for role:
// the hub's roles/<role>.md (the project's fork of the contract), then the
// hub's roles/<role>.local.md, then the working repo's roles/<role>.local.md
// when it is a spoke (spokeDir non-empty). hubRead fetches a hub path —
// from disk when the working repo is the hub, through the tracker from a
// spoke. Missing extensions are skipped; a missing contract is an error.
func composedContract(role string, hubRead func(string) ([]byte, error), spokeDir string) (string, error) {
	base, err := hubRead("roles/" + role + ".md")
	if err != nil {
		return "", fmt.Errorf("no roles/%s.md in the hub", role)
	}
	var locals []localPart
	if data, err := hubRead("roles/" + role + localSuffix); err == nil {
		locals = append(locals, localPart{Source: "roles/" + role + localSuffix + " (hub)", Body: string(data)})
	}
	if spokeDir != "" {
		if data, err := os.ReadFile(filepath.Join(spokeDir, "roles", role+localSuffix)); err == nil {
			locals = append(locals, localPart{Source: "roles/" + role + localSuffix + " (spoke)", Body: string(data)})
		}
	}
	return composeContract(string(base), locals), nil
}

// rolesShow prints the contract a dispatched session loads — the hub's
// contract with its local extensions appended in order — or, with latest,
// the embedded contract whole: the virtual "latest", materialized on demand
// instead of shadow-copied into the repo.
func rolesShow(w io.Writer, role string, latest bool, contracts fs.FS, hubRead func(string) ([]byte, error), spokeDir string) error {
	if latest {
		data, err := fs.ReadFile(contracts, "roles/"+role+".md")
		if err != nil {
			return fmt.Errorf("no embedded contract for role %q", role)
		}
		_, err = w.Write(data)
		return err
	}
	composed, err := composedContract(role, hubRead, spokeDir)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, composed)
	return err
}

// rolesCmd dispatches the roles subverbs against the installed binary's
// embedded contracts and the local hub checkout.
func rolesCmd(w io.Writer, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: codecrew roles diff <role> | codecrew roles show <role> [--latest]")
	}
	sub, role := args[0], args[1]
	cfg, err := config.Load(".")
	if err != nil {
		return err
	}
	switch sub {
	case "diff":
		return rolesDiff(w, cfg.Dir, codecrew.Roles, role)
	case "show":
		fs := flag.NewFlagSet("roles show", flag.ContinueOnError)
		latest := fs.Bool("latest", false, "print the contract embedded in the installed binary")
		if err := fs.Parse(args[2:]); err != nil {
			return err
		}
		// In the hub the contract and its extension are on disk; from a
		// spoke they live in the hub repo (fetched from its default branch)
		// and only the spoke's own extension is local.
		hubRead := func(path string) ([]byte, error) {
			return os.ReadFile(filepath.Join(cfg.Dir, filepath.FromSlash(path)))
		}
		spokeDir := ""
		if cfg.Hub != "self" {
			hub := cfg.Hub
			hubRead = func(path string) ([]byte, error) { return tracker.GitHub{}.FileContent(hub, path) }
			spokeDir = cfg.Dir
		}
		return rolesShow(w, role, *latest, codecrew.Roles, hubRead, spokeDir)
	default:
		return fmt.Errorf("unknown subcommand roles %s", sub)
	}
}
