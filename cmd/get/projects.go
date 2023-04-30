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

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Get images",

	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
			cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		token = auth.GetToken(url, u, p, project)
		getProjects()
	},
}

func getProjects() {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ProvidersOpenstackProjects(ctx, auth.SetAuthorizationHeader(token))
	if err != nil {
		fmt.Println("Error getting projects: ", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	projects := generated.OpenstackProjects{}

	err = json.Unmarshal(body, &projects)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range projects {
		fmt.Printf("Name: %s\t", p.Name)
		if p.Description != nil {
			fmt.Printf("Description: %v\t", *p.Description)
		}
		fmt.Printf("ID: %s\n", p.Id)
	}

}
