package delete

import (
	"github.com/spf13/cobra"
)

var (
	url, u, p, project string
	insecure           bool
)

func NewDeleteCommand() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use: "delete",
	}

	commands := []*cobra.Command{
		deleteControlPlaneCmd(),
		deleteClusterCmd(),
	}

	deleteCmd.AddCommand(commands...)
	return deleteCmd
}
