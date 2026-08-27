package wts

import (
	"context"
	"sync"
)

type MessageHandler = func(Context) error
type DisconnectedHandler = func(context.Context, *Hub, *Session) error

type MessageHandlerRegistry struct {
	onConnected    []MessageHandler
	onDisconnected []DisconnectedHandler
	request        map[MessageType][]MessageHandler
	mu             sync.RWMutex
}

type HandlerRegister interface {
	Register(m *MessageHandlerRegistry)
}

type HandlerRegisters = []HandlerRegister

func newMessageHandlerRegistry() *MessageHandlerRegistry {
	return &MessageHandlerRegistry{
		request: make(map[MessageType][]MessageHandler),
	}
}

func (m *MessageHandlerRegistry) OnConncted(handlers ...MessageHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onConnected = append(m.onConnected, handlers...)
}

func (m *MessageHandlerRegistry) OnDisconnected(handler DisconnectedHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onDisconnected = append(m.onDisconnected, handler)
}

func (m *MessageHandlerRegistry) OnRequest(t MessageType, handlers ...MessageHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.request[t] = append(m.request[t], handlers...)
}
