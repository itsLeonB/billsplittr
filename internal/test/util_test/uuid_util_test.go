package util_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestUUIDUtil_ToString(t *testing.T) {
	id := uuid.New()
	result := util.ToString(id)
	
	assert.Equal(t, id.String(), result)
	assert.NotEmpty(t, result)
}
