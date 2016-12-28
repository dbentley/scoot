package main

import (
	"github.com/ChimeraCoder/gitgo"
)

type Expr interface{}

type BlobExpr struct {
	sha gitgo.SHA
}
