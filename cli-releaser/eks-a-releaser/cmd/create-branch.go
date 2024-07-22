/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

/*
	what does this command do?
	this command is responsible for accessing the trigger file and creating a new release branch in 2 repos,
	the trigger file within the "eks-a-releaser" branch is accessed and its release: contents are extracted
	next, a new branch is created using the extracted release value within the eks-a and build-tooling repo

	Release Process Timeline :
	(1) User first updates trigger file contents within "eks-a-releaser" branch
	(2) User commits changes and raises a PR to be merged into "main" branch
	(3) Codebuild/Pipeline will be triggered once this specific PR is created
	(4) This command will be the first one to be executed and the new release branch will be created
	Moving forward from this point on, all further changes will continue to be committed into the "eks-a-releaser" branch but raised PR's will now target the newly created branch

	MISSING : include code to create same branch also in build-tooling repo
*/

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

func createAnywhereBranch() error {
	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	opts := &github.RepositoryContentGetOptions{
		Ref: "eks-a-releaser", // trigger file is accessed within this branch
	}

	// access trigger file and retrieve content
	triggerFileContentBundleNumber, _, _, err := client.Repositories.GetContents(ctx, PersonalforkedRepoOwner, repoName, triggerFilePath, opts)
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
	// holds release value as strin
	desiredPart := parts[1]

	// Create a new reference for the new branch
	newBranch := desiredPart
	ref := "refs/heads/" + newBranch
	baseRef := "eks-a-releaser" //newly created release branch will be based from this branch
	// future ref : once intergrated into aws repo, baseRef var can := desiredPart - 1 , our new release-0.00 value minus one to be based on previous release branch

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
