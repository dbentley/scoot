package main

import (
	"log"

	"github.com/ChimeraCoder/gitgo"
)

type shaAndError struct {
	sha gitgo.SHA
	err error
}

// type SHAScanner struct {
// 	ch shaAndError

// 	sha gitgo.SHA
// 	err error
// }

// func (s *SHAScanner) Scan() bool {
// 	if s.ch == nil {
// 		return false
// 	}

// 	r, ok := <-s.ch
// 	if !ok {
// 		s.ch = nil
// 		return false
// 	}

// 	s.sha = r.sha
// 	if s.err == nil {
// 		s.err = r.err
// 	}
// }

// func (s *SHAScanner) SHA() gitgo.SHA {
// 	return s.sha
// }

// func (s *SHAScanner) Err() error {
// 	return s.err
// }

func Eval(e Expr, r *gitgo.Repository) chan shaAndError {
	ch := make(chan shaAndError)
	go eval(e, r, ch)
	return ch
}

func eval(e Expr, r *gitgo.Repository, ch chan shaAndError) {
	switch e := e.(type) {
	case BlobExpr:
		evalBlob(e, r, ch)
	}
}

func evalBlob(e BlobExpr, r *gitgo.Repository, ch chan shaAndError) {
	obj, err := r.Object(e.sha)
	if err != nil {
		ch <- errOut(err)
		return
	}
	log.Println("BLOB", obj)
	ch <- shaOut(e.sha)
}

func errOut(err error) shaAndError {
	return shaAndError{err: err}
}

func shaOut(sha gitgo.SHA) shaAndError {
	return shaAndError{sha: sha}
}
