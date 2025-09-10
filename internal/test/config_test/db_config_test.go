package config_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestDBAllFields(t *testing.T) {
	db := config.DB{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
	}

	assert.Equal(t, "postgres", db.Driver)
	assert.Equal(t, "localhost", db.Host)
	assert.Equal(t, "5432", db.Port)
	assert.Equal(t, "testuser", db.User)
	assert.Equal(t, "testpass", db.Password)
	assert.Equal(t, "testdb", db.Name)
}

func TestDBZeroValues(t *testing.T) {
	db := config.DB{}
	
	assert.Equal(t, "", db.Driver)
	assert.Equal(t, "", db.Host)
	assert.Equal(t, "", db.Port)
	assert.Equal(t, "", db.User)
	assert.Equal(t, "", db.Password)
	assert.Equal(t, "", db.Name)
}
