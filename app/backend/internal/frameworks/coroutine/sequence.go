package coroutine

import (
	"context"
	"runtime"

	"maps"
	"slices"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/cockroachdb/errors"
	"github.com/panjf2000/ants/v2"
	"github.com/tcard/coro/v2"
)

var (
	ErrAlreadyRunning = errors.New("SequenceCoroutine is already running")
	MaxProcessors     = runtime.NumCPU()
)

type PhaseOfFunc[T any] func(o T) int

type HandlerFunc[T any] func(o T) error

type SequenceCoroutine[T any] struct {
	mu              sync.RWMutex
	grouped         map[int][]T
	wg              sync.WaitGroup
	errCh           chan error
	handler         HandlerFunc[T]
	isRunning       atomic.Bool
	bundleCoroutine *ants.PoolWithFuncGeneric[*bundleProcArg[T]]
	taskCoroutine   *ants.PoolWithFuncGeneric[*taskProcArg[T]]
}

type bundleProcArg[T any] struct {
	done  chan struct{}
	tasks []T
}

type taskProcArg[T any] struct {
	wg   *sync.WaitGroup
	task T
}

func NewSequenceCoroutine[T any](handler HandlerFunc[T]) (*SequenceCoroutine[T], error) {
	var err error
	r := &SequenceCoroutine[T]{
		handler: handler,
		grouped: make(map[int][]T),
		errCh:   make(chan error, 1),
	}
	if r.bundleCoroutine, err = ants.NewPoolWithFuncGeneric(1, r.processBundle, ants.WithPreAlloc(true)); err != nil {
		return nil, err
	}
	if r.taskCoroutine, err = ants.NewPoolWithFuncGeneric(MaxProcessors, r.processTask); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *SequenceCoroutine[T]) Add(phase int, o T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.grouped[phase] = append(s.grouped[phase], o)
}

func (s *SequenceCoroutine[T]) AddByPhaseFunc(o T, phaseOfFunc PhaseOfFunc[T]) {
	s.Add(phaseOfFunc(o), o)
}

func (s *SequenceCoroutine[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.grouped = make(map[int][]T)
}

func (s *SequenceCoroutine[T]) Invoke(ctx context.Context) error {

	if !s.isRunning.CompareAndSwap(false, true) {
		return ErrAlreadyRunning
	}
	defer s.isRunning.Store(false)

	next := coro.Generate(context.Background(), s.yieldBundleHandler)

	var err error
	var bundle []T
	for next(&err, &bundle) {
		arg := &bundleProcArg[T]{
			tasks: bundle,
			done:  make(chan struct{}),
		}
		if err := s.bundleCoroutine.Invoke(arg); err != nil {
			close(arg.done)
			return err
		}
		select {
		case err := <-s.errCh:
			return err
		case <-arg.done:
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}

func (s *SequenceCoroutine[T]) yieldBundleHandler(yield func([]T)) error {
	s.mu.RLock()
	grouped := maps.Clone(s.grouped)
	s.mu.RUnlock()

	phases := slices.Collect(maps.Keys(grouped))
	sort.Sort(sort.Reverse(sort.IntSlice(phases)))
	for _, phase := range phases {
		yield(grouped[phase])
	}
	return nil
}

func (s *SequenceCoroutine[T]) processBundle(arg *bundleProcArg[T]) {
	defer close(arg.done)
	for _, task := range arg.tasks {
		s.wg.Add(1)
		if err := s.taskCoroutine.Invoke(&taskProcArg[T]{
			wg:   &s.wg,
			task: task,
		}); err != nil {
			s.wg.Done()
			s.errCh <- err
			return
		}
	}
	s.wg.Wait()
}

func (s *SequenceCoroutine[T]) processTask(arg *taskProcArg[T]) {
	defer arg.wg.Done()
	if err := s.handler(arg.task); err != nil {
		s.errCh <- err
	}
}
