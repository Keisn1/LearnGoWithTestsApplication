package app

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

const (
	PlayerNotFoundError = StoreError("Player not found")
)

type StoreError string

type Player struct {
	Name string
	Wins int
}

type PlayerStore interface {
	// Return PlayerNotFound error in case that player wasn't found
	GetPlayerScore(name string) (int, StoreError)
	GetLeagueTable() League
	RecordWin(name string)
}

type PlayerServer struct {
	Store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	svr := &PlayerServer{
		Store: store,
	}
	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(svr.getLeague))
	router.Handle("/players/", http.HandlerFunc(svr.playersHandler))
	svr.Handler = router
	return svr
}

func (svr *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	svr.getLeague(w, r)
}

func (s *PlayerServer) getLeague(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	err := json.NewEncoder(w).Encode(s.Store.GetLeagueTable())
	if err != nil {
		slog.Error("Wasn't able to Marshal players into json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred while encoding JSON"))
		return
	}
}

func (svr *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	switch r.Method {
	case http.MethodGet:
		svr.showScore(w, player)
	case http.MethodPost:
		svr.processWin(w, player)
	}
}

func (svr *PlayerServer) processWin(w http.ResponseWriter, player string) {
	svr.Store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func (svr *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score, err := svr.Store.GetPlayerScore(player)

	if err == PlayerNotFoundError {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Fprint(w, score)
}
