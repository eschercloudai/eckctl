package update

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eschercloudai/eckctl/pkg/auth"
	"github.com/eschercloudai/eckctl/pkg/generated"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
)

var updateClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Update an existing ECK-managed cluster from a JSON definition",
	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
			cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		insecure, _ = cmd.Flags().GetBool("insecure")
		token, err := auth.GetToken(url, u, p, project, insecure)
		if err != nil {
			log.Fatalf("Error authenticating: %s", err)
		}
		err = updateCluster(token)
		if err != nil {
			log.Fatalf("Error creating cluster: %s", err)
		}
	},
}

func readClusterDefs(filePath string) (cluster generated.KubernetesCluster, err error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &cluster)
	if err != nil {
		return
	}
	return
}

func updateCluster(token string) (err error) {
	client, err := auth.NewClient(url, token, insecure)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cluster, err := readClusterDefs(clusterDefPath)
	if err != nil {
		return
	}

	cluster.Name = clusterName

	fmt.Printf("Updating cluster %s from the following definition:\n", cluster.Name)

	clusterJson, err := json.Marshal(cluster)
	if err != nil {
		return
	}

	opts := pretty.DefaultOptions
	opts.SortKeys = true

	fmt.Println(string(pretty.Color(pretty.PrettyOptions(clusterJson, opts), nil)))

	resp, err := client.PutApiV1ControlplanesControlPlaneNameClustersClusterName(ctx, controlPlaneName, clusterName, cluster)
	if resp.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("Error submitting cluster definition, %v", resp.StatusCode)
		return
	}
	return
}
