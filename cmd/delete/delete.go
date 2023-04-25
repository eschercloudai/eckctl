package delete

import (
	"github.com/spf13/cobra"
)

var (
	clusterName      string
	controlPlaneName string
)

func NewDeleteCommand() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use: "delete",
	}

	commands := []*cobra.Command{
		deleteControlPlaneCmd(),
		deleteApplicationCredentialCmd(),
		deleteClusterCmd(),
	}

	deleteCmd.AddCommand(commands...)

	return deleteCmd
}
