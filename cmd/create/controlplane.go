package create

import (
	"context"
	"eckctl/pkg/auth"
	"eckctl/pkg/generated"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var createControlPlaneCmd = &cobra.Command{
	Use:   "controlplane",
	Short: "Create a control plane",
	Run: func(cmd *cobra.Command, args []string) {
		url := cmd.Flag("url").Value.String()
		u := cmd.Flag("username").Value.String()
		p := cmd.Flag("password").Value.String()
		project := cmd.Flag("project").Value.String()
		token := auth.GetToken(url, u, p, project)
		createControlPlane(token, url)
	},
}

func createControlPlane(bearer string, url string) {
	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cp := generated.ControlPlane{}

	cp.Name = controlPlaneName
	cp.ApplicationBundle.Name = "control-plane-" + controlPlaneVersion
	cp.ApplicationBundle.Version = controlPlaneVersion

	resp, err := client.PostApiV1Controlplanes(ctx, cp, auth.SetAuthorizationHeader(bearer))
	if resp.StatusCode != http.StatusAccepted {
		fmt.Println(resp.StatusCode)
		log.Fatal(err)
	}

	fmt.Printf("Control Plane %s created\n", cp.Name)

}
