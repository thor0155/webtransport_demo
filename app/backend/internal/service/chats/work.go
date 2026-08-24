package chats

import "api/internal/frameworks/obj"

type asyncWork struct {
	name string
	done chan struct{}
}

var _ obj.PendingWork = (*asyncWork)(nil)

func (a *asyncWork) Name() string          { return a.name }
func (a *asyncWork) Done() <-chan struct{} { return a.done }

func newAsyncWork(name string) *asyncWork {
	return &asyncWork{name: name, done: make(chan struct{})}
}
