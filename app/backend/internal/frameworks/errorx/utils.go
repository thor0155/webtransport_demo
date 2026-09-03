package errorx

import "github.com/cockroachdb/errors"

func GetCode(err error) ErrorCode {
	codeError, ok := err.(CodeError)
	if !ok {
		ok = errors.As(err, &codeError)
	}
	if ok {
		return codeError.Code()
	}
	return ErrorCodeFailure
}
