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
