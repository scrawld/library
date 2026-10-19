package dbinit

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// PostgresConfig PostgreSQL connection config
type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// DSN builds the MySQL DSN string
func (c PostgresConfig) DSN() string {
	sslmode := c.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.Username, c.Password, c.Host, c.Port, c.DBName, sslmode)
}

// OpenPostgres initializes a PostgreSQL connection; set debug=true to log SQL
func OpenPostgres(config PostgresConfig, options ...GormOption) (*gorm.DB, error) {
	// Default GORM configuration.
	gormConfig := &gorm.Config{
		NamingStrategy:  schema.NamingStrategy{SingularTable: true}, // use singular table names
		Logger:          logger.Default,
		CreateBatchSize: 1000,
	}

	// Apply the caller-provided options in order; later options override earlier ones.
	for _, option := range options {
		if err := option(gormConfig); err != nil {
			return nil, fmt.Errorf("apply gorm option: %w", err)
		}
	}
	dsn := config.DSN()

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("db connect fail: %s, dsn=%s", err, dsn)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db fail: %s, dsn=%s", err, dsn)
	}
	sqlDB.SetMaxIdleConns(10)                  // max idle connections
	sqlDB.SetMaxOpenConns(100)                 // max open connections
	sqlDB.SetConnMaxLifetime(15 * time.Minute) // max connection lifetime

	return db, nil
}
