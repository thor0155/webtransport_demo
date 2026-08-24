package wts

import (
	"api/internal/frameworks/utils/nettool"
	"api/internal/model"
	"api/internal/utils"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/quic-go/webtransport-go"
	"go.uber.org/zap"
)

const (
	ClosedErrorCode webtransport.SessionErrorCode = iota
	MaintainationErrorCode
	IllegalErrorCode
)

type Session struct {
	User        model.User
	send        chan OutgoingMessage
	done        chan struct{}
	writeDone   chan struct{}
	connectedAt time.Time
	id          string
	conn        *webtransport.Session
	once        sync.Once
	closed      atomic.Bool
	logger      *zap.Logger
	stream      *webtransport.Stream
	closeEvent  []func()
	rooms       map[string]*SessionJoinedRoom
}

type SessionJoinedRoom struct {
	room     *Room
	joinedAt time.Time
}

type SessionJoinedRooms []*SessionJoinedRoom

func (s *SessionJoinedRooms) ToRooms() Rooms {
	rooms := make(Rooms, 0, len(*s))
	for _, r := range *s {
		rooms = append(rooms, r.room)
	}
	return rooms
}

func (s *SessionJoinedRoom) Room() *Room {
	return s.room
}

func (s *SessionJoinedRoom) JoinedAt() time.Time {
	return s.joinedAt
}

func (s *Session) GetConnectedAt() time.Time {
	return s.connectedAt
}

func (s *Session) GetId() string {
	return s.id
}
func (s *Session) GetUserName() string {
	return s.User.Name
}
func (s *Session) GetUserId() string {
	return s.User.Id
}
func (s *Session) AddClosedEvent(f func()) {
	s.closeEvent = append(s.closeEvent, f)
}

func (s *Session) SendMessage(t MessageType, data []byte) error {

	if s.closed.Load() {
		return errors.New("failed send message, it's closed")
	}

	select {
	case s.send <- OutgoingMessage{Type: t, Data: data}:
		return nil
	default:
		// queue full
		return errors.New("message pipeline full")

	}
}

func (s *Session) Close() {
	s.once.Do(func() {
		if s.closed.Load() {
			return
		}
		close(s.done)
		close(s.send)
		<-s.writeDone
		s.stream.Close()
		s.conn.CloseWithError(ClosedErrorCode, "close")
		for _, f := range s.closeEvent {
			f()
		}
		s.closeEvent = nil
		s.closed.Store(true)
	})
}

func (s *Session) IsClosed() bool {
	return s.closed.Load()
}

func (s *Session) WriteLoop(ctx context.Context) {

	defer close(s.writeDone)
	for {
		select {
		case msg, ok := <-s.send:
			if !ok {
				return
			}

			err := Write(s.stream, msg.Type, msg.Data)
			if err != nil {
				if nettool.IsExpectedDisconnect(err) {
					return
				}
				s.logger.Error("stream write",
					zap.Error(err), zap.Uint8("message-type", uint8(msg.Type)))
			}
		case <-ctx.Done():
			return
		case <-s.done:
			if len(s.send) > 0 {
				continue
			}
			return
		}
	}
}

func (s *Session) ReadLoop(handler func(MessageType, []byte) error) {

	for {

		header, payload, err := Read(s.stream)
		if err != nil {
			if nettool.IsExpectedDisconnect(err) {
				return
			}
			s.logger.Error("read failed", zap.Error(err))
		}

		if err = handler(header.Type, payload); err != nil {
			s.logger.Error("handler fail", zap.Error(err))
			return
		}
	}

}

func (s *Session) Id() string {
	return s.id
}

func (s *Session) Connection() *webtransport.Session {
	return s.conn
}

func (s *Session) ReadMessage(v any) error {
	buf := make([]byte, 4096)
	n, err := s.stream.Read(buf)
	if err != nil {
		return errors.WithStack(err)
	}
	return utils.Decode(buf[:n], v)
}

func (s *Session) GetRooms() SessionJoinedRooms {
	var rooms = make(SessionJoinedRooms, 0, len(s.rooms))
	for _, r := range s.rooms {
		rooms = append(rooms, r)
	}
	return rooms
}

func NewSession(
	id string, wt *webtransport.Session, ctx context.Context, logger *zap.Logger,
) (*Session, error) {
	stream, err := wt.AcceptStream(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Session{
		User: model.User{
			Name: "guest",
		},
		id:          id,
		stream:      stream,
		conn:        wt,
		connectedAt: time.Now(),
		rooms:       make(map[string]*SessionJoinedRoom),
		send:        make(chan OutgoingMessage, 128),
		done:        make(chan struct{}),
		writeDone:   make(chan struct{}),
		logger:      logger,
	}, nil
}
