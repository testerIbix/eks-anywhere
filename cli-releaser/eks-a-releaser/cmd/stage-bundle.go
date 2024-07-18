/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

/*
	what does this command do?
	this command is responsible for accessing the trigger file from the "eks-a-releaser" branch, and extracting the number and version values
	5 seperate functions have been created in total, 3 out of the 5 are responsible for updating the three distinct files and committing changes to "eks-a-releaser" branch
	1 function is responsible for creating a PR containing the 3 commits, the last function runs all 4 of the other funcs

	The pull request is created using the 3 commits from "eks-a-releaser" branch and is intended to be merged into latest release branch
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
	triggerFilePath         = "release/triggers/eks-a-releaser-trigger"
	PersonalforkedRepoOwner = "ibix16"
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

	errFive := createPullRequestStageBundle()
	if errFive != nil {
		return errFive
	}

	return nil
}

func updateBundleNum() error {

	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	opts := &github.RepositoryContentGetOptions{
		Ref: "eks-a-releaser", // Replace with the desired branch name
	}

	// access trigger file and retrieve contents
	triggerFileContentBundleNumber, _, _, err := client.Repositories.GetContents(ctx, PersonalforkedRepoOwner, repoName, triggerFilePath, opts)
	if err != nil {
		fmt.Print("first breakpoint", err)
	}
	content, err := triggerFileContentBundleNumber.GetContent()
	if err != nil {
		fmt.Print("second breakpoint", err)
	}

	// Find the line containing the identifier
	snippetStartIdentifierB := "number: "
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
		//return fmt.Errorf("snippet not found", nil)  // Snippet not found
		log.Panic("snippet not found...")
	}

	// holds full string
	bundleNumberLine := lines[startIndex]

	// split string to isolate bundle number
	parts := strings.Split(bundleNumberLine, ": ")

	// holds bundle number value as string
	desiredPart := parts[1]

	// get latest commit sha
	ref, _, err := client.Git.GetRef(ctx, PersonalforkedRepoOwner, repoName, "heads/eks-a-releaser")
	if err != nil {
		return fmt.Errorf("error getting ref %s", err)
	}
	latestCommitSha := ref.Object.GetSHA()

	entries := []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{Path: github.String(strings.TrimPrefix(bundleNumPath, "/")), Type: github.String("blob"), Content: github.String(string(desiredPart)), Mode: github.String("100644")})
	tree, _, err := client.Git.CreateTree(ctx, PersonalforkedRepoOwner, repoName, *ref.Object.SHA, entries)
	if err != nil {
		return fmt.Errorf("error creating tree %s", err)
	}

	//validate tree sha
	newTreeSHA := tree.GetSHA()

	// create new commit
	author := &github.CommitAuthor{
		Name:  github.String("ibix16"),
		Email: github.String("ibixrivera16@gmail.com"),
	}

	commit := &github.Commit{
		Message: github.String("Increment bundle number file"),
		Tree:    &github.Tree{SHA: github.String(newTreeSHA)},
		Author:  author,
		Parents: []*github.Commit{{SHA: github.String(latestCommitSha)}},
	}

	commitOP := &github.CreateCommitOptions{}
	newCommit, _, err := client.Git.CreateCommit(ctx, PersonalforkedRepoOwner, repoName, commit, commitOP)
	if err != nil {
		return fmt.Errorf("creating commit %s", err)
	}
	newCommitSHA := newCommit.GetSHA()

	// update branch reference
	ref.Object.SHA = github.String(newCommitSHA)

	_, _, err = client.Git.UpdateRef(ctx, PersonalforkedRepoOwner, repoName, ref, false)
	if err != nil {
		return fmt.Errorf("error updating ref %s", err)
	}

	return nil
}

func updateCLIMax() error {

	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	opts := &github.RepositoryContentGetOptions{
		Ref: "eks-a-releaser", // Replace with the desired branch name
	}

	// access trigger file and retrieve contents
	triggerFileContentBundleNumber, _, _, err := client.Repositories.GetContents(ctx, PersonalforkedRepoOwner, repoName, triggerFilePath, opts)
	if err != nil {
		fmt.Print("first breakpoint", err)
	}
	content, err := triggerFileContentBundleNumber.GetContent()
	if err != nil {
		fmt.Print("second breakpoint", err)
	}

	// Find the line containing the identifier
	snippetStartIdentifierB := "version: "
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
		//return fmt.Errorf("snippet not found", nil)  // Snippet not found
		log.Panic("snippet not found...")
	}

	// holds full string
	bundleNumberLine := lines[startIndex]

	// split string to isolate bundle number
	parts := strings.Split(bundleNumberLine, ": ")

	// holds bundle number value as string
	desiredPart := parts[1]

	// get latest commit sha
	ref, _, err := client.Git.GetRef(ctx, PersonalforkedRepoOwner, repoName, "heads/eks-a-releaser")
	if err != nil {
		return fmt.Errorf("error getting ref %s", err)
	}
	latestCommitSha := ref.Object.GetSHA()

	entries := []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{Path: github.String(strings.TrimPrefix(cliMaxVersionPath, "/")), Type: github.String("blob"), Content: github.String(string(desiredPart)), Mode: github.String("100644")})
	tree, _, err := client.Git.CreateTree(ctx, PersonalforkedRepoOwner, repoName, *ref.Object.SHA, entries)
	if err != nil {
		return fmt.Errorf("error creating tree %s", err)
	}

	//validate tree sha
	newTreeSHA := tree.GetSHA()

	// create new commit
	author := &github.CommitAuthor{
		Name:  github.String("ibix16"),
		Email: github.String("ibixrivera16@gmail.com"),
	}

	commit := &github.Commit{
		Message: github.String("Update CLI Max version"),
		Tree:    &github.Tree{SHA: github.String(newTreeSHA)},
		Author:  author,
		Parents: []*github.Commit{{SHA: github.String(latestCommitSha)}},
	}

	commitOP := &github.CreateCommitOptions{}
	newCommit, _, err := client.Git.CreateCommit(ctx, PersonalforkedRepoOwner, repoName, commit, commitOP)
	if err != nil {
		return fmt.Errorf("creating commit %s", err)
	}
	newCommitSHA := newCommit.GetSHA()

	// update branch reference
	ref.Object.SHA = github.String(newCommitSHA)

	_, _, err = client.Git.UpdateRef(ctx, PersonalforkedRepoOwner, repoName, ref, false)
	if err != nil {
		return fmt.Errorf("error updating ref %s", err)
	}

	return nil
}

func updateCLIMin() error {

	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	opts := &github.RepositoryContentGetOptions{
		Ref: "eks-a-releaser", // Replace with the desired branch name
	}

	// access trigger file and retrieve contents
	triggerFileContentBundleNumber, _, _, err := client.Repositories.GetContents(ctx, PersonalforkedRepoOwner, repoName, triggerFilePath, opts)
	if err != nil {
		fmt.Print("first breakpoint", err)
	}
	content, err := triggerFileContentBundleNumber.GetContent()
	if err != nil {
		fmt.Print("second breakpoint", err)
	}

	// Find the line containing the identifier
	snippetStartIdentifierB := "version: "
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
		//return fmt.Errorf("snippet not found", nil)  // Snippet not found
		log.Panic("snippet not found...")
	}

	// holds full string
	bundleNumberLine := lines[startIndex]

	// split string to isolate bundle number
	parts := strings.Split(bundleNumberLine, ": ")

	// holds bundle number value as string
	desiredPart := parts[1]

	// get latest commit sha
	ref, _, err := client.Git.GetRef(ctx, PersonalforkedRepoOwner, repoName, "heads/eks-a-releaser")
	if err != nil {
		return fmt.Errorf("error getting ref %s", err)
	}
	latestCommitSha := ref.Object.GetSHA()

	entries := []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{Path: github.String(strings.TrimPrefix(cliMinVersionPath, "/")), Type: github.String("blob"), Content: github.String(string(desiredPart)), Mode: github.String("100644")})
	tree, _, err := client.Git.CreateTree(ctx, PersonalforkedRepoOwner, repoName, *ref.Object.SHA, entries)
	if err != nil {
		return fmt.Errorf("error creating tree %s", err)
	}

	//validate tree sha
	newTreeSHA := tree.GetSHA()

	// create new commit
	author := &github.CommitAuthor{
		Name:  github.String("ibix16"),
		Email: github.String("ibixrivera16@gmail.com"),
	}

	commit := &github.Commit{
		Message: github.String("Update CLI Min version"),
		Tree:    &github.Tree{SHA: github.String(newTreeSHA)},
		Author:  author,
		Parents: []*github.Commit{{SHA: github.String(latestCommitSha)}},
	}

	commitOP := &github.CreateCommitOptions{}
	newCommit, _, err := client.Git.CreateCommit(ctx, PersonalforkedRepoOwner, repoName, commit, commitOP)
	if err != nil {
		return fmt.Errorf("creating commit %s", err)
	}
	newCommitSHA := newCommit.GetSHA()

	// update branch reference
	ref.Object.SHA = github.String(newCommitSHA)

	_, _, err = client.Git.UpdateRef(ctx, PersonalforkedRepoOwner, repoName, ref, false)
	if err != nil {
		return fmt.Errorf("error updating ref %s", err)
	}

	return nil
}

func createPullRequestStageBundle() error {

	latestReleaseValue := getLatestRelease()

	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	
	base := latestReleaseValue
	head := fmt.Sprintf("%s:%s", PersonalforkedRepoOwner, "eks-a-releaser")
	title := "Update version files to stage bundle release"
	body := "This pull request is responsible for updating the contents of 3 separate files in order to trigger the staging bundle release pipeline"

	newPR := &github.NewPullRequest{
		Title: &title,
		Head:  &head,
		Base:  &base,
		Body:  &body,
	}

	pr, _, err := client.PullRequests.Create(ctx, PersonalforkedRepoOwner, repoName, newPR)
	if err != nil {
		return fmt.Errorf("error creating PR %s", err)
	}

	log.Printf("Pull request created: %s\n", pr.GetHTMLURL())
	return nil
}
