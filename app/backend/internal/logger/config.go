package logger

import "api/internal/frameworks/config"

type LogConfig struct {
	Level string `mapstructure:"level"`
}

func ProvideLogConfig(helper config.ConfigHelper) (*LogConfig, error) {
	var r LogConfig
	err := helper.UnmarshalKey("log", &r)
	return &r, err
}
