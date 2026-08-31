package codectool

import (
	"api/internal/model"
	"time"
)

func EncodeChatMessage(user *model.User, messageId string, text string, createdAt time.Time) ([]byte, error) {

	result := model.ChatMessage{
		Id:         messageId,
		SenderId:   user.Id,
		SenderName: user.Name,
		Text:       text,
		Timestamp:  createdAt.Unix(),
	}

	b, err := Encode(result)
	if err != nil {
		return nil, err
	}
	return b, nil
}
