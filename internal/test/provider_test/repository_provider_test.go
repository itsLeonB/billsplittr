package provider_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/stretchr/testify/assert"
)

func TestProvideRepositoriesNilDB(t *testing.T) {
	googleConfig := config.Google{
		ServiceAccount: "{}",
		BillBucketName: "test-bucket",
	}
	logger := &mockLogger{}

	assert.Panics(t, func() {
		provider.ProvideRepositories(nil, googleConfig, logger)
	})
}

func TestProvideRepositoriesValidDB(t *testing.T) {
	// This test would require setting up a real GORM DB connection
	// For now, we'll just test the panic case above
	// In a real scenario, you'd use an in-memory SQLite database for testing
	t.Skip("Requires GORM DB setup")
}

func TestRepositoriesShutdown(t *testing.T) {
	// Test shutdown with nil storage
	repos := &provider.Repositories{}
	err := repos.Shutdown()
	assert.NoError(t, err)
}
