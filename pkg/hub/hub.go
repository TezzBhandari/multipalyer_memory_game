package hub

import (
	"log"

	"github.com/TezzBhandari/mgs/pkg/room"
	"github.com/gorilla/websocket"
)

type Hub struct {
	rooms map[int]*room.Room
}

func NewHub() *Hub {
	rooms := &Hub{
		rooms: make(map[int]*room.Room),
	}
	return rooms
}

func (h *Hub) JoinRoom(conn *websocket.Conn) {
	if conn == nil {
		log.Println("nil websocket connection")
		return
	}

	var currentRoom *room.Room

	for _, r := range h.rooms {
		if r.Players[1] == nil {
			currentRoom = r
			break
		} else {
			log.Println("created new room")
			currentRoom = room.NewRoom()
			h.rooms[currentRoom.Id] = currentRoom
		}
	}

	if currentRoom == nil {
		log.Println("created new room")
		currentRoom = room.NewRoom()
		h.rooms[currentRoom.Id] = currentRoom
	}

	player := room.NewPlayer(conn, currentRoom)
	log.Printf("created new player %d\n", player.Id)
	currentRoom.AddPlayer(player)
}
