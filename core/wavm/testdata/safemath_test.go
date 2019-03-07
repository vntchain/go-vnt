package tests

import (
	"path/filepath"
	"testing"
)

var safemathJsonPath = filepath.Join("", "safemath.json")

func TestSafeMath(t *testing.T) {
	run(t, safemathJsonPath)
}
