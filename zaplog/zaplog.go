package zaplog

import (
	"fmt"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// RegisterGlobalLogger initializes the global logger
func RegisterGlobalLogger(cfg Config, options ...zap.Option) error {
	logger, err := RegisterLogger(cfg, options...)
	if err != nil {
		return err
	}
	Logger = logger
	return nil
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
func Sync() error {
	if Logger == nil {
		return nil // Logger has not been initialized yet
	}
	return Logger.Sync()
}

type Config struct {
	Level        string            `json:"level"`        // 日志等级: debug/info/warn/error/dpanic/panic/fatal,default:info
	Encoding     string            `json:"encoding"`     // 日志格式: json/console,default:console
	Directory    string            `json:"directory"`    // 日志目录
	MaxAge       int               `json:"maxAge"`       // 保留日志文件的最大天数
	LogInConsole bool              `json:"logInConsole"` // 是否在终端输出
	FileNames    map[string]string `json:"fileNames"`    // 自定义日志文件名映射: 使用日志级别作为键 (如 "main", "warn", "error")，文件名作为值。
}

// RegisterLogger initializes the logger
func RegisterLogger(cfg Config, options ...zap.Option) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("parse level error: %s", err)
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
	mainFilename := getOrDefault(cfg.FileNames["main"], "app.log")
	// add main core
	if err = addCore(mainFilename, level, cfg.LogInConsole); err != nil {
		return nil, fmt.Errorf("add main core error: %s", err)
	}
	// add core
	for _, v := range []zapcore.Level{zapcore.WarnLevel, zapcore.ErrorLevel} {
		if !level.Enabled(v) {
			continue
		}
		filename := getOrDefault(cfg.FileNames[v.String()], "app."+v.String()+".log")

		if err = addCore(filename, GetLevelEnabler(v), false); err != nil {
			return nil, fmt.Errorf("get %s core error: %s", v.String(), err)
		}
	}
	return zap.New(zapcore.NewTee(cores...), options...), nil
}

// getOrDefault returns the value if it's not empty, otherwise returns the defaultValue
func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
