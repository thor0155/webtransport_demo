package handler

import (
	"api/internal/frameworks/db"
	"api/internal/model"
	"api/internal/server/wts"
	"api/internal/service/auths"
	"api/internal/service/chats"
	"api/internal/utils"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

type chatController struct {
	logger      *zap.Logger
	authService auths.AuthService
	chatService chats.ChatService
	mysql       db.MysqlDB
}

func NewChatController(
	logger *zap.Logger, authService auths.AuthService, chatService chats.ChatService, mysql db.MysqlDB) *chatController {
	return &chatController{
		logger:      logger,
		authService: authService,
		chatService: chatService,
		mysql:       mysql,
	}
}

func (c *chatController) Register(registry *wts.MessageHandlerRegistry) {
	registry.OnConncted(c.handshake)
	registry.OnDisconnected(c.leave)
	registry.OnRequest(wts.RequestTypeChat, c.handleChat)
	registry.OnRequest(wts.RequestTypePing, c.handlePing)

}

func (c *chatController) leave(ctx wts.Context) error {

	s := ctx.GetSession()

	// who leaved
	data, err := model.EncodeLeavedMember(s.GetId(), s.User.Id)
	if err != nil {
		return err
	}

	joinedRooms := s.GetRooms()
	if len(joinedRooms) == 0 {
		return nil
	}
	rooms := joinedRooms.ToRooms()
	wts.BroadcastAllRoom(rooms, wts.ResponseTypeLeave, data)

	// update room info. TODO: maybe more rooms
	room := rooms.First()
	room.Leave(s)
	data, err = utils.Encode(room.Info())
	if err != nil {
		return err
	}
	wts.BroadcastAllRoom(rooms, wts.ResponseTypeRoomInfo, data)

	return nil
}

func (c *chatController) handshake(ctx wts.Context) error {

	session := ctx.GetSession()
	dbCtx := db.NewContext(ctx).WithGorm(c.mysql.Session())

	if err := c.handleWelcome(dbCtx, session, ctx.GetPayload()); err != nil {
		return err
	}

	// notify a member list
	if err := c.handleMembers(dbCtx, session); err != nil {
		return err
	}

	// notify a member joined
	if err := c.handleJoin(session); err != nil {
		return err
	}

	// notify a info of room
	if err := c.handleRoomInfo(session); err != nil {
		return err
	}

	return nil
}

func (c *chatController) handleChat(ctx wts.Context) error {

	var request model.ChatRequestPayload
	if err := utils.Decode(ctx.GetPayload(), &request); err != nil {
		return err
	}

	session := ctx.GetSession()
	joinedRooms := session.GetRooms()
	rooms := joinedRooms.ToRooms()
	room := rooms.First()
	dbCtx := db.WithGorm(ctx, c.mysql.Session())

	if err := c.chatService.WriteMessage(dbCtx, room.Name(), ctx.GetSession().User.Id, request.Text); err != nil {
		return err
	}

	b, err := chats.EncodeChatPayload(&ctx.GetSession().User, request.Text)
	if err != nil {
		return err
	}

	wts.BroadcastAllRoom(rooms, wts.ResponseTypeChat, b)

	return nil
}

func (c *chatController) handlePing(ctx wts.Context) error {
	return ctx.GetSession().SendMessage(wts.ResponseTypePong, nil)
}

// just the first room info
func (h *chatController) handleRoomInfo(s *wts.Session) error {

	joinedRooms := s.GetRooms()
	if len(joinedRooms) == 0 {
		return nil
	}
	rooms := joinedRooms.ToRooms()
	payload := rooms.First().Info()
	data, err := utils.Encode(payload)
	if err != nil {
		return err
	}
	wts.BroadcastAllRoom(rooms, wts.ResponseTypeRoomInfo, data)
	return nil
}

func (h *chatController) handleWelcome(dbCtx *db.Context, s *wts.Session, request []byte) error {

	var hello model.HelloPayload
	if err := utils.Decode(request, &hello); err != nil {
		return err
	}

	if hello.Room == "" {
		return errors.New("room required")
	}
	if hello.Id == "" {
		return errors.New("user id required")
	}

	if _, err := h.authService.Login(dbCtx, s, hello); err != nil {
		return err
	}

	payload := model.WelcomePayload{
		UserId:  s.User.Id,
		Session: s.GetId(),
	}

	data, err := utils.Encode(payload)
	if err != nil {
		return err
	}

	return s.SendMessage(wts.ResponseTypeWelcome, data)
}

func (h *chatController) handleJoin(s *wts.Session) error {

	data, err := model.EncodeJoinedMember(s.GetId(), &s.User, true, s.GetConnectedAt())
	if err != nil {
		return err
	}

	joinedRooms := s.GetRooms()
	rooms := joinedRooms.ToRooms()
	wts.BroadcastAllRoom(rooms, wts.ResponseTypeJoin, data)
	return nil
}

func (h *chatController) handleMembers(ctx *db.Context, s *wts.Session) error {

	joinedRooms := s.GetRooms()
	if len(joinedRooms) == 0 {
		return nil
	}
	rooms := joinedRooms.ToRooms()
	room := rooms.First()

	members, err := h.chatService.GetMemberList(ctx, room)
	if err != nil {
		return err
	}

	data, err := utils.Encode(model.MembersPayload{
		Members: members,
	})
	if err != nil {
		return err
	}

	return s.SendMessage(wts.ResponseTypeMembers, data)
}
