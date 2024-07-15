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
	AWSrepoOwner = "aws"
	repoName  = "eks-anywhere"
	makeFilePath  = "/Makefile"
)

// upMakeFileCmd represents the upMakeFile command
var updateMakefileCmd = &cobra.Command{
	Use:   "update-Makefile",
	Short: "Updates BRANCH_NAME variable to match new release branch within the Makefile",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {
		content := updateMakefileContents()
		fmt.Print(content)
	},
}

func updateMakefileContents() (string) {

	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, PersonalforkedRepoOwner, repoName, makeFilePath, nil)
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


	return snippet



	// until method of pointing to new release branch is established 
	// snippet will be altered with a 

	// get latest commit sha
// 	ref, _, err := client.Git.GetRef(ctx, forkedRepoOwner, repoName, "heads/eks-a-releaser")
// 	if err != nil {
// 		return fmt.Errorf("error getting ref %s", err)
// 	}
// 	latestCommitSha := ref.Object.GetSHA()

// 	entries := []*github.TreeEntry{}
// 	entries = append(entries, &github.TreeEntry{Path: github.String(strings.TrimPrefix(makeFilePath, "/")), Type: github.String("blob"), Content: github.String(string(snippet)), Mode: github.String("100644")})
// 	tree, _, err := client.Git.CreateTree(ctx,forkedRepoOwner, repoName, *ref.Object.SHA, entries)
// 	if err != nil {
// 	 	return fmt.Errorf("error creating tree %s", err)
// 	}

// 	//validate tree sha
// 	newTreeSHA := tree.GetSHA()

// 	// create new commit
// 	author := &github.CommitAuthor{
// 	Name:  github.String("ibix16"),
// 	Email: github.String("ibixrivera16@gmail.com"),
// 	}

// 	commit := &github.Commit{
// 	Message: github.String("Update BRANCH_NAME variable within Makefile"),
// 	Tree:    &github.Tree{SHA: github.String(newTreeSHA)},
// 	Author:  author,
// 	Parents: []*github.Commit{{SHA: github.String(latestCommitSha)}},
// 	}

// 	commitOP := &github.CreateCommitOptions{}
// 	newCommit, _, err := client.Git.CreateCommit(ctx, forkedRepoOwner, repoName, commit, commitOP)
// 	if err != nil {
// 	return fmt.Errorf("creating commit %s", err)
// 	}
// 	newCommitSHA := newCommit.GetSHA()
	
// 	// update branch reference
// 	ref.Object.SHA = github.String(newCommitSHA)

// 	_, _, err = client.Git.UpdateRef(ctx, forkedRepoOwner, repoName, ref, false)
// 	if err != nil {
// 	return fmt.Errorf("error updating ref %s", err)
// 	}

// 	// create pull request
//     base := "main"
//     head := fmt.Sprintf("%s:%s", forkedRepoOwner, "eks-a-releaser")
//     title := "Update BRANCH_NAME variable found within the Makefile"
//     body := "This pull request is responsible for updating the BRANCH_NAME variable found within the Makefile in order to make it point to the new release"

//     newPR := &github.NewPullRequest{
//         Title: &title,
//         Head:  &head,
//         Base:  &base,
//         Body:  &body,
//     }s

// 	pr, _, err := client.PullRequests.Create(ctx, forkedRepoOwner, repoName, newPR)
//     if err != nil {
//         return fmt.Errorf("error creating PR %s", err)
//     }

// 	log.Printf("Pull request created: %s\n", pr.GetHTMLURL())
// 	return nil

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
