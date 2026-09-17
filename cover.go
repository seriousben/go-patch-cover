package patchcover

import (
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"golang.org/x/tools/cover"
)

func ProcessFiles(coverageFile, diffFile, prevCovFile string) (CoverageData, error) {
	patch, err := os.Open(diffFile)
	if err != nil {
		return CoverageData{}, err
	}

	files, _, err := gitdiff.Parse(patch)
	if err != nil {
		return CoverageData{}, err
	}

	profiles, err := cover.ParseProfiles(coverageFile)
	if err != nil {
		return CoverageData{}, err
	}

	var prevProfiles []*cover.Profile
	if prevCovFile != "" {
		prevProfiles, err = cover.ParseProfiles(prevCovFile)
		if err != nil {
			return CoverageData{}, err
		}
	}

	d, err := computeCoverage(files, profiles, prevProfiles)
	if err != nil {
		return CoverageData{}, err
	}

	d.HasPrevCoverage = prevCovFile != ""
	return d, nil
}

type CoverageData struct {
	NumStmt         int     `json:"num_stmt"`
	CoverCount      int     `json:"cover_count"`
	Coverage        float64 `json:"coverage"`
	PatchNumStmt    int     `json:"patch_num_stmt"`
	PatchCoverCount int     `json:"patch_cover_count"`
	PatchCoverage   float64 `json:"patch_coverage"`
	HasPrevCoverage bool    `json:"has_prev_coverage"`
	PrevNumStmt     int     `json:"prev_num_stmt"`
	PrevCoverCount  int     `json:"prev_cover_count"`
	PrevCoverage    float64 `json:"prev_coverage"`
}

func RenderTemplateOutput(data CoverageData, tmplOverride string, out io.Writer) error {
	const defaultTmpl = `
{{- if .HasPrevCoverage -}}
	previous coverage: {{printf "%.1f" .PrevCoverage}}% of statements
{{ else -}}
	previous coverage: unknown
{{ end -}}
new coverage: {{printf "%.1f" .Coverage}}% of statements
patch coverage: {{printf "%.1f" .PatchCoverage}}% of changed statements ({{ .PatchCoverCount }}/{{ .PatchNumStmt }})
`
	tmpl := defaultTmpl
	if tmplOverride != "" {
		tmpl = tmplOverride
	}
	t, err := template.New("cover_template").Parse(tmpl)
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(out, "cover_template", data)
}

func computeCoverage(diffFiles []*gitdiff.File, coverProfiles []*cover.Profile, prevCoverProfiles []*cover.Profile) (CoverageData, error) {
	var data CoverageData
	// patch coverage
	for _, p := range coverProfiles {
		for _, f := range diffFiles {
			// Using suffix since profiles are prepended with the go module.
			if !strings.HasSuffix(p.FileName, f.NewName) {
				//fmt.Printf("%s != %s\n", p.FileName, f.NewName)
				continue
			}

		blockloop:
			for _, b := range p.Blocks {
				//fmt.Printf("BLOCK %s:%d %d %d %d\n", p.FileName, b.StartLine, b.EndLine, b.NumStmt, b.Count)
				for _, t := range f.TextFragments {
					// lineNum tracks the new-file line number. It only
					// advances for lines that occupy a line in the new
					// file (context or added); removed lines don't, so
					// they must not shift the count. Indexing by the raw
					// slice position instead double-counts removed lines
					// and misattributes later added lines to the wrong
					// new-file line whenever a hunk removes more lines
					// than it adds (e.g. deleting a stale comment next
					// to the changed statement).
					lineNum := int(t.NewPosition)
					for _, line := range t.Lines {
						if line.Op == gitdiff.OpDelete {
							continue
						}
						if line.Op == gitdiff.OpAdd && b.StartLine <= lineNum && lineNum <= b.EndLine {
							//lineString := strings.ReplaceAll(line.Line, "\n", "")
							// fmt.Printf("DIFF %s:%d %s\n", f.NewName, lineNum, lineString)
							data.PatchNumStmt += b.NumStmt
							//		fmt.Printf("COVER %s:%d %d %d - %s\n", p.FileName, lineNum, b.NumStmt, b.Count, lineString)
							if b.Count > 0 {
								data.PatchCoverCount += b.NumStmt
							}
							continue blockloop
						}
						lineNum++
					}
				}
			}
		}
	}

	// total coverage
	for _, p := range coverProfiles {
		for _, b := range p.Blocks {
			data.NumStmt += b.NumStmt
			if b.Count > 0 {
				data.CoverCount += b.NumStmt
			}
		}
	}

	// prev total coverage
	for _, p := range prevCoverProfiles {
		for _, b := range p.Blocks {
			data.PrevNumStmt += b.NumStmt
			if b.Count > 0 {
				data.PrevCoverCount += b.NumStmt
			}
		}
	}

	// TODO: Previous coverage

	if data.NumStmt != 0 {
		data.Coverage = float64(data.CoverCount) / float64(data.NumStmt) * 100
	}
	if data.PatchNumStmt != 0 {
		data.PatchCoverage = float64(data.PatchCoverCount) / float64(data.PatchNumStmt) * 100
	}
	if data.PrevNumStmt != 0 {
		data.PrevCoverage = float64(data.PrevCoverCount) / float64(data.PrevNumStmt) * 100
	}

	return data, nil
}
