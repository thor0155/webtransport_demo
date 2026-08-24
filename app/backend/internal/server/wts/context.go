package wts

import (
	"context"
)

type Context interface {
	context.Context
	GetHub() *Hub
	GetSession() *Session
	GetPayload() []byte
}

type handlerContext struct {
	context.Context
	hub     *Hub
	session *Session
	payload []byte
}

// GetHub implements [Context].
func (h *handlerContext) GetHub() *Hub {
	return h.hub
}

// GetPayload implements [Context].
func (h *handlerContext) GetPayload() []byte {
	return h.payload
}

// GetSession implements [Context].
func (h *handlerContext) GetSession() *Session {
	return h.session
}

func NewHandlerContext(
	ctx context.Context, hub *Hub, session *Session, payload []byte) Context {

	return &handlerContext{
		Context: ctx,
		hub:     hub,
		session: session,
		payload: payload,
	}
}
