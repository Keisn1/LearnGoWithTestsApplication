package poker

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"log"

	"github.com/gorilla/websocket"
)

const (
	PlayerNotFoundError = StoreError("Player not found")
)

type StoreError string

type Player struct {
	Name string
	Wins int
}

type PlayerStore interface {
	// Return PlayerNotFound error in case that player wasn't found
	GetPlayerScore(name string) (int, StoreError)
	GetLeagueTable() League
	RecordWin(name string)
}

type playerServerWS struct {
	*websocket.Conn
}

func (ws *playerServerWS) Write(p []byte) (n int, err error) {
	err = ws.Conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), err
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("problem upgrading connection to WebSockets %v\n", err)
	}

	return &playerServerWS{conn}
}

func (w *playerServerWS) WaitForMsg() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}
	return string(msg)
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
	game     Game
	template *template.Template
}

const htmlTemplatePath = "game.html"

func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {

	svr := new(PlayerServer)
	tmpl, err := template.ParseFiles(htmlTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("problem loading template %s %v", htmlTemplatePath, err)
	}
	svr.template = tmpl
	svr.store = store
	svr.game = game

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(svr.getLeague))
	router.Handle("/game", http.HandlerFunc(svr.gameHandler))
	router.Handle("/ws", http.HandlerFunc(svr.webSocket))
	router.Handle("/players/", http.HandlerFunc(svr.playersHandler))
	svr.Handler = router
	return svr, nil
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (svr *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)
	nbrOfPlayersMsg := ws.WaitForMsg()
	nbrOfPlayers, _ := strconv.Atoi(string(nbrOfPlayersMsg))
	svr.game.Start(nbrOfPlayers, ws)

	winnerMsg := ws.WaitForMsg()
	svr.game.Finish(string(winnerMsg))

	svr.store.RecordWin(string(winnerMsg))
}

func (svr *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	svr.getLeague(w, r)
}

func (s *PlayerServer) getLeague(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	err := json.NewEncoder(w).Encode(s.store.GetLeagueTable())
	if err != nil {
		slog.Error("Wasn't able to Marshal players into json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred while encoding JSON"))
		return
	}
}

func (svr *PlayerServer) gameHandler(w http.ResponseWriter, r *http.Request) {
	svr.getGame(w, r)
}

func (p *PlayerServer) getGame(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}

func (svr *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	switch r.Method {
	case http.MethodGet:
		svr.showScore(w, player)
	case http.MethodPost:
		svr.processWin(w, player)
	}
}

func (svr *PlayerServer) processWin(w http.ResponseWriter, player string) {
	svr.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func (svr *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score, err := svr.store.GetPlayerScore(player)

	if err == PlayerNotFoundError {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Player wasn't Found")
		return
	}
	fmt.Fprint(w, score)
}
