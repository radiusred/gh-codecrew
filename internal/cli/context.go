package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// ctx is everything a verb needs: the resolved project topology and the
// tracker backend.
type ctx struct {
	cfg     *config.Config
	current string // owner/repo the command runs in
	hub     string // owner/repo of the hub
	t       tracker.Tracker

	roles *config.Config // the routing table (this repo's, or the hub's), settled at load

	// teams memoizes team member sets per run (the #43 ctx pattern):
	// one fetch per team, however many logins are tested against it.
	teams map[string]map[string]bool
}

// loadConfig reads the pointer and checks its protocol version against the
// one this binary implements: a different major refuses
// (PROTOCOL_MISMATCH), and a "0.1" pointer is one of them — an older
// major, refused with the detail that names gh codecrew migrate, since
// 1.0's acceptance of it went with the other shims (M13-R7). A missing
// field is the only version the check proceeds on, with a note on stderr
// (SPEC §5). A repo still on the protocol 1.x layout refuses
// LAYOUT_LEGACY — this binary does not read that layout, it names the verb
// that moves it forward. A routing row whose identity carries no kind
// refuses IDENTITY_UNTYPED, and a spoke pointer carrying a roles: block
// refuses SPOKE_ROUTING — config detects each condition, the CLI names it,
// as with the protocol check.
func loadConfig(dir string, notes io.Writer) (*config.Config, error) {
	cfg, err := config.Load(dir)
	if err != nil {
		var legacy *config.LegacyLayoutError
		if errors.As(err, &legacy) {
			return nil, refuseLegacyLayout(legacy.Dir, legacy.Found, "run gh codecrew migrate to move it to the "+protocolVersion+" layout")
		}
		var untyped *config.UntypedIdentityError
		if errors.As(err, &untyped) {
			return nil, refuse("IDENTITY_UNTYPED", "%v", untyped)
		}
		var spokeRoles *config.SpokeRoutingError
		if errors.As(err, &spokeRoles) {
			return nil, refuse("SPOKE_ROUTING", "%v", spokeRoles)
		}
		return nil, err
	}
	note, err := config.Compatible(cfg.Codecrew, protocolVersion)
	if err != nil {
		return nil, refuse("PROTOCOL_MISMATCH", "%v", err)
	}
	if note != "" {
		fmt.Fprintln(notes, note)
	}
	return cfg, nil
}

// refuseLegacyLayout is the one wording for a protocol 1.x repo met by a
// 2.0 binary, raised from the pointer walk and from init alike: what was
// found, what this binary speaks, and the verb that moves the repo — never
// a hand edit, and never a read of the old layout. fix completes the
// sentence, because the two callers are in different positions: one was
// looking for a pointer, the other was about to write one.
func refuseLegacyLayout(dir string, found []string, fix string) error {
	return refuse("LAYOUT_LEGACY", "%s holds the protocol 1.x layout (%s); this codecrew implements protocol %s — %s (SPEC §3, §5)",
		dir, strings.Join(found, ", "), protocolVersion, fix)
}

// ghFloor is the oldest gh the verbs work with: `gh pr checks --json`
// (cli/cli#9079, 2.50.0) is what task finish and the close's branch sweep
// read. A distribution-packaged 2.46 met it as a parse error inside the
// gate and a silently skipped sweep (#119 findings 21 and 30; #149). Raise
// it here, in one place, when a verb comes to need a newer gh.
const ghFloor = "2.50.0"

// ghVersion is a func var so tests can stand in for the installed gh. It
// reads the version through the venue, like every other call: the check
// runs before any ctx exists, so it holds its own GitHub rather than one
// a verb hands it.
var ghVersion = tracker.GitHub{}.ClientVersion

// checkGH refuses GH_TOO_OLD below the floor. A banner that does not parse
// proceeds with a note — an unexpected build must not lock the operator
// out of every verb.
func checkGH(notes io.Writer) error {
	v, err := ghVersion()
	if err != nil {
		fmt.Fprintf(notes, "note: could not read the gh version (%v); CodeCrew needs gh %s or later\n", err, ghFloor)
		return nil
	}
	if tracker.CompareVersions(v, ghFloor) < 0 {
		return refuse("GH_TOO_OLD", "gh %s installed; CodeCrew needs %s or later (gh pr checks --json, cli/cli#9079) — upgrade gh", v, ghFloor)
	}
	return nil
}

