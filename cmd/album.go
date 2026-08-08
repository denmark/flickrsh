package cmd

import (
	"fmt"

	flickr_photosets "github.com/denmark/flickr/photosets"
	"github.com/spf13/cobra"
)

var albumCmd = &cobra.Command{
	Use:   "album",
	Short: "Get information about albums",
	Long:  "Get information about albums",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("album called")
		initDb()
		albumList()
	},
}

func albumList() {
	client := getFlickrClient()

	pages := 1
	for page := 1; page <= pages; page++ {
		response, err := flickr_photosets.GetList(client, true, "", page)

		if err != nil {
			panic(err)
		}

		pages = response.Photosets.Pages

		for _, photoset := range response.Photosets.Items {
			fmt.Printf("%s [%s]\n", photoset.Id, photoset.Title)
		}

	}
}

func init() {
	rootCmd.AddCommand(albumCmd)
}
