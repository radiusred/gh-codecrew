package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/radiusred/gh-codecrew/internal/config"
)

// holder resolves a role name against a routing table: the routed
// identity, the operator's `~` (explicitly, or because no table is
// declared), or an error for a name absent from a declared table.
//
// `coordinator` is no such name. Until 2.0 it was special-cased to `~`,
// because the row arrived after 1.0 hubs had scaffolded their tables; the
// shim is gone (M13-R7) and a declared table missing the row is the error
// every other missing role is. `init` scaffolds the row and
// `gh codecrew migrate` adds it to a 1.x table (SPEC §6), so the only way
// to reach this error is a table someone edited to remove a seat that
// exists.
func holder(roles map[string]config.Role, name string) (config.Identity, error) {
	role, ok := roles[name]
	if !ok && len(roles) > 0 {
		return config.Identity{}, fmt.Errorf("role %q is not in the routing table", name)
	}
	return role.Identity, nil
}

// roleHolder prints the typed identity a role routes to — `app:<slug>`,
// `user:<login>`, `team:<org>/<slug>`, or `~` when the role is
// operator-held. With --login it prints instead the handle GitHub will
// accept a review request for (the login of a user: seat, org/slug for a
// team: seat) and nothing at all for an App or the operator, neither of
// which can be requested — so a caller's whole decision is whether the
// output is empty (SPEC §6; the implementer contract). Script-consumable,
// and correct from a pointer-only spoke because a spoke's resolution *is*
// the hub's routing table, fetched at load — there is no fallback, and a
// hub that cannot be read refuses rather than answering `~` (SPEC §6).
// parseRoleArgs reads `<name> [--login]` in either order. The role name
// leads in practice, and Go's flag package stops at the first non-flag
// argument, so the name comes off the front the way task's ref does.
func parseRoleArgs(args []string) (name string, login bool, err error) {
	name, args = splitLeadingRef(args)
	fs := flag.NewFlagSet("role", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	flagLogin := fs.Bool("login", false, "print the review-requestable handle, or nothing for an App or the operator")
	if err := fs.Parse(args); err != nil {
		return "", false, err
	}
	rest := fs.NArg()
	if name == "" && rest == 1 {
		name, rest = fs.Arg(0), 0
	}
	if name == "" || rest > 0 {
		return "", false, fmt.Errorf("usage: gh codecrew role <name> [--login]")
	}
	return name, *flagLogin, nil
}

func roleHolder(w io.Writer, args []string) error {
	name, login, err := parseRoleArgs(args)
	if err != nil {
		return err
	}
	c, err := load()
	if err != nil {
		return err
	}
	id, err := holder(c.rolesConfig().Roles, name)
	if err != nil {
		return err
	}
	printHolder(w, id, login)
	return nil
}

// printHolder writes the one line (or, for --login on an unrequestable
// seat, the no line) that `role <name>` is read for.
func printHolder(w io.Writer, id config.Identity, login bool) {
	if login {
		if handle := id.Login(); handle != "" {
			fmt.Fprintln(w, handle)
		}
		return
	}
	fmt.Fprintln(w, id)
}
