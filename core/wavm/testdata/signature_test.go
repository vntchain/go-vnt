package tests

import (
	"path/filepath"
	"testing"
)

var signatureJsonPath = filepath.Join("", "signature.json")

func TestSignature(t *testing.T) {
	run(t, signatureJsonPath)
}
