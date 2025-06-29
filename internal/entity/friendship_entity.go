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
	Profile1   UserProfile `gorm:"foreignKey:ProfileID1"`
	Profile2   UserProfile `gorm:"foreignKey:ProfileID2"`
}

type FriendshipSpecification struct {
	Specification
	Friendship
	Name string
}
