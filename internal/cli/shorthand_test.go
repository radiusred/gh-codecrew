package cli

import (
	"bytes"
	"io/fs"
	"os"
	"regexp"
	"strings"
	"testing"

	codecrew "github.com/radiusred/gh-codecrew"
)

// bareVerb matches the docs' shorthand for the CLI followed by a verb. The
// installed command is `gh codecrew`; the bare form is only ever runnable
// from a source checkout.
var bareVerb = regexp.MustCompile(`codecrew (init|status|milestone|task|checkpoint|role|roles|identity|version|help)\b`)

// bareShorthand returns every `codecrew <verb>` in text that is not written
// as `gh codecrew <verb>`.
func bareShorthand(text string) []string {
	var hits []string
	for _, m := range bareVerb.FindAllStringIndex(text, -1) {
		if !strings.HasSuffix(text[:m[0]], "gh ") {
			hits = append(hits, text[m[0]:m[1]])
		}
	}
	return hits
}

// TestNoBareCodecrewShorthand: the contracts, the scaffold and the CLI's own
// output are loaded verbatim as agent instructions, so every command they
// name must run as written. A dispatched qa agent ran its contract's literal
// `codecrew milestone evidence` into "command not found" and blocked on it
// twice (#146; #119 finding 31).
func TestNoBareCodecrewShorthand(t *testing.T) {
	texts := map[string]string{
		"usage":             usage,
		"hubConfigScaffold": hubConfigScaffold,
		"agentsScaffold":    agentsScaffold,
		"claudeScaffold":    claudeScaffold,
		"roadmapScaffold":   roadmapScaffold,
	}
	for _, role := range []string{"implementer", "reviewer", "qa", "doc-synthesizer"} {
		data, err := fs.ReadFile(codecrew.Roles, "roles/"+role+".md")
		if err != nil {
			t.Fatal(err)
		}
		texts["roles/"+role+".md"] = string(data)
	}

	// init's printed next steps, hub and spoke.
	dir := t.TempDir()
	wd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)
	for _, hub := range []string{"self", "org/hub"} {
		var out bytes.Buffer
		if err := initCmd(&out, []string{"--hub", hub}); err != nil {
			t.Fatal(err)
		}
		texts["init --hub "+hub] = out.String()
	}

	// Per-verb usage errors.
	for _, args := range [][]string{
		{"task", "start"}, {"task", "finish"}, {"checkpoint"},
		{"milestone", "close"}, {"milestone", "evidence"},
		{"roles", "diff"}, {"role"}, {"identity", "new"},
	} {
		err := run(args)
		if err == nil {
			t.Errorf("%v: expected a usage error", args)
			continue
		}
		texts[strings.Join(args, " ")+" (no args)"] = err.Error()
	}

	for name, text := range texts {
		if hits := bareShorthand(text); len(hits) > 0 {
			t.Errorf("%s names the CLI without gh: %v", name, hits)
		}
	}
}

func TestBareShorthandMatcher(t *testing.T) {
	if hits := bareShorthand("run `gh codecrew status`, then gh codecrew task start 3"); len(hits) != 0 {
		t.Errorf("gh-prefixed forms flagged: %v", hits)
	}
	if hits := bareShorthand("`codecrew-token slug`, .codecrew.yml, gh-codecrew"); len(hits) != 0 {
		t.Errorf("non-verb forms flagged: %v", hits)
	}
	if hits := bareShorthand("first act: `codecrew milestone evidence 1`"); len(hits) != 1 {
		t.Errorf("bare form not flagged: %v", hits)
	}
}
