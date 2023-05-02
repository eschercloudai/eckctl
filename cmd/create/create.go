package create

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	applicationCredentialName string
	controlPlaneName          string
	controlPlaneVersion       string
	clusterName               string
	clusterDefPath            string
	url, token, u, p, project string
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
		createApplicationCredentialCmd,
		createClusterCmd,
	}

	createCmd.AddCommand(commands...)

	createApplicationCredentialCmd.Flags().StringVar(&applicationCredentialName, "name", "", "Name of application credential")
	err := createApplicationCredentialCmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatalln(err)
	}
	createControlPlaneCmd.Flags().StringVar(&controlPlaneName, "name", "", "Name of control plane")
	err = createControlPlaneCmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatalln(err)
	}
	createControlPlaneCmd.Flags().StringVar(&controlPlaneVersion, "version", "1.0.1", "Version of control plane")
	createClusterCmd.Flags().StringVar(&clusterName, "name", "", "Name of cluster")
	createClusterCmd.Flags().StringVar(&controlPlaneName, "controlplane", "", "Name of associated control plane")
	createClusterCmd.Flags().StringVar(&clusterDefPath, "json", "", "Path to JSON cluster definition")
	err = createClusterCmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatalln(err)
	}
	err = createClusterCmd.MarkFlagRequired("controlplane")
	if err != nil {
		log.Fatalln(err)
	}
	err = createClusterCmd.MarkFlagRequired("json")
	if err != nil {
		log.Fatalln(err)
	}

	return createCmd
}