// loadPointer is the one path every verb that reads the working repo's
// .codecrew/config.yml takes: the pointer, then the gh floor. status and roles
// read the pointer without building a ctx, so the check lives here, not in
// load() — the reviewer of #153 found it bypassed there.
func loadPointer(notes io.Writer) (*config.Config, error) {
	cfg, err := loadConfig(".", notes)
	if err != nil {
		return nil, err
	}
	if err := checkGH(notes); err != nil {
		return nil, err
	}
	return cfg, nil
}

func load() (*ctx, error) {
	cfg, err := loadPointer(os.Stderr)
	if err != nil {
		return nil, err
	}
	current, err := tracker.GitHub{}.CurrentRepo()
	if err != nil {
		if ghErr := unreachable(err); ghErr != nil {
			return nil, ghErr
		}
		return nil, err
	}
	c := &ctx{
		cfg:     cfg,
		current: current,
		hub:     cfg.HubRepo(current),
		t:       tracker.GitHub{},
	}
	if err := c.resolveRoles(os.Stderr); err != nil {
		return nil, err
	}
	return c, nil
}

// resolveRoles settles the routing table before any verb can consult it,
// so a ctx never exists without one. That is the whole of "routing fails
// closed" (M13-R5): an unreadable hub used to degrade to the local, empty
// table, which turns holder() into `~` for every role — and with it the
// holder-review gate into "any non-author approved" and the QA verdict
// count into "anyone commented" (the Claude scan on #254, finding 2).
// Resolving here rather than lazily at each call site means no predicate
// — the review gate, the verdict count, crewIdentity, roleFor — can
// forget to handle the error.
//
// The table is chosen by topology, not by emptiness:
//
//   - hub: self — the local pointer *is* the table. No fetch, so a hub
//     keeps working with no network, and a hub declaring no table is
//     legitimately `~` everywhere, as it always was.
//   - hub: owner/repo — the table is always the hub's, fetched here. A
//     spoke may not carry one at all (SPOKE_ROUTING), so there is nothing
//     local to prefer.
//
// "Read fine, no table" and "could not read" are told apart by whether the
// fetch and the parse succeeded — never by the table being empty.
func (c *ctx) resolveRoles(notes io.Writer) error {
	if c.roles != nil {
		return nil
	}
	if c.cfg.Hub == "self" {
		c.roles = c.cfg
		return nil
	}
	data, err := c.t.FileContent(c.hub, config.Pointer)
	if err != nil {
		if ghErr := unreachable(err); ghErr != nil {
			return ghErr
		}
		// The shapes have different remediations, and handing over the
		// wrong one sends an operator to a writing verb against a hub
		// that is perfectly healthy (checky's finding on PR #279). A 403
		// is access this seat does not have, named the way
		// NO_CHECKS_PERMISSION names a permission. A 404 has three causes
		// and the detail owns all of them, because GitHub answers 404 —
		// not 403 — for a repo a token cannot see at all, so "absent"
		// alone would send the same operator to `migrate` by the other
		// door (probed: a private hub this App is not installed on
		// answers 404).
		fix := fmt.Sprintf("%s answered but would not hand this seat the file — check the identity this run mints can read it (a private hub needs the App installed there, with contents: read)", c.hub)
		if strings.Contains(err.Error(), "HTTP 404") {
			fix = fmt.Sprintf("%s has no such file that this seat can see, which is three conditions: a hub still on the protocol 1.x layout, moved forward with gh codecrew migrate; a hub: line naming the wrong repo; or a private hub this seat's identity is not installed on, since GitHub answers 404 rather than 403 for a repo a token cannot see", c.hub)
		}
		return refuse("HUB_UNREADABLE", "the hub %s's %s could not be read (%v) — this repo is a spoke and the hub carries the routing table, so no role can be resolved; %s (SPEC §5, §6)", c.hub, config.Pointer, err, fix)
	}
	hubCfg, err := config.Parse(data)
	if err != nil {
		return refuse("HUB_UNREADABLE", "the hub %s's %s does not parse (%v) — this repo is a spoke and the hub carries the routing table, so no role can be resolved (SPEC §5, §6)", c.hub, config.Pointer, err)
	}
	note, err := config.CompatibleHub(c.hub, hubCfg.Codecrew, protocolVersion)
	if err != nil {
		return refuse("PROTOCOL_MISMATCH", "%v", err)
	}
	if note != "" {
		fmt.Fprintln(notes, note)
	}
	c.roles = hubCfg
	return nil
}

