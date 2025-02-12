package room

import (
	"fmt"
	"log"
	"sync"

	"github.com/TezzBhandari/mgs/pkg/message"
	"github.com/gorilla/websocket"
)

var (
	id    int
	mutex sync.Mutex
)

type Player struct {
	Id   int
	conn *websocket.Conn
	room *Room
	msgs chan message.ServerMessage
}

func NewPlayer(conn *websocket.Conn, room *Room) *Player {
	mutex.Lock()
	id++
	playerId := id
	mutex.Unlock()

	return &Player{
		Id:   playerId,
		msgs: make(chan message.ServerMessage, 10),
		conn: conn,
		room: room,
	}
}

func (p *Player) read() {
	defer func() {
		// remove the player from  or check the status to offline so that only that player can reconnect no extra new players
		p.room.remove(p.Id)
		p.conn.Close()
	}()

	for {
		msg := &message.ClientMessage{}
		err := p.conn.ReadJSON(msg)
		if err != nil {
			fmt.Println(err)
			break
		}

		p.room.handleClientMsg(*msg)

		fmt.Printf("msg recieved from clientId:%d\n", p.Id)
	}
}

func (p *Player) write() {
	defer func() {
		p.room.remove(p.Id)
		p.conn.Close()
	}()

	for {
		msg := <-p.msgs
		log.Println("message just before writing", msg)
		err := p.conn.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}

func (p *Player) Send(msg message.ServerMessage) {
	select {
	case p.msgs <- msg:
	default:
		fmt.Printf("write msg dropped for %d", p.Id)
	}
}

func (p *Player) Listen() {
	go p.read()
	p.write()
}
