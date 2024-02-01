package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

type tape struct {
	file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Truncate(0)
	t.file.Seek(0, io.SeekStart)
	return t.file.Write(p)
}

type FileSystemPlayerStore struct {
	DB     *json.Encoder
	league League
}

func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, io.SeekStart)
	info, err := file.Stat()
	if err != nil {
		err = fmt.Errorf("could get fileinfo, %v", err)
		return err
	}

	if info.Size() < 2 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}
	return nil
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	err := initialisePlayerDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("Problem initialising player db file, %v", err)
	}

	league, err := NewLeague(file)
	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}
	return &FileSystemPlayerStore{
		DB:     json.NewEncoder(&tape{file}),
		league: league,
	}, nil
}

func (f *FileSystemPlayerStore) GetLeagueTable() League {
	sort.Sort(f.league)
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) (int, StoreError) {
	player := f.league.Find(name)
	if player != nil {
		return player.Wins, ""
	}
	return 0, PlayerNotFoundError
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.league
	player := league.Find(name) // Find returns a pointer
	if player != nil {
		player.Wins++
	} else {
		f.league = append(league, Player{Name: name, Wins: 1})
	}

	f.DB.Encode(f.league)
}
