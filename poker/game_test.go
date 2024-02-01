package poker_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Keisn1/LearnGoWithTestsApp/poker"
)

var (
	dummyPlayerStore = poker.StubPlayerStore{}
	dummySpyAlerter  = &poker.SpyBlindAlerter{}
)

func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}

		game := poker.NewGame(&dummyPlayerStore, blindAlerter)
		game.Start(5)

		cases := []poker.ScheduledAlert{
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
		checkSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}
		game := *poker.NewGame(&dummyPlayerStore, blindAlerter)

		game.Start(7)
		cases := []poker.ScheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})

}

func TestGame_Finish(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		store := poker.StubPlayerStore{}
		game := *poker.NewGame(&store, dummySpyAlerter)

		game.Finish("Chris")
		poker.AssertPlayerWin(t, &store, "Chris")
	})

	t.Run("record Cleo win from user input", func(t *testing.T) {
		store := poker.StubPlayerStore{}
		game := *poker.NewGame(&store, dummySpyAlerter)

		game.Finish("Cleo")
		poker.AssertPlayerWin(t, &store, "Cleo")
	})
}

func checkSchedulingCases(t *testing.T, cs []poker.ScheduledAlert, b *poker.SpyBlindAlerter) {
	for i, want := range cs {
		t.Run(fmt.Sprint(want), func(t *testing.T) {

			if len(b.Alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, b.Alerts)
			}

			got := b.Alerts[i]
			poker.AssertScheduledAlert(t, got, want)
		})
	}

}
