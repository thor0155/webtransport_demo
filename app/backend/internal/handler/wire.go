package handler

import (
	"api/internal/server/https"
	"api/internal/server/wts"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	ProvideHttpRouters,
	ProvideWtsHandlers,
	NewAdminController,
	NewChatController,
)

func ProvideHttpRouters(ctrl1 *adminController) https.RouterRegistries {
	return https.RouterRegistries{ctrl1}
}

func ProvideWtsHandlers(c1 *chatController) wts.HandlerRegisters {
	return wts.HandlerRegisters{c1}
}
