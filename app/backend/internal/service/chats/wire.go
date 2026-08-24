package chats

import "github.com/google/wire"

var ProvideSet = wire.NewSet(
	NewChatService,
)
