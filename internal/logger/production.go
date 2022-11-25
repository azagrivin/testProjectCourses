package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     SyslogTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func NewProduction() *zap.Logger {
	return newLogger(
		zapcore.NewJSONEncoder(newProductionEncoderConfig()),
		zap.InfoLevel,
	)
}
