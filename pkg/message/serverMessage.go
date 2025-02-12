package message

import "github.com/TezzBhandari/mgs/pkg/gameState"

type ServerMessageType int

const (
	RoomJoined ServerMessageType = iota + 1
	GameUpdate
)

type ServerMessage struct {
	MsgType ServerMessageType   `json:"msgType"`
	Data    gameState.GameState `json:"data"`
}
