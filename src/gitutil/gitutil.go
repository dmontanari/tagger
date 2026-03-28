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
)

type GitTags struct {
	Tags []TagInfo
	Repo string
}

type TagInfo struct {
	Name         string
	Date         time.Time
	MajorVersion int
	MinorVersion int
	PatchVersion int
}

func NewGitTags(gitpath string) (GitTags, error) {

	// We instantiate a new repository targeting the given path (the .git folder)
	tagInfo := []TagInfo{}
	gitTags := GitTags{Repo: gitpath}
	r, err := git.PlainOpen(gitpath)
	if err != nil {
		return gitTags, err
	}

	// List all tag references, both lightweight tags and annotated tags
	tagrefs, err := r.Tags()
	if err != nil {
		return gitTags, err
	}

	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		var tagDate time.Time

		// Tenta obter o objeto Tag (caso seja uma annotated tag)
		if tagObj, err := r.TagObject(t.Hash()); err == nil {
			tagDate = tagObj.Tagger.When
		} else if commitObj, err := r.CommitObject(t.Hash()); err == nil {
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

func (g GitTags) Newer() TagInfo {
	return g.Tags[0]
}

func (g GitTags) Older() TagInfo {
	return g.Tags[len(g.Tags)-1]
}

func (g GitTags) IncrementPatch() string {
	n := g.Newer()
	return fmt.Sprintf("v%d.%d.%d", n.MajorVersion, n.MinorVersion, n.PatchVersion+1)
}

func (g GitTags) IncrementMinor() string {
	n := g.Newer()
	return fmt.Sprintf("v%d.%d.0", n.MajorVersion, n.MinorVersion+1)
}

func (g GitTags) IncrementMajor() string {
	n := g.Newer()
	return fmt.Sprintf("v%d.0.0", n.MajorVersion+1)
}

func (g GitTags) Dump() {

	t := g.Tags

	fmt.Println("Tags ordenadas:")
	for _, tag := range t {
		fmt.Printf("%s  %s\n", tag.Date.Format("2006-01-02 15:04"), tag.Name)
	}

}
