package app

import (
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("League from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
{"Name": "Pepper", "Wins": 20},
{"Name": "Gilly","Wins": 10}]`)
		store := FileSystemPlayerStore{database}

		want := []Player{
			{"Pepper", 20},
			{"Gilly", 10},
		}

		got := store.GetLeagueTable()
		assertLeague(t, got, want)

		// read again
		got = store.GetLeagueTable()
		assertLeague(t, got, want)
	})
}

func TestGetPlayerScore(t *testing.T) {
	t.Run("League from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
{"Name": "Pepper", "Wins": 20},
{"Name": "Gilly","Wins": 10}]`)
		store := FileSystemPlayerStore{database}

		want := 20

		got, _ := store.GetPlayerScore("Pepper")
		assertEqualScores(t, got, want)

		// read again
		got, _ = store.GetPlayerScore("Pepper")
		assertEqualScores(t, got, want)
	})
}

func TestRecordWin(t *testing.T) {
	t.Run("League from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
{"Name": "Pepper", "Wins": 20}]`)
		store := FileSystemPlayerStore{database}

		want := []Player{
			{
				Name: "Pepper",
				Wins: 21,
			},
		}

		store.RecordWin("Pepper")

		got, _ := NewLeague(store.db)

		assertLeague(t, got, want)
	})
}

func assertEqualScores(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}
