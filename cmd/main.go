package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/mitchellh/cli"

	patchcover "github.com/seriousben/go-patch-cover"
)

var (
	version string = "dev"
)

func main() {
	c := &cli.CLI{
		Name: "go-patch-cover",
		// TODO figure out version aligment with release
		Version:      version,
		HelpFunc:     cli.BasicHelpFunc("go-patch-cover"),
		Autocomplete: true,

		HelpWriter:  os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"cover": newCoverCommand,
		"":      newCoverCommand,
	}

	c.HiddenCommands = []string{""}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println("[ERROR] ", err)
	}

	os.Exit(exitStatus)
}

type CoverCommand struct {
	fs *flag.FlagSet

	OutputFlag   string
	TemplateFlag string
}

func newCoverCommand() (cli.Command, error) {
	gc := &CoverCommand{
		fs: flag.NewFlagSet("cover", flag.ContinueOnError),
	}

	gc.fs.StringVar(&gc.OutputFlag, "o", "template", "coverage output format: json, template")
	gc.fs.StringVar(&gc.TemplateFlag, "tmpl", "", "go template string override")
	return gc, nil
}

func (g *CoverCommand) Help() string {
	// TODO: Link to template variable struct on github.
	return `Usage: go-patch-cover cover [flags...] coverage diff [previous_coverage] 

Arguments:
	coverage
		go coverage file for the code after patch was applied.
		Can be generated with any cover mode.
		Example generation:
			go test -coverprofile=coverage.out -covermode=count ./...

	diff
		unified diff file of the patch to compute coverage for.
		Example generation:
			git diff -U0 --no-color origin/${GITHUB_BASE_REF} > patch.diff

	previous_coverage [OPTIONAL]
		go coverage file for the code before the patch was applied.
		When not provided, previous coverage information will not be displayed.

Flags:
	-o string
		output format: json, template; default: template.

	-tmpl string
		go template string to override default template.

Examples:

	Display total and patch coverage percentages to stdout:
		go-patch-cover cover coverage.out patch.diff

	Display previous, total and patch coverage percentages to stdout:
		go-patch-cover cover coverage.out patch.diff prevcoverage.out

	Display previous, total and patch coverage percentages as JSON to stdout:
		go-patch-cover cover -o json coverage.out patch.diff prevcoverage.out

	Display patch coverage percentage to stdout by providing a custom template:
		go-patch-cover cover -tmpl "{{ .PatchCoverage }}" coverage.out patch.diff
`
}

func (g *CoverCommand) Synopsis() string {
	return "Display patch coverage percentages"
}

func (g *CoverCommand) Run(args []string) int {
	if err := g.fs.Parse(args); err != nil {
		log.Printf("[ERROR] %v\n", err)
		return 1
	}

	covFile := g.fs.Arg(0)
	if covFile == "" {
		log.Printf("[ERROR] missing coverage file argument\n")
		return 1
	}
	diffFile := g.fs.Arg(1)
	if diffFile == "" {
		log.Printf("[ERROR] missing diff file argument\n")
		return 1
	}
	prevCovFile := g.fs.Arg(2)

	coverage, err := patchcover.ProcessFiles(covFile, diffFile, prevCovFile)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		return 1
	}

	if g.OutputFlag == "json" {
		enc := json.NewEncoder(os.Stdout)
		err := enc.Encode(coverage)
		if err != nil {
			log.Printf("[ERROR] %v\n", err)
			return 1
		}

		return 0
	}

	err = patchcover.RenderTemplateOutput(coverage, g.TemplateFlag, os.Stdout)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		return 1
	}

	return 0
}

var _ cli.Command = (*CoverCommand)(nil)
