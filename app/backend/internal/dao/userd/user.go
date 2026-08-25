package userd

import (
	"api/internal/model/entity"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

type UserDao interface {
	MustGetUser(gorm *gorm.DB, userId string, userName string) (*entity.User, error)
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

func NewUserDao() UserDao {
	return &userDao{}
}
