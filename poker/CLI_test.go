package poker_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Keisn1/LearnGoWithTestsApp/poker"
)

const PlayerPrompt = "Please enter the number of players: "

type SpyGame struct {
	StartCalls  []int
	FinishCalls []string
}

func (g *SpyGame) Start(nbrOfPlayers int) {
	g.StartCalls = append(g.StartCalls, nbrOfPlayers)
}

func (g *SpyGame) Finish(winner string) {
	g.FinishCalls = append(g.FinishCalls, winner)
}

func TestCLI(t *testing.T) {
	t.Run("PlayPoker", func(t *testing.T) {
		out := &bytes.Buffer{}
		in := strings.NewReader(`7
Chris wins
`)
		spyGame := SpyGame{}
		cli := poker.NewCLI(in, out, &spyGame)

		cli.PlayPoker()
		assertOutput(t, out.String())
		assertLengthCalls(t, spyGame, 1)
		assertStartCalledWith(t, spyGame, 7)
		assertFinishCalledWith(t, spyGame, "Chris")
	})

	// t.Run("Prompt for getting the number of players", func(t *testing.T) {
	// 	stdout := &bytes.Buffer{}
	// 	spyGame := SpyGame{}
	// 	cli := poker.NewCLI(dummyIn, stdout, &spyGame)

	// 	cli.PromptForPlayers()

	// 	got := stdout.String()
	// 	want := PlayerPrompt
	// 	if got != want {
	// 		t.Errorf(`got = %v; want %v`, got, want)
	// 	}
	// })

	// t.Run("Get the number of players from the user", func(t *testing.T) {
	// 	in := strings.NewReader("7\n")
	// 	spyGame := SpyGame{}
	// 	cli := poker.NewCLI(in, dummyOut, &spyGame)

	// 	got, err := cli.GetNbrOfPlayers()
	// 	poker.AssertNoError(t, err)

	// 	want := 7
	// 	if got != want {
	// 		t.Errorf(`got = %v; want %v`, got, want)
	// 	}
	// })

	// t.Run("Get a winner of a game", func(t *testing.T) {
	// 	in := strings.NewReader("Chris wins\n")
	// 	spyGame := SpyGame{}
	// 	cli := poker.NewCLI(in, dummyOut, &spyGame)

	// 	got := cli.GetWinner()

	// 	want := "Chris"
	// 	if got != want {
	// 		t.Errorf(`got = %v; want %v`, got, want)
	// 	}
	// })
}

func assertOutput(t *testing.T, got string) {
	t.Helper()
	want := strings.Join([]string{
		poker.Welcome,
		poker.PlayerPrompt,
		poker.UserInfo}, "\n",
	)
	if got != want {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func assertStartCalledWith(t *testing.T, s SpyGame, want int) {
	t.Helper()
	got := s.StartCalls[0]
	if got != want {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func assertFinishCalledWith(t *testing.T, s SpyGame, want string) {
	t.Helper()
	got := s.FinishCalls[0]
	if got != want {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func assertLengthCalls(t *testing.T, s SpyGame, want int) {
	t.Helper()
	lenStart := len(s.StartCalls)
	lenFinish := len(s.FinishCalls)
	if lenStart != 1 {
		t.Fatalf("Length of StartCalls = %d, not %d", lenStart, want)
	}
	if lenFinish != 1 {
		t.Fatalf("Length of FinishCalls = %d, not %d", lenFinish, want)
	}
}
