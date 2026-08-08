// Copyright © 2018 NAME HERE <denmark@gmail.com>
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
	"bufio"
	"fmt"
	"os"
	"strings"

	flickr "github.com/denmark/flickr"
	flickr_test "github.com/denmark/flickr/test"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var db *gorm.DB

type auth struct {
	UserID           string `gorm:"type:varchar(25);primary key"`
	APIKey           string `gorm:"type:varchar(50)"`
	APISecret        string `gorm:"type:varchar(50)"`
	OauthToken       string `gorm:"type:varchar(50)"`
	OauthTokenSecret string `gorm:"type:varchar(50)"`
}

var flickrAPIKey string
var flickrAPISecret string

func initDb() {
	var err error

	db, err = gorm.Open("sqlite3", "flickrsh.db")
	if err != nil {
		panic(err)
	}

	db.LogMode(true) // true gives SQL output
	db.SingularTable(true)
	db.AutoMigrate(&auth{})
	fmt.Println("Initialized DB")
}

func getFlickrClient() (client *flickr.FlickrClient) {
	var flickrAuth auth

	db.Take(&flickrAuth)

	if flickrAuth.APIKey == "" {
		fmt.Println("Flickr API not initialized, run 'init flickr --key <API Key> --secret <API Secret>'")
		return nil
	}

	client = flickr.NewFlickrClient(flickrAuth.APIKey, flickrAuth.APISecret)
	client.OAuthToken = flickrAuth.OauthToken
	client.OAuthTokenSecret = flickrAuth.OauthTokenSecret

	return
}

func initFlickr() {
	client := flickr.NewFlickrClient(flickrAPIKey, flickrAPISecret)
	requestTok, _ := flickr.GetRequestToken(client)
	url, _ := flickr.GetAuthorizeUrl(client, requestTok)
	fmt.Printf("apiKey: %s; apiSecret: %s\n", flickrAPIKey, flickrAPISecret)
	fmt.Println(requestTok)
	fmt.Println(url)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter oAuth Confirmation Code: ")
	oauthConfirmationCode, _ := reader.ReadString('\n')

	accessToken, err := flickr.GetAccessToken(client, requestTok, strings.TrimSpace(oauthConfirmationCode))
	if err == nil {
		client.OAuthToken = accessToken.OAuthToken
		client.OAuthTokenSecret = accessToken.OAuthTokenSecret

		if resp, err := flickr_test.Login(client); err == nil {
			var flickrAuth auth
			db.Where(auth{UserID: resp.User.ID}).Assign(auth{
				UserID:           resp.User.ID,
				APIKey:           flickrAPIKey,
				APISecret:        flickrAPISecret,
				OauthToken:       client.OAuthToken,
				OauthTokenSecret: client.OAuthTokenSecret,
			}).FirstOrCreate(&flickrAuth)
		} else {
			fmt.Printf("Login to Flickr failed: %s\n", err)
		}
	} else {
		fmt.Println("Failed to initialize the Flickr connection!!!")
	}
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Flickr API connection.",
	Long:  "Initialize the Flickr API connection.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		initDb()
		initFlickr()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&flickrAPIKey, "key", "", "API Key")
	initCmd.Flags().StringVar(&flickrAPISecret, "secret", "", "API Secret")
	initCmd.MarkFlagRequired("key")
	initCmd.MarkFlagRequired("secret")
}
