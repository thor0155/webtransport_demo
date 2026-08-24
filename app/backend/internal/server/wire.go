package server

import (
	"api/internal/frameworks/obj"
	"api/internal/server/https"
	"api/internal/server/wts"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	https.ProvideSet,
	wts.ProvideSet,
	ProvideServerObjects,
)

type Objects obj.Objects

func ProvideServerObjects(s1 *https.HttpServer, s2 *wts.Server) Objects {
	return Objects{s1, s2}
}
