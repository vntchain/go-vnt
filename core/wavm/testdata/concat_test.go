package tests

import (
	"path/filepath"
	"testing"
)

var concatJsonPath = filepath.Join("", "concat.json")

func TestConcat(t *testing.T) {
	run(t, concatJsonPath)
}
