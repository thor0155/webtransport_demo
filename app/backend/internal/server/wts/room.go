package wts

import (
	"api/internal/model"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	defaultMaxMembers int = 128
)

type Room struct {
	name       string
	maxMembers int
	members    map[string]*Session
	mu         sync.RWMutex
	logger     *zap.Logger
}

type Rooms []*Room

func (r *Rooms) First() *Room {
	if len(*r) > 0 {
		return (*r)[0]
	}
	return nil
}

func NewRoom(name string, maxMembers int, logger *zap.Logger) *Room {
	return &Room{
		name:       name,
		maxMembers: maxMembers,
		members:    make(map[string]*Session),
		logger:     logger,
	}
}

func (r *Room) Info() *model.RoomInfoPayload {
	r.mu.Lock()
	defer r.mu.Unlock()
	return &model.RoomInfoPayload{
		Room:        r.name,
		MemberCount: len(r.members),
		MaxMembers:  r.maxMembers,
	}
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) Join(s *Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.members[s.id] = s
	s.rooms[r.name] = &SessionJoinedRoom{
		joinedAt: time.Now(),
		room:     r,
	}
}

func (r *Room) Leave(s *Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.members, s.id)
	delete(s.rooms, r.name)
}

func (r *Room) GetMemberCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.members)
}

func (r *Room) GetMembers() []*Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var results = make([]*Session, 0, len(r.members))
	for _, s := range r.members {
		results = append(results, s)
	}
	return results
}

func (r *Room) Broadcast(t MessageType, msg []byte) {

	r.mu.RLock()
	var members = make([]*Session, 0, len(r.members))
	for _, s := range r.members {
		members = append(members, s)
	}
	r.mu.RUnlock()

	for _, s := range members {
		if err := s.SendMessage(t, msg); err != nil {
			r.logger.Warn("drop message in room",
				zap.String("room", r.name), zap.String("id", s.id))
		}
	}
}

func BroadcastAllRoom(rooms []*Room, t MessageType, data []byte) {
	for _, room := range rooms {
		room.Broadcast(t, data)
	}
}
