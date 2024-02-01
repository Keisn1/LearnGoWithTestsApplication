package poker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
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

func (cli *CLI) scheduleBlindAlerts(nbrOfPlayers int) {
	blindIncrement := time.Duration(5+nbrOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Minute
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += blindIncrement
	}
}

func (cli *CLI) getNbrOfPlayers() (int, error) {
	nbrOfPlayers, err := strconv.Atoi(cli.readLine())
	if err != nil {
		return -1, fmt.Errorf("Could not convert Userinput to int for nbrOfPlayers, %v", err)
	}
	return nbrOfPlayers, err
}

func (cli *CLI) PlayPoker() {
	cli.promptForPlayers()
	nbrOfPlayers, err := cli.getNbrOfPlayers()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	cli.scheduleBlindAlerts(nbrOfPlayers)

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
