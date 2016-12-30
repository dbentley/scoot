package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/libgit2/git2go"
	"github.com/scootdev/scoot/snapshot/git/repo"
)

type shaAndError struct {
	sha string
	err error
}

type evalState struct {
	repo      *git.Repository
	odb       *git.Odb
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
	prev.Trees = prev.Trees || e.Trees
	prev.Blobs = prev.Trees || e.Blobs
	s.seenTrees[sha] = prev
	return false
}

func (s *evalState) visitCommit(sha string, e CommitExpr) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	prev, ok := s.seenCommits[sha]
	if ok && prev.Covers(e) {
		return true
	}
	// TODO(dbentley): how do we merge these?
	s.seenCommits[sha] = e
	return false
}

func Eval(r *git.Repository, specs []EvalSpec) chan shaAndError {
	odb, err := r.Odb()
	if err != nil {
		panic(err)
	}
	s := &evalState{
		repo:        r,
		odb:         odb,
		outCh:       make(chan shaAndError),
		seenBlobs:   make(map[string]BlobExpr),
		seenTrees:   make(map[string]TreeExpr),
		seenCommits: make(map[string]CommitExpr),
		leaseCh:     make(chan struct{}),
		releaseCh:   make(chan struct{}),
	}
	go leaseOnCh(s.leaseCh, s.releaseCh)
	for _, spec := range specs {
		s.Eval(spec.Expr(), spec.SHA)
	}
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
	case CommitExpr:
		s.evalCommit(e, sha)
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
		switch {
		case bytes.Equal(dirent.Mode, repo.ModeBlobExec):
			fallthrough
		case bytes.Equal(dirent.Mode, repo.ModeBlob):
			fallthrough
		case bytes.Equal(dirent.Mode, repo.ModeSymlink):
			if e.Blobs {
				s.Eval(BlobExpr{}, hex.EncodeToString(dirent.Value))
			}
		case bytes.Equal(dirent.Mode, repo.ModeTree):
			if e.Trees {
				s.Eval(e, hex.EncodeToString(dirent.Value))
			}
		case bytes.Equal(dirent.Mode, repo.ModeCommit):
			continue
		default:
			panic(fmt.Errorf("unrecognized mode %q %q %q", dirent.Mode, dirent.Name, hex.EncodeToString(dirent.Value)))
		}
	}

	if e.Trees {
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

	s.Eval(e.Tree, commit.Tree)
	s.Eval(BlobExpr{}, sha)

	if e.History != 0 {
		e2 := e
		if e.History != -1 {
			e2.History--
		}
		for _, sha := range commit.Parents {
			s.Eval(e2, sha)
		}
	}
}

func (s *evalState) wait() {
	s.wg.Wait()
	close(s.outCh)
	close(s.releaseCh)
}

func (s *evalState) emitErr(err error) {
	s.outCh <- shaAndError{err: err}
}

func (s *evalState) emitSha(sha string) {
	s.outCh <- shaAndError{sha: sha}
}

func (s *evalState) ShaExists(sha string) error {
	oid, err := git.NewOid(sha)
	if err != nil {
		return err
	}
	if !s.odb.Exists(oid) {
		return fmt.Errorf("no such object: %q", sha)
	}
	return nil
}

func (s *evalState) LsTree(sha string) (repo.Tree, error) {
	oid, err := git.NewOid(sha)
	if err != nil {
		return repo.Tree{}, err
	}
	obj, err := s.repo.Lookup(oid)
	if err != nil {
		return repo.Tree{}, err
	}
	defer obj.Free()
	tree, err := obj.AsTree()
	if err != nil {
		return repo.Tree{}, err
	}
	defer tree.Free()
	n := tree.EntryCount()
	result := make([]repo.TreeEnt, n)
	for i := uint64(0); i < n; i++ {
		ent := tree.EntryByIndex(i)
		r := repo.TreeEnt{}
		switch ent.Filemode {
		case git.FilemodeTree:
			r.Mode = repo.ModeTree
		case git.FilemodeBlob:
			r.Mode = repo.ModeBlob
		case git.FilemodeBlobExecutable:
			r.Mode = repo.ModeBlobExec
		case git.FilemodeLink:
			r.Mode = repo.ModeSymlink
		case git.FilemodeCommit:
			r.Mode = repo.ModeCommit
		default:
			panic(fmt.Errorf("unknown mode %v %v %v", ent.Filemode, ent.Name, sha))
		}
		r.Name = []byte(ent.Name)
		bs := make([]byte, 20)
		for j, b := range [20]byte(*ent.Id) {
			bs[j] = b
		}
		r.Value = bs
		result[i] = r
	}
	return repo.Tree(result), nil
}

func (s *evalState) LsCommit(sha string) (repo.Commit, error) {
	oid, err := git.NewOid(sha)
	if err != nil {
		return repo.Commit{}, err
	}
	obj, err := s.repo.Lookup(oid)
	if err != nil {
		return repo.Commit{}, err
	}
	defer obj.Free()
	commit, err := obj.AsCommit()
	if err != nil {
		return repo.Commit{}, err
	}
	defer commit.Free()

	result := repo.Commit{}
	result.Tree = commit.TreeId().String()
	for i := uint(0); i < commit.ParentCount(); i++ {
		result.Parents = append(result.Parents, commit.ParentId(i).String())
	}
	return result, nil
}
