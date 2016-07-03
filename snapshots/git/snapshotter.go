package git

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

// GitSnapshotter snapshots using git repo information. This can make snapshotting easier,
// because git does snapshotting.
// (But comes with git's cost, so we do not expect this to perform well enough that
// we can snapshot on each save.
type GitSnapshotter struct {
	snaps    *GitRepoSnapshots
	scootDir string
	runner   Runner
}

func NewSnapshotter(snaps *GitRepoSnapshots, scootDir string, runner Runner) (*GitSnapshotter, error) {
	return &GitSnapshotter{snaps, scootDir, runner}, nil
}

func (s *GitSnapshotter) Snapshot(path string) (string, error) {
	repo, err := makeRepo(path, s.runner)
	if err != nil {
		return "", err
	}
	log.Printf("Creating snapshot for %q", path)

	c := client{repo: repo, runner: s.runner}

	if err := s.copyIndex(&c); err != nil {
		return "", err
	}

	log.Println("Using index file at ", c.indexFile)

	return "", fmt.Errorf("gitsnapshotter.Snapshot not yet implemented")
}

// copies the git index file from the client git directory to a scratch dir in our scoot dir.
// We will then use this index file for operations so that we don't mutate the client's state.
// Copying the index file is an optimization that allows the client to not read every file all at once.
func (s *GitSnapshotter) copyIndex(c *client) error {

	scratchDir := path.Join(s.scootDir, "snapshotter")
	err := os.MkdirAll(scratchDir, 0700)

	if err != nil {
		return err
	}

	indexFile := path.Join(scratchDir, "index")
	c.indexFile = indexFile

	// TODO(dbentley): we'd like to use git rev-parse --git-path index, but
	// that requires git 2.5 and we have to support git 2.4
	src, err := os.Open(path.Join(c.repo.dir, ".git", "index"))
	if err != nil {
		log.Println("Couldn't copy", path.Join(c.repo.dir, ".git", "index"))
		// We couldn't open the index, but that's actually fine, becase copying the index is just to save us work, not required
		return nil
	}
	defer src.Close()

	dst, err := os.Create(indexFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	log.Println("Copying from ", path.Join(c.repo.dir, ".git", "index"), " to ", indexFile)
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}
