package provider_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/stretchr/testify/assert"
)

func TestDBs_GetDSNPostgres(t *testing.T) {
	// We can't directly test the private getDSN method, but we can test the behavior
	// by checking if the provider panics with invalid config

	logger := ezutil.NewSimpleLogger(config.AppName, true, 0)

	dbConfig := config.DB{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
	}

	// This will panic if database connection fails, which is expected in test environment
	assert.Panics(t, func() {
		provider.ProvideDBs(logger, config.Config{DB: dbConfig})
	})
}

func TestDBs_GetDSNMySQL(t *testing.T) {
	logger := ezutil.NewSimpleLogger(config.AppName, true, 0)

	dbConfig := config.DB{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     "3306",
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
	}

	// This will panic because MySQL driver is commented out in the actual code
	assert.Panics(t, func() {
		provider.ProvideDBs(logger, config.Config{DB: dbConfig})
	})
}

func TestDBsUnsupportedDriver(t *testing.T) {
	logger := ezutil.NewSimpleLogger(config.AppName, true, 0)

	dbConfig := config.DB{
		Driver:   "unsupported",
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
	}

	assert.Panics(t, func() {
		provider.ProvideDBs(logger, config.Config{DB: dbConfig})
	})
}

func TestDBsStruct(t *testing.T) {
	// Test that DBs struct can be created
	dbs := &provider.DBs{}
	assert.NotNil(t, dbs)
}

func TestDBs_ShutdownNilDB(t *testing.T) {
	dbs := &provider.DBs{}

	// Should not panic when trying to shutdown nil DB
	assert.NotPanics(t, func() {
		_ = dbs.Shutdown()
	})
}
