package chats

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ResultOfWriteMessage struct {
	MessageId bson.ObjectID
	CreatedAt time.Time
}
