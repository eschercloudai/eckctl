package get

import (
	"context"
	"eckctl/pkg/auth"
	"eckctl/pkg/generated"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

var clustersCmd = &cobra.Command{
	Use:     "clusters",
	Short:   "Get clusters",
	Aliases: []string{"cluster"},
	Run: func(cmd *cobra.Command, args []string) {
		url := cmd.Flag("url").Value.String()
		u := cmd.Flag("username").Value.String()
		p := cmd.Flag("password").Value.String()
		project := cmd.Flag("project").Value.String()
		token := auth.GetToken(url, u, p, project)
		getClusters(token, url)
	},
}

func printClusterDetails(i generated.KubernetesCluster) {
	tree := treeprint.New()
	fmt.Printf("Cluster: %s, version: %s, status: %s", i.Name, i.ControlPlane.Version, i.Status.Status)
	if i.WorkloadPools != nil {
		pool := tree.AddBranch("Pools:")
		for _, p := range i.WorkloadPools {
			pool.AddNode(fmt.Sprintf("Name: %s\tFlavor: %s\tImage: %s", p.Name, p.Machine.FlavorName, p.Machine.ImageName))
		}
		fmt.Println(tree.String())
	}
}

func getClusters(bearer string, url string) {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ControlplanesControlPlaneNameClusters(ctx, controlPlaneName, auth.SetAuthorizationHeader((bearer)))
	if err != nil {
		fmt.Println("Error retrieving clusters: ", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	clusters := generated.KubernetesClusters{}
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range clusters {
		if clusterName == "" {
			fmt.Printf("Name: %s\tVersion: %s\tStatus: %s\n", c.Name, c.ControlPlane.Version, c.Status.Status)
		} else if c.Name == clusterName {
			printClusterDetails(c)
		}
	}
}
