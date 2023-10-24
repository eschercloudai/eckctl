package delete

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eschercloudai/eckctl/pkg/auth"
	"github.com/spf13/cobra"
)

var (
	clusterName      string
	controlPlaneName string
)

func deleteClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Delete an ECK-managed cluster",
		Run: func(cmd *cobra.Command, args []string) {
			url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
				cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
			insecure, _ = cmd.Flags().GetBool("insecure")
			token, err := auth.GetToken(url, u, p, project, insecure)
			if err != nil {
				return
			}
			err = deleteCluster(token)
			if err != nil {
				log.Fatalf("Error deleting cluster: %s", err)
			}
		},
	}
	cmd.Flags().StringVar(&clusterName, "name", "", "The name of the cluster to be deleted")
	cmd.Flags().StringVar(&controlPlaneName, "controlplane", "", "The name of the associated control plane")
	return cmd
}

func deleteCluster(token string) (err error) {
	client, err := auth.NewClient(url, token, insecure)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.DeleteApiV1ControlplanesControlPlaneNameClustersClusterName(ctx, controlPlaneName, clusterName)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("Unpexected response code: %v", resp.StatusCode)
	}

	fmt.Printf("Deleting cluster %s from controlplane %s\n", clusterName, controlPlaneName)

	return
}
