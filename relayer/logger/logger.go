package logger

import (
	"log"

	"go.uber.org/zap"
)

// Logger is an interface that wraps the basic methods of the logger.
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)

	Sync() error
	With(fields ...zap.Field) Logger
	ToStdLog() *log.Logger
}
