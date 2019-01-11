package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoClangDump(t *testing.T) {
	for _, fname := range []string{
		"../../testdata/basicparsing.c",
		"/Users/weisaizhang/Documents/go/src/github.com/vntchain/go-vnt/core/wasm/testdata/precompile/contract/main3.cpp",
	} {
		assert.Equal(t, 0, cmd([]string{"-fname", fname}))
	}
}
