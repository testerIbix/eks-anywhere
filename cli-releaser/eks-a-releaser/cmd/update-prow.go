/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

/*
	what does this command do?
	this command is responsible for updating 4 seperate files within the prow-jobs repo
	Creates a pull request targetting "main", containing the 4 commits

	7/18/24 : function incomplete, further discussion required
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
	prowRepoName = "eks-anywhere-prow-jobs"
	templaterFilePath = "/templater/jobs/periodic/eks-anywhere-build-tooling/eks-anywhere-attribution-periodics-release-0.19.yaml"
)

// upProwCmd represents the upProw command
var updateProwCmd = &cobra.Command{
	Use:   "update-prow",
	Short: "accesses prow-jobs repo and updates version files",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {
		content := updateTemplaterFile()
		fmt.Print(content)
	},
}



func updateTemplaterFile()string{

	// var holds latest release retrieved from trigger file
	latestRelease := getLatestRelease()

	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)


	// access file one and retrieve entire file contents
	templaterFileContent, _, _, err := client.Repositories.GetContents(ctx, PersonalforkedRepoOwner, prowRepoName, templaterFilePath, nil)
	if err != nil {
		fmt.Print("first breakpoint", err)
	}

	// var "content" holds entire string of templater file
	content, err := templaterFileContent.GetContent()
	if err != nil {
		fmt.Print("second breakpoint", err)
	}
	


	// update jobName field , isolate line 
	nameSnippetStartIdentifier := "jobName: "
	Firstlines := strings.Split(content, "\n")
	startIndex := -1
	endIndex := -1

	for i, line := range Firstlines {
		if strings.Contains(line, nameSnippetStartIdentifier) {
			startIndex = i
			endIndex = i // Set endIndex to the same line as startIndex
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		//return fmt.Errorf("snippet not found", nil)  // Snippet not found
		log.Panic("snippet not found...")
	}

	//holds string - name: eks-anywhere-attribution-periodic-release-0-19
	nameLine := Firstlines[startIndex]

	jobNameLineParts := strings.Split(nameLine, "release-")

	//holds string 0-19
	jobNameLineReleasePortion := jobNameLineParts[1]


	// latestRelease var holds release-0.00 from trigger file
	// we want to isolate the numerical portion 
	// and convert it from 0.00 ---> 0-00
	splitLatestRelease := strings.Split(latestRelease, ".")
	targetLatestReleaseValue := splitLatestRelease[1]

	// var holds 0-21
	convertedTargetLatestReleaseValue := "0-" + targetLatestReleaseValue

	firstUpdatedFile := strings.ReplaceAll(content, jobNameLineReleasePortion, convertedTargetLatestReleaseValue)

	

	// update jobName field , isolate line 
	nameSnippetStartIdentifier = "value: "
	Firstlines = strings.Split(content, "\n")
	startIndex = -1
	endIndex = -1

	for i, line := range Firstlines {
		if strings.Contains(line, nameSnippetStartIdentifier) {
			startIndex = i
			endIndex = i // Set endIndex to the same line as startIndex
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		//return fmt.Errorf("snippet not found", nil)  // Snippet not found
		log.Panic("snippet not found...")
	}

	//holds value: release-0.00 from templater file
	nameLine = Firstlines[startIndex]


	// isolates release-0.00 portion from templater file
	valueLine := strings.Split(nameLine, ": ")
	valueLinePortion := valueLine[1]


	// replaces all instances of "release-0.00" with var valueLinePortion, updating both value: line and baseRef: line
	secondUpdatedFile := strings.ReplaceAll(firstUpdatedFile, valueLinePortion, latestRelease)

	return secondUpdatedFile

	// all required fields successfully get updated 

	// missing : create commit and PR with changes 
	// as well as speak with abhay surrounding changing file names
	// github api does not provide a direct way to update file names
	// proposed solutions online state to delete file and create new one with new name but prob not ideal
	// other solutions suggest using UpdateFile method from github api 
}
