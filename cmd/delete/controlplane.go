package delete

import (
	"context"
	"eckctl/pkg/auth"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

func deleteControlPlaneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "controlplane",
		Short: "Delete an ECK control plane",
		Run: func(cmd *cobra.Command, args []string) {
			url := cmd.Flag("url").Value.String()
			u := cmd.Flag("username").Value.String()
			p := cmd.Flag("password").Value.String()
			project := cmd.Flag("project").Value.String()
			token := auth.GetToken(url, u, p, project)
			deleteControlPlane(token, url)
		},
	}
	cmd.Flags().StringVar(&controlPlaneName, "name", "", "The name of the control plane to be deleted")
	cmd.MarkFlagRequired("name")
	return cmd
}

func deleteControlPlane(bearer string, url string) {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.DeleteApiV1ControlplanesControlPlaneName(ctx, controlPlaneName, auth.SetAuthorizationHeader(bearer))
	if err != nil {
		log.Fatal(resp, err)
	}

	if resp.StatusCode != http.StatusAccepted {
		log.Fatalf("Error deleting control plane %s, %v", controlPlaneName, resp.StatusCode)
	}
}
