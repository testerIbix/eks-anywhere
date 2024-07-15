/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
)

// createBranchCmd represents the createBranch command
var createBranchCmd = &cobra.Command{
	Use:   "create-branch",
	Short: "Creates new release branch from updated trigger file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {

		err := createAnywhereBranch()
			if err != nil {
				fmt.Print(err)
		}
	},
}


func createAnywhereBranch()(error){
	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	// access trigger file and retrieve content
	triggerFileContentBundleNumber,_,_, err := client.Repositories.GetContents(ctx, PersonalforkedRepoOwner, repoName, triggerFilePath, nil)
	if err != nil {
		return fmt.Errorf("first breakpoint %s", err)
	}
	content, err := triggerFileContentBundleNumber.GetContent()
	if err != nil {
		return fmt.Errorf("second breakpoint %s", err)
	}
	// Find the line containing the identifier
	snippetStartIdentifierB := "release: "
	lines := strings.Split(content, "\n")
	startIndex := -1
	endIndex := -1
	for i, line := range lines {
		if strings.Contains(line, snippetStartIdentifierB) {
			startIndex = i
			endIndex = i // Set endIndex to the same line as startIndex
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		log.Print("snippet not found")
	}
	// holds full string 
	bundleNumberLine := lines[startIndex]
	// split string to isolate bundle number
	parts := strings.Split(bundleNumberLine, ": ")
	// holds bundle number value as string 
	desiredPart := parts[1]

	// Create a new reference for the new branch
	newBranch := desiredPart
	ref := "refs/heads/" + newBranch
	baseRef := "eks-a-releaser" //base branch from which new branch will be created 

	// Get the reference for the base branch
	baseRefObj, _, err := client.Git.GetRef(ctx, PersonalforkedRepoOwner, repoName, "heads/"+baseRef)
	if err != nil {
    	return fmt.Errorf("error getting base branch reference: %v", err)
	}

	// Create a new branch
	newBranchRef, _, err := client.Git.CreateRef(ctx, PersonalforkedRepoOwner, repoName, &github.Reference{
 	Ref: &ref,
    Object: &github.GitObject{
        SHA: baseRefObj.Object.SHA,
    	},
	})
	if err != nil {
    return fmt.Errorf("error creating branch: %v", err)
	}

	fmt.Printf("New branch '%s' created successfully\n", *newBranchRef.Ref)
	return nil

}


/*
User will access the trigger file within the releaser branch and raise a PR to the main branch 
PR creation will invoke pipeline, first triggering this create-branch command 
create-branch accesses the first line of the trigger file and creates a new branch with the retrieved contents 
command creates new release branch within eks-anywhere & eks-build-tooling repository 
*/