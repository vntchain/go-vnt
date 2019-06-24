package tests

import (
	"path/filepath"
	"testing"
)

var boundsVarJsonPath = filepath.Join("", "bounds.json")

func TestBounds(t *testing.T) {
	run(t, boundsVarJsonPath)
}
