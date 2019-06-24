package tests

import (
	"path/filepath"
	"testing"
)

var mappingJsonPath = filepath.Join("", "mapping.json")

func TestMAPPING(t *testing.T) {
	run(t, mappingJsonPath)
}
