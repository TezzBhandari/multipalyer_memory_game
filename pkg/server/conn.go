package server

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Conn struct {
	id    int
	conn  *websocket.Conn
	msgs  chan []byte
	relay *Relay
}

func (c *Conn) Read() {
	defer func() {
		c.relay.remove(c.id)
		c.conn.Close()
	}()

	for {
		mt, data, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		if mt == websocket.CloseMessage {
			fmt.Println("closing message", string(data))
		}

		c.relay.relay(data)
		fmt.Printf("clientId:%d msg:%s\n", c.id, data)

	}
}

func (c *Conn) Write() {
	defer func() {
		c.relay.remove(c.id)
		c.conn.Close()
	}()

	for {
		msg := <-c.msgs
		err := c.conn.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			break
		}
	}
}

func (c *Conn) msg(msg []byte) {
	select {
	case c.msgs <- msg:
	default:

		fmt.Printf("msg dropped for %d, msg: %s", c.id, msg)
		time.Sleep(15 * time.Second)
	}
}
