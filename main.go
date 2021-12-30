package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"golang.org/x/tools/cover"
)

func main() {
	patch, err := os.Open("testdata/scenarios/new_file/diff.diff")
	if err != nil {
		log.Fatal(err)
	}

	// files is a slice of *gitdiff.File describing the files changed in the patch
	// preamble is a string of the content of the patch before the first file
	files, _, err := gitdiff.Parse(patch)
	if err != nil {
		log.Fatal(err)
	}

	profiles, err := cover.ParseProfiles("testdata/scenarios/new_file/coverage.out")
	if err != nil {
		log.Fatal(err)
	}

	var (
		numStmt         int
		coverCount      int
		patchNumStmt    int
		patchCoverCount int
	)

	// patch coverage
	for _, p := range profiles {
		for _, f := range files {
			// Using suffix since profiles are prepended with the go module.
			if !strings.HasSuffix(p.FileName, f.NewName) {
				//fmt.Printf("%s != %s\n", p.FileName, f.NewName)
				continue
			}

		blockloop:
			for _, b := range p.Blocks {
				//fmt.Printf("BLOCK %s:%d %d %d %d\n", p.FileName, b.StartLine, b.EndLine, b.NumStmt, b.Count)
				patchNumStmt += b.NumStmt
				for _, t := range f.TextFragments {
					for i, line := range t.Lines {
						if line.Op != gitdiff.OpAdd {
							continue
						}
						lineNum := int(t.NewPosition) + i
						//lineString := strings.ReplaceAll(line.Line, "\n", "")
						// fmt.Printf("DIFF %s:%d %s\n", f.NewName, lineNum, lineString)

						if b.StartLine <= lineNum && lineNum <= b.EndLine {
							//		fmt.Printf("COVER %s:%d %d %d - %s\n", p.FileName, lineNum, b.NumStmt, b.Count, lineString)
							patchCoverCount += b.NumStmt * b.Count
							continue blockloop
						}
					}
				}
			}
		}
	}

	// global coverage
	for _, p := range profiles {
		for _, b := range p.Blocks {
			numStmt += b.NumStmt
			coverCount += b.NumStmt * b.Count
		}
	}

	// TODO: Previous coverage

	if numStmt != 0 {
		fmt.Printf("coverage: %.1f%% of statements\n", float64(coverCount)/float64(numStmt)*100)
	} else {
		fmt.Printf("coverage: %d%% of statements\n", 0)
	}
	if patchNumStmt != 0 {
		fmt.Printf("patch coverage: %.1f%% of changed statements\n", float64(patchCoverCount)/float64(patchNumStmt)*100)
	} else {
		fmt.Printf("patch coverage: %d%% of changed statements\n", 0)
	}
}
