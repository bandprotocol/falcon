package logger

import (
	"log"

	"go.uber.org/zap"
)

// ZapLogger is an interface that implements the basic methods of the zap logger.
type ZapLogger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)

	Sync() error
	With(fields ...zap.Field) ZapLogger
	ToStdLog() *log.Logger
}
