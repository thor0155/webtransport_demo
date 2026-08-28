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
	User           model.User
	send           chan outgoingMessage
	done           chan struct{}
	connectedAt    time.Time
	id             string
	conn           *webtransport.Session
	closeOnce      sync.Once
	closed         atomic.Bool
	logger         *zap.Logger
	stream         *webtransport.Stream
	closeEvent     []func()
	rooms          map[string]*SessionJoinedRoom
	writeWaitGroup sync.WaitGroup
	writeMu        sync.RWMutex
}

type SessionJoinedRoom struct {
	room     *Room
	joinedAt time.Time
}

type SessionJoinedRooms []*SessionJoinedRoom

type SessionReadHandlerFunc func(*Header, []byte) error

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
	case s.send <- outgoingMessage{Type: t, Data: data}:
		return nil
	default:
		// queue full
		return errors.New("message pipeline full")

	}
}

func (s *Session) SyncSendMessage(t MessageType, data []byte) error {

	if s.closed.Load() {
		return errors.New("failed send message, it's closed")
	}
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return Write(s.stream, t, data)
}

func (s *Session) Close() {
	s.closeOnce.Do(func() {
		s.closed.Store(true)
		close(s.done)
		close(s.send)
		s.writeWaitGroup.Wait()
		s.stream.Close()
		s.conn.CloseWithError(ClosedErrorCode, "close")
		for _, f := range s.closeEvent {
			f()
		}
		s.closeEvent = nil
		s.logger.Debug("session closed", zap.String("session", s.id))
	})
}

func (s *Session) IsClosed() bool {
	return s.closed.Load()
}

func (s *Session) WriteLoop(ctx context.Context) {

	if s.closed.Load() {
		return
	}

	s.writeMu.Lock()
	s.writeWaitGroup.Add(1)
	s.writeMu.Unlock()

	defer s.writeWaitGroup.Done()
	for {
		select {
		case msg, ok := <-s.send:
			if !ok {
				return
			}

			s.writeMu.Lock()
			err := Write(s.stream, msg.Type, msg.Data)
			s.writeMu.Unlock()

			if err != nil {
				if nettool.IsExpectedDisconnect(err) {
					return
				}
				s.logger.Error("stream write",
					zap.String("session", s.id), zap.Error(err), zap.Uint8("message-type", uint8(msg.Type)))
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

func (s *Session) ReadLoop(handler SessionReadHandlerFunc) {

	for {

		header, payload, err := Read(s.stream)
		if err != nil {
			if nettool.IsExpectedDisconnect(err) {
				return
			}
			s.logger.Error("read failed", zap.Error(err))
		}

		if err = handler(header, payload); err != nil {
			s.logger.Error("handler fail and terminate", zap.String("session", s.id), zap.Error(err))
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

func (s *Session) GetRoom() *Room {

	if len(s.rooms) > 0 {
		for _, room := range s.rooms {
			return room.Room()
		}
	}
	return nil
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
		send:        make(chan outgoingMessage, 128),
		done:        make(chan struct{}),
		logger:      logger.Named("session"),
	}, nil
}
