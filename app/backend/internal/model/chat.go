package model

import (
	"api/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChatPayload struct {
	Id         string `msgpack:"id"`
	SenderName string `msgpack:"senderName"`
	SenderId   string `msgpack:"senderId"`
	Text       string `msgpack:"text"`
	Timestamp  int64  `msgpack:"timestamp"`
}

type ChatRequestPayload struct {
	Text string `msgpack:"text"`
}

// Cursor: 上一頁最後讀到的 (bucket_seq, message_index)
type ChatCursor struct {
	BucketSeq int64
	MsgId     bson.ObjectID
}

func EncodeChatPayload(messageId string, user *User, text string, createdAt time.Time) ([]byte, error) {

	result := ChatPayload{
		Id:         messageId,
		SenderId:   user.Id,
		SenderName: user.Name,
		Text:       text,
		Timestamp:  createdAt.Unix(),
	}

	b, err := utils.Encode(result)
	if err != nil {
		return nil, err
	}
	return b, nil
}
