package git

import (
	"log"
	"os"
	"os/exec"
)

type client struct {
	repo   repository
	runner Runner

	indexFile string
}

func (c *client) Run(args ...string) (string, error) {
	cmd := command{
		append([]string{"git"}, args...),
		nil,
		c.repo.dir,
	}

	cmd.env = make(map[string]string)
	if c.indexFile != "" {
		cmd.env["GIT_INDEX_FILE"] = c.indexFile
	}

	return c.runner.Run(cmd)
}

type command struct {
	args []string
	env  map[string]string
	cwd  string
}

type Runner interface {
	Run(cmd command) (string, error)
}

func NewExecRunner() Runner {
	return &execRunner{}
}

type execRunner struct{}

func (r *execRunner) Run(command command) (string, error) {
	cmd := exec.Command(command.args[0], command.args[1:]...)
	env := os.Environ()
	for k, v := range command.env {
		// TODO(dbentley): do we have to escape k?
		env = append(env, k+"="+v)
	}
	cmd.Env = env
	cmd.Dir = command.cwd
	log.Println("Running", command)
	out, err := cmd.Output()
	return string(out), err
}
