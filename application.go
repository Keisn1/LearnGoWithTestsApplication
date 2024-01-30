package app

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

const (
	PlayerNotFound = StoreError("Player not found")
)

type StoreError string

type PlayerStore interface {
	// Return PlayerNotFound error in case that player wasn't found
	GetPlayerScore(name string) (int, StoreError)
	GetPlayers() []string
	RecordWin(name string)
}

type PlayerServer struct {
	Store PlayerStore
}

func (svr *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(svr.getLeague))
	router.Handle("/players/", http.HandlerFunc(svr.playersHandler))
	router.ServeHTTP(w, r)
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

func (svr *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	svr.getLeague(w, r)
}

func (s *PlayerServer) getLeague(w http.ResponseWriter, r *http.Request) {
	players := s.Store.GetPlayers()
	b, err := json.Marshal(players)
	if err != nil {
		slog.Error("Wasn't able to Marshal players into json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred while encoding JSON"))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

func (svr *PlayerServer) processWin(w http.ResponseWriter, player string) {
	svr.Store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func (svr *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score, err := svr.Store.GetPlayerScore(player)

	if err == PlayerNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Fprint(w, score)

}
