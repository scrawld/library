package gormzaplog

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func TestMain(m *testing.M) {
	var (
		directory = "logs/orm" // orm日志保存目录
		maxAge    = 30         // 日志保留天数

		ormLogConfig = logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Info,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,                  // Don't include params in the SQL log
			Colorful:                  false,                  // Disable color
		}
	)

	err := RegisterGlobalLogger(directory, maxAge, ormLogConfig)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer Sync()

	os.Exit(m.Run())
}

func TestGormZapLogger(t *testing.T) {
	// Test Info log
	Logger.Info(context.Background(), "test info message %s %s", "key1", "value1")

	// Test Warn log
	Logger.Warn(context.Background(), "test warn message %s %s", "key2", "value2")

	// Test Error log
	Logger.Error(context.Background(), "test error message %s %s", "key3", "value3")

	// Test Trace log
	Logger.Trace(context.Background(), time.Now().Add(-time.Second), func() (string, int64) {
		return "SELECT * FROM users", 5
	}, nil)
}

// User represents a simple user model for testing
type User struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `gorm:"type:varchar(100)"`
	Email string `gorm:"type:varchar(100);uniqueIndex"`
	Num   int64  `gorm:"type:bigint(20)"`
}

func TestFailedUpsert(t *testing.T) {
	dsn := "root:12345678@tcp(127.0.0.1:3306)/bra88_log?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy:  schema.NamingStrategy{SingularTable: true}, // 禁用表名复数
		Logger:          Logger,
		CreateBatchSize: 1000,
	})
	assert.Nil(t, err)

	//err = db.AutoMigrate(&User{})
	//assert.Nil(t, err)

	sql := `INSERT INTO user1 (name, email) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name)`
	result := db.Exec(sql, "zhangsan", "zhangsan@163.com")
	assert.Nil(t, result.Error)

	// select
	u := &User{}
	result = db.Where("name=?", "zhangsan1").Take(u)
	assert.Nil(t, result.Error)

	// upsert
	u = &User{
		ID:    7,
		Name:  "zhangsan",
		Email: "zhangsan@163.com",
		Num:   1,
	}
	result = db.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"role": "user",
		}),
	}).Create(u)
	assert.Nil(t, result.Error)
}
