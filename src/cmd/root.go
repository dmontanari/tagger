// Copyright (c) 2026 Daniel Montanari. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version bool
var versionCode = "v0.0.1"

var rootCmd = &cobra.Command{
	Use:   "tagger",
	Short: "\n\t. : Git tag Swiss Army Knife : .",
	Long:  "\n\t. : Git tag Swiss Army Knife : .",
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Printf("\n\n\t. : Git Tag Swiss Army Knife : .\n\n")
			fmt.Printf("tagger version %s (c) 2026 Daniel Montanari. All rights reserved.\n\n", versionCode)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Oops... " + err.Error())
		os.Exit(1)
	}
}

var verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Modo verboso")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Show version")
}
