/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"

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
// ibix16/eks-anywhere-prow-jobs/templater
func runMakeCmd(){

	ctx := build.Default

	remoteRepo, err := ctx.Import("github.com/ibix16/eks-anywhere-prow-jobs", "", 0)
	if err != nil {
		log.Fatalf("Error importing remote repo: %v",err)
		os.Exit(1)
	}

	templaterDir := filepath.Join(remoteRepo.Dir, "templater")


	makeCmd := exec.Command("make", "prowjobs", "-C", templaterDir)
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr

	err = makeCmd.Run()
	if err != nil {
		log.Fatalf("Error running make command: %v", err)
		os.Exit(1)
	}

}