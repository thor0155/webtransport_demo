package service

import (
	"api/internal/service/auths"
	"api/internal/service/chats"
	"api/internal/service/users"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	auths.ProvideSet,
	chats.ProvideSet,
	users.ProvideSet,
)
