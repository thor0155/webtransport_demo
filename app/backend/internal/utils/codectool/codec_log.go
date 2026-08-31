package codectool

import (
	"api/internal/model"
)

func EncodeLogPayload(level model.LogLevel, message string) ([]byte, error) {
	result := model.LogPayload{
		Level:   level,
		Message: message,
	}
	b, err := Encode(result)
	if err != nil {
		return nil, err
	}
	return b, nil
}
