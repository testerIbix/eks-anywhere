/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

/*
	what does this command do?

	this command is responsible for staging bundle release

	Fetches env variables passed in from pipeline UI

	The pull request is created using the 3 commits from "eks-a-releaser" branch forked repo and is intended to be merged into latest release branch on the upstream repo
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

var (
	bundleNumPath           = "release/triggers/bundle-release/development/BUNDLE_NUMBER"
	cliMaxVersionPath       = "release/triggers/bundle-release/development/CLI_MAX_VERSION"
	cliMinVersionPath       = "release/triggers/bundle-release/development/CLI_MIN_VERSION"
	//triggerFilePath         = "release/triggers/eks-a-releaser-trigger"
	forkedRepoAccount = getAuthenticatedUsername()
)

// stageBundleCmd represents the stageBundle command
var stageBundleCmd = &cobra.Command{
	Use:   "stage-bundle",
	Short: "creates a PR containing 3 commits, each updating the contents of a singular file intended for staging bundle release",
	Long: `Retrieves updated content for development : bundle number, cli max version, and cli min version. 
	Writes the updated changes to the 3 files and raises a PR with the 3 commits.`,

	Run: func(cmd *cobra.Command, args []string) {
		err := runAllStagebundle()
		if err != nil {
			log.Print(err)
		}
	},
}

func runAllStagebundle() error {
	errOne := updateBundleNum()
	if errOne != nil {
		return errOne
	}

	errTwo := updateCLIMax()
	if errTwo != nil {
		return errTwo
	}

	errThree := updateCLIMin()
	if errThree != nil {
		return errThree
	}

	errFive := createPullRequestStageBundleTwo()
	if errFive != nil {
		return errFive
	}

	return nil
}

func updateBundleNum() error {

	//create client
	// secretName := "Secret"
	// accessToken, err := getSecretValue(secretName)
	// if err != nil {
	// 	fmt.Print("error getting secret", err)
	// }

	accessToken := os.Getenv("SECRET_PAT")

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	bundleNumber := os.Getenv("RELEASE_NUMBER")


	// get latest commit sha from branch "eks-a-releaser"
	ref, _, err := client.Git.GetRef(ctx, forkedRepoAccount, EKSAnyrepoName, "heads/eks-a-releaser")
	if err != nil {
		return fmt.Errorf("error getting ref %s", err)
	}
	latestCommitSha := ref.Object.GetSHA()

	entries := []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{Path: github.String(strings.TrimPrefix(bundleNumPath, "/")), Type: github.String("blob"), Content: github.String(string(bundleNumber)), Mode: github.String("100644")})
	tree, _, err := client.Git.CreateTree(ctx, forkedRepoAccount, EKSAnyrepoName, *ref.Object.SHA, entries)
	if err != nil {
		return fmt.Errorf("error creating tree %s", err)
	}

	//validate tree sha
	newTreeSHA := tree.GetSHA()

	// create new commit, update email address
	author := &github.CommitAuthor{
		Name:  github.String("ibix16"),
		Email: github.String("fake@wtv.com"),
	}

	commit := &github.Commit{
		Message: github.String("Increment bundle number file"),
		Tree:    &github.Tree{SHA: github.String(newTreeSHA)},
		Author:  author,
		Parents: []*github.Commit{{SHA: github.String(latestCommitSha)}},
	}

	commitOP := &github.CreateCommitOptions{}
	newCommit, _, err := client.Git.CreateCommit(ctx, forkedRepoAccount, EKSAnyrepoName, commit, commitOP)
	if err != nil {
		return fmt.Errorf("creating commit %s", err)
	}
	newCommitSHA := newCommit.GetSHA()

	// update branch reference
	ref.Object.SHA = github.String(newCommitSHA)

	_, _, err = client.Git.UpdateRef(ctx, forkedRepoAccount, EKSAnyrepoName, ref, false)
	if err != nil {
		return fmt.Errorf("error updating ref %s", err)
	}

	return nil
}

func updateCLIMax() error {

	//create client
	accessToken := os.Getenv("SECRET_PAT")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	latestVersion := os.Getenv("LATEST_VERSION")


	// get latest commit sha from branch "eks-a-releaser"
	ref, _, err := client.Git.GetRef(ctx, forkedRepoAccount, EKSAnyrepoName, "heads/eks-a-releaser")
	if err != nil {
		return fmt.Errorf("error getting ref %s", err)
	}
	latestCommitSha := ref.Object.GetSHA()

	entries := []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{Path: github.String(strings.TrimPrefix(cliMaxVersionPath, "/")), Type: github.String("blob"), Content: github.String(string(latestVersion)), Mode: github.String("100644")})
	tree, _, err := client.Git.CreateTree(ctx, forkedRepoAccount, EKSAnyrepoName, *ref.Object.SHA, entries)
	if err != nil {
		return fmt.Errorf("error creating tree %s", err)
	}

	//validate tree sha
	newTreeSHA := tree.GetSHA()

	// create new commit, update email address
	author := &github.CommitAuthor{
		Name:  github.String("ibix16"),
		Email: github.String("fake@wtv.com"),
	}

	commit := &github.Commit{
		Message: github.String("Update CLI Max version"),
		Tree:    &github.Tree{SHA: github.String(newTreeSHA)},
		Author:  author,
		Parents: []*github.Commit{{SHA: github.String(latestCommitSha)}},
	}

	commitOP := &github.CreateCommitOptions{}
	newCommit, _, err := client.Git.CreateCommit(ctx, forkedRepoAccount, EKSAnyrepoName, commit, commitOP)
	if err != nil {
		return fmt.Errorf("creating commit %s", err)
	}
	newCommitSHA := newCommit.GetSHA()

	// update branch reference
	ref.Object.SHA = github.String(newCommitSHA)

	_, _, err = client.Git.UpdateRef(ctx, forkedRepoAccount, EKSAnyrepoName, ref, false)
	if err != nil {
		return fmt.Errorf("error updating ref %s", err)
	}
	return nil
}

func updateCLIMin() error {

	//create client
	accessToken := os.Getenv("SECRET_PAT")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	latestVersion := os.Getenv("LATEST_VERSION")


	// get latest commit sha from branch "eks-a-releaser"
	ref, _, err := client.Git.GetRef(ctx, forkedRepoAccount, EKSAnyrepoName, "heads/eks-a-releaser")
	if err != nil {
		return fmt.Errorf("error getting ref %s", err)
	}
	latestCommitSha := ref.Object.GetSHA()

	entries := []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{Path: github.String(strings.TrimPrefix(cliMinVersionPath, "/")), Type: github.String("blob"), Content: github.String(string(latestVersion)), Mode: github.String("100644")})
	tree, _, err := client.Git.CreateTree(ctx, forkedRepoAccount, EKSAnyrepoName, *ref.Object.SHA, entries)
	if err != nil {
		return fmt.Errorf("error creating tree %s", err)
	}

	//validate tree sha
	newTreeSHA := tree.GetSHA()

	// create new commit, update email address
	author := &github.CommitAuthor{
		Name:  github.String("ibix16"),
		Email: github.String("fake@wtv.com"),
	}

	commit := &github.Commit{
		Message: github.String("Update CLI Min version"),
		Tree:    &github.Tree{SHA: github.String(newTreeSHA)},
		Author:  author,
		Parents: []*github.Commit{{SHA: github.String(latestCommitSha)}},
	}

	commitOP := &github.CreateCommitOptions{}
	newCommit, _, err := client.Git.CreateCommit(ctx, forkedRepoAccount, EKSAnyrepoName, commit, commitOP)
	if err != nil {
		return fmt.Errorf("creating commit %s", err)
	}
	newCommitSHA := newCommit.GetSHA()

	// update branch reference
	ref.Object.SHA = github.String(newCommitSHA)

	_, _, err = client.Git.UpdateRef(ctx, forkedRepoAccount, EKSAnyrepoName, ref, false)
	if err != nil {
		return fmt.Errorf("error updating ref %s", err)
	}
	return nil
}

func createPullRequestStageBundleTwo() error{

	latestRelease := os.Getenv("LATEST_RELEASE")

	// create client
	accessToken := os.Getenv("SECRET_PAT")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)


	base := latestRelease // target branch for upstream repo
	head := fmt.Sprintf("%s:%s", forkedRepoAccount, "eks-a-releaser")
	title := "Update version files to stage bundle release"
	body := "This pull request is responsible for updating the contents of 3 separate files in order to trigger the staging bundle release pipeline"

	newPR := &github.NewPullRequest{
		Title: &title,
		Head:  &head,
		Base:  &base,
		Body:  &body,
	}

	pr, _, err := client.PullRequests.Create(ctx, upStreamRepoOwner, EKSAnyrepoName, newPR)
	if err != nil {
		return fmt.Errorf("error creating PR %s", err)
	}

	log.Printf("Pull request created: %s\n", pr.GetHTMLURL())
	return nil
}
