package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ChimeraCoder/gitgo"
)

func main() {
	expr := BlobExpr{sha: gitgo.SHA("03039c18b52661b04ecc5d381140222914525a7f")}
	dir, err := os.Open(".")
	if err != nil {
		log.Fatal(err)
	}
	r := &gitgo.Repository{Basedir: *dir}

	ch := Eval(expr, r)

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
