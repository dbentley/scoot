package main

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/scootdev/scoot/snapshot/git/repo"
)

type shaAndError struct {
	sha string
	err error
}

type evalState struct {
	repo      *repo.Repository
	outCh     chan shaAndError
	leaseCh   chan struct{}
	releaseCh chan struct{}

	mu          sync.Mutex
	seenBlobs   map[string]BlobExpr
	seenTrees   map[string]TreeExpr
	seenCommits map[string]CommitExpr

	wg sync.WaitGroup
}

func (s *evalState) visitBlob(sha string, e BlobExpr) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.seenBlobs[sha]
	s.seenBlobs[sha] = e
	return ok
}

func (s *evalState) visitTree(sha string, e TreeExpr) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	prev := s.seenTrees[sha]
	if prev.Covers(e) {
		return true
	}
	prev.trees = prev.trees || e.trees
	prev.blobs = prev.trees || e.blobs
	s.seenTrees[sha] = prev
	return false
}

func (s *evalState) visitCommit(sha string, e CommitExpr) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	prev, _ := s.seenCommits[sha]
	if prev.Covers(e) {
		return true
	}
	// TODO(dbentley): how do we merge these?
	s.seenCommits[sha] = e
	return false
}

func Eval(e Expr, r *repo.Repository, sha string) chan shaAndError {
	s := &evalState{
		repo:        r,
		outCh:       make(chan shaAndError),
		seenBlobs:   make(map[string]BlobExpr),
		seenTrees:   make(map[string]TreeExpr),
		seenCommits: make(map[string]CommitExpr),
		leaseCh:     make(chan struct{}),
		releaseCh:   make(chan struct{}),
	}
	go leaseOnCh(s.leaseCh, s.releaseCh)
	s.Eval(e, sha)
	go s.wait()
	return s.outCh
}

func (s *evalState) Eval(e Expr, sha string) {
	s.wg.Add(1)
	go s.eval(e, sha)
}

func (s *evalState) eval(e Expr, sha string) {
	<-s.leaseCh
	defer func() {
		s.releaseCh <- struct{}{}
		s.wg.Done()
	}()
	switch e := e.(type) {
	case BlobExpr:
		s.evalBlob(e, sha)
	case TreeExpr:
		s.evalTree(e, sha)
	default:
		panic(fmt.Errorf("unknown type %T %v", e, e))
	}
}

func (s *evalState) evalBlob(e BlobExpr, sha string) {
	if s.visitBlob(sha, e) {
		return
	}
	if err := s.ShaExists(sha); err != nil {
		s.emitErr(err)
		return
	}
	s.emitSha(sha)
}

func (s *evalState) evalTree(e TreeExpr, sha string) {
	if s.visitTree(sha, e) {
		return
	}

	dirents, err := s.LsTree(sha)
	if err != nil {
		s.emitErr(err)
		return
	}

	for _, dirent := range dirents {
		if dirent.IsBlob() {
			if e.blobs {
				s.Eval(BlobExpr{}, hex.EncodeToString(dirent.Value))
			}
		} else {
			if e.trees {
				s.Eval(e, hex.EncodeToString(dirent.Value))
			}
		}
	}

	if e.trees {
		s.Eval(BlobExpr{}, sha)
	}
}

func (s *evalState) evalCommit(e CommitExpr, sha string) {
	if s.visitCommit(sha, e) {
		return
	}

	commit, err := s.LsCommit(sha)
	if err != nil {
		s.emitErr(err)
		return
	}

	s.Eval(e.tree, commit.Tree)

	if e.history > 0 {
		for _, sha := range commit.Parents {
			e2 := e
			e2.history--
			s.Eval(e2, sha)
		}
	}
}

func (s *evalState) wait() {
	s.wg.Wait()
	close(s.outCh)
	close(s.leaseCh)
}

func (s *evalState) emitErr(err error) {
	s.outCh <- shaAndError{err: err}
}

func (s *evalState) emitSha(sha string) {
	s.outCh <- shaAndError{sha: sha}
}

func (s *evalState) ShaExists(sha string) error {
	_, err := s.repo.Run("rev-parse", "--verify", sha+"^{object}")
	return err
}

func (s *evalState) LsTree(sha string) (repo.Tree, error) {
	contents, err := s.repo.Run("cat-file", "tree", sha)
	if err != nil {
		return nil, err
	}

	return repo.ParseTree([]byte(contents))

}

func (s *evalState) LsCommit(sha string) (repo.Commit, error) {
	return repo.Commit{}, fmt.Errorf("not yet implemented")
}
