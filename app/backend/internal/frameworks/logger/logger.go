package logger

import (
	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

func New(cfg *LogConfig) (*zap.Logger, error) {
	zapcfg := zap.NewProductionConfig()
	switch cfg.Level {
	case "info":
		zapcfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapcfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapcfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "debug":
		zapcfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		return nil, errors.Newf("unknown log level: %s", cfg.Level)
	}
	return zapcfg.Build()
}
