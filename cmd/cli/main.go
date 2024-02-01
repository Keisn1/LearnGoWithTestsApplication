package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Keisn1/LearnGoWithTestsApp/poker"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Unable to open database file %v with err: %v", dbFileName, err)
	}
	store, err := poker.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	alerter := poker.BlindAlerterFunc(poker.StdOutAlerter)
	cli := poker.NewCLI(os.Stdin, store, alerter)
	cli.PlayPoker()
}
