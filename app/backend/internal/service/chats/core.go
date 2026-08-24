package chats

import (
	"api/internal/dao/chatd"
	"api/internal/frameworks/db"
	"api/internal/frameworks/obj"
	"api/internal/model"
	"api/internal/model/entity"
	"api/internal/server/wts"
	"api/internal/utils"
	"context"
	"runtime"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/panjf2000/ants/v2"
)

type ChatService interface {
	GetMemberList(ctx *db.Context, room *wts.Room) ([]*model.Member, error)
	WriteMessage(ctx context.Context, roomName string, userId string, message string) error
}

type chatService struct {
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
			if s.User.Id == user.ID {
				c.setMemberModelBySession(&member, s)
				break
			}
		}
		members = append(members, &member)
	}
	return members, nil
}

func (c *chatService) WriteMessage(ctx context.Context, roomName string, userId string, message string) error {

	work := newAsyncWork("write_chat_message")
	c.workTracker.Track(work)

	if err := c.writeChatMessageCoroutine.Submit(func() {
		defer close(work.done)
		gorm := db.GetGorm(ctx)
		room, err := c.chatRoomDao.GetRoom(gorm, roomName)
		if err != nil {
			return
		}
		c.chatMessageDao.AppendMessage(ctx, entity.NewChatMessage(room.ID, userId, message))
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func NewChatService(chatRoomDao chatd.ChatRoomDao, chatMessageDao chatd.ChatMessageDao, workTracker obj.WorkTracker) (ChatService, error) {
	svc := &chatService{
		workTracker:    workTracker,
		chatRoomDao:    chatRoomDao,
		chatMessageDao: chatMessageDao,
	}
	var err error
	if svc.writeChatMessageCoroutine, err = ants.NewPool(runtime.NumCPU()); err != nil {
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
	member.Id = user.ID
	member.Name = user.Name
}

func EncodeChatPayload(u *model.User, text string) ([]byte, error) {

	result := model.ChatPayload{
		Id:         uuid.NewString(),
		SenderId:   u.Id,
		SenderName: u.Name,
		Text:       text,
		Timestamp:  time.Now().Unix(),
	}

	b, err := utils.Encode(result)
	if err != nil {
		return nil, err
	}

	return b, nil
}
