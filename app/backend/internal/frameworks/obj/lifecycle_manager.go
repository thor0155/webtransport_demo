package obj

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type Manager struct {
	runners      []Runner
	pendingWorks []PendingWork
	servers      Servers
	logger       *zap.Logger
	workTracker  WorkTracker

	*ShutdownEmitter
}

func NewManager(workTracker WorkTracker, logger *zap.Logger, shutdownEmitter *ShutdownEmitter) *Manager {
	m := &Manager{
		logger:          logger,
		workTracker:     workTracker,
		ShutdownEmitter: shutdownEmitter,
	}
	return m
}

func (m *Manager) FindServers(filter Predicate[Server]) Servers {
	var results = make(Servers, 0, len(m.servers))
	for _, server := range m.servers {
		if filter(server) {
			results = append(results, server)
		}
	}
	return results
}

func (m *Manager) listenSignalsToShuttingdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	m.Shutdown(fmt.Sprintf("os signal: %s", sig))
}
