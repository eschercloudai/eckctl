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
)

var networksCmd = &cobra.Command{
	Use:   "networks",
	Short: "Get networks",

	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
			cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		token = auth.GetToken(url, u, p, project)
		getNetworks()
	},
}

func getNetworks() {

	client := auth.InitClient(url)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ProvidersOpenstackExternalNetworks(ctx, auth.SetAuthorizationHeader(token))
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	networks := generated.OpenstackExternalNetworks{}
	err = json.Unmarshal(body, &networks)
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range networks {
		fmt.Printf("Name: %s\tID: %s\n", i.Name, i.Id)
	}
}
