package get

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/eschercloudai/eckctl/pkg/auth"
	"github.com/spf13/cobra"
)

var kubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Get kubeconfig",

	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(), cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		token, err := auth.GetToken(url, u, p, project)
		if err != nil {
			log.Fatalf("Error authenticating: %s", err)
		}
		kubeConfig, err := getKubeConfig(token)
		if err != nil {
			log.Fatalf("Error retrieving kubeconfig: %s", err)
		}
		fmt.Println(kubeConfig)
	},
}

func getKubeConfig(token string) (kubeconfig string, err error) {
	client, err := auth.NewClient(url, token)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ControlplanesControlPlaneNameClustersClusterNameKubeconfig(ctx, controlPlaneName, clusterName)
	if err != nil {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	kubeconfig = string(body)

	return
}
