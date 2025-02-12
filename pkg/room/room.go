package room

import (
	"log"

	"github.com/TezzBhandari/mgs/pkg/gameState"
	"github.com/TezzBhandari/mgs/pkg/message"
)

type RoomMessageType int

type RoomStatus int

const (
	WaitingForPlayerOne RoomStatus = iota + 1
	WaitingForPlayerTwo
	GameStarted
)

const (
	RoomJoined = iota + 1
	RoomGameStarted
	RoomDisconnected
)

type Room struct {
	id         int
	roomStatus RoomStatus
	players    [2]*Player
	gameState  *gameState.GameState
}

func NewRoom() *Room {
	mutex.Lock()
	id++
	roomId := id
	mutex.Unlock()

	return &Room{
		id:         roomId,
		roomStatus: WaitingForPlayerOne,
	}
}

func (r *Room) addPlayer(player *Player) {
	if r.roomStatus == WaitingForPlayerOne {
		r.players[0] = player
		log.Println("player 1 joined")
		r.roomStatus = WaitingForPlayerTwo
		player.Send(message.ServerMessage{
			MsgType: message.RoomJoined,
		})
	} else {
		r.players[1] = player
		log.Println("player 2 joined")
		player.Send(message.ServerMessage{
			MsgType: message.RoomJoined,
		})
		// room is full
		// start the game
		r.roomStatus = GameStarted
		r.gameState = gameState.NewGameState()
		log.Println("game started")
		r.broadcast(message.ServerMessage{
			MsgType: message.GameUpdate,
			Data:    *r.gameState,
		})
	}
	go player.Listen()
}

func (r *Room) handleClientMsg(msg message.ClientMessage) {
	switch msg.MsgType {
	case message.ClientMessageType(message.Inc):
		r.gameState.IncCounter()
	case message.ClientMessageType(message.Dec):
		r.gameState.DecCounter()
	default:
		log.Println("invalid client message type")
		return
	}

	stateMsg := message.ServerMessage{
		MsgType: message.GameUpdate,
		Data:    *r.gameState,
	}
	r.broadcast(stateMsg)
}

func (r *Room) remove(playerId int) {
	log.Println("player removed", playerId)
}

func (r *Room) broadcast(msg message.ServerMessage) {

	log.Println("boardcasting message to all the player in the room", msg)
	for _, player := range r.players {
		player.Send(msg)
	}
}
