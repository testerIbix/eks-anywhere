/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

/*
	what does this command do?
	this command is responsible for updating the homebrew release version file

	retrievesLatestVersion() - accesses trigger file in "eks-a-releaser" branch, bot fork / ibix16 fork
	returns version: v0.0.0 field value

	updateHomebrew() - retrieves the latest version value using the function above
	accesses homebrew cli version file in "eks-a-releaser" branch, bot fork / ibix16 fork
	updates file contents with retrieved latest version value, commits changes to "eks-a-releaser" branch, bot fork / ibix16 fork

	PR is then raised from "eks-a-releaser" branch, bot fork / ibix16 fork, targgetting new release branch on upstream repo
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
	homebrewPath = "release/triggers/brew-version-release/CLI_RELEASE_VERSION"
)

// updateHomebrewCmd represents the updateHomebrew command
var updateHomebrewCmd = &cobra.Command{
	Use:   "update-homebrew",
	Short: "Updates homebrew with latest version in eks-a-releaser branch, PR targets release branch",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {
		runAllHomebrew()
	},
}

func runAllHomebrew(){
	errOne := updateHomebrew()
	if errOne != nil {
		log.Panic(errOne)
	}

	errTwo := createPullRequestHomebrew()
	if errTwo != nil {
		log.Panic(errTwo)
	}
}



func updateHomebrew()error{
	
	// value we will use to update 
	latestVersionValue := retrieveLatestVersion()


	// create client 
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
	triggerFileContentBundleNumber, _, _, err := client.Repositories.GetContents(ctx, botForkAccount, repoName, homebrewPath, opts)
	if err != nil {
		fmt.Print("first breakpoint", err)
	}

	// holds content of homebrew cli version file
	content, err := triggerFileContentBundleNumber.GetContent()
	if err != nil {
		fmt.Print("second breakpoint", err)
	}

	// update instances of previous release with new
	updatedFile := strings.ReplaceAll(content, content, latestVersionValue)


	// get latest commit sha
	ref, _, err := client.Git.GetRef(ctx, botForkAccount, repoName, "heads/eks-a-releaser")
	if err != nil {
		return fmt.Errorf("error getting ref %s", err)
	}
	latestCommitSha := ref.Object.GetSHA()

	entries := []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{Path: github.String(strings.TrimPrefix(homebrewPath, "/")), Type: github.String("blob"), Content: github.String(string(updatedFile)), Mode: github.String("100644")})
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
		Message: github.String("Update brew-version value to point to new release"),
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



// retrieves latest version from trigger file, eks-a-releaser branch
func retrieveLatestVersion()string{

	// create client 
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
			endIndex = i // Set endIndex to the same line as startIndex
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		log.Panic("snippet not found...")
		return ""
	}

	// holds full string
	bundleNumberLine := lines[startIndex]

	// split string to isolate bundle number
	parts := strings.Split(bundleNumberLine, ": ")

	// holds bundle number value as string
	desiredPart := parts[1]

	return desiredPart
}


func createPullRequestHomebrew()error{

	latestReleaseValue := getLatestRelease()

	secretName := "Secret"
	accessToken, err := getSecretValue(secretName)
	if err != nil {
		fmt.Print("error getting secret", err)
	}
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	// targgetting latest release branch
	targetOwner := "testerIbix"
	base := latestReleaseValue // branch PR will target
	head := fmt.Sprintf("%s:%s", botForkAccount, "eks-a-releaser")
	title := "Update homebrew cli version value to point to new release"
	body := "This pull request is responsible for updating the contents of the home brew cli version file"

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