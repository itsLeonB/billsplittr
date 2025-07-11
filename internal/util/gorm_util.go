package util

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func DefaultOrder() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}
}

func ForUpdate(enable bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if enable {
			return db.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
		}
		return db
	}
}
