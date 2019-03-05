package tests

import (
	"path/filepath"
	"testing"
)

var mutableJsonPath = filepath.Join("", "mutable.json")

func TestMutable(t *testing.T) {
	run(t, mutableJsonPath)
}
