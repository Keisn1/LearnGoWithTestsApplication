package poker

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

type SpyBlindAlerter struct {
	Alerts []ScheduledAlert
}

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	s.Alerts = append(s.Alerts, ScheduledAlert{duration, amount})
}

func AssertScheduledAlert(t *testing.T, got, want ScheduledAlert) {
	amountGot := got.Amount
	if amountGot != want.Amount {
		t.Errorf("got amount %d, want %d", amountGot, want.Amount)
	}

	gotScheduledTime := got.At
	if gotScheduledTime != want.At {
		t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, want.At)
	}
}

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetLeagueTable() League {
	var players []Player
	for p, w := range s.scores {
		players = append(players, Player{p, w})
	}
	return players
}

func (s *StubPlayerStore) RecordWin(name string) {

	s.winCalls = append(s.winCalls, name)
	fmt.Println(s.winCalls)
}

func (s *StubPlayerStore) GetPlayerScore(name string) (int, StoreError) {
	score, exists := s.scores[name]
	if !exists {
		return -1, PlayerNotFoundError
	}
	return score, ""
}

func NewStubPlayerStore(scores map[string]int) StubPlayerStore {
	return StubPlayerStore{scores: scores}
}

func CreateTempFile(t testing.TB, initialData string) (*os.File, func()) {
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

func AssertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}

func AssertStatus(t *testing.T, r *httptest.ResponseRecorder, want int) {
	t.Helper()
	if r.Code != want {
		t.Errorf("got = \"%v\"; want \"%v\"", r.Code, want)
	}
}

func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got = \"%v\"; want \"%v\"", got, want)
	}
}
func NewGetScoreRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)
	return request
}

func NewPostWinRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil)
	return request
}

func AssertLeague(t *testing.T, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func GetLeagueFromResponse(t *testing.T, body io.Reader) League {
	t.Helper()
	var league []Player
	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("Unable to Unmarshal response body into []string")
	}
	return league
}

func AssertEqualScores(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}

// one internal test
func TestTape_Write(t *testing.T) {
	file, clean := CreateTempFile(t, "12345")
	defer clean()

	tape := tape{file}

	tape.Write([]byte("abc"))

	file.Seek(0, 0)
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func NewGameRequest() *http.Request {
	r, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return r
}

func MustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("Could not open a ws connection on %s %v", url, err)
	}
	return ws

}
func WriteWSMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}
