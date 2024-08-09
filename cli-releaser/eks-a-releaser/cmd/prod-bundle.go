/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

/*
	what does this command do?

	this command is responsible for accessing the trigger file from the "eks-a-releaser" branch, ibix16 personal fork repo

	4 distinct functions have been created, 3 out of the 4 update a file and commit the changes to the "eks-a-releaser" branch, ibix16 fork repo

	Additionally, the first update function handles the logic of creating the PR

	the idea is that changes are first committed into my fork / bot fork of the upstream repo
	then a PR is created with those commits targetting the latest release branch in the upstream repo
	therefore with this approach, our cli is never directly touching and updating the upstream repo
*/

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
)

var (
	prodBundleNumPath     = "release/triggers/bundle-release/production/BUNDLE_NUMBER"
	prodCliMaxVersionPath = "release/triggers/bundle-release/production/CLI_MAX_VERSION"
	prodCliMinVersionPath = "release/triggers/bundle-release/production/CLI_MIN_VERSION"
)

// prodBundleCmd represents the prodBundle command
var prodBundleCmd = &cobra.Command{
	Use:   "prod-bundle",
	Short: "creates a PR containing 3 commits, each updating the contents of a singular file intended for prod bundle release",
	Long: `Retrieves updated content for production : bundle number, cli max version, and cli min version. 
	Writes the updated changes to the 3 files and raises a PR with the 3 commits.`,

	Run: func(cmd *cobra.Command, args []string) {
		err := updateAllProdBundleFiles()
		if err != nil {
			log.Panic(err)
		}
	},
}

func updateAllProdBundleFiles() error {

	errOne := updateProdBundleNumber()
	if errOne != nil {
		return errOne
	}

	errTwo := updateProdMaxVersion()
	if errTwo != nil {
		return errTwo
	}

	errThree := updateProdMinVersion()
	if errThree != nil {
		return errThree
	}

	return nil
}

func updateProdBundleNumber() error {

	//create client
	secretName := "Secret"
	accessToken, err := getSecretValue(secretName)
	if err != nil {
		fmt.Print("error getting secret", err)
	}
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)


	opts := &github.RepositoryContentGetOptions{
		Ref: "eks-a-releaser", 
	}
	// access trigger file
	triggerFileContentBundleNumber, _, _, err := client.Repositories.GetContents(ctx, botForkAccount, repoName, triggerFilePath, opts)
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
			endIndex = i 
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		log.Panic("snippet not found...")
	}

	// holds full string
	bundleNumberLine := lines[startIndex]

	// split string to isolate bundle number
	parts := strings.Split(bundleNumberLine, ": ")

	// holds bundle number value as string
	desiredPart := parts[1]

	// get latest commit sha
	ref, _, err := client.Git.GetRef(ctx, botForkAccount, repoName, "heads/eks-a-releaser")
	if err != nil {
		return fmt.Errorf("error getting ref %s", err)
	}
	latestCommitSha := ref.Object.GetSHA()

	entries := []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{Path: github.String(strings.TrimPrefix(prodBundleNumPath, "/")), Type: github.String("blob"), Content: github.String(string(desiredPart)), Mode: github.String("100644")})
	tree, _, err := client.Git.CreateTree(ctx, botForkAccount, repoName, *ref.Object.SHA, entries)
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
		Message: github.String("Update production bundle number file"),
		Tree:    &github.Tree{SHA: github.String(newTreeSHA)},
		Author:  author,
		Parents: []*github.Commit{{SHA: github.String(latestCommitSha)}},
	}

	commitOP := &github.CreateCommitOptions{}
	newCommit, _, err := client.Git.CreateCommit(ctx, botForkAccount, repoName, commit, commitOP)
	if err != nil {
		return fmt.Errorf("creating commit %s", err)
	}
	newCommitSHA := newCommit.GetSHA()

	// update branch reference
	ref.Object.SHA = github.String(newCommitSHA)

	_, _, err = client.Git.UpdateRef(ctx, botForkAccount, repoName, ref, false)
	if err != nil {
		return fmt.Errorf("error updating ref %s", err)
	}

	
	// create pull request
	targetOwner := "testerIbix"
	latestRelease := getLatestRelease()
	base := latestRelease
	head := fmt.Sprintf("%s:%s", botForkAccount, "eks-a-releaser")
	title := "Update version files to stage production bundle release"
	body := "This pull request is responsible for updating the contents of 3 seperate files in order to trigger the production bundle release pipeline"

	newPR := &github.NewPullRequest{
		Title: &title,
		Head:  &head,
		Base:  &base,
		Body:  &body,
	}

	pr, _, err := client.PullRequests.Create(ctx, targetOwner, repoName, newPR)
	if err != nil {
		return fmt.Errorf("error creating PR %s", err)
	}

	log.Printf("Pull request created: %s\n", pr.GetHTMLURL())
	return nil

}

