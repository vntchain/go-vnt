package tests

import (
	"path/filepath"
	"testing"
)

var gasJsonPath = filepath.Join("", "gas.json")

func TestGas(t *testing.T) {
	run(t, gasJsonPath)
}
