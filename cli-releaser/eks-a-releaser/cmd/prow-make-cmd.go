/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/command"
	"github.com/spf13/cobra"
)

// prowMakeCmdCmd represents the prowMakeCmd command
var prowMakeCmdCmd = &cobra.Command{
	Use:   "prow-make-cmd",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {
		runMakeCmd()
	},
}


// func that runs make command 
func runMakeCmd() error{

	projectRootFilePath := "github.com/testerIbix/eks-anywhere-prow-jobs/templater"

	commandSequence := fmt.Sprintf("make -C %s prowjobs templater", projectRootFilePath)
	makefileCmd := exec.Command("bash", "-c", commandSequence)
	_, err := command.ExecCommand(makefileCmd)
	if err != nil {
		return fmt.Errorf("error running make command: %v", err)
	}

	return nil

}