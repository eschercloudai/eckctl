package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/eschercloudai/eckctl/pkg/auth"
	"github.com/eschercloudai/eckctl/pkg/generated"
	"github.com/spf13/cobra"
)

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "Get versions (application bundles)",

	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
			cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		insecure, _ = cmd.Flags().GetBool("insecure")
		token, err := auth.GetToken(url, u, p, project, insecure)
		if err != nil {
			log.Fatalf("Error authenticating: %s", err)
		}
		err = getVersions(token)
		if err != nil {
			log.Fatalf("Error retrieving versions: %s", err)
		}
	},
}

func getControlPlaneBundles(token string) (versions []generated.ApplicationBundle, err error) {

	client, err := auth.NewClient(url, token, insecure)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ApplicationbundlesControlPlane(ctx)
	if err != nil {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	versions = generated.ApplicationBundles{}
	err = json.Unmarshal(body, &versions)
	if err != nil {
		return
	}

	return
}

func getClusterBundles(token string) (versions []generated.ApplicationBundle, err error) {

	client, err := auth.NewClient(url, token, insecure)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ApplicationbundlesCluster(ctx)
	if err != nil {
		log.Fatal(err)
	}

	versions = generated.ApplicationBundles{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &versions)
	if err != nil {
		return
	}

	return
}

func printBundle(bundle generated.ApplicationBundle) {
	fmt.Printf("Name: %s\tVersion: %s", bundle.Name, bundle.Version)
	if bundle.EndOfLife != nil {
		fmt.Printf("\tEOL: %v", bundle.EndOfLife.Format(time.RFC822))
	}
	if bundle.Preview != nil && *bundle.Preview {
		fmt.Print("\tPreview: True\n")
	} else {
		fmt.Println()
	}
}

func getVersions(token string) (err error) {
	fmt.Println("Cluster Bundles:")
	clusterBundles, err := getClusterBundles(token)
	if err != nil {
		return
	}
	for _, i := range clusterBundles {
		printBundle(i)
	}
	controlPlaneBundles, err := getControlPlaneBundles(token)
	if err != nil {
		return
	}
	fmt.Println("Control Plane Bundles:")
	for _, i := range controlPlaneBundles {
		printBundle(i)
	}
	return err
}
