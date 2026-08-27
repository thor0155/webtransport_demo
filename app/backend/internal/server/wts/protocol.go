package wts

const (
	HeaderSize        = 9
	Magic      uint16 = 0x5754
	Version    uint8  = 1
)

type MessageType = byte

type outgoingMessage struct {
	Type MessageType
	Data []byte
}
