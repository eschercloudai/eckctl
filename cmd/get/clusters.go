package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/eschercloudai/eckctl/pkg/auth"
	"github.com/eschercloudai/eckctl/pkg/generated"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

var (
	clustersCmd = &cobra.Command{
		Use:     "clusters",
		Short:   "Get clusters",
		Aliases: []string{"cluster"},
		Run: func(cmd *cobra.Command, args []string) {
			url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
				cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
			clusterName = cmd.Flag("name").Value.String()
			insecure, _ = cmd.Flags().GetBool("insecure")
			token, err := auth.GetToken(url, u, p, project, insecure)
			if err != nil {
				log.Fatalf("Error authenticating: %s", err)
			}
			err = printClusters(token, clusterName)
			if err != nil {
				log.Fatalf("Error getting clusters: %s", err)
			}
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

func getClusters(controlplane string, clustername string, token string) ([]generated.KubernetesCluster, error) {

	client, err := auth.NewClient(url, token, insecure)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var resp *http.Response

	if clustername != "" {
		resp, err = client.GetApiV1ControlplanesControlPlaneNameClustersClusterName(ctx, controlplane, clustername)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err = client.GetApiV1ControlplanesControlPlaneNameClusters(ctx, controlplane)
		if err != nil {
			return nil, err
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			err = fmt.Errorf("Cluster not found, %v", resp.StatusCode)
		case http.StatusInternalServerError:
			err = fmt.Errorf("Server error, %v", resp.StatusCode)
		default:
			err = fmt.Errorf("Error retrieving cluster information, %v", resp.StatusCode)
		}
		return nil, err
	}

	var clusters generated.KubernetesClusters
	var cluster generated.KubernetesCluster

	if clustername != "" {
		err = json.Unmarshal(body, &cluster)
		clusters = append(clusters, cluster)
		return clusters, err
	}

	err = json.Unmarshal(body, &clusters)
	if err != nil {
		return clusters, err
	}

	return clusters, err
}

func printClusters(token string, name string) (err error) {
	controlPlanes, err := getControlPlanes(token)
	if err != nil {
		return err
	}
	if allFlag {
		if err != nil {
			return
		}
		for _, c := range controlPlanes {
			clusters, err := getClusters(c.Name, "", token)
			if err != nil {
				return err
			}
			for _, s := range clusters {
				printClusterDetails(c.Name, s)
			}
		}
	} else if clusterControlPlaneName != "" {
		var clusters []generated.KubernetesCluster
		if name != "" {
			clusters, err = getClusters(clusterControlPlaneName, name, token)
		} else {
			clusters, err = getClusters(clusterControlPlaneName, "", token)
		}
		if err != nil {
			return err
		}
		for _, c := range clusters {
			printClusterDetails(clusterControlPlaneName, c)
		}
	} else {
		log.Fatal("Error: Either --controlplane or --all must be specified")
		return err
	}

	return
}
