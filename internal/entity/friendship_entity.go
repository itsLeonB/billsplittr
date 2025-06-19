package entity

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
)

type Friendship struct {
	BaseEntity
	ProfileID1 uuid.UUID
	ProfileID2 uuid.UUID
	Type       appconstant.FriendshipType
	Profile1   UserProfile
	Profile2   UserProfile
}

type FriendshipSpecification struct {
	Friendship
	ProfileID uuid.UUID
	Name      string
}
