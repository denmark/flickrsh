package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	flickr_photos "github.com/denmark/flickr/photos"
	flickr_photosets "github.com/denmark/flickr/photosets"
	"github.com/spf13/cobra"
)

var albumId string
var minHeight int
var downloadDirectory string

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download photo/video files from Flickr",
	Long:  "Download photo/video files from Flickr",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("download called")
		download()
	},
}

func download() {
	client := getFlickrClient()

	if client == nil {
		return
	}

	pages := 1
	photoIdx := 0
	for page := 1; page <= pages; page++ {
		response, err := flickr_photosets.GetPhotos(client, true, albumId, "", page)

		if err != nil {
			panic(err)
		}

		pages = response.Photoset.Pages

		for _, photoInfo := range response.Photoset.Photos {
			fmt.Printf("%s [%s]\n", photoInfo.Id, photoInfo.Title)
			photo, err := flickr_photos.GetInfo(client, photoInfo.Id, "")

			if err != nil {
				panic(err)
			}
			dateTaken := strings.ReplaceAll(photo.Photo.Dates.Taken, " ", "_")
			fmt.Println(dateTaken)

			photoSizesResponse, err := flickr_photos.GetSizes(client, photoInfo.Id)

			if err != nil {
				panic(err)
			}

			imageURL := ""
			for _, size := range photoSizesResponse.PhotoSizes.Sizes {
				// fmt.Printf("[%s] URL:%s (%sx%s) %s\n", size.Label, size.Source, size.Height, size.Width, size.Media)
				height, err := strconv.Atoi(size.Height)
				if err == nil {
					imageURL = size.Source
					if height >= minHeight {
						fmt.Printf("found size: %d >= %d [%s]\n", height, minHeight, size.Source)
						imageURL = size.Source
						break
					}
				} else {
					break
				}
			}
			fmt.Printf("imageURL: %s\n", imageURL)

			err = downloadFileFromURL(imageURL, downloadDirectory, dateTaken, photoIdx)
			if err != nil {
				panic(err)
			}
			photoIdx += 1
		}
	}
}

func downloadFileFromURL(fileURL string, targetDir string, filePrefix string, photoIdx int) error {
	url, err := url.Parse(fileURL)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/%04d_%s_%s", targetDir, photoIdx, filePrefix, filepath.Base(url.Path))
	fmt.Printf("filename: %s\n", filename)

	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("File exists\n")
		return nil
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	tmpFile := "./___tmpdownloadfile__"
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	err = os.Rename(tmpFile, filename)
	return err
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVar(&albumId, "album", "", "The Album from which we should download photos")
	downloadCmd.Flags().IntVar(&minHeight, "minheight", 0, "The size to download")
	downloadCmd.Flags().StringVar(&downloadDirectory, "dir", ".", "Where to download the files")
}
