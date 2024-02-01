package poker_test

import (
	"fmt"
	"github.com/Keisn1/LearnGoWithTestsApp/poker"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLeagueTable(t *testing.T) {
	store := poker.StubPlayerStore{}
	server := poker.NewPlayerServer(&store)

	t.Run("Returns correct status codes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		poker.AssertStatusCode(t, response.Code, http.StatusOK)

		got := response.Header()["Content-Type"]
		if len(got) != 1 {
			t.Errorf("Got %d entries in response.Header()['content-type']); want %d", len(got), 1)
		}
		want := "application/json"
		if got[0] != want {
			t.Errorf(`response.Header()["content-type"] = "%v"; want "%v"`, got[0], want)
		}
	})
	t.Run("Returns empty list of players", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		poker.AssertStatusCode(t, response.Code, http.StatusOK)

		got := poker.GetLeagueFromResponse(t, response.Body)
		want := poker.League{}
		poker.AssertStatusCode(t, response.Code, http.StatusOK)
		if len(want) != len(got) {
			t.Errorf("got has length %d and want has length %d", len(got), len(want))
		}
	})

	t.Run("Returns json with list of players", func(t *testing.T) {
		store := poker.NewStubPlayerStore(
			map[string]int{
				"Pepper": 20,
				"Billy":  10,
			},
		)

		server := poker.NewPlayerServer(&store)

		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := []poker.Player{
			{"Pepper", 20}, {"Billy", 10},
		}

		got := poker.GetLeagueFromResponse(t, response.Body)
		poker.AssertStatusCode(t, response.Code, http.StatusOK)
		poker.AssertLeague(t, got, want)
	})
}

func TestGetPlayers(t *testing.T) {
	testCases := []struct {
		player string
		want   string
	}{
		{player: "Pepper", want: "20"},
		{player: "Billy", want: "10"},
	}

	store := poker.NewStubPlayerStore(map[string]int{
		"Pepper": 20,
		"Billy":  10,
	})

	server := poker.NewPlayerServer(&store)

	for _, tc := range testCases {
		testName := fmt.Sprintf("Get %s score", tc.player)
		t.Run(testName, func(t *testing.T) {
			request := poker.NewGetScoreRequest(tc.player)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)
			poker.AssertStatusCode(t, response.Code, http.StatusOK)
			poker.AssertResponseBody(t, response.Body.String(), tc.want)
		})
	}

	t.Run("returns 404 on missing player", func(t *testing.T) {
		request := poker.NewGetScoreRequest("missing")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		if response.Code != http.StatusNotFound {
			t.Errorf("Missing player Did not response with 404")
		}
	})
}

func TestStoreWins(t *testing.T) {
	store := poker.NewStubPlayerStore(map[string]int{})
	server := poker.NewPlayerServer(&store)

	t.Run("It returns accepted on Post", func(t *testing.T) {
		request := poker.NewPostWinRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		poker.AssertStatusCode(t, response.Code, http.StatusAccepted)
		poker.AssertPlayerWin(t, &store, "Pepper")
	})
}
