package obj

import (
	"context"
	"sync"
)

type Work interface {
	Done() <-chan struct{}
}

type WorkTracker interface {
	// Track adds [works] to the current tracking batch.
	//
	// A work tracked after Wait() takes its snapshot belongs to
	// the next batch and does not affect that Wait().
	Track(works ...Work)

	// Wait waits until all works that were tracked when Wait() took
	// its snapshot are completed.
	//
	// [ctx] only controls how long the caller waits.
	// It does not cancel the tracked works.
	Wait(ctx context.Context) error
}

type workGroup struct {
	wg sync.WaitGroup
}

type workTracker struct {
	mu sync.Mutex

	// current contains works that belong to the current batch.
	current *workGroup
}

func (t *workTracker) Track(works ...Work) {

	if len(works) == 0 {
		return
	}

	t.mu.Lock()

	group := t.current

	// Add while holding the same mutex that protects the
	// Wait() snapshot boundary.
	group.wg.Add(1)

	t.mu.Unlock()

	go func() {
		defer group.wg.Done()
		for _, w := range works {
			<-w.Done()
		}
	}()

}

func (t *workTracker) Wait(ctx context.Context) error {
	t.mu.Lock()

	// Snapshot the current batch.
	group := t.current

	// Create a new batch for works tracked after this point.
	t.current = &workGroup{}

	t.mu.Unlock()

	done := make(chan struct{})

	go func() {
		group.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func NewWorkTracker() WorkTracker {
	return &workTracker{
		current: &workGroup{},
	}
}
