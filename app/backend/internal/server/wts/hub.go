package wts

import (
	"sync"

	"go.uber.org/zap"
)

type Hub struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	rooms    map[string]*Room
	logger   *zap.Logger
}

type UserEvent struct {
	Session string `msgpack:"session"`
	Name    string `msgpack:"name"`
}

func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		sessions: make(map[string]*Session),
		rooms:    make(map[string]*Room),
		logger:   logger,
	}
}

func (h *Hub) Register(s *Session) {

	h.mu.Lock()
	defer h.mu.Unlock()

	h.sessions[s.id] = s
}

func (h *Hub) Unregister(s *Session) {

	h.mu.Lock()
	delete(h.sessions, s.id)
	h.mu.Unlock()

	for _, room := range h.rooms {
		room.Leave(s)

		h.mu.Lock()
		if room.GetMemberCount() == 0 {
			delete(h.rooms, room.name)
		}
		h.mu.Unlock()
	}
}

func (h *Hub) Room(name string) *Room {

	h.mu.RLock()
	room, ok := h.rooms[name]
	h.mu.RUnlock()
	if ok {
		return room
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	room = h.rooms[name]
	if room == nil {
		room = NewRoom(name, defaultMaxMembers, h.logger)
		h.rooms[name] = room
	}

	return room
}

func (h *Hub) Join(room string, s *Session) {
	h.Room(room).Join(s)
}

func (h *Hub) Leave(room string, s *Session) {
	m := h.Room(room)
	m.Leave(s)
	if m.GetMemberCount() == 0 {
		delete(h.rooms, room)
	}
}

func (h *Hub) Broadcast(t MessageType, msg []byte) {

	h.mu.RLock()
	var members = make([]*Session, 0, len(h.sessions))
	for _, s := range h.sessions {
		members = append(members, s)
	}
	h.mu.RUnlock()

	for _, s := range members {
		if err := s.SendMessage(t, msg); err != nil {
			h.logger.Warn("drop message in hub", zap.String("id", s.id))
		}
	}
}

func (h *Hub) FilterSessions(filter func(s *Session) bool) []*Session {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var filtered = make([]*Session, 0, len(h.sessions))
	for _, s := range h.sessions {
		if filter(s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func (h *Hub) AnySessions(filter func(s *Session) bool) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, s := range h.sessions {
		if filter(s) {
			return true
		}
	}
	return false
}
