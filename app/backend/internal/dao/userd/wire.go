package userd

import "github.com/google/wire"

var ProvideSet = wire.NewSet(
	NewUserDao,
)
