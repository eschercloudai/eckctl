package get

import (
	"context"
	"eckctl/pkg/auth"
	"eckctl/pkg/generated"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:     "images",
	Short:   "Get images",
	Aliases: []string{"image"},

	Run: func(cmd *cobra.Command, args []string) {
		url := cmd.Flag("url").Value.String()
		u := cmd.Flag("username").Value.String()
		p := cmd.Flag("password").Value.String()
		project := cmd.Flag("project").Value.String()
		token := auth.GetToken(url, u, p, project)
		getImages(token, url)
	},
}

func printImageDetails(i generated.OpenstackImage) {
	fmt.Printf("Name: %s\tUUID: %s\tCreated: %s\tKubernetes version: %s\tNVIDIA driver version: %s\n", i.Name, i.Id, i.Created, i.Versions.Kubernetes, i.Versions.NvidiaDriver)
}

func getImages(bearer string, url string) {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ProvidersOpenstackImages(ctx, auth.SetAuthorizationHeader(bearer))
	if err != nil {
		fmt.Println("Error getting images: ", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	images := []generated.OpenstackImage{}
	err = json.Unmarshal(body, &images)
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(images, func(j, k int) bool { return images[k].Created.After(images[j].Created) })

	for _, i := range images {
		if (imageId != "" && i.Id == imageId) || (imageName != "" && i.Name == imageName) || (imageId == "" && imageName == "") {
			printImageDetails(i)
		}
	}

}
