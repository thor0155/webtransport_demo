package v202608150000

import "time"

type User struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"size:24;not null"`
	CreatedAt time.Time
}

func (*User) TableName() string {
	return "user"
}

type ChatRoom struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:32;uniqueIndex;not null"`
}

type ChatRoomUsers struct {
	ChatRoomId uint   `gorm:"primaryKey"`
	UserId     string `gorm:"primaryKey"`
}

func (*ChatRoom) TableName() string {
	return "chat_room"
}

func (*ChatRoomUsers) TableName() string {
	return "chat_room_users"
}
