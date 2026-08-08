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
	"path/filepath"
	"runtime"
	"strings"

	flickr "github.com/denmark/flickr"
	flickr_test "github.com/denmark/flickr/test"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	yaml "go.yaml.in/yaml/v3"
)

type flickrConfig struct {
	UserID           string `yaml:"user_id"`
	APIKey           string `yaml:"api_key"`
	APISecret        string `yaml:"api_secret"`
	OauthToken       string `yaml:"oauth_token"`
	OauthTokenSecret string `yaml:"oauth_token_secret"`
}

var flickrAPIKey string
var flickrAPISecret string

// configFilePath returns the OS-appropriate path to the Flickr credentials
// file, which lives in the user's home directory.
func configFilePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	filename := ".flickrsh.yml"
	if runtime.GOOS == "windows" {
		filename = "flickrsh.yml"
	}

	return filepath.Join(home, filename), nil
}

func loadFlickrConfig() (*flickrConfig, error) {
	path, err := configFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg flickrConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func saveFlickrConfig(cfg *flickrConfig) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func getFlickrClient() (client *flickr.FlickrClient) {
	cfg, err := loadFlickrConfig()
	if err != nil || cfg.APIKey == "" {
		fmt.Println("Flickr API not initialized, run 'flickrsh init --key <API Key> --secret <API Secret>'")
		return nil
	}

	client = flickr.NewFlickrClient(cfg.APIKey, cfg.APISecret)
	client.OAuthToken = cfg.OauthToken
	client.OAuthTokenSecret = cfg.OauthTokenSecret

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
	if err != nil {
		fmt.Println("Failed to initialize the Flickr connection!!!")
		return
	}

	client.OAuthToken = accessToken.OAuthToken
	client.OAuthTokenSecret = accessToken.OAuthTokenSecret

	resp, err := flickr_test.Login(client)
	if err != nil {
		fmt.Printf("Login to Flickr failed: %s\n", err)
		return
	}

	cfg := &flickrConfig{
		UserID:           resp.User.ID,
		APIKey:           flickrAPIKey,
		APISecret:        flickrAPISecret,
		OauthToken:       client.OAuthToken,
		OauthTokenSecret: client.OAuthTokenSecret,
	}

	if err := saveFlickrConfig(cfg); err != nil {
		fmt.Printf("Failed to save Flickr configuration: %s\n", err)
		return
	}

	path, _ := configFilePath()
	fmt.Printf("Flickr credentials saved to %s\n", path)
}

// promptOverwrite asks the user whether an existing configuration file
// should be overwritten.
func promptOverwrite(path string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("A Flickr configuration already exists at %s. Overwrite it? [y/N]: ", path)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Flickr API connection.",
	Long:  "Initialize the Flickr API connection.",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := configFilePath()
		if err != nil {
			fmt.Printf("Unable to determine configuration file location: %s\n", err)
			return
		}

		if _, err := os.Stat(path); err == nil {
			if !promptOverwrite(path) {
				fmt.Println("Keeping existing configuration; aborting.")
				return
			}
		}

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
