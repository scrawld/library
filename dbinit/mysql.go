package dbinit

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// MySQLConfig MySQL connection config
type MySQLConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

// DSN builds the MySQL DSN string
func (c MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.DBName)
}

// OpenMySQL initializes a MySQL connection; set debug=true to log SQL
func OpenMySQL(config MySQLConfig, debug bool) (*gorm.DB, error) {
	logLevel := logger.Silent
	if debug {
		logLevel = logger.Info
	}
	dsn := config.DSN()

	db, err := gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			NamingStrategy:  schema.NamingStrategy{SingularTable: true}, // 禁用表名复数
			Logger:          logger.Default.LogMode(logLevel),
			CreateBatchSize: 1000,
		})
	if err != nil {
		return nil, fmt.Errorf("db connect fail: %s, %s", err, dsn)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db fail: %s, %s", err, dsn)
	}
	sqlDB.SetMaxIdleConns(10)                  // max idle connections
	sqlDB.SetMaxOpenConns(100)                 // max open connections
	sqlDB.SetConnMaxLifetime(15 * time.Minute) // max connection lifetime

	return db, nil
}