func updateProdMaxVersion() error {

	//create client
	secretName := "Secret"
	accessToken, err := getSecretValue(secretName)
	if err != nil {
		fmt.Print("error getting secret", err)
	}
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)


	opts := &github.RepositoryContentGetOptions{
		Ref: "eks-a-releaser",
	}

	// access trigger file 
	triggerFileContentBundleNumber, _, _, err := client.Repositories.GetContents(ctx, botForkAccount, repoName, triggerFilePath, opts)
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
			endIndex = i 
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		log.Panic("snippet not found...")
	}

	// holds full string
	bundleNumberLine := lines[startIndex]

	// split string to isolate bundle number
	parts := strings.Split(bundleNumberLine, ": ")

	// holds bundle number value as string
	desiredPart := parts[1]

	// get latest commit sha
	ref, _, err := client.Git.GetRef(ctx, botForkAccount, repoName, "heads/eks-a-releaser")
	if err != nil {
		return fmt.Errorf("error getting ref %s", err)
	}
	latestCommitSha := ref.Object.GetSHA()

	entries := []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{Path: github.String(strings.TrimPrefix(prodCliMaxVersionPath, "/")), Type: github.String("blob"), Content: github.String(string(desiredPart)), Mode: github.String("100644")})
	tree, _, err := client.Git.CreateTree(ctx, botForkAccount, repoName, *ref.Object.SHA, entries)
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
		Message: github.String("Update production max version file"),
		Tree:    &github.Tree{SHA: github.String(newTreeSHA)},
		Author:  author,
		Parents: []*github.Commit{{SHA: github.String(latestCommitSha)}},
	}

	commitOP := &github.CreateCommitOptions{}
	newCommit, _, err := client.Git.CreateCommit(ctx, botForkAccount, repoName, commit, commitOP)
	if err != nil {
		return fmt.Errorf("creating commit %s", err)
	}
	newCommitSHA := newCommit.GetSHA()

	// update branch reference
	ref.Object.SHA = github.String(newCommitSHA)

	_, _, err = client.Git.UpdateRef(ctx, botForkAccount, repoName, ref, false)
	if err != nil {
		return fmt.Errorf("error updating ref %s", err)
	}

	return nil

}

func updateProdMinVersion() error {

	//create client
	secretName := "Secret"
	accessToken, err := getSecretValue(secretName)
	if err != nil {
		fmt.Print("error getting secret", err)
	}
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	opts := &github.RepositoryContentGetOptions{
		Ref: "eks-a-releaser", 
	}
	// access trigger file
	triggerFileContentBundleNumber, _, _, err := client.Repositories.GetContents(ctx, botForkAccount, repoName, triggerFilePath, opts)
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
			endIndex = i 
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		log.Panic("snippet not found...")
	}

	// holds full string
	bundleNumberLine := lines[startIndex]

	// split string to isolate bundle number
	parts := strings.Split(bundleNumberLine, ": ")

	// holds bundle number value as string
	desiredPart := parts[1]

	// get latest commit sha
	ref, _, err := client.Git.GetRef(ctx, botForkAccount, repoName, "heads/eks-a-releaser")
	if err != nil {
		return fmt.Errorf("error getting ref %s", err)
	}
	latestCommitSha := ref.Object.GetSHA()

	entries := []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{Path: github.String(strings.TrimPrefix(prodCliMinVersionPath, "/")), Type: github.String("blob"), Content: github.String(string(desiredPart)), Mode: github.String("100644")})
	tree, _, err := client.Git.CreateTree(ctx, botForkAccount, repoName, *ref.Object.SHA, entries)
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
		Message: github.String("Update production min version file"),
		Tree:    &github.Tree{SHA: github.String(newTreeSHA)},
		Author:  author,
		Parents: []*github.Commit{{SHA: github.String(latestCommitSha)}},
	}

	commitOP := &github.CreateCommitOptions{}
	newCommit, _, err := client.Git.CreateCommit(ctx, botForkAccount, repoName, commit, commitOP)
	if err != nil {
		return fmt.Errorf("creating commit %s", err)
	}
	newCommitSHA := newCommit.GetSHA()

	// update branch reference
	ref.Object.SHA = github.String(newCommitSHA)

	_, _, err = client.Git.UpdateRef(ctx, botForkAccount, repoName, ref, false)
	if err != nil {
		return fmt.Errorf("error updating ref %s", err)
	}

	return nil

}
