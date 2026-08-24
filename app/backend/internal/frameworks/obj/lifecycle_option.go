package obj

import "time"

type LifecycleOption func(*Lifecycle)

func WithShutdownTimeout(d time.Duration) LifecycleOption {
	return func(l *Lifecycle) { l.shutdownTimeout = d }
}

func setOption(l *Lifecycle, opts []LifecycleOption) {
	for _, opt := range opts {
		opt(l)
	}
}
