package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/radiusred/gh-codecrew/internal/config"
	"github.com/radiusred/gh-codecrew/internal/gh"
	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// ctx is everything a verb needs: the resolved project topology and the
// tracker backend.
type ctx struct {
	cfg     *config.Config
	current string // owner/repo the command runs in
	hub     string // owner/repo of the hub
	t       tracker.Tracker

	roles *config.Config // memoized routing table (local or hub)

	// teams memoizes team member sets per run (the #43 ctx pattern):
	// one fetch per team, however many logins are tested against it.
	teams map[string]map[string]bool
}

// loadConfig reads the pointer and checks its protocol version against the
// one this binary implements: a different major refuses
// (PROTOCOL_MISMATCH); "0.1" and a missing field proceed with a note on
// stderr (SPEC §5).
func loadConfig(dir string, notes io.Writer) (*config.Config, error) {
	cfg, err := config.Load(dir)
	if err != nil {
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

// ghFloor is the oldest gh the verbs work with: `gh pr checks --json`
// (cli/cli#9079, 2.50.0) is what task finish and the close's branch sweep
// read. A distribution-packaged 2.46 met it as a parse error inside the
// gate and a silently skipped sweep (#119 findings 21 and 30; #149). Raise
// it here, in one place, when a verb comes to need a newer gh.
const ghFloor = "2.50.0"

// ghVersion is a func var so tests can stand in for the installed gh.
var ghVersion = gh.Version

// checkGH refuses GH_TOO_OLD below the floor. A banner that does not parse
// proceeds with a note — an unexpected build must not lock the operator
// out of every verb.
func checkGH(notes io.Writer) error {
	v, err := ghVersion()
	if err != nil {
		fmt.Fprintf(notes, "note: could not read the gh version (%v); CodeCrew needs gh %s or later\n", err, ghFloor)
		return nil
	}
	if gh.CompareVersions(v, ghFloor) < 0 {
		return refuse("GH_TOO_OLD", "gh %s installed; CodeCrew needs %s or later (gh pr checks --json, cli/cli#9079) — upgrade gh", v, ghFloor)
	}
	return nil
}

// loadPointer is the one path every verb that reads the working repo's
// .codecrew.yml takes: the pointer, then the gh floor. status and roles
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
	current, err := gh.CurrentRepo()
	if err != nil {
		return nil, err
	}
	return &ctx{
		cfg:     cfg,
		current: current,
		hub:     cfg.HubRepo(current),
		t:       tracker.GitHub{},
	}, nil
}

// rolesConfig returns the config whose routing table governs role
// resolution. Spokes carry only the pointer config (SPEC §5), so when the
// local file has no roles the hub's .codecrew.yml is fetched (memoized); an
// unreadable hub config degrades to the local one rather than failing the
// verb — routing is advisory.
func (c *ctx) rolesConfig() *config.Config {
	if c.roles != nil {
		return c.roles
	}
	c.roles = c.cfg
	if len(c.cfg.Roles) == 0 {
		if data, err := c.t.FileContent(c.hub, ".codecrew.yml"); err == nil {
			if hubCfg, err := config.Parse(data); err == nil {
				c.roles = hubCfg
			}
		}
	}
	return c.roles
}

// teamMembers fetches a team's member logins — child-team members
// included, per the API contract. A func var so tests stub it (the
// convertManifest pattern).
var teamMembers = func(org, team string) (map[string]bool, error) {
	var members []struct {
		Login string `json:"login"`
	}
	if err := gh.JSON(&members, "api", "--paginate", fmt.Sprintf("/orgs/%s/teams/%s/members", org, team)); err != nil {
		return nil, err
	}
	set := make(map[string]bool, len(members))
	for _, m := range members {
		set[m.Login] = true
	}
	return set, nil
}

// inTeam reports whether login is a member of the team identity, through
// the per-run memo. Bot logins are never team members; an unreadable team
// resolves to no members (routing is advisory — the verbs that gate on a
// holder then refuse for absence, which fails closed).
func (c *ctx) inTeam(identity, login string) bool {
	if strings.HasSuffix(login, "[bot]") {
		return false
	}
	if c.teams == nil {
		c.teams = map[string]map[string]bool{}
	}
	set, ok := c.teams[identity]
	if !ok {
		org, team, valid := config.TeamIdentity(identity)
		if !valid {
			return false
		}
		set, _ = teamMembers(org, team)
		c.teams[identity] = set
	}
	return set[login]
}

// roleFor resolves a viewer login to its role name via the routing table —
// team-held roles resolve through membership (#44), everything else stays
// pure string matching.
func (c *ctx) roleFor(login string) string {
	if role := c.rolesConfig().RoleFor(login); role != "" {
		return role
	}
	for name, role := range c.rolesConfig().Roles {
		if _, _, isTeam := config.TeamIdentity(role.Identity); isTeam && c.inTeam(role.Identity, login) {
			return name
		}
	}
	return ""
}

// holdsRole reports whether login holds the named role: membership for a
// team-held seat (#44 — any member holds the role), config.HoldsRole for
// usernames and Apps.
func (c *ctx) holdsRole(login, role string) bool {
	if id := c.rolesConfig().Roles[role].Identity; id != "" {
		if _, _, isTeam := config.TeamIdentity(id); isTeam {
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
