package main

import (
	"bytes"
	"fmt"
	"text/scanner"
	// "github.com/ChimeraCoder/gitgo"
)

type Spec interface{}

type BlobSpec struct {
}

type TreeSpec struct {
	subtrees bool
	subblobs bool
}

// type CommitSpec struct {
// 	tree      TreeSpec
// 	ancestors AncestorsSpec
// }

// type AncestorsSpec struct {
// 	commit CommitSpec
// }

// type AncestorsVal struct {
// 	num int
// 	AncestorsSpec
// }

// Spec can be recursive
// Val isn't
func Parse(specText string) (Spec, error) {
	rdr := bytes.NewBufferString(specText)
	var s scanner.Scanner
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.SkipComments
	s.Init(rdr)

	return nil, fmt.Errorf("not yet implemented")
}
