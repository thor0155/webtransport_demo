package https

import (
	"api/internal/frameworks/obj"
	"api/internal/frameworks/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type HttpServer struct {
	httpServer *http.Server
	logger     *zap.Logger
	cfg        *HttpServerConfig
	isRunning  bool
	name       string
}

var _ obj.Server = (*HttpServer)(nil)

func New(cfg *HttpServerConfig, handler http.Handler, logger *zap.Logger) *HttpServer {
	name := "http"
	s := &HttpServer{
		cfg:    cfg,
		name:   name,
		logger: logger.Named(name),
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
	s.logger.Info("server configured", utils.ObjectToZapFields(cfg)...)
	return s
}

func (s *HttpServer) IsRunning() bool {
	return s.isRunning
}

func (s *HttpServer) Start() {
	s.logger.Info("server started", zap.String("addr", s.httpServer.Addr))
	s.isRunning = true
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("server stopped", zap.Error(err))
	}
	s.isRunning = false
}

func (s *HttpServer) Close(ctx context.Context) error {
	defer func() {
		s.logger.Info("server shutting down")
		s.isRunning = false
	}()
	return s.httpServer.Shutdown(ctx)
}

func (s *HttpServer) Name() string {
	return s.name
}
