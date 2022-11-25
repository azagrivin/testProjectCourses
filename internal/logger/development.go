package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newDevelopmentEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     SyslogTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func NewDevelopment() *zap.Logger {
	return newLogger(
		zapcore.NewConsoleEncoder(newDevelopmentEncoderConfig()),
		zap.DebugLevel,
		zap.Development(),
	)
}
