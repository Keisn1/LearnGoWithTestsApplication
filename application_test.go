package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayers() []Player {
	var players []Player
	for p, w := range s.scores {
		players = append(players, Player{p, w})
	}
	return players
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetPlayerScore(name string) (int, StoreError) {
	score, exists := s.scores[name]
	if !exists {
		return -1, PlayerNotFound
	}
	return score, ""
}

func TestLeagueTable(t *testing.T) {
	store := StubPlayerStore{}
	server := NewPlayerServer(&store)

	t.Run("Returns correct status codes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusOK)

		got := response.Header()["Content-Type"]
		if len(got) != 1 {
			t.Errorf("Got %d entries in response.Header()['Content-Type']); want %d", len(got), 1)
		}
		want := "application/json"
		if got[0] != want {
			t.Errorf(`response.Header()["Content-Type"] = "%v"; want "%v"`, got[0], want)
		}
	})
	t.Run("Returns empty list of players", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusOK)

		var got []Player
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to Unmarshal response body into []string")
		}
		want := []Player{}
		assertEqualListPlayers(t, got, want)

	})

	t.Run("Returns json with list of players", func(t *testing.T) {
		store := StubPlayerStore{
			scores: map[string]int{
				"Pepper": 20,
				"Billy":  10,
			},
		}

		server := NewPlayerServer(&store)

		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := []Player{
			{"Pepper", 20}, {"Billy", 10},
		}

		var got []Player
		err := json.Unmarshal(response.Body.Bytes(), &got)
		if err != nil {
			t.Fatalf("Unable to Unmarshal response body into []string")
		}

		if len(got) != len(want) {
			t.Errorf(`len(got) = %v; len(want) %v`, len(got), len(want))
		}

		assertEqualListPlayers(t, got, want)
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

	store := StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Billy":  10,
		},
	}
	server := NewPlayerServer(&store)

	for _, tc := range testCases {
		testName := fmt.Sprintf("Get %s score", tc.player)
		t.Run(testName, func(t *testing.T) {
			request := newGetScoreRequest(tc.player)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)
			assertStatusCode(t, response.Code, http.StatusOK)
			assertResponseBody(t, response.Body.String(), tc.want)
		})
	}

	t.Run("returns 404 on missing player", func(t *testing.T) {
		request := newGetScoreRequest("missing")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		if response.Code != http.StatusNotFound {
			t.Errorf("Missing player Did not response with 404")
		}
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{scores: map[string]int{}}
	server := NewPlayerServer(&store)

	t.Run("It returns accepted on Post", func(t *testing.T) {
		request := newPostWinRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusAccepted)

		want := []string{"Pepper"}
		got := store.winCalls
		if len(got) != 1 {
			t.Errorf("Got %d calls to RecordWin want %d", len(got), 1)
		}

		if !reflect.DeepEqual(store.winCalls, want) {
			t.Errorf("Got = \"%v\"; want \"%v\"", got, want)
		}
	})
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

func newGetScoreRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)
	return request
}

func newPostWinRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil)
	return request
}

func assertEqualListPlayers(t *testing.T, got, want []Player) {
	for _, p1 := range got {
		present := false
		for _, p2 := range want {
			if p1.Name == p2.Name && p1.Wins == p2.Wins {
				present = true
			}
		}
		if !present {
			t.Errorf("%v of got not in want = %v", p1, want)
		}
	}
}
