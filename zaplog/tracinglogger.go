package zaplog

import (
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
func New(traceId ...string) (r *TracingLogger) {
	var tid string
	if len(traceId) != 0 && len(traceId[0]) != 0 {
		tid = traceId[0]
	} else {
		uid, _ := uuid.NewRandom()
		tid = uid.String()
	}
	o := &TracingLogger{
		TraceId: tid, prefix: "[tid:" + tid + "]",
		logger: Logger.Sugar().WithOptions(zap.AddCallerSkip(1)),
	}
	return o
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

// Debugf logs a debug message with the formatted string and arguments.
func (l *TracingLogger) Debugf(f string, v ...interface{}) {
	l.logger.Debugf(l.prefix+" "+f, v...)
}

// Infof logs an info message with the formatted string and arguments.
func (l *TracingLogger) Infof(f string, v ...interface{}) {
	l.logger.Infof(l.prefix+" "+f, v...)
}

// Warnf logs a warning message with the formatted string and arguments.
func (l *TracingLogger) Warnf(f string, v ...interface{}) {
	l.logger.Warnf(l.prefix+" "+f, v...)
}

// Errorf logs an error message with the formatted string and arguments.
func (l *TracingLogger) Errorf(f string, v ...interface{}) {
	l.logger.Errorf(l.prefix+" "+f, v...)
}
