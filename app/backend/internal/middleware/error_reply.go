package middleware

import (
	"api/internal/model"
	"api/internal/protocol"
	"api/internal/server/wts"
)

func WebtransportErrorReply(ctx wts.Context) error {

	ctx.Next()
	session := ctx.GetSession()
	for _, err := range ctx.GetErrors() {
		if b, err := model.EncodeLogPayload(model.LogLevelError, err.Error()); err != nil {
			return err
		} else {
			if err := session.SyncSendMessage(protocol.ResponseTypeLog, b); err != nil {
				return err
			}
		}
	}
	return nil
}
