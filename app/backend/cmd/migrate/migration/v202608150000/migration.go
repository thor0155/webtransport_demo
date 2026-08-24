package v202608150000

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var Migration *gormigrate.Migration = &gormigrate.Migration{
	ID: "202608150000",
	Migrate: func(tx *gorm.DB) error {
		var err = tx.AutoMigrate(
			User{}, ChatRoom{}, ChatRoomUsers{})
		if err != nil {
			return err
		}

		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable(
			User{}, ChatRoom{}, ChatRoomUsers{},
		)
	},
}
