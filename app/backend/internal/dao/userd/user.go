package userd

import (
	"api/internal/model/entity"
	"api/internal/utils/types"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

type UserDao interface {
	MustGetUser(gorm *gorm.DB, userId string, userName string) (*entity.User, error)
	GetNameMapByIds(gorm *gorm.DB, ids []string) (map[types.UserId]types.UserName, error)
}

type userDao struct {
}

func (*userDao) MustGetUser(gorm *gorm.DB, userId string, userName string) (*entity.User, error) {

	var result entity.User
	if err := gorm.Where("id=?", userId).Attrs(entity.User{
		Id:   userId,
		Name: userName,
	}).FirstOrCreate(&result).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (*userDao) GetNameMapByIds(gorm *gorm.DB, ids []string) (map[types.UserId]types.UserName, error) {
	var users []*entity.User
	if err := gorm.Select("id", "name").Where("id in ?", ids).Find(&users).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	result := make(map[types.UserId]types.UserName, len(users))
	for _, user := range users {
		result[user.Id] = user.Name
	}
	return result, nil
}

func NewUserDao() UserDao {
	return &userDao{}
}
