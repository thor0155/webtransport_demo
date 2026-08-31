package db

import (
	"api/internal/frameworks/obj"

	"github.com/google/wire"
)

type Objects obj.Objects

var ProvideSet = wire.NewSet(
	NewRedisConfig,
	NewRedisDB,
	NewMysqlConfig,
	NewMysqlDB,
	NewMongoConfig,
	NewMongoDB,
	ProvideObjects,
)

func ProvideObjects(mysqlDB MysqlDB, mongoDB MongoDB, redisDB RedisDB) Objects {
	return Objects{mysqlDB, mongoDB, redisDB}
}
