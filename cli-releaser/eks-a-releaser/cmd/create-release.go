/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

/*
	what does this command do?

	this command is responsible for creating a release tag with the commit hash that triggered the prod CLI release

	func retrieveLatestProdCLIHash() - retrieves the latest commit hash from the prod release version file, "eks-a-releaser" branch - bot's forked repo

	func createTag() - takes in commit hash and creates a tag

	func createGitHubRelease() - creates a release on GitHub using the tag created in createTag()

	func runBothTag() - runs both createTag() and createGitHubRelease()
*/

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// createReleaseCmd represents the createRelease command
var createReleaseCmd = &cobra.Command{
	Use:   "create-release",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {
		runBothTag()
	},
}

func runBothTag(){

	//retrieve commit hash 
	commitHash := retrieveLatestProdCLIHash()
	
	//create tag with commit hash
	tag, errOne := createTag(commitHash)
	if errOne != nil {
		log.Panic(errOne)
	}

	rel, errTwo := createGitHubRelease(tag)
	if errTwo != nil {
		log.Panic(errTwo)
	}

	err := createReleasePR(rel)
	if err != nil {
		log.Panic(err)
	}

	//print release object
	fmt.Print(rel)
}


// creates tag using retrieved commit hash
func createTag(commitHash string) (*github.RepositoryRelease, error){
	
	// retrieve tag name "v0.0.00" from trigger file, "eks-a-releaser" branch 
	latestVersionValue := os.Getenv("LATEST_VERSION")

	//create client
	accessToken := os.Getenv("SECRET_PAT")
	ctx := context.Background()

	// Create a new GitHub client instance with the token type set to "Bearer"
	ts := oauth2.StaticTokenSource(&oauth2.Token{
    	AccessToken: accessToken,
    	TokenType:   "Bearer",
	})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	releaseName := latestVersionValue
	releaseDesc := "EKS-Anywhere " + latestVersionValue + " release"
	commitSHA := commitHash
	release := &github.RepositoryRelease{
    	TagName: github.String(releaseName),
    	Name:    github.String(releaseName),
    	Body:    github.String(releaseDesc),
		TargetCommitish: github.String(commitSHA),
	}

	rel, _, err := client.Repositories.CreateRelease(ctx, forkedRepoAccount, EKSAnyrepoName, release)
	if err != nil {
		fmt.Printf("error creating release: %v", err)
	}

	fmt.Printf("Release tag %s created successfully!\n", rel.GetTagName())
	return rel, nil
}



func retrieveLatestProdCLIHash() string {
	
	//create client
	accessToken := os.Getenv("SECRET_PAT")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	opts := &github.CommitsListOptions{
        Path: prodCliReleaseVerPath, // file to check
        SHA:  "eks-a-releaser", // branch to check
    }

	
	commits, _, err := client.Repositories.ListCommits(ctx, forkedRepoAccount, EKSAnyrepoName, opts)
    if err != nil {
        return "error fetching commits list"
    }


	if len(commits) > 0 {
        latestCommit := commits[0]
        return latestCommit.GetSHA()
    }

    return "no commits found for file"
}


func createGitHubRelease(releaseTag *github.RepositoryRelease) (*github.RepositoryRelease, error){
	
	latestVersionValue := os.Getenv("LATEST_VERSION")

	//create client
	accessToken := os.Getenv("SECRET_PAT")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	release, _, err := client.Repositories.GetReleaseByTag(ctx, forkedRepoAccount, EKSAnyrepoName, latestVersionValue)
    if err == nil {
        fmt.Printf("Release %s already exists!\n", latestVersionValue)
        return release, nil
    }

	release = &github.RepositoryRelease{
        TagName: releaseTag.TagName,
        Name:    &latestVersionValue,
        Body:    releaseTag.Body,
    }

    rel, _, err := client.Repositories.CreateRelease(ctx, forkedRepoAccount, EKSAnyrepoName, release)
    if err != nil {
        return nil, err
    }

    return rel, nil
}



func createReleasePR(release *github.RepositoryRelease) error {

	//create client
	latestRelease := os.Getenv("LATEST_RELEASE")
	accessToken := os.Getenv("SECRET_PAT")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)


	
	
	// Prepare pull request details
	title := fmt.Sprintf("Release %s", release.GetTagName())
	body := fmt.Sprintf("This pull request contains the release %s", release.GetTagName())
	head := fmt.Sprintf("%s:%s", forkedRepoAccount, "eks-a-releaser")
	base := latestRelease

	fmt.Printf("Head parameter value: %s\n", head)
	// Create a pull request
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