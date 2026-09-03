package errorx

type ErrorCode = string

const (
	ErrorCodeFailure         ErrorCode = "Failure"
	ErrorCodeInvalidRequest  ErrorCode = "InvalidRequest"
	ErrorCodeUnauthorized    ErrorCode = "Unauthorized"
	ErrorCodeForbidden       ErrorCode = "Forbidden"
	ErrorCodeTokenExpired    ErrorCode = "TokenExpired"
	ErrorCodeInvalidPassword ErrorCode = "InvalidPassword"
	ErrorCodeInvalidOTP      ErrorCode = "InvalidOTP"
	ErrorCodeInvalidToken    ErrorCode = "InvalidToken"
)
