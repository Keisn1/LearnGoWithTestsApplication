package poker

import (
	"io"
	"time"
)

type Game interface {
	Start(nbrOfPlayers int, alertsDestination io.Writer)
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

func (game *TexasHoldem) Start(nbrOfPlayers int, to io.Writer) {
	blindIncrement := time.Duration(5+nbrOfPlayers) * time.Second

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		game.alerter.ScheduleAlertAt(blindTime, blind, to)
		blindTime += blindIncrement
	}
}

func (game *TexasHoldem) Finish(winner string) {
	game.store.RecordWin(winner)
}
