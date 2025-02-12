package message

type ClientMessageType int

const (
	Inc ClientMessageType = iota + 1
	Dec
)

type ClientMessage struct {
	MsgType  ClientMessageType `json:"msgType"`
	PlayerId int64             `json:"playerId"`
	RoomId   int64             `json:"roomId"`
}
