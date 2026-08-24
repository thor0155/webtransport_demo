package cmder

import (
	"api/internal/frameworks/obj"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewCmderConfig,
	NewCmderServer,
	ProvideObjects,
)

type Objects obj.Objects

func ProvideObjects(cmderServer *Server) Objects {
	return Objects{cmderServer}
}
