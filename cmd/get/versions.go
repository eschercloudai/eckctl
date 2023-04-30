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

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "Get versions (application bundles)",

	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
			cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		token = auth.GetToken(url, u, p, project)
		getVersions()
	},
}

func getControlPlaneBundles() []generated.ApplicationBundle {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ApplicationBundlesControlPlane(ctx, auth.SetAuthorizationHeader(token))
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	versions := generated.ApplicationBundles{}
	err = json.Unmarshal(body, &versions)
	if err != nil {
		log.Fatal(err)
	}

	return versions
}

func getClusterBundles() []generated.ApplicationBundle {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ApplicationBundlesCluster(ctx, auth.SetAuthorizationHeader(token))
	if err != nil {
		log.Fatal(err)
	}

	versions := generated.ApplicationBundles{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &versions)
	if err != nil {
		log.Fatal(err)
	}

	return versions
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

func getVersions() {
	fmt.Println("Cluster Bundles:")
	for _, i := range getClusterBundles() {
		printBundle(i)
	}
	fmt.Println("Control Plane Bundles:")
	for _, i := range getControlPlaneBundles() {
		printBundle(i)
	}
}
