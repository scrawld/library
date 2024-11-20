package zaplog

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TracingLogger represents a logger with tracing capabilities.
type TracingLogger struct {
	TraceId string
	prefix  string
	logger  *zap.SugaredLogger
}

// New creates a new TracingLogger instance with an optional trace ID.
func New(traceId ...string) *TracingLogger {
	tid := generateTraceID(traceId...)
	o := &TracingLogger{
		TraceId: tid,
		prefix:  "[tid:" + tid + "] ",
		logger:  Logger.Sugar().WithOptions(zap.AddCallerSkip(1)),
	}
	return o
}

// NewWithLogger creates a new TracingLogger instance with an existing zap.Logger and an optional trace ID.
func NewWithLogger(logger *zap.Logger, traceId ...string) *TracingLogger {
	tid := generateTraceID(traceId...)
	o := &TracingLogger{
		TraceId: tid,
		prefix:  "[tid:" + tid + "] ",
		logger:  logger.Sugar().WithOptions(zap.AddCallerSkip(1)),
	}
	return o
}

// generateTraceID generates a new trace ID if one is not provided.
func generateTraceID(traceId ...string) string {
	if len(traceId) != 0 && len(traceId[0]) != 0 {
		return traceId[0]
	}
	uid, _ := uuid.NewRandom()
	return uid.String()
}

// clone creates a copy of the TracingLogger instance.
func (l *TracingLogger) clone() *TracingLogger {
	copy := *l
	return &copy
}

// WithOptions sets additional options for the logger.
func (l *TracingLogger) WithOptions(opts ...zap.Option) *TracingLogger {
	l.logger = l.clone().logger.WithOptions(opts...)
	return l
}

// AddCallerSkip adjusts the number of caller frames skipped by the logger.
func (l *TracingLogger) AddCallerSkip(skip int) *TracingLogger {
	return l.clone().WithOptions(zap.AddCallerSkip(skip))
}

// CloseCaller disables caller information in the logger.
func (l *TracingLogger) CloseCaller() *TracingLogger {
	return l.clone().WithOptions(zap.WithCaller(false))
}

// Named adds a sub-logger with the given name.
func (l *TracingLogger) Named(name string) *TracingLogger {
	l.logger = l.logger.Named(name)
	return l
}

// Debug logs a debug message with the given arguments.
func (l *TracingLogger) Debug(args ...interface{}) {
	l.logger.Debug(l.prefix, fmt.Sprint(args...))
}

// Info logs an info message with the given arguments.
func (l *TracingLogger) Info(args ...interface{}) {
	l.logger.Info(l.prefix, fmt.Sprint(args...))
}

// Warn logs a warning message with the given arguments.
func (l *TracingLogger) Warn(args ...interface{}) {
	l.logger.Warn(l.prefix, fmt.Sprint(args...))
}

// Error logs an error message with the given arguments.
func (l *TracingLogger) Error(args ...interface{}) {
	l.logger.Error(l.prefix, fmt.Sprint(args...))
}

// Debugf logs a debug message with the formatted string and arguments.
func (l *TracingLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(l.prefix+template, args...)
}

// Infof logs an info message with the formatted string and arguments.
func (l *TracingLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(l.prefix+template, args...)
}

// Warnf logs a warning message with the formatted string and arguments.
func (l *TracingLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(l.prefix+template, args...)
}

// Errorf logs an error message with the formatted string and arguments.
func (l *TracingLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(l.prefix+template, args...)
}
