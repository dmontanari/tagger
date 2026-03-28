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

var incMajor bool
var incMinor bool
var incPatch bool

var incCmd = &cobra.Command{

	Use:   "inc [repository path] [--dry-run|-d]",
	Short: "Create new tag incrementing version number.",
	Long: `inc [repository path] [flats] Create new tag incrementing version number. 
	Tags must follow the pattern vM.m.p.
	Incrementing a higher version level resets lower ones (e.g., a Major bump on v2.1.35 results in v3.0.0).
	`,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		from := args[0]
		if !incMajor && !incMinor && !incPatch {
			fmt.Println("What version number do you want me to increment? Read the help!")
			os.Exit(1)
		}

		tags, err := gitutil.NewGitTags(from)

		if err != nil {
			panic(err)
		}

		IncVersion(tags)

		if !dryRun {
			fmt.Println("Pushing to remote not implemented yet. WIP")
		}

		os.Exit(0)
	},
}

func IncVersion(tags gitutil.GitTags) {

	if verbose {
		fmt.Printf("%s -> ", tags.Newer().Name)
	}

	if incMajor {
		fmt.Println(tags.IncrementMajor())
	} else if incMinor {
		fmt.Println(tags.IncrementMinor())
	} else if incPatch {
		fmt.Println(tags.IncrementPatch())
	}

}

func init() {

	incCmd.Flags().BoolVarP(&incMajor, "major", "M", false, "Increment major version")
	incCmd.Flags().BoolVarP(&incMinor, "minor", "m", false, "Increment minor version")
	incCmd.Flags().BoolVarP(&incPatch, "patch", "p", false, "Increment patch version")

	rootCmd.AddCommand(incCmd)
}
