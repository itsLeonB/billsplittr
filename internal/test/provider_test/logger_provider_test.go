package provider_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/stretchr/testify/assert"
)

func TestProvideLoggerDebugMode(t *testing.T) {
	appConfig := config.App{
		Name: "test-app",
		Env:  "debug",
	}

	logger := provider.ProvideLogger(appConfig)

	assert.NotNil(t, logger)
	// Logger should be created successfully for debug mode
}

func TestProvideLoggerReleaseMode(t *testing.T) {
	appConfig := config.App{
		Name: "test-app",
		Env:  "release",
	}

	logger := provider.ProvideLogger(appConfig)

	assert.NotNil(t, logger)
	// Logger should be created successfully for release mode
}

func TestProvideLoggerDefaultMode(t *testing.T) {
	appConfig := config.App{
		Name: "test-app",
		Env:  "production", // Not "release", so should use debug level
	}

	logger := provider.ProvideLogger(appConfig)

	assert.NotNil(t, logger)
	// Logger should be created successfully for any other mode
}

func TestProvideLoggerEmptyName(t *testing.T) {
	appConfig := config.App{
		Name: "",
		Env:  "debug",
	}

	logger := provider.ProvideLogger(appConfig)

	assert.NotNil(t, logger)
	// Logger should still be created even with empty name
}
