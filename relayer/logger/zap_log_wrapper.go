package logger

import (
	"log"

	"go.uber.org/zap"
)

var _ Logger = (*ZapLogWrapper)(nil)

// ZapLogWrapper is a wrapper around the zap logger so that it align with an Logger interface.
// Note: zap.Logger doesn't implement Logger.With in which return *zap.Logger itself (external package)
type ZapLogWrapper struct {
	inner *zap.Logger
}

// NewZapLogWrapper creates a new ZapLogWrapper.
func NewZapLogWrapper(inner *zap.Logger) *ZapLogWrapper {
	return &ZapLogWrapper{inner: inner}
}

// With returns a new Logger with the given fields.
func (w *ZapLogWrapper) With(fields ...zap.Field) Logger {
	return NewZapLogWrapper(w.inner.With(fields...))
}

// Debug logs a message at debug level.
func (w *ZapLogWrapper) Debug(msg string, fields ...zap.Field) {
	w.inner.Debug(msg, fields...)
}

// Info logs a message at info level.
func (w *ZapLogWrapper) Info(msg string, fields ...zap.Field) {
	w.inner.Info(msg, fields...)
}

// Warn logs a message at warn level.
func (w *ZapLogWrapper) Warn(msg string, fields ...zap.Field) {
	w.inner.Warn(msg, fields...)
}

// Error logs a message at error level.
func (w *ZapLogWrapper) Error(msg string, fields ...zap.Field) {
	w.inner.Error(msg, fields...)
}

// Fatal logs a message at fatal level.
func (w *ZapLogWrapper) Fatal(msg string, fields ...zap.Field) {
	w.inner.Fatal(msg, fields...)
}

// Panic logs a message at panic level.
func (w *ZapLogWrapper) Panic(msg string, fields ...zap.Field) {
	w.inner.Panic(msg, fields...)
}

// Sync flushes any buffered log entries.
func (w *ZapLogWrapper) Sync() error {
	return w.inner.Sync()
}

// ToStdLog returns a new stdlib logger that writes to the zap logger.
func (w *ZapLogWrapper) ToStdLog() *log.Logger {
	return zap.NewStdLog(w.inner)
}
