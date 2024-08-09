package cmd

/*
	what does this command do?

	this command is responsible for accessing and updating the Makefile with the latest release value

	first, the trigger file within the "eks-a-releaser" branch, ibix16 fork is accessed and its release contents are retrieved e.g "release-0.00"

	secondly, returnUpdatedFile() takes in the entire makefile content string and the retrieved release string from the trigger file, returning the updated makefile as a string

	lastly, the updated makefile is committed to the "eks-a-releaser" branch, ibix16 fork and a pull request is raised targetting the "upstream", currently testerIbix fork
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
)

var (
	AWSrepoOwner = "aws"
	repoName     = "eks-anywhere"
	makeFilePath = "/Makefile"
)

// upMakeFileCmd represents the upMakeFile command
var updateMakefileCmd = &cobra.Command{
	Use:   "update-makefile",
	Short: "Updates BRANCH_NAME?= variable to match new release branch within the Makefile",
	Long:  `A longer description.`,

	Run: func(cmd *cobra.Command, args []string) {
		content := updateMakefile()
		fmt.Print(content)
	},
}

// commits changes into "releaser" branch + raises PR to be merged into latest release branch
func updateMakefile() error {

	//create client
	secretName := "Secret"
	accessToken, err := getSecretValue(secretName)
	if err != nil {
		fmt.Print("error getting secret", err)
	}
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	// string variable holding latest release
	newestRelease := getLatestRelease()

	opts := &github.RepositoryContentGetOptions{
		Ref: "eks-a-releaser", // trigger file is accessed within this branch
	}

	// access makefile and retrieve entire file contents
	triggerFileContentBundleNumber, _, _, err := client.Repositories.GetContents(ctx, botForkAccount, repoName, makeFilePath, opts)
	if err != nil {
		fmt.Print("first breakpoint", err)
	}
	// holds makefile 
	content, err := triggerFileContentBundleNumber.GetContent()
	if err != nil {
		fmt.Print("second breakpoint", err)
	}

	// stores entire updated Makefile as a string
	updatedContent := returnUpdatedMakeFile(content, newestRelease)

	// get latest commit sha
	ref, _, err := client.Git.GetRef(ctx, botForkAccount, repoName, "heads/eks-a-releaser")
	if err != nil {
		return fmt.Errorf("error getting ref %s", err)
	}
	latestCommitSha := ref.Object.GetSHA()

	entries := []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{Path: github.String(strings.TrimPrefix(makeFilePath, "/")), Type: github.String("blob"), Content: github.String(string(updatedContent)), Mode: github.String("100644")})
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
		Message: github.String("Update Makefile"),
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
	targetOwner := "testerIbix" // repo owner 
	base := newestRelease // branch PR will be merged into
	head := fmt.Sprintf("%s:%s", botForkAccount, "eks-a-releaser")
	title := "Updates Makefile to point to new release"
	body := "This pull request is responsible for updating the contents of the Makefile"

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

// returns release value from trigger file, "releaser" branch
func getLatestRelease() string {

	//create client
	secretName := "Secret"
	accessToken, err := getSecretValue(secretName)
	if err != nil {
		fmt.Print("error getting secret", err)
	}
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)


	opts := &github.RepositoryContentGetOptions{
		Ref: "eks-a-releaser", // Replace with the desired branch name
	}

	// access trigger file and retrieve contents
	triggerFileContentBundleNumber, _, _, err := client.Repositories.GetContents(ctx, botForkAccount, repoName, triggerFilePath, opts)
	if err != nil {
		fmt.Print("first breakpoint", err)
	}
	content, err := triggerFileContentBundleNumber.GetContent()
	if err != nil {
		fmt.Print("second breakpoint", err)
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

// updates Makefile with new release, returns entire file updated
func returnUpdatedMakeFile(fileContent, newRelease string) string {
	snippetStartIdentifierB := "BRANCH_NAME?="
	lines := strings.Split(fileContent, "\n")
	var updatedLines []string

	for _, line := range lines {
		if strings.Contains(line, snippetStartIdentifierB) {
			parts := strings.Split(line, "=")
			varNamePart := parts[0] // holds "BRANCH_NAME?"
			updatedLine := varNamePart + "=" + newRelease
			updatedLines = append(updatedLines, updatedLine)
		} else {
			updatedLines = append(updatedLines, line)
		}
	}

	return strings.Join(updatedLines, "\n")

}



func getSecretValue(secretName string)(string, error){

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
	)
	if err != nil {

		return "", fmt.Errorf("failed to load SDK config, %v", err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret value, %v", err)
	}

	secretString := *result.SecretString

	var secretMap map[string]string
	if err := json.Unmarshal([]byte(secretString), &secretMap); err == nil {
		if value, exists := secretMap["PAT"]; exists{
			return value, nil
		}
		return "", fmt.Errorf("PAT value not found in secret")
	}

	return secretString, nil
}