// unreachable turns a gh failure that never reached GitHub — offline, no
// DNS, no credentials — into its own refusal, and returns nil for
// everything else so the caller can name its own condition. GitHub
// answering 403 or 404 is not this: the API was reached (SPEC §6).
func unreachable(err error) error {
	if !tracker.Unreachable(err) {
		return nil
	}
	return refuse("GH_UNREACHABLE", "GitHub could not be reached (%v) — check the network and that gh is authenticated (gh auth status), or mint the seat's token with gh codecrew identity token <slug>; codecrew version, help, and roles show/diff in a hub need no network (SPEC §6)", err)
}

// rolesConfig returns the routing table that governs role resolution,
// settled once by resolveRoles at load. It cannot fail: a ctx that reached
// a verb has a table.
func (c *ctx) rolesConfig() *config.Config {
	return c.roles
}

// inTeam reports whether login is a member of the team identity, through
// the per-run memo. Bot logins are never team members; an unreadable team
// resolves to no members, so the verbs that gate on a holder refuse for
// absence — the same direction resolveRoles fails in.
func (c *ctx) inTeam(identity config.Identity, login string) bool {
	if strings.HasSuffix(login, "[bot]") {
		return false
	}
	if c.teams == nil {
		c.teams = map[string]map[string]bool{}
	}
	set, ok := c.teams[identity.Value]
	if !ok {
		org, team, valid := identity.Team()
		if !valid {
			return false
		}
		set = c.teamMembers(org, team)
		c.teams[identity.Value] = set
	}
	return set[login]
}

// teamMembers reads a team's members through the venue and returns them
// as the set inTeam tests against. An unreadable team is no members: the
// verbs that gate on a holder then refuse for absence, the same direction
// resolveRoles fails in.
func (c *ctx) teamMembers(org, team string) map[string]bool {
	logins, err := c.t.TeamMembers(org, team)
	if err != nil {
		return nil
	}
	set := make(map[string]bool, len(logins))
	for _, l := range logins {
		set[l] = true
	}
	return set
}

// roleFor resolves a viewer login to its role name via the routing table —
// team-held roles resolve through membership (#44), everything else stays
// pure string matching.
func (c *ctx) roleFor(login string) string {
	if role := c.rolesConfig().RoleFor(login); role != "" {
		return role
	}
	for name, role := range c.rolesConfig().Roles {
		if role.Identity.Kind == config.KindTeam && c.inTeam(role.Identity, login) {
			return name
		}
	}
	return ""
}

// holdsRole reports whether login holds the named role: membership for a
// team-held seat (#44 — any member holds the role), config.HoldsRole for
// usernames and Apps.
func (c *ctx) holdsRole(login, role string) bool {
	if id := c.rolesConfig().Roles[role].Identity; !id.Operator() {
		if id.Kind == config.KindTeam {
			return c.inTeam(id, login)
		}
		return c.rolesConfig().HoldsRole(login, role)
	}
	// Unrouted (~): the operator fallback must be crew-free through the
	// team-aware lens too — config.RoleFor is team-blind by design, so a
	// member of a role-holding team would otherwise hold every unrouted
	// role (checky's finding on PR #101).
	return c.rolesConfig().HoldsRole(login, role) && c.roleFor(login) == ""
}

// refusal is a blocked gate: a machine-readable code plus a human detail.
// Verbs exit nonzero with "refused[CODE]: detail" so agents can act on the
// specific unmet condition (SPEC.md §6).
type refusal struct {
	Code   string
	Detail string
}

func (r refusal) Error() string {
	return fmt.Sprintf("refused[%s]: %s", r.Code, r.Detail)
}

func refuse(code, format string, args ...any) error {
	return refusal{Code: code, Detail: fmt.Sprintf(format, args...)}
}
