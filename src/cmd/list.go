// Copyright (c) 2026 Daniel Montanari. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"os"
	"tagger/gitutil"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{

	Use:   "list [repository path]",
	Short: "List all tags in repository path.",
	Long:  "list [repository path] List all tags in repository path.",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		from := args[0]
		tags, err := gitutil.NewGitTags(from)

		if err != nil {
			panic(err)
		}

		tags.Dump()

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
