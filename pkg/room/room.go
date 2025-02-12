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

type RoomJoinedMsg struct {
	RoomId   int `json:"roomId"`
	PlayerId int `json:"playerId"`
}

type Room struct {
	Id         int
	roomStatus RoomStatus
	Players    [2]*Player
	gameState  *gameState.GameState
}

func NewRoom() *Room {
	mutex.Lock()
	id++
	roomId := id
	mutex.Unlock()

	return &Room{
		Id:         roomId,
		roomStatus: WaitingForPlayerOne,
	}
}

func (r *Room) AddPlayer(player *Player) {
	if r.roomStatus == WaitingForPlayerOne {
		r.Players[0] = player
		log.Println("player 1 joined")
		r.roomStatus = WaitingForPlayerTwo
		player.Send(message.ServerMessage{
			MsgType: message.RoomJoined,
			Data: RoomJoinedMsg{
				RoomId:   r.Id,
				PlayerId: player.Id,
			},
		})
	} else {
		r.Players[1] = player
		log.Println("player 2 joined")
		player.Send(message.ServerMessage{
			MsgType: message.RoomJoined,
			Data: RoomJoinedMsg{
				RoomId:   r.Id,
				PlayerId: player.Id,
			},
		})
		// room is full
		// start the game
		r.roomStatus = GameStarted
		r.gameState = gameState.NewGameState()
		r.gameState.Turn = r.Players[0].Id
		log.Println("game started")
		r.broadcast(message.ServerMessage{
			MsgType: message.GameUpdate,
			Data:    r.gameState,
		})
	}
	go player.Listen()
}

func (r *Room) handleClientMsg(msg message.ClientMessage) {
	if msg.PlayerId != r.gameState.Turn {
		return
	}

	var nextPlayer int
	if msg.PlayerId == r.Players[0].Id {
		nextPlayer = r.Players[1].Id
	} else {
		nextPlayer = r.Players[0].Id
	}

	switch msg.MsgType {
	case message.ClientMessageType(message.Inc):
		r.gameState.IncCounter(nextPlayer)
	case message.ClientMessageType(message.Dec):
		r.gameState.DecCounter(nextPlayer)
	default:
		log.Println("invalid client message type")
		return
	}

	stateMsg := message.ServerMessage{
		MsgType: message.GameUpdate,
		Data:    r.gameState,
	}
	r.broadcast(stateMsg)
}

func (r *Room) remove(playerId int) {
	log.Println("player removed", playerId)
}

func (r *Room) broadcast(msg message.ServerMessage) {

	log.Println("boardcasting message to all the player in the room", msg)
	for _, player := range r.Players {
		player.Send(msg)
	}
}
