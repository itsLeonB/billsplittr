package entity

import "github.com/google/uuid"

type UserProfile struct {
	BaseEntity
	UserID uuid.UUID
	Name   string
}

func (up UserProfile) IsZero() bool {
	return up.ID == uuid.Nil
}
