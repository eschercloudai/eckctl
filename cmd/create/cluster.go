package create

import (
	"context"
	"eckctl/pkg/auth"
	"eckctl/pkg/generated"
	"encoding/json"
	"fmt"
	"log"
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
		url := cmd.Flag("url").Value.String()
		u := cmd.Flag("username").Value.String()
		p := cmd.Flag("password").Value.String()
		project := cmd.Flag("project").Value.String()
		token := auth.GetToken(url, u, p, project)
		createCluster(token, url)
	},
}

func readClusterDefs(filePath string) (generated.KubernetesCluster, error) {
	var clusters generated.KubernetesCluster
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return clusters, fmt.Errorf("error opening file: %w", err)
	}

	err = json.Unmarshal(bytes, &clusters)
	if err != nil {
		return clusters, fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	return clusters, err
}

func createCluster(bearer string, url string) {
	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cluster, err := readClusterDefs(clusterDefPath)
	if err != nil {
		log.Fatal(err)
	}

	ac := createApplicationCredential(bearer, url)

	cluster.Name = clusterName
	cluster.Openstack.ApplicationCredentialID = ac.Id
	cluster.Openstack.ApplicationCredentialSecret = ac.Secret

	fmt.Printf("Creating cluster %s with the following definition:\n", cluster.Name)

	clusterJson, err := json.Marshal(cluster)
	if err != nil {
		fmt.Println("Error marshalling JSON object")
	}

	opts := pretty.DefaultOptions
	opts.SortKeys = true

	fmt.Println(string(pretty.Color(pretty.PrettyOptions(clusterJson, opts), nil)))

	resp, err := client.PostApiV1ControlplanesControlPlaneNameClusters(ctx, controlPlaneName, cluster, auth.SetAuthorizationHeader((bearer)))
	if resp.StatusCode != http.StatusAccepted {
		log.Fatal(err)
	}
}
