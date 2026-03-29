package gitutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh"
	"github.com/kevinburke/ssh_config"
	cryptossh "golang.org/x/crypto/ssh"
)

func createAuthMethod(endpoint transport.Endpoint) (transport.AuthMethod, error) {

	host := endpoint.Host

	// Get IdentityFile from .ssh/config
	// Expand "~" to full path
	keyPath := ssh_config.Get(host, "IdentityFile")
	if keyPath == "" {
		return nil, fmt.Errorf("publik key not set")
	}
	if strings.HasPrefix(keyPath, "~/") {
		home, _ := os.UserHomeDir()
		keyPath = filepath.Join(home, keyPath[2:])
	}

	// Get username from .ssh/config
	user := ssh_config.Get(host, "User")
	if user == "" {
		// Default fallback: user git
		user = "git"
	}

	// Load key
	if _, err := os.Stat(keyPath); err == nil {
		auth, err := ssh.NewPublicKeysFromFile(user, keyPath, "")
		if err == nil {
			auth.HostKeyCallback = cryptossh.InsecureIgnoreHostKey()
			return auth, nil
		}
	}

	// Fallback to ssh-agent
	auth, err := ssh.NewSSHAgentAuth(user)
	if err == nil {
		auth.HostKeyCallback = cryptossh.InsecureIgnoreHostKey()
		return auth, nil
	}

	return nil, fmt.Errorf("error creating auth method")

}

func (g GitTags) Push(tagName string) error {

	url := g.GetRemote()

	endpoint, _ := transport.NewEndpoint(url)

	authMethod, err := createAuthMethod(*endpoint)
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
