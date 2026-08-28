package handler

import (
	"api/internal/frameworks/db"
	"api/internal/middleware"
	"api/internal/model"
	"api/internal/protocol"
	"api/internal/server/wts"
	"api/internal/service/auths"
	"api/internal/service/chats"
	"api/internal/utils"
	"context"

	"github.com/cockroachdb/errors"
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
	registry.OnRequest(protocol.RequestTypeChat, middleware.WebtransportErrorReply, c.handleChat)
	registry.OnRequest(protocol.RequestTypePing, c.handlePing)
	registry.OnRequest(protocol.RequestTypeChatHistory, c.handleHistory)

}

func (c *chatController) handleHistory(ctx wts.Context) error {
	session := ctx.GetSession()
	requestPayload := ctx.GetPayload()

	// TODO: get history logic
	_ = requestPayload
	responsePayload := []byte{}

	return session.SendMessage(protocol.ResponseTypeChatHistory, responsePayload)
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

	for _, room := range rooms {
		room.Leave(session)
		data, err = utils.Encode(room.Info())
		if err != nil {
			return err
		}
		wts.BroadcastAllRoom(rooms, protocol.ResponseTypeRoomInfo, data)
	}

	return nil
}

func (c *chatController) handleChat(ctx wts.Context) error {

	var request model.ChatMessageRequest
	if err := utils.Decode(ctx.GetPayload(), &request); err != nil {
		return err
	}

	session := ctx.GetSession()
	room := session.FindRoom(request.Room)

	if room == nil {
		return errors.New("room not found")
	}

	dbCtx := db.WithGorm(ctx, c.mysql.Session())

	result, err := c.chatService.WriteMessage(dbCtx, request.Room, session.User.Id, request.Text)
	if err != nil {
		return err
	}

	b, err := model.EncodeChatMessage(result.MessageId.Hex(), &session.User, request.Text, result.CreatedAt)
	if err != nil {
		return err
	}
	// c.logger.Debug("broadcast chat", zap.String("messageId", messageId.Hex()), zap.Time("createdAt", createdAt))

	room.Room().Broadcast(protocol.ResponseTypeChat, b)
	return nil
}

func (c *chatController) handlePing(ctx wts.Context) error {
	return ctx.GetSession().SendMessage(protocol.ResponseTypePong, nil)
}
