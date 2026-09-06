// Package config locates and parses the .codecrew/config.yml pointer file
// that every repo in a CodeCrew project carries (SPEC.md §3, §5).
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// IdentityKind is the kind of GitHub principal a routing row names. The
// grammar types the principal, not who is at the keyboard (SPEC §5).
type IdentityKind uint8

const (
	// KindOperator is `~`, an empty value, or an absent identity key: the
	// human operator, and any session acting under the operator's own
	// auth. It is the zero value because a YAML null never reaches
	// Identity's unmarshaler — an unrouted row decodes to the zero
	// Identity, and that must mean the operator.
	KindOperator IdentityKind = iota
	// KindApp is `app:<slug>` — a GitHub App, whose viewer login is
	// "<slug>[bot]".
	KindApp
	// KindUser is `user:<login>` — one named human.
	KindUser
	// KindTeam is `team:<org>/<slug>` — any member of the team, child
	// teams included.
	KindTeam
	// KindUntyped is anything else: a 1.0 table's bare string. Parse
	// refuses it rather than guessing which of the three it meant.
	KindUntyped
)

// Identity is a routing row's identity, parsed once at load into a kind
// and its name, so no caller has to infer the kind from the value's shape
// (SPEC §5). The zero Identity is the operator.
type Identity struct {
	Kind IdentityKind
	// Value is the App slug, the login, or "org/team-slug"; empty for the
	// operator, and the raw text for an untyped value.
	Value string
}

// ParseIdentity classifies a routing table's identity value. Anything the
// grammar does not name is KindUntyped carrying the raw text — the role
// row it sits on is only known to Parse, which is where the refusal is
// raised.
func ParseIdentity(s string) Identity {
	s = strings.TrimSpace(s)
	if s == "" || s == "~" {
		return Identity{}
	}
	kind, value, ok := strings.Cut(s, ":")
	if !ok {
		return Identity{Kind: KindUntyped, Value: s}
	}
	switch kind {
	case "app":
		if value != "" && !strings.Contains(value, "/") {
			return Identity{Kind: KindApp, Value: value}
		}
	case "user":
		if value != "" && !strings.Contains(value, "/") {
			return Identity{Kind: KindUser, Value: value}
		}
	case "team":
		if org, team, cut := strings.Cut(value, "/"); cut && org != "" && team != "" && !strings.Contains(team, "/") {
			return Identity{Kind: KindTeam, Value: value}
		}
	}
	return Identity{Kind: KindUntyped, Value: s}
}

// String renders the identity back into the grammar — the form `role
// <name>` prints and a routing table carries.
func (i Identity) String() string {
	switch i.Kind {
	case KindApp:
		return "app:" + i.Value
	case KindUser:
		return "user:" + i.Value
	case KindTeam:
		return "team:" + i.Value
	case KindUntyped:
		return i.Value
	default:
		return "~"
	}
}

// Operator reports whether the seat is held by the operator — an
// unrouted row.
func (i Identity) Operator() bool { return i.Kind == KindOperator }

// Team splits a team identity into its org and team slug. Every other
// kind reports false.
func (i Identity) Team() (org, team string, ok bool) {
	if i.Kind != KindTeam {
		return "", "", false
	}
	org, team, _ = strings.Cut(i.Value, "/")
	return org, team, true
}

// Login is the handle GitHub will accept a review request for: the login
// of a `user:` seat, "org/slug" for a `team:` seat — the form `gh pr
// create --reviewer` wants — and empty for an App, which cannot be
// requested, or the operator, who is not a name (SPEC §6; the implementer
// contract branches on it being empty).
func (i Identity) Login() string {
	switch i.Kind {
	case KindUser, KindTeam:
		return i.Value
	default:
		return ""
	}
}

// UnmarshalYAML parses the identity value as it decodes, so the routing
// table is typed exactly once. A YAML null — `identity: ~`, or no
// identity key at all — never reaches here: it leaves the zero Identity,
// which is the operator.
func (i *Identity) UnmarshalYAML(n *yaml.Node) error {
	var s string
	if err := n.Decode(&s); err != nil {
		return fmt.Errorf("identity must be a string or ~: %w", err)
	}
	*i = ParseIdentity(s)
	return nil
}

