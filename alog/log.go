package alog

import (
	"github.com/rs/zerolog"
	"os"
)

var DefaultLogger = &Logger{zerolog.New(os.Stderr).With().CallerWithSkipFrameCount(3).Logger()}

type Logger struct {
	zerolog.Logger
}

func NewLogger(opts ...LogOption) *Logger {
	l := DefaultLogger

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Logger) Trace(id string) {
	l.Logger = l.With().Str("trace_id", id).Logger()
}

func (l *Logger) Debug(msg string) {
	l.Logger.Debug().Msg(msg)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Logger.Debug().Msgf(format, v...)
}

func (l *Logger) Info(msg string) {
	l.Logger.Info().Msg(msg)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Logger.Info().Msgf(format, v...)
}

func (l *Logger) Warn(msg string) {
	l.Logger.Warn().Msg(msg)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Logger.Warn().Msgf(format, v...)
}

func (l *Logger) Error(msg string) {
	l.Logger.Error().Msg(msg)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Logger.Error().Msgf(format, v...)
}
