package chatd

import (
	"api/internal/model/entity"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestAppendMessage(t *testing.T) {

	ctx := context.Background()
	var (
		userId string = "tester"
	)

	t.Run("first message creates a new bucket", func(t *testing.T) {
		var roomId uint = 1

		require.NoError(t, test.chatMessageDao.AppendMessage(ctx, entity.NewChatMessage(roomId, userId, "hello")))

		var bucket entity.ChatMessageBucket
		err := test.mongodb.Database().Collection(entity.TableNameChatMessageBucket).
			FindOne(ctx, bson.D{{Key: "room_id", Value: roomId}}).
			Decode(&bucket)
		require.NoError(t, err)

		assert.Equal(t, 1, bucket.Count)
		assert.Len(t, bucket.Messages, 1)
		assert.Equal(t, "hello", bucket.Messages[0].Content)
		assert.False(t, bucket.IsFull)
	})

	t.Run("bucket marks full at threshold", func(t *testing.T) {
		var roomId uint = 2
		for i := 0; i < BucketMaxCount; i++ {
			require.NoError(t,
				test.chatMessageDao.AppendMessage(ctx, entity.NewChatMessage(roomId, fmt.Sprint("foo", i), "msg")))
		}

		var bucket entity.ChatMessageBucket
		messageBucketCollection := test.mongodb.Database().Collection(entity.TableNameChatMessageBucket)
		messageBucketCollection.
			FindOne(ctx, bson.D{{Key: "room_id", Value: roomId}}).
			Decode(&bucket)

		assert.True(t, bucket.IsFull)
		assert.Equal(t, BucketMaxCount, bucket.Count)

		// 再寫一則應該開新 bucket,而不是塞進滿的那個
		err := test.chatMessageDao.AppendMessage(ctx, entity.NewChatMessage(roomId, "bar", "overflow"))
		require.NoError(t, err)

		count, _ := messageBucketCollection.
			CountDocuments(ctx, bson.D{{Key: "room_id", Value: roomId}})
		assert.Equal(t, int64(2), count)
	})
}

func TestGetRecentMessages_Pagination(t *testing.T) {

	ctx := context.Background()
	var roomId uint = 1

	// 塞 150 則,跨 2 個 bucket
	for i := 0; i < 150; i++ {
		err := test.chatMessageDao.AppendMessage(ctx, entity.NewChatMessage(
			roomId, fmt.Sprintf("user%d", i), fmt.Sprintf("msg-%d", i)))
		require.NoError(t, err)
		time.Sleep(time.Millisecond)
	}

	// 第一頁
	page1, cursor, err := test.chatMessageDao.GetRecentMessages(ctx, roomId, 50, nil)
	require.NoError(t, err)
	require.Len(t, page1, 50)
	assert.Equal(t, "msg-149", page1[0].Content) // 最新的在最前面
	require.NotNil(t, cursor)

	// 第二頁,用第一頁回傳的 cursor
	page2, _, err := test.chatMessageDao.GetRecentMessages(ctx, roomId, 50, cursor)
	require.NoError(t, err)
	require.Len(t, page2, 50)
	assert.Equal(t, "msg-99", page2[0].Content)

	// 確認兩頁沒有重複、沒有遺漏
	seen := map[string]bool{}
	for _, m := range append(page1, page2...) {
		assert.False(t, seen[m.Content], "duplicate message across pages: %s", m.Content)
		seen[m.Content] = true
	}
}
