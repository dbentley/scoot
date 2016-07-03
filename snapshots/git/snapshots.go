// Package git uses git to implement snapshot interfaces
package git

// A Snapshots implementation that uses a git repository as a backing store
type GitRepoSnapshots struct {
	repo repository
}

func NewSnapshots(gitRepoDir string, runner Runner) (*GitRepoSnapshots, error) {
	cmd := command{
		[]string{"git", "init"},
		nil,
		gitRepoDir,
	}
	_, err := runner.Run(cmd)
	if err != nil {
		return nil, err
	}
	repo, err := makeRepo(gitRepoDir, runner)
	if err != nil {
		return nil, err
	}
	return &GitRepoSnapshots{repo}, nil
}
