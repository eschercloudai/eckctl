package create

import (
	"context"
	"eckctl/pkg/auth"
	"eckctl/pkg/generated"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var createApplicationCredentialCmd = &cobra.Command{
	Use:   "applicationcredential",
	Short: "Create an application credential",
	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
			cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		token = auth.GetToken(url, u, p, project)
		createApplicationCredential()
	},
}

func createApplicationCredential() ApplicationCredential {
	client := auth.InitClient(url)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	aco := new(generated.ApplicationCredentialOptions)

	aco.Name = controlPlaneName + "-" + clusterName

	resp, err := client.PostApiV1ProvidersOpenstackApplicationCredentials(ctx, *aco, auth.SetAuthorizationHeader(token))
	if err != nil {
		log.Fatal(err)
	}

	body, _ := io.ReadAll(resp.Body)
	ac := ApplicationCredential{}
	err = json.Unmarshal(body, &ac)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error creating application credential: ", resp.StatusCode)
		log.Fatal(err)
	}

	return ac

}
