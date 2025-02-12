package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TezzBhandari/mgs/pkg/hub"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type GameServer struct {
	addr   string
	s      http.Server
	router *mux.Router
	hub    *hub.Hub
}

func NewServer(addr string) *GameServer {
	router := mux.NewRouter()
	hub := hub.NewHub()

	return &GameServer{
		router: router,
		addr:   addr,
		hub:    hub,
		s: http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (gs *GameServer) Start() error {
	gs.router.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("server status ok"))
	})

	gs.router.HandleFunc("/game", func(rw http.ResponseWriter, r *http.Request) {
		gs.handleWsConnection(rw, r)

	})

	fmt.Println("server listening")
	if err := gs.s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (gs *GameServer) handleWsConnection(rw http.ResponseWriter, r *http.Request) {
	wsc, err := upgrader.Upgrade(rw, r, nil)

	if err != nil {
		http.Error(rw, "updgrade failed", http.StatusBadRequest)
	}

	log.Println("connection upgraded")

	gs.hub.JoinRoom(wsc)
}
