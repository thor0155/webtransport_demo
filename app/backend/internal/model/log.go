package model

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type LogPayload struct {
	Level   LogLevel `msgpack:"level"`
	Message string   `msgpack:"message"`
}
