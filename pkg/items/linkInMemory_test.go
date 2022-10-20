package items

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestStructMemory struct {
	function func(string) string
	argument string
	resStr   string
	resErr   error
}

func TestLinkRepoMemory(t *testing.T) {
	r := NewLinkRepo(nil, "inMemory")

	tests := []TestStructMemory{
		{
			function: r.NewLink,
			argument: "link",
			resStr:   "4f0aa52d65",
		},
		{
			function: r.NewLink,
			argument: "link",
			resStr:   "4f0aa52d65",
		},
		{
			function: r.NewLink,
			argument: "link1",
			resStr:   "862c61fd25",
		},
		{
			function: r.GetOrigLink,
			argument: "4f0aa52d65",
			resStr:   "link",
		},
		{
			function: r.GetOrigLink,
			argument: "asdf",
			resStr:   "",
		},
	}

	for _, test := range tests {
		resStr := test.function(test.argument)

		assert.Equal(t, test.resStr, resStr)
	}
}
