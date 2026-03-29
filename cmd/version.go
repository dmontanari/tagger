// Copyright (c) 2026 Daniel Montanari. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCode string = "v0.0.1"

var versionCmd = &cobra.Command{

	Use:   "version",
	Short: "Show version.",
	Long:  "Show version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n\n\t. : Git Tag Swiss Army Knife : .\n\n")
		fmt.Printf("tagger version %s (c) 2026 Daniel Montanari. All rights reserved.\n\n", versionCode)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
