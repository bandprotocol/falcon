package cmd

import (
	"fmt"
	"os"
	"time"

	zaplogfmt "github.com/jsternberg/zap-logfmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// newLogger creates a new root logger with the given log format and log level.
func newLogger(format string, logLevel string) (*zap.Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format("2006-01-02T15:04:05.000000Z07:00"))
	}
	config.LevelKey = "lvl"
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	var enc zapcore.Encoder
	switch format {
	case "json":
		enc = zapcore.NewJSONEncoder(config)
	case "auto", "console":
		enc = zapcore.NewConsoleEncoder(config)
	case "logfmt":
		enc = zaplogfmt.NewEncoder(config)
	default:
		return nil, fmt.Errorf("unrecognized log format %q", format)
	}

	level := zap.InfoLevel
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	}

	logger := zap.New(zapcore.NewCore(enc, os.Stderr, level))
	if logLevel == "debug" {
		logger = logger.WithOptions(zap.AddCaller())
	}

	return logger, nil
}

// initLogger initializes the logger with the given default log level.
func initLogger(defaultLogLevel string) (log *zap.Logger, err error) {
	logFormat := viper.GetString("log-format")

	logLevel := viper.GetString("log-level")
	if viper.GetBool("debug") {
		logLevel = "debug"
	}
	if logLevel == "" && defaultLogLevel != "" {
		logLevel = defaultLogLevel
	}

	// initialize logger only if user run command "start" or log level is "debug"
	if os.Args[1] == "start" || logLevel == "debug" {
		log, err = newLogger(logFormat, logLevel)
		if err != nil {
			return nil, err
		}
	} else {
		log = zap.NewNop()
	}

	return log, nil
}
