/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

/*
	what does this command do?

	this command is responsible for accessing the trigger file and creating a new release branch in 2 repos, bot's fork of eks-Anywhere and build-tooling
	the trigger file within the "eks-a-releaser" branch of the bot's fork is accessed and its "release: release-0.00" contents are extracted
	next, a new branch is created using the extracted release value within the eks-anywhere and build-tooling repo

	branches are created on bot's fork
	PR fails to reflect newly created branches on upstream repo
*/

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
)

var (
	buildToolingRepoName = "eks-anywhere-build-tooling"
	upStreamRepoOwner = "testerIbix" // will eventually be replaced by actual upstream owner, aws
)

// createBranchCmd represents the createBranch command
var createBranchCmd = &cobra.Command{
	Use:   "create-branch",
	Short: "Creates new release branch from updated trigger file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {

		err := createBoth()
		if err != nil {
			fmt.Print(err)
		}
	},
}

func createBoth()error{
	errOne := createAnywhereBranch()
	if errOne != nil {
		return fmt.Errorf("error calling createAnywhereBranch %s", errOne)
	}

	errTwo := createBuildToolingBranch()
	if errTwo != nil{
		return fmt.Errorf("error calling createBuildToolingBranch %s", errTwo)
	}

	return nil
}


func createAnywhereBranch() error {

	//create client
	secretName := "Secret"
	accessToken, err := getSecretValue(secretName)
	if err != nil {
		fmt.Print("error getting secret", err)
	}
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)


	
	latestRelease := os.Getenv("LATEST_RELEASE")

	ref := "refs/heads/" + latestRelease
	baseRef := "eks-a-releaser" //newly created release branch will be based from this branch
	

	// Get the reference for the base branch
	baseRefObj, _, err := client.Git.GetRef(ctx, forkedRepoAccount, EKSAnyrepoName, "heads/"+baseRef)
	if err != nil {
		return fmt.Errorf("error getting base branch reference: %v", err)
	}

	// Create a new branch
	newBranchRef, _, err := client.Git.CreateRef(ctx, forkedRepoAccount, EKSAnyrepoName, &github.Reference{
		Ref: &ref,
		Object: &github.GitObject{
			SHA: baseRefObj.Object.SHA,
		},
	})
	if err != nil {
		return fmt.Errorf("error creating branch: %v", err)
	}

	
	fmt.Printf("New branch '%s' created successfully\n", *newBranchRef.Ref)
	

	// create pull request targeting upstream eks-A repo
	title := fmt.Sprintf("Release branch %s", latestRelease)
	body := "This is a pull request for the new release branch."
	head := fmt.Sprintf("%s:%s", forkedRepoAccount, latestRelease)
	base := "main"
	pr, _, err := client.PullRequests.Create(ctx, upStreamRepoOwner, EKSAnyrepoName, &github.NewPullRequest{
		Title: &title,
		Body:  &body,
		Head:  &head,
		Base:  &base,
	})
	if err != nil {
		return fmt.Errorf("error creating pull request: %v", err)
	}

	fmt.Printf("Pull request created: %s\n", pr.GetHTMLURL())
	return nil

}




// build tooling branch created is based off "main"
func createBuildToolingBranch() error {
	
	//create client
	secretName := "Secret"
	accessToken, err := getSecretValue(secretName)
	if err != nil {
		fmt.Print("error getting secret", err)
	}
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)


	// Create a new reference for the new branch
	latestRelease := os.Getenv("LATEST_RELEASE")
	ref := "refs/heads/" + latestRelease
	baseRef := "main" //newly created release branch will be based from this branch
	

	// Get the reference for the base branch
	baseRefObj, _, err := client.Git.GetRef(ctx, forkedRepoAccount, buildToolingRepoName, "heads/"+baseRef)
	if err != nil {
		return fmt.Errorf("error getting base branch reference: %v", err)
	}

	// Create a new branch
	newBranchRef, _, err := client.Git.CreateRef(ctx, forkedRepoAccount, buildToolingRepoName, &github.Reference{
		Ref: &ref,
		Object: &github.GitObject{
			SHA: baseRefObj.Object.SHA,
		},
	})
	if err != nil {
		return fmt.Errorf("error creating branch: %v", err)
	}

	
	fmt.Printf("New branch '%s' created successfully\n", *newBranchRef.Ref)
	



	// create pull request targeting upstream build-tooling repo
	title := fmt.Sprintf("Release branch %s", latestRelease)
	body := "This is a pull request for the new release branch."
	head := fmt.Sprintf("%s:%s", forkedRepoAccount, latestRelease)
	base := "main"
	pr, _, err := client.PullRequests.Create(ctx, upStreamRepoOwner, buildToolingRepoName, &github.NewPullRequest{
		Title: &title,
		Body:  &body,
		Head:  &head,
		Base:  &base,
	})
	if err != nil {
		return fmt.Errorf("error creating pull request: %v", err)
	}

	fmt.Printf("Pull request created: %s\n", pr.GetHTMLURL())
	return nil

}