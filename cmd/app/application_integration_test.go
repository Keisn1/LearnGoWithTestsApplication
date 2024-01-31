package main

import (
	"app"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
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

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, "")
	defer cleanDatabase()
	store := app.FileSystemPlayerStore{DB: database}

	server := app.NewPlayerServer(&store)
	player := "Pepper"

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil)
	server.ServeHTTP(httptest.NewRecorder(), request)
	request, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil)
	server.ServeHTTP(httptest.NewRecorder(), request)
	request, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil)
	server.ServeHTTP(httptest.NewRecorder(), request)

	response := httptest.NewRecorder()
	request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)
	server.ServeHTTP(response, request)

	assertStatusCode(t, response.Code, http.StatusOK)
	assertResponseBody(t, response.Body.String(), "3")
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got = \"%v\"; want \"%v\"", got, want)
	}

}

func assertStatusCode(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got = \"%v\"; want \"%v\"", got, want)
	}
}
