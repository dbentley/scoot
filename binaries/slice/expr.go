package main

import ()

type Expr interface{}

type BlobExpr struct {
}

func (e BlobExpr) Covers(o BlobExpr) bool {
	return true
}

type TreeExpr struct {
	Trees bool
	Blobs bool
}

func (e TreeExpr) Covers(o TreeExpr) bool {
	return (e.Trees || !o.Trees) && (e.Blobs || !o.Blobs)
}

type CommitExpr struct {
	Tree    TreeExpr
	History int
}

func (e CommitExpr) Covers(o CommitExpr) bool {
	if o.History == -1 && e.History != -1 {
		return false
	}
	if e.History < o.History {
		return false
	}

	return e.Tree.Covers(o.Tree)
}

type EvalSpec struct {
	SHA    string
	Blob   *BlobExpr
	Tree   *TreeExpr
	Commit *CommitExpr
}
