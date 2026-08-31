package users

import "github.com/google/wire"

var ProvideSet = wire.NewSet(
	NewRedisNameResolver,
)
