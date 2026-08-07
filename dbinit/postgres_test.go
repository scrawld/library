package dbinit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestOpenPostgres
func TestOpenPostgres(t *testing.T) {
	cfg := PostgresConfig{
		Host:     "127.0.0.1",
		Port:     "5432",
		Username: "max",
		Password: "",
		DBName:   "postgres",
	}

	db, err := OpenPostgres(cfg, true)
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	require.NoError(t, sqlDB.Ping())
}
