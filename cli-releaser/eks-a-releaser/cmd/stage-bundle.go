/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
)

var (
	bundleNumPath  = "release/triggers/bundle-release/development/BUNDLE_NUMBER"
	cliMaxVersionPath  = "release/triggers/bundle-release/development/CLI_MAX_VERSION"
	cliMinVersionPath  = "release/triggers/bundle-release/development/CLI_MIN_VERSION"
)

// stageBundleCmd represents the stageBundle command
var stageBundleCmd = &cobra.Command{
	Use:   "stage-bundle",
	Short: "increments version files to trigger staging bundle release",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {
		returnedContent, err := updateFileContentsC()
		if err != nil {
			fmt.Print(err)
		}	
		fmt.Print(returnedContent)
	},
}

func updateFileContentsA() (string, error) {

	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, bundleNumPath, nil)
	if err != nil {
		fmt.Print(err)
	}

	// variable content holds content for entire file
	content, err := fileContent.GetContent()
	if err != nil {
		fmt.Print(err)
	}
	

	// returns file value currently 50
	return content, nil
}


func updateFileContentsB() (string, error) {

	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, cliMaxVersionPath, nil)
	if err != nil {
		fmt.Print(err)
	}

	// variable content holds content for entire file
	content, err := fileContent.GetContent()
	if err != nil {
		fmt.Print(err)
	}
	

	// returns file value currently v0.18.0
	return content, nil
}

func updateFileContentsC() (string, error) {

	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, cliMinVersionPath, nil)
	if err != nil {
		fmt.Print(err)
	}

	// variable content holds content for entire file
	content, err := fileContent.GetContent()
	if err != nil {
		fmt.Print(err)
	}
	

	// returns file value currently v0.18.0
	return content, nil
}



