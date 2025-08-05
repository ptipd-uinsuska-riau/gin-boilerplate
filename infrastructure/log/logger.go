package log

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func Init(format, level string) {
	logger = logrus.New()

	// Set log format
	switch format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// Set output
	logger.SetOutput(os.Stdout)
}

func GetLogger() *logrus.Logger {
	if logger == nil {
		Init("text", "info")
	}
	return logger
}

// Context-aware logging functions
func Debug(ctx context.Context, args ...interface{}) {
	GetLogger().Debug(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	GetLogger().Info(args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	GetLogger().Warn(args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	GetLogger().Error(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	GetLogger().Fatal(args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// Non-context logging functions for backward compatibility
func InfoLog(args ...interface{}) {
	GetLogger().Info(args...)
}

func ErrorLog(args ...interface{}) {
	GetLogger().Error(args...)
}

func WarnLog(args ...interface{}) {
	GetLogger().Warn(args...)
}

func DebugLog(args ...interface{}) {
	GetLogger().Debug(args...)
}
