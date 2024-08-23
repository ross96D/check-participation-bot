package server_test

import (
	"testing"

	"github.com/ross96D/cw_participation_bot/server"
	"github.com/stretchr/testify/require"
)

func TestGetCommand(t *testing.T) {
	type Test struct {
		text string
		cmd  string
		ok   bool
	}

	tests := []Test{
		{
			text: "/command",
			cmd:  "/command",
			ok:   true,
		},
		{
			text: "/command with some spaces",
			cmd:  "/command",
			ok:   true,
		},
		{
			text: "/co23cand",
			cmd:  "/co23cand",
			ok:   true,
		},
		{
			text: "",
			cmd:  "",
			ok:   false,
		},
		{
			text: "no a command",
			cmd:  "",
			ok:   false,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			cmd, ok := server.GetCommand(test.text)
			require.Equal(t, ok, test.ok)
			require.Equal(t, cmd, test.cmd)
		})
	}
}
