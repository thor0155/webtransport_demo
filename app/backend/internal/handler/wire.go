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
	NewLoginController,
)

func ProvideHttpRouters(c1 *adminController) https.RouterRegistries {
	return https.RouterRegistries{c1}
}

func ProvideWtsHandlers(c1 *chatController, c2 *loginController) wts.HandlerRegisters {
	return wts.HandlerRegisters{c1, c2}
}
