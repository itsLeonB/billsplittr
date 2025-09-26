package config_test

import (
	"os"
	"testing"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadWithEnvironmentVariables(t *testing.T) {
	// Set required environment variables
	_ = os.Setenv("DB_HOST", "localhost")
	_ = os.Setenv("DB_PORT", "5432")
	_ = os.Setenv("DB_USER", "testuser")
	_ = os.Setenv("DB_PASSWORD", "testpass")
	_ = os.Setenv("DB_NAME", "testdb")
	_ = os.Setenv("STORAGE_BUCKET_NAME_EXPENSE_BILL", "bills")

	defer func() {
		_ = os.Unsetenv("DB_HOST")
		_ = os.Unsetenv("DB_PORT")
		_ = os.Unsetenv("DB_USER")
		_ = os.Unsetenv("DB_PASSWORD")
		_ = os.Unsetenv("DB_NAME")
		_ = os.Unsetenv("STORAGE_BUCKET_NAME_EXPENSE_BILL")
	}()

	cfg := config.Load()

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "5432", cfg.DB.Port)
	assert.Equal(t, "testuser", cfg.User)
	assert.Equal(t, "testpass", cfg.DB.Password)
	assert.Equal(t, "testdb", cfg.Name)
	assert.Equal(t, "bills", cfg.BucketNameExpenseBill)
}
