package migrator

import (
	"api/internal/frameworks/cmder"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewGormMigrationCmder,
	NewMongoMigrationCmder,
	ProvideCmders,
)

func ProvideCmders(gorm *GormMigrationCmder, mongo *MongoMigrationCmder) cmder.Cmders {
	return cmder.Cmders{gorm, mongo}
}
