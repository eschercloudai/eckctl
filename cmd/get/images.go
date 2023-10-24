package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sort"
	"time"

	"github.com/eschercloudai/eckctl/pkg/auth"
	"github.com/eschercloudai/eckctl/pkg/generated"
	"github.com/spf13/cobra"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:     "images",
	Short:   "Get images",
	Aliases: []string{"image"},

	Run: func(cmd *cobra.Command, args []string) {
		url, u, p, project = cmd.Flag("url").Value.String(), cmd.Flag("username").Value.String(),
			cmd.Flag("password").Value.String(), cmd.Flag("project").Value.String()
		insecure, _ = cmd.Flags().GetBool("insecure")
		token, err := auth.GetToken(url, u, p, project, insecure)
		if err != nil {
			log.Fatalf("Error authenticating: %s", err)
		}
		err = getImages(token)
		if err != nil {
			log.Fatalf("Error retrieving images: %s", err)
		}
	},
}

func printImageDetails(i generated.OpenstackImage) {
	fmt.Printf("Name: %s\tUUID: %s\tCreated: %s\tKubernetes version: %s\tNVIDIA driver version: %s\n", i.Name, i.Id, i.Created, i.Versions.Kubernetes, i.Versions.NvidiaDriver)
}

func getImages(token string) (err error) {
	client, err := auth.NewClient(url, token, insecure)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ProvidersOpenstackImages(ctx)
	if err != nil {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	images := []generated.OpenstackImage{}
	err = json.Unmarshal(body, &images)
	if err != nil {
		return
	}
	sort.Slice(images, func(j, k int) bool { return images[k].Created.After(images[j].Created) })

	for _, i := range images {
		if (imageId != "" && i.Id == imageId) || (imageName != "" && i.Name == imageName) || (imageId == "" && imageName == "") {
			printImageDetails(i)
		}
	}

	return

}
