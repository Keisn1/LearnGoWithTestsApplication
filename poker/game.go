package poker

import (
	"bufio"
	"io"
	"strings"
	"time"
)

type Game struct {
	nbrOfPlayers int
	in           *bufio.Scanner
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

func (game *Game) PlayPoker() {
	game.scheduleBlindAlerts()
	userInput := game.readLine()
	game.playerStore.RecordWin(extractWinner(userInput))
}

func (game *Game) readLine() string {
	game.in.Scan()
	return game.in.Text()
}

func NewGame(nbrOfPlayers int, in io.Reader, s PlayerStore, a BlindAlerter) *Game {
	return &Game{
		playerStore:  s,
		alerter:      a,
		in:           bufio.NewScanner(in),
		nbrOfPlayers: nbrOfPlayers,
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
