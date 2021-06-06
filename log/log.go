package log

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"
	"os"
)

var DefaultLogger = &Logger{zerolog.New(os.Stderr).With().CallerWithSkipFrameCount(3).Logger(), context.Background()}

type Logger struct {
	zerolog.Logger
	ctx context.Context
}

func NewLogger(opts ...LogOption) *Logger {
	l := DefaultLogger

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func NewContextLogger(ctx context.Context, ops ...LogOption) (*Logger, context.Context) {
	var logger = NewLogger(ops...)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		if id, err := uuid.NewUUID(); err == nil {
			logger = logger.Trace(id.String())
			md := metadata.Pairs("trace_id", id.String())
			ctx = metadata.NewOutgoingContext(ctx, md)
		}
		return logger, ctx
	}

	traceIds, ok := md["trace_id"]
	if !ok || len(traceIds) == 0 {
		if id, err := uuid.NewUUID(); err == nil {
			logger = logger.Trace(id.String())
			ctx = metadata.AppendToOutgoingContext(ctx, "trace_id", id.String())
		}
		return logger, ctx
	}
	logger = logger.Trace(traceIds[0])

	return logger, ctx
}

func (l *Logger) Trace(traceId string) *Logger {
	return &Logger{l.With().Str("trace_id", traceId).Logger(), l.ctx}
}

func (l *Logger) Name(name string) Logger {
	return Logger{l.With().Str("name", name).Logger(), l.ctx}
}

func (l *Logger) WithStr(key, value string) Logger {
	return Logger{l.With().Str(key, value).Logger(), l.ctx}
}

//func (l *Logger) StrLogger(key, value string) *Logger {
//	return &Logger{l.With().Str(key, value).Logger()}
//}

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
