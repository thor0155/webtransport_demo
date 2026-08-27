package wts

import (
	"context"
	"math"
)

const (
	abortIndex int8 = math.MaxInt8 >> 1
)

type Context interface {
	context.Context
	GetHub() *Hub
	GetSession() *Session
	GetPayload() []byte
	GetType() MessageType
	GetErrors() []error
	// Next should be used only inside middleware.
	// It executes the pending handlers in the chain inside the calling handler.
	Next()
	// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
	// Let's say you have an authorization middleware that validates that the current request is authorized.
	// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
	// for this request are not called.
	Abort()
}

type handlerContext struct {
	context.Context
	hub      *Hub
	session  *Session
	payload  []byte
	header   *Header
	index    int8
	handlers []MessageHandler
	errors   []error
}

func (h *handlerContext) GetErrors() []error {
	return h.errors
}

// Abort implements [Context].
func (h *handlerContext) Abort() {
	h.index = abortIndex
}

// Next implements [Context].
func (h *handlerContext) Next() {
	h.index++
	handlerCount := safeInt8(len(h.handlers))
	for h.index < handlerCount {
		if err := h.handlers[h.index](h); err != nil {
			h.errors = append(h.errors, err)
		}
		h.index++
	}
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

func (h *handlerContext) GetType() MessageType {
	return h.header.Type
}

func NewHandlerContext(
	ctx context.Context, hub *Hub, session *Session, header *Header, payload []byte, handlers []MessageHandler) *handlerContext {

	return &handlerContext{
		Context:  ctx,
		hub:      hub,
		session:  session,
		payload:  payload,
		header:   header,
		handlers: handlers,
		index:    -1,
	}
}
