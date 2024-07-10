/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
)

var (
	repoNameB   = "eks-anywhere-prow-jobs"
	filePathOne = "/jobs/aws/eks-anywhere-build-tooling/eks-anywhere-attribution-periodics-release-0.19.yaml"
	//filePathTwo  = "/jobs/aws/eks-anywhere-build-tooling/eks-anywhere-checksum-periodics-release-0.14.yaml"
	filePathThree = "/templater/jobs/periodic/eks-anywhere-build-tooling/eks-anywhere-attribution-periodics-release-0.19.yaml"
	//filePathFour = "/templater/jobs/periodic/eks-anywhere-build-tooling/eks-anywhere-checksum-periodics-release-0.14.yaml"
)

// upProwCmd represents the upProw command
var updateProwCmd = &cobra.Command{
	Use:   "update-prow",
	Short: "accesses prow-jobs repo and updates version files",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,


	Run: func(cmd *cobra.Command, args []string) {
		entireFile, err := updateFileContentsOne()
		if err != nil {
			fmt.Print(err)
		}
		fmt.Print(entireFile)
	},
}


func updateFileContentsOne() (string, error) {

	//create client
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoNameB, filePathOne, nil)
	if err != nil {
		fmt.Print(err)
	}

	content, err := fileContent.GetContent() // holds entire file content
	if err != nil {
		fmt.Print(err)
	}
	

	// Find the line containing the identifier
	snippetStartIdentifierB := "- name:"
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
		return "error code snippet not found", nil // Snippet not found
	}

	// holds string for name: 
	snippet := lines[startIndex]

	splitSnippet := strings.Split(snippet, "-")

	// holds last release value (currently 19) --> left to do, increment value remerge string (convert from string to int then back to string)
	desiredSplitSnippet := splitSnippet[7]
	stringToInt, err := strconv.Atoi(desiredSplitSnippet)
	if err != nil {
		fmt.Print(err)
	}
	// holds string incremented "20"
	stringToInt++
	intToString := strconv.Itoa(stringToInt)
	aVersionPart := splitSnippet[len(splitSnippet)-2:]
	oldVersion := strings.Join(aVersionPart, "-")
	newVersion := strings.Replace(oldVersion, aVersionPart[1], intToString, 1)
	newVersionParts := strings.Split(newVersion, "-")
	splitSnippet[len(splitSnippet)-2] = newVersionParts[0]
	splitSnippet[len(splitSnippet)-1] = newVersionParts[1]
	//finalString := strings.Join(splitSnippet, "-")
	//left to do, commit/write new updated finalString back to file





	// Find the line containing the identifier
	snippetStartIdentifierB = "base_ref: release-0."
	lines = strings.Split(content, "\n")
	startIndex = -1
	endIndex = -1

	for i, line := range lines {
		if strings.Contains(line, snippetStartIdentifierB) {
			startIndex = i
			endIndex = i // Set endIndex to the same line as startIndex
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		return "error code snippet not found", nil // Snippet not found
	}

	// holds string for base_ref: release-0.19
	snippetB := lines[startIndex]
	// holds 0.19
	parts := strings.Split(snippetB, "-")
	desiredValueSet := strings.Split(parts[len(parts)-1], ".")
	// holds 19
	isolatedValue := desiredValueSet[1]
	intVal, err:= strconv.Atoi(isolatedValue)
	if err != nil {
		fmt.Print(err)
	}
	// holds incremented string value 20
	intVal++
	
	// convert int back to string
	isolatedValue = strconv.Itoa(intVal)
	versionPart := parts[len(parts)-1]
	versionParts := strings.Split(versionPart, ".")
	versionParts[1] = isolatedValue
	newVersionPart := strings.Join(versionParts, ".")
	parts[len(parts)-1] = newVersionPart
	// holds updated string base_ref: release-0.20
	//updatedString := strings.Join(parts, "-")
	// left to do, commit/write new updated string back to file

	







	// Find the line containing the identifier
	snippetStartIdentifierB = "value:"
	lines = strings.Split(content, "\n")
	startIndex = -1
	endIndex = -1

	for i, line := range lines {
		if strings.Contains(line, snippetStartIdentifierB) {
			startIndex = i
			endIndex = i // Set endIndex to the same line as startIndex
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		return "error code snippet not found", nil // Snippet not found
	}


	// holds string for value: "release-0.19"
	snippetC := lines[startIndex]
	partsC := strings.Split(snippetC, "-")
	versionPartC := partsC[len(partsC)-1]
	versionPartsC := strings.Split(versionPartC, ".")
	minorVerStrC := strings.Trim(versionPartsC[1], `"`)
	minorVersionC, err := strconv.Atoi(minorVerStrC)
	if err != nil {
		fmt.Print(err)
	}
	minorVersionC++
	updatedStringC := strconv.Itoa(minorVersionC)
	newVersionStatementC := strings.Join([]string{versionPartsC[0], updatedStringC}, ".")
	partsC[len(partsC)-1] = "\"" + newVersionStatementC + "\""
	updatedStringC = strings.Join(partsC, "-")
	


	return updatedStringC, nil // snippet found

}






func updateFileContentsTwo()(string, error){

	// create client 
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN2")
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(accessToken)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoNameB, filePathThree, nil)
	if err != nil {
		fmt.Print(err)
	}

	// holds entire file content
	content, err := fileContent.GetContent() 
	if err != nil {
		fmt.Print(err)
	}


	

	// Find the line containing the identifier
	snippetStartIdentifierB := "jobName"
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
		return "error code snippet not found", nil // Snippet not found
	}

	// holds string for name: 
	snippet := lines[startIndex]

	splitSnippet := strings.Split(snippet, "-")

	// holds last release value (currently 19) --> left to do, increment value remerge string (convert from string to int then back to string)
	desiredSplitSnippet := splitSnippet[6]
	stringToInt, err := strconv.Atoi(desiredSplitSnippet)
	if err != nil {
		fmt.Print(err)
	}
	// holds string incremented "20"
	stringToInt++
	intToString := strconv.Itoa(stringToInt)
	aVersionPart := splitSnippet[len(splitSnippet)-2:]
	oldVersion := strings.Join(aVersionPart, "-")
	newVersion := strings.Replace(oldVersion, aVersionPart[1], intToString, 1)
	newVersionParts := strings.Split(newVersion, "-")
	splitSnippet[len(splitSnippet)-2] = newVersionParts[0]
	splitSnippet[len(splitSnippet)-1] = newVersionParts[1]
	//finalString := strings.Join(splitSnippet, "-")
	//left to do, commit/write new updated finalString back to file







	snippetStartIdentifierB = "value:"
	lines = strings.Split(content, "\n")
	startIndex = -1
	endIndex = -1

	for i, line := range lines {
		if strings.Contains(line, snippetStartIdentifierB) {
			startIndex = i
			endIndex = i // Set endIndex to the same line as startIndex
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		return "error code snippet not found", nil // Snippet not found
	}


	// holds string for value: "release-0.19"
	snippetB := lines[startIndex]
	partsB := strings.Split(snippetB, "-")
	versionPartB := partsB[len(partsB)-1]
	versionPartsB := strings.Split(versionPartB, ".")
	minorVerStrB := strings.Trim(versionPartsB[1], `"`)
	minorVersionB, err := strconv.Atoi(minorVerStrB)
	if err != nil {
		fmt.Print(err)
	}
	minorVersionB++
	updatedString := strconv.Itoa(minorVersionB)
	newVersionStatementB := strings.Join([]string{versionPartsB[0], updatedString}, ".")
	partsB[len(partsB)-1] = newVersionStatementB
	//updatedString = strings.Join(partsB, "-")
	//left to do, commit/write new updated finalString back to file





	// Find the line containing the identifier
	snippetStartIdentifierC := "baseRef:"
	lines = strings.Split(content, "\n")
	startIndex = -1
	endIndex = -1

	for i, line := range lines {
		if strings.Contains(line, snippetStartIdentifierC) {
			startIndex = i
			endIndex = i // Set endIndex to the same line as startIndex
			break
		}
	}
	if startIndex == -1 && endIndex == -1 {
		return "error code snippet not found!!!", nil // Snippet not found
	}

	// holds string for base_ref: release-0.19
	snippetC := lines[startIndex]
	// holds 0.19
	partsC := strings.Split(snippetC, "-")
	desiredValueSetC := strings.Split(partsC[len(partsC)-1], ".")
	// holds 19
	isolatedValueC := desiredValueSetC[1]
	intValC, err:= strconv.Atoi(isolatedValueC)
	if err != nil {
		fmt.Print(err)
	}
	// holds incremented string value 20
	intValC++
	
	// convert int back to string
	isolatedValueC = strconv.Itoa(intValC)
	versionPart := partsC[len(partsC)-1]
	versionParts := strings.Split(versionPart, ".")
	versionParts[1] = isolatedValueC
	newVersionPart := strings.Join(versionParts, ".")
	partsC[len(partsC)-1] = newVersionPart
	// holds updated string base_ref: release-0.20
	updatedStringC := strings.Join(partsC, "-")
	// left to do, commit/write new updated string back to file
	return updatedStringC, nil
	
}




/*
this command is responsible for accessing the eks-anywhere-prow-jobs repo
and updating 4 files:
1. jobs/aws/eks-anywhere-build-tooling/eks-anywhere-attribution-periodics-release-0.19.yaml
2. jobs/aws/eks-anywhere-build-tooling/eks-anywhere-checksum-periodics-release-0.14.yaml (cannot find as of 2024, confirm)
3. templater/jobs/periodic/eks-anywhere-build-tooling/eks-anywhere-attribution-periodics-release-0.19.yaml
4. templater/jobs/periodic/eks-anywhere-build-tooling/eks-anywhere-checksum-periodics-release-0.14.yaml (renamed as of 2024?)

depending on the file, the following will be updated :
	- name
	- base_ref
	- value
	- jobName
*/
