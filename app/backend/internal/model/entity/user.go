package entity

import "time"

const (
	TableNameUser string = "user"
)

type User struct {
	Id        string `gorm:"primaryKey"`
	Name      string `gorm:"size:24;not null"`
	CreatedAt time.Time
}

func (*User) TableName() string {
	return TableNameUser
}
