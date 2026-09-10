package migrator

import (
	"api/internal/frameworks/cmder"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type MigrateCmder struct {
	logger     *zap.Logger
	gormCmder  *GormMigrationCmder
	mongoCmder *MongoMigrationCmder
}

var (
	_ cmder.Cmder = (*MigrateCmder)(nil)
)

func NewMigrateCmder(logger *zap.Logger, gormCmder *GormMigrationCmder, mongoCmder *MongoMigrationCmder) *MigrateCmder {
	return &MigrateCmder{
		logger:     logger.Named("migrate_cmder"),
		gormCmder:  gormCmder,
		mongoCmder: mongoCmder,
	}
}

func (cmder *MigrateCmder) Command() string {
	return "migrate"
}

func (cmder *MigrateCmder) Description() string {
	return "migration features"
}

func (cmder *MigrateCmder) ProvideArgs() cobra.PositionalArgs {
	return cobra.NoArgs
}

func (cmder *MigrateCmder) FlagSet(f *pflag.FlagSet) {
}

func (cmder *MigrateCmder) Run(ctx cmder.CmderContext) {

	cmder.logger.Debug("run gorm migration cmder")
	cmder.gormCmder.Run(ctx)
	cmder.logger.Debug("run mongo migration cmder")
	cmder.mongoCmder.Run(ctx)
}
