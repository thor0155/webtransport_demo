package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	TableNameChatMessageBucket string = "chat_message_buckets"
)

type ChatMessage struct {
	Id      bson.ObjectID `bson:"_id,omitempty"` // ObjectID 自帶時間戳,天然可排序
	RoomId  uint          `bson:"room_id"`
	UserId  string        `bson:"user_id"`
	Content string        `bson:"content"`
	// Attachments []Attachment        `bson:"attachments,omitempty"`
	// Reactions   []Reaction          `bson:"reactions,omitempty"`
	// ReplyToID   *primitive.ObjectID `bson:"reply_to_id,omitempty"`
	// Mentions    []primitive.ObjectID `bson:"mentions,omitempty"`
	EditedAt  *time.Time `bson:"edited_at,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty"` // 軟刪除
	CreatedAt time.Time  `bson:"created_at"`
}

type ChatMessageBucket struct {
	Id         bson.ObjectID  `bson:"_id,omitempty"`
	RoomId     uint           `bson:"room_id"`
	BucketSeq  int64          `bson:"bucket_seq"` // 該頻道第 N 個 bucket,遞增
	Count      int            `bson:"count"`
	FirstMsgAt time.Time      `bson:"first_msg_at"`
	LastMsgAt  time.Time      `bson:"last_msg_at"`
	IsFull     bool           `bson:"is_full"` // 達上限即標記,不再寫入
	Messages   []*ChatMessage `bson:"messages"`
}

func NewChatMessage(messageId bson.ObjectID, roomId uint, userId string, content string, createdAt time.Time) *ChatMessage {
	return &ChatMessage{
		Id:        messageId,
		RoomId:    roomId,
		UserId:    userId,
		Content:   content,
		CreatedAt: createdAt}
}
func NewChatMessageSimply(roomId uint, userId string, content string) *ChatMessage {
	return &ChatMessage{
		Id:        bson.NewObjectID(),
		RoomId:    roomId,
		UserId:    userId,
		Content:   content,
		CreatedAt: time.Now()}
}
