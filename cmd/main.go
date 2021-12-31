package main

import (
	"fmt"
	"log"

	patchcover "github.com/seriousben/go-patch-cover"
)

func main() {
	coverage, err := patchcover.ProcessFiles("diff.diff", "coverage.out")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("coverage: %.1f%% of statements\n", coverage.Coverage)
	fmt.Printf("patch coverage: %.1f%% of changed statements\n", coverage.PatchCoverage)
}
