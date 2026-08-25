package migration

import (
	v202608150000 "api/cmd/migrate/migration/v202608150000"

	"github.com/go-gormigrate/gormigrate/v2"
	migrate "github.com/xakep666/mongo-migrate"
)

func GetGormMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		v202608150000.GormMigration,
	}
}

func GetMongoMigrations() []migrate.Migration {
	return []migrate.Migration{
		v202608150000.MongoMigration,
	}
}
