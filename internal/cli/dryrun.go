package cli

import (
	"fmt"
	"io"
	"strings"
)

// plan is what a verb would do: its gates in the order the live verb
// checks them, each with its outcome, then the actions a clean pass takes.
// One code path builds it for both modes — --dry-run prints it and writes
// nothing, the live verb executes it — so a preview can never disagree
// with the run (M7-R5, #133).
type plan struct {
	gates   []gateResult
	actions []string
	refusal error // the first refused gate — the dry run's exit status too
}

type gateStatus int

const (
	gateOK gateStatus = iota
	gateRefused
	gateUnreached // the live verb stops at a refusal; later gates are never asked
	gateNA        // not applicable on this path (a flag not given, a role unrouted)
)

type gateResult struct {
	name   string
	status gateStatus
	err    error
}

// gate records a gate's outcome and reports whether the verb continues.
func (p *plan) gate(name string, err error) bool {
	if err != nil {
		p.gates = append(p.gates, gateResult{name: name, status: gateRefused, err: err})
		if p.refusal == nil {
			p.refusal = err
		}
		return false
	}
	p.gates = append(p.gates, gateResult{name: name, status: gateOK})
	return true
}

func (p *plan) na(name string) { p.gates = append(p.gates, gateResult{name: name, status: gateNA}) }

// stop marks every gate in order after the last recorded one as unreached
// and returns the plan — the shape of a refusal.
func (p *plan) stop(order []string) *plan {
	last := ""
	if len(p.gates) > 0 {
		last = p.gates[len(p.gates)-1].name
	}
	after := last == ""
	for _, n := range order {
		if after {
			p.gates = append(p.gates, gateResult{name: n, status: gateUnreached})
		}
		if n == last {
			after = true
		}
	}
	return p
}

func (p *plan) would(format string, args ...any) {
	p.actions = append(p.actions, fmt.Sprintf(format, args...))
}

func (p *plan) print(w io.Writer) {
	for _, g := range p.gates {
		switch g.status {
		case gateOK:
			fmt.Fprintf(w, "gate %s: ok\n", g.name)
		case gateRefused:
			fmt.Fprintf(w, "gate %s: %v\n", g.name, g.err)
		case gateUnreached:
			fmt.Fprintf(w, "gate %s: not reached\n", g.name)
		case gateNA:
			fmt.Fprintf(w, "gate %s: not applicable\n", g.name)
		}
	}
	if p.refusal != nil {
		fmt.Fprintln(w, "dry run: nothing written — the live verb stops at the first refusal above")
		return
	}
	for _, a := range p.actions {
		fmt.Fprintf(w, "would %s\n", a)
	}
	fmt.Fprintln(w, "dry run: nothing written")
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}
