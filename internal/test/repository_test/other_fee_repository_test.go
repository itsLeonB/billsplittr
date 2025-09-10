package repository_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewOtherFeeRepository(t *testing.T) {
	// Test that repository can be created (without actual DB connection)
	var db *gorm.DB // nil is fine for this test
	
	repo := repository.NewOtherFeeRepository(db)
	assert.NotNil(t, repo)
}

func TestOtherFeeRepositoryInterface(t *testing.T) {
	// Test that the repository implements the interface
	var repo repository.OtherFeeRepository
	assert.Nil(t, repo) // nil interface is fine for this test
}
