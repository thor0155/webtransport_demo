package entity

import "time"

const (
	TableNameUser string = "user"
)

type User struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"size:24;not null"`
	CreatedAt time.Time
}

func (*User) TableName() string {
	return TableNameUser
}
