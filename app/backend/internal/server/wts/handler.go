package wts

import (
	"api/internal/model"
	"context"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/quic-go/webtransport-go"
	"go.uber.org/zap"
)

type Handler struct {
	wt                     *webtransport.Server
	hub                    *Hub
	logger                 *zap.Logger
	cancelCtx              context.Context
	cancelFunc             context.CancelFunc
	messageHandlerRegistry *MessageHandlerRegistry
}

func NewHandler(wt *webtransport.Server, hub *Hub, logger *zap.Logger, messageHandlerRegistry *MessageHandlerRegistry) *Handler {
	h := &Handler{
		wt:                     wt,
		hub:                    hub,
		logger:                 logger,
		messageHandlerRegistry: messageHandlerRegistry,
	}
	h.cancelCtx, h.cancelFunc = context.WithCancel(context.Background())
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	session, err := h.wt.Upgrade(w, r)
	if err != nil {
		h.logger.Error("failed to upgrade session", zap.Error(err))
		return
	}

	s, err := NewSession(
		uuid.NewString(),
		session,
		h.cancelCtx,
		h.logger,
	)
	if err != nil {
		session.CloseWithError(ClosedErrorCode, "failed to new session")
		h.logger.Error("new session failed", zap.Error(err))
		return
	}
	h.logger.Debug("session created", zap.String("session_id", s.id))

	h.hub.Register(s)
	s.AddClosedEvent(func() {
		h.hub.Unregister(s)
		h.logger.Debug("session disconnected", zap.String("session_id", s.id))
	})

	defer h.handleDisconnect(s)

	if err = h.handleConnect(s); err != nil {
		if b, err := model.EncodeLogPayload(model.LogLevelError, err.Error()); err == nil {
			if err := s.SendMessage(ResponseTypeLog, b); err == nil {
				ctx, cancelFunc := context.WithTimeout(h.cancelCtx, time.Second)
				defer cancelFunc()
				s.WriteLoop(ctx)
			}
		}

		h.logger.Error("failed on connected", zap.Error(err))
		return
	}

	go s.WriteLoop(h.cancelCtx)
	s.ReadLoop(h.handleReadLoop(s))
}

func (h *Handler) handleReadLoop(s *Session) func(MessageType, []byte) error {
	return func(t MessageType, b []byte) error {
		if f, ok := h.messageHandlerRegistry.request[t]; ok {
			return f(NewHandlerContext(h.cancelCtx, h.hub, s, b))
		}
		return errors.Newf("%v handler does not exist", t)
	}
}

func (h *Handler) handleConnect(s *Session) error {
	header, data, err := Read(s.stream)
	if err != nil {
		return err
	}

	if header.Type != RequestTypeHello {
		return errors.New("first message must be hello")
	}

	ctx := NewHandlerContext(h.cancelCtx, h.hub, s, data)
	for _, f := range h.messageHandlerRegistry.onConnected {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) handleDisconnect(s *Session) {

	defer s.Close()
	ctx := NewHandlerContext(h.cancelCtx, h.hub, s, nil)

	for _, f := range h.messageHandlerRegistry.onDisconnected {
		if err := f(ctx); err != nil {
			h.logger.Error("failed on disconnected", zap.Error(err))
		}
	}
}
