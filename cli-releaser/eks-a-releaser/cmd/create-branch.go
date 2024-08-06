/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

/*
	what does this command do?

	this command is responsible for accessing the trigger file and creating a new release branch in 2 repos, bot fork of eks-A and build-tooling
	the trigger file within the "eks-a-releaser" branch is accessed and its "release: release-0.00" contents are extracted
	next, a new branch is created using the extracted release value within the eks-anywhere and build-tooling repo

	Release Process Timeline :
	(1) User first updates trigger file contents within "eks-a-releaser" branch ~ bot's fork of eks-A
	(2) Codebuild/Pipeline pulls the latest release version to update/create branch from the trigger file
	(4) This command will be the first one to be executed and the new release branch will be created
*/

import (
	"context"
	"fmt"

	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
)

var (
	buildToolingRepoName = "eks-anywhere-build-tooling"
	upStreamOwner = "testerIbix" // will eventually be replaced by actual upstream owner, aws
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


	
	newBranch := getLatestRelease()
	ref := "refs/heads/" + newBranch
	baseRef := "eks-a-releaser" //newly created release branch will be based from this branch
	// future ref : once intergrated into aws repo, baseRef var can := desiredPart - 1 , our new release-0.00 value minus one to be based on previous release branch

	// Get the reference for the base branch
	baseRefObj, _, err := client.Git.GetRef(ctx, botForkAccount, repoName, "heads/"+baseRef)
	if err != nil {
		return fmt.Errorf("error getting base branch reference: %v", err)
	}

	// Create a new branch
	newBranchRef, _, err := client.Git.CreateRef(ctx, botForkAccount, repoName, &github.Reference{
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
	title := fmt.Sprintf("Release branch %s", newBranch)
	body := "This is a pull request for the new release branch."
	head := fmt.Sprintf("%s:%s", botForkAccount, newBranch)
	base := "main"
	pr, _, err := client.PullRequests.Create(ctx, upStreamOwner, repoName, &github.NewPullRequest{
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
	newBranch := getLatestRelease()
	ref := "refs/heads/" + newBranch
	baseRef := "main" //newly created release branch will be based from this branch
	

	// Get the reference for the base branch
	baseRefObj, _, err := client.Git.GetRef(ctx, botForkAccount, buildToolingRepoName, "heads/"+baseRef)
	if err != nil {
		return fmt.Errorf("error getting base branch reference: %v", err)
	}

	// Create a new branch
	newBranchRef, _, err := client.Git.CreateRef(ctx, botForkAccount, buildToolingRepoName, &github.Reference{
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
	title := fmt.Sprintf("Release branch %s", newBranch)
	body := "This is a pull request for the new release branch."
	head := fmt.Sprintf("%s:%s", botForkAccount, newBranch)
	base := "main"
	pr, _, err := client.PullRequests.Create(ctx, upStreamOwner, buildToolingRepoName, &github.NewPullRequest{
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

// release branches are correctly created on bot's fork
// PR targeting upstream not created correctly 