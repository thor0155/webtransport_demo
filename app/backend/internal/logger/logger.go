package logger

import (
	"go.uber.org/zap"
)

func New(cfg *LogConfig) (*zap.Logger, error) {
	zapcfg := zap.NewProductionConfig()
	if cfg.Level == "debug" {
		zapcfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	return zapcfg.Build()
}
