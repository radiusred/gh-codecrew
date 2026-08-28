package gh

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var versionLine = regexp.MustCompile(`gh version (\d+)\.(\d+)\.(\d+)`)

// Version returns the installed gh's release, e.g. "2.46.0", parsed from
// the first line of `gh --version`.
func Version() (string, error) {
	out, err := Run("--version")
	if err != nil {
		return "", err
	}
	return ParseVersionLine(string(out))
}

// ParseVersionLine extracts x.y.z from gh's version banner
// ("gh version 2.46.0 (2024-03-20)\nhttps://…"). A banner that does not
// carry one is an error the caller reports, not a version.
func ParseVersionLine(s string) (string, error) {
	m := versionLine.FindStringSubmatch(s)
	if m == nil {
		return "", fmt.Errorf("gh --version: unrecognised output %q", strings.TrimSpace(strings.SplitN(s, "\n", 2)[0]))
	}
	return m[1] + "." + m[2] + "." + m[3], nil
}

// CompareVersions orders two x.y.z strings numerically: -1, 0 or 1. Parts
// beyond the third and any suffix are ignored; a missing part is 0.
func CompareVersions(a, b string) int {
	pa, pb := parts(a), parts(b)
	for i := 0; i < 3; i++ {
		if pa[i] != pb[i] {
			if pa[i] < pb[i] {
				return -1
			}
			return 1
		}
	}
	return 0
}

func parts(v string) [3]int {
	var out [3]int
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	for i, p := range strings.SplitN(v, ".", 3) {
		if i > 2 {
			break
		}
		n, _ := strconv.Atoi(strings.TrimFunc(p, func(r rune) bool { return r < '0' || r > '9' }))
		out[i] = n
	}
	return out
}
