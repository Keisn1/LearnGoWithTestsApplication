package poker

import (
	"time"
)

type Game struct {
	nbrOfPlayers int
	playerStore  PlayerStore
	alerter      BlindAlerter
}

func (game *Game) scheduleBlindAlerts() {
	blindIncrement := time.Duration(5+game.nbrOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Minute
	for _, blind := range blinds {
		game.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += blindIncrement
	}
}

func (game *Game) RecordWinner(winner string) {
	game.playerStore.RecordWin(winner)
}

func (game *Game) PlayPoker() {
	game.scheduleBlindAlerts()
}

func NewGame(nbrOfPlayers int, s PlayerStore, a BlindAlerter) *Game {
	return &Game{
		nbrOfPlayers: nbrOfPlayers,
		playerStore:  s,
		alerter:      a,
	}
}
