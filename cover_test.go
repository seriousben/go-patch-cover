package patchcover

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

func TestProcessFiles(t *testing.T) {
	scenarioDir := "./testdata/scenarios"
	fis, err := os.ReadDir(scenarioDir)
	assert.NilError(t, err)

	assert.Assert(t, len(fis) > 0) // Some scenarios exist.

	for _, fi := range fis {
		if !fi.IsDir() {
			t.FailNow() // only directories allowed.
		}

		t.Run(fi.Name(), func(t *testing.T) {
			cov, err := ProcessFiles(path.Join(scenarioDir, fi.Name(), "diff.diff"), path.Join(scenarioDir, fi.Name(), "coverage.out"))
			assert.NilError(t, err)

			covJSON, err := json.MarshalIndent(cov, "", "  ")
			assert.NilError(t, err)

			golden.Assert(t, string(covJSON), path.Join("scenarios", fi.Name(), "golden.json"))
		})
	}
}
