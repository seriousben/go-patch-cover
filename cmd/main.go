package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	patchcover "github.com/seriousben/go-patch-cover"
)

var (
	version string = "dev"
)

func main() {
	c := newCoverCommand(version)
	if err := c.Run(os.Args[1:]); err != nil {
		log.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

type CoverCommand struct {
	fs *flag.FlagSet

	VersionFlag  bool
	HelpFlag     bool
	OutputFlag   string
	TemplateFlag string

	version string
}

func newCoverCommand(version string) *CoverCommand {
	c := &CoverCommand{
		fs:      flag.NewFlagSet("", flag.ContinueOnError),
		version: version,
	}

	c.fs.Usage = c.Usage

	c.fs.BoolVar(&c.VersionFlag, "version", false, "print go-patch-cover version")
	c.fs.BoolVar(&c.HelpFlag, "help", false, "print go-patch-cover help")
	c.fs.StringVar(&c.OutputFlag, "o", "template", "coverage output format: json, template")
	c.fs.StringVar(&c.TemplateFlag, "tmpl", "", "go template string override")
	return c
}

func (c *CoverCommand) Usage() {
	// TODO: Link to template variable struct on github.
	usage := `Usage: go-patch-cover [--version] [--help] [flags...] coverage_file diff_file [previous_coverage_file] 

Arguments:
	coverage_file
		go coverage file for the code after patch was applied.
		Can be generated with any cover mode.
		Example generation:
			go test -coverprofile=coverage.out -covermode=count ./...

	diff_file
		unified diff file of the patch to compute coverage for.
		Example generation:
			git diff -U0 --no-color origin/${GITHUB_BASE_REF} > patch.diff

	previous_coverage_file [OPTIONAL]
		go coverage file for the code before the patch was applied.
		When not provided, previous coverage information will not be displayed.

Flags:
	--version
		display go-patch-cover version.

	--help
		display this help message.

	-o string
		output format: json, template; default: template.

	-tmpl string
		go template string to override default template.

Examples:

	Display total and patch coverage percentages to stdout:
		go-patch-cover coverage.out patch.diff

	Display previous, total and patch coverage percentages to stdout:
		go-patch-cover coverage.out patch.diff prevcoverage.out

	Display previous, total and patch coverage percentages as JSON to stdout:
		go-patch-cover -o json coverage.out patch.diff prevcoverage.out

	Display patch coverage percentage to stdout by providing a custom template:
		go-patch-cover -tmpl "{{ .PatchCoverage }}" coverage.out patch.diff
`

	_, _ = fmt.Fprint(os.Stdout, usage)
}

func (c *CoverCommand) Run(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return fmt.Errorf("flag parse error: %v", err)
	}

	if c.HelpFlag {
		c.fs.Usage()
		return nil
	}

	if c.VersionFlag {
		fmt.Println(c.version)
		return nil
	}

	covFile := c.fs.Arg(0)
	if covFile == "" {
		return fmt.Errorf("missing coverage file argument")
	}
	diffFile := c.fs.Arg(1)
	if diffFile == "" {
		return fmt.Errorf("missing diff file argument")
	}
	prevCovFile := c.fs.Arg(2)

	coverage, err := patchcover.ProcessFiles(covFile, diffFile, prevCovFile)
	if err != nil {
		return fmt.Errorf("processing error: %w", err)
	}

	if c.OutputFlag == "json" {
		enc := json.NewEncoder(os.Stdout)
		err := enc.Encode(coverage)
		if err != nil {
			return fmt.Errorf("json output error: %w", err)
		}
		return nil
	}

	err = patchcover.RenderTemplateOutput(coverage, c.TemplateFlag, os.Stdout)
	if err != nil {
		return fmt.Errorf("json output error: %w", err)
	}

	return nil
}
