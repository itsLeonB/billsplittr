package provider_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/stretchr/testify/assert"
)

func TestProviderStruct(t *testing.T) {
	// Test that Provider struct can be created
	p := &provider.Provider{}
	assert.NotNil(t, p)
}

func TestProvider_ShutdownNilComponents(t *testing.T) {
	p := &provider.Provider{}
	err := p.Shutdown()
	assert.NoError(t, err)
}

// Test provider creation with minimal valid config
func TestAllValidConfig(t *testing.T) {
	// Skip this test as it requires actual database connection
	t.Skip("Skipping integration test that requires database")
	
	configs := config.Config{
		App: config.App{
			Name: "test-app",
			Env:  "test",
			Port: "8080",
		},
		DB: config.DB{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     "5432",
			User:     "test",
			Password: "test",
			Name:     "test",
		},
		Google: config.Google{
			ServiceAccount: `{"type": "service_account"}`,
			BillBucketName: "test-bucket",
		},
	}

	assert.NotPanics(t, func() {
		provider.All(configs)
	})
}
