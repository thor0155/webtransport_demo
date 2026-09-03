package errorx

import "fmt"

var (
	_ CodeError = (*errorWithCode)(nil)
	_ CodeError = (*codeError)(nil)
)

type CodeError interface {
	error
	Code() ErrorCode
}

type codeError struct {
	code ErrorCode
	msg  string
}

type errorWithCode struct {
	error
	code ErrorCode
}

// Code implements [CodeError].
func (e *errorWithCode) Code() ErrorCode {
	return e.code
}

func (e *errorWithCode) Format(s fmt.State, verb rune) {
	if f, ok := e.error.(fmt.Formatter); ok {
		f.Format(s, verb)
	} else {
		fmt.Fprintf(s, "%v", e.error)
	}
}

func WithCode(err error, code ErrorCode) error {
	return &errorWithCode{
		error: err,
		code:  code,
	}
}

func (c *codeError) Code() ErrorCode {
	return c.code
}

// Error implements [error].
func (c *codeError) Error() string {
	return c.msg
}

func NewCodeError(code ErrorCode, msg string) error {
	return &codeError{
		code: code,
		msg:  msg,
	}
}
