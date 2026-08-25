package entity

const (
	TableNameChatRoom      string = "chat_room"
	TableNameChatRoomUsers string = "chat_room_users"
)

type ChatRoom struct {
	Id   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;size:32;not null"`
}

type ChatRoomUsers struct {
	ChatRoomId uint   `gorm:"primaryKey"`
	UserId     string `gorm:"primaryKey"`
}

func (*ChatRoom) TableName() string {
	return TableNameChatRoom
}

func (*ChatRoomUsers) TableName() string {
	return TableNameChatRoomUsers
}
