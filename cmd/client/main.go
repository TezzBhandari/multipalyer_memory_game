package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:42069", "http service address")
var msg = flag.String("msg", "hi", "msg")

var (
	roomId   int
	playerId int
	send     bool = false
)

type ServerMessage struct {
	MsgType int `json:"msgType"`
	Data    any `json:"data,omitempty"`
}

type ClientMessage struct {
	MsgType  int `json:"msgType"`
	RoomId   int `json:"roomId"`
	PlayerId int `json:"playerId"`
}

type JoinRoom struct {
	RoomId   int `json:"roomId"`
	PlayerId int `json:"playerId"`
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/game"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		msg := &ServerMessage{}

		for {
			err := c.ReadJSON(msg)
			if err != nil {
				log.Println("read:", err)
				return
			}

			if msg.MsgType == 1 {
				log.Println("joined a room")
				roomId, playerId = extractRoomData(msg.Data)
				log.Println("set roomId and playerId", roomId, playerId)
				if roomId == -1 || playerId == -1 {
					log.Println("didn't get room and player id")
					break
				}
			} else if msg.MsgType == 2 {
				send = true
			}

			fmt.Printf("recv: %v %d %d\n", msg, playerId, roomId)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if send {
				clientMsg := ClientMessage{
					MsgType:  1,
					PlayerId: playerId,
					RoomId:   roomId,
				}
				log.Println("sending msg to the server", clientMsg)
				err := c.WriteJSON(clientMsg)
				if err != nil {
					log.Println("write:", err)
					return
				}
			} else {
				log.Println("not ready yet")
			}

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "client wanna close"))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// Extracts roomId and playerId from `Data`
func extractRoomData(data any) (int, int) {
	// Convert `data` (map) to JSON bytes
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error re-marshaling Data:", err)
		return -1, -1
	}

	var roomData JoinRoom
	if err := json.Unmarshal(dataBytes, &roomData); err != nil {
		fmt.Println("Error decoding Data into RoomData struct:", err)
		return -1, -1
	}

	fmt.Printf("Extracted Data → RoomID: %d, PlayerID: %d\n", roomData.RoomId, roomData.PlayerId)
	return roomData.RoomId, roomData.PlayerId
}
