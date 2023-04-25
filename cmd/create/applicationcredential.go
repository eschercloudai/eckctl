package create

import (
	"context"
	"eckctl/pkg/auth"
	"eckctl/pkg/generated"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var createApplicationCredentialCmd = &cobra.Command{
	Use:   "applicationcredential",
	Short: "Create an application credential",
	Run: func(cmd *cobra.Command, args []string) {
		url := cmd.Flag("url").Value.String()
		u := cmd.Flag("username").Value.String()
		p := cmd.Flag("password").Value.String()
		project := cmd.Flag("project").Value.String()
		token := auth.GetToken(url, u, p, project)
		createApplicationCredential(token, url)
	},
}

func createApplicationCredential(bearer string, url string) ApplicationCredential {
	client := auth.InitClient(url)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	aco := new(generated.ApplicationCredentialOptions)

	aco.Name = controlPlaneName + "-" + clusterName

	resp, err := client.PostApiV1ProvidersOpenstackApplicationCredentials(ctx, *aco, auth.SetAuthorizationHeader(bearer))
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	ac := ApplicationCredential{}
	err = json.Unmarshal(body, &ac)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error creating application credential: ", resp.StatusCode)
		log.Fatal(err)
	}

	return ac

}
