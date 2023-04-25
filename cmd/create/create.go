package create

import (
	"github.com/spf13/cobra"
)

var (
	applicationCredentialName string
	controlPlaneName          string
	controlPlaneVersion       string
	clusterName               string
	clusterDefPath            string
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

	createApplicationCredentialCmd.Flags().StringVar(&controlPlaneName, "controlplane", "", "Name of associated control plane")
	createApplicationCredentialCmd.Flags().StringVar(&applicationCredentialName, "name", "", "Name of application credential")
	createApplicationCredentialCmd.MarkFlagsRequiredTogether("controlplane", "name")
	createControlPlaneCmd.Flags().StringVar(&controlPlaneName, "name", "", "Name of control plane")
	createControlPlaneCmd.MarkFlagRequired("name")
	createControlPlaneCmd.Flags().StringVar(&controlPlaneVersion, "version", "1.0.1", "Version of control plane")
	createClusterCmd.Flags().StringVar(&clusterName, "name", "", "Name of cluster")
	createClusterCmd.Flags().StringVar(&controlPlaneName, "controlplane", "", "Name of associated control plane")
	createClusterCmd.Flags().StringVar(&clusterDefPath, "json", "", "Path to JSON cluster definition")
	createClusterCmd.MarkFlagRequired("name")
	createClusterCmd.MarkFlagRequired("controlplane")
	createClusterCmd.MarkFlagRequired("json")

	return createCmd
}
