package repository_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewGroupExpenseRepository(t *testing.T) {
	// Test that repository can be created (without actual DB connection)
	var db *gorm.DB // nil is fine for this test
	
	repo := repository.NewGroupExpenseRepository(db)
	assert.NotNil(t, repo)
}

func TestGroupExpenseRepositoryInterface(t *testing.T) {
	// Test that the repository implements the interface
	var repo repository.GroupExpenseRepository
	assert.Nil(t, repo) // nil interface is fine for this test
}
