package get

import (
	"context"
	"eckctl/pkg/auth"
	"eckctl/pkg/generated"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var clustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "Get clusters",

	Run: func(cmd *cobra.Command, args []string) {
		url := cmd.Flag("url").Value.String()
		u := cmd.Flag("username").Value.String()
		p := cmd.Flag("password").Value.String()
		project := cmd.Flag("project").Value.String()
		token := auth.GetToken(url, u, p, project)
		getClusters(token, url)
	},
}

func getClusters(bearer string, url string) {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ControlplanesControlPlaneNameClusters(ctx, controlPlaneName, auth.SetAuthorizationHeader((bearer)))
	if err != nil {
		fmt.Println("Error retrieving clusters: ", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	clusters := generated.KubernetesClusters{}
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range clusters {
		fmt.Printf("Name: %s\t", c.Name)
		fmt.Printf("Status: %s\n", c.Status.Status)
	}

}