// UntypedIdentityError is a routing row whose identity names no kind of
// principal — a 1.0 table read by a 2.0 binary. The CLI gives it the
// refusal code (IDENTITY_UNTYPED); this type carries the row so the
// detail can name it.
type UntypedIdentityError struct {
	Role  string
	Value string
}

func (e *UntypedIdentityError) Error() string {
	return fmt.Sprintf("roles.%s.identity: %q names no kind of principal — a routing identity is `~` (the operator), `app:<slug>` (a GitHub App), `user:<login>` (a named human) or `team:<org>/<slug>` (any member of the team); `gh codecrew migrate` rewrites a 1.0 table (SPEC §5)", e.Role, e.Value)
}

// SpokeRoutingError is a spoke pointer carrying a `roles:` block. The hub
// holds the one routing table for the whole topology (SPEC §5); a copy in
// a spoke is either a stale duplicate that silently outranks it or a
// promise the spoke cannot keep, and 2.0 refuses both rather than picking
// a winner. The CLI gives it the refusal code (SPOKE_ROUTING).
type SpokeRoutingError struct {
	Hub   string   // the hub this pointer names
	Roles []string // the rows it declares, in name order
}

func (e *SpokeRoutingError) Error() string {
	return fmt.Sprintf("%s names the hub %s and carries a roles: block (%s) — the hub carries the routing table for the whole project, and a spoke's copy would outrank it while going stale; if this repo is itself %s, say hub: self, because a hub's own pointer is where the table belongs; otherwise delete the block here and declare the routing in %s's %s (SPEC §5)",
		Pointer, e.Hub, strings.Join(e.Roles, ", "), e.Hub, e.Hub, Pointer)
}

// Pointer, RolesDir and AgentsFile are the protocol 2.0 layout, relative to
// a repo's root: every CodeCrew-owned operational file lives under
// .codecrew/ (SPEC §3), so nothing the framework writes collides with a
// project's own tree. AgentsFile holds the instructions a dispatched agent
// reads; the root AGENTS.md an adopter already owns only points at it.
// Spelled once, here, and used everywhere the layout is named.
const (
	Pointer    = ".codecrew/config.yml"
	RolesDir   = ".codecrew/roles"
	AgentsFile = ".codecrew/AGENTS.md"
)

// legacyPointer and legacyRolesDir are the 1.x layout: recognised only to
// refuse it. Nothing reads them.
const (
	legacyPointer  = ".codecrew.yml"
	legacyRolesDir = "roles"
)

// contractNames are the five role contracts (SPEC §7). A 1.x roles/
// directory is identified by holding at least one of them — a project's own
// roles/ (Ansible's, say) is not the layout this binary refuses.
var contractNames = []string{"implementer.md", "reviewer.md", "qa.md", "doc-synthesizer.md", "coordinator.md"}

// LegacyLayoutError reports a protocol 1.x layout found where the 2.0
// pointer should be. The CLI turns it into refused[LAYOUT_LEGACY]; there is
// no dual-read — a 1.x repo is migrated (`gh codecrew migrate`), never
// interpreted.
type LegacyLayoutError struct {
	Dir   string   // the directory the 1.x files sit in
	Found []string // their paths, relative to Dir
}

func (e *LegacyLayoutError) Error() string {
	return fmt.Sprintf("%s holds the protocol 1.x layout (%s) and no %s", e.Dir, strings.Join(e.Found, ", "), Pointer)
}

