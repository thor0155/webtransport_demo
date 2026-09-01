package v202608150000

import (
	"api/internal/frameworks/utils/mongotool"
	"api/internal/model/entity"
	"context"

	"github.com/go-gormigrate/gormigrate/v2"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorm.io/gorm"
)

var GormMigration *gormigrate.Migration = &gormigrate.Migration{
	ID: "202608150000",
	Migrate: func(tx *gorm.DB) error {
		var err = tx.AutoMigrate(
			User{}, ChatRoom{}, ChatRoomUsers{})
		if err != nil {
			return err
		}

		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable(
			User{}, ChatRoom{}, ChatRoomUsers{},
		)
	},
}

var MongoMigration migrate.Migration = migrate.Migration{
	Version:     202608150000,
	Description: "establish index of chat message",
	Up: func(ctx context.Context, db *mongo.Database) error {
		coll := db.Collection(entity.TableNameChatMessageBucket)
		_, err := coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "room_id", Value: mongotool.ASC}, {Key: "bucket_seq", Value: mongotool.DESC}},
				Options: options.Index().SetUnique(true),
			},
			{Keys: bson.D{{Key: "messages._id", Value: mongotool.ASC}}},
			{
				// 加速「找出目前未滿的 bucket」這個寫入路徑的查詢
				Keys: bson.D{
					{Key: "room_id", Value: mongotool.ASC}, {Key: "is_full", Value: mongotool.ASC}, {Key: "bucket_seq", Value: mongotool.DESC}},
			},
			{Keys: bson.D{{Key: "messages.reply_to_id", Value: mongotool.DESC}}},
		})
		return err
	},
	Down: func(ctx context.Context, db *mongo.Database) error {
		coll := db.Collection(entity.TableNameChatMessageBucket)
		return coll.Drop(ctx)
	},
}
