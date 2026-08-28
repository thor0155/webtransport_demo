package model

import (
	"api/internal/utils"
	"time"
)

type ChatMessage struct {
	Id         string `msgpack:"id"`
	SenderName string `msgpack:"senderName"`
	SenderId   string `msgpack:"senderId"`
	Text       string `msgpack:"text"`
	Timestamp  int64  `msgpack:"timestamp"`
}

type ChatMessageRequest struct {
	Room string `msgpack:"room"`
	Text string `msgpack:"text"`
}

// Cursor: 上一頁最後讀到的 (bucket_seq, message_index)
type ChatCursor struct {
	BucketSeq int64  `msgpack:"bucketSeq"`
	MsgId     string `msgpack:"msgId"`
}

type ChatHistoryRequest struct {
	RoomId string      `msgpack:"roomId"`
	Cursor *ChatCursor `msgpack:"cursor"` // 第一次請求不帶,之後每次帶上一次回應的游標
	Limit  int         `msgpack:"limit"`  // 預設 50
}

type ChatHistoryResponse struct {
	Messages   []ChatMessage `msgpack:"messages"`
	NextCursor *ChatCursor   `msgpack:"nextCursor"`
	Error      string        `msgpack:"error"`
}

func EncodeChatMessage(messageId string, user *User, text string, createdAt time.Time) ([]byte, error) {

	result := ChatMessage{
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
