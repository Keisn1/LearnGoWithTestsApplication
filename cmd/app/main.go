package main

import (
	"app"
	"log"
	"net/http"
	"sync"
)

type InMemoryPlayerStore struct {
	lock   sync.Mutex
	scores map[string]int
}

func (s *InMemoryPlayerStore) GetPlayerScore(name string) (int, app.StoreError) {
	score, exists := s.scores[name]
	if !exists {
		return -1, app.PlayerNotFound
	}
	return score, ""
}

func (s *InMemoryPlayerStore) RecordWin(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.scores[name]++
	return
}

func (s *InMemoryPlayerStore) GetLeagueTable() []app.Player {
	var players []app.Player
	for p, w := range s.scores {
		players = append(players, app.Player{Name: p, Wins: w})
	}
	return players
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{lock: sync.Mutex{}, scores: map[string]int{}}
}

func main() {
	svr := app.NewPlayerServer(NewInMemoryPlayerStore())
	log.Fatal(http.ListenAndServe(":5000", svr))
}
