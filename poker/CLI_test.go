package poker_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Keisn1/LearnGoWithTestsApp/poker"
)

const PlayerPrompt = "Please enter the number of players: "

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{duration, amount})
}

func TestCLI(t *testing.T) {
	var (
		dummyStdout = bytes.Buffer{}
		// dummyStdIn       = bytes.Buffer{}
		dummyPlayerStore = poker.StubPlayerStore{}
		dummySpyAlerter  = &SpyBlindAlerter{}
	)

	t.Run("Get the number of players from the user", func(t *testing.T) {
		stdout := bytes.Buffer{}
		in := strings.NewReader("7\n")
		blindAlerter := &SpyBlindAlerter{}
		cli := *poker.NewCLI(in, &stdout, &dummyPlayerStore, blindAlerter)

		cli.PlayPoker()

		got := stdout.String()
		want := PlayerPrompt
		if got != want {
			t.Errorf(`got = %v; want %v`, got, want)
		}

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {

				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nChris wins\n")
		playerStore := poker.StubPlayerStore{}
		cli := *poker.NewCLI(in, &dummyStdout, &playerStore, dummySpyAlerter)

		cli.PlayPoker()
		poker.AssertPlayerWin(t, &playerStore, "Chris")
	})

	t.Run("record Cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nCleo wins\n")
		playerStore := poker.StubPlayerStore{}

		cli := *poker.NewCLI(in, &dummyStdout, &playerStore, dummySpyAlerter)

		cli.PlayPoker()
		poker.AssertPlayerWin(t, &playerStore, "Cleo")
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("5\nChris wins\n")
		playerStore := poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewCLI(in, &dummyStdout, &playerStore, blindAlerter)
		cli.PlayPoker()

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {

				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}

	})
}

func assertScheduledAlert(t *testing.T, got, want scheduledAlert) {
	amountGot := got.amount
	if amountGot != want.amount {
		t.Errorf("got amount %d, want %d", amountGot, want.amount)
	}

	gotScheduledTime := got.at
	if gotScheduledTime != want.at {
		t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, want.at)
	}
}
