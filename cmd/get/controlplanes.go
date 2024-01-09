package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/eschercloudai/eckctl/pkg/auth"
	"github.com/eschercloudai/eckctl/pkg/generated"
	"github.com/spf13/cobra"
)

func controlPlaneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "controlplanes",
		Aliases: []string{"control-planes", "controlplane", "control-plane", "cp"},
		Short:   "Get control planes",
		Run: func(cmd *cobra.Command, args []string) {
			url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
				cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
			insecure, _ = cmd.Flags().GetBool("insecure")
			token, err := auth.GetToken(url, u, p, project, insecure)
			if err != nil {
				log.Fatalf("Error authenticating: %s", err)
			}
			err = printControlPlanes(token)
			if err != nil {
				log.Fatalf("Error retrieving control planes: %s", err)
			}
		},
	}
	cmd.Flags().StringVar(&controlPlaneName, "name", "", "The name of the control plane to list")
	return cmd
}

func printControlPlaneDetails(i generated.ControlPlane) {
	fmt.Printf("Name: %s\tStatus: %s\tVersion: %s\n", i.Name, i.Status.Status, i.ApplicationBundle.Version)
}

func getControlPlanes(token string) (controlPlanes []generated.ControlPlane, err error) {
	client, err := auth.NewClient(url, token, insecure)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1Controlplanes(ctx)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusInternalServerError:
			err = fmt.Errorf("Server error, %v", resp.StatusCode)
		default:
			err = fmt.Errorf("Error retrieving control plane information, %v", resp.StatusCode)
		}
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	controlPlanes = generated.ControlPlanes{}
	err = json.Unmarshal(body, &controlPlanes)
	if err != nil {
		return
	}

	return
}

func printControlPlanes(token string) (err error) {
	cps, err := getControlPlanes(token)
	if err != nil {
		return
	}
	for _, i := range cps {
		if i.Name == controlPlaneName {
			printControlPlaneDetails(i)
			break
		} else {
			printControlPlaneDetails(i)
		}
	}
	return
}
