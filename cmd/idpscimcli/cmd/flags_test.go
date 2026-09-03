//go:build unit

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

// The --filter flag is only useful on the leaf commands that actually issue a
// request. awsGroupsListCmd registered it correctly on the `list` subcommand,
// but the users equivalent registered it on awsUsersCmd — the `users` grouping
// command, which has no RunE.
//
// Cobra's Flags() are local to the command they are set on, so the flag ended up
// advertised by `idpscimcli aws users --help`, where nothing consumes it, and
// rejected by `idpscimcli aws users list --filter …` with "unknown flag".
func TestAWSListCommands_acceptFilterFlag(t *testing.T) {
	for _, tc := range []struct {
		name string
		cmd  *cobra.Command
	}{
		{"aws groups list", awsGroupsListCmd},
		{"aws users list", awsUsersListCmd},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := tc.cmd.Flags().Lookup("filter")
			if f == nil {
				t.Fatalf("%s has no --filter flag; it cannot filter results", tc.name)
			}
			if f.Shorthand != "q" {
				t.Errorf("%s --filter shorthand = %q, want %q", tc.name, f.Shorthand, "q")
			}
		})
	}
}

// The grouping commands have no RunE, so a local flag on them is unreachable.
// Guard against the flag drifting back up the tree.
func TestAWSGroupingCommands_haveNoFilterFlag(t *testing.T) {
	for _, tc := range []struct {
		name string
		cmd  *cobra.Command
	}{
		{"aws groups", awsGroupsCmd},
		{"aws users", awsUsersCmd},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cmd.RunE != nil || tc.cmd.Run != nil {
				t.Skipf("%s now has a Run function; a local flag would be reachable", tc.name)
			}
			if f := tc.cmd.Flags().Lookup("filter"); f != nil {
				t.Errorf("%s advertises --filter but has no Run function, so nothing consumes it", tc.name)
			}
		})
	}
}
