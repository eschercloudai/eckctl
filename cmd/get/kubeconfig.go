package get

import (
	"context"
	"eckctl/pkg/auth"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var kubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Get kubeconfig",

	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(), cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		token = auth.GetToken(url, u, p, project)
		kubeConfig, err := getKubeConfig()
		if err != nil {
			log.Fatalf("Error retrieving kubeconfig: %s", err)
		}
		fmt.Println(kubeConfig)
	},
}

func getKubeConfig() (kubeconfig string, err error) {
	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ControlplanesControlPlaneNameClustersClusterNameKubeconfig(ctx, controlPlaneName, clusterName, auth.SetAuthorizationHeader((token)))
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