// LegacyLayout reports the 1.x files present directly in dir: the root
// pointer, and each of the five contracts sitting in a root roles/. An
// empty result means dir carries no 1.x layout. It is asked about one
// directory — a repo root — never about every level of a walk: a nested
// roles/ belongs to the project, not to CodeCrew. init calls it too, since
// it scaffolds rather than loading a pointer and its refusal cannot come
// through Load.
func LegacyLayout(dir string) []string {
	var found []string
	if isFile(filepath.Join(dir, legacyPointer)) {
		found = append(found, legacyPointer)
	}
	for _, name := range contractNames {
		if isFile(filepath.Join(dir, legacyRolesDir, name)) {
			found = append(found, legacyRolesDir+"/"+name)
		}
	}
	return found
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// Role is one entry of the routing table: which principal holds a seat,
// and the harness and model an orchestrator should dispatch it on. The
// dispatch fields are advisory — CodeCrew dispatches nothing — but the
// identity is read by the gates, so the table is fetched and checked
// (SPEC §5, §6).
type Role struct {
	Harness  string   `yaml:"harness"`
	Model    string   `yaml:"model"`
	Identity Identity `yaml:"identity"`
}

// Config is the parsed .codecrew/config.yml.
type Config struct {
	Codecrew string          `yaml:"codecrew"`
	Hub      string          `yaml:"hub"`
	Roles    map[string]Role `yaml:"roles"`

	// Dir is the directory the pointer file was found in.
	Dir string `yaml:"-"`
}

// Load walks upward from dir until it finds a .codecrew/config.yml and
// parses it. Dir is set to the directory that *contains* .codecrew/ — the
// repo root — so everything downstream keeps reading paths from the root.
//
// Only when the whole walk finds no pointer does the 1.x layout come into
// question, and then at one level: the repo root, or the starting directory
// when there is no repository. The pointer always wins, wherever it sits
// above the caller — otherwise a 2.0 repo carrying an ordinary nested
// playbook/roles/qa.md would refuse LAYOUT_LEGACY the moment a verb ran
// from playbook/, which is exactly the collision the layout move exists to
// end (checky's finding on PR #277).
func Load(dir string) (*Config, error) {
	start, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	for d := start; ; {
		path := filepath.Join(d, filepath.FromSlash(Pointer))
		if data, err := os.ReadFile(path); err == nil {
			cfg, err := Parse(data)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}
			cfg.Dir = d
			return cfg, nil
		}
		parent := filepath.Dir(d)
		if parent == d {
			break
		}
		d = parent
	}
	root := gitRoot(start)
	if legacy := LegacyLayout(root); len(legacy) > 0 {
		return nil, &LegacyLayoutError{Dir: root, Found: legacy}
	}
	return nil, fmt.Errorf("no %s found (not a CodeCrew repo?)", Pointer)
}

// gitRoot is the nearest ancestor of dir (dir itself included) holding a
// .git entry — the repository root, found the way the pointer walk finds
// the pointer, without shelling out. Outside a repository it is dir: the
// caller's own directory is the only level a 1.x layout could sensibly be
// looked for.
func gitRoot(dir string) string {
	for d := dir; ; {
		if _, err := os.Stat(filepath.Join(d, ".git")); err == nil {
			return d
		}
		parent := filepath.Dir(d)
		if parent == d {
			return dir
		}
		d = parent
	}
}

// Compatible checks a pointer's protocol version against the one the
// binary implements (SPEC §5). Same major: compatible. "0.1" — the
// pre-1.0 form of the same conventions — is compatible with a 1.x binary,
// with a note to update. A missing field is compatible with a note. Any
// other major is an error, and the two directions read differently: a
// pointer ahead of the binary is met with "upgrade the extension", one
// behind it with the migration that moves the repo forward. Neither ever
// suggests editing the version field by hand.
func Compatible(pointer, implemented string) (note string, err error) {
	implMajor := major(implemented)
	switch {
	case pointer == "":
		return fmt.Sprintf("note: %s has no codecrew: protocol version — assuming %s; add codecrew: \"%s\" (SPEC §5)", Pointer, implemented, implemented), nil
	case pointer == "0.1" && implMajor == "1":
		return fmt.Sprintf("note: %s says protocol 0.1, the pre-1.0 form of 1.0 — update it to codecrew: \"%s\" (SPEC §5)", Pointer, implemented), nil
	case major(pointer) == implMajor:
		return "", nil
	case olderMajor(pointer, implemented):
		// The repo predates this binary's protocol. Never "update the
		// pointer": the version field is a statement about the repo's
		// layout and conventions, and editing it by hand would make the
		// file lie. Migration is the act that makes it true.
		return "", fmt.Errorf("%s speaks protocol %s, which predates the protocol this codecrew implements (%s) — move the repo forward with gh codecrew migrate (SPEC §5)", Pointer, pointer, implemented)
	default:
		return "", fmt.Errorf("%s speaks protocol %s; this codecrew implements protocol %s — upgrade the extension (SPEC §5)", Pointer, pointer, implemented)
	}
}

