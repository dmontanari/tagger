// Copyright (c) 2026 Daniel Montanari. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tagger",
	Short: "\n\t. : Git tag Swiss Army Knife : .",
	Long:  "\n\t. : Git tag Swiss Army Knife : .",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Oops... " + err.Error())
		os.Exit(1)
	}
}

var verbose bool
var dryRun bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run")
}
