// Copyright (c) 2026 Daniel Montanari. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
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

		newTag := IncVersion(gitTags)

		if !dryRun {
			CreateTag(gitTags, newTag)
			err := gitTags.Push(newTag)
			if err != nil {
				fmt.Println(err)
			}
		}

	},
}

func IncVersion(tags gitutil.GitTags) string {

	if verbose {
		fmt.Printf("%s -> ", tags.Newer().Name)
	}

	var newTag string

	if incMajor {
		newTag = tags.IncrementMajor()
	} else if incMinor {
		newTag = tags.IncrementMinor()
	} else if incPatch {
		newTag = tags.IncrementPatch()
	}

	fmt.Println(newTag)

	return newTag

}

func CreateTag(tags gitutil.GitTags, newTag string) {

	if verbose {
		fmt.Printf("%s -> ", tags.Newer().Name)
	}

	fmt.Println(newTag)
	tags.CreateTag(newTag)

}

func init() {

	incCmd.Flags().BoolVarP(&incMajor, "major", "M", false, "Increment major version")
	incCmd.Flags().BoolVarP(&incMinor, "minor", "m", false, "Increment minor version")
	incCmd.Flags().BoolVarP(&incPatch, "patch", "p", false, "Increment patch version")

	rootCmd.AddCommand(incCmd)
}
