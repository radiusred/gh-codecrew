// Package config locates and parses the .codecrew.yml pointer file that
// every repo in a CodeCrew project carries (SPEC.md §3, §5).
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

// Role is one entry of the advisory role routing table.
type Role struct {
	Harness  string   `yaml:"harness"`
	Model    string   `yaml:"model"`
	Identity Identity `yaml:"identity"`
}

// Config is the parsed .codecrew.yml.
type Config struct {
	Codecrew string          `yaml:"codecrew"`
	Hub      string          `yaml:"hub"`
	Roles    map[string]Role `yaml:"roles"`

	// Dir is the directory the pointer file was found in.
	Dir string `yaml:"-"`
}

// Load walks upward from dir until it finds a .codecrew.yml and parses it.
func Load(dir string) (*Config, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	for {
		path := filepath.Join(dir, ".codecrew.yml")
		if data, err := os.ReadFile(path); err == nil {
			cfg, err := Parse(data)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}
			cfg.Dir = dir
			return cfg, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil, fmt.Errorf("no .codecrew.yml found (not a CodeCrew repo?)")
		}
		dir = parent
	}
}

// Compatible checks a pointer's protocol version against the one the
// binary implements (SPEC §5). Same major: compatible. "0.1" — the
// pre-1.0 form of the same conventions — is compatible with a 1.x binary,
// with a note to update. A missing field is compatible with a note. Any
// other major is an error: the pointer speaks conventions this binary does
// not.
func Compatible(pointer, implemented string) (note string, err error) {
	implMajor := major(implemented)
	switch {
	case pointer == "":
		return fmt.Sprintf("note: .codecrew.yml has no codecrew: protocol version — assuming %s; add codecrew: \"%s\" (SPEC §5)", implemented, implemented), nil
	case pointer == "0.1" && implMajor == "1":
		return fmt.Sprintf("note: .codecrew.yml says protocol 0.1, the pre-1.0 form of 1.0 — update it to codecrew: \"%s\" (SPEC §5)", implemented), nil
	case major(pointer) == implMajor:
		return "", nil
	default:
		return "", fmt.Errorf(".codecrew.yml speaks protocol %s; this codecrew implements protocol %s — upgrade the extension or the pointer (SPEC §5)", pointer, implemented)
	}
}

func major(v string) string {
	if i := strings.Index(v, "."); i >= 0 {
		return v[:i]
	}
	return v
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
