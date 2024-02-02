package poker_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"io"

	"github.com/Keisn1/LearnGoWithTestsApp/poker"
)

type SpyGame struct {
	StartCalled     bool
	StartCalledWith int
	BlindAlert      []byte

	FinishCalled     bool
	FinishCalledWith string
}

func (g *SpyGame) Start(nbrOfPlayers int, out io.Writer) {
	g.StartCalled = true
	g.StartCalledWith = nbrOfPlayers
	out.Write(g.BlindAlert)
}

func (g *SpyGame) Finish(winner string) {
	g.FinishCalledWith = winner
}

func TestCLI(t *testing.T) {
	t.Run("PlayPoker", func(t *testing.T) {
		out := &bytes.Buffer{}
		spyGame := SpyGame{}

		in := userSends("8", "Cleo wins")
		cli := poker.NewCLI(in, out, &spyGame)

		cli.PlayPoker()

		// checking prompt
		wantMessages := []string{poker.Welcome, poker.PlayerPrompt, poker.UserInfo}
		AssertMessageSentToUser(t, out, wantMessages...)

		// checking calls to Start and Finish
		assertStartCalledWith(t, &spyGame, 8)
		assertFinishCalledWith(t, &spyGame, "Cleo")
	})

	t.Run("PlayPoker", func(t *testing.T) {
		out := &bytes.Buffer{}
		spyGame := SpyGame{}

		in := userSends("7", "Chris wins")
		cli := poker.NewCLI(in, out, &spyGame)

		cli.PlayPoker()

		AssertMessageSentToUser(t, out, poker.Welcome, poker.PlayerPrompt, poker.UserInfo)
		assertStartCalledWith(t, &spyGame, 7)
		assertFinishCalledWith(t, &spyGame, "Chris")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		out := &bytes.Buffer{}
		in := userSends("Pies")
		spyGame := SpyGame{}

		cli := poker.NewCLI(in, out, &spyGame)
		cli.PlayPoker()

		assertGameNotStarted(t, spyGame)
		AssertMessageSentToUser(t, out, poker.Welcome, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		out := &bytes.Buffer{}
		in := userSends("Lloyd is a killer")
		spyGame := SpyGame{}

		cli := poker.NewCLI(in, out, &spyGame)
		cli.PlayPoker()

		assertGameNotStarted(t, spyGame)
		AssertMessageSentToUser(t, out, poker.Welcome, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})
}

func AssertMessageSentToUser(t *testing.T, out *bytes.Buffer, messages ...string) {
	t.Helper()
	wantPrompt := strings.Join(messages, "\n")
	gotPrompt := out.String()
	if gotPrompt != wantPrompt {
		t.Errorf(`got = %v; want %v`, gotPrompt, wantPrompt)
	}
}

func assertStartCalledWith(t *testing.T, game *SpyGame, want int) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.StartCalledWith == want
	})

	if !passed {
		t.Errorf("expected start called with %d but got %d", want, game.StartCalledWith)
	}
}

func assertFinishCalledWith(t *testing.T, game *SpyGame, winner string) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.FinishCalledWith == winner
	})

	if !passed {
		t.Errorf("expected finish called with %q but got %q", winner, game.FinishCalledWith)
	}
}

func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}

func assertGameNotStarted(t *testing.T, g SpyGame) {
	t.Helper()
	if g.StartCalled {
		t.Errorf("game should not have started")
	}
}

func userSends(msgs ...string) *strings.Reader {
	msg := strings.Join(msgs, "\n") + "\n"
	in := strings.NewReader(msg)
	return in
}
