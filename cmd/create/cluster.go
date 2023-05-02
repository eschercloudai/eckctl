package create

import (
	"context"
	"eckctl/pkg/auth"
	"eckctl/pkg/generated"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
)

var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create an ECK-managed cluster from a JSON definition",
	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
			cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		token = auth.GetToken(url, u, p, project)
		createCluster()
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

func createCluster() (err error) {
	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cluster, err := readClusterDefs(clusterDefPath)
	if err != nil {
		return
	}

	ac, err := createApplicationCredential(controlPlaneName + "-" + clusterName)
	if err != nil {
		return
	}

	cluster.Name = clusterName
	cluster.Openstack.ApplicationCredentialID = ac.Id
	cluster.Openstack.ApplicationCredentialSecret = ac.Secret

	fmt.Printf("Creating cluster %s with the following definition:\n", cluster.Name)

	clusterJson, err := json.Marshal(cluster)
	if err != nil {
		return
	}

	opts := pretty.DefaultOptions
	opts.SortKeys = true

	fmt.Println(string(pretty.Color(pretty.PrettyOptions(clusterJson, opts), nil)))

	resp, err := client.PostApiV1ControlplanesControlPlaneNameClusters(ctx, controlPlaneName, cluster, auth.SetAuthorizationHeader((token)))
	if resp.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("Error submitting cluster definition, %v", resp.StatusCode)
		return
	}
	return
}
