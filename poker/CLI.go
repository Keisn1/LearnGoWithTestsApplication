package poker

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

const PlayerPrompt = "Please enter the number of players: "

type CLI struct {
	in          *bufio.Scanner
	out         io.Writer
	playerStore PlayerStore
	alerter     BlindAlerter
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Minute
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += 10 * time.Minute
	}
}

func (cli *CLI) PlayPoker() {
	cli.promptForPlayers()
	cli.scheduleBlindAlerts()
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *CLI) promptForPlayers() {
	fmt.Fprint(cli.out, PlayerPrompt)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func NewCLI(in io.Reader, out io.Writer, s PlayerStore, a BlindAlerter) *CLI {
	return &CLI{
		in:          bufio.NewScanner(in),
		out:         out,
		playerStore: s,
		alerter:     a,
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
