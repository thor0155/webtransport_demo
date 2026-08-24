//go:build wireinject
// +build wireinject

package main

import (
	"api/internal/dao"
	"api/internal/frameworks/config"
	"api/internal/frameworks/db"
	"api/internal/frameworks/obj"
	"api/internal/handler"
	"api/internal/logger"
	"api/internal/server"
	"api/internal/service"

	"github.com/google/wire"
	"go.uber.org/zap"
)

func GenLogger(configHelper config.ConfigHelper) (*zap.Logger, error) {
	wire.Build(
		logger.ProvideSet,
	)
	return nil, nil
}

func GenLifecycle(logger *zap.Logger, configHelper config.ConfigHelper) (*obj.Lifecycle, error) {
	wire.Build(
		handler.ProvideSet,
		server.ProvideSet,
		dao.ProvideSet,
		service.ProvideSet,
		db.ProvideSet,
		obj.LifecycleWithObjectsWireSet,
		obj.ProvideEmptyLifecycleOptions,
		provideObjects,
	)
	return nil, nil
}

func provideObjects(servers server.Objects, db db.Objects, dao dao.Objects) obj.Objects {
	var results obj.Objects
	results = append(results, servers...)
	results = append(results, db...)
	results = append(results, dao...)
	return results
}
