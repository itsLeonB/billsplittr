package util_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestCompareUUID(t *testing.T) {
	uuid1 := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	uuid2 := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
	uuid3 := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000") // same as uuid1

	tests := []struct {
		name     string
		a        uuid.UUID
		b        uuid.UUID
		expected int
	}{
		{"UUID1 < UUID2", uuid1, uuid2, -1},
		{"UUID2 > UUID1", uuid2, uuid1, 1},
		{"UUID1 == UUID3", uuid1, uuid3, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.CompareUUID(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "CompareUUID(%v, %v)", tt.a, tt.b)
		})
	}
}
