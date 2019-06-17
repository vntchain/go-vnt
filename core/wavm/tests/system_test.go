package tests

import (
	"path/filepath"
	"testing"
)

var systemJsonPath = filepath.Join("", "system.json")

func TestSystem(t *testing.T) {
	run(t, systemJsonPath)
}
