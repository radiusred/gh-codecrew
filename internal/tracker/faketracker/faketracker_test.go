package faketracker

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/radiusred/gh-codecrew/internal/tracker"
)

// The fake is a whole venue: every method of the interface is on it. The
// compile-time assertion in the package says the same thing, but only this
// counts them, so a method added to the interface and forgotten here fails
// with a number rather than a type error a reader has to decode.
func TestFakeImplementsEveryVenueMethod(t *testing.T) {
	iface := reflect.TypeOf((*tracker.Tracker)(nil)).Elem()
	fake := reflect.TypeOf(&Venue{})
	for i := 0; i < iface.NumMethod(); i++ {
		name := iface.Method(i).Name
		m, ok := fake.MethodByName(name)
		if !ok {
			t.Errorf("the fake is missing %s", name)
			continue
		}
		// And one scripting field per method, named <Method>Fn.
		if _, ok := reflect.TypeOf(Venue{}).FieldByName(name + "Fn"); !ok {
			t.Errorf("%s has no %sFn field to script it with", name, name)
		}
		_ = m
	}
	if iface.NumMethod() == 0 {
		t.Fatal("no venue methods found")
	}
}

// The zero value answers: a test scripts only what it exercises, and every
// other method returns the zero value of its results with no error rather
// than panicking half way through the verb under test.
func TestUnscriptedMethodsAnswerWithZeroValues(t *testing.T) {
	v := &Venue{}
	if ms, err := v.OpenMilestones("o/hub"); ms != nil || err != nil {
		t.Errorf("OpenMilestones = %v, %v", ms, err)
	}
	if task, err := v.Task(tracker.IssueRef{Repo: "o/r", Number: 7}); task.Title != "" || err != nil {
		t.Errorf("Task = %+v, %v", task, err)
	}
	if err := v.Comment(tracker.IssueRef{Repo: "o/r", Number: 7}, "hello"); err != nil {
		t.Errorf("Comment = %v", err)
	}
	if login, err := v.Viewer(); login != "" || err != nil {
		t.Errorf("Viewer = %q, %v", login, err)
	}
	if status, data, err := v.AppRequest(http.DefaultClient, "jwt", "GET", "/app", nil); status != 0 || data != nil || err != nil {
		t.Errorf("AppRequest = %d, %q, %v", status, data, err)
	}
}

// Scripted returns come back exactly, errors included.
func TestScriptedReturns(t *testing.T) {
	boom := errors.New("403")
	v := &Venue{
		ViewerFn:      func() (string, error) { return "radiusred-cody[bot]", nil },
		TeamMembersFn: func(org, team string) ([]string, error) { return []string{"alice"}, nil },
		LabelsFn:      func(string) ([]tracker.Label, error) { return nil, boom },
	}
	if login, err := v.Viewer(); login != "radiusred-cody[bot]" || err != nil {
		t.Errorf("Viewer = %q, %v", login, err)
	}
	if members, err := v.TeamMembers("myorg", "review-crew"); err != nil || len(members) != 1 || members[0] != "alice" {
		t.Errorf("TeamMembers = %v, %v", members, err)
	}
	if _, err := v.Labels("o/r"); !errors.Is(err, boom) {
		t.Errorf("Labels error = %v, want the scripted one", err)
	}
}

// Every call is recorded in order, with its arguments — scripted or not.
func TestRecordsEveryCallInOrder(t *testing.T) {
	ref := tracker.IssueRef{Repo: "o/r", Number: 7}
	v := &Venue{ViewerFn: func() (string, error) { return "alice", nil }}
	v.Viewer()
	v.AddLabel(ref, tracker.LabelNeedsDecision)
	v.Comment(ref, "the body")
	v.AddLabel(ref, tracker.LabelTask)

	calls := v.Calls()
	if len(calls) != 4 {
		t.Fatalf("recorded %d calls: %v", len(calls), calls)
	}
	if calls[0].Method != "Viewer" || len(calls[0].Args) != 0 {
		t.Errorf("first call = %v", calls[0])
	}
	if got := calls[2].String(); got != "Comment(o/r#7, the body)" {
		t.Errorf("Comment renders as %q", got)
	}
	if n := v.Count("AddLabel"); n != 2 {
		t.Errorf("Count(AddLabel) = %d, want 2", n)
	}
	labels := v.CallsTo("AddLabel")
	if labels[0].Args[1] != tracker.LabelNeedsDecision || labels[1].Args[1] != tracker.LabelTask {
		t.Errorf("AddLabel arguments = %v", labels)
	}
	if v.Count("MergePR") != 0 {
		t.Error("a method never called was counted")
	}

	// Calls returns a copy: a caller cannot edit the record.
	calls[0].Method = "Tampered"
	if v.Calls()[0].Method != "Viewer" {
		t.Error("Calls handed out the fake's own slice")
	}

	v.Reset()
	if len(v.Calls()) != 0 {
		t.Error("Reset left calls behind")
	}
	if login, _ := v.Viewer(); login != "alice" {
		t.Error("Reset dropped the scripted answers")
	}
}
