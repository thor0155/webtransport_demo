package obj

import "sync"

type ShutdownEmitter struct {
	shutdownOnce sync.Once
	shutdownCh   chan string
}

func (s *ShutdownEmitter) Shutdown(reason string) {
	s.shutdownOnce.Do(func() {
		s.shutdownCh <- reason
	})
}

func NewShutdownEmitter() *ShutdownEmitter {
	return &ShutdownEmitter{
		shutdownCh: make(chan string, 1),
	}
}
