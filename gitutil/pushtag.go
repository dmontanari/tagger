package gitutil

import (
	"fmt"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
)

func (g GitTags) Push(tagName string) error {

	authMethod, err := g.createAuthMethod()
	if err != nil {
		return err
	}

	err = g.repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       authMethod,
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("refs/tags/%s:refs/tags/%s", tagName, tagName)),
		},
	})

	if err != nil {
		return err
	}

	return nil
}
