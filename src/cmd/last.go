// Copyright (c) 2026 Daniel Montanari (dmontanari@gmail.com). All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"
	"tagger/gitutil"

	"github.com/spf13/cobra"
)

var fullOutput bool

var lastCmd = &cobra.Command{

	Use:   "last [repository path]",
	Short: "Return last tag in repository path.",
	Long:  "last [repository path] Return last tag in repository path.",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		from := args[0]
		tags, err := gitutil.NewGitTags(from)

		if err != nil {
			panic(err)
		}

		tag := tags.Newer()
		if fullOutput {
			fmt.Printf("%s  %s\n", tag.Date.Format("2006-01-02 15:04"), tag.Name)
		} else {
			fmt.Println(tag.Name)
		}
		os.Exit(0)
	},
}

func init() {
	lastCmd.Flags().BoolVarP(&fullOutput, "full", "f", false, "Full output - date and tag")
	rootCmd.AddCommand(lastCmd)
}
