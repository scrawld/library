package gormzaplog

import (
	"fmt"
	"strings"

	"github.com/scrawld/library/zaplog"

	"gorm.io/gorm/logger"
)

var Logger *GormZapLogger

// RegisterGlobalLogger initializes the global logger
func RegisterGlobalLogger(directory string, maxAge int, conf logger.Config) error {
	zapLogger, err := zaplog.RegisterLogger(zaplog.Config{
		Level:     "debug", // 不使用zap的日志等级控制,调到最低
		Directory: directory,
		MaxAge:    maxAge,
	})
	if err != nil {
		return err
	}
	Logger = NewGormZapLogger(zapLogger, conf)
	return nil
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
func Sync() error {
	if Logger == nil {
		return nil // Logger has not been initialized yet
	}
	return Logger.logger.Sync()
}

// ParseLevel parses a string into a LogLevel. It returns an error if the input does not match any known log level.
func ParseLevel(text string) (logger.LogLevel, error) {
	lower := strings.ToLower(strings.TrimSpace(text))

	switch lower {
	case "silent":
		return logger.Silent, nil
	case "error":
		return logger.Error, nil
	case "warn", "": // make the zero value useful
		return logger.Warn, nil
	case "info":
		return logger.Info, nil
	}
	return 0, fmt.Errorf("unrecognized level: %s", text)
}
