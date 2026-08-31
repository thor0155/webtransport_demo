package handler

import (
	"api/internal/frameworks/db"
	"api/internal/middleware"
	"api/internal/model"
	"api/internal/protocol"
	"api/internal/server/wts"
	"api/internal/service/auths"
	"api/internal/service/chats"
	"api/internal/utils/codectool"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

type loginController struct {
	logger      *zap.Logger
	authService auths.AuthService
	chatService chats.ChatService
	mysql       db.MysqlDB
}

func NewLoginController(
	logger *zap.Logger, authService auths.AuthService, chatService chats.ChatService, mysql db.MysqlDB) *loginController {
	return &loginController{
		logger:      logger.Named("login-controller"),
		authService: authService,
		chatService: chatService,
		mysql:       mysql,
	}
}

func (ctrl *loginController) Register(registry *wts.MessageHandlerRegistry) {
	registry.OnConncted(middleware.WebtransportErrorReply, ctrl.handshake)
}

func (ctrl *loginController) handshake(ctx wts.Context) error {

	if ctx.GetType() != protocol.RequestTypeHello {
		return wts.WithErrorTerminates(errors.New("first message must be hello"))
	}

	session := ctx.GetSession()
	dbCtx := db.NewContext(ctx).WithGorm(ctrl.mysql.Session())

	if err := ctrl.handleWelcome(dbCtx, session, ctx.GetPayload()); err != nil {
		return wts.WithErrorTerminates(err)
	}

	// notify a member list
	if err := ctrl.handleMembers(dbCtx, session); err != nil {
		return wts.WithErrorTerminates(err)
	}

	// notify a member joined
	if err := ctrl.handleJoin(session); err != nil {
		return wts.WithErrorTerminates(err)
	}

	// notify a info of room
	if err := ctrl.handleRoomInfo(session); err != nil {
		return wts.WithErrorTerminates(err)
	}

	return nil
}

// just the first room info
func (ctrl *loginController) handleRoomInfo(s *wts.Session) error {

	joinedRooms := s.GetRooms()
	if len(joinedRooms) == 0 {
		return nil
	}
	rooms := joinedRooms.ToRooms()
	payload := rooms.First().Info()
	data, err := codectool.Encode(payload)
	if err != nil {
		return err
	}
	wts.BroadcastAllRoom(rooms, protocol.ResponseTypeRoomInfo, data)
	return nil
}

func (ctrl *loginController) handleWelcome(dbCtx *db.Context, s *wts.Session, request []byte) error {

	var hello model.HelloPayload
	if err := codectool.Decode(request, &hello); err != nil {
		return err
	}

	if hello.Room == "" {
		return errors.New("room required")
	}
	if hello.Id == "" {
		return errors.New("user id required")
	}

	if _, err := ctrl.authService.Login(dbCtx, s, hello); err != nil {
		return err
	}

	payload := model.WelcomePayload{
		UserId:  s.User.Id,
		Session: s.GetId(),
	}

	data, err := codectool.Encode(payload)
	if err != nil {
		return err
	}

	return s.SendMessage(protocol.ResponseTypeWelcome, data)
}

func (ctrl *loginController) handleJoin(s *wts.Session) error {

	data, err := codectool.EncodeJoinedMember(&s.User, s.GetId(), true, s.GetConnectedAt())
	if err != nil {
		return err
	}

	joinedRooms := s.GetRooms()
	rooms := joinedRooms.ToRooms()
	wts.BroadcastAllRoom(rooms, protocol.ResponseTypeJoin, data)
	return nil
}

func (ctrl *loginController) handleMembers(ctx *db.Context, s *wts.Session) error {

	joinedRooms := s.GetRooms()
	if len(joinedRooms) == 0 {
		return nil
	}
	rooms := joinedRooms.ToRooms()
	room := rooms.First()

	members, err := ctrl.chatService.GetMemberList(ctx, room)
	if err != nil {
		return err
	}

	data, err := codectool.Encode(model.MembersPayload{
		Members: members,
	})
	if err != nil {
		return err
	}

	return s.SendMessage(protocol.ResponseTypeMembers, data)
}
