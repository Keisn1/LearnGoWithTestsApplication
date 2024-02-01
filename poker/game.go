package poker

type Game struct {
	playerStore  PlayerStore
	alerter      BlindAlerter
	nbrOfPlayers int
}
