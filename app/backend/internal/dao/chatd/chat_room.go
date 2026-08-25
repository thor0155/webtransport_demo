package chatd

import (
	"api/internal/frameworks/db"
	"api/internal/frameworks/utils/gormtool"
	"api/internal/model/entity"
	"fmt"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ChatRoomDao interface {
	GetRoom(gorm *gorm.DB, room string) (*entity.ChatRoom, error)
	GetUserList(gorm *gorm.DB, room string) ([]*entity.User, error)
	CreateRoomAndUserList(gorm *gorm.DB, room string, userId string) error
}

type chatRoomDao struct {
	chatRoomUserInsertIgnoreClause clause.OnConflict
}

func (c *chatRoomDao) GetRoom(gorm *gorm.DB, room string) (*entity.ChatRoom, error) {
	var result entity.ChatRoom
	if err := gorm.Take(&result, "name=?", room).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (c *chatRoomDao) CreateRoomAndUserList(gormDb *gorm.DB, roomName string, userId string) error {
	var room entity.ChatRoom
	if err := gormDb.Where("name=?", roomName).Attrs(entity.ChatRoom{Name: roomName}).FirstOrCreate(&room).Error; err != nil {
		return errors.WithStack(err)
	}
	if err := gormDb.Clauses(c.chatRoomUserInsertIgnoreClause).Create(entity.ChatRoomUsers{
		ChatRoomId: room.Id,
		UserId:     userId,
	}).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *chatRoomDao) GetUserList(gorm *gorm.DB, room string) ([]*entity.User, error) {
	var results []*entity.User
	if err := gorm.Clauses(gormtool.LockForUpdate).Table(entity.TableNameChatRoom+" cr").
		Select("u.id, u.name, u.created_at").
		Joins(fmt.Sprintf("JOIN %s cru ON cru.chat_room_id = cr.id", entity.TableNameChatRoomUsers)).
		Joins(fmt.Sprintf("JOIN %s u ON u.id = cru.user_id", entity.TableNameUser)).
		Where("cr.name = ?", room).
		Find(&results).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return results, nil
}

func NewChatRoomDao(mongodb db.MongoDB) ChatRoomDao {
	return &chatRoomDao{
		chatRoomUserInsertIgnoreClause: clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "chat_room_id"}},
			DoNothing: true,
		},
	}
}
