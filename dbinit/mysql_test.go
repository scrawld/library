package dbinit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestOpenMySQL
func TestOpenMySQL(t *testing.T) {
	cfg := MySQLConfig{
		Host:     "127.0.0.1",
		Port:     "3306",
		Username: "root",
		Password: "12345678",
		DBName:   "account",
	}

	db, err := OpenMySQL(cfg, true)
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	require.NoError(t, sqlDB.Ping())
}
