package wts

import (
	"context"
	"net/http"

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
		logger:                 logger.Named("handler"),
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
	})

	defer h.handleDisconnect(s)

	if err = h.handleConnect(s); err != nil {
		h.logger.Error("failed on connected and terminate", zap.Error(err))
		return
	}

	go s.WriteLoop(h.cancelCtx)
	s.ReadLoop(h.handleReadLoop(s))
}

func (h *Handler) handleReadLoop(s *Session) SessionReadHandlerFunc {
	return func(header *Header, b []byte) error {
		if heandlers, ok := h.messageHandlerRegistry.request[header.Type]; ok {
			ctx := NewHandlerContext(h.cancelCtx, h.hub, s, header, b, heandlers)
			ctx.Next()
			for _, err := range ctx.errors {
				if errors.HasType(err, (*errorTerminates)(nil)) {
					return err
				} else {
					h.logger.Error("handler error", zap.Uint8("request-type", header.Type), zap.Error(err))
				}
			}
			return nil
		}
		return errors.Newf("%v reader handler does not exist", header.Type)
	}
}

func (h *Handler) handleConnect(s *Session) error {
	header, data, err := Read(s.stream)
	if err != nil {
		return err
	}

	ctx := NewHandlerContext(h.cancelCtx, h.hub, s, header, data, h.messageHandlerRegistry.onConnected)
	ctx.Next()

	for _, err := range ctx.errors {
		if errors.HasType(err, (*errorTerminates)(nil)) {
			return err
		} else {
			h.logger.Error("connecting error", zap.Uint8("request-type", header.Type), zap.Error(err))
		}
	}
	return nil
}

func (h *Handler) handleDisconnect(s *Session) {
	defer s.Close()
	for _, f := range h.messageHandlerRegistry.onDisconnected {
		if err := f(h.cancelCtx, h.hub, s); err != nil {
			h.logger.Error("failed on disconnected", zap.Error(err))
		}
	}
}
