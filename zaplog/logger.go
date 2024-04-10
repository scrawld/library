package zaplog

import (
	"fmt"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
)

type Config struct {
	Level     string `json:"level"`     // 日志等级: debug/info/warn/error/dpanic/panic/fatal,default:info
	Encoding  string `json:"encoding"`  // 日志格式: json/console,default:console
	Directory string `json:"directory"` // 日志目录
	MaxAge    int    `json:"maxAge"`    // 保留日志文件的最大天数
}

// RegisterLogger 初始化日志
func RegisterLogger(cfg Config, options ...zap.Option) error {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("parse level error: %s", err)
	}
	if cfg.Encoding == "" {
		cfg.Encoding = "console"
	}
	var (
		cores   = []zapcore.Core{}
		encoder = GetEncoder(cfg.Encoding)
	)
	addCore := func(filename string, enab zapcore.LevelEnabler, logInConsole bool) error {
		writer, err := GetWriteSyncer(path.Join(cfg.Directory, filename), cfg.MaxAge, logInConsole)
		if err != nil {
			return fmt.Errorf("get write syncer error: %s", err.Error())
		}
		cores = append(cores, zapcore.NewCore(encoder, writer, enab))
		return nil
	}
	// add main core
	err = addCore("app.log", level, true)
	if err != nil {
		return fmt.Errorf("add main core error: %s", err)
	}
	// add core
	for _, v := range []zapcore.Level{zapcore.WarnLevel, zapcore.ErrorLevel} {
		if level.Enabled(v) {
			err = addCore("app."+v.String()+".log", GetLevelEnabler(v), false)
			if err != nil {
				return fmt.Errorf("get %s core error: %s", v.String(), err)
			}
		}
	}
	Logger = zap.New(zapcore.NewTee(cores...), options...)
	return nil
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
func Sync() error {
	return Logger.Sync()
}
