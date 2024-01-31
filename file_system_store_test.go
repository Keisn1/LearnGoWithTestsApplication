package app

import (
	"fmt"
	"os"
	"testing"
)

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}
	return tmpfile, removeFile
}

func TestFileSystemStore(t *testing.T) {
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()
		_, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
	})

	t.Run("League from a reader", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
{"Name": "Gilly","Wins": 10},
{"Name": "Pepper", "Wins": 20}]`)

		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		fmt.Println(store.league)
		got := store.GetLeagueTable()

		want := []Player{
			{"Pepper", 20},
			{"Gilly", 10},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeagueTable()
		assertLeague(t, got, want)
	})
}

func TestGetPlayerScore(t *testing.T) {
	t.Run("League from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
{"Name": "Pepper", "Wins": 20},
{"Name": "Gilly","Wins": 10}]`)

		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		want := 20

		got, _ := store.GetPlayerScore("Pepper")
		assertEqualScores(t, got, want)

		// read again
		got, _ = store.GetPlayerScore("Pepper")
		assertEqualScores(t, got, want)
	})
}

func TestRecordWin(t *testing.T) {
	t.Run("store wins for player", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
{"Name": "Pepper", "Wins": 20},
{"Name": "Gilly","Wins": 10}]`)

		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin("Pepper")
		want := 21
		got, _ := store.GetPlayerScore("Pepper")
		if got != want {
			t.Errorf(`got = %v; want %v`, got, want)
		}
	})
	t.Run("store wins for player", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
{"Name": "Pepper", "Wins": 20},
{"Name": "Gilly","Wins": 10}]`)

		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin("Johnny")
		want := 1
		got, _ := store.GetPlayerScore("Johnny")
		if got != want {
			t.Errorf(`got = %v; want %v`, got, want)
		}
	})
}

func assertEqualScores(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}
