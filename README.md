# go-patch-cover [![Github Action](https://img.shields.io/badge/Github%20Action-go--patch--cover--action-brightgreen)](https://github.com/seriousben/go-patch-cover-action)

Report coverage on code that changed.

## Example

```
> go-patch-cover coverage.out patch.diff prevcoverage.out
previous coverage: 90% of statements
new coverage: 91.7% of statements
patch coverage: 96% of changed statements (48/50)
```

## Usage

```
Usage: go-patch-cover [--version] [--help] [flags...] coverage_file diff_file [previous_coverage_file]

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
```

# Generating diff file

`git diff -U0 --no-color origin/main`
