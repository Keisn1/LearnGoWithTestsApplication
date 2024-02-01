package poker_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Keisn1/LearnGoWithTestsApp/poker"
)

const PlayerPrompt = "Please enter the number of players: "

func TestCLI(t *testing.T) {
	var (
		dummyOut = &bytes.Buffer{}
		dummyIn  = strings.NewReader("")
	)
	t.Run("Prompt for getting the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		cli := poker.NewCLI(dummyIn, stdout)

		cli.PromptForPlayers()

		got := stdout.String()
		want := PlayerPrompt
		if got != want {
			t.Errorf(`got = %v; want %v`, got, want)
		}
	})
	t.Run("Get the number of players from the user", func(t *testing.T) {
		in := strings.NewReader("7\n")
		cli := poker.NewCLI(in, dummyOut)

		got, err := cli.GetNbrOfPlayers()
		poker.AssertNoError(t, err)

		want := 7
		if got != want {
			t.Errorf(`got = %v; want %v`, got, want)
		}
	})
	t.Run("Get a winner of a game", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		cli := poker.NewCLI(in, dummyOut)

		got := cli.GetWinner()

		want := "Chris"
		if got != want {
			t.Errorf(`got = %v; want %v`, got, want)
		}
	})
}
