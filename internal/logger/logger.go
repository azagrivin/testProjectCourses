package logger

import (
	"net/http"
)

type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
}

type sentryLogger struct {
	Logger
	r *http.Request
}

type HttpLogger interface {
	Logger
	WithRequest(req *http.Request) HttpLogger
}

func New(logger Logger) *sentryLogger {
	return &sentryLogger{
		Logger: logger,
	}
}

func (l *sentryLogger) WithRequest(req *http.Request) HttpLogger {
	log := *l
	log.r = req
	return &log
}

func (l *sentryLogger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *sentryLogger) Infof(pattern string, args ...interface{}) {
	l.Logger.Infof(pattern, args...)
}

func (l *sentryLogger) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}

func (l *sentryLogger) Warnf(pattern string, args ...interface{}) {
	l.Logger.Warnf(pattern, args...)
}

func (l *sentryLogger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

func (l *sentryLogger) Debugf(pattern string, args ...interface{}) {
	l.Logger.Debugf(pattern, args...)
}

func (l *sentryLogger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l *sentryLogger) Errorf(pattern string, args ...interface{}) {
	l.Logger.Errorf(pattern, args...)
}
