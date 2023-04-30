package delete

import (
	"context"
	"eckctl/pkg/auth"
	"fmt"
	"log"
	"net/http"
	"time"

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
			token = auth.GetToken(url, u, p, project)
			deleteCluster(clusterName, controlPlaneName)
		},
	}
	cmd.Flags().StringVar(&clusterName, "name", "", "The name of the cluster to be deleted")
	cmd.Flags().StringVar(&controlPlaneName, "controlplane", "", "The name of the associated control plane")
	return cmd
}

func deleteCluster(clusterName string, controlPlaneName string) {
	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.DeleteApiV1ControlplanesControlPlaneNameClustersClusterName(ctx, controlPlaneName, clusterName, auth.SetAuthorizationHeader(token))
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		log.Fatalf("Error deleting cluster %s from control plane %s, %v", clusterName, controlPlaneName, resp.StatusCode)
	}

	fmt.Printf("Deleting cluster %s from controlplane %s\n", clusterName, controlPlaneName)
}
