package model

import "api/internal/utils"

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

func EncodeLogPayload(level LogLevel, message string) ([]byte, error) {
	result := LogPayload{
		Level:   level,
		Message: message,
	}
	b, err := utils.Encode(result)
	if err != nil {
		return nil, err
	}
	return b, nil
}
