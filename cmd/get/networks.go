package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/eschercloudai/eckctl/pkg/auth"
	"github.com/eschercloudai/eckctl/pkg/generated"
	"github.com/spf13/cobra"
)

var networksCmd = &cobra.Command{
	Use:   "networks",
	Short: "Get networks",

	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
			cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		token, err := auth.GetToken(url, u, p, project)
		if err != nil {
			log.Fatalf("Error authenticating: %s", err)
		}
		err = getNetworks(token)
		if err != nil {
			log.Fatalf("Error getting networks, %s", err)
		}
	},
}

func getNetworks(token string) (err error) {

	client, err := auth.NewClient(url, token)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ProvidersOpenstackExternalNetworks(ctx)
	if err != nil {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	networks := generated.OpenstackExternalNetworks{}
	err = json.Unmarshal(body, &networks)
	if err != nil {
		return
	}

	for _, i := range networks {
		fmt.Printf("Name: %s\tID: %s\n", i.Name, i.Id)
	}

	return
}
