package logger

import (
	"log"
)

// Logger is an interface that implements the basic methods of the zap logger.
type Logger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
	Fatal(msg string, fields ...any)
	Panic(msg string, fields ...any)

	Sync() error
	With(fields ...any) Logger
	ToStdLog() *log.Logger
}
