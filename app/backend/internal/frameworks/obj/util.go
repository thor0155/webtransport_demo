package obj

type Predicate[T Object] func(T) bool

func phaseOf(o Object) int {
	if p, ok := o.(Phaseable); ok {
		return p.Phase()
	}
	return 0
}
