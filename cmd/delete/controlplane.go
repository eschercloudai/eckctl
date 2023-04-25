package delete

import (
	"eckctl/pkg/auth"
	"eckctl/pkg/generated"
	"fmt"
	"log"
	"net/http"

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
	return cmd
}

func deleteControlPlane(b string, u string) {
	client := &http.Client{}

	req, err := generated.NewDeleteApiV1ControlplanesControlPlaneNameRequest(u, controlPlaneName)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+b)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(resp.StatusCode)
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		fmt.Println(resp.StatusCode)
		log.Fatal(err)
	}
}
