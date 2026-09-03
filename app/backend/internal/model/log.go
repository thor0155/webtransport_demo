package model

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type LogCode = string

type LogResponse struct {
	Level   LogLevel `msgpack:"level"`
	Code    LogCode  `msgpack:"code"`
	Message string   `msgpack:"message"`
}
