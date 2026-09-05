// Package cli dispatches the codecrew workflow verbs (SPEC.md §6).
package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const usage = `usage: gh codecrew <verb>

verbs:
  init [--hub owner/repo]                    scaffold a new hub (or spoke pointer)
  status                                     where the project is
  milestone new --title T [--goal G]         create a milestone tracking issue
           [--requirement R]...              (repeatable: numbered M<n>-R1, R2, … in order, under ## Requirements)
           [--dry-run]                       (print the number, title and requirement IDs it would get; create nothing)
  milestone close <milestone number>         close a milestone (gates: milestone open, no gate raised, tasks closed, requirements,
                                             QA verdicts, doc merged)
           [--dry-run]                       (print every gate and the sweep; write nothing)
  milestone evidence <milestone number>      verify the record's citations resolve: a dead github.com link refuses, a dead external
                                             link warns; URLs inside code are content, not citations (M2 → 2)
  task new --milestone N --title T           create a task issue, linked into the milestone
           [--repo owner/repo] [--goal G] [--requirements IDs]
  task start <ref>                           assign, verify plan, create linked branch
  task finish <ref> [--operator-confirm]     enforce gates, then rebase-merge
           [--bypass]                        (recorded admin merge when GitHub won't count the approval)
           [--dry-run]                       (print every gate and what it would do; write nothing)
  checkpoint <ref> --question "..."          raise a human gate (cc:needs-decision)
  role <name>                                who holds a role (identity, or ~ for the operator)
  roles diff <role>                          local contract vs the one embedded in the CLI
  roles show <role> [--latest]               the contract a session loads (with local extensions), or the embedded one
  identity new <role> --name N               mint the role's App identity via the manifest flow
           [--owner O] [--with-webhook --webhook-url U]
           [--with-approval-permission]      (reviewer only: its approvals satisfy required reviews)
           [--events E,E] [--webhook-secret S]  (with --with-webhook: default pull_request,pull_request_review;
                                             set the receiver's secret right after creation)
  identity webhook <slug> [--show]           the App's hook: print it, --url U / --secret S to set,
           [--url U] [--secret S|--rotate-secret]  --rotate-secret to mint one (events change on the settings page)
  identity token [<slug>]                    mint an installation token: env bindings first
           [--installation ID]               (GITHUB_APP_ID|GITHUB_CLIENT_ID + GITHUB_PRIVATE_KEY|GITHUB_PEM),
                                             else ~/.config/codecrew/<slug>; the token alone on stdout
  version                                    installed release tag (dev for source builds)

Blocked gates exit nonzero with "refused[CODE]: detail".
`

// Run executes one verb. --help on any verb prints its usage and is not
// a failure: flag.ErrHelp maps to a clean exit.
func Run(args []string) error {
	// Asking for help never runs the verb: a --help anywhere in the
	// arguments prints usage and exits 0 — before any gate, refusal or
	// side effect could be mistaken for it.
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Fprint(os.Stderr, usage)
			return nil
		}
	}
	err := run(args)
	if errors.Is(err, flag.ErrHelp) {
		return nil
	}
	return err
}

func run(args []string) error {
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, usage)
		return fmt.Errorf("no verb given")
	}
	verb, rest := args[0], args[1:]
	switch verb {
	case "init":
		return initCmd(os.Stdout, rest)
	case "status":
		return status(os.Stdout)
	case "identity":
		if len(rest) == 0 {
			fmt.Fprint(os.Stderr, usage)
			return fmt.Errorf("identity: missing subcommand")
		}
		switch rest[0] {
		case "new":
			return identityNew(os.Stdout, rest[1:])
		case "token":
			return identityToken(os.Stdout, rest[1:])
		case "webhook":
			return identityWebhook(os.Stdout, rest[1:])
		default:
			fmt.Fprint(os.Stderr, usage)
			return fmt.Errorf("identity: unknown subcommand %q", rest[0])
		}
	case "milestone", "task":
		if len(rest) == 0 {
			fmt.Fprint(os.Stderr, usage)
			return fmt.Errorf("%s: missing subcommand", verb)
		}
		switch verb + " " + rest[0] {
		case "milestone new":
			return milestoneNew(os.Stdout, rest[1:])
		case "milestone close":
			return milestoneClose(os.Stdout, rest[1:])
		case "milestone evidence":
			return milestoneEvidence(os.Stdout, rest[1:])
		case "task new":
			return taskNew(os.Stdout, rest[1:])
		case "task start":
			return taskStart(os.Stdout, rest[1:])
		case "task finish":
			return taskFinish(os.Stdout, rest[1:])
		default:
			fmt.Fprint(os.Stderr, usage)
			return fmt.Errorf("unknown subcommand %q", verb+" "+rest[0])
		}
	case "checkpoint":
		return checkpoint(os.Stdout, rest)
	case "role":
		return roleHolder(os.Stdout, rest)
	case "roles":
		return rolesCmd(os.Stdout, rest)
	case "version", "--version", "-v":
		return versionCmd(os.Stdout)
	case "help", "-h", "--help":
		fmt.Print(usage)
		return nil
	default:
		fmt.Fprint(os.Stderr, usage)
		return fmt.Errorf("unknown verb %q", verb)
	}
}
