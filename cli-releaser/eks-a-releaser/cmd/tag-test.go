/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

// works!

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// tagTestCmd represents the tagTest command
var tagTestCmd = &cobra.Command{
	Use:   "tag-test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runBothTagT()
	},
}


func runBothTagT(){

	//retrieve commit hash 
	commitHash := retrieveLatestProdCLIHash()
	
	//create tag with commit hash
	tag, errOne := createTagT(commitHash)
	if errOne != nil {
		log.Panic(errOne)
	}

	rel, errTwo := createGitHubReleaseT(tag)
	if errTwo != nil {
		log.Panic(errTwo)
	}

	//print release object
	fmt.Print(rel)
}


// creates tag using retrieved commit hash
func createTagT(commitHash string) (*github.GitObject, error){
	

	// retrieve tag name "v0.0.00" from trigger file, "eks-a-releaser" branch 
	version := retrieveLatestVersion()

	
	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()

	// Create a new GitHub client instance with the token type set to "Bearer"
	ts := oauth2.StaticTokenSource(&oauth2.Token{
    	AccessToken: accessToken,
    	TokenType:   "Bearer",
	})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)


	
	commit, _, err := client.Git.GetCommit(ctx, PersonalforkedRepoOwner, repoName, commitHash)
    if err != nil {
        return nil, fmt.Errorf("error getting commit: %v", err)
    }

    tagger := &github.CommitAuthor{
        Name:  github.String("ibix16"),
        Email: github.String("ibixrivera16@gmail.com"),
        Date:  &github.Timestamp{Time: time.Now()},
    }

    commitSHA := github.String(*commit.SHA)
    gitObject := &github.GitObject{
        Type: github.String("commit"),
        SHA:  commitSHA,
        URL:  nil,
    }

    tagMessage := "EKS-Anywhere " + version + " release"
    tag, _, err := client.Git.CreateTag(ctx, PersonalforkedRepoOwner, repoName, &github.Tag{
        Tag:     github.String(version),
        Message: github.String(tagMessage),
        Object:  gitObject,
        Tagger:  tagger,
    })
    if err != nil {
        return nil, fmt.Errorf("error creating tag: %v", err)
    }

    fmt.Printf("Created tag %s for commit %s\n", *tag.Tag, commitHash)
    return gitObject, nil
}


func createGitHubReleaseT(gitObject *github.GitObject) (*github.RepositoryRelease, error){
	// Implement logic to create a release on GitHub using the tag created in createTag()

	// Retrieve version
    version := retrieveLatestVersion()

    accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	// Create a new GitHub client instance with the token type set to "Bearer"
	ts := oauth2.StaticTokenSource(&oauth2.Token{
    	AccessToken: accessToken,
    	TokenType:   "Bearer",
	})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)


    // Check if a release with the given tag name already exists
    release, _, err := client.Repositories.GetReleaseByTag(ctx, PersonalforkedRepoOwner, repoName, version)
    if err == nil {
        // Release with the given tag name already exists
        fmt.Printf("Release %s already exists!\n", version)
        return release, nil
    }

    // Create a new release
    releaseDesc := "EKS-Anywhere " + version + " release"
    release = &github.RepositoryRelease{
        TagName: github.String(version),
        Name:    github.String(version),
        Body:    github.String(releaseDesc),
        TargetCommitish: gitObject.SHA,
    }

    rel, _, err := client.Repositories.CreateRelease(ctx, PersonalforkedRepoOwner, repoName, release)
    if err != nil {
        return nil, err
    }

    fmt.Printf("Release %s created successfully!\n", *rel.TagName)
    return rel, nil
	
}




