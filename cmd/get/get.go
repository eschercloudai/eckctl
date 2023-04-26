package get

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	controlPlaneName string
	clusterName      string
	imageName        string
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

	imagesCmd.Flags().StringVar(&imageName, "name", "", "Name of image")
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
