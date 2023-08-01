package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var levelEnablerFuncMap = map[zapcore.Level]zap.LevelEnablerFunc{
	zapcore.DebugLevel:  func(level zapcore.Level) bool { return level == zap.DebugLevel },  // 调试级别
	zapcore.InfoLevel:   func(level zapcore.Level) bool { return level == zap.InfoLevel },   // 日志级别
	zapcore.WarnLevel:   func(level zapcore.Level) bool { return level == zap.WarnLevel },   // 警告级别
	zapcore.ErrorLevel:  func(level zapcore.Level) bool { return level == zap.ErrorLevel },  // 错误级别
	zapcore.DPanicLevel: func(level zapcore.Level) bool { return level == zap.DPanicLevel }, // dpanic级别
	zapcore.PanicLevel:  func(level zapcore.Level) bool { return level == zap.PanicLevel },  // panic级别
	zapcore.FatalLevel:  func(level zapcore.Level) bool { return level == zap.FatalLevel },  // 终止级别
}

// GetLevelEnabler 获取日志强等过滤器
func GetLevelEnabler(level zapcore.Level) zapcore.LevelEnabler {
	fn, found := levelEnablerFuncMap[level]
	if !found {
		return level
	}
	return fn
}
