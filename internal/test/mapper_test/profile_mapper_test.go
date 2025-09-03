package mapper_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFromProfileProto_Success(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	protoResponse := &profile.ProfileResponse{
		Id:          id.String(),
		UserId:      userID.String(),
		Name:        "Test User",
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
		IsAnonymous: false,
	}

	result, err := mapper.FromProfileProto(protoResponse)

	assert.NoError(t, err)
	assert.Equal(t, id, result.ID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "Test User", result.Name)
	assert.False(t, result.IsAnonymous)
}

func TestFromProfileProto_AnonymousUser(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	protoResponse := &profile.ProfileResponse{
		Id:          id.String(),
		UserId:      "",
		Name:        "Anonymous User",
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
		IsAnonymous: true,
	}

	result, err := mapper.FromProfileProto(protoResponse)

	assert.NoError(t, err)
	assert.Equal(t, id, result.ID)
	assert.Equal(t, uuid.Nil, result.UserID)
	assert.Equal(t, "Anonymous User", result.Name)
	assert.True(t, result.IsAnonymous)
}

func TestFromProfileProto_NilProto(t *testing.T) {
	result, err := mapper.FromProfileProto(nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "proto is nil")
	assert.Equal(t, uuid.Nil, result.ID)
}

func TestFromProfileProto_InvalidID(t *testing.T) {
	protoResponse := &profile.ProfileResponse{
		Id:          "invalid-uuid",
		Name:        "Test User",
		IsAnonymous: false,
	}

	result, err := mapper.FromProfileProto(protoResponse)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, result.ID)
}

func TestFromProfileProto_AnonymousWithUserID_Error(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()

	protoResponse := &profile.ProfileResponse{
		Id:          id.String(),
		UserId:      userID.String(),
		Name:        "Test User",
		IsAnonymous: true,
	}

	result, err := mapper.FromProfileProto(protoResponse)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "anonymous user has userID")
	assert.Equal(t, uuid.Nil, result.ID)
}
