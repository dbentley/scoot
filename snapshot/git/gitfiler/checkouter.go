package gitfiler

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"sync"

	"github.com/scootdev/scoot/os/temp"
	"github.com/scootdev/scoot/snapshot"
	"github.com/scootdev/scoot/snapshot/git/repo"
)

// Utilities for Reference Repositories.
// A Reference Repository is a way to clone repos locally so that the clone takes less time and disk space.
// By passing --reference <local path> to a git clone, the clone will not copy the whole ODB but instead
// hardlink. This means the clone is much faster and also takes very little extra hard disk space.
// Cf. https://git-scm.com/docs/git-clone

// RefRepoGetter lets a client get a Repository to use as a Reference Repository.
type RefRepoGetter interface {
	// Get gets the Repository to use as a reference repository.
	Get() (*repo.Repository, error)
}

// RefRepoCloningCheckouter checks out by cloning a Ref Repo.
type RefRepoCloningCheckouter struct {
	getter    RefRepoGetter
	clonesDir *temp.TempDir

	ref  *repo.Repository
	busy []*repo.Repository
	free []*repo.Repository

	mu sync.Mutex
}

func NewRefRepoCloningCheckouter(getter RefRepoGetter, tmp *temp.TempDir) *RefRepoCloningCheckouter {
	r := &RefRepoCloningCheckouter{
		getter:    getter,
		clonesDir: tmp,
		ref:       nil,
	}

	r.findClones()
	return r
}

// findClones finds all the valid clones in clonesDir
func (c *RefRepoCloningCheckouter) findClones() {
	fis, err := ioutil.ReadDir(c.clonesDir.Dir)
	if err != nil {
		return
	}

	for _, fi := range fis {
		if clone, err := repo.NewRepository(path.Join(c.clonesDir.Dir, fi.Name())); err == nil {
			c.free = append(c.free, clone)
		}
	}
}

// Checkout checks out id (a raw git sha) into a Checkout.
// It does this by making a new clone (via reference) and checking out id.
func (c *RefRepoCloningCheckouter) Checkout(id string) (snapshot.Checkout, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ref == nil {
		ref, err := c.getter.Get()
		if err != nil {
			return nil, fmt.Errorf("gitfiler.RefRepoCloningCheckouter.clone: error getting: %v", err)
		}
		c.ref = ref
	}

	if len(c.free) == 0 {
		clone, err := c.clone()
		if err != nil {
			return nil, err
		}
		c.free = []*repo.Repository{clone}
	}

	clone := c.free[0]

	if err := c.checkout(clone, id); err != nil {
		return nil, fmt.Errorf("gitfiler.RefRepoCloningCheckouter.Checkout: could not git checkout: %v", err)
	}

	// move c.free[0] to the end of c.busy
	c.busy, c.free = append(c.busy, c.free[0]), c.free[1:]

	log.Println("gitfiler.RefRepoCloningCheckouter.Checkout done: ", clone.Dir())
	return &RefRepoCloningCheckout{repo: clone, id: id, checkouter: c}, nil
}

func (c *RefRepoCloningCheckouter) clone() (*repo.Repository, error) {
	cloneDir, err := c.clonesDir.TempDir("clone-")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("git", "clone", "-n", "--reference", c.ref.Dir(), c.ref.Dir(), cloneDir.Dir)
	log.Println("gitfiler.RefRepoCloningCheckouter.clone: Cloning", cmd)
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("gitfiler.RefRepoCloningCheckouter.clone: error cloning: %v", err)
	}

	return repo.NewRepository(cloneDir.Dir)
}

func (c *RefRepoCloningCheckouter) checkout(clone *repo.Repository, id string) error {
	// -d removes directories. -x ignores gitignore and removes everything.
	// -f is force. -f the second time removes directories even if they're git repos themselves
	cmds := [][]string{
		{"clean", "-f", "-f", "-d", "-x"},
		{"checkout", id},
	}

	for _, argv := range cmds {
		if _, err := clone.Run(argv...); err != nil {
			return err
		}

	}
	return nil
}

func (c *RefRepoCloningCheckouter) release(release *repo.Repository) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, r := range c.busy {
		if r == release {
			c.busy = append(c.busy[:i], c.busy[i+1:]...)
		}
	}
	c.free = append(c.free, release)
	return nil
}

type RefRepoCloningCheckout struct {
	repo       *repo.Repository
	id         string
	checkouter *RefRepoCloningCheckouter
}

func (c *RefRepoCloningCheckout) Path() string {
	return c.repo.Dir()
}

func (c *RefRepoCloningCheckout) ID() string {
	return c.id
}

func (c *RefRepoCloningCheckout) Release() error {
	c.checkouter.checkout(c.repo, c.id)
	return c.checkouter.release(c.repo)
}
