package get

import (
	"context"
	"eckctl/pkg/auth"
	"eckctl/pkg/generated"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
)

func controlPlaneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "controlplanes",
		Aliases: []string{"control-planes", "controlplane", "control-plane", "cp"},
		Short:   "Get control planes",
		Run: func(cmd *cobra.Command, args []string) {
			url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(), cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
			project = cmd.Flag("project").Value.String()
			token = auth.GetToken(url, u, p, project)
			printControlPlanes()
		},
	}
	cmd.Flags().StringVar(&controlPlaneName, "name", "", "The name of the control plane to list")
	return cmd
}

func printControlPlaneDetails(i generated.ControlPlane) {
	fmt.Printf("Name: %s\tStatus: %s\tVersion: %s\n", i.Name, i.Status.Status, i.ApplicationBundle.Version)
}

func getControlPlanes() []generated.ControlPlane {
	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1Controlplanes(ctx, auth.SetAuthorizationHeader(token))
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	controlPlanes := generated.ControlPlanes{}
	err = json.Unmarshal(body, &controlPlanes)
	if err != nil {
		log.Fatal(err)
	}

	return controlPlanes
}

func printControlPlanes() {
	for _, i := range getControlPlanes() {
		if (controlPlaneName != "" && i.Name == controlPlaneName) || (controlPlaneName == "") {
			printControlPlaneDetails(i)
		}
	}
}
