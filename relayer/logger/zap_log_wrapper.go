package logger

import (
	"log"

	"go.uber.org/zap"
)

var _ Logger = (*ZapLogWrapper)(nil)

// ZapLogWrapper is a wrapper around the zap logger so that it align with an Logger interface.
// Note: zap.Logger doesn't implement Logger.With in which return *zap.Logger itself (external package)
type ZapLogWrapper struct {
	inner *zap.SugaredLogger
}

// NewZapLogWrapper creates a new ZapLogWrapper.
func NewZapLogWrapper(inner *zap.SugaredLogger) *ZapLogWrapper {
	return &ZapLogWrapper{inner: inner}
}

// With returns a new Logger with the given fields.
func (w *ZapLogWrapper) With(fields ...any) Logger {
	return NewZapLogWrapper(w.inner.With(fields...))
}

// Debug logs a message at debug level.
func (w *ZapLogWrapper) Debug(msg string, fields ...any) {
	w.inner.Debugw(msg, fields)
}

// Info logs a message at info level.
func (w *ZapLogWrapper) Info(msg string, fields ...any) {
	w.inner.Infow(msg, fields...)
}

// Warn logs a message at warn level.
func (w *ZapLogWrapper) Warn(msg string, fields ...any) {
	w.inner.Warnw(msg, fields...)
}

// Error logs a message at error level.
func (w *ZapLogWrapper) Error(msg string, fields ...any) {
	w.inner.Errorw(msg, fields...)
}

// Fatal logs a message at fatal level.
func (w *ZapLogWrapper) Fatal(msg string, fields ...any) {
	w.inner.Fatalw(msg, fields...)
}

// Panic logs a message at panic level.
func (w *ZapLogWrapper) Panic(msg string, fields ...any) {
	w.inner.Panicw(msg, fields...)
}

// Sync flushes any buffered log entries.
func (w *ZapLogWrapper) Sync() error {
	return w.inner.Sync()
}

// ToStdLog returns a new stdlib logger that writes to the zap logger.
func (w *ZapLogWrapper) ToStdLog() *log.Logger {
	return zap.NewStdLog(w.inner.Desugar())
}

type ZapLogField struct {
	inner zap.Field
}
