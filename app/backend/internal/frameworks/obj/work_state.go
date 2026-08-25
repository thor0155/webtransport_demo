package obj

type WorkState struct {
	name string
	done chan struct{}
}

var _ PendingWork = (*WorkState)(nil)

func (w *WorkState) Name() string          { return w.name }
func (w *WorkState) Done() <-chan struct{} { return w.done }
func (w *WorkState) Close()                { close(w.done) }

func NewWorkState(name string) *WorkState {
	return &WorkState{name: name, done: make(chan struct{})}
}
