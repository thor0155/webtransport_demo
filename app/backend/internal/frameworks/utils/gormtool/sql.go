package gormtool

import (
	"strings"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

const (
	BatchSize int = 1000
)

type ScopeFunc = func(q *gorm.DB) *gorm.DB

func InRange(column string, start, end any) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where(column+" BETWEEN ? AND ?", start, end)
	}
}

func IfInRange(column string, start, end any) ScopeFunc {
	if start != nil && end != nil {
		return InRange(column, start, end)
	} else if start != nil {
		return GEqThan(column, start)
	} else if end != nil {
		return LEqThan(column, end)
	} else {
		return func(q *gorm.DB) *gorm.DB {
			return q
		}
	}
}

func Clamp(column string, min, max any) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where(column+" > ? AND "+column+" < ?", min, max)
	}
}

func IfClamp(column string, min, max any) ScopeFunc {
	if min != nil && max != nil {
		return Clamp(column, min, max)
	} else if min != nil {
		return GreaterThan(column, min)
	} else if max != nil {
		return LessThan(column, max)
	} else {
		return func(q *gorm.DB) *gorm.DB {
			return q
		}
	}
}

func GEqThan(column string, value any) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where(column+" >= ?", value)
	}
}

func IfGEqThan(column string, value any) ScopeFunc {
	if value != nil {
		return GEqThan(column, value)
	} else {
		return func(q *gorm.DB) *gorm.DB {
			return q
		}
	}
}

func GreaterThan(column string, value any) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where(column+" > ?", value)
	}
}

func IfGreaterThan(column string, value any) ScopeFunc {
	if value != nil {
		return GreaterThan(column, value)
	} else {
		return func(q *gorm.DB) *gorm.DB {
			return q
		}
	}
}

func LEqThan(column string, value any) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where(column+" <= ?", value)
	}
}

func IfLEqThan(column string, value any) ScopeFunc {
	if value != nil {
		return LEqThan(column, value)
	} else {
		return func(q *gorm.DB) *gorm.DB {
			return q
		}
	}
}

func LessThan(column string, value any) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where(column+" < ?", value)
	}
}

func IfLessThan(column string, value any) ScopeFunc {
	if value != nil {
		return LessThan(column, value)
	} else {
		return func(q *gorm.DB) *gorm.DB {
			return q
		}
	}
}

func Equal(column string, value any) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where(column+" = ?", value)
	}
}

func IfEqual(column string, value any) ScopeFunc {
	if value != nil {
		return Equal(column, value)
	} else {
		return func(q *gorm.DB) *gorm.DB {
			return q
		}
	}
}

func In(column string, values any) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where(column+" IN ?", values)
	}
}

func IfIn(column string, values any) ScopeFunc {
	if values != nil {
		return In(column, values)
	} else {
		return func(q *gorm.DB) *gorm.DB {
			return q
		}
	}
}

func Join(table string, on string, args ...any) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Joins("JOIN "+table+" ON "+on, args...)
	}
}

func RightJoin(table string, on string, args ...any) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Joins("RIGHT JOIN "+table+" ON "+on, args...)
	}
}

func LeftJoin(table string, on string, args ...any) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Joins("LEFT JOIN "+table+" ON "+on, args...)
	}
}

func CreateTempTable(db *gorm.DB, table string, attrs ...string) error {
	if err := db.Exec("CREATE TEMPORARY TABLE IF NOT EXISTS " + table + ` (
	` + strings.Join(attrs, ",") + `
	)`).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func InsertSelect(db *gorm.DB,
	insertTable string, insertColumns []string,
	subQuery string, subQueryArgs []any, batchSize int) error {
	for page := 0; ; page++ {
		var args = subQueryArgs
		args = append(args, batchSize, page*batchSize)
		result := db.Exec("INSERT INTO "+insertTable+" ("+strings.Join(insertColumns, ",")+`)
		`+subQuery+`
		LIMIT ? OFFSET ?
		`, args...)
		if result.Error != nil {
			return errors.WithStack(result.Error)
		}
		if result.RowsAffected < int64(batchSize) {
			break
		}
	}
	return nil
}

func Truncate(db *gorm.DB, table string) error {
	return errors.WithStack(db.Exec("TRUNCATE ?", table).Error)
}
