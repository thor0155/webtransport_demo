package codectool

import (
	"api/internal/model"
)

func EncodeLogPayload(level model.LogLevel, code model.LogCode, message string) ([]byte, error) {
	result := model.LogResponse{
		Level:   level,
		Code:    code,
		Message: message,
	}
	b, err := Encode(result)
	if err != nil {
		return nil, err
	}
	return b, nil
}
