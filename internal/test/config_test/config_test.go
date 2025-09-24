package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigStruct(t *testing.T) {
	cfg := config.Config{
		App: config.App{
			Env:     "test",
			Port:    "8080",
			Timeout: 10 * time.Second,
		},
		DB: config.DB{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     "5432",
			User:     "testuser",
			Password: "testpass",
			Name:     "testdb",
		},
	}

	assert.Equal(t, "test", cfg.Env)
	assert.Equal(t, "8080", cfg.App.Port)
	assert.Equal(t, 10*time.Second, cfg.Timeout)
	assert.Equal(t, "postgres", cfg.Driver)
	assert.Equal(t, "localhost", cfg.Host)
}

func TestAppDefaultValues(t *testing.T) {
	app := config.App{}

	// Test that struct can be created with zero values
	assert.Equal(t, "", app.Env)
	assert.Equal(t, "", app.Port)
	assert.Equal(t, time.Duration(0), app.Timeout)
}

func TestDBFields(t *testing.T) {
	db := config.DB{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
	}

	assert.Equal(t, "postgres", db.Driver)
	assert.Equal(t, "localhost", db.Host)
	assert.Equal(t, "5432", db.Port)
	assert.Equal(t, "testuser", db.User)
	assert.Equal(t, "testpass", db.Password)
	assert.Equal(t, "testdb", db.Name)
}

func TestLoadWithEnvironmentVariables(t *testing.T) {
	// Set required environment variables
	_ = os.Setenv("DB_HOST", "localhost")
	_ = os.Setenv("DB_PORT", "5432")
	_ = os.Setenv("DB_USER", "testuser")
	_ = os.Setenv("DB_PASSWORD", "testpass")
	_ = os.Setenv("DB_NAME", "testdb")

	defer func() {
		_ = os.Unsetenv("DB_HOST")
		_ = os.Unsetenv("DB_PORT")
		_ = os.Unsetenv("DB_USER")
		_ = os.Unsetenv("DB_PASSWORD")
		_ = os.Unsetenv("DB_NAME")
	}()

	cfg := config.Load()

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "5432", cfg.DB.Port)
	assert.Equal(t, "testuser", cfg.User)
	assert.Equal(t, "testpass", cfg.DB.Password)
	assert.Equal(t, "testdb", cfg.Name)
}
