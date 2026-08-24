package auths

import (
	"api/internal/dao/chatd"
	"api/internal/dao/userd"
	"api/internal/frameworks/db"
	"api/internal/model"
	"api/internal/model/entity"
	"api/internal/server/wts"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx *db.Context, session *wts.Session, hello model.HelloPayload) (*entity.User, error)
}

type authService struct {
	webtransportHub *wts.Hub
	userDao         userd.UserDao
	chatRoomDao     chatd.ChatRoomDao
}

// Login implements [AuthService].
func (a *authService) Login(ctx *db.Context, session *wts.Session, hello model.HelloPayload) (*entity.User, error) {
	if a.webtransportHub.AnySessions(func(s *wts.Session) bool {
		return s.User.Id == hello.Id
	}) {
		return nil, errors.Newf("User %q has logged in.", hello.Id)
	}

	var (
		err  error
		user = &entity.User{}
	)

	if err = ctx.GetGorm().Transaction(func(tx *gorm.DB) error {

		user, err = a.userDao.MustGetUser(ctx.GetGorm(), hello.Id, hello.Name)
		if err != nil {
			return err
		}

		a.setUserModel(&session.User, user)

		if err := a.chatRoomDao.CreateRoomAndUserList(ctx.GetGorm(), hello.Room, session.User.Id); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	a.webtransportHub.Join(hello.Room, session)

	return user, nil
}

func NewAuthService(webtransportHub *wts.Hub, userDao userd.UserDao, chatRoomDao chatd.ChatRoomDao) AuthService {
	return &authService{
		webtransportHub: webtransportHub,
		userDao:         userDao,
		chatRoomDao:     chatRoomDao,
	}
}

func (a *authService) setUserModel(targetModel *model.User, user *entity.User) {
	targetModel.Id = user.ID
	targetModel.Name = user.Name
}
