package obj

type WorkState struct {
	done chan struct{}
}

var _ Work = (*WorkState)(nil)

func (w *WorkState) Done() <-chan struct{} { return w.done }
func (w *WorkState) Close()                { close(w.done) }

func NewWorkState() *WorkState {
	return &WorkState{done: make(chan struct{})}
}
