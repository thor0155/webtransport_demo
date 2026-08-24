package dao

import (
	"api/internal/dao/chatd"
	"api/internal/dao/userd"
	"api/internal/frameworks/obj"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	chatd.ProvideSet,
	userd.ProvideSet,
	ProvideObjects,
)

type Objects obj.Objects

func ProvideObjects(chat chatd.Objects) Objects {
	results := append(Objects{}, chat...)
	return results
}
