package room

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/TezzBhandari/mgs/pkg/message"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
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

	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(string) error {
		log.Println("pong")
		p.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		msg := &message.ClientMessage{}
		err := p.conn.ReadJSON(msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v\n", err)
				break
			}

			log.Printf("error: %v\n", err)
			break
		}

		p.room.handleClientMsg(*msg)

		fmt.Printf("msg recieved from clientId:%d\n", p.Id)
	}
}

func (p *Player) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.room.remove(p.Id)
		p.conn.Close()

	}()

	for {
		select {
		case msg, ok := <-p.msgs:
			if !ok {
				// The hub closed the channel.
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Println("message just before writing", msg)
			err := p.conn.WriteJSON(msg)
			if err != nil {
				break
			}

		case <-ticker.C:
			log.Println("ping")
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				break
			}
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
