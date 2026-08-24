package gormtool

import "gorm.io/gorm/clause"

var (
	LockForUpdate = clause.Locking{Strength: "UPDATE"}
	InsertIgnore  = clause.Insert{Modifier: "IGNORE"}
)
