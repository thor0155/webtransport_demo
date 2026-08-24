package migrator

import (
	"api/internal/frameworks/cmder"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewGormMigrationCmder,
	ProvideCmders,
)

func ProvideCmders(gorm *GormMigrationCmder) cmder.Cmders {
	return cmder.Cmders{gorm}
}
