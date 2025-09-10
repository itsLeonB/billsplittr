package provider_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/stretchr/testify/assert"
)

// Mock logger for testing
type mockLogger struct{}

func (m *mockLogger) Debug(args ...interface{})                 {}
func (m *mockLogger) Debugf(format string, args ...interface{}) {}
func (m *mockLogger) Info(args ...interface{})                  {}
func (m *mockLogger) Infof(format string, args ...interface{})  {}
func (m *mockLogger) Warn(args ...interface{})                  {}
func (m *mockLogger) Warnf(format string, args ...interface{})  {}
func (m *mockLogger) Error(args ...interface{})                 {}
func (m *mockLogger) Errorf(format string, args ...interface{}) {}
func (m *mockLogger) Fatal(args ...interface{})                 {}
func (m *mockLogger) Fatalf(format string, args ...interface{}) {}

func TestProvideServicesNilRepositories(t *testing.T) {
	googleConfig := config.Google{
		BillBucketName: "test-bucket",
	}
	logger := &mockLogger{}

	assert.Panics(t, func() {
		provider.ProvideServices(googleConfig, nil, logger)
	})
}

func TestProvideServicesValidRepositories(t *testing.T) {
	// This test would require setting up mock repositories
	// For now, we'll just test the panic case above
	// In a real scenario, you'd create mock implementations of all repository interfaces
	t.Skip("Requires mock repository implementations")
}
