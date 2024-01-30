package app

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

type FileSystemPlayerStore struct {
	db io.ReadSeeker
}

func (f *FileSystemPlayerStore) GetLeagueTable() []Player {
	league, _ := NewLeague(f.db)
	f.db.Seek(0, io.SeekStart)
	return league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) (int, StoreError) {
	league, _ := NewLeague(f.db)
	f.db.Seek(0, io.SeekStart)
	for _, p := range league {
		if p.Name == name {
			return p.Wins, ""
		}
	}
	return -1, PlayerNotFoundError
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league, _ := NewLeague(f.db)
	var newLeague []Player
	for _, p := range league {
		if p.Name == name {
			p.Wins++
		}
		newLeague = append(newLeague, p)
	}

	buf := bytes.Buffer{}
	json.NewEncoder(&buf).Encode(newLeague)
	f.db = strings.NewReader(buf.String())
}
