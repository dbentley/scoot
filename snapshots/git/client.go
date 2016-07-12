package git

import (
	"fmt"
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

func (c *client) RunSha(args ...string) (string, error) {
	out, err := c.Run(args...)
	if err != nil {
		return out, err
	}
	return validateSha(out)
}

func validateSha(sha string) (string, error) {
	if len(sha) != 41 || sha[40] != '\n' {
		return "", fmt.Errorf("sha not 41 characters: %q", sha)
	}
	return sha[0:40], nil
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
