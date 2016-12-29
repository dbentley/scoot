package main

import ()

type Expr interface{}

type BlobExpr struct {
}

func (e BlobExpr) Covers(o BlobExpr) bool {
	return true
}

type TreeExpr struct {
	trees bool
	blobs bool
}

func (e TreeExpr) Covers(o TreeExpr) bool {
	return (e.trees || !o.trees) && (e.blobs || !o.blobs)
}

type CommitExpr struct {
	tree    TreeExpr
	history int
}

func (e CommitExpr) Covers(o CommitExpr) bool {
	if o.history == -1 && e.history != -1 {
		return false
	}
	if e.history < o.history {
		return false
	}

	return e.tree.Covers(o.tree)
}
