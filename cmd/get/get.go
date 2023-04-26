package get

import (
	"github.com/spf13/cobra"
	"log"
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
	err := clustersCmd.MarkFlagRequired("controlplane")
	if err != nil {
		log.Fatalln(err)
	}
	kubeconfigCmd.Flags().StringVar(&controlPlaneName, "controlplane", "", "Name of control plane")
	kubeconfigCmd.Flags().StringVar(&clusterName, "cluster", "", "Name of cluster")
	err = kubeconfigCmd.MarkFlagRequired("controlplane")
	if err != nil {
		log.Fatalln(err)
	}
	err = kubeconfigCmd.MarkFlagRequired("cluster")
	if err != nil {
		log.Fatalln(err)
	}

	return getCmd

}
