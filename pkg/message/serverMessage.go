package message

type ServerMessageType int

const (
	RoomJoined ServerMessageType = iota + 1
	GameUpdate
)

type ServerMessage struct {
	MsgType ServerMessageType `json:"msgType"`
	Data    any               `json:"data,omitempty"`
}
