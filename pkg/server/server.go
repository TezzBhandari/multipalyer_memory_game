package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TezzBhandari/mgs/pkg/room"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type GameServer struct {
	addr   string
	s      http.Server
	router *mux.Router
	rooms  *room.Rooms
}

func NewServer(addr string) *GameServer {
	router := mux.NewRouter()
	rooms := room.NewRooms()

	return &GameServer{
		router: router,
		addr:   addr,
		rooms:  rooms,
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

	gs.rooms.JoinRoom(wsc)
}
