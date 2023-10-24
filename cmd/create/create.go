package create

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	controlPlaneName    string
	controlPlaneVersion string
	clusterName         string
	clusterDefPath      string
	insecure            bool
	url, u, p, project  string
)

type ApplicationCredential struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	Secret string `json:"secret"`
}

func NewCreateCommand() *cobra.Command {
	createCmd := &cobra.Command{
		Use: "create",
	}

	commands := []*cobra.Command{
		createControlPlaneCmd,
		createClusterCmd,
	}

	createCmd.AddCommand(commands...)

	createControlPlaneCmd.Flags().StringVar(&controlPlaneName, "name", "", "Name of control plane")
	if err := createControlPlaneCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalln(err)
	}
	createControlPlaneCmd.Flags().StringVar(&controlPlaneVersion, "version", "1.0.1", "Version of control plane")
	createClusterCmd.Flags().StringVar(&clusterName, "name", "", "Name of cluster")
	createClusterCmd.Flags().StringVar(&controlPlaneName, "controlplane", "", "Name of associated control plane")
	createClusterCmd.Flags().StringVar(&clusterDefPath, "json", "", "Path to JSON cluster definition")
	if err := createClusterCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalln(err)
	}
	if err := createClusterCmd.MarkFlagRequired("controlplane"); err != nil {
		log.Fatalln(err)
	}
	if err := createClusterCmd.MarkFlagRequired("json"); err != nil {
		log.Fatalln(err)
	}

	return createCmd
}
