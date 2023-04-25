package delete

import (
	"context"
	"eckctl/pkg/auth"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var applicationCredentialName string

func deleteApplicationCredentialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "applicationcredential",
		Short: "Delete an application credential",
		Run: func(cmd *cobra.Command, args []string) {
			url := cmd.Flag("url").Value.String()
			u := cmd.Flag("username").Value.String()
			p := cmd.Flag("password").Value.String()
			project := cmd.Flag("project").Value.String()
			token := auth.GetToken(url, u, p, project)
			deleteApplicationCredential(token, url)
		},
	}
	cmd.Flags().StringVar(&applicationCredentialName, "name", "", "The name of the application credential to be deleted")
	cmd.MarkFlagRequired("name")
	return cmd
}

func deleteApplicationCredential(bearer string, url string) {
	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.DeleteApiV1ProvidersOpenstackApplicationCredentialsApplicationCredential(ctx, applicationCredentialName, auth.SetAuthorizationHeader(bearer))

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		log.Fatal(err)
	}
}
