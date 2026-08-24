package obj

import (
	"context"
)

type Object interface {
	Name() string
}
type Objects = []Object

// PendingWork represents an asynchronous unit of work that needs to wait for completion in its lifecycle.
// It does not describe "how to execute", but only "when it is considered complete".
type PendingWork interface {
	Object
	Done() <-chan struct{}
}

type Component interface {
	Object
	Activity
	Init() error
}

type Runner interface {
	Object
	Start()
}

type Closeable interface {
	Object
	Close(ctx context.Context) error
}

type Server interface {
	Runner
	Closeable
}

type Servers = []Server

type Activity interface {
	Object
	Active() bool
}

type Phaseable interface {
	Phase() int
}
