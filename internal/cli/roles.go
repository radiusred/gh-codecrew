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

// rolesShow prints the embedded contract whole: the virtual "latest",
// materialized on demand instead of shadow-copied into the repo.
func rolesShow(w io.Writer, contracts fs.FS, role string, latest bool) error {
	if !latest {
		return fmt.Errorf("usage: codecrew roles show <role> --latest (the local file is readable directly)")
	}
	data, err := fs.ReadFile(contracts, "roles/"+role+".md")
	if err != nil {
		return fmt.Errorf("no embedded contract for role %q", role)
	}
	_, err = w.Write(data)
	return err
}

// rolesCmd dispatches the roles subverbs against the installed binary's
// embedded contracts and the local hub checkout.
func rolesCmd(w io.Writer, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: codecrew roles diff <role> | codecrew roles show <role> --latest")
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
		return rolesShow(w, codecrew.Roles, role, *latest)
	default:
		return fmt.Errorf("unknown subcommand roles %s", sub)
	}
}
