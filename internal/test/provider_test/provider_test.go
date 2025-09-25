package provider_test

import (
	"testing"

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
