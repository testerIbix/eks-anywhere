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
	prowFilePathOne  = "/jobs/aws/eks-anywhere-build-tooling/eks-anywhere-attribution-periodics-release-0.19.yaml"
	prowFilePathThree = "/templater/jobs/periodic/eks-anywhere-build-tooling/eks-anywhere-attribution-periodics-release-0.19.yaml"

	// update templater file, run make command to generate jobs files

	// currently, both functions go in and directly update files, no call to make command
)

// upProwCmd represents the upProw command
var updateProwCmd = &cobra.Command{
	Use:   "update-prow",
	Short: "accesses prow-jobs repo and updates version files",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {
		runAllProw()
	},
}

func runAllProw(){
	contentOne := updateProwFileOne()
	fmt.Print(contentOne)
	fmt.Println("--------------------------------------------")
	contentFileThree := updateProwFileThree()
	fmt.Print(contentFileThree)
}


func updateProwFileOne()string{

	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	// access file one and retrieve entire file contents
	prowFileOneContent, _, _, err := client.Repositories.GetContents(ctx, PersonalforkedRepoOwner, prowRepoName, prowFilePathOne, nil)
	if err != nil {
		fmt.Print("first breakpoint", err)
	}
	content, err := prowFileOneContent.GetContent()
	if err != nil {
		fmt.Print("second breakpoint", err)
	}
	// var "content" holds entire string of unedited file



	// isolate base_ref: line
	baseRefSnippetStart := "base_ref: "
	Firstlines := strings.Split(content, "\n")
	startIndex := -1
	endIndex := -1

	for i, line := range Firstlines {
		if strings.Contains(line, baseRefSnippetStart) {
			startIndex = i
			endIndex = i // Set endIndex to the same line as startIndex
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		//return fmt.Errorf("snippet not found", nil)  // Snippet not found
		log.Panic("snippet not found...")
	}

	// holds string base_ref: release-0.00
	baseRefLine := Firstlines[startIndex]

	
	splitBaseRefLine := strings.Split(baseRefLine, ": ")
	// var holds release-0.00 portion
	previousRelease := splitBaseRefLine[1]
	// var holds latest release retrieved from trigger file
	latestRelease := getLatestRelease()

	// replace all instances of previousRelease with latestRelease ~ successfully updates base_ref & value fields
	updatedBaseRefFileContent := strings.ReplaceAll(content, previousRelease, latestRelease)




	// update - name: field 
	// isolate line 
	nameSnippetStartIdentifier := "- name: "
	Firstlines = strings.Split(updatedBaseRefFileContent, "\n")
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

	//holds string - name: eks-anywhere-attribution-periodic-release-0-19
	nameLine := Firstlines[startIndex]

	nameLineParts := strings.Split(nameLine, "release-")

	//holds string 0-19
	nameLineReleasePortion := nameLineParts[1]


	// latestRelease var holds release-0.21 
	// we want to isolate the numerical portion 
	// and convert it from 0.21 ---> 0-21
	splitLatestRelease := strings.Split(latestRelease, ".")
	targetLatestReleaseValue := splitLatestRelease[1]

	// var holds 0-21
	convertedTargetLatestReleaseValue := "0-" + targetLatestReleaseValue

	fullyUpdatedFileOne := strings.ReplaceAll(updatedBaseRefFileContent, nameLineReleasePortion, convertedTargetLatestReleaseValue)

	return fullyUpdatedFileOne

	/*
	this function updates 3 seperate fields within the 1st prow-jobs file
	- name:
	- base_ref:
	- value:
	it updates name: and base_ref: on line 93 by replacing all instances of the previous release, with the latest one
	additionally, it updates the name: field by isolating the numerical portion of the latest release and converting it to match the format of the file
	it then uses the same method ReplaceAll() to update outdated instances 
	*/

	// Missing : add logic to create commit 
}



// func updateProwFileTwo()string{

// 	// create client
// 	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
// 	ctx := context.Background()
// 	client := github.NewClient(nil).WithAuthToken(accessToken)

	
// 	// access file one and retrieve entire file contents
// 	prowFileOneContent, _, _, err := client.Repositories.GetContents(ctx, PersonalforkedRepoOwner, prowRepoName, prowFilePathTwo, nil)
// 	if err != nil {
// 		fmt.Print("first breakpoint", err)
// 	}
// 	content, err := prowFileOneContent.GetContent()
// 	if err != nil {
// 		fmt.Print("second breakpoint", err)
// 	}

// 	return content
// }



func updateProwFileThree()string{

	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	// access file one and retrieve entire file contents
	prowFileOneContent, _, _, err := client.Repositories.GetContents(ctx, PersonalforkedRepoOwner, prowRepoName, prowFilePathThree, nil)
	if err != nil {
		fmt.Print("first breakpoint", err)
	}
	content, err := prowFileOneContent.GetContent()
	if err != nil {
		fmt.Print("second breakpoint", err)
	}

	baseRefSnippetStart := "value: "
	Firstlines := strings.Split(content, "\n")
	startIndex := -1
	endIndex := -1

	for i, line := range Firstlines {
		if strings.Contains(line, baseRefSnippetStart) {
			startIndex = i
			endIndex = i // Set endIndex to the same line as startIndex
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		//return fmt.Errorf("snippet not found", nil)  // Snippet not found
		log.Panic("snippet not found.....")
	}

	// holds string value: release-0.00
	valueLine := Firstlines[startIndex]
	splitValueLine := strings.Split(valueLine, ": ")
	previousRelease := splitValueLine[1]
	latestRelease := getLatestRelease()

	updatedVerOneFile := strings.ReplaceAll(content, previousRelease, latestRelease)

	
	

	// update - name: field 
	// isolate line 
	jobNameSnippetStartIdentifier := "jobName: "
	Firstlines = strings.Split(updatedVerOneFile, "\n")
	startIndex = -1
	endIndex = -1

	for i, line := range Firstlines {
		if strings.Contains(line, jobNameSnippetStartIdentifier) {
			startIndex = i
			endIndex = i // Set endIndex to the same line as startIndex
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		//return fmt.Errorf("snippet not found", nil)  // Snippet not found
		log.Panic("ERROR : snippet not found...")
	}

	//holds string - jobName: eks-anywhere-attribution-periodic-release-0-19
	nameLine := Firstlines[startIndex]

	nameLineParts := strings.Split(nameLine, "release-")

	//holds string 0-19
	nameLineReleasePortion := nameLineParts[1]

	// latestRelease var holds release-0.21 
	// we want to isolate the numerical portion 
	// and convert it from 0.21 ---> 0-21
	splitLatestRelease := strings.Split(latestRelease, ".")
	targetLatestReleaseValue := splitLatestRelease[1]

	// var holds 0-21
	convertedTargetLatestReleaseValue := "0-" + targetLatestReleaseValue

	fullyUpdatedFileOne := strings.ReplaceAll(updatedVerOneFile, nameLineReleasePortion, convertedTargetLatestReleaseValue)

	return fullyUpdatedFileOne

	/*
	this function updates 3 seperate fields within the 3rd prow-jobs file
	- jobName:
	- baseRef:
	- value:
	it updates name: and baseRef: by replacing all instances of the previous release, with the latest one
	additionally, it updates the jobName: field by isolating the numerical portion of the latest release and converting it to match the format of the file
	it then uses the same method ReplaceAll() to update outdated instances 
	*/

	// Missing : add logic to create commit 
}

