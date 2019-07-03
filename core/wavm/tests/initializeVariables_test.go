package tests

import (
	"path/filepath"
	"testing"
)

var initVarJsonPath = filepath.Join("", "initializeVariables.json")

func TestInitializeVariables(t *testing.T) {
	run(t, initVarJsonPath)
}
