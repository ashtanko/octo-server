package config_test

import (
	"github.com/ashtanko/octo-server/app/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	_ = os.Setenv("SERVER_PORT", "8000")
	_ = os.Setenv("DB_HOST", "localhost")
	_ = os.Setenv("DB_DATABASE", "db_name")
	_ = os.Setenv("DB_PORT", "5432")
	_ = os.Setenv("DB_USER", "postgres")
	_ = os.Setenv("DB_PASSWORD", "admin")

	defer func() {
		_ = os.Unsetenv("SERVER_PORT")
	}()

	cfg, err := config.LoadConfig()
	assert.Nil(t, err)
	assert.Equal(t, 8000, cfg.Port)
}

func TestLoad_WithError(t *testing.T) {
	_ = os.Setenv("SERVER_PORT", "invalid-value")
	defer func() {
		_ = os.Unsetenv("SERVER_PORT")
	}()

	_, err := config.LoadConfig()

	assert.NotNil(t, err)
}

