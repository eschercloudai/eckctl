package get

import (
	"github.com/spf13/cobra"
)

func NewGetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use: "get",
	}

	commands := []*cobra.Command{
		imagesCmd,
		networksCmd,
		controlPlaneCmd(),
		versionsCmd,
		clustersCmd,
		projectsCmd,
		kubeconfigCmd,
	}

	getCmd.AddCommand(commands...)

	clustersCmd.Flags().StringVar(&controlPlaneName, "controlplane", "", "Name of control plane")
	kubeconfigCmd.Flags().StringVar(&controlPlaneName, "controlplane", "", "Name of control plane")
	kubeconfigCmd.Flags().StringVar(&clusterName, "cluster", "", "Name of cluster")

	return getCmd

}
