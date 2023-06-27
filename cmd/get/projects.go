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

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Get images",

	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
			cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		token, err := auth.GetToken(url, u, p, project)
		if err != nil {
			log.Fatalf("Error authenticating: %s", err)
		}
		err = getProjects(token)
		if err != nil {
			log.Fatalf("Error retrieving projects: %s", err)
		}
	},
}

func getProjects(token string) (err error) {

	client, err := auth.NewClient(url, token)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ProvidersOpenstackProjects(ctx)
	if err != nil {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	projects := generated.OpenstackProjects{}

	err = json.Unmarshal(body, &projects)
	if err != nil {
		return
	}

	for _, p := range projects {
		fmt.Printf("Name: %s\t", p.Name)
		if p.Description != nil {
			fmt.Printf("Description: %v\t", *p.Description)
		}
		fmt.Printf("ID: %s\n", p.Id)
	}

	return

}
