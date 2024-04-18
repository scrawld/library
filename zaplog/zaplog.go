package zaplog

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TracingLogger struct {
	TraceId string
	prefix  string
	logger  *zap.SugaredLogger
}

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

func (l *TracingLogger) clone() *TracingLogger {
	copy := *l
	return &copy
}

func (l *TracingLogger) WithOptions(opts ...zap.Option) *TracingLogger {
	l.logger = l.clone().logger.WithOptions(opts...)
	return l
}

func (l *TracingLogger) AddCallerSkip(skip int) *TracingLogger {
	return l.clone().WithOptions(zap.AddCallerSkip(skip))
}

// 关闭行号
func (l *TracingLogger) CloseCaller() *TracingLogger {
	return l.clone().WithOptions(zap.WithCaller(false))
}

func (l *TracingLogger) Named(name string) *TracingLogger {
	l.logger = l.logger.Named(name)
	return l
}

func (l *TracingLogger) Debugf(f string, v ...interface{}) {
	l.logger.Debugf(l.prefix+" "+f, v...)
}

func (l *TracingLogger) Infof(f string, v ...interface{}) {
	l.logger.Infof(l.prefix+" "+f, v...)
}

func (l *TracingLogger) Warnf(f string, v ...interface{}) {
	l.logger.Warnf(l.prefix+" "+f, v...)
}

func (l *TracingLogger) Errorf(f string, v ...interface{}) {
	l.logger.Errorf(l.prefix+" "+f, v...)
}
