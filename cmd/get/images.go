package get

import (
	"context"
	"eckctl/pkg/auth"
	"eckctl/pkg/generated"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Get images",

	Run: func(cmd *cobra.Command, args []string) {
		url := cmd.Flag("url").Value.String()
		u := cmd.Flag("username").Value.String()
		p := cmd.Flag("password").Value.String()
		project := cmd.Flag("project").Value.String()
		token := auth.GetToken(url, u, p, project)
		getImages(token, url)
	},
}

func getImages(bearer string, url string) {

	client := auth.InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetApiV1ProvidersOpenstackImages(ctx, auth.SetAuthorizationHeader(bearer))
	if err != nil {
		fmt.Println("Error getting images: ", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
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
		fmt.Printf("Name: %s\t", i.Name)
		fmt.Printf("Created: %s\t", i.Created)
		fmt.Printf("Kubernetes version: %s\t", i.Versions.Kubernetes)
		fmt.Printf("NVIDIA driver version: %s \n", i.Versions.NvidiaDriver)
	}

}
