package errorx

var (
	FailureError         = NewCodeError(ErrorCodeFailure, "failure")
	InvalidRequestError  = NewCodeError(ErrorCodeInvalidRequest, "invalid request")
	InvalidPasswordError = NewCodeError(ErrorCodeInvalidPassword, "invalid password")
	InvalidOTPError      = NewCodeError(ErrorCodeInvalidOTP, "invalid otp")
	InvalidTokenError    = NewCodeError(ErrorCodeInvalidToken, "invalid token")
)
