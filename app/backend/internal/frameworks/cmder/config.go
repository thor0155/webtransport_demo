package cmder

import "api/internal/frameworks/config"

type CmderConfig struct {
	EnableReadEvalPrintLoop bool `mapstructure:"enableReadEvalPrintLoop"`
}

func NewCmderConfig(configHelper config.ConfigHelper) (*CmderConfig, error) {
	var cfg CmderConfig
	if err := configHelper.UnmarshalKey("servers.cmder", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
