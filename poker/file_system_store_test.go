package poker_test

import (
	"github.com/Keisn1/LearnGoWithTestsApp/poker"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, "")
		defer cleanDatabase()
		_, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)
	})

	t.Run("League from a reader", func(t *testing.T) {

		database, cleanDatabase := poker.CreateTempFile(t, `[
{"Name": "Gilly","Wins": 10},
{"Name": "Pepper", "Wins": 20}]`)

		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		got := store.GetLeagueTable()

		want := poker.League{
			{"Pepper", 20},
			{"Gilly", 10},
		}

		poker.AssertLeague(t, got, want)

		// read again
		got = store.GetLeagueTable()
		poker.AssertLeague(t, got, want)
	})
}

func TestGetPlayerScore(t *testing.T) {
	t.Run("League from a reader", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
{"Name": "Pepper", "Wins": 20},
{"Name": "Gilly","Wins": 10}]`)

		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		want := 20

		got, _ := store.GetPlayerScore("Pepper")
		poker.AssertEqualScores(t, got, want)

		// read again
		got, _ = store.GetPlayerScore("Pepper")
		poker.AssertEqualScores(t, got, want)
	})
}

func TestRecordWin(t *testing.T) {
	t.Run("store wins for player", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
{"Name": "Pepper", "Wins": 20},
{"Name": "Gilly","Wins": 10}]`)

		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		store.RecordWin("Pepper")
		want := 21
		got, _ := store.GetPlayerScore("Pepper")
		if got != want {
			t.Errorf(`got = %v; want %v`, got, want)
		}
	})
	t.Run("store wins for player", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
{"Name": "Pepper", "Wins": 20},
{"Name": "Gilly","Wins": 10}]`)

		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		store.RecordWin("Johnny")
		want := 1
		got, _ := store.GetPlayerScore("Johnny")
		if got != want {
			t.Errorf(`got = %v; want %v`, got, want)
		}
	})
}
