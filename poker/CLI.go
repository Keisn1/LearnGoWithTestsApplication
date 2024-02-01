package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	Welcome      = "Let's play poker"
	PlayerPrompt = "Please enter the number of players: "
	UserInfo     = "Type {Name} wins to record a win"
)

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, Welcome, "\n")
	fmt.Fprint(cli.out, PlayerPrompt, "\n")
	nbrOfPlayers, _ := cli.GetNbrOfPlayers()

	cli.game.Start(nbrOfPlayers)

	fmt.Fprint(cli.out, UserInfo)

	winner := cli.GetWinner()

	cli.game.Finish(winner)
}

func (cli *CLI) GetNbrOfPlayers() (int, error) {
	nbrOfPlayers, err := strconv.Atoi(cli.readLine())
	if err != nil {
		return -1, fmt.Errorf("Could not convert Userinput to int for nbrOfPlayers, %v", err)
	}
	return nbrOfPlayers, err
}

func (cli *CLI) PromptForPlayers() {
	fmt.Fprint(cli.out, PlayerPrompt)
}

func (cli *CLI) GetWinner() string {
	return extractWinner(cli.readLine())
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
