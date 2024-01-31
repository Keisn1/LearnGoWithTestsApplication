package poker

import (
	"encoding/json"
	"fmt"
	"os"
)

type League []Player

func NewLeague(file *os.File) (League, error) {
	var league League
	err := json.NewDecoder(file).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league, initialising with empty league: %v", err)
	}
	return league, err
}

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}

func (l League) Len() int {
	return len(l)
}

func (l League) Less(i, j int) bool {
	return l[i].Wins > l[j].Wins
}

func (l League) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
