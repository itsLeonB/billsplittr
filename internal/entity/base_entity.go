package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt sql.NullTime
}

func (be BaseEntity) IsZero() bool {
	return be.ID == uuid.Nil
}

func (be BaseEntity) IsDeleted() bool {
	return be.DeletedAt.Valid
}

func (be BaseEntity) GetID() string {
	return be.ID.String()
}

type Specification struct {
	PreloadRelations []string
	ForUpdate        bool
}

type GenericSpec[T any] struct {
	Model            T
	PreloadRelations []string
	ForUpdate        bool
}

type Entity interface {
	SimpleName() string
	GetID() string
}
