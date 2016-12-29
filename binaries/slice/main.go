package main

import (
	"fmt"
	"log"

	"github.com/scootdev/scoot/snapshot/git/repo"
)

func main() {
	expr := TreeExpr{blobs: false, trees: true}

	r, err := repo.NewRepository(".")
	if err != nil {
		log.Fatal(err)
	}

	ch := Eval(expr, r, "f3756251c62f28a89e56bc3cc0df4a4002f693d1")

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
