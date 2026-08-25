package migrator

import (
	"api/internal/frameworks/cmder"

	"github.com/cockroachdb/errors"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormMigrationCmder struct {
	logger           *zap.Logger
	gormProducerFunc GormProducerFunc
	rollback         bool
	migrations       []*gormigrate.Migration
}

type GormProducerFunc func() *gorm.DB

var (
	_ cmder.Cmder = (*GormMigrationCmder)(nil)
)

func NewGormMigrationCmder(logger *zap.Logger, gormProducerFunc GormProducerFunc, migrations []*gormigrate.Migration) *GormMigrationCmder {
	return &GormMigrationCmder{
		logger:           logger.Named("gorm_migration_cmder"),
		gormProducerFunc: gormProducerFunc,
		migrations:       migrations,
	}
}

func (cmder *GormMigrationCmder) Command() string {
	return "migrate gorm"
}

func (cmder *GormMigrationCmder) Description() string {
	return "gorm migration features"
}

func (cmder *GormMigrationCmder) ProvideArgs() cobra.PositionalArgs {
	return cobra.MaximumNArgs(1)
}

func (cmder *GormMigrationCmder) FlagSet(f *pflag.FlagSet) {
	f.BoolVarP(&cmder.rollback, "rollback", "r", false, "")
}

func (cmder *GormMigrationCmder) Run(ctx cmder.CmderContext) {
	gorm := cmder.gormProducerFunc()
	if gorm == nil {
		ctx.Error(errors.New("gorm is nil"))
		return
	}
	if len(cmder.migrations) == 0 {
		ctx.Error(errors.New("there is no migrations"))
		return
	}
	migrate := gormigrate.New(gorm, gormigrate.DefaultOptions, cmder.migrations)
	args := ctx.Args()
	logger := cmder.logger

	var version string
	if len(args) > 0 {
		version = args[0]
		logger = logger.With(zap.String("version", version))
	}

	var err error
	if cmder.rollback {
		logger.Debug("rollback")
		err = cmder.migrateRollback(migrate, version)
	} else {
		logger.Debug("up")
		err = cmder.migrateUp(migrate, version)
	}
	if err != nil {
		ctx.Error(err)
	}
}

func (*GormMigrationCmder) migrateUp(m *gormigrate.Gormigrate, version string) error {
	if version != "" {
		return m.MigrateTo(version)
	} else {
		return m.Migrate()
	}
}

func (*GormMigrationCmder) migrateRollback(m *gormigrate.Gormigrate, version string) error {
	if version != "" {
		return m.RollbackTo(version)
	} else {
		return m.RollbackLast()
	}
}