// CompatibleHub checks the protocol version of a hub's pointer, fetched
// from a spoke, against the version this binary implements — which the
// spoke's own pointer has already been checked against, so this is what
// makes the version check topology-wide rather than repo-local (SPEC §5,
// §6). A skew names both sides, and reads in the same two directions
// Compatible does: a hub behind the binary is moved forward with
// `migrate`, a hub ahead of it wants a newer extension here. A hub with no
// codecrew: field is assumed current, with a note, exactly as a local
// pointer is.
func CompatibleHub(hub, hubVersion, implemented string) (note string, err error) {
	switch {
	case hubVersion == "":
		return fmt.Sprintf("note: the hub %s's %s has no codecrew: protocol version — assuming %s (SPEC §5)", hub, Pointer, implemented), nil
	case major(hubVersion) == major(implemented):
		return "", nil
	case olderMajor(hubVersion, implemented):
		return "", fmt.Errorf("the hub %s speaks protocol %s in its %s; this repo speaks protocol %s — one project speaks one protocol major, so move the hub forward with gh codecrew migrate (SPEC §5)", hub, hubVersion, Pointer, implemented)
	default:
		return "", fmt.Errorf("the hub %s speaks protocol %s in its %s; this repo speaks protocol %s — one project speaks one protocol major, so upgrade the extension here (SPEC §5)", hub, hubVersion, Pointer, implemented)
	}
}

func major(v string) string {
	if i := strings.Index(v, "."); i >= 0 {
		return v[:i]
	}
	return v
}

// olderMajor reports whether pointer's major is numerically below
// implemented's. An unparseable major is not treated as older: the
// unknown-version case reads better as "upgrade the extension" than as an
// instruction to migrate.
func olderMajor(pointer, implemented string) bool {
	p, err1 := strconv.Atoi(major(pointer))
	i, err2 := strconv.Atoi(major(implemented))
	return err1 == nil && err2 == nil && p < i
}

// Parse decodes pointer-file content.
func Parse(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Hub == "" {
		return nil, fmt.Errorf("missing required field: hub")
	}
	// Named in role order so a table with several 1.0 rows always refuses
	// on the same one — a refusal that moves between runs is a refusal
	// nobody can act on.
	names := make([]string, 0, len(cfg.Roles))
	for name := range cfg.Roles {
		names = append(names, name)
	}
	sort.Strings(names)
	if cfg.Hub != "self" && len(names) > 0 {
		return nil, &SpokeRoutingError{Hub: cfg.Hub, Roles: names}
	}
	for _, name := range names {
		if id := cfg.Roles[name].Identity; id.Kind == KindUntyped {
			return nil, &UntypedIdentityError{Role: name, Value: id.Value}
		}
	}
	return &cfg, nil
}

// matches reports whether login acts as the identity, by its kind: an
// `app:` slug matches the App's "<slug>[bot]" viewer login and the bare
// slug alike, a `user:` login matches exactly (a human login never
// carries the suffix). A `team:` seat never matches here — membership is
// a network question, answered by the CLI's team-aware lens — and the
// operator has no name to match.
func (i Identity) matches(login string) bool {
	switch i.Kind {
	case KindApp:
		return strings.TrimSuffix(login, "[bot]") == i.Value
	case KindUser:
		return login == i.Value
	default:
		return false
	}
}

// RoleFor resolves a viewer login to the role name whose identity it acts
// as. Team-held seats resolve through membership, which needs the API, so
// they are the CLI's business, not this table's; unmatched or empty logins
// resolve to "".
func (c *Config) RoleFor(login string) string {
	if strings.TrimSuffix(login, "[bot]") == "" {
		return ""
	}
	for name, role := range c.Roles {
		if role.Identity.matches(login) {
			return name
		}
	}
	return ""
}

// HoldsRole reports whether login holds the named role. A routed role is
// held exactly by its identity, per its kind. An unrouted role (`~`) is
// held by the human operator: any login not routed to a different role.
// Solo is a routing configuration, not a degraded tier — every role is
// always staffed (SPEC §5, decided at the gate on gh-codecrew#42).
func (c *Config) HoldsRole(login, role string) bool {
	if strings.TrimSuffix(login, "[bot]") == "" {
		return false
	}
	id := c.Roles[role].Identity
	if !id.Operator() {
		return id.matches(login)
	}
	return !strings.HasSuffix(login, "[bot]") && c.RoleFor(login) == ""
}

// HubRepo resolves the hub to an owner/repo string. current is the
// owner/repo of the repository the command runs in, used when hub is "self".
func (c *Config) HubRepo(current string) string {
	if c.Hub == "self" {
		return current
	}
	return c.Hub
}
