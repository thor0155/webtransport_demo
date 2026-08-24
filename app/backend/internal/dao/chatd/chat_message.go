package chatd

import (
	"api/internal/frameworks/db"
	"api/internal/frameworks/obj"
	"api/internal/frameworks/utils/mongotool"
	"api/internal/model"
	"api/internal/model/entity"
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	BucketMaxCount int = 100
)

type ChatMessageDao interface {
	obj.Component
	AppendMessage(ctx context.Context, msg *entity.ChatMessage) error
	GetRecentMessages(ctx context.Context, roomId uint, limit int, cursor *model.ChatCursor) ([]*entity.ChatMessage, *model.ChatCursor, error)
}

type chatMessageDao struct {
	name                        string
	chatMessageBucketCollection *mongo.Collection
	mongodb                     db.MongoDB
}

// Active implements [ChatMessageDao].
func (c *chatMessageDao) Active() bool {
	return true
}

// Init implements [ChatMessageDao].
func (c *chatMessageDao) Init() error {
	mongodatabase := c.mongodb.Database()
	if mongodatabase == nil {
		return errors.New("mongodb database is nil")
	}
	c.chatMessageBucketCollection = mongodatabase.Collection(entity.TableNameChatMessageBucket)
	return nil
}

// Name implements [ChatMessageDao].
func (c *chatMessageDao) Name() string {
	return c.name
}

func (c *chatMessageDao) AppendMessage(ctx context.Context, msg *entity.ChatMessage) error {
	filter := bson.D{
		{Key: "room_id", Value: msg.RoomId},
		{Key: "is_full", Value: bson.D{{Key: mongotool.OperatorNe, Value: true}}},
	}
	update := mongo.Pipeline{
		bson.D{{Key: mongotool.OperatorSet, Value: bson.D{
			{Key: "messages", Value: bson.D{{Key: mongotool.OperatorConcatArrays, Value: bson.A{
				bson.D{{Key: mongotool.OperatorIfNull, Value: bson.A{"$messages", bson.A{}}}},
				bson.A{msg},
			}}}},
			{Key: "count", Value: bson.D{{Key: mongotool.OperatorAdd, Value: bson.A{
				bson.D{{Key: mongotool.OperatorIfNull, Value: bson.A{"$count", 0}}}, 1,
			}}}},
			{Key: "last_msg_at", Value: msg.CreatedAt},
		}}},
		bson.D{{Key: mongotool.OperatorSet, Value: bson.D{
			{Key: "is_full", Value: bson.D{{Key: mongotool.OperatorGte, Value: bson.A{"$count", BucketMaxCount}}}},
		}}},
	}
	opts := options.FindOneAndUpdate().
		SetSort(mongotool.NewSortBuilder().Add("bucket_seq", mongotool.DESC).Build())

	res := c.chatMessageBucketCollection.FindOneAndUpdate(ctx, filter, update, opts)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.createBucket(ctx, msg.RoomId, msg)
		}
		return errors.WithStack(err)
	}
	return nil
}

func (c *chatMessageDao) GetRecentMessages(ctx context.Context, roomId uint, limit int, cursor *model.ChatCursor) ([]*entity.ChatMessage, *model.ChatCursor, error) {
	filter := bson.D{{Key: "room_id", Value: roomId}}
	if cursor != nil {
		filter = append(filter, bson.E{Key: "bucket_seq", Value: bson.D{{Key: mongotool.OperatorLte, Value: cursor.BucketSeq}}})
	}

	// 抓比需要多一點的 bucket 數,避免最後一頁不足 limit
	estBuckets := limit/BucketMaxCount + 2
	opts := options.Find().
		SetSort(mongotool.NewSortBuilder().Add("bucket_seq", mongotool.DESC).Build()).
		SetLimit(int64(estBuckets))

	cur, err := c.chatMessageBucketCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, err
	}
	defer cur.Close(ctx)

	var result []*entity.ChatMessage
	var nextCursor *model.ChatCursor
	skipUntilCursor := cursor != nil

	for cur.Next(ctx) {
		var bucket entity.ChatMessageBucket
		if err := cur.Decode(&bucket); err != nil {
			return nil, nil, err
		}

		msgs := bucket.Messages
		startIdx := len(msgs) - 1

		// 如果這是 cursor 指向的那個 bucket,要跳過已經讀過的部分
		if skipUntilCursor && bucket.BucketSeq == cursor.BucketSeq {
			for i, m := range msgs {
				if m.Id == cursor.MsgId {
					startIdx = i - 1
					break
				}
			}
			skipUntilCursor = false
		}

		for i := startIdx; i >= 0 && len(result) < limit; i-- {
			if msgs[i].DeletedAt != nil {
				continue
			}
			result = append(result, msgs[i])
		}

		if len(result) >= limit {
			nextCursor = &model.ChatCursor{BucketSeq: bucket.BucketSeq, MsgId: msgs[startIdx-len(result)+1].Id}
			break
		}
	}

	return result, nextCursor, nil
}

func NewChatMessageDao(mongodb db.MongoDB) ChatMessageDao {
	return &chatMessageDao{
		name:    "chatMessageDao",
		mongodb: mongodb,
	}
}

func (c *chatMessageDao) createBucket(ctx context.Context, roomId uint, msg *entity.ChatMessage) error {
	bucket := entity.ChatMessageBucket{
		RoomId:     roomId,
		BucketSeq:  time.Now().UnixNano(),
		Count:      1,
		FirstMsgAt: msg.CreatedAt,
		LastMsgAt:  msg.CreatedAt,
		Messages:   []*entity.ChatMessage{msg},
	}
	_, err := c.chatMessageBucketCollection.InsertOne(ctx, bucket)
	// 極端 race 下兩個請求同時判斷「沒有未滿 bucket」而各自嘗試新建,
	// 可加 unique index {channel_id, bucket_seq} 搭配重試邏輯處理衝突
	return errors.WithStack(err)
}
