package chatd

import (
	"api/internal/frameworks/db"
	"api/internal/frameworks/obj"
	"api/internal/frameworks/utils/mongotool"
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
	GetHistory(ctx context.Context, roomId uint, limit int, cursor *HistoryCursor) (*HistoryResult, error)
}

// HistoryCursor 代表「上一頁讀到的最後一則訊息的位置」
// 前端把它原樣存起來,下一頁請求時帶回來即可,不需要理解內部結構
type HistoryCursor struct {
	BucketSeq int64         `json:"bucket_seq"`
	MsgId     bson.ObjectID `json:"msg_id"`
}

type HistoryResult struct {
	Messages   []*entity.ChatMessage
	NextCursor *HistoryCursor // nil 代表沒有更多歷史了
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
	if msg.Id.IsZero() {
		return errors.New("message id is zero")
	}
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

func (c *chatMessageDao) GetHistory(ctx context.Context, roomId uint, limit int, cursor *HistoryCursor) (*HistoryResult, error) {
	filter := bson.D{{Key: "room_id", Value: roomId}}
	if cursor != nil {
		filter = append(filter, bson.E{Key: "bucket_seq", Value: bson.D{{Key: mongotool.OperatorLte, Value: cursor.BucketSeq}}})
	}

	// Fetch from multiple buckets at once to avoid having to perform a second query
	// if a single bucket is insufficient to meet the limit.
	estBuckets := limit/BucketMaxCount + 2
	opts := options.Find().
		SetSort(mongotool.NewSortBuilder().Add("bucket_seq", mongotool.DESC).Build()).
		SetLimit(int64(estBuckets))

	cur, err := c.chatMessageBucketCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var result []*entity.ChatMessage
	var nextCursor *HistoryCursor
	var lastProcessedBucketSeq int64 // Remember which bucket the process was last processed in, so you don't have to look it up later.
	skipDone := cursor == nil

	for cur.Next(ctx) {
		var bucket entity.ChatMessageBucket
		if err := cur.Decode(&bucket); err != nil {
			return nil, err
		}
		lastProcessedBucketSeq = bucket.BucketSeq

		msgs := bucket.Messages
		startIdx := len(msgs) - 1

		if !skipDone && bucket.BucketSeq == cursor.BucketSeq {
			for i, m := range msgs {
				if m.Id == cursor.MsgId {
					startIdx = i - 1
					break
				}
			}
			skipDone = true
		}

		for i := startIdx; i >= 0; i-- {
			if msgs[i].DeletedAt != nil {
				continue
			}
			result = append(result, msgs[i])
			if len(result) >= limit {
				if i > 0 {
					// There's still some data left in the bucket; the cursor is pointing to the next message to read.
					nextCursor = &HistoryCursor{BucketSeq: bucket.BucketSeq, MsgId: msgs[i-1].Id}
				}
				break
			}
		}
		if len(result) >= limit {
			break
		}
	}

	// The bucket limit was just reached at the bucket boundary (all messages in this bucket were used up, but the nextCursor wasn't set).
	// It needs to be checked whether there are any older buckets later on.
	if len(result) >= limit && nextCursor == nil && len(result) > 0 {
		exists, err := c.hasEarlierBucket(ctx, roomId, lastProcessedBucketSeq)
		if err != nil {
			return nil, err
		}
		if exists {
			lastMsg := result[len(result)-1]
			nextCursor = &HistoryCursor{BucketSeq: lastProcessedBucketSeq, MsgId: lastMsg.Id}
		}
	}

	return &HistoryResult{Messages: result, NextCursor: nextCursor}, nil
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
	// 可加 unique index {room_id, bucket_seq} 搭配重試邏輯處理衝突
	return errors.WithStack(err)
}

func (c *chatMessageDao) hasEarlierBucket(ctx context.Context, roomId uint, currentSeq int64) (bool, error) {
	count, err := c.chatMessageBucketCollection.CountDocuments(ctx,
		bson.D{
			{Key: "room_id", Value: roomId},
			{Key: "bucket_seq", Value: bson.D{{Key: mongotool.OperatorLt, Value: currentSeq}}},
		},
		options.Count().SetLimit(1),
	)
	return count > 0, err
}
