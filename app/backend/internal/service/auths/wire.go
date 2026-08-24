package auths

import "github.com/google/wire"

var ProvideSet = wire.NewSet(
	NewAuthService,
)
