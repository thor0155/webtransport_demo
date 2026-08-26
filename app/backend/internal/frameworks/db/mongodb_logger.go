package db

import "go.uber.org/zap"

// ZapLogSink 包裝 zap.Logger，實作 mongo driver 的 LogSink 介面
type ZapLogSink struct {
	logger *zap.SugaredLogger
}

func NewZapLogSink(l *zap.Logger) *ZapLogSink {
	return &ZapLogSink{logger: l.Sugar()}
}

// Info 對應 driver 的一般日誌訊息
func (z *ZapLogSink) Info(level int, msg string, keysAndValues ...interface{}) {
	z.logger.Infow(msg, append([]interface{}{"driver_level", level}, keysAndValues...)...)
}

// Error 對應 driver 的錯誤日誌
func (z *ZapLogSink) Error(err error, msg string, keysAndValues ...interface{}) {
	z.logger.Errorw(msg, append([]interface{}{"error", err}, keysAndValues...)...)
}
