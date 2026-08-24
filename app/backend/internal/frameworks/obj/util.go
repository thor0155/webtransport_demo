package obj

import (
	"runtime"
)

var MaxProcessors = runtime.NumCPU()

type Predicate[T Object] func(T) bool

// type groupedInfo[T Object] struct {
// 	phase   []int
// 	grouped map[int][]T
// }

// func initGroupedInfo[T Object](g *groupedInfo[T]) {
// 	g.phase = g.phase[:0]
// 	g.grouped = make(map[int][]T)
// }

func phaseOf(o Object) int {
	if p, ok := o.(Phaseable); ok {
		return p.Phase()
	}
	return 0
}

// func groupObjectByPhaseDesc[T Object](objs []T) ([]int, map[int][]T) {
// 	grouped := make(map[int][]T)
// 	for _, o := range objs {
// 		p := phaseOf(o)
// 		grouped[p] = append(grouped[p], o)
// 	}
// 	phases := slices.Collect(maps.Keys(grouped))
// 	sort.Sort(sort.Reverse(sort.IntSlice(phases)))
// 	return phases, grouped
// }
