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
	cliReleaseNumPath  = "release/triggers/eks-a-release/development/RELEASE_NUMBER"
	cliReleaseVerPath  = "release/triggers/eks-a-release/development/RELEASE_VERSION"
)

// stageCliCmd represents the stageCli command
var stageCliCmd = &cobra.Command{
	Use:   "stage-cli",
	Short: "increments version files to trigger staging bundle release",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. `,

	Run: func(cmd *cobra.Command, args []string) {
		returnedContent, err := updateFileContentsTwoCLI()
		if err != nil {
			fmt.Print(err)
		}
		fmt.Print(returnedContent)
	},
}


func updateFileContentsCLI()(string,error){
	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, cliReleaseNumPath, nil)
	if err != nil {
		fmt.Print(err)
	}

	// variable content holds content for entire file
	content, err := fileContent.GetContent()
	if err != nil {
		fmt.Print(err)
	}

	// returns file value currently 39
	return content, nil
}

func updateFileContentsTwoCLI()(string,error){
	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, cliReleaseVerPath, nil)
	if err != nil {
		fmt.Print(err)
	}

	// variable content holds content for entire file
	content, err := fileContent.GetContent()
	if err != nil {
		fmt.Print(err)
	}

	// returns file value currently v0.16.0
	return content, nil
}

