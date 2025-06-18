package entity

import "github.com/google/uuid"

type TransferMethod struct {
	BaseEntity
	Name    string
	Display string
}

func (tm TransferMethod) IsZero() bool {
	return tm.ID == uuid.Nil
}
