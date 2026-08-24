package obj

import (
	"context"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/panjf2000/ants/v2"
)

var (
	ErrTrackUnavailable = errors.New("track unavailable")
)

type WorkTracker interface {
	Track(w PendingWork) error
	Wait(ctx context.Context) error
}

type workTracker struct {
	pendingWorkCoroutine *ants.PoolWithFuncGeneric[PendingWork]
	waittingCoroutine    *ants.PoolWithFuncGeneric[*waittingArg]
	wg                   sync.WaitGroup
}

type waittingArg struct {
	done chan struct{}
}

func NewWorkTracker() (WorkTracker, error) {
	t := &workTracker{}
	var err error
	if t.pendingWorkCoroutine, err = ants.NewPoolWithFuncGeneric(MaxProcessors, t.processPendingWork); err != nil {
		return nil, err
	}
	if t.waittingCoroutine, err = ants.NewPoolWithFuncGeneric(1, t.processWaitting, ants.WithPreAlloc(true)); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *workTracker) Track(w PendingWork) error {
	if t.pendingWorkCoroutine.Free() <= 0 {
		return ErrTrackUnavailable
	}

	t.wg.Add(1)
	if err := t.pendingWorkCoroutine.Invoke(w); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (t *workTracker) Wait(ctx context.Context) error {

	arg := &waittingArg{
		done: make(chan struct{}),
	}

	if err := t.waittingCoroutine.Invoke(arg); err != nil {
		return errors.WithStack(err)
	}

	select {
	case <-arg.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *workTracker) processPendingWork(w PendingWork) {
	defer t.wg.Done()
	<-w.Done()
}

func (t *workTracker) processWaitting(arg *waittingArg) {
	t.wg.Wait()
	close(arg.done)
}
