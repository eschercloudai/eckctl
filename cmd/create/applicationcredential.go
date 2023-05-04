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
		token, err := auth.GetToken(url, u, p, project)
		if err != nil {
			log.Fatalf("Error authenticating: %s", err)
		}
		_, err = createApplicationCredential(applicationCredentialName, token)
		if err != nil {
			log.Fatalf("Error creating application credential: %s", err)
		}
	},
}

func createApplicationCredential(name string, token string) (ac ApplicationCredential, err error) {
	client, err := auth.InitClient(url)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	aco := new(generated.ApplicationCredentialOptions)

	aco.Name = name

	resp, err := client.PostApiV1ProvidersOpenstackApplicationCredentials(ctx, *aco, auth.SetAuthorizationHeader(token))
	if err != nil {
		return
	}

	body, _ := io.ReadAll(resp.Body)
	ac = ApplicationCredential{}
	err = json.Unmarshal(body, &ac)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error creating application credential: %v", resp.StatusCode)
		return
	}

	return

}
