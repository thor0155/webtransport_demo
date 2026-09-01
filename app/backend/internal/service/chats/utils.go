package chats

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	DefaultHistoryCount int = 50
	MaxHistoryCount     int = 100
)

type ResultOfWriteMessage struct {
	MessageId bson.ObjectID
	CreatedAt time.Time
}
