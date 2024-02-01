package poker_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Keisn1/LearnGoWithTestsApp/poker"
)

const PlayerPrompt = "Please enter the number of players: "

func TestCLI(t *testing.T) {
	t.Run("Get the number of players from the user", func(t *testing.T) {
		stdout := bytes.Buffer{}
		in := strings.NewReader("7\n")
		cli := poker.NewCLI(in, &stdout)

		cli.PromptForPlayers()

		got := stdout.String()
		want := PlayerPrompt
		if got != want {
			t.Errorf(`got = %v; want %v`, got, want)
		}
	})
}
