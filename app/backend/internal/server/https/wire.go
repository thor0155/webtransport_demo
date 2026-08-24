package https

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewGinRouter,
	ConvGinEngineToHandler,
	ProvideHttpServerConfig,
	New,
)
