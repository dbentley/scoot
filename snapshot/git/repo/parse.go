package repo

import (
	"fmt"
)

type Commit struct {
	Tree    string   // sha1 of the Tree for this commit
	Parents []string // sha1 of the parents of this commit
}

func ParseCommit(text string) (Commit, error) {
	var tree string
	var empty Commit
	_, err := fmt.Sscanf(text, "tree %s", &tree)
	if err != nil {
		return empty, err
	}
	if len(tree) != 40 {
		return empty, fmt.Errorf("sha1 in commit %s not 40 chars long: %v; actually %v",
			text, tree, len(tree))
	}
	return Commit{Tree: tree}, nil
}

type Tree []TreeEnt

type TreeEnt struct {
	Mode  []byte
	Name  []byte
	Value []byte // NB: 20-byte sha1, ie. dense, not in hexadecimal
}

var (
	ModeBlobExec = []byte("100755")
	ModeBlob     = []byte("100644")
	ModeSymlink  = []byte("120000")
	ModeTree     = []byte("040000")
	ModeCommit   = []byte("160000")
)

func ParseTree(data []byte) (Tree, error) {
	if len(data) == 0 {
		return Tree([]TreeEnt{}), nil
	}
	idx := 0
	for i, v := range data {
		if v == ' ' {
			idx = i
			break
		}
	}
	mode := data[0:idx]
	if len(mode) == 5 {
		// For mode 040000, git omits the 0 and stores 40000
		// Not sure why, but okay.
		mode = []byte("0" + string(mode))
	}
	if len(mode) != 6 {
		return nil, fmt.Errorf("No valid mode found: %v", data)
	}
	data = data[idx+1:] // +1 to skip the space

	idx = 0
	for i, v := range data {
		if v == 0 {
			idx = i
			break
		}
	}
	if idx == 0 {
		return nil, fmt.Errorf("Found no name: %v", data)
	}
	name := data[0:idx]
	data = data[idx+1:] // +1 to skip the nul byte
	if len(data) < 20 {
		return nil, fmt.Errorf("No SHA found: %v", data)
	}
	sha1 := data[0:20]
	data = data[20:]
	rest, err := ParseTree(data)
	if err != nil {
		return nil, err
	}

	return append(rest, TreeEnt{mode, name, sha1}), nil
}
