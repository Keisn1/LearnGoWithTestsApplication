package main

import (
	"app"
	"fmt"
	"log"
	"net/http"
	"os"
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

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{lock: sync.Mutex{}, scores: map[string]int{}}
}

func main() {
	store := NewInMemoryPlayerStore()
	svr := &app.PlayerServer{Store: store}
	log.Fatal(http.ListenAndServe(":5000", svr))
}
