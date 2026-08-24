package db

import (
	"api/internal/frameworks/obj"

	"github.com/google/wire"
)

type Objects obj.Objects

var ProvideSet = wire.NewSet(
	NewMysqlConfig,
	NewMysqlDB,
	NewMongoConfig,
	NewMongoDB,
	ProvideObjects,
)

func ProvideObjects(mysqlDB MysqlDB, mongoDB MongoDB) Objects {
	return Objects{mysqlDB, mongoDB}
}
