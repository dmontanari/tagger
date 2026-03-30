// Copyright (c) 2026 Daniel Montanari. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"
	"tagger/gitutil"

	"github.com/spf13/cobra"
)

var gitTags gitutil.GitTags
var remoteName string
var verbose bool
var dryRun bool

var rootCmd = &cobra.Command{
	Use:   "tagger",
	Short: "\n\t. : Git tag Swiss Army Knife : .",
	Long:  "\n\t. : Git tag Swiss Army Knife : .",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Oops... " + err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run")
	rootCmd.PersistentFlags().StringVarP(&remoteName, "remote", "r", "origin", "Remote name to use")

	// This PersistentPreRunE will run before any sub-command, ensuring gitTags is populated.
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Commands that do not need a repository path arg
		noPathCmds := map[string]bool{
			"tagger":     true,
			"help":       true,
			"version":    true,
			"completion": true,
		}
		if noPathCmds[cmd.Name()] {
			return nil
		}

		if len(args) < 1 {
			return fmt.Errorf("repository path argument is required for this command")
		}

		var err error
		from := args[0]
		gitTags, err = gitutil.NewGitTags(from, remoteName)
		if err != nil {
			// If there are no tags, that's not a fatal error for all commands (e.g. inc can start from 0)
			// but NewGitTags handles this gracefully. We'll let the specific commands decide.
			return err
		}
		return nil
	}
}
