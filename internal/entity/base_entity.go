package entity

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt time.Time
}

func (be BaseEntity) IsZero() bool {
	return be.ID == uuid.Nil
}

func (be BaseEntity) IsDeleted() bool {
	return be.DeletedAt.IsZero()
}

type Specification struct {
	PreloadRelations []string
}
