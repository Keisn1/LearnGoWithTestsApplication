package poker

import (
	"time"
)

type Game interface {
	Start(int)
	Finish(string)
}

type TexasHoldem struct {
	store   PlayerStore
	alerter BlindAlerter
}

func NewGame(s PlayerStore, a BlindAlerter) *TexasHoldem {
	return &TexasHoldem{
		store:   s,
		alerter: a,
	}
}

func (game *TexasHoldem) Start(nbrOfPlayers int) {
	blindIncrement := time.Duration(5+nbrOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Minute
	for _, blind := range blinds {
		game.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += blindIncrement
	}
}

func (game *TexasHoldem) Finish(winner string) {
	game.store.RecordWin(winner)
}
