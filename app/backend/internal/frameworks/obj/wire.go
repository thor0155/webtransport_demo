package obj

import "github.com/google/wire"

var LifecycleWireSet = wire.NewSet(
	NewLifecycle,
	dependencyWireSet,
)

var LifecycleWithObjectsWireSet = wire.NewSet(
	NewLifecycleWithObjects,
	dependencyWireSet,
)

var dependencyWireSet = wire.NewSet(
	NewManager,
	NewWorkTracker,
	NewShutdownEmitter,
)
