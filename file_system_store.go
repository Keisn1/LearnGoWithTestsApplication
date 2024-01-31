package app

import (
	"encoding/json"
	"io"
)

type FileSystemPlayerStore struct {
	DB io.ReadWriteSeeker
}

func (f *FileSystemPlayerStore) GetLeagueTable() League {
	f.DB.Seek(0, io.SeekStart)
	league, _ := NewLeague(f.DB)
	return league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) (int, StoreError) {
	player := f.GetLeagueTable().Find(name)
	if player != nil {
		return player.Wins, ""
	}
	return 0, PlayerNotFoundError
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeagueTable()
	player := league.Find(name)
	if player != nil {
		player.Wins++
	} else {
		league = append(league, Player{Name: name, Wins: 1})
	}

	f.DB.Seek(0, 0)
	json.NewEncoder(f.DB).Encode(league)
}
