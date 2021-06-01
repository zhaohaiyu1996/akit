package alog

import (
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"io"
	"strings"
)

func stringToLevel(level string) zerolog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return zerolog.DebugLevel
	case "INFO":
		return zerolog.InfoLevel
	case "WARN":
		return zerolog.WarnLevel
	case "ERROR":
		return zerolog.DebugLevel
	default:
		return zerolog.InfoLevel
	}
}

func getLogWriter(filename string, maxSize, maxBackups, maxAge int, compress bool) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
}
