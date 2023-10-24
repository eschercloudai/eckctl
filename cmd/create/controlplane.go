package create

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eschercloudai/eckctl/pkg/auth"
	"github.com/eschercloudai/eckctl/pkg/generated"
	"github.com/spf13/cobra"
)

var createControlPlaneCmd = &cobra.Command{
	Use:   "controlplane",
	Short: "Create a control plane",
	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
			cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		insecure, _ = cmd.Flags().GetBool("insecure")
		token, err := auth.GetToken(url, u, p, project, insecure)
		if err != nil {
			log.Fatalf("Error authenticating: %s", err)
		}
		err = createControlPlane(token)
		if err != nil {
			log.Fatalf("Error creating control plane: %s", err)
		}
	},
}

func createControlPlane(token string) (err error) {
	client, err := auth.NewClient(url, token, insecure)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cp := generated.ControlPlane{}

	cp.Name = controlPlaneName
	cp.ApplicationBundle.Name = "control-plane-" + controlPlaneVersion
	cp.ApplicationBundle.Version = controlPlaneVersion

	// Create the Unikorn Project if it doesn't already exist, 409s are OK
	resp, err := client.PostApiV1Project(ctx)
	if (resp.StatusCode != http.StatusConflict) && (resp.StatusCode != http.StatusAccepted) {
		err = fmt.Errorf("Error creating project, response code: %v", resp.StatusCode)
		return
	}

	resp, err = client.PostApiV1Controlplanes(ctx, cp)
	if resp.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("Unexpected response code: %v", resp.StatusCode)
		return
	}

	fmt.Printf("Control Plane %s is being created\n", cp.Name)

	return
}
