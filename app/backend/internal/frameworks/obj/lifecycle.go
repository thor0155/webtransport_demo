package obj

import (
	"api/internal/frameworks/coroutine"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

var (
	ErrAlreadyRunning      = errors.New("lifecycle is already running")
	DefaultShutdownTimeout = 60 * time.Second
)

type Lifecycle struct {
	logger               *zap.Logger
	shutdownTimeout      time.Duration
	isRunning            atomic.Bool
	componentCoroutine   *coroutine.SequenceCoroutine[Component]
	closeableCoroutine   *coroutine.SequenceCoroutine[Closeable]
	closeServerCoroutine *coroutine.SequenceCoroutine[Server]
	*Manager
}

func ProvideEmptyLifecycleOptions() []LifecycleOption {
	return nil
}

func NewLifecycleWithObjects(logger *zap.Logger, manager *Manager, objects Objects, opts ...LifecycleOption) (*Lifecycle, error) {
	if l, err := NewLifecycle(logger, manager, opts...); err != nil {
		return nil, err
	} else {
		if err = l.Register(objects...); err != nil {
			return nil, err
		}
		return l, nil
	}
}

func NewLifecycle(logger *zap.Logger, manager *Manager, opts ...LifecycleOption) (*Lifecycle, error) {
	l := &Lifecycle{
		logger:          logger.Named("lifecycle"),
		shutdownTimeout: DefaultShutdownTimeout,
		Manager:         manager,
	}
	var err error
	if l.componentCoroutine, err = coroutine.NewSequenceCoroutine(func(c Component) error {
		l.logger.Debug("initialize component", zap.String("name", c.Name()))
		return c.Init()
	}); err != nil {
		return nil, err
	}
	if l.closeableCoroutine, err = coroutine.NewSequenceCoroutine(func(c Closeable) error {
		l.logger.Debug("close", zap.String("name", c.Name()))
		c.Close(context.Background())
		return nil
	}); err != nil {
		return nil, err
	}
	if l.closeServerCoroutine, err = coroutine.NewSequenceCoroutine(func(s Server) error {
		l.logger.Debug("server close", zap.String("name", s.Name()))
		s.Close(context.Background())
		return nil
	}); err != nil {
		return nil, err
	}
	setOption(l, opts)
	return l, nil
}

func (l *Lifecycle) Register(objs ...Object) error {
	if l.isRunning.Load() {
		return errors.Wrap(ErrAlreadyRunning, "Unable to register")
	}

	for _, o := range objs {

		if a, ok := o.(Activity); ok {
			if !a.Active() {
				continue
			}
		}

		matched := false
		phase := phaseOf(o)
		server, isServer := o.(Server)

		if isServer {
			l.servers = append(l.servers, server)
			l.closeServerCoroutine.Add(phase, server)
			matched = true
		}
		if runner, ok := o.(Runner); ok {
			l.runners = append(l.runners, runner)
			matched = true
		}
		if closeable, ok := o.(Closeable); ok && !isServer {
			l.closeableCoroutine.Add(phase, closeable)
			matched = true
		}
		if pendingWork, ok := o.(PendingWork); ok {
			l.pendingWorks = append(l.pendingWorks, pendingWork)
			matched = true
		}
		if component, ok := o.(Component); ok {
			l.componentCoroutine.Add(phase, component)
			matched = true
		}
		if !matched {
			l.logger.Warn("object does not implement any known lifecycle interface, ignored",
				zap.String("name", o.Name()))
		}
	}
	return nil
}

// Run 啟動所有 Server,並 block 直到 shutdown 完成。應該只被呼叫一次。
func (l *Lifecycle) Run() error {

	if !l.isRunning.CompareAndSwap(false, true) {
		return errors.WithStack(ErrAlreadyRunning)
	}
	defer l.isRunning.Store(false)

	// 1. init all components
	if err := l.componentCoroutine.Invoke(context.Background()); err != nil {
		return errors.Wrap(err, "init component failed")
	}

	// 2. start goroutine with all runner
	var wg sync.WaitGroup
	l.handleRunners(&wg)

	// 3. listen sys signal to shutting down
	go l.listenSignalsToShuttingdown()

	reason := <-l.shutdownCh
	l.logger.Info("shutdown triggered", zap.String("reason", reason))

	shutdownTimeoutCtx, shutdownCancelFunc := context.WithTimeout(context.Background(), l.shutdownTimeout)
	defer shutdownCancelFunc()

	// 4. close all servers
	if err := l.closeServerCoroutine.Invoke(shutdownTimeoutCtx); err != nil && errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	// 5. wait pending works
	if err := l.handlePendingWorks(shutdownTimeoutCtx); err != nil && errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	// 6. close all closeables
	if err := l.closeableCoroutine.Invoke(shutdownTimeoutCtx); err != nil && errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	wg.Wait()
	l.logger.Info("lifecycle shutdown complete")
	return nil
}

func (l *Lifecycle) handleRunners(wg *sync.WaitGroup) {
	for _, runner := range l.runners {
		wg.Add(1)
		go func(r Runner) {
			defer wg.Done()
			l.logger.Debug("starting", zap.String("name", r.Name()))
			r.Start()
			l.logger.Debug("stopped", zap.String("name", r.Name()))
		}(runner)
	}
}

func (l *Lifecycle) handlePendingWorks(ctx context.Context) error {
	if len(l.pendingWorks) > 0 {
		l.logger.Debug("waiting for pending work to finish", zap.Int("count", len(l.pendingWorks)))

		for _, work := range l.pendingWorks {
			if err := l.workTracker.Track(work); err != nil {
				return err
			}
		}
	}
	if err := l.workTracker.Wait(ctx); err != nil {
		return err
	}
	return nil
}
