package tracker

import (
	"reflect"
	"testing"
)

const decisionURL = "https://github.com/o/r/issues/5#issuecomment-111"

// The struck line is its own grammar beside the verdicts': the ID and the
// word inside the bold, the decision link after it. A QA comment that
// writes "struck" is neither a verdict nor — unless its author holds the
// coordinator seat, which callers decide — a counted strike (M18-R4).
func TestParseStrikes(t *testing.T) {
	comments := []Comment{
		{Author: "davison", URL: "https://github.com/o/r/issues/5#issuecomment-1", Body: "**M8-R2 — struck.** " + decisionURL},
		{Author: "davison", URL: "https://github.com/o/r/issues/5#issuecomment-2", Body: "Reason first.\n\n**M8-R3 - Struck** see [the decision](" + decisionURL + ").\r\n**M8-R3 — reinstated.** a second line for the same ID in one comment does not count"},
		{Author: "davison", URL: "https://github.com/o/r/issues/5#issuecomment-3", Body: "Quoted: `**M8-R4 — struck.** " + decisionURL + "`\n\n```\n**M8-R4 — struck.** x\n```"},
		{Author: "davison", URL: "https://github.com/o/r/issues/5#issuecomment-4", Body: "**M8-R5 — struck.** no link at all"},
		{Author: "davison", URL: "https://github.com/o/r/issues/5#issuecomment-5", Body: "**M8-R2 — reinstated.** " + decisionURL},
	}
	got := ParseStrikes(comments)
	want := []Strike{
		{ID: "M8-R2", Word: Struck, Decision: decisionURL, Author: "davison", URL: "https://github.com/o/r/issues/5#issuecomment-1"},
		{ID: "M8-R3", Word: Struck, Decision: decisionURL, Author: "davison", URL: "https://github.com/o/r/issues/5#issuecomment-2"},
		{ID: "M8-R5", Word: Struck, Decision: "", Author: "davison", URL: "https://github.com/o/r/issues/5#issuecomment-4"},
		{ID: "M8-R2", Word: Reinstated, Decision: decisionURL, Author: "davison", URL: "https://github.com/o/r/issues/5#issuecomment-5"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseStrikes =\n%+v\nwant\n%+v", got, want)
	}
	if v := ParseVerdicts(comments); len(v) != 0 {
		t.Errorf("a struck line parsed as a QA verdict: %+v", v)
	}
	if s := ParseStrikes([]Comment{{Body: "**M8-R1 — satisfied.** ran it"}}); len(s) != 0 {
		t.Errorf("a verdict parsed as a strike: %+v", s)
	}
}

func TestStrikeLineRoundTrips(t *testing.T) {
	line := StrikeLine("M8-R2", Struck, decisionURL)
	if line != "**M8-R2 — struck.** "+decisionURL {
		t.Fatalf("StrikeLine = %q", line)
	}
	got := ParseStrikes([]Comment{{Body: line}})
	if len(got) != 1 || got[0].ID != "M8-R2" || got[0].Word != Struck || got[0].Decision != decisionURL {
		t.Errorf("the posted line does not read back: %+v", got)
	}
}

func TestParseCommentURL(t *testing.T) {
	ref, id, ok := ParseCommentURL(" https://github.com/Owner/re.po/issues/42#issuecomment-987 ")
	if !ok || ref != (IssueRef{Repo: "Owner/re.po", Number: 42}) || id != "987" {
		t.Errorf("got %v %q %v", ref, id, ok)
	}
	for _, bad := range []string{
		"https://github.com/o/r/issues/42",
		"https://github.com/o/r/pull/42#issuecomment-1",
		"https://github.com/o/r/pull/42#discussion_r1",
		"o/r#42",
		"see https://github.com/o/r/issues/42#issuecomment-1",
	} {
		if _, _, ok := ParseCommentURL(bad); ok {
			t.Errorf("%q parsed as a comment URL", bad)
		}
	}
	if CommentID("https://github.com/o/r/issues/5#issuecomment-111") != "111" || CommentID("https://github.com/o/r/issues/5") != "" {
		t.Error("CommentID")
	}
}

func TestNamesRequirement(t *testing.T) {
	for text, want := range map[string]bool{
		"strike M8-R2.":           true,
		"M8-R2 is out of scope":   true,
		"(**M8-R2**)":             true,
		"strike M8-R21":           false,
		"strike XM8-R2":           false,
		"strike M18-R2":           false,
		"nothing about it at all": false,
	} {
		if got := NamesRequirement(text, "M8-R2"); got != want {
			t.Errorf("NamesRequirement(%q) = %v, want %v", text, got, want)
		}
	}
}

// A still-bold ID struck through in the body is still a requirement — the
// body's strikethrough is not a record — and StruckThroughIDs lists it
// only so status can say so. md-notes' un-bolded `~~M8-R2~~` stays what it
// always was: no requirement at all (#348).
func TestStillBoldStruckThroughIDIsStillARequirement(t *testing.T) {
	body := "## Goal\n~~**M8-R9**~~ is not in the section\n\n## Requirements\r\n- **M8-R1** — one\n- ~~**M8-R2** — two~~\n- ~~M8-R3 — three, un-bolded~~\n- ~~ **M8-R4**~~ — four\n\n## Gates\nx\n"
	if got := RequirementIDs(body); !reflect.DeepEqual(got, []string{"M8-R1", "M8-R2", "M8-R4"}) {
		t.Errorf("RequirementIDs = %v", got)
	}
	if got := StruckThroughIDs(body); !reflect.DeepEqual(got, []string{"M8-R2", "M8-R4"}) {
		t.Errorf("StruckThroughIDs = %v", got)
	}
}
