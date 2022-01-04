package main

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_newCoverCommand(t *testing.T) {
	c := newCoverCommand("1.0.0")
	assert.Assert(t, c != nil)
}
