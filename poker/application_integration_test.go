package poker_test

import (
	"fmt"
	"github.com/Keisn1/LearnGoWithTestsApp/poker"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := poker.CreateTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := poker.NewFileSystemPlayerStore(database)
	poker.AssertNoError(t, err)

	server, err := poker.NewPlayerServer(store, nil)
	poker.AssertNoError(t, err)

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

	poker.AssertStatus(t, response, http.StatusOK)
	poker.AssertResponseBody(t, response.Body.String(), "3")
}
