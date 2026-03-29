package gitutil

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

	rawRemote := cfg.Raw.Section("remote").Subsection("origin")

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
