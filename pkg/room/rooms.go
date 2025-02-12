package room

import (
	"log"

	"github.com/gorilla/websocket"
)

type Rooms struct {
	rooms map[int]*Room
}

func NewRooms() *Rooms {
	rooms := &Rooms{
		rooms: make(map[int]*Room),
	}
	return rooms
}

func (r *Rooms) JoinRoom(conn *websocket.Conn) {
	if conn == nil {
		log.Println("nil websocket connection")
		return
	}

	var currentRoom *Room

	for _, room := range r.rooms {
		if room.players[1] == nil {
			currentRoom = room
			break
		} else {
			log.Println("created new room")
			currentRoom = NewRoom()
			r.rooms[currentRoom.id] = currentRoom
		}
	}

	if currentRoom == nil {
		log.Println("created new room")
		currentRoom = NewRoom()
		r.rooms[currentRoom.id] = currentRoom
	}

	player := NewPlayer(conn, currentRoom)
	log.Printf("created new player %d\n", player.id)
	currentRoom.addPlayer(player)
}
