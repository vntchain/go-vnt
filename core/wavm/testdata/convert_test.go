package tests

import (
	"path/filepath"
	"testing"
)

var convertVarJsonPath = filepath.Join("", "convert.json")

func TestConvert(t *testing.T) {
	run(t, convertVarJsonPath)
}
