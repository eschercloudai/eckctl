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
		url := cmd.Flag("url").Value.String()
		u := cmd.Flag("username").Value.String()
		p := cmd.Flag("password").Value.String()
		project := cmd.Flag("project").Value.String()
		token := auth.GetToken(url, u, p, project)
		getVersions(token, url)
	},
}

func getControlPlaneBundles(bearer string, url string) []generated.ApplicationBundle {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ApplicationBundlesControlPlane(ctx, auth.SetAuthorizationHeader(bearer))
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

func getClusterBundles(bearer string, url string) []generated.ApplicationBundle {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ApplicationBundlesCluster(ctx, auth.SetAuthorizationHeader(bearer))
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
	if bundle.Preview != nil {
		if *bundle.Preview {
			fmt.Print("\tPreview: True\n")
		}
	} else {
		println()
	}
}

func getVersions(b string, u string) {
	fmt.Println("Cluster Bundles:")
	for _, i := range getClusterBundles(b, u) {
		printBundle(i)
	}
	fmt.Println("Control Plane Bundles:")
	for _, i := range getControlPlaneBundles(b, u) {
		printBundle(i)
	}
}
