package logger

import "github.com/google/wire"

var ProvideSet = wire.NewSet(
	ProvideLogConfig,
	New,
)
