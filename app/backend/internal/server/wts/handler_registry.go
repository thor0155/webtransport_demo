package wts

import "sync"

type MessageHandler = func(Context) error

type MessageHandlerRegistry struct {
	onConnected    []MessageHandler
	onDisconnected []MessageHandler
	request        map[MessageType]MessageHandler
	mu             sync.RWMutex
}

type HandlerRegister interface {
	Register(m *MessageHandlerRegistry)
}

type HandlerRegisters = []HandlerRegister

func newMessageHandlerRegistry() *MessageHandlerRegistry {
	return &MessageHandlerRegistry{
		request: make(map[MessageType]MessageHandler),
	}
}

func (m *MessageHandlerRegistry) OnConncted(handler MessageHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onConnected = append(m.onConnected, handler)
}

func (m *MessageHandlerRegistry) OnDisconnected(handler MessageHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onDisconnected = append(m.onDisconnected, handler)
}

func (m *MessageHandlerRegistry) OnRequest(t MessageType, handler MessageHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.request[t] = handler
}
