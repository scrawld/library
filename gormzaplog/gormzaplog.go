package gormzaplog

import (
	"test_gorm/library/zaplog"

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
