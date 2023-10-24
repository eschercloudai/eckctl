package get

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	controlPlaneName   string
	clusterName        string
	imageName          string
	imageId            string
	allFlag            bool
	insecure           bool
	url, u, p, project string
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
	imagesCmd.Flags().StringVar(&imageId, "id", "", "ID of image")
	clustersCmd.Flags().StringVar(&controlPlaneName, "controlplane", "", "Name of control plane")
	clustersCmd.Flags().BoolVar(&allFlag, "all", false, "Return all clusters across all control planes")
	clustersCmd.Flags().StringVar(&clusterName, "name", "", "Name of cluster")
	clustersCmd.MarkFlagsMutuallyExclusive("controlplane", "all")
	kubeconfigCmd.Flags().StringVar(&controlPlaneName, "controlplane", "", "Name of control plane")
	kubeconfigCmd.Flags().StringVar(&clusterName, "cluster", "", "Name of cluster")
	err := kubeconfigCmd.MarkFlagRequired("controlplane")
	if err != nil {
		log.Fatalln(err)
	}
	err = kubeconfigCmd.MarkFlagRequired("cluster")
	if err != nil {
		log.Fatalln(err)
	}

	return getCmd

}
