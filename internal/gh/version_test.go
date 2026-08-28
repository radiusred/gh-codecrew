package gh

import "testing"

func TestParseVersionLine(t *testing.T) {
	for in, want := range map[string]string{
		"gh version 2.46.0 (2024-03-20)\nhttps://github.com/cli/cli/releases/tag/v2.46.0\n": "2.46.0",
		"gh version 2.98.0 (2026-08-01)":         "2.98.0",
		"gh version 2.50.0-preview (2024-05-29)": "2.50.0",
	} {
		got, err := ParseVersionLine(in)
		if err != nil || got != want {
			t.Errorf("%q = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, in := range []string{"", "gh: command not found", "version 2.46"} {
		if got, err := ParseVersionLine(in); err == nil {
			t.Errorf("%q parsed as %q; want an error", in, got)
		}
	}
}

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"2.46.0", "2.50.0", -1},
		{"2.50.0", "2.50.0", 0},
		{"2.98.0", "2.50.0", 1},
		{"3.0.0", "2.99.9", 1},
		{"2.50", "2.50.0", 0},
		{"v2.51.1", "2.50.0", 1},
		{"2.50.0-preview", "2.50.0", 0},
		{"2.5.0", "2.50.0", -1},
	}
	for _, c := range cases {
		if got := CompareVersions(c.a, c.b); got != c.want {
			t.Errorf("CompareVersions(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}
