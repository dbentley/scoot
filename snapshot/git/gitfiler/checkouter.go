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

type checkoutRequest struct {
	id       string
	resultCh chan checkoutAndError
}

type checkoutAndError struct {
	checkout snapshot.Checkout
	err      error
}

// RefRepoCloningCheckouter checks out by cloning a Ref Repo.
type RefRepoCloningCheckouter struct {
	reqCh  chan checkoutRequest
	freeCh chan *repo.Repository

	getter    RefRepoGetter
	clonesDir *temp.TempDir

	ref *repo.Repository
	err error

	free []*repo.Repository
	reqs []checkoutRequest
}

func NewRefRepoCloningCheckouter(getter RefRepoGetter, tmp *temp.TempDir) *RefRepoCloningCheckouter {
	r := &RefRepoCloningCheckouter{
		getter:    getter,
		clonesDir: tmp,
		ref:       nil,
	}

	go r.loop()
	return r
}

func (c *RefRepoCloningCheckouter) loop() {
	c.ref, c.err = c.getter.Get()

	c.findClones()

	for {
		// Get input
		select {
		case <-c.doneCh:
			return
		case req := <-c.reqCh:
			c.reqs = append(c.reqs, req)
		case cloneResult := <-c.freeCh:
			c.free = append(c.free, cloneResult.repo)
			if c.err == nil {
				c.err = cloneResult.err
			}
		}

		// Serve requests we can serve now
		if len(c.free) > 0 && len(c.waiting) > 0 {
			clone, req := c.free[0], c.waiting[0]
			c.free, c.waiting = c.free[1:], c.waiting[1:]
			go c.checkoutAndSend(req, clone)
		}

		if len(c.free) == 0 {
			go c.cloneNewRepo()
		}
	}
}

func (c *RefRepoCloningCheckouter) cloneNewRepo() {
	repo, err := c.clone()
	c.cloneCh <- clone{repo, err}
}

func (c *RefRepoCloningCheckouter) checkoutAndSend(req checkoutRequest, clone *repo.Respository) {
	err := c.checkout(clone, req.id)
	if err != nil {
		req.resultCh <- checkoutAndError{nil, err}
		c.freeCh <- clone
	}
	req.resultCh <- checkoutAndError{&RefRepoCloningCheckout{repo: clone, id: req.id, checkouter: c}, nil}
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

func (c *RefRepoCloningCheckouter) Checkout(id string) (snapshot.Checkout, error) {
	resultCh := new(chan checkoutAndError)
	c.reqCh <- checkoutRequest{id, resultCh}
	result := <-resultCh
	return result.checkout, result.err
}

// Checkout checks out id (a raw git sha) into a Checkout.
// It does this by making a new clone (via reference) and checking out id.
func (c *RefRepoCloningCheckouter) Checkout(id string) (snapshot.Checkout, error) {
	clone, err := c.clone()
	if err != nil {
		return nil, err
	}

	if err := c.checkout(clone, id); err != nil {
		c.freeCh <- clone
		return nil, fmt.Errorf("gitfiler.RefRepoCloningCheckouter.Checkout: could not git checkout: %v", err)
	}

	return &RefRepoCloningCheckout{repo: clone, id: id, checkouter: c}, nil
}

// clone our reference repo into a new clone.
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

// checkout clone to be at id.
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

// release releases a repo so it can be used again.
func (c *RefRepoCloningCheckouter) release(release *repo.Repository) error {
	c.reqCh <- freeRepo{release}
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
	return c.checkouter.release(c.repo)
}
