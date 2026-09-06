package cli

import (
	"bytes"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/config"
)

func TestHolder(t *testing.T) {
	table := map[string]config.Role{
		"implementer": {Identity: config.ParseIdentity("app:radiusred-cody")},
		"reviewer":    {Identity: config.ParseIdentity("user:davison")},
		"qa":          {}, // identity: ~ — operator-held
	}
	cases := []struct {
		name    string
		roles   map[string]config.Role
		role    string
		want    string
		wantErr bool
	}{
		{"routed app", table, "implementer", "app:radiusred-cody", false},
		{"routed human", table, "reviewer", "user:davison", false},
		{"explicitly operator-held", table, "qa", "~", false},
		{"absent from a declared table", table, "doc-synthesizer", "", true}, // the identity is meaningless beside the error
		{"no table declared: any role is operator-held", nil, "reviewer", "~", false},
	}
	for _, c := range cases {
		got, err := holder(c.roles, c.role)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: err = %v, wantErr %v", c.name, err, c.wantErr)
			continue
		}
		if c.wantErr {
			continue
		}
		if got.String() != c.want {
			t.Errorf("%s: holder = %q, want %q", c.name, got, c.want)
		}
	}
}

// Until 2.0 a declared table with no coordinator row resolved the seat to
// ~, because the row arrived after 1.0 hubs scaffolded their tables. The
// shim is gone (M13-R7): the row is scaffolded by init and written into a
// 1.x table by migrate, so its absence from a declared table is the same
// error every other missing role is, and nothing infers a holder for a
// seat the table does not name.
func TestCoordinatorAbsentFromDeclaredTableRefuses(t *testing.T) {
	table := map[string]config.Role{
		"implementer":     {Identity: config.ParseIdentity("app:myorg-coder")},
		"reviewer":        {Identity: config.ParseIdentity("app:myorg-reviewy")},
		"qa":              {},
		"doc-synthesizer": {},
	}
	if got, err := holder(table, "coordinator"); err == nil {
		t.Errorf("coordinator absent from a declared table = %q, want an error", got)
	}
	table["coordinator"] = config.Role{Identity: config.ParseIdentity("app:myorg-loopy")}
	if got, _ := holder(table, "coordinator"); got.String() != "app:myorg-loopy" {
		t.Errorf("routed coordinator = %q", got)
	}
	// A table declaring nothing at all is untouched by the removal: every
	// seat, coordinator included, is the operator's.
	if got, err := holder(nil, "coordinator"); err != nil || got.String() != "~" {
		t.Errorf("coordinator with no table declared = %q, %v; want ~", got, err)
	}
	if _, err := holder(table, "navigator"); err == nil {
		t.Error("an unknown role still resolves")
	}
}

// `role <name>` prints the typed value; `--login` prints the handle a
// review request can name, and nothing at all for a seat that cannot be
// requested — the emptiness is the implementer contract's whole branch
// (M13-R4, Decision on #258).
func TestRoleHolderOutput(t *testing.T) {
	cases := []struct {
		identity   string
		typed      string
		loginPrint string
	}{
		{"app:radiusred-checky", "app:radiusred-checky\n", ""},
		{"user:alice", "user:alice\n", "alice\n"},
		{"team:myorg/review-crew", "team:myorg/review-crew\n", "myorg/review-crew\n"},
		{"~", "~\n", ""},
	}
	for _, c := range cases {
		id := config.ParseIdentity(c.identity)
		var typed, login bytes.Buffer
		printHolder(&typed, id, false)
		printHolder(&login, id, true)
		if typed.String() != c.typed || login.String() != c.loginPrint {
			t.Errorf("%s: role = %q (want %q), role --login = %q (want %q)",
				c.identity, typed.String(), c.typed, login.String(), c.loginPrint)
		}
	}
}

// --login after the role name is the shape the contract types; Go's flag
// package stops at the first positional, so it is handled explicitly.
func TestParseRoleArgs(t *testing.T) {
	cases := []struct {
		args    []string
		name    string
		login   bool
		wantErr bool
	}{
		{[]string{"reviewer"}, "reviewer", false, false},
		{[]string{"reviewer", "--login"}, "reviewer", true, false},
		{[]string{"--login", "reviewer"}, "reviewer", true, false},
		{[]string{"reviewer", "-login"}, "reviewer", true, false},
		{nil, "", false, true},
		{[]string{"--login"}, "", false, true},
		{[]string{"reviewer", "qa"}, "", false, true},
		{[]string{"reviewer", "--nonesuch"}, "", false, true},
	}
	for _, c := range cases {
		name, login, err := parseRoleArgs(c.args)
		if (err != nil) != c.wantErr || name != c.name || login != c.login {
			t.Errorf("parseRoleArgs(%q) = %q, %v, %v", c.args, name, login, err)
		}
	}
}
