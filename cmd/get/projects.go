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

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Get images",

	Run: func(cmd *cobra.Command, args []string) {
		url := cmd.Flag("url").Value.String()
		u := cmd.Flag("username").Value.String()
		p := cmd.Flag("password").Value.String()
		project := cmd.Flag("project").Value.String()
		token := auth.GetToken(url, u, p, project)
		getProjects(token, url)
	},
}

func getProjects(bearer string, url string) {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ProvidersOpenstackProjects(ctx, auth.SetAuthorizationHeader(bearer))
	if err != nil {
		fmt.Println("Error getting projects: ", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	projects := generated.OpenstackProjects{}

	err = json.Unmarshal(body, &projects)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range projects {
		fmt.Printf("ID: %s\n", p.Id)
		fmt.Printf("Name: %s\n", p.Name)
		if p.Description != nil {
			fmt.Printf("Description: %v\n", *p.Description)
		}
	}

}
