package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const PlayerPrompt = "Please enter the number of players: "

type CLI struct {
	in  *bufio.Scanner
	out io.Writer
}

func (g *CLI) getNbrOfPlayers() (int, error) {
	nbrOfPlayers, err := strconv.Atoi(g.readLine())
	if err != nil {
		return -1, fmt.Errorf("Could not convert Userinput to int for nbrOfPlayers, %v", err)
	}
	return nbrOfPlayers, err
}

func (g *CLI) StartGame() {
	g.promptForPlayers()
}

func NewCLI(in io.Reader, out io.Writer) CLI {
	return CLI{
		in:  bufio.NewScanner(in),
		out: out,
	}
}

func (g *CLI) promptForPlayers() {
	fmt.Fprint(g.out, PlayerPrompt)
}

func (g *CLI) readLine() string {
	g.in.Scan()
	return g.in.Text()
}
