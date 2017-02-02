package gitdb

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/scootdev/scoot/snapshot/git/repo"
)

func (db *DB) ingestDir(dir string) (snapshot, error) {
	// We ingest a dir using git commands:
	// First, create a new index file.
	// Second, add all the files in the work tree.
	// Third, write the tree.
	// This doesn't create a commit, or otherwise mess with repo state.
	indexDir, err := db.tmp.TempDir("git-index")
	if err != nil {
		return nil, err
	}

	indexFilename := filepath.Join(indexDir.Dir, "index")
	defer os.RemoveAll(indexDir.Dir)

	env := append(os.Environ(), "GIT_INDEX_FILE="+indexFilename, "GIT_WORK_TREE="+dir)

	// TODO(dbentley): should we use update-index instead of add? Maybe add looks at repo state
	// (e.g., HEAD) and we should just use the lower-level plumbing command?
	cmd := db.dataRepo.Command("add", ".")
	cmd.Env = env
	_, err = db.dataRepo.RunCmd(cmd)
	if err != nil {
		return nil, err
	}

	cmd = db.dataRepo.Command("write-tree")
	cmd.Env = env
	sha, err := db.dataRepo.RunCmdSha(cmd)
	if err != nil {
		return nil, err
	}

	return &localSnapshot{sha: sha, kind: kindFSSnapshot}, nil
}

const tempRef = "refs/heads/scoot/__temp_for_writing"

func (db *DB) ingestGitCommit(ingestRepo *repo.Repository, commitish string) (snapshot, error) {
	sha, err := ingestRepo.RunSha("rev-parse", "--verify", fmt.Sprintf("%s^{commit}", commitish))
	if err != nil {
		return nil, fmt.Errorf("not a valid commit: %s, %v", commitish, err)
	}

	if err := db.shaPresent(sha); err == nil {
		return &localSnapshot{sha: sha, kind: kindGitCommitSnapshot}, nil
	}

	if err := moveCommit(ingestRepo, db.dataRepo, sha); err != nil {
		return nil, err
	}

	return &localSnapshot{sha: sha, kind: kindGitCommitSnapshot}, nil
}

func (db *DB) shaPresent(sha string) error {
	_, err := db.dataRepo.Run("rev-parse", "--verify", sha+"^{object}")
	return err
}
