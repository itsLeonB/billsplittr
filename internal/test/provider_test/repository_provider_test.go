package provider_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/stretchr/testify/assert"
)

func TestProvideRepositoriesNilDB(t *testing.T) {
	logger := &mockLogger{}

	assert.Panics(t, func() {
		provider.ProvideRepositories(nil, logger)
	})
}

func TestProvideRepositoriesValidDB(t *testing.T) {
	// This test would require setting up a real GORM DB connection
	// For now, we'll just test the panic case above
	// In a real scenario, you'd use an in-memory SQLite database for testing
	t.Skip("Requires GORM DB setup")
}
