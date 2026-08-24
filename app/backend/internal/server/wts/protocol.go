package wts

const (
	HeaderSize        = 9
	Magic      uint16 = 0x5754
	Version    uint8  = 1
)

type MessageType byte
type RequestMessageType = MessageType

const (
	RequestTypeNone RequestMessageType = iota
	RequestTypePing
	RequestTypeHello
	RequestTypeChat
)

type ResponseMessageType = MessageType

const (
	ResponseTypeNone ResponseMessageType = iota
	ResponseTypeLog
	ResponseTypePong
	ResponseTypeLeave
	ResponseTypeWelcome
	ResponseTypeMembers
	ResponseTypeJoin
	ResponseTypeChat
	ResponseTypeRoomInfo
)

type OutgoingMessage struct {
	Type MessageType
	Data []byte
}
