package get

import (
	"github.com/spf13/cobra"
)

var (
	controlPlaneName string
	clusterName      string
)

type Images struct {
	Created  string `json:"created"`
	Id       string `json:"id"`
	Versions struct {
		Kubernetes   string `json:"kubernetes"`
		NvidiaDriver string `json:"nvidiaDriver"`
	} `json:"versions"`
}

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
	clustersCmd.MarkFlagRequired("controlplane")
	kubeconfigCmd.Flags().StringVar(&controlPlaneName, "controlplane", "", "Name of control plane")
	kubeconfigCmd.Flags().StringVar(&clusterName, "cluster", "", "Name of cluster")
	kubeconfigCmd.MarkFlagRequired("controlplane")
	kubeconfigCmd.MarkFlagRequired("cluster")

	return getCmd

}
