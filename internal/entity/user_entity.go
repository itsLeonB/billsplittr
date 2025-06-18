package entity

import "github.com/google/uuid"

type User struct {
	BaseEntity
	Email    string
	Password string
}

func (u *User) IsZero() bool {
	return u.ID == uuid.Nil
}
