// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"time"

	flickr "github.com/denmark/flickr"

	"github.com/spf13/cobra"
)

var album string
var tags []string

var isFamily bool
var isFriend bool
var isPublic bool
var isSearchable bool

const maxRetryAttempts = 3

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload one or more files to Flickr",
	Long:  "Upload one or more files to Flickr",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("upload called")
		upload(args)
	},
}

func upload(files []string) {
	// var progressBar *pb.ProgressBar

	fmt.Printf("files: %v\n", files)

	for i := 0; i < len(tags); i++ {
		fmt.Printf("tag: %s\n", tags[i])
	}
	client := getFlickrClient()

	if client == nil {
		return
	}

	uploadParams := flickr.UploadParams{
		Tags: tags,
		Hidden: func() int {
			if isSearchable {
				return 1
			}
			return 2
		}(),
		IsFamily: isFamily,
		IsFriend: isFriend,
		IsPublic: isPublic,
	}

	totalFiles := len(files)

	fmt.Printf("Beginning upload of %d file%s...\n", totalFiles, (func() string {
		if totalFiles != 1 {
			return "s"
		} else {
			return ""
		}
	})())

	retryAttempts := maxRetryAttempts

	for i := 0; i < totalFiles; i++ {
		startTime := time.Now()
		fmt.Printf("Uploading [%s]...", files[i])
		if resp, err := flickr.UploadFile(client, files[i], &uploadParams); err == nil {
			percentComplete := (float64(i+1) / float64(totalFiles)) * 100
			uploadTime := time.Since(startTime) / time.Second
			fmt.Printf("[%s] [%ds] [%.2f%%]\n", resp.ID, uploadTime, percentComplete)
			retryAttempts = maxRetryAttempts
		} else {
			if retryAttempts > 0 {
				retryAttempts--
				fmt.Printf("\n  Error uploading [%s], retrying [%d]...\n", files[i], retryAttempts)
				i--
				continue
			} else {
				fmt.Printf("  Error uploading [%s], exiting...\n", files[i])
				break
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringArrayVar(&tags, "tag", nil, "Tag to apply to uploaded photos/videos")
	uploadCmd.Flags().StringVar(&album, "album", "", "The Album into which these files should go")
	uploadCmd.Flags().BoolVar(&isSearchable, "searchable", false, "Make photos available to public search")
	uploadCmd.Flags().BoolVar(&isFamily, "family", false, "Make uploaded photos available to family")
	uploadCmd.Flags().BoolVar(&isFriend, "friend", false, "Make uploaded photos available to your friends")
	uploadCmd.Flags().BoolVar(&isPublic, "public", false, "Make uploaded photos available to the public")
}
