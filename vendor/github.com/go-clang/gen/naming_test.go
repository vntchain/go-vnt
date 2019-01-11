package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpperFirstCharacter(t *testing.T) {
	for _, tc := range []struct {
		Data   string
		Expect string
	}{
		{
			Data:   "a",
			Expect: "A",
		},
		{
			Data:   "A",
			Expect: "A",
		},
		{
			Data:   "abc",
			Expect: "Abc",
		},
		{
			Data:   "Abc",
			Expect: "Abc",
		},
		{
			Data:   "ABC",
			Expect: "ABC",
		},
	} {
		assert.Equal(t, tc.Expect, UpperFirstCharacter(tc.Data))
	}
}
