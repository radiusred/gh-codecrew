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
func holder(roles map[string]config.Role, name string) (config.Identity, error) {
	role, ok := roles[name]
	if !ok && len(roles) > 0 {
		// The coordinator row arrived after 1.0 hubs scaffolded their
		// tables; a table without it still has a coordinator — the
		// operator, as every unrouted seat is (SPEC §5, §7).
		if name == "coordinator" {
			return config.Identity{}, nil
		}
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
// and correct from a pointer-only spoke because resolution falls back to
// the hub's routing table.
func roleHolder(w io.Writer, args []string) error {
	fs := flag.NewFlagSet("role", flag.ContinueOnError)
	login := fs.Bool("login", false, "print the review-requestable handle, or nothing for an App or the operator")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: gh codecrew role <name> [--login]")
	}
	c, err := load()
	if err != nil {
		return err
	}
	id, err := holder(c.rolesConfig().Roles, fs.Arg(0))
	if err != nil {
		return err
	}
	printHolder(w, id, *login)
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
