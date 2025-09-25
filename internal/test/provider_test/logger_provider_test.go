package provider_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/stretchr/testify/assert"
)

func TestProvideLoggerDebugMode(t *testing.T) {
	logger := provider.ProvideLogger(config.AppName, "debug")

	assert.NotNil(t, logger)
	// Logger should be created successfully for debug mode
}

func TestProvideLoggerReleaseMode(t *testing.T) {
	logger := provider.ProvideLogger(config.AppName, "release")

	assert.NotNil(t, logger)
	// Logger should be created successfully for release mode
}

func TestProvideLoggerDefaultMode(t *testing.T) {
	logger := provider.ProvideLogger(config.AppName, "production")

	assert.NotNil(t, logger)
	// Logger should be created successfully for any other mode
}
