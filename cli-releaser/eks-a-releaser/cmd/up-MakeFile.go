/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
)

var (
	repoOwner = "aws"
	repoName  = "eks-anywhere"
	makeFilePath  = "/Makefile"
)

// upMakeFileCmd represents the upMakeFile command
var upMakeFileCmd = &cobra.Command{
	Use:   "up-MakeFile",
	Short: "accesses MakeFile & updates BRANCH_NAME variable to match new release branch",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {
		returnedContent, err := updateFileContents()
		if err != nil {
			fmt.Print(err)
		}

		fmt.Print(returnedContent)
	},
}

/*
approach : retrieve entire file contents
		   extract desired code snippet
		   modify code snippet
		   replace old code snippet with new in entire file
		   commit changes + raise PR
*/


func updateFileContents() (string, error) {

	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, makeFilePath, nil)
	if err != nil {
		fmt.Print(err)
	}


	// variable content holds content for entire file
	content, err := fileContent.GetContent()
	if err != nil {
		fmt.Print(err)
	}

	// variable snippet holds desired code snippet
	snippet := extractCode(content)

	return snippet, nil
}

// retrieves extracted code snippet, called in getFileContents
func extractCode(content string) string {
	// extract code snippet from content
	snippetStartIdentifier := "BRANCH_NAME?="
	lines := strings.Split(content, "\n")
	startIndex := -1
	endIndex := -1

	for i, line := range lines {
		if strings.Contains(line, snippetStartIdentifier) {
			startIndex = i
			endIndex = i // Set endIndex to the same line as startIndex
			break
		}
	}

	// return code snippet
	if startIndex != -1 && endIndex != -1 {
		return lines[startIndex]
	}
	return "error code snippet not found" // Snippet not found
}

/*
Notes : 7/5/24
		have constructed functions for extracting code snippet

Future work :
	1. create functions for modifying code snippet
		- establish where we will be pulling new release info from in order
		to update the variable
		- last update to the makefile was to point to "release-0.20"
		manifest file contains the latest release version but in the form of v0.20.0
			- could potentially pull from manifest but would require further parsing and editing
		- UPDATE : this function will pull down the new release version to point to from the release branch created
			- cannot complete upMakeFile command until the cut release branch command is complete since we want our
			variable to be pulled from the new release branch name
*/
