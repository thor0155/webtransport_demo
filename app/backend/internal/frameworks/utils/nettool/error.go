package nettool

import (
	"context"
	"errors"
	"io"
	"net"

	"github.com/quic-go/quic-go"
	webtransport "github.com/quic-go/webtransport-go"
)

func IsExpectedDisconnect(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) {
		return true
	}

	if errors.Is(err, io.EOF) {
		return true
	}

	if errors.Is(err, net.ErrClosed) {
		return true
	}

	var sessionErr *webtransport.SessionError
	if errors.As(err, &sessionErr) {

		// Client 正常關閉 Session
		if sessionErr.Remote && sessionErr.ErrorCode == 0 {
			return true
		}

		// Server 自己關閉 Session
		if !sessionErr.Remote && sessionErr.ErrorCode == 0 {
			return true
		}
	}

	var appErr *quic.ApplicationError
	if errors.As(err, &appErr) {
		if appErr.ErrorCode == 0 {
			return true
		}
	}

	var idleErr *quic.IdleTimeoutError
	if errors.As(err, &idleErr) {
		return true
	}

	var statelessReset *quic.StatelessResetError
	if errors.As(err, &statelessReset) {
		return true
	}

	return false
}
