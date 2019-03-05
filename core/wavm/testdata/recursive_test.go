package tests

import (
	"path/filepath"
	"testing"
)

var recursiveJsonPath = filepath.Join("", "recursive.json")

func TestRecursive(t *testing.T) {
	run(t, recursiveJsonPath)
}
