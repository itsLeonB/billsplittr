package delivery_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/config"
	grpcdelivery "github.com/itsLeonB/billsplittr/internal/delivery/grpc"
	"github.com/stretchr/testify/assert"
)

func TestGrpcServerSetup(t *testing.T) {
	// Skip this test as it requires actual database connection and external dependencies
	t.Skip("Skipping integration test that requires database and external services")

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
		server := grpcdelivery.Setup(configs)
		assert.NotNil(t, server)
	})
}

func TestGrpcServer_SetupInvalidConfig(t *testing.T) {
	// Test with invalid config should panic
	configs := config.Config{
		App: config.App{
			Port: "invalid-port",
		},
		DB: config.DB{
			Driver: "unsupported",
		},
	}

	assert.Panics(t, func() {
		grpcdelivery.Setup(configs)
	})
}
