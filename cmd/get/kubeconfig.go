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
		getKubeConfig()
	},
}

func getKubeConfig() {
	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ControlplanesControlPlaneNameClustersClusterNameKubeconfig(ctx, controlPlaneName, clusterName, auth.SetAuthorizationHeader((token)))
	if err != nil {
		fmt.Println("Error retrieving kubeconfig: ", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
