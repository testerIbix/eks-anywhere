/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

/*

what does this function do?

this function is used to get the user's github account username

the username is extracted by locating the account associated with the passed in PAT
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// userTestCmd represents the userTest command
var userTestCmd = &cobra.Command{
	Use:   "user-test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		user := getAuthenticatedUsername()
		fmt.Print(user)
	},
}

// User represents the user's GitHub account information.
type User struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	Name              string `json:"name"`
	Company           string `json:"company"`
	Blog              string `json:"blog"`
	Location          string `json:"location"`
	Email             string `json:"email"`
	Hireable          bool   `json:"hireable"`
	Bio               string `json:"bio"`
	TwitterUsername   string `json:"twitter_username"`
	PublicRepos       int    `json:"public_repos"`
	PublicGists       int    `json:"public_gists"`
	Followers         int    `json:"followers"`
	Following         int    `json:"following"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

func getAuthenticatedUsername() string {

	// Get the GitHub personal access token from aws secrets manager
	// secretName := "Secret"
	// accessToken, err := getSecretValue(secretName)
	// if err != nil {
	// 	fmt.Print("error getting secret", err)
	// }

	accessToken := os.Getenv("SECRET_PAT")

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return "error creating HTTP request"
	}

	// Set the authorization header with the personal access token
	req.Header.Set("Authorization", "token "+accessToken)

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "error sending HTTP request"
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "error reading response body"
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "failed to retrieve user information"
	}

	// Unmarshal the response body into a User struct
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return "error unmarshalling"
	}

	stringUser := user.Login
	return stringUser
}

func init() {
	rootCmd.AddCommand(userTestCmd)
}
