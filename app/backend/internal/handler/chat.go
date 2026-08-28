package handler

import (
	"api/internal/frameworks/db"
	"api/internal/model"
	"api/internal/protocol"
	"api/internal/server/wts"
	"api/internal/service/auths"
	"api/internal/service/chats"
	"api/internal/utils"
	"context"

	"go.uber.org/zap"
)

type chatController struct {
	logger      *zap.Logger
	chatService chats.ChatService
	mysql       db.MysqlDB
}

func NewChatController(
	logger *zap.Logger, authService auths.AuthService, chatService chats.ChatService, mysql db.MysqlDB) *chatController {
	return &chatController{
		logger:      logger.Named("chat-controller"),
		chatService: chatService,
		mysql:       mysql,
	}
}

func (c *chatController) Register(registry *wts.MessageHandlerRegistry) {
	registry.OnDisconnected(c.leave)
	registry.OnRequest(protocol.RequestTypeChat, c.handleChat)
	registry.OnRequest(protocol.RequestTypePing, c.handlePing)

}

func (c *chatController) leave(ctx context.Context, hub *wts.Hub, session *wts.Session) error {

	// who leaved
	data, err := model.EncodeLeavedMember(session.GetId(), session.User.Id)
	if err != nil {
		return err
	}

	joinedRooms := session.GetRooms()
	if len(joinedRooms) == 0 {
		return nil
	}
	rooms := joinedRooms.ToRooms()
	wts.BroadcastAllRoom(rooms, protocol.ResponseTypeLeave, data)

	// update room info. TODO: maybe more rooms
	room := rooms.First()
	room.Leave(session)
	data, err = utils.Encode(room.Info())
	if err != nil {
		return err
	}
	wts.BroadcastAllRoom(rooms, protocol.ResponseTypeRoomInfo, data)

	return nil
}

func (c *chatController) handleChat(ctx wts.Context) error {

	var request model.ChatMessageRequest
	if err := utils.Decode(ctx.GetPayload(), &request); err != nil {
		return err
	}

	session := ctx.GetSession()
	room := session.GetRoom()
	dbCtx := db.WithGorm(ctx, c.mysql.Session())

	messageId, createdAt, err := c.chatService.WriteMessage(dbCtx, room.Name(), ctx.GetSession().User.Id, request.Text)
	if err != nil {
		// TODO: if an ants.ErrPoolOverload happened, just log error & send to client only
		return err
	}

	b, err := model.EncodeChatPayload(messageId.Hex(), &ctx.GetSession().User, request.Text, createdAt)
	if err != nil {
		return err
	}
	// c.logger.Debug("broadcast chat", zap.String("messageId", messageId.Hex()), zap.Time("createdAt", createdAt))

	wts.BroadcastAllRoom([]*wts.Room{room}, protocol.ResponseTypeChat, b)

	return nil
}

func (c *chatController) handlePing(ctx wts.Context) error {
	return ctx.GetSession().SendMessage(protocol.ResponseTypePong, nil)
}
