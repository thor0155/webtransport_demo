//go:build wireinject
// +build wireinject

package main

import (
	"api/cmd/migrate/migration"
	"api/internal/frameworks/cmder"
	"api/internal/frameworks/config"
	"api/internal/frameworks/db"
	"api/internal/frameworks/migrator"
	"api/internal/frameworks/obj"
	"api/internal/logger"

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
		cmder.ProvideSet,
		migrator.ProvideSet,
		migration.GetMigrations,
		db.ProvideSet,
		obj.LifecycleWithObjectsWireSet,
		obj.ProvideEmptyLifecycleOptions,
		provideObjects,
		provideGormProducer,
	)
	return nil, nil
}

func provideObjects(db db.Objects, cmder cmder.Objects) obj.Objects {
	var results obj.Objects
	results = append(results, cmder...)
	results = append(results, db...)
	return results
}

func provideGormProducer(mysql db.MysqlDB) migrator.GormProducerFunc {
	return mysql.Client
}
