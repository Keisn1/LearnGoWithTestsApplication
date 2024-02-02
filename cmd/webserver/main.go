package main

import (
	"github.com/Keisn1/LearnGoWithTestsApp/poker"
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Unable to open database file %v with err: %v", dbFileName, err)
	}
	store, err := poker.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	alerter := poker.BlindAlerterFunc(poker.Alerter)
	game := poker.NewGame(store, alerter)

	svr, err := poker.NewPlayerServer(store, game)
	if err != nil {
		panic(err)
	}
	if err := http.ListenAndServe(":5000", svr); err != nil {
		log.Fatalf("Could not listen on port 5000 with err: %v", err)
	}
}
