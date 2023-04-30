package delete

import (
	"context"
	"eckctl/pkg/auth"
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
			url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
				cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
			token = auth.GetToken(url, u, p, project)
			deleteApplicationCredential()
		},
	}
	cmd.Flags().StringVar(&applicationCredentialName, "name", "", "The name of the application credential to be deleted")
	err := cmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatalln(err)
	}
	return cmd
}

func deleteApplicationCredential() {
	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.DeleteApiV1ProvidersOpenstackApplicationCredentialsApplicationCredential(ctx, applicationCredentialName, auth.SetAuthorizationHeader(token))

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error deleting application credential %s, %v", applicationCredentialName, resp.StatusCode)
	}
}
