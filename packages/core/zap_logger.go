package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new Logger implementation using zap
func NewZapLogger(development bool) (Logger, error) {
	var logger *zap.Logger
	var err error

	// Get the caller's directory (service directory)
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("failed to get caller info")
	}

	// Get the service root directory (the directory containing the service's main.go)
	serviceDir := filepath.Dir(filename)
	for {
		// Check if we're in a service directory by looking for main.go
		if _, err := os.Stat(filepath.Join(serviceDir, "cmd", "main.go")); err == nil {
			break
		}
		// Move up one directory
		parent := filepath.Dir(serviceDir)
		if parent == serviceDir {
			return nil, fmt.Errorf("could not find service root directory")
		}
		serviceDir = parent
	}

	// Create tmp directory in the service root directory
	tmpDir := filepath.Join(serviceDir, "tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create tmp directory: %w", err)
	}

	// Create log file
	logFile := filepath.Join(tmpDir, "service.log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create file encoder
	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	})

	// Create console encoder
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// Create core for file output
	fileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(file),
		zapcore.DebugLevel,
	)

	// Create core for console output
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	// Create the logger with both cores
	logger = zap.New(zapcore.NewTee(fileCore, consoleCore))

	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(logger)

	return &ZapLogger{
		logger: logger,
	}, nil
}

// getContextFields extracts relevant fields from context
func (l *ZapLogger) getContextFields(ctx context.Context) []Field {
	if ctx == nil {
		return nil
	}

	var fields []Field

	// Get trace ID from context
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, NewField("trace_id", traceID))
	}

	// Get span ID from OpenTelemetry context
	if span := trace.SpanFromContext(ctx); span != nil {
		if spanID := span.SpanContext().SpanID().String(); spanID != "" {
			fields = append(fields, NewField("span_id", spanID))
		}
	}

	// Add service name from context if available
	if serviceName := ctx.Value("service_name"); serviceName != nil {
		fields = append(fields, NewField("service", serviceName))
	}

	return fields
}

// Debug logs a debug message
func (l *ZapLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	ctxFields := l.getContextFields(ctx)
	l.logger.Debug(msg, l.convertFields(append(ctxFields, fields...)...)...)
}

// Debugf logs a formatted debug message
func (l *ZapLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.Debug(ctx, fmt.Sprintf(format, args...))
}

// Info logs an info message
func (l *ZapLogger) Info(ctx context.Context, msg string, fields ...Field) {
	ctxFields := l.getContextFields(ctx)
	l.logger.Info(msg, l.convertFields(append(ctxFields, fields...)...)...)
}

// Infof logs a formatted info message
func (l *ZapLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.Info(ctx, fmt.Sprintf(format, args...))
}

// Warn logs a warning message
func (l *ZapLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	ctxFields := l.getContextFields(ctx)
	l.logger.Warn(msg, l.convertFields(append(ctxFields, fields...)...)...)
}

// Warnf logs a formatted warning message
func (l *ZapLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.Warn(ctx, fmt.Sprintf(format, args...))
}

// Error logs an error message
func (l *ZapLogger) Error(ctx context.Context, msg string, fields ...Field) {
	ctxFields := l.getContextFields(ctx)
	l.logger.Error(msg, l.convertFields(append(ctxFields, fields...)...)...)
}

// Errorf logs a formatted error message
func (l *ZapLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.Error(ctx, fmt.Sprintf(format, args...))
}

// Fatal logs a fatal message
func (l *ZapLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	ctxFields := l.getContextFields(ctx)
	l.logger.Fatal(msg, l.convertFields(append(ctxFields, fields...)...)...)
}

// Fatalf logs a formatted fatal message
func (l *ZapLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.Fatal(ctx, fmt.Sprintf(format, args...))
}

// With returns a new logger with the given fields
func (l *ZapLogger) With(fields ...Field) Logger {
	return &ZapLogger{
		logger: l.logger.With(l.convertFields(fields...)...),
	}
}

// convertFields converts our Field type to zap.Field
func (l *ZapLogger) convertFields(fields ...Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

func (l *ZapLogger) GetZapLogger() *zap.Logger {
	return l.logger
}
