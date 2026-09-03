package handler

import (
	"api/internal/frameworks/db"
	"api/internal/frameworks/errorx"
	"api/internal/middleware"
	"api/internal/model"
	"api/internal/protocol"
	"api/internal/server/wts"
	"api/internal/service/auths"
	"api/internal/service/chats"
	"api/internal/utils/codectool"
	"context"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

type chatController struct {
	logger      *zap.Logger
	chatService chats.ChatService
	mysql       db.MysqlDB
	redis       db.RedisDB
}

func NewChatController(
	logger *zap.Logger, authService auths.AuthService, chatService chats.ChatService,
	mysql db.MysqlDB, redis db.RedisDB) *chatController {
	return &chatController{
		logger:      logger.Named("chat-controller"),
		chatService: chatService,
		mysql:       mysql,
		redis:       redis,
	}
}

func (c *chatController) Register(registry *wts.MessageHandlerRegistry) {
	registry.OnDisconnected(c.leave)
	registry.OnRequest(protocol.RequestTypeChat, middleware.WebtransportErrorReply, c.handleChat)
	registry.OnRequest(protocol.RequestTypePing, c.handlePing)
	registry.OnRequest(protocol.RequestTypeChatHistory, middleware.WebtransportErrorReply, c.handleGetHistory)

}

func (c *chatController) handleGetHistory(ctx wts.Context) error {

	requestPayload := ctx.GetPayload()

	var request model.ChatHistoryRequest
	if err := codectool.Decode(requestPayload, &request); err != nil {
		return errorx.WithCode(err, protocol.ErrorCodeGetHistoryFailure)
	}
	if request.Room == "" {
		return errorx.NewCodeError(protocol.ErrorCodeGetHistoryFailure, "room is empty")
	}

	dbCtx := db.NewContext(ctx).WithRedis(c.redis.Client()).WithGorm(c.mysql.Session())

	response, err := c.chatService.GetHistory(dbCtx, &request)
	if err != nil {
		return errorx.WithCode(err, protocol.ErrorCodeGetHistoryFailure)
	}

	c.logger.Debug("handleGetHistory",
		zap.String("room", request.Room),
		zap.Int("limit", request.Limit),
		zap.Any("cursor", request.Cursor),
		zap.Int("messageCount", len(response.Messages)),
		zap.Any("nextCursor", response.NextCursor))

	if responsePayload, err := codectool.Encode(response); err != nil {
		return errorx.WithCode(err, protocol.ErrorCodeGetHistoryFailure)
	} else {
		return ctx.GetSession().SendMessage(protocol.ResponseTypeChatHistory, responsePayload)
	}
}

func (c *chatController) leave(ctx context.Context, hub *wts.Hub, session *wts.Session) error {

	// who leaved
	data, err := codectool.EncodeLeavedMember(session.GetId(), session.User.Id)
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
		data, err = codectool.Encode(room.Info())
		if err != nil {
			return err
		}
		wts.BroadcastAllRoom(rooms, protocol.ResponseTypeRoomInfo, data)
	}

	return nil
}

func (c *chatController) handleChat(ctx wts.Context) error {

	var request model.ChatMessageRequest
	if err := codectool.Decode(ctx.GetPayload(), &request); err != nil {
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

	b, err := codectool.EncodeChatMessage(&session.User, result.MessageId.Hex(), request.Text, result.CreatedAt)
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
