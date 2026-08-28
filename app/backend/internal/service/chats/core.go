package chats

import (
	"api/internal/dao/chatd"
	"api/internal/frameworks/db"
	"api/internal/frameworks/obj"
	"api/internal/model"
	"api/internal/model/entity"
	"api/internal/server/wts"
	"context"
	"runtime"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/panjf2000/ants/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
)

type ChatService interface {
	GetMemberList(ctx *db.Context, room *wts.Room) ([]*model.Member, error)
	WriteMessage(ctx context.Context, roomName string, userId string, message string) (ResultOfWriteMessage, error)
}

type chatService struct {
	logger                    *zap.Logger
	chatMessageDao            chatd.ChatMessageDao
	chatRoomDao               chatd.ChatRoomDao
	workTracker               obj.WorkTracker
	writeChatMessageCoroutine *ants.Pool
}

func (c *chatService) GetMemberList(ctx *db.Context, room *wts.Room) ([]*model.Member, error) {
	users, err := c.chatRoomDao.GetUserList(ctx.GetGorm(), room.Name())
	if err != nil {
		return nil, err
	}

	sessions := room.GetMembers()
	members := make([]*model.Member, 0, len(users))
	for _, user := range users {
		var member model.Member
		c.setMemberModelByUserEntity(&member, user)
		for _, s := range sessions {
			if s.User.Id == user.Id {
				c.setMemberModelBySession(&member, s)
				break
			}
		}
		members = append(members, &member)
	}
	return members, nil
}

func (c *chatService) WriteMessage(ctx context.Context, roomName string, userId string, message string) (ResultOfWriteMessage, error) {

	state := obj.NewWorkState()
	c.workTracker.Track(state)
	messageId := bson.NewObjectID()
	createdAt := time.Now()

	if err := c.writeChatMessageCoroutine.Submit(func() {
		defer state.Close()
		room, err := c.chatRoomDao.GetRoom(roomName)
		if err != nil {
			c.logger.Error("write chat message error", zap.Error(err))
			return
		}
		if err := c.chatMessageDao.AppendMessage(ctx, entity.NewChatMessage(messageId, room.Id, userId, message, createdAt)); err != nil {
			c.logger.Error("write chat message error", zap.Error(err))
		}
	}); err != nil {
		state.Close()
		return ResultOfWriteMessage{}, errors.WithStack(err)
	}

	return ResultOfWriteMessage{MessageId: messageId, CreatedAt: createdAt}, nil
}

func NewChatService(logger *zap.Logger, chatRoomDao chatd.ChatRoomDao, chatMessageDao chatd.ChatMessageDao, workTracker obj.WorkTracker) (ChatService, error) {
	name := "chat_service"
	svc := &chatService{
		workTracker:    workTracker,
		chatRoomDao:    chatRoomDao,
		chatMessageDao: chatMessageDao,
		logger:         logger.Named(name),
	}
	var err error
	if svc.writeChatMessageCoroutine, err = ants.NewPool(runtime.NumCPU(), ants.WithMaxBlockingTasks(10000)); err != nil {
		return nil, err
	}

	return svc, err
}

func (c *chatService) setMemberModelBySession(member *model.Member, session *wts.Session) {
	member.Session = session.GetId()
	member.Online = true
	member.JoinedAt = session.GetConnectedAt().Unix()
}

func (c *chatService) setMemberModelByUserEntity(member *model.Member, user *entity.User) {
	member.Id = user.Id
	member.Name = user.Name
}
