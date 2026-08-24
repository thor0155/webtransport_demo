package wts

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewHub,
	NewServer,
	ProvideWebtransportServerConfig,
)
