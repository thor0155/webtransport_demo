package migrator

import (
	"api/internal/frameworks/cmder"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type MongoMigrationCmder struct {
	logger                *zap.Logger
	mongoDatabaseProducer MongoDatabaseProducerFunc
	rollback              bool
	migrations            []migrate.Migration
}

type MongoDatabaseProducerFunc func() *mongo.Database

var (
	_ cmder.Cmder = (*MongoMigrationCmder)(nil)
)

// Command implements [cmder.Cmder].
func (m *MongoMigrationCmder) Command() string {
	return "migrate mongodb"
}

// Description implements [cmder.Cmder].
func (m *MongoMigrationCmder) Description() string {
	return "mongodb migration features"
}

// FlagSet implements [cmder.Cmder].
func (m *MongoMigrationCmder) FlagSet(f *pflag.FlagSet) {
	f.BoolVarP(&m.rollback, "rollback", "r", false, "")
}

// ProvideArgs implements [cmder.Cmder].
func (m *MongoMigrationCmder) ProvideArgs() cobra.PositionalArgs {
	return cobra.MaximumNArgs(1)
}

// Run implements [cmder.Cmder].
func (m *MongoMigrationCmder) Run(ctx cmder.CmderContext) {

	var (
		database = m.mongoDatabaseProducer()
		args     = ctx.Args()
		logger   = m.logger
		err      error
		version  uint64
	)
	if database == nil {
		ctx.Error(errors.New("database is nil"))
		return
	}
	if len(m.migrations) == 0 {
		ctx.Error(errors.New("there is no migrations"))
		return
	}
	logger.Info("migration list", zap.Int("count", len(m.migrations)))

	migrator := migrate.NewMigrate(database, m.migrations...)

	if len(args) > 0 {
		if version, err = strconv.ParseUint(args[0], 10, 64); err != nil {
			ctx.Error(err)
			return
		}
		logger = logger.With(zap.Uint64("version", version))
	}

	if version > 0 {
		if err = migrator.SetVersion(ctx, version, "specify a version"); err != nil {
			ctx.Error(err)
			return
		}
	}

	if m.rollback {
		logger.Debug("rollback")
		err = migrator.Down(ctx, migrate.AllAvailable)

	} else {
		logger.Debug("up")
		err = migrator.Up(ctx, migrate.AllAvailable)
	}
	if err != nil {
		ctx.Error(err)
	}
}

func NewMongoMigrationCmder(logger *zap.Logger, mongoDatabaseProducer MongoDatabaseProducerFunc, migrations []migrate.Migration) *MongoMigrationCmder {
	return &MongoMigrationCmder{
		logger:                logger.Named("mongo_migration_cmder"),
		mongoDatabaseProducer: mongoDatabaseProducer,
		migrations:            migrations,
	}
}
