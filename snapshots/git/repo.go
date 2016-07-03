package git

import (
	"log"
	"strings"
)

// Utilities for operating on a git repo.
// Scoot may end up dealing with several git repos. E.g.: the snapshot store is
// owned by Scoot, but we may wish to ingest data from the user's source git
// repo, which is not owned by Scoot. We use the same code to run commands
// in a repo, but which commands we run vary.

type repository struct {
	dir string
}

func makeRepo(dir string, runner Runner) (repository, error) {
	// TODO(dbentley): make sure we handle the case that path is in a git repo,
	// but is not the root of a git repo
	repo := repository{dir}
	c := client{repo: repo, runner: runner}
	// TODO(dbentley): we'd prefer to use features present in git 2.5+, but are stuck on 2.4 for now
	topLevel, err := c.Run("rev-parse", "--show-toplevel")
	if err != nil {
		return repository{}, err
	}
	topLevel = strings.Replace(topLevel, "\n", "", -1)
	log.Println("Made a repo at ", dir, topLevel)
	repo.dir = topLevel
	return repo, nil
}
