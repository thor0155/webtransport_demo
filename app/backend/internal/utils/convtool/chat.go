package convtool

import (
	"api/internal/dao/chatd"
	"api/internal/model"
	"api/internal/utils/types"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ConvChatDaoHistroyCursor(cursor *model.ChatCursor) (*chatd.HistoryCursor, error) {
	result := &chatd.HistoryCursor{
		BucketSeq: cursor.BucketSeq,
	}
	var err error
	if result.MsgId, err = bson.ObjectIDFromHex(cursor.MsgId); err != nil {
		return nil, err
	}
	return result, nil
}

func ConvChatDtoHistoryCursor(cursor *chatd.HistoryCursor) *model.ChatCursor {
	return &model.ChatCursor{
		BucketSeq: cursor.BucketSeq,
		MsgId:     cursor.MsgId.Hex(),
	}

}

func ConvChatHistoryResponse(result *chatd.HistoryResult, nameMap map[types.UserId]types.UserName) model.ChatHistoryResponse {
	messagesDto := make([]model.ChatMessage, 0, len(result.Messages))
	for _, m := range result.Messages {
		messagesDto = append(messagesDto, model.ChatMessage{
			Id:         m.Id.Hex(),
			SenderId:   m.UserId,
			SenderName: nameMap[m.UserId],
			Text:       m.Content,
			Timestamp:  m.CreatedAt.UnixMilli(),
		})
	}
	response := model.ChatHistoryResponse{
		Messages: messagesDto,
	}
	if result.NextCursor != nil {
		response.NextCursor = ConvChatDtoHistoryCursor(result.NextCursor)
	}
	return response
}
