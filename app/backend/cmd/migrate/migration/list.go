package migration

import (
	v202608150000 "api/cmd/migrate/migration/v202608150000"

	"github.com/go-gormigrate/gormigrate/v2"
)

func GetMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		v202608150000.Migration,
	}
}
