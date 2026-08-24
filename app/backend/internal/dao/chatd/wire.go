package chatd

import (
	"api/internal/frameworks/obj"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewChatMessageDao,
	NewChatRoomDao,
	ProvideObjects,
)

type Objects obj.Objects

func ProvideObjects(chatMessageDao ChatMessageDao) Objects {
	return Objects{chatMessageDao}
}
