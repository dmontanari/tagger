package gitutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh"
	"github.com/kevinburke/ssh_config"
	cryptossh "golang.org/x/crypto/ssh"
)

// Check if the repo have a remote
func (g GitTags) HaveRemote() bool {

	r, err := g.repository.Remotes()
	if err != nil {
		return false
	}

	if len(r) == 0 {
		return false
	}

	for _, remote := range r {
		if remote.Config().Name == "origin" {
			return true
		}
	}

	return true
}

func (g GitTags) GetRemote() string {

	cfg, _ := g.repository.Config()

	rawRemote := cfg.Raw.Section("remote").Subsection(g.remoteName)

	// Note
	//
	// .git/config options for remote
	// [remote "origin"]
	// url = git@github.com:usuario/projeto.git
	// fetch = +refs/heads/*:refs/remotes/origin/*
	// pushurl = git@github.com:usuario/projeto-push.git
	// push = refs/heads/main:refs/heads/main
	// tagOpt = --tags

	pushUrl := rawRemote.Option("pushurl")
	if pushUrl != "" {
		return pushUrl
	}

	pushURL := rawRemote.Option("url")

	return pushURL

}

// Create AuthMethod to connect to remote
func (g GitTags) createAuthMethod() (transport.AuthMethod, error) {

	url := g.GetRemote()

	endpoint, _ := transport.NewEndpoint(url)

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
