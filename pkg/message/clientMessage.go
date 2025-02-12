package message

type ClientMessageType int

const (
	Inc ClientMessageType = iota + 1
	Dec
)

type ClientMessage struct {
	MsgType  ClientMessageType `json:"msgType"`
	PlayerId int             `json:"playerId"`
	RoomId   int             `json:"roomId"`
}
