// Package cli dispatches the codecrew workflow verbs (SPEC.md §6).
package cli

import (
	"fmt"
	"os"
)

const usage = `usage: codecrew <verb>

verbs:
  init [--hub owner/repo]                    scaffold a new hub (or spoke pointer)
  status                                     where the project is
  milestone new --title T [--goal G]         create a milestone tracking issue
  milestone close <n>                        close a milestone (gates: tasks closed, doc merged)
  milestone evidence <n>                     verify every cited link in the milestone's record resolves
  task new --milestone N --title T           create a task issue, linked into the milestone
           [--repo owner/repo] [--goal G] [--requirements IDs]
  task start <ref>                           assign, verify plan, create linked branch
  task finish <ref> [--operator-confirm]     enforce gates, then rebase-merge
           [--bypass]                        (recorded admin merge when GitHub won't count the approval)
  checkpoint <ref> --question "..."          raise a human gate (cc:needs-decision)
  role <name>                                who holds a role (identity, or ~ for the operator)
  roles diff <role>                          local contract vs the one embedded in the CLI
  roles show <role> [--latest]               the contract a session loads (with local extensions), or the embedded one
  identity new <role> --name N               mint the role's App identity via the manifest flow
           [--owner O] [--with-webhook --webhook-url U]
           [--with-approval-permission]      (reviewer only: its approvals satisfy required reviews)
  version                                    installed release tag (dev for source builds)

Blocked gates exit nonzero with "refused[CODE]: detail".
`

// Run executes one verb.
func Run(args []string) error {
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
		if len(rest) == 0 || rest[0] != "new" {
			fmt.Fprint(os.Stderr, usage)
			return fmt.Errorf("identity: unknown subcommand")
		}
		return identityNew(os.Stdout, rest[1:])
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
