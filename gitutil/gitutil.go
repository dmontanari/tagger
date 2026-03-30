// Copyright (c) 2026 Daniel Montanari. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitutil

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
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

	r, err := git.PlainOpen(gitpath)
	if err != nil {
		return gitTags, err
	}
	gitTags.repository = r

	// Create auth method for remote operations
	auth, err := gitTags.createAuthMethod()
	if err != nil {
		return gitTags, fmt.Errorf("failed to create auth method: %w", err)
	}

	// Fetch all tags from the remote to ensure our local repo is up-to-date
	err = r.Fetch(&git.FetchOptions{
		RemoteName: remoteName,
		Auth:       auth,
		Tags:       plumbing.AllTags,
		Prune:      true,
	})
	if err != nil && err.Error() != "already up-to-date" {
		return gitTags, fmt.Errorf("failed to fetch tags from remote: %w", err)
	}

	// Now that we have fetched, we can iterate over local tags, which are now in sync with the remote
	tagrefs, err := r.Tags()
	if err != nil {
		return gitTags, err
	}

	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		var tagDate time.Time

		// Get tag object. For annotated tags, this hash points to a tag object.
		// For lightweight tags, it points to a commit object.
		obj, err := r.Object(plumbing.AnyObject, t.Hash())
		if err != nil {
			return nil // Should not happen, but let's be safe
		}

		switch o := obj.(type) {
		case *object.Tag:
			// It's an annotated tag
			tagDate = o.Tagger.When
		case *object.Commit:
			// It's a lightweight tag
			tagDate = o.Committer.When
		default:
			// Not a tag or commit, maybe a tree or blob? Skip.
			return nil
		}

		name := t.Name().Short()

		// Search for tags with "v" pattern (vM.m.p)
		if !strings.HasPrefix(name, "v") {
			return nil // continue
		}

		var major, minor, patch int
		parts := strings.Split(name[1:], ".") // name[1:] remove 'v'
		if len(parts) < 3 {
			// Using return nil to continue ForEach loop
			return nil
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
		return nil
	})

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
