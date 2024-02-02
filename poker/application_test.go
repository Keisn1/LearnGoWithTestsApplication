package poker_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"time"

	"github.com/Keisn1/LearnGoWithTestsApp/poker"
	"github.com/gorilla/websocket"
)

var (
	dummyGame = &SpyGame{}
)

func TestLeagueTable(t *testing.T) {
	store := poker.StubPlayerStore{}
	server, err := poker.NewPlayerServer(&store, nil)
	poker.AssertNoError(t, err)
	t.Run("Returns correct status codes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		poker.AssertStatus(t, response, http.StatusOK)

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
		poker.AssertStatus(t, response, http.StatusOK)

		got := poker.GetLeagueFromResponse(t, response.Body)
		want := poker.League{}
		poker.AssertStatus(t, response, http.StatusOK)
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

		server, err := poker.NewPlayerServer(&store, nil)
		poker.AssertNoError(t, err)
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := []poker.Player{
			{"Pepper", 20}, {"Billy", 10},
		}

		got := poker.GetLeagueFromResponse(t, response.Body)
		poker.AssertStatus(t, response, http.StatusOK)
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

	server, err := poker.NewPlayerServer(&store, nil)
	poker.AssertNoError(t, err)

	for _, tc := range testCases {
		testName := fmt.Sprintf("Get %s score", tc.player)
		t.Run(testName, func(t *testing.T) {
			request := poker.NewGetScoreRequest(tc.player)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)
			poker.AssertStatus(t, response, http.StatusOK)
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
	server, err := poker.NewPlayerServer(&store, nil)
	poker.AssertNoError(t, err)

	t.Run("It returns accepted on Post", func(t *testing.T) {
		request := poker.NewPostWinRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		poker.AssertStatus(t, response, http.StatusAccepted)
		poker.AssertPlayerWin(t, &store, "Pepper")
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := MustMakePlayerServer(t, &poker.StubPlayerStore{}, dummyGame)

		request := poker.NewGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusOK)
	})

	t.Run("start game with 3 players and finish fame with 'Crhis' as winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		winner := "Ruth"

		spyGame := SpyGame{BlindAlert: []byte(wantedBlindAlert)}
		svr := httptest.NewServer(MustMakePlayerServer(t, dummyPlayerStore, &spyGame))
		ws := poker.MustDialWS(t, "ws"+strings.TrimPrefix(svr.URL, "http")+"/ws")

		defer svr.Close()
		defer ws.Close()

		poker.WriteWSMessage(t, ws, "3")
		poker.WriteWSMessage(t, ws, winner)

		assertStartCalledWith(t, &spyGame, 3)
		assertFinishCalledWith(t, &spyGame, winner)
		within(t, 10*time.Millisecond, func() { assertWebSocketGotMsg(t, ws, wantedBlindAlert) })
	})
}

func MustMakePlayerServer(t *testing.T, store poker.PlayerStore, game poker.Game) *poker.PlayerServer {
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func assertWebSocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, msg, _ := ws.ReadMessage()
	if string(msg) != want {
		t.Errorf(`got "%s", want "%s"`, string(msg), want)
	}
}

func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)
	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}

}
