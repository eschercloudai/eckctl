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
			url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
				cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
			token = auth.GetToken(url, u, p, project)
			deleteControlPlane()
		},
	}
	cmd.Flags().StringVar(&controlPlaneName, "name", "", "The name of the control plane to be deleted")
	err := cmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatalln(err)
	}
	return cmd
}

func deleteControlPlane() {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.DeleteApiV1ControlplanesControlPlaneName(ctx, controlPlaneName, auth.SetAuthorizationHeader(token))
	if err != nil {
		log.Fatal(resp, err)
	}

	if resp.StatusCode != http.StatusAccepted {
		log.Fatalf("Error deleting control plane %s, %v", controlPlaneName, resp.StatusCode)
	}
}
