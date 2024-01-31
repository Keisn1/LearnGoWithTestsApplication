package main

import (
	"app"
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
	store := &app.FileSystemPlayerStore{DB: db}
	svr := app.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", svr); err != nil {
		log.Fatalf("Could not listen on port 5000 with err: %v", err)
	}
}
