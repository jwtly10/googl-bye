package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

// Logger is our custom logger interface
type Logger interface {
	Debug(msg string, args ...interface{})
	Debugf(format string, args ...interface{})
	Info(msg string, args ...interface{})
	Infof(format string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// CustomLogger implements the Logger interface
type CustomLogger struct {
	slogger *zap.SugaredLogger
}

// NewLogger creates a new CustomLogger
func NewLogger(useJSON bool, level zapcore.Level) Logger {
	var config zap.Config
	if useJSON {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.Level = zap.NewAtomicLevelAt(level)
	logger, _ := config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	defer logger.Sync() // flushes buffer, if any

	return &CustomLogger{
		slogger: logger.Sugar(),
	}
}

// Debug logs a debug message
func (l *CustomLogger) Debug(msg string, args ...interface{}) {
	l.slogger.Debugw(msg, args...)
}

// Debugf logs a debug message with formatting.
func (l *CustomLogger) Debugf(format string, args ...interface{}) {
	l.slogger.Debugf(format, args...)
}

// Info logs an info message
func (l *CustomLogger) Info(msg string, args ...interface{}) {
	l.slogger.Infow(msg, args...)
}

// Infof logs an info message with formatting.
func (l *CustomLogger) Infof(format string, args ...interface{}) {
	l.slogger.Infof(format, args...)
}

// Warn logs a warning message
func (l *CustomLogger) Warn(msg string, args ...interface{}) {
	l.slogger.Warnw(msg, args...)
}

// Error logs an error message
func (l *CustomLogger) Error(msg string, args ...interface{}) {
	l.slogger.Errorw(msg, args...)
}

// Errorf logs an error message with formatting.
func (l *CustomLogger) Errorf(format string, args ...interface{}) {
	l.slogger.Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func (l *CustomLogger) Fatal(msg string, args ...interface{}) {
	l.slogger.Fatalw(msg, args...)
}

// Fatalf logs a fatal message with formatting and exits
func (l *CustomLogger) Fatalf(format string, args ...interface{}) {
	l.slogger.Fatalf(format, args...)
}

// Global logger instance
var (
	globalLogger Logger
	once         sync.Once
)

// InitLogger initializes the global logger
func InitLogger(useJSON bool, level zapcore.Level) {
	once.Do(func() {
		globalLogger = NewLogger(useJSON, level)
	})
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	if globalLogger == nil {
		InitLogger(false, zapcore.InfoLevel)
	}
	return globalLogger
}
