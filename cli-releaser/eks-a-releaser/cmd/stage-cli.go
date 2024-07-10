/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cliReleaseNumPath  = "release/triggers/eks-a-release/development/RELEASE_NUMBER"
	cliReleaseVerPath  = "release/triggers/eks-a-release/development/RELEASE_VERSION"
)

// stageCliCmd represents the stageCli command
var stageCliCmd = &cobra.Command{
	Use:   "stage-cli",
	Short: "increments version files to trigger staging bundle release",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. `,

	Run: func(cmd *cobra.Command, args []string) {
		returnedContent, err := updateAllStageCliFiles()
		if err != nil {
			fmt.Print(err)
		}
		fmt.Print(returnedContent)
	},
}


// this function is responsible for updating the release number file 
// the function accesses the trigger file and retrieves the value assigned to bundle number : #
// a new commit and PR is then created using the retrieved value from the trigger file
func updateAllStageCliFiles()(error){


}



func updateReleaseNumber()(){

}

func updateReleaseVersion()(){

}