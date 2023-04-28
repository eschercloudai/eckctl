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

var (
	clustersCmd = &cobra.Command{
		Use:     "clusters",
		Short:   "Get clusters",
		Aliases: []string{"cluster"},
		Run: func(cmd *cobra.Command, args []string) {
			url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(), cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
			token = auth.GetToken(url, u, p, project)
			printClusters()
		},
	}
)

func printClusterDetails(controlPlane string, i generated.KubernetesCluster) {
	tree := treeprint.New()
	fmt.Printf("Cluster: %s, version: %s, control plane: %s, status: %s", i.Name, i.ControlPlane.Version, controlPlane, i.Status.Status)
	if i.WorkloadPools != nil {
		pools := tree.AddBranch("Pools:")
		for _, p := range i.WorkloadPools {
			pool := pools.AddBranch(fmt.Sprintf("Name: %s\tFlavor: %s\tImage: %s\tReplicas: %v", p.Name, p.Machine.FlavorName, p.Machine.ImageName, p.Machine.Replicas))
			if p.Autoscaling != nil {
				pool.AddNode(fmt.Sprintf("Autoscaling: Minimum: %v, maximum: %v", p.Autoscaling.MinimumReplicas, p.Autoscaling.MaximumReplicas))
			}
		}
		fmt.Println(tree.String())
	}
}

func getClusters(controlplane string) []generated.KubernetesCluster {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ControlplanesControlPlaneNameClusters(ctx, controlplane, auth.SetAuthorizationHeader((token)))
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

	return clusters
}

func printClusters() {
	var clusters = []generated.KubernetesCluster{}
	if allFlag {
		controlPlanes := getControlPlanes()
		for _, c := range controlPlanes {
			for _, s := range getClusters(c.Name) {
				printClusterDetails(c.Name, s)
			}
		}
	} else {
		clusters = getClusters(controlPlaneName)
		for _, c := range clusters {
			printClusterDetails(controlPlaneName, c)
		}
	}
}
