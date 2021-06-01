package alog

import (
	"github.com/rs/zerolog"
)

type LogOption func(o *Logger)

func WithLogLevel(level string) LogOption {
	return func(o *Logger) {
		zerolog.SetGlobalLevel(stringToLevel(level))
	}
}

// WithFileLog is print log to file
// filename log's file path
// maxSize 在进行切割之前，日志文件的最大大小（以MB为单位）
// maxBackups 保留旧文件的最大个数
// maxAge：保留旧文件的最大天数
// compress：是否压缩/归档旧文件
func WithFileLog(filename string, maxSize, maxBackups, maxAge int, compress, caller bool) LogOption {
	return func(o *Logger) {
		logger := zerolog.New(getLogWriter(filename, maxSize, maxBackups, maxAge, compress))
		if caller {
			logger = logger.With().CallerWithSkipFrameCount(3).Logger()
		}
		o.Logger = logger
	}
}
