package dbinit

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormOption customizes a gorm.Config before the connection is opened.
type GormOption func(*gorm.Config) error

// WithLogger sets a custom GORM logger.
func WithLogger(l logger.Interface) GormOption {
	return func(cfg *gorm.Config) error {
		cfg.Logger = l
		return nil
	}
}

// WithDebug enables debug mode by setting the log level to Info,
// which prints every executed SQL statement.
func WithDebug() GormOption {
	return func(cfg *gorm.Config) error {
		cfg.Logger = cfg.Logger.LogMode(logger.Info)
		return nil
	}
}
