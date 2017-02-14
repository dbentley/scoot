package cli

// package cli implements a cli for the SnapshotDB
// This is our first Scoot CLI that works very well with Cobra and also flags that
// main wants to set. How?
//
// First, main.go (either open-source, closed-source, or some future one) runs.
//
// main.go defines its own impl of DBInjector and constructs it; call it DBIImpl.
//
// main.go calls MakeDBCLI with DBIImpl
//
// [not yet needed/implemented] MakeDBCLI calls DBIImpl.RegisterFlags, which registers
//   the flags that main needs. These may be related to closed-source impls; e.g., which
//   backend server to use.

// MakeDBCLI creates the cobra commands and subcommands
//   (for each cobra command, there will be a dbCommand)
//   creatings the cobra command involves:
//     calling dbCommand.register(), which will register the common functionality flags
//     creating the cobra command with RunE as a wrapper function that will call the DBInjector()
//
// MakeDBCLI returns the root *cobra.Command
//
// main.go calls cmd.Execute()
//
// cobra will parse the command-line flags
//
// cobra will call cmd's RunE, which includes the wrapper defined in MakeDBCLI
//
// the wrapper will call DBInjector.Inject(), which will be DBIImpl.Inject()
//
// DBIImpl.Inject() will construct a SnapshotDB
// the wrapper will call dbCommand.run() with the db, the cobra command (which holds the
//   registered flags) and the additional command-line args
//
// dbCommand.run() does the work of calling a function on the SnapshotDB
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/scootdev/scoot/os/temp"
	"github.com/scootdev/scoot/snapshot"
	"github.com/scootdev/scoot/snapshot/git/repo"
)

type DBInjector interface {
	// TODO(dbentley): we probably want a way to register flags
	RegisterFlags(cmd *cobra.Command)
	Inject() (snapshot.DB, error)
}

func MakeDBCLI(injector DBInjector) *cobra.Command {
	rootCobraCmd := &cobra.Command{
		Use:   "scoot-snapshot-db",
		Short: "scoot snapshot db CLI",
	}

	injector.RegisterFlags(rootCobraCmd)

	add := func(subCmd dbCommand, parentCobraCmd *cobra.Command) {
		cmd := subCmd.register()
		cmd.RunE = func(innerCmd *cobra.Command, args []string) error {
			db, err := injector.Inject()
			if err != nil {
				return err
			}
			return subCmd.run(db, innerCmd, args)
		}
		parentCobraCmd.AddCommand(cmd)
	}

	createCobraCmd := &cobra.Command{
		Use:   "create",
		Short: "create a snapshot",
	}
	rootCobraCmd.AddCommand(createCobraCmd)

	add(&ingestGitCommitCommand{}, createCobraCmd)
	add(&ingestDirCommand{}, createCobraCmd)
	add(&mergeLoopCommand{}, createCobraCmd)
	add(&mergeTestCommand{}, createCobraCmd)

	readCobraCmd := &cobra.Command{
		Use:   "read",
		Short: "read data from a snapshot",
	}
	rootCobraCmd.AddCommand(readCobraCmd)

	add(&catCommand{}, readCobraCmd)

	exportCobraCmd := &cobra.Command{
		Use:   "export",
		Short: "export a snapshot",
	}
	rootCobraCmd.AddCommand(exportCobraCmd)

	add(&exportGitCommitCommand{}, exportCobraCmd)

	return rootCobraCmd
}

type dbCommand interface {
	register() *cobra.Command
	run(db snapshot.DB, cmd *cobra.Command, args []string) error
}

type ingestGitCommitCommand struct {
	commit string
}

func (c *ingestGitCommitCommand) register() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingest_git_commit",
		Short: "ingests a git commit into the repo in cwd",
	}
	cmd.Flags().StringVar(&c.commit, "commit", "", "commit to ingest")
	return cmd
}

func (c *ingestGitCommitCommand) run(db snapshot.DB, _ *cobra.Command, _ []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get working directory: wd")
	}

	ingestRepo, err := repo.NewRepository(wd)
	if err != nil {
		return fmt.Errorf("not a valid repo dir: %v, %v", wd, err)
	}

	id, err := db.IngestGitCommit(ingestRepo, c.commit)
	if err != nil {
		return err
	}

	fmt.Println(id)
	return nil
}

type exportGitCommitCommand struct {
	id string
}

func (c *exportGitCommitCommand) register() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "to_git_commit",
		Short: "exports a GitCommitSnapshot identified by id into the repo in cwd",
	}
	cmd.Flags().StringVar(&c.id, "id", "", "id to export")
	return cmd
}

func (c *exportGitCommitCommand) run(db snapshot.DB, _ *cobra.Command, _ []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get working directory: wd")
	}

	exportRepo, err := repo.NewRepository(wd)
	if err != nil {
		return fmt.Errorf("not a valid repo dir: %v, %v", wd, err)
	}

	commit, err := db.ExportGitCommit(snapshot.ID(c.id), exportRepo)
	if err != nil {
		return err
	}

	fmt.Println(commit)
	return nil
}

type ingestDirCommand struct {
	dir string
}

func (c *ingestDirCommand) register() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingest_dir",
		Short: "ingests a directory into the repo in cwd",
	}
	cmd.Flags().StringVar(&c.dir, "dir", "", "dir to ingest")
	return cmd
}

func (c *ingestDirCommand) run(db snapshot.DB, _ *cobra.Command, _ []string) error {
	id, err := db.IngestDir(c.dir)
	if err != nil {
		return err
	}

	fmt.Println(id)
	return nil
}

