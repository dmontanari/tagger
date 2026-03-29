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
}

type TagInfo struct {
	Name         string
	Date         time.Time
	Hash         plumbing.Hash
	MajorVersion int
	MinorVersion int
	PatchVersion int
}

func NewGitTags(gitpath string) (GitTags, error) {

	// We instantiate a new repository targeting the given path (the .git folder)
	var err error

	tagInfo := []TagInfo{}
	gitTags := GitTags{Repo: gitpath}

	gitTags.repository, err = git.PlainOpen(gitpath)
	if err != nil {
		return gitTags, err
	}

	// List all tag references, both lightweight tags and annotated tags
	tagrefs, err := gitTags.repository.Tags()
	if err != nil {
		return gitTags, err
	}

	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		var tagDate time.Time

		// Tenta obter o objeto Tag (caso seja uma annotated tag)
		if tagObj, err := gitTags.repository.TagObject(t.Hash()); err == nil {
			tagDate = tagObj.Tagger.When
		} else if commitObj, err := gitTags.repository.CommitObject(t.Hash()); err == nil {
			// Se for uma lightweight tag, pega a data do commit referenciado
			tagDate = commitObj.Committer.When
		}

		name := t.Name().Short()

		// Ignora tags que não começam com "v"
		if !strings.HasPrefix(name, "v") {
			return nil
		}

		var major, minor, patch int
		parts := strings.Split(name[1:], ".") // name[1:] remove o 'v' inicial
		if len(parts) < 3 {
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

	// Ordena a fatia (slice) usando a data. O ".After" garante ordem decrescente (mais recente primeiro).
	// Para ordem crescente (mais antiga primeiro), troque ".After" por ".Before".
	sort.Slice(tagInfo, func(i, j int) bool {
		return tagInfo[i].Date.After(tagInfo[j].Date)
	})

	gitTags.Tags = tagInfo

	return gitTags, nil

}
