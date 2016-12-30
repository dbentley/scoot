package main

import (
	"fmt"
	"log"
	"os"

	"github.com/libgit2/git2go"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Must supply an Eval spec")
	}

	var specs []EvalSpec
	for _, arg := range os.Args[1:] {
		spec, err := Parse(arg)
		if err != nil {
			log.Fatal(err)
		}
		specs = append(specs, spec)
	}

	r, err := git.OpenRepository("/Users/dbentley/workspace/source")
	if err != nil {
		log.Fatal(err)
	}

	ch := Eval(r, specs)

	for sAe := range ch {
		sha, err := sAe.sha, sAe.err
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(sha)
	}
}

const val = `
blob("6c6441e209c6318a562ae3e61e0f4d7bd4bcda05") {}
`
