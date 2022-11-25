package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger(encoder zapcore.Encoder, levelEnabler zapcore.LevelEnabler, options ...zap.Option) *zap.Logger {
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
		),
		levelEnabler,
	)
	options = append(options, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))

	return zap.New(core, options...)
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.1234 -07:00"))
}