type catCommand struct {
	id string
}

func (c *catCommand) register() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cat",
		Short: "concatenate files from an FSSnapshot to stdout",
	}
	cmd.Flags().StringVar(&c.id, "id", "", "Snapshot ID to read from")
	return cmd
}

func (c *catCommand) run(db snapshot.DB, _ *cobra.Command, filenames []string) error {
	id := snapshot.ID(c.id)
	for _, filename := range filenames {
		data, err := db.ReadFileAll(id, filename)
		if err != nil {
			return err
		}
		fmt.Printf("%s", data)
	}
	return nil
}

type mergeLoopCommand struct {
	jenkinsURL string
}

func (c *mergeLoopCommand) register() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge-loop",
		Short: "try other stuff",
	}
	cmd.Flags().StringVar(&c.jenkinsURL, "jenkins_url", "", "base URL for jenkins")
	return cmd
}

func (c *mergeLoopCommand) run(db snapshot.DB, _ *cobra.Command, _ []string) error {
	for {
		log.Println("merge loop iteration")
		cmd, err := c.findParams(db)
		if err != nil {
			return err
		}
		err = cmd.run(db, nil, nil)
		if err != nil {
			return err
		}
	}
}

type JenkinsJobInfo struct {
	Builds []BuildEntry
}

type BuildEntry struct {
	URL string
}

func (c *mergeLoopCommand) findParams(db snapshot.DB) (*mergeTestCommand, error) {
	jsonURL := c.jenkinsURL + "/api/json"
	log.Println("fetching", jsonURL)
	resp, err := http.Get(jsonURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("job info text %s", body)

	var data JenkinsJobInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	log.Println("job info", data)

	jsonURL = data.Builds[0].URL + "/api/json"
	resp, err = http.Get(jsonURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("f info text", body)

	var j map[string]json.RawMessage
	err = json.Unmarshal(body, &j)
	entries := make([]map[string]string, 0)
	err = json.Unmarshal(j["actions"], &entries)
	params := make(map[string]string)
	for _, entry := range entries {
		params[entry["name"]] = entry["value"]
	}
	log.Println("Hmm", params)

	branch := params["branch"]
	repoURL := params["temp_branch_repo_url"]
	return &mergeTestCommand{branch: branch, repoURL: repoURL}, nil
}

type mergeTestCommand struct {
	repoURL string
	branch  string
}

func (c *mergeTestCommand) register() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge-test",
		Short: "try some stuff",
	}
	cmd.Flags().StringVar(&c.repoURL, "repo_url", "", "base URL for repo")
	cmd.Flags().StringVar(&c.branch, "branch", "", "branch to get")
	return cmd
}

func (c *mergeTestCommand) run(db snapshot.DB, _ *cobra.Command, _ []string) error {
	log.Println("Merge test", c.repoURL, c.branch)

	tmp, err := temp.TempDirDefault()
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get working directory: wd")
	}

	ingestRepo, err := repo.NewRepository(wd)
	if err != nil {
		return fmt.Errorf("not a valid repo dir: %v, %v", wd, err)
	}

	localBranch := fmt.Sprintf("fetched/%s", c.branch)
	refSpec := fmt.Sprintf("+%s:%s", c.branch, localBranch)
	if _, err := ingestRepo.Run("fetch", c.repoURL, refSpec); err != nil {
		return err
	}

	proposalSHA, err := ingestRepo.RunSha("rev-parse", localBranch)
	if err != nil {
		return err
	}

	if _, err := ingestRepo.Run("fetch"); err != nil {
		return err
	}

	masterSHA, err := ingestRepo.RunSha("rev-parse", "origin/master")
	if err != nil {
		return err
	}

	mergeBaseSHA, err := ingestRepo.RunSha("merge-base", proposalSHA, masterSHA)
	if err != nil {
		return err
	}

	f, err := tmp.TempFile("merge-index-")
	if err != nil {
		return err
	}

	fname := f.Name()
	f.Close()
	os.Remove(fname)

	log.Println("ba-dum", mergeBaseSHA, masterSHA, proposalSHA)
	cmd := ingestRepo.Command("read-tree", "-i", "-m", "--aggressive", mergeBaseSHA, masterSHA, proposalSHA)
	indexEnv := fmt.Sprintf("GIT_INDEX_FILE=%s", fname)
	cmd.Env = append(os.Environ(), indexEnv)
	log.Println("env", indexEnv)

	if _, err := ingestRepo.RunCmd(cmd); err != nil {
		return err
	}

	cmd = ingestRepo.Command("write-tree")
	cmd.Env = append(os.Environ(), indexEnv)
	treeSHA, err := ingestRepo.RunCmdSha(cmd)
	if err != nil {
		return err
	}

	log.Println("merged", treeSHA)

	msg, err := ingestRepo.Run("log", "--format=%B", "-n", "1", proposalSHA)
	if err != nil {
		return err
	}

	name, err := ingestRepo.Run("show", "-s", "--format=%aN", proposalSHA)
	if err != nil {
		return err
	}

	email, err := ingestRepo.Run("show", "-s", "--format=%aE", proposalSHA)
	if err != nil {
		return err
	}

	cmd = ingestRepo.Command("commit-tree", treeSHA, "-p", masterSHA, "-m", msg)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME="+name,
		"GIT_AUTHOR_EMAIL="+email,
	)

	created, err := ingestRepo.RunCmdSha(cmd)
	if err != nil {
		return err
	}

	log.Println("Created new one", created)
	return nil
}
