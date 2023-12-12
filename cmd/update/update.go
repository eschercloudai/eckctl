package update

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	controlPlaneName   string
	clusterName        string
	clusterDefPath     string
	insecure           bool
	url, u, p, project string
)

func NewUpdateCommand() *cobra.Command {
	updateCmd := &cobra.Command{
		Use: "update",
	}

	commands := []*cobra.Command{
		updateClusterCmd,
	}

	updateCmd.AddCommand(commands...)

	updateClusterCmd.Flags().StringVar(&clusterName, "name", "", "Name of cluster")
	updateClusterCmd.Flags().StringVar(&controlPlaneName, "controlplane", "default", "Name of associated control plane")
	updateClusterCmd.Flags().StringVar(&clusterDefPath, "json", "", "Path to JSON cluster definition")
	if err := updateClusterCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalln(err)
	}
	if err := updateClusterCmd.MarkFlagRequired("json"); err != nil {
		log.Fatalln(err)
	}

	return updateCmd
}
