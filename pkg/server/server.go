package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type GameServer struct {
	addr   string
	s      http.Server
	router *mux.Router
	relay  *Relay
}

func NewServer(addr string) *GameServer {
	relay := NewRelay()
	router := mux.NewRouter()
	return &GameServer{
		router: router,
		addr:   addr,
		s: http.Server{
			Addr:    addr,
			Handler: router,
		},
		relay: relay,
	}
}

func (gs *GameServer) Start() error {
	gs.router.HandleFunc("/game", func(rw http.ResponseWriter, r *http.Request) {
		gs.handleWsConnection(rw, r)
	})

	if err := gs.s.ListenAndServe(); err != nil {
		return err
	}
	fmt.Println("server started listening")
	return nil
}

func (gs *GameServer) handleWsConnection(rw http.ResponseWriter, r *http.Request) {
	wsc, err := upgrader.Upgrade(rw, r, nil)

	if err != nil {
		http.Error(rw, "updgrade failed", http.StatusBadRequest)
	}

	gs.relay.add(wsc)
}
