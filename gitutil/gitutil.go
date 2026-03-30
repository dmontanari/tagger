// Copyright (c) 2026 Daniel Montanari. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitutil

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type GitTags struct {
	Tags       []TagInfo
	Repo       string
	repository *git.Repository
	remoteName string
}

type TagInfo struct {
	Name         string
	Date         time.Time
	Hash         plumbing.Hash
	MajorVersion int
	MinorVersion int
	PatchVersion int
}

func NewGitTags(gitpath string, remoteName string) (GitTags, error) {

	// We instantiate a new repository targeting the given path (the .git folder)
	var err error

	tagInfo := []TagInfo{}
	gitTags := GitTags{Repo: gitpath, remoteName: remoteName}

	gitTags.repository, err = git.PlainOpen(gitpath)
	if err != nil {
		return gitTags, err
	}

	// Go to remote
	remote, err := gitTags.repository.Remote(gitTags.remoteName)
	if err != nil {
		return gitTags, err
	}
	auth, err := gitTags.createAuthMethod()
	if err != nil {
		return gitTags, err
	}

	// List all tag references, both lightweight tags and annotated tags
	// tagrefs, err := gitTags.repository.Tags()
	tagrefs, err := remote.List(&git.ListOptions{Auth: auth})
	if err != nil {
		return gitTags, err
	}

	for _, t := range tagrefs {
		var tagDate time.Time

		// Get tag object
		if tagObj, err := gitTags.repository.TagObject(t.Hash()); err == nil {
			tagDate = tagObj.Tagger.When
		} else if commitObj, err := gitTags.repository.CommitObject(t.Hash()); err == nil {
			tagDate = commitObj.Committer.When
		}

		name := t.Name().Short()

		// Search for tags with "v" pattern (vM.m.p)
		if !strings.HasPrefix(name, "v") {
			continue
		}

		var major, minor, patch int
		parts := strings.Split(name[1:], ".") // name[1:] remove 'v'
		if len(parts) < 3 {
			continue
		}

		major, _ = strconv.Atoi(parts[0])
		minor, _ = strconv.Atoi(parts[1])
		patch, _ = strconv.Atoi(parts[2])

		tagInfo = append(tagInfo, TagInfo{
			Name:         name,
			Date:         tagDate,
			Hash:         t.Hash(),
			MajorVersion: major,
			MinorVersion: minor,
			PatchVersion: patch,
		})
	}

	if err != nil {
		return gitTags, err
	}

	// Order by date
	sort.Slice(tagInfo, func(i, j int) bool {
		return tagInfo[i].Date.After(tagInfo[j].Date)
	})

	gitTags.Tags = tagInfo

	return gitTags, nil

}
