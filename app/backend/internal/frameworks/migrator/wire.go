package migrator

import (
	"api/internal/frameworks/cmder"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewMigrateCmder,
	NewGormMigrationCmder,
	NewMongoMigrationCmder,
	ProvideCmders,
)

func ProvideCmders(root *MigrateCmder, gorm *GormMigrationCmder, mongo *MongoMigrationCmder) cmder.Cmders {
	return cmder.Cmders{root, gorm, mongo}
}
